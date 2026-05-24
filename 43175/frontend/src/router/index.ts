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
    component: () => import('@/layouts/MainLayout.vue'),
    redirect: '/dashboard',
    meta: { requiresAuth: true },
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/Dashboard.vue'),
        meta: { title: '首页概览' }
      },
      {
        path: 'devices',
        name: 'Devices',
        component: () => import('@/views/Devices.vue'),
        meta: { title: '设备管理' }
      },
      {
        path: 'devices/:id',
        name: 'DeviceDetail',
        component: () => import('@/views/DeviceDetail.vue'),
        meta: { title: '设备详情' }
      },
      {
        path: 'groups',
        name: 'Groups',
        component: () => import('@/views/Groups.vue'),
        meta: { title: '设备分组' }
      },
      {
        path: 'energy',
        name: 'Energy',
        component: () => import('@/views/Energy.vue'),
        meta: { title: '能耗监控' }
      },
      {
        path: 'scenes',
        name: 'Scenes',
        component: () => import('@/views/Scenes.vue'),
        meta: { title: '场景联动' }
      },
      {
        path: 'schedules',
        name: 'Schedules',
        component: () => import('@/views/Schedules.vue'),
        meta: { title: '定时任务' }
      },
      {
        path: 'families',
        name: 'Families',
        component: () => import('@/views/Families.vue'),
        meta: { title: '家庭管理' }
      },
      {
        path: 'notifications',
        name: 'Notifications',
        component: () => import('@/views/Notifications.vue'),
        meta: { title: '通知消息' }
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/Profile.vue'),
        meta: { title: '个人中心' }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const userStore = useUserStore()

  if (to.meta.requiresAuth !== false && !userStore.token) {
    next({ name: 'Login' })
  } else if (to.name === 'Login' && userStore.token) {
    next({ name: 'Dashboard' })
  } else {
    next()
  }
})

export default router
