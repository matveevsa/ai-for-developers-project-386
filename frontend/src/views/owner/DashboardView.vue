<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import dayjs from 'dayjs'
import { useEventTypes } from '@/composables/useEventTypes'
import { listUpcomingBookings, deleteBooking } from '@/api/bookings'
import EventTypeForm from '@/components/owner/EventTypeForm.vue'
import type { BookingEnriched, EventType, EventTypeCreate } from '@/types/api'

const route = useRoute()
const router = useRouter()

const activeTab = computed(() => {
  const tab = route.query.tab
  if (tab === 'bookings') return 'bookings'
  return 'event-types'
})

function switchTab(tab: string) {
  router.push({ query: { tab } })
}

const {
  eventTypes, loading: etLoading, error: etError,
  create, update, remove,
} = useEventTypes()

const showForm = ref(false)
const editTarget = ref<EventType | null>(null)

function startCreate() {
  editTarget.value = null
  showForm.value = true
}

function startEdit(et: EventType) {
  editTarget.value = et
  showForm.value = true
}

function cancelForm() {
  showForm.value = false
  editTarget.value = null
}

async function handleSubmit(data: EventTypeCreate) {
  if (editTarget.value) {
    await update(editTarget.value.id, data)
  } else {
    await create(data)
  }
  showForm.value = false
  editTarget.value = null
}

async function handleDelete(id: string, name: string) {
  if (confirm(`Удалить тип "${name}"? Все бронирования будут удалены.`)) {
    await remove(id)
  }
}

const ownerId = 'owner-1'
const bookings = ref<BookingEnriched[]>([])
const bkLoading = ref(false)
const bkError = ref<string | null>(null)

const bookingsFilter = ref<'upcoming' | 'past'>('upcoming')

async function loadBookings() {
  bkLoading.value = true
  try {
    bookings.value = await listUpcomingBookings(ownerId)
  } catch (e: any) {
    bkError.value = e?.response?.data?.message ?? 'Ошибка загрузки'
  } finally {
    bkLoading.value = false
  }
}

async function handleDeleteBooking(id: string) {
  if (confirm('Удалить бронирование?')) {
    await deleteBooking(id)
    await loadBookings()
  }
}

onMounted(loadBookings)

const now = () => dayjs()

const filteredBookings = computed(() => {
  const filtered = bookings.value.filter(b => {
    const t = dayjs(b.startTime)
    return bookingsFilter.value === 'upcoming' ? t.isAfter(now()) : t.isBefore(now())
  })
  return filtered.sort((a, b) => {
    const diff = dayjs(a.startTime).unix() - dayjs(b.startTime).unix()
    return bookingsFilter.value === 'upcoming' ? diff : -diff
  })
})
</script>

<template>
  <div class="owner-page">
    <h1>Панель владельца</h1>

    <div class="tabs">
      <button
        class="tab"
        :class="{ active: activeTab === 'event-types' }"
        @click="switchTab('event-types')"
      >
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <rect x="3" y="4" width="18" height="18" rx="2" ry="2"/>
          <line x1="16" y1="2" x2="16" y2="6"/>
          <line x1="8" y1="2" x2="8" y2="6"/>
          <line x1="3" y1="10" x2="21" y2="10"/>
        </svg>
        Типы событий
      </button>
      <button
        class="tab"
        :class="{ active: activeTab === 'bookings' }"
        @click="switchTab('bookings')"
      >
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <rect x="3" y="4" width="18" height="18" rx="2" ry="2"/>
          <line x1="16" y1="2" x2="16" y2="6"/>
          <line x1="8" y1="2" x2="8" y2="6"/>
          <line x1="3" y1="10" x2="21" y2="10"/>
        </svg>
        Бронирования
      </button>
    </div>

    <div v-if="activeTab === 'event-types'" class="tab-content">
      <div class="toolbar">
        <button class="btn btn-primary" @click="startCreate">
          + Создать
        </button>
      </div>

      <EventTypeForm
        v-if="showForm"
        :event-type="editTarget"
        @submit="handleSubmit"
        @cancel="cancelForm"
      />

      <div v-if="etLoading" class="state-msg">Загрузка...</div>
      <div v-else-if="etError" class="state-msg error">{{ etError }}</div>
      <div v-else-if="eventTypes.length === 0 && !showForm" class="state-msg empty">
        Пока нет ни одного типа события. Нажмите «Создать», чтобы добавить первый.
      </div>
      <table v-else class="data-table">
        <thead>
          <tr>
            <th>Название</th>
            <th style="width:120px">Длительность</th>
            <th>Описание</th>
            <th style="width:100px">Действия</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="et in eventTypes" :key="et.id">
            <td>{{ et.name }}</td>
            <td>{{ et.duration }} мин</td>
            <td class="muted">{{ et.description ?? '—' }}</td>
            <td>
              <div class="row-actions">
                <button class="btn-icon" title="Редактировать" @click="startEdit(et)">
                  ✎
                </button>
                <button class="btn-icon danger" title="Удалить" @click="handleDelete(et.id, et.name)">
                  ✕
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="activeTab === 'bookings'" class="tab-content">
      <div class="toolbar">
        <div class="filter-group">
          <button
            class="filter-btn"
            :class="{ active: bookingsFilter === 'upcoming' }"
            @click="bookingsFilter = 'upcoming'"
          >Предстоящие</button>
          <button
            class="filter-btn"
            :class="{ active: bookingsFilter === 'past' }"
            @click="bookingsFilter = 'past'"
          >Прошедшие</button>
        </div>
      </div>

      <div v-if="bkLoading" class="state-msg">Загрузка...</div>
      <div v-else-if="bkError" class="state-msg error">{{ bkError }}</div>
      <div v-else-if="bookings.length === 0" class="state-msg empty">
        Нет бронирований
      </div>
      <table v-else class="data-table">
        <thead>
          <tr>
            <th>Дата</th>
            <th>Время</th>
            <th>Тип события</th>
            <th>Гость</th>
            <th>Email</th>
            <th style="width:60px"></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="b in filteredBookings" :key="b.id">
            <td class="nowrap">{{ dayjs(b.startTime).format('D MMM YYYY') }}</td>
            <td class="nowrap">{{ dayjs(b.startTime).format('HH:mm') }} — {{ dayjs(b.endTime).format('HH:mm') }}</td>
            <td>{{ b.eventTypeName }}</td>
            <td>{{ b.guestName }}</td>
            <td class="muted">{{ b.guestEmail }}</td>
            <td>
              <button class="btn-icon danger" title="Удалить" @click="handleDeleteBooking(b.id)">✕</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<style scoped>
