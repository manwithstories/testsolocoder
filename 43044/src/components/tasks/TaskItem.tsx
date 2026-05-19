import { GripVertical, Edit2, Trash2 } from 'lucide-react'
import { useSortable } from '@dnd-kit/sortable'
import { CSS } from '@dnd-kit/utilities'
import { useAppStore } from '@/store'
import { Badge } from '@/components/ui/Badge'
import { getStatusColor, getStatusLabel, getPriorityColor, getPriorityLabel, isOverdue, isDueToday } from '@/utils/progress'
import { formatDateReadable } from '@/utils/date'
import type { Task, TaskStatus } from '@/types'

interface TaskItemProps {
  task: Task
  onEdit: (id: string) => void
  onDelete: (id: string) => void
}

export function TaskItem({ task, onEdit, onDelete }: TaskItemProps) {
  const { updateTaskStatus, milestones } = useAppStore()
  
  const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({
    id: `task-${task.id}`,
  })
  
  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
  }
  
  const taskMilestones = milestones.filter((m) => task.milestoneIds.includes(m.id))
  const overdue = isOverdue(task)
  const dueToday = isDueToday(task)
  
  const handleStatusChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    updateTaskStatus(task.id, e.target.value as TaskStatus)
  }
  
  return (
    <div
      ref={setNodeRef}
      style={style}
      className={`bg-white rounded-lg border border-gray-200 p-4 ${
        task.status === 'completed' ? 'opacity-70' : ''
      }`}
    >
      <div className="flex items-start gap-3">
        <div
          {...attributes}
          {...listeners}
          className="cursor-grab active:cursor-grabbing pt-1"
        >
          <GripVertical className="w-4 h-4 text-gray-400" />
        </div>
        
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2 mb-1">
            <h4 className={`font-medium text-gray-800 ${
              task.status === 'completed' ? 'line-through text-gray-500' : ''
            }`}>
              {task.title}
            </h4>
            <Badge
              text={getStatusLabel(task.status)}
              color={`${getStatusColor(task.status)} text-white`}
              size="sm"
            />
          </div>
          
          {task.description && (
            <p className="text-sm text-gray-500 mb-2">{task.description}</p>
          )}
          
          <div className="flex flex-wrap items-center gap-2 text-xs">
            <Badge text={getPriorityLabel(task.priority)} color={getPriorityColor(task.priority)} size="sm" />
            
            {task.dueDate && (
              <span className={`px-2 py-0.5 rounded-full ${
                overdue ? 'bg-red-100 text-red-700' :
                dueToday ? 'bg-yellow-100 text-yellow-700' :
                'bg-gray-100 text-gray-600'
              }`}>
                {overdue ? '已逾期' : dueToday ? '今天截止' : formatDateReadable(task.dueDate)}
              </span>
            )}
            
            {taskMilestones.length > 0 && (
              <span className="text-gray-500">
                里程碑: {taskMilestones.map((m) => m.title).join(', ')}
              </span>
            )}
          </div>
        </div>
        
        <div className="flex items-center gap-2">
          <select
            value={task.status}
            onChange={handleStatusChange}
            className="text-sm border border-gray-300 rounded-lg px-2 py-1 focus:outline-none focus:ring-1 focus:ring-blue-500"
          >
            <option value="pending">待开始</option>
            <option value="in-progress">进行中</option>
            <option value="completed">已完成</option>
            <option value="on-hold">已搁置</option>
          </select>
          
          <button
            onClick={() => onEdit(task.id)}
            className="p-1.5 hover:bg-gray-100 rounded-lg"
          >
            <Edit2 className="w-4 h-4 text-gray-400" />
          </button>
          
          <button
            onClick={() => onDelete(task.id)}
            className="p-1.5 hover:bg-red-50 rounded-lg"
          >
            <Trash2 className="w-4 h-4 text-red-400" />
          </button>
        </div>
      </div>
    </div>
  )
}
