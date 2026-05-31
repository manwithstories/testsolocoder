import { useQuery } from '@tanstack/react-query'
import { Package, Truck, CheckCircle, XCircle } from 'lucide-react'
import { orderApi } from '@/api'
import { useState } from 'react'

export default function OrdersPage() {
  const [filterStatus, setFilterStatus] = useState('')

  const { data, isLoading } = useQuery({
    queryKey: ['orders', filterStatus],
    queryFn: () =>
      orderApi.getAll({
        status: filterStatus || undefined,
      }),
  })

  const orders = data?.data?.orders || []

  const statusConfig: Record<string, { label: string; color: string; icon: any }> = {
    pending: { label: '待处理', color: 'bg-amber-100 text-amber-700', icon: Package },
    paid: { label: '已付款', color: 'bg-blue-100 text-blue-700', icon: CheckCircle },
    shipped: { label: '已发货', color: 'bg-purple-100 text-purple-700', icon: Truck },
    delivered: { label: '已送达', color: 'bg-green-100 text-green-700', icon: CheckCircle },
    cancelled: { label: '已取消', color: 'bg-gray-100 text-gray-700', icon: XCircle },
  }

  const statuses = ['', 'pending', 'paid', 'shipped', 'delivered', 'cancelled']

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">我的订单</h1>
        <p className="text-gray-500">查看和管理您的订单</p>
      </div>

      {/* Filter */}
      <div className="flex gap-2 overflow-x-auto pb-2">
        {statuses.map((status) => (
          <button
            key={status}
            onClick={() => setFilterStatus(status)}
            className={`px-4 py-2 rounded-lg text-sm font-medium whitespace-nowrap transition-colors ${
              filterStatus === status
                ? 'bg-garden-600 text-white'
                : 'bg-white border border-gray-200 text-gray-600 hover:bg-gray-50'
            }`}
          >
            {status === '' ? '全部' : statusConfig[status]?.label}
          </button>
        ))}
      </div>

      {/* Orders */}
      {isLoading ? (
        <div className="text-center py-12 text-gray-500">加载中...</div>
      ) : orders.length === 0 ? (
        <div className="card text-center py-12">
          <Package className="w-16 h-16 text-gray-300 mx-auto mb-4" />
          <p className="text-gray-500">还没有订单</p>
        </div>
      ) : (
        <div className="space-y-4">
          {orders.map((order: any) => {
            const config = statusConfig[order.status] || statusConfig.pending
            const StatusIcon = config.icon
            return (
              <div key={order.id} className="card">
                <div className="card-header flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <span className="text-sm text-gray-500">订单号：</span>
                    <span className="font-medium">{order.order_no}</span>
                  </div>
                  <span className={`badge ${config.color}`}>
                    <StatusIcon className="w-3 h-3 mr-1" />
                    {config.label}
                  </span>
                </div>
                <div className="card-body">
                  {/* Order Items */}
                  {order.order_items?.map((item: any) => (
                    <div key={item.id} className="flex items-center gap-4 py-2">
                      <div className="w-12 h-12 bg-gray-100 rounded-lg flex items-center justify-center">
                        <Package className="w-6 h-6 text-gray-400" />
                      </div>
                      <div className="flex-1">
                        <p className="font-medium text-gray-900">{item.product_name}</p>
                        <p className="text-sm text-gray-500">
                          ¥{item.price?.toFixed(2)} x {item.quantity}
                        </p>
                      </div>
                      <p className="font-medium">¥{item.subtotal?.toFixed(2)}</p>
                    </div>
                  ))}

                  {/* Order Info */}
                  <div className="mt-4 pt-4 border-t border-gray-100 grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
                    <div>
                      <p className="text-gray-500">下单时间</p>
                      <p className="font-medium">
                        {new Date(order.created_at).toLocaleString('zh-CN')}
                      </p>
                    </div>
                    <div>
                      <p className="text-gray-500">收货人</p>
                      <p className="font-medium">{order.shipping_name}</p>
                    </div>
                    <div>
                      <p className="text-gray-500">联系电话</p>
                      <p className="font-medium">{order.shipping_phone}</p>
                    </div>
                    <div>
                      <p className="text-gray-500">收货地址</p>
                      <p className="font-medium truncate">{order.shipping_address}</p>
                    </div>
                  </div>

                  {/* Tracking */}
                  {order.tracking_number && (
                    <div className="mt-4 p-3 bg-garden-50 rounded-lg">
                      <p className="text-sm font-medium text-garden-800">
                        物流信息：{order.tracking_number}
                      </p>
                      {order.shipping_status && (
                        <p className="text-sm text-garden-600 mt-1">
                          状态：{order.shipping_status}
                        </p>
                      )}
                    </div>
                  )}

                  {/* Total */}
                  <div className="mt-4 flex items-center justify-between">
                    <span className="text-gray-600">订单总额</span>
                    <span className="text-xl font-bold text-garden-600">
                      ¥{order.total_amount?.toFixed(2)}
                    </span>
                  </div>

                  {/* Actions */}
                  {order.status === 'pending' && (
                    <div className="mt-4 flex gap-2">
                      <button
                        onClick={() => {
                          if (confirm('确定要取消这个订单吗？')) {
                            orderApi.cancel(order.id)
                          }
                        }}
                        className="btn-outline text-sm"
                      >
                        取消订单
                      </button>
                    </div>
                  )}
                </div>
              </div>
            )
          })}
        </div>
      )}
    </div>
  )
}
