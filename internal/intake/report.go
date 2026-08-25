package intake

import (
	"fmt"
	"strconv"
)

const turbidityKey = "intake:turbidity"

type FlowState struct {
	Flow      float64 `json:"flow"`
	Turbidity float64 `json:"turbidity"`
	Present   bool    `json:"present"`
}

func (r *FlowRepository) Record(sensor Sensor) error {
	if err := r.PersistFlow(sensor.FlowValue); err != nil {
		return err
	}
	return r.store.Put(turbidityKey, strconv.FormatFloat(sensor.Turbidity, 'f', 4, 64))
}

func (r *FlowRepository) Turbidity() float64 {
	raw, ok := r.store.Get(turbidityKey)
	if !ok {
		return 0
	}
	value, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0
	}
	return value
}

func (r *FlowRepository) State() FlowState {
	flow, ok := r.LoadFlow()
	return FlowState{Flow: flow, Turbidity: r.Turbidity(), Present: ok}
}

func ValidateFlow(value float64) error {
	if value < 0 {
		return fmt.Errorf("flow must be non-negative")
	}
	if value > 1000000 {
		return fmt.Errorf("flow exceeds the supported range")
	}
	return nil
}

func (r *FlowRepository) Describe() string {
	state := r.State()
	return fmt.Sprintf("intake flow=%.4f turbidity=%.4f present=%t", state.Flow, state.Turbidity, state.Present)
}
