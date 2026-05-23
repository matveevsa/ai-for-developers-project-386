<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue'
import { listPublicEventTypes, getAvailability, getAvailableSlots, createPublicBooking } from '@/api/public'
import type { EventType, Slot } from '@/types/api'
import type { DayAvailability } from '@/api/public'
import CalendarGrid from '@/components/public/CalendarGrid.vue'
import dayjs from 'dayjs'

const eventTypes = ref<EventType[]>([])
const availability = ref<DayAvailability[]>([])
const selectedTypeId = ref('')
const selectedDate = ref('')
const slots = ref<Slot[]>([])
const loadingSlots = ref(false)
const selectedSlot = ref<Slot | null>(null)
const bookingForm = ref(false)
const guestName = ref('')
const guestEmail = ref('')
const bookingLoading = ref(false)
const bookingError = ref('')
const bookingSuccess = ref(false)

const availMap = computed(() => {
  const map: Record<string, { available: number; total: number }> = {}
  for (const a of availability.value) {
    map[a.date] = { available: a.available, total: a.total }
  }
  return map
})

onMounted(async () => {
  try {
    const types = await listPublicEventTypes()
    eventTypes.value = types
    if (types.length > 0) {
      selectedTypeId.value = types[0].id
      const avail = await getAvailability(types[0].id)
      availability.value = avail
    }
  } catch (e: any) {
    console.error(e)
  }
})

watch(selectedTypeId, async (id) => {
  if (!id) return
  selectedSlot.value = null
  slots.value = []
  bookingForm.value = false
  bookingSuccess.value = false
  try {
    const avail = await getAvailability(id)
    availability.value = avail
  } catch (e: any) {
    console.error(e)
  }
  if (selectedDate.value) {
    await loadSlots()
  }
})

async function selectDate(date: string) {
  selectedDate.value = date
  selectedSlot.value = null
  slots.value = []
  bookingForm.value = false
  bookingSuccess.value = false
  if (selectedTypeId.value) {
    await loadSlots()
  }
}

async function loadSlots() {
  if (!selectedTypeId.value || !selectedDate.value) return
  selectedSlot.value = null
  bookingForm.value = false
  bookingSuccess.value = false
  loadingSlots.value = true
  try {
    slots.value = await getAvailableSlots(selectedTypeId.value, selectedDate.value)
  } catch (e: any) {
    console.error(e)
  } finally {
    loadingSlots.value = false
  }
}

function selectSlot(slot: Slot) {
  selectedSlot.value = slot
  bookingForm.value = true
  bookingSuccess.value = false
}

async function submitBooking() {
  if (!selectedSlot.value || !selectedTypeId.value) return
  bookingLoading.value = true
  bookingError.value = ''
  try {
    await createPublicBooking({
      slotId: selectedSlot.value.id,
      eventTypeId: selectedTypeId.value,
      guestName: guestName.value,
      guestEmail: guestEmail.value,
    })
    bookingSuccess.value = true
    bookingForm.value = false
    const avail = await getAvailability(selectedTypeId.value)
    availability.value = avail
  } catch (e: any) {
    bookingError.value = e?.response?.data?.message ?? 'Ошибка бронирования'
  } finally {
    bookingLoading.value = false
  }
}
</script>

