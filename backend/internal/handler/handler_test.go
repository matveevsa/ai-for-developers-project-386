package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/matveevsa/ai-for-developers-project-386/backend/internal/model"
	"github.com/matveevsa/ai-for-developers-project-386/backend/internal/store"
)

func setupTest(t *testing.T) (*store.Store, *httptest.Server) {
	t.Helper()
	s := store.New()
	mux := http.NewServeMux()
	RegisterAll(mux, s)
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)
	return s, ts
}

func newServer(s *store.Store) *httptest.Server {
	mux := http.NewServeMux()
	RegisterAll(mux, s)
	return httptest.NewServer(mux)
}

func createEventType(t *testing.T, s *store.Store, name string, duration int) model.EventType {
	t.Helper()
	et := model.EventType{
		ID:       s.NextEventID(),
		OwnerID:  "owner-1",
		Name:     name,
		Duration: duration,
	}
	s.CreateEventType(et)
	return et
}

func nextWeekdayDate(daysFromNow int) time.Time {
	t := time.Now().UTC().Truncate(24 * time.Hour).AddDate(0, 0, daysFromNow)
	for t.Weekday() == time.Saturday || t.Weekday() == time.Sunday {
		t = t.AddDate(0, 0, 1)
	}
	return t
}

func TestPublicEventTypesList(t *testing.T) {
	s := store.New()
	createEventType(t, s, "Test Meeting", 15)
	createEventType(t, s, "Another Meeting", 30)
	ts := newServer(s)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/public/event-types")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var types []model.EventType
	if err := json.NewDecoder(resp.Body).Decode(&types); err != nil {
		t.Fatal(err)
	}

	if len(types) != 2 {
		t.Fatalf("expected 2 event types, got %d", len(types))
	}
}

func TestPublicAvailability_Success(t *testing.T) {
	s := store.New()
	et := createEventType(t, s, "Test", 15)
	ts := newServer(s)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/public/event-types/" + et.ID + "/availability")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var availability []map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&availability); err != nil {
		t.Fatal(err)
	}

	if len(availability) != 14 {
		t.Fatalf("expected 14 days, got %d", len(availability))
	}
}

func TestPublicAvailability_NotFound(t *testing.T) {
	_, ts := setupTest(t)

	resp, err := http.Get(ts.URL + "/api/public/event-types/non-existent/availability")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}

func TestPublicSlots_Success(t *testing.T) {
	s := store.New()
	et := createEventType(t, s, "Test", 15)
	ts := newServer(s)
	defer ts.Close()

	date := nextWeekdayDate(1)
	dateStr := date.Format("2006-01-02")

	resp, err := http.Get(ts.URL + "/api/public/event-types/" + et.ID + "/slots?date=" + dateStr)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var slots []model.Slot
	if err := json.NewDecoder(resp.Body).Decode(&slots); err != nil {
		t.Fatal(err)
	}

	now := time.Now().UTC()
	filtered := 0
	for _, sl := range slots {
		st, _, _ := parseSlotID(sl.ID)
		if !st.Before(now) {
			filtered++
		}
	}
	if filtered == 0 {
		t.Fatal("expected at least one available slot")
	}
}

func TestPublicSlots_NotFound(t *testing.T) {
	_, ts := setupTest(t)

	resp, err := http.Get(ts.URL + "/api/public/event-types/non-existent/slots?date=2026-05-25")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}

