package verifycase

import (
	"path/filepath"
	"testing"

	"waterplant/internal/coag"
	"waterplant/internal/flow"
	"waterplant/internal/store"
)

func TestWpFlowRatioFresh(t *testing.T) {
	s, err := store.Open(filepath.Join(t.TempDir(), "store.json"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	doser := coag.NewDoser(s)
	calib := flow.NewCalibration(s)
	if err := calib.Replace(1.5); err != nil {
		t.Fatal(err)
	}
	dose := doser.DoseForFlow(10)
	if dose != 15.0 {
		t.Fatalf("coagulant ratio is stale: got %.4f, want 15.0000", dose)
	}
}
