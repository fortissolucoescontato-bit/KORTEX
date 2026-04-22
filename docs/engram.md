# Referência de Comandos do Engram

← [Voltar para o README](../README.md)

---

O Engram funciona automaticamente. Seu agente de IA salva decisões, descobertas e contexto em uma memória persistente sem que você precise fazer nada. Você não precisa memorizar comandos ou gerenciar a memória manualmente.

Esta página existe para quando você quiser inspecionar, compartilhar ou corrigir suas memórias manualmente.

---

## Comandos do Dia a Dia

Estes são os únicos comandos que a maioria das pessoas precisa.

```bash
# Navegue por suas memórias visualmente — busque, filtre e detalhe observações
engram tui

# Busque pelo terminal sem abrir a TUI
engram search "refatoração de auth"

# Exporta as memórias do projeto para .engram/ para que você possa commitá-las no git
engram sync
```

O `engram tui` é a maneira mais rápida de ver o que seu agente tem salvado. Comece por lá.

---

## Gerenciamento de Projetos

O Engram agrupa memórias por nome de projeto, detectado automaticamente a partir do remoto git desde a v1.11.0. Às vezes, os projetos acabam com nomes duplicados (ex: "meu-app" vs "Meu-App" vs "meu-app-frontend"). Estes comandos corrigem isso.

```bash
# Lista todos os projetos com contagem de observações
engram projects list

# Mescla nomes de projetos duplicados em um só de forma interativa
engram projects consolidate
```

O `projects list` mostra todos os projetos que o Engram conhece e quantas observações cada um possui. Se você vir o mesmo projeto sob vários nomes, execute `projects consolidate` para mesclá-los.

O equivalente MCP é `mem_merge_projects`, que o agente de IA pode chamar diretamente quando detecta divergência de nomes.

---

## Compartilhamento em Equipe

As memórias do Engram vivem localmente por padrão. Para compartilhá-las com sua equipe via git:

```bash
# Após uma sessão de trabalho — exporta as memórias para .engram/ no seu repositório
engram sync

# Em outra máquina — importa as memórias após clonar o repositório
engram sync --import
```

Adicione o diretório `.engram/` ao seu repositório e faça o commit. Quando um colega clonar e executar `engram sync --import`, ele receberá todo o contexto do projeto. Isso é especialmente útil para o onboarding — novos contribuidores começam com o conhecimento acumulado da equipe.

---

## Referência de Ferramentas MCP

Estas são as ferramentas que o agente de IA usa nos bastidores. Você nunca as chama diretamente, mas entendê-las ajuda a saber o que seu agente está fazendo.

### Ferramentas Principais

| Ferramenta | O que faz |
|------|--------------|
| `mem_save` | Salva uma decisão, correção de bug, descoberta ou convenção na memória |
| `mem_search` | Busca na memória por palavras-chave — retorna observações correspondentes |
| `mem_context` | Obtém o histórico recente da sessão (chamado no início da sessão) |
| `mem_session_summary` | Salva um resumo do fim da sessão para que a próxima tenha contexto |
| `mem_get_observation` | Recupera o conteúdo completo e não truncado de uma observação específica pelo ID |
| `mem_save_prompt` | Salva o prompt do usuário para contexto adicional |

### Ferramentas Avançadas

<details>
<summary>Clique para expandir — raramente necessárias, mas disponíveis</summary>

| Ferramenta | O que faz |
|------|--------------|
| `mem_update` | Atualiza uma observação existente pelo ID |
| `mem_suggest_topic_key` | Sugere uma chave de tópico estável para tópicos em evolução |
| `mem_session_start` / `mem_session_end` | Gerenciamento do ciclo de vida da sessão |
| `mem_stats` | Estatísticas de memória (contagem de observações, detalhamento por projeto) |
| `mem_delete` | Exclui uma observação pelo ID |
| `mem_timeline` | Visualização cronológica das observações |
| `mem_capture_passive` | Extrai aprendizados da conversa de forma passiva |
| `mem_merge_projects` | Mescla variantes de nomes de projetos (equivalente na CLI: `engram projects consolidate`) |

</details>

---

## Como a Detecção de Projeto Funciona

Desde a v1.11.0, o Engram lê a URL remota do git na inicialização, normaliza para minúsculas e usa isso como o nome do projeto. Se encontrar nomes de projetos existentes semelhantes, ele avisa. Isso evita o problema mais comum — o mesmo projeto acumulando memórias sob nomes ligeiramente diferentes.

Se você estiver trabalhando fora de um repositório git, o Engram usa o nome do diretório como alternativa.

---

## Documentação Completa

Para o código-fonte completo, opções de configuração e guia de contribuição: [github.com/fortissolucoescontato-bit/engram](https://github.com/fortissolucoescontato-bit/engram)
