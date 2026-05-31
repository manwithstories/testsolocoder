import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAppStore = defineStore('app', () => {
  const sidebarCollapsed = ref(false)
  const loading = ref(false)
  const languagePairs = ref<any[]>([])
  const expertiseTags = ref<any[]>([])

  function toggleSidebar() {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }

  function setLanguagePairs(pairs: any[]) {
    languagePairs.value = pairs
  }

  function setExpertiseTags(tags: any[]) {
    expertiseTags.value = tags
  }

  return {
    sidebarCollapsed,
    loading,
    languagePairs,
    expertiseTags,
    toggleSidebar,
    setLanguagePairs,
    setExpertiseTags
  }
})
