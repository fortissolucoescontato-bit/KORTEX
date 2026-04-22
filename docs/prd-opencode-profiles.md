# PRD: Perfis SDD do OpenCode

> **Crie perfis de modelos intercambiáveis para o OpenCode — alterne entre configurações do orquestrador com um único Tab.**

**Versão**: 0.1.0-draft
**Autor**: Kortex Programming
**Data**: 03-04-2026
**Status**: Rascunho

---

## 1. Declaração do Problema

Hoje, o OpenCode permite ter apenas UM `sdd-orchestrator` com UM único conjunto de modelos atribuídos aos subagentes SDD. Isso força o usuário a escolher entre:

- **Qualidade máxima** (Opus em tudo → caro e lento)
- **Equilíbrio** (Opus orquestrador + Sonnet subagentes → o padrão atual)
- **Economia** (Sonnet/Haiku em tudo → rápido e barato, mas menos potente)

O problema: **você não pode alternar entre essas configurações sem editar manualmente o `opencode.json`** cada vez que quer passar de um modo para outro. Na prática, um desenvolvedor precisa de diferentes perfis para diferentes momentos:

- **"Vou fazer algo pesado"** → orquestrador Opus, subagentes Sonnet
- **"É algo simples, não quero queimar tokens"** → tudo Haiku
- **"Quero testar um modelo novo do Google"** → orquestrador Gemini, subagentes mistos
- **"Estou apenas revisando um PR"** → perfil leve

Hoje isso é uma dor de cabeça manual. Este recurso resolve isso.

---

## 2. Visão

**O usuário cria N perfis de modelos a partir da TUI. Cada perfil gera seu próprio `sdd-orchestrator-{nome}` com seus próprios subagentes no `opencode.json`. No OpenCode, ele pressiona Tab e vê todos os orquestradores disponíveis — alternando entre perfis como quem troca de marcha.**

```
┌─────────────────────────────────────────────────────────────┐
│  opencode.json                                               │
│                                                              │
│  ┌──────────────────────┐   ┌──────────────────────────────┐ │
│  │  sdd-orchestrator    │   │  sdd-orchestrator-barato     │ │
│  │  (opus + sonnet)     │   │  (haiku em tudo)             │ │
│  │                      │   │                              │ │
│  │  sdd-init     sonnet │   │  sdd-init-barato     haiku   │ │
│  │  sdd-explore  sonnet │   │  sdd-explore-barato  haiku   │ │
│  │  sdd-apply    sonnet │   │  sdd-apply-barato    haiku   │ │
│  │  ...                 │   │  ...                        │ │
│  └──────────────────────┘   └──────────────────────────────┘ │
│                                                              │
│  Tab no OpenCode → escolha qual orquestrador usar            │
└─────────────────────────────────────────────────────────────┘
```

---

## 3. Usuários-Alvo

| Usuário | Ponto de Dor | Como os Perfis Ajudam |
|------|-----------|-------------------|
| **Power user com múltiplos provedores** | Quer testar Anthropic vs Google vs OpenAI para SDD sem mexer na config | Cria um perfil por provedor, alterna com Tab |
| **Desenvolvedor atento aos custos** | Quer um modo "barato" para tarefas simples | Perfil "barato" com Haiku/Flash, perfil "premium" com Opus |
| **Team lead** | Quer padronizar perfis para a equipe | Os perfis vivem no `opencode.json`, sincronizáveis |
| **Experimentador** | Quer testar modelos novos sem quebrar sua config padrão | Perfil experimental, o padrão permanece intacto |

---

## 4. Escopo

### No Escopo (v1)
- Criação de perfis a partir da TUI (nova tela)
- Visualização de perfis existentes
- **Edição de perfis existentes a partir da TUI** (selecionar perfil → modificar modelos → sync)
- **Exclusão de perfis a partir da TUI** (selecionar perfil → confirmar → remove orchestrator + subagentes do JSON → sync)
- Geração de N orchestrators + N×9 subagentes no `opencode.json`
- Atualização de perfis existentes durante Sincronização / Atualização+Sincronização
- Prompts compartilhados: um arquivo por fase, reutilizado por todos os perfis
- Flag da CLI para criar perfis (`--profile`)

### Fora do Escopo (permanentemente)
- **Perfis para o Claude Code** — NÃO SE APLICA. O Claude Code usa um mecanismo completamente diferente (CLAUDE.md + ferramenta Task). O recurso de perfis é exclusivo do OpenCode porque depende do sistema de agents/sub-agents do `opencode.json` e da seleção por Tab. Isso NÃO é "futuro" — é uma decisão de arquitetura.

