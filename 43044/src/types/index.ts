export type Priority = 'low' | 'medium' | 'high' | 'urgent'
export type TaskStatus = 'pending' | 'in-progress' | 'completed' | 'on-hold'

export interface Tag {
  id: string
  name: string
  color: string
}

export interface Task {
  id: string
  title: string
  description?: string
  status: TaskStatus
  priority: Priority
  dueDate?: string
  milestoneIds: string[]
  order: number
  createdAt: string
  updatedAt: string
  completedAt?: string
}

export interface Milestone {
  id: string
  title: string
  description?: string
  dueDate?: string
  order: number
  createdAt: string
  updatedAt: string
}

export interface Goal {
  id: string
  title: string
  description?: string
  priority: Priority
  dueDate?: string
  tagIds: string[]
  milestoneIds: string[]
  taskIds: string[]
  createdAt: string
  updatedAt: string
}

export interface AppState {
  goals: Goal[]
  milestones: Milestone[]
  tasks: Task[]
  tags: Tag[]
  selectedGoalId?: string
  selectedDate: string
}

export interface ProgressInfo {
  totalTasks: number
  completedTasks: number
  inProgressTasks: number
  pendingTasks: number
  onHoldTasks: number
  percentage: number
}
