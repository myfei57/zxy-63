package clearwell

import "waterplant/internal/store"

const levelKey = "clearwell:level"
const residualKey = "clearwell:residual-target"

type Well struct {
	store *store.Store
}

func NewWell(s *store.Store) *Well {
	return &Well{store: s}
}

func (w *Well) Level() float64 {
	value, _ := store.LoadFloat(w.store, levelKey)
	return value
}

func (w *Well) SetLevel(value float64) error {
	return store.SaveFloat(w.store, levelKey, value)
}

func (w *Well) ResidualTarget() float64 {
	value, ok := store.LoadFloat(w.store, residualKey)
	if !ok {
		return 0.5
	}
	return value
}

func (w *Well) SetResidualTarget(value float64) error {
	return store.SaveFloat(w.store, residualKey, value)
}
