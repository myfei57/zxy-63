package console

import (
	"net/http"

	"waterplant/internal/ns"
)

type simulateRequest struct {
	Ticks int `json:"ticks"`
}

type simulateTick struct {
	Index     int      `json:"index"`
	Dirtiest  string   `json:"dirtiest"`
	Cleanest  string   `json:"cleanest"`
	Eligible  []string `json:"eligible"`
	Dose      float64  `json:"dose"`
	Target    float64  `json:"target"`
	Remaining float64  `json:"remaining"`
	AuditKind string   `json:"audit_kind"`
}

type simulateReport struct {
	Ticks           []simulateTick `json:"ticks"`
	FinalAuditCount int            `json:"final_audit_count"`
	LastStage       ns.Stage       `json:"last_stage"`
}

func (s *Server) handleSimulate(w http.ResponseWriter, r *http.Request) {
	var body simulateRequest
	if err := readJSON(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	ticks := body.Ticks
	if ticks <= 0 {
		ticks = 5
	}
	if ticks > 100 {
		ticks = 100
	}
	report := simulateReport{LastStage: ns.TreatmentLine().Last()}
	for i := 0; i < ticks; i++ {
		report.Ticks = append(report.Ticks, simulateTick{
			Index:     i,
			Dirtiest:  s.rt.bank.Dirtiest(),
			Cleanest:  s.rt.bank.Cleanest(),
			Eligible:  s.rt.backwash.Eligible(),
			Dose:      s.rt.coagDoser.DosePlan(float64(i+1), float64(i)),
			Target:    s.rt.chlorDoser.TargetForDemand(float64(i)),
			Remaining: s.rt.acc.Remaining(100.0),
			AuditKind: s.lastAuditKind(),
		})
	}
	report.FinalAuditCount = len(s.rt.audit.Entries())
	writeJSON(w, http.StatusOK, report)
}

func (s *Server) lastAuditKind() string {
	last, ok := s.rt.audit.Last()
	if !ok {
		return ""
	}
	return last.Kind
}
