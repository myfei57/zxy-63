package chlor

import (
	"waterplant/internal/audit"
	"waterplant/internal/clearwell"
	"waterplant/internal/store"
)

type Doser struct {
	well  *clearwell.Well
	audit *audit.Auditor
}

func NewDoser(s *store.Store) *Doser {
	return &Doser{
		well:  clearwell.NewWell(s),
		audit: audit.NewAuditor(s),
	}
}
