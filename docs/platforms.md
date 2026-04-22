# Plataformas Suportadas

← [Voltar para o README](../README.md)

---

| Plataforma | Gerenciador de Pacotes | Status |
|----------|----------------|--------|
| macOS (Apple Silicon + Intel) | Homebrew | Suportado |
| Linux (Ubuntu/Debian) | apt | Suportado |
| Linux (Arch) | pacman | Suportado |
| Linux (Família Fedora/RHEL) | dnf | Suportado |
| Windows 10/11 | winget | Suportado |

Derivados são detectados via `ID_LIKE` em `/etc/os-release` (Linux Mint, Pop!_OS, Manjaro, EndeavourOS, CentOS Stream, Rocky Linux, AlmaLinux, etc.).

Os binários de lançamento são compilados para `linux`, `darwin` e `windows` tanto para `amd64` quanto para `arm64`.

---

## Notas sobre Windows

- **winget** é usado como o gerenciador de pacotes padrão (pré-instalado no Windows 10/11).
- **instalações globais npm** não requerem `sudo` no Windows (graváveis pelo usuário por padrão).
- **curl** já vem pré-instalado no Windows 10+ e não requer instalação separada.
- **PowerShell** é o shell padrão quando `$SHELL` não está definido.
- Arquivos de lançamento usam o formato `.zip` no Windows (`.tar.gz` no macOS/Linux).
- **Kortex CLI no Windows** funciona tanto no Git Bash quanto no PowerShell. o kortex instala um shim `kortex.ps1` que delega automaticamente para o Git Bash, portanto não é necessária a troca manual de shell.

---

## Verificação de Segurança no Windows

Alguns produtos antivírus podem sinalizar binários Go não assinados de forma heurística.

Use o checksum do lançamento para verificar a integridade:

```powershell
# 1) Baixe o arquivo checksums.txt da mesma tag de lançamento
# 2) Calcule o hash local
Get-FileHash .\kortex_<VERSÃO>_windows_amd64.zip -Algorithm SHA256

# 3) Compare o hash com a entrada correspondente no arquivo checksums.txt
```

Se o hash coincidir com o do `checksums.txt`, o arquivo é autêntico para aquele lançamento.

---

## Caminhos de Configuração no Windows

| Agente | Caminho de Configuração no Windows |
|-------|-------------------|
| Claude Code | `%USERPROFILE%\.claude\` |
| OpenCode | `%USERPROFILE%\.config\opencode\` |
| Gemini CLI | `%USERPROFILE%\.gemini\` |
| Cursor | `%USERPROFILE%\.cursor\` |
| VS Code Copilot | `%APPDATA%\Code\User\` (configurações, MCP, prompts) + `%USERPROFILE%\.copilot\` (skills) |
| Codex | `%USERPROFILE%\.codex\` |
| Windsurf | `%USERPROFILE%\.codeium\windsurf\` (skills, MCP, regras) + `%APPDATA%\Windsurf\User\` (configurações) |
| Kimi | `%USERPROFILE%\.kimi\` (inclui `config.toml`, prompt de sistema, agentes, MCP) |
| Antigravity | `%USERPROFILE%\.gemini\antigravity\` |
| Kiro IDE | `%USERPROFILE%\.kiro\steering\` (prompts) + `%USERPROFILE%\.kiro\skills\` (skills) + `%USERPROFILE%\.kiro\agents\` (agentes SDD) + `%APPDATA%\kiro\User\settings.json` (configurações) + `%USERPROFILE%\.kiro\settings\mcp.json` (MCP) |
