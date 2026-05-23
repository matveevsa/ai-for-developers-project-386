<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()
</script>

<template>
  <div class="home">
    <h1>Booking System</h1>
    <div class="actions">
      <button class="btn btn-primary" @click="router.push('/owner')">
        Владелец
      </button>
      <button class="btn btn-secondary" @click="router.push('/book')">
        Записаться
      </button>
    </div>
  </div>
</template>

<style scoped>
.home {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 60vh;
  gap: 2rem;
}
.home h1 {
  font-size: 2rem;
  margin: 0;
}
.actions {
  display: flex;
  gap: 1rem;
}
</style>
