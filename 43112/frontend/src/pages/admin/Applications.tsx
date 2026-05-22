import React, { useEffect, useState } from 'react'
import { Table, Tag, Button, Space, Modal, Form, Input, message, Select, Descriptions } from 'antd'
import { userApi } from '@/services'
import { InstructorApplication } from '@/types'
import dayjs from 'dayjs'

const statusMap: Record<string, { color: string; text: string }> = {
  pending: { color: 'orange', text: '待审核' },
  approved: { color: 'green', text: '已通过' },
  rejected: { color: 'red', text: '已拒绝' },
}

const AdminApplications: React.FC = () => {
  const [data, setData] = useState<InstructorApplication[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize] = useState(10)
  const [loading, setLoading] = useState(false)
  const [statusFilter, setStatusFilter] = useState('')
  const [detailVisible, setDetailVisible] = useState(false)
  const [selectedApp, setSelectedApp] = useState<InstructorApplication | null>(null)
  const [reviewForm] = Form.useForm()

  const loadApplications = async (p = page) => {
    setLoading(true)
    try {
      const params: any = { page: p, page_size: pageSize }
      if (statusFilter) params.status = statusFilter
      const res = await userApi.listApplications(params)
      if (res.code === 0 && res.data) {
        setData(res.data.items)
        setTotal(res.data.total)
      }
    } catch (error) {
      console.error('Failed to load applications:', error)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadApplications()
  }, [])

  const handleReview = async (status: string) => {
    if (!selectedApp) return
    try {
      const values = await reviewForm.validateFields()
      const res = await userApi.reviewApplication(selectedApp.id, { status, ...values })
      if (res.code === 0) {
        message.success(status === 'approved' ? '审核通过' : '已拒绝')
        setDetailVisible(false)
        reviewForm.resetFields()
        loadApplications()
      }
    } catch (error: any) {
      if (error.errorFields) return
      message.error(error.message || '操作失败')
    }
  }

  const columns = [
    {
      title: '申请人',
      key: 'user',
      render: (_: any, record: InstructorApplication) => record.user?.username || '-',
    },
    { title: '真实姓名', dataIndex: 'real_name', key: 'real_name' },
    {
      title: '资质',
      dataIndex: 'qualification',
      key: 'qualification',
      ellipsis: true,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => {
        const info = statusMap[status] || { color: 'default', text: status }
        return <Tag color={info.color}>{info.text}</Tag>
      },
    },
    {
      title: '申请时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (val: string) => dayjs(val).format('YYYY-MM-DD HH:mm'),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: InstructorApplication) => (
        <Space>
          <Button type="link" onClick={() => { setSelectedApp(record); setDetailVisible(true) }}>
            查看详情
          </Button>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <h2>讲师资质审核</h2>
      <Space style={{ marginBottom: 16 }}>
        <Select
          placeholder="全部状态"
          style={{ width: 140 }}
          allowClear
          value={statusFilter || undefined}
          onChange={(val) => { setStatusFilter(val || ''); loadApplications(1) }}
        >
          {Object.entries(statusMap).map(([key, val]) => (
            <Select.Option key={key} value={key}>{val.text}</Select.Option>
          ))}
        </Select>
      </Space>
      <Table
        rowKey="id"
        columns={columns}
        dataSource={data}
        loading={loading}
        pagination={{
          current: page,
          total,
          pageSize,
          onChange: (p) => { setPage(p); loadApplications(p) },
        }}
      />
      <Modal
        title="讲师资质详情"
        open={detailVisible}
        onCancel={() => setDetailVisible(false)}
        footer={selectedApp?.status === 'pending' ? [
          <Button key="reject" danger onClick={() => handleReview('rejected')}>
            拒绝
          </Button>,
          <Button key="approve" type="primary" onClick={() => handleReview('approved')}>
            通过
          </Button>,
        ] : null}
        width={600}
      >
        {selectedApp && (
          <div>
            <Descriptions column={1} bordered size="small">
              <Descriptions.Item label="用户名">{selectedApp.user?.username}</Descriptions.Item>
              <Descriptions.Item label="邮箱">{selectedApp.user?.email}</Descriptions.Item>
              <Descriptions.Item label="真实姓名">{selectedApp.real_name}</Descriptions.Item>
              <Descriptions.Item label="资质说明">{selectedApp.qualification}</Descriptions.Item>
              <Descriptions.Item label="教学经历">{selectedApp.experience || '-'}</Descriptions.Item>
              <Descriptions.Item label="证书链接">{selectedApp.certificates || '-'}</Descriptions.Item>
              <Descriptions.Item label="申请时间">{dayjs(selectedApp.created_at).format('YYYY-MM-DD HH:mm')}</Descriptions.Item>
              <Descriptions.Item label="状态">
                <Tag color={statusMap[selectedApp.status]?.color}>{statusMap[selectedApp.status]?.text}</Tag>
              </Descriptions.Item>
              {selectedApp.review_remark && (
                <Descriptions.Item label="审核备注">{selectedApp.review_remark}</Descriptions.Item>
              )}
            </Descriptions>
            {selectedApp.status === 'pending' && (
              <Form form={reviewForm} layout="vertical" style={{ marginTop: 16 }}>
                <Form.Item name="review_remark" label="审核备注">
                  <Input.TextArea rows={3} placeholder="请输入审核备注..." />
                </Form.Item>
              </Form>
            )}
          </div>
        )}
      </Modal>
    </div>
  )
}

export default AdminApplications
