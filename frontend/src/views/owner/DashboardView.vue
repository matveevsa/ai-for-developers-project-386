<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()
</script>

<template>
  <div class="dashboard">
    <h1>Панель владельца</h1>
    <div class="grid">
      <button class="card-btn" @click="router.push('/owner/event-types')">
        <strong>Типы событий</strong>
        <span class="muted">Создание и управление</span>
      </button>
      <button class="card-btn" @click="router.push('/owner/bookings')">
        <strong>Предстоящие встречи</strong>
        <span class="muted">Все бронирования</span>
      </button>
    </div>
  </div>
</template>

<style scoped>
.dashboard h1 {
  margin-bottom: 1.5rem;
}
.grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 1rem;
}
.card-btn {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  padding: 1.5rem;
  border: 1px solid var(--mantine-color-gray-3);
  border-radius: var(--mantine-radius-md);
  background: white;
  cursor: pointer;
  text-align: left;
  font-size: 1rem;
  transition: 0.15s;
}
.card-btn:hover {
  border-color: var(--mantine-color-blue-4);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}
.muted {
  font-size: 0.875rem;
  color: var(--mantine-color-gray-6);
}
</style>
