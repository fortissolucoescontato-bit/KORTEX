# Guia de Backup e Rollback

O sistema de backup cria automaticamente snapshots dos seus arquivos de configuração antes de cada operação de instalação, sincronização e upgrade. Os backups são compactados, deduplicados e limpos automaticamente para manter o uso do disco sob controle.

## Como funciona

Toda vez que você executa `kortex install`, `sync` ou `upgrade`, o sistema:

1.  **Calcula um checksum** de todos os arquivos que serão incluídos no backup.
2.  **Pula o backup** se ele for idêntico ao mais recente (deduplicação).
3.  **Cria um snapshot compactado** (`snapshot.tar.gz`) com todos os seus arquivos de configuração.
4.  **Limpa backups antigos** — mantém os 5 mais recentes e exclui o restante.

## Conteúdo do snapshot

- `manifest.json` — metadados (origem, timestamp, contagem de arquivos, checksum, status de fixação).
- `snapshot.tar.gz` — arquivo compactado de todos os arquivos do backup.
- Para caminhos que não existiam antes da operação, o manifesto rastreia como `existed=false`.

Backups legados (anteriores à v1.16) usam um diretório `files/` com cópias simples em vez de um arquivo tar.gz. Ambos os formatos são totalmente suportados para restauração.

## Política de retenção

| Configuração | Padrão | Comportamento |
|---------|---------|----------|
| Contagem mantida | 5 | Os 5 backups não fixados mais recentes são mantidos |
| Backups fixados | Nunca excluídos | Sobrevivem à limpeza independentemente da contagem |
| Duplicatas | Ignoradas | Se a configuração não mudou, nenhum novo backup é criado |
| Compactação | Sempre | Novos backups usam tar.gz (~75% menores) |

## Fixando backups

Você pode marcar qualquer backup como "fixado" na TUI para protegê-lo da limpeza automática:

1.  Execute o `kortex` e navegue até a tela de **Backups**.
2.  Use `j`/`k` para selecionar um backup.
3.  Pressione **`p`** para alternar entre fixar/desafixar.
4.  Backups fixados exibem o indicador `[fixado]`.

Backups fixados nunca são excluídos automaticamente, mesmo quando o limite de retenção é excedido.

## Gerenciando backups (TUI)

| Tecla | Ação |
|-----|--------|
| `j` / `k` | Navegar para cima/baixo |
| `Enter` | Restaurar backup selecionado |
| `p` | Fixar/desafixar (protege da limpeza) |
| `r` | Renomear (adicionar uma descrição) |
| `d` | Excluir |
| `Esc` | Voltar |

## Comportamento da restauração

- Se `existed=true`: restaura o arquivo do snapshot para seu caminho original.
- Se `existed=false`: remove o arquivo (revertendo arquivos criados durante a instalação).
- A restauração é atômica por gravação de arquivo — sem restaurações parciais.
- Funciona tanto com backups compactados (tar.gz) quanto com legados (arquivos simples).

## Se a verificação falhar

1.  Revise as verificações que falharam no relatório de verificação.
2.  Restaure a partir do último snapshot via TUI ou com `kortex restore latest`.
3.  Execute a instalação novamente com `--dry-run` para validar o plano.
4.  Execute a instalação real após corrigir as dependências externas.

## O que o rollback NÃO cobre

- Pacotes instalados via `brew install`, `apt-get install` ou `pacman -S` não são desinstalados durante o rollback. O sistema de snapshot lida apenas com arquivos de configuração.
- Se você precisar desfazer a instalação de um pacote, use o gerenciador de pacotes da sua plataforma diretamente (ex: `brew uninstall`, `sudo apt-get remove`, `sudo pacman -R`).
