package store

import "strconv"

func LoadFloat(s *Store, key string) (float64, bool) {
	raw, ok := s.Get(key)
	if !ok {
		return 0, false
	}
	value, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0, false
	}
	return value, true
}

func SaveFloat(s *Store, key string, value float64) error {
	return s.Put(key, strconv.FormatFloat(value, 'f', 4, 64))
}

func AdvanceEpoch(s *Store, key string, amount, cap float64) (float64, error) {
	current, _ := LoadFloat(s, key)
	next := current + amount
	if cap > 0 && next >= cap {
		next -= cap
	}
	return next, SaveFloat(s, key, next)
}
