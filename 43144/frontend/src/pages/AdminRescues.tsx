import React, { useEffect, useState } from 'react'
import { Table, Tag, Card, Button, Modal, message, Space, Select, Empty, Spin } from 'antd'
import { CheckOutlined, CloseOutlined, EyeOutlined } from '@ant-design/icons'
import { listRescueStations, reviewRescueStation } from '../api/rescue'
import { RescueStation, ReviewRescueRequest } from '../types'

const { Option } = Select

const AdminRescues: React.FC = () => {
  const [rescues, setRescues] = useState<RescueStation[]>([])
  const [loading, setLoading] = useState(false)
  const [statusFilter, setStatusFilter] = useState<string>('')
  const [detailModal, setDetailModal] = useState(false)
  const [selectedRescue, setSelectedRescue] = useState<RescueStation | null>(null)

  useEffect(() => {
    loadRescues()
  }, [statusFilter])

  const loadRescues = async () => {
    setLoading(true)
    try {
      const response = await listRescueStations({
        status: statusFilter || undefined,
        page_size: 100,
      })
      if (response.code === 0 && response.data) {
        setRescues((response.data as any).items || [])
      }
    } catch (error) {
      console.error('Failed to load rescues:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleReview = async (id: number, action: 'approve' | 'reject') => {
    let rejectReason = ''
    if (action === 'reject') {
      rejectReason = await new Promise<string>((resolve) => {
        Modal.confirm({
          title: '拒绝原因',
          content: (
            <input
              id="reject-reason-input"
              placeholder="请输入拒绝原因"
              style={{ width: '100%', height: 32, padding: '0 8px' }}
            />
          ),
          onOk: () => {
            const input = document.getElementById('reject-reason-input') as HTMLInputElement
            resolve(input?.value || '资质不符合要求')
          },
          onCancel: () => resolve(''),
        })
      })
      if (!rejectReason) return
    }

    try {
      const req: ReviewRescueRequest = { action }
      if (rejectReason) req.reject_reason = rejectReason
      await reviewRescueStation(id, req)
      message.success(action === 'approve' ? '已通过审核' : '已拒绝审核')
      loadRescues()
    } catch (error: any) {
      message.error(error.message || '操作失败')
    }
  }

  const columns = [
    { title: '编号', dataIndex: 'id', key: 'id' },
    { title: '名称', dataIndex: 'name', key: 'name' },
    { title: '联系人', dataIndex: 'contact_person', key: 'contact_person' },
    { title: '联系电话', dataIndex: 'contact_phone', key: 'contact_phone' },
    { title: '联系邮箱', dataIndex: 'contact_email', key: 'contact_email' },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => {
        const colorMap: Record<string, string> = {
          pending: 'orange',
          approved: 'green',
          rejected: 'red',
        }
        return <Tag color={colorMap[status]}>{status}</Tag>
      },
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: RescueStation) => (
        <Space>
          <Button
            type="link"
            size="small"
            icon={<EyeOutlined />}
            onClick={() => {
              setSelectedRescue(record)
              setDetailModal(true)
            }}
          >
            查看
          </Button>
          {record.status === 'pending' && (
            <>
              <Button
                type="link"
                size="small"
                icon={<CheckOutlined />}
                style={{ color: 'green' }}
                onClick={() => handleReview(record.id, 'approve')}
              >
                通过
              </Button>
              <Button
                type="link"
                size="small"
                danger
                icon={<CloseOutlined />}
                onClick={() => handleReview(record.id, 'reject')}
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
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h2 style={{ margin: 0 }}>救助站管理</h2>
        <Select
          placeholder="筛选状态"
          allowClear
          style={{ width: 150 }}
          value={statusFilter || undefined}
          onChange={setStatusFilter}
        >
          <Option value="pending">待审核</Option>
          <Option value="approved">已通过</Option>
          <Option value="rejected">已拒绝</Option>
        </Select>
      </div>

      <Spin spinning={loading}>
        {rescues.length > 0 ? (
          <Card>
            <Table
              dataSource={rescues}
              columns={columns}
              rowKey="id"
              pagination={{ pageSize: 10 }}
            />
          </Card>
        ) : (
          !loading && (
            <Card>
              <Empty description="暂无救助站" />
            </Card>
          )
        )}
      </Spin>

      <Modal
        title="救助站详情"
        open={detailModal}
        onCancel={() => setDetailModal(false)}
        footer={null}
        width={600}
      >
        {selectedRescue && (
          <div>
            <p><strong>名称:</strong> {selectedRescue.name}</p>
            <p><strong>地址:</strong> {selectedRescue.address || '-'}</p>
            <p><strong>联系人:</strong> {selectedRescue.contact_person || '-'}</p>
            <p><strong>联系电话:</strong> {selectedRescue.contact_phone || '-'}</p>
            <p><strong>联系邮箱:</strong> {selectedRescue.contact_email || '-'}</p>
            <p><strong>执照编号:</strong> {selectedRescue.license_number || '-'}</p>
            <p><strong>描述:</strong> {selectedRescue.description || '-'}</p>
            <p><strong>状态:</strong> <Tag>{selectedRescue.status}</Tag></p>
            {selectedRescue.reject_reason && (
              <p><strong>拒绝原因:</strong> {selectedRescue.reject_reason}</p>
            )}
          </div>
        )}
      </Modal>
    </div>
  )
}

export default AdminRescues
