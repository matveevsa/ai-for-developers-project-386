import { ref } from 'vue'
import { createPublicBooking } from '@/api/public'
import type { Booking, BookingCreate } from '@/types/api'

export function useBooking() {
  const booking = ref<Booking | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function create(data: BookingCreate) {
    loading.value = true
    error.value = null
    try {
      booking.value = await createPublicBooking(data)
    } catch (e: any) {
      error.value = e?.response?.data?.message ?? 'Ошибка бронирования'
    } finally {
      loading.value = false
    }
  }

  return { booking, loading, error, create }
}
