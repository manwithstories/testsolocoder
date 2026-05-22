import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/stores/user'

const routes: RouteRecordRaw[] = [
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
    path: '/',
    component: () => import('@/layouts/MainLayout.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/Dashboard.vue'),
        meta: { title: '首页', icon: 'HomeFilled' }
      },
      {
        path: 'customers',
        name: 'Customers',
        component: () => import('@/views/customer/CustomerList.vue'),
        meta: { title: '顾客管理', icon: 'UserFilled', roles: ['admin', 'technician'] }
      },
      {
        path: 'customers/:id',
        name: 'CustomerDetail',
        component: () => import('@/views/customer/CustomerDetail.vue'),
        meta: { title: '顾客详情', hidden: true }
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/customer/MyProfile.vue'),
        meta: { title: '我的档案', icon: 'User', roles: ['customer'] }
      },
      {
        path: 'technicians',
        name: 'Technicians',
        component: () => import('@/views/technician/TechnicianList.vue'),
        meta: { title: '技师管理', icon: 'Avatar', roles: ['admin'] }
      },
      {
        path: 'technicians/:id',
        name: 'TechnicianDetail',
        component: () => import('@/views/technician/TechnicianDetail.vue'),
        meta: { title: '技师详情', hidden: true }
      },
      {
        path: 'services',
        name: 'Services',
        component: () => import('@/views/service/ServiceList.vue'),
        meta: { title: '服务项目', icon: 'Service', roles: ['admin'] }
      },
      {
        path: 'services/:id',
        name: 'ServiceDetail',
        component: () => import('@/views/service/ServiceDetail.vue'),
        meta: { title: '服务详情', hidden: true }
      },
      {
        path: 'appointments',
        name: 'Appointments',
        component: () => import('@/views/appointment/AppointmentList.vue'),
        meta: { title: '预约管理', icon: 'Calendar', roles: ['admin', 'technician', 'customer'] }
      },
      {
        path: 'appointments/create',
        name: 'CreateAppointment',
        component: () => import('@/views/appointment/AppointmentCreate.vue'),
        meta: { title: '创建预约', icon: 'Plus', roles: ['admin', 'customer'] }
      },
      {
        path: 'my-appointments',
        name: 'MyAppointments',
        component: () => import('@/views/appointment/MyAppointments.vue'),
        meta: { title: '我的预约', icon: 'Tickets', roles: ['customer'] }
      },
      {
        path: 'payments',
        name: 'Payments',
        component: () => import('@/views/payment/PaymentList.vue'),
        meta: { title: '支付管理', icon: 'Money', roles: ['admin'] }
      },
      {
        path: 'member-cards',
        name: 'MemberCards',
        component: () => import('@/views/payment/MemberCardList.vue'),
        meta: { title: '会员卡', icon: 'CreditCard', roles: ['admin', 'customer'] }
      },
      {
        path: 'products',
        name: 'Products',
        component: () => import('@/views/product/ProductList.vue'),
        meta: { title: '库存管理', icon: 'Goods', roles: ['admin'] }
      },
      {
        path: 'product-records',
        name: 'ProductRecords',
        component: () => import('@/views/product/ProductRecordList.vue'),
        meta: { title: '库存记录', icon: 'Document', roles: ['admin'] }
      },
      {
        path: 'reports',
        name: 'Reports',
        component: () => import('@/views/report/ReportDashboard.vue'),
        meta: { title: '报表分析', icon: 'DataAnalysis', roles: ['admin'] }
      },
      {
        path: 'notifications',
        name: 'Notifications',
        component: () => import('@/views/notification/NotificationList.vue'),
        meta: { title: '消息通知', icon: 'Bell', roles: ['admin', 'technician', 'customer'] }
      },
      {
        path: 'audits',
        name: 'Audits',
        component: () => import('@/views/audit/AuditList.vue'),
        meta: { title: '审计日志', icon: 'Notebook', roles: ['admin'] }
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/NotFound.vue'),
    meta: { title: '页面未找到', requiresAuth: false }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  
  document.title = `${to.meta.title || '美容美发系统'} - 美容美发预约管理系统`

  if (to.meta.requiresAuth !== false && !userStore.isLoggedIn) {
    next('/login')
    return
  }

  if (to.meta.roles && userStore.user && !to.meta.roles.includes(userStore.user.role)) {
    next('/dashboard')
    return
  }

  next()
})

export default router
