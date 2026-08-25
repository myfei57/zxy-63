package verifycase

import (
	"path/filepath"
	"testing"

	"waterplant/internal/backwash"
	"waterplant/internal/filter"
	"waterplant/internal/store"
)

func TestWpFilterZoneMappingFresh(t *testing.T) {
	s, err := store.Open(filepath.Join(t.TempDir(), "store.json"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	bank := filter.NewBank()
	bank.AddBed("A", 4, 0)
	controller := backwash.NewController(bank, s)
	if err := bank.Renumber("A", 3); err != nil {
		t.Fatal(err)
	}
	bank.AddBed("B", 4, 0)
	bedID, err := controller.Select(4)
	if err != nil {
		t.Fatal(err)
	}
	if bedID != "B" {
		t.Fatalf("backwash uses stale zone mapping: got %s, want B", bedID)
	}
}
