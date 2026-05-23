<script setup lang="ts">
</script>

<template>
  <div id="app-shell">
    <header class="header">
      <div class="header-inner">
        <router-link to="/" class="logo">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="var(--mantine-color-orange-6)" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <rect x="3" y="4" width="18" height="18" rx="2" ry="2"/>
            <line x1="16" y1="2" x2="16" y2="6"/>
            <line x1="8" y1="2" x2="8" y2="6"/>
            <line x1="3" y1="10" x2="21" y2="10"/>
          </svg>
          Calendar
        </router-link>
        <nav class="nav">
          <router-link to="/book">Записаться</router-link>
          <router-link to="/owner">Админка</router-link>
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
  --mantine-color-orange-0: #fff3e0;
  --mantine-color-orange-1: #ffe0b2;
  --mantine-color-orange-2: #ffcc80;
  --mantine-color-orange-3: #ffb74d;
  --mantine-color-orange-4: #ffa726;
  --mantine-color-orange-6: #fb8c00;
  --mantine-color-orange-7: #f57c00;
  --mantine-color-orange-8: #ef6c00;
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
  background: linear-gradient(135deg, #dbeafe 0%, #93c5fd 20%, #fff3e0 60%, #f5f5f5 100%);
  min-height: 100vh;
}

.header {
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(8px);
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
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
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  font-weight: 700;
  font-size: 1.125rem;
  text-decoration: none;
  color: var(--mantine-color-gray-9);
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
  color: var(--mantine-color-orange-6);
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
  background: var(--mantine-color-orange-6);
  color: white;
  border-color: var(--mantine-color-orange-6);
}

.btn-primary:hover {
  background: var(--mantine-color-orange-7);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-secondary {
  background: white;
  color: var(--mantine-color-orange-6);
  border-color: var(--mantine-color-orange-6);
}

.btn-secondary:hover {
  background: var(--mantine-color-orange-0);
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
