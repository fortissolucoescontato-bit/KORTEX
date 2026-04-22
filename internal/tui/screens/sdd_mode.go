package screens

import (
	"strings"

	"github.com/fortissolucoescontato-bit/kortex/internal/model"
	"github.com/fortissolucoescontato-bit/kortex/internal/tui/styles"
)

func SDDModeOptions() []model.SDDModeID {
	return []model.SDDModeID{model.SDDModeSingle, model.SDDModeMulti}
}

var sddModeDescriptions = map[model.SDDModeID]string{
	model.SDDModeSingle: "Orquestrador Único — um agente gerencia todas as fases do SDD",
	model.SDDModeMulti:  "Multi-Agente — sub-agentes dedicados para cada fase (9 agentes ocultos)",
}

var sddModeLabels = map[model.SDDModeID]string{
	model.SDDModeSingle: "Simples",
	model.SDDModeMulti:  "Multi-Agente",
}

func RenderSDDMode(selected model.SDDModeID, cursor int) string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Selecione o Modo SDD"))
	b.WriteString("\n\n")
	b.WriteString(styles.SubtextStyle.Render("Como o orquestrador SDD deve ser configurado para o OpenCode?"))
	b.WriteString("\n\n")

	for idx, mode := range SDDModeOptions() {
		isSelected := mode == selected
		focused := idx == cursor
		label := sddModeLabels[mode]
		b.WriteString(renderRadio(label, isSelected, focused))
		b.WriteString(styles.SubtextStyle.Render("    "+sddModeDescriptions[mode]) + "\n")
	}

	b.WriteString("\n")
	b.WriteString(renderOptions([]string{"Voltar"}, cursor-len(SDDModeOptions())))
	b.WriteString("\n")
	b.WriteString(styles.HelpStyle.Render("j/k: navegar • enter: selecionar • esc: voltar"))

	return b.String()
}
