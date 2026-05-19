import { useState, useEffect } from 'react'
import { Button } from '@/components/ui/Button'
import { Input, TextArea, Select } from '@/components/ui/Form'
import { Modal } from '@/components/ui/Modal'
import { useAppStore } from '@/store'
import type { Priority } from '@/types'

interface GoalFormProps {
  isOpen: boolean
  onClose: () => void
  goalId?: string | null
}

export function GoalForm({ isOpen, onClose, goalId }: GoalFormProps) {
  const { goals, tags, addGoal, updateGoal } = useAppStore()
  const existingGoal = goalId ? goals.find((g) => g.id === goalId) : null
  
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [priority, setPriority] = useState<Priority>('medium')
  const [dueDate, setDueDate] = useState('')
  const [selectedTags, setSelectedTags] = useState<string[]>([])
  
  useEffect(() => {
    if (existingGoal) {
      setTitle(existingGoal.title)
      setDescription(existingGoal.description || '')
      setPriority(existingGoal.priority)
      setDueDate(existingGoal.dueDate || '')
      setSelectedTags(existingGoal.tagIds)
    } else {
      setTitle('')
      setDescription('')
      setPriority('medium')
      setDueDate('')
      setSelectedTags([])
    }
  }, [existingGoal, isOpen])
  
  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!title.trim()) return
    
    if (existingGoal) {
      updateGoal(existingGoal.id, {
        title: title.trim(),
        description: description.trim() || undefined,
        priority,
        dueDate: dueDate || undefined,
        tagIds: selectedTags,
      })
    } else {
      addGoal({
        title: title.trim(),
        description: description.trim() || undefined,
        priority,
        dueDate: dueDate || undefined,
        tagIds: selectedTags,
      })
    }
    
    onClose()
  }
  
  const toggleTag = (tagId: string) => {
    setSelectedTags((prev) =>
      prev.includes(tagId) ? prev.filter((id) => id !== tagId) : [...prev, tagId]
    )
  }
  
  return (
    <Modal isOpen={isOpen} onClose={onClose} title={existingGoal ? '编辑目标' : '创建新目标'}>
      <form onSubmit={handleSubmit} className="space-y-4">
        <Input
          label="目标名称"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          placeholder="请输入目标名称"
          required
        />
        
        <TextArea
          label="目标描述"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          placeholder="请输入目标描述（可选）"
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
        
        <Input
          label="截止日期"
          type="date"
          value={dueDate}
          onChange={(e) => setDueDate(e.target.value)}
        />
        
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">标签</label>
          <div className="flex flex-wrap gap-2">
            {tags.map((tag) => (
              <button
                key={tag.id}
                type="button"
                onClick={() => toggleTag(tag.id)}
                className={`px-3 py-1 rounded-full text-sm font-medium transition-colors ${
                  selectedTags.includes(tag.id)
                    ? 'text-white'
                    : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                }`}
                style={selectedTags.includes(tag.id) ? { backgroundColor: tag.color } : {}}
              >
                {tag.name}
              </button>
            ))}
          </div>
        </div>
        
        <div className="flex gap-3 pt-4">
          <Button type="button" variant="secondary" onClick={onClose} className="flex-1">
            取消
          </Button>
          <Button type="submit" className="flex-1">
            {existingGoal ? '保存修改' : '创建目标'}
          </Button>
        </div>
      </form>
    </Modal>
  )
}
