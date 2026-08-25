package filter

import "fmt"

func (b *Bank) Close(id string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	bed, ok := b.beds[id]
	if !ok {
		return fmt.Errorf("filter bed %s not found", id)
	}
	bed.Closed = true
	return nil
}

func (b *Bank) Open(id string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	bed, ok := b.beds[id]
	if !ok {
		return fmt.Errorf("filter bed %s not found", id)
	}
	bed.Closed = false
	return nil
}

func (b *Bank) IsClosed(id string) (bool, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	bed, ok := b.beds[id]
	if !ok {
		return false, fmt.Errorf("filter bed %s not found", id)
	}
	return bed.Closed, nil
}
