# Proposal: bug-fixes

**Change ID**: `bug-fixes`
**Project**: kortex
**Status**: Proposed
**Author**: sdd-propose (Opus 4.7)
**Date**: 2026-04-24

---

## 1. Intent

Eliminar uma safra de 10 bugs catalogados na fase de exploração que comprometem **correção de dados**, **robustez operacional** e **higiene de código** do Kortex. Três deles (CRITICAL) têm potencial direto de mascarar falhas de I/O, vazar recursos (goroutines/contextos pendurados) e deixar o sistema em estado inconsistente durante operações de rollback. Os demais (SERIOUS / WARNING) degradam silenciosamente a experiência do usuário, a manutenção e a portabilidade entre SOs.

Corrigir tudo num único change permite:

- Auditar o resultado em um só PR revisável (evita churn em múltiplos PRs mecânicos).
- Aplicar o mesmo padrão de testes (TDD estrito) de forma coerente.
- Atualizar os goldens uma única vez, minimizando ruído de regressão.

## 2. Scope

### Arquivos afetados (9 arquivos + testes associados)

| # | Severidade | Arquivo | Função / Símbolo | Problema |
|---|------------|---------|------------------|----------|
| 1 | CRITICAL   | `internal/state/state.go` | `listXxx` / qualquer `sql.Rows` consumer | `rows.Err()` não é checado após `rows.Next()` loop — erros de iteração são silenciados. |
| 2 | CRITICAL   | `internal/tui/model.go` (ou `tui/install*.go`) | Mensagem/command de install | `context` com timeout é criado mas descartado (ctx cancelado imediatamente ou nunca propagado). |
| 3 | CRITICAL   | `internal/pipeline/rollback.go` | `Rollback(...)` | `return` no primeiro erro impede desfazer demais operações; sistema fica parcialmente rollback-ado. |
| 4 | SERIOUS    | `internal/update/instructions.go` | `reH2Section` (regex) | Regex declarada mas nunca usada (código morto / intenção perdida). |
| 5 | SERIOUS    | `internal/agentbuilder/installer.go` | Rollback do installer | Remove arquivos, mas não remove diretórios criados durante install. |
| 6 | SERIOUS    | `internal/components/kortex-engram/*.go` | Lookup duplicado | Duas chamadas equivalentes de resolver kortex-engram no mesmo caminho de execução. |
| 7 | SERIOUS    | `internal/cli/run.go` (ou similar) | Config do server | Atribuições redundantes (campo setado duas vezes com o mesmo valor ou sobrescrito inutilmente). |
| 8 | WARNING    | `internal/components/kortex-engram/download.go` | Verificação de checksum | Falha de checksum é engolida (`_ =` ou log-only). |
| 9 | WARNING    | `internal/components/sdd/inject.go` (ou similar) | Construção de path | Usa `s1 + "/" + s2` em vez de `filepath.Join`. |
| 10| WARNING    | `internal/components/kortex-engram/verify*.go` | `Verify(...)` | `context.Background()` sem timeout → operação pode pendurar indefinidamente. |

> Os caminhos exatos serão confirmados pela fase `sdd-spec` a partir da exploração. Nenhum arquivo fora dessa lista deve ser tocado.

### Fora de escopo

- Refatorações estruturais não relacionadas aos 10 bugs.
- Introdução de novas dependências ou libs.
- Mudanças em API pública / contratos TUI além do mínimo necessário.
- Alterações em specs funcionais do Kortex (apenas correção de defeito).

## 3. Approach

Estratégia **bug-a-bug, TDD estrito** (modo ativo no projeto). Para cada bug:

1. Escrever teste que reproduz o defeito (red).
2. Corrigir o código (green).
3. Refatorar se necessário (refactor), mantendo suite verde.
4. Atualizar golden files apenas quando o comportamento observável mudar.

### Detalhe por bug

**CRITICAL #1 — `rows.Err()` ausente**
Após todo `for rows.Next() { ... }`, acrescentar `if err := rows.Err(); err != nil { return ..., err }`. Teste: injetar driver/mock que retorna erro no meio da iteração e assertar que o erro é propagado.

**CRITICAL #2 — context timeout descartado**
Identificar `ctx, cancel := context.WithTimeout(...)`. Garantir que (a) `cancel` é chamado via `defer`, (b) `ctx` é efetivamente passado para a operação de install. Teste: usar `context.WithDeadline` controlável e assertar cancelamento.

