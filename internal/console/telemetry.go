package console

import "net/http"

type telemetryReport struct {
	StoreKeys      int     `json:"store_keys"`
	FilterBeds     int     `json:"filter_beds"`
	FilterActive   int     `json:"filter_active"`
	FilterLoad     float64 `json:"filter_load"`
	BackwashQueue  int     `json:"backwash_queue"`
	BackwashDrains int     `json:"backwash_drains"`
	LastDrain      string  `json:"last_drain"`
	AuditEntries   int     `json:"audit_entries"`
	QuotaValue     float64 `json:"quota_value"`
	ClearwellLevel float64 `json:"clearwell_level"`
	CoagRatio      float64 `json:"coag_ratio"`
	ChlorTarget    float64 `json:"chlor_target"`
	Turbidity      float64 `json:"turbidity_verdict"`
}

func (s *Server) collectTelemetry() telemetryReport {
	lastDrain, _ := s.rt.backwash.LastDrain()
	return telemetryReport{
		StoreKeys:      s.rt.store.Count(),
		FilterBeds:     s.rt.bank.Count(),
		FilterActive:   s.rt.bank.ActiveCount(),
		FilterLoad:     s.rt.bank.TotalLoad(),
		BackwashQueue:  s.rt.backwash.PendingCount(),
		BackwashDrains: s.rt.backwash.DrainCount(),
		LastDrain:      lastDrain,
		AuditEntries:   s.rt.audit.Count(),
		QuotaValue:     s.rt.acc.Value(),
		ClearwellLevel: s.rt.well.Level(),
		CoagRatio:      s.rt.coagDoser.CurrentRatio(),
		ChlorTarget:    s.rt.chlorDoser.CurrentTarget(),
		Turbidity:      s.rt.sampler.State().LastVerdict,
	}
}

func (s *Server) handleTelemetry(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.collectTelemetry())
}
