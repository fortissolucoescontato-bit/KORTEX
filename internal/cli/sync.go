package cli

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/fortissolucoescontato-bit/kortex/internal/agents"
	"github.com/fortissolucoescontato-bit/kortex/internal/backup"
	"github.com/fortissolucoescontato-bit/kortex/internal/components/kortexengram"
	"github.com/fortissolucoescontato-bit/kortex/internal/components/kortex-cli"
	"github.com/fortissolucoescontato-bit/kortex/internal/components/mcp"
	"github.com/fortissolucoescontato-bit/kortex/internal/components/permissions"
	"github.com/fortissolucoescontato-bit/kortex/internal/components/sdd"
	"github.com/fortissolucoescontato-bit/kortex/internal/components/skills"
	"github.com/fortissolucoescontato-bit/kortex/internal/components/theme"
	"github.com/fortissolucoescontato-bit/kortex/internal/model"
	"github.com/fortissolucoescontato-bit/kortex/internal/pipeline"
	"github.com/fortissolucoescontato-bit/kortex/internal/state"
	"github.com/fortissolucoescontato-bit/kortex/internal/verify"
)

// SyncFlags holds parsed CLI flags for the sync command.
type SyncFlags struct {
	Agents             []string
	Skills             []string
	SDDMode            string
	SDDProfileStrategy string
	StrictTDD          bool
	IncludePermissions bool
	IncludeTheme       bool
	DryRun             bool
	// Profiles holds named SDD profiles parsed from --profile flags.
	// Each entry is populated by parseProfileFlag and augmented by
	// parseProfilePhaseFlag.
	Profiles []model.Profile
	// rawProfiles and rawProfilePhases hold the raw string values from
	// --profile and --profile-phase flags before parsing into model.Profile.
	rawProfiles      []string
	rawProfilePhases []string
}

// SyncResult holds the outcome of a sync execution.
type SyncResult struct {
	Agents    []model.AgentID
	Selection model.Selection
	Plan      pipeline.StagePlan
	Execution pipeline.ExecutionResult
	Verify    verify.Report
	DryRun    bool
	// NoOp is true when no managed asset changes were needed:
	// either no agents were discovered/provided, or all managed assets
	// were already current (idempotent re-sync).
	NoOp bool
	// FilesChanged is the number of managed files actually written or updated
	// during this sync. Zero means all assets were already current.
	FilesChanged int
}

// ParseSyncFlags parses the CLI arguments for the sync subcommand.
func ParseSyncFlags(args []string) (SyncFlags, error) {
	var opts SyncFlags

	fs := flag.NewFlagSet("sync", flag.ContinueOnError)
	fs.SetOutput(ioDiscard{})
	registerListFlag(fs, "agent", &opts.Agents)
	registerListFlag(fs, "agents", &opts.Agents)
	registerListFlag(fs, "skill", &opts.Skills)
	registerListFlag(fs, "skills", &opts.Skills)
	fs.StringVar(&opts.SDDMode, "sdd-mode", "", "modo do orquestrador SDD: single ou multi (padrão: single)")
	fs.StringVar(&opts.SDDProfileStrategy, "sdd-profile-strategy", "", "estratégia de sincronização de perfil OpenCode SDD: generated-multi ou external-single-active (padrão: auto-detecção)")
	fs.BoolVar(&opts.StrictTDD, "strict-tdd", false, "ativar modo TDD rigoroso para agentes SDD (RED → GREEN → REFACTOR)")
	fs.BoolVar(&opts.IncludePermissions, "include-permissions", false, "incluir componente de permissões na sincronização")
	fs.BoolVar(&opts.IncludeTheme, "include-theme", false, "incluir componente de tema na sincronização")
	fs.BoolVar(&opts.DryRun, "dry-run", false, "pré-visualizar o plano sem executar")
	registerListFlag(fs, "profile", &opts.rawProfiles)
	registerListFlag(fs, "profile-phase", &opts.rawProfilePhases)

	if err := fs.Parse(args); err != nil {
		return SyncFlags{}, err
	}

	if fs.NArg() > 0 {
		return SyncFlags{}, fmt.Errorf("argumento de sincronização inesperado %q", fs.Arg(0))
	}

	strategy, err := parseProfileSyncStrategy(opts.SDDProfileStrategy)
	if err != nil {
		return SyncFlags{}, err
	}
	opts.SDDProfileStrategy = string(strategy)

	// Parse --profile flags into model.Profile values.
	if len(opts.rawProfiles) > 0 || len(opts.rawProfilePhases) > 0 {
		profiles, err := parseProfileFlags(opts.rawProfiles, opts.rawProfilePhases)
		if err != nil {
			return SyncFlags{}, err
		}
		opts.Profiles = profiles
	}

	return opts, nil
}

