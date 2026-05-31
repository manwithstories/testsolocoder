import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { ShoppingCart, Trash2, Plus, Minus, Package } from 'lucide-react'
import { cartApi, orderApi } from '@/api'
import { useCartStore } from '@/store/cart'

export default function CartPage() {
  const queryClient = useQueryClient()
  const navigate = useNavigate()
  const setItems = useCartStore((state) => state.setItems)
  const totalAmount = useCartStore((state) => state.totalAmount)
  const totalCount = useCartStore((state) => state.totalCount)

  const [showCheckout, setShowCheckout] = useState(false)
  const [shippingInfo, setShippingInfo] = useState({
    shipping_name: '',
    shipping_phone: '',
    shipping_address: '',
  })

  const { data, isLoading } = useQuery({
    queryKey: ['cart'],
    queryFn: () => cartApi.getAll(),
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, quantity }: { id: string; quantity: number }) => cartApi.update(id, { quantity }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['cart'] })
    },
  })

  const removeMutation = useMutation({
    mutationFn: (id: string) => cartApi.remove(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['cart'] })
    },
  })

  const clearMutation = useMutation({
    mutationFn: () => cartApi.clear(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['cart'] })
    },
  })

  const orderMutation = useMutation({
    mutationFn: (data: object) => orderApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['cart'] })
      queryClient.invalidateQueries({ queryKey: ['orders'] })
      setShowCheckout(false)
      navigate('/orders')
    },
  })

  const cartItems = data?.data?.cart_items || []

  const handleQuantityChange = (id: string, newQuantity: number) => {
    if (newQuantity <= 0) {
      removeMutation.mutate(id)
    } else {
      updateMutation.mutate({ id, quantity: newQuantity })
    }
  }

  const handleCheckout = (e: React.FormEvent) => {
    e.preventDefault()
    orderMutation.mutate(shippingInfo)
  }

  if (isLoading) {
    return <div className="text-center py-12 text-gray-500">加载中...</div>
  }

  if (cartItems.length === 0) {
    return (
      <div className="space-y-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">购物车</h1>
          <p className="text-gray-500">管理您的购物车商品</p>
        </div>
        <div className="card text-center py-16">
          <ShoppingCart className="w-16 h-16 text-gray-300 mx-auto mb-4" />
          <p className="text-gray-500 mb-4">购物车是空的</p>
          <Link to="/shop" className="btn-primary">
            去购物
          </Link>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">购物车</h1>
          <p className="text-gray-500">{totalCount} 件商品</p>
        </div>
        <button
          onClick={() => {
            if (confirm('确定要清空购物车吗？')) {
              clearMutation.mutate()
            }
          }}
          className="btn-outline"
        >
          清空购物车
        </button>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Cart Items */}
        <div className="lg:col-span-2 space-y-4">
          {cartItems.map((item: any) => (
            <div key={item.id} className="card">
              <div className="card-body flex items-center gap-4">
                <div className="w-20 h-20 bg-gray-100 rounded-lg flex items-center justify-center flex-shrink-0">
                  {item.product?.image_urls ? (
                    <img
                      src={item.product.image_urls.split(',')[0]}
                      alt={item.product?.name}
                      className="w-full h-full object-cover rounded-lg"
                    />
                  ) : (
                    <Package className="w-8 h-8 text-gray-400" />
                  )}
                </div>
                <div className="flex-1 min-w-0">
                  <h3 className="font-semibold text-gray-900 truncate">
                    {item.product?.name}
                  </h3>
                  <p className="text-garden-600 font-bold">
                    ¥{item.product?.price?.toFixed(2)}
                  </p>
                </div>
                <div className="flex items-center gap-2">
                  <button
                    onClick={() => handleQuantityChange(item.id, item.quantity - 1)}
                    className="p-1 hover:bg-gray-100 rounded-lg"
                  >
                    <Minus className="w-4 h-4" />
                  </button>
                  <span className="w-12 text-center font-medium">{item.quantity}</span>
                  <button
                    onClick={() => handleQuantityChange(item.id, item.quantity + 1)}
                    className="p-1 hover:bg-gray-100 rounded-lg"
                  >
                    <Plus className="w-4 h-4" />
                  </button>
                </div>
                <div className="text-right">
                  <p className="font-bold text-gray-900">
                    ¥{(item.product?.price * item.quantity)?.toFixed(2)}
                  </p>
                  <button
                    onClick={() => removeMutation.mutate(item.id)}
                    className="p-1 hover:bg-red-50 rounded-lg mt-1"
                  >
                    <Trash2 className="w-4 h-4 text-red-400" />
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>

        {/* Order Summary */}
        <div className="card h-fit">
          <div className="card-header">
            <h2 className="text-lg font-semibold text-gray-900">订单摘要</h2>
          </div>
          <div className="card-body space-y-4">
            <div className="flex justify-between">
              <span className="text-gray-600">商品数量</span>
              <span className="font-medium">{totalCount} 件</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-600">商品金额</span>
              <span className="font-medium">¥{totalAmount?.toFixed(2)}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-600">运费</span>
              <span className="font-medium text-garden-600">免运费</span>
            </div>
            <div className="border-t border-gray-200 pt-4 flex justify-between">
              <span className="font-semibold text-gray-900">合计</span>
              <span className="text-xl font-bold text-garden-600">
                ¥{totalAmount?.toFixed(2)}
              </span>
            </div>
            <button
              onClick={() => setShowCheckout(true)}
              className="btn-primary w-full"
              disabled={orderMutation.isPending}
            >
              去结算
            </button>
          </div>
        </div>
      </div>

      {/* Checkout Modal */}
      {showCheckout && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-2xl w-full max-w-md">
            <div className="p-6 border-b border-gray-200">
              <h2 className="text-lg font-semibold">填写收货信息</h2>
            </div>
            <form onSubmit={handleCheckout} className="p-6 space-y-4">
              <div>
                <label className="label">收货人 *</label>
                <input
                  type="text"
                  className="input"
                  required
                  value={shippingInfo.shipping_name}
                  onChange={(e) =>
                    setShippingInfo({ ...shippingInfo, shipping_name: e.target.value })
                  }
                />
              </div>
              <div>
                <label className="label">联系电话 *</label>
                <input
                  type="tel"
                  className="input"
                  required
                  value={shippingInfo.shipping_phone}
                  onChange={(e) =>
                    setShippingInfo({ ...shippingInfo, shipping_phone: e.target.value })
                  }
                />
              </div>
              <div>
                <label className="label">收货地址 *</label>
                <textarea
                  className="input h-24"
                  required
                  value={shippingInfo.shipping_address}
                  onChange={(e) =>
                    setShippingInfo({ ...shippingInfo, shipping_address: e.target.value })
                  }
                />
              </div>
              <div className="flex gap-3 pt-4">
                <button
                  type="button"
                  onClick={() => setShowCheckout(false)}
                  className="btn-outline flex-1"
                >
                  取消
                </button>
                <button
                  type="submit"
                  disabled={orderMutation.isPending}
                  className="btn-primary flex-1"
                >
                  {orderMutation.isPending ? '提交中...' : '提交订单'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}