### Fora do Escopo (v1, consideração futura)
- Exportar/importar perfis entre máquinas

---

## 5. Requisitos Detalhados

### 5.1 TUI: Tela de Criação de Perfis

**R-PROF-01**: A tela de Boas-vindas DEVE incluir uma nova opção **"OpenCode SDD Profiles"** abaixo de "Configure Models".

**R-PROF-02**: Se já existirem perfis criados, a opção DEVE mostrar a contagem: `"OpenCode SDD Profiles (2)"`.

**R-PROF-03**: A tela de perfis DEVE mostrar os perfis existentes com ações disponíveis:

```
┌─────────────────────────────────────────────────────────┐
│  OpenCode SDD Profiles                                   │
│                                                          │
│  Perfis existentes:                                      │
│    ✦ default ─── anthropic/claude-opus-4                 │
│    • barato ──── anthropic/claude-haiku-3.5              │
│    • gemini ──── google/gemini-2.5-pro                   │
│                                                          │
│  ▸ Criar novo perfil                                     │
│    Voltar                                                │
│                                                          │
│  j/k: navegar • enter: editar • n: novo • d: excluir     │
│  esc: voltar                                             │
└─────────────────────────────────────────────────────────┘
```

**R-PROF-04**: Ao selecionar "Criar novo perfil" (ou pressionar `n`), o usuário DEVE:
1. **Inserir um nome** para o perfil (texto livre, validado para slug: minúsculas, sem espaços, alfanumérico + hifens)
2. **Selecionar o modelo do orquestrador** (reutilizando o ModelPicker existente — provedor → modelo)
3. **Selecionar modelos para os subagentes** (reutilizando o ModelPicker existente com as 9+1 linhas: Definir todos + 9 fases)
4. **Confirmar** → o perfil é gerado e a sincronização é executada

**R-PROF-05**: O nome "default" ESTÁ RESERVADO para o orquestrador base (`sdd-orchestrator`). O usuário NÃO pode criar um perfil chamado "default".

**R-PROF-06**: Se o usuário inserir um nome que já existe, DEVE-SE perguntar se ele deseja sobrescrever.

### 5.1b TUI: Edição de Perfil

**R-PROF-07**: Ao pressionar `enter` sobre um perfil existente na lista, o usuário entra no modo de edição. O fluxo é IDÊNTICO ao de criação, mas:
- O nome NÃO pode ser alterado (exibido como cabeçalho fixo)
- O modelo do orquestrador vem pré-selecionado com o valor atual
- Os modelos de subagentes vêm pré-selecionados com os valores atuais
- Ao confirmar, o perfil existente é sobrescrito e a sincronização é executada

**R-PROF-07b**: O perfil `default` também PODE ser editado — é o `sdd-orchestrator` base. Editar o default é equivalente ao que hoje faz "Configure Models → OpenCode", mas integrado no fluxo de perfis.

### 5.1c TUI: Exclusão de Perfil

**R-PROF-08**: Ao pressionar `d` sobre um perfil existente na lista, DEVE-SE mostrar uma tela de confirmação:

```
┌─────────────────────────────────────────────────────────┐
│  Excluir Perfil                                          │
│                                                          │
│  Tem certeza que deseja excluir o perfil "barato"?       │
│                                                          │
│  Isso removerá do opencode.json:                         │
│    • sdd-orchestrator-barato                             │
│    • sdd-init-barato                                     │
│    • sdd-explore-barato                                  │
│    • ... (10 agentes no total)                           │
│                                                          │
│  ▸ Excluir                                               │
│    Cancelar                                              │
│                                                          │
│  enter: selecionar • esc: cancelar                       │
└─────────────────────────────────────────────────────────┘
```

**R-PROF-08b**: Ao confirmar a exclusão:
1. Todos os agent keys do perfil são excluídos do `opencode.json` (`sdd-orchestrator-{nome}` + 10 subagentes `sdd-{fase}-{nome}`)
2. É executada uma gravação atômica do JSON atualizado
3. O resultado é exibido (sucesso ou erro)
4. Volta-se para a lista de perfis (com o perfil removido)

**R-PROF-08c**: O perfil `default` NÃO pode ser excluído. Pressionar `d` sobre o default NÃO faz nada (o atalho é ignorado). O default é o orquestrador base que sempre deve existir.

