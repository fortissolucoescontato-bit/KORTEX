package cli

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fortissolucoescontato-bit/kortex/internal/model"
	"github.com/fortissolucoescontato-bit/kortex/internal/system"
)

// TestKortexCLIFixInstallErrorWhenAlreadyAvailable tests that when KortexCLI install
// command fails but KortexCLI is already available on the system, the error is
// swallowed and the pipeline continues instead of failing.
// This simulates the Windows scenario where install.sh fails due to TTY
// issues but KortexCLI is already present.
func TestKortexCLIFixInstallErrorWhenAlreadyAvailable(t *testing.T) {
	home := t.TempDir()

	// Save original function references
	origHome := osUserHomeDir
	origCmdLookPath := cmdLookPath
	origRunCmd := runCommand
	origKortexCLIAvailableCheck := kortexAvailableCheck

	t.Cleanup(func() {
		osUserHomeDir = origHome
		cmdLookPath = origCmdLookPath
		runCommand = origRunCmd
		kortexAvailableCheck = origKortexCLIAvailableCheck
	})

	// Setup mocks
	osUserHomeDir = func() (string, error) { return home, nil }
	cmdLookPath = func(string) (string, error) {
		return "", errors.New("not found")
	}

	// Track if runCommand was called and capture its error
	runCommandCalled := false
	runCommand = func(name string, args ...string) error {
		runCommandCalled = true
		// Simulate install.sh failing due to TTY issue
		return errors.New("exit status 1: read: open /dev/tty: no such device or address")
	}

	// Make kortexAvailable return false initially (simulating install needed),
	// then return true after the "install" (simulating KortexCLI was already there)
	kortexAvailableCheck = func(profile system.PlatformProfile) bool {
		// After install command runs, simulate KortexCLI being available
		// (this is the fix scenario: install failed but KortexCLI is there)
		return runCommandCalled
	}

	// Create a minimal config so the test can run
	configPath := filepath.Join(home, ".config", "opencode", "opencode.json")
	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(configPath, []byte("{}"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	// Create the runtime manually to test the component step
	profile := system.PlatformProfile{OS: "windows", PackageManager: "winget"}
	step := componentApplyStep{
		id:           "component:kortex",
		component:    model.ComponentKortexCLI,
		homeDir:      home,
		workspaceDir: home,
		agents:       []model.AgentID{model.AgentOpenCode},
		selection:    model.Selection{},
		profile:      profile,
	}

	err := step.Run()

	// Verify: no error should be returned (fix: error swallowed when KortexCLI available)
	if err != nil {
		t.Fatalf("Run() error = %v, want nil (error should be swallowed when KortexCLI is available)", err)
	}

	// Verify: runCommand was called (we attempted install)
	if !runCommandCalled {
		t.Fatal("runCommand was not called, expected install to be attempted")
	}
}

// TestKortexCLIFixInstallErrorWhenNotAvailable tests that when KortexCLI install
// command fails and KortexCLI is NOT available, the error is NOT swallowed
// and is returned to the caller. This ensures we don't mask real errors.
func TestKortexCLIFixInstallErrorWhenNotAvailable(t *testing.T) {
	home := t.TempDir()

	origHome := osUserHomeDir
	origRunCmd := runCommand
	origKortexCLIAvailableCheck := kortexAvailableCheck
	origCmdLookPath := cmdLookPath

	t.Cleanup(func() {
		osUserHomeDir = origHome
		runCommand = origRunCmd
		kortexAvailableCheck = origKortexCLIAvailableCheck
		cmdLookPath = origCmdLookPath
	})

	osUserHomeDir = func() (string, error) { return home, nil }
	cmdLookPath = func(string) (string, error) {
		return "", errors.New("not found")
	}

	// Make kortexAvailable ALWAYS return false (KortexCLI is not available)
	kortexAvailableCheck = func(profile system.PlatformProfile) bool {
		return false
	}

	// Simulate a REAL install error (not the TTY issue)
	runCommand = func(name string, args ...string) error {
		return errors.New("network error: connection refused")
	}

	configPath := filepath.Join(home, ".config", "opencode", "opencode.json")
	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(configPath, []byte("{}"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	profile := system.PlatformProfile{OS: "windows", PackageManager: "winget"}
	step := componentApplyStep{
		id:           "component:kortex",
		component:    model.ComponentKortexCLI,
		homeDir:      home,
		workspaceDir: home,
		agents:       []model.AgentID{model.AgentOpenCode},
		selection:    model.Selection{},
		profile:      profile,
	}

	err := step.Run()

	// Verify: error should be returned (not swallowed)
	if err == nil {
		t.Fatal("Run() expected error when KortexCLI is not available and install fails, got nil")
	}

	if !strings.Contains(err.Error(), "network error") {
		t.Fatalf("Expected network error in message, got: %v", err)
	}
}
