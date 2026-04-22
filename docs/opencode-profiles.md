# Perfis SDD do OpenCode

← [Voltar para o README](../README.md)

---

Você configurou seus modelos SDD uma vez e agora toda tarefa — seja barata ou cara, experimental ou testada em batalha — passa pelo mesmo orquestrador. Os perfis resolvem isso: **crie configurações de modelos nomeadas e alterne entre elas com Tab dentro do OpenCode.**

O Kortex suporta **duas maneiras** de trabalhar com perfis do OpenCode:

1.  **Modo multi-perfil gerado** — o fluxo clássico do Kortex. Cada perfil nomeado gera seu próprio `sdd-orchestrator-{nome}` mais 10 subagentes com sufixo no `opencode.json`, e você alterna entre eles com **Tab**.
2.  **Modo externo de perfil único ativo** — para ferramentas da comunidade que mantêm arquivos de perfil fora do `opencode.json` e ativam um perfil de cada vez em tempo de execução.

Isso significa que você pode continuar com o overlay multi-perfil embutido ou conectar o Kortex a um gerenciador de perfis externo sem que os dois sistemas conflitem.

---

## Início Rápido (TUI)

1.  Inicie o instalador: `kortex` (ou `go run ./cmd/kortex`).
2.  Selecione **"Perfis SDD do OpenCode"** na tela de boas-vindas.
3.  Selecione **"Criar novo perfil"** (ou pressione `n`).
4.  Digite um nome para o perfil em formato slug (letras minúsculas e hifens). Exemplo: `economico`.
5.  Escolha o modelo do orquestrador (provedor e depois o modelo — usa o seletor de modelos existente).
6.  Atribua modelos aos subagentes (use "Definir todas as fases" para uma configuração uniforme ou configure cada fase individualmente).
7.  Confirme — o instalador grava o perfil no `opencode.json` e executa a sincronização.

Abra o OpenCode e pressione **Tab** — seu novo orquestrador aparecerá ao lado do padrão.

## Início Rápido (CLI)

Crie um perfil durante a sincronização com `--profile nome:provedor/modelo`:

```bash
kortex sync --profile economico:anthropic/claude-haiku-3.5-20241022
```

Múltiplos perfis em um comando:

```bash
kortex sync \
  --profile economico:anthropic/claude-haiku-3.5-20241022 \
  --profile premium:anthropic/claude-opus-4-20250514
```

Sobrescreva uma fase específica com `--profile-phase nome:fase:provedor/modelo`:

```bash
kortex sync \
  --profile economico:anthropic/claude-haiku-3.5-20241022 \
  --profile-phase economico:sdd-apply:anthropic/claude-sonnet-4-20250514
```

Isso cria um perfil "economico" onde tudo roda no Haiku, exceto o `sdd-apply`, que usa o Sonnet.

## Gerenciadores de Perfis Externos

Se você estiver usando uma ferramenta da comunidade que armazena perfis em `~/.config/opencode/profiles/*.json` e os ativa em tempo de execução, o Kortex agora pode sincronizar o OpenCode em um modo de compatibilidade.

### Autodetecção

Ao executar `kortex sync`, se existirem arquivos de perfil do OpenCode em:

```text
~/.config/opencode/profiles/*.json
```

O Kortex alterna automaticamente para a estratégia **`external-single-active`** para a sincronização do OpenCode.

### Sobrescrita manual

Você também pode forçar a estratégia explicitamente:

```bash
kortex sync --agent opencode --sdd-profile-strategy external-single-active
```

Ou forçar o comportamento clássico de overlay gerado:

```bash
kortex sync --agent opencode --sdd-profile-strategy generated-multi
```

### O que o modo de compatibilidade faz

No modo `external-single-active`, o Kortex:

