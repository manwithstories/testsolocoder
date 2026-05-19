<script setup lang="ts">
import { ref } from 'vue'
import TrendChart from '@/components/TrendChart.vue'
import PieChart from '@/components/PieChart.vue'
import HeatmapChart from '@/components/HeatmapChart.vue'

const trendType = ref<'week' | 'month'>('week')
const pieRange = ref<'week' | 'month' | 'all'>('week')
</script>

<template>
  <div class="space-y-6">
    <div>
      <h1 class="text-2xl lg:text-3xl font-bold text-gray-900 dark:text-white mb-2">
        数据统计
      </h1>
      <p class="text-gray-500 dark:text-gray-400">
        查看你的专注趋势和分析
      </p>
    </div>

    <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-lg p-6">
      <div class="flex items-center justify-between mb-6">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white">专注趋势</h3>
        <div class="flex gap-2">
          <button
            @click="trendType = 'week'"
            :class="[
              'px-4 py-2 rounded-lg text-sm font-medium transition-colors',
              trendType === 'week'
                ? 'bg-indigo-500 text-white'
                : 'bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600'
            ]"
          >
            本周
          </button>
          <button
            @click="trendType = 'month'"
            :class="[
              'px-4 py-2 rounded-lg text-sm font-medium transition-colors',
              trendType === 'month'
                ? 'bg-indigo-500 text-white'
                : 'bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600'
            ]"
          >
            本月
          </button>
        </div>
      </div>
      <TrendChart :type="trendType" />
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-lg p-6">
        <div class="flex items-center justify-between mb-6">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white">分类占比</h3>
          <div class="flex gap-2">
            <button
              @click="pieRange = 'week'"
              :class="[
                'px-3 py-1.5 rounded-lg text-xs font-medium transition-colors',
                pieRange === 'week'
                  ? 'bg-indigo-500 text-white'
                  : 'bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600'
              ]"
            >
              本周
            </button>
            <button
              @click="pieRange = 'month'"
              :class="[
                'px-3 py-1.5 rounded-lg text-xs font-medium transition-colors',
                pieRange === 'month'
                  ? 'bg-indigo-500 text-white'
                  : 'bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600'
              ]"
            >
              本月
            </button>
            <button
              @click="pieRange = 'all'"
              :class="[
                'px-3 py-1.5 rounded-lg text-xs font-medium transition-colors',
                pieRange === 'all'
                  ? 'bg-indigo-500 text-white'
                  : 'bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600'
              ]"
            >
              全部
            </button>
          </div>
        </div>
        <PieChart :range="pieRange" />
      </div>

      <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-lg p-6">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-6">时段分布热力图</h3>
        <p class="text-sm text-gray-500 dark:text-gray-400 mb-4">
          显示近4周的每小时专注时长分布
        </p>
        <HeatmapChart />
      </div>
    </div>
  </div>
</template>
