import React, { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Row, Col, Card, Button, Rate, Descriptions, Tag, Divider, Modal, Form, Select, DatePicker, InputNumber, Input, message, Avatar } from 'antd'
import { serviceApi } from '@/services/service'
import { orderApi } from '@/services/order'
import { addressApi } from '@/services/auth'
import { ServiceItem, Address } from '@/types'
import { formatPrice, formatDuration } from '@/utils'
import dayjs, { Dayjs } from 'dayjs'

const { TextArea } = Input

const ServiceDetail: React.FC = () => {
  const { id } = useParams()
  const navigate = useNavigate()
  const [service, setService] = useState<ServiceItem | null>(null)
  const [loading, setLoading] = useState(false)
  const [orderModalVisible, setOrderModalVisible] = useState(false)
  const [addresses, setAddresses] = useState<Address[]>([])
  const [form] = Form.useForm()
  const [submitting, setSubmitting] = useState(false)

  useEffect(() => {
    if (id) {
      loadService()
      loadAddresses()
    }
  }, [id])

  const loadService = async () => {
    setLoading(true)
    try {
      const res = await serviceApi.getDetail(Number(id))
      setService(res)
    } catch (error) {
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  const loadAddresses = async () => {
    try {
      const res = await addressApi.getList()
      setAddresses(res)
    } catch (error) {
      console.error(error)
    }
  }

  const disabledDate = (current: Dayjs) => {
    return current && current < dayjs().startOf('day')
  }

  const handleSubmitOrder = async () => {
    try {
      const values = await form.validateFields()
      setSubmitting(true)

      const res = await orderApi.create({
        service_item_id: Number(id),
        address_id: values.address_id,
        appointment_time: values.appointment_time.toISOString(),
        duration: values.duration,
        remark: values.remark,
      })

      message.success('预约成功！等待服务人员接单')
      setOrderModalVisible(false)
      navigate(`/orders/${res.id}`)
    } catch (error) {
      console.error(error)
    } finally {
      setSubmitting(false)
    }
  }

  if (!service) {
    return <div>加载中...</div>
  }

  return (
    <div>
      <Button icon={<ArrowLeft />} onClick={() => navigate(-1)} style={{ marginBottom: 16 }}>
        返回
      </Button>

      <Row gutter={24}>
        <Col span={16}>
          <Card loading={loading}>
            <h1 style={{ marginBottom: 16 }}>{service.name}</h1>
            <div style={{ marginBottom: 16 }}>
              <Rate disabled value={service.rating} allowHalf />
              <span style={{ marginLeft: 8, color: '#999' }}>
                {service.rating} 分 | {service.review_count} 条评价 | 已售 {service.order_count}
              </span>
            </div>

            <Descriptions column={2} bordered size="middle">
              <Descriptions.Item label="服务价格">
                <span style={{ color: '#ff4d4f', fontSize: 20, fontWeight: 600 }}>
                  {formatPrice(service.base_price)}
                </span>
                <span style={{ color: '#999' }}>/{service.price_unit}</span>
              </Descriptions.Item>
              <Descriptions.Item label="服务时长">
                {formatDuration(service.min_duration)} - {formatDuration(service.max_duration)}
              </Descriptions.Item>
              <Descriptions.Item label="服务分类">
                {service.category?.name}
              </Descriptions.Item>
              <Descriptions.Item label="服务区域">
                {service.service_areas?.map((area) => (
                  <Tag key={area.id}>{area.district}</Tag>
                ))}
              </Descriptions.Item>
            </Descriptions>

            <Divider />

            <h3>服务详情</h3>
            <p style={{ whiteSpace: 'pre-wrap', lineHeight: 1.8 }}>
              {service.description || '暂无详细描述'}
            </p>

            {service.provider && (
              <>
                <Divider />
                <h3>服务人员</h3>
                <div style={{ display: 'flex', alignItems: 'center', gap: 16 }}>
                  <Avatar size={64} src={service.provider.avatar}>
                    {service.provider.nickname?.charAt(0)}
                  </Avatar>
                  <div>
                    <div style={{ fontSize: 16, fontWeight: 500 }}>
                      {service.provider.nickname}
                    </div>
                    <div style={{ color: '#999' }}>
                      评分: {service.provider.rating} | 接单: {service.provider.order_count}
                    </div>
                  </div>
                </div>
              </>
            )}
          </Card>
        </Col>

        <Col span={8}>
          <Card>
            <div style={{ textAlign: 'center', marginBottom: 24 }}>
              <div style={{ fontSize: 32, fontWeight: 600, color: '#ff4d4f' }}>
                {formatPrice(service.base_price)}
              </div>
              <div style={{ color: '#999' }}>/{service.price_unit}</div>
            </div>
            <Button
              type="primary"
              size="large"
              block
              onClick={() => setOrderModalVisible(true)}
            >
              立即预约
            </Button>
          </Card>
        </Col>
      </Row>

      <Modal
        title="预约服务"
        open={orderModalVisible}
        onCancel={() => setOrderModalVisible(false)}
        footer={null}
        width={500}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="address_id"
            label="服务地址"
            rules={[{ required: true, message: '请选择服务地址' }]}
          >
            <Select
              placeholder="选择服务地址"
              options={addresses.map((addr) => ({
                label: `${addr.contact_name} ${addr.contact_phone} ${addr.province}${addr.city}${addr.district}${addr.address}`,
                value: addr.id,
              }))}
            />
          </Form.Item>
          <Form.Item
            name="appointment_time"
            label="预约时间"
            rules={[{ required: true, message: '请选择预约时间' }]}
          >
            <DatePicker
              showTime={{ format: 'HH:mm' }}
              format="YYYY-MM-DD HH:mm"
              style={{ width: '100%' }}
              disabledDate={disabledDate}
              disabledTime={() => ({
                disabledHours: () => Array.from({ length: 8 }, (_, i) => i),
              })}
            />
          </Form.Item>
          <Form.Item
            name="duration"
            label="服务时长(分钟)"
            rules={[{ required: true, message: '请输入服务时长' }]}
            initialValue={service.min_duration}
          >
            <InputNumber
              min={service.min_duration}
              max={service.max_duration}
              step={30}
              style={{ width: '100%' }}
            />
          </Form.Item>
          <Form.Item name="remark" label="备注">
            <TextArea rows={3} placeholder="有什么特殊要求请填写" />
          </Form.Item>
          <Form.Item>
            <Button
              type="primary"
              block
              loading={submitting}
              onClick={handleSubmitOrder}
            >
              提交预约
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

function ArrowLeft() {
  return <span>←</span>
}

export default ServiceDetail
