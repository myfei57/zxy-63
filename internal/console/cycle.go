package console

import (
	"net/http"

	"waterplant/internal/filter"
	"waterplant/internal/flow"
	"waterplant/internal/intake"
	"waterplant/internal/quota"
)

type cycleRequest struct {
	Flow    float64   `json:"flow"`
	Samples []float64 `json:"samples"`
	Demand  float64   `json:"demand"`
	Level   float64   `json:"level"`
	BedID   string    `json:"bed_id"`
	Zone    int       `json:"zone"`
	Amount  float64   `json:"amount"`
}

type cycleResult struct {
	CoagDose  float64  `json:"coag_dose"`
	TurbDose  float64  `json:"turb_dose"`
	ChlorDose float64  `json:"chlor_dose"`
	MinLevel  float64  `json:"min_level"`
	Level     float64  `json:"level"`
	Quota     float64  `json:"quota"`
	Rotation  []string `json:"rotation"`
	OnDuty    string   `json:"on_duty"`
	Backwash  string   `json:"backwash"`
	Drains    int      `json:"drains"`
	Audit     int      `json:"audit_count"`
}

func average(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	total := 0.0
	for _, value := range values {
		total += value
	}
	return total / float64(len(values))
}

func (s *Server) handleCycle(w http.ResponseWriter, r *http.Request) {
	var body cycleRequest
	if err := readJSON(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if err := intake.ValidateFlow(body.Flow); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if err := flow.ValidateFactor(s.rt.calib.Current()); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if err := filter.ValidateZone(body.Zone); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if err := quota.ValidateAmount(body.Amount); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	sensor := intake.Sensor{FlowValue: body.Flow, Turbidity: average(body.Samples)}
	if err := s.rt.flowRepo.Record(sensor); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	coagDose, err := s.rt.coagDoser.UpdateFlowAndDose(body.Flow)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	turbDose, err := s.rt.sampler.Judge(body.Samples)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if err := s.rt.well.UpdateResidualDemand(body.Demand); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if err := s.rt.chlorDoser.ApplyResidual(); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	chlorDose := s.rt.chlorDoser.DoseForResidual()
	minLevel, err := s.rt.well.AdjustLevel(body.Level, s.rt.inlet, s.rt.outlet)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	rotation, err := s.rt.backwash.OrderRotation()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	duty, _ := s.rt.bank.OnDuty()
	if body.BedID != "" {
		if err := s.rt.backwash.Start(body.BedID); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
	}
	quotaValue, err := s.rt.acc.Add(body.Amount)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	result := cycleResult{
		CoagDose:  coagDose,
		TurbDose:  turbDose,
		ChlorDose: chlorDose,
		MinLevel:  minLevel,
		Level:     s.rt.well.Level(),
		Quota:     quotaValue,
		Rotation:  rotation,
		OnDuty:    duty,
		Backwash:  body.BedID,
		Drains:    s.rt.backwash.DrainCount(),
		Audit:     len(s.rt.audit.Entries()),
	}
	writeJSON(w, http.StatusOK, result)
}
