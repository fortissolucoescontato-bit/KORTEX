package screens

import (
	"fmt"
	"strings"

	"github.com/fortissolucoescontato-bit/kortex/internal/catalog"
	componentuninstall "github.com/fortissolucoescontato-bit/kortex/internal/components/uninstall"
	"github.com/fortissolucoescontato-bit/kortex/internal/model"
	"github.com/fortissolucoescontato-bit/kortex/internal/tui/styles"
)

type UninstallModeOption struct {
	Mode        model.UninstallMode
	Label       string
	Description string
}

type UninstallKortexEngramScopeOption struct {
	Scope       model.KortexEngramUninstallScope
	Label       string
	Description string
}

func UninstallModeOptions() []UninstallModeOption {
	return []UninstallModeOption{
		{
			Mode:        model.UninstallModePartial,
			Label:       "Desinstalação Parcial",
			Description: "Selecionar agentes e componentes específicos para remover",
		},
		{
			Mode:        model.UninstallModeFull,
			Label:       "Desinstalação Completa",
			Description: "Remover toda a configuração gerenciada pelo kortex de todos os agentes",
		},
		{
			Mode:        model.UninstallModeFullRemove,
			Label:       "Desinstalação Completa e Remover Binário",
			Description: "Remover toda a configuração E excluir o próprio binário do kortex",
		},
		{
			Mode:        model.UninstallModeCleanInstall,
			Label:       "Desinstalação Completa + Instalação Limpa",
			Description: "Remover toda a configuração e sincronizar tudo do zero novamente",
		},
	}
}

func RenderUninstallMode(cursor int) string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Seleção do Modo de Desinstalação"))
	b.WriteString("\n\n")
	b.WriteString(styles.SubtextStyle.Render("Escolha como você deseja desinstalar o kortex:"))
	b.WriteString("\n\n")

	options := UninstallModeOptions()
	for idx, opt := range options {
		focused := idx == cursor
		if focused {
			b.WriteString(styles.SelectedStyle.Render("▸ " + opt.Label))
		} else {
			b.WriteString(styles.UnselectedStyle.Render("  " + opt.Label))
		}
		b.WriteString("\n")
		b.WriteString(styles.SubtextStyle.Render("  " + opt.Description))
		b.WriteString("\n")
		if opt.Mode == model.UninstallModeFullRemove {
			b.WriteString(styles.ErrorStyle.Render("  ⚠ AVISO: Isso não pode ser desfeito sem reinstalar"))
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	b.WriteString(renderOptions([]string{"Voltar"}, cursor-len(options)))
	b.WriteString("\n")
	b.WriteString(styles.HelpStyle.Render("j/k: navegar • enter: selecionar • esc: voltar"))

	return b.String()
}

func UninstallAgentOptions() []catalog.Agent {
	return catalog.AllAgents()
}

func UninstallComponentOptions() []catalog.Component {
	return catalog.MVPComponents()
}

func RenderUninstall(selected []model.AgentID, cursor int) string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Desinstalar Configurações Gerenciadas"))
	b.WriteString("\n\n")
	b.WriteString(styles.HelpStyle.Render("Use j/k para mover, espaço para marcar, enter para continuar."))
	b.WriteString("\n\n")
	b.WriteString(styles.SubtextStyle.Render("Selecione os agentes cujas configurações gerenciadas pelo kortex devem ser removidas."))
	b.WriteString("\n\n")

	selectedSet := make(map[model.AgentID]struct{}, len(selected))
	for _, agent := range selected {
		selectedSet[agent] = struct{}{}
	}

	for idx, agent := range UninstallAgentOptions() {
		_, checked := selectedSet[agent.ID]
		focused := idx == cursor
		b.WriteString(renderCheckbox(agent.Name, checked, focused))
	}

	b.WriteString("\n")
	agentCount := len(UninstallAgentOptions())
	relCursor := cursor - agentCount
	if len(selected) == 0 {
		// Render Continue as dimmed with an inline hint when nothing is selected.
		if relCursor == 0 {
			b.WriteString(styles.SelectedStyle.Render(styles.Cursor+styles.SubtextStyle.Render("Continuar")+" "+styles.HelpStyle.Render("(selecione ao menos um agente)")) + "\n")
		} else {
			b.WriteString(styles.SubtextStyle.Render("  Continuar") + "\n")
		}
		b.WriteString(renderOptions([]string{"Voltar"}, relCursor-1))
	} else {
		b.WriteString(renderOptions([]string{"Continuar", "Voltar"}, relCursor))
	}
	b.WriteString("\n")
	b.WriteString(styles.HelpStyle.Render("espaço: marcar • enter: confirmar • esc: voltar"))

	return b.String()
}