**CRITICAL #3 — rollback para no primeiro erro**
Trocar `return err` por acumulação: iterar todos os passos, coletar erros em `errors.Join(...)` (Go 1.20+) e retornar composto. Teste: rollback com 3 passos onde o segundo falha; assertar que o terceiro também foi executado e que o erro resultante menciona ambos os passos falhos.

**SERIOUS #4 — regex morta `reH2Section`**
Duas opções: (a) remover se de fato obsoleta; (b) religá-la ao fluxo original se a intenção se perdeu. Decisão será tomada na fase `sdd-design` após confirmar o histórico git. Default: remover, já que cobertura atual passa sem ela.

**SERIOUS #5 — installer rollback não remove diretórios**
Após remoção de arquivos, tentar `os.Remove(dir)` (não recursivo) em ordem inversa de criação; ignorar `ErrNotExist`/`ErrNotEmpty` silenciosamente. Teste: install cria `a/b/c/file`, rollback remove file + c + b + a.

**SERIOUS #6 — lookup duplicado kortex-engram**
Extrair para variável local cacheada dentro do escopo da função. Teste: spy no resolver assertando uma única chamada.

**SERIOUS #7 — atribuições redundantes no server**
Remover a atribuição dominada. Cobertura existente deve continuar verde; adicionar teste-guarda apenas se a redundância representar bug latente (ex.: valores diferentes em ordens diferentes).

**WARNING #8 — checksum engolido**
Propagar erro de checksum como `fmt.Errorf("checksum mismatch: %w", err)`. Teste: download com checksum corrompido retorna erro.

**WARNING #9 — concat de path**
Substituir `"a/" + b` por `filepath.Join("a", b)`. Teste: asserção em Windows-style path (`filepath.ToSlash`) opcional — mínimo: test unitário existente continua verde.

**WARNING #10 — `context.Background()` sem timeout em verify**
Aceitar `ctx` como parâmetro ou aplicar `context.WithTimeout(ctx, N)`. Valor N definido no design (sugestão: 30s, alinhado com o resto do módulo).

## 4. Risks

| Risco | Probabilidade | Impacto | Mitigação |
|-------|---------------|---------|-----------|
| Golden files quebram em cascata após correções de TUI/install | Média | Baixo | Atualizar goldens num commit dedicado; revisar diff visualmente. |
| `errors.Join` muda formato de mensagens que testes atuais asseguram | Baixa | Médio | Ajustar assertions para `errors.Is` em vez de string-compare. |
| Remoção de `reH2Section` esconder feature parcialmente implementada | Baixa | Médio | Verificar `git log -S reH2Section` antes de remover; se houver intenção viva, religar em vez de deletar. |
| Timeout em verify muito curto em máquinas lentas (CI) | Média | Baixo | Tornar configurável via variável de ambiente ou flag; default conservador. |
| Novo comportamento de rollback (continua após erro) pode mascarar falha original no log | Baixa | Médio | Log estruturado por passo + erro agregado com `errors.Join` (preserva cada causa). |

## 5. Rollback plan

Por tratar-se de correções independentes, cada bug é um commit atômico seguindo Conventional Commits (`fix(scope): ...`). Se qualquer correção se mostrar problemática:

1. **Rollback granular**: `git revert <sha>` do commit específico. Como os commits são atômicos e independentes, reverter um não desfaz os demais.
2. **Rollback total do PR**: se muitos bugs convergirem para o mesmo módulo e o revert granular conflitar, reverter o merge commit do PR inteiro (`git revert -m 1 <merge-sha>`).
3. **Feature flag**: não aplicável — correções de defeito não entram sob flag; são sempre-ativas.
4. **Validação pós-rollback**: rodar `go test ./...` após qualquer revert para garantir que a suite permanece verde (os testes TDD adicionados podem precisar também ser revertidos junto com a correção).

Marcos de verificação antes do merge:

- [ ] `go test ./...` 100% verde localmente.
- [ ] `go vet ./...` sem warnings novos.
- [ ] Goldens atualizados quando aplicável, com diff explicado no PR.
- [ ] Nenhum arquivo fora da tabela de escopo modificado.
- [ ] Issue vinculada com `status:approved` e label `type:bug` (governança `kortex-branch-pr`).

---

## Next phase

→ `sdd-spec` — escrever delta specs por bug (cenários Given/When/Then) e `sdd-design` caso haja decisão arquitetural (#3 agregação de erros, #10 timeout config).
