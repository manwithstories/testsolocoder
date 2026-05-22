import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/store/user'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
    meta: {
      title: '登录',
      requiresAuth: false
    }
  },
  {
    path: '/',
    component: () => import('@/layout/index.vue'),
    redirect: '/dashboard',
    meta: {
      requiresAuth: true
    },
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: {
          title: '仪表盘',
          icon: 'Odometer'
        }
      },
      {
        path: 'room-status',
        name: 'RoomStatusBoard',
        component: () => import('@/views/dashboard/RoomStatusBoard.vue'),
        meta: {
          title: '房间状态看板',
          icon: 'View'
        }
      },
      {
        path: 'booking',
        name: 'Booking',
        meta: {
          title: '预订管理',
          icon: 'Booking'
        },
        children: [
          {
            path: 'list',
            name: 'BookingList',
            component: () => import('@/views/booking/BookingList.vue'),
            meta: {
              title: '预订列表',
              icon: 'List'
            }
          },
          {
            path: 'create',
            name: 'BookingCreate',
            component: () => import('@/views/booking/BookingCreate.vue'),
            meta: {
              title: '创建预订',
              icon: 'Plus'
            }
          }
        ]
      },
      {
        path: 'checkin',
        name: 'CheckIn',
        meta: {
          title: '入住管理',
          icon: 'HomeFilled'
        },
        children: [
          {
            path: 'list',
            name: 'CheckInList',
            component: () => import('@/views/checkin/CheckInList.vue'),
            meta: {
              title: '入住列表',
              icon: 'List'
            }
          }
        ]
      },
      {
        path: 'room',
        name: 'Room',
        meta: {
          title: '房间管理',
          icon: 'OfficeBuilding'
        },
        children: [
          {
            path: 'list',
            name: 'RoomList',
            component: () => import('@/views/room/RoomList.vue'),
            meta: {
              title: '房间列表',
              icon: 'List'
            }
          },
          {
            path: 'type',
            name: 'RoomType',
            component: () => import('@/views/room/RoomTypeList.vue'),
            meta: {
              title: '房型管理',
              icon: 'Menu'
            }
          }
        ]
      },
      {
        path: 'payment',
        name: 'Payment',
        meta: {
          title: '支付管理',
          icon: 'Money'
        },
        children: [
          {
            path: 'list',
            name: 'PaymentList',
            component: () => import('@/views/payment/PaymentList.vue'),
            meta: {
              title: '支付记录',
              icon: 'List'
            }
          }
        ]
      },
      {
        path: 'member',
        name: 'Member',
        meta: {
          title: '会员管理',
          icon: 'UserFilled'
        },
        children: [
          {
            path: 'list',
            name: 'MemberList',
            component: () => import('@/views/member/MemberList.vue'),
            meta: {
              title: '会员列表',
              icon: 'List'
            }
          },
          {
            path: 'level',
            name: 'MemberLevel',
            component: () => import('@/views/member/MemberLevelList.vue'),
            meta: {
              title: '会员等级',
              icon: 'Medal'
            }
          }
        ]
      },
      {
        path: 'report',
        name: 'Report',
        meta: {
          title: '报表统计',
          icon: 'DataLine'
        },
        children: [
          {
            path: 'occupancy',
            name: 'OccupancyReport',
            component: () => import('@/views/report/OccupancyReport.vue'),
            meta: {
              title: '入住率报表',
              icon: 'TrendCharts'
            }
          },
          {
            path: 'revenue',
            name: 'RevenueReport',
            component: () => import('@/views/report/RevenueReport.vue'),
            meta: {
              title: '营收报表',
              icon: 'Money'
            }
          }
        ]
      },
      {
        path: 'system',
        name: 'System',
        meta: {
          title: '系统管理',
          icon: 'Setting'
        },
        children: [
          {
            path: 'user',
            name: 'SystemUser',
            component: () => import('@/views/user/UserList.vue'),
            meta: {
              title: '用户管理',
              icon: 'User'
            }
          }
        ]
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/profile/Profile.vue'),
        meta: {
          title: '个人中心',
          icon: 'User'
        }
      }
    ]
  },
  {
    path: '/404',
    name: 'NotFound',
    component: () => import('@/views/error/404.vue'),
    meta: {
      title: '页面不存在'
    }
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/404'
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior: () => ({ left: 0, top: 0 })
})

router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  const token = userStore.token

  if (to.meta.title) {
    document.title = `${to.meta.title} - 酒店管理系统`
  }

  if (to.path === '/login') {
    if (token) {
      next('/')
    } else {
      next()
    }
    return
  }

  if (!token) {
    next({
      path: '/login',
      query: { redirect: to.fullPath }
    })
    return
  }

  if (to.meta.requiresAuth && !userStore.userInfo) {
    userStore.getProfile()
      .then(() => {
        next()
      })
      .catch(() => {
        userStore.logout()
        next({
          path: '/login',
          query: { redirect: to.fullPath }
        })
      })
    return
  }

  next()
})

router.afterEach(() => {
  // 可以在这里添加页面加载完成后的逻辑
})

export default router
