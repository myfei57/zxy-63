package turb

import (
	"waterplant/internal/coag"
	"waterplant/internal/intake"
)

type Sampler struct {
	doser *coag.Doser
	last  float64
}

func NewSampler(d *coag.Doser) *Sampler {
	return &Sampler{doser: d}
}

func (s *Sampler) Judge(raw []float64) (float64, error) {
	mixed := intake.Mix(raw)
	verdict := verdictFor(mixed)
	s.last = verdict
	return s.doser.DoseForTurbidity(verdict), nil
}

func verdictFor(turbidity float64) float64 {
	if turbidity < 1.0 {
		return 0.0
	}
	return turbidity
}
