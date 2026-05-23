<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useEventTypes } from '@/composables/useEventTypes'
import EventTypeCard from '@/components/public/EventTypeCard.vue'

const router = useRouter()
const { eventTypes, loading, error } = useEventTypes()
</script>

<template>
  <div class="page">
    <h1>Выберите тип встречи</h1>

    <div v-if="loading">Загрузка...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <div v-else class="grid">
      <EventTypeCard
        v-for="et in eventTypes"
        :key="et.id"
        :event-type="et"
      >
        <template #action>
          <button
            class="btn btn-primary btn-sm"
            @click="router.push(`/book/${et.id}`)"
          >
            Выбрать
          </button>
        </template>
      </EventTypeCard>
    </div>
  </div>
</template>

<style scoped>
.grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1rem;
}
.error {
  color: var(--mantine-color-red-7);
}
</style>
