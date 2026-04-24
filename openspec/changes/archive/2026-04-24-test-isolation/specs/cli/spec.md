# Especificação: Isolamento de Testes CLI (Verify)

## Propósito

Define os requisitos de isolamento para os testes de integração do pacote `internal/cli`, garantindo que as validações de integridade pós-instalação (`VerifyInstalled`, `VerifyHealth`) não consumam estado vazado da máquina hospedeira (host) e operem 100% de forma determinística no diretório de instalação virtual do teste.

## Requisitos

### Requisito: Mockabilidade das Verificações KortexEngram

O sistema DEVE permitir que o pacote de testes (`cli_test` ou similar) instrua o componente KortexEngram a ignorar as validações reais do SO e da Rede, retornando um falso-positivo (sucesso) ou um erro controlado, garantindo pureza de teste.

#### Cenário: Teste em ambiente limpo sem KortexEngram global

- DADO que o runner de testes não possui o binário `kortex-engram` no `$PATH`
- E nenhuma porta 7437 está aberta
- QUANDO o teste de integração executa o fluxo `RunInstall` em modo mock
- ENTÃO o `kortexengram.VerifyInstalled()` e `VerifyHealth()` DEVEM retornar `nil` (sucesso)
- E a suíte de testes passa sem depender do estado global

#### Cenário: Prevenção de Falsos Positivos de Test Pollution

- DADO que o runner de testes POSSUI o binário KortexEngram rodando globalmente
- QUANDO o teste de integração falha internamente antes do mock
- ENTÃO a verificação real NÃO DEVE mascara a falha usando a instância global, pois o mock isolado já interceptou a chamada

### Requisito: Injeção Segura e Isolada

A injeção do comportamento mockado NÃO DEVE expor APIs inseguras para a runtime de produção. As variáveis ou funções de override DEVEM ser visíveis apenas para testes, preferencialmente isoladas por meio de injeção direta no `context.Context` ou via fechamentos controlados nas declarações de verificação.

#### Cenário: Execução Normal (Produção)

- DADO que a CLI Kortex está rodando em ambiente de produção
- QUANDO o `RunChecks` é acionado
- ENTÃO as funções físicas `exec.LookPath` e `http.Get` são usadas para validar o binário e a saúde

#### Cenário: Setup e Teardown Seguros (Testes)

- DADO um teste individual de integração (`TestRunInstall*`)
- QUANDO o teste é iniciado, ele configura o `t.Cleanup` para restaurar o estado das validações
- ENTÃO nenhum teste que ocorra subsequentemente no mesmo processo sofrerá interferência de estados de mocks vazados
