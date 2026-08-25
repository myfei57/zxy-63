package console

import (
	"net/http"

	"waterplant/internal/flow"
	"waterplant/internal/chlor"
	"waterplant/internal/clearwell"
	"waterplant/internal/filter"
	"waterplant/internal/intake"
	"waterplant/internal/ns"
	"waterplant/internal/quota"
)

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status": "ok",
		"keys":   s.rt.store.Keys(),
		"size":   s.rt.store.Count(),
		"data":   s.rt.store.Export(),
	})
}

func (s *Server) handlePipeline(w http.ResponseWriter, r *http.Request) {
	line := ns.TreatmentLine()
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"name":      line.Name,
		"steps":     line.Steps(),
		"ordered":   line.Before(ns.StageIntake, ns.StageCoag),
		"lastStage": line.Stages[len(line.Stages)-1],
	})
}

func (s *Server) handleIntakeFlow(w http.ResponseWriter, r *http.Request) {
	var sensor intake.Sensor
	if err := readJSON(r, &sensor); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if err := intake.ValidateFlow(sensor.FlowValue); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if err := s.rt.flowRepo.Record(sensor); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"flow": sensor.FlowValue})
}

func (s *Server) handleIntakeFlowGet(w http.ResponseWriter, r *http.Request) {
	value, ok := s.rt.flowRepo.LoadFlow()
	writeJSON(w, http.StatusOK, map[string]interface{}{"flow": value, "ok": ok})
}

func (s *Server) handleCoagDose(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Flow float64 `json:"flow"`
	}
	if err := readJSON(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if err := intake.ValidateFlow(body.Flow); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	dose, err := s.rt.coagDoser.UpdateFlowAndDose(body.Flow)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"dose": dose})
}

func (s *Server) handleCoagTurbidity(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Samples []float64 `json:"samples"`
	}
	if err := readJSON(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	dose, err := s.rt.sampler.Judge(body.Samples)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"dose": dose})
}

func (s *Server) handleCoagRatio(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{"ratio": s.rt.coagDoser.CurrentRatio()})
}

func (s *Server) handleFlowReplace(w http.ResponseWriter, r *http.Request) {
	var meter flow.Meter
	if err := readJSON(r, &meter); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if err := flow.ValidateFactor(meter.Factor); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	calibrated := flow.CalibrateMeter(meter.Serial, meter.Factor)
	if err := s.rt.calib.Replace(calibrated.Factor); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"meter": calibrated})
}

func (s *Server) handleChlorTarget(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Demand float64 `json:"demand"`
	}
	if err := readJSON(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if err := chlor.ValidateDemand(body.Demand); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if err := s.rt.well.UpdateResidualDemand(body.Demand); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"target": s.rt.well.ResidualTarget()})
}

func (s *Server) handleChlorDose(w http.ResponseWriter, r *http.Request) {
	if err := s.rt.chlorDoser.ApplyResidual(); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"dose": s.rt.chlorDoser.DoseForResidual()})
}

func (s *Server) handleClearwellLevel(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Target float64 `json:"target"`
	}
	if err := readJSON(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if err := clearwell.ValidateLevel(body.Target); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	minSeen, err := s.rt.well.AdjustLevel(body.Target, s.rt.inlet, s.rt.outlet)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"min_level": minSeen, "level": s.rt.well.Level()})
}

func (s *Server) handleFilterAdd(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID   string  `json:"id"`
		Zone int     `json:"zone"`
		Load float64 `json:"load"`
	}
	if err := readJSON(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if err := filter.ValidateZone(body.Zone); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	s.rt.bank.AddBed(body.ID, body.Zone, body.Load)
	bed, _ := s.rt.bank.Bed(body.ID)
	writeJSON(w, http.StatusOK, map[string]interface{}{"id": bed.ID, "zone": bed.Zone, "load": bed.Load})
}

func (s *Server) handleFilterClose(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID string `json:"id"`
	}
	if err := readJSON(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if err := s.rt.bank.Close(body.ID); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"closed": body.ID})
}

