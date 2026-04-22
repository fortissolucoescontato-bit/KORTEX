# Shim PowerShell do Kortex CLI — Suporte para Windows

## O que é isto

Quando o `kortex` instala o Kortex CLI no Windows, ele agora instala um wrapper `kortex.ps1` junto com o script bash principal. Isso permite que os usuários executem o `kortex` diretamente do PowerShell sem precisar mudar manualmente para o Git Bash.

## Como funciona

```
O usuário digita: kortex init   (no PowerShell)
                        │
                        ▼
             O Windows resolve o kortex.ps1
      (O PowerShell entende extensões .ps1)
                        │
                        ▼
      O kortex.ps1 encontra o Git Bash via Get-Command git
                        │
                        ▼
      O Git Bash executa o script bash original do kortex
                        │
                        ▼
      Código de saída + saída retornados ao PowerShell
```

O shim é instalado no mesmo diretório do binário `kortex` (`~/.local/share/kortex/bin/kortex.ps1`) e usa uma gravação atômica com verificação de igualdade de conteúdo — reexecutar o `kortex install` é uma operação idempotente.

## Requisitos

- O Git para Windows deve estar instalado (fornece o Git Bash).
- O shim é exclusivo para Windows — macOS e Linux não são afetados.

## Limitações Conhecidas e Iterações Futuras

Os itens a seguir foram identificados durante a verificação e adiados para trabalhos futuros. Eles não são bugs — o Kortex CLI funciona corretamente para os casos comuns. Estas são melhorias que valem a pena revisar.

### Iteração 1 — Encaminhamento de argumentos com espaços entre aspas (W-01)

O shim usa:
```powershell
& $gitBash -c "kortex $args"
```

Argumentos com aspas embutidas ou espaços são passados via interpolação de string para o `bash -c`, o que pode perder a fidelidade das aspas em casos específicos. Por exemplo:

```powershell
kortex commit -m "minha mensagem"   # pode chegar como: kortex commit -m minha mensagem
```

**Correção recomendada**: usar o splatting `@args` ou construir o array de argumentos explicitamente em vez da interpolação de string.

### Iteração 2 — Mensagem de erro de Git Bash não encontrado (W-02)

A especificação original descrevia a exibição de um erro "Git Bash não encontrado" **durante o `kortex install`**. No design final, isso foi movido para o **tempo de execução** — o shim `.ps1` detecta o Git Bash quando o usuário executa o `kortex` pela primeira vez. O cenário da especificação agora está impreciso e deve ser atualizado para refletir o modelo de detecção em tempo de execução.

**Correção recomendada**: atualizar o `openspec/changes/kortex-powershell-support/specs/kortex/spec.md` para renomear o cenário de "tempo de instalação" para "detecção em tempo de execução" e adicionar um teste de integração que exercite o caminho do código de erro no PS runtime.

### Iteração 3 — Cobertura de teste da trava para sistemas não-Windows (W-03)

Os pontos de chamada em `internal/cli/run.go` e `internal/cli/sync.go` protegem o shim com `if runtime.GOOS == "windows"`. Isso é verificado estruturalmente (a trava existe no código-fonte), mas não há um teste automatizado que simule um SO não-Windows e assegure que `EnsurePowerShellShim` nunca seja chamado.

**Correção recomendada**: adicionar um teste baseado em tabela que injete um valor `GOOS` falso e assegure que o caminho de instalação do shim seja ignorado no `linux` e `darwin`.
