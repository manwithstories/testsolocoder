import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { Goal, Milestone, Task, Tag, TaskStatus } from '@/types'
import { generateId } from '@/utils/progress'

interface AppStore {
  goals: Goal[]
  milestones: Milestone[]
  tasks: Task[]
  tags: Tag[]
  selectedGoalId: string | null
  selectedDate: string
  
  addGoal: (goal: Omit<Goal, 'id' | 'createdAt' | 'updatedAt' | 'milestoneIds' | 'taskIds'>) => void
  updateGoal: (id: string, updates: Partial<Goal>) => void
  deleteGoal: (id: string) => void
  setSelectedGoalId: (id: string | null) => void
  
  addMilestone: (goalId: string, milestone: Omit<Milestone, 'id' | 'createdAt' | 'updatedAt'>) => void
  updateMilestone: (id: string, updates: Partial<Milestone>) => void
  deleteMilestone: (id: string) => void
  
  addTask: (goalId: string, milestoneIds: string[], task: Omit<Task, 'id' | 'milestoneIds' | 'order' | 'createdAt' | 'updatedAt'>) => void
  updateTask: (id: string, updates: Partial<Task>) => void
  deleteTask: (id: string) => void
  updateTaskStatus: (id: string, status: TaskStatus) => void
  reorderTasks: (goalId: string, taskIds: string[]) => void
  reorderMilestones: (goalId: string, milestoneIds: string[]) => void
  
  addTag: (tag: Omit<Tag, 'id'>) => void
  updateTag: (id: string, updates: Partial<Tag>) => void
  deleteTag: (id: string) => void
  
  setSelectedDate: (date: string) => void
}

