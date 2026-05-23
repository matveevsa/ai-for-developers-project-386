<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { listUpcomingBookings } from '@/api/bookings'
import BookingListItem from '@/components/owner/BookingListItem.vue'
import type { Booking } from '@/types/api'

const ownerId = 'owner-1'
const bookings = ref<Booking[]>([])
const loading = ref(false)
const error = ref<string | null>(null)

onMounted(async () => {
  loading.value = true
  try {
    bookings.value = await listUpcomingBookings(ownerId)
  } catch (e: any) {
    error.value = e?.response?.data?.message ?? 'Ошибка загрузки'
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="page">
    <h1>Предстоящие встречи</h1>

    <div v-if="loading">Загрузка...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <div v-else-if="bookings.length === 0" class="empty">
      Нет предстоящих встреч
    </div>
    <div v-else class="list">
      <BookingListItem
        v-for="b in bookings"
        :key="b.id"
        :booking="b"
      />
    </div>
  </div>
</template>

<style scoped>
.list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}
.error {
  color: var(--mantine-color-red-7);
}
.empty {
  color: var(--mantine-color-gray-5);
}
</style>
