package main

import (
	"fmt"
	"os"

	"github.com/fortissolucoescontato-bit/kortex/internal/app"
)

// version is set by GoReleaser via ldflags at build time.
var version = "dev"

func main() {
	app.Version = app.ResolveVersion(version)

	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Erro: %v\n", err)
		os.Exit(1)
	}
}
