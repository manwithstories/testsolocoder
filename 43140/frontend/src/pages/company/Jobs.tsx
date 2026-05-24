import { useState, useEffect } from 'react'
import { Plus, Edit2, Trash2, MapPin, Briefcase, DollarSign, Search } from 'lucide-react'
import { jobApi, departmentApi, Job, Department } from '@/api/jobs'
import dayjs from 'dayjs'

export default function CompanyJobs() {
  const [jobs, setJobs] = useState<Job[]>([])
  const [departments, setDepartments] = useState<Department[]>([])
  const [loading, setLoading] = useState(true)
  const [showModal, setShowModal] = useState(false)
  const [editingJob, setEditingJob] = useState<Job | null>(null)
  const [searchTerm, setSearchTerm] = useState('')

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    try {
      const [jobsRes, deptsRes] = await Promise.all([
        jobApi.listCompanyJobs(),
        departmentApi.list(),
      ])
      setJobs(jobsRes.data.data)
      setDepartments(deptsRes.data)
    } catch (err) {
      console.error('Failed to load data:', err)
    } finally {
      setLoading(false)
    }
  }

  const handleDelete = async (id: number) => {
    if (confirm('Are you sure you want to delete this job?')) {
      try {
        await jobApi.delete(id)
        loadData()
      } catch (err) {
        alert('Failed to delete job')
      }
    }
  }

  const handleStatusChange = async (id: number, status: string) => {
    try {
      await jobApi.updateStatus(id, status)
      loadData()
    } catch (err) {
      alert('Failed to update status')
    }
  }

  const filteredJobs = jobs.filter(
    (job) =>
      job.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
      job.location.toLowerCase().includes(searchTerm.toLowerCase())
  )

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-gray-900">Job Management</h1>
        <button onClick={() => setShowModal(true)} className="btn-primary flex items-center gap-2">
          <Plus className="w-5 h-5" />
          Add Job
        </button>
      </div>

      <div className="relative max-w-md">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
        <input
          type="text"
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          placeholder="Search jobs..."
          className="input-field pl-10"
        />
      </div>

      {loading ? (
        <div className="text-center py-12">
          <div className="animate-spin w-8 h-8 border-4 border-primary-600 border-t-transparent rounded-full mx-auto"></div>
        </div>
      ) : filteredJobs.length === 0 ? (
        <div className="text-center py-12">
          <Briefcase className="w-16 h-16 text-gray-300 mx-auto" />
          <p className="mt-4 text-gray-500">No jobs found</p>
        </div>
      ) : (
        <div className="space-y-4">
          {filteredJobs.map((job) => (
            <div key={job.id} className="card">
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <h3 className="text-lg font-semibold">{job.title}</h3>
                  <div className="flex flex-wrap items-center gap-4 mt-2 text-sm text-gray-500">
                    <span className="flex items-center gap-1">
                      <MapPin className="w-4 h-4" />
                      {job.location}
                    </span>
                    <span className="flex items-center gap-1">
                      <Briefcase className="w-4 h-4" />
                      {job.job_type}
                    </span>
                    {job.salary_min > 0 && (
                      <span className="flex items-center gap-1 text-green-600">
                        <DollarSign className="w-4 h-4" />
                        ${job.salary_min.toLocaleString()} - ${job.salary_max.toLocaleString()}
                      </span>
                    )}
                    <span className="text-xs text-gray-400">
                      Updated {dayjs(job.updated_at).fromNow()}
                    </span>
                  </div>
                </div>
                <div className="flex items-center gap-2">
                  <select
                    value={job.status}
                    onChange={(e) => handleStatusChange(job.id, e.target.value)}
                    className="text-sm border rounded-lg px-2 py-1"
                  >
                    <option value="open">Open</option>
                    <option value="paused">Paused</option>
                    <option value="closed">Closed</option>
                  </select>
                  <button
                    onClick={() => {
                      setEditingJob(job)
                      setShowModal(true)
                    }}
                    className="p-2 text-gray-600 hover:text-primary-600"
                  >
                    <Edit2 className="w-5 h-5" />
                  </button>
                  <button
                    onClick={() => handleDelete(job.id)}
                    className="p-2 text-gray-600 hover:text-red-600"
                  >
                    <Trash2 className="w-5 h-5" />
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      {showModal && (
        <JobModal
          job={editingJob}
          departments={departments}
          onClose={() => {
            setShowModal(false)
            setEditingJob(null)
          }}
          onSave={() => {
            setShowModal(false)
            setEditingJob(null)
            loadData()
          }}
        />
      )}
    </div>
  )
}

function JobModal({
  job,
  departments,
  onClose,
  onSave,
}: {
  job: Job | null
  departments: Department[]
  onClose: () => void
  onSave: () => void
}) {
  const [formData, setFormData] = useState({
    title: job?.title || '',
    location: job?.location || '',
    salary_min: job?.salary_min || 0,
    salary_max: job?.salary_max || 0,
    salary_type: job?.salary_type || 'monthly',
    description: job?.description || '',
    requirements: job?.requirements || '',
    skills: job?.skills || '',
    job_type: job?.job_type || 'full-time',
    department_id: job?.department_id || undefined,
    status: job?.status || 'open',
  })
  const [submitting, setSubmitting] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setSubmitting(true)
    try {
      if (job) {
        await jobApi.update(job.id, formData)
      } else {
        await jobApi.create(formData)
      }
      onSave()
    } catch (err: any) {
      alert(err.message || 'Failed to save job')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-lg shadow-xl max-w-2xl w-full max-h-[90vh] overflow-y-auto">
        <div className="p-6">
          <h2 className="text-2xl font-bold mb-6">{job ? 'Edit Job' : 'Create Job'}</h2>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Title *</label>
                <input
                  type="text"
                  value={formData.title}
                  onChange={(e) => setFormData({ ...formData, title: e.target.value })}
                  className="input-field"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Location *</label>
                <input
                  type="text"
                  value={formData.location}
                  onChange={(e) => setFormData({ ...formData, location: e.target.value })}
                  className="input-field"
                  required
                />
              </div>
            </div>

            <div className="grid grid-cols-3 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Salary Min</label>
                <input
                  type="number"
                  value={formData.salary_min}
                  onChange={(e) => setFormData({ ...formData, salary_min: parseFloat(e.target.value) })}
                  className="input-field"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Salary Max</label>
                <input
                  type="number"
                  value={formData.salary_max}
                  onChange={(e) => setFormData({ ...formData, salary_max: parseFloat(e.target.value) })}
                  className="input-field"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Job Type</label>
                <select
                  value={formData.job_type}
                  onChange={(e) => setFormData({ ...formData, job_type: e.target.value as any })}
                  className="input-field"
                >
                  <option value="full-time">Full Time</option>
                  <option value="part-time">Part Time</option>
                  <option value="contract">Contract</option>
                  <option value="internship">Internship</option>
                  <option value="remote">Remote</option>
                </select>
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Department</label>
              <select
                value={formData.department_id || ''}
                onChange={(e) => setFormData({ ...formData, department_id: e.target.value ? parseInt(e.target.value) : undefined })}
                className="input-field"
              >
                <option value="">No Department</option>
                {departments.map((dept) => (
                  <option key={dept.id} value={dept.id}>
                    {dept.name}
                  </option>
                ))}
              </select>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Description *</label>
              <textarea
                value={formData.description}
                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                rows={4}
                className="input-field"
                required
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Requirements</label>
              <textarea
                value={formData.requirements}
                onChange={(e) => setFormData({ ...formData, requirements: e.target.value })}
                rows={3}
                className="input-field"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Skills (comma separated)</label>
              <input
                type="text"
                value={formData.skills}
                onChange={(e) => setFormData({ ...formData, skills: e.target.value })}
                className="input-field"
                placeholder="e.g. React, TypeScript, Node.js"
              />
            </div>

            <div className="flex gap-3 pt-4">
              <button type="button" onClick={onClose} className="btn-secondary flex-1">
                Cancel
              </button>
              <button type="submit" disabled={submitting} className="btn-primary flex-1">
                {submitting ? 'Saving...' : job ? 'Update Job' : 'Create Job'}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  )
}
