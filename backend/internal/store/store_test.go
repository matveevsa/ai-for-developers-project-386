package store

import (
	"testing"

	"github.com/matveevsa/ai-for-developers-project-386/backend/internal/model"
)

func TestNew(t *testing.T) {
	s := New()
	if s == nil {
		t.Fatal("New() returned nil")
	}
}

func TestNextEventID(t *testing.T) {
	s := New()
	id1 := s.NextEventID()
	id2 := s.NextEventID()

	if id1 == id2 {
		t.Fatal("expected unique event IDs")
	}

	if id1 != "event-1" {
		t.Fatalf("expected event-1, got %s", id1)
	}
	if id2 != "event-2" {
		t.Fatalf("expected event-2, got %s", id2)
	}
}

func TestNextBookingID(t *testing.T) {
	s := New()
	id1 := s.NextBookingID()
	id2 := s.NextBookingID()

	if id1 == id2 {
		t.Fatal("expected unique booking IDs")
	}

	if id1 != "booking-1" {
		t.Fatalf("expected booking-1, got %s", id1)
	}
	if id2 != "booking-2" {
		t.Fatalf("expected booking-2, got %s", id2)
	}
}

func TestCreateAndGetEventType(t *testing.T) {
	s := New()
	et := model.EventType{
		ID:          s.NextEventID(),
		OwnerID:     "owner-1",
		Name:        "Test Meeting",
		Description: "A test",
		Duration:    15,
	}
	s.CreateEventType(et)

	got, ok := s.GetEventType(et.ID)
	if !ok {
		t.Fatal("event type not found")
	}
	if got.Name != "Test Meeting" {
		t.Fatalf("expected name 'Test Meeting', got %s", got.Name)
	}
	if got.Duration != 15 {
		t.Fatalf("expected duration 15, got %d", got.Duration)
	}
}

func TestGetEventType_NotFound(t *testing.T) {
	s := New()
	_, ok := s.GetEventType("non-existent")
	if ok {
		t.Fatal("expected false for non-existent event type")
	}
}

func TestListEventTypes(t *testing.T) {
	s := New()
	if len(s.ListEventTypes()) != 0 {
		t.Fatal("expected empty list initially")
	}

	s.CreateEventType(model.EventType{ID: s.NextEventID(), OwnerID: "owner-1", Name: "A", Duration: 15})
	s.CreateEventType(model.EventType{ID: s.NextEventID(), OwnerID: "owner-1", Name: "B", Duration: 30})

	types := s.ListEventTypes()
	if len(types) != 2 {
		t.Fatalf("expected 2 event types, got %d", len(types))
	}
}

func TestUpdateEventType(t *testing.T) {
	s := New()
	et := model.EventType{
		ID:       s.NextEventID(),
		OwnerID:  "owner-1",
		Name:     "Original",
		Duration: 15,
	}
	s.CreateEventType(et)

	updated, err := s.UpdateEventType(et.ID, model.EventTypeCreate{
		Name:     "Updated",
		Duration: 30,
	})
	if err != nil {
		t.Fatal(err)
	}

	if updated.Name != "Updated" {
		t.Fatalf("expected 'Updated', got %s", updated.Name)
	}
	if updated.Duration != 30 {
		t.Fatalf("expected duration 30, got %d", updated.Duration)
	}
}

func TestUpdateEventType_NotFound(t *testing.T) {
	s := New()
	_, err := s.UpdateEventType("non-existent", model.EventTypeCreate{Name: "X", Duration: 15})
	if err == nil {
		t.Fatal("expected error for non-existent event type")
	}
}

func TestDeleteEventType(t *testing.T) {
	s := New()
	et := model.EventType{ID: s.NextEventID(), OwnerID: "owner-1", Name: "To Delete", Duration: 15}
	s.CreateEventType(et)

	err := s.DeleteEventType(et.ID)
	if err != nil {
		t.Fatal(err)
	}

	_, ok := s.GetEventType(et.ID)
	if ok {
		t.Fatal("event type should have been deleted")
	}
}

func TestDeleteEventType_NotFound(t *testing.T) {
	s := New()
	err := s.DeleteEventType("non-existent")
	if err == nil {
		t.Fatal("expected error for non-existent event type")
	}
}

