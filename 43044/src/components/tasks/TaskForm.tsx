import { useState, useEffect } from 'react'
import { Button } from '@/components/ui/Button'
import { Input, TextArea, Select } from '@/components/ui/Form'
import { Modal } from '@/components/ui/Modal'
import { useAppStore } from '@/store'
import type { Priority, TaskStatus } from '@/types'

interface TaskFormProps {
  isOpen: boolean
  onClose: () => void
  goalId: string
  taskId?: string | null
}

export function TaskForm({ isOpen, onClose, goalId, taskId }: TaskFormProps) {
  const { goals, milestones, tasks, addTask, updateTask } = useAppStore()
  const goal = goals.find((g) => g.id === goalId)
  const existingTask = taskId ? tasks.find((t) => t.id === taskId) : null
  
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [status, setStatus] = useState<TaskStatus>('pending')
  const [priority, setPriority] = useState<Priority>('medium')
  const [dueDate, setDueDate] = useState('')
  const [selectedMilestones, setSelectedMilestones] = useState<string[]>([])
  
  useEffect(() => {
    if (existingTask) {
      setTitle(existingTask.title)
      setDescription(existingTask.description || '')
      setStatus(existingTask.status)
      setPriority(existingTask.priority)
      setDueDate(existingTask.dueDate || '')
      setSelectedMilestones(existingTask.milestoneIds)
    } else {
      setTitle('')
      setDescription('')
      setStatus('pending')
      setPriority('medium')
      setDueDate('')
      setSelectedMilestones([])
    }
  }, [existingTask, isOpen])
  
  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!title.trim()) return
    
    if (existingTask) {
      updateTask(existingTask.id, {
        title: title.trim(),
        description: description.trim() || undefined,
        status,
        priority,
        dueDate: dueDate || undefined,
        milestoneIds: selectedMilestones,
      })
    } else {
      addTask(goalId, selectedMilestones, {
        title: title.trim(),
        description: description.trim() || undefined,
        status,
        priority,
        dueDate: dueDate || undefined,
      })
    }
    
    onClose()
  }
  
  const toggleMilestone = (milestoneId: string) => {
    setSelectedMilestones((prev) =>
      prev.includes(milestoneId)
        ? prev.filter((id) => id !== milestoneId)
        : [...prev, milestoneId]
    )
  }
  
  const goalMilestones = milestones.filter((m) => goal?.milestoneIds.includes(m.id))
  
  return (
    <Modal isOpen={isOpen} onClose={onClose} title={existingTask ? '编辑任务' : '添加任务'}>
      <form onSubmit={handleSubmit} className="space-y-4">
        <Input
          label="任务名称"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          placeholder="请输入任务名称"
          required
        />
        
        <TextArea
          label="任务描述"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          placeholder="请输入任务描述（可选）"
        />
        
        <div className="grid grid-cols-2 gap-4">
          <Select
            label="状态"
            value={status}
            onChange={(e) => setStatus(e.target.value as TaskStatus)}
            options={[
              { value: 'pending', label: '待开始' },
              { value: 'in-progress', label: '进行中' },
              { value: 'completed', label: '已完成' },
              { value: 'on-hold', label: '已搁置' },
            ]}
          />
          
          <Select
            label="优先级"
            value={priority}
            onChange={(e) => setPriority(e.target.value as Priority)}
            options={[
              { value: 'low', label: '低' },
              { value: 'medium', label: '中' },
              { value: 'high', label: '高' },
              { value: 'urgent', label: '紧急' },
            ]}
          />
        </div>
        
        <Input
          label="截止日期"
          type="date"
          value={dueDate}
          onChange={(e) => setDueDate(e.target.value)}
        />
        
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">所属里程碑</label>
          {goalMilestones.length > 0 ? (
            <div className="flex flex-wrap gap-2">
              {goalMilestones.map((milestone) => (
                <button
                  key={milestone.id}
                  type="button"
                  onClick={() => toggleMilestone(milestone.id)}
                  className={`px-3 py-1 rounded-full text-sm font-medium transition-colors ${
                    selectedMilestones.includes(milestone.id)
                      ? 'bg-blue-500 text-white'
                      : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                  }`}
                >
                  {milestone.title}
                </button>
              ))}
            </div>
          ) : (
            <p className="text-sm text-gray-500">暂无里程碑，请先创建里程碑</p>
          )}
        </div>
        
        <div className="flex gap-3 pt-4">
          <Button type="button" variant="secondary" onClick={onClose} className="flex-1">
            取消
          </Button>
          <Button type="submit" className="flex-1">
            {existingTask ? '保存修改' : '添加任务'}
          </Button>
        </div>
      </form>
    </Modal>
  )
}
