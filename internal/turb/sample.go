package turb

import (
	"waterplant/internal/coag"
)

type Sampler struct {
	doser *coag.Doser
	last  float64
}

func NewSampler(d *coag.Doser) *Sampler {
	return &Sampler{doser: d}
}

func (s *Sampler) Judge(raw []float64) (float64, error) {
	verdict := verdictFor(maxRaw(raw))
	s.last = verdict
	return s.doser.DoseForTurbidity(verdict), nil
}

func maxRaw(values []float64) float64 {
	highest := 0.0
	for _, value := range values {
		if value > highest {
			highest = value
		}
	}
	return highest
}

func verdictFor(turbidity float64) float64 {
	if turbidity < 1.0 {
		return 0.0
	}
	return turbidity
}