<template>
  <div class="page">
    <div class="card-center">
      <div class="host-info">
        <div class="avatar-wrap">
          <svg viewBox="0 0 40 50" width="48" height="60">
            <ellipse cx="20" cy="15" rx="14" ry="14" fill="#f5cba7" />
            <circle cx="13" cy="13" r="2.5" fill="#333" />
            <circle cx="27" cy="13" r="2.5" fill="#333" />
            <ellipse cx="20" cy="39" rx="16" ry="11" fill="#2ecc71" />
          </svg>
        </div>
        <div class="host-meta">
          <span class="host-label">John Doe</span>
          <span class="host-role">host</span>
        </div>
      </div>

      <div class="section">
        <label class="section-label">Выберите дату</label>
        <CalendarGrid
          :selected-date="selectedDate"
          :availability="availMap"
          @select-date="selectDate"
        />
      </div>

      <div v-if="selectedDate" class="section">
        <div class="field">
          <label for="event-type">Тип события</label>
          <select id="event-type" v-model="selectedTypeId">
            <option v-for="et in eventTypes" :key="et.id" :value="et.id">
              {{ et.name }} ({{ et.duration }} мин)
            </option>
          </select>
        </div>
      </div>

      <div v-if="loadingSlots" class="loading-text">Загрузка слотов...</div>

      <div v-else-if="selectedTypeId && selectedDate && slots.length > 0" class="section">
        <label class="section-label">Свободное время — {{ dayjs(selectedDate).format('DD.MM.YYYY') }}</label>
        <div class="slots-row">
          <button
            v-for="slot in slots"
            :key="slot.id"
            class="slot-btn"
            :class="{ selected: selectedSlot?.id === slot.id }"
            @click="selectSlot(slot)"
          >
            {{ dayjs(slot.startTime).format('HH:mm') }}—{{ dayjs(slot.endTime).format('HH:mm') }}
          </button>
        </div>
      </div>

      <div v-else-if="selectedTypeId && selectedDate && slots.length === 0 && !loadingSlots" class="empty-text">
        Нет свободных слотов на эту дату
      </div>

      <div v-if="bookingForm && selectedSlot" class="section form-section">
        <div class="selected-slot-info">
          {{ dayjs(selectedSlot.startTime).format('DD.MM.YYYY') }}
          {{ dayjs(selectedSlot.startTime).format('HH:mm') }}—
          {{ dayjs(selectedSlot.endTime).format('HH:mm') }}
        </div>

        <form @submit.prevent="submitBooking">
          <div class="field">
            <label for="name">Имя</label>
            <input id="name" v-model="guestName" required minlength="1" maxlength="100" />
          </div>
          <div class="field">
            <label for="email">Email</label>
            <input id="email" v-model="guestEmail" type="email" required />
          </div>
          <div v-if="bookingError" class="error-text">{{ bookingError }}</div>
          <button type="submit" class="btn btn-primary" :disabled="bookingLoading">
            {{ bookingLoading ? 'Отправка...' : 'Забронировать' }}
          </button>
        </form>
      </div>

      <div v-if="bookingSuccess" class="success-block">
        <h2>Бронирование подтверждено!</h2>
        <p>Спасибо, {{ guestName }}.</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.page {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding-top: 2rem;
}
.card-center {
  background: white;
  border-radius: var(--mantine-radius-md);
  padding: 2rem;
  max-width: 700px;
  width: 100%;
}
.host-info {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-bottom: 1.5rem;
}
.avatar-wrap {
  background: #dbeafe;
  padding: 0.5rem;
  border-radius: var(--mantine-radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
}
.host-meta {
  display: flex;
  flex-direction: column;
}
.host-label {
  font-weight: 700;
  font-size: 1rem;
  color: var(--mantine-color-gray-9);
}
.host-role {
  font-size: 0.8125rem;
  color: var(--mantine-color-gray-6);
  opacity: 0.5;
}
.section {
  margin-top: 1.5rem;
}
.section-label {
  display: block;
  font-weight: 600;
  font-size: 0.9375rem;
  margin-bottom: 0.75rem;
  color: var(--mantine-color-gray-8);
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
.field select {
  width: 100%;
  padding: 0.5rem 0.75rem;
  border: 1px solid var(--mantine-color-gray-4);
  border-radius: var(--mantine-radius-sm);
  font-size: 0.875rem;
}
.field select {
  background: white;
  cursor: pointer;
}
.slots-row {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}
.slot-btn {
  padding: 0.5rem 1rem;
  border: 1px solid var(--mantine-color-orange-3);
  border-radius: var(--mantine-radius-sm);
  background: white;
  color: var(--mantine-color-orange-8);
  cursor: pointer;
  font-size: 0.875rem;
  transition: 0.15s;
}
.slot-btn:hover {
  background: var(--mantine-color-orange-0);
}
.slot-btn.selected {
  background: var(--mantine-color-orange-6);
  color: white;
  border-color: var(--mantine-color-orange-6);
}
.form-section {
  border-top: 1px solid var(--mantine-color-gray-3);
  padding-top: 1.5rem;
  margin-top: 1.5rem;
}
.selected-slot-info {
  font-weight: 600;
  font-size: 1rem;
  margin-bottom: 1rem;
  color: var(--mantine-color-orange-8);
  padding: 0.5rem 0.75rem;
  background: var(--mantine-color-orange-1);
  border-radius: var(--mantine-radius-sm);
}
.loading-text,
.empty-text {
  color: var(--mantine-color-gray-5);
  font-size: 0.875rem;
  margin-top: 1rem;
}
.error-text {
  color: var(--mantine-color-red-7);
  font-size: 0.875rem;
  margin-bottom: 0.5rem;
}
.success-block {
  text-align: center;
  padding: 2rem 0;
}
.success-block h2 {
  margin: 0 0 0.5rem;
}
.success-block p {
  color: var(--mantine-color-gray-6);
}
</style>
