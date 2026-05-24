import { useState } from 'react'
import { Download, FileText, Calendar, Briefcase } from 'lucide-react'
import { exportApi } from '@/api/jobs'

export default function CompanyExport() {
  const [exporting, setExporting] = useState<string | null>(null)
  const [filterJobId, setFilterJobId] = useState('')

  const handleExport = async (type: 'applications' | 'interviews' | 'jobs') => {
    setExporting(type)
    try {
      const params: Record<string, string> = {}
      if (type === 'applications' && filterJobId) {
        params.job_id = filterJobId
      }

      let response
      if (type === 'applications') {
        response = await exportApi.applications(params)
      } else if (type === 'interviews') {
        response = await exportApi.interviews()
      } else {
        response = await exportApi.jobs()
      }

      const blob = new Blob([response], {
        type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
      })
      const url = window.URL.createObjectURL(blob)
      const link = document.createElement('a')
      link.href = url
      link.download = `${type}_${new Date().toISOString().split('T')[0]}.xlsx`
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
      window.URL.revokeObjectURL(url)
    } catch (err: any) {
      alert(err.message || 'Export failed')
    } finally {
      setExporting(null)
    }
  }

  const exportItems = [
    {
      id: 'applications' as const,
      title: 'Export Applications',
      description: 'Download all application records as Excel file',
      icon: FileText,
      color: 'blue',
    },
    {
      id: 'interviews' as const,
      title: 'Export Interviews',
      description: 'Download all interview records as Excel file',
      icon: Calendar,
      color: 'purple',
    },
    {
      id: 'jobs' as const,
      title: 'Export Jobs',
      description: 'Download all job postings as Excel file',
      icon: Briefcase,
      color: 'green',
    },
  ]

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-gray-900">Data Export</h1>

      <p className="text-gray-500">
        Export your recruitment data to Excel format for reporting and analysis purposes.
      </p>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {exportItems.map((item) => (
          <div key={item.id} className="card hover:shadow-lg transition-shadow">
            <div className={`p-3 rounded-lg w-fit mb-4 bg-${item.color}-100`}>
              <item.icon className={`w-8 h-8 text-${item.color}-600`} />
            </div>
            <h3 className="text-lg font-semibold mb-2">{item.title}</h3>
            <p className="text-gray-500 text-sm mb-4">{item.description}</p>
            <button
              onClick={() => handleExport(item.id)}
              disabled={exporting === item.id}
              className="btn-primary w-full flex items-center justify-center gap-2 disabled:opacity-50"
            >
              <Download className="w-4 h-4" />
              {exporting === item.id ? 'Exporting...' : 'Export'}
            </button>
          </div>
        ))}
      </div>

      <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4 mt-8">
        <h3 className="font-semibold text-yellow-800 mb-2">Note</h3>
        <p className="text-yellow-700 text-sm">
          Exports are generated in Excel (.xlsx) format. The data includes all relevant records
          based on your current filters. Please ensure you have the necessary permissions before
          sharing exported data.
        </p>
      </div>
    </div>
  )
}
