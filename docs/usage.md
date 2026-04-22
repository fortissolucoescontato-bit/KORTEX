# Uso

← [Voltar para o README](../README.md)

---

## Modos de Persona

| Persona | ID | Descrição |
|---------|-----|-------------|
| Nexo-Fortis | `carbon` | Persona de mentor orientada ao ensino — corrige más práticas e explica o "porquê" |
| Neutra | `neutral` | Mesma filosofia de ensino, sem linguagem regional — calorosa e profissional |
| Personalizada| `custom` | Mantém sua persona/configuração atual não gerenciada — o kortex não injeta uma persona |

`custom` é uma escolha de compatibilidade/propriedade, não um editor de persona. Use-o quando você já tiver suas próprias instruções de persona e quiser que o kortex não as altere.

---

## TUI Interativa

Basta executar — a TUI (Interface de Terminal) guia você pela seleção de agentes, componentes, skills, presets e fluxos de desinstalação gerenciada:

```bash
kortex
```

O fluxo de desinstalação também está disponível no menu da TUI. Ele permite que você:

- selecione um ou mais agentes configurados
- selecione quais componentes gerenciados remover (por exemplo `sdd`, `persona` ou `context7`)
- confirme o escopo exato da desinstalação antes de aplicar as alterações

Antes de qualquer arquivo gerenciado ser modificado, o `kortex` cria um snapshot de backup para que a configuração possa ser restaurada posteriormente, se necessário.

---

## Comandos CLI

### install (instalar)

Configuração inicial — detecta suas ferramentas, configura agentes e injeta todos os componentes:

```bash
# Ecossistema completo para múltiplos agentes
kortex install \
  --agent claude-code,opencode,gemini-cli \
  --preset full-carbon

# Configuração mínima para o Cursor
kortex install \
  --agent cursor \
  --preset minimal

# Escolher componentes e skills específicos
kortex install \
  --agent claude-code \
  --component engram,sdd,skills,context7,persona,permissions \
  --skill go-testing,skill-creator,branch-pr,issue-creation \
  --persona carbon

# Simulação primeiro (visualizar o plano sem aplicar alterações)
kortex install --dry-run \
  --agent claude-code,opencode \
  --preset full-carbon
```

### sync (sincronizar)

Atualiza os assets gerenciados para a versão atual. Use após `brew upgrade kortex` ou quando quiser que suas configurações locais estejam alinhadas com o último lançamento. NÃO reinstala binários (engram, Kortex CLI) — apenas atualiza o conteúdo de prompts, skills, configurações MCP e orquestradores SDD.

```bash
# Sincronizar todos os agentes instalados
kortex sync

# Sincronizar apenas agentes específicos
kortex sync --agent cursor --agent windsurf

# Sincronizar um componente específico
kortex sync --component sdd
kortex sync --component skills
kortex sync --component engram
```

O sync é seguro e idempotente — executá-lo duas vezes não produz alterações na segunda vez.

### uninstall (desinstalar)

Remove apenas a configuração gerenciada pelo `kortex` de um ou mais agentes. Isso não desinstala pacotes externos ou binários — remove seções de prompt gerenciadas, entradas MCP, fragmentos de skills/configuração e outros arquivos gerenciados, atualizando o `state.json` adequadamente.

Antes de qualquer alteração ser aplicada, o `kortex` cria um snapshot de backup dos arquivos afetados.

```bash
# Desinstalação parcial para agentes específicos
kortex uninstall \
  --agent claude-code \
  --agent opencode

# Desinstalação parcial apenas para componentes específicos
kortex uninstall \
  --agent claude-code \
  --component sdd,persona,context7

# Desinstalação completa da configuração gerenciada de todos os agentes suportados
kortex uninstall --all

# Pular a confirmação
kortex uninstall --agent cursor --component skills --yes
```

Se nenhuma flag `--component` for fornecida para uma desinstalação parcial, o `kortex` remove todos os componentes desinstaláveis gerenciados para o conjunto de agentes selecionado.

### update / upgrade (atualizar)

Verifica e instala novas versões do próprio `kortex`:

```bash
# Verificar se uma versão mais recente está disponível
kortex update

# Atualizar para a versão mais recente (baixa o novo binário e substitui o atual)
kortex upgrade
```