func (s *Server) handleFilterOpen(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID string `json:"id"`
	}
	if err := readJSON(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if err := s.rt.bank.Open(body.ID); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"opened": body.ID})
}

func (s *Server) handleFilterRenumber(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID   string `json:"id"`
		Zone int    `json:"zone"`
	}
	if err := readJSON(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if err := filter.ValidateZone(body.Zone); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if err := s.rt.bank.Renumber(body.ID, body.Zone); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"renumbered": body.ID, "zone": body.Zone})
}

func (s *Server) handleFilterLoad(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID   string  `json:"id"`
		Load float64 `json:"load"`
	}
	if err := readJSON(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if body.Load < 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "load must be non-negative"})
		return
	}
	if err := s.rt.bank.SetLoad(body.ID, body.Load); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"id": body.ID, "load": body.Load})
}

func (s *Server) handleFilterReset(w http.ResponseWriter, r *http.Request) {
	s.rt.bank.ResetClosed()
	writeJSON(w, http.StatusOK, map[string]interface{}{"reset": true})
}

func (s *Server) handleFilterRemove(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID string `json:"id"`
	}
	if err := readJSON(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if err := s.rt.bank.RemoveBed(body.ID); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"removed": body.ID})
}

func (s *Server) handleFilterZone(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID string `json:"id"`
	}
	if err := readJSON(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	zone, ok := s.rt.bank.BedZone(body.ID)
	writeJSON(w, http.StatusOK, map[string]interface{}{"id": body.ID, "zone": zone, "ok": ok})
}

func (s *Server) handleBackwashStart(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID string `json:"id"`
	}
	if err := readJSON(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if err := s.rt.backwash.Start(body.ID); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"started": body.ID})
}

func (s *Server) handleBackwashOrder(w http.ResponseWriter, r *http.Request) {
	order, err := s.rt.backwash.OrderRotation()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	duty, _ := s.rt.bank.OnDuty()
	writeJSON(w, http.StatusOK, map[string]interface{}{"order": order, "on_duty": duty})
}

func (s *Server) handleBackwashSelect(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Zone int `json:"zone"`
	}
	if err := readJSON(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	bedID, err := s.rt.backwash.Select(body.Zone)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"bed": bedID, "zone": body.Zone})
}

func (s *Server) handleBackwashEnqueue(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID string `json:"id"`
	}
	if err := readJSON(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if err := s.rt.backwash.Enqueue(body.ID); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"enqueued": body.ID})
}

func (s *Server) handleBackwashRecover(w http.ResponseWriter, r *http.Request) {
	if err := s.rt.backwash.Recover(); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"recovered": true})
}

func (s *Server) handleBackwashReplay(w http.ResponseWriter, r *http.Request) {
	commands, err := s.rt.backwash.Replay()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"replayed": commands})
}

func (s *Server) handleQuotaAdd(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Amount float64 `json:"amount"`
	}
	if err := readJSON(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if err := quota.ValidateAmount(body.Amount); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	value, err := s.rt.acc.Add(body.Amount)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"value": value})
}

func (s *Server) handleQuotaGet(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{"value": s.rt.acc.Value()})
}

func (s *Server) handleQuotaCheck(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Chemical string  `json:"chemical"`
		Used     float64 `json:"used"`
		Limit    float64 `json:"limit"`
	}
	if err := readJSON(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	quotaValue := quota.Quota{Chemical: body.Chemical, Limit: body.Limit}
	remaining, ok := quota.CheckQuota(quotaValue.Chemical, body.Used, quotaValue.Limit)
	writeJSON(w, http.StatusOK, map[string]interface{}{"remaining": remaining, "ok": ok})
}

func (s *Server) handleAudit(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{"entries": s.rt.audit.Entries()})
}

func (s *Server) handleAuditSummary(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"count":   len(s.rt.audit.Entries()),
		"by_kind": s.rt.audit.CountByKind(),
	})
}

func (s *Server) handleAuditFilter(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Kind string `json:"kind"`
	}
	if err := readJSON(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"entries": s.rt.audit.Filter(body.Kind)})
}
