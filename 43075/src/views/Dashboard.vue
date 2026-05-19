<script setup lang="ts">
import FocusTimer from '@/components/FocusTimer.vue'
import StatsPanel from '@/components/StatsPanel.vue'
import CategorySelector from '@/components/CategorySelector.vue'
</script>

<template>
  <div class="space-y-6 lg:space-y-8">
    <div>
      <h1 class="text-2xl lg:text-3xl font-bold text-gray-900 dark:text-white mb-2">
        开始专注
      </h1>
      <p class="text-gray-500 dark:text-gray-400">
        选择分类，开始你的专注之旅
      </p>
    </div>

    <div class="grid grid-cols-1 xl:grid-cols-3 gap-6">
      <div class="xl:col-span-2">
        <FocusTimer />
      </div>
      <div class="space-y-6">
        <CategorySelector />
      </div>
    </div>

    <StatsPanel />
  </div>
</template>
