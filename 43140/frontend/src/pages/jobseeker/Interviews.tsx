import { useState, useEffect } from 'react'
import { Calendar, MapPin, User, Clock, Check, X } from 'lucide-react'
import { interviewApi, Interview } from '@/api/jobs'
import dayjs from 'dayjs'

export default function JobSeekerInterviews() {
  const [interviews, setInterviews] = useState<Interview[]>([])
  const [loading, setLoading] = useState(true)
  const [filterStatus, setFilterStatus] = useState('')
  const [showDeclineModal, setShowDeclineModal] = useState<number | null>(null)
  const [declineReason, setDeclineReason] = useState('')

  useEffect(() => {
    loadInterviews()
  }, [])

  const loadInterviews = async () => {
    try {
      const response = await interviewApi.listJobSeeker()
      setInterviews(response.data)
    } catch (err) {
      console.error('Failed to load interviews:', err)
    } finally {
      setLoading(false)
    }
  }

  const handleConfirm = async (id: number) => {
    try {
      await interviewApi.confirm(id)
      loadInterviews()
    } catch (err) {
      alert('Failed to confirm interview')
    }
  }

  const handleDecline = async (id: number) => {
    try {
      await interviewApi.decline(id, { reason: declineReason })
      setShowDeclineModal(null)
      setDeclineReason('')
      loadInterviews()
    } catch (err) {
      alert('Failed to decline interview')
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
      <h1 className="text-2xl font-bold text-gray-900">My Interviews</h1>

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
                  <p className="text-gray-600">
                    {interview.application?.job?.company?.company_name}
                  </p>
                  <div className="grid grid-cols-2 gap-4 mt-3 text-sm">
                    <div className="flex items-center gap-2 text-gray-600">
                      <Calendar className="w-4 h-4" />
                      {dayjs(interview.scheduled_at).format('MMM D, YYYY')}
                    </div>
                    <div className="flex items-center gap-2 text-gray-600">
                      <Clock className="w-4 h-4" />
                      {dayjs(interview.scheduled_at).format('HH:mm')} ({interview.duration} min)
                    </div>
                    <div className="flex items-center gap-2 text-gray-600">
                      <MapPin className="w-4 h-4" />
                      {interview.location}
                    </div>
                    <div className="flex items-center gap-2 text-gray-600">
                      <User className="w-4 h-4" />
                      {interview.interviewer}
                    </div>
                  </div>
                  {interview.notes && (
                    <p className="mt-2 text-sm text-gray-500 italic">{interview.notes}</p>
                  )}
                </div>
                {interview.status === 'scheduled' && (
                  <div className="flex items-center gap-2">
                    <button
                      onClick={() => handleConfirm(interview.id)}
                      className="btn-primary flex items-center gap-2"
                    >
                      <Check className="w-4 h-4" />
                      Confirm
                    </button>
                    <button
                      onClick={() => setShowDeclineModal(interview.id)}
                      className="btn-danger flex items-center gap-2"
                    >
                      <X className="w-4 h-4" />
                      Decline
                    </button>
                  </div>
                )}
              </div>
            </div>
          ))}
        </div>
      )}

      {showDeclineModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-lg shadow-xl max-w-md w-full p-6">
            <h2 className="text-xl font-bold mb-4">Decline Interview</h2>
            <p className="text-gray-600 mb-4">Please provide a reason for declining:</p>
            <textarea
              value={declineReason}
              onChange={(e) => setDeclineReason(e.target.value)}
              rows={3}
              className="input-field mb-4"
              placeholder="Reason for declining..."
            />
            <div className="flex gap-3">
              <button
                onClick={() => {
                  setShowDeclineModal(null)
                  setDeclineReason('')
                }}
                className="btn-secondary flex-1"
              >
                Cancel
              </button>
              <button
                onClick={() => handleDecline(showDeclineModal)}
                className="btn-danger flex-1"
              >
                Confirm Decline
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
