package verifycase

import (
	"path/filepath"
	"testing"

	"waterplant/internal/backwash"
	"waterplant/internal/filter"
	"waterplant/internal/store"
)

func TestWpNoBackwashReplay(t *testing.T) {
	s, err := store.Open(filepath.Join(t.TempDir(), "store.json"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	bank := filter.NewBank()
	bank.AddBed("A", 1, 0)
	controller := backwash.NewController(bank, s)
	if err := controller.Enqueue("A"); err != nil {
		t.Fatal(err)
	}
	if err := controller.Recover(); err != nil {
		t.Fatal(err)
	}
	replayed, err := controller.Replay()
	if err != nil {
		t.Fatal(err)
	}
	if len(replayed) != 0 {
		t.Fatalf("stale backwash command replayed after recovery: %v", replayed)
	}
}
