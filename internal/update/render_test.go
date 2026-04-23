package update

import (
	"fmt"
	"strings"
	"testing"
)

func TestRenderCLI_IncompleteCheckDoesNotClaimUpToDate(t *testing.T) {
	results := []UpdateResult{
		{Tool: ToolInfo{Name: "kortex"}, InstalledVersion: "1.0.0", LatestVersion: "1.0.0", Status: UpToDate},
		{Tool: ToolInfo{Name: "engram"}, Status: CheckFailed, Err: fmt.Errorf("timeout")},
	}

	out := RenderCLI(results)

	if strings.Contains(out, "Todas as ferramentas estão atualizadas!") {
		t.Fatalf("RenderCLI must not claim all tools are up to date when checks fail:\n%s", out)
	}
	if !strings.Contains(out, "Verificação incompleta") {
		t.Fatalf("RenderCLI must mention incomplete checks:\n%s", out)
	}
	if !strings.Contains(out, "falha na verificação") {
		t.Fatalf("RenderCLI must surface failed rows:\n%s", out)
	}
}

func TestCheckFailures(t *testing.T) {
	results := []UpdateResult{
		{Tool: ToolInfo{Name: "kortex"}, Status: UpToDate},
		{Tool: ToolInfo{Name: "engram"}, Status: CheckFailed},
		{Tool: ToolInfo{Name: "kortex"}, Status: CheckFailed},
	}

	failed := CheckFailures(results)
	if len(failed) != 2 {
		t.Fatalf("len(CheckFailures) = %d, want 2", len(failed))
	}
	if failed[0] != "engram" || failed[1] != "kortex" {
		t.Fatalf("CheckFailures() = %v, want [engram kortex]", failed)
	}
	if !HasCheckFailures(results) {
		t.Fatalf("HasCheckFailures() = false, want true")
	}
}
