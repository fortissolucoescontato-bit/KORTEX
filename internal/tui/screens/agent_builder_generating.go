package screens

import (
	"strings"

	"github.com/fortissolucoescontato-bit/kortex/internal/tui/styles"
)

// RenderABGenerating renders the generation-in-progress (or error) screen.
func RenderABGenerating(engineName string, spinnerFrame int, genErr error) string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Gerando seu Agente..."))
	b.WriteString("\n\n")

	if genErr != nil {
		b.WriteString(styles.ErrorStyle.Render("✗ Falha na geração"))
		b.WriteString("\n")
		b.WriteString(styles.SubtextStyle.Render("  Motor: " + engineName))
		b.WriteString("\n")
		b.WriteString(styles.ErrorStyle.Render("  Erro: " + genErr.Error()))
		b.WriteString("\n\n")
		b.WriteString(renderOptions([]string{"Tentar novamente", "Voltar"}, 0))
		b.WriteString("\n")
		b.WriteString(styles.HelpStyle.Render("enter: selecionar • esc: voltar"))
		return b.String()
	}

	b.WriteString(styles.WarningStyle.Render(SpinnerChar(spinnerFrame) + "  Executando " + engineName + "..."))
	b.WriteString("\n\n")
	b.WriteString(styles.SubtextStyle.Render("Compondo prompt e chamando o motor de geração. Isso pode levar um momento."))
	b.WriteString("\n\n")
	b.WriteString(styles.HelpStyle.Render("esc: cancelar"))

	return b.String()
}
