package screens

import (
	"strings"

	"github.com/fortissolucoescontato-bit/kortex/internal/catalog"
	"github.com/fortissolucoescontato-bit/kortex/internal/model"
	"github.com/fortissolucoescontato-bit/kortex/internal/tui/styles"
)

func AgentOptions() []model.AgentID {
	agents := catalog.AllAgents()
	ids := make([]model.AgentID, 0, len(agents))
	for _, agent := range agents {
		ids = append(ids, agent.ID)
	}
	return ids
}

func RenderAgents(selected []model.AgentID, cursor int) string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Selecione os Agentes de IA"))
	b.WriteString("\n\n")
	b.WriteString(styles.HelpStyle.Render("Use j/k para mover, espaço para marcar, enter para continuar."))
	b.WriteString("\n\n")

	selectedSet := make(map[model.AgentID]struct{}, len(selected))
	for _, agent := range selected {
		selectedSet[agent] = struct{}{}
	}

	agents := AgentOptions()
	for idx, agent := range agents {
		_, checked := selectedSet[agent]
		focused := idx == cursor
		b.WriteString(renderCheckbox(string(agent), checked, focused))
	}

	b.WriteString("\n")
	actions := []string{"Continuar", "Voltar"}
	b.WriteString(renderOptions(actions, cursor-len(agents)))
	b.WriteString("\n")
	b.WriteString(styles.HelpStyle.Render("espaço: marcar • enter: confirmar • esc: voltar"))

	return b.String()
}
