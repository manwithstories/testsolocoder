import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

const SIDEBAR_COLLAPSED_KEY = 'hotel_sidebar_collapsed'
const THEME_KEY = 'hotel_theme'

export const useAppStore = defineStore('app', () => {
  const sidebarCollapsed = ref<boolean>(localStorage.getItem(SIDEBAR_COLLAPSED_KEY) === 'true')
  const theme = ref<string>(localStorage.getItem(THEME_KEY) || 'light')
  const loading = ref<boolean>(false)
  const locale = ref<string>('zh-CN')

  const sidebarWidth = computed(() => sidebarCollapsed.value ? '64px' : '220px')

  const toggleSidebar = () => {
    sidebarCollapsed.value = !sidebarCollapsed.value
    localStorage.setItem(SIDEBAR_COLLAPSED_KEY, String(sidebarCollapsed.value))
  }

  const setSidebarCollapsed = (collapsed: boolean) => {
    sidebarCollapsed.value = collapsed
    localStorage.setItem(SIDEBAR_COLLAPSED_KEY, String(collapsed))
  }

  const toggleTheme = () => {
    theme.value = theme.value === 'light' ? 'dark' : 'light'
    localStorage.setItem(THEME_KEY, theme.value)
    applyTheme(theme.value)
  }

  const setTheme = (newTheme: string) => {
    theme.value = newTheme
    localStorage.setItem(THEME_KEY, newTheme)
    applyTheme(newTheme)
  }

  const applyTheme = (newTheme: string) => {
    const html = document.documentElement
    if (newTheme === 'dark') {
      html.classList.add('dark')
    } else {
      html.classList.remove('dark')
    }
  }

  const setLoading = (isLoading: boolean) => {
    loading.value = isLoading
  }

  const setLocale = (newLocale: string) => {
    locale.value = newLocale
  }

  return {
    sidebarCollapsed,
    sidebarWidth,
    theme,
    loading,
    locale,
    toggleSidebar,
    setSidebarCollapsed,
    toggleTheme,
    setTheme,
    setLoading,
    setLocale
  }
})
