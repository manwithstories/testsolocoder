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
    component: () => import('@/layout/MainLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'Home',
        component: () => import('@/views/dashboard/Home.vue'),
        meta: { title: '首页' }
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/auth/Profile.vue'),
        meta: { title: '个人中心' }
      },
      {
        path: 'drones',
        name: 'Drones',
        component: () => import('@/views/drone/DroneList.vue'),
        meta: { title: '设备列表' }
      },
      {
        path: 'my-drones',
        name: 'MyDrones',
        component: () => import('@/views/drone/MyDrones.vue'),
        meta: { title: '我的设备', roles: ['owner'] }
      },
      {
        path: 'drone/:id',
        name: 'DroneDetail',
        component: () => import('@/views/drone/DroneDetail.vue'),
        meta: { title: '设备详情' }
      },
      {
        path: 'drone/create',
        name: 'CreateDrone',
        component: () => import('@/views/drone/DroneForm.vue'),
        meta: { title: '添加设备', roles: ['owner'] }
      },
      {
        path: 'orders',
        name: 'Orders',
        component: () => import('@/views/order/OrderList.vue'),
        meta: { title: '租赁订单' }
      },
      {
        path: 'order/:id',
        name: 'OrderDetail',
        component: () => import('@/views/order/OrderDetail.vue'),
        meta: { title: '订单详情' }
      },
      {
        path: 'services',
        name: 'Services',
        component: () => import('@/views/service/ServiceList.vue'),
        meta: { title: '航拍服务' }
      },
      {
        path: 'service/create',
        name: 'CreateService',
        component: () => import('@/views/service/ServiceForm.vue'),
        meta: { title: '发布需求', roles: ['client'] }
      },
      {
        path: 'service/:id',
        name: 'ServiceDetail',
        component: () => import('@/views/service/ServiceDetail.vue'),
        meta: { title: '服务详情' }
      },
      {
        path: 'flights',
        name: 'Flights',
        component: () => import('@/views/flight/FlightList.vue'),
        meta: { title: '飞行记录' }
      },
      {
        path: 'flight/create',
        name: 'CreateFlight',
        component: () => import('@/views/flight/FlightForm.vue'),
        meta: { title: '添加飞行记录', roles: ['pilot'] }
      },
      {
        path: 'insurance',
        name: 'Insurance',
        component: () => import('@/views/insurance/ClaimList.vue'),
        meta: { title: '保险理赔' }
      },
      {
        path: 'insurance/create',
        name: 'CreateClaim',
        component: () => import('@/views/insurance/ClaimForm.vue'),
        meta: { title: '申请理赔' }
      },
      {
        path: 'reviews',
        name: 'Reviews',
        component: () => import('@/views/review/ReviewList.vue'),
        meta: { title: '评价中心' }
      },
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/Dashboard.vue'),
        meta: { title: '数据统计', roles: ['owner', 'pilot'] }
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
  userStore.restoreSession()

  if (to.meta.requiresAuth !== false && !userStore.isLoggedIn) {
    next({ name: 'Login', query: { redirect: to.fullPath } })
    return
  }

  const roles = to.meta.roles as string[] | undefined
  if (roles && userStore.role && !roles.includes(userStore.role)) {
    next({ name: 'Home' })
    return
  }

  if ((to.name === 'Login' || to.name === 'Register') && userStore.isLoggedIn) {
    next({ name: 'Home' })
    return
  }

  next()
})

export default router
