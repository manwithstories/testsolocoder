<script setup lang="ts">
import { ref, watch } from 'vue'
import { Download, Upload, RotateCcw, Trash2, Volume2, VolumeX, Sun, Moon } from 'lucide-vue-next'
import { useSettingsStore } from '@/stores/settings'
import { useToastStore } from '@/stores/toast'
import { storage } from '@/utils/storage'
import { logger } from '@/utils/logger'
import { validateNumberRange } from '@/utils/validator'

const settingsStore = useSettingsStore()
const toastStore = useToastStore()

const localSettings = ref({ ...settingsStore.settings })
const showResetConfirm = ref(false)
const showClearConfirm = ref(false)

watch(
  () => settingsStore.settings,
  (newVal) => {
    localSettings.value = { ...newVal }
  },
  { deep: true }
)

const updateSetting = <K extends keyof typeof localSettings.value>(
  key: K,
  value: typeof localSettings.value[K]
) => {
  settingsStore.saveSettings({ [key]: value })
}

const handlePomodoroDurationChange = (e: Event) => {
  const value = parseInt((e.target as HTMLInputElement).value)
  const validation = validateNumberRange(value, 1, 120)
  if (validation.valid) {
    updateSetting('pomodoroDuration', value)
  } else {
    toastStore.error(validation.error || '时长无效')
  }
}

const handleShortBreakChange = (e: Event) => {
  const value = parseInt((e.target as HTMLInputElement).value)
  const validation = validateNumberRange(value, 1, 60)
  if (validation.valid) {
    updateSetting('shortBreakDuration', value)
  } else {
    toastStore.error(validation.error || '时长无效')
  }
}

const handleLongBreakChange = (e: Event) => {
  const value = parseInt((e.target as HTMLInputElement).value)
  const validation = validateNumberRange(value, 1, 120)
  if (validation.valid) {
    updateSetting('longBreakDuration', value)
  } else {
    toastStore.error(validation.error || '时长无效')
  }
}

const handleLongBreakIntervalChange = (e: Event) => {
  const value = parseInt((e.target as HTMLInputElement).value)
  const validation = validateNumberRange(value, 1, 10)
  if (validation.valid) {
    updateSetting('longBreakInterval', value)
  } else {
    toastStore.error(validation.error || '数量无效')
  }
}

const exportData = () => {
  const data = storage.exportData()
  if (!data) {
    toastStore.error('导出数据失败')
    return
  }

  const blob = new Blob([data], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `focus-backup-${new Date().toISOString().split('T')[0]}.json`
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)

  toastStore.success('数据导出成功')
  logger.info('导出数据备份')
}

const importData = () => {
  const input = document.createElement('input')
  input.type = 'file'
  input.accept = '.json'

  input.onchange = (e) => {
    const file = (e.target as HTMLInputElement).files?.[0]
    if (!file) return

    const reader = new FileReader()
    reader.onload = (event) => {
      const content = event.target?.result as string
      const result = storage.importData(content)

      if (result.success) {
        toastStore.success('数据导入成功，页面即将刷新')
        setTimeout(() => window.location.reload(), 1500)
      } else {
        toastStore.error(result.error || '导入失败')
      }
    }
    reader.onerror = () => {
      toastStore.error('读取文件失败')
    }
    reader.readAsText(file)
  }

  input.click()
}

const resetSettings = () => {
  settingsStore.resetSettings()
  showResetConfirm.value = false
}

const clearAllData = () => {
  if (storage.clearAll()) {
    toastStore.success('数据已清空，页面即将刷新')
    setTimeout(() => window.location.reload(), 1500)
  } else {
    toastStore.error('清空数据失败')
  }
  showClearConfirm.value = false
}
</script>

