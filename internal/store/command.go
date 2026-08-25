package store

import "encoding/json"

func AppendCommand(s *Store, listKey, value string) error {
	items := LoadCommands(s, listKey)
	items = append(items, value)
	raw, err := json.Marshal(items)
	if err != nil {
		return err
	}
	return s.Put(listKey, string(raw))
}

func LoadCommands(s *Store, listKey string) []string {
	raw, ok := s.Get(listKey)
	if !ok {
		return []string{}
	}
	var items []string
	if err := json.Unmarshal([]byte(raw), &items); err != nil {
		return []string{}
	}
	return items
}

func ClearCommands(s *Store, listKey string) error {
	return s.Delete(listKey)
}
