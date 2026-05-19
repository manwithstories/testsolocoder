<script setup lang="ts">
import { computed } from 'vue'
import { Play, Pause, Square, Clock, Timer } from 'lucide-vue-next'
import { useTimerStore } from '@/stores/timer'
import { useCategoriesStore } from '@/stores/categories'
import { useSettingsStore } from '@/stores/settings'
import { formatTimeDigits } from '@/utils/date'
import * as LucideIcons from 'lucide-vue-next'

const timerStore = useTimerStore()
const categoriesStore = useCategoriesStore()
const settingsStore = useSettingsStore()

const progressStyle = computed(() => {
  const circumference = 2 * Math.PI * 120
  const offset = circumference - (timerStore.progress / 100) * circumference
  return {
    strokeDasharray: `${circumference} ${circumference}`,
    strokeDashoffset: offset
  }
})

const getIcon = (iconName: string) => {
  return (LucideIcons as any)[iconName] || LucideIcons.Circle
}
</script>

<template>
  <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-lg p-6 lg:p-8">
    <div class="flex items-center justify-center gap-4 mb-8">
      <button
        @click="timerStore.setMode('pomodoro')"
        :class="[
          'flex items-center gap-2 px-4 py-2 rounded-lg font-medium transition-all',
          timerStore.mode === 'pomodoro'
            ? 'bg-indigo-500 text-white shadow-md'
            : 'bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600'
        ]"
        :disabled="timerStore.isRunning"
      >
        <Timer class="w-4 h-4" />
        番茄钟
      </button>
      <button
        @click="timerStore.setMode('custom')"
        :class="[
          'flex items-center gap-2 px-4 py-2 rounded-lg font-medium transition-all',
          timerStore.mode === 'custom'
            ? 'bg-indigo-500 text-white shadow-md'
            : 'bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600'
        ]"
        :disabled="timerStore.isRunning"
      >
        <Clock class="w-4 h-4" />
        自定义
      </button>
    </div>

    <div v-if="timerStore.mode === 'custom' && timerStore.isIdle" class="flex justify-center mb-6">
      <div class="flex items-center gap-2">
        <button
          @click="timerStore.setCustomDuration(Math.max(1, timerStore.customDuration - 5))"
          class="w-10 h-10 rounded-full bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 flex items-center justify-center text-gray-600 dark:text-gray-300 font-bold"
        >
          -
        </button>
        <div class="w-24 text-center">
          <input
            type="number"
            :value="timerStore.customDuration"
            @input="(e) => timerStore.setCustomDuration(parseInt((e.target as HTMLInputElement).value) || 1)"
            class="w-full text-center text-2xl font-bold bg-transparent text-gray-900 dark:text-white focus:outline-none"
            min="1"
            max="180"
          />
          <span class="text-sm text-gray-500 dark:text-gray-400">分钟</span>
        </div>
        <button
          @click="timerStore.setCustomDuration(Math.min(180, timerStore.customDuration + 5))"
          class="w-10 h-10 rounded-full bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 flex items-center justify-center text-gray-600 dark:text-gray-300 font-bold"
        >
          +
        </button>
      </div>
    </div>

    <div class="flex justify-center mb-6">
      <div class="relative w-64 h-64 lg:w-72 lg:h-72">
        <svg class="w-full h-full transform -rotate-90">
          <circle
            cx="128"
            cy="128"
            r="120"
            fill="none"
            stroke="currentColor"
            stroke-width="8"
            class="text-gray-200 dark:text-gray-700"
          />
          <circle
            cx="128"
            cy="128"
            r="120"
            fill="none"
            stroke="currentColor"
            stroke-width="8"
            stroke-linecap="round"
            class="text-indigo-500 transition-all duration-1000"
            :style="progressStyle"
          />
        </svg>
        <div class="absolute inset-0 flex flex-col items-center justify-center">
          <div class="text-5xl lg:text-6xl font-bold text-gray-900 dark:text-white mb-2">
            {{ formatTimeDigits(timerStore.remainingTime) }}
          </div>
          <div v-if="categoriesStore.selectedCategory" class="flex items-center gap-2 text-gray-600 dark:text-gray-400">
            <component
              :is="getIcon(categoriesStore.selectedCategory.icon)"
              class="w-4 h-4"
              :style="{ color: categoriesStore.selectedCategory.color }"
            />
            <span class="text-sm">{{ categoriesStore.selectedCategory.name }}</span>
          </div>
          <div v-else class="text-sm text-gray-400 dark:text-gray-500">
            请选择分类
          </div>
        </div>
      </div>
    </div>

    <div class="flex justify-center gap-4">
      <button
        v-if="timerStore.isIdle || timerStore.isPaused"
        @click="timerStore.start()"
        :disabled="!categoriesStore.selectedCategoryId"
        class="w-16 h-16 rounded-full bg-indigo-500 hover:bg-indigo-600 disabled:bg-gray-300 dark:disabled:bg-gray-600 text-white flex items-center justify-center shadow-lg transition-all hover:scale-105 active:scale-95"
      >
        <Play class="w-8 h-8 ml-1" />
      </button>
      <button
        v-if="timerStore.isRunning"
        @click="timerStore.pause()"
        class="w-16 h-16 rounded-full bg-amber-500 hover:bg-amber-600 text-white flex items-center justify-center shadow-lg transition-all hover:scale-105 active:scale-95"
      >
        <Pause class="w-8 h-8" />
      </button>
      <button
        v-if="!timerStore.isIdle"
        @click="timerStore.stop()"
        class="w-16 h-16 rounded-full bg-red-500 hover:bg-red-600 text-white flex items-center justify-center shadow-lg transition-all hover:scale-105 active:scale-95"
      >
        <Square class="w-6 h-6" />
      </button>
    </div>

    <div v-if="timerStore.mode === 'pomodoro'" class="mt-6 text-center text-sm text-gray-500 dark:text-gray-400">
      番茄时长: {{ settingsStore.settings.pomodoroDuration }} 分钟
      <span v-if="settingsStore.settings.soundEnabled" class="ml-2">🔊</span>
      <span v-else class="ml-2">🔇</span>
    </div>
  </div>
</template>
