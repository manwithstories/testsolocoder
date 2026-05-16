import { reactive, computed, watch } from 'vue'
import type { Task, Settings, PomodoroSession, DailyStats, TimerPhase, TimerStatus } from '../types'
import { storage, getTodayKey, generateId } from '../utils/storage'
import { playCompleteSound } from '../utils/audio'

const STORAGE_KEYS = {
  TASKS: 'focus-timer:tasks',
  SETTINGS: 'focus-timer:settings',
  STATS: 'focus-timer:stats',
  CURRENT_SESSION: 'focus-timer:current-session'
}

const defaultSettings: Settings = {
  workDuration: 25,
  shortBreakDuration: 5,
  longBreakDuration: 15,
  pomodorosBeforeLongBreak: 4,
  soundEnabled: true,
  autoStartBreak: false,
  autoStartWork: false
}

interface State {
  tasks: Task[]
  settings: Settings
  stats: DailyStats[]
  currentTaskId: string | null
  timerPhase: TimerPhase
  timerStatus: TimerStatus
  timeRemaining: number
  completedPomodorosInCycle: number
  currentSession: PomodoroSession | null
}

const state = reactive<State>({
  tasks: storage.get<Task[]>(STORAGE_KEYS.TASKS, []),
  settings: storage.get<Settings>(STORAGE_KEYS.SETTINGS, defaultSettings),
  stats: storage.get<DailyStats[]>(STORAGE_KEYS.STATS, []),
  currentTaskId: null,
  timerPhase: 'work',
  timerStatus: 'idle',
  timeRemaining: defaultSettings.workDuration * 60,
  completedPomodorosInCycle: 0,
  currentSession: null
})

watch(
  () => state.tasks,
  (tasks) => storage.set(STORAGE_KEYS.TASKS, tasks),
  { deep: true }
)

watch(
  () => state.settings,
  (settings) => storage.set(STORAGE_KEYS.SETTINGS, settings),
  { deep: true }
)

watch(
  () => state.stats,
  (stats) => storage.set(STORAGE_KEYS.STATS, stats),
  { deep: true }
)

function getTodayStats(): DailyStats {
  const todayKey = getTodayKey()
  let todayStats = state.stats.find(s => s.date === todayKey)
  
  if (!todayStats) {
    todayStats = {
      date: todayKey,
      completedPomodoros: 0,
      totalWorkMinutes: 0,
      sessions: []
    }
    state.stats.push(todayStats)
  }
  
  return todayStats
}

