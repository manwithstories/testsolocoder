import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/stores/user'

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
    meta: { title: '登录', guestOnly: true }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/Register.vue'),
    meta: { title: '注册', guestOnly: true }
  },
  {
    path: '/equipments',
    name: 'EquipmentList',
    component: () => import('@/views/EquipmentList.vue'),
    meta: { title: '设备列表' }
  },
  {
    path: '/equipments/:id',
    name: 'EquipmentDetail',
    component: () => import('@/views/EquipmentDetail.vue'),
    meta: { title: '设备详情' }
  },
  {
    path: '/my-equipments',
    name: 'MyEquipments',
    component: () => import('@/views/MyEquipments.vue'),
    meta: { title: '我的设备', requiresAuth: true, requiresRole: ['owner', 'admin'] }
  },
  {
    path: '/equipments/create',
    name: 'CreateEquipment',
    component: () => import('@/views/CreateEquipment.vue'),
    meta: { title: '添加设备', requiresAuth: true, requiresRole: ['owner', 'admin'] }
  },
  {
    path: '/equipments/:id/edit',
    name: 'EditEquipment',
    component: () => import('@/views/EditEquipment.vue'),
    meta: { title: '编辑设备', requiresAuth: true, requiresRole: ['owner', 'admin'] }
  },
  {
    path: '/orders',
    name: 'MyOrders',
    component: () => import('@/views/MyOrders.vue'),
    meta: { title: '我的订单', requiresAuth: true }
  },
  {
    path: '/orders/:id',
    name: 'OrderDetail',
    component: () => import('@/views/OrderDetail.vue'),
    meta: { title: '订单详情', requiresAuth: true }
  },
  {
    path: '/reviews/:equipmentId',
    name: 'Reviews',
    component: () => import('@/views/Reviews.vue'),
    meta: { title: '评价列表' }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('@/views/Profile.vue'),
    meta: { title: '个人中心', requiresAuth: true }
  },
  {
    path: '/admin/users',
    name: 'UserManagement',
    component: () => import('@/views/UserManagement.vue'),
    meta: { title: '用户管理', requiresAuth: true, requiresRole: ['admin'] }
  },
  {
    path: '/export',
    name: 'Export',
    component: () => import('@/views/Export.vue'),
    meta: { title: '数据导出', requiresAuth: true, requiresRole: ['owner', 'admin'] }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, _from, next) => {
  const userStore = useUserStore()

  document.title = `${to.meta.title || '摄影器材租赁平台'} - 摄影器材租赁平台`

  if (to.meta.guestOnly && userStore.isLoggedIn) {
    next('/')
    return
  }

  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    next('/login')
    return
  }

  if (to.meta.requiresRole && userStore.isLoggedIn) {
    const requiredRoles = to.meta.requiresRole as string[]
    if (!requiredRoles.includes(userStore.userRole)) {
      next('/')
      return
    }
  }

  next()
})

export default router
