# Design: Kortex CLI PowerShell Support

## Technical Approach

Embed `kortex.ps1` as a Go asset under `internal/assets/kortex/`, then extend `runtime.go` with a new exported function `EnsurePowerShellShim(homeDir string) error` that writes the shim using the same `assets.Read` + `filemerge.WriteFileAtomic` pattern already used for `pr_mode.sh`. The install trigger is added in `resolver.go` `resolveKortex CLIInstall` for `winget`, appending a post-install step that invokes a new `internal/components/kortex` exported helper via a dedicated `WriteKortex CLIShim` command — but because `CommandSequence` is shell-level and the shim write is a Go-level side-effect, it does NOT go into `CommandSequence`. Instead, `EnsurePowerShellShim` is called from the same call-site that calls `EnsureRuntimeAssets`, guarded by `runtime.GOOS == "windows"` (via an injected predicate for testability).

## Architecture Decisions

| Option | Tradeoff | Decision |
|--------|----------|----------|
| Write shim inside `EnsureRuntimeAssets` with OS guard | Keeps all runtime asset logic in one function; couples Linux/macOS path to Windows code | Rejected — violates single-responsibility |
| New `EnsurePowerShellShim(homeDir string) error` in `runtime.go` | Mirrors `EnsureRuntimeAssets` exactly; caller guards with OS check; easy to test in isolation | **Chosen** |
| Add PS1 step to `CommandSequence` in `resolver.go` | Keeps all install steps in one place; but shell-level commands cannot call Go's `os.Rename` atomic write | Rejected — shim write must be Go-level |
| Embed `kortex.ps1` template with `gitBashPath()` expanded at runtime | Dynamic path avoids hardcoding; but `gitBashPath()` lives in `installcmd` package, not accessible from assets | Rejected — `kortex.ps1` is a static shim that calls `git.exe`-relative `bash.exe` via a lookup at run time inside the script itself |

**Final choice for shim content:** `kortex.ps1` resolves Git Bash by deriving its path from `(Get-Command git).Source` at PowerShell runtime (not at install time). This eliminates the cross-package dependency on `gitBashPath()` and produces a shim that self-heals if Git is reinstalled to a different path. The shim is therefore a fully static asset — no templating needed.

## Data Flow

```
kortex install (Windows)
  │
  ├─ resolveKortex CLIInstall(winget)  ──→  CommandSequence (powershell cleanup + git clone + bash install.sh)
  │                                   [existing — unchanged]
  │
  └─ EnsurePowerShellShim(homeDir)   [NEW — called after Kortex CLI install completes]
       │
       ├─ assets.Read("kortex/kortex.ps1")          ← embedded asset
       ├─ RuntimeBinDir(homeDir)               ← ~/.local/share/kortex/bin
       └─ filemerge.WriteFileAtomic(path, content, 0o755)
            ├─ no-op if content matches
            └─ atomic rename if stale/missing
```

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `internal/assets/kortex/kortex.ps1` | Create | Static PowerShell shim; resolves bash via `git` on PATH at runtime |
| `internal/components/kortex/runtime.go` | Modify | Add `RuntimeBinDir`, `RuntimePS1Path`, `EnsurePowerShellShim` |
| `internal/components/kortex/runtime_test.go` | Modify | Add tests for shim: missing, stale, idempotent, git-bash-not-found |
| `docs/platforms.md` | Modify | Remove Windows/PowerShell limitation note |

`install.go` and `resolver.go` are NOT modified — the shim install is a runtime-asset concern, not a package-manager command concern.

## Interfaces / Contracts

```go
// RuntimeBinDir returns ~/.local/share/kortex/bin — where Kortex CLI's bash script lives on Linux/Windows.
func RuntimeBinDir(homeDir string) string

// RuntimePS1Path returns the expected kortex.ps1 path.
func RuntimePS1Path(homeDir string) string

// EnsurePowerShellShim writes kortex.ps1 to the Kortex CLI bin directory.
// Uses WriteFileAtomic: no-op when content matches, atomic replace otherwise.
// Must only be called on Windows (caller is responsible for the OS guard).
func EnsurePowerShellShim(homeDir string) error
```

`kortex.ps1` content (static asset — no templating):

```powershell
$gitCmd = Get-Command git -ErrorAction SilentlyContinue
if (-not $gitCmd) {
    Write-Error "Git not found on PATH. Install Git for Windows to use kortex from PowerShell."
    exit 1
}
$bash = Join-Path (Split-Path (Split-Path $gitCmd.Source)) "bin\bash.exe"
if (-not (Test-Path $bash)) {
    Write-Error "Git Bash not found at '$bash'. Reinstall Git for Windows."
    exit 1
}
& $bash -c "kortex $args"
exit $LASTEXITCODE
```

**Note on argument forwarding:** `$args` in PowerShell is the automatic array of unbound arguments. Using `"kortex $args"` passes them as a space-joined string to bash `-c`. For arguments with spaces, callers must quote them in PowerShell as they normally would — this matches standard shell forwarding behavior and is consistent with how similar PS1 shims (e.g., nvm.ps1) work.

## Testing Strategy

| Layer | What to Test | Approach |
|-------|-------------|----------|
| Unit — asset | `kortex.ps1` is embedded and readable | `assets.Read("kortex/kortex.ps1")` returns non-empty content |
| Unit — `EnsurePowerShellShim` | Creates file when missing | `t.TempDir()` home, call once, assert file exists with expected content |
| Unit — `EnsurePowerShellShim` | Overwrites stale shim | Write sentinel content, call, assert content replaced |
| Unit — `EnsurePowerShellShim` | No-op when content matches | Call twice, assert `ModTime` unchanged (mirrors `TestEnsureRuntimeAssetsIsNoOpWhenContentMatches`) |
| Unit — path helpers | `RuntimeBinDir` / `RuntimePS1Path` | Table test with known homeDir, assert expected suffix |

No integration or E2E tests are required for this change — the `filemerge.WriteFileAtomic` path is already covered by its own test suite, and PowerShell execution is an OS-level concern outside unit test scope.

## Migration / Rollout

No migration required. On first run after update, `EnsurePowerShellShim` creates the file. On non-Windows, the function is never called. Rollback: delete `~/.local/share/kortex/bin/kortex.ps1`.

## Open Questions

- [ ] Confirm where `EnsurePowerShellShim` is called from — the tasks phase must identify the exact call-site (likely the same place `EnsureRuntimeAssets` is invoked) and add the Windows OS guard there.
- [ ] Verify that Kortex CLI's `install.sh` on Windows places the `kortex` bash script under `~/.local/share/kortex/bin/` (same as Linux) so `RuntimeBinDir` is the correct target directory.
