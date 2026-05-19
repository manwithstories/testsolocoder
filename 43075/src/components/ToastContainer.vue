<script setup lang="ts">
import { useToastStore } from '@/stores/toast'
import { CheckCircle, XCircle, AlertTriangle, Info, X } from 'lucide-vue-next'

const toastStore = useToastStore()

const iconMap = {
  success: CheckCircle,
  error: XCircle,
  warning: AlertTriangle,
  info: Info
}

const colorMap = {
  success: 'bg-green-50 border-green-200 text-green-800 dark:bg-green-900/30 dark:border-green-800 dark:text-green-300',
  error: 'bg-red-50 border-red-200 text-red-800 dark:bg-red-900/30 dark:border-red-800 dark:text-red-300',
  warning: 'bg-amber-50 border-amber-200 text-amber-800 dark:bg-amber-900/30 dark:border-amber-800 dark:text-amber-300',
  info: 'bg-blue-50 border-blue-200 text-blue-800 dark:bg-blue-900/30 dark:border-blue-800 dark:text-blue-300'
}

const iconColorMap = {
  success: 'text-green-500',
  error: 'text-red-500',
  warning: 'text-amber-500',
  info: 'text-blue-500'
}
</script>

<template>
  <div class="fixed top-4 right-4 z-[100] space-y-2 w-80 max-w-[calc(100vw-2rem)]">
    <TransitionGroup name="toast">
      <div
        v-for="toast in toastStore.messages"
        :key="toast.id"
        :class="[
          'flex items-start gap-3 p-4 rounded-lg border shadow-lg',
          colorMap[toast.type]
        ]"
      >
        <component
          :is="iconMap[toast.type]"
          :class="['w-5 h-5 flex-shrink-0 mt-0.5', iconColorMap[toast.type]]"
        />
        <p class="flex-1 text-sm font-medium">{{ toast.message }}</p>
        <button
          @click="toastStore.removeMessage(toast.id)"
          class="flex-shrink-0 p-1 hover:bg-black/5 dark:hover:bg-white/10 rounded transition-colors"
        >
          <X class="w-4 h-4" />
        </button>
      </div>
    </TransitionGroup>
  </div>
</template>

<style scoped>
.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}

.toast-enter-from {
  opacity: 0;
  transform: translateX(100%);
}

.toast-leave-to {
  opacity: 0;
  transform: translateX(100%);
}

.toast-move {
  transition: transform 0.3s ease;
}
</style>
