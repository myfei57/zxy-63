package chlor

import "fmt"

type ChlorState struct {
	Target float64 `json:"target"`
	Dose   float64 `json:"dose"`
}

func (d *Doser) State() ChlorState {
	return ChlorState{Target: d.CurrentTarget(), Dose: d.DoseForResidual()}
}

func ValidateDemand(demand float64) error {
	if demand < 0 {
		return fmt.Errorf("demand must be non-negative")
	}
	if demand > 1000 {
		return fmt.Errorf("demand exceeds the supported range")
	}
	return nil
}

func (d *Doser) Describe() string {
	state := d.State()
	return fmt.Sprintf("chlorine target=%.4f dose=%.4f", state.Target, state.Dose)
}

func (d *Doser) TargetForDemand(demand float64) float64 {
	return 0.3 + demand*0.7
}
