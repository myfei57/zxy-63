package flow

import "waterplant/internal/store"

const calibrationKey = "flow:calibration"

type Calibration struct {
	store *store.Store
}

func NewCalibration(s *store.Store) *Calibration {
	return &Calibration{store: s}
}

func (c *Calibration) Current() float64 {
	value, ok := store.LoadFloat(c.store, calibrationKey)
	if !ok {
		return 1.0
	}
	return value
}

func (c *Calibration) Replace(factor float64) error {
	return store.SaveFloat(c.store, calibrationKey, factor)
}
