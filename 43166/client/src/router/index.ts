import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/store/user'
import { UserRoles, type UserRole } from '@/types'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/Login.vue'),
    meta: { title: '登录', requiresAuth: false }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/auth/Register.vue'),
    meta: { title: '注册', requiresAuth: false }
  },
  {
    path: '/',
    redirect: '/dashboard'
  },
  {
    path: '/',
    component: () => import('@/layouts/DefaultLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/Dashboard.vue'),
        meta: { title: '工作台', icon: 'HomeFilled' }
      },
      {
        path: 'applications',
        name: 'Applications',
        component: () => import('@/views/application/ApplicationList.vue'),
        meta: { title: '注册申请', icon: 'Document' }
      },
      {
        path: 'applications/new',
        name: 'NewApplication',
        component: () => import('@/views/application/NewApplication.vue'),
        meta: { title: '新建申请', icon: 'Plus', hidden: true }
      },
      {
        path: 'applications/:id',
        name: 'ApplicationDetail',
        component: () => import('@/views/application/ApplicationDetail.vue'),
        meta: { title: '申请详情', icon: 'View', hidden: true }
      },
      {
        path: 'fees',
        name: 'Fees',
        component: () => import('@/views/fee/FeeList.vue'),
        meta: { title: '费用管理', icon: 'Money' }
      },
      {
        path: 'notifications',
        name: 'Notifications',
        component: () => import('@/views/notification/NotificationList.vue'),
        meta: { title: '消息中心', icon: 'Bell' }
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/user/Profile.vue'),
        meta: { title: '个人中心', icon: 'User', hidden: true }
      },
      {
        path: 'admin/agents',
        name: 'AgentManagement',
        component: () => import('@/views/admin/AgentManagement.vue'),
        meta: { title: '专员管理', icon: 'UserFilled', roles: [UserRoles.ADMIN] }
      },
      {
        path: 'admin/fee-standards',
        name: 'FeeStandards',
        component: () => import('@/views/admin/FeeStandards.vue'),
        meta: { title: '费用标准', icon: 'Setting', roles: [UserRoles.ADMIN] }
      },
      {
        path: 'admin/discounts',
        name: 'DiscountPolicies',
        component: () => import('@/views/admin/DiscountPolicies.vue'),
        meta: { title: '优惠策略', icon: 'Discount', roles: [UserRoles.ADMIN] }
      },
      {
        path: 'admin/notification-templates',
        name: 'NotificationTemplates',
        component: () => import('@/views/admin/NotificationTemplates.vue'),
        meta: { title: '通知模板', icon: 'MessageBox', roles: [UserRoles.ADMIN] }
      },
      {
        path: 'admin/statistics',
        name: 'Statistics',
        component: () => import('@/views/admin/Statistics.vue'),
        meta: { title: '统计分析', icon: 'DataLine', roles: [UserRoles.ADMIN] }
      },
      {
        path: 'admin/exports',
        name: 'DataExport',
        component: () => import('@/views/admin/DataExport.vue'),
        meta: { title: '数据导出', icon: 'Download', roles: [UserRoles.ADMIN] }
      }
    ]
  },
  {
    path: '/403',
    name: 'Forbidden',
    component: () => import('@/views/error/403.vue'),
    meta: { title: '无权限', requiresAuth: false }
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/error/404.vue'),
    meta: { title: '页面未找到', requiresAuth: false }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const userStore = useUserStore()

  document.title = `${to.meta.title || '企业工商注册平台'} - 代办服务平台`

  if (to.meta.requiresAuth !== false && !userStore.isLoggedIn) {
    next({ name: 'Login', query: { redirect: to.fullPath } })
    return
  }

  if (to.meta.roles && userStore.userRole) {
    const allowedRoles = to.meta.roles as UserRole[]
    if (!allowedRoles.includes(userStore.userRole)) {
      next({ name: 'Forbidden' })
      return
    }
  }

  next()
})

router.afterEach(() => {
})

export default router
