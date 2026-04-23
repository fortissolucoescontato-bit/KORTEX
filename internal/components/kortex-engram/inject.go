package kortexengram

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fortissolucoescontato-bit/kortex/internal/agents"
	"github.com/fortissolucoescontato-bit/kortex/internal/assets"
	"github.com/fortissolucoescontato-bit/kortex/internal/components/filemerge"
	"github.com/fortissolucoescontato-bit/kortex/internal/model"
)

type InjectionResult struct {
	Changed bool
	Files   []string
}

// bootstrapper is an optional adapter capability: if an adapter implements
// this interface, any injector that writes Jinja modules will first ensure
// the base template (entry point) exists.
type bootstrapper interface {
	BootstrapTemplate(homeDir string) error
}

// kortexEngramLookPath is the function used to resolve the KortexEngram binary path.
// It is a package-level variable so it can be replaced in tests — both from
// within the KortexEngram package and from external test packages (e.g. golden_test.go).
// In production it is set to exec.LookPath.
var kortexEngramLookPath = exec.LookPath

// SetLookPathForTest replaces kortexEngramLookPath with a mock for the duration of
// a test and restores the original after the test completes. Exported so that
// external test packages (e.g. golden_test.go in components) can control the
// resolved KortexEngram path.
func SetLookPathForTest(t interface {
	Helper()
	Cleanup(func())
}, result, errMsg string) {
	t.Helper()
	orig := kortexEngramLookPath
	kortexEngramLookPath = func(string) (string, error) {
		if errMsg != "" {
			return "", fmt.Errorf("%s", errMsg)
		}
		return result, nil
	}
	t.Cleanup(func() { kortexEngramLookPath = orig })
}

// resolveKortexEngramCommand attempts to resolve the KortexEngram binary to an absolute
// path using exec.LookPath. If found, it returns the absolute path and true.
// If not found (e.g. binary not yet installed), it returns "kortex-engram" and false.
// This is used to write the most stable command possible into MCP configs:
// an absolute path survives across environments where PATH is not fully
// inherited (e.g. Windsurf, IDEs that launch without a login shell).
func resolveKortexEngramCommand() (string, bool) {
	// Try 'KortexEngram' first (new dedicated server name)
	if p, err := kortexEngramLookPath("kortex-engram"); err == nil && p != "" {
		return p, true
	}
	// Fallback to 'kortex' (previous rebranding attempt)
	if p, err := kortexEngramLookPath("kortex"); err == nil && p != "" {
		return p, true
	}
	// Fallback to 'KortexEngram' (legacy)
	if p, err := kortexEngramLookPath("kortex-engram"); err == nil && p != "" {
		return p, true
	}
	return "kortex-engram", false
}

// kortexEngramServerJSON returns the MCP server config bytes, using the absolute
// path to the KortexEngram binary if it can be resolved via PATH.
func kortexEngramServerJSON() []byte {
	cmd, _ := resolveKortexEngramCommand()
	return kortexEngramServerJSONWithCmd(cmd)
}

// kortexEngramServerJSONWithCmd returns the MCP server config bytes for a specific
// command.
func kortexEngramServerJSONWithCmd(cmd string) []byte {
	cfg := map[string]any{
		"command": cmd,
		"args":    []string{"mcp", "--tools=agent"},
	}
	b, _ := json.MarshalIndent(cfg, "", "  ")
	return append(b, '\n')
}

// kortexEngramOverlayJSON returns the settings overlay JSON (used for merge-into-settings
// and MCPConfigFile strategies), with the resolved KortexEngram command.
func kortexEngramOverlayJSON(agentID model.AgentID, cmd string) []byte {
	var cfg map[string]any
	if agentID == model.AgentOpenCode || agentID == model.AgentKilocode {
		// OpenCode 1.3.3+ requires command as an array for type:local servers.
		// The separate "args" field is not accepted; all args must be in the
		// command array itself.
		//
		// Use the __replace__ sentinel so that MergeJSONObjects replaces the
		// entire mcp.KortexEngram object atomically instead of deep-merging into it.
		// Without this, users upgrading from v1.11.3 (which had a separate
		// "args" key) would end up with both "args" and the new array "command"
		// in their config, which is invalid for OpenCode 1.3.3.
		cfg = map[string]any{
			"mcp": map[string]any{
				"kortex-engram": map[string]any{
					"__replace__": map[string]any{
						"command": []string{cmd, "mcp", "--tools=agent"},
						"type":    "local",
					},
				},
			},
		}
	} else {
		cfg = map[string]any{
			"mcpServers": map[string]any{
				"kortex-engram": map[string]any{
					"command": cmd,
					"args":    []string{"mcp", "--tools=agent"},
				},
			},
		}
	}
	b, _ := json.MarshalIndent(cfg, "", "  ")
	return append(b, '\n')
}

