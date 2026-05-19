import { STORAGE_KEYS, DEFAULT_SETTINGS, DEFAULT_CATEGORIES, type Category, type FocusRecord, type Settings, type LogEntry } from '@/types'
import { logger } from './logger'

const generateId = (): string => {
  return Date.now().toString(36) + Math.random().toString(36).substr(2, 9)
}

const safeParse = <T>(jsonString: string | null, defaultValue: T): T => {
  if (!jsonString) return defaultValue
  try {
    return JSON.parse(jsonString) as T
  } catch (error) {
    logger.error('JSON解析失败', { error, jsonString })
    return defaultValue
  }
}

const safeStringify = (value: any): string | null => {
  try {
    return JSON.stringify(value)
  } catch (error) {
    logger.error('JSON序列化失败', { error, value })
    return null
  }
}

const checkQuota = (): boolean => {
  try {
    const testKey = '__quota_test__'
    const testValue = 'x'.repeat(1024)
    localStorage.setItem(testKey, testValue)
    localStorage.removeItem(testKey)
    return true
  } catch (error) {
    logger.error('存储空间不足', { error })
    return false
  }
}

export const storage = {
  getCategories: (): Category[] => {
    const data = localStorage.getItem(STORAGE_KEYS.CATEGORIES)
    const categories = safeParse<Category[]>(data, [])
    if (categories.length === 0) {
      const defaultCats: Category[] = DEFAULT_CATEGORIES.map(cat => ({
        ...cat,
        id: generateId(),
        createdAt: Date.now()
      }))
      storage.saveCategories(defaultCats)
      return defaultCats
    }
    return categories
  },

  saveCategories: (categories: Category[]): boolean => {
    if (!checkQuota()) return false
    const json = safeStringify(categories)
    if (!json) return false
    try {
      localStorage.setItem(STORAGE_KEYS.CATEGORIES, json)
      logger.info('保存分类成功', { count: categories.length })
      return true
    } catch (error) {
      logger.error('保存分类失败', { error })
      return false
    }
  },

  getRecords: (): FocusRecord[] => {
    const data = localStorage.getItem(STORAGE_KEYS.RECORDS)
    return safeParse<FocusRecord[]>(data, [])
  },

  saveRecords: (records: FocusRecord[]): boolean => {
    if (!checkQuota()) return false
    const json = safeStringify(records)
    if (!json) return false
    try {
      localStorage.setItem(STORAGE_KEYS.RECORDS, json)
      return true
    } catch (error) {
      logger.error('保存记录失败', { error })
      return false
    }
  },

  addRecord: (record: Omit<FocusRecord, 'id'>): FocusRecord | null => {
    const records = storage.getRecords()
    const newRecord: FocusRecord = {
      ...record,
      id: generateId()
    }
    records.push(newRecord)
    if (storage.saveRecords(records)) {
      logger.info('添加专注记录', { record: newRecord })
      return newRecord
    }
    return null
  },

  deleteRecord: (recordId: string): boolean => {
    const records = storage.getRecords()
    const filtered = records.filter(r => r.id !== recordId)
    if (storage.saveRecords(filtered)) {
      logger.info('删除专注记录', { recordId })
      return true
    }
    return false
  },

  getSettings: (): Settings => {
    const data = localStorage.getItem(STORAGE_KEYS.SETTINGS)
    const saved = safeParse<Partial<Settings>>(data, {})
    return { ...DEFAULT_SETTINGS, ...saved }
  },

  saveSettings: (settings: Settings): boolean => {
    if (!checkQuota()) return false
    const json = safeStringify(settings)
    if (!json) return false
    try {
      localStorage.setItem(STORAGE_KEYS.SETTINGS, json)
      logger.info('保存设置成功', { settings })
      return true
    } catch (error) {
      logger.error('保存设置失败', { error })
      return false
    }
  },

  getLogs: (): LogEntry[] => {
    const data = localStorage.getItem(STORAGE_KEYS.LOGS)
    return safeParse<LogEntry[]>(data, [])
  },

  saveLogs: (logs: LogEntry[]): boolean => {
    if (!checkQuota()) return false
    const limitedLogs = logs.slice(-1000)
    const json = safeStringify(limitedLogs)
    if (!json) return false
    try {
      localStorage.setItem(STORAGE_KEYS.LOGS, json)
      return true
    } catch (error) {
      console.error('保存日志失败', error)
      return false
    }
  },

  addLog: (entry: Omit<LogEntry, 'id'>): LogEntry => {
    const logs = storage.getLogs()
    const newEntry: LogEntry = {
      ...entry,
      id: generateId()
    }
    logs.push(newEntry)
    storage.saveLogs(logs)
    return newEntry
  },

  clearLogs: (): boolean => {
    try {
      localStorage.removeItem(STORAGE_KEYS.LOGS)
      return true
    } catch (error) {
      logger.error('清空日志失败', { error })
      return false
    }
  },

  exportData: (): string | null => {
    try {
      const data = {
        version: '1.0',
        exportedAt: Date.now(),
        categories: storage.getCategories(),
        records: storage.getRecords(),
        settings: storage.getSettings(),
        logs: storage.getLogs()
      }
      return JSON.stringify(data, null, 2)
    } catch (error) {
      logger.error('导出数据失败', { error })
      return null
    }
  },

  importData: (jsonString: string): { success: boolean; error?: string } => {
    try {
      const data = JSON.parse(jsonString)

      if (!data.categories || !data.records || !data.settings) {
        return { success: false, error: '数据格式不完整' }
      }

      if (!Array.isArray(data.categories) || !Array.isArray(data.records)) {
        return { success: false, error: '数据格式不正确' }
      }

      if (!storage.saveCategories(data.categories)) {
        return { success: false, error: '导入分类失败' }
      }
      if (!storage.saveRecords(data.records)) {
        return { success: false, error: '导入记录失败' }
      }
      if (!storage.saveSettings(data.settings)) {
        return { success: false, error: '导入设置失败' }
      }
      if (data.logs && Array.isArray(data.logs)) {
        storage.saveLogs(data.logs)
      }

      logger.info('导入数据成功', { recordCount: data.records.length })
      return { success: true }
    } catch (error) {
      logger.error('导入数据失败', { error })
      return { success: false, error: '数据解析失败，请检查文件格式' }
    }
  },

  clearAll: (): boolean => {
    try {
      localStorage.removeItem(STORAGE_KEYS.CATEGORIES)
      localStorage.removeItem(STORAGE_KEYS.RECORDS)
      localStorage.removeItem(STORAGE_KEYS.SETTINGS)
      localStorage.removeItem(STORAGE_KEYS.LOGS)
      logger.info('清空所有数据')
      return true
    } catch (error) {
      logger.error('清空数据失败', { error })
      return false
    }
  }
}
