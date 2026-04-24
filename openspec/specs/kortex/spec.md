# Especificação da CLI Kortex

## Propósito

Define o comportamento de instalação e execução do componente CLI Kortex, cobrindo o asset shim do PowerShell e sua etapa de instalação específica para Windows.

## Requisitos

### Requisito: Integridade de Identificadores E2E

A infraestrutura de testes E2E DEVE utilizar identificadores compatíveis com a sintaxe POSIX Shell/Bash em todas as suas variáveis locais e globais. Nomes de variáveis NÃO DEVEM conter hífens ou caracteres especiais além de sublinhados (`_`).

#### Cenário: Validação de Ordem de Componentes no E2E

- DADO o script de teste `e2e/e2e_test.sh`
- QUANDO a função `test_dry_run_full_preset_persona_before_sdd` é executada
- ENTÃO a variável que rastreia o índice do `kortex-engram` DEVE ser nomeada `kortex_engram_idx` ou similar sem hífens
- E o Bash deve interpretar a declaração `local` sem erros de "not a valid identifier"

---

### Requisito: Asset Shim do PowerShell

O sistema DEVE embutir um arquivo `kortex.ps1` como um asset Go em `internal/assets/kortex/`. O shim DEVE delegar a execução ao binário Git Bash resolvido por `gitBashPath()`, encaminhando todos os argumentos literalmente e propagando o código de saída.

#### Cenário: Shim delega para o Git Bash

- DADO que o `kortex.ps1` embutido está instalado em uma máquina Windows com Git Bash presente
- QUANDO o usuário executa `kortex <subcomando>` a partir do PowerShell
- ENTÃO o shim invoca o Git Bash com o caminho do binário bash resolvido e todos os argumentos fornecidos
- E o processo termina com o mesmo código retornado pelo comando bash da CLI Kortex subjacente

#### Cenário: Argumentos contendo espaços são encaminhados corretamente

- DADO que o `kortex.ps1` está instalado
- QUANDO o usuário executa `kortex commit -m "minha mensagem"` a partir do PowerShell
- ENTÃO o argumento `"minha mensagem"` chega à CLI Kortex como um único token (não dividido)

#### Cenário: Propagação de código de saída em caso de erro

- DADO que o `kortex.ps1` está instalado
- QUANDO o comando da CLI Kortex subjacente termina com um código diferente de zero
- ENTÃO o `$LASTEXITCODE` do PowerShell reflete exatamente esse valor diferente de zero

---

### Requisito: Etapa de Instalação no Windows

No Windows, o instalador DEVE gravar o `kortex.ps1` no mesmo diretório que o script bash da CLI Kortex após a conclusão do `install.sh` próprio da CLI Kortex. A gravação DEVE usar um padrão atômico de no-op: se o arquivo já existir com conteúdo idêntico, o instalador NÃO DEVE sobrescrevê-lo.

#### Cenário: Primeira instalação no Windows

- DADO que a CLI Kortex concluiu sua própria instalação
- E o `kortex.ps1` ainda não existe no diretório de instalação
- QUANDO a etapa de instalação do Windows é executada
- ENTÃO o `kortex.ps1` é gravado no diretório de instalação com o conteúdo correto

#### Cenário: Reinstalação idempotente (conteúdo inalterado)

- DADO que o `kortex.ps1` já existe com conteúdo correspondente ao asset embutido atual
- QUANDO o instalador é executado novamente
- ENTÃO o arquivo NÃO é sobrescrito (não ocorre E/S de gravação)

#### Cenário: Shim desatualizado é atualizado

- DADO que o `kortex.ps1` existe, mas seu conteúdo difere do asset embutido atual
- QUANDO o instalador é executado
- ENTÃO o arquivo é substituído atomicamente pelo novo conteúdo

#### Cenário: Git Bash não encontrado no momento da instalação

- DADO que o Git Bash não está instalado na máquina Windows de destino
- QUANDO a etapa de instalação tenta resolver `gitBashPath()`
- ENTÃO o instalador apresenta uma mensagem de erro clara e acionável
- E a instalação é interrompida sem gravar um shim corrompido

---

### Requisito: Sistemas não-Windows não afetados

Em plataformas que não sejam Windows, o instalador NÃO DEVE tentar gravar o `kortex.ps1` ou invocar a etapa do shim do PowerShell.

#### Cenário: Fluxo de instalação no Linux/macOS inalterado

- DADO um host Linux ou macOS
- QUANDO a instalação da CLI Kortex é executada
- ENTÃO nenhum arquivo `.ps1` é criado e nenhum caminho de código específico do Windows é executado
- E a suíte de testes E2E valida esse isolamento em containers Docker (Ubuntu, Arch, Fedora)

---

## Nota de Documentação

`docs/platforms.md` DEVE remover qualquer nota de limitação do Windows que afirme que o PowerShell não é suportado assim que esta alteração for lançada. Esta é uma atualização apenas de documentação, sem requisito de comportamento além de manter o documento preciso.
