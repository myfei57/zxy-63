package verifycase

import (
	"path/filepath"
	"testing"

	"waterplant/internal/backwash"
	"waterplant/internal/filter"
	"waterplant/internal/store"
)

func TestWpBackwashOrder(t *testing.T) {
	s, err := store.Open(filepath.Join(t.TempDir(), "store.json"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	bank := filter.NewBank()
	bank.AddBed("A", 1, 0)
	controller := backwash.NewController(bank, s)
	if err := controller.Start("A"); err != nil {
		t.Fatal(err)
	}
	if _, ok := s.Get("backwash:spill"); ok {
		t.Fatalf("filter bed drained while open")
	}
}
