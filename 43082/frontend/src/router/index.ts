import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/stores/user'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/Register.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/',
    component: () => import('@/views/Layout.vue'),
    meta: { requiresAuth: true },
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/Dashboard.vue'),
        meta: { title: '数据概览' }
      },
      {
        path: 'members',
        name: 'Members',
        component: () => import('@/views/member/MemberList.vue'),
        meta: { title: '会员管理' }
      },
      {
        path: 'members/:id',
        name: 'MemberDetail',
        component: () => import('@/views/member/MemberDetail.vue'),
        meta: { title: '会员详情' }
      },
      {
        path: 'courses',
        name: 'Courses',
        component: () => import('@/views/course/CourseList.vue'),
        meta: { title: '课程管理' }
      },
      {
        path: 'courses/:id',
        name: 'CourseDetail',
        component: () => import('@/views/course/CourseDetail.vue'),
        meta: { title: '课程详情' }
      },
      {
        path: 'schedules',
        name: 'Schedules',
        component: () => import('@/views/course/ScheduleList.vue'),
        meta: { title: '课程排期' }
      },
      {
        path: 'bookings',
        name: 'Bookings',
        component: () => import('@/views/booking/BookingList.vue'),
        meta: { title: '预约管理' }
      },
      {
        path: 'my-bookings',
        name: 'MyBookings',
        component: () => import('@/views/booking/MyBookings.vue'),
        meta: { title: '我的预约' }
      },
      {
        path: 'check-ins',
        name: 'CheckIns',
        component: () => import('@/views/checkin/CheckInList.vue'),
        meta: { title: '签到记录' }
      },
      {
        path: 'stats',
        name: 'Stats',
        component: () => import('@/views/stats/StatsCenter.vue'),
        meta: { title: '数据统计' }
      },
      {
        path: 'coaches',
        name: 'Coaches',
        component: () => import('@/views/course/CoachList.vue'),
        meta: { title: '教练管理' }
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/NotFound.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  
  if (to.meta.requiresAuth && !userStore.isLoggedIn()) {
    next('/login')
  } else {
    next()
  }
})

export default router
