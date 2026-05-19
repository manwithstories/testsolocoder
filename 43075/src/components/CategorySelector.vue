<script setup lang="ts">
import { useCategoriesStore } from '@/stores/categories'
import { useTimerStore } from '@/stores/timer'
import * as LucideIcons from 'lucide-vue-next'

const categoriesStore = useCategoriesStore()
const timerStore = useTimerStore()

const getIcon = (iconName: string) => {
  return (LucideIcons as any)[iconName] || LucideIcons.Circle
}
</script>

<template>
  <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-lg p-6">
    <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">选择分类</h3>
    <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-5 gap-3">
      <button
        v-for="category in categoriesStore.categories"
        :key="category.id"
        @click="categoriesStore.selectCategory(category.id)"
        :disabled="timerStore.isRunning"
        :class="[
          'flex flex-col items-center gap-2 p-4 rounded-xl transition-all duration-200 border-2',
          categoriesStore.selectedCategoryId === category.id
            ? 'border-indigo-500 bg-indigo-50 dark:bg-indigo-900/20'
            : 'border-transparent bg-gray-50 dark:bg-gray-700/50 hover:bg-gray-100 dark:hover:bg-gray-700',
          timerStore.isRunning ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'
        ]"
      >
        <div
          class="w-12 h-12 rounded-full flex items-center justify-center"
          :style="{ backgroundColor: category.color + '20' }"
        >
          <component
            :is="getIcon(category.icon)"
            class="w-6 h-6"
            :style="{ color: category.color }"
          />
        </div>
        <span class="text-sm font-medium text-gray-700 dark:text-gray-300">{{ category.name }}</span>
      </button>
    </div>
  </div>
</template>
