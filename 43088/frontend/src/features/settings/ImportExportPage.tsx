import { useState, useRef } from 'react'
import { Upload, Download, FileText, Radio, FileSpreadsheet } from 'lucide-react'
import axios from 'axios'
import ProgressBar from '@/components/common/ProgressBar'

interface ProgressState {
  isUploading: boolean
  uploadProgress: number
  isDownloading: boolean
  downloadProgress: number
  currentOperation: string
}

export default function ImportExportPage() {
  const fileInputRef = useRef<HTMLInputElement>(null)
  const [importResult, setImportResult] = useState<string>('')
  const [progress, setProgress] = useState<ProgressState>({
    isUploading: false,
    uploadProgress: 0,
    isDownloading: false,
    downloadProgress: 0,
    currentOperation: '',
  })

  const handleExportOPML = async () => {
    try {
      setProgress({
        isUploading: false,
        uploadProgress: 0,
        isDownloading: true,
        downloadProgress: 0,
        currentOperation: '正在导出 OPML...',
      })

      const response = await axios.get('/api/export/opml', {
        responseType: 'blob',
        onDownloadProgress: (progressEvent) => {
          if (progressEvent.total) {
            const percentCompleted = Math.round(
              (progressEvent.loaded * 100) / progressEvent.total
            )
            setProgress((prev) => ({
              ...prev,
              downloadProgress: percentCompleted,
            }))
          }
        },
      })

      const blob = new Blob([response.data], { type: 'application/xml' })
      const url = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `podcasts_${new Date().toISOString().split('T')[0]}.opml`
      a.click()
      window.URL.revokeObjectURL(url)
    } catch (err) {
      console.error('Failed to export OPML:', err)
    } finally {
      setProgress({
        isUploading: false,
        uploadProgress: 0,
        isDownloading: false,
        downloadProgress: 0,
        currentOperation: '',
      })
    }
  }

  const handleImportOPML = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (!file) return

    const formData = new FormData()
    formData.append('file', file)

    try {
      setProgress({
        isUploading: true,
        uploadProgress: 0,
        isDownloading: false,
        downloadProgress: 0,
        currentOperation: '正在导入 OPML...',
      })
      setImportResult('')

      const response = await axios.post('/api/import/opml', formData, {
        onUploadProgress: (progressEvent) => {
          if (progressEvent.total) {
            const percentCompleted = Math.round(
              (progressEvent.loaded * 100) / progressEvent.total
            )
            setProgress((prev) => ({
              ...prev,
              uploadProgress: percentCompleted,
            }))
          }
        },
      })

      const result = response.data.data
      setImportResult(
        `成功导入 ${result.imported_count} 个播客，跳过 ${result.skipped_count || 0} 个已存在的播客`
      )
      setTimeout(() => setImportResult(''), 8000)
    } catch (err) {
      setImportResult('导入失败，请检查文件格式')
      setTimeout(() => setImportResult(''), 5000)
    } finally {
      setProgress({
        isUploading: false,
        uploadProgress: 0,
        isDownloading: false,
        downloadProgress: 0,
        currentOperation: '',
      })
    }

    if (fileInputRef.current) {
      fileInputRef.current.value = ''
    }
  }

  const handleExportHistoryCSV = async () => {
    try {
      setProgress({
        isUploading: false,
        uploadProgress: 0,
        isDownloading: true,
        downloadProgress: 0,
        currentOperation: '正在导出收听历史...',
      })

      const response = await axios.get('/api/export/history/csv', {
        responseType: 'blob',
        onDownloadProgress: (progressEvent) => {
          if (progressEvent.total) {
            const percentCompleted = Math.round(
              (progressEvent.loaded * 100) / progressEvent.total
            )
            setProgress((prev) => ({
              ...prev,
              downloadProgress: percentCompleted,
            }))
          }
        },
      })

      const blob = new Blob([response.data], { type: 'text/csv;charset=utf-8' })
      const url = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `listening_history_${new Date().toISOString().split('T')[0]}.csv`
      a.click()
      window.URL.revokeObjectURL(url)
    } catch (err) {
      console.error('Failed to export history CSV:', err)
    } finally {
      setProgress({
        isUploading: false,
        uploadProgress: 0,
        isDownloading: false,
        downloadProgress: 0,
        currentOperation: '',
      })
    }
  }

  const handleExportNotesCSV = async () => {
    try {
      setProgress({
        isUploading: false,
        uploadProgress: 0,
        isDownloading: true,
        downloadProgress: 0,
        currentOperation: '正在导出笔记...',
      })

      const response = await axios.get('/api/export/notes/csv', {
        responseType: 'blob',
        onDownloadProgress: (progressEvent) => {
          if (progressEvent.total) {
            const percentCompleted = Math.round(
              (progressEvent.loaded * 100) / progressEvent.total
            )
            setProgress((prev) => ({
              ...prev,
              downloadProgress: percentCompleted,
            }))
          }
        },
      })

      const blob = new Blob([response.data], { type: 'text/csv;charset=utf-8' })
      const url = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `notes_${new Date().toISOString().split('T')[0]}.csv`
      a.click()
      window.URL.revokeObjectURL(url)
    } catch (err) {
      console.error('Failed to export notes CSV:', err)
    } finally {
      setProgress({
        isUploading: false,
        uploadProgress: 0,
        isDownloading: false,
        downloadProgress: 0,
        currentOperation: '',
      })
    }
  }

  const isBusy = progress.isUploading || progress.isDownloading

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-gray-900">数据导入导出</h1>

      {progress.currentOperation && (
        <div className="card p-4">
          <p className="text-sm text-gray-600 mb-2">{progress.currentOperation}</p>
          {progress.isUploading && (
            <ProgressBar progress={progress.uploadProgress} label="上传进度" />
          )}
          {progress.isDownloading && (
            <ProgressBar progress={progress.downloadProgress} label="下载进度" />
          )}
        </div>
      )}

      <div className="card p-6">
        <h2 className="text-lg font-semibold mb-4 flex items-center gap-2">
          <Radio className="w-5 h-5" />
          播客订阅 (OPML)
        </h2>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <p className="text-sm text-gray-600 mb-3">
              导出所有播客订阅为 OPML 文件，可用于在其他播客应用中导入
            </p>
            <button
              onClick={handleExportOPML}
              disabled={isBusy}
              className="btn btn-secondary flex items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <Download className="w-4 h-4" />
              导出 OPML
            </button>
          </div>
          <div>
            <p className="text-sm text-gray-600 mb-3">
              从 OPML 文件导入播客订阅
            </p>
            <input
              type="file"
              ref={fileInputRef}
              accept=".opml,.xml"
              onChange={handleImportOPML}
              className="hidden"
            />
            <button
              onClick={() => fileInputRef.current?.click()}
              disabled={isBusy}
              className="btn btn-primary flex items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <Upload className="w-4 h-4" />
              导入 OPML
            </button>
            {importResult && (
              <p
                className={`text-sm mt-2 ${
                  importResult.includes('失败') ? 'text-red-600' : 'text-green-600'
                }`}
              >
                {importResult}
              </p>
            )}
          </div>
        </div>
      </div>

      <div className="card p-6">
        <h2 className="text-lg font-semibold mb-4 flex items-center gap-2">
          <FileSpreadsheet className="w-5 h-5" />
          收听历史 (CSV)
        </h2>
        <p className="text-sm text-gray-600 mb-3">
          导出所有收听历史记录为 CSV 文件
        </p>
        <button
          onClick={handleExportHistoryCSV}
          disabled={isBusy}
          className="btn btn-secondary flex items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <Download className="w-4 h-4" />
          导出收听历史 CSV
        </button>
      </div>

      <div className="card p-6">
        <h2 className="text-lg font-semibold mb-4 flex items-center gap-2">
          <FileText className="w-5 h-5" />
          笔记数据 (CSV)
        </h2>
        <p className="text-sm text-gray-600 mb-3">
          导出所有笔记为 CSV 文件
        </p>
        <button
          onClick={handleExportNotesCSV}
          disabled={isBusy}
          className="btn btn-secondary flex items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <Download className="w-4 h-4" />
          导出笔记 CSV
        </button>
      </div>

      <div className="card p-6 bg-amber-50 border border-amber-200">
        <h3 className="text-amber-800 font-medium mb-2">注意事项</h3>
        <ul className="text-sm text-amber-700 space-y-1">
          <li>• OPML 导入会自动跳过已存在的播客订阅</li>
          <li>• CSV 文件使用 UTF-8 编码，支持在 Excel 或其他表格软件中打开</li>
          <li>• 建议定期导出数据进行备份</li>
        </ul>
      </div>
    </div>
  )
}