- Continua gravando os assets base do SDD do OpenCode e arquivos de prompt compartilhados.
- **Não** regenera automaticamente perfis nomeados com sufixo no `opencode.json`.
- **Preserva o prompt atual do `sdd-orchestrator`** durante a sincronização para que as ferramentas externas possam manter suas políticas de runtime / blocos de fallback intactos.

Este é o ponto importante: o Kortex ainda mantém a base do SDD, mas para de agir como se o `opencode.json` fosse a única fonte da verdade para cada perfil.

## Usando Perfis no OpenCode

Após criar perfis no modo gerado, cada um aparece como um orquestrador selecionável no OpenCode:

| O que você vê no Tab | O que ele executa |
|---|---|
| `sdd-orchestrator` | Perfil padrão (sua configuração original) |
| `sdd-orchestrator-economico` | Perfil "economico" — Haiku em tudo |
| `sdd-orchestrator-premium` | Perfil "premium" — Opus em tudo |

Pressione **Tab** para alternar entre os orquestradores. Todos os comandos slash do SDD (`/sdd-new`, `/sdd-ff`, `/sdd-explore`, etc.) rodam no orquestrador selecionado. O orquestrador delega para seus próprios subagentes com sufixo (ex: `sdd-apply-economico`), garantindo que os perfis nunca interfiram uns nos outros.

Se você estiver usando um gerenciador externo de perfil único ativo, geralmente continuará trabalhando com o `sdd-orchestrator` base, enquanto a ferramenta externa troca as atribuições de modelos ativos em tempo de execução.

## Gerenciando Perfis

Na tela de lista de perfis da TUI:

| Ação | Tecla | Notas |
|---|---|---|
| Editar um perfil | `Enter` no perfil | Altere os modelos e depois sincronize |
| Excluir um perfil | `d` no perfil | Remove o orquestrador + todos os subagentes do JSON |
| Criar um novo perfil | `n` (ou selecione "Criar novo perfil") | Fluxo completo de criação |

O perfil `default` (o `sdd-orchestrator` sem sufixo) pode ser editado, mas não excluído — ele sempre existe quando o SDD está configurado.

### Regras para nomes de perfis

| Entrada | Válido? | Motivo |
|---|---|---|
| `economico` | Sim | Slug simples |
| `premium-v2` | Sim | Hifens são permitidos |
| `meu perfil` | Não | Espaços não são permitidos |
| `default` | Não | Reservado para o orquestrador base |
| `ALTO` | Torna-se `alto` | Convertido automaticamente para minúsculas |

---

<details>
<summary><strong>Como Funciona</strong></summary>

No modo multi-perfil gerado, cada perfil gera 11 entradas de agente no `opencode.json`: um orquestrador (`sdd-orchestrator-{nome}`, modo `primary`) e 10 subagentes (`sdd-{fase}-{nome}`, modo `subagent`, ocultos). As permissões do orquestrador são limitadas para que ele possa delegar apenas para seus próprios subagentes com sufixo.

Os prompts dos subagentes são compartilhados entre todos os perfis como arquivos em `~/.config/opencode/prompts/sdd/` (ex: `sdd-apply.md`). Cada entrada de agente referencia o arquivo compartilhado via `{file:~/.config/opencode/prompts/sdd/sdd-apply.md}` — apenas o campo `model` difere entre os perfis. Os prompts do orquestrador são inseridos em linha por perfil porque contêm tabelas de atribuição de modelos e referências de subagentes específicas de cada perfil.

Durante a sincronização ou atualização, o Kortex usa uma das duas estratégias:

- **`generated-multi`** — varre o `opencode.json` em busca de `sdd-orchestrator-*`, atualiza prompts compartilhados, regenera orquestradores de perfil e preserva atribuições de modelos.
- **`external-single-active`** — detecta arquivos de perfil externos, mantém os assets compartilhados do SDD atualizados e preserva o prompt do orquestrador base existente em vez de sobrescrever extensões externas de runtime.

</details>

---

← [Voltar para o README](../README.md)
