import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/stores/user'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { requiresAuth: false, title: '登录' },
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/Register.vue'),
    meta: { requiresAuth: false, title: '注册' },
  },
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/Home.vue'),
    meta: { requiresAuth: false, title: '首页' },
  },
  {
    path: '/items',
    name: 'Items',
    component: () => import('@/views/ItemList.vue'),
    meta: { requiresAuth: false, title: '拍卖品列表' },
  },
  {
    path: '/items/:id',
    name: 'ItemDetail',
    component: () => import('@/views/ItemDetail.vue'),
    meta: { requiresAuth: false, title: '拍卖品详情' },
  },
  {
    path: '/sessions',
    name: 'Sessions',
    component: () => import('@/views/SessionList.vue'),
    meta: { requiresAuth: false, title: '拍卖会列表' },
  },
  {
    path: '/sessions/:id',
    name: 'SessionDetail',
    component: () => import('@/views/SessionDetail.vue'),
    meta: { requiresAuth: false, title: '拍卖会详情' },
  },
  {
    path: '/user',
    name: 'UserLayout',
    component: () => import('@/views/user/UserLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        redirect: '/user/profile',
      },
      {
        path: 'profile',
        name: 'UserProfile',
        component: () => import('@/views/user/Profile.vue'),
        meta: { title: '个人中心' },
      },
      {
        path: 'items',
        name: 'UserItems',
        component: () => import('@/views/user/MyItems.vue'),
        meta: { title: '我的拍卖品' },
      },
      {
        path: 'bids',
        name: 'UserBids',
        component: () => import('@/views/user/MyBids.vue'),
        meta: { title: '我的出价' },
      },
      {
        path: 'auto-bids',
        name: 'UserAutoBids',
        component: () => import('@/views/user/MyAutoBids.vue'),
        meta: { title: '自动出价' },
      },
      {
        path: 'orders',
        name: 'UserOrders',
        component: () => import('@/views/user/MyOrders.vue'),
        meta: { title: '我的订单' },
      },
      {
        path: 'notifications',
        name: 'UserNotifications',
        component: () => import('@/views/user/Notifications.vue'),
        meta: { title: '消息中心' },
      },
    ],
  },
  {
    path: '/admin',
    name: 'AdminLayout',
    component: () => import('@/views/admin/AdminLayout.vue'),
    meta: { requiresAuth: true, requiresAdmin: true },
    children: [
      {
        path: '',
        redirect: '/admin/dashboard',
      },
      {
        path: 'dashboard',
        name: 'AdminDashboard',
        component: () => import('@/views/admin/Dashboard.vue'),
        meta: { title: '数据统计' },
      },
      {
        path: 'users',
        name: 'AdminUsers',
        component: () => import('@/views/admin/UserManagement.vue'),
        meta: { title: '用户管理' },
      },
      {
        path: 'sessions',
        name: 'AdminSessions',
        component: () => import('@/views/admin/SessionManagement.vue'),
        meta: { title: '拍卖会管理' },
      },
      {
        path: 'orders',
        name: 'AdminOrders',
        component: () => import('@/views/admin/OrderManagement.vue'),
        meta: { title: '订单管理' },
      },
    ],
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/NotFound.vue'),
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  userStore.initFromStorage()

  document.title = (to.meta.title as string) || '在线拍卖系统'

  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    next({ path: '/login', query: { redirect: to.fullPath } })
    return
  }

  if (to.meta.requiresAdmin && !userStore.isAdmin) {
    next('/')
    return
  }

  next()
})

export default router
