package clearwell

import "fmt"

type WellState struct {
	Level          float64 `json:"level"`
	ResidualTarget float64 `json:"residual_target"`
}

func (w *Well) State() WellState {
	return WellState{Level: w.Level(), ResidualTarget: w.ResidualTarget()}
}

func ValidateLevel(level float64) error {
	if level < 0 {
		return fmt.Errorf("level must be non-negative")
	}
	if level > 1000000 {
		return fmt.Errorf("level exceeds the supported range")
	}
	return nil
}

func (w *Well) Describe() string {
	state := w.State()
	return fmt.Sprintf("clearwell level=%.4f residual_target=%.4f", state.Level, state.ResidualTarget)
}
