<script setup lang="ts">
import { ref } from 'vue'
import { useEventTypes } from '@/composables/useEventTypes'
import EventTypeCard from '@/components/owner/EventTypeCard.vue'
import EventTypeForm from '@/components/owner/EventTypeForm.vue'
import type { EventTypeCreate } from '@/types/api'

const { eventTypes, loading, error, create } = useEventTypes()
const showForm = ref(false)

async function handleCreate(data: EventTypeCreate) {
  await create(data)
  showForm.value = false
}
</script>

<template>
  <div class="page">
    <div class="page-header">
      <h1>Типы событий</h1>
      <button class="btn btn-primary" @click="showForm = !showForm">
        {{ showForm ? 'Отмена' : 'Создать' }}
      </button>
    </div>

    <EventTypeForm v-if="showForm" @submit="handleCreate" />

    <div v-if="loading">Загрузка...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <div v-else-if="eventTypes.length === 0" class="empty">
      Нет типов событий
    </div>
    <div v-else class="grid">
      <EventTypeCard
        v-for="et in eventTypes"
        :key="et.id"
        :event-type="et"
      />
    </div>
  </div>
</template>

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}
.page-header h1 {
  margin: 0;
}
.grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 1rem;
}
.error {
  color: var(--mantine-color-red-7);
}
.empty {
  color: var(--mantine-color-gray-5);
}
</style>
