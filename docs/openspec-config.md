# Configuração OpenSpec para SDD

O arquivo `openspec/config.yaml` é uma convenção documentada no nível do projeto para o SDD no `kortex` ao trabalhar nos modos de persistência `openspec` ou `hybrid`.

Ele é útil, e partes do prompt e da pilha de skills do SDD o procuram hoje, mas isso não deve ser lido como um schema de tempo de execução formalmente versionado e imposto pelo código Go.

## O que "Suporte" Significa Hoje

No repositório atual, o suporte ao `openspec/config.yaml` é baseado principalmente em prompts:

- As skills do SDD e os prompts do orquestrador dizem aos agentes para ler ou gravar este arquivo.
- O `sdd-init` e exemplos de convenções compartilhadas mostram os formatos de arquivo que se espera que os agentes criem.
- Fases posteriores podem reutilizar valores como `context`, `strict_tdd`, `rules` e `testing`.

O que NÃO é verdade hoje:

- Não existe um parser ou validador do lado do Go que imponha um schema canônico para o `openspec/config.yaml`.
- Não existe um contrato de compatibilidade forte que garanta que cada campo documentado seja consumido uniformemente em todas as fases.
- O formato exato ainda é melhor compreendido como a convenção atual do repositório, não como uma especificação pública bloqueada.

## O que Este Arquivo Pode Customizar

Este arquivo é usado pelas skills do SDD como contexto de projeto compartilhado e como um lugar para declarar regras específicas de cada fase.

O `openspec/config.yaml` pode ser usado para customizar o comportamento do SDD por convenções de projeto, especificamente:

- Contexto do projeto reutilizado entre as fases.
- Ativação de TDD estrito.
- Regras específicas de fases para proposal, specs, design, tasks, apply, verify e archive.
- Sobrescrita de comandos e cobertura usados pelos prompts de apply/verify.
- Capacidades de teste em cache para fluxos de apply/verify.

## Quais Fases o Referenciam

As seguintes fases do SDD referenciam explicitamente o `openspec/config.yaml` nos assets atuais de prompt/skill:

| Fase | Como usa a configuração |
|-------|-------------------------|
| `sdd-init` | No modo OpenSpec, as instruções do prompt dizem ao agente para criar o arquivo e gravar as seções detectadas de `context`, `rules` e `testing`. |
| `sdd-explore` | Lê o arquivo como parte da descoberta do contexto do projeto. |
| `sdd-propose` | Aplica `rules.proposal` se presente. |
| `sdd-design` | Aplica `rules.design` se presente. |
| `sdd-spec` | Aplica `rules.specs` se presente. |
| `sdd-tasks` | Aplica `rules.tasks` se presente. |
| `sdd-apply` | Lê `strict_tdd`, `testing` e `rules.apply` se presente. |
| `sdd-verify` | Lê `strict_tdd`, `testing` e `rules.verify` se presente. |
| `sdd-archive` | Aplica `rules.archive` se presente. |

## Exemplo de Convenção Sintetizada

Combinando o documento de convenção compartilhada atual, a orientação do `sdd-init` e as referências de apply/verify, a estrutura prática de alto nível se parece com isto:

```yaml
schema: spec-driven

context: |
  Tech stack: ...
  Arquitetura: ...
  Testes: ...
  Estilo: ...

strict_tdd: true

rules:
  proposal:
    - Incluir plano de rollback para mudanças arriscadas
  specs:
    - Usar Dado/Quando/Então para cenários
  design:
    - Documentar decisões de arquitetura com a justificativa
  tasks:
    - Manter as tarefas concluíveis em uma única sessão
  apply:
    - Seguir padrões de código existentes
  verify:
    test_command: ""
    build_command: ""
    coverage_threshold: 0
  archive:
    - Avisar antes de mesclar deltas destrutivos

testing:
  strict_tdd: true
  detected: "AAAA-MM-DD"
  runner:
    command: "go test ./..."
    framework: "Go standard testing"
```

Trate isso como uma síntese prática de campos que a camada de prompt atualmente pode ler ou emitir, não como uma definição estrita de schema.

## Referência de Campos

### `schema`

- Valor esperado nos exemplos: `spec-driven`
- Propósito: identifica o arquivo como uma configuração SDD/OpenSpec nos exemplos atuais e convenções compartilhadas.

### `context`

- Tipo: string multi-linha
- Propósito: contexto do projeto em cache para fases posteriores do SDD.
- Conteúdo típico: stack, arquitetura, testes, estilo e outras convenções do projeto.

### `strict_tdd`

- Tipo: booleano
- Referenciado por: `sdd-init`, prompts do orquestrador, `sdd-apply`, `sdd-verify`
- Propósito: ativa ou desativa o comportamento de TDD estrito quando o suporte a testes existe.

### `rules`

- Tipo: mapa indexado por fase
- Propósito: anexar convenções de projeto a cada fase do SDD.
- Chaves de fase conhecidas:
  - `proposal`
  - `specs`
  - `design`
  - `tasks`
  - `apply`
  - `verify`
  - `archive`

### `testing`

- Tipo: objeto estruturado
- Geralmente escrito por: `sdd-init`
- Referenciado por: `sdd-apply`, `sdd-verify`
- Propósito: armazena em cache as capacidades de teste detectadas para que as fases não precisem redescobri-las todas as vezes.

## Ressalvas

### O Formato Não é Totalmente Uniforme

Existem inconsistências importantes nos exemplos atuais e referências de skills:

- O `sdd-init` mostra `rules.apply` e `rules.verify` como listas simples de instruções.
- O documento de convenção OpenSpec compartilhada mostra `rules.apply` com `tdd` e `test_command`, e `rules.verify` com chaves estruturadas de sobrescrita.
- O arquivo `sdd-apply/strict-tdd.md` refere-se a `rules.apply.test_command`.
- O `sdd-verify` refere-se a `rules.verify.test_command`, `rules.verify.build_command` e `rules.verify.coverage_threshold`.

Isso significa que `rules.apply` e `rules.verify` são atualmente tratados como se pudessem conter chaves estruturadas, enquanto outros exemplos mostram essas mesmas regras de fase como listas simples.

### Prefira "Convenção Atual" em vez de "Contrato Estável"

Se você estiver documentando ou dependendo deste arquivo, a abordagem mais segura hoje é:

- `openspec/config.yaml` faz parte do fluxo de trabalho OpenSpec documentado.
- Várias skills do SDD e assets de prompt o procuram e podem gravar nele.
- Os campos documentados refletem as convenções atuais do repositório.
- Consumidores não devem assumir um schema de tempo de execução formalmente versionado e imposto, a menos que a implementação adicione parsing/validação explícitos posteriormente.
