<script setup lang="ts">
import { ref } from 'vue'
import type { BookingCreate } from '@/types/api'

const emit = defineEmits<{
  submit: [data: BookingCreate]
}>()

const guestName = ref('')
const guestEmail = ref('')
const guestNotes = ref('')

defineProps<{
  slotId: string
  eventTypeId: string
}>()

function onSubmit() {
  emit('submit', {
    slotId: '',
    eventTypeId: '',
    guestName: guestName.value,
    guestEmail: guestEmail.value,
    guestNotes: guestNotes.value || undefined,
  })
}
</script>

<template>
  <form @submit.prevent="onSubmit">
    <div class="field">
      <label>Имя</label>
      <input v-model="guestName" required minlength="1" maxlength="100" />
    </div>
    <div class="field">
      <label>Email</label>
      <input v-model="guestEmail" type="email" required />
    </div>
    <div class="field">
      <label>Заметки (необязательно)</label>
      <textarea v-model="guestNotes" />
    </div>
    <button type="submit" class="btn btn-primary">Забронировать</button>
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
