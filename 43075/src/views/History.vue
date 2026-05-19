<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ChevronDown, ChevronRight, Filter, Calendar, Clock } from 'lucide-vue-next'
import { useRecordsStore } from '@/stores/records'
import { useCategoriesStore } from '@/stores/categories'
import { formatDate, formatDuration, getStartOfDay, getEndOfDay } from '@/utils/date'
import HistoryItem from '@/components/HistoryItem.vue'
import * as LucideIcons from 'lucide-vue-next'

const recordsStore = useRecordsStore()
const categoriesStore = useCategoriesStore()

const selectedCategoryId = ref<string | null>(null)
const startDate = ref('')
const endDate = ref('')
const expandedDates = ref<Set<string>>(new Set())

const getIcon = (iconName: string) => {
  return (LucideIcons as any)[iconName] || LucideIcons.Circle
}

const filteredRecords = computed(() => {
  let start = 0
  let end = Date.now()

  if (startDate.value) {
    start = new Date(startDate.value).getTime()
  }
  if (endDate.value) {
    end = getEndOfDay(new Date(endDate.value).getTime())
  }

  return recordsStore.getRecordsByDateRange(start, end, selectedCategoryId.value || undefined)
})

const groupedRecords = computed(() => {
  const groups = new Map<string, typeof filteredRecords.value>()

  filteredRecords.value.forEach(record => {
    const dateKey = formatDate(record.startTime, 'YYYY-MM-DD')
    if (!groups.has(dateKey)) {
      groups.set(dateKey, [])
    }
    groups.get(dateKey)!.push(record)
  })

  return Array.from(groups.entries())
    .sort((a, b) => b[0].localeCompare(a[0]))
    .map(([date, records]) => ({
      date,
      records,
      totalDuration: records.reduce((sum, r) => sum + r.duration, 0),
      totalSessions: records.length
    }))
})

const toggleDate = (date: string) => {
  if (expandedDates.value.has(date)) {
    expandedDates.value.delete(date)
  } else {
    expandedDates.value.add(date)
  }
}

const handleDelete = (recordId: string) => {
  if (confirm('确定要删除这条记录吗？')) {
    recordsStore.deleteRecord(recordId)
  }
}

const resetFilters = () => {
  selectedCategoryId.value = null
  startDate.value = ''
  endDate.value = ''
}
</script>

<template>
  <div class="space-y-6">
    <div>
      <h1 class="text-2xl lg:text-3xl font-bold text-gray-900 dark:text-white mb-2">
        历史记录
      </h1>
      <p class="text-gray-500 dark:text-gray-400">
        查看和管理你的专注历史
      </p>
    </div>

    <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-lg p-6">
      <div class="flex flex-col lg:flex-row lg:items-center gap-4 mb-6">
        <div class="flex items-center gap-2 text-gray-700 dark:text-gray-300">
          <Filter class="w-5 h-5" />
          <span class="font-medium">筛选条件</span>
        </div>
        <div class="flex flex-wrap gap-4">
          <select
            v-model="selectedCategoryId"
            class="px-4 py-2 rounded-xl border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-indigo-500"
          >
            <option :value="null">全部分类</option>
            <option
              v-for="cat in categoriesStore.categories"
              :key="cat.id"
              :value="cat.id"
            >
              {{ cat.name }}
            </option>
          </select>
          <input
            v-model="startDate"
            type="date"
            class="px-4 py-2 rounded-xl border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-indigo-500"
          />
          <span class="text-gray-400 self-center">至</span>
          <input
            v-model="endDate"
            type="date"
            class="px-4 py-2 rounded-xl border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-indigo-500"
          />
          <button
            @click="resetFilters"
            class="px-4 py-2 text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-xl transition-colors"
          >
            重置
          </button>
        </div>
      </div>

      <div v-if="groupedRecords.length === 0" class="text-center py-16">
        <Clock class="w-16 h-16 text-gray-300 dark:text-gray-600 mx-auto mb-4" />
        <p class="text-gray-500 dark:text-gray-400">暂无专注记录</p>
        <p class="text-sm text-gray-400 dark:text-gray-500 mt-1">开始你的第一次专注吧</p>
      </div>

      <div v-else class="space-y-4">
        <div
          v-for="group in groupedRecords"
          :key="group.date"
          class="border border-gray-200 dark:border-gray-700 rounded-xl overflow-hidden"
        >
          <button
            @click="toggleDate(group.date)"
            class="w-full flex items-center justify-between p-4 bg-gray-50 dark:bg-gray-700/50 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
          >
            <div class="flex items-center gap-4">
              <div class="flex items-center gap-2">
                <Calendar class="w-5 h-5 text-indigo-500" />
                <span class="font-semibold text-gray-900 dark:text-white">
                  {{ group.date }}
                </span>
              </div>
              <span class="text-sm text-gray-500 dark:text-gray-400">
                {{ group.totalSessions }} 次专注 · {{ formatDuration(group.totalDuration) }}
              </span>
            </div>
            <ChevronDown
              v-if="expandedDates.has(group.date)"
              class="w-5 h-5 text-gray-400"
            />
            <ChevronRight
              v-else
              class="w-5 h-5 text-gray-400"
            />
          </button>

          <Transition name="expand">
            <div v-if="expandedDates.has(group.date)" class="p-4 space-y-3">
              <HistoryItem
                v-for="record in group.records"
                :key="record.id"
                :record="record"
                @delete="handleDelete"
              />
            </div>
          </Transition>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.expand-enter-active,
.expand-leave-active {
  transition: all 0.3s ease;
  overflow: hidden;
}

.expand-enter-from,
.expand-leave-to {
  opacity: 0;
  max-height: 0;
}

.expand-enter-to,
.expand-leave-from {
  opacity: 1;
  max-height: 1000px;
}
</style>
