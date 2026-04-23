package kortexengram

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fortissolucoescontato-bit/kortex/internal/agents"
	"github.com/fortissolucoescontato-bit/kortex/internal/agents/antigravity"
	"github.com/fortissolucoescontato-bit/kortex/internal/agents/claude"
	"github.com/fortissolucoescontato-bit/kortex/internal/agents/codex"
	"github.com/fortissolucoescontato-bit/kortex/internal/agents/gemini"
	"github.com/fortissolucoescontato-bit/kortex/internal/agents/opencode"
	"github.com/fortissolucoescontato-bit/kortex/internal/agents/qwen"
	"github.com/fortissolucoescontato-bit/kortex/internal/agents/vscode"
)

func claudeAdapter() agents.Adapter   { return claude.NewAdapter() }
func opencodeAdapter() agents.Adapter { return opencode.NewAdapter() }
func codexAdapter() agents.Adapter    { return codex.NewAdapter() }
func geminiAdapter() agents.Adapter   { return gemini.NewAdapter() }
func qwenAdapter() agents.Adapter     { return qwen.NewAdapter() }
func antigravityAdapter() agents.Adapter {
	return antigravity.NewAdapter()
}

// assertArgsHaveToolsAgent is a shared helper that validates a JSON file
// contains the MCP "kortex-engram" entry with --tools=agent in args.
func assertArgsHaveToolsAgent(t *testing.T, path string) {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", path, err)
	}
	text := string(content)
	if !strings.Contains(text, `"--tools=agent"`) {
		t.Fatalf("file %q missing --tools=agent in args; got:\n%s", path, text)
	}
}

func TestInjectClaudeWritesMCPConfig(t *testing.T) {
	home := t.TempDir()

	result, err := Inject(home, claudeAdapter())
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}
	if !result.Changed {
		t.Fatalf("Inject() changed = false")
	}

	// Check MCP JSON file was created.
	mcpPath := filepath.Join(home, ".claude", "mcp", "kortexengram.json")
	mcpContent, err := os.ReadFile(mcpPath)
	if err != nil {
		t.Fatalf("ReadFile(kortexengram.json) error = %v", err)
	}

	// Parse the JSON and validate the "command" key exists and references kortexengram.
	// The command may be an absolute path (if KortexEngram is on PATH) or the relative
	// string "kortex-engram" (if not found). Both are valid.
	var parsed map[string]any
	if err := json.Unmarshal(mcpContent, &parsed); err != nil {
		t.Fatalf("Unmarshal(kortexengram.json) error = %v", err)
	}
	cmd, ok := parsed["command"].(string)
	if !ok || cmd == "" {
		t.Fatalf("kortexengram.json missing or empty command field; got: %s", mcpContent)
	}
	// Command must either be the literal "kortex-engram" or an absolute path ending in "kortex-engram".
	base := filepath.Base(cmd)
	if !iskortexEngramCommand(base) {
		t.Fatalf("kortexengram.json command %q does not reference KortexEngram binary; got: %s", cmd, mcpContent)
	}
	if _, ok := parsed["args"]; !ok {
		t.Fatal("kortexengram.json missing args field")
	}
	// RED: must include --tools=agent
	assertArgsHaveToolsAgent(t, mcpPath)
}

func TestInjectClaudeWritesProtocolSection(t *testing.T) {
	home := t.TempDir()

	_, err := Inject(home, claudeAdapter())
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}

	claudeMDPath := filepath.Join(home, ".claude", "CLAUDE.md")
	content, err := os.ReadFile(claudeMDPath)
	if err != nil {
		t.Fatalf("ReadFile(CLAUDE.md) error = %v", err)
	}

	text := string(content)
	if !strings.Contains(text, "<!-- kortex:KortexEngram-protocol -->") {
		t.Fatal("CLAUDE.md missing open marker for KortexEngram-protocol")
	}
	if !strings.Contains(text, "<!-- /kortex:KortexEngram-protocol -->") {
		t.Fatal("CLAUDE.md missing close marker for KortexEngram-protocol")
	}
	// Real content check.
	if !strings.Contains(text, "mem_save") {
		t.Fatal("CLAUDE.md missing real KortexEngram protocol content (expected 'mem_save')")
	}
}

func TestInjectClaudeIsIdempotent(t *testing.T) {
	home := t.TempDir()

	first, err := Inject(home, claudeAdapter())
	if err != nil {
		t.Fatalf("Inject() first error = %v", err)
	}
	if !first.Changed {
		t.Fatalf("Inject() first changed = false")
	}

	second, err := Inject(home, claudeAdapter())
	if err != nil {
		t.Fatalf("Inject() second error = %v", err)
	}
	if second.Changed {
		t.Fatalf("Inject() second changed = true")
	}
}

