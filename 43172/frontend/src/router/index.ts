import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/stores/user'

const routes: RouteRecordRaw[] = [
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
    path: '/register/authenticator',
    name: 'RegisterAuthenticator',
    component: () => import('@/views/RegisterAuthenticator.vue'),
    meta: { title: '鉴定师注册' }
  },
  {
    path: '/',
    component: () => import('@/layouts/DefaultLayout.vue'),
    children: [
      {
        path: '',
        name: 'Home',
        component: () => import('@/views/Home.vue'),
        meta: { title: '首页' }
      },
      {
        path: 'products',
        name: 'Products',
        component: () => import('@/views/ProductList.vue'),
        meta: { title: '商品列表' }
      },
      {
        path: 'products/:id',
        name: 'ProductDetail',
        component: () => import('@/views/ProductDetail.vue'),
        meta: { title: '商品详情' }
      },
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/Dashboard.vue'),
        meta: { title: '数据仪表盘', requiresAuth: true, roles: ['admin'] }
      },
      {
        path: 'my-orders',
        name: 'MyOrders',
        component: () => import('@/views/MyOrders.vue'),
        meta: { title: '我的订单', requiresAuth: true }
      },
      {
        path: 'seller/products',
        name: 'SellerProducts',
        component: () => import('@/views/SellerProducts.vue'),
        meta: { title: '我的商品', requiresAuth: true, roles: ['seller', 'admin'] }
      },
      {
        path: 'seller/products/create',
        name: 'CreateProduct',
        component: () => import('@/views/CreateProduct.vue'),
        meta: { title: '发布商品', requiresAuth: true, roles: ['seller', 'admin'] }
      },
      {
        path: 'seller/products/:id/edit',
        name: 'EditProduct',
        component: () => import('@/views/EditProduct.vue'),
        meta: { title: '编辑商品', requiresAuth: true, roles: ['seller', 'admin'] }
      },
      {
        path: 'authenticator/tasks',
        name: 'AuthenticatorTasks',
        component: () => import('@/views/AuthenticatorTasks.vue'),
        meta: { title: '鉴定任务', requiresAuth: true, roles: ['authenticator', 'admin'] }
      },
      {
        path: 'authentications/:id',
        name: 'AuthenticationDetail',
        component: () => import('@/views/AuthenticationDetail.vue'),
        meta: { title: '鉴定详情', requiresAuth: true }
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/Profile.vue'),
        meta: { title: '个人中心', requiresAuth: true }
      },
      {
        path: 'admin/users',
        name: 'AdminUsers',
        component: () => import('@/views/AdminUsers.vue'),
        meta: { title: '用户管理', requiresAuth: true, roles: ['admin'] }
      },
      {
        path: 'admin/authenticators',
        name: 'AdminAuthenticators',
        component: () => import('@/views/AdminAuthenticators.vue'),
        meta: { title: '鉴定师审核', requiresAuth: true, roles: ['admin'] }
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/NotFound.vue'),
    meta: { title: '404' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  document.title = `${to.meta.title || '奢侈品交易平台'} - 奢侈品二手交易与鉴定平台`

  const userStore = useUserStore()

  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    next('/login')
    return
  }

  if (to.meta.roles && userStore.isLoggedIn) {
    const roles = to.meta.roles as string[]
    if (!roles.includes(userStore.userRole)) {
      next('/')
      return
    }
  }

  next()
})

export default router
