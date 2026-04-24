# Tasks: test-isolation

## Phase 1: Foundation / Hooks

- [x] 1.1 Modificar `internal/components/kortex-engram/verify.go`: Exportar variáveis globais ou setters para substituição das dependências de validação (ex: `VerifyInstalledOverride func() error`, `VerifyHealthOverride func(context.Context, string) error`).
- [x] 1.2 Atualizar chamadas em `VerifyInstalled`, `VerifyVersion` e `VerifyHealth` para checarem se os overrides estão presentes. Se estiverem, delegar execução para eles; caso contrário, seguir com o comportamento padrão real.

## Phase 2: Test Isolation (cli package)

- [x] 2.1 Atualizar `internal/cli/run_integration_test.go`: No bloco de setup global dos testes ou em helpers (`t.Cleanup`), injetar os mocks em `kortexengram.VerifyInstalledOverride` e `VerifyHealthOverride` para sempre retornarem sucesso simulado (`nil`).
- [x] 2.2 Atualizar `internal/cli/run_kortexengram_download_test.go`: Aplicar o mesmo padrão de isolamento, garantindo restauração limpa ao final de cada teste.

## Phase 3: Verification

- [x] 3.1 Executar `go test ./internal/cli/...` com ambiente local "limpo" (KortexEngram físico desinstalado ou fora do PATH) e confirmar que os testes passam 100%.
- [x] 3.2 Executar `go vet ./...` para garantir integridade estrutural.
