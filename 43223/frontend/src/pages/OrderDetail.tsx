import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { orderApi } from '@/api/order'
import { Order, OrderStatus } from '@/types'

export default function OrderDetail() {
  const { id } = useParams()
  const navigate = useNavigate()
  const [order, setOrder] = useState<Order | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (id) {
      loadOrder(Number(id))
    }
  }, [id])

  const loadOrder = async (orderId: number) => {
    try {
      const res = await orderApi.get(orderId)
      setOrder(res.data || null)
    } catch {
      setOrder(null)
    } finally {
      setLoading(false)
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

  if (loading) {
    return <div className="text-center py-12">加载中...</div>
  }

  if (!order) {
    return <div className="text-center py-12">订单不存在</div>
  }

  return (
    <div className="space-y-6">
      <button onClick={() => navigate(-1)} className="text-coffee-600 hover:underline">
        ← 返回
      </button>

      <div className="card overflow-hidden">
        <div className="p-4 bg-gray-50 flex items-center justify-between">
          <div>
            <h1 className="text-xl font-bold">订单详情</h1>
            <p className="text-sm text-gray-500">订单号: {order.order_no}</p>
          </div>
          {getStatusBadge(order.status)}
        </div>

        <div className="p-4 space-y-4">
          <div className="grid grid-cols-2 gap-4">
            <div>
              <p className="text-sm text-gray-500">收货人</p>
              <p className="font-semibold">{order.receiver_name}</p>
            </div>
            <div>
              <p className="text-sm text-gray-500">联系电话</p>
              <p className="font-semibold">{order.receiver_phone}</p>
            </div>
            <div className="col-span-2">
              <p className="text-sm text-gray-500">收货地址</p>
              <p className="font-semibold">{order.address}</p>
            </div>
            {order.remark && (
              <div className="col-span-2">
                <p className="text-sm text-gray-500">备注</p>
                <p className="font-semibold">{order.remark}</p>
              </div>
            )}
          </div>

          <div className="border-t pt-4">
            <h2 className="font-semibold mb-2">商品清单</h2>
            {order.items?.map((item) => (
              <div key={item.id} className="flex items-center justify-between py-2 border-b last:border-0">
                <span>{item.product_name}</span>
                <span>x{item.quantity}</span>
                <span className="text-coffee-600">¥{item.subtotal.toFixed(2)}</span>
              </div>
            ))}
            <div className="flex justify-between pt-4 font-bold">
              <span>合计</span>
              <span className="text-xl text-coffee-600">¥{order.total_amount.toFixed(2)}</span>
            </div>
          </div>

          <div className="border-t pt-4">
            <h2 className="font-semibold mb-2">订单状态</h2>
            <div className="flex flex-wrap gap-4 text-sm">
              <div>
                <span className="text-gray-500">下单时间: </span>
                {order.created_at}
              </div>
              {order.paid_at && (
                <div>
                  <span className="text-gray-500">支付时间: </span>
                  {order.paid_at}
                </div>
              )}
              {order.shipped_at && (
                <div>
                  <span className="text-gray-500">发货时间: </span>
                  {order.shipped_at}
                </div>
              )}
              {order.delivered_at && (
                <div>
                  <span className="text-gray-500">送达时间: </span>
                  {order.delivered_at}
                </div>
              )}
              {order.cancelled_at && (
                <div>
                  <span className="text-gray-500">取消时间: </span>
                  {order.cancelled_at}
                </div>
              )}
            </div>
          </div>
        </div>

        <div className="p-4 bg-gray-50 flex justify-end gap-2">
          {order.status === 'pending' && (
            <>
              <button
                onClick={async () => {
                  try {
                    await orderApi.cancel(order.id)
                    alert('订单已取消')
                    loadOrder(order.id)
                  } catch (err: any) {
                    alert(err.message)
                  }
                }}
                className="btn btn-secondary"
              >
                取消订单
              </button>
              <button
                onClick={async () => {
                  try {
                    await orderApi.pay(order.id, 'alipay')
                    alert('支付成功')
                    loadOrder(order.id)
                  } catch (err: any) {
                    alert(err.message)
                  }
                }}
                className="btn btn-primary"
              >
                立即支付
              </button>
            </>
          )}
        </div>
      </div>
    </div>
  )
}
