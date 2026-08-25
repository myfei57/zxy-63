package console

import "net/http"

type healthCheck struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Detail string `json:"detail"`
}

func (s *Server) runChecks() []healthCheck {
	checks := make([]healthCheck, 0, 5)
	checks = append(checks, healthCheck{Name: "store", Status: "ok", Detail: "keys=" + itoa(s.rt.store.Count())})
	if s.rt.calib.Current() <= 0 {
		checks = append(checks, healthCheck{Name: "calibration", Status: "fail", Detail: "factor must be positive"})
	} else {
		checks = append(checks, healthCheck{Name: "calibration", Status: "ok", Detail: "factor set"})
	}
	if s.rt.well.ResidualTarget() <= 0 {
		checks = append(checks, healthCheck{Name: "residual", Status: "fail", Detail: "target must be positive"})
	} else {
		checks = append(checks, healthCheck{Name: "residual", Status: "ok", Detail: "target set"})
	}
	if s.rt.acc.Value() < 0 {
		checks = append(checks, healthCheck{Name: "quota", Status: "fail", Detail: "accumulator is negative"})
	} else {
		checks = append(checks, healthCheck{Name: "quota", Status: "ok", Detail: "accumulator set"})
	}
	if _, ok := s.rt.flowRepo.LoadFlow(); !ok {
		checks = append(checks, healthCheck{Name: "intake", Status: "warn", Detail: "no flow reading yet"})
	} else {
		checks = append(checks, healthCheck{Name: "intake", Status: "ok", Detail: "flow reading present"})
	}
	return checks
}

func (s *Server) handleChecks(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{"checks": s.runChecks()})
}

func itoa(value int) string {
	if value == 0 {
		return "0"
	}
	negative := value < 0
	if negative {
		value = -value
	}
	var digits []byte
	for value > 0 {
		digits = append([]byte{byte('0' + value%10)}, digits...)
		value /= 10
	}
	if negative {
		return "-" + string(digits)
	}
	return string(digits)
}
