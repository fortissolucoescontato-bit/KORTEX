package planner

import (
	"errors"
	"reflect"
	"testing"

	"github.com/fortissolucoescontato-bit/kortex/internal/model"
)

func TestTopologicalSortOrdersDependenciesFirst(t *testing.T) {
	deps := map[model.ComponentID][]model.ComponentID{
		model.ComponentSkills:   {model.ComponentSDD},
		model.ComponentSDD:      {model.ComponentKortexEngram},
		model.ComponentKortexEngram:   nil,
		model.ComponentPersona:  nil,
		model.ComponentContext7: nil,
	}

	ordered, err := TopologicalSort(deps)
	if err != nil {
		t.Fatalf("TopologicalSort() returned error: %v", err)
	}

	if !reflect.DeepEqual(ordered, []model.ComponentID{
		model.ComponentContext7,
		model.ComponentKortexEngram,
		model.ComponentPersona,
		model.ComponentSDD,
		model.ComponentSkills,
	}) {
		t.Fatalf("TopologicalSort() order = %v", ordered)
	}
}

func TestApplySoftOrderingReordersWithoutAddingDependencies(t *testing.T) {
	ordered := []model.ComponentID{
		model.ComponentContext7,
		model.ComponentKortexEngram,
		model.ComponentPersona,
		model.ComponentSDD,
	}

	result := applySoftOrdering(ordered, [][2]model.ComponentID{{model.ComponentPersona, model.ComponentKortexEngram}})

	if !reflect.DeepEqual(result, []model.ComponentID{
		model.ComponentContext7,
		model.ComponentPersona,
		model.ComponentKortexEngram,
		model.ComponentSDD,
	}) {
		t.Fatalf("applySoftOrdering() = %v", result)
	}

	// If the first component is absent, nothing should be added.
	result = applySoftOrdering([]model.ComponentID{model.ComponentKortexEngram}, [][2]model.ComponentID{{model.ComponentPersona, model.ComponentKortexEngram}})
	if !reflect.DeepEqual(result, []model.ComponentID{model.ComponentKortexEngram}) {
		t.Fatalf("applySoftOrdering() should not add missing components (first absent), got %v", result)
	}
}

func TestApplySoftOrderingEdgeCases(t *testing.T) {
	pair := [][2]model.ComponentID{{model.ComponentPersona, model.ComponentKortexEngram}}

	// Second absent — no-op, no panic
	result := applySoftOrdering([]model.ComponentID{model.ComponentPersona}, pair)
	if !reflect.DeepEqual(result, []model.ComponentID{model.ComponentPersona}) {
		t.Fatalf("second absent: expected [persona], got %v", result)
	}

	// Both absent — no-op
	result = applySoftOrdering([]model.ComponentID{model.ComponentSDD}, pair)
	if !reflect.DeepEqual(result, []model.ComponentID{model.ComponentSDD}) {
		t.Fatalf("both absent: expected [sdd], got %v", result)
	}

	// Already correct order — no-op (must not mutate)
	already := []model.ComponentID{model.ComponentPersona, model.ComponentKortexEngram}
	result = applySoftOrdering(already, pair)
	if !reflect.DeepEqual(result, []model.ComponentID{model.ComponentPersona, model.ComponentKortexEngram}) {
		t.Fatalf("already correct: expected [persona, KortexEngram], got %v", result)
	}

	// Input slice must NOT be mutated
	input := []model.ComponentID{model.ComponentKortexEngram, model.ComponentPersona}
	_ = applySoftOrdering(input, pair)
	if !reflect.DeepEqual(input, []model.ComponentID{model.ComponentKortexEngram, model.ComponentPersona}) {
		t.Fatalf("input slice was mutated")
	}
}

func TestApplySoftOrderingBothMVPPairsWithFullSelection(t *testing.T) {
	// Simulates the real scenario: topo gives [context7, KortexEngram, persona, sdd, skills]
	// Both MVPGraph soft pairs should result in persona before KortexEngram AND sdd.
	ordered := []model.ComponentID{
		model.ComponentContext7,
		model.ComponentKortexEngram,
		model.ComponentPersona,
		model.ComponentSDD,
		model.ComponentSkills,
	}

	result := applySoftOrdering(ordered, SoftOrderingConstraints())

	// Persona must appear before both KortexEngram and SDD.
	personaIdx, KortexEngramIdx, sddIdx := -1, -1, -1
	for i, c := range result {
		switch c {
		case model.ComponentPersona:
			personaIdx = i
		case model.ComponentKortexEngram:
			KortexEngramIdx = i
		case model.ComponentSDD:
			sddIdx = i
		}
	}

	if personaIdx < 0 || KortexEngramIdx < 0 || sddIdx < 0 {
		t.Fatalf("missing components in result: %v", result)
	}
	if personaIdx > KortexEngramIdx {
		t.Fatalf("Persona (%d) must be before KortexEngram (%d), got %v", personaIdx, KortexEngramIdx, result)
	}
	if personaIdx > sddIdx {
		t.Fatalf("Persona (%d) must be before SDD (%d), got %v", personaIdx, sddIdx, result)
	}
	// Hard dep: KortexEngram must still be before SDD
	if KortexEngramIdx > sddIdx {
		t.Fatalf("KortexEngram (%d) must remain before SDD (%d) after soft reorder, got %v", KortexEngramIdx, sddIdx, result)
	}
	// Skills must remain last
	if result[len(result)-1] != model.ComponentSkills {
		t.Fatalf("Skills must remain last, got %v", result)
	}
}

func TestTopologicalSortDetectsCycles(t *testing.T) {
	deps := map[model.ComponentID][]model.ComponentID{
		model.ComponentKortexEngram: {model.ComponentSDD},
		model.ComponentSDD:    {model.ComponentKortexEngram},
	}

	_, err := TopologicalSort(deps)
	if err == nil {
		t.Fatalf("TopologicalSort() expected cycle error")
	}

	if !errors.Is(err, ErrDependencyCycle) {
		t.Fatalf("TopologicalSort() error = %v, want ErrDependencyCycle", err)
	}
}
