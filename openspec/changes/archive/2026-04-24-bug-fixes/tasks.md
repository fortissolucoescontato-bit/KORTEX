# Tasks: bug-fixes
**Change**: bug-fixes
**Project**: kortex
**TDD Mode**: STRICT — `go test ./...`
**Created**: 2026-04-24

---

## Overview

10 bugs targetados em 6 pacotes. Ordenados por severidade: CRITICAL → SERIOUS → WARNING.
Cada ciclo TDD segue RED → GREEN → REFACTOR antes de avançar ao próximo bug.

**Total de tarefas**: 32

---

## Fase 1 — Infrastructure (Preparação)

### 1.0 — [DONE] Verificar estado inicial dos testes

- **Objetivo**: confirmar que `go test ./...` passa em baseline antes de qualquer alteração.
- **Comando**: `go test ./...`
- **Critério de aceite**: todos os testes passam (ou falhas pré-existentes são documentadas).
- **Arquivos**: N/A (somente execução)

---

## Fase 2 — Implementation: CRITICAL Bugs (RED → GREEN → REFACTOR)

### 2.0 — [DONE] [RED] Bug #1: `rows.Err()` ausente em `GetInstalledAgents`

- **Objetivo**: escrever teste que falha quando `rows.Err()` não é verificado.
- **Arquivo**: `internal/state/state_test.go`
- **O que fazer**: adicionar `TestGetInstalledAgents_RowsErr` que injeta um DB driver com erro de iteração e verifica que o erro é propagado.
- **Critério de aceite**: teste compila e falha com a implementação atual.

### 2.1 — [DONE] [GREEN] Bug #1: Propagar `rows.Err()` em `GetInstalledAgents`

- **Objetivo**: adicionar `if err := rows.Err(); err != nil { return nil, err }` após o loop.
- **Arquivo**: `internal/state/state.go`, linhas 40–48
- **Critério de aceite**: `TestGetInstalledAgents_RowsErr` passa.

### 2.2 — [DONE] [RED] Bug #1: `rows.Err()` ausente em `GetAssignments`

- **Objetivo**: escrever `TestGetAssignments_RowsErr` com a mesma mecânica de injeção.
- **Arquivo**: `internal/state/state_test.go`
- **Critério de aceite**: teste compila e falha.

### 2.3 — [DONE] [GREEN] Bug #1: Propagar `rows.Err()` em `GetAssignments`

- **Objetivo**: adicionar `if err := rows.Err(); err != nil { return nil, err }` após o loop em `GetAssignments`.
- **Arquivo**: `internal/state/state.go`, linhas 82–90
- **Critério de aceite**: `TestGetAssignments_RowsErr` passa; `go test ./internal/state/...` verde.

---

### 2.4 — [DONE] [RED] Bug #2: Context leak em `startInstallation`

- **Objetivo**: escrever teste `TestStartInstallation_ContextPropagated` que verifica que o contexto com timeout é passado para `Install` (e não descartado).
- **Arquivo**: `internal/tui/model_test.go` (ou arquivo de teste existente)
- **Critério de aceite**: teste compila e falha porque o ctx é descartado (`_ = ctx`).
- **Nota**: requer que `agentbuilder.InstallContext` exista (será implementado em 2.5).

### 2.5 — [GREEN] Bug #2: Adicionar `InstallContext` + corrigir `startInstallation`

- **Objetivo**: implementar `InstallContext(ctx context.Context, agent *GeneratedAgent, adapters []AdapterInfo) ([]InstallResult, error)` e tornar `Install` um thin shim; passar `ctx` em vez de `_ = ctx` em `startInstallation`.
- **Arquivos**:
  - `internal/agentbuilder/installer.go` — adicionar `InstallContext`; `Install` chama `InstallContext(context.Background(), ...)`
  - `internal/tui/model.go`, linhas 3390–3403 — substituir `_ = ctx` e chamar `agentbuilder.InstallContext(ctx, ...)`
- **Critério de aceite**: `TestStartInstallation_ContextPropagated` passa.

### 2.6 — [DONE] [REFACTOR] Bug #2: Limpeza de `startInstallation`

- **Objetivo**: garantir que a assinatura pública de `Install` (compatibilidade retroativa) não quebra nenhum teste existente.
- **Critério de aceite**: `go test ./internal/agentbuilder/... ./internal/tui/...` verde; sem regressões.

---

### 2.7 — [DONE] [RED] Bug #3: `ExecuteRollback` para na primeira falha

- **Objetivo**: escrever `TestExecuteRollback_ContinuesOnError` com dois steps de rollback onde o primeiro falha e verificar que o segundo ainda é executado e ambos os erros aparecem no resultado.
- **Arquivo**: `internal/pipeline/rollback_test.go` (arquivo existente ou novo)
- **Critério de aceite**: teste compila e falha com a implementação atual (retorna na primeira falha).

### 2.8 — [DONE] [GREEN] Bug #3: `ExecuteRollback` acumula erros com `errors.Join`