func TestDeleteEventType_CascadesBookings(t *testing.T) {
	s := New()
	et := model.EventType{ID: s.NextEventID(), OwnerID: "owner-1", Name: "Test", Duration: 15}
	s.CreateEventType(et)

	s.CreateBooking(model.Booking{
		ID:          s.NextBookingID(),
		SlotID:      "event-1-12345",
		EventTypeID: et.ID,
		GuestName:   "Guest",
		GuestEmail:  "guest@example.com",
	})

	err := s.DeleteEventType(et.ID)
	if err != nil {
		t.Fatal(err)
	}

	if s.IsSlotBooked("event-1-12345") {
		t.Fatal("slot should have been freed after event type deletion")
	}

	bookings := s.ListBookingsByOwner("owner-1")
	if len(bookings) != 0 {
		t.Fatalf("expected 0 bookings after cascade, got %d", len(bookings))
	}
}

func TestCreateBooking(t *testing.T) {
	s := New()
	et := model.EventType{ID: s.NextEventID(), OwnerID: "owner-1", Name: "Test", Duration: 15}
	s.CreateEventType(et)

	booking := model.Booking{
		ID:          s.NextBookingID(),
		SlotID:      "event-1-12345",
		EventTypeID: et.ID,
		GuestName:   "Иван",
		GuestEmail:  "ivan@example.com",
	}
	s.CreateBooking(booking)

	if !s.IsSlotBooked("event-1-12345") {
		t.Fatal("slot should be marked as booked")
	}
}

func TestIsSlotBooked(t *testing.T) {
	s := New()
	if s.IsSlotBooked("non-existent") {
		t.Fatal("expected false for non-booked slot")
	}

	s.CreateBooking(model.Booking{
		ID:          s.NextBookingID(),
		SlotID:      "slot-1",
		EventTypeID: "event-1",
		GuestName:   "Guest",
		GuestEmail:  "g@example.com",
	})

	if !s.IsSlotBooked("slot-1") {
		t.Fatal("expected true for booked slot")
	}
}

func TestDeleteBooking(t *testing.T) {
	s := New()
	et := model.EventType{ID: s.NextEventID(), OwnerID: "owner-1", Name: "Test", Duration: 15}
	s.CreateEventType(et)

	booking := model.Booking{
		ID:          s.NextBookingID(),
		SlotID:      "slot-1",
		EventTypeID: et.ID,
		GuestName:   "Guest",
		GuestEmail:  "g@example.com",
	}
	s.CreateBooking(booking)

	err := s.DeleteBooking(booking.ID)
	if err != nil {
		t.Fatal(err)
	}

	if s.IsSlotBooked("slot-1") {
		t.Fatal("slot should be freed after booking deletion")
	}
}

func TestDeleteBooking_NotFound(t *testing.T) {
	s := New()
	err := s.DeleteBooking("non-existent")
	if err == nil {
		t.Fatal("expected error for non-existent booking")
	}
}

func TestListBookingsByOwner(t *testing.T) {
	s := New()
	et := model.EventType{ID: s.NextEventID(), OwnerID: "owner-1", Name: "Test", Duration: 15}
	s.CreateEventType(et)

	s.CreateBooking(model.Booking{
		ID:          s.NextBookingID(),
		SlotID:      "slot-1",
		EventTypeID: et.ID,
		GuestName:   "Guest 1",
		GuestEmail:  "g1@example.com",
	})
	s.CreateBooking(model.Booking{
		ID:          s.NextBookingID(),
		SlotID:      "slot-2",
		EventTypeID: et.ID,
		GuestName:   "Guest 2",
		GuestEmail:  "g2@example.com",
	})

	bookings := s.ListBookingsByOwner("owner-1")
	if len(bookings) != 2 {
		t.Fatalf("expected 2 bookings, got %d", len(bookings))
	}

	otherBookings := s.ListBookingsByOwner("other-owner")
	if len(otherBookings) != 0 {
		t.Fatalf("expected 0 bookings for other owner, got %d", len(otherBookings))
	}
}

func TestListBookingsByOwner_Empty(t *testing.T) {
	s := New()
	bookings := s.ListBookingsByOwner("owner-1")
	if len(bookings) != 0 {
		t.Fatalf("expected 0 bookings, got %d", len(bookings))
	}
}

func TestConcurrency(t *testing.T) {
	s := New()
	done := make(chan bool)

	for range 10 {
		go func() {
			et := model.EventType{ID: s.NextEventID(), OwnerID: "owner-1", Name: "Concurrent", Duration: 15}
			s.CreateEventType(et)
			s.GetEventType(et.ID)
			s.ListEventTypes()
			done <- true
		}()
	}

	for range 10 {
		<-done
	}

	types := s.ListEventTypes()
	if len(types) != 10 {
		t.Fatalf("expected 10 event types, got %d", len(types))
	}
}
