package verifycase

import (
	"path/filepath"
	"testing"

	"waterplant/internal/backwash"
	"waterplant/internal/filter"
	"waterplant/internal/store"
)

func TestWpFilterDutyOrder(t *testing.T) {
	s, err := store.Open(filepath.Join(t.TempDir(), "store.json"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	bank := filter.NewBank()
	bank.AddBed("A", 1, 1)
	bank.AddBed("B", 2, 5)
	bank.AddBed("C", 3, 3)
	controller := backwash.NewController(bank, s)
	order, err := controller.OrderRotation()
	if err != nil {
		t.Fatal(err)
	}
	if order[0] != "B" {
		t.Fatalf("dirtiest filter not rotated first: %v", order)
	}
}
