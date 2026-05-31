import { useEffect } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useCartStore } from '@/store/cart'
import { useAuthStore } from '@/store/auth'

export default function Cart() {
  const { items, totalAmount, totalCount, loadCart, updateQuantity, removeFromCart, isLoading } = useCartStore()
  const { isAuthenticated } = useAuthStore()
  const navigate = useNavigate()

  useEffect(() => {
    if (isAuthenticated) {
      loadCart()
    }
  }, [isAuthenticated])

  if (isLoading) {
    return <div className="text-center py-12">加载中...</div>
  }

  if (items.length === 0) {
    return (
      <div className="text-center py-12">
        <div className="text-6xl mb-4">🛒</div>
        <p className="text-gray-500 mb-4">购物车是空的</p>
        <Link to="/products" className="btn btn-primary">去购物</Link>
      </div>
    )
  }

  const handleQuantityChange = async (id: number, qty: number) => {
    if (qty < 1) return
    await updateQuantity(id, qty)
  }

  const handleRemove = async (id: number) => {
    await removeFromCart(id)
  }

  const handleCheckout = () => {
    navigate('/orders', { state: { fromCart: true } })
  }

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-gray-800">购物车 ({totalCount})</h1>

      <div className="space-y-4">
        {items.map((item) => (
          <div key={item.id} className="card p-4 flex items-center gap-4">
            <Link to={`/products/${item.product_id}`} className="w-24 h-24 bg-gray-100 rounded-lg overflow-hidden flex-shrink-0">
              {item.product?.images?.[0] ? (
                <img src={item.product.images[0].url} alt="" className="w-full h-full object-cover" />
              ) : (
                <div className="w-full h-full flex items-center justify-center text-coffee-300 text-3xl">☕</div>
              )}
            </Link>
            <div className="flex-1 min-w-0">
              <Link to={`/products/${item.product_id}`} className="font-semibold text-gray-800 hover:text-coffee-600">
                {item.product?.name || `商品 #${item.product_id}`}
              </Link>
              <p className="text-sm text-gray-500">{item.product?.origin}</p>
              <p className="text-coffee-600 font-bold">¥{item.product?.price}</p>
            </div>
            <div className="flex items-center border rounded-lg">
              <button
                onClick={() => handleQuantityChange(item.id, item.quantity - 1)}
                className="px-3 py-1 text-gray-600 hover:bg-gray-100"
              >
                -
              </button>
              <span className="px-4 py-1 border-x min-w-12 text-center">{item.quantity}</span>
              <button
                onClick={() => handleQuantityChange(item.id, item.quantity + 1)}
                className="px-3 py-1 text-gray-600 hover:bg-gray-100"
              >
                +
              </button>
            </div>
            <div className="text-right w-24">
              <p className="font-bold text-coffee-600">¥{((item.product?.price || 0) * item.quantity).toFixed(2)}</p>
              <button
                onClick={() => handleRemove(item.id)}
                className="text-red-500 text-sm hover:underline"
              >
                删除
              </button>
            </div>
          </div>
        ))}
      </div>

      <div className="card p-6">
        <div className="flex items-center justify-between mb-4">
          <span className="text-gray-600">共 {totalCount} 件商品</span>
          <span className="text-gray-600">合计: <span className="text-2xl font-bold text-coffee-600">¥{totalAmount.toFixed(2)}</span></span>
        </div>
        <button onClick={handleCheckout} className="btn btn-primary w-full py-3">
          去结算
        </button>
      </div>
    </div>
  )
}
