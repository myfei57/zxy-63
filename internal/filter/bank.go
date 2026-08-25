package filter

import (
	"sort"
	"sync"
)

type Bed struct {
	ID     string
	Zone   int
	Load   float64
	Closed bool
	Duty   bool
}

type Bank struct {
	mu   sync.Mutex
	beds map[string]*Bed
}

func NewBank() *Bank {
	return &Bank{beds: map[string]*Bed{}}
}

func (b *Bank) AddBed(id string, zone int, load float64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.beds[id] = &Bed{ID: id, Zone: zone, Load: load}
}

func (b *Bank) Bed(id string) (*Bed, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	bed, ok := b.beds[id]
	return bed, ok
}

func (b *Bank) Beds() []*Bed {
	b.mu.Lock()
	defer b.mu.Unlock()
	beds := make([]*Bed, 0, len(b.beds))
	for _, bed := range b.beds {
		beds = append(beds, bed)
	}
	sort.Slice(beds, func(i, j int) bool { return beds[i].ID < beds[j].ID })
	return beds
}
