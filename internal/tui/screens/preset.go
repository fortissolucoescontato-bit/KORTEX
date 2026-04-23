package screens

import (
	"strings"

	"github.com/fortissolucoescontato-bit/kortex/internal/model"
	"github.com/fortissolucoescontato-bit/kortex/internal/tui/styles"
)

func PresetOptions() []model.PresetID {
	return []model.PresetID{
		model.PresetFullKortex,
		model.PresetEcosystemOnly,
		model.PresetMinimal,
		model.PresetCustom,
	}
}

var presetLabels = map[model.PresetID]string{
	model.PresetFullKortex:    "Ecossistema Completo",
	model.PresetEcosystemOnly: "Apenas Ecossistema",
	model.PresetMinimal:       "Mínimo",
	model.PresetCustom:        "Personalizado",
}

var presetDescriptions = map[model.PresetID]string{
	model.PresetFullKortex:    "Completo: memória, SDD, skills, docs, persona e segurança",
	model.PresetEcosystemOnly: "Apenas ferramentas core: memória, SDD, skills e docs (sem persona/segurança)",
	model.PresetMinimal:       "Apenas memória persistente KortexEngram",
	model.PresetCustom:        "Escolher componentes manualmente; mantém persona e configurações atuais",
}

func RenderPreset(selected model.PresetID, cursor int) string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Selecione o Preset do Ecossistema"))
	b.WriteString("\n\n")

	for idx, preset := range PresetOptions() {
		isSelected := preset == selected
		focused := idx == cursor
		label := presetLabels[preset]
		b.WriteString(renderRadio(label, isSelected, focused))
		b.WriteString(styles.SubtextStyle.Render("    "+presetDescriptions[preset]) + "\n")
	}

	b.WriteString("\n")
	b.WriteString(renderOptions([]string{"Voltar"}, cursor-len(PresetOptions())))
	b.WriteString("\n")
	b.WriteString(styles.HelpStyle.Render("j/k: navegar • enter: selecionar • esc: voltar"))

	return b.String()
}
