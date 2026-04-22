package screens

import (
	"fmt"
	"strings"

	"github.com/fortissolucoescontato-bit/kortex/internal/tui/styles"
)

const maxErrorLines = 15

type FailedStep struct {
	ID    string
	Error string
}

type MissingDep struct {
	Name        string
	InstallHint string
}

// UpdateInfo holds version update information for a single tool.
type UpdateInfo struct {
	Name             string
	InstalledVersion string
	LatestVersion    string
	UpdateHint       string
}

type CompletePayload struct {
	ConfiguredAgents    int
	InstalledComponents int
	KortexCLIInstalled        bool
	FailedSteps         []FailedStep
	RollbackPerformed   bool
	MissingDeps         []MissingDep
	AvailableUpdates    []UpdateInfo
}

func RenderComplete(data CompletePayload) string {
	if len(data.FailedSteps) > 0 {
		return renderCompleteFailed(data)
	}
	return renderCompleteSuccess(data)
}

func renderCompleteSuccess(data CompletePayload) string {
	var b strings.Builder

	b.WriteString(styles.SuccessStyle.Render("Pronto! Seus agentes de IA estão configurados."))
	b.WriteString("\n\n")

	b.WriteString("  " + styles.HeadingStyle.Render("Agentes configurados") + "  " + styles.SuccessStyle.Render(fmt.Sprintf("%d", data.ConfiguredAgents)) + "\n")
	b.WriteString("  " + styles.HeadingStyle.Render("Componentes instalados") + "  " + styles.SuccessStyle.Render(fmt.Sprintf("%d", data.InstalledComponents)) + "\n")
	b.WriteString("\n")

	renderMissingDeps(&b, data.MissingDeps)
	renderAvailableUpdates(&b, data.AvailableUpdates)

	b.WriteString(styles.HeadingStyle.Render("Próximos passos"))
	b.WriteString("\n")
	b.WriteString(styles.UnselectedStyle.Render("  1. Configure suas chaves de API"))
	b.WriteString("\n")
	b.WriteString(styles.UnselectedStyle.Render("  2. Inicie o agente de sua escolha"))
	b.WriteString("\n")
	b.WriteString(styles.UnselectedStyle.Render("  3. Experimente /sdd-new minha-funcionalidade"))
	b.WriteString("\n\n")

	if data.KortexCLIInstalled {
		b.WriteString(styles.HeadingStyle.Render("KortexCLI (por projeto)"))
		b.WriteString("\n")
		b.WriteString(styles.UnselectedStyle.Render("  O KortexCLI foi instalado globalmente."))
		b.WriteString("\n")
		b.WriteString(styles.UnselectedStyle.Render("  Em cada repositório execute: kortex init"))
		b.WriteString("\n")
		b.WriteString(styles.UnselectedStyle.Render("  Depois execute: kortex install"))
		b.WriteString("\n\n")
	}

	b.WriteString(styles.HelpStyle.Render("Pressione Enter para sair."))

	return b.String()
}

func renderMissingDeps(b *strings.Builder, deps []MissingDep) {
	if len(deps) == 0 {
		return
	}

	b.WriteString(styles.WarningStyle.Render(fmt.Sprintf("Faltam %d dependência(s):", len(deps))))
	b.WriteString("\n")
	for _, dep := range deps {
		b.WriteString("  " + styles.WarningStyle.Render(dep.Name) + "  " + styles.SubtextStyle.Render(dep.InstallHint))
		b.WriteString("\n")
	}
	b.WriteString("\n")
}

func renderAvailableUpdates(b *strings.Builder, updates []UpdateInfo) {
	if len(updates) == 0 {
		return
	}

	b.WriteString(styles.HeadingStyle.Render("Atualizações Disponíveis"))
	b.WriteString("\n")
	for _, u := range updates {
		line := fmt.Sprintf("  %s %s -> %s", u.Name, u.InstalledVersion, u.LatestVersion)
		b.WriteString(styles.WarningStyle.Render(line))
		if u.UpdateHint != "" {
			b.WriteString("  " + styles.SubtextStyle.Render(u.UpdateHint))
		}
		b.WriteString("\n")
	}
	b.WriteString("\n")
}

func renderCompleteFailed(data CompletePayload) string {
	var b strings.Builder

	b.WriteString(styles.ErrorStyle.Render("Instalação concluída com erros."))
	b.WriteString("\n\n")

	b.WriteString(styles.HeadingStyle.Render("Etapas que falharam"))
	b.WriteString("\n")
	for _, step := range data.FailedSteps {
		b.WriteString("  " + styles.ErrorStyle.Render("✗ "+step.ID))
		b.WriteString("\n")
		lines := strings.Split(step.Error, "\n")
		if len(lines) > maxErrorLines {
			lines = lines[:maxErrorLines]
			lines = append(lines, "... (truncado)")
		}
		for _, line := range lines {
			b.WriteString("    " + styles.SubtextStyle.Render(line))
			b.WriteString("\n")
		}
	}
	b.WriteString("\n")

	if data.RollbackPerformed {
		b.WriteString(styles.WarningStyle.Render("Rollback realizado — configuração anterior restaurada."))
		b.WriteString("\n\n")
	}

	renderMissingDeps(&b, data.MissingDeps)
	renderAvailableUpdates(&b, data.AvailableUpdates)

	b.WriteString(styles.HeadingStyle.Render("O que fazer"))
	b.WriteString("\n")
	b.WriteString(styles.UnselectedStyle.Render("  1. Verifique as mensagens de erro acima"))
	b.WriteString("\n")
	b.WriteString(styles.UnselectedStyle.Render("  2. Corrija o problema (dependências, permissões, etc.)"))
	b.WriteString("\n")
	b.WriteString(styles.UnselectedStyle.Render("  3. Execute o kortex novamente para tentar de novo"))
	b.WriteString("\n\n")

	b.WriteString(styles.HelpStyle.Render("Pressione Enter para sair."))

	return b.String()
}
