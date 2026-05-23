import client from './client'
import type { BookingEnriched } from '@/types/api'

export function deleteBooking(id: string): Promise<void> {
  return client.delete(`/bookings/${id}`).then(() => {})
}

export function listUpcomingBookings(ownerId: string): Promise<BookingEnriched[]> {
  return client.get(`/owners/${ownerId}/bookings`).then((r) => r.data)
}
