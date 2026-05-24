import { useState, useEffect } from 'react'
import { Star, FileText, User } from 'lucide-react'
import { reviewApi, Review } from '@/api/jobs'
import dayjs from 'dayjs'

export default function JobSeekerReviews() {
  const [reviews, setReviews] = useState<Review[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadReviews()
  }, [])

  const loadReviews = async () => {
    try {
      const response = await reviewApi.listJobSeeker()
      setReviews(response.data)
    } catch (err) {
      console.error('Failed to load reviews:', err)
    } finally {
      setLoading(false)
    }
  }

  const getStatusBadge = (status: string) => {
    const colors: Record<string, string> = {
      offer: 'bg-green-100 text-green-800',
      pass: 'bg-blue-100 text-blue-800',
      reject: 'bg-red-100 text-red-800',
      pending: 'bg-yellow-100 text-yellow-800',
    }
    return colors[status] || 'bg-gray-100 text-gray-800'
  }

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-gray-900">Interview Reviews</h1>

      <p className="text-gray-500">
        Review feedback from your completed interviews will appear here.
      </p>

      {loading ? (
        <div className="text-center py-12">
          <div className="animate-spin w-8 h-8 border-4 border-primary-600 border-t-transparent rounded-full mx-auto"></div>
        </div>
      ) : reviews.length === 0 ? (
        <div className="text-center py-12">
          <Star className="w-16 h-16 text-gray-300 mx-auto" />
          <p className="mt-4 text-gray-500">No reviews yet</p>
          <p className="text-sm text-gray-400">Complete interviews to receive feedback</p>
        </div>
      ) : (
        <div className="space-y-4">
          {reviews.map((review) => (
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
                  <p className="text-gray-600">
                    {review.interview?.application?.job?.company?.company_name}
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
                    <span className="text-sm text-gray-500 ml-2">
                      {review.rating}/5 rating
                    </span>
                  </div>
                  {review.feedback && (
                    <div className="mt-3 bg-gray-50 rounded-lg p-4">
                      <p className="text-sm font-medium text-gray-700 mb-1">Feedback:</p>
                      <p className="text-gray-600">{review.feedback}</p>
                    </div>
                  )}
                  {review.strengths && (
                    <div className="mt-2">
                      <p className="text-sm font-medium text-green-700">Strengths:</p>
                      <p className="text-gray-600 text-sm">{review.strengths}</p>
                    </div>
                  )}
                  {review.weaknesses && (
                    <div className="mt-2">
                      <p className="text-sm font-medium text-orange-700">Areas to improve:</p>
                      <p className="text-gray-600 text-sm">{review.weaknesses}</p>
                    </div>
                  )}
                  <p className="text-xs text-gray-400 mt-3">
                    Reviewed {dayjs(review.created_at).fromNow()}
                  </p>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
