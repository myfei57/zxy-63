package audit

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"waterplant/internal/store"
)

const auditKey = "audit:entries"

type Auditor struct {
	store *store.Store
}

type Entry struct {
	ID     string    `json:"id"`
	Time   time.Time `json:"time"`
	Kind   string    `json:"kind"`
	Detail string    `json:"detail"`
}

func NewAuditor(s *store.Store) *Auditor {
	return &Auditor{store: s}
}

func (a *Auditor) Record(kind, detail string) error {
	entry := Entry{ID: uuid.NewString(), Time: time.Now().UTC(), Kind: kind, Detail: detail}
	raw, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	return store.AppendCommand(a.store, auditKey, string(raw))
}

func (a *Auditor) Entries() []Entry {
	var entries []Entry
	for _, item := range store.LoadCommands(a.store, auditKey) {
		var entry Entry
		if err := json.Unmarshal([]byte(item), &entry); err == nil {
			entries = append(entries, entry)
		}
	}
	return entries
}