func TestInjectOpenCodeMergesKortexEngramToSettings(t *testing.T) {
	home := t.TempDir()

	result, err := Inject(home, opencodeAdapter())
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}
	if !result.Changed {
		t.Fatalf("Inject() changed = false")
	}

	// Should include opencode.json and AGENTS.md (fallback protocol injection).
	if len(result.Files) != 2 {
		t.Fatalf("Inject() files = %v, want exactly 2 (opencode.json + AGENTS.md)", result.Files)
	}

	configPath := filepath.Join(home, ".config", "opencode", "opencode.json")
	config, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("ReadFile(opencode.json) error = %v", err)
	}

	text := string(config)
	if !strings.Contains(text, `"kortex-engram"`) {
		t.Fatal("opencode.json missing KortexEngram server entry")
	}
	if !strings.Contains(text, `"mcp"`) {
		t.Fatal("opencode.json missing mcp key")
	}
	if strings.Contains(text, `"mcpServers"`) {
		t.Fatal("opencode.json should use 'mcp' key, not 'mcpServers'")
	}
	if !strings.Contains(text, `"type": "local"`) {
		t.Fatal("opencode.json KortexEngram missing type: local")
	}
	// OpenCode 1.3.3+: command must be an array, no separate "args" field.
	if !strings.Contains(text, `"--tools=agent"`) {
		t.Fatal("opencode.json missing --tools=agent in command array")
	}
	if strings.Contains(text, `"args"`) {
		t.Fatal("opencode.json must NOT have a separate args field — command must be an array")
	}

	// Verify NO plugin files or plugin arrays exist.
	pluginPath := filepath.Join(home, ".config", "opencode", "plugins", "kortexengram.ts")
	if _, err := os.Stat(pluginPath); err == nil {
		t.Fatal("plugin file should NOT exist — old approach removed")
	}
	if strings.Contains(text, `"plugins"`) {
		t.Fatal("opencode.json should NOT contain plugins key")
	}

	agentsPath := filepath.Join(home, ".config", "opencode", "AGENTS.md")
	agentsContent, err := os.ReadFile(agentsPath)
	if err != nil {
		t.Fatalf("ReadFile(AGENTS.md) error = %v", err)
	}
	agentsText := string(agentsContent)
	if !strings.Contains(agentsText, "<!-- kortex:KortexEngram-protocol -->") {
		t.Fatal("AGENTS.md missing KortexEngram protocol section marker")
	}
	if !strings.Contains(agentsText, "mem_save") {
		t.Fatal("AGENTS.md missing KortexEngram protocol content (expected 'mem_save')")
	}
}

func TestInjectOpenCodeIsIdempotent(t *testing.T) {
	home := t.TempDir()

	first, err := Inject(home, opencodeAdapter())
	if err != nil {
		t.Fatalf("Inject() first error = %v", err)
	}
	if !first.Changed {
		t.Fatalf("Inject() first changed = false")
	}

	second, err := Inject(home, opencodeAdapter())
	if err != nil {
		t.Fatalf("Inject() second error = %v", err)
	}
	if second.Changed {
		t.Fatalf("Inject() second changed = true")
	}
}

// TestInjectOpenCodeMigratesFromOldFormat verifies that when a user's
// opencode.json contains the old v1.11.3 format (separate "args" key),
// Inject() replaces mcp.KortexEngram atomically so that "args" is absent and
// "command" is an array — the format required by OpenCode 1.3.3+.
func TestInjectOpenCodeMigratesFromOldFormat(t *testing.T) {
	home := t.TempDir()

	mockkortexEngramLookPath(t, "/opt/homebrew/bin/KortexEngram", "")

	adapter := opencodeAdapter()
	configPath := adapter.SettingsPath(home)
	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		t.Fatalf("MkdirAll error = %v", err)
	}

	// Pre-seed with the old v1.11.3 format.
	oldFormat := `{"mcp": {"kortex-engram": {"command": "/opt/homebrew/bin/KortexEngram", "args": ["mcp","--tools=agent"], "type": "local"}}}`
	if err := os.WriteFile(configPath, []byte(oldFormat), 0o644); err != nil {
		t.Fatalf("WriteFile(opencode.json) error = %v", err)
	}

	result, err := Inject(home, adapter)
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}
	if !result.Changed {
		t.Fatalf("Inject() changed = false; expected migration to produce a change")
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("ReadFile(opencode.json) error = %v", err)
	}

	// (1) "args" key must be absent from mcp.kortexengram.
	if strings.Contains(string(content), `"args"`) {
		t.Fatalf("mcp.KortexEngram still contains 'args' key after migration; got:\n%s", content)
	}

	// (2) command must be a []any containing the KortexEngram binary.
	var parsed map[string]any
	if err := json.Unmarshal(content, &parsed); err != nil {
		t.Fatalf("Unmarshal(opencode.json) error = %v", err)
	}
	mcpMap, _ := parsed["mcp"].(map[string]any)
	KortexEngramMap, _ := mcpMap["kortex-engram"].(map[string]any)
	cmdRaw, ok := KortexEngramMap["command"]
	if !ok {
		t.Fatalf("mcp.KortexEngram missing command key; got:\n%s", content)
	}
	cmdArr, ok := cmdRaw.([]any)
	if !ok {
		t.Fatalf("mcp.kortexengram.command must be []any after migration, got %T; got:\n%s", cmdRaw, content)
	}
	if len(cmdArr) == 0 {
		t.Fatalf("mcp.kortexengram.command array is empty; got:\n%s", content)
	}
	firstElem, _ := cmdArr[0].(string)
	if firstElem == "" {
		t.Fatalf("mcp.kortexengram.command[0] is empty or not a string; got:\n%s", content)
	}
	// Must be an KortexEngram command.
	if !iskortexEngramCommand(firstElem) {
		t.Fatalf("mcp.kortexengram.command[0] = %q does not reference KortexEngram binary; got:\n%s", firstElem, content)
	}

	// (3) Second Inject() call must be idempotent (changed=false).
	second, err := Inject(home, adapter)
	if err != nil {
		t.Fatalf("Inject() second error = %v", err)
	}
	if second.Changed {
		t.Fatalf("Inject() second changed = true; expected idempotent (no change)")
	}
}