func parseProfileSyncStrategy(raw string) (model.SDDProfileStrategyID, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return "", nil
	}

	switch model.SDDProfileStrategyID(value) {
	case model.SDDProfileStrategyGeneratedMulti, model.SDDProfileStrategyExternalSingleActive:
		return model.SDDProfileStrategyID(value), nil
	default:
		return "", fmt.Errorf("estratégia sdd-profile-strategy %q não suportada (válidas: generated-multi, external-single-active)", raw)
	}
}

// parseProfileFlags converts the raw --profile and --profile-phase string values
// into a slice of model.Profile. Returns an error if any value is malformed.
//
// --profile format:  name:provider/model
// --profile-phase format: name:phase:provider/model
func parseProfileFlags(rawProfiles, rawProfilePhases []string) ([]model.Profile, error) {
	// Build a map of profile name → profile so we can merge phase assignments.
	profileMap := make(map[string]*model.Profile)
	profileOrder := make([]string, 0, len(rawProfiles))

	for _, raw := range rawProfiles {
		p, err := parseProfileFlag(raw)
		if err != nil {
			return nil, err
		}
		profileMap[p.Name] = &p
		profileOrder = append(profileOrder, p.Name)
	}

	for _, raw := range rawProfilePhases {
		name, phase, assignment, err := parseProfilePhaseFlag(raw)
		if err != nil {
			return nil, err
		}
		entry, exists := profileMap[name]
		if !exists {
			// Profile referenced in --profile-phase but not declared in --profile.
			// Create a minimal entry so phase assignments are not lost.
			newProfile := model.Profile{Name: name, PhaseAssignments: make(map[string]model.ModelAssignment)}
			profileMap[name] = &newProfile
			profileOrder = append(profileOrder, name)
			entry = profileMap[name]
		}
		if entry.PhaseAssignments == nil {
			entry.PhaseAssignments = make(map[string]model.ModelAssignment)
		}
		entry.PhaseAssignments[phase] = assignment
	}

	profiles := make([]model.Profile, 0, len(profileOrder))
	seen := make(map[string]bool)
	for _, name := range profileOrder {
		if seen[name] {
			continue
		}
		seen[name] = true
		profiles = append(profiles, *profileMap[name])
	}
	return profiles, nil
}

// parseProfileFlag parses a single --profile value of the form "name:provider/model".
// Returns an error for empty name, reserved names, or missing separator.
func parseProfileFlag(raw string) (model.Profile, error) {
	colonIdx := strings.Index(raw, ":")
	if colonIdx <= 0 {
		return model.Profile{}, fmt.Errorf("--profile %q: formato inválido, esperado nome:provedor/modelo", raw)
	}
	name := raw[:colonIdx]
	modelSpec := raw[colonIdx+1:]

	if err := sdd.ValidateProfileName(name); err != nil {
		return model.Profile{}, fmt.Errorf("--profile %q: %w", raw, err)
	}

	assignment, err := parseModelSpec(modelSpec)
	if err != nil {
		return model.Profile{}, fmt.Errorf("--profile %q: %w", raw, err)
	}

	return model.Profile{
		Name:              name,
		OrchestratorModel: assignment,
		PhaseAssignments:  make(map[string]model.ModelAssignment),
	}, nil
}

