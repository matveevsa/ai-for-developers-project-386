package model

type EventType struct {
	ID          string `json:"id"`
	OwnerID     string `json:"ownerId"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Duration    int    `json:"duration"`
}

type EventTypeCreate struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Duration    int    `json:"duration"`
}

type Slot struct {
	ID          string `json:"id"`
	EventTypeID string `json:"eventTypeId"`
	StartTime   string `json:"startTime"`
	EndTime     string `json:"endTime"`
}

type Booking struct {
	ID          string `json:"id"`
	SlotID      string `json:"slotId"`
	EventTypeID string `json:"eventTypeId"`
	GuestName   string `json:"guestName"`
	GuestEmail  string `json:"guestEmail"`
	GuestNotes  string `json:"guestNotes,omitempty"`
}

type BookingCreate struct {
	SlotID      string `json:"slotId"`
	EventTypeID string `json:"eventTypeId"`
	GuestName   string `json:"guestName"`
	GuestEmail  string `json:"guestEmail"`
	GuestNotes  string `json:"guestNotes,omitempty"`
}

type BookingEnriched struct {
	ID            string `json:"id"`
	SlotID        string `json:"slotId"`
	EventTypeID   string `json:"eventTypeId"`
	GuestName     string `json:"guestName"`
	GuestEmail    string `json:"guestEmail"`
	GuestNotes    string `json:"guestNotes,omitempty"`
	StartTime     string `json:"startTime"`
	EndTime       string `json:"endTime"`
	EventTypeName string `json:"eventTypeName"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
