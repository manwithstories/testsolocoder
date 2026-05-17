import { useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { ArrowLeft, Plus, Edit2, Trash2, GripVertical } from 'lucide-react'
import {
  DndContext,
  closestCenter,
  KeyboardSensor,
  PointerSensor,
  useSensor,
  useSensors,
  DragEndEvent,
} from '@dnd-kit/core'
import {
  arrayMove,
  SortableContext,
  sortableKeyboardCoordinates,
  verticalListSortingStrategy,
  useSortable,
} from '@dnd-kit/sortable'
import { CSS } from '@dnd-kit/utilities'
import { useAppStore } from '@/store'
import { Button } from '@/components/ui/Button'
import { ProgressBar } from '@/components/ui/ProgressBar'
import { Badge } from '@/components/ui/Badge'
import { TaskItem } from '@/components/tasks/TaskItem'
import { TaskForm } from '@/components/tasks/TaskForm'
import { MilestoneForm } from '@/components/milestones/MilestoneForm'
import { calculateGoalProgress, calculateMilestoneProgress, getPriorityColor, getPriorityLabel } from '@/utils/progress'
import { formatDateReadable } from '@/utils/date'


interface SortableMilestoneProps {
  milestoneId: string
  goalId: string
  allTaskIds: string[]
  onEdit: (id: string) => void
  onDelete: (id: string) => void
  onAddTask: (milestoneId: string) => void
  onEditTask: (taskId: string) => void
  onDeleteTask: (taskId: string) => void
}

function SortableMilestone({ milestoneId, goalId, allTaskIds, onEdit, onDelete, onAddTask, onEditTask, onDeleteTask }: SortableMilestoneProps) {
  const { milestones, tasks, goals } = useAppStore()
  const milestone = milestones.find((m) => m.id === milestoneId)
  const goal = goals.find((g) => g.id === goalId)
  
  const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({
    id: `milestone-${milestoneId}`,
  })
  
  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
  }
  
  if (!milestone || !goal) return null
  
  const progress = calculateMilestoneProgress(milestone, tasks, goal.taskIds)
  const milestoneTaskIds = allTaskIds.filter((taskId) => {
    const task = tasks.find((t) => t.id === taskId)
    return task?.milestoneIds.includes(milestoneId)
  })
  
  return (
    <div
      ref={setNodeRef}
      style={style}
      className="bg-white rounded-lg border border-gray-200 p-5"
    >
      <div className="flex items-start justify-between mb-4">
        <div className="flex items-start gap-2">
          <div
            {...attributes}
            {...listeners}
            className="cursor-grab active:cursor-grabbing pt-1"
          >
            <GripVertical className="w-4 h-4 text-gray-400" />
          </div>
          <div>
            <h3 className="text-lg font-semibold text-gray-800">{milestone.title}</h3>
            {milestone.description && (
              <p className="text-sm text-gray-500 mt-1">{milestone.description}</p>
            )}
          </div>
        </div>
        <div className="flex items-center gap-2">
          <button
            onClick={() => onAddTask(milestoneId)}
            className="p-1.5 hover:bg-blue-50 rounded-lg text-blue-600"
            title="添加任务"
          >
            <Plus className="w-4 h-4" />
          </button>
          <button
            onClick={() => onEdit(milestoneId)}
            className="p-1.5 hover:bg-gray-100 rounded-lg"
          >
            <Edit2 className="w-4 h-4 text-gray-400" />
          </button>
          <button
            onClick={() => onDelete(milestoneId)}
            className="p-1.5 hover:bg-red-50 rounded-lg"
          >
            <Trash2 className="w-4 h-4 text-red-400" />
          </button>
        </div>
      </div>
      
      <div className="mb-4">
        <ProgressBar percentage={progress.percentage} />
      </div>
      
      {milestone.dueDate && (
        <div className="text-sm text-gray-500 mb-3">
          截止: {formatDateReadable(milestone.dueDate)}
        </div>
      )}
      
      <div className="space-y-2">
        {milestoneTaskIds.length === 0 ? (
          <p className="text-sm text-gray-400 text-center py-4">暂无任务，点击 + 添加</p>
        ) : (
          <SortableContext
            items={milestoneTaskIds.map((id) => `task-${id}`)}
            strategy={verticalListSortingStrategy}
          >
            {milestoneTaskIds.map((taskId) => {
              const task = tasks.find((t) => t.id === taskId)
              if (!task) return null
              return (
                <TaskItem
                  key={task.id}
                  task={task}
                  onEdit={onEditTask}
                  onDelete={onDeleteTask}
                />
              )
            })}
          </SortableContext>
        )}
      </div>
    </div>
  )
}

