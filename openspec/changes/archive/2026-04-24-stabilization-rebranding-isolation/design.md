# Design: Estabilização e Isolamento de Testes

## Technical Approach

A estratégia técnica foca em duas frentes:
1. **Remediação Imediata**: Correção da sintaxe Bash no script de E2E para restaurar a integridade do pipeline.
2. **Isolamento de Testes (Arquitetura)**: Implementação do padrão *Functional Overrides* para desacoplar a lógica de negócio do KORTEX das chamadas físicas ao sistema operacional (Filesystem, PATH, Execução). Isso permitirá que os testes rodem de forma determinística, mesmo em ambientes com restrições de permissão ou ausência de binários externos.

## Architecture Decisions

### Decision: Padrão de Injeção para Isolamento
**Choice**: Functional Overrides (Package-level variables).
**Alternatives considered**: Interfaces/Mocking frameworks (excessivo para este projeto), Dependency Injection via Context (complexo para funções utilitárias).
**Rationale**: Mantém o código de produção limpo e performático, permitindo que arquivos `_test.go` substituam as implementações de sistema em tempo de teste de forma trivial.

### Decision: Sintaxe de Variáveis no E2E
**Choice**: Underscore naming convention (`kortex_engram_idx`).
**Alternatives considered**: CamelCase (menos idiomático em Bash), quoting (complexo de manter).
**Rationale**: Total conformidade com o padrão POSIX/Bash para identificadores válidos em comandos `local`.

## Data Flow

Os testes não devem mais tocar o sistema de arquivos real ou invocar binários do host sem controle:

    Teste (Go/Bash) ──→ Mocks/Overrides ──→ Resposta Determinística
         │                                       │
         └──────── (Interrupção do fluxo real) ──┘

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `e2e/e2e_test.sh` | Modify | Renomear `kortex_kortex-engram_idx` para `kortex_engram_idx`. |
| `internal/system/os_overrides.go` | Create | Definição de `Stat`, `LookPath`, `ReadFile` e `UserHomeDir` como variáveis injetáveis. |
| `internal/system/detect.go` | Modify | Substituir chamadas diretas a `os` e `exec` pelos overrides definidos. |
| `internal/components/kortex-engram/verify.go` | Modify | Adicionar `execCommand` como override para `exec.Command`. |
| `internal/components/kortex-engram/verify_test.go` | Modify | Implementar mocks para `VerifyVersion` usando o novo override. |

## Interfaces / Contracts

```go
// No pacote internal/system
var (
    Stat        = os.Stat
    LookPath    = exec.LookPath
    ReadFile    = os.ReadFile
    UserHomeDir = os.UserHomeDir
    Command     = exec.Command
)
```

## Testing Strategy

| Layer | What to Test | Approach |
|-------|-------------|----------|
| Unit | Detecção de Distro/Tools | Sobrescrever `ReadFile` e `LookPath` para simular diferentes sistemas. |
| Integration | Fluxo de Instalação | Mockar `VerifyHealth` e `VerifyInstalled` para garantir que o instalador prossiga sem binários reais. |
| E2E | Ordem de Componentes | Execução do script corrigido via Docker em 3 plataformas. |

## Migration / Rollout

No migration required. Esta é uma mudança de infraestrutura de desenvolvimento.

## Open Questions

- [ ] A falha de `permission denied` nos testes Go em `/tmp` pode persistir se for uma restrição do host sobre o binário de teste compilado. Se isso ocorrer, será necessário investigar o uso de `GOTMPDIR` apontando para um diretório com permissões de execução.
