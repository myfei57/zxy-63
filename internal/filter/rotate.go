package filter

import "fmt"

func (b *Bank) Rotate(id string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	bed, ok := b.beds[id]
	if !ok {
		return fmt.Errorf("filter bed %s not found", id)
	}
	for _, candidate := range b.beds {
		candidate.Duty = false
	}
	bed.Duty = true
	return nil
}

func (b *Bank) OnDuty() (string, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for id, bed := range b.beds {
		if bed.Duty {
			return id, true
		}
	}
	return "", false
}