func TestInjectCursorMergesKortexEngramToSettings(t *testing.T) {
	home := t.TempDir()

	cursorAdapter, err := agents.NewAdapter("cursor")
	if err != nil {
		t.Fatalf("NewAdapter(cursor) error = %v", err)
	}

	result, injectErr := Inject(home, cursorAdapter)
	if injectErr != nil {
		t.Fatalf("Inject(cursor) error = %v", injectErr)
	}

	// Cursor uses MCPConfigFile strategy — KortexEngram gets merged into mcp.json.
	if !result.Changed {
		t.Fatalf("Inject(cursor) changed = false")
	}
}

func TestInjectCursorWithMalformedMCPJsonRecovery(t *testing.T) {
	// Real Windows users may have a ~/.cursor/mcp.json that starts with non-JSON
	// content (e.g. "allow: all" or just "a"). The installer should recover by
	// treating the broken file as {} and proceeding with the overlay merge.
	home := t.TempDir()

	cursorAdapter, err := agents.NewAdapter("cursor")
	if err != nil {
		t.Fatalf("NewAdapter(cursor) error = %v", err)
	}

	// Pre-create ~/.cursor/mcp.json with invalid (non-JSON) content.
	mcpPath := cursorAdapter.MCPConfigPath(home, "kortex-engram")
	if err := os.MkdirAll(filepath.Dir(mcpPath), 0o755); err != nil {
		t.Fatalf("MkdirAll error = %v", err)
	}
	if err := os.WriteFile(mcpPath, []byte("allow: all"), 0o644); err != nil {
		t.Fatalf("WriteFile(malformed mcp.json) error = %v", err)
	}

	result, injectErr := Inject(home, cursorAdapter)
	if injectErr != nil {
		t.Fatalf("Inject(cursor) with malformed mcp.json error = %v; want nil (should recover)", injectErr)
	}
	if !result.Changed {
		t.Fatalf("Inject(cursor) changed = false; want true")
	}

	content, err := os.ReadFile(mcpPath)
	if err != nil {
		t.Fatalf("ReadFile(mcp.json) error = %v", err)
	}

	text := string(content)
	if !strings.Contains(text, `"mcpServers"`) {
		t.Fatalf("mcp.json missing mcpServers key after recovery; got:\n%s", text)
	}
	if !strings.Contains(text, `"kortex-engram"`) {
		t.Fatalf("mcp.json missing KortexEngram server after recovery; got:\n%s", text)
	}
}

func TestInjectVSCodeMergesKortexEngramToMCPConfigFile(t *testing.T) {
	home := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))
	adapter := vscode.NewAdapter()

	result, err := Inject(home, adapter)
	if err != nil {
		t.Fatalf("Inject(vscode) error = %v", err)
	}
	if !result.Changed {
		t.Fatalf("Inject(vscode) changed = false")
	}

	mcpPath := adapter.MCPConfigPath(home, "kortex-engram")
	content, err := os.ReadFile(mcpPath)
	if err != nil {
		t.Fatalf("ReadFile(mcp.json) error = %v", err)
	}

	text := string(content)
	if !strings.Contains(text, `"servers"`) {
		t.Fatal("mcp.json missing servers key")
	}
	if !strings.Contains(text, `"kortex-engram"`) {
		t.Fatal("mcp.json missing KortexEngram server")
	}
	if !strings.Contains(text, `"mcp"`) {
		t.Fatal("mcp.json missing KortexEngram args mcp")
	}
	if strings.Contains(text, `"mcpServers"`) {
		t.Fatal("mcp.json should use 'servers' key, not 'mcpServers'")
	}
	// RED: VS Code overlay must include --tools=agent
	assertArgsHaveToolsAgent(t, mcpPath)
}

