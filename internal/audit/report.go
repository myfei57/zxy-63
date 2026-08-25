package audit

import "fmt"

type AuditState struct {
	Entries []Entry `json:"entries"`
	Count   int     `json:"count"`
}

func (a *Auditor) State() AuditState {
	entries := a.Entries()
	return AuditState{Entries: entries, Count: len(entries)}
}

func (a *Auditor) Describe() string {
	return fmt.Sprintf("audit entries=%d", len(a.Entries()))
}

func (a *Auditor) Last() (Entry, bool) {
	entries := a.Entries()
	if len(entries) == 0 {
		return Entry{}, false
	}
	return entries[len(entries)-1], true
}

func (a *Auditor) CountByKind() map[string]int {
	counts := map[string]int{}
	for _, entry := range a.Entries() {
		counts[entry.Kind]++
	}
	return counts
}

func (a *Auditor) Filter(kind string) []Entry {
	var out []Entry
	for _, entry := range a.Entries() {
		if entry.Kind == kind {
			out = append(out, entry)
		}
	}
	return out
}

func (a *Auditor) Count() int {
	return len(a.Entries())
}

func (a *Auditor) Kinds() []string {
	seen := map[string]bool{}
	var kinds []string
	for _, entry := range a.Entries() {
		if !seen[entry.Kind] {
			seen[entry.Kind] = true
			kinds = append(kinds, entry.Kind)
		}
	}
	return kinds
}
