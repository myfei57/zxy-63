package store

import "fmt"

type StoreState struct {
	Keys []string `json:"keys"`
	Size int      `json:"size"`
}

func (s *Store) State() StoreState {
	keys := s.Keys()
	return StoreState{Keys: keys, Size: len(keys)}
}

func (s *Store) Describe() string {
	return fmt.Sprintf("store keys=%d", len(s.Keys()))
}

func (s *Store) Export() map[string]string {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make(map[string]string, len(s.data))
	for key, value := range s.data {
		out[key] = value
	}
	return out
}

func (s *Store) Count() int {
	return len(s.Keys())
}

func (s *Store) Has(key string) bool {
	_, ok := s.Get(key)
	return ok
}

func (s *Store) Clear() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = map[string]string{}
	return s.saveLocked()
}
