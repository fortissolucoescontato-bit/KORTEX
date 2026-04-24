---
name: sdd-init
description: >
  Initialize Spec-Driven Development context in any project. Detects stack, conventions, testing capabilities, and bootstraps the active persistence backend.
  Trigger: When user wants to initialize SDD in a project, or says "sdd init", "iniciar sdd", "openspec init".
license: MIT
metadata:
  author: carbon-programming
  version: "3.0"
---

## Purpose

You are a sub-agent responsible for initializing the Spec-Driven Development (SDD) context in a project. You detect the project stack, conventions, and testing capabilities, then bootstrap the active persistence backend.

You are an EXECUTOR for this phase, not the orchestrator. Do the initialization work yourself. Do NOT launch sub-agents, do NOT call `delegate` or `task`, and do NOT hand execution back unless you hit a real blocker that must be reported upstream.

## Execution and Persistence Contract

- If mode is `kortex-engram`:
  Do NOT create `openspec/` directory.

  **Save project context**:
  ```
  mem_save(
    title: "sdd-init/{project-name}",
    topic_key: "sdd-init/{project-name}",
    type: "architecture",
    project: "{project-name}",
    content: "{detected project context markdown}"
  )
  ```
  `topic_key` enables upserts — re-running init updates the existing context, not duplicates.

  (See `skills/_shared/kortex-engram-convention.md` for full naming conventions.)
- If mode is `openspec`: Read and follow `skills/_shared/openspec-convention.md`. Run full bootstrap.
- If mode is `hybrid`: Read and follow BOTH convention files. Run openspec bootstrap AND persist context to Kortex-Engram.
- If mode is `none`: Return detected context without writing project files.

## What to Do

### Step 1: Detect Project Context

Read the project to understand:
- Tech stack (check package.json, go.mod, pyproject.toml, etc.)
- Existing conventions (linters, test frameworks, CI)
- Architecture patterns in use

### Step 2: Detect Testing Capabilities

Scan the project for ALL testing infrastructure. This determines what testing modes are available.

```
Detect testing capabilities:
├── Test Runner
│   ├── package.json → devDependencies: vitest, jest, mocha, ava
│   ├── package.json → scripts.test (what command it runs)
│   ├── pyproject.toml / pytest.ini / setup.cfg → pytest
│   ├── go.mod → go test (built-in)
│   ├── Cargo.toml → cargo test (built-in)
│   ├── Makefile → make test
│   └── Result: {framework name, command} or NOT FOUND
│
├── Test Layers
│   ├── Unit: test runner exists → AVAILABLE
│   ├── Integration:
│   │   ├── JS/TS: @testing-library/* in dependencies
│   │   ├── Python: pytest + httpx/requests-mock/factory-boy
│   │   ├── Go: net/http/httptest (built-in)
│   │   ├── .NET: xUnit/NUnit + WebApplicationFactory
│   │   └── Result: AVAILABLE or NOT INSTALLED
│   ├── E2E:
│   │   ├── playwright, cypress, selenium in dependencies
│   │   ├── Python: playwright, selenium
│   │   ├── Go: chromedp
│   │   └── Result: AVAILABLE or NOT INSTALLED
│   └── Each layer → record tool name
│
├── Coverage Tool
│   ├── JS/TS: vitest --coverage, jest --coverage, c8, istanbul/nyc
│   ├── Python: coverage.py, pytest-cov
│   ├── Go: go test -cover (built-in)
│   ├── .NET: coverlet
│   └── Result: {command} or NOT AVAILABLE
│
└── Quality Tools
    ├── Linter: eslint, pylint, ruff, golangci-lint, clippy
    ├── Type checker: tsc --noEmit, mypy, pyright, go vet
    ├── Formatter: prettier, black, gofmt, rustfmt
    └── Each: {command} or NOT AVAILABLE
```

### Step 3: Resolve STRICT TDD MODE

Determine whether Strict TDD Mode should be enabled. The resolution follows a priority chain — first match wins:

```
1. Read from system prompt / agent config (highest priority):
   ├── Search for "strict-tdd-mode" marker in the agent's system prompt file
   │   (e.g., CLAUDE.md, GEMINI.md, .cursorrules, etc.)
   ├── If found and says "enabled" → strict_tdd: true
   ├── If found and says "disabled" → strict_tdd: false
   └── This is the preference set by the user in the kortex TUI

2. If no marker found, check openspec config:
   ├── Read openspec/config.yaml → strict_tdd field
   └── If found → use that value

3. If nothing found AND test runner was detected in Step 2:
   ├── Default: strict_tdd: true (enable if the project CAN do TDD)
   └── This ensures TDD is active even without kortex TUI setup

4. If no test runner detected:
   ├── strict_tdd: false (cannot enable without test runner)
   └── Include NOTE in summary: "Strict TDD Mode unavailable — no test runner detected"
```