func RenderUninstallComponents(selected []model.ComponentID, cursor int) string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Desinstalar Componentes Gerenciados"))
	b.WriteString("\n\n")
	b.WriteString(styles.HelpStyle.Render("Use j/k para mover, espaço para marcar, enter para continuar."))
	b.WriteString("\n\n")
	b.WriteString(styles.SubtextStyle.Render("Selecione quais componentes gerenciados devem ser removidos dos agentes selecionados."))
	b.WriteString("\n\n")

	selectedSet := make(map[model.ComponentID]struct{}, len(selected))
	for _, component := range selected {
		selectedSet[component] = struct{}{}
	}

	for idx, component := range UninstallComponentOptions() {
		_, checked := selectedSet[component.ID]
		focused := idx == cursor
		b.WriteString(renderCheckbox(component.Name, checked, focused))
		b.WriteString(styles.SubtextStyle.Render("    "+component.Description) + "\n")
	}

	b.WriteString("\n")
	compCount := len(UninstallComponentOptions())
	relCursor := cursor - compCount
	if len(selected) == 0 {
		// Render Continue as dimmed with an inline hint when nothing is selected.
		if relCursor == 0 {
			b.WriteString(styles.SelectedStyle.Render(styles.Cursor+styles.SubtextStyle.Render("Continuar")+" "+styles.HelpStyle.Render("(selecione ao menos um componente)")) + "\n")
		} else {
			b.WriteString(styles.SubtextStyle.Render("  Continuar") + "\n")
		}
		b.WriteString(renderOptions([]string{"Voltar"}, relCursor-1))
	} else {
		b.WriteString(renderOptions([]string{"Continuar", "Voltar"}, relCursor))
	}
	b.WriteString("\n")
	b.WriteString(styles.HelpStyle.Render("espaço: marcar • enter: continuar • esc: voltar"))

	return b.String()
}

func uninstallKortexEngramScopeOptions(projectScopeAvailable bool) []UninstallKortexEngramScopeOption {
	options := make([]UninstallKortexEngramScopeOption, 0, 2)
	if projectScopeAvailable {
		options = append(options, UninstallKortexEngramScopeOption{
			Scope:       model.KortexEngramUninstallScopeProject,
			Label:       "Limpeza apenas do projeto",
			Description: "Excluir apenas .KortexEngram/ no projeto atual",
		})
	}
	options = append(options, UninstallKortexEngramScopeOption{
		Scope:       model.KortexEngramUninstallScopeGlobal,
		Label:       "Limpeza global",
		Description: "Remover integração global KortexEngram MCP/prompt de sistema",
	})
	return options
}

