import { createRouter, createWebHashHistory, RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes: RouteRecordRaw[] = [
  { path: '/login', component: () => import('@/views/Login.vue') },
  { path: '/register', component: () => import('@/views/Register.vue') },
  {
    path: '/',
    component: () => import('@/views/Layout.vue'),
    redirect: '/watches',
    children: [
      { path: 'watches', component: () => import('@/views/Watches.vue'), meta: { title: '手表市场' } },
      { path: 'watches/:id', component: () => import('@/views/WatchDetail.vue'), meta: { title: '手表详情' } },
      { path: 'my-watches', component: () => import('@/views/MyWatches.vue'), meta: { title: '我的手表', requiresRole: 'seller' } },
      { path: 'publish', component: () => import('@/views/PublishWatch.vue'), meta: { title: '发布手表', requiresRole: 'seller' } },
      { path: 'auth-orders', component: () => import('@/views/AuthOrders.vue'), meta: { title: '鉴定申请' } },
      { path: 'auth-review', component: () => import('@/views/AuthReview.vue'), meta: { title: '鉴定审核', requiresRole: 'appraiser' } },
      { path: 'trades', component: () => import('@/views/Trades.vue'), meta: { title: '我的交易' } },
      { path: 'favorites', component: () => import('@/views/Favorites.vue'), meta: { title: '收藏夹' } },
      { path: 'messages', component: () => import('@/views/Messages.vue'), meta: { title: '消息中心' } },
      { path: 'stats', component: () => import('@/views/Stats.vue'), meta: { title: '数据统计' } },
      { path: 'profile', component: () => import('@/views/Profile.vue'), meta: { title: '个人中心' } }
    ]
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

router.beforeEach((to, _from, next) => {
  const auth = useAuthStore()
  if (to.path === '/login' || to.path === '/register') {
    return next()
  }
  if (!auth.token) {
    return next('/login')
  }
  if (to.meta.requiresRole && auth.user?.role !== to.meta.requiresRole) {
    return next('/')
  }
  next()
})

export default router
