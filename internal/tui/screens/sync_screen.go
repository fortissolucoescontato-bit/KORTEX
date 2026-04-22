package screens

// Note: this file is intentionally named sync_screen.go instead of sync.go
// because sync.go would conflict with the Go standard library "sync" package name.

import (
	"fmt"
	"strings"

	"github.com/fortissolucoescontato-bit/kortex/internal/tui/styles"
)

// RenderSync handles all states of the sync screen.
//
// State logic:
//  1. operationRunning → "Syncing configurations..." with spinner
//  2. hasSyncRun && (filesChanged > 0 || syncErr != nil) → show result
//  3. Otherwise → show confirmation screen
func RenderSync(filesChanged int, syncErr error, operationRunning bool, hasSyncRun bool, spinnerFrame int) string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Sincronizar Configurações"))
	b.WriteString("\n\n")

	// State 1: sync is running
	if operationRunning {
		b.WriteString(styles.WarningStyle.Render(SpinnerChar(spinnerFrame) + "  Sincronizando configurações..."))
		b.WriteString("\n\n")
		b.WriteString(styles.HelpStyle.Render("Por favor, aguarde..."))
		return b.String()
	}

	// State 2: sync has run — show result
	if hasSyncRun {
		b.WriteString(renderSyncResult(filesChanged, syncErr))
		return b.String()
	}

	// State 3: confirmation screen
	b.WriteString(renderSyncConfirm())
	return b.String()
}

func renderSyncConfirm() string {
	var b strings.Builder

	b.WriteString(styles.UnselectedStyle.Render("A sincronização irá reaplicar suas configurações"))
	b.WriteString("\n")
	b.WriteString(styles.UnselectedStyle.Render("a todos os agentes de IA detectados nesta máquina."))
	b.WriteString("\n\n")

	b.WriteString(styles.SubtextStyle.Render("Esta operação:"))
	b.WriteString("\n")
	b.WriteString(styles.SubtextStyle.Render("  • Lê suas seleções atuais de agentes"))
	b.WriteString("\n")
	b.WriteString(styles.SubtextStyle.Render("  • Reescreve os arquivos de config a partir dos templates"))
	b.WriteString("\n")
	b.WriteString(styles.SubtextStyle.Render("  • Não modifica seus arquivos globais (dotfiles)"))
	b.WriteString("\n\n")

	b.WriteString(styles.HeadingStyle.Render("Pressione Enter para sincronizar"))
	b.WriteString("\n\n")
	b.WriteString(styles.HelpStyle.Render("enter: confirmar • esc: voltar • q: sair"))

	return b.String()
}

func renderSyncResult(filesChanged int, syncErr error) string {
	var b strings.Builder

	if syncErr != nil {
		b.WriteString(styles.ErrorStyle.Render("✗ Falha na sincronização"))
		b.WriteString("\n\n")
		b.WriteString(styles.SubtextStyle.Render(syncErr.Error()))
		b.WriteString("\n\n")
		b.WriteString(styles.HelpStyle.Render("Verifique sua configuração e tente novamente."))
	} else if filesChanged == 0 {
		b.WriteString(styles.SuccessStyle.Render("✓ Sincronização concluída"))
		b.WriteString("\n\n")
		b.WriteString(styles.SubtextStyle.Render("Nenhum agente detectado ou nenhum arquivo precisava de atualização."))
	} else {
		b.WriteString(styles.SuccessStyle.Render("✓ Sincronização concluída"))
		b.WriteString("\n\n")
		b.WriteString(fmt.Sprintf("%s %s", styles.HeadingStyle.Render(fmt.Sprintf("%d arquivo(s)", filesChanged)), styles.UnselectedStyle.Render("sincronizados")))
	}

	b.WriteString("\n\n")
	b.WriteString(styles.HelpStyle.Render("enter: retornar • esc: voltar • q: sair"))

	return b.String()
}