Após o upgrade, execute `kortex sync` para atualizar todos os assets gerenciados para o conteúdo da nova versão.

### version (versão)

```bash
kortex version
kortex --version
kortex -v
```

---

## Flags da CLI (install)

| Flag | Descrição |
|------|-------------|
| `--agent`, `--agents` | Agentes a configurar (separados por vírgula) |
| `--component`, `--components` | Componentes a instalar (separados por vírgula) |
| `--skill`, `--skills` | Skills a instalar (separados por vírgula) |
| `--persona` | Modo de persona: `carbon`, `neutral`, `custom` (`custom` mantém sua persona atual não gerenciada) |
| `--preset` | Preset: `full-carbon`, `ecosystem-only`, `minimal`, `custom` (`custom` permite seleção manual de componentes/skills) |
| `--dry-run` | Visualiza o plano de instalação sem aplicar alterações |

## Flags da CLI (sync)

| Flag | Descrição |
|------|-------------|
| `--agent`, `--agents` | Agentes a sincronizar (padrão: todos os agentes instalados) |
| `--component` | Sincroniza apenas um componente específico: `sdd`, `engram`, `context7`, `skills`, `kortex`, `permissions`, `theme` |
| `--profile` | Cria ou atualiza um perfil SDD: `nome:provedor/modelo` (define o modelo padrão para todas as fases) |
| `--profile-phase` | Sobrescreve uma fase específica em um perfil: `nome:fase:provedor/modelo` |
| `--sdd-profile-strategy` | Estratégia de sincronização de perfis OpenCode: `generated-multi` ou `external-single-active` |
| `--include-permissions` | Inclui sincronização de permissões (opcional) |
| `--include-theme` | Inclui sincronização de tema (opcional) |

**Exemplos de perfis:**

```bash
# Criar um perfil "econômico" usando um modelo gratuito para todas as fases
kortex sync --profile economico:openrouter/qwen/qwen3-30b-a3b:free

# Sobrescrever a fase de design para usar um modelo mais forte
kortex sync --profile-phase economico:sdd-design:anthropic/claude-sonnet-4-20250514

# Criar múltiplos perfis em um comando
kortex sync \
  --profile economico:openrouter/qwen/qwen3-30b-a3b:free \
  --profile premium:anthropic/claude-sonnet-4-20250514

# Usar modo de compatibilidade com um gerenciador de perfis externo do OpenCode
kortex sync --agent opencode --sdd-profile-strategy external-single-active
```

Consulte [Perfis SDD do OpenCode](opencode-profiles.md) para o guia completo.

## Flags da CLI (uninstall)

| Flag | Descrição |
|------|-------------|
| `--agent`, `--agents` | Agentes dos quais remover a configuração gerenciada (obrigatório, a menos que use `--all`) |
| `--component`, `--components` | Componentes gerenciados a remover apenas dos agentes selecionados |
| `--all` | Remove a configuração gerenciada de todos os agentes suportados |
| `--yes`, `-y` | Pula a confirmação |

---

## Fluxo de Trabalho Típico

```bash
# Primeira vez: instalar tudo
brew install carbon-programming/tap/kortex
kortex install --agent claude-code,cursor --preset full-carbon

# Após um novo lançamento: upgrade + sync
brew upgrade kortex
kortex sync

# Remover apenas as configurações gerenciadas de SDD + persona de um agente
kortex uninstall --agent claude-code --component sdd,persona

# Adicionando um novo agente depois
kortex install --agent windsurf --preset full-carbon
```

---

## Gerenciamento de Dependências

O `kortex` detecta automaticamente os pré-requisitos antes da instalação e fornece orientação específica para cada plataforma:

- **Ferramentas detectadas**: git, curl, node, npm, brew, go
- **Verificações de versão**: valida versões mínimas quando aplicável
- **Dicas baseadas na plataforma**: sugere `brew install`, `apt install`, `pacman -S`, `dnf install` ou `winget install` dependendo do seu SO
- **Alinhamento Node LTS**: em sistemas apt/dnf, as dicas de Node.js usam o bootstrap do NodeSource LTS antes da instalação do pacote
- **Abordagem "dependência primeiro"**: detecta o que está instalado, calcula o que é necessário, mostra a árvore completa de dependências antes de instalar qualquer coisa e, em seguida, verifica cada dependência após a instalação
