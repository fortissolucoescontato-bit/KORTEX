# Componentes, Skills e Presets

← [Voltar para o README](../README.md)

---

## Componentes

| Componente | ID | Descrição |
|-----------|-----|-------------|
| Kortex-Engram | `kortex-engram` | Memória persistente entre sessões via MCP — detecção automática de nome de projeto, busca em texto completo, sincronização git, consolidação de projeto. Veja o [repositório do kortex-engram](https://github.com/fortissolucoescontato-bit/kortex-engram) |
| SDD | `sdd` | Fluxo de trabalho de Desenvolvimento Orientado a Especificações (9 fases) — o agente lida com o SDD organicamente quando a tarefa exige, ou quando você solicita; não é necessário aprender os comandos |
| Skills | `skills` | Biblioteca curada de habilidades de codificação |
| Context7 | `context7` | Servidor MCP para documentação ao vivo de frameworks/bibliotecas |
| Persona | `persona` | Injeção gerenciada da persona Kortex/neutra, ou modo de persona personalizada não gerenciada |
| Permissões | `permissions` | Padrões e proteções focados em segurança |
| Kortex CLI | `kortex` | Kortex Guardian Angel — alternador de provedores de IA |
| Tema | `theme` | Sobreposição do tema Kortex Kanagawa |

## Comportamento do Kortex CLI

`kortex --component kortex` instala/provisiona o binário `kortex` globalmente em sua máquina.

Ele **não** executa a configuração automática de hooks no nível do projeto (`kortex init` / `kortex install`) porque essa deve ser uma decisão explícita por repositório.

Após a instalação global, ative o Kortex CLI por projeto com:

```bash
kortex init
kortex install
```

---

## Skills

### Skills Incluídas (instaladas pelo kortex)

14 arquivos de skill organizados por categoria, embutidos no binário e injetados na configuração do seu agente:

#### SDD (Desenvolvimento Orientado a Especificações)

| Skill | ID | Descrição |
|-------|-----|-------------|
| SDD Init | `sdd-init` | Inicializa o contexto SDD em um projeto |
| SDD Explore | `sdd-explore` | Investiga a base de código antes de se comprometer com uma mudança |
| SDD Propose | `sdd-propose` | Cria proposta de mudança com intenção, escopo e abordagem |
| SDD Spec | `sdd-spec` | Escreve especificações com requisitos e cenários |
| SDD Design | `sdd-design` | Design técnico com decisões de arquitetura |
| SDD Tasks | `sdd-tasks` | Divide uma mudança em tarefas de implementação |
| SDD Apply | `sdd-apply` | Implementa as tarefas seguindo as specs e o design |
| SDD Verify | `sdd-verify` | Valida se a implementação corresponde às specs |
| SDD Archive | `sdd-archive` | Sincroniza as specs delta com as specs principais e arquiva |
| Judgment Day | `judgment-day` | Revisão adversária paralela — dois juízes independentes revisam o mesmo alvo |

#### Base

| Skill | ID | Descrição |
|-------|-----|-------------|
| Go Testing | `go-testing` | Padrões de teste em Go, incluindo testes de TUI Bubbletea |
| Skill Creator | `skill-creator` | Cria novas skills para agentes de IA seguindo a especificação de Agent Skills |
| Branch & PR | `branch-pr` | Fluxo de criação de PR com commits convencionais, nomenclatura de branch e exigência de issue primeiro |
| Criação de Issues| `issue-creation` | Fluxo de abertura de issues com templates de bug report e solicitações de recursos |

Essas skills de base são instaladas por padrão tanto nos presets `full-carbon` quanto no `ecosystem-only`.

### Skills de Codificação (repositório separado)

Para skills específicas de frameworks (React 19, Angular, TypeScript, Tailwind 4, Zod 4, Playwright, etc.), consulte [fortissolucoescontato-bit/Kortex-Skills](https://github.com/fortissolucoescontato-bit/Kortex-Skills). Estas são mantidas pela comunidade e instaladas separadamente clonando o repositório e copiando as skills para o diretório de skills do seu agente.

---

## Presets

| Preset | ID | O Que Está Incluído |
|--------|-----|-------------------|
| Kortex Completo | `full-carbon` | Todos os componentes (Kortex-Engram + SDD + Skills + Context7 + Kortex CLI + Persona + Permissões + Tema) + todas as skills + persona carbon |
| Apenas Ecossistema| `ecosystem-only` | Componentes principais (Kortex-Engram + SDD + Skills + Context7 + Kortex CLI) + todas as skills + persona carbon |
| Mínimo | `minimal` | Apenas skills de Kortex-Engram + SDD |
| Personalizado | `custom` | Você escolhe os componentes e skills manualmente, mantendo qualquer persona/configuração existente não gerenciada |
