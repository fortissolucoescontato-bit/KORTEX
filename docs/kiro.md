# Kiro IDE

← [Voltar para o README](../README.md)

---

Este documento explica como o kortex se integra ao **Kiro IDE** e o que é instalado em sua configuração local do Kiro.

## Visão Geral

O kortex suporta o Kiro como uma plataforma de **subagentes nativos** (`kiro-ide`).

Quando configurado, o kortex instala:

| Artefato | Caminho |
|----------|------|
| Arquivo de Steering (direção) | `~/.kiro/steering/kortex.md` |
| Agentes SDD nativos | `~/.kiro/agents/sdd-{fase}.md` *(10 arquivos)* |
| Diretório de Skills | `~/.kiro/skills/` |
| Configuração MCP | `~/.kiro/settings/mcp.json` *(raiz separada — veja nota abaixo)* |

> **A instalação automática não é suportada.** O Kiro deve ser instalado manualmente antes de executar o kortex.
> Baixe em: [kiro.dev/downloads](https://kiro.dev/downloads)

---

## Detecção

O kortex usa **dois sinais** para detectar o Kiro:

1. **Presença do diretório `~/.kiro`** — usado por `system.ScanConfigs` para o fluxo de autodetecção da instalação/TUI. Se `~/.kiro` existir no disco, o Kiro é exibido como detectado no instalador, independentemente do binário estar no `PATH`.
2. **Binário `kiro` no `PATH`** — usado por `adapter.Detect()` para o fluxo de sincronização/upgrade e para confirmar que a IDE é realmente executável.

Na prática: **o instalador detecta o Kiro a partir de `~/.kiro`**, não pelo `PATH`. Se você tem o Kiro instalado, mas o diretório `~/.kiro` ainda não foi criado (ex: antes do primeiro lançamento), execute o Kiro uma vez para inicializar seu diretório de configuração e depois execute o `kortex install` novamente.

---

## Modelo de Execução do SDD

O Kiro funciona com **delegação nativa de subagentes** via `~/.kiro/agents/`.

O orquestrador permanece no arquivo de steering e coordena a execução das fases, enquanto cada fase roda em seu arquivo de agente do Kiro dedicado:

```
sdd-init → sdd-explore → sdd-propose → sdd-spec → sdd-design → sdd-tasks → sdd-apply → sdd-verify → sdd-archive (+ sdd-onboard)
```

Isso segue a mesma arquitetura SDD usada no kortex: o orquestrador coordena, os agentes de fase executam e o Engram persiste os artefatos entre as fases.

**Portões de aprovação** continuam sendo obrigatórios antes das fases `apply` (aplicar) e `archive` (arquivar).

---

## Integração Nativa de Specs do Kiro

O Kiro possui um fluxo de trabalho de specs embutido que o kortex aproveita. Para mudanças médias e grandes, o orquestrador usará artefatos nativos do Kiro em:

```
.kiro/specs/<recurso>/
├── requirements.md
├── design.md
└── tasks.md
```

**Arquivos de Steering** em `.kiro/steering/*.md` fornecem contexto persistente do workspace entre as sessões — trate-os como um contexto de sistema sempre ativo para as convenções do seu projeto, decisões de arquitetura e regras da equipe.

**Classificação de tamanho** encaminha as tarefas através de caminhos Pequeno / Médio / Grande para decidir a profundidade do planejamento:

| Tamanho | Abordagem |
|------|----------|
| Pequeno | Em linha — sem fases formais de SDD |
| Médio | Specs nativas do Kiro (`.kiro/specs/`) + Engram |
| Grande | Ciclo completo de SDD: explore → propose → spec → design → tasks → apply → verify → archive |

---

## Formato do Arquivo de Steering

O arquivo de steering gravado pelo kortex usa o seguinte frontmatter:

```yaml
---
inclusion: always
---
```

`inclusion: always` garante que o Kiro carregue este contexto em cada conversa automaticamente, independentemente do workspace ou tipo de arquivo.

## Frontmatter do Agente Nativo

Os agentes de fase SDD do Kiro são gerados com frontmatter YAML incluindo:

- `name`
- `description`
- `tools`
- `model`
- `includeMcpJson: true`

O valor de `model` é injetado durante a sincronização a partir das atribuições de alias do Claude (`opus|sonnet|haiku`) para IDs de modelos nativos do Kiro.

---

## Caminhos de Configuração por Plataforma

### macOS

| Artefato | Caminho |
|----------|------|
| Dir de config global | `~/Library/Application Support/Kiro/User` |
| Arquivo de Steering | `~/.kiro/steering/kortex.md` |
| Dir de Skills | `~/.kiro/skills/` |
| Caminho de Settings | `~/Library/Application Support/Kiro/User/settings.json` |
| Configuração MCP | `~/.kiro/settings/mcp.json` |

### Windows

| Artefato | Caminho |
|----------|------|
| Dir de config global | `%APPDATA%\kiro\User` |
| Arquivo de Steering | `%USERPROFILE%\.kiro\steering\kortex.md` |
| Dir de Skills | `%USERPROFILE%\.kiro\skills\` |
| Caminho de Settings | `%APPDATA%\kiro\User\settings.json` |
| Configuração MCP | `%USERPROFILE%\.kiro\settings\mcp.json` |

### Linux (XDG)

| Artefato | Caminho |
|----------|------|
| Dir de config global | `$XDG_CONFIG_HOME/kiro/user` *(fallback: `~/.config/kiro/user`)* |
| Arquivo de Steering | `~/.kiro/steering/kortex.md` |
| Dir de Skills | `~/.kiro/skills/` |
| Caminho de Settings | `$XDG_CONFIG_HOME/kiro/user/settings.json` |
| Configuração MCP | `~/.kiro/settings/mcp.json` |

---

## ⚠️ Layout de Raiz Dividida (Split-Root)

O Kiro usa um **layout de raiz dividida** — os arquivos gerenciados pelo kortex e as configurações da IDE vivem em diretórios diferentes:

- **Steering, skills e agentes nativos** → `~/.kiro/` (ou `%USERPROFILE%\.kiro\` no Windows)
  - `~/.kiro/steering/kortex.md` — persona do orquestrador
  - `~/.kiro/skills/` — arquivos de skill SDD
  - `~/.kiro/agents/` — subagentes de fase SDD
- **Configurações da IDE** → diretório de Usuário do Kiro nativo da plataforma (apenas `settings.json`)
  - macOS: `~/Library/Application Support/Kiro/User/settings.json`
  - Windows: `%APPDATA%\kiro\User\settings.json`
  - Linux: `$XDG_CONFIG_HOME/kiro/user/settings.json`
- **Configuração MCP** → sempre `~/.kiro/settings/mcp.json` (ou `%USERPROFILE%\.kiro\settings\mcp.json` no Windows)

Se as ferramentas MCP não estiverem carregando, verifique `~/.kiro/settings/mcp.json`.  
Se as configurações do app Kiro não estiverem sendo aplicadas, verifique o diretório de Usuário nativo da plataforma (`settings.json`).  
Se as skills ou o steering do kortex estiverem faltando, verifique `~/.kiro/skills/` e `~/.kiro/steering/`.

---

## Resumo de Capacidades

| Capacidade | Status |
|------------|--------|
| Skills | ✅ Sim |
| Prompt de sistema | ✅ Sim |
| MCP | ✅ Sim |
| Estilos de saída | ❌ Não |
| Comandos Slash | ❌ Não |
| Modelo de delegação | Total (subagentes nativos) |
| Instalação automática | ❌ Não — instalação manual necessária |
