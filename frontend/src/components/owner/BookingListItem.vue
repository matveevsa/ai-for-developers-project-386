<script setup lang="ts">
import type { Booking } from '@/types/api'
import dayjs from 'dayjs'

defineProps<{
  booking: Booking
}>()
</script>

<template>
  <div class="item">
    <div class="item-info">
      <strong>{{ booking.guestName }}</strong>
      <span class="muted">{{ booking.guestEmail }}</span>
    </div>
    <div class="item-notes" v-if="booking.guestNotes">
      {{ booking.guestNotes }}
    </div>
  </div>
</template>

<style scoped>
.item {
  padding: 0.75rem 1rem;
  border: 1px solid var(--mantine-color-gray-3);
  border-radius: var(--mantine-radius-md);
}
.item-info {
  display: flex;
  gap: 0.75rem;
  align-items: center;
}
.muted {
  color: var(--mantine-color-gray-6);
  font-size: 0.875rem;
}
.item-notes {
  margin-top: 0.25rem;
  font-size: 0.875rem;
  color: var(--mantine-color-gray-6);
}
</style>