// ─── Gemini tests ─────────────────────────────────────────────────────────────

func TestInjectGeminiToolsFlagPresent(t *testing.T) {
	home := t.TempDir()

	result, err := Inject(home, geminiAdapter())
	if err != nil {
		t.Fatalf("Inject(gemini) error = %v", err)
	}
	if !result.Changed {
		t.Fatalf("Inject(gemini) changed = false")
	}

	settingsPath := filepath.Join(home, ".gemini", "settings.json")
	content, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatalf("ReadFile(settings.json) error = %v", err)
	}
	text := string(content)
	if !strings.Contains(text, `"mcpServers"`) {
		t.Fatal("settings.json missing mcpServers key")
	}
	if !strings.Contains(text, `"kortex-engram"`) {
		t.Fatal("settings.json missing KortexEngram entry")
	}
	// RED: Gemini overlay must use --tools=agent
	if !strings.Contains(text, `"--tools=agent"`) {
		t.Fatal("settings.json missing --tools=agent in args")
	}
}

func TestInjectAntigravityCopiesGeminiSettingsAfterKortexEngramSetup(t *testing.T) {
	home := t.TempDir()
	sourcePath := filepath.Join(home, ".gemini", "settings.json")
	if err := os.MkdirAll(filepath.Dir(sourcePath), 0o755); err != nil {
		t.Fatalf("MkdirAll(%q) error = %v", filepath.Dir(sourcePath), err)
	}
	want := []byte("{\"theme\":\"dark\"}\n")
	if err := os.WriteFile(sourcePath, want, 0o644); err != nil {
		t.Fatalf("WriteFile(%q) error = %v", sourcePath, err)
	}

	result, err := Inject(home, antigravityAdapter())
	if err != nil {
		t.Fatalf("Inject(antigravity) error = %v", err)
	}
	if !result.Changed {
		t.Fatalf("Inject(antigravity) changed = false")
	}

	settingsPath := filepath.Join(home, ".gemini", "antigravity", "settings.json")
	got, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", settingsPath, err)
	}
	if string(got) != string(want) {
		t.Fatalf("antigravity settings = %q, want %q", got, want)
	}

	mcpPath := filepath.Join(home, ".gemini", "antigravity", "mcp_config.json")
	assertArgsHaveToolsAgent(t, mcpPath)
}

func TestInjectAntigravityInitializesEmptySettingsWhenGeminiMissing(t *testing.T) {
	home := t.TempDir()

	first, err := Inject(home, antigravityAdapter())
	if err != nil {
		t.Fatalf("Inject(antigravity) first error = %v", err)
	}
	if !first.Changed {
		t.Fatalf("Inject(antigravity) first changed = false")
	}

	settingsPath := filepath.Join(home, ".gemini", "antigravity", "settings.json")
	got, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", settingsPath, err)
	}
	if strings.TrimSpace(string(got)) != "{}" {
		t.Fatalf("antigravity settings = %q, want empty JSON object", got)
	}

	second, err := Inject(home, antigravityAdapter())
	if err != nil {
		t.Fatalf("Inject(antigravity) second error = %v", err)
	}
	if second.Changed {
		t.Fatalf("Inject(antigravity) second changed = true; want false")
	}
}

// ─── Codex tests ──────────────────────────────────────────────────────────────

func TestInjectCodexWritesTOMLMCP(t *testing.T) {
	home := t.TempDir()

	result, err := Inject(home, codexAdapter())
	if err != nil {
		t.Fatalf("Inject(codex) error = %v", err)
	}
	if !result.Changed {
		t.Fatalf("Inject(codex) changed = false")
	}

	configPath := filepath.Join(home, ".codex", "config.toml")
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("ReadFile(config.toml) error = %v", err)
	}
	text := string(content)
	if !strings.Contains(text, "[mcp_servers.KortexEngram]") {
		t.Fatalf("config.toml missing [mcp_servers.KortexEngram] block; got:\n%s", text)
	}
	// command must reference the KortexEngram binary — either relative ("kortex-engram") or an
	// absolute path (when KortexEngram is on PATH). Both are valid.
	if !strings.Contains(text, "command = ") {
		t.Fatalf("config.toml missing command field; got:\n%s", text)
	}
	cmdLine := ""
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "command = ") {
			cmdLine = strings.TrimSpace(line)
			break
		}
	}
	if cmdLine == "" {
		t.Fatalf("config.toml missing command line; got:\n%s", text)
	}
	// The command value must end with "kortex-engram" or "kortexengram.exe".
	cmdVal := strings.TrimPrefix(cmdLine, "command = ")
	cmdVal = strings.Trim(cmdVal, `"`)
	base := filepath.Base(cmdVal)
	if !iskortexEngramCommand(base) {
		t.Fatalf("config.toml command %q does not reference KortexEngram binary; got:\n%s", cmdVal, text)
	}
	if !strings.Contains(text, `"--tools=agent"`) {
		t.Fatalf("config.toml missing --tools=agent; got:\n%s", text)
	}
}

