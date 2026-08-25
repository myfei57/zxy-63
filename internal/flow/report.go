package flow

import "fmt"

type CalibState struct {
	Factor float64 `json:"factor"`
}

func (c *Calibration) State() CalibState {
	return CalibState{Factor: c.Current()}
}

func ValidateFactor(factor float64) error {
	if factor <= 0 {
		return fmt.Errorf("factor must be positive")
	}
	if factor > 100 {
		return fmt.Errorf("factor exceeds the supported range")
	}
	return nil
}

func (c *Calibration) Describe() string {
	return fmt.Sprintf("flow calibration factor=%.4f", c.Current())
}
