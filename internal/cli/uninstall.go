package cli

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fortissolucoescontato-bit/kortex/internal/catalog"
	componentuninstall "github.com/fortissolucoescontato-bit/kortex/internal/components/uninstall"
	"github.com/fortissolucoescontato-bit/kortex/internal/model"
)

type UninstallFlags struct {
	Agents     []string
	Components []string
	All        bool
	Yes        bool
}

func ParseUninstallFlags(args []string) (UninstallFlags, error) {
	var opts UninstallFlags

	fs := flag.NewFlagSet("uninstall", flag.ContinueOnError)
	fs.SetOutput(ioDiscard{})
	registerListFlag(fs, "agent", &opts.Agents)
	registerListFlag(fs, "agents", &opts.Agents)
	registerListFlag(fs, "component", &opts.Components)
	registerListFlag(fs, "components", &opts.Components)
	fs.BoolVar(&opts.All, "all", false, "remover configuração gerenciada para todos os agentes suportados")
	fs.BoolVar(&opts.Yes, "yes", false, "pular prompt de confirmação")
	fs.BoolVar(&opts.Yes, "y", false, "pular prompt de confirmação")

	if err := fs.Parse(args); err != nil {
		return UninstallFlags{}, err
	}
	if fs.NArg() > 0 {
		return UninstallFlags{}, fmt.Errorf("argumento de desinstalação inesperado %q", fs.Arg(0))
	}
	if opts.All && (len(opts.Agents) > 0 || len(opts.Components) > 0) {
		return UninstallFlags{}, fmt.Errorf("--all não pode ser combinado com --agent/--agents ou --component/--components")
	}
	if !opts.All && len(opts.Agents) == 0 {
		return UninstallFlags{}, fmt.Errorf("desinstalação parcial requer ao menos um --agent/--agents ou o uso de --all")
	}

	return opts, nil
}

func RunUninstall(args []string, stdout io.Writer) (componentuninstall.Result, error) {
	return runUninstallWithInput(args, stdout, os.Stdin)
}

func RunUninstallWithSelection(homeDir, workspaceDir string, agentIDs []model.AgentID, componentIDs []model.ComponentID) (componentuninstall.Result, error) {
	agents := make([]string, 0, len(agentIDs))
	for _, agentID := range agentIDs {
		agents = append(agents, string(agentID))
	}
	components := make([]string, 0, len(componentIDs))
	for _, componentID := range componentIDs {
		components = append(components, string(componentID))
	}
	return componentuninstall.PartialUninstall(homeDir, workspaceDir, AppVersion, agents, components)
}

func RunUninstallWithSelectionAndProfiles(homeDir, workspaceDir string, agentIDs []model.AgentID, componentIDs []model.ComponentID, profileNames []string, engramScope model.EngramUninstallScope) (componentuninstall.Result, error) {
	agents := make([]string, 0, len(agentIDs))
	for _, agentID := range agentIDs {
		agents = append(agents, string(agentID))
	}
	components := make([]string, 0, len(componentIDs))
	for _, componentID := range componentIDs {
		components = append(components, string(componentID))
	}
	return componentuninstall.PartialUninstallWithProfileSelection(homeDir, workspaceDir, AppVersion, agents, components, profileNames, engramScope)
}

func RenderUninstallReport(result componentuninstall.Result) string {
	var b strings.Builder

	_, _ = fmt.Fprintln(&b, "Desinstalação gerenciada concluída")
	if result.Manifest.ID != "" {
		_, _ = fmt.Fprintf(&b, "Backup: %s (%s)\n", result.Manifest.ID, result.Manifest.DisplayLabel())
		_, _ = fmt.Fprintf(&b, "Caminho do backup: %s\n", result.BackupPath)
	}
	_, _ = fmt.Fprintf(&b, "Arquivos alterados: %d\n", len(result.ChangedFiles))
	_, _ = fmt.Fprintf(&b, "Arquivos removidos: %d\n", len(result.RemovedFiles))
	_, _ = fmt.Fprintf(&b, "Diretórios removidos: %d\n", len(result.RemovedDirectories))
	if len(result.AgentsRemovedFromState) > 0 {
		_, _ = fmt.Fprintf(&b, "Estado (state.json) atualizado: removido %s\n", strings.Join(agentLabels(result.AgentsRemovedFromState), ", "))
	}
	appendPathSection(&b, "Arquivos reescritos", result.ChangedFiles)
	appendPathSection(&b, "Arquivos excluídos", result.RemovedFiles)
	appendPathSection(&b, "Diretórios excluídos", result.RemovedDirectories)
	appendPathSection(&b, "Limpeza manual necessária", result.ManualActions)

	return strings.TrimRight(b.String(), "\n")
}

