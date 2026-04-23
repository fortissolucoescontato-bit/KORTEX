# SDD Phase — Common Protocol

Boilerplate identical across all SDD phase skills. Sub-agents MUST load this alongside their phase-specific SKILL.md.

Executor boundary: every SDD phase agent is an EXECUTOR, not an orchestrator. Do the phase work yourself. Do NOT launch sub-agents, do NOT call `delegate`/`task`, and do NOT bounce work back unless the phase skill explicitly says to stop and report a blocker.

## A. Skill Loading

1. Check if the orchestrator injected a `## Project Standards (auto-resolved)` block in your launch prompt. If yes, follow those rules — they are pre-digested compact rules from the skill registry. **Do NOT read any SKILL.md files.**
2. If no Project Standards block was provided, check for `SKILL: Load` instructions. If present, load those exact skill files.
3. If neither was provided, search for the skill registry as a fallback:
   a. `mem_search(query: "skill-registry", project: "{project}")` — if found, `mem_get_observation(id)` for full content
   b. Fallback: read `.atl/skill-registry.md` from the project root if it exists
   c. From the registry's **Compact Rules** section, apply rules whose triggers match your current task.
4. If no registry exists, proceed with your phase skill only.

NOTE: the preferred path is (1) — compact rules pre-injected by the orchestrator. Paths (2) and (3) are fallbacks for backwards compatibility. Searching the registry is SKILL LOADING, not delegation. If `## Project Standards` is present, IGNORE any `SKILL: Load` instructions — they are redundant.

## B. Artifact Retrieval (Kortex-Engram Mode)

**CRITICAL**: `mem_search` returns 300-char PREVIEWS, not full content. You MUST call `mem_get_observation(id)` for EVERY artifact. **Skipping this produces wrong output.**

**Run all searches in parallel** — do NOT search sequentially.

```
mem_search(query: "sdd/{change-name}/{artifact-type}", project: "{project}") → save ID
```

Then **run all retrievals in parallel**:

```
mem_get_observation(id: {saved_id}) → full content (REQUIRED)
```

Do NOT use search previews as source material.

## C. Artifact Persistence

Every phase that produces an artifact MUST persist it. Skipping this BREAKS the pipeline — downstream phases will not find your output.

### Kortex-Engram mode

```
mem_save(
  title: "sdd/{change-name}/{artifact-type}",
  topic_key: "sdd/{change-name}/{artifact-type}",
  type: "architecture",
  project: "{project}",
  content: "{your full artifact markdown}"
)
```

`topic_key` enables upserts — saving again updates, not duplicates.

### OpenSpec mode

File was already written during the phase's main step. No additional action needed.

### Hybrid mode

Do BOTH: write the file to the filesystem AND call `mem_save` as above.

### None mode

Return result inline only. Do not write any files or call `mem_save`.

## D. Return Envelope

Every phase MUST return a structured envelope to the orchestrator. Localize the display labels to the user's language (default: Portuguese), but keep the intent of the keys for the orchestrator.

Example (Portuguese):

```markdown
**Status**: sucesso
**Resumo**: Proposta criada para `{change-name}`. Escopo, abordagem e plano de rollback definidos.
**Artefatos**: Kortex-Engram `sdd/{change-name}/proposal` | `openspec/changes/{change-name}/proposal.md`
**Próximo**: sdd-spec ou sdd-design
**Riscos**: Nenhum
**Resolução de Skills**: injected — 3 skills (react-19, typescript, tailwind-4)
```

## E. Language Policy

- **Output Language**: ALWAYS respond in the language the user is communicating (default: Portuguese).
- **Artifact Language**: The content of specifications, designs, proposals, and tasks MUST be written in the user's language.
- **Technical Terms**: Keep standard technical terms (e.g., "Middleware", "Frontend", "Endpoint") if they are commonly used in the local tech community, but ensure the surrounding context is localized.
