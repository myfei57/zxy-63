package verifycase

import (
	"path/filepath"
	"testing"

	"waterplant/internal/clearwell"
	"waterplant/internal/intake"
	"waterplant/internal/store"
)

func TestWpClearwellLevelOrder(t *testing.T) {
	s, err := store.Open(filepath.Join(t.TempDir(), "store.json"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	well := clearwell.NewWell(s)
	if err := well.SetLevel(10); err != nil {
		t.Fatal(err)
	}
	inlet := intake.NewInletController()
	outlet := intake.NewInletController()
	minSeen, err := well.AdjustLevel(12, inlet, outlet)
	if err != nil {
		t.Fatal(err)
	}
	if minSeen < 10.0 {
		t.Fatalf("clear well level dipped during adjustment: got %.4f", minSeen)
	}
}