.owner-page h1 {
  margin-bottom: 1.5rem;
}

.tabs {
  display: flex;
  gap: 0;
  border-bottom: 2px solid var(--mantine-color-gray-3);
  margin-bottom: 1.5rem;
}

.tab {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.625rem 1.25rem;
  font-size: 0.9375rem;
  font-weight: 500;
  color: var(--mantine-color-gray-6);
  background: none;
  border: none;
  border-bottom: 2px solid transparent;
  margin-bottom: -2px;
  cursor: pointer;
  transition: 0.15s;
}

.tab:hover {
  color: var(--mantine-color-gray-8);
}

.tab.active {
  color: var(--mantine-color-orange-6);
  border-bottom-color: var(--mantine-color-orange-6);
}

.tab-content {
  background: white;
  border: 1px solid var(--mantine-color-gray-3);
  border-radius: var(--mantine-radius-md);
  padding: 1.5rem;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.filter-group {
  display: flex;
  gap: 0.25rem;
  background: var(--mantine-color-gray-1);
  border-radius: var(--mantine-radius-sm);
  padding: 0.25rem;
}

.filter-btn {
  padding: 0.375rem 0.75rem;
  font-size: 0.8125rem;
  font-weight: 500;
  border: none;
  border-radius: var(--mantine-radius-sm);
  cursor: pointer;
  background: transparent;
  color: var(--mantine-color-gray-6);
  transition: 0.15s;
}

.filter-btn.active {
  background: white;
  color: var(--mantine-color-gray-9);
  box-shadow: 0 1px 2px rgba(0,0,0,0.06);
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.875rem;
}

.data-table th {
  text-align: left;
  padding: 0.625rem 0.75rem;
  font-weight: 600;
  color: var(--mantine-color-gray-7);
  border-bottom: 2px solid var(--mantine-color-gray-3);
  white-space: nowrap;
}

.data-table td {
  padding: 0.625rem 0.75rem;
  border-bottom: 1px solid var(--mantine-color-gray-2);
}

.data-table tbody tr:hover {
  background: var(--mantine-color-gray-0);
}

.row-actions {
  display: flex;
  gap: 0.375rem;
}

.btn-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border: 1px solid var(--mantine-color-gray-3);
  border-radius: var(--mantine-radius-sm);
  background: white;
  cursor: pointer;
  font-size: 0.875rem;
  color: var(--mantine-color-gray-6);
  transition: 0.15s;
}

.btn-icon:hover {
  border-color: var(--mantine-color-orange-4);
  color: var(--mantine-color-orange-6);
}

.btn-icon.danger:hover {
  border-color: var(--mantine-color-red-7);
  color: var(--mantine-color-red-7);
}

.state-msg {
  padding: 2rem;
  text-align: center;
  font-size: 0.9375rem;
}

.state-msg.error {
  color: var(--mantine-color-red-7);
}

.state-msg.empty {
  color: var(--mantine-color-gray-5);
}

.muted {
  color: var(--mantine-color-gray-5);
}
.nowrap {
  white-space: nowrap;
}
</style>