**R-PROF-08d**: A exclusão de um perfil NÃO exclui os arquivos de prompt compartilhados (`~/.config/opencode/prompts/sdd/*.md`) — eles são compartilhados por todos os perfis e mantidos enquanto existir pelo menos um perfil.

### 5.2 Convenção de Nomenclatura

**R-PROF-10**: O perfil DEFAULT (sem sufixo) gera os agentes com os nomes atuais:
- `sdd-orchestrator`
- `sdd-init`, `sdd-explore`, `sdd-propose`, `sdd-spec`, `sdd-design`, `sdd-tasks`, `sdd-apply`, `sdd-verify`, `sdd-archive`

**R-PROF-11**: Um perfil com nome `barato` gera agentes com sufixo:
- `sdd-orchestrator-barato`
- `sdd-init-barato`, `sdd-explore-barato`, ..., `sdd-archive-barato`

**R-PROF-12**: O `sdd-orchestrator-{nome}` DEVE ter `"mode": "primary"` para aparecer como selecionável com Tab no OpenCode. Os subagentes `sdd-{fase}-{nome}` DEVEM ter `"mode": "subagent"` e `"hidden": true`.

**R-PROF-13**: As permissões do orquestrador de um perfil DEVEM ser restritas aos seus próprios subagentes:
```json
{
  "permission": {
    "task": {
      "*": "deny",
      "sdd-*-barato": "allow"
    }
  }
}
```

### 5.3 Arquitetura de Prompt Compartilhado

**R-PROF-20**: Os prompts de cada fase SDD DEVEM residir em arquivos separados sob `~/.config/opencode/prompts/sdd/`:
```
~/.config/opencode/prompts/sdd/
├── orchestrator.md
├── sdd-init.md
├── sdd-explore.md
├── sdd-propose.md
├── sdd-spec.md
├── sdd-design.md
├── sdd-tasks.md
├── sdd-apply.md
├── sdd-verify.md
├── sdd-archive.md
└── sdd-onboard.md
```

**R-PROF-21**: O `prompt` de cada agente no opencode.json DEVE referenciar o arquivo compartilhado usando a sintaxe do OpenCode `{file:caminho}`:
```json
{
  "sdd-apply": {
    "mode": "subagent",
    "hidden": true,
    "model": "anthropic/claude-sonnet-4-20250514",
    "prompt": "{file:~/.config/opencode/prompts/sdd/sdd-apply.md}"
  },
  "sdd-apply-barato": {
    "mode": "subagent",
    "hidden": true,
    "model": "anthropic/claude-haiku-3.5-20241022",
    "prompt": "{file:~/.config/opencode/prompts/sdd/sdd-apply.md}"
  }
}
```

**R-PROF-22**: O conteúdo destes arquivos de prompt DEVE ser EXATAMENTE o mesmo que hoje é inserido em linha no overlay JSON. O refactor é uma extração sem mudança de comportamento.

**R-PROF-23**: O prompt do orquestrador (`orchestrator.md`) DEVE incluir um bloco `<!-- kortex:sdd-model-assignments -->` que é injetado dinamicamente com a tabela de modelos desse perfil específico.

**R-PROF-24**: Para o orquestrador de um perfil, o prompt DEVE referenciar os subagentes COM SUFIXO.

**Decisão arquitetônica**: O prompt do orquestrador NÃO é compartilhado entre perfis — cada perfil gera sua própria versão com:
- A tabela de atribuições de modelos desse perfil
- As referências a `sdd-{fase}-{sufixo}` corretas

Os prompts dos subagentes SÃO compartilhados porque são idênticos entre perfis (muda apenas o modelo, não o prompt).

### 5.4 Comportamento de Sincronização e Atualização

**R-PROF-30**: Durante a `Sincronização` ou `Atualização+Sincronização`, o sistema DEVE:
1. Detectar TODOS os perfis existentes no `opencode.json` (padrão: `sdd-orchestrator-*`)
2. Atualizar os prompts compartilhados em `~/.config/opencode/prompts/sdd/`
3. Regenerar os prompts de orquestrador de cada perfil (para injetar atribuições de modelos atualizadas)
4. NÃO modificar as atribuições de modelos dos perfis — apenas os prompts

