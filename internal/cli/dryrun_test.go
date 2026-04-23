package cli

import (
	"strings"
	"testing"

	"github.com/fortissolucoescontato-bit/kortex/internal/model"
	"github.com/fortissolucoescontato-bit/kortex/internal/planner"
)

func TestRenderDryRunIncludesPlatformDecision(t *testing.T) {
	result := InstallResult{
		Selection: model.Selection{Persona: model.PersonaKortex, Preset: model.PresetFullKortex},
		Resolved: planner.ResolvedPlan{
			Agents:            []model.AgentID{model.AgentClaudeCode},
			OrderedComponents: []model.ComponentID{model.ComponentEngram},
		},
		Review: planner.ReviewPayload{
			PlatformDecision: planner.PlatformDecision{
				OS:             "linux",
				LinuxDistro:    "ubuntu",
				PackageManager: "apt",
				Supported:      true,
			},
		},
	}

	output := RenderDryRun(result)

	want := "Decisão de plataforma: os=linux distro=ubuntu package-manager=apt status=suportado"
	if !strings.Contains(output, want) {
		t.Fatalf("RenderDryRun() missing platform decision\noutput=%s", output)
	}
}
