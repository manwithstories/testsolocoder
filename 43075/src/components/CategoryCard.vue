<script setup lang="ts">
import { Pencil, Trash2 } from 'lucide-vue-next'
import type { Category } from '@/types'
import * as LucideIcons from 'lucide-vue-next'

interface Props {
  category: Category
}

defineProps<Props>()
const emit = defineEmits<{
  edit: [category: Category]
  delete: [id: string]
}>()

const getIcon = (iconName: string) => {
  return (LucideIcons as any)[iconName] || LucideIcons.Circle
}
</script>

<template>
  <div class="bg-white dark:bg-gray-800 rounded-xl p-5 shadow-md hover:shadow-lg transition-all group">
    <div class="flex items-start justify-between">
      <div class="flex items-center gap-3">
        <div
          class="w-12 h-12 rounded-xl flex items-center justify-center"
          :style="{ backgroundColor: category.color + '20' }"
        >
          <component
            :is="getIcon(category.icon)"
            class="w-6 h-6"
            :style="{ color: category.color }"
          />
        </div>
        <div>
          <h4 class="font-semibold text-gray-900 dark:text-white">{{ category.name }}</h4>
          <p class="text-sm text-gray-500 dark:text-gray-400">{{ category.color }}</p>
        </div>
      </div>
      <div class="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
        <button
          @click="emit('edit', category)"
          class="p-2 text-gray-400 hover:text-indigo-500 hover:bg-indigo-50 dark:hover:bg-indigo-900/30 rounded-lg transition-colors"
        >
          <Pencil class="w-4 h-4" />
        </button>
        <button
          @click="emit('delete', category.id)"
          class="p-2 text-gray-400 hover:text-red-500 hover:bg-red-50 dark:hover:bg-red-900/30 rounded-lg transition-colors"
        >
          <Trash2 class="w-4 h-4" />
        </button>
      </div>
    </div>
  </div>
</template>
