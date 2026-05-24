import { useState, useEffect } from 'react'
import { toast } from 'sonner'
import { reviewApi } from '@/services/api'
import { Review } from '@/types'
import { Star, User as UserIcon, MessageSquare, Flag } from 'lucide-react'

export default function Reviews() {
  const [reviews, setReviews] = useState<Review[]>([])
  const [loading, setLoading] = useState(true)
  const [selectedReview, setSelectedReview] = useState<Review | null>(null)
  const [replyText, setReplyText] = useState('')

  useEffect(() => {
    loadReviews()
  }, [])

  const loadReviews = async () => {
    try {
      setLoading(true)
      const res = await reviewApi.getAll()
      setReviews(res.data)
    } catch (error) {
      console.error('Failed to load reviews:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleReply = async () => {
    if (!selectedReview || !replyText.trim()) return

    try {
      await reviewApi.reply(selectedReview.id, { reply: replyText })
      toast.success('回复成功')
      setSelectedReview(null)
      setReplyText('')
      loadReviews()
    } catch (error: any) {
      toast.error(error.response?.data?.error || '回复失败')
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">评价管理</h1>
        <p className="text-gray-500">查看和管理您的评价</p>
      </div>

      <div className="space-y-4">
        {reviews.map((review) => (
          <div key={review.id} className="card">
            <div className="flex items-start justify-between mb-4">
              <div className="flex items-center gap-3">
                {review.reviewer?.avatarUrl ? (
                  <img
                    src={review.reviewer.avatarUrl}
                    alt=""
                    className="h-10 w-10 rounded-full object-cover"
                  />
                ) : (
                  <div className="h-10 w-10 rounded-full bg-primary-600 flex items-center justify-center text-white">
                    <UserIcon className="h-5 w-5" />
                  </div>
                )}
                <div>
                  <div className="font-medium text-gray-900">
                    {review.isAnonymous ? '匿名用户' : `${review.reviewer?.firstName} ${review.reviewer?.lastName}`}
                  </div>
                  <div className="flex items-center gap-1 text-sm text-gray-500">
                    {[...Array(5)].map((_, i) => (
                      <Star
                        key={i}
                        className={`h-4 w-4 ${i < review.rating ? 'text-yellow-500 fill-current' : 'text-gray-300'}`}
                      />
                    ))}
                    <span className="ml-2">
                      {new Date(review.createdAt).toLocaleDateString('zh-CN')}
                    </span>
                  </div>
                </div>
              </div>
              <button className="p-2 text-gray-400 hover:text-gray-600">
                <Flag className="h-4 w-4" />
              </button>
            </div>

            <p className="text-gray-600 mb-4">{review.content}</p>

            {review.tags && (
              <div className="flex flex-wrap gap-1 mb-4">
                {review.tags.split(',').map((tag, index) => (
                  <span key={index} className="badge bg-gray-100 text-gray-600">
                    {tag.trim()}
                  </span>
                ))}
              </div>
            )}

            {review.teacherReply && (
              <div className="p-4 bg-gray-50 rounded-lg">
                <div className="flex items-center gap-2 mb-2">
                  <MessageSquare className="h-4 w-4 text-primary-600" />
                  <span className="text-sm font-medium text-gray-900">老师回复</span>
                  {review.teacherRepliedAt && (
                    <span className="text-xs text-gray-400">
                      {new Date(review.teacherRepliedAt).toLocaleDateString('zh-CN')}
                    </span>
                  )}
                </div>
                <p className="text-sm text-gray-600">{review.teacherReply}</p>
              </div>
            )}

            {!review.teacherReply && (
              <button
                onClick={() => setSelectedReview(review)}
                className="text-primary-600 hover:text-primary-700 text-sm font-medium"
              >
                回复评价
              </button>
            )}
          </div>
        ))}

        {reviews.length === 0 && (
          <div className="text-center py-12">
            <Star className="h-12 w-12 text-gray-300 mx-auto mb-4" />
            <p className="text-gray-500">暂无评价</p>
          </div>
        )}
      </div>

      {selectedReview && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-2xl p-6 w-full max-w-md mx-4">
            <h3 className="text-xl font-bold text-gray-900 mb-4">回复评价</h3>
            <textarea
              value={replyText}
              onChange={(e) => setReplyText(e.target.value)}
              className="input-field h-32 resize-none"
              placeholder="输入您的回复..."
            />
            <div className="flex gap-3 mt-4">
              <button
                onClick={() => setSelectedReview(null)}
                className="btn-secondary flex-1"
              >
                取消
              </button>
              <button
                onClick={handleReply}
                className="btn-primary flex-1"
              >
                发送回复
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
