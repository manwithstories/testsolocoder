import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/store/user'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/',
    component: () => import('@/views/Layout.vue'),
    meta: { requiresAuth: true },
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/Dashboard.vue'),
        meta: { title: '数据概览' }
      },
      {
        path: 'properties',
        name: 'Properties',
        component: () => import('@/views/PropertyList.vue'),
        meta: { title: '房源管理' }
      },
      {
        path: 'properties/create',
        name: 'PropertyCreate',
        component: () => import('@/views/PropertyEdit.vue'),
        meta: { title: '发布房源' }
      },
      {
        path: 'properties/:id/edit',
        name: 'PropertyEdit',
        component: () => import('@/views/PropertyEdit.vue'),
        meta: { title: '编辑房源' }
      },
      {
        path: 'properties/:id',
        name: 'PropertyDetail',
        component: () => import('@/views/PropertyDetail.vue'),
        meta: { title: '房源详情' }
      },
      {
        path: 'tenants',
        name: 'Tenants',
        component: () => import('@/views/TenantList.vue'),
        meta: { title: '租户管理' }
      },
      {
        path: 'appointments',
        name: 'Appointments',
        component: () => import('@/views/AppointmentList.vue'),
        meta: { title: '看房预约' }
      },
      {
        path: 'contracts',
        name: 'Contracts',
        component: () => import('@/views/ContractList.vue'),
        meta: { title: '合同管理' }
      },
      {
        path: 'contracts/create',
        name: 'ContractCreate',
        component: () => import('@/views/ContractEdit.vue'),
        meta: { title: '创建合同' }
      },
      {
        path: 'contracts/:id',
        name: 'ContractDetail',
        component: () => import('@/views/ContractDetail.vue'),
        meta: { title: '合同详情' }
      },
      {
        path: 'rent',
        name: 'Rent',
        component: () => import('@/views/RentList.vue'),
        meta: { title: '租金管理' }
      },
      {
        path: 'repairs',
        name: 'Repairs',
        component: () => import('@/views/RepairList.vue'),
        meta: { title: '维修工单' }
      },
      {
        path: 'repairs/create',
        name: 'RepairCreate',
        component: () => import('@/views/RepairEdit.vue'),
        meta: { title: '提交报修' }
      },
      {
        path: 'repairs/:id',
        name: 'RepairDetail',
        component: () => import('@/views/RepairDetail.vue'),
        meta: { title: '工单详情' }
      },
      {
        path: 'fees',
        name: 'Fees',
        component: () => import('@/views/FeeList.vue'),
        meta: { title: '公共费用' }
      },
      {
        path: 'notices',
        name: 'Notices',
        component: () => import('@/views/NoticeList.vue'),
        meta: { title: '公告通知' }
      },
      {
        path: 'notices/create',
        name: 'NoticeCreate',
        component: () => import('@/views/NoticeEdit.vue'),
        meta: { title: '发布公告' }
      },
      {
        path: 'notices/:id',
        name: 'NoticeDetail',
        component: () => import('@/views/NoticeDetail.vue'),
        meta: { title: '公告详情' }
      },
      {
        path: 'users',
        name: 'Users',
        component: () => import('@/views/UserList.vue'),
        meta: { title: '用户管理', roles: ['admin'] }
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/Profile.vue'),
        meta: { title: '个人中心' }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const userStore = useUserStore()

  if (to.meta.requiresAuth !== false && !userStore.isLoggedIn) {
    next('/login')
  } else if (to.meta.requiresAuth === false && userStore.isLoggedIn) {
    next('/')
  } else {
    next()
  }
})

export default router
