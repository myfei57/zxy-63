package verifycase

import (
	"path/filepath"
	"testing"

	"waterplant/internal/chlor"
	"waterplant/internal/clearwell"
	"waterplant/internal/store"
)

func TestWpChlorDoseFreshTarget(t *testing.T) {
	s, err := store.Open(filepath.Join(t.TempDir(), "store.json"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	well := clearwell.NewWell(s)
	if err := well.SetResidualTarget(0.5); err != nil {
		t.Fatal(err)
	}
	doser := chlor.NewDoser(s)
	if err := well.UpdateResidualDemand(1.0); err != nil {
		t.Fatal(err)
	}
	dose := doser.DoseForResidual()
	if dose != 1.0 {
		t.Fatalf("chlorine dose uses stale target: got %.4f, want 1.0000", dose)
	}
}
