# Spec: bug-fixes

**Change**: bug-fixes  
**Date**: 2026-04-24  
**Status**: draft  
**TDD Mode**: strict (RED → GREEN → REFACTOR per bug)

---

## Overview

This spec covers 10 bugs across 5 packages. Each section defines the exact behavioral contract
(RFC 2119 keywords) and Given/When/Then scenarios that drive the TDD cycle.
Bugs are ordered by severity: CRITICAL → SERIOUS → WARNING.

---

## Bug #1 — `rows.Err()` not checked after scan loops

**File**: `internal/state/state.go`  
**Symbols**: `GetInstalledAgents` (lines 33–49), `GetAssignments` (lines 73–91)  
**Severity**: CRITICAL

### Requirements

- `GetInstalledAgents` MUST call `rows.Err()` after the `rows.Next()` loop and return that error
  if non-nil.
- `GetAssignments` MUST call `rows.Err()` after the `rows.Next()` loop and return that error
  if non-nil.
- Both functions SHALL return `nil, err` (or `nil, err`) when `rows.Err()` is non-nil.
- The existing early-return paths on `rows.Scan` failure MUST NOT be removed.

### Scenarios

#### Scenario 1.1 — GetInstalledAgents propagates cursor error

```
Given  a database driver that injects a non-nil error into rows.Err()
  after yielding zero or more valid rows
When   GetInstalledAgents() is called
Then   it MUST return (nil, <the injected error>)
  And  no partial agent slice MUST be returned
```

#### Scenario 1.2 — GetInstalledAgents succeeds when rows.Err() is nil

```
Given  a database with two installed agents ["agent-a", "agent-b"]
  And  rows.Err() returns nil
When   GetInstalledAgents() is called
Then   it MUST return (["agent-a", "agent-b"], nil)
```

#### Scenario 1.3 — GetAssignments propagates cursor error

```
Given  a database driver that injects a non-nil error into rows.Err()
  after yielding zero or more valid assignment rows
When   GetAssignments("claude") is called
Then   it MUST return (nil, <the injected error>)
```

#### Scenario 1.4 — GetAssignments succeeds when rows.Err() is nil

```
Given  a database with one assignment {phase:"orchestrator", provider:"anthropic", model:"opus"}
  for agentID "claude"
  And  rows.Err() returns nil
When   GetAssignments("claude") is called
Then   it MUST return (map{"orchestrator": {ProviderID:"anthropic", ModelID:"opus"}}, nil)
```

---

## Bug #2 — Context created but unused in `startInstallation()`

**File**: `internal/tui/model.go`  
**Symbol**: `startInstallation()` (lines 3390–3392)  
**Severity**: CRITICAL

### Requirements

- The 5-minute `context.WithTimeout` MUST be passed into — or honoured by — the installation
  call chain.
- If `agentbuilder.Install` does not accept a `context.Context`, a goroutine watchdog MUST be
  used: when `<-ctx.Done()` fires before the install goroutine finishes, the goroutine MUST be
  signalled (via a `done` channel or `context.CancelCauseFunc`) and the installation MUST be
  marked as failed with an appropriate timeout error.
- The `_ = ctx` suppression line MUST be removed.
- The `cancel()` defer MUST remain to avoid context leak.

### Scenarios

#### Scenario 2.1 — Timeout fires before install completes

```
Given  a 5-minute context timeout
  And  agentbuilder.Install blocks longer than the timeout
When   startInstallation() runs
Then   the install attempt MUST be cancelled/interrupted
  And  an error describing the timeout MUST be surfaced to the TUI result channel
  And  cancel() MUST have been called (no goroutine leak)
```

#### Scenario 2.2 — Install completes within timeout

```
Given  a 5-minute context timeout
  And  agentbuilder.Install returns successfully before the timeout
When   startInstallation() runs
Then   the successful InstallResult MUST be sent to the result channel
  And  no timeout error MUST be raised
```

---

## Bug #3 — `ExecuteRollback` aborts on first error

**File**: `internal/pipeline/rollback.go`  
**Symbol**: `ExecuteRollback` (lines 21–55)  
**Severity**: CRITICAL

### Requirements

- `ExecuteRollback` MUST attempt to roll back ALL previously-succeeded steps regardless of
  individual rollback failures.
- Errors from individual rollback steps MUST be accumulated using `errors.Join`.
- The returned `StageResult.Err` MUST contain the joined error of all rollback failures.
- `StageResult.Success` MUST be `false` if any rollback step failed.
- Steps that fail to roll back MUST still have `StepStatusFailed` in `StageResult.Steps`.
- Steps that succeed in rolling back MUST have `StepStatusRolledBack` in `StageResult.Steps`.
- The function MUST NOT return early after a single rollback failure.

