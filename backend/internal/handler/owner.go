package handler

import (
	"net/http"
	"time"

	"github.com/matveevsa/ai-for-developers-project-386/backend/internal/model"
	"github.com/matveevsa/ai-for-developers-project-386/backend/internal/store"
)

func RegisterOwnerBookings(mux *http.ServeMux, s *store.Store) {
	mux.HandleFunc("DELETE /api/bookings/{bookingId}", func(w http.ResponseWriter, r *http.Request) {
		bookingID := r.PathValue("bookingId")
		if err := s.DeleteBooking(bookingID); err != nil {
			writeError(w, 404, err.Error())
			return
		}
		writeJSON(w, 200, map[string]string{"status": "deleted"})
	})

	mux.HandleFunc("GET /api/owners/{ownerId}/bookings", func(w http.ResponseWriter, r *http.Request) {
		ownerID := r.PathValue("ownerId")
		bookings := s.ListBookingsByOwner(ownerID)

		enriched := make([]model.BookingEnriched, 0, len(bookings))
		for _, b := range bookings {
			startTime, _, err := parseSlotID(b.SlotID)
			if err != nil {
				continue
			}

			et, _ := s.GetEventType(b.EventTypeID)
			dur := time.Duration(et.Duration) * time.Minute

			enriched = append(enriched, model.BookingEnriched{
				ID:            b.ID,
				SlotID:        b.SlotID,
				EventTypeID:   b.EventTypeID,
				GuestName:     b.GuestName,
				GuestEmail:    b.GuestEmail,
				GuestNotes:    b.GuestNotes,
				StartTime:     startTime.Format(time.RFC3339),
				EndTime:       startTime.Add(dur).Format(time.RFC3339),
				EventTypeName: et.Name,
			})
		}

		writeJSON(w, 200, enriched)
	})
}
