import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useState } from 'react'
import { Plus, RefreshCw, MapPin, X, MessageSquare } from 'lucide-react'
import { exchangeApi } from '@/api'

export default function ExchangePage() {
  const queryClient = useQueryClient()
  const [showModal, setShowModal] = useState(false)
  const [selectedExchange, setSelectedExchange] = useState<any>(null)
  const [showOfferModal, setShowOfferModal] = useState(false)
  const [formData, setFormData] = useState({
    title: '',
    seed_name: '',
    description: '',
    image_urls: '',
    quantity: 1,
    exchange_type: 'exchange',
    want_seeds: '',
    location: '',
  })
  const [offerData, setOfferData] = useState({
    offer_seeds: '',
    message: '',
  })

  const { data, isLoading } = useQuery({
    queryKey: ['seed-exchanges'],
    queryFn: () => exchangeApi.getAll(),
  })

  const { data: offersData } = useQuery({
    queryKey: ['exchange-offers', selectedExchange?.id],
    queryFn: () => exchangeApi.getOffers(selectedExchange?.id),
    enabled: !!selectedExchange,
  })

  const createMutation = useMutation({
    mutationFn: (data: object) => exchangeApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['seed-exchanges'] })
      setShowModal(false)
      resetForm()
    },
  })

  const offerMutation = useMutation({
    mutationFn: (data: object) => exchangeApi.createOffer(selectedExchange?.id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['exchange-offers', selectedExchange?.id] })
      setShowOfferModal(false)
      setOfferData({ offer_seeds: '', message: '' })
    },
  })

  const updateOfferMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: object }) => exchangeApi.updateOffer(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['exchange-offers', selectedExchange?.id] })
      queryClient.invalidateQueries({ queryKey: ['seed-exchanges'] })
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => exchangeApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['seed-exchanges'] })
    },
  })

  const resetForm = () => {
    setFormData({
      title: '',
      seed_name: '',
      description: '',
      image_urls: '',
      quantity: 1,
      exchange_type: 'exchange',
      want_seeds: '',
      location: '',
    })
  }

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    createMutation.mutate(formData)
  }

  const handleOfferSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    offerMutation.mutate(offerData)
  }

  const handleAcceptOffer = (id: string) => {
    updateOfferMutation.mutate({ id, data: { status: 'accepted' } })
  }

  const handleDelete = (id: string) => {
    if (confirm('确定要删除这条交换信息吗？')) {
      deleteMutation.mutate(id)
    }
  }

  const exchanges = data?.data?.exchanges || []
  const offers = offersData?.data?.offers || []

  const statusColors: Record<string, string> = {
    available: 'bg-green-100 text-green-700',
    exchanged: 'bg-gray-100 text-gray-700',
    pending: 'bg-amber-100 text-amber-700',
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">种子交换</h1>
          <p className="text-gray-500">与其他园艺爱好者交换多余的种子</p>
        </div>
        <button onClick={() => setShowModal(true)} className="btn-primary">
          <Plus className="w-4 h-4 mr-2" />
          发布交换
        </button>
      </div>

      {/* Exchange List */}
      {isLoading ? (
        <div className="text-center py-12 text-gray-500">加载中...</div>
      ) : exchanges.length === 0 ? (
        <div className="card text-center py-12">
          <RefreshCw className="w-16 h-16 text-gray-300 mx-auto mb-4" />
          <p className="text-gray-500">还没有种子交换信息</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {exchanges.map((exchange: any) => (
            <div key={exchange.id} className="card hover:shadow-md transition-shadow">
              <div className="card-body">
                <div className="flex items-start justify-between mb-3">
                  <h3 className="font-semibold text-gray-900">{exchange.title}</h3>
                  <span className={`badge ${statusColors[exchange.status] || 'bg-gray-100'}`}>
                    {exchange.status === 'available' ? '可交换' : exchange.status}
                  </span>
                </div>

                {exchange.image_urls && (
                  <div className="flex gap-2 mb-3">
                    {exchange.image_urls.split(',').slice(0, 2).map((url: string, i: number) => (
                      <img
                        key={i}
                        src={url.trim()}
                        alt=""
                        className="w-20 h-20 object-cover rounded-lg"
                      />
                    ))}
                  </div>
                )}

                <div className="space-y-2 text-sm">
                  <p className="text-gray-700">
                    <span className="text-gray-500">提供：</span>
                    {exchange.seed_name} ({exchange.quantity}份)
                  </p>
                  {exchange.want_seeds && (
                    <p className="text-gray-700">
                      <span className="text-gray-500">想要：</span>
                      {exchange.want_seeds}
                    </p>
                  )}
                  {exchange.location && (
                    <div className="flex items-center gap-1 text-gray-500">
                      <MapPin className="w-4 h-4" />
                      <span>{exchange.location}</span>
                    </div>
                  )}
                </div>

                <div className="mt-4 flex gap-2">
                  <button
                    onClick={() => {
                      setSelectedExchange(exchange)
                      setShowOfferModal(true)
                    }}
                    className="btn-outline flex-1 text-sm"
                    disabled={exchange.status !== 'available'}
                  >
                    <MessageSquare className="w-4 h-4 mr-1" />
                    发起交换
                  </button>
                </div>

                <div className="mt-3 pt-3 border-t border-gray-100 flex items-center justify-between">
                  <span className="text-xs text-gray-500">
                    {exchange.owner?.nickname || exchange.owner?.username}
                  </span>
                  <span className="text-xs text-gray-400">
                    {new Date(exchange.created_at).toLocaleDateString('zh-CN')}
                  </span>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Create Modal */}
      {showModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-2xl w-full max-w-lg max-h-[90vh] overflow-y-auto">
            <div className="flex items-center justify-between p-6 border-b border-gray-200">
              <h2 className="text-lg font-semibold">发布种子交换</h2>
              <button
                onClick={() => {
                  setShowModal(false)
                  resetForm()
                }}
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
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="label">种子名称 *</label>
                  <input
                    type="text"
                    className="input"
                    required
                    value={formData.seed_name}
                    onChange={(e) => setFormData({ ...formData, seed_name: e.target.value })}
                  />
                </div>
                <div>
                  <label className="label">数量</label>
                  <input
                    type="number"
                    className="input"
                    min="1"
                    value={formData.quantity}
                    onChange={(e) => setFormData({ ...formData, quantity: parseInt(e.target.value) })}
                  />
                </div>
              </div>
              <div>
                <label className="label">描述</label>
                <textarea
                  className="input h-24"
                  value={formData.description}
                  onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                />
              </div>
              <div>
                <label className="label">想要交换的种子</label>
                <input
                  type="text"
                  className="input"
                  placeholder="请描述您想要的种子"
                  value={formData.want_seeds}
                  onChange={(e) => setFormData({ ...formData, want_seeds: e.target.value })}
                />
              </div>
              <div>
                <label className="label">位置</label>
                <input
                  type="text"
                  className="input"
                  placeholder="您所在的城市或地区"
                  value={formData.location}
                  onChange={(e) => setFormData({ ...formData, location: e.target.value })}
                />
              </div>
              <div>
                <label className="label">图片链接（多个用逗号分隔）</label>
                <input
                  type="text"
                  className="input"
                  value={formData.image_urls}
                  onChange={(e) => setFormData({ ...formData, image_urls: e.target.value })}
                />
              </div>
              <div className="flex gap-3 pt-4">
                <button
                  type="button"
                  onClick={() => {
                    setShowModal(false)
                    resetForm()
                  }}
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

      {/* Offer Modal */}
      {showOfferModal && selectedExchange && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-2xl w-full max-w-lg">
            <div className="flex items-center justify-between p-6 border-b border-gray-200">
              <h2 className="text-lg font-semibold">发起交换</h2>
              <button
                onClick={() => setShowOfferModal(false)}
                className="p-2 hover:bg-gray-100 rounded-lg"
              >
                <X className="w-5 h-5" />
              </button>
            </div>
            <form onSubmit={handleOfferSubmit} className="p-6 space-y-4">
              <div className="p-3 bg-garden-50 rounded-lg">
                <p className="text-sm text-garden-800">
                  <span className="font-medium">对方提供：</span>
                  {selectedExchange.seed_name}
                </p>
              </div>
              <div>
                <label className="label">您提供的种子 *</label>
                <textarea
                  className="input h-24"
                  required
                  placeholder="请描述您能提供的种子"
                  value={offerData.offer_seeds}
                  onChange={(e) => setOfferData({ ...offerData, offer_seeds: e.target.value })}
                />
              </div>
              <div>
                <label className="label">留言</label>
                <textarea
                  className="input h-20"
                  placeholder="给对方的留言"
                  value={offerData.message}
                  onChange={(e) => setOfferData({ ...offerData, message: e.target.value })}
                />
              </div>
              <div className="flex gap-3 pt-4">
                <button
                  type="button"
                  onClick={() => setShowOfferModal(false)}
                  className="btn-outline flex-1"
                >
                  取消
                </button>
                <button
                  type="submit"
                  disabled={offerMutation.isPending}
                  className="btn-primary flex-1"
                >
                  {offerMutation.isPending ? '发送中...' : '发送交换请求'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}
