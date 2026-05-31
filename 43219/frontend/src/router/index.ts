import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useUserStore } from '../stores/user'

const routes: RouteRecordRaw[] = [
  { path: '/login', name: 'Login', component: () => import('../views/Login.vue') },
  { path: '/register', name: 'Register', component: () => import('../views/Register.vue') },
  { path: '/', redirect: '/dashboard' },
  { path: '/services', name: 'Services', component: () => import('../views/Services.vue') },
  { path: '/staff/:id', name: 'StaffDetail', component: () => import('../views/StaffDetail.vue'), props: true },
  { path: '/dashboard', name: 'Dashboard', component: () => import('../views/Dashboard.vue'), meta: { requiresAuth: true } },
  { path: '/booking/new', name: 'BookingCreate', component: () => import('../views/customer/BookingCreate.vue'), meta: { requiresAuth: true, roles: ['customer'] } },
  { path: '/bookings', name: 'Bookings', component: () => import('../views/BookingList.vue'), meta: { requiresAuth: true } },
  { path: '/orders', name: 'Orders', component: () => import('../views/OrderList.vue'), meta: { requiresAuth: true } },
  { path: '/orders/:id/review', name: 'ReviewWrite', component: () => import('../views/customer/ReviewWrite.vue'), props: true, meta: { requiresAuth: true, roles: ['customer'] } },
  { path: '/tickets', name: 'Tickets', component: () => import('../views/TicketList.vue'), meta: { requiresAuth: true } },
  { path: '/tickets/new', name: 'TicketCreate', component: () => import('../views/customer/TicketCreate.vue'), meta: { requiresAuth: true, roles: ['customer'] } },
  { path: '/company/services', name: 'CompanyServices', component: () => import('../views/company/ServiceManage.vue'), meta: { requiresAuth: true, roles: ['company'] } },
  { path: '/company/bookings', name: 'CompanyBookings', component: () => import('../views/company/BookingReview.vue'), meta: { requiresAuth: true, roles: ['company'] } },
  { path: '/company/finance', name: 'CompanyFinance', component: () => import('../views/company/Finance.vue'), meta: { requiresAuth: true, roles: ['company'] } },
  { path: '/company/stats', name: 'CompanyStats', component: () => import('../views/company/Stats.vue'), meta: { requiresAuth: true, roles: ['company', 'admin'] } },
  { path: '/staff/profile', name: 'StaffProfile', component: () => import('../views/staff/Profile.vue'), meta: { requiresAuth: true, roles: ['staff'] } },
  { path: '/staff/schedule', name: 'StaffSchedule', component: () => import('../views/staff/Schedule.vue'), meta: { requiresAuth: true, roles: ['staff'] } },
  { path: '/staff/orders', name: 'StaffOrders', component: () => import('../views/staff/OrderWork.vue'), meta: { requiresAuth: true, roles: ['staff'] } },
  { path: '/staff/earnings', name: 'StaffEarnings', component: () => import('../views/staff/Earnings.vue'), meta: { requiresAuth: true, roles: ['staff'] } },
  { path: '/admin/stats', name: 'AdminStats', component: () => import('../views/admin/Stats.vue'), meta: { requiresAuth: true, roles: ['admin'] } },
  { path: '/admin/tickets', name: 'AdminTickets', component: () => import('../views/admin/Tickets.vue'), meta: { requiresAuth: true, roles: ['admin'] } },
  { path: '/admin/staff', name: 'AdminStaff', component: () => import('../views/admin/StaffManage.vue'), meta: { requiresAuth: true, roles: ['admin'] } },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to) => {
  const userStore = useUserStore()
  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    return { name: 'Login', query: { redirect: to.fullPath } }
  }
  if (to.meta.roles && Array.isArray(to.meta.roles)) {
    if (!userStore.role || !(to.meta.roles as string[]).includes(userStore.role)) {
      return { name: 'Dashboard' }
    }
  }
  return true
})

export default router
