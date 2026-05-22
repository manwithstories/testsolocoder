import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useNoticeStore = defineStore('notice', () => {
  const unreadCount = ref(0)
  return { unreadCount }
})
