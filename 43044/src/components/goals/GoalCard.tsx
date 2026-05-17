import { MoreHorizontal, Edit2, Trash2, ChevronRight } from 'lucide-react'
import { useState } from 'react'
import { useAppStore } from '@/store'
import { ProgressBar } from '@/components/ui/ProgressBar'
import { Badge } from '@/components/ui/Badge'
import { calculateGoalProgress, getPriorityColor, getPriorityLabel } from '@/utils/progress'
import { formatDateReadable } from '@/utils/date'
import type { Goal } from '@/types'

interface GoalCardProps {
  goal: Goal
  onSelect: (id: string) => void
  onEdit: (id: string) => void
  onDelete: (id: string) => void
}

export function GoalCard({ goal, onSelect, onEdit, onDelete }: GoalCardProps) {
  const { tasks, tags } = useAppStore()
  const [showMenu, setShowMenu] = useState(false)
  
  const progress = calculateGoalProgress(goal, tasks)
  const goalTags = tags.filter((t) => goal.tagIds.includes(t.id))
  
  return (
    <div
      className="bg-white rounded-lg shadow-sm border border-gray-200 p-5 hover:shadow-md transition-shadow cursor-pointer"
      onClick={() => onSelect(goal.id)}
    >
      <div className="flex items-start justify-between mb-3">
        <div className="flex-1">
          <div className="flex items-center gap-2 mb-1">
            <h3 className="text-lg font-semibold text-gray-800">{goal.title}</h3>
            <Badge text={getPriorityLabel(goal.priority)} color={getPriorityColor(goal.priority)} size="sm" />
          </div>
          {goal.description && (
            <p className="text-sm text-gray-500 line-clamp-2">{goal.description}</p>
          )}
        </div>
        <div className="relative" onClick={(e) => e.stopPropagation()}>
          <button
            onClick={() => setShowMenu(!showMenu)}
            className="p-1 hover:bg-gray-100 rounded-full"
          >
            <MoreHorizontal className="w-5 h-5 text-gray-400" />
          </button>
          {showMenu && (
            <div className="absolute right-0 mt-1 bg-white rounded-lg shadow-lg border border-gray-200 py-1 z-10 min-w-[120px]">
              <button
                onClick={() => { onEdit(goal.id); setShowMenu(false) }}
                className="w-full px-4 py-2 text-left text-sm text-gray-700 hover:bg-gray-50 flex items-center gap-2"
              >
                <Edit2 className="w-4 h-4" />
                编辑
              </button>
              <button
                onClick={() => { onDelete(goal.id); setShowMenu(false) }}
                className="w-full px-4 py-2 text-left text-sm text-red-600 hover:bg-red-50 flex items-center gap-2"
              >
                <Trash2 className="w-4 h-4" />
                删除
              </button>
            </div>
          )}
        </div>
      </div>
      
      {goalTags.length > 0 && (
        <div className="flex flex-wrap gap-1.5 mb-3">
          {goalTags.map((tag) => (
            <span
              key={tag.id}
              className="px-2 py-0.5 rounded-full text-xs font-medium text-white"
              style={{ backgroundColor: tag.color }}
            >
              {tag.name}
            </span>
          ))}
        </div>
      )}
      
      <div className="mb-3">
        <ProgressBar percentage={progress.percentage} showLabel />
      </div>
      
      <div className="flex items-center justify-between text-sm text-gray-500">
        <div className="flex items-center gap-4">
          <span>{goal.milestoneIds.length} 个里程碑</span>
          <span>{goal.taskIds.length} 个任务</span>
        </div>
        {goal.dueDate && (
          <span>截止: {formatDateReadable(goal.dueDate)}</span>
        )}
      </div>
      
      <div className="mt-3 flex items-center justify-end">
        <span className="text-blue-600 text-sm font-medium flex items-center gap-1">
          查看详情 <ChevronRight className="w-4 h-4" />
        </span>
      </div>
    </div>
  )
}
