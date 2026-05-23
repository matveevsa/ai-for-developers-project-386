<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useSlots } from '@/composables/useSlots'
import CalendarGrid from '@/components/public/CalendarGrid.vue'
import dayjs from 'dayjs'

const route = useRoute()
const router = useRouter()
const { slots, loading, error, load } = useSlots()

const eventTypeId = route.params.eventTypeId as string
const selectedDate = ref(dayjs().format('YYYY-MM-DD'))
const selectedSlotId = ref<string | null>(null)

onMounted(() => {
  load(eventTypeId, selectedDate.value)
})

function selectDate(date: string) {
  selectedDate.value = date
  load(eventTypeId, date)
}

function confirmSlot() {
  if (!selectedSlotId.value) return
  router.push(`/book/${eventTypeId}/confirm?slotId=${selectedSlotId.value}`)
}
</script>

<template>
  <div class="page">
    <h1>Выберите время</h1>

    <div v-if="loading">Загрузка...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <template v-else>
      <CalendarGrid
        :slots="slots"
        :selected-slot-id="selectedSlotId ?? undefined"
        @select-slot="selectedSlotId = $event"
      />

      <div class="actions">
        <button
          class="btn btn-primary"
          :disabled="!selectedSlotId"
          @click="confirmSlot"
        >
          Далее
        </button>
      </div>
    </template>
  </div>
</template>

<style scoped>
.actions {
  margin-top: 1.5rem;
  display: flex;
  justify-content: center;
}
.error {
  color: var(--mantine-color-red-7);
}
</style>
