<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  LayoutDashboard,
  Clock,
  BarChart3,
  Tags,
  Settings,
  Menu,
  X
} from 'lucide-vue-next'
import ToastContainer from '@/components/ToastContainer.vue'

const route = useRoute()
const router = useRouter()
const isMobileMenuOpen = ref(false)

const navItems = [
  { path: '/', name: '仪表盘', icon: LayoutDashboard },
  { path: '/history', name: '历史', icon: Clock },
  { path: '/statistics', name: '统计', icon: BarChart3 },
  { path: '/categories', name: '分类', icon: Tags },
  { path: '/settings', name: '设置', icon: Settings }
]

const navigate = (path: string) => {
  router.push(path)
  isMobileMenuOpen.value = false
}
</script>

<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 transition-colors duration-300">
    <div class="flex">
      <aside class="hidden lg:flex lg:flex-col w-64 min-h-screen bg-white dark:bg-gray-800 border-r border-gray-200 dark:border-gray-700 fixed left-0 top-0">
        <div class="p-6 border-b border-gray-200 dark:border-gray-700">
          <h1 class="text-xl font-bold text-gray-900 dark:text-white flex items-center gap-2">
            <Clock class="w-6 h-6 text-indigo-500" />
            专注仪表盘
          </h1>
        </div>
        <nav class="flex-1 p-4 space-y-1">
          <button
            v-for="item in navItems"
            :key="item.path"
            @click="navigate(item.path)"
            :class="[
              'w-full flex items-center gap-3 px-4 py-3 rounded-lg transition-all duration-200',
              route.path === item.path
                ? 'bg-indigo-50 dark:bg-indigo-900/30 text-indigo-600 dark:text-indigo-400'
                : 'text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700'
            ]"
          >
            <component :is="item.icon" class="w-5 h-5" />
            <span class="font-medium">{{ item.name }}</span>
          </button>
        </nav>
      </aside>

      <header class="lg:hidden fixed top-0 left-0 right-0 bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 z-50">
        <div class="flex items-center justify-between px-4 py-3">
          <h1 class="text-lg font-bold text-gray-900 dark:text-white flex items-center gap-2">
            <Clock class="w-5 h-5 text-indigo-500" />
            专注仪表盘
          </h1>
          <button
            @click="isMobileMenuOpen = !isMobileMenuOpen"
            class="p-2 text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg"
          >
            <Menu v-if="!isMobileMenuOpen" class="w-6 h-6" />
            <X v-else class="w-6 h-6" />
          </button>
        </div>
        <div
          v-if="isMobileMenuOpen"
          class="bg-white dark:bg-gray-800 border-t border-gray-200 dark:border-gray-700"
        >
          <nav class="p-2 space-y-1">
            <button
              v-for="item in navItems"
              :key="item.path"
              @click="navigate(item.path)"
              :class="[
                'w-full flex items-center gap-3 px-4 py-3 rounded-lg transition-all duration-200',
                route.path === item.path
                  ? 'bg-indigo-50 dark:bg-indigo-900/30 text-indigo-600 dark:text-indigo-400'
                  : 'text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700'
              ]"
            >
              <component :is="item.icon" class="w-5 h-5" />
              <span class="font-medium">{{ item.name }}</span>
            </button>
          </nav>
        </div>
      </header>

      <nav class="lg:hidden fixed bottom-0 left-0 right-0 bg-white dark:bg-gray-800 border-t border-gray-200 dark:border-gray-700 z-50">
        <div class="flex justify-around py-2">
          <button
            v-for="item in navItems.slice(0, 5)"
            :key="item.path"
            @click="navigate(item.path)"
            :class="[
              'flex flex-col items-center gap-1 px-3 py-2 rounded-lg transition-all duration-200',
              route.path === item.path
                ? 'text-indigo-600 dark:text-indigo-400'
                : 'text-gray-500 dark:text-gray-400'
            ]"
          >
            <component :is="item.icon" class="w-5 h-5" />
            <span class="text-xs">{{ item.name }}</span>
          </button>
        </div>
      </nav>

      <main class="flex-1 lg:ml-64 pt-16 lg:pt-0 pb-20 lg:pb-0">
        <div class="p-4 lg:p-8">
          <router-view />
        </div>
      </main>
    </div>

    <ToastContainer />
  </div>
</template>
