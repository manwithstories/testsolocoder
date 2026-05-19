import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/dashboard'
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: () => import('@/views/Dashboard.vue'),
    meta: { title: '仪表盘', icon: 'Odometer' }
  },
  {
    path: '/books',
    name: 'Books',
    component: () => import('@/views/BookList.vue'),
    meta: { title: '图书管理', icon: 'Reading' }
  },
  {
    path: '/books/:id',
    name: 'BookDetail',
    component: () => import('@/views/BookDetail.vue'),
    meta: { title: '图书详情', hidden: true }
  },
  {
    path: '/reading',
    name: 'Reading',
    component: () => import('@/views/ReadingProgress.vue'),
    meta: { title: '阅读进度', icon: 'Clock' }
  },
  {
    path: '/notes',
    name: 'Notes',
    component: () => import('@/views/Notes.vue'),
    meta: { title: '读书笔记', icon: 'EditPen' }
  },
  {
    path: '/tags',
    name: 'Tags',
    component: () => import('@/views/TagsCategories.vue'),
    meta: { title: '标签分类', icon: 'PriceTag' }
  },
  {
    path: '/borrows',
    name: 'Borrows',
    component: () => import('@/views/Borrows.vue'),
    meta: { title: '借阅管理', icon: 'Share' }
  },
  {
    path: '/stats',
    name: 'Stats',
    component: () => import('@/views/Statistics.vue'),
    meta: { title: '统计分析', icon: 'DataLine' }
  },
  {
    path: '/goals',
    name: 'Goals',
    component: () => import('@/views/Goals.vue'),
    meta: { title: '阅读目标', icon: 'Flag' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, _from, next) => {
  document.title = (to.meta.title as string) ? `${to.meta.title} - 图书管理系统` : '图书管理系统'
  next()
})

export { routes }
export default router
