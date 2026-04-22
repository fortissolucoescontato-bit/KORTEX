package screens

import (
	"strings"

	"github.com/fortissolucoescontato-bit/kortex/internal/tui/styles"
)

// RenderABInstalling renders the installation-in-progress (or error) screen.
func RenderABInstalling(engineName string, spinnerFrame int, installErr error) string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Instalando seu Agente..."))
	b.WriteString("\n\n")

	if installErr != nil {
		b.WriteString(styles.ErrorStyle.Render("✗ Falha na instalação"))
		b.WriteString("\n")
		b.WriteString(styles.SubtextStyle.Render("  Motor: " + engineName))
		b.WriteString("\n")
		b.WriteString(styles.ErrorStyle.Render("  Erro: " + installErr.Error()))
		b.WriteString("\n\n")
		b.WriteString(renderOptions([]string{"Tentar novamente", "Voltar"}, 0))
		b.WriteString("\n")
		b.WriteString(styles.HelpStyle.Render("enter: selecionar • esc: voltar"))
		return b.String()
	}

	b.WriteString(styles.WarningStyle.Render(SpinnerChar(spinnerFrame) + "  Gravando arquivos de skill..."))
	b.WriteString("\n\n")
	b.WriteString(styles.SubtextStyle.Render("Instalando SKILL.md em todos os agentes detectados. Isso deve ser rápido."))
	b.WriteString("\n\n")
	b.WriteString(styles.HelpStyle.Render("por favor, aguarde..."))

	return b.String()
}
