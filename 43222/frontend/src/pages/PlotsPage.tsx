import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useState } from 'react'
import { Plus, Edit2, Trash2, MapPin, Sun, Droplets, X } from 'lucide-react'
import { plotApi } from '@/api'
import type { Plot } from '@/types'

export default function PlotsPage() {
  const queryClient = useQueryClient()
  const [showModal, setShowModal] = useState(false)
  const [editingPlot, setEditingPlot] = useState<Plot | null>(null)
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    soil_type: '',
    sunlight: '',
    area: 0,
    location: '',
    grid_config: '',
    irrigation_device: '',
    sensor_data: '',
  })

  const { data, isLoading } = useQuery({
    queryKey: ['plots'],
    queryFn: () => plotApi.getAll(),
  })

  const createMutation = useMutation({
    mutationFn: (data: object) => plotApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['plots'] })
      setShowModal(false)
      resetForm()
    },
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: object }) => plotApi.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['plots'] })
      setShowModal(false)
      resetForm()
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => plotApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['plots'] })
    },
  })

  const resetForm = () => {
    setFormData({
      name: '',
      description: '',
      soil_type: '',
      sunlight: '',
      area: 0,
      location: '',
      grid_config: '',
      irrigation_device: '',
      sensor_data: '',
    })
    setEditingPlot(null)
  }

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (editingPlot) {
      updateMutation.mutate({ id: editingPlot.id, data: formData })
    } else {
      createMutation.mutate(formData)
    }
  }

  const handleEdit = (plot: Plot) => {
    setEditingPlot(plot)
    setFormData({
      name: plot.name,
      description: plot.description,
      soil_type: plot.soil_type,
      sunlight: plot.sunlight,
      area: plot.area,
      location: plot.location,
      grid_config: plot.grid_config,
      irrigation_device: plot.irrigation_device,
      sensor_data: plot.sensor_data,
    })
    setShowModal(true)
  }

  const handleDelete = (id: string) => {
    if (confirm('确定要删除这个地块吗？')) {
      deleteMutation.mutate(id)
    }
  }

  const plots = data?.data?.plots || []

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">我的菜园</h1>
          <p className="text-gray-500">管理您的种植区域</p>
        </div>
        <button
          onClick={() => {
            resetForm()
            setShowModal(true)
          }}
          className="btn-primary"
        >
          <Plus className="w-4 h-4 mr-2" />
          新建地块
        </button>
      </div>

      {isLoading ? (
        <div className="text-center py-12 text-gray-500">加载中...</div>
      ) : plots.length === 0 ? (
        <div className="card text-center py-12">
          <MapPin className="w-16 h-16 text-gray-300 mx-auto mb-4" />
          <p className="text-gray-500 mb-4">还没有创建地块</p>
          <button
            onClick={() => setShowModal(true)}
            className="btn-primary"
          >
            创建第一个地块
          </button>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {plots.map((plot: Plot) => (
            <div key={plot.id} className="card hover:shadow-md transition-shadow">
              <div className="card-body">
                <div className="flex items-start justify-between mb-4">
                  <div>
                    <h3 className="font-semibold text-gray-900">{plot.name}</h3>
                    {plot.description && (
                      <p className="text-sm text-gray-500 mt-1">{plot.description}</p>
                    )}
                  </div>
                  <div className="flex gap-1">
                    <button
                      onClick={() => handleEdit(plot)}
                      className="p-1.5 hover:bg-gray-100 rounded-lg"
                    >
                      <Edit2 className="w-4 h-4 text-gray-400" />
                    </button>
                    <button
                      onClick={() => handleDelete(plot.id)}
                      className="p-1.5 hover:bg-red-50 rounded-lg"
                    >
                      <Trash2 className="w-4 h-4 text-red-400" />
                    </button>
                  </div>
                </div>

                <div className="space-y-2 text-sm">
                  <div className="flex items-center gap-2 text-gray-600">
                    <MapPin className="w-4 h-4" />
                    <span>{plot.location || '未设置位置'}</span>
                  </div>
                  <div className="flex items-center gap-2 text-gray-600">
                    <Sun className="w-4 h-4" />
                    <span>{plot.sunlight || '未设置光照'}</span>
                  </div>
                  <div className="flex items-center gap-2 text-gray-600">
                    <Droplets className="w-4 h-4" />
                    <span>{plot.soil_type || '未设置土壤'}</span>
                  </div>
                  {plot.area > 0 && (
                    <div className="text-gray-600">
                      面积：{plot.area} 平方米
                    </div>
                  )}
                </div>

                {plot.planting_records && plot.planting_records.length > 0 && (
                  <div className="mt-4 pt-4 border-t border-gray-100">
                    <p className="text-sm text-gray-500">
                      种植记录：{plot.planting_records.length} 条
                    </p>
                  </div>
                )}
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Modal */}
      {showModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-2xl w-full max-w-lg max-h-[90vh] overflow-y-auto">
            <div className="flex items-center justify-between p-6 border-b border-gray-200">
              <h2 className="text-lg font-semibold">
                {editingPlot ? '编辑地块' : '新建地块'}
              </h2>
              <button
                onClick={() => {
                  setShowModal(false)
                  resetForm()
                }}
                className="p-2 hover:bg-gray-100 rounded-lg"
              >
                <X className="w-5 h-5" />
              </button>
            </div>
            <form onSubmit={handleSubmit} className="p-6 space-y-4">
              <div>
                <label className="label">地块名称 *</label>
                <input
                  type="text"
                  className="input"
                  required
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                />
              </div>
              <div>
                <label className="label">描述</label>
                <textarea
                  className="input h-24"
                  value={formData.description}
                  onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                />
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="label">土壤类型</label>
                  <select
                    className="input"
                    value={formData.soil_type}
                    onChange={(e) => setFormData({ ...formData, soil_type: e.target.value })}
                  >
                    <option value="">选择土壤类型</option>
                    <option value="sandy">沙质土</option>
                    <option value="loamy">壤土</option>
                    <option value="clay">粘土</option>
                    <option value="silty">粉土</option>
                    <option value="peaty">泥炭土</option>
                  </select>
                </div>
                <div>
                  <label className="label">光照条件</label>
                  <select
                    className="input"
                    value={formData.sunlight}
                    onChange={(e) => setFormData({ ...formData, sunlight: e.target.value })}
                  >
                    <option value="">选择光照条件</option>
                    <option value="full_sun">全日照</option>
                    <option value="partial_sun">半日照</option>
                    <option value="partial_shade">半阴</option>
                    <option value="full_shade">全阴</option>
                  </select>
                </div>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="label">面积（平方米）</label>
                  <input
                    type="number"
                    className="input"
                    value={formData.area}
                    onChange={(e) => setFormData({ ...formData, area: parseFloat(e.target.value) })}
                  />
                </div>
                <div>
                  <label className="label">位置</label>
                  <input
                    type="text"
                    className="input"
                    value={formData.location}
                    onChange={(e) => setFormData({ ...formData, location: e.target.value })}
                  />
                </div>
              </div>
              <div>
                <label className="label">灌溉设备</label>
                <input
                  type="text"
                  className="input"
                  value={formData.irrigation_device}
                  onChange={(e) => setFormData({ ...formData, irrigation_device: e.target.value })}
                />
              </div>
              <div className="flex gap-3 pt-4">
                <button
                  type="button"
                  onClick={() => {
                    setShowModal(false)
                    resetForm()
                  }}
                  className="btn-outline flex-1"
                >
                  取消
                </button>
                <button
                  type="submit"
                  disabled={createMutation.isPending || updateMutation.isPending}
                  className="btn-primary flex-1"
                >
                  {createMutation.isPending || updateMutation.isPending
                    ? '保存中...'
                    : editingPlot
                    ? '保存修改'
                    : '创建'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}
