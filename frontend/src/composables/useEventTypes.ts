import { ref, onMounted } from 'vue'
import { listEventTypes, createEventType, updateEventType, deleteEventType } from '@/api/eventTypes'
import type { EventType, EventTypeCreate } from '@/types/api'

export function useEventTypes() {
  const eventTypes = ref<EventType[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function load() {
    loading.value = true
    error.value = null
    try {
      eventTypes.value = await listEventTypes()
    } catch (e: any) {
      error.value = e?.response?.data?.message ?? 'Ошибка загрузки'
    } finally {
      loading.value = false
    }
  }

  async function create(data: EventTypeCreate) {
    const created = await createEventType(data)
    eventTypes.value.push(created)
    return created
  }

  async function update(id: string, data: EventTypeCreate) {
    const updated = await updateEventType(id, data)
    const idx = eventTypes.value.findIndex(et => et.id === id)
    if (idx !== -1) {
      eventTypes.value[idx] = updated
    }
    return updated
  }

  async function remove(id: string) {
    await deleteEventType(id)
    eventTypes.value = eventTypes.value.filter(et => et.id !== id)
  }

  onMounted(load)

  return { eventTypes, loading, error, load, create, update, remove }
}
