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
}

func NewDoser(s *store.Store) *Doser {
	return &Doser{
		flow:  intake.NewFlowRepository(s),
		calib: flow.NewCalibration(s),
		audit: audit.NewAuditor(s),
	}
}