// parseProfilePhaseFlag parses a single --profile-phase value of the form
// "name:phase:provider/model".
func parseProfilePhaseFlag(raw string) (name, phase string, assignment model.ModelAssignment, err error) {
	parts := strings.SplitN(raw, ":", 3)
	if len(parts) != 3 {
		return "", "", model.ModelAssignment{}, fmt.Errorf("--profile-phase %q: formato inválido, esperado nome:fase:provedor/modelo", raw)
	}
	name = parts[0]
	phase = parts[1]
	modelSpec := parts[2]

	if name == "" {
		return "", "", model.ModelAssignment{}, fmt.Errorf("--profile-phase %q: o nome do perfil não deve estar vazio", raw)
	}
	if err = sdd.ValidateProfileName(name); err != nil {
		return "", "", model.ModelAssignment{}, fmt.Errorf("--profile-phase %q: %w", raw, err)
	}
	if phase == "" {
		return "", "", model.ModelAssignment{}, fmt.Errorf("--profile-phase %q: a fase não deve estar vazia", raw)
	}
	// Validate that the phase is a known SDD phase name.
	knownPhases := sdd.ProfilePhaseOrder()
	validPhase := false
	for _, p := range knownPhases {
		if p == phase {
			validPhase = true
			break
		}
	}
	if !validPhase {
		return "", "", model.ModelAssignment{}, fmt.Errorf("--profile-phase %q: fase desconhecida %q; as fases válidas são: %v", raw, phase, knownPhases)
	}

	assignment, err = parseModelSpec(modelSpec)
	if err != nil {
		return "", "", model.ModelAssignment{}, fmt.Errorf("--profile-phase %q: %w", raw, err)
	}
	return name, phase, assignment, nil
}

// parseModelSpec parses a "provider/model" or "provider:model" string into a
// ModelAssignment. Returns an error if the spec is empty or has no separator.
func parseModelSpec(spec string) (model.ModelAssignment, error) {
	// Try slash separator first (common CLI format: anthropic/claude-haiku-3-5),
	// then colon (opencode internal format: anthropic:claude-haiku-3-5).
	sep := -1
	for i, c := range spec {
		if c == '/' || c == ':' {
			sep = i
			break
		}
	}
	if sep <= 0 {
		return model.ModelAssignment{}, fmt.Errorf("especificação de modelo %q inválida: esperado provedor/modelo ou provedor:modelo", spec)
	}
	providerID := spec[:sep]
	modelID := spec[sep+1:]
	if providerID == "" || modelID == "" {
		return model.ModelAssignment{}, fmt.Errorf("especificação de modelo %q inválida: tanto o provedor quanto o modelo devem ser preenchidos", spec)
	}
	return model.ModelAssignment{ProviderID: providerID, ModelID: modelID}, nil
}

// BuildSyncSelection builds a model.Selection for the sync command.
//
// Default sync scope: SDD, KortexEngram, Context7, KortexCLI, Skills.
// Excluded by default: Persona, Permissions, Theme (user-config-adjacent).
// Permissions and Theme can be opted-in via flags.
//
// This is the reusable managed-asset sync contract. A future `upgrade --sync`
// flow can call this function to get the same managed-only selection semantics.
func BuildSyncSelection(flags SyncFlags, agentIDs []model.AgentID) model.Selection {
	components := []model.ComponentID{
		model.ComponentSDD,
		model.ComponentKortexEngram,
		model.ComponentContext7,
		model.ComponentKortexCLI,
		model.ComponentSkills,
	}

	if flags.IncludePermissions {
		components = append(components, model.ComponentPermission)
	}
	if flags.IncludeTheme {
		components = append(components, model.ComponentTheme)
	}

	sddMode := model.SDDModeID(flags.SDDMode)

	var skillIDs []model.SkillID
	for _, raw := range flags.Skills {
		skillIDs = append(skillIDs, model.SkillID(raw))
	}

	return model.Selection{
		Agents:             agentIDs,
		Components:         components,
		SDDMode:            sddMode,
		SDDProfileStrategy: model.SDDProfileStrategyID(flags.SDDProfileStrategy),
		StrictTDD:          flags.StrictTDD,
		Skills:             skillIDs,
		Profiles:           flags.Profiles,
		// Preset is set to full-carbon so selectedSkillIDs() returns the
		// correct default skill set when no explicit skills are provided.
		Preset: model.PresetFullKortex,
	}
}

