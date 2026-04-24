# Proposal: Estabilização de Testes e Correção de Rebranding

## Intent

Esta proposta visa restaurar a integridade do pipeline de testes do KORTEX após o rebranding sistemático de `engram` para `kortex-engram`. Atualmente, o projeto apresenta falhas críticas de sintaxe no Bash (E2E) e falhas de permissão no Go Test (`permission denied`), indicando um acoplamento indesejado com o ambiente host.

## Scope

### In Scope
- Correção do identificador Bash inválido em `e2e/e2e_test.sh`.
- Implementação de injeção de dependência (DI) para `os.Stat` e `exec.LookPath` nos pacotes falhos.
- Garantia de isolamento total dos testes de integração em relação ao estado do host.
- Restauração do status "GREEN" em todas as plataformas de E2E (Ubuntu, Arch, Fedora).

### Out of Scope
- Mudanças funcionais nas lógicas de instalação ou gerenciamento de memória.
- Adição de novos agentes ou componentes.

## Capabilities

### New Capabilities
- None

### Modified Capabilities
- `kortex`: Correção da infraestrutura de validação E2E para refletir o novo branding.
- `cli`: Garantia de que comandos de instalação sejam testáveis em ambientes restritos (CI/CD).

## Approach

1. **Correção Mecânica (Bash)**: Renomear a variável `kortex_kortex-engram_idx` para `kortex_engram_idx` (removendo o hífen) em `e2e/e2e_test.sh`.
2. **Isolamento de Testes (Go)**:
   - Identificar pontos de `fork/exec` ou acesso a arquivos globais que causam `permission denied`.
   - Utilizar o padrão de *Functional Overrides* (variáveis globais exportadas no pacote, sobrescritas em `func init()` dos arquivos `_test.go`) para interceptar chamadas ao sistema.
   - Migrar caminhos temporários remanescentes para `t.TempDir()`.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `e2e/e2e_test.sh` | Modified | Correção de sintaxe Bash. |
| `internal/components/kortex-engram/` | Modified | Implementação de DI para isolamento de testes. |
| `internal/cli/` | Modified | Ajuste nos testes de integração de instalação. |
| `internal/app/` | Modified | Correção de falhas de permissão em testes de inicialização. |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Injeção de DI afetar produção | Low | Usar variáveis com defaults de produção e sobrescrever apenas em arquivos de teste. |
| Regressão em outras plataformas E2E | Med | Rodar a suíte completa via Docker (ubuntu, arch, fedora) antes do merge. |

## Rollback Plan

Reverter as alterações nos arquivos `.go` e `.sh` via Git. O impacto é nulo para o usuário final, pois afeta apenas a infraestrutura de testes.

## Dependencies

- Docker (para validação dos testes E2E).

## Success Criteria

- [ ] Todos os testes unitários e de integração (`go test ./...`) passam sem erros de permissão.
- [ ] O script `./e2e/e2e_test.sh` executa sem erros de sintaxe.
- [ ] O resumo do E2E Docker mostra PASSED: 3 / 3 (ubuntu, arch, fedora).
