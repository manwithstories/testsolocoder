import { useState, useEffect } from 'react'
import { Button } from '@/components/ui/Button'
import { Input, TextArea } from '@/components/ui/Form'
import { Modal } from '@/components/ui/Modal'
import { useAppStore } from '@/store'

interface MilestoneFormProps {
  isOpen: boolean
  onClose: () => void
  goalId: string
  milestoneId?: string | null
}

export function MilestoneForm({ isOpen, onClose, goalId, milestoneId }: MilestoneFormProps) {
  const { milestones, addMilestone, updateMilestone } = useAppStore()
  const existingMilestone = milestoneId ? milestones.find((m) => m.id === milestoneId) : null
  
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [dueDate, setDueDate] = useState('')
  
  useEffect(() => {
    if (existingMilestone) {
      setTitle(existingMilestone.title)
      setDescription(existingMilestone.description || '')
      setDueDate(existingMilestone.dueDate || '')
    } else {
      setTitle('')
      setDescription('')
      setDueDate('')
    }
  }, [existingMilestone, isOpen])
  
  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!title.trim()) return
    
    if (existingMilestone) {
      updateMilestone(existingMilestone.id, {
        title: title.trim(),
        description: description.trim() || undefined,
        dueDate: dueDate || undefined,
      })
    } else {
      addMilestone(goalId, {
        title: title.trim(),
        description: description.trim() || undefined,
        dueDate: dueDate || undefined,
        order: 0,
      })
    }
    
    onClose()
  }
  
  return (
    <Modal isOpen={isOpen} onClose={onClose} title={existingMilestone ? '编辑里程碑' : '添加里程碑'}>
      <form onSubmit={handleSubmit} className="space-y-4">
        <Input
          label="里程碑名称"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          placeholder="请输入里程碑名称"
          required
        />
        
        <TextArea
          label="里程碑描述"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          placeholder="请输入里程碑描述（可选）"
        />
        
        <Input
          label="截止日期"
          type="date"
          value={dueDate}
          onChange={(e) => setDueDate(e.target.value)}
        />
        
        <div className="flex gap-3 pt-4">
          <Button type="button" variant="secondary" onClick={onClose} className="flex-1">
            取消
          </Button>
          <Button type="submit" className="flex-1">
            {existingMilestone ? '保存修改' : '添加里程碑'}
          </Button>
        </div>
      </form>
    </Modal>
  )
}
