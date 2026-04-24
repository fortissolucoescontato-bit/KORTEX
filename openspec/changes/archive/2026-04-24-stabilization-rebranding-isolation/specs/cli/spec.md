# Delta for cli

## MODIFIED Requirements

### Requirement: Injeção Segura e Isolada

A injeção do comportamento mockado NÃO DEVE expor APIs inseguras para a runtime de produção. As variáveis ou funções de override DEVEM ser visíveis apenas para testes, utilizando o padrão de *Functional Overrides* (variáveis globais exportadas iniciadas por `Verify...Override` ou `Default...Stat`). O sistema DEVE garantir que os testes não dependam de permissões de execução em diretórios globais (como `/tmp`) fora do controle do Go Test Runner.
(Previously: Requirement focusing on context-based or closure-based injection for production safety.)

#### Scenario: Execução Normal (Produção)

- DADO que a CLI Kortex está rodando em ambiente de produção
- QUANDO o `RunChecks` é acionado
- ENTÃO as funções físicas `exec.LookPath`, `os.Stat` e `http.Get` são usadas para validar o binário e a saúde

#### Scenario: Setup e Teardown Seguros (Testes)

- DADO um teste individual de integração (`TestRunInstall*`)
- QUANDO o teste é iniciado, ele configura o `t.Cleanup` para restaurar o estado das validações
- ENTÃO as chamadas de sistema são interceptadas pelos mocks
- AND o teste NÃO falha com `permission denied` ao tentar executar binários de verificação, pois os mocks evitam o `fork/exec` real de componentes externos

### Requirement: Mockabilidade das Verificações KortexEngram

O sistema DEVE permitir que o pacote de testes (`cli_test` ou similar) instrua o componente KortexEngram a ignorar as validações reais do SO e da Rede, retornando um falso-positivo (sucesso) ou um erro controlado, garantindo pureza de teste. O isolamento DEVE incluir a capacidade de mockar o `os.Stat` do diretório de configuração do agente.
(Previously: Basic requirement for mocking SO/Network validations in KortexEngram.)

#### Scenario: Teste em ambiente limpo sem KortexEngram global

- DADO que o runner de testes não possui o binário `kortex-engram` no `$PATH`
- E nenhuma porta 7437 está aberta
- QUANDO o teste de integração executa o fluxo `RunInstall` em modo mock
- ENTÃO o `kortexengram.VerifyInstalled()` e `VerifyHealth()` DEVEM retornar `nil` (sucesso)
- AND o `os.Stat` do diretório `~/.kortex-engram` DEVE ser mockado para retornar sucesso sem acessar o disco real
- E a suíte de testes passa sem depender do estado global
