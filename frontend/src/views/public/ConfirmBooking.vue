<script setup lang="ts">
import { useRoute } from 'vue-router'
import { useBooking } from '@/composables/useBooking'
import BookingForm from '@/components/public/BookingForm.vue'
import type { BookingCreate } from '@/types/api'

const route = useRoute()
const slotId = route.query.slotId as string
const eventTypeId = route.params.eventTypeId as string

const { booking, loading, error, create } = useBooking()

async function handleSubmit(data: BookingCreate) {
  await create({
    ...data,
    slotId,
    eventTypeId,
  })
}
</script>

<template>
  <div class="page">
    <h1>Подтверждение бронирования</h1>

    <div v-if="booking" class="success">
      <h2>Бронирование подтверждено!</h2>
      <p>Спасибо, {{ booking.guestName }}.</p>
    </div>

    <div v-else>
      <BookingForm
        v-if="slotId && eventTypeId"
        :slot-id="slotId"
        :event-type-id="eventTypeId"
        @submit="handleSubmit"
      />

      <div v-if="loading" class="loading">Отправка...</div>
      <div v-if="error" class="error">{{ error }}</div>
    </div>
  </div>
</template>

<style scoped>
.success {
  text-align: center;
  padding: 2rem;
}
.loading {
  text-align: center;
  color: var(--mantine-color-gray-5);
}
.error {
  color: var(--mantine-color-red-7);
}
</style>
