import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/stores/user'
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
    component: () => import('@/views/auth/Login.vue'),
    meta: { title: '登录', guest: true }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/auth/Register.vue'),
    meta: { title: '注册', guest: true }
  },
  {
    path: '/questions',
    name: 'Questions',
    component: () => import('@/views/question/QuestionList.vue'),
    meta: { title: '问题列表' }
  },
  {
    path: '/questions/:id',
    name: 'QuestionDetail',
    component: () => import('@/views/question/QuestionDetail.vue'),
    meta: { title: '问题详情' }
  },
  {
    path: '/questions/ask',
    name: 'AskQuestion',
    component: () => import('@/views/question/AskQuestion.vue'),
    meta: { title: '提问', requiresAuth: true }
  },
  {
    path: '/user/profile',
    name: 'UserProfile',
    component: () => import('@/views/user/Profile.vue'),
    meta: { title: '个人中心', requiresAuth: true }
  },
  {
    path: '/user/points',
    name: 'UserPoints',
    component: () => import('@/views/user/Points.vue'),
    meta: { title: '积分中心', requiresAuth: true }
  },
  {
    path: '/user/favorites',
    name: 'UserFavorites',
    component: () => import('@/views/user/Favorites.vue'),
    meta: { title: '我的收藏', requiresAuth: true }
  },
  {
    path: '/user/notifications',
    name: 'UserNotifications',
    component: () => import('@/views/notifications/NotificationList.vue'),
    meta: { title: '消息通知', requiresAuth: true }
  },
  {
    path: '/users/:id',
    name: 'UserDetail',
    component: () => import('@/views/user/UserDetail.vue'),
    meta: { title: '用户主页' }
  },
  {
    path: '/search',
    name: 'Search',
    component: () => import('@/views/Search.vue'),
    meta: { title: '搜索' }
  },
  {
    path: '/admin',
    name: 'Admin',
    component: () => import('@/views/admin/AdminLayout.vue'),
    meta: { title: '管理后台', requiresAuth: true, requiresAdmin: true },
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
        component: () => import('@/views/admin/UserList.vue'),
        meta: { title: '用户管理' }
      },
      {
        path: 'audit',
        name: 'AdminAudit',
        component: () => import('@/views/admin/AuditList.vue'),
        meta: { title: '内容审核' }
      },
      {
        path: 'reports',
        name: 'AdminReports',
        component: () => import('@/views/admin/ReportList.vue'),
        meta: { title: '举报处理' }
      },
      {
        path: 'sensitive-words',
        name: 'AdminSensitiveWords',
        component: () => import('@/views/admin/SensitiveWords.vue'),
        meta: { title: '敏感词管理' }
      },
      {
        path: 'rewards',
        name: 'AdminRewards',
        component: () => import('@/views/admin/RewardList.vue'),
        meta: { title: '奖品管理' }
      },
      {
        path: 'expert-applications',
        name: 'AdminExpertApplications',
        component: () => import('@/views/admin/ExpertApplications.vue'),
        meta: { title: '专家审核' }
      },
      {
        path: 'stats',
        name: 'AdminStats',
        component: () => import('@/views/admin/StatsReport.vue'),
        meta: { title: '统计报表' }
      }
    ]
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
  NProgress.start()
  document.title = `${to.meta.title || '问答平台'} - 知识问答社区`

  const userStore = useUserStore()
  userStore.loadUserInfo()

  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    next({ name: 'Login', query: { redirect: to.fullPath } })
    return
  }

  if (to.meta.requiresAdmin && !userStore.isAdmin) {
    next({ name: 'Home' })
    return
  }

  if (to.meta.guest && userStore.isLoggedIn) {
    next({ name: 'Home' })
    return
  }

  next()
})

router.afterEach(() => {
  NProgress.done()
})

export default router
