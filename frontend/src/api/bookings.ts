import client from './client'
import type { Booking } from '@/types/api'

export function listUpcomingBookings(ownerId: string): Promise<Booking[]> {
  return client.get(`/owners/${ownerId}/bookings`).then((r) => r.data)
}
