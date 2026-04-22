package cli

import (
	"fmt"
	"strings"

	"github.com/fortissolucoescontato-bit/kortex/internal/model"
	"github.com/fortissolucoescontato-bit/kortex/internal/planner"
	"github.com/fortissolucoescontato-bit/kortex/internal/system"
)

func RenderDryRun(result InstallResult) string {
	b := &strings.Builder{}

	_, _ = fmt.Fprintln(b, "Kortex Stack dry-run (Simulação)")
	_, _ = fmt.Fprintln(b, "=================================")
	_, _ = fmt.Fprintf(b, "Agentes: %s\n", joinAgentIDs(result.Resolved.Agents))
	_, _ = fmt.Fprintf(b, "Agentes não suportados: %s\n", joinAgentIDs(result.Resolved.UnsupportedAgents))
	_, _ = fmt.Fprintf(b, "Persona: %s\n", result.Selection.Persona)
	_, _ = fmt.Fprintf(b, "Preset: %s\n", result.Selection.Preset)
	if result.Selection.SDDMode != "" {
		_, _ = fmt.Fprintf(b, "Modo SDD: %s\n", result.Selection.SDDMode)
	}
	_, _ = fmt.Fprintf(b, "Ordem dos componentes: %s\n", joinComponentIDs(result.Resolved.OrderedComponents))
	_, _ = fmt.Fprintf(b, "Dependências auto-adicionadas: %s\n", joinComponentIDs(result.Resolved.AddedDependencies))
	_, _ = fmt.Fprintf(b, "Decisão de plataforma: %s\n", formatPlatformDecision(result.Review.PlatformDecision))
	_, _ = fmt.Fprintf(b, "Etapas de preparação: %d\n", len(result.Plan.Prepare))
	_, _ = fmt.Fprintf(b, "Etapas de aplicação: %d\n", len(result.Plan.Apply))

	if len(result.Dependencies.Dependencies) > 0 {
		_, _ = fmt.Fprintln(b, "")
		_, _ = fmt.Fprintln(b, system.RenderDependencyReport(result.Dependencies))
	}

	return strings.TrimRight(b.String(), "\n")
}

func joinAgentIDs(values []model.AgentID) string {
	if len(values) == 0 {
		return "nenhum"
	}

	parts := make([]string, 0, len(values))
	for _, value := range values {
		parts = append(parts, string(value))
	}
	return strings.Join(parts, ",")
}

func joinComponentIDs(values []model.ComponentID) string {
	if len(values) == 0 {
		return "nenhum"
	}

	parts := make([]string, 0, len(values))
	for _, value := range values {
		parts = append(parts, string(value))
	}
	return strings.Join(parts, ",")
}

func formatPlatformDecision(decision planner.PlatformDecision) string {
	osName := decision.OS
	if strings.TrimSpace(osName) == "" {
		osName = "desconhecido"
	}

	distro := decision.LinuxDistro
	if strings.TrimSpace(distro) == "" {
		distro = "n/a"
	}

	manager := decision.PackageManager
	if strings.TrimSpace(manager) == "" {
		manager = "n/a"
	}

	status := "não suportado"
	if decision.Supported {
		status = "suportado"
	}

	return fmt.Sprintf("os=%s distro=%s package-manager=%s status=%s", osName, distro, manager, status)
}
