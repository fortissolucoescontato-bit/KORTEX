package upgrade

import (
	"fmt"
	"strings"
)

// RenderUpgradeReport produces a plain-text report of upgrade results.
// Designed for the CLI path — no lipgloss, no color.
func RenderUpgradeReport(report UpgradeReport) string {
	var b strings.Builder

	if report.DryRun {
		b.WriteString("Upgrade (Simulação)\n")
	} else {
		b.WriteString("Upgrade\n")
	}
	b.WriteString("=======\n\n")

	// Preamble: explain what upgrade does (binary-only, no install/sync).
	b.WriteString("  Atualiza apenas os binários das ferramentas gerenciadas.\n")
	b.WriteString("  As configurações dos agentes são preservadas — não é feita instalação ou sync.\n\n")

	if len(report.Results) == 0 {
		b.WriteString("  Nenhum upgrade disponível. Todas as ferramentas estão atualizadas.\n")
		return b.String()
	}

	succeeded := 0
	failed := 0
	skipped := 0

	for _, r := range report.Results {
		icon := upgradeIcon(r.Status)
		fmt.Fprintf(&b, "  %s %-12s", icon, r.ToolName)

		switch r.Status {
		case UpgradeSucceeded:
			fmt.Fprintf(&b, "  %s → %s\n", r.OldVersion, r.NewVersion)
			succeeded++
		case UpgradeFailed:
			errMsg := ""
			if r.Err != nil {
				errMsg = r.Err.Error()
			}
			fmt.Fprintf(&b, "  FALHOU: %s\n", errMsg)
			failed++
		case UpgradeSkipped:
			if r.ManualHint != "" {
				fmt.Fprintf(&b, "  atualização manual necessária: %s\n", r.ManualHint)
			} else if report.DryRun {
				fmt.Fprintf(&b, "  %s → %s  (simulação)\n", r.OldVersion, r.NewVersion)
			} else {
				fmt.Fprintf(&b, "  pulado\n")
			}
			skipped++
		}
	}

	b.WriteString("\n")

	if report.BackupID != "" {
		fmt.Fprintf(&b, "  Backup da config: %s\n", report.BackupID)
	}
	if report.BackupWarning != "" {
		fmt.Fprintf(&b, "  AVISO: %s\n", report.BackupWarning)
	}

	if report.DryRun {
		// Count only actionable upgrades (no ManualHint) as pending.
		// Manual-hint items (DevBuild, VersionUnknown) will not run even
		// without --dry-run, so counting them as "pending" is misleading.
		actionable := 0
		for _, r := range report.Results {
			if r.Status == UpgradeSkipped && r.ManualHint == "" {
				actionable++
			}
		}
		if actionable > 0 {
			fmt.Fprintf(&b, "  %d upgrade(s) pendente(s). Execute sem --dry-run para aplicar.\n", actionable)
		}
		if skipped-actionable > 0 {
			fmt.Fprintf(&b, "  %d ferramenta(s) requerem atenção manual (veja dicas acima).\n", skipped-actionable)
		}
		if actionable == 0 && skipped == 0 {
			b.WriteString("  Nenhum upgrade acionável encontrado.\n")
		}
	} else {
		fmt.Fprintf(&b, "  %d com sucesso, %d falharam, %d pulados.\n", succeeded, failed, skipped)
	}

	return b.String()
}

// upgradeIcon returns a status indicator for upgrade CLI output.
func upgradeIcon(status ToolUpgradeStatus) string {
	switch status {
	case UpgradeSucceeded:
		return "[ok]"
	case UpgradeFailed:
		return "[!!]"
	case UpgradeSkipped:
		return "[--]"
	default:
		return "[  ]"
	}
}
