# Arquitetura e Desenvolvimento

← [Voltar para o README](../README.md)

---

## Arquitetura

```
cmd/kortex/             Ponto de entrada da CLI
internal/
  app/                     Despacho de comandos + fiação de runtime
  model/                   Tipos de domínio (agentes, componentes, skills, presets, personas)
  catalog/                 Definições de registro (agentes, skills, componentes)
  system/                  Detecção de SO/distro, verificações de dependência, travas de plataforma
  cli/                     Flags de instalação, validação, orquestração, dry-run
  planner/                 Gráfico de dependências, resolução, ordenação, payloads de revisão
  installcmd/              Resolvido de comandos ciente de perfil (brew/apt/pacman/dnf/winget/go install)
  pipeline/                Execução em estágios + orquestração de rollback
  backup/                  Snapshot de configuração + restauração
  assets/                  Arquivos de skill embutidos + templates de persona
  components/              Lógica de instalação/injeção por componente
    engram/  sdd/  skills/  mcp/  persona/  theme/  permissions/  kortex/
    filemerge/             Mesclagem de arquivos baseada em marcadores (injeção sem sobrescrever)
  agents/                  Adaptadores de agentes (estratégia de config por agente)
    claude/  opencode/  gemini/  cursor/  vscode/  codex/  windsurf/  antigravity/
  opencode/                Utilitários de parsing de modelo/config do OpenCode
  state/                   Rastreamento de estado da instalação
  update/                  Lógica de auto-atualização + upgrade
  verify/                  Verificações de integridade pós-aplicação + relatórios
  tui/                     Interface Bubbletea TUI (tema Rose Pine)
    styles/  screens/
scripts/                   Scripts de instalação (bash + PowerShell)
e2e/                       Testes E2E baseados em Docker (Ubuntu + Arch)
testdata/                  Arquivos fixos (fixtures) para testes
```

---

## Testes

```bash
# Testes unitários
go test ./...

# Docker E2E (Ubuntu + Arch, requer Docker)
RUN_FULL_E2E=1 RUN_BACKUP_TESTS=1 ./e2e/docker-test.sh

# Teste de fumaça dry-run (macOS/Linux)
kortex install --dry-run --agent claude-code --preset minimal

# Teste de fumaça dry-run (Windows PowerShell)
kortex.exe install --dry-run --agent claude-code --preset minimal
```

Cobertura de testes:

- **26 pacotes de teste** em toda a base de código
- **260+ funções de teste** cobrindo todos os adaptadores de agentes, componentes e detecção de sistema
- **78 funções de teste E2E** rodando em containers Docker (Ubuntu + Arch)
- **17 arquivos golden** para testes de snapshot da saída dos componentes
- Pipeline completo testado: detecção, planejamento, execução, backup, restauração, verificação
- Todos os 8 adaptadores de agentes possuem testes unitários com validação de caminho multiplataforma

---

## Relacionamento com Kortex.Dots

| | Kortex.Dots | Kortex Stack |
|--|---------------|-----------------|
| **Propósito** | Ambiente de dev (editores, shells, terminais) | Camada de desenvolvimento com IA (agentes, memória, skills) |
| **Instalações** | Neovim, Fish/Zsh, Tmux/Zellij, Ghostty | Configura Claude Code, OpenCode, Gemini CLI, Cursor, VS Code Copilot, Codex, Windsurf, Antigravity |
| **Sobreposição** | Nenhuma — complementares | Nenhuma — camada diferente |

Instale o Kortex.Dots primeiro para seu ambiente de desenvolvimento, depois o Kortex Stack para a camada de IA por cima.

---

## Licença

MIT