**R-PROF-31**: Se um perfil tiver um subagente que referencia um modelo que não existe mais no cache do OpenCode, a sincronização DEVE:
- Emitir um **Aviso** ao usuário (não erro)
- Preservar a atribuição existente (o usuário pode ter configurado manualmente)

**R-PROF-32**: Os arquivos de prompt compartilhados DEVEM estar cobertos pelo backup pré-sincronização, assim como o `opencode.json`.

**R-PROF-33**: A sincronização DEVE ser idempotente: se os prompts já estiverem atualizados, o `filesChanged` NÃO deve aumentar.

### 5.5 Detecção de Perfil e Estado

**R-PROF-40**: Os perfis DEVEM ser detectados lendo o `opencode.json` existente, NÃO a partir de um arquivo de estado separado. O `opencode.json` É a fonte da verdade.

**R-PROF-41**: Um perfil é detectado pela presença de uma chave de agente que corresponda a `sdd-orchestrator-{nome}` com `"mode": "primary"`.

**R-PROF-42**: Ao detectar perfis existentes, o sistema DEVE inferir:
- **Nome**: o sufixo após `sdd-orchestrator-`
- **Modelo do orquestrador**: o campo `"model"` do orquestrador
- **Modelos de subagentes**: os campos `"model"` de `sdd-{fase}-{nome}`

**R-PROF-43**: O perfil padrão (`sdd-orchestrator` sem sufixo) SEMPRE existe quando o SDD está configurado. Perfis adicionais são opcionais.

### 5.6 Suporte CLI

**R-PROF-50**: O comando `sync` DEVE aceitar uma flag `--profile <nome>:<modelo-do-orquestrador>` que cria/atualiza um perfil durante a sincronização:
```bash
kortex sync --profile barato:anthropic/claude-haiku-3.5-20241022
```

**R-PROF-51**: DEVE-SE poder especificar múltiplas flags `--profile`:
```bash
kortex sync \
  --profile barato:anthropic/claude-haiku-3.5-20241022 \
  --profile premium:anthropic/claude-opus-4-20250514
```

**R-PROF-52**: O formato da flag é `nome:provedor/modelo`. Para atribuir modelos individuais a subagentes via CLI, usa-se a sintaxe estendida:
```bash
kortex sync --profile barato:anthropic/claude-haiku-3.5-20241022 \
  --profile-phase barato:sdd-apply:anthropic/claude-sonnet-4-20250514
```

---

## 6. Design Técnico

### 6.1 Modelo de Dados

```go
// Profile representa uma configuração nomeada de orquestrador SDD com atribuições de modelos.
type Profile struct {
    Name                string                       // ex: "barato", "premium"
    OrchestratorModel   model.ModelAssignment         // modelo do orquestrador
    PhaseAssignments    map[string]model.ModelAssignment // modelos por fase (sobrescritas opcionais)
}
```

### 6.2 Estrutura JSON do OpenCode (por perfil)

Para um perfil chamado "barato" com Haiku:

```json
{
  "agent": {
    "sdd-orchestrator-barato": {
      "mode": "primary",
      "description": "SDD Orchestrator (perfil barato) — haiku em tudo",
      "model": "anthropic/claude-haiku-3.5-20241022",
      "prompt": "... prompt do orquestrador com tabela de modelos específica do perfil e referências de subagentes ...",
      "permission": {
        "task": {
          "*": "deny",
          "sdd-*-barato": "allow"
        }
      },
      "tools": {
        "read": true,
        "write": true,
        "edit": true,
        "bash": true,
        "delegate": true,
        "delegation_read": true,
        "delegation_list": true
      }
    },
    "sdd-init-barato": {
      "mode": "subagent",
      "hidden": true,
      "model": "anthropic/claude-haiku-3.5-20241022",
      "description": "Bootstrap do contexto SDD (perfil barato)",
      "prompt": "{file:~/.config/opencode/prompts/sdd/sdd-init.md}"
    },
    "sdd-explore-barato": {
      "mode": "subagent",
      "hidden": true,
      "model": "anthropic/claude-haiku-3.5-20241022",
      "description": "Investigação da base de código (perfil barato)",
      "prompt": "{file:~/.config/opencode/prompts/sdd/sdd-explore.md}"
    }
    // ... restantes 7 subagentes com sufixo -barato
  }
}
```

### 6.3 Arquitetura de Arquivos de Prompt

