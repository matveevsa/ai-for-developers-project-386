import client from './client'
import type { EventType, EventTypeCreate } from '@/types/api'

export function listEventTypes(): Promise<EventType[]> {
  return client.get('/event-types').then((r) => r.data)
}

export function createEventType(data: EventTypeCreate): Promise<EventType> {
  return client.post('/event-types', data).then((r) => r.data)
}

export function updateEventType(id: string, data: EventTypeCreate): Promise<EventType> {
  return client.put(`/event-types/${id}`, data).then((r) => r.data)
}

export function deleteEventType(id: string): Promise<void> {
  return client.delete(`/event-types/${id}`)
}
