<script setup lang="ts">
import dayjs from 'dayjs'

const props = defineProps<{
  selectedDate?: string
  availability?: Record<string, { available: number; total: number }>
}>()

const emit = defineEmits<{
  selectDate: [date: string]
}>()

const dayNames = ['Вс', 'Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб']

const days = Array.from({ length: 14 }, (_, i) => {
  const date = dayjs().add(i, 'day')
  return {
    date: date.format('YYYY-MM-DD'),
    label: i === 0 ? 'Сегодня' : date.format('DD.MM'),
    isPast: date.isBefore(dayjs(), 'day'),
    isWeekend: date.day() === 0 || date.day() === 6,
  }
})

function dayClass(day: (typeof days)[0]) {
  if (day.date === props.selectedDate) return 'day selected'

  if (props.availability) {
    const a = props.availability[day.date]
    if (a && a.total > 0 && a.available === 0) return 'day full'
    if (a && a.available > 0) return 'day free'
  }

  if (day.isPast || day.isWeekend) return 'day dim'
  return 'day'
}

function dayCount(day: (typeof days)[0]) {
  if (!props.availability) return
  const a = props.availability[day.date]
  if (!a || a.total === 0) return
  if (a.available > 0) return `слотов: ${a.available}`
  return 'занято'
}
</script>

<template>
  <div class="calendar">
    <div v-for="d in dayNames" :key="d" class="day-header">{{ d }}</div>
    <div
      v-for="day in days"
      :key="day.date"
      :class="dayClass(day)"
      @click="emit('selectDate', day.date)"
    >
      <div class="day-label">{{ day.label }}</div>
      <div v-if="dayCount(day)" class="day-count">{{ dayCount(day) }}</div>
    </div>
  </div>
</template>

<style scoped>
.calendar {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 0.5rem;
}
.day {
  padding: 0.75rem 0.5rem;
  border-radius: var(--mantine-radius-md);
  text-align: center;
  cursor: pointer;
  transition: 0.15s;
  border: 2px solid var(--mantine-color-gray-3);
  background: white;
}
.day:hover {
  border-color: var(--mantine-color-orange-4);
}
.day-label {
  font-weight: 600;
  font-size: 0.875rem;
}
.day-header {
  font-size: 0.75rem;
  font-weight: 600;
  text-align: center;
  color: var(--mantine-color-gray-5);
  padding: 0.25rem 0;
}
.day-count {
  font-size: 0.65rem;
  margin-top: 0.125rem;
}
.day.dim {
  opacity: 0.35;
  cursor: default;
}
.day.selected {
  background: var(--mantine-color-orange-1);
  border-color: var(--mantine-color-orange-6);
}
.day.free {
  background: #d1fae5;
  border-color: #a7f3d0;
}
.day.free:hover {
  border-color: #059669;
}
.day.full {
  background: #fee2e2;
  border-color: #fecaca;
}
.day.full:hover {
  border-color: #dc2626;
}
</style>