```
~/.config/opencode/
├── opencode.json          (agentes com referências de modelo + prompt)
├── prompts/
│   └── sdd/
│       ├── sdd-init.md        (compartilhado por todos os perfis)
│       ├── sdd-explore.md     (compartilhado por todos os perfis)
│       ├── sdd-propose.md     (compartilhado)
│       ├── sdd-spec.md        (compartilhado)
│       ├── sdd-design.md      (compartilhado)
│       ├── sdd-tasks.md       (compartilhado)
│       ├── sdd-apply.md       (compartilhado)
│       ├── sdd-verify.md      (compartilhado)
│       ├── sdd-archive.md     (compartilhado)
│       └── sdd-onboard.md     (compartilhado)
├── skills/                (skills SDD existentes)
├── commands/              (comandos slash existentes)
├── plugins/               (plugins existentes)
└── settings/              (configurações)
```

**Ponto chave**: Os prompts do orquestrador NÃO são compartilhados como arquivos externos porque cada perfil precisa de sua própria tabela de atribuições de modelos e referências a subagentes com sufixo. Eles são inseridos em linha no JSON de cada orquestrador durante a geração.

Os prompts dos subagentes SÃO compartilhados como arquivos `{file:...}` porque são idênticos entre perfis — apenas o campo `"model"` muda.

### 6.4 Arquivos Afetados (Mapa de Implementação)

| Área | Arquivo | Mudanças |
|------|------|---------|
| **Modelo de domínio** | `internal/model/types.go` | Adicionar struct `Profile` |
| **Modelo de domínio** | `internal/model/selection.go` | Adicionar `Profiles []Profile` a `Selection` e `SyncOverrides` |
| **TUI: telas** | `internal/tui/screens/profiles.go` | NOVO — tela de lista de perfis (ações de lista + editar + excluir) |
| **TUI: telas** | `internal/tui/screens/profile_create.go` | NOVO — fluxo de criação/edição de perfil (nome → modelos → confirmar) |
| **TUI: telas** | `internal/tui/screens/profile_delete.go` | NOVO — tela de confirmação de exclusão de perfil |
| **TUI: modelo** | `internal/tui/model.go` | Adicionar `ScreenProfiles`, `ScreenProfileCreate`, `ScreenProfileEdit`, `ScreenProfileDelete`, `ScreenProfileResult` |
| **TUI: roteador** | `internal/tui/router.go` | Adicionar rotas para todas as telas de perfil |
| **TUI: boas-vindas** | `internal/tui/screens/welcome.go` | Adicionar opção "OpenCode SDD Profiles" |
| **Injeção SDD** | `internal/components/sdd/inject.go` | Extrair prompts para arquivos, gerar agentes de perfil |
| **Injeção SDD** | `internal/components/sdd/profiles.go` | NOVO — CRUD de perfil: gerar, detectar, excluir agentes do JSON |
| **Injeção SDD** | `internal/components/sdd/prompts.go` | NOVO — gerenciamento de arquivos de prompt compartilhados |
| **Injeção SDD** | `internal/components/sdd/read_assignments.go` | Adicionar detecção de perfil a partir do opencode.json |
| **Sincronização** | `internal/cli/sync.go` | Atualizar sincronização para lidar com perfis, adicionar flag `--profile` |
| **Assets** | `internal/assets/opencode/sdd-overlay-multi.json` | Refatorar para usar referências `{file:...}` |
| **Modelos OpenCode**| `internal/opencode/models.go` | Sem mudanças (reutilizar existente) |

### 6.5 Fluxo de Sincronização (Atualizado)

```
Início da Sincronização
  │
  ├─ 1. Ler opencode.json → detectar perfis existentes
  │     (padrão: sdd-orchestrator-*)
  │
  ├─ 2. Gravar/atualizar arquivos de prompt compartilhados
  │     ~/.config/opencode/prompts/sdd/*.md
  │     (a partir dos assets embutidos)
  │
  ├─ 3. Atualizar orquestrador PADRÃO + subagentes
  │     (sdd-orchestrator, sdd-init, ..., sdd-archive)
  │     - Atualizar prompts (em linha para orquestrador, {file:} para subagentes)
  │     - Preservar atribuições de modelos
  │
  ├─ 4. Para CADA perfil detectado:
  │     ├─ Atualizar prompts de subagentes (usam {file:}, atualizados no passo 2)
  │     ├─ Regenerar prompt do orquestrador (em linha, com a tabela de modelos do perfil)
  │     └─ Preservar atribuições de modelos
  │
  └─ 5. Verificar: todos os orquestradores de perfil + subagentes presentes
```