// vsCodekortexEngramOverlayJSON is the VS Code mcp.json overlay using the "servers" key.
// Uses --tools=agent per KortexEngram contract.
// VS Code uses a fixed "servers" key structure rather than mcpServers, so it
// is kept as a separate helper.
func vsCodekortexEngramOverlayJSON(cmd string) []byte {
	cfg := map[string]any{
		"servers": map[string]any{
			"kortex-engram": map[string]any{
				"command": cmd,
				"args":    []string{"mcp", "--tools=agent"},
			},
		},
	}
	b, _ := json.MarshalIndent(cfg, "", "  ")
	return append(b, '\n')
}

func Inject(homeDir string, adapter agents.Adapter) (InjectionResult, error) {
	if !adapter.SupportsMCP() {
		return InjectionResult{}, nil
	}

	files := make([]string, 0, 2)
	changed := false

	// 1. Write MCP server config using the adapter's strategy.
	switch adapter.MCPStrategy() {
	case model.StrategySeparateMCPFiles:
		// KortexEngram v1.10.3+ writes an absolute path for the command field when
		// `KortexEngram setup <agent>` is invoked. kortex's Inject() runs after
		// KortexEngram setup, so we must preserve any absolute command path already
		// present instead of silently overwriting it with the relative "kortex-engram".
		// See: https://github.com/fortissolucoescontato-bit/kortex/issues (KortexEngram absolute path regression)
		mcpPath := adapter.MCPConfigPath(homeDir, "kortex-engram")
		cmd := stablekortexEngramCommandForMergedConfig(mcpPath, adapter.Agent())
		content := buildSeparateMCPContent(mcpPath, kortexEngramServerJSONWithCmd(cmd))
		mcpWrite, err := filemerge.WriteFileAtomic(mcpPath, content, 0o644)
		if err != nil {
			return InjectionResult{}, err
		}
		changed = changed || mcpWrite.Changed
		files = append(files, mcpPath)

	case model.StrategyMergeIntoSettings:
		settingsPath := adapter.SettingsPath(homeDir)
		if settingsPath == "" {
			break
		}
		overlay := kortexEngramOverlayJSON(adapter.Agent(), stablekortexEngramCommandForMergedConfig(settingsPath, adapter.Agent()))
		settingsWrite, err := mergeJSONFile(settingsPath, overlay)
		if err != nil {
			return InjectionResult{}, err
		}
		changed = changed || settingsWrite.Changed
		files = append(files, settingsPath)

	case model.StrategyMCPConfigFile:
		mcpPath := adapter.MCPConfigPath(homeDir, "kortex-engram")
		if mcpPath == "" {
			break
		}
		var overlay []byte
		if adapter.Agent() == model.AgentVSCodeCopilot {
			overlay = vsCodekortexEngramOverlayJSON(stablekortexEngramCommandForMergedConfig(mcpPath, adapter.Agent()))
		} else {
			overlay = kortexEngramOverlayJSON(adapter.Agent(), stablekortexEngramCommandForMergedConfig(mcpPath, adapter.Agent()))
		}

		mcpWrite, err := mergeJSONFile(mcpPath, overlay)
		if err != nil {
			return InjectionResult{}, err
		}
		changed = changed || mcpWrite.Changed
		files = append(files, mcpPath)

		if adapter.Agent() == model.AgentAntigravity {
			settingsWrite, err := ensureAntigravitySettings(homeDir, adapter)
			if err != nil {
				return InjectionResult{}, err
			}
			changed = changed || settingsWrite.Changed
			if settingsWrite.Path != "" {
				files = append(files, settingsWrite.Path)
			}
		}

	case model.StrategyTOMLFile:
		// Codex: upsert [mcp_servers.KortexEngram] block and instruction-file keys
		// in ~/.codex/config.toml, then write instruction files.
		// All TOML mutations are composed in a single pass before writing to
		// ensure idempotency (no intermediate states that differ on re-run).
		configPath := adapter.MCPConfigPath(homeDir, "kortex-engram")
		if configPath == "" {
			break
		}

		// Determine instruction file paths before mutating the config.
		instructionsPath, compactPath, instrErr := writeCodexInstructionFiles(homeDir)
		if instrErr != nil {
			return InjectionResult{}, instrErr
		}

		// Read existing config and apply all mutations in one pass.
		existing, err := readFileOrEmpty(configPath)
		if err != nil {
			return InjectionResult{}, err
		}
		KortexEngramCmd := stablekortexEngramCommandForMergedConfig(configPath, adapter.Agent())
		withMCP := filemerge.UpsertCodexKortexEngramBlock(existing, KortexEngramCmd)
		withInstr := filemerge.UpsertTopLevelTOMLString(withMCP, "model_instructions_file", instructionsPath)
		withCompact := filemerge.UpsertTopLevelTOMLString(withInstr, "experimental_compact_prompt_file", compactPath)

		tomlWrite, err := filemerge.WriteFileAtomic(configPath, []byte(withCompact), 0o644)
		if err != nil {
			return InjectionResult{}, err
		}
		changed = changed || tomlWrite.Changed
		files = append(files, configPath)
	}

	// 2. Inject KortexEngram memory protocol into system prompt (if supported).
	if adapter.SupportsSystemPrompt() {
		switch adapter.SystemPromptStrategy() {
		case model.StrategyMarkdownSections:
			promptPath := adapter.SystemPromptFile(homeDir)
			protocolContent := assets.MustRead("claude/KortexEngram-protocol.md")

			existing, err := readFileOrEmpty(promptPath)
			if err != nil {
				return InjectionResult{}, err
			}

			updated := filemerge.InjectMarkdownSection(existing, "KortexEngram-protocol", protocolContent)

			mdWrite, err := filemerge.WriteFileAtomic(promptPath, []byte(updated), 0o644)
			if err != nil {
				return InjectionResult{}, err
			}
			changed = changed || mdWrite.Changed
			files = append(files, promptPath)

		case model.StrategyJinjaModules:
			// Ensure the base template exists for Jinja-based agents.
			if bs, ok := adapter.(bootstrapper); ok {
				if err := bs.BootstrapTemplate(homeDir); err != nil {
					return InjectionResult{}, fmt.Errorf("bootstrap template: %w", err)
				}
			}

			// Write the KortexEngram protocol as a standalone Jinja include module.
			// The static KIMI.md template references it via {% include "KortexEngram-protocol.md" %}.
			configDir := adapter.GlobalConfigDir(homeDir)
			protocolContent := assets.MustRead("claude/KortexEngram-protocol.md")
			modulePath := filepath.Join(configDir, "KortexEngram-protocol.md")
			mdWrite, err := filemerge.WriteFileAtomic(modulePath, []byte(protocolContent), 0o644)
			if err != nil {
				return InjectionResult{}, err
			}
			changed = changed || mdWrite.Changed
			files = append(files, modulePath)

		default:
			promptPath := adapter.SystemPromptFile(homeDir)
			protocolContent := assets.MustRead("claude/KortexEngram-protocol.md")

			existing, err := readFileOrEmpty(promptPath)
			if err != nil {
				return InjectionResult{}, err
			}

			updated := filemerge.InjectMarkdownSection(existing, "KortexEngram-protocol", protocolContent)

			mdWrite, err := filemerge.WriteFileAtomic(promptPath, []byte(updated), 0o644)
			if err != nil {
				return InjectionResult{}, err
			}
			changed = changed || mdWrite.Changed
			files = append(files, promptPath)
		}
	}

	// 3. Write shared convention files (if the agent supports skills).
	if adapter.SupportsSkills() {
		skillDir := adapter.SkillsDir(homeDir)
		if skillDir != "" {
			sharedFiles := []string{
				"KortexEngram-convention.md",
				"kortex-convention.md",
			}

			for _, fileName := range sharedFiles {
				assetPath := "skills/_shared/" + fileName
				content, readErr := assets.Read(assetPath)
				if readErr != nil {
					continue
				}
				if len(content) == 0 {
					continue
				}

				path := filepath.Join(skillDir, "_shared", fileName)
				writeResult, err := filemerge.WriteFileAtomic(path, []byte(content), 0o644)
				if err != nil {
					return InjectionResult{}, err
				}
				changed = changed || writeResult.Changed
				files = append(files, path)
			}
		}
	}

	return InjectionResult{Changed: changed, Files: files}, nil
}

