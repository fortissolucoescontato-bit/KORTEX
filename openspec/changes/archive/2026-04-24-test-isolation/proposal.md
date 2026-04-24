# Proposta: Test Isolation for KortexEngram Verify Checks

## Intenção
Garantir que os testes de integração e e2e do pacote `internal/cli` operem de forma 100% pura e isolada, sem "vazar" para o ambiente real do host (evitando test pollution). Os testes não devem falhar se o KortexEngram não estiver instalado no ambiente que executa os testes (ex: containers de CI recém-inicializados).

## Escopo
**O que entra:**
- Modificação no pacote `kortexengram` para expor ganchos (hooks) ou setters que permitam mockar o comportamento do `VerifyInstalled`, `VerifyVersion` e `VerifyHealth` durante testes.
- Atualização nos testes de integração do `internal/cli` (`run_integration_test.go`, `run_kortexengram_download_test.go`) para invocar os mocks e blindar os testes do ambiente do SO.
- Criação de interface abstrata (ou variável injetável) para as funções de sistema e rede consumidas pelas verificações.

**O que sai:**
- Qualquer alteração na lógica real de pós-instalação para usuários (`RunChecks` continuará validando estritamente em produção).
- Modificação nos testes unitários isolados de outros pacotes.

## Abordagem Técnica
1. **Pacote `kortex-engram`**:
   - Expor as dependências do Verify (como `lookPath`, e o cliente/requisição HTTP do `VerifyHealth`) através de variáveis ou de uma função `SetMockMode(bool)`.
   - Alternativa melhor: Adicionar variáveis de exportação seletiva (ex: `var VerifyInstalledFunc = defaultVerifyInstalled`) ou uma função `MockVerify()` que sobrescreve o comportamento padrão, para que os testes do pacote de cima (`cli`) possam intervir.

2. **Pacote `cli` (Testes de Integração)**:
   - Em todos os `TestRunInstall*`, adicionar no bloco `t.Cleanup(...)` a restauração do estado do mock do KortexEngram.
   - Forçar as verificações a retornarem sucesso simulado durante o assert final, ou suprimir a chamada `VerifyChecks` em modo `DryRun` / Integração simulada.

## Riscos
- Introduzir acoplamento entre pacotes só por causa de testes (precisamos manter as exportações isoladas ou utilizar tags de build `//go:build !test` se necessário, embora `internal/` já proteja contra uso externo).
- Mascarar falhas reais na inicialização do serviço caso os mocks se tornem muito abrangentes.