export function GoalDetailPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const { goals, milestones, tasks, deleteGoal, deleteMilestone, deleteTask, reorderMilestones, reorderTasks } = useAppStore()
  
  const goal = goals.find((g) => g.id === id)
  
  const [isMilestoneFormOpen, setIsMilestoneFormOpen] = useState(false)
  const [editingMilestoneId, setEditingMilestoneId] = useState<string | null>(null)
  const [isTaskFormOpen, setIsTaskFormOpen] = useState(false)
  const [editingTaskId, setEditingTaskId] = useState<string | null>(null)
  
  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: {
        distance: 8,
      },
    }),
    useSensor(KeyboardSensor, {
      coordinateGetter: sortableKeyboardCoordinates,
    })
  )
  
  if (!goal) {
    return (
      <div className="text-center py-12">
        <p className="text-gray-500">目标不存在</p>
        <Button variant="secondary" className="mt-4" onClick={() => navigate('/')}>
          返回目标列表
        </Button>
      </div>
    )
  }
  
  const goalProgress = calculateGoalProgress(goal, tasks)
  const sortedMilestoneIds = [...goal.milestoneIds].sort((a, b) => {
    const m1 = milestones.find((m) => m.id === a)
    const m2 = milestones.find((m) => m.id === b)
    return (m1?.order || 0) - (m2?.order || 0)
  })
  
  const sortedTaskIds = [...goal.taskIds].sort((a, b) => {
    const t1 = tasks.find((t) => t.id === a)
    const t2 = tasks.find((t) => t.id === b)
    return (t1?.order || 0) - (t2?.order || 0)
  })
  
  const unassignedTaskIds = sortedTaskIds.filter((taskId) => {
    const task = tasks.find((t) => t.id === taskId)
    return task?.milestoneIds.length === 0
  })
  
  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event
    
    if (!over || active.id === over.id) return
    
    const activeIdStr = String(active.id)
    const overIdStr = String(over.id)
    
    if (activeIdStr.startsWith('milestone-') && overIdStr.startsWith('milestone-')) {
      const activeMilestoneId = activeIdStr.replace('milestone-', '')
      const overMilestoneId = overIdStr.replace('milestone-', '')
      
      const oldIndex = sortedMilestoneIds.indexOf(activeMilestoneId)
      const newIndex = sortedMilestoneIds.indexOf(overMilestoneId)
      
      if (oldIndex !== -1 && newIndex !== -1) {
        const newOrder = arrayMove(sortedMilestoneIds, oldIndex, newIndex)
        reorderMilestones(goal.id, newOrder)
      }
    } else if (activeIdStr.startsWith('task-') && overIdStr.startsWith('task-')) {
      const activeTaskId = activeIdStr.replace('task-', '')
      const overTaskId = overIdStr.replace('task-', '')
      
      const oldIndex = sortedTaskIds.indexOf(activeTaskId)
      const newIndex = sortedTaskIds.indexOf(overTaskId)
      
      if (oldIndex !== -1 && newIndex !== -1) {
        const newOrder = arrayMove(sortedTaskIds, oldIndex, newIndex)
        reorderTasks(goal.id, newOrder)
      }
    }
  }
  
  const handleEditMilestone = (milestoneId: string) => {
    setEditingMilestoneId(milestoneId)
    setIsMilestoneFormOpen(true)
  }
  
  const handleDeleteMilestone = (milestoneId: string) => {
    if (confirm('确定要删除这个里程碑吗？')) {
      deleteMilestone(milestoneId)
    }
  }
  
  const handleAddTaskToMilestone = (_milestoneId: string) => {
    setEditingTaskId(null)
    setIsTaskFormOpen(true)
  }
  
  const handleEditTask = (taskId: string) => {
    setEditingTaskId(taskId)
    setIsTaskFormOpen(true)
  }
  
  const handleDeleteTask = (taskId: string) => {
    if (confirm('确定要删除这个任务吗？')) {
      deleteTask(taskId)
    }
  }
  
  const allDraggableIds = [
    ...sortedMilestoneIds.map((id) => `milestone-${id}`),
    ...sortedTaskIds.map((id) => `task-${id}`),
  ]
  
  return (
    <DndContext
      sensors={sensors}
      collisionDetection={closestCenter}
      onDragEnd={handleDragEnd}
    >
      <SortableContext
        items={allDraggableIds}
        strategy={verticalListSortingStrategy}
      >
        <div>
          <div className="mb-6">
            <button
              onClick={() => navigate('/')}
              className="flex items-center gap-2 text-gray-600 hover:text-gray-800 mb-4"
            >
              <ArrowLeft className="w-4 h-4" />
              返回目标列表
            </button>
            
            <div className="bg-white rounded-lg border border-gray-200 p-6">
              <div className="flex items-start justify-between mb-4">
                <div>
                  <h1 className="text-2xl font-bold text-gray-800">{goal.title}</h1>
                  {goal.description && (
                    <p className="text-gray-500 mt-2">{goal.description}</p>
                  )}
                </div>
                <div className="flex items-center gap-2">
                  <Badge text={getPriorityLabel(goal.priority)} color={getPriorityColor(goal.priority)} />
                  {goal.dueDate && (
                    <span className="text-sm text-gray-500">
                      截止: {formatDateReadable(goal.dueDate)}
                    </span>
                  )}
                  <Button
                    variant="danger"
                    size="sm"
                    onClick={() => {
                      if (confirm('确定要删除这个目标吗？')) {
                        deleteGoal(goal.id)
                        navigate('/')
                      }
                    }}
                  >
                    删除目标
                  </Button>
                </div>
              </div>
              
              <div className="mb-4">
                <ProgressBar percentage={goalProgress.percentage} size="lg" />
                <div className="flex justify-between text-sm text-gray-500 mt-2">
                  <span>
                    {goalProgress.completedTasks} / {goalProgress.totalTasks} 任务已完成
                  </span>
                  <span className="font-medium">{goalProgress.percentage}%</span>
                </div>
              </div>
              
              <div className="flex gap-6 text-sm">
                <div>
                  <span className="text-gray-500">里程碑: </span>
                  <span className="font-medium">{goal.milestoneIds.length}</span>
                </div>
                <div>
                  <span className="text-gray-500">任务: </span>
                  <span className="font-medium">{goal.taskIds.length}</span>
                </div>
                <div>
                  <span className="text-gray-500">进行中: </span>
                  <span className="font-medium text-blue-600">{goalProgress.inProgressTasks}</span>
                </div>
                <div>
                  <span className="text-gray-500">已搁置: </span>
                  <span className="font-medium text-yellow-600">{goalProgress.onHoldTasks}</span>
                </div>
              </div>
            </div>
          </div>
          
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-xl font-semibold text-gray-800">里程碑</h2>
            <Button onClick={() => { setEditingMilestoneId(null); setIsMilestoneFormOpen(true) }}>
              <Plus className="w-4 h-4 mr-2" />
              添加里程碑
            </Button>
          </div>
          
          <div className="space-y-4">
            {sortedMilestoneIds.map((milestoneId) => (
              <SortableMilestone
                key={milestoneId}
                milestoneId={milestoneId}
                goalId={goal.id}
                allTaskIds={sortedTaskIds}
                onEdit={handleEditMilestone}
                onDelete={handleDeleteMilestone}
                onAddTask={handleAddTaskToMilestone}
                onEditTask={handleEditTask}
                onDeleteTask={handleDeleteTask}
              />
            ))}
          </div>
          
          {unassignedTaskIds.length > 0 && (
            <div className="mt-8">
              <div className="flex items-center justify-between mb-4">
                <h2 className="text-xl font-semibold text-gray-800">未分配里程碑的任务</h2>
                <Button onClick={() => { setEditingTaskId(null); setIsTaskFormOpen(true) }}>
                  <Plus className="w-4 h-4 mr-2" />
                  添加任务
                </Button>
              </div>
              
              <div className="bg-white rounded-lg border border-gray-200 p-5 space-y-2">
                <SortableContext
                  items={unassignedTaskIds.map((id) => `task-${id}`)}
                  strategy={verticalListSortingStrategy}
                >
                  {unassignedTaskIds.map((taskId) => {
                    const task = tasks.find((t) => t.id === taskId)
                    if (!task) return null
                    return (
                      <TaskItem
                        key={task.id}
                        task={task}
                        onEdit={handleEditTask}
                        onDelete={handleDeleteTask}
                      />
                    )
                  })}
                </SortableContext>
              </div>
            </div>
          )}
          
          {sortedMilestoneIds.length === 0 && unassignedTaskIds.length === 0 && (
            <div className="bg-white rounded-lg border border-gray-200 p-12 text-center">
              <p className="text-gray-500 mb-4">还没有里程碑和任务</p>
              <div className="flex justify-center gap-4">
                <Button onClick={() => { setEditingMilestoneId(null); setIsMilestoneFormOpen(true) }}>
                  添加里程碑
                </Button>
                <Button variant="secondary" onClick={() => { setEditingTaskId(null); setIsTaskFormOpen(true) }}>
                  添加任务
                </Button>
              </div>
            </div>
          )}
          
          <MilestoneForm
            isOpen={isMilestoneFormOpen}
            onClose={() => { setIsMilestoneFormOpen(false); setEditingMilestoneId(null) }}
            goalId={goal.id}
            milestoneId={editingMilestoneId}
          />
          
          <TaskForm
            isOpen={isTaskFormOpen}
            onClose={() => { setIsTaskFormOpen(false); setEditingTaskId(null) }}
            goalId={goal.id}
            taskId={editingTaskId}
          />
        </div>
      </SortableContext>
    </DndContext>
  )
}
