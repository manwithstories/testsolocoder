import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/store/user'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { requiresAuth: false, title: '登录' }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/Register.vue'),
    meta: { requiresAuth: false, title: '注册' }
  },
  {
    path: '/',
    component: () => import('@/layout/MainLayout.vue'),
    redirect: '/dashboard',
    meta: { requiresAuth: true },
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/Dashboard.vue'),
        meta: { title: '仪表盘', icon: 'DataAnalysis' }
      },
      {
        path: 'weddings',
        name: 'Weddings',
        component: () => import('@/views/WeddingList.vue'),
        meta: { title: '婚礼管理', icon: 'House' }
      },
      {
        path: 'weddings/:id',
        name: 'WeddingDetail',
        component: () => import('@/views/WeddingDetail.vue'),
        meta: { title: '婚礼详情', hidden: true }
      },
      {
        path: 'vendors',
        name: 'Vendors',
        component: () => import('@/views/VendorList.vue'),
        meta: { title: '供应商管理', icon: 'OfficeBuilding' }
      },
      {
        path: 'guests',
        name: 'Guests',
        component: () => import('@/views/GuestList.vue'),
        meta: { title: '嘉宾管理', icon: 'User' }
      },
      {
        path: 'budget',
        name: 'Budget',
        component: () => import('@/views/Budget.vue'),
        meta: { title: '预算管理', icon: 'Money' }
      },
      {
        path: 'tasks',
        name: 'Tasks',
        component: () => import('@/views/TaskList.vue'),
        meta: { title: '任务清单', icon: 'List' }
      },
      {
        path: 'documents',
        name: 'Documents',
        component: () => import('@/views/DocumentList.vue'),
        meta: { title: '合同文档', icon: 'Document' }
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/Profile.vue'),
        meta: { title: '个人中心', icon: 'UserFilled', hidden: true }
      },
      {
        path: 'users',
        name: 'Users',
        component: () => import('@/views/UserList.vue'),
        meta: { title: '用户管理', icon: 'UserGroup', roles: ['admin'] }
      },
      {
        path: 'logs',
        name: 'Logs',
        component: () => import('@/views/OperationLog.vue'),
        meta: { title: '操作日志', icon: 'DocumentCopy', roles: ['admin'] }
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

  document.title = `${to.meta.title || '婚礼策划管理系统'} - 婚礼策划管理系统`

  if (to.meta.requiresAuth !== false && !userStore.isLoggedIn) {
    next('/login')
    return
  }

  if ((to.path === '/login' || to.path === '/register') && userStore.isLoggedIn) {
    next('/')
    return
  }

  if (to.meta.roles && !userStore.hasRole(to.meta.roles as string[])) {
    next('/')
    return
  }

  next()
})

export default router
