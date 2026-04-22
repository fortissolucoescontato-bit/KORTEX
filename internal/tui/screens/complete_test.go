package screens

import (
	"strings"
	"testing"
)

func TestRenderCompleteSuccessShowsKortexCLINotesWhenInstalled(t *testing.T) {
	out := RenderComplete(CompletePayload{
		ConfiguredAgents:    1,
		InstalledComponents: 1,
		KortexCLIInstalled:        true,
	})

	if !strings.Contains(out, "KortexCLI (per project)") {
		t.Fatalf("missing KortexCLI section: %q", out)
	}
	if !strings.Contains(out, "kortex init") || !strings.Contains(out, "kortex install") {
		t.Fatalf("missing KortexCLI repo commands: %q", out)
	}
}

func TestRenderCompleteSuccessHidesKortexCLINotesWhenNotInstalled(t *testing.T) {
	out := RenderComplete(CompletePayload{
		ConfiguredAgents:    1,
		InstalledComponents: 1,
		KortexCLIInstalled:        false,
	})

	if strings.Contains(out, "KortexCLI (per project)") {
		t.Fatalf("unexpected KortexCLI section: %q", out)
	}
}
