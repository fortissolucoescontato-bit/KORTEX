package screens

import (
	"fmt"
	"strings"

	"github.com/fortissolucoescontato-bit/kortex/internal/tui/styles"
	"github.com/fortissolucoescontato-bit/kortex/internal/update"
)

// WelcomeOptions returns the welcome menu options.
// When showProfiles is true, an "OpenCode SDD Profiles" option is inserted
// between "Configure models" and "Manage backups".
// profileCount is used to show a badge with the current profile count.
// When hasEngines is false, "Create your own Agent" is shown as disabled
// (labelled "(no agents)") to signal that no supported AI engine is installed.
func WelcomeOptions(updateResults []update.UpdateResult, updateCheckDone bool, showProfiles bool, profileCount int, hasEngines bool) []string {
	upgradeLabel := "Atualizar ferramentas"
	if updateCheckDone && update.HasUpdates(updateResults) {
		upgradeLabel = "Atualizar ferramentas ★"
	} else if updateCheckDone && !update.HasUpdates(updateResults) {
		upgradeLabel = "Atualizar ferramentas (atualizado)"
	}

	agentLabel := "Criar seu próprio Agente"
	if !hasEngines {
		agentLabel = "Criar seu próprio Agente (sem motores)"
	}

	opts := []string{
		"Iniciar instalação",
		upgradeLabel,
		"Sincronizar configurações",
		"Atualizar + Sincronizar",
		"Configurar modelos",
		agentLabel,
	}

	if showProfiles {
		profilesLabel := "Perfis de SDD OpenCode"
		if profileCount > 0 {
			profilesLabel = fmt.Sprintf("Perfis de SDD OpenCode (%d)", profileCount)
		}
		opts = append(opts, profilesLabel)
	}

	opts = append(opts, "Gerenciar backups")
	opts = append(opts, "Desinstalação gerenciada")
	opts = append(opts, "Sair")

	return opts
}

func RenderWelcome(cursor int, updateBanner string, updateResults []update.UpdateResult, updateCheckDone bool, showProfiles bool, profileCount int, hasEngines bool) string {
	var b strings.Builder

	b.WriteString(styles.RenderLogo())
	b.WriteString("\n\n")
	b.WriteString(styles.SubtextStyle.Render(styles.Tagline()))
	b.WriteString("\n")

	if updateBanner != "" {
		b.WriteString(styles.WarningStyle.Render(updateBanner))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(styles.HeadingStyle.Render("Menu"))
	b.WriteString("\n\n")
	b.WriteString(renderOptions(WelcomeOptions(updateResults, updateCheckDone, showProfiles, profileCount, hasEngines), cursor))
	b.WriteString("\n")
	b.WriteString(styles.HelpStyle.Render("j/k: navegar • enter: selecionar • q: sair"))

	return styles.FrameStyle.Render(b.String())
}
