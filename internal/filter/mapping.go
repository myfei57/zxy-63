package filter

import "fmt"

func (b *Bank) Renumber(id string, zone int) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	bed, ok := b.beds[id]
	if !ok {
		return fmt.Errorf("filter bed %s not found", id)
	}
	bed.Zone = zone
	return nil
}

func (b *Bank) Mapping() map[int]string {
	b.mu.Lock()
	defer b.mu.Unlock()
	mapping := map[int]string{}
	for id, bed := range b.beds {
		mapping[bed.Zone] = id
	}
	return mapping
}