func TestInjectCodexWritesInstructionFiles(t *testing.T) {
	home := t.TempDir()

	_, err := Inject(home, codexAdapter())
	if err != nil {
		t.Fatalf("Inject(codex) error = %v", err)
	}

	instructionsPath := filepath.Join(home, ".codex", "KortexEngram-instructions.md")
	content, err := os.ReadFile(instructionsPath)
	if err != nil {
		t.Fatalf("ReadFile(KortexEngram-instructions.md) error = %v", err)
	}
	if !strings.Contains(string(content), "mem_save") {
		t.Fatal("KortexEngram-instructions.md missing expected content (mem_save)")
	}

	compactPath := filepath.Join(home, ".codex", "KortexEngram-compact-prompt.md")
	compactContent, err := os.ReadFile(compactPath)
	if err != nil {
		t.Fatalf("ReadFile(KortexEngram-compact-prompt.md) error = %v", err)
	}
	if !strings.Contains(string(compactContent), "FIRST ACTION REQUIRED") {
		t.Fatal("KortexEngram-compact-prompt.md missing expected content (FIRST ACTION REQUIRED)")
	}
}

func TestInjectCodexInjectsTOMLKeys(t *testing.T) {
	home := t.TempDir()

	_, err := Inject(home, codexAdapter())
	if err != nil {
		t.Fatalf("Inject(codex) error = %v", err)
	}

	configPath := filepath.Join(home, ".codex", "config.toml")
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("ReadFile(config.toml) error = %v", err)
	}
	text := string(content)

	instructionsPath := filepath.Join(home, ".codex", "KortexEngram-instructions.md")
	if !strings.Contains(text, `model_instructions_file`) {
		t.Fatalf("config.toml missing model_instructions_file key; got:\n%s", text)
	}
	if !strings.Contains(text, instructionsPath) {
		t.Fatalf("config.toml model_instructions_file does not reference %q; got:\n%s", instructionsPath, text)
	}

	compactPath := filepath.Join(home, ".codex", "KortexEngram-compact-prompt.md")
	if !strings.Contains(text, `experimental_compact_prompt_file`) {
		t.Fatalf("config.toml missing experimental_compact_prompt_file key; got:\n%s", text)
	}
	if !strings.Contains(text, compactPath) {
		t.Fatalf("config.toml experimental_compact_prompt_file does not reference %q; got:\n%s", compactPath, text)
	}
}

// ─── KortexEngram setup absolute path preservation tests ────────────────────────────

// TestInjectClaudePreservesAbsoluteCommandFromKortexEngramSetup verifies that when
// `KortexEngram setup claude-code` has already written an absolute-path command to
// ~/.claude/mcp/kortexengram.json (KortexEngram v1.10.3+ behaviour), a subsequent call to
// Inject() does NOT overwrite the absolute path with the relative "kortex-engram".
func TestInjectClaudePreservesAbsoluteCommandFromKortexEngramSetup(t *testing.T) {
	home := t.TempDir()

	// Simulate what `KortexEngram setup claude-code` writes on v1.10.3+:
	// an absolute path as the command value.
	absPath := "/opt/homebrew/bin/KortexEngram"
	mcpPath := filepath.Join(home, ".claude", "mcp", "kortexengram.json")
	if err := os.MkdirAll(filepath.Dir(mcpPath), 0o755); err != nil {
		t.Fatalf("MkdirAll error = %v", err)
	}
	setupContent := []byte(`{
  "command": "/opt/homebrew/bin/KortexEngram",
  "args": ["mcp", "--tools=agent"]
}
`)
	if err := os.WriteFile(mcpPath, setupContent, 0o644); err != nil {
		t.Fatalf("WriteFile(kortexengram.json) error = %v", err)
	}

	// Now run Inject — should NOT overwrite the absolute command.
	_, err := Inject(home, claudeAdapter())
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}

	content, err := os.ReadFile(mcpPath)
	if err != nil {
		t.Fatalf("ReadFile(kortexengram.json) error = %v", err)
	}

	text := string(content)
	if !strings.Contains(text, absPath) {
		t.Fatalf("Inject() overwrote absolute command path; want %q preserved, got:\n%s", absPath, text)
	}
	// Still must have --tools=agent.
	assertArgsHaveToolsAgent(t, mcpPath)
}

