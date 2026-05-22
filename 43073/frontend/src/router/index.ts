import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/store/user'

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
    path: '/',
    component: () => import('@/layout/Layout.vue'),
    redirect: '/dashboard',
    meta: { requiresAuth: true },
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/Dashboard.vue'),
        meta: { title: '首页', icon: 'HomeFilled' }
      },
      {
        path: 'activities',
        name: 'Activities',
        component: () => import('@/views/ActivityList.vue'),
        meta: { title: '活动管理', icon: 'Calendar' }
      },
      {
        path: 'activities/create',
        name: 'ActivityCreate',
        component: () => import('@/views/ActivityForm.vue'),
        meta: { title: '创建活动', icon: 'Plus' }
      },
      {
        path: 'activities/:id/edit',
        name: 'ActivityEdit',
        component: () => import('@/views/ActivityForm.vue'),
        meta: { title: '编辑活动', icon: 'Edit' }
      },
      {
        path: 'activities/:id',
        name: 'ActivityDetail',
        component: () => import('@/views/ActivityDetail.vue'),
        meta: { title: '活动详情', icon: 'View' }
      },
      {
        path: 'ticket-types',
        name: 'TicketTypes',
        component: () => import('@/views/TicketTypeList.vue'),
        meta: { title: '票务管理', icon: 'Ticket' }
      },
      {
        path: 'coupons',
        name: 'Coupons',
        component: () => import('@/views/CouponList.vue'),
        meta: { title: '优惠券管理', icon: 'Discount' }
      },
      {
        path: 'orders',
        name: 'Orders',
        component: () => import('@/views/OrderList.vue'),
        meta: { title: '订单管理', icon: 'List' }
      },
      {
        path: 'checkin',
        name: 'CheckIn',
        component: () => import('@/views/CheckIn.vue'),
        meta: { title: '签到管理', icon: 'Camera' }
      },
      {
        path: 'statistics',
        name: 'Statistics',
        component: () => import('@/views/Statistics.vue'),
        meta: { title: '统计报表', icon: 'DataLine' }
      },
      {
        path: 'users',
        name: 'Users',
        component: () => import('@/views/UserList.vue'),
        meta: { title: '用户管理', icon: 'User', admin: true }
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
  document.title = (to.meta.title as string) || '活动票务管理系统'

  if (to.meta.requiresAuth && !userStore.isLogin) {
    next({ path: '/login', query: { redirect: to.fullPath } })
  } else if (to.meta.admin && !userStore.isAdmin) {
    next('/dashboard')
  } else {
    next()
  }
})

export default router
