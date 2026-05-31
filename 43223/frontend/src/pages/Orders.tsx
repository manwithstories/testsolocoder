import { useEffect, useState } from 'react'
import { Link, useLocation } from 'react-router-dom'
import { orderApi, cartApi } from '@/api/order'
import { Order, OrderStatus } from '@/types'

export default function Orders() {
  const location = useLocation()
  const [orders, setOrders] = useState<Order[]>([])
  const [loading, setLoading] = useState(true)
  const [showCheckout, setShowCheckout] = useState(false)
  const [checkoutForm, setCheckoutForm] = useState({
    receiver_name: '',
    receiver_phone: '',
    address: '',
    remark: '',
  })
  const [cartItems, setCartItems] = useState<any[]>([])
  const [cartTotal, setCartTotal] = useState(0)

  useEffect(() => {
    loadOrders()
    if ((location.state as any)?.fromCart) {
      loadCart()
      setShowCheckout(true)
    }
  }, [])

  const loadOrders = async () => {
    try {
      const res = await orderApi.list({ page: 1, page_size: 20 })
      setOrders(res.data?.items || [])
    } catch {
      setOrders([])
    } finally {
      setLoading(false)
    }
  }

  const loadCart = async () => {
    try {
      const res = await cartApi.get()
      setCartItems(res.data?.items || [])
      setCartTotal(res.data?.total_amount || 0)
    } catch {
      setCartItems([])
    }
  }

  const handleCheckout = async () => {
    if (!checkoutForm.receiver_name || !checkoutForm.receiver_phone || !checkoutForm.address) {
      alert('请填写完整的收货信息')
      return
    }
    try {
      const items = cartItems.map((item) => ({
        product_id: item.product_id,
        quantity: item.quantity,
      }))
      await orderApi.create({ ...checkoutForm, items })
      await cartApi.clear()
      alert('订单创建成功')
      setShowCheckout(false)
      loadOrders()
    } catch (err: any) {
      alert(err.message || '下单失败')
    }
  }

  const getStatusBadge = (status: OrderStatus) => {
    const styles: Record<OrderStatus, string> = {
      pending: 'bg-yellow-100 text-yellow-700',
      paid: 'bg-blue-100 text-blue-700',
      processing: 'bg-purple-100 text-purple-700',
      shipped: 'bg-indigo-100 text-indigo-700',
      delivered: 'bg-green-100 text-green-700',
      cancelled: 'bg-gray-100 text-gray-700',
      refunded: 'bg-red-100 text-red-700',
    }
    const labels: Record<OrderStatus, string> = {
      pending: '待支付',
      paid: '已支付',
      processing: '处理中',
      shipped: '已发货',
      delivered: '已送达',
      cancelled: '已取消',
      refunded: '已退款',
    }
    return <span className={`badge ${styles[status]}`}>{labels[status]}</span>
  }

  if (showCheckout) {
    return (
      <div className="max-w-2xl mx-auto space-y-6">
        <h1 className="text-2xl font-bold text-gray-800">确认订单</h1>

        <div className="card p-4 space-y-4">
          <h2 className="font-semibold">商品清单</h2>
          {cartItems.map((item) => (
            <div key={item.id} className="flex items-center justify-between py-2 border-b last:border-0">
              <span>{item.product?.name}</span>
              <span>x{item.quantity}</span>
              <span className="text-coffee-600">¥{((item.product?.price || 0) * item.quantity).toFixed(2)}</span>
            </div>
          ))}
          <div className="flex justify-between pt-2 font-bold">
            <span>合计</span>
            <span className="text-coffee-600">¥{cartTotal.toFixed(2)}</span>
          </div>
        </div>

        <div className="card p-4 space-y-4">
          <h2 className="font-semibold">收货信息</h2>
          <div>
            <label className="label">收货人 *</label>
            <input
              type="text"
              value={checkoutForm.receiver_name}
              onChange={(e) => setCheckoutForm({ ...checkoutForm, receiver_name: e.target.value })}
              className="input"
              required
            />
          </div>
          <div>
            <label className="label">联系电话 *</label>
            <input
              type="tel"
              value={checkoutForm.receiver_phone}
              onChange={(e) => setCheckoutForm({ ...checkoutForm, receiver_phone: e.target.value })}
              className="input"
              required
            />
          </div>
          <div>
            <label className="label">收货地址 *</label>
            <textarea
              value={checkoutForm.address}
              onChange={(e) => setCheckoutForm({ ...checkoutForm, address: e.target.value })}
              className="input"
              rows={3}
              required
            />
          </div>
          <div>
            <label className="label">备注</label>
            <input
              type="text"
              value={checkoutForm.remark}
              onChange={(e) => setCheckoutForm({ ...checkoutForm, remark: e.target.value })}
              className="input"
            />
          </div>
        </div>

        <div className="flex gap-4">
          <button onClick={() => setShowCheckout(false)} className="btn btn-secondary flex-1">
            返回
          </button>
          <button onClick={handleCheckout} className="btn btn-primary flex-1">
            提交订单
          </button>
        </div>
      </div>
    )
  }

  if (loading) {
    return <div className="text-center py-12">加载中...</div>
  }

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-gray-800">我的订单</h1>

      {orders.length === 0 ? (
        <div className="text-center py-12">
          <div className="text-6xl mb-4">📦</div>
          <p className="text-gray-500 mb-4">暂无订单</p>
          <Link to="/products" className="btn btn-primary">去购物</Link>
        </div>
      ) : (
        <div className="space-y-4">
          {orders.map((order) => (
            <div key={order.id} className="card overflow-hidden">
              <div className="p-4 bg-gray-50 flex items-center justify-between">
                <div className="flex items-center gap-4">
                  <span className="text-gray-500">订单号: {order.order_no}</span>
                  <span className="text-gray-500">{order.created_at}</span>
                </div>
                {getStatusBadge(order.status)}
              </div>
              <div className="p-4">
                {order.items?.map((item) => (
                  <div key={item.id} className="flex items-center justify-between py-2">
                    <span>{item.product_name}</span>
                    <span>x{item.quantity}</span>
                    <span className="text-coffee-600">¥{item.subtotal.toFixed(2)}</span>
                  </div>
                ))}
                <div className="flex justify-between pt-4 border-t mt-2">
                  <span className="text-gray-600">合计</span>
                  <span className="text-xl font-bold text-coffee-600">¥{order.total_amount.toFixed(2)}</span>
                </div>
              </div>
              <div className="p-4 bg-gray-50 flex justify-end gap-2">
                <Link to={`/orders/${order.id}`} className="btn btn-outline text-sm">
                  查看详情
                </Link>
                {order.status === 'pending' && (
                  <button
                    onClick={async () => {
                      try {
                        await orderApi.pay(order.id, 'alipay')
                        alert('支付成功')
                        loadOrders()
                      } catch (err: any) {
                        alert(err.message)
                      }
                    }}
                    className="btn btn-primary text-sm"
                  >
                    立即支付
                  </button>
                )}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
