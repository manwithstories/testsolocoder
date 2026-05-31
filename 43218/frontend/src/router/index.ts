import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/store/user'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'

NProgress.configure({ showSpinner: false })

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/Home.vue'),
    meta: { title: '首页' }
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { title: '登录', noAuth: true }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/Register.vue'),
    meta: { title: '注册', noAuth: true }
  },
  {
    path: '/products',
    name: 'ProductList',
    component: () => import('@/views/products/ProductList.vue'),
    meta: { title: '商品列表' }
  },
  {
    path: '/products/:id',
    name: 'ProductDetail',
    component: () => import('@/views/products/ProductDetail.vue'),
    meta: { title: '商品详情' }
  },
  {
    path: '/services',
    name: 'ServiceList',
    component: () => import('@/views/repairs/ServiceList.vue'),
    meta: { title: '维修服务' }
  },
  {
    path: '/services/:id',
    name: 'ServiceDetail',
    component: () => import('@/views/repairs/ServiceDetail.vue'),
    meta: { title: '服务详情' }
  },
  {
    path: '/user',
    name: 'UserLayout',
    component: () => import('@/layouts/UserLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        redirect: '/user/profile'
      },
      {
        path: 'profile',
        name: 'UserProfile',
        component: () => import('@/views/user/Profile.vue'),
        meta: { title: '个人中心' }
      },
      {
        path: 'orders',
        name: 'UserOrders',
        component: () => import('@/views/user/Orders.vue'),
        meta: { title: '我的订单' }
      },
      {
        path: 'repair-orders',
        name: 'UserRepairOrders',
        component: () => import('@/views/user/RepairOrders.vue'),
        meta: { title: '维修订单' }
      },
      {
        path: 'favorites',
        name: 'UserFavorites',
        component: () => import('@/views/user/Favorites.vue'),
        meta: { title: '我的收藏' }
      },
      {
        path: 'wallet',
        name: 'UserWallet',
        component: () => import('@/views/user/Wallet.vue'),
        meta: { title: '我的钱包' }
      },
      {
        path: 'messages',
        name: 'UserMessages',
        component: () => import('@/views/user/Messages.vue'),
        meta: { title: '消息中心' }
      },
      {
        path: 'reviews',
        name: 'UserReviews',
        component: () => import('@/views/user/Reviews.vue'),
        meta: { title: '评价管理' }
      }
    ]
  },
  {
    path: '/seller',
    name: 'SellerLayout',
    component: () => import('@/layouts/SellerLayout.vue'),
    meta: { requiresAuth: true, role: 'seller' },
    children: [
      {
        path: '',
        redirect: '/seller/products'
      },
      {
        path: 'products',
        name: 'SellerProducts',
        component: () => import('@/views/seller/Products.vue'),
        meta: { title: '商品管理' }
      },
      {
        path: 'products/create',
        name: 'SellerProductCreate',
        component: () => import('@/views/seller/ProductCreate.vue'),
        meta: { title: '发布商品' }
      },
      {
        path: 'orders',
        name: 'SellerOrders',
        component: () => import('@/views/seller/Orders.vue'),
        meta: { title: '订单管理' }
      },
      {
        path: 'negotiations',
        name: 'SellerNegotiations',
        component: () => import('@/views/seller/Negotiations.vue'),
        meta: { title: '议价管理' }
      },
      {
        path: 'stats',
        name: 'SellerStats',
        component: () => import('@/views/seller/Stats.vue'),
        meta: { title: '数据统计' }
      }
    ]
  },
  {
    path: '/technician',
    name: 'TechnicianLayout',
    component: () => import('@/layouts/TechnicianLayout.vue'),
    meta: { requiresAuth: true, role: 'technician' },
    children: [
      {
        path: '',
        redirect: '/technician/services'
      },
      {
        path: 'services',
        name: 'TechnicianServices',
        component: () => import('@/views/technician/Services.vue'),
        meta: { title: '服务管理' }
      },
      {
        path: 'services/create',
        name: 'TechnicianServiceCreate',
        component: () => import('@/views/technician/ServiceCreate.vue'),
        meta: { title: '发布服务' }
      },
      {
        path: 'orders',
        name: 'TechnicianOrders',
        component: () => import('@/views/technician/Orders.vue'),
        meta: { title: '维修订单' }
      },
      {
        path: 'stats',
        name: 'TechnicianStats',
        component: () => import('@/views/technician/Stats.vue'),
        meta: { title: '数据统计' }
      }
    ]
  },
  {
    path: '/admin',
    name: 'AdminLayout',
    component: () => import('@/layouts/AdminLayout.vue'),
    meta: { requiresAuth: true, role: 'admin' },
    children: [
      {
        path: '',
        redirect: '/admin/dashboard'
      },
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
        path: 'products',
        name: 'AdminProducts',
        component: () => import('@/views/admin/Products.vue'),
        meta: { title: '商品审核' }
      },
      {
        path: 'orders',
        name: 'AdminOrders',
        component: () => import('@/views/admin/Orders.vue'),
        meta: { title: '订单管理' }
      },
      {
        path: 'reports',
        name: 'AdminReports',
        component: () => import('@/views/admin/Reports.vue'),
        meta: { title: '举报处理' }
      },
      {
        path: 'warranties',
        name: 'AdminWarranties',
        component: () => import('@/views/admin/Warranties.vue'),
        meta: { title: '质保管理' }
      },
      {
        path: 'transactions',
        name: 'AdminTransactions',
        component: () => import('@/views/admin/Transactions.vue'),
        meta: { title: '交易记录' }
      },
      {
        path: 'logs',
        name: 'AdminLogs',
        component: () => import('@/views/admin/Logs.vue'),
        meta: { title: '操作日志' }
      }
    ]
  },
  {
    path: '/403',
    name: 'Forbidden',
    component: () => import('@/views/error/403.vue'),
    meta: { title: '无权访问' }
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/error/404.vue'),
    meta: { title: '页面未找到' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  NProgress.start()

  document.title = `${to.meta.title || '二手电子产品交易平台'} - 二手电子产品交易与维修服务平台`

  const userStore = useUserStore()

  if (!userStore.token && !to.meta.noAuth) {
    userStore.initUser()
  }

  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    next({ name: 'Login', query: { redirect: to.fullPath } })
    return
  }

  if (to.meta.role && userStore.userRole !== to.meta.role && userStore.userRole !== 'admin') {
    next({ name: 'Forbidden' })
    return
  }

  next()
})

router.afterEach(() => {
  NProgress.done()
})

export default router
