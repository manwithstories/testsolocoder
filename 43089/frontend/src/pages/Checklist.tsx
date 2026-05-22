import { useEffect, useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { checklistAPI, planAPI } from '@/services/api'
import type { Checklist, ChecklistItem, PlanDetail } from '@/types'

const checklistTypes = [
  { value: 'packing', label: '行李清单', icon: '🧳' },
  { value: 'preparation', label: '准备清单', icon: '✅' },
  { value: 'other', label: '其他清单', icon: '📋' },
]

export default function ChecklistPage() {
  const { id } = useParams<{ id: string }>()
  const [plan, setPlan] = useState<PlanDetail | null>(null)
  const [checklists, setChecklists] = useState<Checklist[]>([])
  const [loading, setLoading] = useState(true)
  const [showChecklistModal, setShowChecklistModal] = useState(false)
  const [showItemModal, setShowItemModal] = useState(false)
  const [selectedChecklistId, setSelectedChecklistId] = useState<string>('')
  const [editingItem, setEditingItem] = useState<ChecklistItem | null>(null)
  const [checklistForm, setChecklistForm] = useState({ title: '', type: 'packing' })
  const [itemForm, setItemForm] = useState({
    title: '',
    description: '',
    category: '',
    quantity: 1,
  })

  useEffect(() => {
    if (id) {
      loadData()
    }
  }, [id])

  const loadData = async () => {
    try {
      setLoading(true)
      const [planData, checklistsData] = await Promise.all([
        planAPI.getPlan(id!),
        checklistAPI.getChecklists(id!),
      ])
      setPlan(planData)
      setChecklists(checklistsData)
    } catch (error) {
      console.error('Failed to load data:', error)
    } finally {
      setLoading(false)
    }
  }

  const getChecklistTypeInfo = (type: string) => {
    return checklistTypes.find((t) => t.value === type) || checklistTypes[2]
  }

  const getCompletionStats = (checklist: Checklist) => {
    const total = checklist.items.length
    const completed = checklist.items.filter((item) => item.is_completed).length
    const percentage = total > 0 ? (completed / total) * 100 : 0
    return { total, completed, percentage }
  }

  const handleCreateChecklist = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await checklistAPI.createChecklist(id!, checklistForm)
      setShowChecklistModal(false)
      setChecklistForm({ title: '', type: 'packing' })
      loadData()
    } catch (error: any) {
      alert(error.message || '创建失败')
    }
  }

  const handleDeleteChecklist = async (checklistId: string) => {
    if (!confirm('确定要删除这个清单吗？所有清单项也会被删除。')) return
    try {
      await checklistAPI.deleteChecklist(id!, checklistId)
      loadData()
    } catch (error: any) {
      alert(error.message || '删除失败')
    }
  }

  const handleAddItem = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!selectedChecklistId) return
    try {
      if (editingItem) {
        await checklistAPI.updateItem(id!, selectedChecklistId, editingItem.id, itemForm)
      } else {
        await checklistAPI.addItem(id!, selectedChecklistId, itemForm)
      }
      setShowItemModal(false)
      resetItemForm()
      loadData()
    } catch (error: any) {
      alert(error.message || '保存失败')
    }
  }

  const handleEditItem = (checklistId: string, item: ChecklistItem) => {
    setSelectedChecklistId(checklistId)
    setEditingItem(item)
    setItemForm({
      title: item.title,
      description: item.description || '',
      category: item.category || '',
      quantity: item.quantity || 1,
    })
    setShowItemModal(true)
  }

  const handleDeleteItem = async (checklistId: string, itemId: string) => {
    if (!confirm('确定要删除这个项目吗？')) return
    try {
      await checklistAPI.deleteItem(id!, checklistId, itemId)
      loadData()
    } catch (error: any) {
      alert(error.message || '删除失败')
    }
  }

  const handleToggleItem = async (checklistId: string, item: ChecklistItem) => {
    try {
      await checklistAPI.updateItem(id!, checklistId, item.id, {
        is_completed: !item.is_completed,
      })
      loadData()
    } catch (error: any) {
      alert(error.message || '更新失败')
    }
  }

  const resetItemForm = () => {
    setEditingItem(null)
    setItemForm({ title: '', description: '', category: '', quantity: 1 })
  }

  const openAddItemModal = (checklistId: string) => {
    setSelectedChecklistId(checklistId)
    resetItemForm()
    setShowItemModal(true)
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      </div>
    )
  }

  return (
    <div className="p-6">
      <div className="mb-6">
        <Link to={`/plans/${id}`} className="text-primary-600 hover:underline text-sm mb-2 inline-block">
          ← 返回计划详情
        </Link>
        <div className="flex items-center justify-between">
          <h1 className="text-2xl font-bold text-gray-900">清单管理</h1>
          <button
            onClick={() => setShowChecklistModal(true)}
            className="px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition-colors"
          >
            + 创建清单
          </button>
        </div>
      </div>

      {checklists.length === 0 ? (
        <div className="bg-white rounded-xl shadow-sm p-12 text-center">
          <p className="text-4xl mb-4">📋</p>
          <p className="text-gray-500 mb-4">暂无清单</p>
          <button
            onClick={() => setShowChecklistModal(true)}
            className="text-primary-600 hover:underline"
          >
            + 创建第一个清单
          </button>
        </div>
      ) : (
        <div className="space-y-6">
          {checklists.map((checklist) => {
            const typeInfo = getChecklistTypeInfo(checklist.type)
            const stats = getCompletionStats(checklist)
            return (
              <div key={checklist.id} className="bg-white rounded-xl shadow-sm">
                <div className="p-6">
                  <div className="flex items-start justify-between mb-4">
                    <div>
                      <div className="flex items-center gap-3">
                        <span className="text-2xl">{typeInfo.icon}</span>
                        <h2 className="text-xl font-bold text-gray-900">{checklist.title}</h2>
                        <span className="text-sm text-gray-500">
                          ({stats.completed}/{stats.total})
                        </span>
                      </div>
                      <p className="text-sm text-gray-500 mt-1">{typeInfo.label}</p>
                    </div>
                    <div className="flex items-center gap-2">
                      <button
                        onClick={() => openAddItemModal(checklist.id)}
                        className="px-3 py-1.5 text-sm bg-primary-50 text-primary-600 rounded-lg hover:bg-primary-100 transition-colors"
                      >
                        + 添加项目
                      </button>
                      <button
                        onClick={() => handleDeleteChecklist(checklist.id)}
                        className="p-2 text-gray-500 hover:text-red-600 hover:bg-gray-100 rounded"
                      >
                        🗑️
                      </button>
                    </div>
                  </div>

                  <div className="w-full bg-gray-200 rounded-full h-2 mb-4">
                    <div
                      className={`h-2 rounded-full transition-all ${
                        stats.percentage === 100 ? 'bg-green-500' : 'bg-primary-500'
                      }`}
                      style={{ width: `${stats.percentage}%` }}
                    ></div>
                  </div>

                  {checklist.items.length === 0 ? (
                    <div className="text-center py-8 text-gray-500">
                      <p>暂无项目</p>
                      <button
                        onClick={() => openAddItemModal(checklist.id)}
                        className="text-primary-600 hover:underline text-sm"
                      >
                        + 添加第一个项目
                      </button>
                    </div>
                  ) : (
                    <div className="space-y-2">
                      {checklist.items
                        .sort((a, b) => (a.order_index || 0) - (b.order_index || 0))
                        .map((item) => (
                          <div
                            key={item.id}
                            className={`flex items-start gap-3 p-3 rounded-lg border transition-colors ${
                              item.is_completed
                                ? 'bg-green-50 border-green-200'
                                : 'bg-white border-gray-200 hover:bg-gray-50'
                            }`}
                          >
                            <button
                              onClick={() => handleToggleItem(checklist.id, item)}
                              className={`w-5 h-5 rounded border-2 flex items-center justify-center flex-shrink-0 mt-0.5 ${
                                item.is_completed
                                  ? 'bg-green-500 border-green-500 text-white'
                                  : 'border-gray-300 hover:border-primary-500'
                              }`}
                            >
                              {item.is_completed && '✓'}
                            </button>
                            <div className="flex-1 min-w-0">
                              <p
                                className={`font-medium ${
                                  item.is_completed
                                    ? 'text-gray-400 line-through'
                                    : 'text-gray-900'
                                }`}
                              >
                                {item.title}
                                {item.quantity > 1 && (
                                  <span className="text-sm text-gray-500 ml-2">
                                    ×{item.quantity}
                                  </span>
                                )}
                              </p>
                              {item.description && (
                                <p className="text-sm text-gray-500 mt-0.5">{item.description}</p>
                              )}
                              {item.category && (
                                <span className="inline-block mt-1 px-2 py-0.5 bg-gray-100 text-gray-600 text-xs rounded">
                                  {item.category}
                                </span>
                              )}
                            </div>
                            <div className="flex gap-1">
                              <button
                                onClick={() => handleEditItem(checklist.id, item)}
                                className="p-1.5 text-gray-500 hover:text-primary-600 hover:bg-gray-100 rounded"
                              >
                                ✏️
                              </button>
                              <button
                                onClick={() => handleDeleteItem(checklist.id, item.id)}
                                className="p-1.5 text-gray-500 hover:text-red-600 hover:bg-gray-100 rounded"
                              >
                                🗑️
                              </button>
                            </div>
                          </div>
                        ))}
                    </div>
                  )}
                </div>
              </div>
            )
          })}
        </div>
      )}

      {showChecklistModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-xl w-full max-w-md">
            <div className="p-6">
              <h2 className="text-xl font-bold mb-6">创建清单</h2>
              <form onSubmit={handleCreateChecklist} className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">清单名称 *</label>
                  <input
                    type="text"
                    required
                    value={checklistForm.title}
                    onChange={(e) => setChecklistForm({ ...checklistForm, title: e.target.value })}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
                    placeholder="输入清单名称"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">清单类型</label>
                  <select
                    value={checklistForm.type}
                    onChange={(e) => setChecklistForm({ ...checklistForm, type: e.target.value })}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
                  >
                    {checklistTypes.map((type) => (
                      <option key={type.value} value={type.value}>
                        {type.icon} {type.label}
                      </option>
                    ))}
                  </select>
                </div>
                <div className="flex gap-3 pt-4">
                  <button
                    type="button"
                    onClick={() => setShowChecklistModal(false)}
                    className="flex-1 px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
                  >
                    取消
                  </button>
                  <button
                    type="submit"
                    className="flex-1 px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition-colors"
                  >
                    创建
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>
      )}

      {showItemModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-xl w-full max-w-md">
            <div className="p-6">
              <h2 className="text-xl font-bold mb-6">
                {editingItem ? '编辑项目' : '添加项目'}
              </h2>
              <form onSubmit={handleAddItem} className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">项目名称 *</label>
                  <input
                    type="text"
                    required
                    value={itemForm.title}
                    onChange={(e) => setItemForm({ ...itemForm, title: e.target.value })}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
                    placeholder="输入项目名称"
                  />
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">分类</label>
                    <input
                      type="text"
                      value={itemForm.category}
                      onChange={(e) => setItemForm({ ...itemForm, category: e.target.value })}
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
                      placeholder="如：衣物、证件"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">数量</label>
                    <input
                      type="number"
                      min="1"
                      value={itemForm.quantity}
                      onChange={(e) => setItemForm({ ...itemForm, quantity: Number(e.target.value) })}
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
                    />
                  </div>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">描述</label>
                  <textarea
                    value={itemForm.description}
                    onChange={(e) => setItemForm({ ...itemForm, description: e.target.value })}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
                    rows={2}
                    placeholder="项目描述（可选）"
                  />
                </div>
                <div className="flex gap-3 pt-4">
                  <button
                    type="button"
                    onClick={() => {
                      setShowItemModal(false)
                      resetItemForm()
                    }}
                    className="flex-1 px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
                  >
                    取消
                  </button>
                  <button
                    type="submit"
                    className="flex-1 px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition-colors"
                  >
                    {editingItem ? '保存修改' : '添加'}
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
