package quota

import (
	"fmt"

	"waterplant/internal/store"
)

type AccumState struct {
	Value float64 `json:"value"`
}

func (a *Accumulator) State() AccumState {
	return AccumState{Value: a.Value()}
}

func ValidateAmount(amount float64) error {
	if amount < 0 {
		return fmt.Errorf("amount must be non-negative")
	}
	if amount > 1000000 {
		return fmt.Errorf("amount exceeds the supported range")
	}
	return nil
}

func (a *Accumulator) Describe() string {
	return fmt.Sprintf("quota accumulator value=%.4f", a.Value())
}

func (a *Accumulator) Remaining(limit float64) float64 {
	remaining := limit - a.Value()
	if remaining < 0 {
		return 0
	}
	return remaining
}

func (a *Accumulator) Reset() error {
	return store.SaveFloat(a.store, accumKey, 0)
}
