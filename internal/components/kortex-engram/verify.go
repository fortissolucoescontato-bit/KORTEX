package kortexengram

import (
	"context"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

var (
	lookPath    = exec.LookPath
	execCommand = exec.Command
)

func VerifyInstalled() error {
	if _, err := lookPath("kortex-engram"); err == nil {
		return nil
	}
	if _, err := lookPath("kortex"); err == nil {
		return nil
	}
	if _, err := lookPath("kortex-engram"); err != nil {
		return fmt.Errorf("neither 'KortexEngram', 'kortex' nor 'KortexEngram' binary found in PATH: %w", err)
	}

	return nil
}

// VerifyVersion runs "KortexEngram version" and returns the trimmed output.
// Returns an error if the command fails or produces no output.
func VerifyVersion() (string, error) {
	cmdName := "kortex-engram"
	if _, err := lookPath(cmdName); err != nil {
		cmdName = "kortex"
		if _, err := lookPath(cmdName); err != nil {
			cmdName = "kortex-engram"
		}
	}

	cmd := execCommand(cmdName, "version")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("%s version command failed: %w", cmdName, err)
	}

	version := strings.TrimSpace(string(out))
	if version == "" {
		return "", fmt.Errorf("%s version returned empty output", cmdName)
	}

	return version, nil
}

func VerifyHealth(ctx context.Context, baseURL string) error {
	if strings.TrimSpace(baseURL) == "" {
		baseURL = "http://127.0.0.1:7437"
	}

	client := &http.Client{Timeout: 2 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, strings.TrimRight(baseURL, "/")+"/health", nil)
	if err != nil {
		return fmt.Errorf("build KortexEngram health request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("KortexEngram health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("KortexEngram health check returned status %d", resp.StatusCode)
	}

	return nil
}
