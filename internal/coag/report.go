package coag

import "fmt"

type DoseState struct {
	Ratio         float64 `json:"ratio"`
	PersistedFlow float64 `json:"persisted_flow"`
	FlowPresent   bool    `json:"flow_present"`
}

func (d *Doser) State() DoseState {
	flow, ok := d.flow.LoadFlow()
	return DoseState{Ratio: d.CurrentRatio(), PersistedFlow: flow, FlowPresent: ok}
}

func (d *Doser) Describe() string {
	state := d.State()
	return fmt.Sprintf("coagulant ratio=%.4f persisted_flow=%.4f present=%t", state.Ratio, state.PersistedFlow, state.FlowPresent)
}

func (d *Doser) DosePlan(flow, turbidity float64) float64 {
	return (flow*0.6 + turbidity*0.4) * d.CurrentRatio()
}
