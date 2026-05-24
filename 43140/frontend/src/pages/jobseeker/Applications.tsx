import { useState, useEffect } from 'react'
import { FileText, MapPin, Trash2, Eye } from 'lucide-react'
import { applicationApi, Application } from '@/api/jobs'
import dayjs from 'dayjs'

export default function JobSeekerApplications() {
  const [applications, setApplications] = useState<Application[]>([])
  const [loading, setLoading] = useState(true)
  const [filterStatus, setFilterStatus] = useState('')

  useEffect(() => {
    loadApplications()
  }, [])

  const loadApplications = async () => {
    try {
      const response = await applicationApi.listJobSeeker()
      setApplications(response.data)
    } catch (err) {
      console.error('Failed to load applications:', err)
    } finally {
      setLoading(false)
    }
  }

  const handleDelete = async (id: number) => {
    if (confirm('Are you sure you want to withdraw this application?')) {
      try {
        await applicationApi.delete(id)
        loadApplications()
      } catch (err) {
        alert('Failed to withdraw application')
      }
    }
  }

  const filteredApplications = applications.filter(
    (app) => !filterStatus || app.status === filterStatus
  )

  const getStatusBadge = (status: string) => {
    const colors: Record<string, string> = {
      pending: 'bg-yellow-100 text-yellow-800',
      reviewed: 'bg-blue-100 text-blue-800',
      interview: 'bg-purple-100 text-purple-800',
      accepted: 'bg-green-100 text-green-800',
      rejected: 'bg-red-100 text-red-800',
      hold: 'bg-gray-100 text-gray-800',
    }
    return colors[status] || 'bg-gray-100 text-gray-800'
  }

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-gray-900">My Applications</h1>

      <select
        value={filterStatus}
        onChange={(e) => setFilterStatus(e.target.value)}
        className="input-field w-auto"
      >
        <option value="">All Statuses</option>
        <option value="pending">Pending</option>
        <option value="reviewed">Reviewed</option>
        <option value="interview">Interview</option>
        <option value="accepted">Accepted</option>
        <option value="rejected">Rejected</option>
        <option value="hold">Hold</option>
      </select>

      {loading ? (
        <div className="text-center py-12">
          <div className="animate-spin w-8 h-8 border-4 border-primary-600 border-t-transparent rounded-full mx-auto"></div>
        </div>
      ) : filteredApplications.length === 0 ? (
        <div className="text-center py-12">
          <FileText className="w-16 h-16 text-gray-300 mx-auto" />
          <p className="mt-4 text-gray-500">No applications yet</p>
        </div>
      ) : (
        <div className="space-y-4">
          {filteredApplications.map((app) => (
            <div key={app.id} className="card">
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <div className="flex items-center gap-3">
                    <h3 className="text-lg font-semibold">{app.job?.title}</h3>
                    <span className={`badge ${getStatusBadge(app.status)}`}>
                      {app.status}
                    </span>
                  </div>
                  <p className="text-gray-600">{app.job?.company?.company_name}</p>
                  <div className="flex items-center gap-4 mt-2 text-sm text-gray-500">
                    <span className="flex items-center gap-1">
                      <MapPin className="w-4 h-4" />
                      {app.job?.location}
                    </span>
                    <span>Applied {dayjs(app.created_at).fromNow()}</span>
                  </div>
                </div>
                <div className="flex items-center gap-2">
                  {app.status === 'pending' && (
                    <button
                      onClick={() => handleDelete(app.id)}
                      className="p-2 text-gray-600 hover:text-red-600"
                      title="Withdraw application"
                    >
                      <Trash2 className="w-5 h-5" />
                    </button>
                  )}
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
