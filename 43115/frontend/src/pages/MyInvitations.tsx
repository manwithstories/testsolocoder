import React, { useState, useEffect } from 'react'
import { Table, Tag, Button, Modal, message, Card, Input } from 'antd'
import { invitationApi } from '@/services/order'
import { OrderInvitation } from '@/types'
import { formatDate, getInvitationStatusText, formatPrice } from '@/utils'
import { usePagination } from '@/hooks'

const MyInvitations: React.FC = () => {
  const { page, pageSize, total, setPage, setPageSize, setTotal } = usePagination()
  const [invitations, setInvitations] = useState<OrderInvitation[]>([])
  const [loading, setLoading] = useState(false)
  const [rejectModalVisible, setRejectModalVisible] = useState(false)
  const [rejectInvitation, setRejectInvitation] = useState<OrderInvitation | null>(null)
  const [rejectReason, setRejectReason] = useState('')

  useEffect(() => {
    loadInvitations()
  }, [page, pageSize])

  const loadInvitations = async () => {
    setLoading(true)
    try {
      const res = await invitationApi.getMyInvitations({
        page,
        page_size: pageSize,
      })
      setInvitations(res.list)
      setTotal(res.total)
    } catch (error) {
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  const handleAccept = async (invitation: OrderInvitation) => {
    try {
      await invitationApi.respond(invitation.id, { accepted: true })
      message.success('已接受邀请')
      loadInvitations()
    } catch (error) {
      console.error(error)
    }
  }

  const handleReject = async () => {
    if (!rejectInvitation || !rejectReason) {
      message.warning('请填写拒绝原因')
      return
    }
    try {
      await invitationApi.respond(rejectInvitation.id, {
        accepted: false,
        reject_reason: rejectReason,
      })
      message.success('已拒绝邀请')
      setRejectModalVisible(false)
      setRejectInvitation(null)
      setRejectReason('')
      loadInvitations()
    } catch (error) {
      console.error(error)
    }
  }

  const columns = [
    {
      title: '邀请ID',
      dataIndex: 'id',
      key: 'id',
    },
    {
      title: '订单号',
      dataIndex: 'order_no',
      key: 'order_no',
      render: (_: any, record: OrderInvitation) => record.order?.order_no || '-',
    },
    {
      title: '服务名称',
      dataIndex: 'service_name',
      key: 'service_name',
      render: (_: any, record: OrderInvitation) => record.order?.service_item?.name || '-',
    },
    {
      title: '预约时间',
      dataIndex: 'appointment_time',
      key: 'appointment_time',
      render: (_: any, record: OrderInvitation) =>
        formatDate(record.order?.appointment_time || ''),
    },
    {
      title: '订单金额',
      dataIndex: 'total_amount',
      key: 'total_amount',
      render: (_: any, record: OrderInvitation) =>
        formatPrice(record.order?.total_amount || 0),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (text: string) => <Tag>{getInvitationStatusText(text)}</Tag>,
    },
    {
      title: '过期时间',
      dataIndex: 'expired_at',
      key: 'expired_at',
      render: (text: string) => formatDate(text),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: OrderInvitation) =>
        record.status === 'pending' && (
          <div style={{ display: 'flex', gap: 8 }}>
            <Button type="primary" size="small" onClick={() => handleAccept(record)}>
              接受
            </Button>
            <Button
              size="small"
              danger
              onClick={() => {
                setRejectInvitation(record)
                setRejectModalVisible(true)
              }}
            >
              拒绝
            </Button>
          </div>
        ),
    },
  ]

  return (
    <div>
      <div className="page-header">
        <h1 className="page-title">预约邀请</h1>
      </div>

      <Card>
        <Table
          rowKey="id"
          loading={loading}
          dataSource={invitations}
          columns={columns}
          pagination={{
            current: page,
            pageSize,
            total,
            showSizeChanger: true,
            onChange: (p, ps) => {
              setPage(p)
              setPageSize(ps)
            },
          }}
        />
      </Card>

      <Modal
        title="拒绝邀请"
        open={rejectModalVisible}
        onOk={handleReject}
        onCancel={() => {
          setRejectModalVisible(false)
          setRejectInvitation(null)
          setRejectReason('')
        }}
      >
        <p>确定要拒绝该邀请吗？</p>
        <Input.TextArea
          rows={3}
          placeholder="请输入拒绝原因"
          value={rejectReason}
          onChange={(e) => setRejectReason(e.target.value)}
        />
      </Modal>
    </div>
  )
}

export default MyInvitations