func RenderUninstallProfiles(available []string, selected []string, KortexEngramProjectScopeAvailable bool, selectedKortexEngramScope model.KortexEngramUninstallScope, cursor int) string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Seleção do Escopo de Desinstalação"))
	b.WriteString("\n\n")
	b.WriteString(styles.HelpStyle.Render("Use j/k para mover, espaço para marcar/selecionar, enter para continuar."))
	b.WriteString("\n\n")

	if len(available) > 0 {
		b.WriteString(styles.SubtextStyle.Render("Escolha quais perfis de SDD OpenCode devem ser removidos do opencode.json."))
		b.WriteString("\n\n")
	}

	selectedSet := make(map[string]struct{}, len(selected))
	for _, profile := range selected {
		selectedSet[profile] = struct{}{}
	}

	for idx, profileName := range available {
		_, checked := selectedSet[profileName]
		focused := idx == cursor
		b.WriteString(renderCheckbox(profileName, checked, focused))
	}

	KortexEngramScopeOptions := uninstallKortexEngramScopeOptions(KortexEngramProjectScopeAvailable)
	KortexEngramScopeDisplayed := 0
	if len(KortexEngramScopeOptions) > 1 {
		KortexEngramScopeDisplayed = len(KortexEngramScopeOptions)
		if len(available) > 0 {
			b.WriteString("\n")
		}
		b.WriteString(styles.SubtextStyle.Render("Selecione o escopo de limpeza do KortexEngram:"))
		b.WriteString("\n")
		for idx, option := range KortexEngramScopeOptions {
			focused := len(available)+idx == cursor
			checked := selectedKortexEngramScope == option.Scope
			b.WriteString(renderCheckbox(option.Label, checked, focused))
			b.WriteString(styles.SubtextStyle.Render("    " + option.Description))
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	relCursor := cursor - (len(available) + KortexEngramScopeDisplayed)
	b.WriteString(renderOptions([]string{"Continuar", "Voltar"}, relCursor))
	b.WriteString("\n")
	b.WriteString(styles.HelpStyle.Render("espaço: marcar/selecionar • enter: continuar • esc: voltar"))

	return b.String()
}

func RenderUninstallConfirm(mode model.UninstallMode, selected []model.AgentID, components []model.ComponentID, profilesToRemove []string, KortexEngramScope model.KortexEngramUninstallScope, KortexEngramProjectScopeAvailable bool, cursor int, operationRunning bool, spinnerFrame int) string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Confirmar Desinstalação"))
	b.WriteString("\n\n")

	if operationRunning {
		b.WriteString(styles.WarningStyle.Render(SpinnerChar(spinnerFrame) + "  Removendo configurações gerenciadas..."))
		b.WriteString("\n\n")
		b.WriteString(styles.HelpStyle.Render("Por favor, aguarde..."))
		return b.String()
	}

	// Render mode-specific information
	switch mode {
	case model.UninstallModePartial:
		if len(selected) == 0 {
			b.WriteString(styles.WarningStyle.Render("Nenhum agente selecionado."))
			b.WriteString("\n\n")
			b.WriteString(styles.HelpStyle.Render("enter: voltar • esc: voltar"))
			return b.String()
		}
		b.WriteString(styles.SubtextStyle.Render("Modo: Desinstalação Parcial"))
		b.WriteString("\n\n")
		b.WriteString(styles.SubtextStyle.Render("Agentes:"))
		b.WriteString("\n")
		for _, label := range uninstallAgentLabels(selected) {
			b.WriteString(styles.UnselectedStyle.Render("  • " + label))
			b.WriteString("\n")
		}
		b.WriteString("\n")
		b.WriteString(styles.SubtextStyle.Render("Componentes:"))
		b.WriteString("\n")
		for _, label := range uninstallComponentLabels(components) {
			b.WriteString(styles.UnselectedStyle.Render("  • " + label))
			b.WriteString("\n")
		}
	case model.UninstallModeFull:
		b.WriteString(styles.SubtextStyle.Render("Modo: Desinstalação Completa"))
		b.WriteString("\n\n")
		b.WriteString(styles.UnselectedStyle.Render("Isso irá remover toda a configuração gerenciada pelo kortex de todos os agentes suportados."))
		b.WriteString("\n")
	case model.UninstallModeFullRemove:
		b.WriteString(styles.ErrorStyle.Render("Modo: Desinstalação Completa e Remover Binário"))
		b.WriteString("\n\n")
		b.WriteString(styles.UnselectedStyle.Render("Isso irá remover toda a configuração gerenciada de todos os agentes"))
		b.WriteString("\n")
		b.WriteString(styles.ErrorStyle.Render("E excluir o próprio binário do kortex."))
		b.WriteString("\n\n")
		b.WriteString(styles.ErrorStyle.Render("  ⚠ AVISO: Esta ação não pode ser desfeita sem reinstalar!"))
		b.WriteString("\n")
	case model.UninstallModeCleanInstall:
		b.WriteString(styles.SuccessStyle.Render("Modo: Desinstalação Completa + Instalação Limpa"))
		b.WriteString("\n\n")
		b.WriteString(styles.UnselectedStyle.Render("Isso irá remover toda a configuração de todos os agentes"))
		b.WriteString("\n")
		b.WriteString(styles.SuccessStyle.Render("e imediatamente sincronizar todos os ativos gerenciados do zero."))
		b.WriteString("\n\n")
		b.WriteString(styles.SubtextStyle.Render("Use isso para corrigir configurações corrompidas ou resetar o estado."))
		b.WriteString("\n")
	}

	if len(profilesToRemove) > 0 {
		b.WriteString("\n")
		b.WriteString(styles.SubtextStyle.Render("Perfis a serem removidos:"))
		b.WriteString("\n")
		for _, profile := range profilesToRemove {
			b.WriteString(styles.UnselectedStyle.Render("  • " + profile))
			b.WriteString("\n")
		}
	}

	if hasSelectedComponent(components, model.ComponentKortexEngram) {
		b.WriteString("\n")
		b.WriteString(styles.SubtextStyle.Render("Escopo de limpeza do KortexEngram:"))
		b.WriteString("\n")
		scopeLabel := "Global"
		detail := "  • Remove a integração global KortexEngram MCP/prompt de sistema"
		if KortexEngramScope == model.KortexEngramUninstallScopeProject && KortexEngramProjectScopeAvailable {
			scopeLabel = "Apenas Projeto"
			detail = "  • Exclui .KortexEngram/ apenas no projeto atual"
		}
		b.WriteString(styles.UnselectedStyle.Render("  • " + scopeLabel))
		b.WriteString("\n")
		b.WriteString(styles.SubtextStyle.Render(detail))
		b.WriteString("\n")
	}

	b.WriteString("\n")

	// Workspace-scoped assets warning
	hasWorkspaceAssets := false
	for _, comp := range components {
		if comp == model.ComponentSDD || comp == model.ComponentSkills {
			hasWorkspaceAssets = true
			break
		}
	}
	if (mode == model.UninstallModeFull || mode == model.UninstallModeFullRemove) || hasWorkspaceAssets {
		b.WriteString(styles.WarningStyle.Render("  ⚠ Aviso de Ativos do Workspace:"))
		b.WriteString("\n")
		b.WriteString(styles.SubtextStyle.Render("  Remover SDD ou Skills excluirá arquivos do projeto como:"))
		b.WriteString("\n")
		b.WriteString(styles.SubtextStyle.Render("  • .windsurf/workflows/ (fluxos SDD)"))
		b.WriteString("\n")
		b.WriteString(styles.SubtextStyle.Render("  • .KortexEngram/ (contexto de memória persistente)"))
		b.WriteString("\n")
		b.WriteString(styles.SubtextStyle.Render("  • Diretórios de Skills"))
		b.WriteString("\n\n")
		b.WriteString(styles.ErrorStyle.Render("  Se você realizar o commit dessas deleções, TODOS os colaboradores perderão este contexto!"))
		b.WriteString("\n\n")
	}

	b.WriteString(styles.WarningStyle.Render("Um backup será criado antes que qualquer arquivo seja modificado."))
	b.WriteString("\n\n")
	b.WriteString(renderOptions([]string{"Desinstalar", "Cancelar"}, cursor))
	b.WriteString("\n")
	b.WriteString(styles.HelpStyle.Render("j/k: navegar • enter: selecionar • esc: voltar"))

	return b.String()
}

