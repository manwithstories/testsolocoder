import React, { useEffect, useState } from 'react'
import { Table, Tag, Button, Modal, Card, message, Space, Select, Empty, Spin } from 'antd'
import { CheckOutlined, CloseOutlined, EyeOutlined } from '@ant-design/icons'
import { listAdoptionApplications, reviewAdoptionApplication } from '../api/adoption'
import { AdoptionApplication, ReviewAdoptionRequest } from '../types'

const { Option } = Select

const AdoptionApplications: React.FC = () => {
  const [applications, setApplications] = useState<AdoptionApplication[]>([])
  const [loading, setLoading] = useState(false)
  const [statusFilter, setStatusFilter] = useState<string>('')
  const [detailModal, setDetailModal] = useState(false)
  const [selectedApp, setSelectedApp] = useState<AdoptionApplication | null>(null)

  useEffect(() => {
    loadApplications()
  }, [statusFilter])

  const loadApplications = async () => {
    setLoading(true)
    try {
      const response = await listAdoptionApplications({
        status: statusFilter || undefined,
        page_size: 100,
      })
      if (response.code === 0 && response.data) {
        setApplications((response.data as any).items || [])
      }
    } catch (error) {
      console.error('Failed to load applications:', error)
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
            resolve(input?.value || '申请不符合要求')
          },
          onCancel: () => resolve(''),
        })
      })
      if (!rejectReason) return
    }

    try {
      const req: ReviewAdoptionRequest = { action }
      if (rejectReason) req.reject_reason = rejectReason
      await reviewAdoptionApplication(id, req)
      message.success(action === 'approve' ? '已通过申请' : '已拒绝申请')
      loadApplications()
    } catch (error: any) {
      message.error(error.message || '操作失败')
    }
  }

  const columns = [
    { title: '申请编号', dataIndex: 'id', key: 'id' },
    {
      title: '宠物',
      key: 'pet',
      render: (_: any, r: AdoptionApplication) => r.pet?.name || '-',
    },
    {
      title: '申请人',
      key: 'adopter',
      render: (_: any, r: AdoptionApplication) => r.adopter?.name || '-',
    },
    { title: '申请原因', dataIndex: 'reason', key: 'reason', ellipsis: true },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => {
        const colorMap: Record<string, string> = {
          pending: 'orange',
          approved: 'blue',
          rejected: 'red',
          signed: 'green',
          completed: 'cyan',
          cancelled: 'default',
        }
        return <Tag color={colorMap[status]}>{status}</Tag>
      },
    },
    {
      title: '申请时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date: string) => new Date(date).toLocaleDateString(),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: AdoptionApplication) => (
        <Space>
          <Button
            type="link"
            size="small"
            icon={<EyeOutlined />}
            onClick={() => {
              setSelectedApp(record)
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
        <h2 style={{ margin: 0 }}>领养申请管理</h2>
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
          <Option value="signed">已签署</Option>
          <Option value="completed">已完成</Option>
        </Select>
      </div>

      <Spin spinning={loading}>
        {applications.length > 0 ? (
          <Card>
            <Table
              dataSource={applications}
              columns={columns}
              rowKey="id"
              pagination={{ pageSize: 10 }}
            />
          </Card>
        ) : (
          !loading && (
            <Card>
              <Empty description="暂无领养申请" />
            </Card>
          )
        )}
      </Spin>

      <Modal
        title="申请详情"
        open={detailModal}
        onCancel={() => setDetailModal(false)}
        footer={null}
        width={600}
      >
        {selectedApp && (
          <div>
            <p><strong>宠物:</strong> {selectedApp.pet?.name}</p>
            <p><strong>申请人:</strong> {selectedApp.adopter?.name}</p>
            <p><strong>联系邮箱:</strong> {selectedApp.adopter?.email}</p>
            <p><strong>申请原因:</strong> {selectedApp.reason}</p>
            <p><strong>居住情况:</strong> {selectedApp.living_situation}</p>
            <p><strong>住房类型:</strong> {selectedApp.housing_type}</p>
            <p><strong>收入水平:</strong> {selectedApp.income_level}</p>
            <p><strong>家庭成员:</strong> {selectedApp.family_members}人</p>
            <p><strong>有小孩:</strong> {selectedApp.has_children ? '是' : '否'}</p>
            <p><strong>有其他宠物:</strong> {selectedApp.has_other_pets ? '是' : '否'}</p>
            {selectedApp.other_pets_desc && <p><strong>其他宠物:</strong> {selectedApp.other_pets_desc}</p>}
            <p><strong>养宠经验:</strong> {selectedApp.pet_experience || '-'}</p>
            <p><strong>可承担医疗费用:</strong> {selectedApp.can_afford_vet ? '是' : '否'}</p>
            <p><strong>同意定期回访:</strong> {selectedApp.agree_to_visit ? '是' : '否'}</p>
            <p><strong>申请状态:</strong> <Tag>{selectedApp.status}</Tag></p>
            {selectedApp.reject_reason && <p><strong>拒绝原因:</strong> {selectedApp.reject_reason}</p>}
          </div>
        )}
      </Modal>
    </div>
  )
}

export default AdoptionApplications
