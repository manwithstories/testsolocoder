export interface Task {
  id: string
  title: string
  completedPomodoros: number
  createdAt: number
  completed: boolean
  completedAt?: number
}

export interface Settings {
  workDuration: number
  shortBreakDuration: number
  longBreakDuration: number
  pomodorosBeforeLongBreak: number
  soundEnabled: boolean
  autoStartBreak: boolean
  autoStartWork: boolean
}

export interface PomodoroSession {
  id: string
  taskId: string | null
  startTime: number
  endTime: number
  duration: number
  type: 'work' | 'shortBreak' | 'longBreak'
  completed: boolean
}

export interface DailyStats {
  date: string
  completedPomodoros: number
  totalWorkMinutes: number
  sessions: PomodoroSession[]
}

export type TimerPhase = 'work' | 'shortBreak' | 'longBreak'

export type TimerStatus = 'idle' | 'running' | 'paused'
