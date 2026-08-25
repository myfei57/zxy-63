package turb

import "fmt"

type TurbState struct {
	LastVerdict float64 `json:"last_verdict"`
}

func (s *Sampler) State() TurbState {
	return TurbState{LastVerdict: s.last}
}

func (s *Sampler) Describe() string {
	return fmt.Sprintf("turbidity last_verdict=%.4f", s.last)
}
