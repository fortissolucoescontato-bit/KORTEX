package verify

import (
	"fmt"
	"strings"
)

const readyMessage = "Você está pronto. Execute `claude` ou `opencode` e comece a construir. ⚡"

type Report struct {
	Checks    []CheckResult
	Passed    int
	Failed    int
	Skipped   int
	Warnings  int
	Ready     bool
	FinalNote string
}

func BuildReport(results []CheckResult) Report {
	report := Report{Checks: append([]CheckResult(nil), results...)}
	for _, result := range results {
		switch result.Status {
		case CheckStatusPassed:
			report.Passed++
		case CheckStatusFailed:
			report.Failed++
		case CheckStatusSkipped:
			report.Skipped++
		case CheckStatusWarning:
			report.Warnings++
		}
	}

	report.Ready = report.Failed == 0
	if report.Ready {
		report.FinalNote = readyMessage
	} else {
		report.FinalNote = "A instalação foi concluída com problemas de verificação. Corrija as falhas indicadas."
	}

	return report
}

func RenderReport(report Report) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Verificações de integridade: %d passaram, %d falharam, %d avisos, %d puladas\n", report.Passed, report.Failed, report.Warnings, report.Skipped)

	for _, check := range report.Checks {
		line := "[ ]"
		switch check.Status {
		case CheckStatusPassed:
			line = "[ok]"
		case CheckStatusFailed:
			line = "[!!]"
		case CheckStatusWarning:
			line = "[??]"
		case CheckStatusSkipped:
			line = "[--]"
		}

		fmt.Fprintf(&b, "%s %s", line, check.ID)
		if check.Description != "" {
			fmt.Fprintf(&b, " - %s", check.Description)
		}
		if check.Error != "" {
			fmt.Fprintf(&b, " (%s)", check.Error)
		}
		b.WriteString("\n")
	}

	b.WriteString(report.FinalNote)
	if !strings.HasSuffix(report.FinalNote, "\n") {
		b.WriteString("\n")
	}

	return b.String()
}
