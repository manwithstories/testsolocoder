import type { Goal, Milestone, Task, ProgressInfo, TaskStatus } from '@/types'

export function generateId(): string {
  return Date.now().toString(36) + Math.random().toString(36).substr(2)
}

export function calculateGoalProgress(
  goal: Goal,
  allTasks: Task[]
): ProgressInfo {
  const goalTasks = allTasks.filter(task => goal.taskIds.includes(task.id))
  
  const uniqueTaskIds = new Set(goalTasks.map(t => t.id))
  const uniqueTasks = goalTasks.filter(t => uniqueTaskIds.has(t.id))

  return calculateProgressFromTasks(uniqueTasks)
}

export function calculateMilestoneProgress(
  milestone: Milestone,
  allTasks: Task[],
  goalTaskIds: string[]
): ProgressInfo {
  const milestoneTasks = allTasks.filter(
    task => task.milestoneIds.includes(milestone.id) && goalTaskIds.includes(task.id)
  )

  const uniqueTaskIds = new Set(milestoneTasks.map(t => t.id))
  const uniqueTasks = milestoneTasks.filter(t => uniqueTaskIds.has(t.id))

  return calculateProgressFromTasks(uniqueTasks)
}

function calculateProgressFromTasks(tasks: Task[]): ProgressInfo {
  const totalTasks = tasks.length
  const completedTasks = tasks.filter(t => t.status === 'completed').length
  const inProgressTasks = tasks.filter(t => t.status === 'in-progress').length
  const pendingTasks = tasks.filter(t => t.status === 'pending').length
  const onHoldTasks = tasks.filter(t => t.status === 'on-hold').length

  const percentage = totalTasks > 0 
    ? Math.round((completedTasks / totalTasks) * 100) 
    : 0

  return {
    totalTasks,
    completedTasks,
    inProgressTasks,
    pendingTasks,
    onHoldTasks,
    percentage,
  }
}

export function isOverdue(task: Task): boolean {
  if (!task.dueDate || task.status === 'completed') return false
  return new Date(task.dueDate) < new Date()
}

export function isDueToday(task: Task): boolean {
  if (!task.dueDate) return false
  const today = new Date()
  const dueDate = new Date(task.dueDate)
  return (
    today.getFullYear() === dueDate.getFullYear() &&
    today.getMonth() === dueDate.getMonth() &&
    today.getDate() === dueDate.getDate()
  )
}

export function getTasksByDate(tasks: Task[], date: Date): Task[] {
  return tasks.filter(task => {
    if (!task.dueDate) return false
    const taskDate = new Date(task.dueDate)
    return (
      taskDate.getFullYear() === date.getFullYear() &&
      taskDate.getMonth() === date.getMonth() &&
      taskDate.getDate() === date.getDate()
    )
  })
}

export function getTasksInDateRange(
  tasks: Task[],
  startDate: Date,
  endDate: Date
): Task[] {
  return tasks.filter(task => {
    if (!task.dueDate) return false
    const taskDate = new Date(task.dueDate)
    return taskDate >= startDate && taskDate <= endDate
  })
}

export function getStatusColor(status: TaskStatus): string {
  const colors: Record<TaskStatus, string> = {
    pending: 'bg-gray-400',
    'in-progress': 'bg-blue-500',
    completed: 'bg-green-500',
    'on-hold': 'bg-yellow-500',
  }
  return colors[status]
}

export function getStatusLabel(status: TaskStatus): string {
  const labels: Record<TaskStatus, string> = {
    pending: '待开始',
    'in-progress': '进行中',
    completed: '已完成',
    'on-hold': '已搁置',
  }
  return labels[status]
}

export function getPriorityColor(priority: string): string {
  const colors: Record<string, string> = {
    low: 'bg-gray-300 text-gray-700',
    medium: 'bg-blue-100 text-blue-700',
    high: 'bg-orange-100 text-orange-700',
    urgent: 'bg-red-100 text-red-700',
  }
  return colors[priority] || colors.medium
}

export function getPriorityLabel(priority: string): string {
  const labels: Record<string, string> = {
    low: '低',
    medium: '中',
    high: '高',
    urgent: '紧急',
  }
  return labels[priority] || '中'
}
