import client from './client'
import type { EventType, Slot, Booking, BookingCreate } from '@/types/api'

export function listPublicEventTypes(): Promise<EventType[]> {
  return client.get('/public/event-types').then((r) => r.data)
}

export function getAvailableSlots(
  eventTypeId: string,
  date: string,
): Promise<Slot[]> {
  return client
    .get(`/public/event-types/${eventTypeId}/slots`, { params: { date } })
    .then((r) => r.data)
}

export function createPublicBooking(data: BookingCreate): Promise<Booking> {
  return client.post('/public/bookings', data).then((r) => r.data)
}
