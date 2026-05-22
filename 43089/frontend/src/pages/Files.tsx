import { useEffect, useState, useRef } from 'react'
import { useParams, Link } from 'react-router-dom'
import { fileAPI, planAPI } from '@/services/api'
import type { File, PlanDetail } from '@/types'

const fileCategories = [
  { value: 'all', label: '全部', icon: '📁' },
  { value: 'ticket', label: '机票/车票', icon: '✈️' },
  { value: 'hotel', label: '酒店订单', icon: '🏨' },
  { value: 'insurance', label: '保险单据', icon: '🛡️' },
  { value: 'receipt', label: '发票收据', icon: '🧾' },
  { value: 'other', label: '其他', icon: '📄' },
]

const fileTypeIcons: Record<string, string> = {
  'application/pdf': '📕',
  'image/jpeg': '🖼️',
  'image/png': '🖼️',
  'image/gif': '🖼️',
  'application/msword': '📘',
  'application/vnd.openxmlformats-officedocument.wordprocessingml.document': '📘',
  'application/vnd.ms-excel': '📗',
  'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet': '📗',
  'default': '📄',
}

export default function Files() {
  const { id } = useParams<{ id: string }>()
  const [plan, setPlan] = useState<PlanDetail | null>(null)
  const [files, setFiles] = useState<File[]>([])
  const [loading, setLoading] = useState(true)
  const [showUploadModal, setShowUploadModal] = useState(false)
  const [selectedCategory, setSelectedCategory] = useState('all')
  const [uploadProgress, setUploadProgress] = useState(0)
  const [isUploading, setIsUploading] = useState(false)
  const [uploadForm, setUploadForm] = useState({
    category: 'other',
    description: '',
  })
  const fileInputRef = useRef<HTMLInputElement>(null)
  const [selectedFile, setSelectedFile] = useState<globalThis.File | null>(null)

  useEffect(() => {
    if (id) {
      loadData()
    }
  }, [id])

  const loadData = async () => {
    try {
      setLoading(true)
      const [planData, filesData] = await Promise.all([
        planAPI.getPlan(id!),
        fileAPI.getFiles(id!),
      ])
      setPlan(planData)
      setFiles(filesData)
    } catch (error) {
      console.error('Failed to load data:', error)
    } finally {
      setLoading(false)
    }
  }

  const getFilteredFiles = () => {
    if (selectedCategory === 'all') return files
    return files.filter((f) => f.category === selectedCategory)
  }

  const getFileIcon = (fileType: string) => {
    return fileTypeIcons[fileType] || fileTypeIcons['default']
  }

  const formatFileSize = (bytes: number) => {
    if (bytes < 1024) return bytes + ' B'
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
    return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
  }

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (file) {
      setSelectedFile(file)
    }
  }

  const handleUpload = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!selectedFile) {
      alert('请选择要上传的文件')
      return
    }

    const maxSize = 10 * 1024 * 1024 // 10MB
    if (selectedFile.size > maxSize) {
      alert('文件大小不能超过10MB')
      return
    }

    const allowedTypes = ['application/pdf', 'image/jpeg', 'image/png', 'image/gif']
    if (!allowedTypes.includes(selectedFile.type) && !selectedFile.name.match(/\.(pdf|jpg|jpeg|png|gif)$/i)) {
      alert('只支持PDF和图片文件')
      return
    }

    try {
      setIsUploading(true)
      setUploadProgress(0)

      const formData = new FormData()
      formData.append('file', selectedFile)
      formData.append('category', uploadForm.category)
      formData.append('description', uploadForm.description)

      await fileAPI.uploadFile(id!, formData, (progress) => {
        setUploadProgress(progress)
      })

      setShowUploadModal(false)
      resetUploadForm()
      loadData()
    } catch (error: any) {
      alert(error.message || '上传失败')
    } finally {
      setIsUploading(false)
      setUploadProgress(0)
    }
  }

  const resetUploadForm = () => {
    setSelectedFile(null)
    setUploadForm({ category: 'other', description: '' })
    if (fileInputRef.current) {
      fileInputRef.current.value = ''
    }
  }

  const handleDelete = async (fileId: string) => {
    if (!confirm('确定要删除这个文件吗？')) return
    try {
      await fileAPI.deleteFile(id!, fileId)
      loadData()
    } catch (error: any) {
      alert(error.message || '删除失败')
    }
  }

  const handleDownload = async (file: File) => {
    try {
      const blob = await fileAPI.downloadFile(file.id)
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = file.original_name
      a.click()
      URL.revokeObjectURL(url)
    } catch (error: any) {
      alert(error.message || '下载失败')
    }
  }

  const handlePreview = (file: File) => {
    window.open(file.file_url, '_blank')
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      </div>
    )
  }

  const filteredFiles = getFilteredFiles()

  return (
    <div className="p-6">
      <div className="mb-6">
        <Link to={`/plans/${id}`} className="text-primary-600 hover:underline text-sm mb-2 inline-block">
          ← 返回计划详情
        </Link>
        <div className="flex items-center justify-between">
          <h1 className="text-2xl font-bold text-gray-900">文件附件</h1>
          <button
            onClick={() => setShowUploadModal(true)}
            className="px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition-colors"
          >
            + 上传文件
          </button>
        </div>
      </div>

      <div className="flex gap-2 mb-6 overflow-x-auto pb-2">
        {fileCategories.map((cat) => (
          <button
            key={cat.value}
            onClick={() => setSelectedCategory(cat.value)}
            className={`flex-shrink-0 px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
              selectedCategory === cat.value
                ? 'bg-primary-600 text-white'
                : 'bg-white border border-gray-200 hover:bg-gray-50'
            }`}
          >
            <span className="mr-2">{cat.icon}</span>
            {cat.label}
          </button>
        ))}
      </div>

      <div className="bg-white rounded-xl shadow-sm">
        <div className="p-6">
          {filteredFiles.length === 0 ? (
            <div className="text-center py-12 text-gray-500">
              <p className="text-4xl mb-2">📁</p>
              <p>暂无文件</p>
              <button
                onClick={() => setShowUploadModal(true)}
                className="mt-4 text-primary-600 hover:underline"
              >
                + 上传第一个文件
              </button>
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {filteredFiles.map((file) => (
                <div
                  key={file.id}
                  className="border border-gray-200 rounded-lg p-4 hover:shadow-md transition-shadow"
                >
                  <div className="flex items-start gap-3">
                    <div className="text-4xl">{getFileIcon(file.file_type)}</div>
                    <div className="flex-1 min-w-0">
                      <h3 className="font-medium text-gray-900 truncate">{file.original_name}</h3>
                      <p className="text-sm text-gray-500">
                        {formatFileSize(file.file_size)}
                      </p>
                      {file.description && (
                        <p className="text-xs text-gray-400 mt-1 line-clamp-2">{file.description}</p>
                      )}
                    </div>
                  </div>
                  <div className="flex items-center gap-2 mt-4 pt-3 border-t border-gray-100">
                    <button
                      onClick={() => handlePreview(file)}
                      className="flex-1 px-3 py-1.5 text-sm text-primary-600 hover:bg-primary-50 rounded transition-colors"
                    >
                      👁️ 预览
                    </button>
                    <button
                      onClick={() => handleDownload(file)}
                      className="flex-1 px-3 py-1.5 text-sm text-gray-600 hover:bg-gray-50 rounded transition-colors"
                    >
                      ⬇️ 下载
                    </button>
                    <button
                      onClick={() => handleDelete(file.id)}
                      className="px-3 py-1.5 text-sm text-red-600 hover:bg-red-50 rounded transition-colors"
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

      {showUploadModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-xl w-full max-w-lg">
            <div className="p-6">
              <h2 className="text-xl font-bold mb-6">上传文件</h2>
              <form onSubmit={handleUpload} className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">选择文件 *</label>
                  <input
                    ref={fileInputRef}
                    type="file"
                    onChange={handleFileSelect}
                    accept=".pdf,.jpg,.jpeg,.png,.gif,application/pdf,image/*"
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500 file:mr-4 file:py-1.5 file:px-3 file:rounded file:border-0 file:text-sm file:font-medium file:bg-primary-50 file:text-primary-700 hover:file:bg-primary-100"
                  />
                  {selectedFile && (
                    <p className="mt-2 text-sm text-gray-500">
                      已选择: {selectedFile.name} ({formatFileSize(selectedFile.size)})
                    </p>
                  )}
                  <p className="mt-1 text-xs text-gray-400">
                    支持PDF和图片文件，最大10MB
                  </p>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">文件分类</label>
                  <select
                    value={uploadForm.category}
                    onChange={(e) => setUploadForm({ ...uploadForm, category: e.target.value })}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
                  >
                    {fileCategories.filter((c) => c.value !== 'all').map((cat) => (
                      <option key={cat.value} value={cat.value}>
                        {cat.icon} {cat.label}
                      </option>
                    ))}
                  </select>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">描述</label>
                  <textarea
                    value={uploadForm.description}
                    onChange={(e) => setUploadForm({ ...uploadForm, description: e.target.value })}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
                    rows={2}
                    placeholder="文件描述（可选）"
                  />
                </div>
                {isUploading && (
                  <div>
                    <div className="w-full bg-gray-200 rounded-full h-2">
                      <div
                        className="bg-primary-600 h-2 rounded-full transition-all"
                        style={{ width: `${uploadProgress}%` }}
                      ></div>
                    </div>
                    <p className="text-center text-sm text-gray-500 mt-1">上传中... {uploadProgress}%</p>
                  </div>
                )}
                <div className="flex gap-3 pt-4">
                  <button
                    type="button"
                    onClick={() => {
                      setShowUploadModal(false)
                      resetUploadForm()
                    }}
                    disabled={isUploading}
                    className="flex-1 px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors disabled:opacity-50"
                  >
                    取消
                  </button>
                  <button
                    type="submit"
                    disabled={isUploading || !selectedFile}
                    className="flex-1 px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition-colors disabled:opacity-50"
                  >
                    {isUploading ? '上传中...' : '上传'}
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
