import { ref } from 'vue'
import { getAvailableSlots } from '@/api/public'
import type { Slot } from '@/types/api'

export function useSlots() {
  const slots = ref<Slot[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function load(eventTypeId: string, date: string) {
    loading.value = true
    error.value = null
    try {
      slots.value = await getAvailableSlots(eventTypeId, date)
    } catch (e: any) {
      error.value = e?.response?.data?.message ?? 'Ошибка загрузки слотов'
    } finally {
      loading.value = false
    }
  }

  return { slots, loading, error, load }
}