func runUninstallWithInput(args []string, stdout io.Writer, stdin io.Reader) (componentuninstall.Result, error) {
	flags, err := ParseUninstallFlags(args)
	if err != nil {
		return componentuninstall.Result{}, err
	}

	homeDir, err := osUserHomeDir()
	if err != nil {
		return componentuninstall.Result{}, fmt.Errorf("falha ao resolver diretório home: %w", err)
	}
	workspaceDir, err := os.Getwd()
	if err != nil {
		return componentuninstall.Result{}, fmt.Errorf("falha ao resolver diretório de trabalho: %w", err)
	}

	if !flags.Yes {
		confirmed, err := promptUninstallConfirm(flags, stdout, stdin)
		if err != nil {
			return componentuninstall.Result{}, err
		}
		if !confirmed {
			_, _ = fmt.Fprintln(stdout, "desinstalação cancelada")
			return componentuninstall.Result{}, nil
		}
	}

	if flags.All {
		return componentuninstall.CompleteUninstall(homeDir, workspaceDir, AppVersion)
	}
	return componentuninstall.PartialUninstall(homeDir, workspaceDir, AppVersion, flags.Agents, flags.Components)
}

func promptUninstallConfirm(flags UninstallFlags, stdout io.Writer, stdin io.Reader) (bool, error) {
	if flags.All {
		_, _ = fmt.Fprintln(stdout, "Isso irá remover a configuração gerenciada pelo kortex de todos os agentes suportados.")
	} else {
		_, _ = fmt.Fprintf(stdout, "Isso irá remover a configuração gerenciada pelo kortex de: %s\n", strings.Join(agentLabelsFromStrings(flags.Agents), ", "))
	}
	if len(flags.Components) > 0 {
		_, _ = fmt.Fprintf(stdout, "Componentes: %s\n", strings.Join(flags.Components, ", "))
	} else {
		_, _ = fmt.Fprintln(stdout, "Componentes: todos os componentes desinstaláveis gerenciados")
	}
	_, _ = fmt.Fprintln(stdout, "Um backup será criado antes que qualquer arquivo seja modificado.")
	_, _ = fmt.Fprint(stdout, "Digite 'sim' para confirmar: ")

	scanner := bufio.NewScanner(stdin)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return false, fmt.Errorf("erro ao ler confirmação de desinstalação: %w", err)
		}
		return false, fmt.Errorf("nenhuma confirmação fornecida (use --yes para pular o prompt)")
	}
	return strings.EqualFold(strings.TrimSpace(scanner.Text()), "sim"), nil
}

func appendPathSection(b *strings.Builder, title string, paths []string) {
	if len(paths) == 0 {
		return
	}

	sorted := append([]string(nil), paths...)
	sort.Strings(sorted)
	cwd, cwdErr := os.Getwd()
	_, _ = fmt.Fprintf(b, "\n%s:\n", title)
	for _, path := range sorted {
		rel := path
		if cwdErr == nil {
			if r, relErr := filepath.Rel(cwd, path); relErr == nil && !strings.HasPrefix(r, "..") {
				rel = r
			}
		}
		_, _ = fmt.Fprintf(b, "  - %s\n", rel)
	}
}

func agentLabels(agentIDs []model.AgentID) []string {
	labels := make([]string, 0, len(agentIDs))
	for _, agentID := range agentIDs {
		labels = append(labels, agentLabel(agentID))
	}
	return labels
}

func agentLabelsFromStrings(agentIDs []string) []string {
	labels := make([]string, 0, len(agentIDs))
	for _, agentID := range agentIDs {
		labels = append(labels, agentLabel(model.AgentID(agentID)))
	}
	return labels
}

func agentLabel(agentID model.AgentID) string {
	for _, agent := range catalog.AllAgents() {
		if agent.ID == agentID {
			return agent.Name
		}
	}
	return string(agentID)
}
