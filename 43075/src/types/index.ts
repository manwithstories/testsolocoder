export interface Category {
  id: string
  name: string
  color: string
  icon: string
  createdAt: number
}

export type TimerMode = 'pomodoro' | 'custom'

export type TimerStatus = 'idle' | 'running' | 'paused'

export interface FocusRecord {
  id: string
  categoryId: string
  mode: TimerMode
  duration: number
  plannedDuration: number
  startTime: number
  endTime: number
  completed: boolean
}

export interface Settings {
  pomodoroDuration: number
  shortBreakDuration: number
  longBreakDuration: number
  longBreakInterval: number
  soundEnabled: boolean
  darkMode: boolean
  autoStartBreak: boolean
  autoStartFocus: boolean
}

export interface LogEntry {
  id: string
  timestamp: number
  level: 'info' | 'warn' | 'error'
  message: string
  data?: any
}

export interface ToastMessage {
  id: string
  type: 'success' | 'error' | 'warning' | 'info'
  message: string
  duration: number
}

export interface DailyStats {
  date: string
  totalDuration: number
  totalSessions: number
  categoryStats: Array<{
    categoryId: string
    duration: number
    sessions: number
  }>
}

export const STORAGE_KEYS = {
  CATEGORIES: 'focus_categories',
  RECORDS: 'focus_records',
  SETTINGS: 'focus_settings',
  LOGS: 'focus_logs'
} as const

export const DEFAULT_SETTINGS: Settings = {
  pomodoroDuration: 25,
  shortBreakDuration: 5,
  longBreakDuration: 15,
  longBreakInterval: 4,
  soundEnabled: true,
  darkMode: false,
  autoStartBreak: false,
  autoStartFocus: false
}

export const DEFAULT_CATEGORIES: Omit<Category, 'id' | 'createdAt'>[] = [
  { name: '工作', color: '#6366f1', icon: 'Briefcase' },
  { name: '学习', color: '#10b981', icon: 'BookOpen' },
  { name: '阅读', color: '#f59e0b', icon: 'Book' },
  { name: '运动', color: '#ef4444', icon: 'Dumbbell' },
  { name: '其他', color: '#8b5cf6', icon: 'MoreHorizontal' }
]

export const ICON_OPTIONS = [
  'Briefcase', 'BookOpen', 'Book', 'Dumbbell', 'Code', 'Music',
  'Pencil', 'Coffee', 'Headphones', 'Target', 'Zap', 'Heart',
  'Star', 'Flag', 'Trophy', 'Medal', 'Clock', 'Calendar',
  'Layers', 'Folder', 'Tag', 'MoreHorizontal'
]

export const COLOR_OPTIONS = [
  '#6366f1', '#8b5cf6', '#ec4899', '#ef4444', '#f59e0b',
  '#10b981', '#06b6d4', '#3b82f6', '#64748b', '#84cc16',
  '#f97316', '#14b8a6', '#a855f7', '#22c55e', '#0ea5e9'
]
