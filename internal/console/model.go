package console

import (
	"encoding/json"
	"net/http"
)

func writeJSON(w http.ResponseWriter, status int, value interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func readJSON(r *http.Request, value interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(value)
}
