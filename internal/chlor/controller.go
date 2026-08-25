package chlor

import (
	"waterplant/internal/audit"
	"waterplant/internal/clearwell"
	"waterplant/internal/store"
)

type Doser struct {
	well   *clearwell.Well
	audit  *audit.Auditor
	target float64
}

func NewDoser(s *store.Store) *Doser {
	well := clearwell.NewWell(s)
	return &Doser{
		well:   well,
		audit:  audit.NewAuditor(s),
		target: well.ResidualTarget(),
	}
}
