import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/stores/user'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/tasks'
  },
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
    path: '/tasks',
    name: 'Tasks',
    component: () => import('@/views/task/TaskList.vue'),
    meta: { title: '任务大厅' }
  },
  {
    path: '/tasks/create',
    name: 'CreateTask',
    component: () => import('@/views/task/CreateTask.vue'),
    meta: { title: '发布任务', requiresRole: ['publisher', 'admin'] }
  },
  {
    path: '/tasks/:id',
    name: 'TaskDetail',
    component: () => import('@/views/task/TaskDetail.vue'),
    meta: { title: '任务详情' }
  },
  {
    path: '/my-tasks',
    name: 'MyTasks',
    component: () => import('@/views/task/MyTasks.vue'),
    meta: { title: '我的任务' }
  },
  {
    path: '/orders',
    name: 'Orders',
    component: () => import('@/views/order/OrderList.vue'),
    meta: { title: '我的订单' }
  },
  {
    path: '/orders/:id',
    name: 'OrderDetail',
    component: () => import('@/views/order/OrderDetail.vue'),
    meta: { title: '订单详情' }
  },
  {
    path: '/payments',
    name: 'Payments',
    component: () => import('@/views/payment/PaymentHistory.vue'),
    meta: { title: '交易记录' }
  },
  {
    path: '/wallet',
    name: 'Wallet',
    component: () => import('@/views/payment/Wallet.vue'),
    meta: { title: '我的钱包' }
  },
  {
    path: '/reviews',
    name: 'Reviews',
    component: () => import('@/views/review/ReviewList.vue'),
    meta: { title: '评价中心' }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('@/views/user/Profile.vue'),
    meta: { title: '个人中心' }
  },
  {
    path: '/verification',
    name: 'Verification',
    component: () => import('@/views/user/Verification.vue'),
    meta: { title: '实名认证' }
  },
  {
    path: '/admin',
    name: 'Admin',
    component: () => import('@/views/admin/AdminPanel.vue'),
    meta: { title: '管理后台', requiresRole: ['admin'] }
  },
  {
    path: '/admin/users',
    name: 'UserManagement',
    component: () => import('@/views/admin/UserManagement.vue'),
    meta: { title: '用户管理', requiresRole: ['admin'] }
  },
  {
    path: '/admin/couriers',
    name: 'CourierManagement',
    component: () => import('@/views/admin/CourierManagement.vue'),
    meta: { title: '跑腿员审核', requiresRole: ['admin'] }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, _from, next) => {
  const userStore = useUserStore()
  const requiresAuth = to.meta.requiresAuth !== false

  document.title = `${to.meta.title || '跑腿服务'} - 同城配送平台`

  if (requiresAuth && !userStore.isLoggedIn) {
    next({ name: 'Login', query: { redirect: to.fullPath } })
    return
  }

  if (to.meta.requiresRole && userStore.userInfo) {
    const requiredRoles = to.meta.requiresRole as string[]
    if (!requiredRoles.includes(userStore.userInfo.role)) {
      next({ name: 'Tasks' })
      return
    }
  }

  if ((to.name === 'Login' || to.name === 'Register') && userStore.isLoggedIn) {
    next({ name: 'Tasks' })
    return
  }

  next()
})

export default router
