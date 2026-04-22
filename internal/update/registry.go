package update

// Tools is the static registry of managed tools that can be checked for updates.
//
// InstallMethod controls which upgrade strategy the executor uses:
//   - InstallBrew: managed via homebrew (macOS/Linux with brew)
//   - InstallGoInstall: installed via `go install <GoImportPath>@version`
//   - InstallBinary: downloaded binary from GitHub Releases (atomic replace)
//
// For brew-managed platforms the executor picks brew regardless of the
// field here; InstallMethod represents the non-brew fallback strategy.
var Tools = []ToolInfo{
	{
		Name:          "kortex",
		Owner:         "fortissolucoescontato-bit",
		Repo:          "kortex",
		DetectCmd:     nil, // version comes from build-time ldflags (app.Version)
		VersionPrefix: "v",
		// kortex: brew on macOS, binary release download on Linux/Windows.
		// Self-upgrade of the running binary on Windows is deferred to Phase 2.
		InstallMethod: InstallBinary,
	},
	{
		Name:          "engram",
		Owner:         "fortissolucoescontato-bit",
		Repo:          "engram",
		DetectCmd:     []string{"engram", "version"},
		VersionPrefix: "v",
		// engram: brew on macOS/Linux-brew, binary download elsewhere.
		InstallMethod: InstallBinary,
	},
}