export const useAppStore = create<AppStore>()(
  persist(
    (set, get) => ({
      goals: [],
      milestones: [],
      tasks: [],
      tags: [],
      selectedGoalId: null,
      selectedDate: new Date().toISOString(),
      
      addGoal: (goalData) => {
        const now = new Date().toISOString()
        const newGoal: Goal = {
          ...goalData,
          id: generateId(),
          milestoneIds: [],
          taskIds: [],
          createdAt: now,
          updatedAt: now,
        }
        set((state) => ({ goals: [...state.goals, newGoal] }))
      },
      
      updateGoal: (id, updates) => {
        set((state) => ({
          goals: state.goals.map((g) =>
            g.id === id ? { ...g, ...updates, updatedAt: new Date().toISOString() } : g
          ),
        }))
      },
      
      deleteGoal: (id) => {
        const state = get()
        const goal = state.goals.find((g) => g.id === id)
        if (!goal) return
        
        set((state) => ({
          goals: state.goals.filter((g) => g.id !== id),
          milestones: state.milestones.filter((m) => !goal.milestoneIds.includes(m.id)),
          tasks: state.tasks.filter((t) => !goal.taskIds.includes(t.id)),
          selectedGoalId: state.selectedGoalId === id ? null : state.selectedGoalId,
        }))
      },
      
      setSelectedGoalId: (id) => set({ selectedGoalId: id }),
      
      addMilestone: (goalId, milestoneData) => {
        const state = get()
        const goal = state.goals.find((g) => g.id === goalId)
        if (!goal) return
        
        const now = new Date().toISOString()
        const newMilestone: Milestone = {
          ...milestoneData,
          id: generateId(),
          order: goal.milestoneIds.length,
          createdAt: now,
          updatedAt: now,
        }
        
        set((state) => ({
          milestones: [...state.milestones, newMilestone],
          goals: state.goals.map((g) =>
            g.id === goalId
              ? { ...g, milestoneIds: [...g.milestoneIds, newMilestone.id], updatedAt: now }
              : g
          ),
        }))
      },
      
      updateMilestone: (id, updates) => {
        set((state) => ({
          milestones: state.milestones.map((m) =>
            m.id === id ? { ...m, ...updates, updatedAt: new Date().toISOString() } : m
          ),
        }))
      },
      
      deleteMilestone: (id) => {
        set((state) => ({
          milestones: state.milestones.filter((m) => m.id !== id),
          goals: state.goals.map((g) => ({
            ...g,
            milestoneIds: g.milestoneIds.filter((mid) => mid !== id),
          })),
          tasks: state.tasks.map((t) => ({
            ...t,
            milestoneIds: t.milestoneIds.filter((mid) => mid !== id),
          })),
        }))
      },
      
      addTask: (goalId, milestoneIds, taskData) => {
        const state = get()
        const goal = state.goals.find((g) => g.id === goalId)
        if (!goal) return
        
        const now = new Date().toISOString()
        const newTask: Task = {
          ...taskData,
          id: generateId(),
          milestoneIds,
          order: goal.taskIds.length,
          createdAt: now,
          updatedAt: now,
        }
        
        set((state) => ({
          tasks: [...state.tasks, newTask],
          goals: state.goals.map((g) =>
            g.id === goalId
              ? { ...g, taskIds: [...g.taskIds, newTask.id], updatedAt: now }
              : g
          ),
        }))
      },
      
      updateTask: (id, updates) => {
        const now = new Date().toISOString()
        set((state) => ({
          tasks: state.tasks.map((t) =>
            t.id === id
              ? {
                  ...t,
                  ...updates,
                  updatedAt: now,
                  completedAt: updates.status === 'completed' && t.status !== 'completed'
                    ? now
                    : t.completedAt,
                }
              : t
          ),
        }))
      },
      
      deleteTask: (id) => {
        set((state) => ({
          tasks: state.tasks.filter((t) => t.id !== id),
          goals: state.goals.map((g) => ({
            ...g,
            taskIds: g.taskIds.filter((tid) => tid !== id),
          })),
        }))
      },
      
      updateTaskStatus: (id, status) => {
        const now = new Date().toISOString()
        set((state) => ({
          tasks: state.tasks.map((t) =>
            t.id === id
              ? {
                  ...t,
                  status,
                  updatedAt: now,
                  completedAt: status === 'completed' && t.status !== 'completed'
                    ? now
                    : t.completedAt,
                }
              : t
          ),
        }))
      },
      
      reorderTasks: (goalId, taskIds) => {
        set((state) => {
          const updatedTasks = state.tasks.map((t) => {
            const newOrder = taskIds.indexOf(t.id)
            return newOrder !== -1 ? { ...t, order: newOrder } : t
          })
          const updatedGoals = state.goals.map((g) =>
            g.id === goalId ? { ...g, taskIds, updatedAt: new Date().toISOString() } : g
          )
          return { tasks: updatedTasks, goals: updatedGoals }
        })
      },
      
      reorderMilestones: (goalId, milestoneIds) => {
        set((state) => {
          const updatedMilestones = state.milestones.map((m) => {
            const newOrder = milestoneIds.indexOf(m.id)
            return newOrder !== -1 ? { ...m, order: newOrder } : m
          })
          const updatedGoals = state.goals.map((g) =>
            g.id === goalId ? { ...g, milestoneIds, updatedAt: new Date().toISOString() } : g
          )
          return { milestones: updatedMilestones, goals: updatedGoals }
        })
      },
      
      addTag: (tagData) => {
        const newTag: Tag = {
          ...tagData,
          id: generateId(),
        }
        set((state) => ({ tags: [...state.tags, newTag] }))
      },
      
      updateTag: (id, updates) => {
        set((state) => ({
          tags: state.tags.map((t) => (t.id === id ? { ...t, ...updates } : t)),
        }))
      },
      
      deleteTag: (id) => {
        set((state) => ({
          tags: state.tags.filter((t) => t.id !== id),
          goals: state.goals.map((g) => ({
            ...g,
            tagIds: g.tagIds.filter((tid) => tid !== id),
          })),
        }))
      },
      
      setSelectedDate: (date) => set({ selectedDate: date }),
    }),
    {
      name: 'goal-tracker-storage',
    }
  )
)
