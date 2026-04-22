package screens

import (
	"strings"

	"github.com/fortissolucoescontato-bit/kortex/internal/tui/styles"
)

// StrictTDDOptionEnable is the index of the "Enable" option.
const StrictTDDOptionEnable = 0

// StrictTDDOptionDisable is the index of the "Disable" option.
const StrictTDDOptionDisable = 1

// StrictTDDOptions returns the list of option labels for the Strict TDD screen.
func StrictTDDOptions() []string {
	return []string{"Ativar", "Desativar"}
}

// RenderStrictTDD renders the Strict TDD Mode selection screen.
// enabled indicates whether Strict TDD Mode is currently active.
// cursor is the current cursor position.
func RenderStrictTDD(enabled bool, cursor int) string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("MODO TDD ESTRITO"))
	b.WriteString("\n\n")
	b.WriteString(styles.SubtextStyle.Render("Os agentes devem seguir o TDD Estrito (RED → GREEN → REFACTOR) para cada tarefa?"))
	b.WriteString("\n")
	b.WriteString(styles.SubtextStyle.Render("Quando ativado, o agente sdd-apply escreve os testes primeiro, confirma a falha,"))
	b.WriteString("\n")
	b.WriteString(styles.SubtextStyle.Render("e então implementa o código mínimo para passar antes da refatoração."))
	b.WriteString("\n\n")

	options := StrictTDDOptions()
	for idx, opt := range options {
		isSelected := (idx == StrictTDDOptionEnable && enabled) || (idx == StrictTDDOptionDisable && !enabled)
		focused := idx == cursor
		b.WriteString(renderRadio(opt, isSelected, focused))
	}

	b.WriteString("\n")
	b.WriteString(renderOptions([]string{"Voltar"}, cursor-len(options)))
	b.WriteString("\n")
	b.WriteString(styles.HelpStyle.Render("j/k: navegar • enter: selecionar • esc: voltar"))

	return b.String()
}
