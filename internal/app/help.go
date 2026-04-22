package app

import (
	"fmt"
	"io"
)

func printHelp(w io.Writer, version string) {
	fmt.Fprintf(w, `kortex — Orquestrador de Elite Kortex (%s)

USO
  kortex                     Iniciar interface interativa (TUI)
  kortex <comando> [flags]

COMANDOS
  install      Configurar agentes de IA nesta máquina
  uninstall    Remover arquivos gerenciados pelo Kortex desta máquina
  sync         Sincronizar configs e skills para a versão atual
  update       Verificar atualizações disponíveis
  upgrade      Aplicar atualizações às ferramentas gerenciadas
  restore      Restaurar um backup de configuração
  version      Exibir versão

FLAGS
  --help, -h    Exibir esta ajuda

Execute 'kortex help' para ver esta mensagem.
Documentação: https://github.com/fortissolucoescontato-bit/kortex
`, version)
}
