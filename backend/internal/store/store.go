package store

import (
	"fmt"
	"sync"

	"github.com/matveevsa/ai-for-developers-project-386/backend/internal/model"
)

type Store struct {
	mu            sync.RWMutex
	eventTypes    map[string]model.EventType
	bookings      map[string]model.Booking
	bookedSlots   map[string]bool
	nextEventID   int
	nextBookingID int
}

func New() *Store {
	return &Store{
		eventTypes:    make(map[string]model.EventType),
		bookings:      make(map[string]model.Booking),
		bookedSlots:   make(map[string]bool),
		nextEventID:   1,
		nextBookingID: 1,
	}
}

func (s *Store) CreateEventType(et model.EventType) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.eventTypes[et.ID] = et
}

func (s *Store) ListEventTypes() []model.EventType {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]model.EventType, 0, len(s.eventTypes))
	for _, et := range s.eventTypes {
		result = append(result, et)
	}
	return result
}

func (s *Store) GetEventType(id string) (model.EventType, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	et, ok := s.eventTypes[id]
	return et, ok
}

func (s *Store) NextEventID() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := s.nextEventID
	s.nextEventID++
	return fmt.Sprintf("event-%d", id)
}

func (s *Store) NextBookingID() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := s.nextBookingID
	s.nextBookingID++
	return fmt.Sprintf("booking-%d", id)
}

func (s *Store) CreateBooking(b model.Booking) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.bookings[b.ID] = b
	s.bookedSlots[b.SlotID] = true
}

func (s *Store) IsSlotBooked(slotID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.bookedSlots[slotID]
}

func (s *Store) UpdateEventType(id string, update model.EventTypeCreate) (model.EventType, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	et, ok := s.eventTypes[id]
	if !ok {
		return model.EventType{}, fmt.Errorf("event type %s not found", id)
	}

	et.Name = update.Name
	et.Description = update.Description
	et.Duration = update.Duration
	s.eventTypes[id] = et
	return et, nil
}

func (s *Store) DeleteEventType(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.eventTypes[id]; !ok {
		return fmt.Errorf("event type %s not found", id)
	}

	delete(s.eventTypes, id)

	for bid, b := range s.bookings {
		if b.EventTypeID == id {
			delete(s.bookedSlots, b.SlotID)
			delete(s.bookings, bid)
		}
	}

	return nil
}

func (s *Store) DeleteBooking(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	b, ok := s.bookings[id]
	if !ok {
		return fmt.Errorf("booking %s not found", id)
	}

	delete(s.bookedSlots, b.SlotID)
	delete(s.bookings, id)
	return nil
}

func (s *Store) ListBookingsByOwner(ownerID string) []model.Booking {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ownerEventTypes := make(map[string]bool)
	for _, et := range s.eventTypes {
		if et.OwnerID == ownerID {
			ownerEventTypes[et.ID] = true
		}
	}
	result := make([]model.Booking, 0)
	for _, b := range s.bookings {
		if ownerEventTypes[b.EventTypeID] {
			result = append(result, b)
		}
	}
	return result
}
