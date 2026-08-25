package console

import "net/http"

const buildName = "waterplant-control"
const buildVersion = "1.0.0"

type versionReport struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func Version() string {
	return buildName + " " + buildVersion
}

func (s *Server) handleVersion(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, versionReport{Name: buildName, Version: buildVersion})
}
