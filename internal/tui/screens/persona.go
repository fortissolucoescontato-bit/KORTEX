package screens

import (
	"strings"

	"github.com/fortissolucoescontato-bit/kortex/internal/model"
	"github.com/fortissolucoescontato-bit/kortex/internal/tui/styles"
)

func PersonaOptions() []model.PersonaID {
	return []model.PersonaID{model.PersonaKortex, model.PersonaNeutral, model.PersonaCustom}
}

var personaLabels = map[model.PersonaID]string{
	model.PersonaKortex:  "Elite (Analista Nexo-Fortis)",
	model.PersonaNeutral: "Neutro (Execução Técnica)",
	model.PersonaCustom:  "Customizado",
}

var personaDescriptions = map[model.PersonaID]string{
	model.PersonaKortex:  "Analista de Elite Nexo-Fortis: Mentor estratégico 360°. Panorama completo com foco em aprendizado e proatividade.",
	model.PersonaNeutral: "Persona Neutra: Tom profissional e polido para execução técnica direta.",
	model.PersonaCustom:  "Persona Customizada: Mantém sua persona atual; o Kortex não injeta personalidade.",
}

func RenderPersona(selected model.PersonaID, cursor int) string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Escolha sua Persona"))
	b.WriteString("\n\n")
	b.WriteString(styles.SubtextStyle.Render("O seu próprio Kortex! Análise profunda antes da execução."))
	b.WriteString("\n\n")

	for idx, persona := range PersonaOptions() {
		isSelected := persona == selected
		focused := idx == cursor
		label := personaLabels[persona]
		b.WriteString(renderRadio(label, isSelected, focused))
		b.WriteString(styles.SubtextStyle.Render("    " + personaDescriptions[persona]))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(renderOptions([]string{"Voltar"}, cursor-len(PersonaOptions())))
	b.WriteString("\n")
	b.WriteString(styles.HelpStyle.Render("j/k: navegar • enter: selecionar • esc: voltar"))

	return b.String()
}
