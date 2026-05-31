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
        name: 'Dashboard',
        component: () => import('@/views/Dashboard.vue'),
        meta: { title: '首页' }
      },
      {
        path: 'ships',
        name: 'ShipList',
        component: () => import('@/views/ships/ShipList.vue'),
        meta: { title: '船舶列表' }
      },
      {
        path: 'ships/:id',
        name: 'ShipDetail',
        component: () => import('@/views/ships/ShipDetail.vue'),
        meta: { title: '船舶详情' }
      },
      {
        path: 'my-ships',
        name: 'MyShips',
        component: () => import('@/views/ships/MyShips.vue'),
        meta: { title: '我的船只', roles: ['owner', 'admin'] }
      },
      {
        path: 'ship-create',
        name: 'ShipCreate',
        component: () => import('@/views/ships/ShipCreate.vue'),
        meta: { title: '发布船只', roles: ['owner', 'admin'] }
      },
      {
        path: 'docks',
        name: 'DockList',
        component: () => import('@/views/docks/DockList.vue'),
        meta: { title: '码头列表' }
      },
      {
        path: 'docks/:id',
        name: 'DockDetail',
        component: () => import('@/views/docks/DockDetail.vue'),
        meta: { title: '码头详情' }
      },
      {
        path: 'berths',
        name: 'BerthList',
        component: () => import('@/views/docks/BerthList.vue'),
        meta: { title: '泊位管理', roles: ['admin'] }
      },
      {
        path: 'berth-reservations',
        name: 'BerthReservations',
        component: () => import('@/views/docks/BerthReservations.vue'),
        meta: { title: '泊位预约' }
      },
      {
        path: 'rentals',
        name: 'RentalList',
        component: () => import('@/views/rentals/RentalList.vue'),
        meta: { title: '租赁订单' }
      },
      {
        path: 'rentals/:id',
        name: 'RentalDetail',
        component: () => import('@/views/rentals/RentalDetail.vue'),
        meta: { title: '订单详情' }
      },
      {
        path: 'rental-create',
        name: 'RentalCreate',
        component: () => import('@/views/rentals/RentalCreate.vue'),
        meta: { title: '创建租赁', roles: ['tenant', 'admin'] }
      },
      {
        path: 'my-rentals',
        name: 'MyRentals',
        component: () => import('@/views/rentals/MyRentals.vue'),
        meta: { title: '我的租赁' }
      },
      {
        path: 'voyage-logs',
        name: 'VoyageLogList',
        component: () => import('@/views/voyage/VoyageLogList.vue'),
        meta: { title: '航海日志' }
      },
      {
        path: 'voyage-logs/:id',
        name: 'VoyageLogDetail',
        component: () => import('@/views/voyage/VoyageLogDetail.vue'),
        meta: { title: '日志详情' }
      },
      {
        path: 'maintenance',
        name: 'MaintenanceList',
        component: () => import('@/views/maintenance/MaintenanceList.vue'),
        meta: { title: '维修保养', roles: ['owner', 'admin'] }
      },
      {
        path: 'maintenance/:id',
        name: 'MaintenanceDetail',
        component: () => import('@/views/maintenance/MaintenanceDetail.vue'),
        meta: { title: '保养详情' }
      },
      {
        path: 'finance',
        name: 'FinanceDashboard',
        component: () => import('@/views/finance/FinanceDashboard.vue'),
        meta: { title: '财务结算' }
      },
      {
        path: 'transactions',
        name: 'TransactionList',
        component: () => import('@/views/finance/TransactionList.vue'),
        meta: { title: '交易记录' }
      },
      {
        path: 'settlements',
        name: 'SettlementList',
        component: () => import('@/views/finance/SettlementList.vue'),
        meta: { title: '结算单' }
      },
      {
        path: 'reviews',
        name: 'ReviewList',
        component: () => import('@/views/reviews/ReviewList.vue'),
        meta: { title: '评价管理' }
      },
      {
        path: 'my-reviews',
        name: 'MyReviews',
        component: () => import('@/views/reviews/MyReviews.vue'),
        meta: { title: '我的评价' }
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/user/Profile.vue'),
        meta: { title: '个人中心' }
      },
      {
        path: 'admin/users',
        name: 'UserManagement',
        component: () => import('@/views/admin/UserManagement.vue'),
        meta: { title: '用户管理', roles: ['admin'] }
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/NotFound.vue'),
    meta: { requiresAuth: false }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const userStore = useUserStore()

  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    next({ name: 'Login', query: { redirect: to.fullPath } })
    return
  }

  if (to.meta.roles && !userStore.hasRole(to.meta.roles as string[])) {
    next({ name: 'Dashboard' })
    return
  }

  if ((to.name === 'Login' || to.name === 'Register') && userStore.isLoggedIn) {
    next({ name: 'Dashboard' })
    return
  }

  next()
})

export default router
