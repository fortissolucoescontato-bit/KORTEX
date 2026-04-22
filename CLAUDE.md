<!-- kortex:engram-protocol -->
## Engram Persistent Memory — Protocol

You have access to Engram, a persistent memory system that survives across sessions and compactions.
This protocol is MANDATORY and ALWAYS ACTIVE.

### PROACTIVE SAVE TRIGGERS (mandatory)
Call `mem_save` IMMEDIATELY after:
- Architecture or design decisions.
- Bug fixes (include root cause).
- New patterns or conventions established.
- Feature implementations with non-obvious approaches.

### SESSION START (MANDATORY)
At the beginning of EVERY session, call `mem_search(query: "session-summary", project: "kortex", limit: 3)` to recover context.

### SESSION CLOSE PROTOCOL (mandatory)
Before ending, call `mem_session_summary` with Goal, Discoveries, Accomplished, and Next Steps.
<!-- /kortex:engram-protocol -->

<!-- kortex:project-standards -->
## Padrões do Projeto Kortex

- **Idioma**: Toda a interface (TUI) e documentação principal devem ser em **Português (PT-BR)**.
- **Identidade**: O ícone de seleção padrão é o cérebro (`🧠`) e a marca é **KORTEX**.
- **Localização**: Ao adicionar novas telas ou mensagens, siga os padrões de tradução estabelecidos em `internal/tui/screens/`.
- **Skills**: As instruções dos agentes devem priorizar a clareza e didática para usuários que estão aprendendo a programar com IA.
<!-- /kortex:project-standards -->

<!-- kortex:sdd-orchestrator -->
## SDD Workflow (Spec-Driven Development)
Este projeto utiliza o fluxo SDD para mudanças substanciais.
- Use `/sdd-init` para inicializar.
- Siga o fluxo: Proposta -> Specs -> Design -> Tasks -> Apply -> Verify.
- Priorize a persistência `hybrid` (arquivos + engram).
<!-- /kortex:sdd-orchestrator -->