type settingsBootstrapResult struct {
	Changed bool
	Path    string
}

func ensureAntigravitySettings(homeDir string, adapter agents.Adapter) (settingsBootstrapResult, error) {
	settingsPath := adapter.SettingsPath(homeDir)
	if settingsPath == "" {
		return settingsBootstrapResult{}, nil
	}

	if _, err := os.Stat(settingsPath); err == nil {
		return settingsBootstrapResult{Path: settingsPath}, nil
	} else if !os.IsNotExist(err) {
		return settingsBootstrapResult{}, fmt.Errorf("stat antigravity settings %q: %w", settingsPath, err)
	}

	sourcePath := filepath.Join(homeDir, ".gemini", "settings.json")
	content, err := os.ReadFile(sourcePath)
	if err != nil {
		if !os.IsNotExist(err) {
			return settingsBootstrapResult{}, fmt.Errorf("read gemini settings %q: %w", sourcePath, err)
		}
		content = []byte("{}")
	}

	writeResult, err := filemerge.WriteFileAtomic(settingsPath, content, 0o644)
	if err != nil {
		return settingsBootstrapResult{}, err
	}

	return settingsBootstrapResult{Changed: writeResult.Changed, Path: settingsPath}, nil
}

// writeCodexInstructionFiles writes the KortexEngram memory protocol and compact prompt
// files to ~/.codex/ and returns their paths.
func writeCodexInstructionFiles(homeDir string) (instructionsPath, compactPath string, err error) {
	codexDir := homeDir + "/.codex"
	instructionsPath = codexDir + "/KortexEngram-instructions.md"
	compactPath = codexDir + "/KortexEngram-compact-prompt.md"

	instrContent := assets.MustRead("codex/KortexEngram-instructions.md")
	instrWrite, err := filemerge.WriteFileAtomic(instructionsPath, []byte(instrContent), 0o644)
	if err != nil {
		return "", "", fmt.Errorf("write codex KortexEngram-instructions.md: %w", err)
	}
	_ = instrWrite

	compactContent := assets.MustRead("codex/KortexEngram-compact-prompt.md")
	compactWrite, err := filemerge.WriteFileAtomic(compactPath, []byte(compactContent), 0o644)
	if err != nil {
		return "", "", fmt.Errorf("write codex KortexEngram-compact-prompt.md: %w", err)
	}
	_ = compactWrite

	return instructionsPath, compactPath, nil
}

