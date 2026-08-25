package console

import (
	"net/http"

	"waterplant/internal/store"
)

const eventKey = "console:events"

type historyAppendRequest struct {
	Kind  string `json:"kind"`
	Value string `json:"value"`
}

func (s *Server) handleHistoryAppend(w http.ResponseWriter, r *http.Request) {
	var body historyAppendRequest
	if err := readJSON(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if body.Kind == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "kind is required"})
		return
	}
	if err := store.AppendEvent(s.rt.store, eventKey, body.Kind, body.Value); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"appended": true})
}

func (s *Server) handleHistory(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"events":     store.ListEvents(s.rt.store, eventKey),
		"count":      store.EventCount(s.rt.store, eventKey),
		"has_events": s.rt.store.Has(eventKey),
	})
}

func (s *Server) handleHistoryClear(w http.ResponseWriter, r *http.Request) {
	if err := store.ClearCommands(s.rt.store, eventKey); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"cleared": true})
}
