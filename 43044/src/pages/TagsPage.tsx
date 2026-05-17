import { useState } from 'react'
import { Plus } from 'lucide-react'
import { useAppStore } from '@/store'
import { Button } from '@/components/ui/Button'
import { Input } from '@/components/ui/Form'
import { Modal } from '@/components/ui/Modal'

const PRESET_COLORS = [
  '#EF4444', '#F97316', '#F59E0B', '#EAB308', '#84CC16',
  '#22C55E', '#10B981', '#14B8A6', '#06B6D4', '#0EA5E9',
  '#3B82F6', '#6366F1', '#8B5CF6', '#A855F7', '#D946EF',
  '#EC4899', '#F43F5E', '#78716C', '#6B7280', '#4B5563',
]

export function TagsPage() {
  const { tags, addTag, updateTag, deleteTag, goals, tasks } = useAppStore()
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [editingTagId, setEditingTagId] = useState<string | null>(null)
  const [tagName, setTagName] = useState('')
  const [tagColor, setTagColor] = useState(PRESET_COLORS[0])
  
  const handleOpenModal = (tagId?: string) => {
    if (tagId) {
      const tag = tags.find((t) => t.id === tagId)
      if (tag) {
        setEditingTagId(tagId)
        setTagName(tag.name)
        setTagColor(tag.color)
      }
    } else {
      setEditingTagId(null)
      setTagName('')
      setTagColor(PRESET_COLORS[0])
    }
    setIsModalOpen(true)
  }
  
  const handleCloseModal = () => {
    setIsModalOpen(false)
    setEditingTagId(null)
    setTagName('')
    setTagColor(PRESET_COLORS[0])
  }
  
  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!tagName.trim()) return
    
    if (editingTagId) {
      updateTag(editingTagId, {
        name: tagName.trim(),
        color: tagColor,
      })
    } else {
      addTag({
        name: tagName.trim(),
        color: tagColor,
      })
    }
    
    handleCloseModal()
  }
  
  const handleDeleteTag = (tagId: string) => {
    const tag = tags.find((t) => t.id === tagId)
    if (!tag) return
    
    const usedInGoals = goals.filter((g) => g.tagIds.includes(tagId)).length
    if (usedInGoals > 0) {
      if (!confirm(`标签"${tag.name}"已被 ${usedInGoals} 个目标使用，确定要删除吗？`)) {
        return
      }
    } else if (!confirm(`确定要删除标签"${tag.name}"吗？`)) {
      return
    }
    
    deleteTag(tagId)
  }
  
  const getTagUsage = (tagId: string) => {
    const tagGoals = goals.filter((g) => g.tagIds.includes(tagId))
    const tagTasks = tasks.filter((t) => {
      const goal = goals.find((g) => g.taskIds.includes(t.id))
      return goal?.tagIds.includes(tagId)
    })
    return { goals: tagGoals.length, tasks: tagTasks.length }
  }
  
  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-800">标签管理</h1>
          <p className="text-gray-500 mt-1">管理目标分类标签</p>
        </div>
        <Button onClick={() => handleOpenModal()}>
          <Plus className="w-4 h-4 mr-2" />
          新建标签
        </Button>
      </div>
      
      {tags.length === 0 ? (
        <div className="bg-white rounded-lg border border-gray-200 p-12 text-center">
          <div className="text-gray-400 mb-4">
            <Plus className="w-12 h-12 mx-auto" />
          </div>
          <h3 className="text-lg font-medium text-gray-700 mb-2">还没有标签</h3>
          <p className="text-gray-500 mb-4">创建标签来分类管理你的目标</p>
          <Button onClick={() => handleOpenModal()}>创建标签</Button>
        </div>
      ) : (
        <div className="bg-white rounded-lg border border-gray-200 overflow-hidden">
          <table className="w-full">
            <thead className="bg-gray-50 border-b border-gray-200">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  标签
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  关联目标
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  关联任务
                </th>
                <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                  操作
                </th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200">
              {tags.map((tag) => {
                const usage = getTagUsage(tag.id)
                return (
                  <tr key={tag.id} className="hover:bg-gray-50">
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="flex items-center gap-3">
                        <span
                          className="w-4 h-4 rounded-full"
                          style={{ backgroundColor: tag.color }}
                        />
                        <span className="font-medium text-gray-800">{tag.name}</span>
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                      {usage.goals} 个
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                      {usage.tasks} 个
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                      <button
                        onClick={() => handleOpenModal(tag.id)}
                        className="text-blue-600 hover:text-blue-800 mr-3"
                      >
                        编辑
                      </button>
                      <button
                        onClick={() => handleDeleteTag(tag.id)}
                        className="text-red-600 hover:text-red-800"
                      >
                        删除
                      </button>
                    </td>
                  </tr>
                )
              })}
            </tbody>
          </table>
        </div>
      )}
      
      <Modal
        isOpen={isModalOpen}
        onClose={handleCloseModal}
        title={editingTagId ? '编辑标签' : '新建标签'}
      >
        <form onSubmit={handleSubmit} className="space-y-4">
          <Input
            label="标签名称"
            value={tagName}
            onChange={(e) => setTagName(e.target.value)}
            placeholder="请输入标签名称"
            required
          />
          
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">选择颜色</label>
            <div className="grid grid-cols-10 gap-2">
              {PRESET_COLORS.map((color) => (
                <button
                  key={color}
                  type="button"
                  onClick={() => setTagColor(color)}
                  className={`w-8 h-8 rounded-full transition-transform hover:scale-110 ${
                    tagColor === color ? 'ring-2 ring-offset-2 ring-gray-400' : ''
                  }`}
                  style={{ backgroundColor: color }}
                />
              ))}
            </div>
          </div>
          
          <div className="flex items-center gap-2 p-3 bg-gray-50 rounded-lg">
            <span
              className="w-4 h-4 rounded-full"
              style={{ backgroundColor: tagColor }}
            />
            <span className="text-sm text-gray-600">
              预览: {tagName || '标签名称'}
            </span>
          </div>
          
          <div className="flex gap-3 pt-4">
            <Button type="button" variant="secondary" onClick={handleCloseModal} className="flex-1">
              取消
            </Button>
            <Button type="submit" className="flex-1">
              {editingTagId ? '保存修改' : '创建标签'}
            </Button>
          </div>
        </form>
      </Modal>
    </div>
  )
}
