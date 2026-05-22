import { useState, useEffect } from 'react'
import {
  Table, Button, Card, Typography, Tag, Space, Pagination,
  message, Modal, Form, Select
} from 'antd'
import { CheckOutlined, CloseOutlined } from '@ant-design/icons'
import { adminAPI } from '@/api'
import { Dispute } from '@/types'

const { Title } = Typography
const { Option } = Select

const AdminDisputes = () => {
  const [disputes, setDisputes] = useState<Dispute[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize] = useState(10)
  const [status, setStatus] = useState<string>('')

  useEffect(() => {
    loadDisputes()
  }, [page, status])

  const loadDisputes = async () => {
    try {
      const params: any = { page, pageSize }
      if (status) params.status = status
      const res = await adminAPI.getDisputes(params)
      setDisputes(res.data.data)
      setTotal(res.data.pagination.total)
    } catch (err) {
      console.error('加载纠纷失败', err)
    }
  }

  const handleResolve = async (id: number, resolution: string) => {
    try {
      await adminAPI.resolveDispute(id, { result: resolution })
      message.success('处理成功')
      loadDisputes()
    } catch (err: any) {
      message.error(err.message || '处理失败')
    }
  }

  const getStatusTag = (status: string) => {
    const statusMap: Record<string, { color: string; text: string }> = {
      pending: { color: 'orange', text: '待处理' },
      resolved: { color: 'green', text: '已解决' },
      rejected: { color: 'red', text: '已驳回' },
    }
    const info = statusMap[status] || { color: 'default', text: status }
    return <Tag color={info.color}>{info.text}</Tag>
  }

  const columns = [
    {
      title: '纠纷ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '订单号',
      dataIndex: 'orderNo',
      key: 'orderNo',
      width: 180,
    },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
      width: 120,
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
      ellipsis: true,
    },
    {
      title: '申请人',
      dataIndex: 'userName',
      key: 'userName',
      width: 120,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => getStatusTag(status),
    },
    {
      title: '创建时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
    },
    {
      title: '操作',
      key: 'action',
      width: 180,
      render: (_: any, record: Dispute) => (
        <Space>
          {record.status === 'pending' && (
            <>
              <Button
                type="primary"
                size="small"
                icon={<CheckOutlined />}
                onClick={() => handleResolve(record.id, '同意退款')}
              >
                同意
              </Button>
              <Button
                danger
                size="small"
                icon={<CloseOutlined />}
                onClick={() => handleResolve(record.id, '纠纷不成立')}
              >
                驳回
              </Button>
            </>
          )}
        </Space>
      ),
    },
  ]

  return (
    <div>
      <Title level={3}>纠纷处理</Title>

      <Card>
        <div style={{ marginBottom: 16 }}>
          <Button.Group>
            <Button type={status === '' ? 'primary' : 'default'} onClick={() => { setStatus(''); setPage(1) }}>
              全部
            </Button>
            <Button type={status === 'pending' ? 'primary' : 'default'} onClick={() => { setStatus('pending'); setPage(1) }}>
              待处理
            </Button>
            <Button type={status === 'resolved' ? 'primary' : 'default'} onClick={() => { setStatus('resolved'); setPage(1) }}>
              已解决
            </Button>
          </Button.Group>
        </div>

        <Table
          columns={columns}
          dataSource={disputes}
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
    </div>
  )
}

export default AdminDisputes
