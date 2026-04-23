package screens

import (
	"strings"
	"testing"

	componentuninstall "github.com/fortissolucoescontato-bit/kortex/internal/components/uninstall"
	"github.com/fortissolucoescontato-bit/kortex/internal/model"
)

func TestRenderUninstallResultIncludesManualCleanup(t *testing.T) {
	out := RenderUninstallResult(componentuninstall.Result{
		RemovedDirectories: []string{"/tmp/skills"},
		ManualActions: []string{
			"Remove manually if no longer needed: /tmp/skills (directory still contains non-managed files)",
		},
	}, nil, "", nil, model.KortexEngramUninstallScopeGlobal, false, 0, nil)

	if !strings.Contains(out, "Limpeza manual necessária") {
		t.Fatalf("RenderUninstallResult() should include manual cleanup heading; got:\n%s", out)
	}
	if !strings.Contains(out, "/tmp/skills") {
		t.Fatalf("RenderUninstallResult() should include manual cleanup item; got:\n%s", out)
	}
}

func TestRenderUninstallConfirmIncludesSelectedProfiles(t *testing.T) {
	out := RenderUninstallConfirm(
		model.UninstallModePartial,
		[]model.AgentID{model.AgentOpenCode},
		[]model.ComponentID{model.ComponentSDD},
		[]string{"cheap"},
		model.KortexEngramUninstallScopeGlobal,
		false,
		0,
		false,
		0,
	)

	if !strings.Contains(out, "Perfis a serem removidos") {
		t.Fatalf("RenderUninstallConfirm() should include profile section; got:\n%s", out)
	}
	if !strings.Contains(out, "cheap") {
		t.Fatalf("RenderUninstallConfirm() should include selected profile name; got:\n%s", out)
	}
}

func TestRenderUninstallConfirmIncludesKortexEngramProjectScopeDetails(t *testing.T) {
	out := RenderUninstallConfirm(
		model.UninstallModePartial,
		[]model.AgentID{model.AgentOpenCode},
		[]model.ComponentID{model.ComponentKortexEngram},
		nil,
		model.KortexEngramUninstallScopeProject,
		true,
		0,
		false,
		0,
	)

	if !strings.Contains(out, "Escopo de limpeza do KortexEngram") {
		t.Fatalf("RenderUninstallConfirm() should include KortexEngram cleanup scope heading; got:\n%s", out)
	}
	if !strings.Contains(out, "Apenas Projeto") {
		t.Fatalf("RenderUninstallConfirm() should include project-only scope label; got:\n%s", out)
	}
	if !strings.Contains(out, ".KortexEngram/") {
		t.Fatalf("RenderUninstallConfirm() should mention .KortexEngram project data removal; got:\n%s", out)
	}
}

func TestRenderUninstallResultIncludesSelectedProfiles(t *testing.T) {
	out := RenderUninstallResult(componentuninstall.Result{}, nil, model.UninstallModePartial, []string{"cheap", "fast"}, model.KortexEngramUninstallScopeGlobal, false, 0, nil)

	if !strings.Contains(out, "Perfis removidos") {
		t.Fatalf("RenderUninstallResult() should include profile summary heading; got:\n%s", out)
	}
	if !strings.Contains(out, "cheap") || !strings.Contains(out, "fast") {
		t.Fatalf("RenderUninstallResult() should include selected profile names; got:\n%s", out)
	}
}

func TestRenderUninstallResultIncludesKortexEngramScopeSummary(t *testing.T) {
	out := RenderUninstallResult(componentuninstall.Result{
		RemovedDirectories: []string{"/tmp/workspace/.KortexEngram"},
	}, nil, model.UninstallModePartial, nil, model.KortexEngramUninstallScopeProject, true, 0, nil)

	if !strings.Contains(out, "Escopo KortexEngram: Apenas Projeto") {
		t.Fatalf("RenderUninstallResult() should include KortexEngram project scope summary; got:\n%s", out)
	}
}
