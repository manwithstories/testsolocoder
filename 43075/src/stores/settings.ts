import { defineStore } from 'pinia'
import { ref, watch } from 'vue'
import type { Settings } from '@/types'
import { DEFAULT_SETTINGS } from '@/types'
import { storage } from '@/utils/storage'
import { logger } from '@/utils/logger'
import { useToastStore } from './toast'

export const useSettingsStore = defineStore('settings', () => {
  const settings = ref<Settings>({ ...DEFAULT_SETTINGS })
  const toast = useToastStore()

  const loadSettings = () => {
    try {
      const saved = storage.getSettings()
      settings.value = saved
      applyDarkMode(saved.darkMode)
      logger.info('加载设置成功', saved)
    } catch (error) {
      logger.error('加载设置失败', { error })
      toast.error('加载设置失败，使用默认配置')
    }
  }

  const saveSettings = (newSettings: Partial<Settings>) => {
    const updated = { ...settings.value, ...newSettings }
    if (storage.saveSettings(updated)) {
      settings.value = updated
      if (newSettings.darkMode !== undefined) {
        applyDarkMode(newSettings.darkMode)
      }
      logger.info('更新设置', newSettings)
      toast.success('设置已保存')
      return true
    }
    toast.error('保存设置失败')
    return false
  }

  const applyDarkMode = (dark: boolean) => {
    if (dark) {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
  }

  const resetSettings = () => {
    if (storage.saveSettings(DEFAULT_SETTINGS)) {
      settings.value = { ...DEFAULT_SETTINGS }
      applyDarkMode(DEFAULT_SETTINGS.darkMode)
      logger.info('重置设置为默认')
      toast.success('已重置为默认设置')
      return true
    }
    toast.error('重置设置失败')
    return false
  }

  const toggleDarkMode = () => {
    return saveSettings({ darkMode: !settings.value.darkMode })
  }

  const toggleSound = () => {
    return saveSettings({ soundEnabled: !settings.value.soundEnabled })
  }

  watch(
    () => settings.value.darkMode,
    (newVal) => {
      applyDarkMode(newVal)
    }
  )

  return {
    settings,
    loadSettings,
    saveSettings,
    resetSettings,
    toggleDarkMode,
    toggleSound
  }
})
