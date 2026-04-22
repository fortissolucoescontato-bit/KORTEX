package screens

import (
	"strings"

	"github.com/fortissolucoescontato-bit/kortex/internal/model"
	"github.com/fortissolucoescontato-bit/kortex/internal/planner"
	"github.com/fortissolucoescontato-bit/kortex/internal/tui/styles"
)

func ReviewOptions() []string {
	return []string{"Instalar", "Voltar"}
}

func RenderReview(payload planner.ReviewPayload, cursor int) string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Revisar e Confirmar"))
	b.WriteString("\n\n")

	b.WriteString("  " + styles.HeadingStyle.Render("Agentes") + "  " + styles.UnselectedStyle.Render(joinIDs(payload.Agents)) + "\n")
	b.WriteString("  " + styles.HeadingStyle.Render("Persona") + "  " + styles.UnselectedStyle.Render(reviewPersonaLabel(payload.Persona)) + "\n")
	b.WriteString("  " + styles.HeadingStyle.Render("Preset") + "  " + styles.UnselectedStyle.Render(reviewPresetLabel(payload.Preset)) + "\n")
	b.WriteString("\n")

	if len(payload.Components) > 0 {
		autoSet := make(map[model.ComponentID]struct{}, len(payload.AddedDependencies))
		for _, dep := range payload.AddedDependencies {
			autoSet[dep] = struct{}{}
		}

		b.WriteString(styles.HeadingStyle.Render("Componentes"))
		b.WriteString("\n")
		for _, comp := range payload.Components {
			badge := styles.SubtextStyle.Render("selecionado")
			if _, isAuto := autoSet[comp.ID]; isAuto {
				badge = styles.WarningStyle.Render("dependência automática")
			}
			b.WriteString("  " + styles.UnselectedStyle.Render(string(comp.ID)) + " " + badge + "\n")
		}

		// Issue #145: show individual skill names when the Skills component is selected.
		if len(payload.Skills) > 0 {
			b.WriteString(styles.HeadingStyle.Render("  Skills"))
			b.WriteString("\n")
			for _, skill := range payload.Skills {
				b.WriteString("    " + styles.SubtextStyle.Render(string(skill)) + "\n")
			}
		}

		// Issue #149: show Strict TDD status when SDD is in the plan.
		if payload.HasSDD {
			strictLabel := "Desativado"
			if payload.StrictTDD {
				strictLabel = "Ativado"
			}
			b.WriteString("  " + styles.HeadingStyle.Render("Strict TDD") + "  " + styles.UnselectedStyle.Render(strictLabel) + "\n")
		}

		b.WriteString("\n")
	}

	if len(payload.UnsupportedAgents) > 0 {
		b.WriteString(styles.WarningStyle.Render("Agentes não suportados: " + joinIDs(payload.UnsupportedAgents)))
		b.WriteString("\n\n")
	}

	b.WriteString(renderOptions(ReviewOptions(), cursor))
	b.WriteString("\n")
	b.WriteString(styles.HelpStyle.Render("enter: instalar • esc: voltar"))

	return b.String()
}

func joinIDs[T ~string](values []T) string {
	if len(values) == 0 {
		return "nenhum"
	}

	parts := make([]string, 0, len(values))
	for _, value := range values {
		parts = append(parts, string(value))
	}

	return strings.Join(parts, ", ")
}

func reviewPersonaLabel(persona model.PersonaID) string {
	if persona == model.PersonaCustom {
		return "manter persona atual (não gerenciada)"
	}

	return string(persona)
}

func reviewPresetLabel(preset model.PresetID) string {
	if preset == model.PresetCustom {
		return "escolher componentes e skills manualmente"
	}

	return string(preset)
}
