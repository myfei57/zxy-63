package verifycase

import (
	"path/filepath"
	"testing"

	"waterplant/internal/coag"
	"waterplant/internal/store"
	"waterplant/internal/turb"
)

func TestWpTurbMixOrder(t *testing.T) {
	s, err := store.Open(filepath.Join(t.TempDir(), "store.json"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	doser := coag.NewDoser(s)
	sampler := turb.NewSampler(doser)
	dose, err := sampler.Judge([]float64{8, 1, 1})
	if err != nil {
		t.Fatal(err)
	}
	if dose >= 5.0 {
		t.Fatalf("turbidity judged before mixing: got %.4f", dose)
	}
}
