# Agentes Suportados

← [Voltar para o README](../README.md)

---

## Matriz de Agentes

| Agente          | ID               | Skills       | MCP | Delegação                    | Estilos de Saída | Comandos Slash | Caminho de Configuração             |
| --------------- | ---------------- | ------------ | --- | ---------------------------- | ---------------- | --------------- | ----------------------------------- |
| Claude Code     | `claude-code`    | Sim          | Sim | Total (ferramenta Task)      | Sim              | Não             | `~/.claude`                         |
| OpenCode        | `opencode`       | Sim          | Sim | Total (overlay multi-modo)   | Não              | Sim             | `~/.config/opencode`                |
| Gemini CLI      | `gemini-cli`     | Sim          | Sim | Total (experimental)         | Não              | Não             | `~/.gemini`                         |
| Cursor          | `cursor`         | Sim          | Sim | Total (subagentes nativos)   | Não              | Não             | `~/.cursor`                         |
| VS Code Copilot | `vscode-copilot` | Sim          | Sim | Total (runSubagent)          | Não              | Não             | `~/.copilot` + Perfil de Usuário VS Code |
| Codex           | `codex`          | Sim          | Sim | Agente solo                  | Não              | Não             | `~/.codex`                          |
| Windsurf        | `windsurf`       | Sim (nativo) | Sim | Agente solo                  | Não              | Não             | `~/.codeium/windsurf`               |
| Antigravity     | `antigravity`    | Sim (nativo) | Sim | Agente solo + Mission Control| Não              | Não             | `~/.gemini/antigravity`             |
| Kimi            | `kimi`           | Sim          | Sim | Total (agentes personalizados nativos) | Não      | Não             | `~/.kimi`                           |
| Qwen Code       | `qwen-code`      | Sim          | Sim | Total (subagentes nativos)   | Não              | Sim             | `~/.qwen`                           |
| Kiro IDE        | `kiro-ide`       | Sim          | Sim | Total (subagentes nativos)   | Não              | Não             | `~/.kiro`                           |

Todos os agentes recebem o **orquestrador SDD completo** injetado em seu prompt de sistema, além de arquivos de skill gravados em seus respectivos diretórios de skills. O agente lida com o SDD automaticamente quando a tarefa é grande o suficiente ou quando o usuário solicita explicitamente — sem necessidade de configuração manual.

---

## Modelos de Delegação

| Modelo                | Como Funciona                                                                                                                         | Agentes                                                           |
| --------------------- | ------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------- |
| **Total (subagentes)**| Cada fase do SDD roda em uma janela de contexto isolada via delegação nativa de subagentes. O orquestrador coordena; subagentes executam. | Claude Code, OpenCode, Gemini CLI, Cursor, VS Code Copilot, Kimi, Kiro IDE, Qwen Code |
| **Agente solo**       | Todas as fases do SDD rodam na mesma conversa. O orquestrador É o executor. O Engram fornece persistência entre as fases.             | Codex, Windsurf, Antigravity                                      |

### Subagentes Nativos do Cursor

O Cursor usa seu sistema embutido em `.cursor/agents/`. O `kortex` grava 10 arquivos de agente em `~/.cursor/agents/sdd-{fase}.md` — um para cada fase do SDD. O Agente do Cursor delega automaticamente para o subagente correto com base no campo `description` do frontmatter YAML de cada arquivo.

- `sdd-explore` e `sdd-verify` rodam com `readonly: true`
- Cada subagente recebe sua própria janela de contexto (contexto limpo, sem poluição)
- O orquestrador resolve regras compactas do registro de skills e as passa na mensagem de invocação

### Windsurf Cascade

O Windsurf funciona como um agente solo (sem subagentes personalizados). O orquestrador aproveita as funcionalidades nativas do Windsurf:

- **Modo Plano (Plan Mode)** — cria documentos de plano persistentes que podem ser @mencionados em várias sessões; ideal para especificações e artefatos de design em mudanças grandes
- **Modo Código (Code Mode)** — modo padrão de execução agentica
- **Workflows Nativos** — `sdd-new` está disponível como um workflow em `.windsurf/workflows/sdd-new.md`
- **Classificação de Tamanho** — o orquestrador encaminha as tarefas através de caminhos de decisão Pequeno/Médio/Grande