**Do NOT ask the user interactively.** The preference is resolved from existing config. If the user wants to change it, they run `kortex sync` with the TUI or set `strict_tdd` in `openspec/config.yaml`.

### Step 4: Initialize Persistence Backend

If mode resolves to `openspec`, create this directory structure:

```
openspec/
├── config.yaml              ← Project-specific SDD config
├── specs/                   ← Source of truth (empty initially)
└── changes/                 ← Active changes
    └── archive/             ← Completed changes
```

### Step 5: Generate Config (openspec mode)

Based on what you detected, create the config when in `openspec` mode:

```yaml
# openspec/config.yaml
schema: spec-driven

context: |
  Tech stack: {detected stack}
  Architecture: {detected patterns}
  Testing: {detected test framework}
  Style: {detected linting/formatting}

strict_tdd: {true/false}

rules:
  proposal:
    - Include rollback plan for risky changes
    - Identify affected modules/packages
  specs:
    - Use Given/When/Then format for scenarios
    - Use RFC 2119 keywords (MUST, SHALL, SHOULD, MAY)
  design:
    - Include sequence diagrams for complex flows
    - Document architecture decisions with rationale
  tasks:
    - Group tasks by phase (infrastructure, implementation, testing)
    - Use hierarchical numbering (1.1, 1.2, etc.)
    - Keep tasks small enough to complete in one session
  apply:
    - Follow existing code patterns and conventions
    - Load relevant coding skills for the project stack
  verify:
    - Run tests if test infrastructure exists
    - Compare implementation against every spec scenario
  archive:
    - Warn before merging destructive deltas (large removals)
```

### Step 6: Persist Testing Capabilities

**This step is MANDATORY — do NOT skip it.**

Persist detected testing capabilities as a separate Kortex-Engram observation (or section in config.yaml for openspec). This cache prevents re-detection on every `sdd-apply` and `sdd-verify` run.

If mode is `kortex-engram` or `hybrid`:
```
mem_save(
  title: "sdd/{project-name}/testing-capabilities",
  topic_key: "sdd/{project-name}/testing-capabilities",
  type: "config",
  project: "{project-name}",
  content: "{testing capabilities markdown — see format below}"
)
```

**Testing Capabilities format**:

```markdown
## Testing Capabilities

**Strict TDD Mode**: {enabled/disabled}
**Detected**: {date}

### Test Runner
- Command: `{command}`
- Framework: {name}

### Test Layers
| Layer | Available | Tool |
|-------|-----------|------|
| Unit | ✅ / ❌ | {tool or —} |
| Integration | ✅ / ❌ | {tool or —} |
| E2E | ✅ / ❌ | {tool or —} |

### Coverage
- Available: ✅ / ❌
- Command: `{command or —}`

### Quality Tools
| Tool | Available | Command |
|------|-----------|---------|
| Linter | ✅ / ❌ | {command or —} |
| Type checker | ✅ / ❌ | {command or —} |
| Formatter | ✅ / ❌ | {command or —} |
```

If mode is `openspec` or `hybrid`, also write this as a section in `openspec/config.yaml` under `testing:`.

### Step 7: Build Skill Registry

Follow the same logic as the `skill-registry` skill (`skills/skill-registry/SKILL.md`):

1. Scan user skills: glob `*/SKILL.md` across ALL known skill directories. **User-level**: `~/.claude/skills/`, `~/.config/opencode/skills/`, `~/.gemini/skills/`, `~/.cursor/skills/`, `~/.copilot/skills/`, parent of this skill file. **Project-level**: `.claude/skills/`, `.gemini/skills/`, `.agent/skills/`, `skills/`. Skip `sdd-*`, `_shared`, `skill-registry`. Deduplicate by name (project-level wins). Read frontmatter triggers.
2. Scan project conventions: check for `agents.md`, `AGENTS.md`, `CLAUDE.md` (project-level), `.cursorrules`, `GEMINI.md`, `copilot-instructions.md` in the project root. If an index file is found (e.g., `agents.md`), READ it and extract all referenced file paths — include both the index and its referenced files in the registry.
3. **ALWAYS write `.atl/skill-registry.md`** in the project root (create `.atl/` if needed). This file is mode-independent — it's infrastructure, not an SDD artifact.
4. If kortex-engram is available, **ALSO save to kortex-engram**: `mem_save(title: "skill-registry", topic_key: "skill-registry", type: "config", project: "{project}", content: "{registry markdown}")`

See `skills/skill-registry/SKILL.md` for the full registry format and scanning details.

### Step 8: Persist Project Context

**This step is MANDATORY — do NOT skip it.**

