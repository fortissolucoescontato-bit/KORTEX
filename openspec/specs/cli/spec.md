# Especificação: Isolamento de Testes CLI (Verify)

## Propósito

Define os requisitos de isolamento para os testes de integração do pacote `internal/cli`, garantindo que as validações de integridade pós-instalação (`VerifyInstalled`, `VerifyHealth`) não consumam estado vazado da máquina hospedeira (host) e operem 100% de forma determinística no diretório de instalação virtual do teste.

## Requisitos

### Requisito: Mockabilidade das Verificações KortexEngram

O sistema DEVE permitir que o pacote de testes (`cli_test` ou similar) instrua o componente KortexEngram a ignorar as validações reais do SO e da Rede, retornando um falso-positivo (sucesso) ou um erro controlado, garantindo pureza de teste. O isolamento DEVE incluir a capacidade de mockar o `os.Stat` do diretório de configuração do agente.

#### Cenário: Teste em ambiente limpo sem KortexEngram global

- DADO que o runner de testes não possui o binário `kortex-engram` no `$PATH`
- E nenhuma porta 7437 está aberta
- QUANDO o teste de integração executa o fluxo `RunInstall` em modo mock
- ENTÃO o `kortexengram.VerifyInstalled()` e `VerifyHealth()` DEVEM retornar `nil` (sucesso)
- E o `os.Stat` do diretório `~/.kortex-engram` DEVE ser mockado para retornar sucesso sem acessar o disco real
- E a suíte de testes passa sem depender do estado global

### Requisito: Injeção Segura e Isolada

A injeção do comportamento mockado NÃO DEVE expor APIs inseguras para a runtime de produção. As variáveis ou funções de override DEVEM ser visíveis apenas para testes, utilizando o padrão de *Functional Overrides* (variáveis globais exportadas iniciadas por `Verify...Override` ou `Default...Stat`). O sistema DEVE garantir que os testes não dependam de permissões de execução em diretórios globais (como `/tmp`) fora do controle do Go Test Runner.

#### Cenário: Execução Normal (Produção)

- DADO que a CLI Kortex está rodando em ambiente de produção
- QUANDO o `RunChecks` é acionado
- ENTÃO as funções físicas `exec.LookPath`, `os.Stat` e `http.Get` são usadas para validar o binário e a saúde

#### Cenário: Setup e Teardown Seguros (Testes)

- DADO um teste individual de integração (`TestRunInstall*`)
- QUANDO o teste é iniciado, ele configura o `t.Cleanup` para restaurar o estado das validações
- ENTÃO as chamadas de sistema são interceptadas pelos mocks
- E o teste NÃO falha com `permission denied` ao tentar executar binários de verificação, pois os mocks evitam o `fork/exec` real de componentes externos