<template>
  <div class="space-y-6">
    <div>
      <h1 class="text-2xl lg:text-3xl font-bold text-gray-900 dark:text-white mb-2">
        应用设置
      </h1>
      <p class="text-gray-500 dark:text-gray-400">
        配置你的专注体验
      </p>
    </div>

    <div class="space-y-6">
      <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-lg p-6">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-6">时间设置</h3>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              番茄时长（分钟）
            </label>
            <input
              type="number"
              :value="localSettings.pomodoroDuration"
              @change="handlePomodoroDurationChange"
              min="1"
              max="120"
              class="w-full px-4 py-3 rounded-xl border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-indigo-500"
            />
            <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">默认 25 分钟</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              短休息时长（分钟）
            </label>
            <input
              type="number"
              :value="localSettings.shortBreakDuration"
              @change="handleShortBreakChange"
              min="1"
              max="60"
              class="w-full px-4 py-3 rounded-xl border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-indigo-500"
            />
            <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">默认 5 分钟</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              长休息时长（分钟）
            </label>
            <input
              type="number"
              :value="localSettings.longBreakDuration"
              @change="handleLongBreakChange"
              min="1"
              max="120"
              class="w-full px-4 py-3 rounded-xl border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-indigo-500"
            />
            <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">默认 15 分钟</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              长休息间隔（番茄数）
            </label>
            <input
              type="number"
              :value="localSettings.longBreakInterval"
              @change="handleLongBreakIntervalChange"
              min="1"
              max="10"
              class="w-full px-4 py-3 rounded-xl border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-indigo-500"
            />
            <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">每完成几个番茄后长休息</p>
          </div>
        </div>
      </div>

      <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-lg p-6">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-6">界面与提醒</h3>
        <div class="space-y-4">
          <div class="flex items-center justify-between py-3">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-xl bg-indigo-100 dark:bg-indigo-900/30 flex items-center justify-center">
                <Sun v-if="!localSettings.darkMode" class="w-5 h-5 text-amber-500" />
                <Moon v-else class="w-5 h-5 text-indigo-500" />
              </div>
              <div>
                <p class="font-medium text-gray-900 dark:text-white">深色模式</p>
                <p class="text-sm text-gray-500 dark:text-gray-400">切换应用主题</p>
              </div>
            </div>
            <button
              @click="updateSetting('darkMode', !localSettings.darkMode)"
              :class="[
                'w-14 h-8 rounded-full transition-colors relative',
                localSettings.darkMode ? 'bg-indigo-500' : 'bg-gray-300 dark:bg-gray-600'
              ]"
            >
              <div
                :class="[
                  'absolute top-1 w-6 h-6 bg-white rounded-full shadow transition-transform',
                  localSettings.darkMode ? 'translate-x-7' : 'translate-x-1'
                ]"
              />
            </button>
          </div>

          <div class="flex items-center justify-between py-3">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-xl bg-green-100 dark:bg-green-900/30 flex items-center justify-center">
                <Volume2 v-if="localSettings.soundEnabled" class="w-5 h-5 text-green-500" />
                <VolumeX v-else class="w-5 h-5 text-gray-400" />
              </div>
              <div>
                <p class="font-medium text-gray-900 dark:text-white">声音提醒</p>
                <p class="text-sm text-gray-500 dark:text-gray-400">计时结束时播放提示音</p>
              </div>
            </div>
            <button
              @click="updateSetting('soundEnabled', !localSettings.soundEnabled)"
              :class="[
                'w-14 h-8 rounded-full transition-colors relative',
                localSettings.soundEnabled ? 'bg-green-500' : 'bg-gray-300 dark:bg-gray-600'
              ]"
            >
              <div
                :class="[
                  'absolute top-1 w-6 h-6 bg-white rounded-full shadow transition-transform',
                  localSettings.soundEnabled ? 'translate-x-7' : 'translate-x-1'
                ]"
              />
            </button>
          </div>
        </div>
      </div>

      <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-lg p-6">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-6">数据管理</h3>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <button
            @click="exportData"
            class="flex items-center justify-center gap-3 p-4 border-2 border-dashed border-gray-300 dark:border-gray-600 rounded-xl hover:border-indigo-500 hover:bg-indigo-50 dark:hover:bg-indigo-900/20 transition-all group"
          >
            <Download class="w-6 h-6 text-gray-400 group-hover:text-indigo-500" />
            <div class="text-left">
              <p class="font-medium text-gray-900 dark:text-white">导出数据</p>
              <p class="text-sm text-gray-500 dark:text-gray-400">备份所有数据到JSON文件</p>
            </div>
          </button>

          <button
            @click="importData"
            class="flex items-center justify-center gap-3 p-4 border-2 border-dashed border-gray-300 dark:border-gray-600 rounded-xl hover:border-indigo-500 hover:bg-indigo-50 dark:hover:bg-indigo-900/20 transition-all group"
          >
            <Upload class="w-6 h-6 text-gray-400 group-hover:text-indigo-500" />
            <div class="text-left">
              <p class="font-medium text-gray-900 dark:text-white">导入数据</p>
              <p class="text-sm text-gray-500 dark:text-gray-400">从备份文件恢复数据</p>
            </div>
          </button>

          <button
            @click="showResetConfirm = true"
            class="flex items-center justify-center gap-3 p-4 border-2 border-dashed border-gray-300 dark:border-gray-600 rounded-xl hover:border-amber-500 hover:bg-amber-50 dark:hover:bg-amber-900/20 transition-all group"
          >
            <RotateCcw class="w-6 h-6 text-gray-400 group-hover:text-amber-500" />
            <div class="text-left">
              <p class="font-medium text-gray-900 dark:text-white">重置设置</p>
              <p class="text-sm text-gray-500 dark:text-gray-400">恢复默认设置</p>
            </div>
          </button>

          <button
            @click="showClearConfirm = true"
            class="flex items-center justify-center gap-3 p-4 border-2 border-dashed border-gray-300 dark:border-gray-600 rounded-xl hover:border-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 transition-all group"
          >
            <Trash2 class="w-6 h-6 text-gray-400 group-hover:text-red-500" />
            <div class="text-left">
              <p class="font-medium text-gray-900 dark:text-white">清空所有数据</p>
              <p class="text-sm text-gray-500 dark:text-gray-400">删除所有记录和设置</p>
            </div>
          </button>
        </div>
      </div>

      <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-lg p-6">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">关于</h3>
        <div class="text-sm text-gray-500 dark:text-gray-400 space-y-2">
          <p>专注仪表盘 v1.0.0</p>
          <p>帮助你建立良好的专注习惯，提升工作效率。</p>
          <p>所有数据存储在本地浏览器中，请定期导出备份。</p>
        </div>
      </div>
    </div>

    <Teleport to="body">
      <Transition name="modal">
        <div
          v-if="showResetConfirm"
          class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50"
          @click.self="showResetConfirm = false"
        >
          <div class="bg-white dark:bg-gray-800 rounded-2xl w-full max-w-sm p-6 shadow-2xl">
            <h3 class="text-lg font-bold text-gray-900 dark:text-white mb-2">重置设置</h3>
            <p class="text-gray-500 dark:text-gray-400 mb-6">
              确定要重置所有设置为默认值吗？这不会影响你的专注记录。
            </p>
            <div class="flex gap-3">
              <button
                @click="showResetConfirm = false"
                class="flex-1 px-4 py-2.5 rounded-xl border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 font-medium hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
              >
                取消
              </button>
              <button
                @click="resetSettings"
                class="flex-1 px-4 py-2.5 rounded-xl bg-amber-500 text-white font-medium hover:bg-amber-600 transition-colors"
              >
                确认重置
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>

    <Teleport to="body">
      <Transition name="modal">
        <div
          v-if="showClearConfirm"
          class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50"
          @click.self="showClearConfirm = false"
        >
          <div class="bg-white dark:bg-gray-800 rounded-2xl w-full max-w-sm p-6 shadow-2xl">
            <h3 class="text-lg font-bold text-gray-900 dark:text-white mb-2">清空所有数据</h3>
            <p class="text-gray-500 dark:text-gray-400 mb-6">
              确定要清空所有数据吗？这将删除所有专注记录、分类和设置，且无法恢复。
            </p>
            <div class="flex gap-3">
              <button
                @click="showClearConfirm = false"
                class="flex-1 px-4 py-2.5 rounded-xl border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 font-medium hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
              >
                取消
              </button>
              <button
                @click="clearAllData"
                class="flex-1 px-4 py-2.5 rounded-xl bg-red-500 text-white font-medium hover:bg-red-600 transition-colors"
              >
                确认清空
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
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
</style>