If mode is `kortex-engram`:
```
mem_save(
  title: "sdd-init/{project-name}",
  topic_key: "sdd-init/{project-name}",
  type: "architecture",
  project: "{project-name}",
  content: "{your detected project context from Steps 1-7}"
)
```

If mode is `openspec` or `hybrid`: the config was already written in Step 5.

If mode is `hybrid`: also call `mem_save` as above (write to BOTH backends).

### Step 9: Return Summary

Return a structured summary adapted to the resolved mode:

#### If mode is `kortex-engram`:

Persist project context following `skills/_shared/kortex-engram-convention.md` with title and topic_key `sdd-init/{project-name}`.

Return:
```
## SDD Inicializado

**Projeto**: {project name}
**Stack**: {detected stack}
**Persistência**: kortex-engram
**Modo TDD Estrito**: {ativado ✅ / desativado ❌ / indisponível (sem runner de testes)}

### Capacidades de Teste
| Capacidade | Status |
|------------|--------|
| Executor de Testes | {tool} ✅ / ❌ Não encontrado |
| Testes Unitários | ✅ / ❌ |
| Testes de Integração | {tool} ✅ / ❌ Não instalado |
| Testes E2E | {tool} ✅ / ❌ Não instalado |
| Cobertura | ✅ / ❌ |
| Linter | {tool} ✅ / ❌ |
| Verificador de Tipos | {tool} ✅ / ❌ |

### Contexto Salvo
O contexto do projeto foi persistido no Kortex-Engram.
- **Kortex-Engram ID**: #{observation-id}
- **Chave do Tópico**: sdd-init/{project-name}
- **ID de Capacidades**: #{capabilities-observation-id}
- **Chave de Capacidades**: sdd/{project-name}/testing-capabilities

Nenhum arquivo de projeto foi criado.

### ⚠️ Notas do Modo Kortex-Engram
O modo Kortex-Engram é ideal para **desenvolvedores solo** que buscam iteração rápida. Esteja ciente:
- **Sem histórico de iteração**: rodar uma fase novamente (ex: `sdd-spec`) sobrescreve a versão anterior. Apenas o artefato mais recente é mantido.
- **Não compartilhável**: o kortex-engram é um banco de dados local — outros membros da equipe não verão seus artefatos SDD.
- **Trilha de auditoria parcial**: a fase de arquivamento salva um relatório de resumo, mas não a pasta completa de artefatos.

Para **projetos em equipe** ou trabalhos que exigem trilha completa de auditoria, considere mudar para `openspec` (baseado em arquivos, amigável ao git) ou `hybrid` (arquivos + recuperação via kortex-engram).

### Próximos Passos
Pronto para /sdd-explore <tópico> ou /sdd-new <nome-da-mudança>.
```

#### If mode is `openspec`:
```
## SDD Inicializado

**Projeto**: {project name}
**Stack**: {detected stack}
**Persistência**: openspec
**Modo TDD Estrito**: {ativado ✅ / desativado ❌ / indisponível (sem runner de testes)}

### Capacidades de Teste
{mesma tabela acima}

### Estrutura Criada
- openspec/config.yaml ← Configuração do projeto com contexto detectado + capacidades de teste
- openspec/specs/      ← Pronto para especificações
- openspec/changes/    ← Pronto para propostas de mudança

### Próximos Passos
Pronto para /sdd-explore <tópico> ou /sdd-new <nome-da-mudança>.
```

#### If mode is `none`:
```
## SDD Inicializado

**Projeto**: {project name}
**Stack**: {detected stack}
**Persistência**: nenhuma (efêmera)
**Modo TDD Estrito**: {ativado ✅ / desativado ❌ / indisponível (sem runner de testes)}

### Capacidades de Teste
{mesma tabela acima}

### Contexto Detectado
{resumo da stack e convenções detectadas}

### Recomendação
Ative o `kortex-engram` ou `openspec` para persistência de artefatos entre sessões. Sem persistência, todos os artefatos SDD serão perdidos quando a conversa terminar.

### Próximos Passos
Pronto para /sdd-explore <tópico> ou /sdd-new <nome-da-mudança>.
```

## Rules

- NEVER create placeholder spec files - specs are created via sdd-spec during a change
- ALWAYS detect the real tech stack, don't guess
- NEVER behave like the orchestrator from this phase - execute directly and return results
- If the project already has an `openspec/` directory, report what exists and ask the orchestrator if it should be updated
- Keep config.yaml context CONCISE - no more than 10 lines
- ALWAYS detect testing capabilities — this is not optional
- ALWAYS persist testing capabilities as a separate observation/section — downstream phases depend on it
- If Strict TDD Mode is requested but no test runner exists, set strict_tdd: false and explain why
- Return a structured envelope with: `status`, `executive_summary`, `detailed_report` (optional), `artifacts`, `next_recommended`, and `risks`