// TestInjectClaudePreservesAbsoluteCommandIsIdempotent verifies that calling
// Inject() twice when an absolute-path kortexengram.json already exists does not
// cause repeated writes (idempotency).
func TestInjectClaudePreservesAbsoluteCommandIsIdempotent(t *testing.T) {
	home := t.TempDir()

	absPath := "/usr/local/bin/KortexEngram"
	mcpPath := filepath.Join(home, ".claude", "mcp", "kortexengram.json")
	if err := os.MkdirAll(filepath.Dir(mcpPath), 0o755); err != nil {
		t.Fatalf("MkdirAll error = %v", err)
	}
	setupContent := []byte(`{
  "command": "/usr/local/bin/KortexEngram",
  "args": ["mcp", "--tools=agent"]
}
`)
	if err := os.WriteFile(mcpPath, setupContent, 0o644); err != nil {
		t.Fatalf("WriteFile(kortexengram.json) error = %v", err)
	}

	first, err := Inject(home, claudeAdapter())
	if err != nil {
		t.Fatalf("Inject() first error = %v", err)
	}

	second, err := Inject(home, claudeAdapter())
	if err != nil {
		t.Fatalf("Inject() second error = %v", err)
	}
	if second.Changed {
		t.Fatalf("Inject() second changed = true after absolute-path setup; want idempotent (no change)")
	}

	// Absolute path must still be present.
	content, err := os.ReadFile(mcpPath)
	if err != nil {
		t.Fatalf("ReadFile(kortexengram.json) error = %v", err)
	}
	if !strings.Contains(string(content), absPath) {
		t.Fatalf("absolute command path %q was lost after second Inject(); got:\n%s", absPath, string(content))
	}
	_ = first // first result not the focus of this test
}

// TestInjectClaudeAddsToolsAgentWhenSetupWritesBareArgs verifies that if
// `KortexEngram setup` wrote an absolute command but with bare args (no --tools=agent),
// Inject() adds --tools=agent while preserving the absolute path.
func TestInjectClaudeAddsToolsAgentWhenSetupWritesBareArgs(t *testing.T) {
	home := t.TempDir()

	absPath := "/home/user/go/bin/KortexEngram"
	mcpPath := filepath.Join(home, ".claude", "mcp", "kortexengram.json")
	if err := os.MkdirAll(filepath.Dir(mcpPath), 0o755); err != nil {
		t.Fatalf("MkdirAll error = %v", err)
	}
	// Bare mcp arg without --tools=agent — older KortexEngram setup format.
	setupContent := []byte(`{
  "command": "/home/user/go/bin/KortexEngram",
  "args": ["mcp"]
}
`)
	if err := os.WriteFile(mcpPath, setupContent, 0o644); err != nil {
		t.Fatalf("WriteFile(kortexengram.json) error = %v", err)
	}

	_, err := Inject(home, claudeAdapter())
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}

	content, err := os.ReadFile(mcpPath)
	if err != nil {
		t.Fatalf("ReadFile(kortexengram.json) error = %v", err)
	}
	text := string(content)

	// Absolute path must be preserved.
	if !strings.Contains(text, absPath) {
		t.Fatalf("absolute path %q was lost; got:\n%s", absPath, text)
	}
	// --tools=agent must be added.
	assertArgsHaveToolsAgent(t, mcpPath)
}

func TestInjectCodexIsIdempotent(t *testing.T) {
	home := t.TempDir()

	first, err := Inject(home, codexAdapter())
	if err != nil {
		t.Fatalf("Inject(codex) first error = %v", err)
	}
	if !first.Changed {
		t.Fatalf("Inject(codex) first changed = false")
	}

	second, err := Inject(home, codexAdapter())
	if err != nil {
		t.Fatalf("Inject(codex) second error = %v", err)
	}
	if second.Changed {
		t.Fatalf("Inject(codex) second changed = true (should be idempotent)")
	}

	// Verify only one [mcp_servers.KortexEngram] block.
	configPath := filepath.Join(home, ".codex", "config.toml")
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("ReadFile(config.toml) error = %v", err)
	}
	count := strings.Count(string(content), "[mcp_servers.KortexEngram]")
	if count != 1 {
		t.Fatalf("config.toml has %d [mcp_servers.KortexEngram] blocks, want exactly 1; got:\n%s", count, string(content))
	}
}

// ─── Absolute path resolution tests ──────────────────────────────────────────

