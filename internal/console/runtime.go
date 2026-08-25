package console

import (
	"waterplant/internal/audit"
	"waterplant/internal/backwash"
	"waterplant/internal/chlor"
	"waterplant/internal/clearwell"
	"waterplant/internal/coag"
	"waterplant/internal/filter"
	"waterplant/internal/flow"
	"waterplant/internal/intake"
	"waterplant/internal/quota"
	"waterplant/internal/store"
	"waterplant/internal/turb"
)

type Runtime struct {
	store      *store.Store
	flowRepo   *intake.FlowRepository
	inlet      *intake.InletController
	outlet     *intake.InletController
	coagDoser  *coag.Doser
	chlorDoser *chlor.Doser
	bank       *filter.Bank
	backwash   *backwash.Controller
	sampler    *turb.Sampler
	calib      *flow.Calibration
	well       *clearwell.Well
	acc        *quota.Accumulator
	audit      *audit.Auditor
}

func NewRuntime(s *store.Store) *Runtime {
	bank := filter.NewBank()
	coagDoser := coag.NewDoser(s)
	return &Runtime{
		store:      s,
		flowRepo:   intake.NewFlowRepository(s),
		inlet:      intake.NewInletController(),
		outlet:     intake.NewInletController(),
		coagDoser:  coagDoser,
		chlorDoser: chlor.NewDoser(s),
		bank:       bank,
		backwash:   backwash.NewController(bank, s),
		sampler:    turb.NewSampler(coagDoser),
		calib:      flow.NewCalibration(s),
		well:       clearwell.NewWell(s),
		acc:        quota.NewAccumulator(s),
		audit:      audit.NewAuditor(s),
	}
}
