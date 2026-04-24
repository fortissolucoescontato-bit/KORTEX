# Tasks: Estabilização e Isolamento de Testes

## Phase 1: Foundation & Quick Fixes

- [x] 1.1 Corrigir sintaxe Bash em `e2e/e2e_test.sh`: renomear `kortex_kortex-engram_idx` para `kortex_engram_idx` (linhas 247, 249, 252, 253).
- [x] 1.2 Criar `internal/system/os_overrides.go`: definir variáveis globais para `Stat`, `LookPath`, `ReadFile`, `UserHomeDir` e `Command`.
- [x] 1.3 Adicionar `execCommand` como variável injetável em `internal/components/kortex-engram/verify.go` (unexported, mapeada para `exec.Command`).

## Phase 2: Core Implementation (Isolation)

- [x] 2.1 Refatorar `internal/system/detect.go`: substituir `os.UserHomeDir`, `os.ReadFile` e `exec.Command` pelos overrides do pacote `system`.
- [x] 2.2 Refatorar `internal/components/kortex-engram/verify.go`: garantir que `VerifyVersion` utilize `execCommand` para execução do binário.
- [x] 2.3 Refatorar acoplamentos em `internal/cli` e `internal/app` detectados durante a fase de Design.

## Phase 3: Testing & Verification

- [x] 3.1 Atualizar `internal/system/detect_test.go`: injetar mocks de `ReadFile` para simular `/etc/os-release` sem acesso ao disco.
- [x] 3.2 Atualizar `internal/components/kortex-engram/verify_test.go`: adicionar teste para `VerifyVersion` mockando `execCommand`.
- [x] 3.3 Executar `go test ./...`: validar que o erro de `permission denied` desapareceu com o uso de mocks.
- [x] 3.4 Executar `./e2e/e2e_test.sh`: validar que o script de teste E2E não apresenta mais erros de "not a valid identifier".
- [x] 3.5 Executar suite completa E2E via Docker (`RUN_FULL_E2E=1 ./e2e/docker-test.sh`): validar Ubuntu, Arch e Fedora.

## Phase 4: Cleanup & Documentation

- [x] 4.1 Remover quaisquer comentários de debug ou logs temporários adicionados durante a remediação.
- [x] 4.2 Atualizar `CLAUDE.md` ou documentação técnica se houver novos padrões de teste estabelecidos.
