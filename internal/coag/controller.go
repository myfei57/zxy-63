package coag

import (
	"waterplant/internal/audit"
	"waterplant/internal/flow"
	"waterplant/internal/intake"
	"waterplant/internal/store"
)

type Doser struct {
	flow  *intake.FlowRepository
	calib *flow.Calibration
	audit *audit.Auditor
	ratio float64
}

func NewDoser(s *store.Store) *Doser {
	calib := flow.NewCalibration(s)
	return &Doser{
		flow:  intake.NewFlowRepository(s),
		calib: calib,
		audit: audit.NewAuditor(s),
		ratio: calib.Current(),
	}
}
