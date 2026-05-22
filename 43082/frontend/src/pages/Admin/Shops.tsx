import { useState, useEffect } from 'react'
import {
  Table, Button, Card, Typography, Tag, Space, Pagination,
  message, Modal, Form, Input, Select
} from 'antd'
import { CheckOutlined, CloseOutlined, EyeOutlined } from '@ant-design/icons'
import { adminAPI, shopAPI } from '@/api'
import { Shop } from '@/types'

const { Title } = Typography
const { Option } = Select
const { TextArea } = Input

const AdminShops = () => {
  const [shops, setShops] = useState<Shop[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize] = useState(10)
  const [status, setStatus] = useState<string>('')
  const [rejectModalVisible, setRejectModalVisible] = useState(false)
  const [currentShop, setCurrentShop] = useState<Shop | null>(null)
  const [rejectForm] = Form.useForm()

  useEffect(() => {
    loadShops()
  }, [page, status])

  const loadShops = async () => {
    try {
      const params: any = { page, pageSize }
      if (status) params.status = status
      const res = await shopAPI.list(params)
      setShops(res.data.data)
      setTotal(res.data.pagination.total)
    } catch (err) {
      console.error('加载店铺失败', err)
    }
  }

  const handleApprove = async (id: number) => {
    try {
      await shopAPI.review(id, { status: 'approved' })
      message.success('审核通过')
      loadShops()
    } catch (err: any) {
      message.error(err.message || '操作失败')
    }
  }

  const handleReject = (shop: Shop) => {
    setCurrentShop(shop)
    setRejectModalVisible(true)
  }

  const handleRejectSubmit = async (values: any) => {
    if (!currentShop) return
    try {
      await shopAPI.review(currentShop.id, { status: 'rejected', rejectReason: values.reason })
      message.success('已拒绝')
      setRejectModalVisible(false)
      rejectForm.resetFields()
      loadShops()
    } catch (err: any) {
      message.error(err.message || '操作失败')
    }
  }

  const getStatusTag = (status: string) => {
    const statusMap: Record<string, { color: string; text: string }> = {
      pending: { color: 'orange', text: '审核中' },
      approved: { color: 'green', text: '已通过' },
      rejected: { color: 'red', text: '已拒绝' },
      banned: { color: 'default', text: '已封禁' },
    }
    const info = statusMap[status] || { color: 'default', text: status }
    return <Tag color={info.color}>{info.text}</Tag>
  }

  const columns = [
    {
      title: '店铺',
      key: 'shop',
      render: (_: any, record: Shop) => (
        <div style={{ display: 'flex', gap: 12, alignItems: 'center' }}>
          <img
            src={record.logo || 'https://via.placeholder.com/60'}
            style={{ width: 60, height: 60, objectFit: 'cover', borderRadius: '50%' }}
          />
          <div>
            <div style={{ fontWeight: 'bold' }}>{record.name}</div>
            <div style={{ color: '#999', fontSize: 12 }}>{record.contactPhone}</div>
          </div>
        </div>
      ),
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
      ellipsis: true,
    },
    {
      title: '营业执照',
      dataIndex: 'businessLicense',
      key: 'businessLicense',
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => getStatusTag(status),
    },
    {
      title: '申请时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
    },
    {
      title: '操作',
      key: 'action',
      width: 200,
      render: (_: any, record: Shop) => (
        <Space>
          {record.status === 'pending' && (
            <>
              <Button
                type="primary"
                size="small"
                icon={<CheckOutlined />}
                onClick={() => handleApprove(record.id)}
              >
                通过
              </Button>
              <Button
                danger
                size="small"
                icon={<CloseOutlined />}
                onClick={() => handleReject(record)}
              >
                拒绝
              </Button>
            </>
          )}
        </Space>
      ),
    },
  ]

  return (
    <div>
      <Title level={3}>商家管理</Title>

      <Card>
        <div style={{ marginBottom: 16 }}>
          <Button.Group>
            <Button type={status === '' ? 'primary' : 'default'} onClick={() => { setStatus(''); setPage(1) }}>
              全部
            </Button>
            <Button type={status === 'pending' ? 'primary' : 'default'} onClick={() => { setStatus('pending'); setPage(1) }}>
              待审核
            </Button>
            <Button type={status === 'approved' ? 'primary' : 'default'} onClick={() => { setStatus('approved'); setPage(1) }}>
              已通过
            </Button>
            <Button type={status === 'rejected' ? 'primary' : 'default'} onClick={() => { setStatus('rejected'); setPage(1) }}>
              已拒绝
            </Button>
          </Button.Group>
        </div>

        <Table
          columns={columns}
          dataSource={shops}
          rowKey="id"
          pagination={false}
        />

        <div style={{ textAlign: 'center', marginTop: 24 }}>
          <Pagination
            current={page}
            pageSize={pageSize}
            total={total}
            onChange={setPage}
          />
        </div>
      </Card>

      <Modal
        title="拒绝入驻申请"
        open={rejectModalVisible}
        onCancel={() => setRejectModalVisible(false)}
        footer={null}
      >
        <Form form={rejectForm} layout="vertical" onFinish={handleRejectSubmit}>
          <Form.Item
            name="reason"
            label="拒绝原因"
            rules={[{ required: true, message: '请输入拒绝原因' }]}
          >
            <TextArea rows={4} placeholder="请输入拒绝原因" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" style={{ width: '100%' }}>
              确认拒绝
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default AdminShops
