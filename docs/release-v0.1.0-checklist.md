# Checklist de Lançamento v0.1.0

## Congelamento de Escopo (Scope Freeze)

- [ ] Confirmar que o escopo do MVP permanece apenas macOS + Claude Code + OpenCode.
- [ ] Confirmar que nenhum recurso pós-MVP foi mesclado.

## Portões de Qualidade (Quality Gates)

- [ ] Executar testes unitários direcionados para cobertura de verificação e golden files.
- [ ] Executar testes de paridade da CLI de instalação.
- [ ] Validar que a saída do dry-run é gerada sem efeitos colaterais de aplicação.

## Teste de Fumaça Manual (Smoke Test)

- [ ] Executar dry-run da instalação no macOS.
- [ ] Executar instalação real em uma conta de teste do macOS.
- [ ] Validar os caminhos de saída principais para Claude Code e OpenCode.
- [ ] Validar que o endpoint de saúde do Kortex-Engram está acessível quando selecionado.

## Documentação

- [ ] O README referencia os documentos de início rápido, modo não interativo e rollback.
- [ ] Os comandos de Início Rápido refletem as flags atuais da CLI.
- [ ] As orientações de Rollback coincidem com o comportamento de backup/restauração.

## Preparação para o Lançamento

- [ ] Criar tag `v0.1.0`.
- [ ] Publicar as notas de lançamento com limitações conhecidas.
- [ ] Anunciar explicitamente as restrições do MVP (apenas macOS, suporte limitado a agentes).