- **Objetivo**: substituir `return result` dentro do loop por acumulação de erros; usar `errors.Join` para consolidar no final.
- **Arquivo**: `internal/pipeline/rollback.go`, linhas 40–55
- **Dependência de import**: adicionar `"errors"` ao bloco de imports.
- **Critério de aceite**: `TestExecuteRollback_ContinuesOnError` passa; `result.Success = false` e `result.Err` contém todos os erros.

---

## Fase 3 — Implementation: SERIOUS Bugs

### 3.0 — [DONE] [RED] Bug #4: `reH2Section` dead regex / regex cacheado

- **Objetivo**: escrever `TestExtractSection_CachedRegex` que verifica que chamadas repetidas a `extractSection` com o mesmo nome retornam o mesmo resultado (e que a regex não é recompilada desnecessariamente — por exemplo, benchmarkando ou verificando comportamento com seções especiais).
- **Arquivo**: `internal/agentbuilder/parser_test.go`
- **Critério de aceite**: teste documenta o comportamento esperado; confirma que `reH2Section` (package-level) nunca é usada por `extractSection`.

### 3.1 — [DONE] [GREEN] Bug #4: Remover `reH2Section` e cachear regex em `extractSection`

- **Objetivo**: remover a variável `reH2Section` do bloco `var`; dentro de `extractSection`, criar a regex compilada uma única vez por chamada (já feito) — confirmar que não há referência a `reH2Section` no código.
- **Arquivo**: `internal/agentbuilder/parser.go`, linha 17
- **Critério de aceite**: `reH2Section` removida; `go vet ./...` sem erros; testes do parser verdes.

---

### 3.2 — [DONE] [RED] Bug #5: `rollback` não remove diretórios

- **Objetivo**: escrever `TestRollback_RemovesDirs` que cria diretórios, chama `Install` forçando falha após criação de dir, e verifica que os diretórios foram removidos.
- **Arquivo**: `internal/agentbuilder/installer_test.go`
- **Critério de aceite**: teste falha porque `rollback` só deleta arquivos, não diretórios.

### 3.3 — [DONE] [GREEN] Bug #5: `rollback` passa a usar `rollbackState{files, dirs}`

- **Objetivo**: mudar a assinatura interna de `rollback` para receber `rollbackState{files []string, dirs []string}`; acumular diretórios criados por `MkdirAll`; removê-los em ordem reversa após remover arquivos.
- **Arquivo**: `internal/agentbuilder/installer.go`, linhas 26–74
- **Critério de aceite**: `TestRollback_RemovesDirs` passa; testes de instalação existentes não regridem.

---

### 3.4 — [DONE] [RED] Bug #6: `resolveKortexEngramCommand` duplica LookPath

- **Objetivo**: escrever `TestResolveKortexEngramCommand_NoDuplicateLookPath` que verifica que a função realiza no máximo 2 lookups (kortex-engram e kortex), não 3.
- **Arquivo**: `internal/components/kortex-engram/inject_test.go`
- **Critério de aceite**: teste conta o número de chamadas ao LookPath mockado e falha se for > 2.

### 3.5 — [DONE] [GREEN] Bug #6: Remover terceiro LookPath duplicado em `resolveKortexEngramCommand`

- **Objetivo**: remover o terceiro bloco `kortexEngramLookPath("kortex-engram")` duplicado (linha 70–73).
- **Arquivo**: `internal/components/kortex-engram/inject.go`, linhas 61–75
- **Critério de aceite**: `TestResolveKortexEngramCommand_NoDuplicateLookPath` passa; função tem exatamente 2 lookups.

---

### 3.6 — [DONE] [RED] Bug #7: `existingMergedKortexEngramCommand` tem assignments redundantes

- **Objetivo**: escrever `TestExistingMergedKortexEngramCommand_NoRedundantAssignments` verificando que cada `case` no switch retorna o valor correto sem duplicação de lógica.
- **Arquivo**: `internal/components/kortex-engram/inject_test.go`
- **Critério de aceite**: teste verifica comportamento; documenta que as linhas `server = mcp["kortex-engram"] // Fallback` são dead code.

### 3.7 — [DONE] [GREEN] Bug #7: Remover assignments redundantes em `existingMergedKortexEngramCommand`

- **Objetivo**: remover as linhas de fallback redundantes (`server = mcp["kortex-engram"] // Fallback`, etc.) em todos os três cases do switch.
- **Arquivo**: `internal/components/kortex-engram/inject.go`, linhas 480–501
- **Critério de aceite**: `TestExistingMergedKortexEngramCommand_NoRedundantAssignments` passa; `go vet` limpo.

---

## Fase 4 — Implementation: WARNING Bugs

### 4.0 — [DONE] [RED] Bug #8: `Snapshotter.Create` silencia erro de checksum

- **Objetivo**: escrever `TestSnapshotter_Create_ChecksumError` que injeta um erro no `ComputeChecksum` (via mock/substituição de função) e verifica que `Create` retorna esse erro.
- **Arquivo**: `internal/backup/snapshot_test.go`
- **Critério de aceite**: teste compila e falha (atualmente o erro é apenas logado).

### 4.1 — [DONE] [GREEN] Bug #8: `Snapshotter.Create` retorna erro de checksum

