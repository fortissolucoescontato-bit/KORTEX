# Technical Design — bug-fixes

**Change ID:** `bug-fixes`
**Scope:** 10 defects across `state`, `pipeline`, `tui`, `agentbuilder`, `components/kortex-engram`, `backup`, and `cli`.
**Stack:** Go 1.25.0 (`errors.Join` available), BubbleTea TUI, SQLite (`modernc.org/sqlite`).
**Mode:** Strict TDD — `go test ./...` is the authoritative oracle.

---

## 1. Guiding Principles

1. **Minimum surface change.** Each bug fix is contained; we avoid cascading refactors.
2. **Preserve public API where feasible.** Only Bug #2 requires a signature change, and it is additive via a new exported function; the old signature is preserved through a thin shim.
3. **Best-effort continuity over fail-fast in cleanup paths.** Rollback and filesystem cleanup must keep running past individual errors, aggregating them with `errors.Join`.
4. **Deterministic timeouts for all background I/O.** Every `context.Background()` reaching an external process (verifier) or long loop must be wrapped with `WithTimeout`.
5. **Surface silent data integrity errors.** Any failure to hash, checksum, or rollback must be reported upward; logs alone are insufficient for backup/rollback correctness.

---

## 2. Per-Bug Architecture Decisions

### Bug #1 — `rows.Err()` not checked (CRITICAL)

**Files:** `internal/state/state.go` — `GetInstalledAgents` (L32-49), `GetAssignments` (L72-92).

**Decision.** After the `for rows.Next()` loop, check `rows.Err()` and propagate.

```go
for rows.Next() { ... }
if err := rows.Err(); err != nil {
    return nil, err
}
```

**Alternatives rejected.**
- *Wrap in a generic helper* (e.g. `scanRows[T]`): premature generalization; only two callsites, and the current code has distinct scan shapes. Reconsider if a third callsite appears.
- *Ignore (status quo):* `database/sql` documents `rows.Next()` returning false on both EOF and error; silent truncation of result sets is a correctness bug on transient SQLite I/O failures.

**Interface changes.** None.

**Interaction.** Both functions are read paths feeding the TUI; the caller already handles errors, so propagation is free.

**Testing.** Inject a mocked `driver.Rows` that returns `ErrBadConn` on the second `Next()`; assert the function returns that error instead of a partial slice. Use `sqlmock` or an in-memory SQLite with a closed connection mid-iteration.

---

### Bug #2 — `ctx` leak in `startInstallation()` (CRITICAL)

**Files:** `internal/tui/model.go:3390-3392`, `internal/agentbuilder/installer.go:20`.

**Current state.** A 5-minute `context.WithTimeout` is created inside the tea.Cmd closure, `defer cancel()` is registered, and the context is assigned to `_`. The timeout is effectively dead: `Install()` never sees the ctx, so nothing enforces the deadline. The `defer cancel()` keeps the goroutine's resources tidy on return but does not cancel anything observable.

**Decision.** Propagate a context through the builder:

1. Add a new function `InstallContext(ctx context.Context, agent *GeneratedAgent, adapters []AdapterInfo) ([]InstallResult, error)`.
2. Keep `Install(agent, adapters, _ string)` as a thin shim that calls `InstallContext(context.Background(), ...)` — preserves backward compatibility for any external callers and existing tests.
3. Inside `InstallContext`, honour `ctx.Err()` at the top of each loop iteration and return early with the installed paths rolled back.
4. The third string parameter (currently `_`) is dropped in the new function (it was never used).

**Why not modify `Install()` in place?** The signature already has an unused third parameter, indicating an earlier aborted refactor; adding a fourth is worse. Two functions with clear names is cleaner.

**Why not use `context` for `os.MkdirAll` / `os.WriteFile`?** The stdlib does not accept contexts for these. Honouring ctx between iterations is the strongest realistic guarantee; fsync-level cancellation would require a filesystem abstraction out of scope.

