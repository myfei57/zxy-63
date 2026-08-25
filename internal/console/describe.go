package console

import (
	"net/http"

	"waterplant/internal/ns"
)

type describeReport struct {
	Pipeline  string `json:"pipeline"`
	Store     string `json:"store"`
	Intake    string `json:"intake"`
	Coag      string `json:"coag"`
	Chlor     string `json:"chlor"`
	Filter    string `json:"filter"`
	Backwash  string `json:"backwash"`
	Turbidity string `json:"turbidity"`
	Flow      string `json:"flow"`
	Clearwell string `json:"clearwell"`
	Quota     string `json:"quota"`
	Audit     string `json:"audit"`
}

func (s *Server) collectDescribe() describeReport {
	return describeReport{
		Pipeline:  ns.TreatmentLine().Describe(),
		Store:     s.rt.store.Describe(),
		Intake:    s.rt.flowRepo.Describe(),
		Coag:      s.rt.coagDoser.Describe(),
		Chlor:     s.rt.chlorDoser.Describe(),
		Filter:    s.rt.bank.Describe(),
		Backwash:  s.rt.backwash.Describe(),
		Turbidity: s.rt.sampler.Describe(),
		Flow:      s.rt.calib.Describe(),
		Clearwell: s.rt.well.Describe(),
		Quota:     s.rt.acc.Describe(),
		Audit:     s.rt.audit.Describe(),
	}
}

func (s *Server) handleDescribe(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.collectDescribe())
}
