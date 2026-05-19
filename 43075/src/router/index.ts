import { createRouter, createWebHistory } from 'vue-router'
import Dashboard from '@/views/Dashboard.vue'
import History from '@/views/History.vue'
import Statistics from '@/views/Statistics.vue'
import Categories from '@/views/Categories.vue'
import Settings from '@/views/Settings.vue'

const routes = [
  {
    path: '/',
    name: 'dashboard',
    component: Dashboard,
    meta: { title: '仪表盘' }
  },
  {
    path: '/history',
    name: 'history',
    component: History,
    meta: { title: '历史记录' }
  },
  {
    path: '/statistics',
    name: 'statistics',
    component: Statistics,
    meta: { title: '数据统计' }
  },
  {
    path: '/categories',
    name: 'categories',
    component: Categories,
    meta: { title: '分类管理' }
  },
  {
    path: '/settings',
    name: 'settings',
    component: Settings,
    meta: { title: '应用设置' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, _from, next) => {
  document.title = `${to.meta.title || '专注计时器'} - 专注仪表盘`
  next()
})

export default router