**TUI callsite update.**
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()
results, err := agentbuilder.InstallContext(ctx, installAgent, adapters)
```
Remove the `_ = ctx` line.

**Interface changes.**
- **NEW:** `agentbuilder.InstallContext(ctx context.Context, agent *GeneratedAgent, adapters []AdapterInfo) ([]InstallResult, error)`
- **UNCHANGED:** `agentbuilder.Install(agent *GeneratedAgent, adapters []AdapterInfo, _ string) ([]InstallResult, error)` becomes a shim that forwards to `InstallContext`.

**Sequence (corrected flow).**
```
TUI (startInstallation Cmd goroutine)
  │
  ├── ctx, cancel := WithTimeout(Background, 5m)
  ├── defer cancel()                           ← now effective: cancels on return
  │
  ▼
agentbuilder.InstallContext(ctx, agent, adapters)
  │
  for each adapter:
    ├── if ctx.Err() != nil → rollback(written) + return ctx.Err()
    ├── os.MkdirAll  ← not ctx-aware
    ├── os.WriteFile ← not ctx-aware
    └── written = append(written, skillFile)
  │
  ▼
return results, nil → AgentBuilderInstallDoneMsg
```

**Testing.** Table-driven test in `installer_test.go`:
- `ctx` already canceled → returns `ctx.Err()`, no files written.
- `ctx` cancels mid-loop (2 of 3 adapters) → first file is rolled back, results length = adapters attempted.
- Happy path → identical behaviour to pre-fix `Install()`.

---

### Bug #3 — `ExecuteRollback` stops on first error (CRITICAL)

**File:** `internal/pipeline/rollback.go:40-48`.

**Current state.** On any `rollbackStep.Rollback()` error, the function appends the failed step and `return`s — remaining previously-applied steps are never rolled back. This leaves the system in an inconsistent half-rolled-back state.

**Decision.** Continue iterating past rollback failures; collect all errors with `errors.Join` (Go 1.20+, already available in Go 1.25). Each failed step still appears individually in `result.Steps` for observability.

```go
var joined error
for i := len(steps) - 1; i >= 0; i-- {
    // ... existing filtering ...
    err := rollbackStep.Rollback()
    item := StepResult{StepID: rollbackStep.ID(), Status: StepStatusRolledBack}
    if err != nil {
        item.Status = StepStatusFailed
        item.Err = err
        result.Success = false
        joined = errors.Join(joined, fmt.Errorf("rollback step %q: %w", rollbackStep.ID(), err))
    }
    result.Steps = append(result.Steps, item)
}
if joined != nil {
    result.Err = joined
}
return result
```

**Alternatives rejected.**
- *Aggregate errors in a slice + custom type:* loses `errors.Is`/`errors.As` semantics. `errors.Join` already gives us wrap traversal.
- *Best-effort without reporting:* rollback failures are operationally critical; silently swallowing them would make partial-state debugging impossible.

**Interface changes.** None. `StageResult.Err` semantics slightly broaden: may now wrap multiple rollback errors via `errors.Join`. Callers using `errors.Is` are unaffected; callers using `== ErrX` would already have been broken because the error was already wrapped with `fmt.Errorf`.

**Sequence.**
```
Applied steps: S1 → S2 → S3 → S4  (reverse order: S4, S3, S2, S1)

Before:
  Rollback S4 ✓
  Rollback S3 ✗ → return (S2, S1 LEAKED)

After:
  Rollback S4 ✓
  Rollback S3 ✗ → record + continue
  Rollback S2 ✓
  Rollback S1 ✗ → record + continue
  return StageResult{Success: false, Err: errors.Join(S3_err, S1_err)}
```

**Testing.** Build four fake `RollbackStep` impls; assert all four are invoked even when the middle two fail, and that `errors.Is(result.Err, S3Err)` and `errors.Is(result.Err, S1Err)` both hold.

---

### Bug #4 — Dead `reH2Section` package regex (SERIOUS)

**File:** `internal/agentbuilder/parser.go:17`.

**Current state.** `reH2Section` is compiled at package init with a literal `%s` in the pattern; nothing uses it. `extractSection()` recompiles an equivalent regex on every call via `regexp.MustCompile(...)` — a measurable hotspot if many agents are parsed.

**Decision.** Remove `reH2Section`. Replace the per-call compilation in `extractSection` with a cached regex that matches any H2 and filters by captured name.

```go
// reH2SectionAny matches any "## name" section, capturing the name and body.
var reH2SectionAny = regexp.MustCompile(`(?ms)^##\s+(.+?)\s*\n(.*?)(?:^##\s|\z)`)

