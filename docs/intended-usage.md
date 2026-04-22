# Uso Pretendido

← [Voltar para o README](../README.md)

---

Esta página explica como o kortex deve ser usado. Não se trata de flags ou arquitetura — trata-se do modelo mental. Se você ler apenas uma página além do README, que seja esta.

---

## Depois de Instalar — Você Está Pronto

Depois de executar o `kortex` e selecionar seus agentes, componentes e preset, tudo estará configurado. Não há mais nada a fazer. Sem comandos para memorizar, sem fluxos de trabalho para aprender, sem arquivos de configuração para editar.

Abra seu agente de IA e comece a trabalhar. Só isso.

---

## Engram (Memória) — Automático, mas Você PODE Usar

O Engram é a memória persistente para o seu agente de IA. Ele salva decisões, descobertas, correções de bugs e contexto entre as sessões — automaticamente. O agente gerencia tudo via ferramentas MCP (`mem_save`, `mem_search`, etc.).

**No dia a dia: você não precisa fazer nada.** O agente lida com a memória automaticamente.

**Mas o Engram tem ferramentas úteis quando você precisar:**

| Comando | Quando usar |
|---------|-------------|
| `engram tui` | Navegue por suas memórias visualmente — busque, filtre e detalhe observações |
| `engram sync` | Exporta as memórias do projeto para `.engram/` para rastreamento via git. Execute após sessões de trabalho significativas |
| `engram sync --import` | Importa memórias em outra máquina após clonar um repositório com `.engram/` |
| `engram projects list` | Veja todos os projetos com contagem de observações |
| `engram projects consolidate` | Corrige divergências de nomes de projetos (ex: "meu-app" vs "Meu-App" vs "meu-app-frontend") |
| `engram search <termo>` | Busca rápida de memória pelo terminal |

Desde a v1.11.0, o Engram detecta automaticamente o nome do projeto a partir do remoto git na inicialização, normaliza para minúsculas e avisa se encontrar nomes de projetos existentes semelhantes. Isso evita que o mesmo projeto acabe com múltiplas variantes de nome.

