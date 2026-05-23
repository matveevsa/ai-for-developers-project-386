<script setup lang="ts">
import type { EventType } from '@/types/api'

const props = defineProps<{
  eventType: EventType
  onDelete?: (id: string) => void
}>()

function handleDelete() {
  if (confirm(`Удалить тип "${props.eventType.name}"? Все бронирования будут удалены.`)) {
    props.onDelete?.(props.eventType.id)
  }
}
</script>

<template>
  <div class="card">
    <div class="card-body">
      <div class="card-top">
        <h3>{{ eventType.name }}</h3>
        <button class="btn-delete" @click="handleDelete" title="Удалить">✕</button>
      </div>
      <p v-if="eventType.description" class="muted">
        {{ eventType.description }}
      </p>
      <span class="badge">{{ eventType.duration }} мин</span>
    </div>
  </div>
</template>

<style scoped>
.card {
  border: 1px solid var(--mantine-color-gray-3);
  border-radius: var(--mantine-radius-md);
  overflow: hidden;
}
.card-body {
  padding: 1rem;
}
.card-top {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 0.5rem;
  margin-bottom: 0.25rem;
}
.card-top h3 {
  margin: 0;
  font-size: 1rem;
}
.btn-delete {
  background: none;
  border: none;
  color: var(--mantine-color-red-6);
  cursor: pointer;
  font-size: 1.125rem;
  padding: 0.125rem 0.375rem;
  border-radius: var(--mantine-radius-sm);
  transition: 0.15s;
}
.btn-delete:hover {
  background: var(--mantine-color-red-1);
}
.muted {
  color: var(--mantine-color-gray-6);
  font-size: 0.875rem;
  margin: 0 0 0.5rem;
}
.badge {
  display: inline-block;
  padding: 0.125rem 0.5rem;
  font-size: 0.75rem;
  background: var(--mantine-color-orange-1);
  color: var(--mantine-color-orange-8);
  border-radius: var(--mantine-radius-xl);
}
</style>
