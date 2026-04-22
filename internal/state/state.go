package state

import (
	"fmt"

	"github.com/fortissolucoescontato-bit/kortex/internal/storage"
)

// ModelAssignmentState is the serialisable form of a provider+model pair.
type ModelAssignmentState struct {
	ProviderID string
	ModelID    string
}

// Manager handles reading and writing the Kortex install state using SQLite.
type Manager struct {
	db *storage.DB
}

// NewManager creates a new state manager from a home directory.
func NewManager(homeDir string) (*Manager, error) {
	db, err := storage.Open(homeDir)
	if err != nil {
		return nil, err
	}
	return &Manager{db: db}, nil
}

// Close closes the underlying database.
func (m *Manager) Close() error {
	return m.db.Close()
}

// GetInstalledAgents retrieves the list of installed agent IDs.
func (m *Manager) GetInstalledAgents() ([]string, error) {
	rows, err := m.db.Query("SELECT agent_id FROM installed_agents")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agents []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		agents = append(agents, id)
	}
	return agents, nil
}

// SetInstalledAgents persists the list of installed agent IDs.
func (m *Manager) SetInstalledAgents(agents []string) error {
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec("DELETE FROM installed_agents"); err != nil {
		return err
	}

	for _, id := range agents {
		if _, err := tx.Exec("INSERT INTO installed_agents (agent_id) VALUES (?)", id); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetAssignments retrieves all model assignments for a specific agent (e.g. "claude").
func (m *Manager) GetAssignments(agentID string) (map[string]ModelAssignmentState, error) {
	rows, err := m.db.Query("SELECT phase, provider_id, model_id FROM model_assignments WHERE agent_id = ?", agentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	assignments := make(map[string]ModelAssignmentState)
	for rows.Next() {
		var phase, provider, model string
		if err := rows.Scan(&phase, &provider, &model); err != nil {
			return nil, err
		}
		assignments[phase] = ModelAssignmentState{
			ProviderID: provider,
			ModelID:    model,
		}
	}
	return assignments, nil
}

// SetAssignment persists a model assignment for a specific agent and phase.
func (m *Manager) SetAssignment(agentID, phase, providerID, modelID string) error {
	_, err := m.db.Exec(`
		INSERT OR REPLACE INTO model_assignments (agent_id, phase, provider_id, model_id)
		VALUES (?, ?, ?, ?)`,
		agentID, phase, providerID, modelID,
	)
	return err
}
