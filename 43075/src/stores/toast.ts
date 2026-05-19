import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { ToastMessage } from '@/types'

export const useToastStore = defineStore('toast', () => {
  const messages = ref<ToastMessage[]>([])

  const showMessage = (type: ToastMessage['type'], message: string, duration: number = 3000) => {
    const id = Date.now().toString(36) + Math.random().toString(36).substr(2, 9)
    const toast: ToastMessage = { id, type, message, duration }
    messages.value.push(toast)

    setTimeout(() => {
      removeMessage(id)
    }, duration)

    return id
  }

  const removeMessage = (id: string) => {
    const index = messages.value.findIndex(m => m.id === id)
    if (index > -1) {
      messages.value.splice(index, 1)
    }
  }

  const success = (message: string, duration?: number) => showMessage('success', message, duration)
  const error = (message: string, duration?: number) => showMessage('error', message, duration)
  const warning = (message: string, duration?: number) => showMessage('warning', message, duration)
  const info = (message: string, duration?: number) => showMessage('info', message, duration)

  return {
    messages,
    showMessage,
    removeMessage,
    success,
    error,
    warning,
    info
  }
})
