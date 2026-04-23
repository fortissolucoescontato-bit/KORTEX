package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fortissolucoescontato-bit/kortex/internal/components/kortex-engram"
	"github.com/fortissolucoescontato-bit/kortex/internal/system"
)

// TestRunInstallLinuxKortexEngramUsesDownloadNotGoInstall verifies that after the fix,
// Linux KortexEngram installation does NOT use "go install" but instead calls
// DownloadLatestBinary (i.e. no "go install" in recorder.get()).
func TestRunInstallLinuxKortexEngramUsesDownloadNotGoInstall(t *testing.T) {
	home := t.TempDir()
	restoreHome := osUserHomeDir
	restoreCommand := runCommand
	restoreLookPath := cmdLookPath
	t.Cleanup(func() {
		osUserHomeDir = restoreHome
		runCommand = restoreCommand
		cmdLookPath = restoreLookPath
	})

	osUserHomeDir = func() (string, error) { return home, nil }
	cmdLookPath = missingBinaryLookPath
	recorder := &commandRecorder{}
	runCommand = recorder.record

	// Override the KortexEngram download function to succeed without hitting GitHub.
	origDownloadFn := KortexEngramDownloadFn
	KortexEngramDownloadFn = func(profile system.PlatformProfile) (string, error) {
		// Simulate a successful binary download to a temp path.
		return "/tmp/fake-KortexEngram", nil
	}
	t.Cleanup(func() { KortexEngramDownloadFn = origDownloadFn })

	detection := linuxDetectionResult(system.LinuxDistroUbuntu, "apt")
	result, err := RunInstall(
		[]string{"--agent", "opencode", "--component", "kortex-engram"},
		detection,
	)
	if err != nil {
		t.Fatalf("RunInstall() error = %v", err)
	}

	if !result.Verify.Ready {
		t.Fatalf("verification ready = false, report = %#v", result.Verify)
	}

	// Must NOT have called "go install" for kortexengram.
	for _, cmd := range recorder.get() {
		if strings.Contains(cmd, "go install") && strings.Contains(cmd, "kortex-engram") {
			t.Fatalf("Linux KortexEngram install should NOT use go install, got command: %s", cmd)
		}
	}
}

// TestRunInstallKortexEngramDownloadAddsBinDirToPath verifies that after downloading
// the KortexEngram binary, its directory is prepended to PATH so that subsequent
// commands (KortexEngram setup, resolveKortexEngramCommand) can find it.
func TestRunInstallKortexEngramDownloadAddsBinDirToPath(t *testing.T) {
	home := t.TempDir()
	restoreHome := osUserHomeDir
	restoreCommand := runCommand
	restoreLookPath := cmdLookPath
	restorePath := os.Getenv("PATH")
	t.Cleanup(func() {
		osUserHomeDir = restoreHome
		runCommand = restoreCommand
		cmdLookPath = restoreLookPath
		os.Setenv("PATH", restorePath)
	})

	osUserHomeDir = func() (string, error) { return home, nil }
	cmdLookPath = missingBinaryLookPath
	recorder := &commandRecorder{}
	runCommand = recorder.record

	fakeBinDir := filepath.Join(home, "KortexEngram-bin")
	os.MkdirAll(fakeBinDir, 0o755)
	fakeBinaryPath := filepath.Join(fakeBinDir, "kortex-engram")

	origDownloadFn := KortexEngramDownloadFn
	KortexEngramDownloadFn = func(profile system.PlatformProfile) (string, error) {
		return fakeBinaryPath, nil
	}
	t.Cleanup(func() { KortexEngramDownloadFn = origDownloadFn })

	detection := linuxDetectionResult(system.LinuxDistroUbuntu, "apt")
	_, err := RunInstall(
		[]string{"--agent", "opencode", "--component", "kortex-engram"},
		detection,
	)
	if err != nil {
		t.Fatalf("RunInstall() error = %v", err)
	}

	currentPath := os.Getenv("PATH")
	if !strings.Contains(currentPath, fakeBinDir) {
		t.Fatalf("PATH should contain KortexEngram bin dir %q after download, got PATH=%q", fakeBinDir, currentPath)
	}
}

