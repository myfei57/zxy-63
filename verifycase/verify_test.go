package verifycase

import (
	"path/filepath"
	"testing"

	"waterplant/internal/coag"
	"waterplant/internal/intake"
	"waterplant/internal/store"
)

func TestWpCoagAfterFlowPersist(t *testing.T) {
	s, err := store.Open(filepath.Join(t.TempDir(), "store.json"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	repo := intake.NewFlowRepository(s)
	if err := repo.PersistFlow(10); err != nil {
		t.Fatal(err)
	}
	doser := coag.NewDoser(s)
	dose, err := doser.UpdateFlowAndDose(20)
	if err != nil {
		t.Fatal(err)
	}
	if dose != 20.0 {
		t.Fatalf("dose follows stale flow: got %.4f, want 20.0000", dose)
	}
}
