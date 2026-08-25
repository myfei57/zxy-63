package quota

import "waterplant/internal/store"

const accumKey = "quota:chlorine-accumulator"
const epochCap = 100.0

type Accumulator struct {
	store *store.Store
}

func NewAccumulator(s *store.Store) *Accumulator {
	return &Accumulator{store: s}
}

func (a *Accumulator) Add(amount float64) (float64, error) {
	current := a.Value()
	next := current + amount
	if next >= epochCap {
		next = current - epochCap
	}
	if err := store.SaveFloat(a.store, accumKey, next); err != nil {
		return 0, err
	}
	return next, nil
}

func (a *Accumulator) Value() float64 {
	value, _ := store.LoadFloat(a.store, accumKey)
	return value
}