// DiscoverAgents returns the agent IDs to sync.
//
// Discovery order:
//  1. Persisted state (~/.kortex/state.json) — written at install time.
//     When present and non-empty, only the agents the user explicitly installed
//     are returned. This prevents sync from injecting into every IDE config dir
//     that happens to exist on the system (issue #107).
//  2. Filesystem fallback — delegates to agents.DiscoverInstalled with the
//     default registry. Used when state.json is absent (users who installed
//     before state persistence was added) or empty.
//
// When --agents is provided explicitly, callers should pass those IDs directly
// instead of calling DiscoverAgents.
func DiscoverAgents(homeDir string) []model.AgentID {
	mgr, err := state.NewManager(homeDir)
	if err == nil {
		defer mgr.Close()
		installed, err := mgr.GetInstalledAgents()
		if err == nil && len(installed) > 0 {
			ids := make([]model.AgentID, 0, len(installed))
			for _, a := range installed {
				ids = append(ids, model.AgentID(a))
			}
			return ids
		}
	}

	// Fallback: filesystem discovery (backward compat for users who installed
	// before state persistence was added).
	reg, err := agents.NewDefaultRegistry()
	if err != nil {
		// Registry construction only fails if a duplicate adapter is registered,
		// which would indicate a programming error. Treat as no agents found
		// rather than propagating — callers treat an empty result as a no-op.
		return nil
	}

	installed := agents.DiscoverInstalled(reg, homeDir)
	ids := make([]model.AgentID, 0, len(installed))
	for _, a := range installed {
		ids = append(ids, a.ID)
	}
	return ids
}

// syncRuntime mirrors installRuntime but builds a sync-scoped StagePlan.
// It reuses backup/rollback infrastructure but only calls inject functions —
// no agentInstallStep, no KortexEngram setup, no persona.
type syncRuntime struct {
	homeDir      string
	workspaceDir string
	selection    model.Selection
	agentIDs     []model.AgentID
	backupRoot   string
	state        *runtimeState
	filesChanged int // accumulates changed-file count across all component steps
}

func newSyncRuntime(homeDir string, selection model.Selection) (*syncRuntime, error) {
	backupRoot := filepath.Join(homeDir, ".kortex", "backups")
	if err := os.MkdirAll(backupRoot, 0o755); err != nil {
		return nil, fmt.Errorf("falha ao criar diretório raiz de backup %q: %w", backupRoot, err)
	}

	workspaceDir, _ := os.Getwd()

	return &syncRuntime{
		homeDir:      homeDir,
		workspaceDir: workspaceDir,
		selection:    selection,
		agentIDs:     selection.Agents,
		backupRoot:   backupRoot,
		state:        &runtimeState{},
	}, nil
}

func (r *syncRuntime) stagePlan() pipeline.StagePlan {
	adapters := resolveAdapters(r.agentIDs)
	targets := syncBackupTargets(r.homeDir, r.selection, adapters)

	prepare := []pipeline.Step{
		prepareBackupStep{
			id:          "prepare:backup-snapshot",
			snapshotter: backup.NewSnapshotter(),
			snapshotDir: filepath.Join(r.backupRoot, time.Now().UTC().Format("20060102150405.000000000")),
			targets:     targets,
			state:       r.state,
			backupRoot:  r.backupRoot,
			source:      backup.BackupSourceSync,
			description: "snapshot pré-sincronização",
			appVersion:  AppVersion,
		},
	}

	apply := []pipeline.Step{
		rollbackRestoreStep{id: "apply:rollback-restore", state: r.state},
	}

	for _, component := range r.selection.Components {
		apply = append(apply, componentSyncStep{
			id:           "sync:component:" + string(component),
			component:    component,
			homeDir:      r.homeDir,
			workspaceDir: r.workspaceDir,
			agents:       r.agentIDs,
			selection:    r.selection,
			filesChanged: &r.filesChanged,
		})
	}

	return pipeline.StagePlan{Prepare: prepare, Apply: apply}
}

