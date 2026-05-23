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
      name: 'owner-event-types',
      component: () => import('@/views/owner/EventTypesView.vue'),
    },
    {
      path: '/owner/bookings',
      name: 'owner-bookings',
      component: () => import('@/views/owner/BookingsView.vue'),
    },
    {
      path: '/book',
      name: 'public-event-types',
      component: () => import('@/views/public/ChooseEventType.vue'),
    },
    {
      path: '/book/:eventTypeId',
      name: 'public-pick-slot',
      component: () => import('@/views/public/PickSlot.vue'),
    },
    {
      path: '/book/:eventTypeId/confirm',
      name: 'public-confirm-booking',
      component: () => import('@/views/public/ConfirmBooking.vue'),
    },
  ],
})

export default router
