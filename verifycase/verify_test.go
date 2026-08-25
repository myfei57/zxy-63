package verifycase

import (
	"path/filepath"
	"testing"

	"waterplant/internal/quota"
	"waterplant/internal/store"
)

func TestWpChlorEpochRollover(t *testing.T) {
	s, err := store.Open(filepath.Join(t.TempDir(), "store.json"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	acc := quota.NewAccumulator(s)
	if _, err := acc.Add(98); err != nil {
		t.Fatal(err)
	}
	value, err := acc.Add(5)
	if err != nil {
		t.Fatal(err)
	}
	if value < 0 {
		t.Fatalf("chlorine accumulator wrapped negative: got %.4f", value)
	}
}