func mergeJSONFile(path string, overlay []byte) (filemerge.WriteResult, error) {
	baseJSON, err := osReadFile(path)
	if err != nil {
		return filemerge.WriteResult{}, err
	}

	merged, err := filemerge.MergeJSONObjects(baseJSON, overlay)
	if err != nil {
		return filemerge.WriteResult{}, err
	}

	return filemerge.WriteFileAtomic(path, merged, 0o644)
}

var osReadFile = func(path string) ([]byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read json file %q: %w", path, err)
	}

	return content, nil
}

func readFileOrEmpty(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("read file %q: %w", path, err)
	}
	return string(data), nil
}

func stablekortexEngramCommandForMergedConfig(path string, agentID model.AgentID) string {
	raw, err := osReadFile(path)
	if err == nil {
		if cmd, ok := existingMergedKortexEngramCommand(raw, agentID); ok {
			return cmd
		}
	}

	if isStandardAgent(agentID) {
		return "kortex-engram"
	}

	cmd, _ := resolveKortexEngramCommand()
	return cmd
}

func existingMergedKortexEngramCommand(raw []byte, agentID model.AgentID) (string, bool) {
	if len(raw) == 0 {
		return "", false
	}

	normalized, err := filemerge.MergeJSONObjects(raw, []byte("{}"))
	if err != nil {
		return "", false
	}

	var root map[string]any
	if err := json.Unmarshal(normalized, &root); err != nil {
		return "", false
	}

	var server any
	switch agentID {
	case model.AgentOpenCode:
		mcp, ok := root["mcp"].(map[string]any)
		if !ok {
			return "", false
		}
		server = mcp["kortex-engram"]
		if server == nil {
			server = mcp["kortex-engram"] // Fallback
		}
	case model.AgentVSCodeCopilot:
		servers, ok := root["servers"].(map[string]any)
		if !ok {
			return "", false
		}
		server = servers["kortex-engram"]
		if server == nil {
			server = servers["kortex-engram"] // Fallback
		}
	default:
		mcpServers, ok := root["mcpServers"].(map[string]any)
		if !ok {
			return "", false
		}
		server = mcpServers["kortex-engram"]
		if server == nil {
			server = mcpServers["kortex-engram"] // Fallback
		}
	}

	serverMap, ok := server.(map[string]any)
	if !ok {
		return "", false
	}

	return executableFromCommandValue(serverMap["command"])
}

