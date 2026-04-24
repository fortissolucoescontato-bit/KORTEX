package system

import (
	"os"
	"os/exec"
)

// Functional Overrides para isolamento de testes.
// Estas variáveis permitem que suítes de teste interceptem chamadas ao SO
// sem depender do estado global ou de permissões do host.
var (
	Stat        = os.Stat
	LookPath    = exec.LookPath
	ReadFile    = os.ReadFile
	UserHomeDir = os.UserHomeDir
	Command     = exec.Command
)
