import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useState } from 'react'
import { Upload, Bug, Plus, FileText, X } from 'lucide-react'
import { diseaseApi, uploadApi } from '@/api'

export default function DiseasePage() {
  const queryClient = useQueryClient()
  const [showModal, setShowModal] = useState(false)
  const [formData, setFormData] = useState({
    plant_name: '',
    image_url: '',
    description: '',
    symptoms: '',
  })
  const [uploading, setUploading] = useState(false)

  const { data, isLoading } = useQuery({
    queryKey: ['disease-diagnoses'],
    queryFn: () => diseaseApi.getAll(),
  })

  const createMutation = useMutation({
    mutationFn: (data: object) => diseaseApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['disease-diagnoses'] })
      setShowModal(false)
      setFormData({
        plant_name: '',
        image_url: '',
        description: '',
        symptoms: '',
      })
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => diseaseApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['disease-diagnoses'] })
    },
  })

  const handleImageUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (!file) return

    setUploading(true)
    try {
      const response = await uploadApi.upload(file)
      setFormData({ ...formData, image_url: response.data.url })
    } catch (err) {
      alert('上传失败')
    } finally {
      setUploading(false)
    }
  }

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    createMutation.mutate(formData)
  }

  const handleDelete = (id: string) => {
    if (confirm('确定要删除这条诊断记录吗？')) {
      deleteMutation.mutate(id)
    }
  }

  const diagnoses = data?.data?.diagnoses || []

  const severityColors: Record<string, string> = {
    轻微: 'bg-green-100 text-green-700',
    中等: 'bg-amber-100 text-amber-700',
    严重: 'bg-red-100 text-red-700',
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">病虫害诊断</h1>
          <p className="text-gray-500">上传病叶照片或描述症状获取诊断建议</p>
        </div>
        <button onClick={() => setShowModal(true)} className="btn-primary">
          <Plus className="w-4 h-4 mr-2" />
          新建诊断
        </button>
      </div>

      {/* Diagnosis List */}
      {isLoading ? (
        <div className="text-center py-12 text-gray-500">加载中...</div>
      ) : diagnoses.length === 0 ? (
        <div className="card text-center py-12">
          <Bug className="w-16 h-16 text-gray-300 mx-auto mb-4" />
          <p className="text-gray-500 mb-4">还没有诊断记录</p>
          <button onClick={() => setShowModal(true)} className="btn-primary">
            开始诊断
          </button>
        </div>
      ) : (
        <div className="space-y-4">
          {diagnoses.map((diagnosis: any) => (
            <div key={diagnosis.id} className="card">
              <div className="card-body">
                <div className="flex items-start justify-between">
                  <div className="flex items-start gap-4">
                    {diagnosis.image_url && (
                      <img
                        src={diagnosis.image_url}
                        alt="病叶"
                        className="w-24 h-24 object-cover rounded-lg"
                      />
                    )}
                    <div>
                      <h3 className="font-semibold text-gray-900">
                        {diagnosis.plant_name || '未命名植物'}
                      </h3>
                      <p className="text-sm text-gray-500">
                        诊断时间：{new Date(diagnosis.created_at).toLocaleString('zh-CN')}
                      </p>
                      {diagnosis.symptoms && (
                        <p className="text-sm text-gray-600 mt-2">
                          症状：{diagnosis.symptoms}
                        </p>
                      )}
                    </div>
                  </div>
                  <button
                    onClick={() => handleDelete(diagnosis.id)}
                    className="p-2 hover:bg-red-50 rounded-lg"
                  >
                    <X className="w-4 h-4 text-red-400" />
                  </button>
                </div>

                {diagnosis.diagnosis && (
                  <div className="mt-4 p-4 bg-garden-50 rounded-lg">
                    <div className="flex items-center gap-2 mb-2">
                      <FileText className="w-5 h-5 text-garden-600" />
                      <span className="font-medium text-garden-800">诊断结果</span>
                      {diagnosis.severity && (
                        <span className={`badge ${severityColors[diagnosis.severity] || 'bg-gray-100'}`}>
                          {diagnosis.severity}
                        </span>
                      )}
                    </div>
                    <p className="text-gray-700 mb-2">{diagnosis.diagnosis}</p>
                    {diagnosis.treatment && (
                      <div className="mt-2 pt-2 border-t border-garden-200">
                        <p className="text-sm font-medium text-garden-800">防治建议：</p>
                        <p className="text-sm text-gray-600">{diagnosis.treatment}</p>
                      </div>
                    )}
                    {diagnosis.confidence > 0 && (
                      <p className="text-xs text-gray-500 mt-2">
                        置信度：{Math.round(diagnosis.confidence * 100)}%
                      </p>
                    )}
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
              <h2 className="text-lg font-semibold">病虫害诊断</h2>
              <button
                onClick={() => setShowModal(false)}
                className="p-2 hover:bg-gray-100 rounded-lg"
              >
                <X className="w-5 h-5" />
              </button>
            </div>
            <form onSubmit={handleSubmit} className="p-6 space-y-4">
              <div>
                <label className="label">植物名称</label>
                <input
                  type="text"
                  className="input"
                  placeholder="请输入植物名称"
                  value={formData.plant_name}
                  onChange={(e) => setFormData({ ...formData, plant_name: e.target.value })}
                />
              </div>
              <div>
                <label className="label">上传病叶照片</label>
                <div className="border-2 border-dashed border-gray-300 rounded-lg p-6 text-center hover:border-garden-400 transition-colors">
                  <input
                    type="file"
                    accept="image/*"
                    onChange={handleImageUpload}
                    className="hidden"
                    id="image-upload"
                  />
                  <label htmlFor="image-upload" className="cursor-pointer">
                    {formData.image_url ? (
                      <img
                        src={formData.image_url}
                        alt="预览"
                        className="max-h-48 mx-auto rounded-lg"
                      />
                    ) : (
                      <div>
                        <Upload className="w-12 h-12 text-gray-400 mx-auto mb-2" />
                        <p className="text-sm text-gray-500">
                          {uploading ? '上传中...' : '点击上传图片'}
                        </p>
                      </div>
                    )}
                  </label>
                </div>
              </div>
              <div>
                <label className="label">症状描述</label>
                <textarea
                  className="input h-24"
                  placeholder="请描述植物的症状，如叶片变黄、出现斑点等"
                  value={formData.symptoms}
                  onChange={(e) => setFormData({ ...formData, symptoms: e.target.value })}
                />
              </div>
              <div>
                <label className="label">其他说明</label>
                <textarea
                  className="input h-20"
                  placeholder="其他需要说明的情况"
                  value={formData.description}
                  onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                />
              </div>
              <div className="flex gap-3 pt-4">
                <button
                  type="button"
                  onClick={() => setShowModal(false)}
                  className="btn-outline flex-1"
                >
                  取消
                </button>
                <button
                  type="submit"
                  disabled={createMutation.isPending}
                  className="btn-primary flex-1"
                >
                  {createMutation.isPending ? '诊断中...' : '提交诊断'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}
