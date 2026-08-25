package console

import (
	"net/http"

	"waterplant/internal/audit"
	"waterplant/internal/backwash"
	"waterplant/internal/chlor"
	"waterplant/internal/clearwell"
	"waterplant/internal/coag"
	"waterplant/internal/filter"
	"waterplant/internal/flow"
	"waterplant/internal/intake"
	"waterplant/internal/ns"
	"waterplant/internal/quota"
	"waterplant/internal/store"
	"waterplant/internal/turb"
)

type snapshot struct {
	Pipeline  []ns.Step              `json:"pipeline"`
	Store     store.StoreState       `json:"store"`
	Intake    intake.FlowState       `json:"intake"`
	Coag      coag.DoseState         `json:"coag"`
	Chlor     chlor.ChlorState       `json:"chlor"`
	Filter    filter.BankState       `json:"filter"`
	Backwash  backwash.BackwashState `json:"backwash"`
	Turbidity turb.TurbState         `json:"turbidity"`
	Flow      flow.CalibState        `json:"flow"`
	Clearwell clearwell.WellState    `json:"clearwell"`
	Quota     quota.AccumState       `json:"quota"`
	Audit     audit.AuditState       `json:"audit"`
}

func (s *Server) collectSnapshot() snapshot {
	return snapshot{
		Pipeline:  ns.TreatmentLine().Steps(),
		Store:     s.rt.store.State(),
		Intake:    s.rt.flowRepo.State(),
		Coag:      s.rt.coagDoser.State(),
		Chlor:     s.rt.chlorDoser.State(),
		Filter:    s.rt.bank.State(),
		Backwash:  s.rt.backwash.State(),
		Turbidity: s.rt.sampler.State(),
		Flow:      s.rt.calib.State(),
		Clearwell: s.rt.well.State(),
		Quota:     s.rt.acc.State(),
		Audit:     s.rt.audit.State(),
	}
}

func (s *Server) handleSnapshot(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.collectSnapshot())
}
