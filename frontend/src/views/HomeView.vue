<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()
</script>

<template>
  <div class="home">
    <div class="left">
      <span class="label">Быстрая запись на звонок</span>
      <h1 class="title">Calendar</h1>
      <p class="subtitle">
        Забронируйте встречу за минуту: выберите тип события и удобное время.
      </p>
      <button class="btn btn-primary" @click="router.push('/book')">
        Записаться →
      </button>
    </div>
    <div class="right">
      <h2 class="features-title">Возможности</h2>
      <ul class="features">
        <li>Выбор типа события и удобного времени</li>
        <li>
          Быстрое бронирование с подтверждением
          и дополнительными заметками
        </li>
        <li>Управление типами встреч и просмотр предстоящих записей в админке</li>
      </ul>
    </div>
  </div>
</template>

<style scoped>
.home {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 3rem;
  align-items: center;
  min-height: 60vh;
}

.left {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  align-items: flex-start;
}

.label {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  background: rgba(255, 255, 255, 0.6);
  border: 1px solid var(--mantine-color-gray-3);
  border-radius: var(--mantine-radius-xl);
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--mantine-color-gray-6);
}

.title {
  font-size: 2.5rem;
  font-weight: 800;
  margin: 0;
  color: var(--mantine-color-gray-9);
}

.subtitle {
  font-size: 1.125rem;
  line-height: 1.6;
  color: var(--mantine-color-gray-6);
  margin: 0;
}

.right {
  background: rgba(255, 255, 255, 0.6);
  border: 1px solid var(--mantine-color-gray-3);
  border-radius: var(--mantine-radius-md);
  padding: 1.5rem 2rem;
}

.features-title {
  font-size: 1rem;
  font-weight: 700;
  margin: 0 0 1rem;
  color: var(--mantine-color-gray-9);
}

.features {
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.features li {
  position: relative;
  padding-left: 1.5rem;
  font-size: 0.9375rem;
  line-height: 1.5;
  color: var(--mantine-color-gray-7);
}

.features li::before {
  content: '•';
  position: absolute;
  left: 0;
  color: var(--mantine-color-orange-6);
  font-weight: 700;
}
</style>
