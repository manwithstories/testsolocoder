<script setup lang="ts">
import { computed } from 'vue'
import { Trash2, Clock, CheckCircle, XCircle } from 'lucide-vue-next'
import type { FocusRecord } from '@/types'
import { useCategoriesStore } from '@/stores/categories'
import { formatDate, formatDuration } from '@/utils/date'
import * as LucideIcons from 'lucide-vue-next'

interface Props {
  record: FocusRecord
}

const props = defineProps<Props>()
const emit = defineEmits<{
  delete: [id: string]
}>()

const categoriesStore = useCategoriesStore()

const category = computed(() => {
  return categoriesStore.getCategoryById(props.record.categoryId)
})

const getIcon = (iconName: string) => {
  return (LucideIcons as any)[iconName] || LucideIcons.Circle
}

const timeRange = computed(() => {
  const start = formatDate(props.record.startTime, 'HH:mm')
  const end = formatDate(props.record.endTime, 'HH:mm')
  return `${start} - ${end}`
})
</script>

<template>
  <div class="flex items-center gap-4 p-4 bg-gray-50 dark:bg-gray-700/50 rounded-xl group hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors">
    <div
      class="w-10 h-10 rounded-lg flex items-center justify-center flex-shrink-0"
      :style="{ backgroundColor: category?.color + '20' }"
    >
      <component
        :is="getIcon(category?.icon || 'Circle')"
        class="w-5 h-5"
        :style="{ color: category?.color }"
      />
    </div>
    <div class="flex-1 min-w-0">
      <div class="flex items-center gap-2">
        <span class="font-medium text-gray-900 dark:text-white truncate">
          {{ category?.name || '未知分类' }}
        </span>
        <span
          :class="[
            'text-xs px-2 py-0.5 rounded-full',
            record.mode === 'pomodoro'
              ? 'bg-indigo-100 dark:bg-indigo-900/30 text-indigo-600 dark:text-indigo-400'
              : 'bg-gray-200 dark:bg-gray-600 text-gray-600 dark:text-gray-300'
          ]"
        >
          {{ record.mode === 'pomodoro' ? '番茄' : '自定义' }}
        </span>
        <CheckCircle
          v-if="record.completed"
          class="w-4 h-4 text-green-500 flex-shrink-0"
        />
        <XCircle
          v-else
          class="w-4 h-4 text-amber-500 flex-shrink-0"
        />
      </div>
      <div class="flex items-center gap-3 mt-1 text-sm text-gray-500 dark:text-gray-400">
        <span class="flex items-center gap-1">
          <Clock class="w-3.5 h-3.5" />
          {{ formatDuration(record.duration) }}
        </span>
        <span>{{ timeRange }}</span>
      </div>
    </div>
    <button
      @click="emit('delete', record.id)"
      class="p-2 text-gray-400 hover:text-red-500 hover:bg-red-50 dark:hover:bg-red-900/30 rounded-lg transition-colors opacity-0 group-hover:opacity-100"
    >
      <Trash2 class="w-4 h-4" />
    </button>
  </div>
</template>
