import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'
import { useUserStore } from '@/stores/user'

NProgress.configure({ showSpinner: false })

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/Login.vue'),
    meta: { title: '登录' }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/login/Register.vue'),
    meta: { title: '注册' }
  },
  {
    path: '/',
    component: () => import('@/layouts/MainLayout.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/common/Dashboard.vue'),
        meta: { title: '首页', requiresAuth: true }
      }
    ]
  },
  {
    path: '/hr',
    component: () => import('@/layouts/MainLayout.vue'),
    meta: { requiresAuth: true, role: 'hr' },
    children: [
      {
        path: 'employees',
        name: 'HREmployees',
        component: () => import('@/views/hr/EmployeeList.vue'),
        meta: { title: '员工管理' }
      },
      {
        path: 'departments',
        name: 'HRDepartments',
        component: () => import('@/views/hr/DepartmentList.vue'),
        meta: { title: '部门管理' }
      },
      {
        path: 'appointments',
        name: 'HRAppointments',
        component: () => import('@/views/hr/AppointmentList.vue'),
        meta: { title: '预约管理' }
      },
      {
        path: 'reports',
        name: 'HRReports',
        component: () => import('@/views/hr/ReportList.vue'),
        meta: { title: '体检报告' }
      },
      {
        path: 'budget',
        name: 'HRBudget',
        component: () => import('@/views/hr/Budget.vue'),
        meta: { title: '预算管理' }
      },
      {
        path: 'department-appointments',
        name: 'HRDeptAppointments',
        component: () => import('@/views/hr/DeptAppointment.vue'),
        meta: { title: '部门预约分配' }
      },
      {
        path: 'statistics',
        name: 'HRStatistics',
        component: () => import('@/views/hr/Statistics.vue'),
        meta: { title: '数据统计' }
      },
      {
        path: 'billings',
        name: 'HRBillings',
        component: () => import('@/views/hr/BillingList.vue'),
        meta: { title: '账单管理' }
      },
      {
        path: 'transactions',
        name: 'HRTransactions',
        component: () => import('@/views/hr/TransactionList.vue'),
        meta: { title: '交易记录' }
      },
      {
        path: 'balance',
        name: 'HRBalance',
        component: () => import('@/views/hr/Balance.vue'),
        meta: { title: '账户余额' }
      }
    ]
  },
  {
    path: '/agency',
    component: () => import('@/layouts/MainLayout.vue'),
    meta: { requiresAuth: true, role: 'agency' },
    children: [
      {
        path: 'packages',
        name: 'AgencyPackages',
        component: () => import('@/views/agency/PackageList.vue'),
        meta: { title: '套餐管理' }
      },
      {
        path: 'appointments',
        name: 'AgencyAppointments',
        component: () => import('@/views/agency/AppointmentList.vue'),
        meta: { title: '预约管理' }
      },
      {
        path: 'reports',
        name: 'AgencyReports',
        component: () => import('@/views/agency/ReportUpload.vue'),
        meta: { title: '报告上传' }
      },
      {
        path: 'billings',
        name: 'AgencyBillings',
        component: () => import('@/views/agency/BillingList.vue'),
        meta: { title: '账单管理' }
      },
      {
        path: 'timeslots',
        name: 'AgencyTimeSlots',
        component: () => import('@/views/agency/TimeSlotManage.vue'),
        meta: { title: '时段管理' }
      }
    ]
  },
  {
    path: '/employee',
    component: () => import('@/layouts/MainLayout.vue'),
    meta: { requiresAuth: true, role: 'employee' },
    children: [
      {
        path: 'appointments',
        name: 'EmployeeAppointments',
        component: () => import('@/views/employee/AppointmentList.vue'),
        meta: { title: '我的预约' }
      },
      {
        path: 'new-appointment',
        name: 'EmployeeNewAppointment',
        component: () => import('@/views/employee/NewAppointment.vue'),
        meta: { title: '新建预约' }
      },
      {
        path: 'reports',
        name: 'EmployeeReports',
        component: () => import('@/views/employee/ReportList.vue'),
        meta: { title: '我的报告' }
      },
      {
        path: 'report/:id',
        name: 'EmployeeReportDetail',
        component: () => import('@/views/employee/ReportDetail.vue'),
        meta: { title: '报告详情' }
      },
      {
        path: 'health-records',
        name: 'EmployeeHealthRecords',
        component: () => import('@/views/employee/HealthRecords.vue'),
        meta: { title: '健康档案' }
      },
      {
        path: 'health-trend',
        name: 'EmployeeHealthTrend',
        component: () => import('@/views/employee/HealthTrend.vue'),
        meta: { title: '趋势分析' }
      },
      {
        path: 'abnormal-items',
        name: 'EmployeeAbnormalItems',
        component: () => import('@/views/employee/AbnormalItems.vue'),
        meta: { title: '异常指标' }
      },
      {
        path: 'reminders',
        name: 'EmployeeReminders',
        component: () => import('@/views/employee/Reminders.vue'),
        meta: { title: '复查提醒' }
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/common/NotFound.vue'),
    meta: { title: '页面未找到' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  NProgress.start()
  document.title = `${to.meta.title || '健康管理平台'} - 企业员工体检预约系统`

  const userStore = useUserStore()

  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    next({ path: '/login', query: { redirect: to.fullPath } })
    return
  }

  if (to.meta.role && userStore.userRole !== to.meta.role && userStore.userRole !== 'admin') {
    next('/dashboard')
    return
  }

  if ((to.path === '/login' || to.path === '/register') && userStore.isLoggedIn) {
    next('/dashboard')
    return
  }

  next()
})

router.afterEach(() => {
  NProgress.done()
})

export default router
