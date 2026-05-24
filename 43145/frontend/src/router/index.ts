import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/store'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/Register.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/',
    component: () => import('@/views/Layout.vue'),
    meta: { requiresAuth: true },
    redirect: '/surveys',
    children: [
      {
        path: 'surveys',
        name: 'SurveyList',
        component: () => import('@/views/SurveyList.vue'),
        meta: { title: '问卷列表' }
      },
      {
        path: 'surveys/create',
        name: 'SurveyCreate',
        component: () => import('@/views/SurveyEditor.vue'),
        meta: { title: '创建问卷' }
      },
      {
        path: 'surveys/:id/edit',
        name: 'SurveyEdit',
        component: () => import('@/views/SurveyEditor.vue'),
        meta: { title: '编辑问卷' }
      },
      {
        path: 'surveys/:id/statistics',
        name: 'SurveyStatistics',
        component: () => import('@/views/SurveyStatistics.vue'),
        meta: { title: '数据统计' }
      },
      {
        path: 'users',
        name: 'UserManagement',
        component: () => import('@/views/UserManagement.vue'),
        meta: { title: '用户管理', roles: ['admin'] }
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/Profile.vue'),
        meta: { title: '个人中心' }
      }
    ]
  },
  {
    path: '/survey/fill/:token',
    name: 'SurveyFill',
    component: () => import('@/views/SurveyFill.vue'),
    meta: { requiresAuth: false, title: '填写问卷' }
  },
  {
    path: '/survey/thank-you',
    name: 'ThankYou',
    component: () => import('@/views/ThankYou.vue'),
    meta: { requiresAuth: false }
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
  const token = userStore.token

  if (to.meta.requiresAuth && !token) {
    next({ path: '/login', query: { redirect: to.fullPath } })
    return
  }

  if (to.meta.roles && userStore.user) {
    const roles = to.meta.roles as string[]
    if (!roles.includes(userStore.user.role)) {
      next('/surveys')
      return
    }
  }

  next()
})

export default router
