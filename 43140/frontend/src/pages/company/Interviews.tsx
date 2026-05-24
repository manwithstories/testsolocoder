import { useState, useEffect } from 'react'
import { Plus, Edit2, Calendar, MapPin, User, Phone, Clock } from 'lucide-react'
import { interviewApi, applicationApi, Interview, Application } from '@/api/jobs'
import dayjs from 'dayjs'

export default function CompanyInterviews() {
  const [interviews, setInterviews] = useState<Interview[]>([])
  const [applications, setApplications] = useState<Application[]>([])
  const [loading, setLoading] = useState(true)
  const [showModal, setShowModal] = useState(false)
  const [editingInterview, setEditingInterview] = useState<Interview | null>(null)
  const [filterStatus, setFilterStatus] = useState('')

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    try {
      const [interviewsRes, appsRes] = await Promise.all([
        interviewApi.listCompany(),
        applicationApi.listCompany(),
      ])
      setInterviews(interviewsRes.data)
      setApplications(appsRes.data)
    } catch (err) {
      console.error('Failed to load data:', err)
    } finally {
      setLoading(false)
    }
  }

  const filteredInterviews = interviews.filter(
    (i) => !filterStatus || i.status === filterStatus
  )

  const getStatusBadge = (status: string) => {
    const colors: Record<string, string> = {
      scheduled: 'bg-blue-100 text-blue-800',
      confirmed: 'bg-green-100 text-green-800',
      declined: 'bg-red-100 text-red-800',
      completed: 'bg-gray-100 text-gray-800',
      cancelled: 'bg-yellow-100 text-yellow-800',
    }
    return colors[status] || 'bg-gray-100 text-gray-800'
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-gray-900">Interviews</h1>
        <button onClick={() => setShowModal(true)} className="btn-primary flex items-center gap-2">
          <Plus className="w-5 h-5" />
          Schedule Interview
        </button>
      </div>

      <select
        value={filterStatus}
        onChange={(e) => setFilterStatus(e.target.value)}
        className="input-field w-auto"
      >
        <option value="">All Statuses</option>
        <option value="scheduled">Scheduled</option>
        <option value="confirmed">Confirmed</option>
        <option value="declined">Declined</option>
        <option value="completed">Completed</option>
        <option value="cancelled">Cancelled</option>
      </select>

      {loading ? (
        <div className="text-center py-12">
          <div className="animate-spin w-8 h-8 border-4 border-primary-600 border-t-transparent rounded-full mx-auto"></div>
        </div>
      ) : filteredInterviews.length === 0 ? (
        <div className="text-center py-12">
          <Calendar className="w-16 h-16 text-gray-300 mx-auto" />
          <p className="mt-4 text-gray-500">No interviews scheduled</p>
        </div>
      ) : (
        <div className="space-y-4">
          {filteredInterviews.map((interview) => (
            <div key={interview.id} className="card">
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <div className="flex items-center gap-3">
                    <h3 className="text-lg font-semibold">
                      {interview.application?.job?.title}
                    </h3>
                    <span className={`badge ${getStatusBadge(interview.status)}`}>
                      {interview.status}
                    </span>
                  </div>
                  <div className="mt-3 space-y-2 text-sm text-gray-600">
                    <div className="flex items-center gap-2">
                      <User className="w-4 h-4" />
                      {interview.application?.jobseeker?.user?.name}
                    </div>
                    <div className="flex items-center gap-2">
                      <Calendar className="w-4 h-4" />
                      {dayjs(interview.scheduled_at).format('MMM D, YYYY HH:mm')}
                    </div>
                    <div className="flex items-center gap-2">
                      <Clock className="w-4 h-4" />
                      {interview.duration} minutes
                    </div>
                    <div className="flex items-center gap-2">
                      <MapPin className="w-4 h-4" />
                      {interview.location}
                    </div>
                    <div className="flex items-center gap-2">
                      <Phone className="w-4 h-4" />
                      {interview.interviewer}
                    </div>
                  </div>
                </div>
                <button
                  onClick={() => {
                    setEditingInterview(interview)
                    setShowModal(true)
                  }}
                  className="p-2 text-gray-600 hover:text-primary-600"
                >
                  <Edit2 className="w-5 h-5" />
                </button>
              </div>
            </div>
          ))}
        </div>
      )}

      {showModal && (
        <InterviewModal
          interview={editingInterview}
          applications={applications.filter((a) => a.status === 'pending' || a.status === 'reviewed' || a.status === 'interview')}
          onClose={() => {
            setShowModal(false)
            setEditingInterview(null)
          }}
          onSave={() => {
            setShowModal(false)
            setEditingInterview(null)
            loadData()
          }}
        />
      )}
    </div>
  )
}