export function useStore() {
  const tasks = computed(() => state.tasks.filter(t => !t.completed))
  const completedTasks = computed(() => state.tasks.filter(t => t.completed))
  const settings = computed(() => state.settings)
  const timerPhase = computed(() => state.timerPhase)
  const timerStatus = computed(() => state.timerStatus)
  const timeRemaining = computed(() => state.timeRemaining)
  const currentTaskId = computed(() => state.currentTaskId)
  const currentTask = computed(() => state.tasks.find(t => t.id === state.currentTaskId) || null)
  
  const todayStats = computed(() => {
    const stats = getTodayStats()
    return {
      completedPomodoros: stats.completedPomodoros,
      totalWorkMinutes: stats.totalWorkMinutes,
      sessions: stats.sessions.filter(s => s.completed)
    }
  })

  const taskStats = computed(() => {
    const stats = getTodayStats()
    const taskMap = new Map<string, number>()
    
    stats.sessions
      .filter(s => s.completed && s.type === 'work' && s.taskId)
      .forEach(s => {
        const current = taskMap.get(s.taskId!) || 0
        taskMap.set(s.taskId!, current + s.duration)
      })
    
    return state.tasks.map(task => ({
      task,
      minutes: taskMap.get(task.id) || 0
    })).filter(s => s.minutes > 0).sort((a, b) => b.minutes - a.minutes)
  })

  function addTask(title: string): void {
    const task: Task = {
      id: generateId(),
      title: title.trim(),
      completedPomodoros: 0,
      createdAt: Date.now(),
      completed: false
    }
    state.tasks.unshift(task)
  }

  function toggleTaskComplete(taskId: string): void {
    const task = state.tasks.find(t => t.id === taskId)
    if (task) {
      task.completed = !task.completed
      task.completedAt = task.completed ? Date.now() : undefined
    }
  }

  function deleteTask(taskId: string): void {
    const index = state.tasks.findIndex(t => t.id === taskId)
    if (index > -1) {
      state.tasks.splice(index, 1)
    }
    if (state.currentTaskId === taskId) {
      state.currentTaskId = null
    }
  }

  function selectTask(taskId: string | null): void {
    state.currentTaskId = taskId
  }

  function updateSettings(newSettings: Partial<Settings>): void {
    Object.assign(state.settings, newSettings)
    if (state.timerStatus === 'idle' && 'workDuration' in newSettings && state.timerPhase === 'work') {
      state.timeRemaining = state.settings.workDuration * 60
    }
  }

  function startTimer(): void {
    if (state.timerStatus === 'running') return
    
    if (state.timerStatus === 'idle') {
      state.timeRemaining = getPhaseDuration(state.timerPhase)
      state.currentSession = {
        id: generateId(),
        taskId: state.currentTaskId,
        startTime: Date.now(),
        endTime: Date.now() + state.timeRemaining * 1000,
        duration: state.timeRemaining / 60,
        type: state.timerPhase,
        completed: false
      }
    }
    
    state.timerStatus = 'running'
  }

  function pauseTimer(): void {
    state.timerStatus = 'paused'
  }

  function resumeTimer(): void {
    state.timerStatus = 'running'
  }

  function resetTimer(): void {
    state.timerStatus = 'idle'
    state.currentSession = null
    state.timeRemaining = getPhaseDuration(state.timerPhase)
  }

  function skipPhase(): void {
    completePhase(false)
  }

  function setPhase(phase: TimerPhase): void {
    state.timerPhase = phase
    state.timerStatus = 'idle'
    state.currentSession = null
    state.timeRemaining = getPhaseDuration(phase)
  }

  function getPhaseDuration(phase: TimerPhase): number {
    switch (phase) {
      case 'work':
        return state.settings.workDuration * 60
      case 'shortBreak':
        return state.settings.shortBreakDuration * 60
      case 'longBreak':
        return state.settings.longBreakDuration * 60
    }
  }

  function completePhase(completed: boolean): void {
    const stats = getTodayStats()
    
    if (state.currentSession && state.timerPhase === 'work' && completed) {
      state.currentSession.completed = true
      state.currentSession.endTime = Date.now()
      stats.sessions.push(state.currentSession)
      stats.completedPomodoros++
      stats.totalWorkMinutes += state.settings.workDuration
      state.completedPomodorosInCycle++
      
      if (state.currentTaskId) {
        const task = state.tasks.find(t => t.id === state.currentTaskId)
        if (task) {
          task.completedPomodoros++
        }
      }
      
      if (state.settings.soundEnabled) {
        playCompleteSound()
      }
    }
    
    let nextPhase: TimerPhase
    if (state.timerPhase === 'work') {
      if (state.completedPomodorosInCycle >= state.settings.pomodorosBeforeLongBreak) {
        nextPhase = 'longBreak'
        state.completedPomodorosInCycle = 0
      } else {
        nextPhase = 'shortBreak'
      }
    } else {
      nextPhase = 'work'
    }
    
    state.timerPhase = nextPhase
    state.timerStatus = 'idle'
    state.currentSession = null
    state.timeRemaining = getPhaseDuration(nextPhase)
    
    if ((nextPhase === 'work' && state.settings.autoStartWork) ||
        (nextPhase !== 'work' && state.settings.autoStartBreak)) {
      setTimeout(() => startTimer(), 500)
    }
  }

  function tick(): void {
    if (state.timerStatus !== 'running') return
    
    state.timeRemaining--
    
    if (state.timeRemaining <= 0) {
      completePhase(true)
    }
  }

  function clearCompletedTasks(): void {
    state.tasks = state.tasks.filter(t => !t.completed)
  }

  return {
    tasks,
    completedTasks,
    settings,
    timerPhase,
    timerStatus,
    timeRemaining,
    currentTaskId,
    currentTask,
    todayStats,
    taskStats,
    addTask,
    toggleTaskComplete,
    deleteTask,
    selectTask,
    updateSettings,
    startTimer,
    pauseTimer,
    resumeTimer,
    resetTimer,
    skipPhase,
    setPhase,
    tick,
    clearCompletedTasks
  }
}