Para documentação completa: [github.com/fortissolucoescontato-bit/engram](https://github.com/fortissolucoescontato-bit/engram)

---

## SDD (Desenvolvimento Orientado a Especificações) — Acontece Organicamente

O SDD é um fluxo de trabalho de planejamento estruturado para recursos substanciais. Ele possui fases (explore, propose, spec, design, implement, verify), mas você NÃO precisa aprender nenhuma delas.

Veja como funciona na prática:

- **Solicitação pequena?** O agente simplesmente executa. Sem cerimônia.
- **Recurso substancial?** O agente sugerirá o uso do SDD para planejá-lo adequadamente — explorando a base de código, propondo uma abordagem, desenhando a arquitetura e implementando passo a passo.
- **Quer o SDD explicitamente?** Apenas diga "use sdd" ou "hazlo con sdd" e o agente inicia o fluxo de trabalho.

O agente gerencia todas as fases internamente. Você apenas revisa e aprova nos pontos-chave de decisão.

Se você quiser a convenção de configuração OpenSpec no nível do projeto que as fases do SDD usam para padrões, TDD estrito e metadados de teste, veja [Configuração OpenSpec para SDD](openspec-config.md).

---

## SDD Multi-modo (Perfis SDD do OpenCode)

O multi-modo permite atribuir diferentes modelos de IA a diferentes fases do SDD — por exemplo, um modelo poderoso para o design e um mais rápido para a implementação. Este é um recurso exclusivo do OpenCode, gerenciado através de **Perfis SDD**.

Para **todos os outros agentes** (Claude Code, Cursor, Gemini CLI, VS Code Copilot), o SDD roda em modo único automaticamente. Um único modelo lida com tudo, e isso funciona perfeitamente bem.

Se você quiser multi-modo no OpenCode:

1. Conecte seus provedores de IA no OpenCode primeiro
2. Crie um perfil via TUI do kortex ("Perfis SDD do OpenCode") ou CLI (flag `--profile`)
3. O perfil gera um orquestrador personalizado + subagentes, cada um atribuído ao modelo escolhido
4. No OpenCode, pressione **Tab** para alternar entre seu orquestrador padrão e os perfis personalizados

Você pode criar múltiplos perfis (ex: "economico" para experimentação, "premium" para produção) e alternar entre eles livremente.

Se você preferir um **gerenciador de perfis em tempo de execução** que mantém os perfis fora do `opencode.json`, o kortex agora também suporta isso. Durante a sincronização, o OpenCode pode detectar automaticamente arquivos de perfil externos em `~/.config/opencode/profiles/*.json` e alternar para um caminho de compatibilidade mais seguro que preserva o prompt do `sdd-orchestrator` ativo em vez de sobrescrevê-lo.

**Guia passo a passo completo**: [Perfis SDD do OpenCode](opencode-profiles.md)

---

## Subagentes — Mais Inteligentes do que Você Pensa

Quando o orquestrador delega o trabalho para um subagente (digamos, `sdd-explore` para investigar uma base de código), esse subagente não é um executor burro rodando um único script. É um agente completo com sua própria sessão, ferramentas e contexto.

O que os torna "super subagentes":

1. **Eles descobrem habilidades por conta própria.** A primeira ação de cada subagente é procurar o registro de skills — via memória do Engram ou pelo arquivo local `.atl/skill-registry.md`. Se encontrar skills relevantes (padrões React, testes em Go, arquitetura Angular, etc.), ele as carrega e as segue. O orquestrador não precisa "dar na boca" os caminhos das skills.

2. **Eles se adaptam ao seu projeto.** Um subagente `sdd-apply` trabalhando em um projeto React carregará padrões do React 19. O mesmo subagente trabalhando em um projeto Go carregará convenções de teste em Go. As skills que ele carrega dependem do que o registro diz ser relevante, não de uma lista fixa.

3. **Eles persistem seu trabalho.** Cada subagente salva seus artefatos no Engram antes de retornar. O próximo subagente no pipeline pode continuar exatamente de onde o anterior parou, mesmo entre diferentes sessões.

Este padrão funciona hoje em:

| Agente | Como os subagentes rodam |
|-------|-------------------|
| **OpenCode** | Sistema de subagentes nativo — cada fase é um agente dedicado com seu próprio modelo, ferramentas e permissões definidos no `opencode.json` |
| **Claude Code** | Através da ferramenta Agent — o orquestrador inicia subagentes que autodescobrem skills a partir do registro |
| **Outros** | O SDD roda em linha (sessão única) — o modelo segue as instruções do orquestrador sem iniciar agentes separados |

Você não precisa configurar nada disso. O instalador configura e o orquestrador gerencia a delegação automaticamente.

---

## Skills — Duas Camadas

O kortex instala as **skills do SDD** e **skills de base** (workflow, padrões de teste) diretamente no diretório de skills do seu agente. Estas estão embutidas no binário e estão sempre atualizadas.

Para **skills de codificação** (React 19, Angular, TypeScript, Tailwind, Zod, Playwright, etc.), a comunidade mantém um repositório separado: [fortissolucoescontato-bit/Kortex-Skills](https://github.com/fortissolucoescontato-bit/Kortex-Skills). Você as instala manualmente clonando o repositório e copiando as skills que desejar:

```bash
git clone https://github.com/fortissolucoescontato-bit/Kortex-Skills.git
cp -r Kortex-Skills/curated/react-19 ~/.claude/skills/
cp -r Kortex-Skills/curated/typescript ~/.claude/skills/
# ... ou copie todo o diretório curated/
```

Uma vez instaladas, seu agente detecta no que você está trabalhando e carrega as skills relevantes automaticamente. Você não precisa ativá-las ou invocá-las.

**O registro de skills.** O registro de skills é um catálogo de todas as skills disponíveis que o orquestrador lê uma vez por sessão para saber o que está disponível e onde. Ele precisa rodar **dentro de cada projeto** em que você trabalha, porque ele também busca por convenções no nível do projeto (como `CLAUDE.md`, `agents.md`, `.cursorrules`, etc.).

Como funciona:

1. **Execute `/skill-registry` dentro do seu projeto** — ele varre todas as suas skills instaladas (nível de usuário e nível de projeto), lê seus frontmatters e constrói um registro em `.atl/skill-registry.md`. Se o Engram estiver disponível, ele também salva o registro na memória para acesso entre sessões.
2. **O orquestrador o usa automaticamente** — uma vez que o registro existe, o orquestrador o lê no início da sessão e passa os caminhos das skills já resolvidos para os subagentes. Você não interage com o registro depois disso.
3. **Execute novamente quando as coisas mudarem** — sempre que você adicionar, remover ou modificar uma skill, execute `/skill-registry` novamente para que o orquestrador perceba as mudanças.

Também há um lado automatizado: o `sdd-init` executa a mesma lógica de registro internamente, então, se você usar o SDD em um novo projeto, o registro é construído como parte desse fluxo.

**Dica de mestre**: Se você costuma atualizar skills com frequência, pode criar uma skill (usando o `/skill-creator`) que dispare automaticamente uma atualização do registro após mudanças nas skills — assim você nunca precisará pensar nisso.

---

## A Regra de Ouro

O Kortex é um **configurador** de ecossistema. Ele prepara seu agente de IA com memória, skills, fluxos de trabalho e uma persona — e depois sai do caminho.

Quanto menos você pensar no kortex após a instalação, melhor ele estará funcionando.

---

## Referência Rápida

| Faça | Não Faça |
|----|-------|
| Execute o instalador, escolha seus agentes e preset | Edite manualmente os arquivos de configuração gerados |
| Apenas comece a programar com seu agente de IA | Memorize as fases ou comandos do SDD |
| Deixe o agente sugerir o SDD quando a tarefa for grande | Force o SDD em cada pequena tarefa |
| Confie que o Engram está salvando o contexto para você | Bisbilhote o armazenamento do Engram, a menos que precise de `engram sync` ou `engram tui` |
| Execute `/skill-registry` após instalar ou alterar skills | Esqueça de atualizar o registro após adicionar novas skills |
| Diga "use sdd" se souber que quer um planejamento estruturado | Se preocupe com qual fase do SDD vem a seguir |
| Execute o instalador novamente para atualizar ou mudar seu setup | Tente "remendar" arquivos de skill ou instruções de persona manualmente |
