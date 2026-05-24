import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: () => import('@/layouts/DefaultLayout.vue'),
    redirect: '/home',
    children: [
      {
        path: 'home',
        name: 'Home',
        component: () => import('@/views/Home.vue'),
        meta: { title: '首页' }
      },
      {
        path: 'ranking',
        name: 'Ranking',
        component: () => import('@/views/Ranking.vue'),
        meta: { title: '排行榜' }
      },
      {
        path: 'works',
        name: 'Works',
        component: () => import('@/views/Works.vue'),
        meta: { title: '作品' }
      },
      {
        path: 'work/:id',
        name: 'WorkDetail',
        component: () => import('@/views/WorkDetail.vue'),
        meta: { title: '作品详情' }
      },
      {
        path: 'artists',
        name: 'Artists',
        component: () => import('@/views/Artists.vue'),
        meta: { title: '音乐人' }
      },
      {
        path: 'artist/:id',
        name: 'ArtistDetail',
        component: () => import('@/views/ArtistDetail.vue'),
        meta: { title: '音乐人详情' }
      },
      {
        path: 'events',
        name: 'Events',
        component: () => import('@/views/Events.vue'),
        meta: { title: '演出' }
      },
      {
        path: 'event/:id',
        name: 'EventDetail',
        component: () => import('@/views/EventDetail.vue'),
        meta: { title: '演出详情' }
      },
      {
        path: 'playlists',
        name: 'Playlists',
        component: () => import('@/views/Playlists.vue'),
        meta: { title: '歌单' }
      },
      {
        path: 'playlist/:id',
        name: 'PlaylistDetail',
        component: () => import('@/views/PlaylistDetail.vue'),
        meta: { title: '歌单详情' }
      }
    ]
  },
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
    path: '/user',
    component: () => import('@/layouts/UserLayout.vue'),
    redirect: '/user/profile',
    meta: { requiresAuth: true },
    children: [
      {
        path: 'profile',
        name: 'UserProfile',
        component: () => import('@/views/user/Profile.vue'),
        meta: { title: '个人中心' }
      },
      {
        path: 'works',
        name: 'UserWorks',
        component: () => import('@/views/user/MyWorks.vue'),
        meta: { title: '我的作品' }
      },
      {
        path: 'upload',
        name: 'UploadWork',
        component: () => import('@/views/user/UploadWork.vue'),
        meta: { title: '上传作品' }
      },
      {
        path: 'albums',
        name: 'UserAlbums',
        component: () => import('@/views/user/MyAlbums.vue'),
        meta: { title: '我的专辑' }
      },
      {
        path: 'events',
        name: 'UserEvents',
        component: () => import('@/views/user/MyEvents.vue'),
        meta: { title: '我的演出' }
      },
      {
        path: 'tickets',
        name: 'MyTickets',
        component: () => import('@/views/user/MyTickets.vue'),
        meta: { title: '我的票' }
      },
      {
        path: 'playlists',
        name: 'MyPlaylists',
        component: () => import('@/views/user/MyPlaylists.vue'),
        meta: { title: '我的歌单' }
      },
      {
        path: 'revenue',
        name: 'MyRevenue',
        component: () => import('@/views/user/MyRevenue.vue'),
        meta: { title: '我的收益' }
      },
      {
        path: 'withdraw',
        name: 'Withdraw',
        component: () => import('@/views/user/Withdraw.vue'),
        meta: { title: '提现' }
      },
      {
        path: 'stats',
        name: 'MyStats',
        component: () => import('@/views/user/MyStats.vue'),
        meta: { title: '数据统计' }
      },
      {
        path: 'follows',
        name: 'MyFollows',
        component: () => import('@/views/user/MyFollows.vue'),
        meta: { title: '关注列表' }
      },
      {
        path: 'notifications',
        name: 'MyNotifications',
        component: () => import('@/views/user/MyNotifications.vue'),
        meta: { title: '消息通知' }
      }
    ]
  },
  {
    path: '/admin',
    component: () => import('@/layouts/AdminLayout.vue'),
    redirect: '/admin/dashboard',
    meta: { requiresAuth: true, requiresAdmin: true },
    children: [
      {
        path: 'dashboard',
        name: 'AdminDashboard',
        component: () => import('@/views/admin/Dashboard.vue'),
        meta: { title: '管理后台' }
      },
      {
        path: 'users',
        name: 'AdminUsers',
        component: () => import('@/views/admin/AdminUsers.vue'),
        meta: { title: '用户管理' }
      },
      {
        path: 'works',
        name: 'AdminWorks',
        component: () => import('@/views/admin/AdminWorks.vue'),
        meta: { title: '作品管理' }
      },
      {
        path: 'events',
        name: 'AdminEvents',
        component: () => import('@/views/admin/AdminEvents.vue'),
        meta: { title: '演出管理' }
      },
      {
        path: 'withdraw',
        name: 'AdminWithdraw',
        component: () => import('@/views/admin/AdminWithdraw.vue'),
        meta: { title: '提现审核' }
      },
      {
        path: 'logs',
        name: 'AdminLogs',
        component: () => import('@/views/admin/AdminLogs.vue'),
        meta: { title: '操作日志' }
      }
    ]
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

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  
  if (to.meta.requiresAuth && !token) {
    next({ name: 'Login', query: { redirect: to.fullPath } })
    return
  }
  
  if (to.meta.requiresAdmin) {
    const user = JSON.parse(localStorage.getItem('user') || '{}')
    if (user.role !== 'admin') {
      next({ name: 'Home' })
      return
    }
  }
  
  document.title = `${to.meta.title || '独立音乐人平台'} - 独立音乐人平台`
  next()
})

export default router