// syncBackupTargets returns the file paths that need to be backed up
// before sync executes. Uses the same componentPaths logic as install.
func syncBackupTargets(homeDir string, selection model.Selection, adapters []agents.Adapter) []string {
	paths := map[string]struct{}{}
	for _, component := range selection.Components {
		for _, path := range componentPaths(homeDir, selection, adapters, component) {
			paths[path] = struct{}{}
		}
	}

	targets := make([]string, 0, len(paths))
	for path := range paths {
		targets = append(targets, path)
	}
	return targets
}

// componentSyncStep is the sync-specific apply step.
// Unlike componentApplyStep, it ONLY calls inject functions —
// no binary install, no KortexEngram setup, no persona injection.
//
// filesChanged is a shared counter pointer. Each step increments it by the
// number of files that were actually written (i.e., whose content changed).
// This lets RunSync detect a true no-op when all assets are already current.
type componentSyncStep struct {
	id           string
	component    model.ComponentID
	homeDir      string
	workspaceDir string
	agents       []model.AgentID
	selection    model.Selection
	filesChanged *int
}

func (s componentSyncStep) ID() string {
	return s.id
}

func (s componentSyncStep) Run() error {
	adapters := resolveAdapters(s.agents)

	switch s.component {
	case model.ComponentKortexEngram:
		// Sync: inject MCP config + system prompt protocol only.
		// NO binary install. NO KortexEngram setup.
		for _, adapter := range adapters {
			res, err := kortexengram.Inject(s.homeDir, adapter)
			if err != nil {
				return fmt.Errorf("falha ao sincronizar KortexEngram para %q: %w", adapter.Agent(), err)
			}
			s.countChanged(boolToInt(res.Changed))
		}
		return nil

	case model.ComponentContext7:
		for _, adapter := range adapters {
			res, err := mcp.Inject(s.homeDir, adapter)
			if err != nil {
				return fmt.Errorf("falha ao sincronizar context7 para %q: %w", adapter.Agent(), err)
			}
			s.countChanged(boolToInt(res.Changed))
		}
		return nil

	case model.ComponentSDD:
		profileStrategy := sdd.ResolveProfileStrategy(s.homeDir, s.selection.SDDProfileStrategy)

		// Resolve profiles for injection:
		// - When profiles are explicitly provided (TUI/CLI), use them directly.
		// - On a regular sync (no explicit profiles), detect existing named profiles
		//   from disk so their orchestrator prompts are refreshed from updated embedded
		//   assets while model assignments are preserved.
		profiles := s.selection.Profiles
		if len(profiles) == 0 && profileStrategy != model.SDDProfileStrategyExternalSingleActive {
			settingsPath := ""
			for _, adapter := range adapters {
				if adapter.Agent() == model.AgentOpenCode {
					settingsPath = adapter.SettingsPath(s.homeDir)
					break
				}
			}
			if settingsPath != "" {
				detected, detectErr := sdd.DetectProfiles(settingsPath)
				if detectErr == nil {
					profiles = detected
				}
				// If detect fails (e.g. file missing), silently skip — no profiles to refresh.
			}
		}

		// If profiles exist (explicit or detected), SDDModeMulti is required:
		// shared prompt files must be written and {file:...} refs must resolve.
		sddMode := s.selection.SDDMode
		if profileStrategy == model.SDDProfileStrategyExternalSingleActive {
			sddMode = model.SDDModeMulti
		} else if len(profiles) > 0 && sddMode == "" {
			sddMode = model.SDDModeMulti
		}

		for _, adapter := range adapters {
			opts := sdd.InjectOptions{
				OpenCodeModelAssignments:           s.selection.ModelAssignments,
				ClaudeModelAssignments:             s.selection.ClaudeModelAssignments,
				KiroModelAssignments:               s.selection.KiroModelAssignments,
				WorkspaceDir:                       s.workspaceDir,
				StrictTDD:                          s.selection.StrictTDD,
				PreserveOpenCodeOrchestratorPrompt: profileStrategy == model.SDDProfileStrategyExternalSingleActive,
				Profiles:                           profiles,
			}
			res, err := sdd.Inject(s.homeDir, adapter, sddMode, opts)
			if err != nil {
				return fmt.Errorf("falha ao sincronizar sdd para %q: %w", adapter.Agent(), err)
			}
			s.countChanged(boolToInt(res.Changed))
		}
		return nil

	case model.ComponentSkills:
		skillIDs := selectedSkillIDs(s.selection)
		if len(skillIDs) == 0 {
			return nil
		}
		for _, adapter := range adapters {
			res, err := skills.Inject(s.homeDir, adapter, skillIDs)
			if err != nil {
				return fmt.Errorf("falha ao sincronizar skills para %q: %w", adapter.Agent(), err)
			}
			s.countChanged(boolToInt(res.Changed))
		}
		return nil

	case model.ComponentKortexCLI:
		// Sync: ensure runtime assets are current and inject config.
		// NO binary install.
		if err := kortex.EnsureRuntimeAssets(s.homeDir); err != nil {
			return fmt.Errorf("falha ao sincronizar ativos de runtime do kortex: %w", err)
		}
		if runtime.GOOS == "windows" {
			if err := kortex.EnsurePowerShellShim(s.homeDir); err != nil {
				return fmt.Errorf("falha ao garantir o shim powershell do kortex: %w", err)
			}
		}
		res, err := kortex.Inject(s.homeDir, s.agents)
		if err != nil {
			return fmt.Errorf("falha ao sincronizar configuração do kortex: %w", err)
		}
		// Count KortexCLI files changed based on individual Changed flags.
		s.countChanged(boolToInt(res.ConfigChanged) + boolToInt(res.AgentsChanged))
		return nil

	case model.ComponentPermission:
		// Opt-in only — reached when --include-permissions is set.
		for _, adapter := range adapters {
			res, err := permissions.Inject(s.homeDir, adapter)
			if err != nil {
				return fmt.Errorf("falha ao sincronizar permissões para %q: %w", adapter.Agent(), err)
			}
			s.countChanged(boolToInt(res.Changed))
		}
		return nil

	case model.ComponentTheme:
		// Opt-in only — reached when --include-theme is set.
		for _, adapter := range adapters {
			res, err := theme.Inject(s.homeDir, adapter)
			if err != nil {
				return fmt.Errorf("falha ao sincronizar tema para %q: %w", adapter.Agent(), err)
			}
			s.countChanged(boolToInt(res.Changed))
		}
		return nil

	default:
		return fmt.Errorf("o componente %q não é suportado no runtime de sincronização", s.component)
	}
}

