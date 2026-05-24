import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/stores/user'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/jobs'
  },
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
    path: '/jobs',
    name: 'JobList',
    component: () => import('@/views/jobs/JobList.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/jobs/:id',
    name: 'JobDetail',
    component: () => import('@/views/jobs/JobDetail.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/my/jobs',
    name: 'MyJobs',
    component: () => import('@/views/jobs/MyJobs.vue'),
    meta: { requiresAuth: true, roles: ['company'] }
  },
  {
    path: '/my/jobs/create',
    name: 'CreateJob',
    component: () => import('@/views/jobs/CreateJob.vue'),
    meta: { requiresAuth: true, roles: ['company'] }
  },
  {
    path: '/my/jobs/:id/edit',
    name: 'EditJob',
    component: () => import('@/views/jobs/EditJob.vue'),
    meta: { requiresAuth: true, roles: ['company'] }
  },
  {
    path: '/resumes',
    name: 'MyResumes',
    component: () => import('@/views/resumes/ResumeList.vue'),
    meta: { requiresAuth: true, roles: ['applicant'] }
  },
  {
    path: '/resumes/create',
    name: 'CreateResume',
    component: () => import('@/views/resumes/CreateResume.vue'),
    meta: { requiresAuth: true, roles: ['applicant'] }
  },
  {
    path: '/resumes/:id',
    name: 'ResumeDetail',
    component: () => import('@/views/resumes/ResumeDetail.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/resumes/:id/edit',
    name: 'EditResume',
    component: () => import('@/views/resumes/EditResume.vue'),
    meta: { requiresAuth: true, roles: ['applicant'] }
  },
  {
    path: '/applications',
    name: 'MyApplications',
    component: () => import('@/views/applications/ApplicationList.vue'),
    meta: { requiresAuth: true, roles: ['applicant'] }
  },
  {
    path: '/applications/:id',
    name: 'ApplicationDetail',
    component: () => import('@/views/applications/ApplicationDetail.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/company/applications',
    name: 'CompanyApplications',
    component: () => import('@/views/applications/CompanyApplicationList.vue'),
    meta: { requiresAuth: true, roles: ['company'] }
  },
  {
    path: '/interviews',
    name: 'MyInterviews',
    component: () => import('@/views/interviews/InterviewList.vue'),
    meta: { requiresAuth: true, roles: ['applicant'] }
  },
  {
    path: '/company/interviews',
    name: 'CompanyInterviews',
    component: () => import('@/views/interviews/CompanyInterviewList.vue'),
    meta: { requiresAuth: true, roles: ['company'] }
  },
  {
    path: '/company/interviews/schedule',
    name: 'ScheduleInterview',
    component: () => import('@/views/interviews/ScheduleInterview.vue'),
    meta: { requiresAuth: true, roles: ['company'] }
  },
  {
    path: '/stats',
    name: 'Statistics',
    component: () => import('@/views/stats/Statistics.vue'),
    meta: { requiresAuth: true, roles: ['company', 'admin'] }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('@/views/profile/Profile.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/NotFound.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, _from, next) => {
  const userStore = useUserStore()

  if (!userStore.user && userStore.token) {
    userStore.initFromStorage()
  }

  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    next({ name: 'Login', query: { redirect: to.fullPath } })
    return
  }

  if (to.meta.roles && userStore.user) {
    const roles = to.meta.roles as string[]
    if (!roles.includes(userStore.userRole)) {
      next({ name: 'JobList' })
      return
    }
  }

  if ((to.name === 'Login' || to.name === 'Register') && userStore.isLoggedIn) {
    next({ name: 'JobList' })
    return
  }

  next()
})

export default router