func extractSection(content, name string) (string, error) {
    for _, m := range reH2SectionAny.FindAllStringSubmatch(content, -1) {
        if strings.EqualFold(strings.TrimSpace(m[1]), name) {
            body := strings.TrimSpace(m[2])
            if body == "" {
                return "", errors.New("parse: '## " + name + "' section is empty")
            }
            return body, nil
        }
    }
    return "", errors.New("parse: missing '## " + name + "' section")
}
```

**Alternatives rejected.**
- *Map of pre-compiled named regexes (e.g. `regexMap[name]`):* unnecessary cache growth keyed by arbitrary section names; the generic-scan variant is O(n) over sections, which is always small (<10).
- *Keep `reH2Section` and use `fmt.Sprintf` at runtime:* defeats the purpose of package-level compilation, and requires escaping user input. Pure O(1)-compile generic regex is simpler.

**Interface changes.** None (both symbols are unexported).

**Testing.** Existing parser tests continue to pass. Add one test confirming a section whose name contains regex metachars (e.g. `## A.B`) is matched literally (via `strings.EqualFold` comparison, not regex).

---

### Bug #5 — `rollback()` leaves empty dirs (SERIOUS)

**File:** `internal/agentbuilder/installer.go:70-73`.

**Current state.** `Install` calls `os.MkdirAll(skillDir, 0755)` for each adapter (where `skillDir = adapters.SkillsDir + "/" + agent.Name`). On rollback we remove the SKILL.md files but leave the per-agent directories — future installs see phantom entries and `filepath.Walk` output is polluted.