### Scenarios

#### Scenario 3.1 — All rollback steps succeed

```
Given  three steps that all previously succeeded (StepStatusSucceeded)
  And  all three implement RollbackStep and their Rollback() returns nil
When   ExecuteRollback() is called
Then   StageResult.Success MUST be true
  And  StageResult.Steps MUST contain three entries with StepStatusRolledBack
  And  StageResult.Err MUST be nil
```

#### Scenario 3.2 — One rollback step fails, others continue

```
Given  three steps that all previously succeeded
  And  step[1].Rollback() returns an error
  And  step[0].Rollback() and step[2].Rollback() return nil
When   ExecuteRollback() is called
Then   StageResult.Success MUST be false
  And  StageResult.Steps MUST contain entries for all three steps
  And  step[0] and step[2] MUST have StepStatusRolledBack
  And  step[1] MUST have StepStatusFailed
  And  StageResult.Err MUST contain step[1]'s error (via errors.Join)
```

#### Scenario 3.3 — All rollback steps fail

```
Given  two steps that all previously succeeded
  And  both Rollback() calls return distinct errors
When   ExecuteRollback() is called
Then   StageResult.Success MUST be false
  And  StageResult.Err MUST contain both errors (errors.Join)
  And  both steps MUST have StepStatusFailed
```

#### Scenario 3.4 — No succeeded steps to roll back

```
Given  a steps slice where all entries have StepStatusFailed (nothing to undo)
When   ExecuteRollback() is called
Then   StageResult.Success MUST be true
  And  StageResult.Steps MUST be empty
  And  StageResult.Err MUST be nil
```

---

## Bug #4 — Dead package-level regex + per-call recompile in `parser.go`

**File**: `internal/agentbuilder/parser.go`  
**Symbol**: `reH2Section` (line 17), `extractSection` (lines 103–114)  
**Severity**: SERIOUS

### Requirements

- The package-level variable `reH2Section` (which contains a literal `%s` placeholder and is
  never referenced) MUST be removed.
- `extractSection` MUST NOT compile a new `*regexp.Regexp` on every invocation.
- The compiled regex for H2 section extraction SHALL be cached. Acceptable approaches:
  - A `sync.Map` keyed by section name, OR
  - A fixed set of package-level compiled regexes for known sections (Description, Trigger,
    Instructions), OR
  - `sync.OnceValue` / `sync.Once` per section name.
- The cached implementation MUST produce identical output to the current per-call implementation
  for all valid and invalid inputs.
- Removing `reH2Section` MUST NOT break any existing test or build target.

### Scenarios

#### Scenario 4.1 — reH2Section is gone

```
Given  the parser package source
When   it is compiled
Then   no identifier named reH2Section MUST exist in the package
```

#### Scenario 4.2 — extractSection returns the same result with cached regex

```
Given  a markdown string with a ## Description section containing "hello world"
When   extractSection(content, "Description") is called 100 times sequentially
Then   every call MUST return ("hello world", nil)
  And  only one regexp.Compile call MUST occur (verified via benchmark or counter)
```

#### Scenario 4.3 — extractSection still errors on missing section

```
Given  a markdown string with no ## Description section
When   extractSection(content, "Description") is called
Then   it MUST return ("", error containing "## Description")
```

---

## Bug #5 — `rollback()` in installer leaves orphaned directories

**File**: `internal/agentbuilder/installer.go`  
**Symbol**: `rollback` (lines 69–74)  
**Severity**: SERIOUS

### Requirements

- `rollback` MUST remove the entire skill directory (created by `os.MkdirAll`) for each
  rolled-back path, not only the `SKILL.md` file within it.
