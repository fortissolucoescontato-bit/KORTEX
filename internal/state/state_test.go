package state

import (
	"reflect"
	"testing"
)

func TestManager_InstalledAgents(t *testing.T) {
	home := t.TempDir()
	m, err := NewManager(home)
	if err != nil {
		t.Fatalf("NewManager() error = %v", err)
	}
	defer m.Close()

	agents := []string{"claude-code", "opencode"}
	if err := m.SetInstalledAgents(agents); err != nil {
		t.Fatalf("SetInstalledAgents() error = %v", err)
	}

	got, err := m.GetInstalledAgents()
	if err != nil {
		t.Fatalf("GetInstalledAgents() error = %v", err)
	}

	if !reflect.DeepEqual(got, agents) {
		t.Errorf("GetInstalledAgents() = %v, want %v", got, agents)
	}
}

func TestManager_Assignments(t *testing.T) {
	home := t.TempDir()
	m, err := NewManager(home)
	if err != nil {
		t.Fatalf("NewManager() error = %v", err)
	}
	defer m.Close()

	agent := "claude"
	phase := "orchestrator"
	provider := "anthropic"
	model := "claude-3-opus-20240229"

	if err := m.SetAssignment(agent, phase, provider, model); err != nil {
		t.Fatalf("SetAssignment() error = %v", err)
	}

	assignments, err := m.GetAssignments(agent)
	if err != nil {
		t.Fatalf("GetAssignments() error = %v", err)
	}

	got, ok := assignments[phase]
	if !ok {
		t.Fatalf("assignment for phase %q not found", phase)
	}

	if got.ProviderID != provider || got.ModelID != model {
		t.Errorf("assignment = %+v, want {ProviderID: %q, ModelID: %q}", got, provider, model)
	}
}

func TestManager_Idempotency(t *testing.T) {
	home := t.TempDir()
	m, err := NewManager(home)
	if err != nil {
		t.Fatalf("NewManager() error = %v", err)
	}
	defer m.Close()

	agents := []string{"claude-code"}
	if err := m.SetInstalledAgents(agents); err != nil {
		t.Fatalf("SetInstalledAgents() first error = %v", err)
	}

	if err := m.SetInstalledAgents(agents); err != nil {
		t.Fatalf("SetInstalledAgents() second error = %v", err)
	}

	got, _ := m.GetInstalledAgents()
	if len(got) != 1 {
		t.Errorf("len(GetInstalledAgents()) = %d, want 1", len(got))
	}
}
