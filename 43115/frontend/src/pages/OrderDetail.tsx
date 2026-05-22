import React, { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Card, Descriptions, Tag, Button, Divider, Rate, Modal, Input, message, Timeline, Avatar } from 'antd'
import { orderApi } from '@/services/order'
import { Order } from '@/types'
import {
  formatPrice,
  formatDate,
  getOrderStatusText,
  getOrderStatusColor,
  formatAddress,
} from '@/utils'
import { useAppSelector } from '@/store/hooks'

const { TextArea } = Input

const OrderDetail: React.FC = () => {
  const { id } = useParams()
  const navigate = useNavigate()
  const { userInfo } = useAppSelector((state) => state.auth)
  const [order, setOrder] = useState<Order | null>(null)
  const [loading, setLoading] = useState(false)
  const [reviewModalVisible, setReviewModalVisible] = useState(false)
  const [complaintModalVisible, setComplaintModalVisible] = useState(false)

  useEffect(() => {
    if (id) {
      loadOrder()
    }
  }, [id])

  const loadOrder = async () => {
    setLoading(true)
    try {
      const res = await orderApi.getDetail(Number(id))
      setOrder(res)
    } catch (error) {
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  const handleStart = async () => {
    try {
      await orderApi.start(Number(id))
      message.success('服务已开始')
      loadOrder()
    } catch (error) {
      console.error(error)
    }
  }

  const handleComplete = async () => {
    try {
      await orderApi.complete(Number(id))
      message.success('服务已完成')
      loadOrder()
    } catch (error) {
      console.error(error)
    }
  }

  const handleCancel = async (reason: string) => {
    try {
      await orderApi.cancel(Number(id), { reason })
      message.success('订单已取消')
      loadOrder()
    } catch (error) {
      console.error(error)
    }
  }

  if (!order) {
    return <div>加载中...</div>
  }

  const getTimelineItems = () => {
    const items: any[] = [
      {
        color: 'blue',
        children: (
          <div>
            <div>订单创建</div>
            <div style={{ color: '#999' }}>{formatDate(order.created_at)}</div>
          </div>
        ),
      },
    ]

    if (order.status !== 'pending') {
      items.push({
        color: 'blue',
        children: (
          <div>
            <div>订单确认</div>
            <div style={{ color: '#999' }}>服务人员已接单</div>
          </div>
        ),
      })
    }

    if (order.actual_start_time) {
      items.push({
        color: 'cyan',
        children: (
          <div>
            <div>服务开始</div>
            <div style={{ color: '#999' }}>{formatDate(order.actual_start_time)}</div>
          </div>
        ),
      })
    }

    if (order.completed_at) {
      items.push({
        color: 'green',
        children: (
          <div>
            <div>服务完成</div>
            <div style={{ color: '#999' }}>{formatDate(order.completed_at)}</div>
          </div>
        ),
      })
    }

    if (order.cancelled_at) {
      items.push({
        color: 'red',
        children: (
          <div>
            <div>订单取消</div>
            <div style={{ color: '#999' }}>{formatDate(order.cancelled_at)}</div>
            <div style={{ color: '#999' }}>原因：{order.cancel_reason}</div>
          </div>
        ),
      })
    }

    return items
  }

  return (
    <div>
      <Button icon={<span>←</span>} onClick={() => navigate(-1)} style={{ marginBottom: 16 }}>
        返回
      </Button>

      <Card loading={loading}>
        <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
          <div>
            <h2 style={{ margin: 0 }}>订单号：{order.order_no}</h2>
            <Tag color={getOrderStatusColor(order.status)} style={{ marginTop: 8 }}>
              {getOrderStatusText(order.status)}
            </Tag>
          </div>
          <div style={{ display: 'flex', gap: 8 }}>
            {userInfo?.role === 'service_provider' && order.status === 'confirmed' && (
              <Button type="primary" onClick={handleStart}>
                开始服务
              </Button>
            )}
            {userInfo?.role === 'service_provider' && order.status === 'in_service' && (
              <Button type="primary" onClick={handleComplete}>
                完成服务
              </Button>
            )}
            {(order.status === 'pending' || order.status === 'confirmed') && (
              <Button danger onClick={() => handleCancel('客户取消')}>
                取消订单
              </Button>
            )}
          </div>
        </div>

        <Descriptions column={2} bordered size="middle">
          <Descriptions.Item label="服务名称">
            {order.service_item?.name}
          </Descriptions.Item>
          <Descriptions.Item label="服务时长">
            {order.duration} 分钟
          </Descriptions.Item>
          <Descriptions.Item label="预约时间">
            {formatDate(order.appointment_time)}
          </Descriptions.Item>
          <Descriptions.Item label="服务地址">
            {order.service_address}
          </Descriptions.Item>
          <Descriptions.Item label="联系人">
            {order.contact_name}
          </Descriptions.Item>
          <Descriptions.Item label="联系电话">
            {order.contact_phone}
          </Descriptions.Item>
          {order.remark && (
            <Descriptions.Item label="备注" span={2}>
              {order.remark}
            </Descriptions.Item>
          )}
        </Descriptions>

        <Divider />

        <Descriptions column={3} bordered size="middle">
          <Descriptions.Item label="服务单价">
            {formatPrice(order.base_price)}/小时
          </Descriptions.Item>
          <Descriptions.Item label="订单金额">
            <span style={{ color: '#ff4d4f', fontWeight: 600 }}>
              {formatPrice(order.total_amount)}
            </span>
          </Descriptions.Item>
          <Descriptions.Item label="平台佣金">
            {formatPrice(order.platform_fee)}
          </Descriptions.Item>
          {order.penalty_amount > 0 && (
            <Descriptions.Item label="违约金" span={3}>
              <span style={{ color: '#ff4d4f' }}>{formatPrice(order.penalty_amount)}</span>
            </Descriptions.Item>
          )}
        </Descriptions>

        {order.customer && (
          <>
            <Divider />
            <h3>客户信息</h3>
            <div style={{ display: 'flex', alignItems: 'center', gap: 16 }}>
              <Avatar src={order.customer.avatar}>
                {order.customer.nickname?.charAt(0)}
              </Avatar>
              <div>
                <div style={{ fontWeight: 500 }}>{order.customer.nickname}</div>
                <div style={{ color: '#999' }}>手机号：{order.customer.phone}</div>
              </div>
            </div>
          </>
        )}

        {order.provider && (
          <>
            <Divider />
            <h3>服务人员</h3>
            <div style={{ display: 'flex', alignItems: 'center', gap: 16 }}>
              <Avatar src={order.provider.avatar}>
                {order.provider.nickname?.charAt(0)}
              </Avatar>
              <div>
                <div style={{ fontWeight: 500 }}>
                  {order.provider.nickname}
                  <Tag color="blue" style={{ marginLeft: 8 }}>
                    评分 {order.provider.rating}
                  </Tag>
                </div>
                <div style={{ color: '#999' }}>已接单 {order.provider.order_count} 次</div>
              </div>
            </div>
          </>
        )}

        <Divider />

        <h3>订单进度</h3>
        <Timeline items={getTimelineItems()} />
      </Card>
    </div>
  )
}

export default OrderDetail