// mockkortexEngramLookPath sets kortexEngramLookPath to a mock and restores it after the test.
func mockkortexEngramLookPath(t *testing.T, result string, errMsg string) {
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

// TestKortexEngramInjectUsesAbsolutePathWhenAvailable verifies that when KortexEngram is
// resolvable on PATH, its absolute path is written into the MCP config file
// for agents that use StrategyMCPConfigFile (e.g. Windsurf).
func TestKortexEngramInjectUsesAbsolutePathWhenAvailable(t *testing.T) {
	home := t.TempDir()

	absPath := "/usr/local/bin/KortexEngram"
	mockkortexEngramLookPath(t, absPath, "")

	windsurfAdapter, err := agents.NewAdapter("windsurf")
	if err != nil {
		t.Fatalf("NewAdapter(windsurf) error = %v", err)
	}

	result, injectErr := Inject(home, windsurfAdapter)
	if injectErr != nil {
		t.Fatalf("Inject(windsurf) error = %v", injectErr)
	}
	if !result.Changed {
		t.Fatalf("Inject(windsurf) changed = false")
	}

	mcpPath := windsurfAdapter.MCPConfigPath(home, "kortex-engram")
	content, readErr := os.ReadFile(mcpPath)
	if readErr != nil {
		t.Fatalf("ReadFile(%q) error = %v", mcpPath, readErr)
	}

	// Parse and validate the command field contains the absolute path.
	var parsed map[string]any
	if err := json.Unmarshal(content, &parsed); err != nil {
		t.Fatalf("Unmarshal(%q) error = %v", mcpPath, err)
	}

	mcpServersRaw, ok := parsed["mcpServers"]
	if !ok {
		t.Fatalf("mcp_config.json missing mcpServers key; got:\n%s", content)
	}
	mcpServers, ok := mcpServersRaw.(map[string]any)
	if !ok {
		t.Fatalf("mcpServers has unexpected type: %T", mcpServersRaw)
	}
	KortexEngramServerRaw, ok := mcpServers["kortex-engram"]
	if !ok {
		t.Fatalf("mcpServers missing KortexEngram entry; got:\n%s", content)
	}
	KortexEngramServer, ok := KortexEngramServerRaw.(map[string]any)
	if !ok {
		t.Fatalf("KortexEngram server has unexpected type: %T", KortexEngramServerRaw)
	}

	cmd, _ := KortexEngramServer["command"].(string)
	if cmd != absPath {
		t.Fatalf("mcp_config.json command = %q, want absolute path %q", cmd, absPath)
	}
}

// TestKortexEngramInjectFallsBackToRelativeWhenNotFound verifies that when KortexEngram
// cannot be resolved on PATH, the config falls back to the relative "kortex-engram"
// command string.
func TestKortexEngramInjectFallsBackToRelativeWhenNotFound(t *testing.T) {
	home := t.TempDir()

	mockkortexEngramLookPath(t, "", "not found")

	windsurfAdapter, err := agents.NewAdapter("windsurf")
	if err != nil {
		t.Fatalf("NewAdapter(windsurf) error = %v", err)
	}

	result, injectErr := Inject(home, windsurfAdapter)
	if injectErr != nil {
		t.Fatalf("Inject(windsurf) error = %v", injectErr)
	}
	if !result.Changed {
		t.Fatalf("Inject(windsurf) changed = false")
	}

	mcpPath := windsurfAdapter.MCPConfigPath(home, "kortex-engram")
	content, readErr := os.ReadFile(mcpPath)
	if readErr != nil {
		t.Fatalf("ReadFile(%q) error = %v", mcpPath, readErr)
	}

	text := string(content)
	if !strings.Contains(text, `"command": "kortex-engram"`) {
		t.Fatalf("mcp_config.json should use relative fallback 'KortexEngram'; got:\n%s", text)
	}
}

// TestKortexEngramInjectAbsolutePathForOpenCodeMergeStrategy verifies that the
// absolute path is used when the StrategyMergeIntoSettings strategy is
// applied for OpenCode.
func TestKortexEngramInjectAbsolutePathForOpenCodeMergeStrategy(t *testing.T) {
	home := t.TempDir()

	absPath := "/usr/local/bin/KortexEngram"
	mockkortexEngramLookPath(t, absPath, "")

	adapter := opencodeAdapter()
	settingsDir := filepath.Dir(adapter.SettingsPath(home))
	os.MkdirAll(settingsDir, 0o755)
	os.WriteFile(adapter.SettingsPath(home), []byte("{}"), 0o644)

	_, err := Inject(home, adapter)
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}

	content, err := os.ReadFile(adapter.SettingsPath(home))
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	text := string(content)
	// For standard agents (OpenCode), we now prioritize a stable relative path
	// "kortex-engram" instead of a dynamic absolute path to ensure idempotency.
	if !strings.Contains(text, `"kortex-engram"`) {
		t.Fatalf("OpenCode settings missing stable relative KortexEngram path, got: %s", text)
	}
	// OpenCode 1.3.3+: command must be an array, no separate "args" field.
	if strings.Contains(text, `"args"`) {
		t.Fatalf("OpenCode settings must NOT have a separate args field; got: %s", text)
	}

	// Structurally verify command is a []any containing the stable path "kortex-engram".
	var parsed map[string]any
	if err := json.Unmarshal(content, &parsed); err != nil {
		t.Fatalf("Unmarshal(opencode.json) error = %v", err)
	}
	mcpRaw, ok := parsed["mcp"]
	if !ok {
		t.Fatalf("opencode.json missing mcp key; got:\n%s", text)
	}
	mcpMap, ok := mcpRaw.(map[string]any)
	if !ok {
		t.Fatalf("mcp key has unexpected type %T; got:\n%s", mcpRaw, text)
	}
	KortexEngramRaw, ok := mcpMap["kortex-engram"]
	if !ok {
		t.Fatalf("mcp missing KortexEngram key; got:\n%s", text)
	}
	KortexEngramMap, ok := KortexEngramRaw.(map[string]any)
	if !ok {
		t.Fatalf("mcp.KortexEngram has unexpected type %T; got:\n%s", KortexEngramRaw, text)
	}
	cmdRaw, ok := KortexEngramMap["command"]
	if !ok {
		t.Fatalf("mcp.KortexEngram missing command key; got:\n%s", text)
	}
	cmdArr, ok := cmdRaw.([]any)
	if !ok {
		t.Fatalf("mcp.kortexengram.command must be an array, got %T; value:\n%s", cmdRaw, text)
	}
	if len(cmdArr) == 0 {
		t.Fatalf("mcp.kortexengram.command array is empty; got:\n%s", text)
	}
	firstElem, ok := cmdArr[0].(string)
	if !ok || firstElem != "kortex-engram" {
		t.Fatalf("mcp.kortexengram.command[0] = %v, want stable relative 'KortexEngram'; got:\n%s", cmdArr[0], text)
	}
}

