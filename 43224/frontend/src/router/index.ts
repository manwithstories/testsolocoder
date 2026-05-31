import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/stores/user'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/Login.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/auth/Register.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/',
    component: () => import('@/layouts/DefaultLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        redirect: '/dashboard'
      },
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/Dashboard.vue'),
        meta: { title: '仪表盘' }
      },
      {
        path: 'projects',
        name: 'Projects',
        component: () => import('@/views/project/ProjectList.vue'),
        meta: { title: '项目管理' }
      },
      {
        path: 'projects/create',
        name: 'CreateProject',
        component: () => import('@/views/project/ProjectCreate.vue'),
        meta: { title: '创建项目', roles: ['client', 'admin'] }
      },
      {
        path: 'projects/:id',
        name: 'ProjectDetail',
        component: () => import('@/views/project/ProjectDetail.vue'),
        meta: { title: '项目详情' }
      },
      {
        path: 'documents',
        name: 'Documents',
        component: () => import('@/views/document/DocumentList.vue'),
        meta: { title: '文档管理' }
      },
      {
        path: 'translation',
        name: 'Translation',
        component: () => import('@/views/translation/TranslationWorkbench.vue'),
        meta: { title: '翻译工作台', roles: ['translator'] }
      },
      {
        path: 'memory',
        name: 'TranslationMemory',
        component: () => import('@/views/memory/MemoryList.vue'),
        meta: { title: '翻译记忆库' }
      },
      {
        path: 'glossary',
        name: 'Glossary',
        component: () => import('@/views/memory/GlossaryList.vue'),
        meta: { title: '术语库' }
      },
      {
        path: 'review',
        name: 'Review',
        component: () => import('@/views/review/ReviewList.vue'),
        meta: { title: '质量审核' }
      },
      {
        path: 'payments',
        name: 'Payments',
        component: () => import('@/views/payment/PaymentList.vue'),
        meta: { title: '费用管理' }
      },
      {
        path: 'statistics',
        name: 'Statistics',
        component: () => import('@/views/statistics/StatisticsDashboard.vue'),
        meta: { title: '数据统计' }
      },
      {
        path: 'users',
        name: 'Users',
        component: () => import('@/views/user/UserList.vue'),
        meta: { title: '用户管理', roles: ['admin'] }
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/user/Profile.vue'),
        meta: { title: '个人中心' }
      },
      {
        path: 'logs',
        name: 'Logs',
        component: () => import('@/views/system/OperationLogs.vue'),
        meta: { title: '操作日志', roles: ['admin'] }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, _from, next) => {
  const userStore = useUserStore()

  if (to.meta.requiresAuth !== false && !userStore.token) {
    next({ name: 'Login' })
    return
  }

  if (to.meta.roles && Array.isArray(to.meta.roles)) {
    const hasRole = to.meta.roles.includes(userStore.role)
    if (!hasRole) {
      next({ name: 'Dashboard' })
      return
    }
  }

  next()
})

export default router
