import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/home'
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { title: '登录' }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/Register.vue'),
    meta: { title: '注册' }
  },
  {
    path: '/home',
    name: 'Home',
    component: () => import('@/views/Home.vue'),
    meta: { title: '首页' }
  },
  {
    path: '/orders',
    name: 'Orders',
    component: () => import('@/views/Orders.vue'),
    meta: { title: '我的工单', requiresAuth: true }
  },
  {
    path: '/orders/create',
    name: 'CreateOrder',
    component: () => import('@/views/CreateOrder.vue'),
    meta: { title: '创建工单', requiresAuth: true, roles: ['customer'] }
  },
  {
    path: '/orders/:id',
    name: 'OrderDetail',
    component: () => import('@/views/OrderDetail.vue'),
    meta: { title: '工单详情', requiresAuth: true }
  },
  {
    path: '/technicians',
    name: 'Technicians',
    component: () => import('@/views/Technicians.vue'),
    meta: { title: '技师列表' }
  },
  {
    path: '/technicians/:id',
    name: 'TechnicianDetail',
    component: () => import('@/views/TechnicianDetail.vue'),
    meta: { title: '技师详情' }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('@/views/Profile.vue'),
    meta: { title: '个人中心', requiresAuth: true }
  },
  {
    path: '/finance',
    name: 'Finance',
    component: () => import('@/views/Finance.vue'),
    meta: { title: '财务管理', requiresAuth: true }
  },
  {
    path: '/parts',
    name: 'Parts',
    component: () => import('@/views/Parts.vue'),
    meta: { title: '配件管理', requiresAuth: true }
  },
  {
    path: '/part-requests',
    name: 'PartRequests',
    component: () => import('@/views/PartRequests.vue'),
    meta: { title: '配件申请', requiresAuth: true }
  },
  {
    path: '/admin',
    name: 'Admin',
    component: () => import('@/views/admin/AdminLayout.vue'),
    meta: { requiresAuth: true, roles: ['admin'] },
    redirect: '/admin/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'AdminDashboard',
        component: () => import('@/views/admin/Dashboard.vue'),
        meta: { title: '数据概览' }
      },
      {
        path: 'users',
        name: 'AdminUsers',
        component: () => import('@/views/admin/Users.vue'),
        meta: { title: '用户管理' }
      },
      {
        path: 'technicians/verify',
        name: 'TechnicianVerify',
        component: () => import('@/views/admin/TechnicianVerify.vue'),
        meta: { title: '技师审核' }
      },
      {
        path: 'orders',
        name: 'AdminOrders',
        component: () => import('@/views/admin/Orders.vue'),
        meta: { title: '工单管理' }
      },
      {
        path: 'refunds',
        name: 'AdminRefunds',
        component: () => import('@/views/admin/Refunds.vue'),
        meta: { title: '退款审核' }
      },
      {
        path: 'parts',
        name: 'AdminParts',
        component: () => import('@/views/admin/Parts.vue'),
        meta: { title: '配件管理' }
      },
      {
        path: 'part-requests',
        name: 'AdminPartRequests',
        component: () => import('@/views/admin/PartRequests.vue'),
        meta: { title: '配件申请审核' }
      },
      {
        path: 'withdraws',
        name: 'AdminWithdraws',
        component: () => import('@/views/admin/Withdraws.vue'),
        meta: { title: '提现审核' }
      },
      {
        path: 'reports',
        name: 'AdminReports',
        component: () => import('@/views/admin/Reports.vue'),
        meta: { title: '财务报表' }
      },
      {
        path: 'reviews',
        name: 'AdminReviews',
        component: () => import('@/views/admin/Reviews.vue'),
        meta: { title: '差评处理' }
      },
      {
        path: 'categories',
        name: 'AdminCategories',
        component: () => import('@/views/admin/Categories.vue'),
        meta: { title: '分类管理' }
      },
      {
        path: 'service-items',
        name: 'AdminServiceItems',
        component: () => import('@/views/admin/ServiceItems.vue'),
        meta: { title: '服务项目管理' }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  document.title = to.meta.title ? `${to.meta.title} - 维修服务平台` : '维修服务平台'

  const token = localStorage.getItem('token')

  if (to.meta.requiresAuth && !token) {
    next({ path: '/login', query: { redirect: to.fullPath } })
    return
  }

  if (to.meta.roles && token) {
    const userInfo = JSON.parse(localStorage.getItem('userInfo') || '{}')
    const roles = to.meta.roles as string[]
    if (!roles.includes(userInfo.role)) {
      next('/home')
      return
    }
  }

  next()
})

export default router
