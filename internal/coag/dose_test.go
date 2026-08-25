package coag

import (
	"path/filepath"
	"testing"

	"waterplant/internal/flow"
	"waterplant/internal/intake"
	"waterplant/internal/store"
)

// openStore opens the store at path and seeds the persisted flow with seedFlow
// (the value landed before the incident). Returns a Doser bound to it so a
// later Close+Open simulates a service restart.
func openStore(t *testing.T, path string, seedFlow float64) (*Doser, *store.Store) {
	t.Helper()
	s, err := store.Open(path)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	if seedFlow > 0 {
		if err := intake.NewFlowRepository(s).PersistFlow(seedFlow); err != nil {
			t.Fatalf("seed flow: %v", err)
		}
	}
	return NewDoser(s), s
}

// TestUpdateFlowAndDosePersistsBeforeDosing reproduces the reported incident:
// flow was 1200, then bumped to 1500. The dosing side must read 1500 (the
// latest landed value), not the stale 1200 — otherwise the coagulant
// underdoses and effluent turbidity climbs.
func TestUpdateFlowAndDosePersistsBeforeDosing(t *testing.T) {
	path := filepath.Join(t.TempDir(), "waterplant.json")

	// Pre-restart state: flow landed at 1200.
	_, s := openStore(t, path, 1200)
	if err := s.Close(); err != nil {
		t.Fatalf("close store: %v", err)
	}

	// "Restart": reopen the same file, then update flow to 1500 and compute dose.
	doser2, s2 := openStore(t, path, 0) // reopen; do not reseed, value on disk is 1200
	defer s2.Close()

	got, err := doser2.UpdateFlowAndDose(1500)
	if err != nil {
		t.Fatalf("UpdateFlowAndDose: %v", err)
	}

	// The dose must be derived from the newly landed 1500, not the stale
	// 1200. With the default calibration ratio of 1.0 the dose equals the flow.
	if want := 1500.0; got != want {
		t.Fatalf("dose = %v, want %v (dosing side read stale flow)", got, want)
	}

	// After the call the persisted flow must be 1500 so a subsequent restart
	// still sees the latest value.
	persisted, ok := doser2.flow.LoadFlow()
	if !ok || persisted != 1500 {
		t.Fatalf("persisted flow = %v ok=%v, want 1500", persisted, ok)
	}
}

// TestUpdateFlowAndDoseSurvivesRestart asserts the core requirement directly:
// after updating flow and reopening the store, a fresh Doser reads the
// latest landed value.
func TestUpdateFlowAndDoseSurvivesRestart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "waterplant.json")

	doser, s := openStore(t, path, 1200)
	if _, err := doser.UpdateFlowAndDose(1500); err != nil {
		t.Fatalf("UpdateFlowAndDose: %v", err)
	}
	if err := s.Close(); err != nil {
		t.Fatalf("close store: %v", err)
	}

	// Reopen — simulates the service restart from the incident report.
	s2, err := store.Open(path)
	if err != nil {
		t.Fatalf("reopen store: %v", err)
	}
	defer s2.Close()
	restarted := NewDoser(s2)

	// The dosing side must now read the freshly landed 1500, not 1200.
	dose := restarted.DoseFromPersistedFlow()
	if want := 1500.0; dose != want {
		t.Fatalf("post-restart dose = %v, want %v (stale value survived)", dose, want)
	}
}

// Ensure the imported flow package's default calibration path is exercised
// (ratio defaults to 1.0 when no calibration is stored), matching the
// assertion above that dose == flow.
func TestDefaultRatioIsOne(t *testing.T) {
	s, err := store.Open(filepath.Join(t.TempDir(), "waterplant.json"))
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	defer s.Close()
	d := NewDoser(s)
	if r := d.CurrentRatio(); r != flow.NewCalibration(s).Current() || r != 1.0 {
		t.Fatalf("default ratio = %v, want 1.0", r)
	}
}
