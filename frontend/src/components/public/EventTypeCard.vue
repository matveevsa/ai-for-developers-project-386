<script setup lang="ts">
import type { EventType } from '@/types/api'

defineProps<{
  eventType: EventType
}>()
</script>

<template>
  <div class="card">
    <div class="card-body">
      <h3>{{ eventType.name }}</h3>
      <p v-if="eventType.description" class="muted">
        {{ eventType.description }}
      </p>
      <span class="badge">{{ eventType.duration }} мин</span>
    </div>
    <div class="card-footer">
      <slot name="action" />
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
.card-body h3 {
  margin: 0 0 0.25rem;
  font-size: 1.125rem;
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
  background: var(--mantine-color-green-1);
  color: var(--mantine-color-green-8);
  border-radius: var(--mantine-radius-xl);
}
.card-footer {
  padding: 0.75rem 1rem;
  border-top: 1px solid var(--mantine-color-gray-2);
}
</style>
