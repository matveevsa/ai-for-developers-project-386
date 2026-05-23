package handler

import (
	"encoding/json"
	"net/http"

	"github.com/matveevsa/ai-for-developers-project-386/backend/internal/model"
	"github.com/matveevsa/ai-for-developers-project-386/backend/internal/store"
)

func RegisterEventTypes(mux *http.ServeMux, s *store.Store) {
	mux.HandleFunc("POST /api/event-types", func(w http.ResponseWriter, r *http.Request) {
		var req model.EventTypeCreate
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, 400, "invalid JSON body")
			return
		}

		if req.Name == "" || len(req.Name) > 100 {
			writeError(w, 400, "name must be non-empty and at most 100 characters")
			return
		}

		if req.Duration != 15 && req.Duration != 30 {
			writeError(w, 400, "duration must be 15 or 30 minutes")
			return
		}

		et := model.EventType{
			ID:          s.NextEventID(),
			OwnerID:     "owner-1",
			Name:        req.Name,
			Description: req.Description,
			Duration:    req.Duration,
		}
		s.CreateEventType(et)
		writeJSON(w, 200, et)
	})

	mux.HandleFunc("GET /api/event-types", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, s.ListEventTypes())
	})

	mux.HandleFunc("PUT /api/event-types/{eventTypeId}", func(w http.ResponseWriter, r *http.Request) {
		eventTypeID := r.PathValue("eventTypeId")

		var req model.EventTypeCreate
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, 400, "invalid JSON body")
			return
		}

		if req.Name == "" || len(req.Name) > 100 {
			writeError(w, 400, "name must be non-empty and at most 100 characters")
			return
		}

		if req.Duration != 15 && req.Duration != 30 {
			writeError(w, 400, "duration must be 15 or 30 minutes")
			return
		}

		updated, err := s.UpdateEventType(eventTypeID, req)
		if err != nil {
			writeError(w, 404, err.Error())
			return
		}

		writeJSON(w, 200, updated)
	})

	mux.HandleFunc("DELETE /api/event-types/{eventTypeId}", func(w http.ResponseWriter, r *http.Request) {
		eventTypeID := r.PathValue("eventTypeId")

		if err := s.DeleteEventType(eventTypeID); err != nil {
			writeError(w, 404, err.Error())
			return
		}

		writeJSON(w, 200, map[string]string{"status": "deleted"})
	})
}
