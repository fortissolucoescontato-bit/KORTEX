# Início Rápido (Quickstart)

## Pré-requisitos

### macOS

- Homebrew instalado e disponível no PATH.
- `git` disponível.

### Ubuntu/Debian (e derivados como Linux Mint, Pop!_OS)

- `apt-get` disponível (padrão nestas distros).
- Acesso ao `sudo` para instalação de pacotes.
- `git` disponível.

### Arch Linux (e derivados como Manjaro, EndeavourOS)

- `pacman` disponível (padrão nestas distros).
- Acesso ao `sudo` para instalação de pacotes.
- `git` disponível.

### Família Fedora / RHEL (Fedora, CentOS Stream, Rocky Linux, AlmaLinux)

- `dnf` disponível (padrão nestas distros).
- Acesso ao `sudo` para instalação de pacotes.
- `git` disponível.
- Instalações de Node.js usam o setup NodeSource LTS + `dnf install -y nodejs` durante a remediação de dependências.

### Todas as plataformas

- Go 1.24+ (para compilação a partir do código-fonte).
- Node.js / npm se for instalar o Claude Code (o agente é instalado via `npm install -g`).

## Execução

```bash
go run ./cmd/kortex install --dry-run
```

Use `--dry-run` primeiro para validar as seleções e o plano de execução sem aplicar alterações. A saída da simulação inclui uma linha de `Decisão da plataforma` mostrando o SO detectado, a distro, o gerenciador de pacotes e o status de suporte.

## Primeira instalação real

```bash
go run ./cmd/kortex install
```

O instalador detecta sua plataforma automaticamente — não são necessárias flags para selecionar macOS vs Linux. Os comandos de instalação são resolvidos através do gerenciador de pacotes apropriado (brew, apt, pacman ou dnf) com base na detecção.

Após a conclusão, verifique se as configurações do agente e os componentes selecionados foram instalados em seus caminhos esperados.

## Resultado da verificação

Quando as verificações passam, o instalador relata:

`Você está pronto. Execute 'claude' ou 'opencode' e comece a construir. ⚡`

## Plataformas não suportadas

Se você executar o instalador em um SO ou distro Linux não suportado, ele encerrará imediatamente com um erro:

- `sistema operacional não suportado: apenas macOS, Linux e Windows são suportados (detectado <os>)`
- `distro linux não suportada: o suporte ao Linux é limitado a Ubuntu/Debian, Arch e família Fedora/RHEL (detectado <distro>)`
