# Delta for kortex

## ADDED Requirements

### Requirement: Integridade de Identificadores E2E

A infraestrutura de testes E2E DEVE utilizar identificadores compatíveis com a sintaxe POSIX Shell/Bash em todas as suas variáveis locais e globais. Nomes de variáveis NÃO DEVEM conter hífens ou caracteres especiais além de sublinhados (`_`).

#### Scenario: Validação de Ordem de Componentes no E2E

- DADO o script de teste `e2e/e2e_test.sh`
- QUANDO a função `test_dry_run_full_preset_persona_before_sdd` é executada
- ENTÃO a variável que rastreia o índice do `kortex-engram` DEVE ser nomeada `kortex_engram_idx` ou similar sem hífens
- E o Bash deve interpretar a declaração `local` sem erros de "not a valid identifier"

## MODIFIED Requirements

### Requirement: Sistemas não-Windows não afetados

Em plataformas que não sejam Windows, o instalador NÃO DEVE tentar gravar o `kortex.ps1` ou invocar a etapa do shim do PowerShell.
(Previously: Requirement describing non-Windows systems are unaffected by Windows-specific installation steps.)

#### Scenario: Fluxo de instalação no Linux/macOS inalterado

- DADO um host Linux ou macOS
- QUANDO a instalação da CLI Kortex é executada
- ENTÃO nenhum arquivo `.ps1` é criado e nenhum caminho de código específico do Windows é executado
- AND a suíte de testes E2E valida esse isolamento em containers Docker (Ubuntu, Arch, Fedora)
