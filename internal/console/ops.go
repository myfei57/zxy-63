package console

import (
	"net/http"
)

type opsReport struct {
	StoreKeys     int  `json:"store_keys"`
	FlowPresent   bool `json:"flow_present"`
	FilterBeds    int  `json:"filter_beds"`
	BackwashQueue int  `json:"backwash_queue"`
	AuditEntries  int  `json:"audit_entries"`
}

func (s *Server) collectOps() opsReport {
	_, flowPresent := s.rt.flowRepo.LoadFlow()
	return opsReport{
		StoreKeys:     s.rt.store.Count(),
		FlowPresent:   flowPresent,
		FilterBeds:    s.rt.bank.Count(),
		BackwashQueue: s.rt.backwash.PendingCount(),
		AuditEntries:  s.rt.audit.Count(),
	}
}

func (s *Server) handleOps(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.collectOps())
}

func (s *Server) handleOpsResetStore(w http.ResponseWriter, r *http.Request) {
	if err := s.rt.store.Clear(); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"reset": true, "keys": s.rt.store.Count()})
}
