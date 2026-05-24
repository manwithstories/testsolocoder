import { useState, useEffect } from 'react'
import { Star, Edit2 } from 'lucide-react'
import { reviewApi, interviewApi, Review, Interview } from '@/api/jobs'
import dayjs from 'dayjs'

export default function CompanyReviews() {
  const [reviews, setReviews] = useState<Review[]>([])
  const [interviews, setInterviews] = useState<Interview[]>([])
  const [loading, setLoading] = useState(true)
  const [showModal, setShowModal] = useState(false)
  const [editingReview, setEditingReview] = useState<Review | null>(null)
  const [filterStatus, setFilterStatus] = useState('')

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    try {
      const [reviewsRes, interviewsRes] = await Promise.all([
        reviewApi.listCompany(),
        interviewApi.listCompany(),
      ])
      setReviews(reviewsRes.data)
      setInterviews(interviewsRes.data)
    } catch (err) {
      console.error('Failed to load data:', err)
    } finally {
      setLoading(false)
    }
  }

  const filteredReviews = reviews.filter(
    (r) => !filterStatus || r.status === filterStatus
  )

  const getStatusBadge = (status: string) => {
    const colors: Record<string, string> = {
      offer: 'bg-green-100 text-green-800',
      pass: 'bg-blue-100 text-blue-800',
      reject: 'bg-red-100 text-red-800',
      pending: 'bg-yellow-100 text-yellow-800',
    }
    return colors[status] || 'bg-gray-100 text-gray-800'
  }

  const completedInterviewsWithoutReview = interviews.filter(
    (i) => i.status === 'completed' && !i.review
  )

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-gray-900">Interview Reviews</h1>
        <button
          onClick={() => setShowModal(true)}
          className="btn-primary flex items-center gap-2"
        >
          <Star className="w-5 h-5" />
          Add Review
        </button>
      </div>

      <select
        value={filterStatus}
        onChange={(e) => setFilterStatus(e.target.value)}
        className="input-field w-auto"
      >
        <option value="">All Statuses</option>
        <option value="offer">Offer</option>
        <option value="pass">Pass</option>
        <option value="reject">Reject</option>
        <option value="pending">Pending</option>
      </select>

      {loading ? (
        <div className="text-center py-12">
          <div className="animate-spin w-8 h-8 border-4 border-primary-600 border-t-transparent rounded-full mx-auto"></div>
        </div>
      ) : filteredReviews.length === 0 ? (
        <div className="text-center py-12">
          <Star className="w-16 h-16 text-gray-300 mx-auto" />
          <p className="mt-4 text-gray-500">No reviews yet</p>
        </div>
      ) : (
        <div className="space-y-4">
          {filteredReviews.map((review) => (
            <div key={review.id} className="card">
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <div className="flex items-center gap-3">
                    <h3 className="text-lg font-semibold">
                      {review.interview?.application?.job?.title}
                    </h3>
                    <span className={`badge ${getStatusBadge(review.status)}`}>
                      {review.status}
                    </span>
                  </div>
                  <p className="text-sm text-gray-500 mt-1">
                    Candidate: {review.interview?.application?.jobseeker?.user?.name}
                  </p>
                  <div className="flex items-center gap-1 mt-2">
                    {[1, 2, 3, 4, 5].map((star) => (
                      <Star
                        key={star}
                        className={`w-5 h-5 ${
                          star <= review.rating
                            ? 'text-yellow-400 fill-yellow-400'
                            : 'text-gray-300'
                        }`}
                      />
                    ))}
                  </div>
                  {review.feedback && (
                    <p className="mt-2 text-gray-600 text-sm">{review.feedback}</p>
                  )}
                </div>
                <button
                  onClick={() => {
                    setEditingReview(review)
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
        <ReviewModal
          review={editingReview}
          interviews={completedInterviewsWithoutReview}
          onClose={() => {
            setShowModal(false)
            setEditingReview(null)
          }}
          onSave={() => {
            setShowModal(false)
            setEditingReview(null)
            loadData()
          }}
        />
      )}
    </div>
  )
}

function ReviewModal({
  review,
  interviews,
  onClose,
  onSave,
}: {
  review: Review | null
  interviews: Interview[]
  onClose: () => void
  onSave: () => void
}) {
  const [formData, setFormData] = useState({
    interview_id: review?.interview_id || 0,
    rating: review?.rating || 5,
    feedback: review?.feedback || '',
    strengths: review?.strengths || '',
    weaknesses: review?.weaknesses || '',
    status: review?.status || 'pending',
  })
  const [submitting, setSubmitting] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!formData.interview_id) {
      alert('Please select an interview')
      return
    }
    setSubmitting(true)
    try {
      if (review) {
        await reviewApi.update(review.id, formData)
      } else {
        await reviewApi.create(formData)
      }
      onSave()
    } catch (err: any) {
      alert(err.message || 'Failed to save review')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-lg shadow-xl max-w-lg w-full max-h-[90vh] overflow-y-auto">
        <div className="p-6">
          <h2 className="text-2xl font-bold mb-6">
            {review ? 'Edit Review' : 'Add Review'}
          </h2>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Interview *
              </label>
              <select
                value={formData.interview_id}
                onChange={(e) =>
                  setFormData({ ...formData, interview_id: parseInt(e.target.value) })
                }
                className="input-field"
                required
                disabled={!!review}
              >
                <option value="">Select interview...</option>
                {interviews.map((interview) => (
                  <option key={interview.id} value={interview.id}>
                    {interview.application?.job?.title} -{' '}
                    {interview.application?.jobseeker?.user?.name} (
                    {dayjs(interview.scheduled_at).format('MMM D')})
                  </option>
                ))}
              </select>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Rating *</label>
              <div className="flex items-center gap-2">
                {[1, 2, 3, 4, 5].map((star) => (
                  <button
                    key={star}
                    type="button"
                    onClick={() => setFormData({ ...formData, rating: star })}
                    className="focus:outline-none"
                  >
                    <Star
                      className={`w-8 h-8 ${
                        star <= formData.rating
                          ? 'text-yellow-400 fill-yellow-400'
                          : 'text-gray-300'
                      }`}
                    />
                  </button>
                ))}
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Feedback *</label>
              <textarea
                value={formData.feedback}
                onChange={(e) => setFormData({ ...formData, feedback: e.target.value })}
                rows={4}
                className="input-field"
                required
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Strengths</label>
              <textarea
                value={formData.strengths}
                onChange={(e) => setFormData({ ...formData, strengths: e.target.value })}
                rows={2}
                className="input-field"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Weaknesses</label>
              <textarea
                value={formData.weaknesses}
                onChange={(e) => setFormData({ ...formData, weaknesses: e.target.value })}
                rows={2}
                className="input-field"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Decision *</label>
              <select
                value={formData.status}
                onChange={(e) => setFormData({ ...formData, status: e.target.value as any })}
                className="input-field"
                required
              >
                <option value="pending">Pending</option>
                <option value="offer">Offer</option>
                <option value="pass">Pass</option>
                <option value="reject">Reject</option>
              </select>
            </div>

            <div className="flex gap-3 pt-4">
              <button type="button" onClick={onClose} className="btn-secondary flex-1">
                Cancel
              </button>
              <button type="submit" disabled={submitting} className="btn-primary flex-1">
                {submitting ? 'Saving...' : review ? 'Update' : 'Submit'}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  )
}
