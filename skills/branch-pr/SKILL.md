---
name: kortex-branch-pr
description: >
  Fluxo de trabalho para criação de PR no Kortex seguindo o sistema de obrigatoriedade de issue.
  Trigger: Ao criar um pull request, abrir um PR ou preparar alterações para revisão.
license: Apache-2.0
metadata:
  author: nexo-fortis
  version: "2.0"
---

# Kortex — Skill de Branch & PR

## Quando Usar

Carregue esta skill sempre que precisar:
- Criar uma branch para uma nova correção ou funcionalidade
- Abrir um pull request em [fortissolucoescontato-bit/kortex](https://github.com/fortissolucoescontato-bit/kortex)
- Preparar alterações para revisão

## Regras Críticas

1. **Todo PR DEVE vincular uma issue aprovada** — Use `Closes/Fixes/Resolves #<N>` no corpo do PR. A issue DEVE ter o `status:approved`. PRs sem isso serão **rejeitados automaticamente** pelo CI.
2. **Exatamente uma etiqueta `type:*`** — Aplique exatamente UMA etiqueta de tipo ao PR. O CI rejeitará PRs com zero ou múltiplas etiquetas de tipo.
3. **5 verificações automatizadas devem passar** — Veja a tabela de Verificações Automatizadas abaixo.
4. **Sem trailers `Co-Authored-By`** — Nunca adicione atribuição de IA aos commits.
5. **Sem force-push para main/master** — Branch protegida.

## Fluxo de Trabalho

```
1. Confirme se a issue tem status:approved
   gh issue view <N> --repo fortissolucoescontato-bit/kortex

2. Crie uma branch a partir da main usando a convenção de nomenclatura abaixo

3. Implemente as alterações seguindo as especificações e o design

4. Execute os testes localmente (unitários + E2E)

5. Faça o commit usando o formato Conventional Commits

6. Abra um PR referenciando a issue
   → Adicione exatamente UMA etiqueta type:*
   → Preencha o corpo do PR usando o template

7. Todas as 5 verificações automatizadas devem passar antes do merge
```

---

## Nomenclatura de Branch

Os nomes das branches **devem** seguir este padrão:

```
^(feat|fix|chore|docs|style|refactor|perf|test|build|ci|revert)\/[a-z0-9._-]+$
```

| Tipo | Exemplo |
|------|---------|
| `feat/` | `feat/login-usuario` |
| `fix/` | `fix/correcao-insercao-duplicada` |
| `docs/` | `docs/atualizacao-referencia-api` |
| `refactor/` | `refactor/extrair-sanitizador-query` |
| `chore/` | `chore/atualizar-bubbletea-v0.26` |
| `style/` | `style/corrigir-avisos-linter` |
| `perf/` | `perf/otimizar-carregamento-catalogo` |
| `test/` | `test/adicionar-cobertura-pipeline` |
| `build/` | `build/atualizar-config-goreleaser` |
| `ci/` | `ci/adicionar-job-docker-e2e` |
| `revert/` | `revert/desfazer-mudanca-seletor-modelo` |

**Regras:**
- Tudo em minúsculas
- Use hífens, pontos ou sublinhados como separadores (sem espaços, sem maiúsculas)
- A descrição deve ser curta e descritiva

---

## Formato do Corpo do PR

O corpo do PR deve seguir o template em `.github/PULL_REQUEST_TEMPLATE.md`. Todas as seções são obrigatórias, a menos que marcadas como opcionais.

```markdown
## 🔗 Issue Vinculada

Closes #<N>

## 🏷️ Tipo de PR

- [ ] `type:bug` — Correção de bug (mudança que não quebra a compatibilidade e corrige um erro)
- [ ] `type:feature` — Nova funcionalidade (mudança que não quebra a compatibilidade e adiciona funcionalidade)
- [ ] `type:docs` — Apenas documentação
- [ ] `type:refactor` — Refatoração de código (sem mudanças funcionais)
- [ ] `type:chore` — Mudanças em build, CI ou ferramentas
- [ ] `type:breaking-change` — Mudança que quebra a compatibilidade

## 📝 Resumo

<!-- Descrição clara do que este PR faz e por quê. -->

## 📂 Alterações

| Arquivo / Área | O que mudou |
|-------------|-------------|
| `caminho/do/arquivo` | Breve descrição |

## 🧪 Plano de Testes

**Testes Unitários**
\`\`\`bash
go test ./...
\`\`\`

**Testes E2E** (Requer Docker)
\`\`\`bash
cd e2e && ./docker-test.sh
\`\`\`

- [ ] Testes unitários passam (`go test ./...`)
- [ ] Testes E2E passam (`cd e2e && ./docker-test.sh`)
- [ ] Testado manualmente localmente

## ✅ Checklist do Contribuidor

- [ ] PR está vinculado a uma issue com `status:approved`
- [ ] Adicionei a etiqueta `type:*` apropriada a este PR
- [ ] Testes unitários passam (`go test ./...`)
- [ ] Testes E2E passam (`cd e2e && ./docker-test.sh`)
- [ ] Atualizei a documentação, se necessário
- [ ] Meus commits seguem o formato Conventional Commits
- [ ] Meus commits não incluem trailers `Co-Authored-By`
```

---

## Verificações Automatizadas

As 5 verificações são executadas em cada PR e **todas devem passar** antes do merge:

| Verificação | O que valida | Como corrigir |
|-------|-----------------|------------|
| **Check Issue Reference** | O corpo do PR contém `Closes/Fixes/Resolves #N` | Adicione `Closes #<N>` ao corpo do PR |
| **Check Issue Has `status:approved`** | A issue vinculada foi aprovada por um mantenedor | Aguarde o mantenedor adicionar `status:approved` à issue |
| **Check PR Has `type:*` Label** | Exatamente uma etiqueta `type:*` foi aplicada ao PR | Peça a um mantenedor para adicionar a etiqueta correta; remova as extras |
| **Unit Tests** | `go test ./...` passa | Corrija os testes que falharam antes de enviar |
| **E2E Tests** | `cd e2e && ./docker-test.sh` passa | Corrija os cenários E2E que falharam antes de enviar |

---

## Conventional Commits

As mensagens de commit **devem** seguir este padrão:

```
^(build|chore|ci|docs|feat|fix|perf|refactor|revert|style|test)(\([a-z0-9\._-]+\))?!?: .+
```

### Formato

```
<tipo>(<escopo-opcional>)!: <descrição>

[corpo opcional]

[rodapé opcional]
```

### Tipos Permitidos

| Tipo | Propósito | Etiqueta de PR |
|------|---------|----------|
| `feat` | Nova funcionalidade | `type:feature` |
| `fix` | Correção de bug | `type:bug` |
| `docs` | Apenas documentação | `type:docs` |
| `refactor` | Mudança de código (sem mudança de comportamento) | `type:refactor` |
| `chore` | Manutenção, dependências, ferramentas | `type:chore` |
| `style` | Formatação, linting (sem mudança de lógica) | `type:chore` |
| `perf` | Melhoria de performance | `type:feature` |
| `test` | Adição ou atualização de testes | `type:chore` |
| `build` | Sistema de build ou dependências externas | `type:chore` |
| `ci` | Configuração de CI | `type:chore` |
| `revert` | Reverte um commit anterior | corresponde ao tipo revertido |

### Breaking Changes (Mudanças de Ruptura)

Adicione `!` após o tipo/escopo:

```
feat(cli)!: renomear flag --config para --config-file

BREAKING CHANGE: a flag --config foi renomeada para --config-file.
```

Mudanças de ruptura mapeiam para a etiqueta `type:breaking-change`.

### Exemplos

```
feat(tui): adicionar barra de progresso nas etapas de instalação
fix(agent): corrigir detecção do Claude Code no macOS
docs: atualizar guia de contribuição
chore(deps): atualizar bubbletea para v0.26
refactor(pipeline): extrair executor de etapas
style: corrigir avisos do linter no pacote catalog
perf(system): cachear resultado da detecção de SO
test(installer): adicionar cobertura para execução de etapa do catálogo
build: atualizar config do goreleaser para arm64
ci: separar jobs de teste unitário e e2e
revert: desfazer redesign do seletor de modelo
feat(cli)!: alterar caminho padrão da configuração
```

---

## Comandos

### Configuração

```bash
# Confirme se a issue está aprovada antes de começar
gh issue view <N> --repo fortissolucoescontato-bit/kortex

# Criar branch
git checkout main && git pull
git checkout -b fix/<descricao-curta>
```

### Testando Localmente

```bash
# Testes unitários
go test ./...

# Testes unitários — pacote específico
go test ./internal/tui/...

# Testes unitários — modo detalhado
go test -v ./...

# Testes E2E (O Docker deve estar rodando)
cd e2e && ./docker-test.sh
```

### Abrir um PR

```bash
gh pr create \
  --repo fortissolucoescontato-bit/kortex \
  --title "fix(agent): corrigir detecção do Claude Code no Linux" \
  --body "$(cat <<'EOF'
## 🔗 Issue Vinculada

Closes #42

## 🏷️ Tipo de PR

- [x] \`type:bug\` — Correção de bug (mudança que não quebra a compatibilidade e corrige um erro)

## 📝 Resumo

Corrige a detecção do binário Claude Code que falhava no Linux quando HOME não estava definido.

## 📂 Alterações

| Arquivo / Área | O que mudou |
|-------------|-------------|
| \`internal/agents/claude.go\` | Adicionado fallback para a variável de ambiente HOME |

## 🧪 Plano de Testes

- [x] Testes unitários passam (\`go test ./...\`)
- [x] Testes E2E passam (\`cd e2e && ./docker-test.sh\`)
- [x] Testado manualmente localmente

## ✅ Checklist do Contribuidor

- [x] PR está vinculado a uma issue com \`status:approved\`
- [x] Adicionei a etiqueta \`type:*\` apropriada a este PR
- [x] Testes unitários passam (\`go test ./...\`)
- [x] Testes E2E passam (\`cd e2e && ./docker-test.sh\`)
- [x] Atualizei a documentação, se necessário
- [x] Meus commits seguem o formato Conventional Commits
- [x] Meus commits não incluem trailers \`Co-Authored-By\`
EOF
)"
```

### Verificar Status do PR

```bash
gh pr checks --repo fortissolucoescontato-bit/kortex <numero-do-PR>
gh pr view --repo fortissolucoescontato-bit/kortex <numero-do-PR>
```

### Adicionar uma Etiqueta

```bash
gh pr edit <numero-do-PR> --repo fortissolucoescontato-bit/kortex --add-label "type:bug"
```
