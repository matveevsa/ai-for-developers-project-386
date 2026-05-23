<script setup lang="ts">
</script>

<template>
  <div id="app-shell">
    <header class="header">
      <div class="header-inner">
        <router-link to="/" class="logo">Booking System</router-link>
        <nav class="nav">
          <router-link to="/owner">Владелец</router-link>
          <router-link to="/book">Записаться</router-link>
        </nav>
      </div>
    </header>
    <main class="main">
      <div class="container">
        <router-view />
      </div>
    </main>
  </div>
</template>

<style>
:root {
  --mantine-color-gray-0: #f8f9fa;
  --mantine-color-gray-1: #f1f3f5;
  --mantine-color-gray-2: #e9ecef;
  --mantine-color-gray-3: #dee2e6;
  --mantine-color-gray-4: #ced4da;
  --mantine-color-gray-5: #adb5bd;
  --mantine-color-gray-6: #868e96;
  --mantine-color-gray-7: #495057;
  --mantine-color-gray-8: #343a40;
  --mantine-color-gray-9: #212529;
  --mantine-color-blue-0: #e7f5ff;
  --mantine-color-blue-1: #d0ebff;
  --mantine-color-blue-2: #a5d8ff;
  --mantine-color-blue-3: #74c0fc;
  --mantine-color-blue-4: #4dabf7;
  --mantine-color-blue-5: #339af0;
  --mantine-color-blue-6: #228be6;
  --mantine-color-blue-7: #1c7ed6;
  --mantine-color-blue-8: #1971c2;
  --mantine-color-blue-9: #1864ab;
  --mantine-color-green-1: #d3f9d8;
  --mantine-color-green-8: #2f9e44;
  --mantine-color-red-7: #f03e3e;
  --mantine-radius-sm: 4px;
  --mantine-radius-md: 8px;
  --mantine-radius-xl: 32px;
}

* {
  box-sizing: border-box;
}

body {
  margin: 0;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  color: var(--mantine-color-gray-9);
  background: var(--mantine-color-gray-0);
}

.header {
  background: white;
  border-bottom: 1px solid var(--mantine-color-gray-3);
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-inner {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0.75rem 1rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.logo {
  font-weight: 700;
  font-size: 1.125rem;
  text-decoration: none;
  color: var(--mantine-color-blue-7);
}

.nav {
  display: flex;
  gap: 1rem;
}

.nav a {
  text-decoration: none;
  color: var(--mantine-color-gray-7);
  font-size: 0.875rem;
}

.nav a:hover {
  color: var(--mantine-color-blue-7);
}

.main {
  padding: 1.5rem 0;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 1rem;
}

.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0.5rem 1rem;
  border-radius: var(--mantine-radius-sm);
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  border: 1px solid transparent;
  transition: 0.15s;
  text-decoration: none;
}

.btn-primary {
  background: var(--mantine-color-blue-6);
  color: white;
  border-color: var(--mantine-color-blue-6);
}

.btn-primary:hover {
  background: var(--mantine-color-blue-7);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-secondary {
  background: white;
  color: var(--mantine-color-blue-6);
  border-color: var(--mantine-color-blue-4);
}

.btn-secondary:hover {
  background: var(--mantine-color-blue-0);
}

.btn-sm {
  padding: 0.375rem 0.75rem;
  font-size: 0.8125rem;
}

.page {
  max-width: 800px;
  margin: 0 auto;
}

.page h1 {
  font-size: 1.5rem;
  margin: 0 0 1.5rem;
}
</style>
