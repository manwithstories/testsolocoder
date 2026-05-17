import { useState } from 'react'
import { Plus } from 'lucide-react'
import { useAppStore } from '@/store'
import { Button } from '@/components/ui/Button'
import { GoalCard } from '@/components/goals/GoalCard'
import { GoalForm } from '@/components/goals/GoalForm'
import { useNavigate } from 'react-router-dom'

export function GoalsPage() {
  const { goals, deleteGoal } = useAppStore()
  const navigate = useNavigate()
  const [isFormOpen, setIsFormOpen] = useState(false)
  const [editingGoalId, setEditingGoalId] = useState<string | null>(null)
  
  const handleSelectGoal = (id: string) => {
    navigate(`/goals/${id}`)
  }
  
  const handleEditGoal = (id: string) => {
    setEditingGoalId(id)
    setIsFormOpen(true)
  }
  
  const handleDeleteGoal = (id: string) => {
    if (confirm('确定要删除这个目标吗？相关的里程碑和任务也会被删除。')) {
      deleteGoal(id)
    }
  }
  
  const handleCloseForm = () => {
    setIsFormOpen(false)
    setEditingGoalId(null)
  }
  
  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-800">目标管理</h1>
          <p className="text-gray-500 mt-1">管理你的长期目标和里程碑</p>
        </div>
        <Button onClick={() => setIsFormOpen(true)}>
          <Plus className="w-4 h-4 mr-2" />
          新建目标
        </Button>
      </div>
      
      {goals.length === 0 ? (
        <div className="bg-white rounded-lg border border-gray-200 p-12 text-center">
          <div className="text-gray-400 mb-4">
            <Plus className="w-12 h-12 mx-auto" />
          </div>
          <h3 className="text-lg font-medium text-gray-700 mb-2">还没有目标</h3>
          <p className="text-gray-500 mb-4">点击上方按钮创建你的第一个目标</p>
          <Button onClick={() => setIsFormOpen(true)}>创建目标</Button>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {goals.map((goal) => (
            <GoalCard
              key={goal.id}
              goal={goal}
              onSelect={handleSelectGoal}
              onEdit={handleEditGoal}
              onDelete={handleDeleteGoal}
            />
          ))}
        </div>
      )}
      
      <GoalForm
        isOpen={isFormOpen}
        onClose={handleCloseForm}
        goalId={editingGoalId}
      />
    </div>
  )
}
