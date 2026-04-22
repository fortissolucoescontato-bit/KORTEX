# Modo Não Interativo

Use o modo não interativo para CI, scripts ou configurações locais reproduzíveis.

## Comando

```bash
go run ./cmd/kortex install [flags]
```

## Flags suportadas

- `--agent`, `--agents`: separados por vírgula e repetíveis.
- `--component`, `--components`: separados por vírgula e repetíveis.
- `--skill`, `--skills`: separados por vírgula e repetíveis.
- `--persona`: ID de persona explícito.
- `--preset`: ID de preset explícito.
- `--dry-run`: renderiza o plano sem executar as alterações.

## Comportamento da plataforma

O instalador detecta a plataforma automaticamente em tempo de execução — não há flag para sobrescrever a seleção da plataforma. O perfil de plataforma detectado determina qual gerenciador de pacotes é usado para os comandos de instalação:

| Plataforma | Gerenciador de pacotes | Exemplo de comando de instalação |
|---|---|---|
| macOS | `brew` | `brew install anomalyco/tap/opencode` |
| Ubuntu/Debian | `apt` | `sudo npm install -g opencode-ai` |
| Arch | `pacman` | `sudo npm install -g opencode-ai` |
| Família Fedora/RHEL | `dnf` | `sudo npm install -g opencode-ai` |

A saída do `--dry-run` inclui uma linha de `Decisão da plataforma` mostrando `os`, `distro`, `package-manager` e `status`.

## Exemplos

macOS (ou qualquer plataforma suportada — as mesmas flags, a plataforma é autodetectada):

```bash
go run ./cmd/kortex install \
  --agent claude-code,opencode \
  --component engram,sdd,skills \
  --skill sdd-apply \
  --persona carbon \
  --preset full-carbon \
  --dry-run
```

As flags são idênticas em todas as plataformas. Apenas os comandos de instalação resolvidos mudam com base na detecção.

## Tratamento de erros

- Opções desconhecidas ou não suportadas falham rapidamente com erros de validação.
- A execução em uma plataforma não suportada encerra imediatamente antes de qualquer trabalho de instalação começar.
