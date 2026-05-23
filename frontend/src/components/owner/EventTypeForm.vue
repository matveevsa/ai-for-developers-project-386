<script setup lang="ts">
import { ref, watch } from 'vue'
import type { EventType, EventTypeCreate } from '@/types/api'

const props = defineProps<{
  eventType?: EventType | null
}>()

const emit = defineEmits<{
  submit: [data: EventTypeCreate]
  cancel: []
}>()

const name = ref('')
const description = ref('')
const duration = ref(30)

watch(() => props.eventType, (et) => {
  if (et) {
    name.value = et.name
    description.value = et.description ?? ''
    duration.value = et.duration
  } else {
    name.value = ''
    description.value = ''
    duration.value = 30
  }
}, { immediate: true })

function onSubmit() {
  emit('submit', {
    name: name.value,
    description: description.value || undefined,
    duration: duration.value,
  })
}
</script>

<template>
  <form @submit.prevent="onSubmit" class="form">
    <div class="field">
      <label>Название</label>
      <input v-model="name" required minlength="1" maxlength="100" />
    </div>
    <div class="field">
      <label>Описание</label>
      <textarea v-model="description" />
    </div>
    <div class="field">
      <label>Длительность (мин)</label>
      <input v-model.number="duration" type="number" min="15" max="480" required />
    </div>
    <div class="actions">
      <button type="submit" class="btn btn-primary">
        {{ eventType ? 'Сохранить' : 'Создать' }}
      </button>
      <button type="button" class="btn btn-secondary" @click="emit('cancel')">
        Отмена
      </button>
    </div>
  </form>
</template>

<style scoped>
.form {
  margin-bottom: 1rem;
  padding: 1rem;
  border: 1px solid var(--mantine-color-gray-3);
  border-radius: var(--mantine-radius-md);
  background: white;
}
.field {
  margin-bottom: 1rem;
}
.field label {
  display: block;
  margin-bottom: 0.25rem;
  font-weight: 500;
  font-size: 0.875rem;
}
.field input,
.field textarea {
  width: 100%;
  padding: 0.5rem 0.75rem;
  border: 1px solid var(--mantine-color-gray-4);
  border-radius: var(--mantine-radius-sm);
  font-size: 0.875rem;
}
.field textarea {
  min-height: 80px;
  resize: vertical;
}
.actions {
  display: flex;
  gap: 0.5rem;
}
</style>
