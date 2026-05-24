import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { getToken } from '@/utils/storage'
import type { Role } from '@/types'

// 角色菜单配置
export const roleMenus: Record<Role, { path: string; title: string; icon: string }[]> = {
  manufacturer: [
    { path: '/dashboard', title: '工作台', icon: 'HomeFilled' },
    { path: '/products', title: '产品管理', icon: 'Goods' },
    { path: '/orders', title: '订单管理', icon: 'List' },
    { path: '/deliveries', title: '交付管理', icon: 'Finished' },
    { path: '/productions', title: '生产进度', icon: 'Operation' }
  ],
  designer: [
    { path: '/dashboard', title: '工作台', icon: 'HomeFilled' },
    { path: '/designs', title: '设计项目', icon: 'EditPen' },
    { path: '/deliveries', title: '交付管理', icon: 'Finished' },
    { path: '/reviews', title: '评审记录', icon: 'ChatDotRound' }
  ],
  owner: [
    { path: '/dashboard', title: '工作台', icon: 'HomeFilled' },
    { path: '/users', title: '用户管理', icon: 'User' },
    { path: '/products', title: '产品管理', icon: 'Goods' },
    { path: '/orders', title: '订单管理', icon: 'List' },
    { path: '/designs', title: '设计项目', icon: 'EditPen' },
    { path: '/tickets', title: '工单管理', icon: 'Tickets' },
    { path: '/statistics', title: '数据统计', icon: 'DataAnalysis' }
  ]
}

// 静态路由
const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/Login.vue'),
    meta: { title: '登录', requiresAuth: false }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/login/Register.vue'),
    meta: { title: '注册', requiresAuth: false }
  },
  {
    path: '/',
    component: () => import('@/layout/AdminLayout.vue'),
    redirect: '/dashboard',
    meta: { requiresAuth: true },
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/Dashboard.vue'),
        meta: { title: '工作台', icon: 'HomeFilled' }
      },
      {
        path: 'users',
        name: 'Users',
        component: () => import('@/views/Users.vue'),
        meta: { title: '用户管理', icon: 'User', roles: ['owner'] as Role[] }
      },
      {
        path: 'products',
        name: 'Products',
        component: () => import('@/views/Products.vue'),
        meta: { title: '产品管理', icon: 'Goods', roles: ['manufacturer', 'owner'] as Role[] }
      },
      {
        path: 'orders',
        name: 'Orders',
        component: () => import('@/views/Orders.vue'),
        meta: { title: '订单管理', icon: 'List', roles: ['manufacturer', 'owner'] as Role[] }
      },
      {
        path: 'designs',
        name: 'Designs',
        component: () => import('@/views/design/DesignList.vue'),
        meta: { title: '设计项目', icon: 'EditPen', roles: ['designer', 'owner'] as Role[] }
      },
      {
        path: 'designs/new',
        name: 'DesignNew',
        component: () => import('@/views/design/DesignForm.vue'),
        meta: { title: '新建方案', icon: 'EditPen', roles: ['designer', 'owner'] as Role[], hidden: true }
      },
      {
        path: 'designs/:id',
        name: 'DesignDetail',
        component: () => import('@/views/design/DesignDetail.vue'),
        meta: { title: '方案详情', icon: 'EditPen', roles: ['designer', 'owner'] as Role[], hidden: true }
      },
      {
        path: 'designs/:id/edit',
        name: 'DesignEdit',
        component: () => import('@/views/design/DesignForm.vue'),
        meta: { title: '编辑方案', icon: 'EditPen', roles: ['designer', 'owner'] as Role[], hidden: true }
      },
      {
        path: 'deliveries',
        name: 'Deliveries',
        component: () => import('@/views/delivery/DeliveryList.vue'),
        meta: { title: '交付管理', icon: 'Finished', roles: ['manufacturer', 'designer'] as Role[] }
      },
      {
        path: 'deliveries/new',
        name: 'DeliveryNew',
        component: () => import('@/views/delivery/DeliveryForm.vue'),
        meta: { title: '预约安装', icon: 'Finished', roles: ['manufacturer', 'designer'] as Role[], hidden: true }
      },
      {
        path: 'deliveries/:id/edit',
        name: 'DeliveryEdit',
        component: () => import('@/views/delivery/DeliveryForm.vue'),
        meta: { title: '编辑预约', icon: 'Finished', roles: ['manufacturer', 'designer'] as Role[], hidden: true }
      },
      {
        path: 'reviews',
        name: 'Reviews',
        component: () => import('@/views/review/ReviewList.vue'),
        meta: { title: '评审记录', icon: 'ChatDotRound', roles: ['designer', 'owner'] as Role[] }
      },
      {
        path: 'reviews/new',
        name: 'ReviewNew',
        component: () => import('@/views/review/ReviewForm.vue'),
        meta: { title: '发表评价', icon: 'ChatDotRound', roles: ['designer', 'owner'] as Role[], hidden: true }
      },
      {
        path: 'tickets',
        name: 'Tickets',
        component: () => import('@/views/ticket/TicketList.vue'),
        meta: { title: '工单管理', icon: 'Tickets', roles: ['owner'] as Role[] }
      },
      {
        path: 'tickets/new',
        name: 'TicketNew',
        component: () => import('@/views/ticket/TicketForm.vue'),
        meta: { title: '创建工单', icon: 'Tickets', roles: ['owner'] as Role[], hidden: true }
      },
      {
        path: 'statistics',
        name: 'Statistics',
        component: () => import('@/views/statistics/Statistics.vue'),
        meta: { title: '数据统计', icon: 'DataAnalysis', roles: ['owner'] as Role[] }
      },
      {
        path: 'productions',
        name: 'Productions',
        component: () => import('@/views/production/ProductionList.vue'),
        meta: { title: '生产进度', icon: 'Operation', roles: ['manufacturer'] as Role[] }
      }
    ]
  },
  {
    path: '/404',
    name: 'NotFound',
    component: () => import('@/views/NotFound.vue'),
    meta: { title: '页面不存在', requiresAuth: false }
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/404'
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const token = getToken()
  const userStore = useUserStore()

  // 设置页面标题
  if (to.meta?.title) {
    document.title = `${to.meta.title} - 管理后台`
  }

  // 不需要认证的路由直接放行
  if (to.meta?.requiresAuth === false) {
    if (to.path === '/login' && token) {
      // 已登录用户访问登录页，跳转到首页
      next('/')
    } else {
      next()
    }
    return
  }

  // 未登录，跳转到登录页
  if (!token) {
    next({ path: '/login', query: { redirect: to.fullPath } })
    return
  }

  // 已登录但没有用户信息，尝试恢复
  if (!userStore.userInfo) {
    userStore.restoreUserInfo()
  }

  // 角色权限校验
  if (to.meta?.roles && userStore.userInfo) {
    const allowedRoles = to.meta.roles as Role[]
    if (!allowedRoles.includes(userStore.userInfo.role)) {
      // 无权限，跳转到首页
      next('/')
      return
    }
  }

  next()
})

export default router
