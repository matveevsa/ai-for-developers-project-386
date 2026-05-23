<script setup lang="ts">
import { computed } from 'vue'
import dayjs from 'dayjs'
import type { Slot } from '@/types/api'
import SlotButton from './SlotButton.vue'

const props = defineProps<{
  slots: Slot[]
  selectedSlotId?: string
}>()

const emit = defineEmits<{
  selectSlot: [slotId: string]
}>()

const days = computed(() => {
  const result: { date: string; label: string; slots: Slot[] }[] = []
  const today = dayjs()
  for (let i = 0; i < 14; i++) {
    const date = today.add(i, 'day')
    const dateStr = date.format('YYYY-MM-DD')
    result.push({
      date: dateStr,
      label: i === 0 ? 'Сегодня' : date.format('DD.MM'),
      slots: props.slots.filter((s) => dayjs(s.startTime).format('YYYY-MM-DD') === dateStr),
    })
  }
  return result
})

function isPast(date: string): boolean {
  return dayjs(date).isBefore(dayjs(), 'day')
}
</script>

<template>
  <div class="calendar">
    <div v-for="day in days" :key="day.date" class="day-col" :class="{ dim: isPast(day.date) && day.slots.length === 0 }">
      <div class="day-header">{{ day.label }}</div>
      <div class="slots">
        <div v-if="day.slots.length === 0" class="no-slots">Нет слотов</div>
        <SlotButton
          v-for="slot in day.slots"
          :key="slot.id"
          :slot="slot"
          :selected="slot.id === selectedSlotId"
          @select="emit('selectSlot', $event)"
        />
      </div>
    </div>
  </div>
</template>

<style scoped>
.calendar {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 0.5rem;
}
.day-col {
  border: 1px solid var(--mantine-color-gray-3);
  border-radius: var(--mantine-radius-md);
  padding: 0.5rem;
  min-height: 120px;
}
.day-col.dim {
  opacity: 0.4;
}
.day-header {
  font-weight: 600;
  font-size: 0.875rem;
  margin-bottom: 0.5rem;
  text-align: center;
}
.slots {
  display: flex;
  flex-direction: column;
  gap: 0.375rem;
}
.no-slots {
  font-size: 0.75rem;
  color: var(--mantine-color-gray-5);
  text-align: center;
}
</style>