- **Objetivo**: substituir o `log.Printf` de checksum por `return Manifest{}, fmt.Errorf("compute checksum: %w", csErr)`.
- **Arquivo**: `internal/backup/snapshot.go`, linhas 79–85
- **Critério de aceite**: `TestSnapshotter_Create_ChecksumError` passa; importação de `"log"` pode ser removida se não usada em outro lugar.

---

### 4.2 — [DONE] [RED] Bug #9: `writeCodexInstructionFiles` usa concatenação de string para path

- **Objetivo**: escrever `TestWriteCodexInstructionFiles_UsesFilepathJoin` — verificar que o path produzido é equivalente ao `filepath.Join(homeDir, ".codex", "kortex-engram-instructions.md")` em todos os SOs (especialmente Windows com separador diferente).
- **Arquivo**: `internal/components/kortex-engram/inject_test.go`
- **Critério de aceite**: teste documenta a expectativa; em Linux sempre passa, mas torna a intenção explícita.

### 4.3 — [DONE] [GREEN] Bug #9: Usar `filepath.Join` em `writeCodexInstructionFiles`

- **Objetivo**: substituir `homeDir + "/.codex"` por `filepath.Join(homeDir, ".codex")` e construir os caminhos filhos com `filepath.Join(codexDir, ...)`.
- **Arquivo**: `internal/components/kortex-engram/inject.go`, linha 384
- **Critério de aceite**: `TestWriteCodexInstructionFiles_UsesFilepathJoin` passa; `go test ./...` verde.

---

### 4.4 — [DONE] [RED] Bug #10: `verify.RunChecks` sem timeout em `sync.go` e `run.go`

- **Objetivo**: escrever `TestRunChecksWithTimeout_SyncAndRun` verificando que as funções que chamam `verify.RunChecks` passam um contexto com timeout (não `context.Background()` nu).
- **Arquivos**: `internal/cli/sync_test.go`, `internal/cli/run_test.go`
- **Critério de aceite**: teste falha ou documenta a ausência de timeout com `context.Background()`.

### 4.5 — [DONE] [GREEN] Bug #10: Adicionar `context.WithTimeout` em `sync.go` e `run.go`

- **Objetivo**: substituir `context.Background()` nas chamadas a `verify.RunChecks` por `context.WithTimeout(context.Background(), 30*time.Second)` com defer cancel.
- **Arquivos**:
  - `internal/cli/sync.go`, linha 848
  - `internal/cli/run.go`, linha 1041
- **Critério de aceite**: `TestRunChecksWithTimeout_SyncAndRun` passa; imports de `"time"` adicionados se ausentes.

---

## Fase 5 — Validation

### 5.0 — [DONE] Executar suite completa

- **Objetivo**: confirmar que TODOS os testes passam após todas as correções.
- **Comando**: `go test ./...`
- **Critério de aceite**: zero falhas; zero regressões.

### 5.1 — [DONE] Executar `go vet` e verificar imports

- **Objetivo**: sem warnings de vet; sem imports não usados.
- **Comando**: `go vet ./...`
- **Critério de aceite**: saída limpa.

### 5.2 — [DONE] Atualizar golden files se necessário

- **Objetivo**: verificar se algum golden file foi impactado pelas correções (especialmente mudanças de comportamento em inject.go).
- **Comando**: identificar testes com `_golden` ou `testdata/` e re-executar com `-update` flag se necessário.
- **Critério de aceite**: golden files atualizados e commitados.

---

## Sumário de Arquivos por Tarefa

| Arquivo | Bugs | Tipo |
|---------|------|------|
| `internal/state/state.go` | #1 | GREEN |
| `internal/state/state_test.go` | #1 | RED |
| `internal/tui/model.go` | #2 | GREEN |
| `internal/agentbuilder/installer.go` | #2, #5 | GREEN |
| `internal/agentbuilder/installer_test.go` | #5 | RED |
| `internal/pipeline/rollback.go` | #3 | GREEN |
| `internal/pipeline/rollback_test.go` | #3 | RED |
| `internal/agentbuilder/parser.go` | #4 | GREEN |
| `internal/agentbuilder/parser_test.go` | #4 | RED |
| `internal/components/kortex-engram/inject.go` | #6, #7, #9 | GREEN |
| `internal/components/kortex-engram/inject_test.go` | #6, #7, #9 | RED |
| `internal/backup/snapshot.go` | #8 | GREEN |
| `internal/backup/snapshot_test.go` | #8 | RED |
| `internal/cli/sync.go` | #10 | GREEN |
| `internal/cli/run.go` | #10 | GREEN |
| `internal/cli/sync_test.go` | #10 | RED |
| `internal/cli/run_test.go` | #10 | RED |

---

## Dependências entre Tarefas

```
2.4 (RED InstallContext) → depende de 2.5 (GREEN InstallContext) existir
2.5 (GREEN InstallContext) → depende de 2.4 para validar
3.2 (RED rollback dirs) → independente
3.4 (RED resolveKortexEngram) → independente
5.2 (golden files) → depende de 3.5 e 3.7 estarem completos
```

Todos os outros ciclos RED→GREEN são independentes entre si e podem ser executados em paralelo por arquivo.
