package store

import (
	"encoding/json"
	"time"
)

type Event struct {
	At    int64  `json:"at"`
	Kind  string `json:"kind"`
	Value string `json:"value"`
}

func AppendEvent(s *Store, key, kind, value string) error {
	event := Event{At: time.Now().Unix(), Kind: kind, Value: value}
	raw, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return AppendCommand(s, key, string(raw))
}

func ListEvents(s *Store, key string) []Event {
	var events []Event
	for _, item := range LoadCommands(s, key) {
		var event Event
		if err := json.Unmarshal([]byte(item), &event); err == nil {
			events = append(events, event)
		}
	}
	return events
}

func EventCount(s *Store, key string) int {
	return len(LoadCommands(s, key))
}
