package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/matveevsa/ai-for-developers-project-386/backend/internal/model"
	"github.com/matveevsa/ai-for-developers-project-386/backend/internal/store"
)

func RegisterPublic(mux *http.ServeMux, s *store.Store) {
	mux.HandleFunc("GET /api/public/event-types", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, s.ListEventTypes())
	})

	mux.HandleFunc("GET /api/public/event-types/{eventTypeId}/availability", func(w http.ResponseWriter, r *http.Request) {
		eventTypeID := r.PathValue("eventTypeId")

		et, ok := s.GetEventType(eventTypeID)
		if !ok {
			writeError(w, 404, "event type not found")
			return
		}

		now := time.Now().UTC()
		today := now.Truncate(24 * time.Hour)
		result := make([]map[string]any, 0, 14)

		for i := 0; i < 14; i++ {
			date := today.AddDate(0, 0, i)
			allSlots := enumerateSlots(et, date)

			available := 0
			for _, slot := range allSlots {
				startTime, _, err := parseSlotID(slot.ID)
				if err != nil {
					continue
				}
				if startTime.Before(now) {
					continue
				}
				if !s.IsSlotBooked(slot.ID) {
					available++
				}
			}

			result = append(result, map[string]any{
				"date":      date.Format("2006-01-02"),
				"available": available,
				"total":     len(allSlots),
			})
		}

		writeJSON(w, 200, result)
	})

	mux.HandleFunc("GET /api/public/event-types/{eventTypeId}/slots", func(w http.ResponseWriter, r *http.Request) {
		eventTypeID := r.PathValue("eventTypeId")
		dateStr := r.URL.Query().Get("date")

		et, ok := s.GetEventType(eventTypeID)
		if !ok {
			writeError(w, 404, "event type not found")
			return
		}

		date, err := time.Parse(time.RFC3339, dateStr)
		if err != nil {
			date, err = time.Parse("2006-01-02", dateStr)
			if err != nil {
				writeError(w, 400, "invalid date format, use YYYY-MM-DD or RFC3339")
				return
			}
		}
		date = date.UTC().Truncate(24 * time.Hour)

		slots := generateSlots(et, date, s)
		writeJSON(w, 200, slots)
	})

	mux.HandleFunc("POST /api/public/bookings", func(w http.ResponseWriter, r *http.Request) {
		var req model.BookingCreate
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, 400, "invalid JSON body")
			return
		}

		if req.GuestName == "" || len(req.GuestName) > 100 {
			writeError(w, 400, "guestName must be non-empty and at most 100 characters")
			return
		}

		if !emailRegex.MatchString(req.GuestEmail) {
			writeError(w, 400, "invalid email format")
			return
		}

		et, ok := s.GetEventType(req.EventTypeID)
		if !ok {
			writeError(w, 404, "event type not found")
			return
		}

		slotStartTime, slotDate, err := parseSlotID(req.SlotID)
		if err != nil {
			writeError(w, 400, "invalid slot ID")
			return
		}

		if !isSlotValid(et, slotStartTime, slotDate, s) {
			writeError(w, 400, "invalid or unavailable slot")
			return
		}

		if s.IsSlotBooked(req.SlotID) {
			writeError(w, 409, "slot is already booked")
			return
		}

		booking := model.Booking{
			ID:          s.NextBookingID(),
			SlotID:      req.SlotID,
			EventTypeID: req.EventTypeID,
			GuestName:   req.GuestName,
			GuestEmail:  req.GuestEmail,
			GuestNotes:  req.GuestNotes,
		}
		s.CreateBooking(booking)
		writeJSON(w, 200, booking)
	})
}

func makeSlotID(eventTypeID string, startTime time.Time) string {
	return fmt.Sprintf("%s-%d", eventTypeID, startTime.Unix())
}

func parseSlotID(slotID string) (startTime time.Time, date time.Time, err error) {
	lastDash := strings.LastIndex(slotID, "-")
	if lastDash < 0 {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid slot ID")
	}
	unixStr := slotID[lastDash+1:]
	unix, err := strconv.ParseInt(unixStr, 10, 64)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid slot ID")
	}
	startTime = time.Unix(unix, 0).UTC()
	date = startTime.Truncate(24 * time.Hour)
	return
}

func isWeekday(t time.Time) bool {
	weekday := t.Weekday()
	return weekday != time.Saturday && weekday != time.Sunday
}

func generateSlots(et model.EventType, date time.Time, s *store.Store) []model.Slot {
	slots := enumerateSlots(et, date)

	now := time.Now().UTC()
	available := make([]model.Slot, 0, len(slots))
	for _, slot := range slots {
		startTime, _, err := parseSlotID(slot.ID)
		if err != nil {
			continue
		}
		if startTime.Before(now) {
			continue
		}
		if s.IsSlotBooked(slot.ID) {
			continue
		}
		available = append(available, slot)
	}
	return available
}

func enumerateSlots(et model.EventType, date time.Time) []model.Slot {
	if !isWeekday(date) {
		return []model.Slot{}
	}

	now := time.Now().UTC()
	today := now.Truncate(24 * time.Hour)
	maxDate := today.AddDate(0, 0, 14)

	if date.Before(today) || date.After(maxDate) {
		return []model.Slot{}
	}

	duration := time.Duration(et.Duration) * time.Minute
	workStart := time.Date(date.Year(), date.Month(), date.Day(), 9, 0, 0, 0, time.UTC)
	workEnd := time.Date(date.Year(), date.Month(), date.Day(), 18, 0, 0, 0, time.UTC)

	slots := make([]model.Slot, 0)
	for t := workStart; t.Add(duration).Compare(workEnd) <= 0; t = t.Add(duration) {
		slots = append(slots, model.Slot{
			ID:          makeSlotID(et.ID, t),
			EventTypeID: et.ID,
			StartTime:   t.Format(time.RFC3339),
			EndTime:     t.Add(duration).Format(time.RFC3339),
		})
	}
	return slots
}

func isSlotValid(et model.EventType, startTime time.Time, date time.Time, s *store.Store) bool {
	slotID := makeSlotID(et.ID, startTime)
	slots := enumerateSlots(et, date)
	for _, slot := range slots {
		if slot.ID == slotID {
			return true
		}
	}
	return false
}