**Decision.** Track directories created by `MkdirAll` alongside written file paths, then `os.Remove` them (not `RemoveAll` — we only want empty dirs we created) in reverse order after removing files. Errors from `os.Remove` are swallowed (dir may not be empty if user added files; that's intentional).

```go
type rollbackState struct {
    files []string
    dirs  []string
}

func rollback(s rollbackState) {
    for _, p := range s.files {
        _ = os.Remove(p)
    }
    for i := len(s.dirs) - 1; i >= 0; i-- {
        _ = os.Remove(s.dirs[i]) // fails silently if non-empty — intentional
    }
}
```

**Why not `RemoveAll`?** If the user or another process placed unrelated files into the agent dir between `MkdirAll` and rollback, `RemoveAll` would destroy them. `os.Remove` on a directory only succeeds when empty — exactly the safety we need.

**Only remove the innermost dir we created?** `MkdirAll` may have created multiple levels, but `SkillsDir` is a precondition (caller owns it). Removing only the agent-name leaf is correct.

**Interface changes.** Internal only — `rollback` changes parameter type (it's unexported).

**Testing.** Integration test: install into a temp dir with two adapters; force the second `WriteFile` to fail by pointing at a read-only subdir; assert both the file and the agent-named subdir of the first adapter are gone after rollback.

---

### Bug #6 — Dead third `LookPath` branch (SERIOUS)

**File:** `internal/components/kortex-engram/inject.go:62-74`.

**Current state.** Line 63 and line 72 both call `LookPath("kortex-engram")` — the third branch is unreachable. The comment says "legacy" but the string is identical. This appears to be a merge/refactor mistake.

**Decision.** Two fallback names exist in reality: the new `kortex-engram` and the interim `kortex`. Remove the redundant third branch entirely.

```go
func resolveKortexEngramCommand() (string, bool) {
    if p, err := kortexEngramLookPath("kortex-engram"); err == nil && p != "" {
        return p, true
    }
    if p, err := kortexEngramLookPath("kortex"); err == nil && p != "" {
        return p, true
    }
    return "kortex-engram", false
}
```

**Alternatives rejected.**
- *Keep the dead branch "just in case":* dead code rots; if a truly legacy name is ever needed, re-add it deliberately.

**Interface changes.** None.

**Testing.** Unit test stubs `kortexEngramLookPath` (package-level var, already injectable) to return results for different names; assert returned path and bool match expectations.

---

### Bug #7 — Redundant switch-branch assignments (SERIOUS)

**File:** `internal/components/kortex-engram/inject.go:480-499`.

**Current state.** Three switch branches follow the pattern:
```go
server = mcp["kortex-engram"]
if server == nil {
    server = mcp["kortex-engram"] // Fallback  ← identical lookup
}
```
The fallback is a no-op.

**Decision.** Remove the redundant `if server == nil` block from all three branches. The key `kortex-engram` is canonical post-rebrand; if legitimate legacy keys exist in user configs (e.g. `graphiti`, `kortex`), those would require explicit fallback strings — but the current code does not do that and introducing them is out of scope for this bugfix.

```go
case model.AgentOpenCode:
    mcp, ok := root["mcp"].(map[string]any)
    if !ok {
        return "", false
    }
    server = mcp["kortex-engram"]
```

**Alternatives rejected.**
- *Fall back to legacy keys (`graphiti`, `kortex`):* would change detection semantics and risks silently migrating unintended configs. File a separate change if desired.

**Interface changes.** None.

**Testing.** Existing detection tests must continue to pass; add coverage for a config missing the `kortex-engram` key to assert `ok == false`.

---

### Bug #8 — Checksum error silenced in snapshot (WARNING)

**File:** `internal/backup/snapshot.go:82`.

**Current state.** If `ComputeChecksum` fails, the error is logged via `log.Printf` and the manifest is written with `checksum = ""`. Empty checksum breaks deduplication (manifest comparison treats all zero-files as duplicates but all checksum-failed snapshots as distinct duplicates).

**Decision.** Return the error. A backup with an unverifiable checksum is worse than no backup: restoration tools cannot trust it, and partial corruption is undetectable.

```go
checksum, csErr = ComputeChecksum(existingPaths)
if csErr != nil {
    return Manifest{}, fmt.Errorf("backup: compute checksum: %w", csErr)
}
```

**Alternatives considered.**
- *Sentinel checksum `"UNVERIFIED"`:* pollutes the manifest schema and pushes the error handling downstream. The upstream caller is a better place to decide whether to retry.

**Interface changes.** `Snapshotter.Create` — already returns `error`; behaviour change only.

**Testing.** Inject a `ComputeChecksum` that returns an error (can be done via a file handle mock or by removing a file between enumeration and checksum); assert `Create` returns a wrapped error and no manifest is written.

---

### Bug #9 — String concat path join (WARNING)

**File:** `internal/components/kortex-engram/inject.go:384`.

**Decision.** Use `filepath.Join(homeDir, ".codex")` instead of `homeDir + "/.codex"`. Correctness on Windows and defensive against trailing separators.

**Interface changes.** None.

**Testing.** Existing tests suffice; add a case with `homeDir` that has a trailing slash to confirm normalization.

---

### Bug #10 — `context.Background()` without timeout in verify (WARNING)

**Files:** `internal/cli/sync.go:848`, `internal/cli/run.go:1041`.

**Current state.** `verify.RunChecks(context.Background(), checks)` runs all verification checks with no timeout. A hung check (network socket, process spawn) blocks sync/run indefinitely.

**Decision.** Introduce a single package-level constant `verifyTimeout = 30 * time.Second` in an existing verify helper or in the caller's package. Wrap each callsite:

```go
// internal/cli/postverify.go (new) or inline:
const postVerifyTimeout = 30 * time.Second

ctx, cancel := context.WithTimeout(context.Background(), postVerifyTimeout)
defer cancel()
report := verify.RunChecks(ctx, checks)
```

**Why a constant and not configurable?** Post-install verification is a bounded set of local checks. 30 s is generous for any realistic check. A CLI flag adds surface area without clear value; if a user needs more, a follow-up change can add `--verify-timeout`.

**Alternatives rejected.**
- *`context.Background()` forever (status quo):* real hang risk on any check that opens sockets or spawns processes.
- *Per-check timeouts:* `RunChecks` already iterates; adding per-check timeouts is a larger refactor belonging in the verify package itself, not in CLI callers.

**Interface changes.** None.

**Testing.** Inject a `Check` that blocks on `<-ctx.Done()` and returns `ctx.Err()`. Assert total elapsed time ≤ `postVerifyTimeout + 1s` and the returned report marks the check as timed out.

---

## 3. Consolidated Interface Changes

| Symbol | Before | After |
|---|---|---|
| `agentbuilder.Install` | `Install(agent *GeneratedAgent, adapters []AdapterInfo, _ string) ([]InstallResult, error)` | Unchanged (shim delegating to `InstallContext`) |
| `agentbuilder.InstallContext` | — | **NEW** `InstallContext(ctx context.Context, agent *GeneratedAgent, adapters []AdapterInfo) ([]InstallResult, error)` |
| `agentbuilder.rollback` (unexported) | `rollback(paths []string)` | `rollback(s rollbackState)` |
| `agentbuilder.reH2Section` | package var, unused | **REMOVED** — replaced by `reH2SectionAny` |
| `state.Manager.GetInstalledAgents` | missing `rows.Err()` | now propagates `rows.Err()` |
| `state.Manager.GetAssignments` | missing `rows.Err()` | now propagates `rows.Err()` |
| `pipeline.ExecuteRollback` | stops on first error | continues, aggregates with `errors.Join` |
| `backup.Snapshotter.Create` | silences checksum error | returns checksum error |

Total: 1 new exported function (`InstallContext`), 1 private signature change (`rollback`), 0 breaking changes.

---

## 4. Testing Strategy (Strict TDD)

For each bug, the work order is Red → Green → Refactor:

1. **Write the failing test first.** Tests go in the package adjacent to the fix (`state_test.go`, `rollback_test.go`, `installer_test.go`, `parser_test.go`, `inject_test.go`, `snapshot_test.go`, and a new `postverify_test.go` for bug #10).
2. **Run `go test ./<pkg>/...`** to confirm red.
3. **Apply the fix.** Re-run; confirm green.
4. **Run full `go test ./...`** before moving to next bug.

Cross-cutting:
- Use `t.TempDir()` for all filesystem tests (bugs #5, #8, #9).
- Use `context.WithCancel` + `context.WithTimeout` helpers for ctx tests (#2, #10).
- For bug #1 prefer an in-memory SQLite with a forced close mid-iteration rather than `sqlmock`; `modernc.org/sqlite` is already a dep.
- TUI test for bug #2 is NOT required at the BubbleTea layer — the fix is in `agentbuilder`; the TUI change is a mechanical callsite update covered by the TUI's existing smoke tests.

**Coverage target.** Each fix must add at least one test exercising the previously broken path. No fix lands without a failing-then-passing test in the same PR.

---

## 5. Risks

1. **`errors.Join` wrapping may surprise a caller that does `==` comparison on `StageResult.Err`.** Mitigation: audit callers; all current callers use `if err != nil` or `errors.Is`.
2. **Bug #5 dir cleanup could remove a dir the user populated between install and rollback.** Mitigation: `os.Remove` (not `RemoveAll`) refuses to delete non-empty dirs.
3. **Bug #8 turning a warning into a hard failure may break CI pipelines that previously tolerated checksum failures.** Mitigation: document in CHANGELOG; the previous behaviour was silent data corruption and must break.
4. **Bug #2 `InstallContext` early-return semantics differ from old `Install` when ctx is already cancelled.** Mitigation: `Install` shim passes `context.Background()`, which never cancels — existing callers see zero behavioural change.
5. **Bug #10 `30s` timeout may be too aggressive for slow CI filesystems.** Mitigation: start at 30 s; if flakiness appears, promote to a var and expose `--verify-timeout`.

---

## 6. Out of Scope

- Full context propagation through `os.MkdirAll` / `os.WriteFile` (requires FS abstraction).
- Configurable verify timeouts via CLI flag.
- Legacy key fallbacks in `existingMergedKortexEngramCommand` (separate migration story).
- Per-check timeouts inside `verify.RunChecks`.
