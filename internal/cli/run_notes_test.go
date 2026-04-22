package cli

import (
	"strings"
	"testing"

	"github.com/fortissolucoescontato-bit/kortex/internal/model"
	"github.com/fortissolucoescontato-bit/kortex/internal/planner"
	"github.com/fortissolucoescontato-bit/kortex/internal/verify"
)

func TestWithPostInstallNotesAddsKortexCLINextSteps(t *testing.T) {
	report := verify.Report{Ready: true, FinalNote: "You're ready."}
	resolved := planner.ResolvedPlan{OrderedComponents: []model.ComponentID{model.ComponentKortexCLI}}

	updated := withPostInstallNotes(report, resolved)
	if !strings.Contains(updated.FinalNote, "KortexCLI is now installed globally") {
		t.Fatalf("FinalNote missing KortexCLI global install note: %q", updated.FinalNote)
	}
	if !strings.Contains(updated.FinalNote, "kortex init") || !strings.Contains(updated.FinalNote, "kortex install") {
		t.Fatalf("FinalNote missing KortexCLI repo setup steps: %q", updated.FinalNote)
	}
}

func TestWithPostInstallNotesDoesNotChangeNonKortexCLI(t *testing.T) {
	// Set GOBIN to a directory already in PATH so that withGoInstallPathNote
	// does not append a PATH guidance note for the Engram component.
	t.Setenv("GOBIN", "/usr/local/bin")

	report := verify.Report{Ready: true, FinalNote: "You're ready."}
	resolved := planner.ResolvedPlan{OrderedComponents: []model.ComponentID{model.ComponentEngram}}

	updated := withPostInstallNotes(report, resolved)
	if updated.FinalNote != report.FinalNote {
		t.Fatalf("FinalNote changed unexpectedly: %q", updated.FinalNote)
	}
}
