package filter

import "fmt"

type BedState struct {
	ID     string  `json:"id"`
	Zone   int     `json:"zone"`
	Load   float64 `json:"load"`
	Closed bool    `json:"closed"`
	Duty   bool    `json:"duty"`
}

type BankState struct {
	Beds   []BedState `json:"beds"`
	OnDuty string     `json:"on_duty"`
	Count  int        `json:"count"`
}

func (b *Bank) State() BankState {
	beds := b.Beds()
	states := make([]BedState, 0, len(beds))
	for _, bed := range beds {
		states = append(states, BedState{ID: bed.ID, Zone: bed.Zone, Load: bed.Load, Closed: bed.Closed, Duty: bed.Duty})
	}
	duty, _ := b.OnDuty()
	return BankState{Beds: states, OnDuty: duty, Count: len(states)}
}

func ValidateZone(zone int) error {
	if zone < 0 {
		return fmt.Errorf("zone must be non-negative")
	}
	if zone > 100000 {
		return fmt.Errorf("zone exceeds the supported range")
	}
	return nil
}

func (b *Bank) Describe() string {
	state := b.State()
	return fmt.Sprintf("filter beds=%d on_duty=%s", state.Count, state.OnDuty)
}

func (b *Bank) Dirtiest() string {
	var dirtiest string
	highest := -1.0
	for _, bed := range b.Beds() {
		if bed.Load > highest {
			highest = bed.Load
			dirtiest = bed.ID
		}
	}
	return dirtiest
}

func (b *Bank) Cleanest() string {
	var cleanest string
	lowest := -1.0
	for _, bed := range b.Beds() {
		if lowest < 0 || bed.Load < lowest {
			lowest = bed.Load
			cleanest = bed.ID
		}
	}
	return cleanest
}

func (b *Bank) SetLoad(id string, load float64) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	bed, ok := b.beds[id]
	if !ok {
		return fmt.Errorf("filter bed %s not found", id)
	}
	bed.Load = load
	return nil
}

func (b *Bank) ResetClosed() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, bed := range b.beds {
		bed.Closed = false
	}
}

func (b *Bank) Count() int {
	return len(b.Beds())
}

func (b *Bank) TotalLoad() float64 {
	total := 0.0
	for _, bed := range b.Beds() {
		total += bed.Load
	}
	return total
}

func (b *Bank) ActiveCount() int {
	count := 0
	for _, bed := range b.Beds() {
		if !bed.Closed {
			count++
		}
	}
	return count
}

func (b *Bank) BedZone(id string) (int, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	bed, ok := b.beds[id]
	if !ok {
		return 0, false
	}
	return bed.Zone, true
}

func (b *Bank) RemoveBed(id string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.beds[id]; !ok {
		return fmt.Errorf("filter bed %s not found", id)
	}
	delete(b.beds, id)
	return nil
}

func (b *Bank) BedIDs() []string {
	beds := b.Beds()
	ids := make([]string, 0, len(beds))
	for _, bed := range beds {
		ids = append(ids, bed.ID)
	}
	return ids
}