### Antigravity + Mission Control

O Antigravity é uma plataforma focada em agentes com subagentes embutidos (Navegador, Terminal) gerenciados pelo Mission Control. No entanto, a criação de subagentes personalizados ainda não está disponível. As fases do SDD rodam em linha, com o Mission Control lidando com a delegação automática para subagentes embutidos quando ferramentas especializadas são necessárias (ex: Navegador para pesquisa durante o `sdd-explore`).

### Subagentes Nativos do Kiro

O Kiro usa agentes personalizados nativos em `~/.kiro/agents/`. O `kortex` grava os 10 agentes de fase (`sdd-init` até `sdd-onboard`) e resolve o campo `model:` durante a injeção a partir de atribuições de alias do Claude (`opus|sonnet|haiku`) para IDs de modelos nativos do Kiro.

- O frontmatter inclui `includeMcpJson: true` para todos os agentes de fase
- Ferramentas específicas de fase são preservadas (`sdd-explore` e `sdd-verify` usam read/shell/context7 conforme necessário)
- O orquestrador permanece na direção (`~/.kiro/steering/kortex.md`) e delega a execução aos subagentes nativos

---

## Suporte ao Modo SDD

| Recurso | Claude Code | OpenCode | Gemini CLI | Cursor | VS Code Copilot | Codex | Windsurf | Antigravity | Kiro IDE | Qwen Code |
|---------|:-----------:|:--------:|:----------:|:------:|:---------------:|:-----:|:--------:|:-----------:|:--------:|:---------:|
| Orquestrador SDD | Sim | Sim | Sim | Sim | Sim | Sim | Sim | Sim | Sim | Sim |
| SDD modo único  | Sim | Sim | Sim | Sim | Sim | Sim | Sim | Sim | Sim | Sim |
| SDD multi-modo  | —   | Sim | —   | —   | —   | —   | —   | —   | Sim*| —   |

**Multi-modo** (atribuir diferentes modelos de IA para cada fase do SDD) é suportado nativamente pelo **OpenCode** (via seu sistema de provedores) e pelo **Kiro IDE** (via frontmatter `model:` do subagente nativo — cada agente de fase roda com seu próprio ID de modelo). Todos os outros agentes rodam em **modo único** — o orquestrador gerencia tudo usando o modelo que o agente já estiver utilizando.

> \* **Kiro multi-modo** atribui modelos por fase através de `KiroModelAssignments` (configurado via *Configurar Modelos → Configurar modelos do Kiro* na TUI). O alias selecionado (`opus|sonnet|haiku`) é resolvido para um ID de modelo nativo do Kiro e carimbado em cada `~/.kiro/agents/sdd-{fase}.md` no momento da sincronização.

---

## Notas por Agente

### Claude Code

- Subagentes via a ferramenta nativa Task com janelas de contexto isoladas
- Servidores MCP configurados como plugins em `~/.claude/mcp/`
- Estilos de saída em `~/.claude/output-styles/`
- Prompt de sistema via seções markdown em `~/.claude/CLAUDE.md`

### OpenCode

- Overlay completo multi-agente com 12 agentes nomeados em `opencode.json`
- Comandos Slash para fases do SDD (`/sdd-new`, `/sdd-explore`, etc.)
- Plugin de agentes em segundo plano para execução paralela
- Pré-requisito multi-modo: conecte seus provedores de IA primeiro, depois execute `opencode models --refresh`

### Gemini CLI

- Subagentes são experimentais: requer `experimental.enableAgents: true` em `settings.json`
- Subagentes personalizados definidos como arquivos markdown em `~/.gemini/agents/`

### Cursor
- Subagentes nativos via `~/.cursor/agents/sdd-{fase}.md` (10 arquivos instalados pelo kortex)
- Skills em `~/.cursor/skills/`
- Prompt de sistema em `~/.cursor/rules/kortex.mdc`
- Configuração MCP em `~/.cursor/mcp.json`

### VS Code Copilot

- Usa a ferramenta `runSubagent` com suporte para execução paralela
- Skills em `~/.copilot/skills/`
- Prompt de sistema em `Code/User/prompts/kortex.instructions.md`
- Configuração MCP em `Code/User/mcp.json`

### Codex

