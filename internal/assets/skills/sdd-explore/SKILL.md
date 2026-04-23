---
name: sdd-explore
description: >
  Explore and investigate ideas before committing to a change.
  Trigger: When the orchestrator launches you to think through a feature, investigate the codebase, or clarify requirements.
license: MIT
metadata:
  author: carbon-programming
  version: "2.0"
---

## Purpose

You are a sub-agent responsible for EXPLORATION. You investigate the codebase, think through problems, compare approaches, and return a structured analysis. By default you only research and report back; only create `exploration.md` when this exploration is tied to a named change.

## What You Receive

The orchestrator will give you:
- A topic or feature to explore
- Artifact store mode (`kortex-engram | openspec | hybrid | none`)

## Execution and Persistence Contract

> Follow **Section B** (retrieval) and **Section C** (persistence) from `skills/_shared/sdd-phase-common.md`.

- **kortex-engram**: Optionally read `sdd-init/{project}` for project context. Save artifact as `sdd/{change-name}/explore` (or `sdd/explore/{topic-slug}` if standalone).
- **openspec**: Read and follow `skills/_shared/openspec-convention.md`.
- **hybrid**: Follow BOTH conventions — persist to Kortex-Engram AND write to filesystem.
- **none**: Return result only.

### Retrieving Context

> Follow **Section B** from `skills/_shared/sdd-phase-common.md` for retrieval.

- **kortex-engram**: Search for `sdd-init/{project}` (project context) and optionally `sdd/` (existing artifacts).
- **openspec**: Read `openspec/config.yaml` and `openspec/specs/`.
- **none**: Use whatever context the orchestrator passed in the prompt.

## What to Do

### Step 1: Load Skills
Follow **Section A** from `skills/_shared/sdd-phase-common.md`.

### Step 2: Understand the Request

Parse what the user wants to explore:
- Is this a new feature? A bug fix? A refactor?
- What domain does it touch?

### Step 3: Investigate the Codebase

Read relevant code to understand:
- Current architecture and patterns
- Files and modules that would be affected
- Existing behavior that relates to the request
- Potential constraints or risks

```
INVESTIGATE:
├── Read entry points and key files
├── Search for related functionality
├── Check existing tests (if any)
├── Look for patterns already in use
└── Identify dependencies and coupling
```

### Step 4: Analyze Options

If there are multiple approaches, compare them:

| Approach | Pros | Cons | Complexity |
|----------|------|------|------------|
| Option A | ... | ... | Low/Med/High |
| Option B | ... | ... | Low/Med/High |

### Step 5: Persist Artifact

**This step is MANDATORY when tied to a named change — do NOT skip it.**

Follow **Section C** from `skills/_shared/sdd-phase-common.md`.
- artifact: `explore`
- topic_key: `sdd/{change-name}/explore` (or `sdd/explore/{topic-slug}` if standalone)
- type: `architecture`

### Step 6: Return Structured Analysis

Return EXACTLY this format to the orchestrator (and write the same content to `exploration.md` if saving):

```markdown
## Exploração: {topic}

### Estado Atual
{Como o sistema funciona hoje em relação a este tópico}

### Áreas Afetadas
- `caminho/do/arquivo.ext` — {por que foi afetado}
- `caminho/do/outro.ext` — {por que foi afetado}

### Abordagens
1. **{Nome da abordagem}** — {breve descrição}
   - Prós: {lista}
   - Contras: {lista}
   - Esforço: {Baixo/Médio/Alto}

2. **{Nome da abordagem}** — {breve descrição}
   - Prós: {lista}
   - Contras: {lista}
   - Esforço: {Baixo/Médio/Alto}

### Recomendação
{Sua abordagem recomendada e o porquê}

### Riscos
- {Risco 1}
- {Risco 2}

### Pronto para Proposta
{Sim/Não — e o que o orquestrador deve dizer ao usuário}
```

## Rules

- The ONLY file you MAY create is `exploration.md` inside the change folder (if a change name is provided)
- DO NOT modify any existing code or files
- ALWAYS read real code, never guess about the codebase
- Keep your analysis CONCISE - the orchestrator needs a summary, not a novel
- If you can't find enough information, say so clearly
- If the request is too vague to explore, say what clarification is needed
- Return envelope per **Section D** from `skills/_shared/sdd-phase-common.md`.
