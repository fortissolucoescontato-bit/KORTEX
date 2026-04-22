package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fortissolucoescontato-bit/kortex/internal/system"
)

// TestKortexCLIAvailableDetectsViaLookPath verifies that kortexAvailable returns true
// when kortex is found on PATH via cmdLookPath.
func TestKortexCLIAvailableDetectsViaLookPath(t *testing.T) {
	origLookPath := cmdLookPath
	cmdLookPath = func(file string) (string, error) {
		if file == "kortex" {
			return "/usr/local/bin/kortex", nil
		}
		return "", os.ErrNotExist
	}
	t.Cleanup(func() { cmdLookPath = origLookPath })

	if !kortexAvailable(system.PlatformProfile{OS: "darwin", PackageManager: "brew"}) {
		t.Fatal("kortexAvailable() = false, want true when kortex is on PATH")
	}
}

// TestKortexCLIAvailableDetectsViaLocalBin verifies that kortexAvailable returns true
// when kortex exists at ~/.local/bin/kortex (default for install.sh on Linux/macOS).
func TestKortexCLIAvailableDetectsViaLocalBin(t *testing.T) {
	tmpHome := t.TempDir()
	localBin := filepath.Join(tmpHome, ".local", "bin")
	if err := os.MkdirAll(localBin, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(localBin, "kortex"), []byte("fake"), 0o755); err != nil {
		t.Fatal(err)
	}

	origLookPath := cmdLookPath
	origHomeDir := osUserHomeDir
	origStat := osStat
	cmdLookPath = func(file string) (string, error) { return "", os.ErrNotExist }
	osUserHomeDir = func() (string, error) { return tmpHome, nil }
	osStat = os.Stat
	t.Cleanup(func() {
		cmdLookPath = origLookPath
		osUserHomeDir = origHomeDir
		osStat = origStat
	})

	if !kortexAvailable(system.PlatformProfile{OS: "linux", PackageManager: "apt"}) {
		t.Fatal("kortexAvailable() = false, want true when kortex is at ~/.local/bin/kortex")
	}
}

// TestKortexCLIAvailableDetectsViaHomebrewOptPrefix verifies that kortexAvailable returns
// true when kortex exists at /opt/homebrew/bin/kortex (Apple Silicon Homebrew default).
func TestKortexCLIAvailableDetectsViaHomebrewOptPrefix(t *testing.T) {
	tmpDir := t.TempDir()
	fakeOptHomebrew := filepath.Join(tmpDir, "opt", "homebrew", "bin", "kortex")
	if err := os.MkdirAll(filepath.Dir(fakeOptHomebrew), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(fakeOptHomebrew, []byte("fake"), 0o755); err != nil {
		t.Fatal(err)
	}

	origLookPath := cmdLookPath
	origHomeDir := osUserHomeDir
	origStat := osStat
	cmdLookPath = func(file string) (string, error) { return "", os.ErrNotExist }
	osUserHomeDir = func() (string, error) { return tmpDir, nil }
	// Override osStat to redirect well-known brew paths to our temp dir.
	osStat = func(name string) (os.FileInfo, error) {
		switch name {
		case "/opt/homebrew/bin/kortex":
			return os.Stat(fakeOptHomebrew)
		case "/usr/local/bin/kortex":
			return nil, os.ErrNotExist
		default:
			return os.Stat(name)
		}
	}
	t.Cleanup(func() {
		cmdLookPath = origLookPath
		osUserHomeDir = origHomeDir
		osStat = origStat
	})

	if !kortexAvailable(system.PlatformProfile{OS: "darwin", PackageManager: "brew"}) {
		t.Fatal("kortexAvailable() = false, want true when kortex is at /opt/homebrew/bin/kortex")
	}
}

// TestKortexCLIAvailableDetectsViaHomebrewUsrLocalPrefix verifies that kortexAvailable
// returns true when kortex exists at /usr/local/bin/kortex (Intel Mac Homebrew default).
func TestKortexCLIAvailableDetectsViaHomebrewUsrLocalPrefix(t *testing.T) {
	origLookPath := cmdLookPath
	origHomeDir := osUserHomeDir
	origStat := osStat
	cmdLookPath = func(file string) (string, error) { return "", os.ErrNotExist }
	osUserHomeDir = func() (string, error) { return t.TempDir(), nil }
	osStat = func(name string) (os.FileInfo, error) {
		switch name {
		case "/opt/homebrew/bin/kortex":
			return nil, os.ErrNotExist
		case "/usr/local/bin/kortex":
			// Simulate kortex present here.
			return os.Stat(os.DevNull)
		default:
			return nil, os.ErrNotExist
		}
	}
	t.Cleanup(func() {
		cmdLookPath = origLookPath
		osUserHomeDir = origHomeDir
		osStat = origStat
	})

	if !kortexAvailable(system.PlatformProfile{OS: "darwin", PackageManager: "brew"}) {
		t.Fatal("kortexAvailable() = false, want true when kortex is at /usr/local/bin/kortex")
	}
}

// TestKortexCLIAvailableReturnsFalseWhenNotFound verifies that kortexAvailable returns
// false when kortex is not found via any detection path.
func TestKortexCLIAvailableReturnsFalseWhenNotFound(t *testing.T) {
	origLookPath := cmdLookPath
	origHomeDir := osUserHomeDir
	origStat := osStat
	cmdLookPath = func(file string) (string, error) { return "", os.ErrNotExist }
	osUserHomeDir = func() (string, error) { return t.TempDir(), nil }
	osStat = func(name string) (os.FileInfo, error) { return nil, os.ErrNotExist }
	t.Cleanup(func() {
		cmdLookPath = origLookPath
		osUserHomeDir = origHomeDir
		osStat = origStat
	})

	if kortexAvailable(system.PlatformProfile{OS: "darwin", PackageManager: "brew"}) {
		t.Fatal("kortexAvailable() = true, want false when kortex is not installed anywhere")
	}
}

// TestKortexCLIAvailableBrewPathsSkippedOnLinux verifies that the Homebrew-specific
// paths (/opt/homebrew/bin/kortex, /usr/local/bin/kortex) are NOT checked on Linux
// even if those paths happen to exist (they never exist there in practice, but
// the guard ensures no cross-platform false positives).
func TestKortexCLIAvailableBrewPathsSkippedOnLinux(t *testing.T) {
	origLookPath := cmdLookPath
	origHomeDir := osUserHomeDir
	origStat := osStat
	cmdLookPath = func(file string) (string, error) { return "", os.ErrNotExist }
	osUserHomeDir = func() (string, error) { return t.TempDir(), nil }

	statCallCount := 0
	osStat = func(name string) (os.FileInfo, error) {
		if name == "/opt/homebrew/bin/kortex" || name == "/usr/local/bin/kortex" {
			statCallCount++
		}
		return nil, os.ErrNotExist
	}
	t.Cleanup(func() {
		cmdLookPath = origLookPath
		osUserHomeDir = origHomeDir
		osStat = origStat
	})

	kortexAvailable(system.PlatformProfile{OS: "linux", PackageManager: "apt"})
	if statCallCount > 0 {
		t.Fatalf("kortexAvailable() checked Homebrew paths on Linux (%d calls), expected 0", statCallCount)
	}
}
