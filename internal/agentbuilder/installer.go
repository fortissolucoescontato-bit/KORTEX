package agentbuilder

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fortissolucoescontato-bit/kortex/internal/model"
)

// AdapterInfo pairs an AgentID with the path to its skills directory.
type AdapterInfo struct {
	AgentID   model.AgentID
	SkillsDir string
}

// Install writes the SKILL.md for agent into each adapter's skills directory.
// On any write failure all previously written files are rolled back (deleted).
// Returns one InstallResult per adapter.
func Install(ctx context.Context, agent *GeneratedAgent, adapters []AdapterInfo, _ string) ([]InstallResult, error) {
	if agent == nil {
		return nil, fmt.Errorf("instalação: o agente não pode ser nulo")
	}

	results := make([]InstallResult, 0, len(adapters))
	writtenFiles := make([]string, 0, len(adapters))
	createdDirs := make([]string, 0, len(adapters))

	for _, adapter := range adapters {
		if err := ctx.Err(); err != nil {
			rollback(writtenFiles, createdDirs)
			markAllFailed(results)
			return results, fmt.Errorf("instalação cancelada: %w", err)
		}
		skillDir := filepath.Join(adapter.SkillsDir, agent.Name)
		skillFile := filepath.Join(skillDir, "SKILL.md")

		// Check if directory exists before trying to create it
		_, statErr := os.Stat(skillDir)
		if os.IsNotExist(statErr) {
			if err := os.MkdirAll(skillDir, 0755); err != nil {
				rollback(writtenFiles, createdDirs)
				markAllFailed(results)
				results = append(results, InstallResult{
					AgentID: adapter.AgentID,
					Path:    skillFile,
					Success: false,
					Err:     fmt.Errorf("falha ao criar diretório %s: %w", skillDir, err),
				})
				return results, fmt.Errorf("falha na instalação para %s: %w", adapter.AgentID, err)
			}
			createdDirs = append(createdDirs, skillDir)
		}

		if err := os.WriteFile(skillFile, []byte(agent.Content), 0644); err != nil {
			rollback(writtenFiles, createdDirs)
			markAllFailed(results)
			results = append(results, InstallResult{
				AgentID: adapter.AgentID,
				Path:    skillFile,
				Success: false,
				Err:     fmt.Errorf("falha ao gravar %s: %w", skillFile, err),
			})
			return results, fmt.Errorf("falha na instalação para %s: %w", adapter.AgentID, err)
		}

		writtenFiles = append(writtenFiles, skillFile)
		results = append(results, InstallResult{
			AgentID: adapter.AgentID,
			Path:    skillFile,
			Success: true,
		})
	}

	return results, nil
}

// rollback removes all files and directories in paths, ignoring errors (best-effort cleanup).
func rollback(files, dirs []string) {
	for _, p := range files {
		_ = os.Remove(p)
	}
	// Remove directories in reverse order of creation.
	for i := len(dirs) - 1; i >= 0; i-- {
		_ = os.Remove(dirs[i])
	}
}

// markAllFailed sets Success=false on every result in the slice.
// Called after a rollback so previously-succeeded results reflect the true outcome.
func markAllFailed(results []InstallResult) {
	for i := range results {
		results[i].Success = false
	}
}