- Agente nativo CLI com configuração TOML em `~/.codex/config.toml`
- Skills em `~/.codex/skills/`
- Prompt de sistema em `~/.codex/agents.md`
- Arquivos de instrução Engram em `~/.codex/engram-instructions.md`

### Windsurf

- Skills em `~/.codeium/windsurf/skills/` (recurso nativo do Windsurf)
- Configuração MCP em `~/.codeium/windsurf/mcp_config.json`
- Regras globais em `~/.codeium/windsurf/memories/global_rules.md`
- Workflows em `.windsurf/workflows/` (escopo do workspace)

### Antigravity

- Skills em `~/.gemini/antigravity/skills/` (recurso nativo do Antigravity)
- Configuração MCP em `~/.gemini/antigravity/mcp_config.json`
- Prompt de sistema anexado ao `~/.gemini/GEMINI.md` (compartilhado com o Gemini CLI — aviso de colisão se ambos estiverem instalados)
- O Mission Control lida com a delegação de subagentes embutidos (Navegador, Terminal) automaticamente
- Configurações gerenciadas via a UI de configurações do Agente na IDE, não via `settings.json`

### Kimi

- A instalação requer o gerenciador de pacotes Python `uv` (`uv tool install kimi-cli`).
- Agente personalizado raiz em `~/.kimi/agents/carbon.yaml` com `system_prompt_path: ../KIMI.md`
- `KIMI.md` é um template Jinja que inclui arquivos de prompt modulares:
  `persona.md`, `output-style.md`, `engram-protocol.md`, `sdd-orchestrator.md`
- Variáveis embutidas do Kimi são preservadas no `KIMI.md`: `${KIMI_AGENTS_MD}` e `${KIMI_SKILLS}`

### Kiro IDE

- **Detecção**: o kortex detecta o Kiro a partir de sua raiz de configuração (`~/.kiro`) durante a descoberta da instalação/TUI — `~/.kiro` deve existir (criado no primeiro lançamento do Kiro). `kiro` no `PATH` também é verificado para fluxos de sincronização/atualização, mas não é obrigatório para a autodetecção da instalação
- **Arquivo de Direção (Steering)** (todas as plataformas): `~/.kiro/steering/kortex.md` com frontmatter `inclusion: always`
- Subagentes nativos em `~/.kiro/agents/sdd-{fase}.md` (10 arquivos)
- Skills (todas as plataformas) em `~/.kiro/skills/`
- **Configuração MCP em uma raiz separada** — sempre `~/.kiro/settings/mcp.json` (macOS/Linux) ou `%USERPROFILE%\.kiro\settings\mcp.json` (Windows), independentemente do GlobalConfigDir
- Workflow nativo de specs do Kiro: `.kiro/specs/<recurso>/requirements.md`, `design.md`, `tasks.md` — com portões de aprovação antes das fases de aplicação e arquivamento
- Apenas instalação manual — baixe em [kiro.dev/downloads](https://kiro.dev/downloads)
- Consulte [docs/kiro.md](kiro.md) para referência completa de caminhos e detalhes de comportamento do SDD

### Qwen Code
- **Detecção**: o kortex detecta o Qwen Code a partir de sua raiz de configuração (`~/.qwen`) e verifica o binário `qwen` no `PATH`
- **Raiz de configuração**: `~/.qwen/` (multiplataforma)
- **Prompt de sistema**: `~/.qwen/QWEN.md` (gerenciado via `StrategyFileReplace`)
- **Skills**: `~/.qwen/skills/`
- **Configuração MCP**: `~/.qwen/settings.json` (gerenciado via `StrategyMergeIntoSettings` com a chave `mcpServers`)
- **Comandos Slash**: `~/.qwen/commands/*.md` — suporta comandos slash personalizados com namespace (ex: `commands/sdd/init.md` → `/sdd:init`)
- **Permissões**: modo `auto_edit` — aprova automaticamente edições de arquivo, aprovação manual para comandos de shell
- **Instalação**: via npm — `npm install -g @qwen-code/qwen-code@latest`
- **Slug Engram**: `"qwen-code"` para integração com `engram setup`
- **Orquestrador SDD**: `internal/assets/qwen/sdd-orchestrator.md` com referências de caminho específicas do Qwen
