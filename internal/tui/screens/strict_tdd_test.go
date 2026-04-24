package screens

import (
	"strings"
	"testing"
)

func TestRenderStrictTDDContainsTitle(t *testing.T) {
	output := RenderStrictTDD(false, 0)
	if !strings.Contains(output, "TDD ESTRITO") && !strings.Contains(output, "STRICT TDD") {
		t.Errorf("RenderStrictTDD output missing title\ngot: %s", output)
	}
}

func TestRenderStrictTDDContainsEnableOption(t *testing.T) {
	output := RenderStrictTDD(false, 0)
	if !strings.Contains(output, "Ativar") && !strings.Contains(output, "Enable") {
		t.Errorf("RenderStrictTDD output missing enable option\ngot: %s", output)
	}
}

func TestRenderStrictTDDContainsDisableOption(t *testing.T) {
	output := RenderStrictTDD(false, 0)
	if !strings.Contains(output, "Desativar") && !strings.Contains(output, "Disable") {
		t.Errorf("RenderStrictTDD output missing disable option\ngot: %s", output)
	}
}

func TestRenderStrictTDDContainsBackOption(t *testing.T) {
	output := RenderStrictTDD(false, 0)
	if !strings.Contains(output, "Voltar") && !strings.Contains(output, "Back") {
		t.Errorf("RenderStrictTDD output missing back option\ngot: %s", output)
	}
}

func TestRenderStrictTDDEnabledState(t *testing.T) {
	output := RenderStrictTDD(true, 0)
	if !strings.Contains(output, "(*) Ativar") && !strings.Contains(output, "(*) Enable") {
		t.Errorf("RenderStrictTDD(enabled=true) should show enable as selected\ngot: %s", output)
	}
}

func TestRenderStrictTDDDisabledState(t *testing.T) {
	output := RenderStrictTDD(false, 0)
	if !strings.Contains(output, "(*) Desativar") && !strings.Contains(output, "(*) Disable") {
		t.Errorf("RenderStrictTDD(enabled=false) should show disable as selected\ngot: %s", output)
	}
}
