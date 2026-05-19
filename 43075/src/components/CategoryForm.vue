<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { X } from 'lucide-vue-next'
import type { Category } from '@/types'
import { COLOR_OPTIONS, ICON_OPTIONS } from '@/types'
import * as LucideIcons from 'lucide-vue-next'

interface Props {
  visible: boolean
  category?: Category | null
}

const props = withDefaults(defineProps<Props>(), {
  category: null
})

const emit = defineEmits<{
  submit: [data: { name: string; color: string; icon: string }]
  cancel: []
}>()

const name = ref('')
const color = ref(COLOR_OPTIONS[0])
const icon = ref(ICON_OPTIONS[0])

const isEdit = computed(() => !!props.category)

const getIcon = (iconName: string) => {
  return (LucideIcons as any)[iconName] || LucideIcons.Circle
}

watch(
  () => props.visible,
  (newVal) => {
    if (newVal && props.category) {
      name.value = props.category.name
      color.value = props.category.color
      icon.value = props.category.icon
    } else if (newVal) {
      name.value = ''
      color.value = COLOR_OPTIONS[0]
      icon.value = ICON_OPTIONS[0]
    }
  }
)

const handleSubmit = () => {
  if (!name.value.trim()) return
  emit('submit', {
    name: name.value.trim(),
    color: color.value,
    icon: icon.value
  })
}
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div
        v-if="visible"
        class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50"
        @click.self="emit('cancel')"
      >
        <div class="bg-white dark:bg-gray-800 rounded-2xl w-full max-w-md p-6 shadow-2xl">
          <div class="flex items-center justify-between mb-6">
            <h3 class="text-xl font-bold text-gray-900 dark:text-white">
              {{ isEdit ? '编辑分类' : '新建分类' }}
            </h3>
            <button
              @click="emit('cancel')"
              class="p-2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-200 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
            >
              <X class="w-5 h-5" />
            </button>
          </div>

          <div class="space-y-6">
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                分类名称
              </label>
              <input
                v-model="name"
                type="text"
                maxlength="20"
                placeholder="输入分类名称"
                class="w-full px-4 py-3 rounded-xl border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
              />
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                选择颜色
              </label>
              <div class="flex flex-wrap gap-2">
                <button
                  v-for="c in COLOR_OPTIONS"
                  :key="c"
                  @click="color = c"
                  :class="[
                    'w-10 h-10 rounded-full transition-all',
                    color === c ? 'ring-2 ring-offset-2 ring-gray-400 scale-110' : 'hover:scale-105'
                  ]"
                  :style="{ backgroundColor: c }"
                />
              </div>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                选择图标
              </label>
              <div class="grid grid-cols-6 gap-2">
                <button
                  v-for="i in ICON_OPTIONS"
                  :key="i"
                  @click="icon = i"
                  :class="[
                    'w-10 h-10 rounded-lg flex items-center justify-center transition-all',
                    icon === i
                      ? 'bg-indigo-100 dark:bg-indigo-900/30 text-indigo-600 dark:text-indigo-400'
                      : 'bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600'
                  ]"
                >
                  <component :is="getIcon(i)" class="w-5 h-5" />
                </button>
              </div>
            </div>

            <div class="flex items-center gap-3 p-4 rounded-xl" :style="{ backgroundColor: color + '15' }">
              <div
                class="w-12 h-12 rounded-xl flex items-center justify-center"
                :style="{ backgroundColor: color + '30' }"
              >
                <component :is="getIcon(icon)" class="w-6 h-6" :style="{ color }" />
              </div>
              <span class="font-medium text-gray-900 dark:text-white">
                {{ name || '分类名称' }}
              </span>
            </div>
          </div>

          <div class="flex gap-3 mt-8">
            <button
              @click="emit('cancel')"
              class="flex-1 px-4 py-3 rounded-xl border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 font-medium hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
            >
              取消
            </button>
            <button
              @click="handleSubmit"
              :disabled="!name.trim()"
              class="flex-1 px-4 py-3 rounded-xl bg-indigo-500 text-white font-medium hover:bg-indigo-600 disabled:bg-gray-300 dark:disabled:bg-gray-600 disabled:cursor-not-allowed transition-colors"
            >
              {{ isEdit ? '保存' : '创建' }}
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.2s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-active .bg-white,
.modal-leave-active .bg-white {
  transition: transform 0.2s ease;
}

.modal-enter-from .bg-white,
.modal-leave-to .bg-white {
  transform: scale(0.95);
}
</style>
