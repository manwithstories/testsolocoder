import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { FocusRecord, DailyStats } from '@/types'
import { storage } from '@/utils/storage'
import { logger } from '@/utils/logger'
import { formatDate, getStartOfDay, getEndOfDay, isToday, isSameDay, getWeekDates, getMonthDates } from '@/utils/date'
import { useToastStore } from './toast'
import { useCategoriesStore } from './categories'

export const useRecordsStore = defineStore('records', () => {
  const records = ref<FocusRecord[]>([])
  const toast = useToastStore()
  const categoriesStore = useCategoriesStore()

  const loadRecords = () => {
    try {
      records.value = storage.getRecords()
      logger.info('加载记录成功', { count: records.value.length })
    } catch (error) {
      logger.error('加载记录失败', { error })
      toast.error('加载记录失败')
    }
  }

  const addRecord = (record: Omit<FocusRecord, 'id'>): FocusRecord | null => {
    const newRecord = storage.addRecord(record)
    if (newRecord) {
      records.value.push(newRecord)
      return newRecord
    }
    toast.error('保存专注记录失败')
    return null
  }

  const deleteRecord = (recordId: string): boolean => {
    if (storage.deleteRecord(recordId)) {
      records.value = records.value.filter(r => r.id !== recordId)
      return true
    }
    return false
  }

  const todayRecords = computed(() => {
    return records.value.filter(r => isToday(r.startTime))
  })

  const todayTotalDuration = computed(() => {
    return todayRecords.value.reduce((sum, r) => sum + r.duration, 0)
  })

  const todayCompletedPomodoros = computed(() => {
    return todayRecords.value.filter(r => r.mode === 'pomodoro' && r.completed).length
  })

  const todayLongestFocus = computed(() => {
    if (todayRecords.value.length === 0) return 0
    return Math.max(...todayRecords.value.map(r => r.duration))
  })

  const currentStreak = computed(() => {
    if (records.value.length === 0) return 0

    const dates = new Set<string>()
    records.value.forEach(r => {
      dates.add(formatDate(r.startTime, 'YYYY-MM-DD'))
    })

    const sortedDates = Array.from(dates).sort().reverse()
    let streak = 0
    let checkDate = new Date()
    checkDate.setHours(0, 0, 0, 0)

    for (const dateStr of sortedDates) {
      const recordDate = new Date(dateStr)
      recordDate.setHours(0, 0, 0, 0)

      const diffDays = Math.round((checkDate.getTime() - recordDate.getTime()) / (1000 * 60 * 60 * 24))

      if (diffDays === streak) {
        streak++
      } else if (diffDays === streak + 1 && streak === 0) {
        const yesterday = new Date()
        yesterday.setDate(yesterday.getDate() - 1)
        yesterday.setHours(0, 0, 0, 0)
        if (recordDate.getTime() === yesterday.getTime()) {
          streak++
        } else {
          break
        }
      } else {
        break
      }
    }

    return streak
  })

  const getRecordsByDateRange = (startDate: number, endDate: number, categoryId?: string): FocusRecord[] => {
    let filtered = records.value.filter(r => r.startTime >= startDate && r.startTime <= endDate)
    if (categoryId) {
      filtered = filtered.filter(r => r.categoryId === categoryId)
    }
    return filtered.sort((a, b) => b.startTime - a.startTime)
  }

  const getDailyStats = (startDate: number, endDate: number, categoryId?: string): DailyStats[] => {
    const filtered = getRecordsByDateRange(startDate, endDate, categoryId)
    const statsMap = new Map<string, DailyStats>()

    filtered.forEach(record => {
      const dateKey = formatDate(record.startTime, 'YYYY-MM-DD')
      if (!statsMap.has(dateKey)) {
        statsMap.set(dateKey, {
          date: dateKey,
          totalDuration: 0,
          totalSessions: 0,
          categoryStats: []
        })
      }

      const stats = statsMap.get(dateKey)!
      stats.totalDuration += record.duration
      stats.totalSessions++

      const catStat = stats.categoryStats.find(c => c.categoryId === record.categoryId)
      if (catStat) {
        catStat.duration += record.duration
        catStat.sessions++
      } else {
        stats.categoryStats.push({
          categoryId: record.categoryId,
          duration: record.duration,
          sessions: 1
        })
      }
    })

    return Array.from(statsMap.values()).sort((a, b) => b.date.localeCompare(a.date))
  }

  const getWeeklyTrendData = () => {
    const weekDates = getWeekDates()
    const data: { date: string; duration: number }[] = []

    weekDates.forEach(date => {
      const start = getStartOfDay(date.getTime())
      const end = getEndOfDay(date.getTime())
      const dayRecords = records.value.filter(r => r.startTime >= start && r.startTime <= end)
      const totalDuration = dayRecords.reduce((sum, r) => sum + r.duration, 0)
      data.push({
        date: `${date.getMonth() + 1}/${date.getDate()}`,
        duration: totalDuration
      })
    })

    return data
  }

  const getMonthlyTrendData = () => {
    const monthDates = getMonthDates()
    const data: { date: string; duration: number }[] = []

    monthDates.forEach(date => {
      const start = getStartOfDay(date.getTime())
      const end = getEndOfDay(date.getTime())
      const dayRecords = records.value.filter(r => r.startTime >= start && r.startTime <= end)
      const totalDuration = dayRecords.reduce((sum, r) => sum + r.duration, 0)
      data.push({
        date: `${date.getMonth() + 1}/${date.getDate()}`,
        duration: totalDuration
      })
    })

    return data
  }

  const getCategoryStats = (startDate: number, endDate: number) => {
    const filtered = getRecordsByDateRange(startDate, endDate)
    const statsMap = new Map<string, { name: string; value: number; color: string }>()

    filtered.forEach(record => {
      const category = categoriesStore.getCategoryById(record.categoryId)
      if (!category) return

      if (!statsMap.has(record.categoryId)) {
        statsMap.set(record.categoryId, {
          name: category.name,
          value: 0,
          color: category.color
        })
      }

      const stat = statsMap.get(record.categoryId)!
      stat.value += record.duration
    })

    return Array.from(statsMap.values())
  }

  const getHourlyDistribution = (startDate: number, endDate: number) => {
    const filtered = getRecordsByDateRange(startDate, endDate)
    const hours = Array(24).fill(0)

    filtered.forEach(record => {
      const startHour = new Date(record.startTime).getHours()
      const endHour = new Date(record.endTime).getHours()

      if (startHour === endHour) {
        hours[startHour] += record.duration
      } else {
        const startOfNextHour = new Date(record.startTime)
        startOfNextHour.setHours(startHour + 1, 0, 0, 0)
        const firstPart = (startOfNextHour.getTime() - record.startTime) / 1000
        hours[startHour] += firstPart

        for (let h = startHour + 1; h < endHour; h++) {
          hours[h] += 3600
        }

        const lastPart = record.duration - firstPart - (endHour - startHour - 1) * 3600
        if (lastPart > 0) {
          hours[endHour] += lastPart
        }
      }
    })

    return hours.map((duration, hour) => ({
      hour: `${hour.toString().padStart(2, '0')}:00`,
      hourNum: hour,
      duration
    }))
  }

  const getHeatmapData = () => {
    const data: { day: string; hour: number; value: number }[] = []
    const days = ['周一', '周二', '周三', '周四', '周五', '周六', '周日']

    const now = Date.now()
    const fourWeeksAgo = now - 28 * 24 * 60 * 60 * 1000

    const filtered = records.value.filter(r => r.startTime >= fourWeeksAgo)

    const weekHourMap = new Map<string, number>()

    filtered.forEach(record => {
      const date = new Date(record.startTime)
      const dayOfWeek = (date.getDay() + 6) % 7
      const hour = date.getHours()
      const key = `${dayOfWeek}-${hour}`
      weekHourMap.set(key, (weekHourMap.get(key) || 0) + record.duration)
    })

    for (let day = 0; day < 7; day++) {
      for (let hour = 0; hour < 24; hour++) {
        const key = `${day}-${hour}`
        data.push({
          day: days[day],
          hour,
          value: weekHourMap.get(key) || 0
        })
      }
    }

    return data
  }

  return {
    records,
    todayRecords,
    todayTotalDuration,
    todayCompletedPomodoros,
    todayLongestFocus,
    currentStreak,
    loadRecords,
    addRecord,
    deleteRecord,
    getRecordsByDateRange,
    getDailyStats,
    getWeeklyTrendData,
    getMonthlyTrendData,
    getCategoryStats,
    getHourlyDistribution,
    getHeatmapData
  }
})
