<script setup lang="ts">
import type { Slot } from '@/types/api'
import dayjs from 'dayjs'

defineProps<{
  slot: Slot
  selected?: boolean
}>()

const emit = defineEmits<{
  select: [slotId: string]
}>()
</script>

<template>
  <button
    class="slot-btn"
    :class="{ selected }"
    @click="emit('select', slot.id)"
  >
    {{ dayjs(slot.startTime).format('HH:mm') }}
  </button>
</template>

<style scoped>
.slot-btn {
  padding: 0.5rem 1rem;
  border: 1px solid var(--mantine-color-blue-3);
  border-radius: var(--mantine-radius-sm);
  background: white;
  color: var(--mantine-color-blue-8);
  cursor: pointer;
  font-size: 0.875rem;
  transition: 0.15s;
}
.slot-btn:hover {
  background: var(--mantine-color-blue-0);
}
.slot-btn.selected {
  background: var(--mantine-color-blue-6);
  color: white;
  border-color: var(--mantine-color-blue-6);
}
</style>
