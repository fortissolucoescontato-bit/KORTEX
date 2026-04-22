package screens

import (
	"fmt"
	"strings"
	"testing"
)

// ─── RenderSync states ─────────────────────────────────────────────────────

// TestRenderSync_ConfirmState verifies the default confirm state — no operation
// running, no result yet — shows sync description and a prompt.
func TestRenderSync_ConfirmState(t *testing.T) {
	out := RenderSync(0, nil, false /*operationRunning*/, false /*hasSyncRun*/, 0)

	lower := strings.ToLower(out)
	if !strings.Contains(lower, "sincroniza") {
		t.Errorf("RenderSync(confirm) should contain 'sincronizar'; got:\n%s", out)
	}
	// Should show a prompt to press enter.
	if !strings.Contains(lower, "enter") && !strings.Contains(lower, "confirmar") {
		t.Errorf("RenderSync(confirm) should show enter/confirm prompt; got:\n%s", out)
	}
}

// TestRenderSync_RunningState verifies that while sync is running the screen
// shows a spinner/progress indicator.
func TestRenderSync_RunningState(t *testing.T) {
	out := RenderSync(0, nil, true /*operationRunning*/, false, 0)

	lower := strings.ToLower(out)
	if !strings.Contains(lower, "sincronizando") && !strings.Contains(lower, "por favor, aguarde") {
		t.Errorf("RenderSync(running) should show 'sincronizando' or 'por favor, aguarde'; got:\n%s", out)
	}
}

// TestRenderSync_ResultWithFilesChanged verifies that after a successful sync
// with changed files, the screen shows the file count.
func TestRenderSync_ResultWithFilesChanged(t *testing.T) {
	const filesChanged = 5
	out := RenderSync(filesChanged, nil, false, true /*hasSyncRun*/, 0)

	if !strings.Contains(out, "5") {
		t.Errorf("RenderSync(filesChanged=5) should show '5'; got:\n%s", out)
	}
	lower := strings.ToLower(out)
	if !strings.Contains(lower, "sincroniza") {
		t.Errorf("RenderSync(result) should mention 'sincronizar'; got:\n%s", out)
	}
}

// TestRenderSync_ResultWithError verifies that a failed sync shows the error
// message.
func TestRenderSync_ResultWithError(t *testing.T) {
	syncErr := fmt.Errorf("connection refused: agent config dir not writable")
	out := RenderSync(0, syncErr, false, true /*hasSyncRun*/, 0)

	lower := strings.ToLower(out)
	if !strings.Contains(lower, "falha") && !strings.Contains(lower, "erro") {
		t.Errorf("RenderSync(error) should show failure indicator; got:\n%s", out)
	}
	if !strings.Contains(out, syncErr.Error()) {
		t.Errorf("RenderSync(error) should show error text %q; got:\n%s", syncErr.Error(), out)
	}
}

// TestRenderSync_TitleAlwaysPresent verifies the screen title is shown in all
// states.
func TestRenderSync_TitleAlwaysPresent(t *testing.T) {
	states := []struct {
		name             string
		filesChanged     int
		syncErr          error
		operationRunning bool
		hasSyncRun       bool
	}{
		{"confirm", 0, nil, false, false},
		{"running", 0, nil, true, false},
		{"success", 3, nil, false, true},
		{"error", 0, fmt.Errorf("fail"), false, true},
	}

	for _, s := range states {
		t.Run(s.name, func(t *testing.T) {
			out := RenderSync(s.filesChanged, s.syncErr, s.operationRunning, s.hasSyncRun, 0)
			if !strings.Contains(out, "Sincronizar") {
				t.Errorf("RenderSync state=%q should contain 'Sincronizar'; got:\n%s", s.name, out)
			}
		})
	}
}

// TestRenderSync_ZeroFilesChangedWithNoError verifies the "nothing to update"
// case (hasSyncRun=true, filesChanged=0, no error) shows a completion message.
func TestRenderSync_ZeroFilesChangedWithNoError(t *testing.T) {
	out := RenderSync(0, nil, false, true /*hasSyncRun*/, 0)

	lower := strings.ToLower(out)
	if !strings.Contains(lower, "sincronização concluída") && !strings.Contains(lower, "concluída") &&
		!strings.Contains(lower, "nenhum agente") {
		t.Errorf("RenderSync(0 files, no error) should show completion; got:\n%s", out)
	}
}