// TestRunInstallWindowsKortexEngramUsesDownloadNotGoInstall verifies Windows path.
func TestRunInstallWindowsKortexEngramUsesDownloadNotGoInstall(t *testing.T) {
	home := t.TempDir()
	restoreHome := osUserHomeDir
	restoreCommand := runCommand
	restoreLookPath := cmdLookPath
	t.Cleanup(func() {
		osUserHomeDir = restoreHome
		runCommand = restoreCommand
		cmdLookPath = restoreLookPath
	})

	osUserHomeDir = func() (string, error) { return home, nil }
	cmdLookPath = missingBinaryLookPath
	recorder := &commandRecorder{}
	runCommand = recorder.record

	origDownloadFn := KortexEngramDownloadFn
	KortexEngramDownloadFn = func(profile system.PlatformProfile) (string, error) {
		return `C:\fake\kortexengram.exe`, nil
	}
	t.Cleanup(func() { KortexEngramDownloadFn = origDownloadFn })

	detection := system.DetectionResult{
		System: system.SystemInfo{
			OS:        "windows",
			Arch:      "amd64",
			Supported: true,
			Profile: system.PlatformProfile{
				OS:             "windows",
				PackageManager: "winget",
				Supported:      true,
			},
		},
	}

	result, err := RunInstall(
		[]string{"--agent", "opencode", "--component", "kortex-engram"},
		detection,
	)
	if err != nil {
		t.Fatalf("RunInstall() error = %v", err)
	}

	if !result.Verify.Ready {
		t.Fatalf("verification ready = false, report = %#v", result.Verify)
	}

	// Must NOT have called "go install" for kortexengram.
	for _, cmd := range recorder.get() {
		if strings.Contains(cmd, "go install") && strings.Contains(cmd, "kortex-engram") {
			t.Fatalf("Windows KortexEngram install should NOT use go install, got command: %s", cmd)
		}
	}
}

// TestRunInstallMacOSKortexEngramStillUsesBrew verifies macOS unchanged.
func TestRunInstallMacOSKortexEngramStillUsesBrew(t *testing.T) {
	home := t.TempDir()
	restoreHome := osUserHomeDir
	restoreCommand := runCommand
	restoreLookPath := cmdLookPath
	t.Cleanup(func() {
		osUserHomeDir = restoreHome
		runCommand = restoreCommand
		cmdLookPath = restoreLookPath
	})

	osUserHomeDir = func() (string, error) { return home, nil }
	cmdLookPath = missingBinaryLookPath
	recorder := &commandRecorder{}
	runCommand = recorder.record

	// DownloadFn should NOT be called for macOS (brew handles it).
	origDownloadFn := KortexEngramDownloadFn
	KortexEngramDownloadFn = func(profile system.PlatformProfile) (string, error) {
		t.Error("DownloadLatestBinary should NOT be called on macOS (brew handles it)")
		return "", nil
	}
	t.Cleanup(func() { KortexEngramDownloadFn = origDownloadFn })

	detection := macOSDetectionResult()
	result, err := RunInstall(
		[]string{"--agent", "opencode", "--component", "kortex-engram"},
		detection,
	)
	if err != nil {
		t.Fatalf("RunInstall() error = %v", err)
	}
	if !result.Verify.Ready {
		t.Fatalf("verification ready = false")
	}

	// Must use brew install kortexengram.
	commands := recorder.get()
	foundBrew := false
	for _, cmd := range commands {
		if strings.Contains(cmd, "brew install KortexEngram") {
			foundBrew = true
		}
	}
	if !foundBrew {
		t.Fatalf("expected brew install KortexEngram on macOS, got commands: %v", commands)
	}
}

// Make sure the KortexEngram package's DownloadLatestBinary is accessible.
var _ = kortexengram.DownloadLatestBinary