function InterviewModal({
  interview,
  applications,
  onClose,
  onSave,
}: {
  interview: Interview | null
  applications: Application[]
  onClose: () => void
  onSave: () => void
}) {
  const [formData, setFormData] = useState({
    application_id: interview?.application_id || 0,
    scheduled_at: interview?.scheduled_at
      ? dayjs(interview.scheduled_at).format('YYYY-MM-DDTHH:mm')
      : dayjs().add(1, 'day').format('YYYY-MM-DDTHH:mm'),
    duration: interview?.duration || 60,
    location: interview?.location || '',
    interviewer: interview?.interviewer || '',
    interviewer_email: interview?.interviewer_email || '',
    notes: interview?.notes || '',
  })
  const [submitting, setSubmitting] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!formData.application_id) {
      alert('Please select an application')
      return
    }
    setSubmitting(true)
    try {
      const data = {
        ...formData,
        scheduled_at: new Date(formData.scheduled_at).toISOString(),
      }
      if (interview) {
        await interviewApi.update(interview.id, data)
      } else {
        await interviewApi.create(data)
      }
      onSave()
    } catch (err: any) {
      alert(err.message || 'Failed to save interview')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-lg shadow-xl max-w-lg w-full max-h-[90vh] overflow-y-auto">
        <div className="p-6">
          <h2 className="text-2xl font-bold mb-6">
            {interview ? 'Edit Interview' : 'Schedule Interview'}
          </h2>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Application *
              </label>
              <select
                value={formData.application_id}
                onChange={(e) =>
                  setFormData({ ...formData, application_id: parseInt(e.target.value) })
                }
                className="input-field"
                required
                disabled={!!interview}
              >
                <option value="">Select application...</option>
                {applications.map((app) => (
                  <option key={app.id} value={app.id}>
                    {app.job?.title} - {app.jobseeker?.user?.name}
                  </option>
                ))}
              </select>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Date & Time *
              </label>
              <input
                type="datetime-local"
                value={formData.scheduled_at}
                onChange={(e) => setFormData({ ...formData, scheduled_at: e.target.value })}
                className="input-field"
                required
              />
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Duration (minutes)</label>
                <input
                  type="number"
                  value={formData.duration}
                  onChange={(e) => setFormData({ ...formData, duration: parseInt(e.target.value) })}
                  className="input-field"
                  min={15}
                  max={240}
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Location *</label>
                <input
                  type="text"
                  value={formData.location}
                  onChange={(e) => setFormData({ ...formData, location: e.target.value })}
                  className="input-field"
                  placeholder="e.g. Zoom Meeting, Office Address"
                  required
                />
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Interviewer *</label>
              <input
                type="text"
                value={formData.interviewer}
                onChange={(e) => setFormData({ ...formData, interviewer: e.target.value })}
                className="input-field"
                required
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Interviewer Email</label>
              <input
                type="email"
                value={formData.interviewer_email}
                onChange={(e) => setFormData({ ...formData, interviewer_email: e.target.value })}
                className="input-field"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Notes</label>
              <textarea
                value={formData.notes}
                onChange={(e) => setFormData({ ...formData, notes: e.target.value })}
                rows={3}
                className="input-field"
              />
            </div>

            <div className="flex gap-3 pt-4">
              <button type="button" onClick={onClose} className="btn-secondary flex-1">
                Cancel
              </button>
              <button type="submit" disabled={submitting} className="btn-primary flex-1">
                {submitting ? 'Saving...' : interview ? 'Update' : 'Schedule'}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  )
}
