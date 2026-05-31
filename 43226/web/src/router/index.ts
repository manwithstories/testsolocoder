import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/stores/user'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/Login.vue'),
    meta: { requiresAuth: false, title: '登录' }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/auth/Register.vue'),
    meta: { requiresAuth: false, title: '注册' }
  },
  {
    path: '/',
    component: () => import('@/layouts/MainLayout.vue'),
    meta: { requiresAuth: false },
    children: [
      {
        path: '',
        name: 'Home',
        component: () => import('@/views/Home.vue'),
        meta: { title: '首页' }
      },
      {
        path: 'exhibitions',
        name: 'ExhibitionList',
        component: () => import('@/views/exhibition/List.vue'),
        meta: { title: '展览列表' }
      },
      {
        path: 'exhibitions/:id',
        name: 'ExhibitionDetail',
        component: () => import('@/views/exhibition/Detail.vue'),
        meta: { title: '展览详情' }
      },
      {
        path: 'collections',
        name: 'CollectionList',
        component: () => import('@/views/collection/List.vue'),
        meta: { title: '藏品列表' }
      },
      {
        path: 'collections/:id',
        name: 'CollectionDetail',
        component: () => import('@/views/collection/Detail.vue'),
        meta: { title: '藏品详情' }
      }
    ]
  },
  {
    path: '/dashboard',
    component: () => import('@/layouts/DashboardLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/Index.vue'),
        meta: { title: '控制面板', roles: ['admin', 'guide'] }
      },
      {
        path: 'collections',
        name: 'DashboardCollections',
        component: () => import('@/views/dashboard/Collections.vue'),
        meta: { title: '藏品管理', roles: ['admin'] }
      },
      {
        path: 'collections/categories',
        name: 'DashboardCategories',
        component: () => import('@/views/dashboard/Categories.vue'),
        meta: { title: '分类管理', roles: ['admin'] }
      },
      {
        path: 'exhibitions',
        name: 'DashboardExhibitions',
        component: () => import('@/views/dashboard/Exhibitions.vue'),
        meta: { title: '展览管理', roles: ['admin'] }
      },
      {
        path: 'exhibitions/:id/time-slots',
        name: 'DashboardTimeSlots',
        component: () => import('@/views/dashboard/TimeSlots.vue'),
        meta: { title: '时段管理', roles: ['admin'] }
      },
      {
        path: 'reservations',
        name: 'DashboardReservations',
        component: () => import('@/views/dashboard/Reservations.vue'),
        meta: { title: '预约管理', roles: ['admin'] }
      },
      {
        path: 'users',
        name: 'DashboardUsers',
        component: () => import('@/views/dashboard/Users.vue'),
        meta: { title: '用户管理', roles: ['admin'] }
      },
      {
        path: 'museums',
        name: 'DashboardMuseums',
        component: () => import('@/views/dashboard/Museums.vue'),
        meta: { title: '博物馆管理', roles: ['admin'] }
      },
      {
        path: 'guide/schedules',
        name: 'DashboardGuideSchedules',
        component: () => import('@/views/dashboard/GuideSchedules.vue'),
        meta: { title: '导览排班', roles: ['guide'] }
      },
      {
        path: 'guide/contents',
        name: 'DashboardGuideContents',
        component: () => import('@/views/dashboard/GuideContents.vue'),
        meta: { title: '导览内容', roles: ['admin', 'guide'] }
      },
      {
        path: 'research',
        name: 'DashboardResearch',
        component: () => import('@/views/dashboard/Research.vue'),
        meta: { title: '学术申请', roles: ['admin'] }
      },
      {
        path: 'statistics',
        name: 'DashboardStatistics',
        component: () => import('@/views/dashboard/Statistics.vue'),
        meta: { title: '统计分析', roles: ['admin'] }
      },
      {
        path: 'profile',
        name: 'DashboardProfile',
        component: () => import('@/views/dashboard/Profile.vue'),
        meta: { title: '个人中心' }
      },
      {
        path: 'my-reservations',
        name: 'DashboardMyReservations',
        component: () => import('@/views/dashboard/MyReservations.vue'),
        meta: { title: '我的预约' }
      },
      {
        path: 'my-visits',
        name: 'DashboardMyVisits',
        component: () => import('@/views/dashboard/MyVisits.vue'),
        meta: { title: '参观记录' }
      },
      {
        path: 'my-research',
        name: 'DashboardMyResearch',
        component: () => import('@/views/dashboard/MyResearch.vue'),
        meta: { title: '我的申请' }
      }
    ]
  },
  {
    path: '/guide/:id',
    name: 'OnlineGuide',
    component: () => import('@/views/guide/OnlineGuide.vue'),
    meta: { title: '在线导览' }
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/NotFound.vue'),
    meta: { title: '页面未找到' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  document.title = `${to.meta.title || '博物馆预约平台'} - 在线博物馆`
  const userStore = useUserStore()
  userStore.initFromStorage()

  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    next({ path: '/login', query: { redirect: to.fullPath } })
    return
  }

  if (to.meta.roles) {
    const roles = to.meta.roles as string[]
    if (userStore.user && !roles.includes(userStore.user.role)) {
      next('/dashboard')
      return
    }
  }

  next()
})

export default router
