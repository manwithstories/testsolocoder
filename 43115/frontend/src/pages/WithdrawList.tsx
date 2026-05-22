import React, { useState, useEffect } from 'react'
import { Table, Card, Tag, Button, Modal, message, Input } from 'antd'
import { withdrawApi } from '@/services/bill'
import { WithdrawRequest } from '@/types'
import { formatPrice, formatDate, getWithdrawStatusText } from '@/utils'
import { usePagination } from '@/hooks'
import { useAppSelector } from '@/store/hooks'

const WithdrawList: React.FC = () => {
  const { userInfo } = useAppSelector((state) => state.auth)
  const { page, pageSize, total, setPage, setPageSize, setTotal } = usePagination()
  const [withdraws, setWithdraws] = useState<WithdrawRequest[]>([])
  const [loading, setLoading] = useState(false)
  const [detailModalVisible, setDetailModalVisible] = useState(false)
  const [selectedWithdraw, setSelectedWithdraw] = useState<WithdrawRequest | null>(null)
  const [rejectReason, setRejectReason] = useState('')

  useEffect(() => {
    loadWithdraws()
  }, [page, pageSize])

  const loadWithdraws = async () => {
    setLoading(true)
    try {
      const res = await withdrawApi.getList({
        page,
        page_size: pageSize,
      })
      setWithdraws(res.list)
      setTotal(res.total)
    } catch (error) {
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  const handleApprove = async (id: number) => {
    try {
      await withdrawApi.handle(id, { approved: true })
      message.success('提现申请已通过')
      loadWithdraws()
    } catch (error) {
      console.error(error)
    }
  }

  const handleReject = async () => {
    if (!selectedWithdraw || !rejectReason) {
      message.warning('请填写拒绝原因')
      return
    }
    try {
      await withdrawApi.handle(selectedWithdraw.id, { approved: false, remark: rejectReason })
      message.success('已拒绝提现申请')
      setDetailModalVisible(false)
      setSelectedWithdraw(null)
      setRejectReason('')
      loadWithdraws()
    } catch (error) {
      console.error(error)
    }
  }

  const columns = [
    {
      title: '申请ID',
      dataIndex: 'id',
      key: 'id',
    },
    {
      title: '申请用户',
      dataIndex: 'bank_name',
      key: 'bank_name',
      render: (_: any, record: WithdrawRequest) => record.bank_name || '-',
    },
    {
      title: '提现金额',
      dataIndex: 'amount',
      key: 'amount',
      render: (text: number) => <span style={{ color: '#ff4d4f', fontWeight: 600 }}>{formatPrice(text)}</span>,
    },
    {
      title: '收款方式',
      key: 'payment',
      render: (_: any, record: WithdrawRequest) => (
        <span>{record.bank_name} - {record.bank_account}</span>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (text: string) => <Tag>{getWithdrawStatusText(text)}</Tag>,
    },
    {
      title: '申请时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (text: string) => formatDate(text),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: WithdrawRequest) => (
        <div>
          {userInfo?.role === 'admin' && record.status === 'pending' && (
            <div style={{ display: 'flex', gap: 8 }}>
              <Button type="link" size="small" onClick={() => handleApprove(record.id)}>
                通过
              </Button>
              <Button
                type="link"
                size="small"
                danger
                onClick={() => {
                  setSelectedWithdraw(record)
                  setDetailModalVisible(true)
                }}
              >
                拒绝
              </Button>
            </div>
          )}
          {record.status !== 'pending' && record.handle_remark && (
            <Tag color="red">{record.handle_remark}</Tag>
          )}
        </div>
      ),
    },
  ]

  return (
    <div>
      <div className="page-header">
        <h1 className="page-title">提现审核</h1>
      </div>

      <Card>
        <Table
          rowKey="id"
          loading={loading}
          dataSource={withdraws}
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
        title="拒绝提现"
        open={detailModalVisible}
        onOk={handleReject}
        onCancel={() => {
          setDetailModalVisible(false)
          setSelectedWithdraw(null)
          setRejectReason('')
        }}
      >
        {selectedWithdraw && (
          <div style={{ marginBottom: 16 }}>
            <div>收款银行：{selectedWithdraw.bank_name}</div>
            <div>提现金额：{formatPrice(selectedWithdraw.amount)}</div>
          </div>
        )}
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

export default WithdrawList