### 6.6 Caminho de Migração

**Retrocompatibilidade**: Usuários sem perfis não veem mudanças. A refatoração de prompts para arquivos é transparente:

1.  **Primeira sincronização após a atualização**:
    - Cria o diretório `~/.config/opencode/prompts/sdd/`.
    - Grava os arquivos de prompt.
    - Migra subagentes do overlay de prompt em linha para a referência `{file:...}`.
    - Resultado: comportamento idêntico, muda apenas onde o prompt reside.

2.  **Usuários com multi-modo existente**:
    - Suas atribuições de modelos são preservadas.
    - Seus subagentes são migrados para `{file:...}` automaticamente.
    - Zero interrupção.

---

## 7. Fluxo de UX

### 7.1 Tela de Boas-vindas (Atualizada)

```
┌─────────────────────────────────────────────────────────┐
│                                                          │
│  ★  Ecossistema de IA Kortex — v0.x.x                   │
│     Potencialize seus agentes de IA.                     │
│                                                          │
│  ▸ Instalar Ecossistema                                  │
│    Atualizar                                             │
│    Sincronizar                                           │
│    Atualizar + Sincronizar                                │
│    Configurar Modelos                                    │
│    Perfis SDD do OpenCode (2)                    ← NOVO  │
│    Gerenciar Backups                                     │
│    Sair                                                  │
│                                                          │
│  j/k: navegar • enter: selecionar • q: sair             │
└─────────────────────────────────────────────────────────┘
```

### 7.2 Tela de Lista de Perfis

```
┌─────────────────────────────────────────────────────────┐
│  Perfis SDD do OpenCode                                  │
│                                                          │
│  Seus perfis de modelos SDD para o OpenCode. Cada perfil │
│  cria seu próprio orquestrador (visível com Tab).       │
│                                                          │
│  Perfis existentes:                                      │
│    ✦ default ─── anthropic/claude-opus-4                 │
│  ▸   barato ──── anthropic/claude-haiku-3.5              │
│      gemini ──── google/gemini-2.5-pro                   │
│                                                          │
│    Criar novo perfil                                     │
│    Voltar                                                │
│                                                          │
│  j/k: navegar • enter: editar • n: novo • d: excluir     │
│  esc: voltar                                             │
└─────────────────────────────────────────────────────────┘
```

Os perfis são itens navegáveis. O cursor pode estar em um perfil OU em "Criar novo perfil" / "Voltar":
- **enter em um perfil** → modo de edição (modificar modelos e sincronizar)
- **d em um perfil** → confirmação de exclusão (exceto default)
- **enter em "Criar novo perfil"** → fluxo de criação
- **n em qualquer lugar** → atalho para "Criar novo perfil"

### 7.3 Fluxo de Edição de Perfil

Idêntico à criação, mas com valores pré-preenchidos:

```
┌─────────────────────────────────────────────────────────┐
│  Editar Perfil "barato"                                  │
│                                                          │
│  Orquestrador atual: anthropic/claude-haiku-3.5          │
│                                                          │
│  ▸ Alterar modelo do orquestrador                        │
│    Alterar modelos dos subagentes                        │
│    Salvar e Sincronizar                                  │
│    Cancelar                                              │
│                                                          │
│  j/k: navegar • enter: selecionar • esc: cancelar       │
└─────────────────────────────────────────────────────────┘
```

### 7.4 Fluxo de Exclusão de Perfil

```
┌─────────────────────────────────────────────────────────┐
│  Excluir Perfil                                          │
│                                                          │
│  Tem certeza que deseja excluir o perfil "barato"?       │
│                                                          │
│  Isso removerá do opencode.json:                         │
│    • sdd-orchestrator-barato                             │
│    • sdd-init-barato ... sdd-archive-barato              │
│    • (11 agentes no total)                               │
│                                                          │
│  ▸ Excluir e Sincronizar                                 │
│    Cancelar                                              │
│                                                          │
│  enter: selecionar • esc: cancelar                       │
└─────────────────────────────────────────────────────────┘
```

### 7.5 Fluxo de Criação de Perfil