func TestPublicSlots_InvalidDate(t *testing.T) {
	s := store.New()
	et := createEventType(t, s, "Test", 15)
	ts := newServer(s)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/public/event-types/" + et.ID + "/slots?date=invalid")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestPublicSlots_MissingDate(t *testing.T) {
	s := store.New()
	et := createEventType(t, s, "Test", 15)
	ts := newServer(s)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/public/event-types/" + et.ID + "/slots")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestCreateBooking_Success(t *testing.T) {
	s := store.New()
	et := createEventType(t, s, "Test", 15)

	slotDate := nextWeekdayDate(1)
	slotTime := time.Date(slotDate.Year(), slotDate.Month(), slotDate.Day(), 10, 0, 0, 0, time.UTC)
	slotID := makeSlotID(et.ID, slotTime)

	ts := newServer(s)
	defer ts.Close()

	body := fmt.Sprintf(`{"slotId":"%s","eventTypeId":"%s","guestName":"Иван","guestEmail":"ivan@example.com"}`, slotID, et.ID)
	resp, err := http.Post(ts.URL+"/api/public/bookings", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var booking model.Booking
	if err := json.NewDecoder(resp.Body).Decode(&booking); err != nil {
		t.Fatal(err)
	}

	if booking.GuestName != "Иван" {
		t.Fatalf("expected name Иван, got %s", booking.GuestName)
	}
	if booking.GuestEmail != "ivan@example.com" {
		t.Fatalf("expected email ivan@example.com, got %s", booking.GuestEmail)
	}
	if booking.SlotID != slotID {
		t.Fatalf("expected slotID %s, got %s", slotID, booking.SlotID)
	}
}

func TestCreateBooking_DoubleBooking(t *testing.T) {
	s := store.New()
	et := createEventType(t, s, "Test", 15)

	slotDate := nextWeekdayDate(1)
	slotTime := time.Date(slotDate.Year(), slotDate.Month(), slotDate.Day(), 10, 0, 0, 0, time.UTC)
	slotID := makeSlotID(et.ID, slotTime)

	ts := newServer(s)
	defer ts.Close()

	body := fmt.Sprintf(`{"slotId":"%s","eventTypeId":"%s","guestName":"Иван","guestEmail":"ivan@example.com"}`, slotID, et.ID)

	resp, err := http.Post(ts.URL+"/api/public/bookings", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	resp, err = http.Post(ts.URL+"/api/public/bookings", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusConflict {
		t.Fatalf("expected 409, got %d", resp.StatusCode)
	}
}

func TestCreateBooking_InvalidEmail(t *testing.T) {
	s := store.New()
	et := createEventType(t, s, "Test", 15)

	slotDate := nextWeekdayDate(1)
	slotTime := time.Date(slotDate.Year(), slotDate.Month(), slotDate.Day(), 10, 0, 0, 0, time.UTC)
	slotID := makeSlotID(et.ID, slotTime)

	ts := newServer(s)
	defer ts.Close()

	body := fmt.Sprintf(`{"slotId":"%s","eventTypeId":"%s","guestName":"Иван","guestEmail":"not-an-email"}`, slotID, et.ID)
	resp, err := http.Post(ts.URL+"/api/public/bookings", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestCreateBooking_EmptyName(t *testing.T) {
	s := store.New()
	et := createEventType(t, s, "Test", 15)
	slotDate := nextWeekdayDate(1)
	slotTime := time.Date(slotDate.Year(), slotDate.Month(), slotDate.Day(), 10, 0, 0, 0, time.UTC)
	slotID := makeSlotID(et.ID, slotTime)

	ts := newServer(s)
	defer ts.Close()

	body := fmt.Sprintf(`{"slotId":"%s","eventTypeId":"%s","guestName":"","guestEmail":"ivan@example.com"}`, slotID, et.ID)
	resp, err := http.Post(ts.URL+"/api/public/bookings", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestCreateBooking_EventTypeNotFound(t *testing.T) {
	_, ts := setupTest(t)

	body := `{"slotId":"slot-123","eventTypeId":"non-existent","guestName":"Иван","guestEmail":"ivan@example.com"}`
	resp, err := http.Post(ts.URL+"/api/public/bookings", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}

func TestCreateBooking_InvalidSlot(t *testing.T) {
	s := store.New()
	et := createEventType(t, s, "Test", 15)

	ts := newServer(s)
	defer ts.Close()

	body := fmt.Sprintf(`{"slotId":"invalid","eventTypeId":"%s","guestName":"Иван","guestEmail":"ivan@example.com"}`, et.ID)
	resp, err := http.Post(ts.URL+"/api/public/bookings", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestCreateBooking_InvalidJSON(t *testing.T) {
	_, ts := setupTest(t)

	resp, err := http.Post(ts.URL+"/api/public/bookings", "application/json", strings.NewReader("not json"))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestEventTypes_List(t *testing.T) {
	s := store.New()
	createEventType(t, s, "Meeting 1", 15)
	createEventType(t, s, "Meeting 2", 30)
	ts := newServer(s)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/event-types")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var types []model.EventType
	json.NewDecoder(resp.Body).Decode(&types)

	if len(types) != 2 {
		t.Fatalf("expected 2 event types, got %d", len(types))
	}
}

func TestEventTypes_Create(t *testing.T) {
	_, ts := setupTest(t)

	body := `{"name":"New Meeting","duration":15,"description":"test desc"}`
	resp, err := http.Post(ts.URL+"/api/event-types", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var et model.EventType
	json.NewDecoder(resp.Body).Decode(&et)

	if et.Name != "New Meeting" {
		t.Fatalf("expected name 'New Meeting', got %s", et.Name)
	}
	if et.Duration != 15 {
		t.Fatalf("expected duration 15, got %d", et.Duration)
	}
	if et.OwnerID != "owner-1" {
		t.Fatalf("expected owner-1, got %s", et.OwnerID)
	}
}

func TestEventTypes_Create_InvalidDuration(t *testing.T) {
	_, ts := setupTest(t)

	body := `{"name":"Bad Meeting","duration":60}`
	resp, err := http.Post(ts.URL+"/api/event-types", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestEventTypes_Create_EmptyName(t *testing.T) {
	_, ts := setupTest(t)

	body := `{"name":"","duration":15}`
	resp, err := http.Post(ts.URL+"/api/event-types", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestEventTypes_Create_InvalidJSON(t *testing.T) {
	_, ts := setupTest(t)

	resp, err := http.Post(ts.URL+"/api/event-types", "application/json", strings.NewReader("not json"))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestEventTypes_Update(t *testing.T) {
	s := store.New()
	et := createEventType(t, s, "Original", 15)
	ts := newServer(s)
	defer ts.Close()

	body := `{"name":"Updated","duration":30}`
	req, err := http.NewRequest(http.MethodPut, ts.URL+"/api/event-types/"+et.ID, strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var updated model.EventType
	json.NewDecoder(resp.Body).Decode(&updated)

	if updated.Name != "Updated" {
		t.Fatalf("expected 'Updated', got %s", updated.Name)
	}
	if updated.Duration != 30 {
		t.Fatalf("expected duration 30, got %d", updated.Duration)
	}
}

func TestEventTypes_Update_NotFound(t *testing.T) {
	_, ts := setupTest(t)

	body := `{"name":"Updated","duration":15}`
	req, err := http.NewRequest(http.MethodPut, ts.URL+"/api/event-types/non-existent", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}

func TestEventTypes_Update_InvalidJSON(t *testing.T) {
	s := store.New()
	et := createEventType(t, s, "Test", 15)
	ts := newServer(s)
	defer ts.Close()

	req, err := http.NewRequest(http.MethodPut, ts.URL+"/api/event-types/"+et.ID, strings.NewReader("not json"))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestEventTypes_Delete(t *testing.T) {
	s := store.New()
	et := createEventType(t, s, "To Delete", 15)
	ts := newServer(s)
	defer ts.Close()

	req, err := http.NewRequest(http.MethodDelete, ts.URL+"/api/event-types/"+et.ID, nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	_, ok := s.GetEventType(et.ID)
	if ok {
		t.Fatal("event type should have been deleted")
	}
}

func TestEventTypes_Delete_CascadesBookings(t *testing.T) {
	s := store.New()
	et := createEventType(t, s, "To Delete", 15)

	slotDate := nextWeekdayDate(1)
	slotTime := time.Date(slotDate.Year(), slotDate.Month(), slotDate.Day(), 10, 0, 0, 0, time.UTC)
	slotID := makeSlotID(et.ID, slotTime)

	s.CreateBooking(model.Booking{
		ID:          s.NextBookingID(),
		SlotID:      slotID,
		EventTypeID: et.ID,
		GuestName:   "Guest",
		GuestEmail:  "guest@example.com",
	})

	ts := newServer(s)
	defer ts.Close()

	req, err := http.NewRequest(http.MethodDelete, ts.URL+"/api/event-types/"+et.ID, nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	if s.IsSlotBooked(slotID) {
		t.Fatal("slot should have been freed after event type deletion")
	}

	bookings := s.ListBookingsByOwner("owner-1")
	if len(bookings) != 0 {
		t.Fatalf("expected 0 bookings after cascade delete, got %d", len(bookings))
	}
}

func TestEventTypes_Delete_NotFound(t *testing.T) {
	_, ts := setupTest(t)

	req, err := http.NewRequest(http.MethodDelete, ts.URL+"/api/event-types/non-existent", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}

func TestOwnerBookings_List(t *testing.T) {
	s := store.New()
	et := createEventType(t, s, "Test", 15)

	slotDate := nextWeekdayDate(1)
	slotTime := time.Date(slotDate.Year(), slotDate.Month(), slotDate.Day(), 10, 0, 0, 0, time.UTC)
	slotID := makeSlotID(et.ID, slotTime)

	s.CreateBooking(model.Booking{
		ID:          s.NextBookingID(),
		SlotID:      slotID,
		EventTypeID: et.ID,
		GuestName:   "Guest",
		GuestEmail:  "guest@example.com",
	})

	ts := newServer(s)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/owners/owner-1/bookings")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var enriched []model.BookingEnriched
	json.NewDecoder(resp.Body).Decode(&enriched)

	if len(enriched) != 1 {
		t.Fatalf("expected 1 booking, got %d", len(enriched))
	}

	if enriched[0].EventTypeName != "Test" {
		t.Fatalf("expected event type name 'Test', got %s", enriched[0].EventTypeName)
	}
	if enriched[0].GuestName != "Guest" {
		t.Fatalf("expected guest name 'Guest', got %s", enriched[0].GuestName)
	}
}

func TestOwnerBookings_Empty(t *testing.T) {
	_, ts := setupTest(t)

	resp, err := http.Get(ts.URL + "/api/owners/owner-1/bookings")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var enriched []model.BookingEnriched
	json.NewDecoder(resp.Body).Decode(&enriched)

	if len(enriched) != 0 {
		t.Fatalf("expected 0 bookings, got %d", len(enriched))
	}
}

func TestOwnerBookings_OtherOwner(t *testing.T) {
	s := store.New()
	et := createEventType(t, s, "Test", 15)
	slotDate := nextWeekdayDate(1)
	slotTime := time.Date(slotDate.Year(), slotDate.Month(), slotDate.Day(), 10, 0, 0, 0, time.UTC)
	slotID := makeSlotID(et.ID, slotTime)

	s.CreateBooking(model.Booking{
		ID:          s.NextBookingID(),
		SlotID:      slotID,
		EventTypeID: et.ID,
		GuestName:   "Guest",
		GuestEmail:  "guest@example.com",
	})

	ts := newServer(s)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/owners/other-owner/bookings")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var enriched []model.BookingEnriched
	json.NewDecoder(resp.Body).Decode(&enriched)

	if len(enriched) != 0 {
		t.Fatalf("expected 0 bookings for other owner, got %d", len(enriched))
	}
}

func TestOwnerDeleteBooking(t *testing.T) {
	s := store.New()
	et := createEventType(t, s, "Test", 15)

	slotDate := nextWeekdayDate(1)
	slotTime := time.Date(slotDate.Year(), slotDate.Month(), slotDate.Day(), 10, 0, 0, 0, time.UTC)
	slotID := makeSlotID(et.ID, slotTime)

	booking := model.Booking{
		ID:          s.NextBookingID(),
		SlotID:      slotID,
		EventTypeID: et.ID,
		GuestName:   "Guest",
		GuestEmail:  "guest@example.com",
	}
	s.CreateBooking(booking)

	ts := newServer(s)
	defer ts.Close()

	req, err := http.NewRequest(http.MethodDelete, ts.URL+"/api/bookings/"+booking.ID, nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	if s.IsSlotBooked(slotID) {
		t.Fatal("slot should have been freed after booking deletion")
	}
}

func TestOwnerDeleteBooking_NotFound(t *testing.T) {
	_, ts := setupTest(t)

	req, err := http.NewRequest(http.MethodDelete, ts.URL+"/api/bookings/non-existent", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}

func TestMakeSlotID(t *testing.T) {
	now := time.Date(2026, 5, 25, 10, 0, 0, 0, time.UTC)
	id := makeSlotID("event-1", now)
	expected := "event-1-" + fmt.Sprintf("%d", now.Unix())
	if id != expected {
		t.Fatalf("expected %s, got %s", expected, id)
	}
}

func TestParseSlotID_Valid(t *testing.T) {
	now := time.Date(2026, 5, 25, 10, 0, 0, 0, time.UTC)
	id := makeSlotID("event-1", now)

	startTime, date, err := parseSlotID(id)
	if err != nil {
		t.Fatal(err)
	}

	if !startTime.Equal(now) {
		t.Fatalf("expected startTime %v, got %v", now, startTime)
	}

	expectedDate := now.Truncate(24 * time.Hour)
	if !date.Equal(expectedDate) {
		t.Fatalf("expected date %v, got %v", expectedDate, date)
	}
}

func TestParseSlotID_Invalid(t *testing.T) {
	cases := []string{
		"",
		"no-dash",
		"event-1-notanumber",
	}
	for _, c := range cases {
		_, _, err := parseSlotID(c)
		if err == nil {
			t.Fatalf("expected error for input %q", c)
		}
	}
}

func TestIsWeekday(t *testing.T) {
	if isWeekday(time.Date(2026, 5, 23, 0, 0, 0, 0, time.UTC)) {
		t.Fatal("2026-05-23 is Saturday, expected false")
	}
	if isWeekday(time.Date(2026, 5, 24, 0, 0, 0, 0, time.UTC)) {
		t.Fatal("2026-05-24 is Sunday, expected false")
	}
	if !isWeekday(time.Date(2026, 5, 25, 0, 0, 0, 0, time.UTC)) {
		t.Fatal("2026-05-25 is Monday, expected true")
	}
}

func TestEnumerateSlots_Weekday(t *testing.T) {
	et := model.EventType{ID: "event-1", Duration: 15}
	date := time.Date(2026, 5, 25, 0, 0, 0, 0, time.UTC)

	slots := enumerateSlots(et, date)
	expected := 9 * 60 / et.Duration
	if len(slots) != expected {
		t.Fatalf("expected %d slots for 15min duration, got %d", expected, len(slots))
	}

	for i, slot := range slots {
		if slot.EventTypeID != "event-1" {
			t.Fatalf("slot %d: expected eventTypeID event-1, got %s", i, slot.EventTypeID)
		}
	}
}

func TestEnumerateSlots_Weekend(t *testing.T) {
	et := model.EventType{ID: "event-1", Duration: 15}

	saturday := time.Date(2026, 5, 23, 0, 0, 0, 0, time.UTC)
	slots := enumerateSlots(et, saturday)
	if len(slots) != 0 {
		t.Fatalf("expected 0 slots on Saturday, got %d", len(slots))
	}

	sunday := time.Date(2026, 5, 24, 0, 0, 0, 0, time.UTC)
	slots = enumerateSlots(et, sunday)
	if len(slots) != 0 {
		t.Fatalf("expected 0 slots on Sunday, got %d", len(slots))
	}
}

func TestEnumerateSlots_30Min(t *testing.T) {
	et := model.EventType{ID: "event-1", Duration: 30}
	date := time.Date(2026, 5, 25, 0, 0, 0, 0, time.UTC)

	slots := enumerateSlots(et, date)
	expected := 9 * 60 / et.Duration
	if len(slots) != expected {
		t.Fatalf("expected %d slots for 30min duration, got %d", expected, len(slots))
	}
}

func TestEnumerateSlots_StartEndTimes(t *testing.T) {
	et := model.EventType{ID: "event-1", Duration: 15}
	date := time.Date(2026, 5, 25, 0, 0, 0, 0, time.UTC)

	slots := enumerateSlots(et, date)
	if len(slots) == 0 {
		t.Fatal("expected at least one slot")
	}

	firstSlot := slots[0]
	expectedStart := "2026-05-25T09:00:00Z"
	expectedEnd := "2026-05-25T09:15:00Z"
	if firstSlot.StartTime != expectedStart {
		t.Fatalf("first slot start: expected %s, got %s", expectedStart, firstSlot.StartTime)
	}
	if firstSlot.EndTime != expectedEnd {
		t.Fatalf("first slot end: expected %s, got %s", expectedEnd, firstSlot.EndTime)
	}

	lastSlot := slots[len(slots)-1]
	expectedLastStart := "2026-05-25T17:45:00Z"
	expectedLastEnd := "2026-05-25T18:00:00Z"
	if lastSlot.StartTime != expectedLastStart {
		t.Fatalf("last slot start: expected %s, got %s", expectedLastStart, lastSlot.StartTime)
	}
	if lastSlot.EndTime != expectedLastEnd {
		t.Fatalf("last slot end: expected %s, got %s", expectedLastEnd, lastSlot.EndTime)
	}
}

func TestIsSlotValid_Valid(t *testing.T) {
	s := store.New()
	et := createEventType(t, s, "Test", 15)
	date := nextWeekdayDate(1)
	slotTime := time.Date(date.Year(), date.Month(), date.Day(), 10, 0, 0, 0, time.UTC)

	if !isSlotValid(et, slotTime, date, s) {
		t.Fatal("expected slot to be valid")
	}
}

func TestIsSlotValid_Invalid(t *testing.T) {
	s := store.New()
	et := createEventType(t, s, "Test", 15)

	date := nextWeekdayDate(1)
	slotTime := time.Date(date.Year(), date.Month(), date.Day(), 18, 30, 0, 0, time.UTC)

	if isSlotValid(et, slotTime, date, s) {
		t.Fatal("expected slot after work hours to be invalid")
	}
}