// countChanged adds n to the shared filesChanged counter (nil-safe).
func (s componentSyncStep) countChanged(n int) {
	if s.filesChanged != nil && n > 0 {
		*s.filesChanged += n
	}
}

// boolToInt converts a boolean to 0 or 1.
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// RunSyncWithSelection is the programmatic entry point for sync.
// It skips flag parsing and agent discovery — the caller provides the homeDir
// and a fully-built Selection (agents + components + options).
// This is the function the TUI calls directly to avoid CLI flag parsing.
func RunSyncWithSelection(homeDir string, selection model.Selection) (SyncResult, error) {
	agentIDs := selection.Agents

	result := SyncResult{
		Agents:    agentIDs,
		Selection: selection,
	}

	// No-op path: no agents were discovered or provided.
	// Per spec: "No managed assets to sync — system completes without modifying
	// unrelated files and reports that no managed sync actions were needed."
	if len(agentIDs) == 0 {
		result.NoOp = true
		return result, nil
	}

	rt, err := newSyncRuntime(homeDir, selection)
	if err != nil {
		return result, err
	}

	stagePlan := rt.stagePlan()
	result.Plan = stagePlan

	orchestrator := pipeline.NewOrchestrator(pipeline.DefaultRollbackPolicy())
	result.Execution = orchestrator.Execute(stagePlan)
	if result.Execution.Err != nil {
		return result, fmt.Errorf("falha ao executar pipeline de sincronização: %w", result.Execution.Err)
	}

	// Capture how many managed assets were actually changed.
	result.FilesChanged = rt.filesChanged

	// True no-op: agents were discovered but all managed assets were already
	// current — no file was written or updated. Per spec scenario:
	// "No managed assets to sync — system completes without modifying files
	// and reports that no managed sync actions were needed."
	if result.FilesChanged == 0 {
		result.NoOp = true
	}

	// Post-apply verification reuses the same component paths as install.
	result.Verify = runPostSyncVerification(homeDir, selection)
	if !result.Verify.Ready {
		return result, fmt.Errorf("verificação pós-sincronização falhou:\n%s", verify.RenderReport(result.Verify))
	}

	return result, nil
}

