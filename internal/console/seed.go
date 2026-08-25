package console

import (
	"waterplant/internal/clearwell"
	"waterplant/internal/flow"
	"waterplant/internal/intake"
	"waterplant/internal/quota"
	"waterplant/internal/store"
)

func Seed(s *store.Store) error {
	calib := flow.NewCalibration(s)
	if calib.Current() <= 0 {
		if err := calib.Replace(1.0); err != nil {
			return err
		}
	}
	well := clearwell.NewWell(s)
	if well.ResidualTarget() <= 0 {
		if err := well.SetResidualTarget(0.5); err != nil {
			return err
		}
	}
	repo := intake.NewFlowRepository(s)
	if _, ok := repo.LoadFlow(); !ok {
		if err := repo.PersistFlow(1.0); err != nil {
			return err
		}
	}
	acc := quota.NewAccumulator(s)
	if acc.Value() < 0 {
		if err := acc.Reset(); err != nil {
			return err
		}
	}
	return store.AppendEvent(s, "console:events", "seed", "defaults")
}
