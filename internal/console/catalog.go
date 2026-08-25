package console

import "net/http"

type routeEntry struct {
	Method string `json:"method"`
	Path   string `json:"path"`
	Note   string `json:"note"`
}

func routeCatalog() []routeEntry {
	return []routeEntry{
		{Method: "GET", Path: "/health", Note: "service health and persisted keys"},
		{Method: "GET", Path: "/health/checks", Note: "detailed component health checks"},
		{Method: "GET", Path: "/snapshot", Note: "full component snapshot"},
		{Method: "GET", Path: "/describe", Note: "human readable component summary"},
		{Method: "GET", Path: "/catalog", Note: "list all HTTP routes"},
		{Method: "GET", Path: "/pipeline", Note: "treatment pipeline stages"},
		{Method: "GET", Path: "/ops", Note: "operator counters"},
		{Method: "GET", Path: "/telemetry", Note: "runtime telemetry counters"},
		{Method: "GET", Path: "/version", Note: "build version information"},
		{Method: "GET", Path: "/system", Note: "runtime system information"},
		{Method: "GET", Path: "/report", Note: "plain text control report"},
		{Method: "POST", Path: "/intake/flow", Note: "record a flow and turbidity reading"},
		{Method: "POST", Path: "/coag/dose", Note: "update persisted flow and dose coagulant"},
		{Method: "POST", Path: "/coag/turbidity", Note: "dose coagulant from mixed turbidity"},
		{Method: "POST", Path: "/flow/replace", Note: "replace flow meter calibration"},
		{Method: "POST", Path: "/chlor/target", Note: "update residual demand target"},
		{Method: "POST", Path: "/chlor/dose", Note: "apply a chlorine dose"},
		{Method: "POST", Path: "/clearwell/level", Note: "adjust clear well level"},
		{Method: "POST", Path: "/cycle", Note: "run a full treatment control cycle"},
		{Method: "POST", Path: "/simulate", Note: "run deterministic simulation ticks"},
		{Method: "POST", Path: "/filter/add", Note: "add a filter bed"},
		{Method: "POST", Path: "/filter/renumber", Note: "renumber a filter bed zone"},
		{Method: "POST", Path: "/filter/remove", Note: "remove a filter bed"},
		{Method: "POST", Path: "/filter/zone", Note: "read a filter bed zone"},
		{Method: "POST", Path: "/filter/load", Note: "set a filter bed load"},
		{Method: "POST", Path: "/filter/reset", Note: "reset all filter bed closed states"},
		{Method: "POST", Path: "/backwash/start", Note: "start a backwash sequence"},
		{Method: "GET", Path: "/backwash/order", Note: "compute filter duty rotation order"},
		{Method: "POST", Path: "/backwash/select", Note: "select a bed by zone"},
		{Method: "POST", Path: "/backwash/enqueue", Note: "enqueue a backwash command"},
		{Method: "POST", Path: "/backwash/recover", Note: "recover and discard stale commands"},
		{Method: "POST", Path: "/backwash/replay", Note: "replay queued backwash commands"},
		{Method: "POST", Path: "/quota/add", Note: "accumulate chlorine consumption"},
		{Method: "POST", Path: "/quota/check", Note: "check chemical quota usage"},
		{Method: "GET", Path: "/audit", Note: "list chemical dosing audit entries"},
		{Method: "GET", Path: "/audit/summary", Note: "audit entry count by kind"},
		{Method: "POST", Path: "/audit/filter", Note: "filter audit entries by kind"},
		{Method: "POST", Path: "/history", Note: "append a console event"},
		{Method: "GET", Path: "/history", Note: "list console events"},
		{Method: "POST", Path: "/history/clear", Note: "clear console events"},
		{Method: "POST", Path: "/ops/reset-store", Note: "factory reset persisted store"},
	}
}

func (s *Server) handleCatalog(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{"routes": routeCatalog()})
}
