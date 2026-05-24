import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAppStore = defineStore('app', () => {
  // 侧边栏是否折叠
  const sidebarCollapsed = ref<boolean>(false)
  // 加载状态
  const loading = ref<boolean>(false)

  // 切换侧边栏折叠状态
  function toggleSidebar(): void {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }

  // 设置侧边栏状态
  function setSidebarCollapsed(collapsed: boolean): void {
    sidebarCollapsed.value = collapsed
  }

  // 显示加载状态
  function showLoading(): void {
    loading.value = true
  }

  // 隐藏加载状态
  function hideLoading(): void {
    loading.value = false
  }

  return {
    sidebarCollapsed,
    loading,
    toggleSidebar,
    setSidebarCollapsed,
    showLoading,
    hideLoading
  }
})