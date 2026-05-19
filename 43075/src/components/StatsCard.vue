<script setup lang="ts">
import { computed } from 'vue'
import type { Component } from 'vue'

interface Props {
  icon: Component
  title: string
  value: string | number
  subtitle?: string
  color?: string
}

const props = withDefaults(defineProps<Props>(), {
  color: '#6366f1'
})

const bgColor = computed(() => {
  const hex = props.color.replace('#', '')
  const r = parseInt(hex.substring(0, 2), 16)
  const g = parseInt(hex.substring(2, 4), 16)
  const b = parseInt(hex.substring(4, 6), 16)
  return `rgba(${r}, ${g}, ${b}, 0.1)`
})
</script>

<template>
  <div class="bg-white dark:bg-gray-800 rounded-xl p-5 shadow-md hover:shadow-lg transition-shadow">
    <div class="flex items-start justify-between">
      <div>
        <p class="text-sm text-gray-500 dark:text-gray-400 mb-1">{{ title }}</p>
        <p class="text-2xl lg:text-3xl font-bold text-gray-900 dark:text-white">{{ value }}</p>
        <p v-if="subtitle" class="text-xs text-gray-400 dark:text-gray-500 mt-1">{{ subtitle }}</p>
      </div>
      <div
        class="w-12 h-12 rounded-xl flex items-center justify-center"
        :style="{ backgroundColor: bgColor }"
      >
        <component :is="icon" class="w-6 h-6" :style="{ color }" />
      </div>
    </div>
  </div>
</template>