func RenderUninstallResult(result componentuninstall.Result, err error, mode model.UninstallMode, selectedProfiles []string, KortexEngramScope model.KortexEngramUninstallScope, KortexEngramProjectScopeAvailable bool, syncFilesChanged int, syncErr error) string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Resultado da Desinstalação"))
	b.WriteString("\n\n")

	if err != nil {
		b.WriteString(styles.ErrorStyle.Render("✗ Falha na desinstalação"))
		b.WriteString("\n\n")
		b.WriteString(styles.HeadingStyle.Render("Erro:"))
		b.WriteString("\n")
		b.WriteString(styles.ErrorStyle.Render("  " + err.Error()))
		b.WriteString("\n\n")
		if result.Manifest.ID != "" {
			b.WriteString(styles.SubtextStyle.Render("Backup criado antes da falha: "))
			b.WriteString(styles.SelectedStyle.Render(result.Manifest.ID))
			b.WriteString("\n")
			b.WriteString(styles.SubtextStyle.Render(result.Manifest.DisplayLabel()))
		}
	} else {
		b.WriteString(styles.SuccessStyle.Render("✓ Desinstalação concluída"))
		b.WriteString("\n\n")
		if result.Manifest.ID != "" {
			b.WriteString(styles.SubtextStyle.Render("Backup: "))
			b.WriteString(styles.SelectedStyle.Render(result.Manifest.ID))
			b.WriteString("\n")
			b.WriteString(styles.SubtextStyle.Render(result.Manifest.DisplayLabel()))
			b.WriteString("\n\n")
		}
		b.WriteString(styles.UnselectedStyle.Render(fmt.Sprintf("Arquivos reescritos: %d", len(result.ChangedFiles))))
		b.WriteString("\n")
		b.WriteString(styles.UnselectedStyle.Render(fmt.Sprintf("Arquivos excluídos: %d", len(result.RemovedFiles))))
		b.WriteString("\n")
		b.WriteString(styles.UnselectedStyle.Render(fmt.Sprintf("Diretórios excluídos: %d", len(result.RemovedDirectories))))
		if len(result.AgentsRemovedFromState) > 0 {
			b.WriteString("\n")
			b.WriteString(styles.UnselectedStyle.Render("Estado atualizado: " + strings.Join(uninstallAgentLabels(result.AgentsRemovedFromState), ", ")))
		}
		if len(result.ManualActions) > 0 {
			b.WriteString("\n\n")
			b.WriteString(styles.WarningStyle.Render("Limpeza manual necessária:"))
			for _, item := range result.ManualActions {
				b.WriteString("\n")
				b.WriteString(styles.UnselectedStyle.Render("  • " + item))
			}
		}

		if len(selectedProfiles) > 0 {
			b.WriteString("\n\n")
			b.WriteString(styles.UnselectedStyle.Render("Perfis removidos: " + strings.Join(selectedProfiles, ", ")))
		}

		if hasKortexEngramArtifacts(result) {
			b.WriteString("\n\n")
			if KortexEngramScope == model.KortexEngramUninstallScopeProject && KortexEngramProjectScopeAvailable {
				b.WriteString(styles.UnselectedStyle.Render("Escopo KortexEngram: Apenas Projeto (.KortexEngram/ removido do workspace atual)"))
			} else {
				b.WriteString(styles.UnselectedStyle.Render("Escopo KortexEngram: Global (Integração MCP/prompt de sistema removida)"))
			}
		}

		// Clean install: show sync results after uninstall stats.
		if mode == model.UninstallModeCleanInstall {
			b.WriteString("\n\n")
			if syncErr != nil {
				b.WriteString(styles.ErrorStyle.Render("✗ Falha na sincronização da instalação limpa"))
				b.WriteString("\n")
				b.WriteString(styles.ErrorStyle.Render("  " + syncErr.Error()))
				b.WriteString("\n\n")
				b.WriteString(styles.WarningStyle.Render("Você pode executar 'kortex sync' manualmente para tentar novamente."))
			} else {
				b.WriteString(styles.SuccessStyle.Render("✓ Sincronização da instalação limpa concluída"))
				b.WriteString("\n")
				b.WriteString(styles.UnselectedStyle.Render(fmt.Sprintf("Arquivos sincronizados: %d", syncFilesChanged)))
			}
		}
	}

	b.WriteString("\n\n")
	b.WriteString(styles.HelpStyle.Render("enter: retornar • esc: voltar • q: sair"))
	return b.String()
}

