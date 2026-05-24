import { useState, useEffect } from 'react'
import { Eye, Filter, User, Mail, Phone, FileText } from 'lucide-react'
import { applicationApi, Application, Job } from '@/api/jobs'
import { jobApi } from '@/api/jobs'
import dayjs from 'dayjs'

export default function CompanyApplications() {
  const [applications, setApplications] = useState<Application[]>([])
  const [jobs, setJobs] = useState<Job[]>([])
  const [loading, setLoading] = useState(true)
  const [filterStatus, setFilterStatus] = useState('')
  const [filterJob, setFilterJob] = useState('')
  const [selectedApplication, setSelectedApplication] = useState<Application | null>(null)
  const [showDetailModal, setShowDetailModal] = useState(false)

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    try {
      const [appsRes, jobsRes] = await Promise.all([
        applicationApi.listCompany(),
        jobApi.listCompanyJobs(),
      ])
      setApplications(appsRes.data)
      setJobs(jobsRes.data.data)
    } catch (err) {
      console.error('Failed to load data:', err)
    } finally {
      setLoading(false)
    }
  }

  const filteredApplications = applications.filter((app) => {
    if (filterStatus && app.status !== filterStatus) return false
    if (filterJob && app.job_id !== parseInt(filterJob)) return false
    return true
  })

  const updateStatus = async (id: number, status: string) => {
    try {
      await applicationApi.updateStatus(id, status)
      loadData()
    } catch (err) {
      alert('Failed to update status')
    }
  }

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
      <h1 className="text-2xl font-bold text-gray-900">Applications</h1>

      <div className="flex flex-wrap items-center gap-4">
        <div className="flex items-center gap-2">
          <Filter className="w-5 h-5 text-gray-400" />
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
        </div>
        <select
          value={filterJob}
          onChange={(e) => setFilterJob(e.target.value)}
          className="input-field w-auto"
        >
          <option value="">All Jobs</option>
          {jobs.map((job) => (
            <option key={job.id} value={job.id}>
              {job.title}
            </option>
          ))}
        </select>
      </div>

      {loading ? (
        <div className="text-center py-12">
          <div className="animate-spin w-8 h-8 border-4 border-primary-600 border-t-transparent rounded-full mx-auto"></div>
        </div>
      ) : filteredApplications.length === 0 ? (
        <div className="text-center py-12">
          <FileText className="w-16 h-16 text-gray-300 mx-auto" />
          <p className="mt-4 text-gray-500">No applications found</p>
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
                  <div className="flex items-center gap-3 mt-2 text-sm text-gray-500">
                    <span className="flex items-center gap-1">
                      <User className="w-4 h-4" />
                      {app.jobseeker?.user?.name}
                    </span>
                    <span className="flex items-center gap-1">
                      <Mail className="w-4 h-4" />
                      {app.jobseeker?.user?.email}
                    </span>
                    <span>Applied {dayjs(app.created_at).fromNow()}</span>
                  </div>
                </div>
                <div className="flex items-center gap-2">
                  <button
                    onClick={() => {
                      setSelectedApplication(app)
                      setShowDetailModal(true)
                    }}
                    className="btn-secondary flex items-center gap-2"
                  >
                    <Eye className="w-4 h-4" />
                    View
                  </button>
                  <select
                    value={app.status}
                    onChange={(e) => updateStatus(app.id, e.target.value)}
                    className="text-sm border rounded-lg px-2 py-1"
                  >
                    <option value="pending">Pending</option>
                    <option value="reviewed">Reviewed</option>
                    <option value="interview">Interview</option>
                    <option value="accepted">Accepted</option>
                    <option value="rejected">Rejected</option>
                    <option value="hold">Hold</option>
                  </select>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      {showDetailModal && selectedApplication && (
        <ApplicationDetailModal
          application={selectedApplication}
          onClose={() => setShowDetailModal(false)}
        />
      )}
    </div>
  )
}

function ApplicationDetailModal({
  application,
  onClose,
}: {
  application: Application
  onClose: () => void
}) {
  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-lg shadow-xl max-w-2xl w-full max-h-[90vh] overflow-y-auto">
        <div className="p-6">
          <h2 className="text-2xl font-bold mb-6">Application Details</h2>

          <div className="space-y-6">
            <div className="border-b pb-4">
              <h3 className="font-semibold text-lg">{application.job?.title}</h3>
              <p className="text-gray-500">{application.job?.company?.company_name}</p>
            </div>

            <div>
              <h4 className="font-semibold mb-2">Applicant Information</h4>
              <div className="grid grid-cols-2 gap-4">
                <div className="flex items-center gap-2">
                  <User className="w-4 h-4 text-gray-400" />
                  <span>{application.jobseeker?.user?.name}</span>
                </div>
                <div className="flex items-center gap-2">
                  <Mail className="w-4 h-4 text-gray-400" />
                  <span>{application.jobseeker?.user?.email}</span>
                </div>
                <div className="flex items-center gap-2">
                  <Phone className="w-4 h-4 text-gray-400" />
                  <span>{application.jobseeker?.user?.phone || 'N/A'}</span>
                </div>
              </div>
            </div>

            {application.resume && (
              <div>
                <h4 className="font-semibold mb-2">Resume</h4>
                <div className="bg-gray-50 rounded-lg p-4">
                  <p className="font-medium">{application.resume.title}</p>
                  <p className="text-sm text-gray-500">{application.resume.full_name}</p>
                  {application.resume.summary && (
                    <p className="text-sm text-gray-600 mt-2">{application.resume.summary}</p>
                  )}
                </div>
              </div>
            )}

            {application.cover_letter && (
              <div>
                <h4 className="font-semibold mb-2">Cover Letter</h4>
                <div className="bg-gray-50 rounded-lg p-4">
                  <p className="text-gray-700 whitespace-pre-wrap">{application.cover_letter}</p>
                </div>
              </div>
            )}

            <div>
              <h4 className="font-semibold mb-2">Status</h4>
              <span className={`badge ${
                application.status === 'accepted' ? 'bg-green-100 text-green-800' :
                application.status === 'rejected' ? 'bg-red-100 text-red-800' :
                application.status === 'interview' ? 'bg-purple-100 text-purple-800' :
                'bg-yellow-100 text-yellow-800'
              }`}>
                {application.status}
              </span>
            </div>
          </div>

          <button onClick={onClose} className="btn-primary w-full mt-6">
            Close
          </button>
        </div>
      </div>
    </div>
  )
}
