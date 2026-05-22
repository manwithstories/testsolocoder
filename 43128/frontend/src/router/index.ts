import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  { path: '/login', name: 'Login', component: () => import('@/views/Login.vue'), meta: { guest: true } },
  { path: '/register', name: 'Register', component: () => import('@/views/Register.vue'), meta: { guest: true } },
  { path: '/', name: 'Home', component: () => import('@/views/Home.vue') },
  { path: '/events', name: 'EventList', component: () => import('@/views/EventList.vue') },
  { path: '/events/:id', name: 'EventDetail', component: () => import('@/views/EventDetail.vue') },
  { path: '/my/registrations', name: 'MyRegistrations', component: () => import('@/views/MyRegistrations.vue'), meta: { auth: true } },
  { path: '/my/scores', name: 'MyScores', component: () => import('@/views/MyScores.vue'), meta: { auth: true } },
  { path: '/my/certificates', name: 'MyCertificates', component: () => import('@/views/MyCertificates.vue'), meta: { auth: true } },
  { path: '/messages', name: 'Messages', component: () => import('@/views/Messages.vue'), meta: { auth: true } },
  { path: '/admin/events', name: 'AdminEvents', component: () => import('@/views/admin/AdminEvents.vue'), meta: { auth: true, admin: true } },
  { path: '/admin/event/create', name: 'AdminEventCreate', component: () => import('@/views/admin/AdminEventCreate.vue'), meta: { auth: true, admin: true } },
  { path: '/admin/event/:id/edit', name: 'AdminEventEdit', component: () => import('@/views/admin/AdminEventCreate.vue'), meta: { auth: true, admin: true } },
  { path: '/admin/scores/entry', name: 'AdminScoreEntry', component: () => import('@/views/admin/AdminScoreEntry.vue'), meta: { auth: true, admin: true } },
  { path: '/admin/stats', name: 'AdminStats', component: () => import('@/views/admin/AdminStats.vue'), meta: { auth: true, admin: true } },
  { path: '/:pathMatch(.*)*', redirect: '/' },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to, _from, next) => {
  const token = localStorage.getItem('token')
  if (to.meta.auth && !token) {
    next({ path: '/login', query: { redirect: to.fullPath } })
  } else {
    next()
  }
})

export default router
