---
name: kortex-issue-creation
description: >
  Fluxo de trabalho para criação de issues no Kortex seguindo o sistema de obrigatoriedade de issue-first.
  Trigger: Ao criar uma issue no GitHub, reportar um bug ou solicitar uma funcionalidade.
license: Apache-2.0
metadata:
  author: nexo-fortis
  version: "1.0"
---

# Kortex — Skill de Criação de Issue

## Quando Usar

Carregue esta skill sempre que precisar:
- Reportar um bug no `kortex`
- Solicitar uma nova funcionalidade ou melhoria
- Abrir qualquer issue no repositório [fortissolucoescontato-bit/kortex](https://github.com/fortissolucoescontato-bit/kortex)

## Regras Críticas

1. **Issues em branco estão DESATIVADAS** — Você DEVE usar um template.
2. **`status:needs-review` é aplicado automaticamente** — Toda nova issue recebe esta etiqueta; você NÃO a adiciona manualmente.
3. **`status:approved` é OBRIGATÓRIO antes de iniciar QUALQUER trabalho** — Um mantenedor deve aprovar a issue antes que você ou qualquer pessoa abra um PR.
4. **Dúvidas vão para Discussions** — Use o [GitHub Discussions](https://github.com/fortissolucoescontato-bit/kortex/discussions), NÃO as issues, para perguntas e conversas gerais.
5. **Sem trailers `Co-Authored-By`** — Nunca adicione atribuição de IA aos commits.

## Fluxo de Trabalho

```
1. Pesquise issues existentes → confirme que não é duplicada
   https://github.com/fortissolucoescontato-bit/kortex/issues

2. Escolha o template correto:
   - Bug   → .github/ISSUE_TEMPLATE/bug_report.yml
   - Feat  → .github/ISSUE_TEMPLATE/feature_request.yml

3. Envie a issue → status:needs-review é aplicado automaticamente

4. Aguarde — um mantenedor revisará e adicionará status:approved (ou fechará)

5. Somente APÓS o status:approved → abra um PR referenciando esta issue
```

> ⚠️ **PARE após o passo 3.** NÃO abra um PR até que a issue tenha o `status:approved`.

---

## Relato de Bug (Bug Report)

**Caminho do Template**: `.github/ISSUE_TEMPLATE/bug_report.yml`
**Etiquetas Automáticas**: `bug`, `status:needs-review`

### Campos Obrigatórios

| Campo | Descrição |
|-------|-------------|
| Pre-flight Checklist | Confirme que não existe duplicata; confirme o entendimento da aprovação do PR |
| Bug Description | Descrição clara do que é o bug |
| Steps to Reproduce | Passos numerados para reproduzir o comportamento |
| Expected Behavior | O que deveria acontecer |
| Actual Behavior | O que realmente acontece |
| Kortex Version | Saída do comando `kortex version` |
| Operating System | macOS / Distro Linux / Windows / WSL |
| AI Agent / Client | Claude Code / OpenCode / Gemini CLI / Cursor / Windsurf / Outro |
| Affected Area | Veja a lista de áreas afetadas abaixo |

### Áreas Afetadas

`CLI (comandos, flags)` · `TUI (interface de terminal)` · `Pipeline de Instalação` · `Detecção de Agentes` · `Detecção de Sistema` · `Catálogo/Etapas` · `Documentação` · `Outro`

### Exemplo de Comando CLI

```bash
gh issue create \
  --repo fortissolucoescontato-bit/kortex \
  --template bug_report.yml \
  --title "fix(agent): Claude Code não detectado no Linux Arch"
```

Ou abra o formulário web diretamente:
```
https://github.com/fortissolucoescontato-bit/kortex/issues/new?template=bug_report.yml
```

---

## Solicitação de Funcionalidade (Feature Request)

**Caminho do Template**: `.github/ISSUE_TEMPLATE/feature_request.yml`
**Etiquetas Automáticas**: `enhancement`, `status:needs-review`

### Campos Obrigatórios

| Campo | Descrição |
|-------|-------------|
| Pre-flight Checklist | Confirme que não existe duplicata; confirme o entendimento da aprovação do PR |
| Affected Area | Qual área do `kortex` esta funcionalidade afeta |
| Problem Statement | Descreva o problema que esta funcionalidade resolve |
| Proposed Solution | Descrição específica — inclua exemplo de comando/saída do `kortex` se relevante |
| Alternatives Considered | (opcional) Outras abordagens que você considerou |
| Additional Context | (opcional) Screenshots, arquivos de config, etc. |

### Exemplo de Comando CLI

```bash
gh issue create \
  --repo fortissolucoescontato-bit/kortex \
  --template feature_request.yml \
  --title "feat(tui): adicionar atalho de teclado para ajuda"
```

Ou abra o formulário web diretamente:
```
https://github.com/fortissolucoescontato-bit/kortex/issues/new?template=feature_request.yml
```

---

## Sistema de Etiquetas (Labels)

### Etiquetas de Status (aplicadas às Issues)

| Etiqueta | Descrição | Quem Aplica |
|-------|-------------|-------------|
| `status:needs-review` | Recém-aberta, aguardando revisão do mantenedor | **Auto** (template) |
| `status:approved` | Aprovada — o trabalho pode começar | Apenas Mantenedor |
| `status:in-progress` | Sendo trabalhada ativamente | Contribuidor |
| `status:blocked` | Bloqueada por outra issue ou dependência externa | Mantenedor / Contribuidor |
| `status:wont-fix` | Fora de escopo ou não será tratada | Apenas Mantenedor |

### Etiquetas de Tipo (aplicadas às Issues e PRs)

| Etiqueta | Descrição |
|-------|-------------|
| `bug` | Relato de defeito |
| `enhancement` | Solicitação de funcionalidade ou melhoria |
| `type:bug` | Correção de bug (usada em PRs) |
| `type:feature` | Nova funcionalidade (usada em PRs) |
| `type:docs` | Apenas documentação (usada em PRs) |
| `type:refactor` | Refatoração, sem mudanças funcionais (usada em PRs) |
| `type:chore` | Build, CI, ferramentas (usada em PRs) |
| `type:breaking-change` | Mudança de ruptura (usada em PRs) |

---

## Fluxo de Aprovação do Mantenedor

```
Issue enviada
      │
      ▼
status:needs-review  ← auto-aplicado pelo template
      │
      ▼
Mantenedor revisa
      │
  ┌───┴────────────────┐
  │                    │
  ▼                    ▼
status:approved    Fechada
(trabalho começa)  (inválida / duplicada / wont-fix)
      │
      ▼
Contribuidor comenta "Eu vou trabalhar nisso"
      │
      ▼
status:in-progress
      │
      ▼
PR aberto com `Closes #<N>`
```

---

## Árvore de Decisão

```
Você tem uma dúvida ou ideia para discutir?
├── SIM → GitHub Discussions (NÃO issues)
│         https://github.com/fortissolucoescontato-bit/kortex/discussions
└── NÃO → É um defeito no kortex?
          ├── SIM → Template de Bug Report
          └── NÃO → Template de Feature Request
                    │
                    ▼
          Já existe uma issue similar?
          ├── SIM → Comente na issue existente em vez de criar uma nova
          └── NÃO → Envie a nova issue → aguarde pelo status:approved
```

---

## Comandos

### Buscar por Issues Existentes

```bash
# Buscar issues abertas
gh issue list --repo fortissolucoescontato-bit/kortex --state open --search "suas palavras-chave"

# Buscar todas as issues incluindo as fechadas
gh issue list --repo fortissolucoescontato-bit/kortex --state all --search "suas palavras-chave"
```

### Criar um Bug Report

```bash
gh issue create \
  --repo fortissolucoescontato-bit/kortex \
  --template bug_report.yml \
  --title "fix(<escopo>): <descrição curta>"
```

### Criar um Feature Request

```bash
gh issue create \
  --repo fortissolucoescontato-bit/kortex \
  --template feature_request.yml \
  --title "feat(<escopo>): <descrição curta>"
```

### Verificar Status da Issue

```bash
gh issue view <numero> --repo fortissolucoescontato-bit/kortex
```

### Escopos Válidos para Títulos de Issue

`tui`, `cli`, `installer`, `catalog`, `system`, `agent`, `e2e`, `ci`, `docs`