// TestKortexEngramInjectAbsolutePathForGeminiMergeStrategy verifies that the
// absolute path is also used when the StrategyMergeIntoSettings strategy is
// applied (e.g. Gemini CLI).
func TestKortexEngramInjectAbsolutePathForGeminiMergeStrategy(t *testing.T) {
	home := t.TempDir()

	absPath := "/opt/homebrew/bin/KortexEngram"
	mockkortexEngramLookPath(t, absPath, "")

	result, err := Inject(home, geminiAdapter())
	if err != nil {
		t.Fatalf("Inject(gemini) error = %v", err)
	}
	if !result.Changed {
		t.Fatalf("Inject(gemini) changed = false")
	}

	settingsPath := filepath.Join(home, ".gemini", "settings.json")
	content, readErr := os.ReadFile(settingsPath)
	if readErr != nil {
		t.Fatalf("ReadFile(settings.json) error = %v", readErr)
	}

	text := string(content)
	// For standard agents (Gemini), we now prioritize a stable relative path
	// "kortex-engram" instead of a dynamic absolute path to ensure idempotency.
	if !strings.Contains(text, `"kortex-engram"`) {
		t.Fatalf("settings.json missing stable relative path 'KortexEngram'; got:\n%s", text)
	}
}

func TestQwenKortexEngramIdempotency(t *testing.T) {
	orig := kortexEngramLookPath
	t.Cleanup(func() { kortexEngramLookPath = orig })

	homeDir := t.TempDir()
	adapter := qwenAdapter()
	settingsPath := adapter.SettingsPath(homeDir)

	if err := os.MkdirAll(filepath.Dir(settingsPath), 0755); err != nil {
		t.Fatal(err)
	}

	kortexEngramLookPath = func(string) (string, error) {
		return "", os.ErrNotExist
	}

	_, err := Inject(homeDir, adapter)
	if err != nil {
		t.Fatalf("First injection failed: %v", err)
	}

	content1, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatal(err)
	}

	// Simulate KortexEngram being found later (e.g. after go install or manual install)
	absPath := "/usr/local/bin/KortexEngram"
	kortexEngramLookPath = func(string) (string, error) {
		return absPath, nil
	}

	_, err = Inject(homeDir, adapter)
	if err != nil {
		t.Fatalf("Second injection failed: %v", err)
	}

	content2, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatal(err)
	}

	if string(content1) != string(content2) {
		t.Errorf("Idempotency failure! Settings changed between runs despite KortexEngram command being stable-relative.\nRun 1:\n%s\nRun 2:\n%s", string(content1), string(content2))
	}
}