// RunSync is the top-level sync entry point, parallel to RunInstall.
// It parses CLI flags, discovers agents, builds the selection, then delegates
// to RunSyncWithSelection for the actual sync execution.
func RunSync(args []string) (SyncResult, error) {
	flags, err := ParseSyncFlags(args)
	if err != nil {
		return SyncResult{}, err
	}

	homeDir, err := osUserHomeDir()
	if err != nil {
		return SyncResult{}, fmt.Errorf("falha ao resolver diretório home do usuário: %w", err)
	}

	// Resolve agents: explicit flag takes precedence over auto-discovery.
	var agentIDs []model.AgentID
	if len(flags.Agents) > 0 {
		agentIDs = asAgentIDs(flags.Agents)
	} else {
		agentIDs = DiscoverAgents(homeDir)
	}
	agentIDs = unique(agentIDs)

	selection := BuildSyncSelection(flags, agentIDs)

	// Load persisted model assignments from state when not provided via flags.
	// This is the key fix: without this, every CLI sync falls back to the
	// "balanced" preset and silently overwrites the user's model choices.
	if len(selection.ModelAssignments) == 0 {
		mgr, err := state.NewManager(homeDir)
		if err == nil {
			defer mgr.Close()
			// Load assignments for all selected agents + "claude" and "kiro" which
			// store global SDD model preferences.
			agentsToLoad := []string{"claude", "kiro"}
			for _, a := range selection.Agents {
				agentsToLoad = append(agentsToLoad, string(a))
			}

			for _, agent := range agentsToLoad {
				dbAssignments, err := mgr.GetAssignments(agent)
				if err == nil && len(dbAssignments) > 0 {
					for phase, mState := range dbAssignments {
						if agent == "claude" {
							if selection.ClaudeModelAssignments == nil {
								selection.ClaudeModelAssignments = make(map[string]model.ClaudeModelAlias)
							}
							selection.ClaudeModelAssignments[phase] = model.ClaudeModelAlias(mState.ModelID)
						} else if agent == "kiro" {
							if selection.KiroModelAssignments == nil {
								selection.KiroModelAssignments = make(map[string]model.ClaudeModelAlias)
							}
							selection.KiroModelAssignments[phase] = model.ClaudeModelAlias(mState.ModelID)
						} else {
							if selection.ModelAssignments == nil {
								selection.ModelAssignments = make(map[string]model.ModelAssignment)
							}
							selection.ModelAssignments[phase] = model.ModelAssignment{
								ProviderID: mState.ProviderID,
								ModelID:    mState.ModelID,
							}
						}
					}
				}
			}
		}
	}

	if flags.DryRun {
		// Build the plan for inspection, skip execution.
		result := SyncResult{
			Agents:    agentIDs,
			Selection: selection,
			DryRun:    true,
		}
		if len(agentIDs) == 0 {
			result.NoOp = true
			return result, nil
		}
		rt, err := newSyncRuntime(homeDir, selection)
		if err != nil {
			return result, err
		}
		result.Plan = rt.stagePlan()
		return result, nil
	}

	result, err := RunSyncWithSelection(homeDir, selection)
	if err != nil {
		return result, err
	}
	result.DryRun = false
	return result, nil
}

