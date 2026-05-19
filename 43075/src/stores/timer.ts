import { defineStore } from 'pinia'
import { ref, computed, watch, onUnmounted } from 'vue'
import type { TimerMode, TimerStatus } from '@/types'
import { storage } from '@/utils/storage'
import { logger } from '@/utils/logger'
import { useToastStore } from './toast'
import { useSettingsStore } from './settings'
import { useCategoriesStore } from './categories'
import { useRecordsStore } from './records'

export const useTimerStore = defineStore('timer', () => {
  const status = ref<TimerStatus>('idle')
  const mode = ref<TimerMode>('pomodoro')
  const remainingTime = ref(0)
  const totalTime = ref(0)
  const customDuration = ref(25)
  const pomodoroCount = ref(0)
  const startTime = ref<number | null>(null)
  const pausedAt = ref<number | null>(null)
  const accumulatedTime = ref(0)

  const toast = useToastStore()
  const settingsStore = useSettingsStore()
  const categoriesStore = useCategoriesStore()
  const recordsStore = useRecordsStore()

  let intervalId: number | null = null

  const progress = computed(() => {
    if (totalTime.value === 0) return 0
    return ((totalTime.value - remainingTime.value) / totalTime.value) * 100
  })

  const isRunning = computed(() => status.value === 'running')
  const isPaused = computed(() => status.value === 'paused')
  const isIdle = computed(() => status.value === 'idle')

  const getDefaultDuration = (): number => {
    if (mode.value === 'pomodoro') {
      return settingsStore.settings.pomodoroDuration * 60
    }
    return customDuration.value * 60
  }

  const setMode = (newMode: TimerMode) => {
    if (status.value !== 'idle') {
      toast.warning('请先结束当前计时')
      return
    }
    mode.value = newMode
    remainingTime.value = getDefaultDuration()
    totalTime.value = remainingTime.value
    logger.info('切换计时模式', { mode: newMode })
  }

  const setCustomDuration = (minutes: number) => {
    if (minutes < 1 || minutes > 180) {
      toast.error('时长需在1-180分钟之间')
      return
    }
    customDuration.value = minutes
    if (mode.value === 'custom' && status.value === 'idle') {
      remainingTime.value = minutes * 60
      totalTime.value = remainingTime.value
    }
  }

  const start = () => {
    if (status.value === 'running') return

    if (!categoriesStore.selectedCategoryId) {
      toast.error('请先选择一个分类')
      return
    }

    if (status.value === 'idle') {
      totalTime.value = getDefaultDuration()
      remainingTime.value = totalTime.value
      startTime.value = Date.now()
      accumulatedTime.value = 0
      logger.info('开始新计时', {
        mode: mode.value,
        duration: totalTime.value,
        category: categoriesStore.selectedCategory?.name
      })
    } else if (status.value === 'paused' && pausedAt.value && startTime.value) {
      const pauseDuration = Date.now() - pausedAt.value
      startTime.value += pauseDuration
      logger.info('继续计时', { pausedFor: pauseDuration })
    }

    status.value = 'running'
    pausedAt.value = null

    intervalId = window.setInterval(() => {
      if (remainingTime.value > 0) {
        remainingTime.value--
      } else {
        complete()
      }
    }, 1000)
  }

  const pause = () => {
    if (status.value !== 'running') return

    status.value = 'paused'
    pausedAt.value = Date.now()

    if (intervalId) {
      clearInterval(intervalId)
      intervalId = null
    }

    accumulatedTime.value = totalTime.value - remainingTime.value
    logger.info('暂停计时', { accumulatedTime: accumulatedTime.value })
  }

  const resume = () => {
    if (status.value !== 'paused') return
    start()
  }

  const stop = () => {
    if (status.value === 'idle') return

    if (intervalId) {
      clearInterval(intervalId)
      intervalId = null
    }

    const actualDuration = totalTime.value - remainingTime.value

    if (actualDuration >= 60) {
      saveRecord(actualDuration, false)
      toast.info(`已记录 ${Math.floor(actualDuration / 60)} 分钟的专注时间`)
    } else {
      toast.info('专注时间不足1分钟，未记录')
    }

    reset()
  }

  const complete = () => {
    if (intervalId) {
      clearInterval(intervalId)
      intervalId = null
    }

    if (mode.value === 'pomodoro') {
      pomodoroCount.value++
    }

    saveRecord(totalTime.value, true)

    if (settingsStore.settings.soundEnabled) {
      playCompletionSound()
    }

    toast.success('专注完成！干得漂亮 👏')
    logger.info('计时完成', { duration: totalTime.value, mode: mode.value })

    reset()
  }

  const saveRecord = (duration: number, completed: boolean) => {
    if (!categoriesStore.selectedCategoryId || !startTime.value) return

    const record = {
      categoryId: categoriesStore.selectedCategoryId,
      mode: mode.value,
      duration,
      plannedDuration: totalTime.value,
      startTime: startTime.value,
      endTime: Date.now(),
      completed
    }

    recordsStore.addRecord(record)
  }

  const playCompletionSound = () => {
    try {
      const audioContext = new (window.AudioContext || (window as any).webkitAudioContext)()
      const oscillator = audioContext.createOscillator()
      const gainNode = audioContext.createGain()

      oscillator.connect(gainNode)
      gainNode.connect(audioContext.destination)

      oscillator.frequency.value = 800
      oscillator.type = 'sine'

      gainNode.gain.setValueAtTime(0.3, audioContext.currentTime)
      gainNode.gain.exponentialRampToValueAtTime(0.01, audioContext.currentTime + 0.5)

      oscillator.start(audioContext.currentTime)
      oscillator.stop(audioContext.currentTime + 0.5)
    } catch (error) {
      logger.warn('播放提示音失败', { error })
    }
  }

  const reset = () => {
    status.value = 'idle'
    remainingTime.value = getDefaultDuration()
    totalTime.value = remainingTime.value
    startTime.value = null
    pausedAt.value = null
    accumulatedTime.value = 0
  }

  watch(
    () => settingsStore.settings.pomodoroDuration,
    () => {
      if (mode.value === 'pomodoro' && status.value === 'idle') {
        remainingTime.value = getDefaultDuration()
        totalTime.value = remainingTime.value
      }
    }
  )

  onUnmounted(() => {
    if (intervalId) {
      clearInterval(intervalId)
    }
  })

  return {
    status,
    mode,
    remainingTime,
    totalTime,
    customDuration,
    pomodoroCount,
    progress,
    isRunning,
    isPaused,
    isIdle,
    setMode,
    setCustomDuration,
    start,
    pause,
    resume,
    stop,
    reset
  }
})