func hasKortexEngramArtifacts(result componentuninstall.Result) bool {
	for _, path := range result.ChangedFiles {
		if strings.Contains(path, "kortex-engram") {
			return true
		}
	}
	for _, path := range result.RemovedFiles {
		if strings.Contains(path, "kortex-engram") {
			return true
		}
	}
	for _, path := range result.RemovedDirectories {
		if strings.Contains(path, ".KortexEngram") {
			return true
		}
	}
	return false
}

func uninstallAgentLabels(agentIDs []model.AgentID) []string {
	labels := make([]string, 0, len(agentIDs))
	for _, selected := range agentIDs {
		labels = append(labels, uninstallAgentLabel(selected))
	}
	return labels
}

func uninstallAgentLabel(agentID model.AgentID) string {
	for _, agent := range UninstallAgentOptions() {
		if agent.ID == agentID {
			return agent.Name
		}
	}
	return string(agentID)
}

func uninstallComponentLabels(componentIDs []model.ComponentID) []string {
	labels := make([]string, 0, len(componentIDs))
	for _, selected := range componentIDs {
		labels = append(labels, uninstallComponentLabel(selected))
	}
	return labels
}

func uninstallComponentLabel(componentID model.ComponentID) string {
	for _, component := range UninstallComponentOptions() {
		if component.ID == componentID {
			return component.Name
		}
	}
	return string(componentID)
}

func hasSelectedComponent(components []model.ComponentID, target model.ComponentID) bool {
	for _, component := range components {
		if component == target {
			return true
		}
	}
	return false
}