```
Passo 1: Nome
┌─────────────────────────────────────────────────────────┐
│  Criar Perfil SDD                                        │
│                                                          │
│  Nome do perfil: barato_                                 │
│                                                          │
│  (minúsculas, hifens permitidos, sem espaços)           │
│  Reservado: "default"                                    │
│                                                          │
│  enter: confirmar • esc: cancelar                        │
└─────────────────────────────────────────────────────────┘

Passo 2: Modelo do Orquestrador
┌─────────────────────────────────────────────────────────┐
│  Perfil "barato" — Selecionar Modelo do Orquestrador     │
│                                                          │
│  ▸ anthropic                                             │
│    google                                                │
│    openai                                                │
│    Voltar                                                │
│                                                          │
│  (reutiliza o ModelPicker existente)                      │
└─────────────────────────────────────────────────────────┘

Passo 3: Modelos dos Subagentes
┌─────────────────────────────────────────────────────────┐
│  Perfil "barato" — Atribuir Modelos dos Subagentes       │
│                                                          │
│  ▸ Definir todas as fases ── (nenhum)                    │
│    sdd-init ─────────────── (nenhum)                     │
│    sdd-explore ──────────── (nenhum)                     │
│    sdd-propose ──────────── (nenhum)                     │
│    sdd-spec ─────────────── (nenhum)                     │
│    sdd-design ───────────── (nenhum)                     │
│    sdd-tasks ────────────── (nenhum)                     │
│    sdd-apply ────────────── (nenhum)                     │
│    sdd-verify ───────────── (nenhum)                     │
│    sdd-archive ──────────── (nenhum)                     │
│    Continuar                                             │
│    Voltar                                                │
│                                                          │
│  (reutiliza o ModelPicker existente)                      │
└─────────────────────────────────────────────────────────┘

Passo 4: Confirmar + Sincronizar
┌─────────────────────────────────────────────────────────┐
│  Perfil "barato" — Pronto para Criar                     │
│                                                          │
│  Orquestrador: anthropic/claude-haiku-3.5-20241022      │
│  Subagentes:   anthropic/claude-haiku-3.5-20241022 (todos)│
│                                                          │
│  Isso irá:                                               │
│  • Adicionar sdd-orchestrator-barato ao opencode.json    │
│  • Adicionar 10 subagentes (sdd-init-barato ... )        │
│  • Rodar sync para aplicar as mudanças                   │
│                                                          │
│  ▸ Criar e Sincronizar                                   │
│    Cancelar                                              │
│                                                          │
│  enter: selecionar • esc: cancelar                       │
└─────────────────────────────────────────────────────────┘
```

---

## 8. Casos de Borda e Decisões

### 8.1 Cache de Modelos do OpenCode não disponível

Se o `~/.cache/opencode/models.json` não existir (OpenCode nunca executado), a tela de criação de perfil DEVE:
- Mostrar uma mensagem explicativa: "Execute o OpenCode pelo menos uma vez para popular o cache de modelos".
- Oferecer apenas "Voltar".
- NÃO bloquear o restante da TUI.

### 8.2 Validação do Nome do Perfil

| Entrada | Válido? | Motivo |
|-------|--------|--------|
| `barato` | ✓ | Slug simples |
| `premium-v2` | ✓ | Hifens permitidos |
| `meu perfil` | ✗ | Espaços não permitidos |
| `default` | ✗ | Reservado |
| `ALTO` | → `alto` | Convertido automaticamente para minúsculas |
| `sdd-orchestrator` | ✗ | Criaria `sdd-orchestrator-sdd-orchestrator` — confuso |
| `a` | ✓ | Mínimo 1 caractere |
| (vazio) | ✗ | Deve ter um nome |

### 8.3 Herança de Modelo para Subagentes

Quando um subagente não tem uma atribuição de modelo explícita:
1. Usar o modelo do orquestrador do mesmo perfil.
2. Se o modelo do orquestrador não estiver definido, usar o `"model"` raiz do opencode.json.
3. Se nada estiver definido, o OpenCode usa seu padrão.

### 8.4 Excluindo um Perfil

A exclusão é totalmente suportada na TUI (pressionar `d` em um perfil → confirmar → agentes removidos do JSON → sincronizar). A operação:
1. Lê o `opencode.json`.
2. Remove TODAS as chaves que correspondam a `sdd-orchestrator-{nome}` e `sdd-{fase}-{nome}` (11 chaves no total).
3. Grava o JSON atualizado atomicamente.
4. Executa a sincronização para garantir a consistência.
5. O perfil `default` NÃO pode ser excluído — o atalho é ignorado sobre ele.

### 8.5 Prompt do Orquestrador — Referências de Subagentes