func executableFromCommandValue(command any) (string, bool) {
	switch value := command.(type) {
	case string:
		if value == "" {
			return "", false
		}
		return value, true
	case []any:
		if len(value) == 0 {
			return "", false
		}
		first, ok := value[0].(string)
		if !ok || first == "" {
			return "", false
		}
		return first, true
	default:
		return "", false
	}
}

func isStandardAgent(id model.AgentID) bool {
	switch id {
	case model.AgentOpenCode, model.AgentQwenCode, model.AgentCodex, model.AgentGeminiCLI, model.AgentAntigravity, model.AgentClaudeCode:
		return true
	default:
		return false
	}
}

// buildSeparateMCPContent returns the content to write to the MCP server JSON
// file for agents that use the StrategySeparateMCPFiles strategy (e.g. Claude
// Code).
//
// KortexEngram v1.10.3+ writes an absolute command path when `KortexEngram setup` is run.
// kortex runs Inject() after setup, so we must not overwrite that absolute
// path with the relative "kortex-engram" string from defaultKortexEngramServerJSON.
//
// Logic:
//   - If the file does not exist yet, return defaultContent unchanged.
//   - If the file exists but cannot be parsed as JSON, return defaultContent.
//   - If the parsed JSON has a "command" value that is an absolute path to the
//     KortexEngram binary, rebuild the config using that command and the canonical
//     args (["mcp", "--tools=agent"]) so that the absolute path is preserved
//     and the correct flags are always present.
//   - Otherwise (relative command or other value), return defaultContent.
func buildSeparateMCPContent(mcpPath string, defaultContent []byte) []byte {
	raw, err := os.ReadFile(mcpPath)
	if err != nil {
		// File does not exist or is not readable — use the default.
		return defaultContent
	}

	var existing map[string]any
	if err := json.Unmarshal(raw, &existing); err != nil {
		// Malformed JSON — use the default.
		return defaultContent
	}

	cmd, ok := executableFromCommandValue(existing["command"])
	if !ok || !iskortexEngramCommand(cmd) {
		// No command, or not an KortexEngram command — use the default.
		return defaultContent
	}

	// Rebuild with the preserved command and the canonical args (["mcp", "--tools=agent"]).
	rebuilt := map[string]any{
		"command": cmd,
		"args":    []string{"mcp", "--tools=agent"},
	}
	encoded, err := json.MarshalIndent(rebuilt, "", "  ")
	if err != nil {
		// Should be impossible with a plain map — use the default as fallback.
		return defaultContent
	}
	return append(encoded, '\n')
}

// iskortexEngramCommand reports whether cmd is either a relative "kortex-engram" command
// or an absolute path pointing to an KortexEngram binary.
func iskortexEngramCommand(cmd string) bool {
	if cmd == "" {
		return false
	}
	base := filepath.Base(cmd)
	if runtime.GOOS == "windows" {
		return strings.EqualFold(base, "kortexengram.exe") || strings.EqualFold(base, "kortex-engram") ||
			strings.EqualFold(base, "kortex.exe") || strings.EqualFold(base, "kortex") ||
			strings.EqualFold(base, "kortexengram.exe") || strings.EqualFold(base, "kortex-engram")
	}
	return base == "kortex-engram" || base == "kortex" || base == "kortex-engram"
}

// isAbsolutekortexEngramPath reports whether path is an absolute filesystem path
// that points to an KortexEngram binary.
func isAbsolutekortexEngramPath(path string) bool {
	return filepath.IsAbs(path) && iskortexEngramCommand(path)
}
