# Testes E2E com Docker

Testes de ponta a ponta (E2E) que validam o binário do instalador `kortex` dentro de containers Docker rodando distribuições Linux reais.

## Arquitetura

```
e2e/
  lib.sh              # Helpers compartilhados: cores, contadores, logging, limpeza
  e2e_test.sh         # Todos os casos de teste, divididos por variáveis de ambiente
  Dockerfile.ubuntu   # Imagem de teste para Ubuntu 22.04
  Dockerfile.arch     # Imagem de teste para Arch Linux
  docker-test.sh      # Orquestrador: compila + executa todas as plataformas
```

## Início Rápido

```bash
# Executa apenas o Tier 1 (binário básico + testes de dry-run)
./e2e/docker-test.sh

# Executa todos os níveis (tiers)
RUN_FULL_E2E=1 RUN_BACKUP_TESTS=1 ./e2e/docker-test.sh
```

## Níveis de Teste (Tiers)

| Nível | Var. de Ambiente | O que testa |
|------|---------|---------------|
| 1 (padrão) | — | Se o binário existe e roda, formato de saída do dry-run, validação de flags |
| 2 | `RUN_FULL_E2E=1` | Instalação completa: opencode+permissions, claude-code+persona, context7, sdd |
| 3 | `RUN_BACKUP_TESTS=1` | Criação de snapshots de backup, conteúdo dos arquivos de backup |

## Plataformas Suportadas

| Plataforma | Dockerfile | Gerenciador de pacotes |
|----------|-----------|-----------------|
| Ubuntu 22.04 | `Dockerfile.ubuntu` | apt |
| Arch Linux | `Dockerfile.arch` | pacman |

## Como funciona

1.  **docker-test.sh** itera sobre a matriz de plataformas.
2.  Para cada plataforma, ele constrói uma imagem Docker que:
    - Instala dependências do sistema (git, curl, sudo, Go).
    - Cria um usuário não-root `testuser` com privilégios sudo sem senha.
    - Copia o código-fonte do projeto e compila o binário (`go build`).
    - Copia os scripts de teste E2E.
3.  Executa o container, que por sua vez roda o `e2e_test.sh` como `testuser`.
4.  Coleta os sucessos/falhas por plataforma e encerra com erro se houver alguma falha.

## Executando plataformas individuais

```bash
# Compila e executa apenas o Ubuntu
docker build -f e2e/Dockerfile.ubuntu -t kortex-e2e-ubuntu .
docker run --rm kortex-e2e-ubuntu

# Executa com E2E completo no Arch
docker build -f e2e/Dockerfile.arch -t kortex-e2e-arch .
docker run --rm -e RUN_FULL_E2E=1 kortex-e2e-arch

# Depuração interativa
docker run --rm -it kortex-e2e-ubuntu /bin/bash
```

## Adicionando uma nova plataforma

1.  Crie o `e2e/Dockerfile.<plataforma>` seguindo o padrão existente.
2.  Adicione a entrada ao array `PLATFORMS` no `docker-test.sh`.
3.  Garanta que o Dockerfile crie o `testuser` com sudo NOPASSWD.
4.  Compile o binário `kortex` para `linux/amd64`.

## Adicionando novos casos de teste

1.  Adicione uma função `test_*` ao `e2e_test.sh`.
2.  Use `log_test`, `log_pass`, `log_fail` do `lib.sh`.
3.  Chame `cleanup_test_env` antes de testes que gravam no sistema de arquivos.
4.  Coloque a chamada da função na seção do nível apropriado.
5.  Proteja os testes de Nível 2/3 com `RUN_FULL_E2E` / `RUN_BACKUP_TESTS`.

## Integração CI

```yaml
# Exemplo do GitHub Actions
jobs:
  e2e-linux:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Executar testes E2E (Tier 1)
        run: ./e2e/docker-test.sh
      - name: Executar testes E2E completos
        if: github.ref == 'refs/heads/main'
        run: RUN_FULL_E2E=1 RUN_BACKUP_TESTS=1 ./e2e/docker-test.sh
```