// RenderSyncReport renders a human-readable summary of a sync execution.
//
// Unlike verify.RenderReport (which shows verification check statuses), this
// function reports the managed sync actions that were executed — matching the
// spec requirement to surface "what was done" rather than "what was checked".
//
// No-op cases:
//   - No agents were discovered or specified (NoOp=true, Agents empty).
//   - All managed assets were already current (NoOp=true, FilesChanged=0).
func RenderSyncReport(result SyncResult) string {
	var b strings.Builder

	if result.NoOp {
		fmt.Fprintln(&b, "kortex sync — nenhuma ação de sincronização necessária")
		if len(result.Agents) == 0 {
			fmt.Fprintln(&b, "Nenhum agente foi descoberto ou especificado. Nada para sincronizar.")
		} else {
			fmt.Fprintf(&b, "Agentes: %s\n", joinAgentIDs(result.Agents))
			fmt.Fprintln(&b, "Todos os ativos gerenciados já estão atualizados. Nenhum arquivo alterado.")
		}
		return strings.TrimRight(b.String(), "\n")
	}

	if result.DryRun {
		fmt.Fprintln(&b, "kortex sync — simulação (dry-run)")
		fmt.Fprintf(&b, "Agentes: %s\n", joinAgentIDs(result.Agents))

		compParts := make([]string, 0, len(result.Selection.Components))
		for _, c := range result.Selection.Components {
			compParts = append(compParts, string(c))
		}
		if len(compParts) > 0 {
			fmt.Fprintf(&b, "Componentes gerenciados: %s\n", strings.Join(compParts, ", "))
		}
		fmt.Fprintf(&b, "Etapas de preparação: %d\n", len(result.Plan.Prepare))
		fmt.Fprintf(&b, "Etapas de aplicação: %d\n", len(result.Plan.Apply))
		return strings.TrimRight(b.String(), "\n")
	}

	fmt.Fprintln(&b, "kortex sync — sincronização gerenciada executada")
	fmt.Fprintf(&b, "Agentes sincronizados: %s\n", joinAgentIDs(result.Agents))

	compParts := make([]string, 0, len(result.Selection.Components))
	for _, c := range result.Selection.Components {
		compParts = append(compParts, string(c))
	}
	if len(compParts) > 0 {
		fmt.Fprintf(&b, "Componentes gerenciados sincronizados: %s\n", strings.Join(compParts, ", "))
	}

	// Report actual files changed — not the count of successful pipeline steps.
	// FilesChanged is 0 only when all assets were already current (no-op path
	// above handles that case). A non-zero value here reflects real writes.
	fmt.Fprintf(&b, "Ações de sincronização executadas: %d arquivos alterados\n", result.FilesChanged)

	if !result.Verify.Ready {
		fmt.Fprintln(&b, "")
		fmt.Fprintln(&b, "Post-sync verification:")
		fmt.Fprint(&b, verify.RenderReport(result.Verify))
	}

	return strings.TrimRight(b.String(), "\n")
}

// runPostSyncVerification verifies that managed files exist after sync.
func runPostSyncVerification(homeDir string, selection model.Selection) verify.Report {
	checks := make([]verify.Check, 0)
	adapters := resolveAdapters(selection.Agents)

	for _, component := range selection.Components {
		for _, path := range componentPaths(homeDir, selection, adapters, component) {
			currentPath := path
			checks = append(checks, verify.Check{
				ID:          "verify:sync:file:" + currentPath,
				Description: "synced file exists",
				Run: func(context.Context) error {
					if _, err := os.Stat(currentPath); err != nil {
						return err
					}
					return nil
				},
			})
		}
	}

	return verify.BuildReport(verify.RunChecks(context.Background(), checks))
}