- The parent directory MUST be removed using `os.RemoveAll(filepath.Dir(p))` or equivalent.
- The rollback MUST remain best-effort (errors ignored), consistent with current behaviour.
- Directories that were NOT created by `Install` (i.e. the adapter's root `SkillsDir`) MUST NOT
  be deleted — only the per-agent subdirectory (one level below `SkillsDir`) MUST be removed.

### Scenarios

#### Scenario 5.1 — Successful install followed by forced rollback cleans directories

```
Given  two adapters with SkillsDirs ["/tmp/ada/skills", "/tmp/adb/skills"]
  And  Install writes SKILL.md into each adapter's "my-agent/" subdirectory
  And  a failure is then simulated causing rollback
When   rollback() completes
Then   "/tmp/ada/skills/my-agent" MUST NOT exist on the filesystem
  And  "/tmp/adb/skills/my-agent" MUST NOT exist on the filesystem
  And  "/tmp/ada/skills" MUST still exist
  And  "/tmp/adb/skills" MUST still exist
```

#### Scenario 5.2 — Rollback on partial install cleans only written paths

```
Given  three adapters
  And  only the first two were written before the failure
When   rollback() is called with only the first two paths
Then   only those two agent subdirectories MUST be removed
  And  the third adapter's skills directory MUST be unaffected
```

---

## Bug #6 — Duplicate `LookPath` call in `resolveKortexEngramCommand()`

**File**: `internal/components/kortex-engram/inject.go`  
**Symbol**: `resolveKortexEngramCommand` (lines 61–75)  
**Severity**: SERIOUS

### Requirements

- The third branch of `resolveKortexEngramCommand` (lines 71–73), which calls
  `kortexEngramLookPath("kortex-engram")` again after the identical call on line 63, MUST be
  removed.
- The function MUST still try `"kortex-engram"` first, then `"kortex"` as fallback, then return
  `("kortex-engram", false)` if neither is found.
- The removal MUST NOT change observable behaviour for any input.

### Scenarios

#### Scenario 6.1 — kortex-engram binary found on first try

```
Given  LookPath("kortex-engram") returns "/usr/local/bin/kortex-engram"
When   resolveKortexEngramCommand() is called
Then   it MUST return ("/usr/local/bin/kortex-engram", true)
  And  LookPath MUST have been called exactly once
```

#### Scenario 6.2 — kortex-engram not found, kortex found

```
Given  LookPath("kortex-engram") returns an error
  And  LookPath("kortex") returns "/usr/bin/kortex"
When   resolveKortexEngramCommand() is called
Then   it MUST return ("/usr/bin/kortex", true)
```

#### Scenario 6.3 — Neither binary found

```
Given  LookPath returns an error for all inputs
When   resolveKortexEngramCommand() is called
Then   it MUST return ("kortex-engram", false)
  And  LookPath MUST have been called exactly twice (kortex-engram, then kortex)
```

---

## Bug #7 — Self-identical no-op fallback in `existingMergedKortexEngramCommand()`

**File**: `internal/components/kortex-engram/inject.go`  
**Symbol**: `existingMergedKortexEngramCommand` (lines 473–501)  
**Severity**: SERIOUS

### Requirements

- In each `case` branch, the pattern:
  ```go
  server = mcp["kortex-engram"]
  if server == nil {
      server = mcp["kortex-engram"] // identical key — dead code
  }
  ```
  MUST be simplified to a single assignment: `server = mcp["kortex-engram"]`.
- This MUST be done for all three branches (`AgentOpenCode`, `AgentVSCodeCopilot`, `default`).
- The simplified code MUST produce identical output for all possible inputs.

### Scenarios

#### Scenario 7.1 — OpenCode branch resolves server

```
Given  agentID == AgentOpenCode
  And  root["mcp"] == map{"kortex-engram": <serverMap>}
When   existingMergedKortexEngramCommand() is called
Then   it MUST return the command extracted from <serverMap>
```

#### Scenario 7.2 — OpenCode branch returns empty when key absent

```
Given  agentID == AgentOpenCode
  And  root["mcp"] == map{} (key "kortex-engram" absent)
When   existingMergedKortexEngramCommand() is called
Then   it MUST return ("", false)
```

#### Scenario 7.3 — Dead fallback assignment is gone (static analysis)

```
Given  the source of existingMergedKortexEngramCommand
When   each case branch is inspected
Then   no branch MUST contain two consecutive assignments to the same map key
```

---

## Bug #8 — Checksum failure silently disables deduplication

**File**: `internal/backup/snapshot.go`  
**Symbol**: `Snapshotter.Create` (line 82)  
**Severity**: WARNING

### Requirements

- When `ComputeChecksum` fails, the failure MUST be surfaced to the caller as a returned error
  rather than silenced with `log.Printf`.
- `Snapshotter.Create` MUST return `(Manifest{}, err)` when checksum computation fails.
- The `log.Printf` fallback that sets `checksum = ""` MUST be removed.
- Callers that previously tolerated a missing checksum MAY need to be updated if they exist; this
  spec does not constrain them.

### Scenarios

#### Scenario 8.1 — Checksum error is propagated

```
Given  one or more existing snapshot files
  And  ComputeChecksum returns a non-nil error
When   Snapshotter.Create() is called
Then   it MUST return (Manifest{}, <the checksum error>)
  And  no manifest file MUST be written to disk
```

#### Scenario 8.2 — Successful checksum is stored in manifest

```
Given  one or more existing snapshot files
  And  ComputeChecksum returns ("abc123", nil)
When   Snapshotter.Create() is called
Then   the returned Manifest.Checksum MUST equal "abc123"
  And  the manifest file MUST be written with Checksum == "abc123"
```

#### Scenario 8.3 — No existing files uses emptyFilesChecksum

```
Given  no existing snapshot files (len(existingPaths) == 0)
When   Snapshotter.Create() is called
Then   Manifest.Checksum MUST equal emptyFilesChecksum
  And  ComputeChecksum MUST NOT be called
```

---

## Bug #9 — String concatenation instead of `filepath.Join` for `.codex` path

**File**: `internal/components/kortex-engram/inject.go`  
**Symbol**: `writeCodexInstructionFiles` (line 384)  
**Severity**: WARNING

### Requirements

- The expression `homeDir + "/.codex"` MUST be replaced with `filepath.Join(homeDir, ".codex")`.
- Subsequent path constructions inside the same function that use string concatenation on
  `codexDir` MUST also be replaced with `filepath.Join` calls.
- The `filepath` package MUST already be imported; no new import is required.
- The change MUST produce identical paths on Linux (the primary target) and correct paths on
  Windows (a secondary target where path separators differ).

### Scenarios

#### Scenario 9.1 — Path is constructed correctly on Linux

```
Given  homeDir == "/home/lucas"
When   writeCodexInstructionFiles("/home/lucas") is called (or the path construction is tested)
Then   codexDir MUST equal "/home/lucas/.codex"
  And  instructionsPath MUST equal "/home/lucas/.codex/kortex-engram-instructions.md"
  And  compactPath MUST equal "/home/lucas/.codex/kortex-engram-compact-prompt.md"
```

#### Scenario 9.2 — filepath.Join is used (no string concat)

```
Given  the source of writeCodexInstructionFiles
When   it is inspected
Then   no string concatenation involving homeDir and a path separator MUST exist
  And  filepath.Join MUST be used for every path construction
```

---

## Bug #10 — `verify.RunChecks` called without timeout

**Files**: `internal/cli/sync.go` (line 848), `internal/cli/run.go` (line 1041)  
**Symbols**: `runPostSyncVerification`, post-install verify  
**Severity**: WARNING

### Requirements

- Every call to `verify.RunChecks(context.Background(), checks)` in CLI code MUST be replaced
  with a call that uses a context carrying a finite timeout.
- The timeout value SHOULD be configurable or at minimum a named constant (e.g.
  `verifyTimeout = 30 * time.Second`).
- Both call sites MUST use the same timeout strategy for consistency.
- If `verify.RunChecks` itself accepts a context, the timeout context MUST be passed directly.
- The cancel function from `context.WithTimeout` MUST be deferred to avoid context leak.

### Scenarios

#### Scenario 10.1 — RunChecks respects timeout (sync.go)

```
Given  verify.RunChecks is called in runPostSyncVerification
  And  checks take longer than the configured timeout
When   the timeout fires
Then   RunChecks MUST return with a context deadline exceeded error (or wrapped equivalent)
  And  the calling function MUST propagate that error to the caller
```

#### Scenario 10.2 — RunChecks completes within timeout (sync.go)

```
Given  verify.RunChecks is called with a 30-second timeout
  And  all checks complete in under 1 second
When   runPostSyncVerification executes
Then   it MUST return a BuildReport with no timeout error
```

#### Scenario 10.3 — RunChecks respects timeout (run.go)

```
Given  verify.RunChecks is called in the post-install verify path (run.go)
  And  checks take longer than the configured timeout
When   the timeout fires
Then   RunChecks MUST return with a context deadline exceeded error
  And  the post-install flow MUST handle it gracefully (log or surface to user)
```

#### Scenario 10.4 — cancel() is always called

```
Given  either call site
When   RunChecks returns (success or error)
Then   the cancel() function from context.WithTimeout MUST have been called
  And  no context goroutine leak MUST occur
```

---

## Test Implementation Notes

All tests MUST follow the project's standard table-driven pattern:

```go
tests := []struct {
    name     string
    input    <type>
    expected <type>
}{...}
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) { ... })
}
```

- Database tests (Bug #1): use an in-memory SQLite via `database/sql` + `modernc.org/sqlite`
  driver; inject cursor errors via a `sqlmock` or a custom `sql.Driver` stub.
- TUI tests (Bug #2): call `model.Update()` directly with synthetic `tea.Msg` values; do not
  use `teatest`.
- Mocks MUST be placed only at I/O boundaries (DB, filesystem, exec).
- Coverage target: `go test -cover ./...` MUST show no regression from pre-fix baseline.
