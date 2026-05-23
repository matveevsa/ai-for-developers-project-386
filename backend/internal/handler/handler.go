package handler

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/matveevsa/ai-for-developers-project-386/backend/internal/model"
	"github.com/matveevsa/ai-for-developers-project-386/backend/internal/store"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, model.ErrorResponse{
		Code:    status,
		Message: message,
	})
}

func RegisterAll(mux *http.ServeMux, s *store.Store) {
	RegisterEventTypes(mux, s)
	RegisterPublic(mux, s)
	RegisterOwnerBookings(mux, s)
}
