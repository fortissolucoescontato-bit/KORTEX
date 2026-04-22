package cli

import (
	"flag"
	"fmt"
	"strings"
)

type InstallFlags struct {
	Agents     []string
	Components []string
	Skills     []string
	Persona    string
	Preset     string
	SDDMode    string
	DryRun     bool
}

func ParseInstallFlags(args []string) (InstallFlags, error) {
	var opts InstallFlags

	fs := flag.NewFlagSet("install", flag.ContinueOnError)
	fs.SetOutput(ioDiscard{})
	registerListFlag(fs, "agent", &opts.Agents)
	registerListFlag(fs, "agents", &opts.Agents)
	registerListFlag(fs, "component", &opts.Components)
	registerListFlag(fs, "components", &opts.Components)
	registerListFlag(fs, "skill", &opts.Skills)
	registerListFlag(fs, "skills", &opts.Skills)
	fs.StringVar(&opts.Persona, "persona", "", "persona a ser aplicada")
	fs.StringVar(&opts.Preset, "preset", "", "preset a ser aplicado")
	fs.StringVar(&opts.SDDMode, "sdd-mode", "", "modo do orquestrador SDD: single ou multi (padrão: single)")
	fs.BoolVar(&opts.DryRun, "dry-run", false, "pré-visualizar o plano sem executar")

	if err := fs.Parse(args); err != nil {
		return InstallFlags{}, err
	}

	if fs.NArg() > 0 {
		return InstallFlags{}, fmt.Errorf("unexpected install argument %q", fs.Arg(0))
	}

	return opts, nil
}

type csvListFlag struct {
	values *[]string
}

func (f csvListFlag) String() string {
	if f.values == nil {
		return ""
	}

	return strings.Join(*f.values, ",")
}

func (f csvListFlag) Set(value string) error {
	for _, part := range strings.Split(value, ",") {
		item := strings.TrimSpace(part)
		if item == "" {
			continue
		}
		*f.values = append(*f.values, item)
	}

	return nil
}

func registerListFlag(fs *flag.FlagSet, name string, values *[]string) {
	fs.Var(csvListFlag{values: values}, name, "lista separada por vírgulas")
}

type ioDiscard struct{}

func (ioDiscard) Write(p []byte) (int, error) {
	return len(p), nil
}
