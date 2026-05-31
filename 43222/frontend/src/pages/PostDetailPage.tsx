import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { Heart, MessageSquare, Eye, User, ArrowLeft, Send, Trash2 } from 'lucide-react'
import { postApi } from '@/api'

export default function PostDetailPage() {
  const { id } = useParams()
  const queryClient = useQueryClient()
  const [comment, setComment] = useState('')

  const { data, isLoading } = useQuery({
    queryKey: ['post', id],
    queryFn: () => postApi.getById(id!),
  })

  const likeMutation = useMutation({
    mutationFn: (postId: string) => postApi.like(postId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['post', id] })
    },
  })

  const commentMutation = useMutation({
    mutationFn: (content: string) => postApi.createComment(id!, { content }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['post', id] })
      setComment('')
    },
  })

  const deleteCommentMutation = useMutation({
    mutationFn: (commentId: string) => postApi.deleteComment(commentId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['post', id] })
    },
  })

  const handleSubmitComment = (e: React.FormEvent) => {
    e.preventDefault()
    if (comment.trim()) {
      commentMutation.mutate(comment)
    }
  }

  const post = data?.data?.post

  if (isLoading) {
    return <div className="text-center py-12 text-gray-500">加载中...</div>
  }

  if (!post) {
    return (
      <div className="card text-center py-12">
        <p className="text-gray-500">帖子不存在</p>
        <Link to="/community" className="text-garden-600 hover:text-garden-700">
          返回社区
        </Link>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <Link
        to="/community"
        className="flex items-center gap-2 text-gray-600 hover:text-garden-600"
      >
        <ArrowLeft className="w-4 h-4" />
        返回社区
      </Link>

      <div className="card">
        <div className="card-body">
          <div className="flex items-center gap-3 mb-4">
            <div className="w-10 h-10 rounded-full bg-garden-100 flex items-center justify-center">
              <User className="w-5 h-5 text-garden-600" />
            </div>
            <div>
              <p className="font-medium text-gray-900">
                {post.user?.nickname || post.user?.username}
              </p>
              <p className="text-xs text-gray-500">
                {new Date(post.created_at).toLocaleString('zh-CN')}
              </p>
            </div>
          </div>

          <h1 className="text-2xl font-bold text-gray-900 mb-4">{post.title}</h1>
          <div className="prose max-w-none text-gray-700 whitespace-pre-wrap">
            {post.content}
          </div>

          {post.image_urls && (
            <div className="mt-4 grid grid-cols-2 md:grid-cols-3 gap-4">
              {post.image_urls.split(',').map((url: string, index: number) => (
                <img
                  key={index}
                  src={url.trim()}
                  alt=""
                  className="w-full aspect-square object-cover rounded-lg"
                />
              ))}
            </div>
          )}

          <div className="flex items-center gap-6 mt-6 pt-4 border-t border-gray-100 text-sm text-gray-500">
            <button
              onClick={() => likeMutation.mutate(post.id)}
              className="flex items-center gap-1 hover:text-red-500"
            >
              <Heart className="w-4 h-4" />
              <span>{post.like_count || 0}</span>
            </button>
            <div className="flex items-center gap-1">
              <MessageSquare className="w-4 h-4" />
              <span>{post.comments?.length || 0}</span>
            </div>
            <div className="flex items-center gap-1">
              <Eye className="w-4 h-4" />
              <span>{post.view_count || 0}</span>
            </div>
          </div>
        </div>
      </div>

      {/* Comments */}
      <div className="card">
        <div className="card-header">
          <h2 className="text-lg font-semibold text-gray-900">
            评论 ({post.comments?.length || 0})
          </h2>
        </div>
        <div className="card-body">
          <form onSubmit={handleSubmitComment} className="mb-6">
            <div className="flex gap-3">
              <input
                type="text"
                className="input flex-1"
                placeholder="写下你的评论..."
                value={comment}
                onChange={(e) => setComment(e.target.value)}
              />
              <button
                type="submit"
                disabled={commentMutation.isPending || !comment.trim()}
                className="btn-primary"
              >
                <Send className="w-4 h-4" />
              </button>
            </div>
          </form>

          {post.comments?.length === 0 ? (
            <p className="text-center text-gray-500 py-8">还没有评论，快来抢沙发吧！</p>
          ) : (
            <div className="space-y-4">
              {post.comments.map((c: any) => (
                <div key={c.id} className="flex gap-3">
                  <div className="w-8 h-8 rounded-full bg-garden-100 flex items-center justify-center flex-shrink-0">
                    <User className="w-4 h-4 text-garden-600" />
                  </div>
                  <div className="flex-1">
                    <div className="flex items-center justify-between">
                      <span className="text-sm font-medium text-gray-900">
                        {c.user?.nickname || c.user?.username}
                      </span>
                      <span className="text-xs text-gray-400">
                        {new Date(c.created_at).toLocaleString('zh-CN')}
                      </span>
                    </div>
                    <p className="text-gray-700 mt-1">{c.content}</p>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
