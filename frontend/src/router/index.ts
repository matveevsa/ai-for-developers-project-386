import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('@/views/HomeView.vue'),
    },
    {
      path: '/owner',
      name: 'owner-dashboard',
      component: () => import('@/views/owner/DashboardView.vue'),
    },
    {
      path: '/owner/event-types',
      redirect: () => ({ path: '/owner', query: { tab: 'event-types' } }),
    },
    {
      path: '/owner/bookings',
      redirect: () => ({ path: '/owner', query: { tab: 'bookings' } }),
    },
    {
      path: '/book',
      name: 'public-event-types',
      component: () => import('@/views/public/ChooseEventType.vue'),
    },
  ],
})

export default router
