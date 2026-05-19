<script setup lang="ts">
import { computed } from 'vue'
import { Clock, Target, Zap, Flame } from 'lucide-vue-next'
import { useRecordsStore } from '@/stores/records'
import { formatDuration, formatDurationShort } from '@/utils/date'
import StatsCard from './StatsCard.vue'

const recordsStore = useRecordsStore()

const todayTotalText = computed(() => {
  return formatDurationShort(recordsStore.todayTotalDuration)
})

const todayLongestText = computed(() => {
  if (recordsStore.todayLongestFocus === 0) return '0m'
  return formatDurationShort(recordsStore.todayLongestFocus)
})

const streakText = computed(() => {
  return `${recordsStore.currentStreak} 天`
})
</script>

<template>
  <div class="space-y-4">
    <h3 class="text-lg font-semibold text-gray-900 dark:text-white">今日统计</h3>
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
      <StatsCard
        :icon="Clock"
        title="专注时长"
        :value="todayTotalText"
        :subtitle="formatDuration(recordsStore.todayTotalDuration)"
        color="#6366f1"
      />
      <StatsCard
        :icon="Target"
        title="完成番茄"
        :value="recordsStore.todayCompletedPomodoros"
        subtitle="个"
        color="#10b981"
      />
      <StatsCard
        :icon="Zap"
        title="最长专注"
        :value="todayLongestText"
        :subtitle="recordsStore.todayLongestFocus > 0 ? formatDuration(recordsStore.todayLongestFocus) : '暂无记录'"
        color="#f59e0b"
      />
      <StatsCard
        :icon="Flame"
        title="连续天数"
        :value="streakText"
        :subtitle="recordsStore.currentStreak > 0 ? '继续保持！' : '开始你的第一天'"
        color="#ef4444"
      />
    </div>
  </div>
</template>
