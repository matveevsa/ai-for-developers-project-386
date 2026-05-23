export interface EventType {
  id: string
  ownerId: string
  name: string
  description?: string
  duration: number
}

export interface EventTypeCreate {
  name: string
  description?: string
  duration: number
}

export interface Slot {
  id: string
  eventTypeId: string
  startTime: string
  endTime: string
}

export interface Booking {
  id: string
  slotId: string
  eventTypeId: string
  guestName: string
  guestEmail: string
  guestNotes?: string
}

export interface BookingEnriched {
  id: string
  slotId: string
  eventTypeId: string
  guestName: string
  guestEmail: string
  guestNotes?: string
  startTime: string
  endTime: string
  eventTypeName: string
}

export interface BookingCreate {
  slotId: string
  eventTypeId: string
  guestName: string
  guestEmail: string
  guestNotes?: string
}

export interface ErrorResponse {
  code: number
  message: string
}
