import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useState } from 'react'
import { Plus, Heart, MessageSquare, Eye, User, Trash2, Edit2, X } from 'lucide-react'
import { Link } from 'react-router-dom'
import { postApi } from '@/api'

export default function CommunityPage() {
  const queryClient = useQueryClient()
  const [showModal, setShowModal] = useState(false)
  const [formData, setFormData] = useState({
    title: '',
    content: '',
    image_urls: '',
    category: '',
    tags: '',
  })

  const { data, isLoading } = useQuery({
    queryKey: ['posts'],
    queryFn: () => postApi.getAll(),
  })

  const createMutation = useMutation({
    mutationFn: (data: object) => postApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['posts'] })
      setShowModal(false)
      setFormData({
        title: '',
        content: '',
        image_urls: '',
        category: '',
        tags: '',
      })
    },
  })

  const likeMutation = useMutation({
    mutationFn: (id: string) => postApi.like(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['posts'] })
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => postApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['posts'] })
    },
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    createMutation.mutate(formData)
  }

  const handleDelete = (id: string) => {
    if (confirm('确定要删除这篇帖子吗？')) {
      deleteMutation.mutate(id)
    }
  }

  const posts = data?.data?.posts || []

  const categories = [
    { value: '', label: '全部分类' },
    { value: 'experience', label: '种植经验' },
    { value: 'question', label: '求助提问' },
    { value: 'show', label: '成果展示' },
    { value: 'knowledge', label: '知识分享' },
  ]

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">社区交流</h1>
          <p className="text-gray-500">与园艺爱好者分享经验、交流心得</p>
        </div>
        <button onClick={() => setShowModal(true)} className="btn-primary">
          <Plus className="w-4 h-4 mr-2" />
          发布帖子
        </button>
      </div>

      {/* Category Filter */}
      <div className="flex gap-2 overflow-x-auto pb-2">
        {categories.map((cat) => (
          <button
            key={cat.value}
            className="px-4 py-2 rounded-lg text-sm font-medium whitespace-nowrap bg-white border border-gray-200 text-gray-600 hover:bg-gray-50"
          >
            {cat.label}
          </button>
        ))}
      </div>

      {/* Posts */}
      {isLoading ? (
        <div className="text-center py-12 text-gray-500">加载中...</div>
      ) : posts.length === 0 ? (
        <div className="card text-center py-12">
          <MessageSquare className="w-16 h-16 text-gray-300 mx-auto mb-4" />
          <p className="text-gray-500">还没有帖子，快来发布第一篇吧！</p>
        </div>
      ) : (
        <div className="space-y-4">
          {posts.map((post: any) => (
            <div key={post.id} className="card">
              <div className="card-body">
                <div className="flex items-start justify-between mb-3">
                  <div className="flex items-center gap-3">
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
                  {post.category && (
                    <span className="badge bg-garden-100 text-garden-700">
                      {categories.find((c) => c.value === post.category)?.label || post.category}
                    </span>
                  )}
                </div>

                <Link to={`/community/${post.id}`}>
                  <h3 className="text-lg font-semibold text-gray-900 hover:text-garden-600 mb-2">
                    {post.title}
                  </h3>
                  <p className="text-gray-600 line-clamp-3 mb-3">{post.content}</p>
                </Link>

                {post.image_urls && (
                  <div className="flex gap-2 mb-3">
                    {post.image_urls.split(',').slice(0, 3).map((url: string, index: number) => (
                      <img
                        key={index}
                        src={url.trim()}
                        alt=""
                        className="w-24 h-24 object-cover rounded-lg"
                      />
                    ))}
                  </div>
                )}

                <div className="flex items-center gap-6 text-sm text-gray-500">
                  <button
                    onClick={() => likeMutation.mutate(post.id)}
                    className="flex items-center gap-1 hover:text-red-500"
                  >
                    <Heart className="w-4 h-4" />
                    <span>{post.like_count || 0}</span>
                  </button>
                  <Link
                    to={`/community/${post.id}`}
                    className="flex items-center gap-1 hover:text-garden-600"
                  >
                    <MessageSquare className="w-4 h-4" />
                    <span>{post.comments?.length || 0}</span>
                  </Link>
                  <div className="flex items-center gap-1">
                    <Eye className="w-4 h-4" />
                    <span>{post.view_count || 0}</span>
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Modal */}
      {showModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-2xl w-full max-w-lg max-h-[90vh] overflow-y-auto">
            <div className="flex items-center justify-between p-6 border-b border-gray-200">
              <h2 className="text-lg font-semibold">发布帖子</h2>
              <button
                onClick={() => setShowModal(false)}
                className="p-2 hover:bg-gray-100 rounded-lg"
              >
                <X className="w-5 h-5" />
              </button>
            </div>
            <form onSubmit={handleSubmit} className="p-6 space-y-4">
              <div>
                <label className="label">标题 *</label>
                <input
                  type="text"
                  className="input"
                  required
                  value={formData.title}
                  onChange={(e) => setFormData({ ...formData, title: e.target.value })}
                />
              </div>
              <div>
                <label className="label">内容 *</label>
                <textarea
                  className="input h-40"
                  required
                  value={formData.content}
                  onChange={(e) => setFormData({ ...formData, content: e.target.value })}
                />
              </div>
              <div>
                <label className="label">图片链接（多个用逗号分隔）</label>
                <input
                  type="text"
                  className="input"
                  placeholder="https://example.com/image1.jpg, https://example.com/image2.jpg"
                  value={formData.image_urls}
                  onChange={(e) => setFormData({ ...formData, image_urls: e.target.value })}
                />
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="label">分类</label>
                  <select
                    className="input"
                    value={formData.category}
                    onChange={(e) => setFormData({ ...formData, category: e.target.value })}
                  >
                    {categories.map((cat) => (
                      <option key={cat.value} value={cat.value}>
                        {cat.label}
                      </option>
                    ))}
                  </select>
                </div>
                <div>
                  <label className="label">标签</label>
                  <input
                    type="text"
                    className="input"
                    placeholder="标签1,标签2"
                    value={formData.tags}
                    onChange={(e) => setFormData({ ...formData, tags: e.target.value })}
                  />
                </div>
              </div>
              <div className="flex gap-3 pt-4">
                <button
                  type="button"
                  onClick={() => setShowModal(false)}
                  className="btn-outline flex-1"
                >
                  取消
                </button>
                <button
                  type="submit"
                  disabled={createMutation.isPending}
                  className="btn-primary flex-1"
                >
                  {createMutation.isPending ? '发布中...' : '发布'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}
