package console

import (
	"fmt"
	"net/http"
	"strings"

	"waterplant/internal/ns"
)

func (s *Server) textReport() string {
	describe := s.collectDescribe()
	telemetry := s.collectTelemetry()
	var builder strings.Builder
	builder.WriteString("waterplant control report\n")
	builder.WriteString(fmt.Sprintf("pipeline: %s\n", describe.Pipeline))
	builder.WriteString(fmt.Sprintf("pipeline stages=%d contains_filter=%t\n", ns.TreatmentLine().Count(), ns.TreatmentLine().Contains(ns.StageFilter)))
	builder.WriteString(fmt.Sprintf("store: %s\n", describe.Store))
	builder.WriteString(fmt.Sprintf("intake: %s\n", describe.Intake))
	builder.WriteString(fmt.Sprintf("coagulant: %s\n", describe.Coag))
	builder.WriteString(fmt.Sprintf("chlorine: %s\n", describe.Chlor))
	builder.WriteString(fmt.Sprintf("filter: %s\n", describe.Filter))
	builder.WriteString(fmt.Sprintf("backwash: %s\n", describe.Backwash))
	builder.WriteString(fmt.Sprintf("turbidity: %s\n", describe.Turbidity))
	builder.WriteString(fmt.Sprintf("flow: %s\n", describe.Flow))
	builder.WriteString(fmt.Sprintf("clearwell: %s\n", describe.Clearwell))
	builder.WriteString(fmt.Sprintf("quota: %s\n", describe.Quota))
	builder.WriteString(fmt.Sprintf("audit: %s\n", describe.Audit))
	builder.WriteString(fmt.Sprintf("audit kinds=%v\n", s.rt.audit.Kinds()))
	builder.WriteString(fmt.Sprintf("filter beds=%d active=%d total_load=%.4f\n", telemetry.FilterBeds, telemetry.FilterActive, telemetry.FilterLoad))
	builder.WriteString(fmt.Sprintf("filter bed_ids=%v\n", s.rt.bank.BedIDs()))
	builder.WriteString(fmt.Sprintf("backwash queue=%d drains=%d\n", telemetry.BackwashQueue, telemetry.BackwashDrains))
	builder.WriteString(fmt.Sprintf("backwash commands=%v\n", s.rt.backwash.CommandList()))
	builder.WriteString(fmt.Sprintf("quota value=%.4f clearwell level=%.4f\n", telemetry.QuotaValue, telemetry.ClearwellLevel))
	builder.WriteString(fmt.Sprintf("coag ratio=%.4f chlorine target=%.4f\n", telemetry.CoagRatio, telemetry.ChlorTarget))
	return builder.String()
}

func (s *Server) handleReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(s.textReport()))
}
