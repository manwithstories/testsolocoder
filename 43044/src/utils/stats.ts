import { eachDayOfInterval, subDays, isWithinInterval, parseISO, startOfDay } from 'date-fns'
import type { Goal, Task, Tag } from '@/types'
import { calculateGoalProgress, isOverdue } from '@/utils/progress'

export interface TagStats {
  tag: Tag
  totalGoals: number
  completedGoals: number
  completionRate: number
  totalTasks: number
  completedTasks: number
}

export interface DailyActivity {
  date: string
  completedTasks: number
  createdTasks: number
}

export function getTagStats(goals: Goal[], tasks: Task[], tags: Tag[]): TagStats[] {
  return tags.map((tag) => {
    const tagGoals = goals.filter((g) => g.tagIds.includes(tag.id))
    const tagTasks = tasks.filter((t) => {
      const goal = goals.find((g) => g.taskIds.includes(t.id))
      return goal?.tagIds.includes(tag.id)
    })
    
    const completedGoals = tagGoals.filter((g) => {
      const progress = calculateGoalProgress(g, tasks)
      return progress.percentage === 100
    }).length
    
    const completedTasks = tagTasks.filter((t) => t.status === 'completed').length
    
    return {
      tag,
      totalGoals: tagGoals.length,
      completedGoals,
      completionRate: tagGoals.length > 0 ? Math.round((completedGoals / tagGoals.length) * 100) : 0,
      totalTasks: tagTasks.length,
      completedTasks,
    }
  })
}

export function getOverdueTasks(tasks: Task[]): Task[] {
  return tasks.filter(isOverdue)
}

export function getWeeklyActivity(tasks: Task[], days: number = 7): DailyActivity[] {
  const endDate = new Date()
  const startDate = subDays(endDate, days - 1)
  
  const daysInRange = eachDayOfInterval({ start: startDate, end: endDate })
  
  return daysInRange.map((day) => {
    const startOfDayDate = startOfDay(day)
    const endOfDayDate = new Date(startOfDayDate)
    endOfDayDate.setHours(23, 59, 59, 999)
    
    const completedTasks = tasks.filter((t) => {
      if (!t.completedAt) return false
      const completedAt = parseISO(t.completedAt)
      return isWithinInterval(completedAt, { start: startOfDayDate, end: endOfDayDate })
    }).length
    
    const createdTasks = tasks.filter((t) => {
      const createdAt = parseISO(t.createdAt)
      return isWithinInterval(createdAt, { start: startOfDayDate, end: endOfDayDate })
    }).length
    
    return {
      date: startOfDayDate.toISOString(),
      completedTasks,
      createdTasks,
    }
  })
}

export function getOverallStats(goals: Goal[], tasks: Task[]) {
  const totalGoals = goals.length
  const completedGoals = goals.filter((g) => {
    const progress = calculateGoalProgress(g, tasks)
    return progress.percentage === 100
  }).length
  
  const totalTasks = tasks.length
  const completedTasks = tasks.filter((t) => t.status === 'completed').length
  const inProgressTasks = tasks.filter((t) => t.status === 'in-progress').length
  const pendingTasks = tasks.filter((t) => t.status === 'pending').length
  const onHoldTasks = tasks.filter((t) => t.status === 'on-hold').length
  const overdueTasks = getOverdueTasks(tasks).length
  
  const overallCompletionRate = totalTasks > 0
    ? Math.round((completedTasks / totalTasks) * 100)
    : 0
  
  const goalsInProgress = totalGoals - completedGoals
  
  return {
    totalGoals,
    completedGoals,
    goalsInProgress,
    totalTasks,
    completedTasks,
    inProgressTasks,
    pendingTasks,
    onHoldTasks,
    overdueTasks,
    overallCompletionRate,
  }
}

export function getPriorityStats(tasks: Task[]) {
  const priorities = ['low', 'medium', 'high', 'urgent'] as const
  
  return priorities.map((priority) => {
    const priorityTasks = tasks.filter((t) => t.priority === priority)
    const completed = priorityTasks.filter((t) => t.status === 'completed').length
    
    return {
      priority,
      total: priorityTasks.length,
      completed,
      rate: priorityTasks.length > 0 ? Math.round((completed / priorityTasks.length) * 100) : 0,
    }
  })
}

export function getMonthlyStats(_goals: Goal[], tasks: Task[], months: number = 6) {
  const result = []
  const now = new Date()
  
  for (let i = months - 1; i >= 0; i--) {
    const date = new Date(now.getFullYear(), now.getMonth() - i, 1)
    const monthStart = new Date(date.getFullYear(), date.getMonth(), 1)
    const monthEnd = new Date(date.getFullYear(), date.getMonth() + 1, 0, 23, 59, 59)
    
    const monthTasks = tasks.filter((t) => {
      const createdAt = parseISO(t.createdAt)
      return isWithinInterval(createdAt, { start: monthStart, end: monthEnd })
    })
    
    const completedTasks = monthTasks.filter((t) => t.status === 'completed').length
    
    result.push({
      month: `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`,
      label: `${date.getFullYear()}年${date.getMonth() + 1}月`,
      created: monthTasks.length,
      completed: completedTasks,
    })
  }
  
  return result
}