O prompt do orquestrador do perfil padrão referencia subagentes como `sdd-apply`. Um perfil "barato" precisa que seu orquestrador referencie `sdd-apply-barato`.

**Solução**: Ao gerar o prompt do orquestrador de um perfil, faz-se a substituição da string do padrão `sdd-{fase}` → `sdd-{fase}-{sufixo}` APENAS dentro das seções que referenciam subagentes (tabela de Model Assignments, regras de delegação). Isso é feito no momento da geração, não no arquivo compartilhado.

---

## 9. Métricas de Sucesso

| Métrica | Meta |
|--------|--------|
| Tempo de criação de perfil (TUI) | < 60 segundos |
| Tempo de sincronização com 3 perfis | < 5 segundos adicionais |
| Zero regressão para usuários sem perfis | 100% retrocompatível |
| Quantidade de perfis suportados | Testado até 10 |
| Arquivos alterados por sincronização (sem mudanças reais) | 0 (idempotente) |

---

## 10. Fases de Implementação

### Fase 1: Refatoração de Prompts Compartilhados (Base)
- Extrair prompts de subagentes para `~/.config/opencode/prompts/sdd/*.md`.
- Atualizar `sdd-overlay-multi.json` para usar referências `{file:...}`.
- Atualizar `inject.go` para gravar os arquivos de prompt.
- Atualizar sincronização para manter os arquivos de prompt.
- **Mudança de comportamento zero** — mesmos prompts, localização diferente.

### Fase 2: Modelo de Dados de Perfil e Geração
- Adicionar o tipo `Profile` ao modelo de domínio.
- Implementar a geração de agentes de perfil (orquestrador + subagentes com sufixo).
- Detecção de perfil a partir do opencode.json existente.
- Atualizar `injectModelAssignments` para lidar com múltiplos perfis.

### Fase 3: Telas da TUI — Criar e Listar
- Tela de lista de perfis (mostra perfis existentes com ações).
- Fluxo de criação de perfil (nome → modelo do orquestrador → modelos dos subagentes → confirmar).
- Integrar na tela de Boas-vindas.
- Integrar com o fluxo de sincronização (auto-sync após criação do perfil).

### Fase 4: Telas da TUI — Editar e Excluir
- Fluxo de edição de perfil (selecionar perfil → modificar modelos → salvar e sincronizar).
- Tela de confirmação de exclusão de perfil + limpeza do JSON.
- Atalho `d` na lista de perfis para exclusão.
- Atalho `enter` no perfil para edição.
- Proteção do perfil padrão (sem exclusão, permite edição).

### Fase 5: Integração da Sincronização
- Atualizar sincronização para detectar e manter todos os perfis.
- Adicionar a flag CLI `--profile`.
- Atualizar alvos de backup para incluir arquivos de prompt.
- Atualizar verificação pós-sincronização para perfis.

### Fase 6: Polimento e Testes
- Testes E2E para criação, edição, exclusão e sincronização de perfis.
- Tratamento de casos de borda (cache faltando, nomes inválidos, etc.).
- Atualização da documentação.

---

## 11. Perguntas em Aberto

1.  **O prompt do orquestrador de cada perfil é inserido em linha no JSON ou salvo como arquivo?**
    → Decisão: EM LINHA no JSON. O prompt do orquestrador é específico do perfil (tabela de modelos + referências de subagentes) e não pode ser compartilhado como arquivo. Os prompts dos subagentes SÃO compartilhados como arquivos.

2.  **O que acontece com o `sdd-onboard` nos perfis?**
    → Decisão: `sdd-onboard-{nome}` é gerado como subagente do perfil, assim como os outros 9 subagentes.

3.  **Os comandos slash do SDD (`/sdd-new`, `/sdd-ff`, etc.) funcionam com perfis personalizados?**
    → Sim. Os comandos estão vinculados ao orquestrador. Quando o usuário seleciona `sdd-orchestrator-barato` com Tab, os comandos são executados contra esse orquestrador, que delega aos subagentes `sdd-*-barato`.

4.  **Como o OpenCode lida com o `{file:...}` nos prompts? Suporta a expansão de `~`?**
    → Validar com os docs do OpenCode. Se não suportar `~`, usar o caminho absoluto expandido durante a geração.

5.  **O agente `carbon` (persona) também precisa de variantes por perfil?**
    → Não. O agente `carbon` é a persona geral, não faz parte do SDD. Apenas o modelo do orquestrador padrão é espelhado.
