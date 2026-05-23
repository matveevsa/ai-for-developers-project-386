<script setup lang="ts">
import { ref } from 'vue'
import type { EventTypeCreate } from '@/types/api'

const emit = defineEmits<{
  submit: [data: EventTypeCreate]
}>()

const name = ref('')
const description = ref('')
const duration = ref(30)

function onSubmit() {
  emit('submit', {
    name: name.value,
    description: description.value || undefined,
    duration: duration.value,
  })
}
</script>

<template>
  <form @submit.prevent="onSubmit">
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
    <button type="submit" class="btn btn-primary">Создать</button>
  </form>
</template>

<style scoped>
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
</style>
