import React, { useState, useEffect } from 'react'
import { Table, Card, Tag, Button, Modal, Form, Input, Select, message, Upload } from 'antd'
import { complaintApi } from '@/services/review'
import { Complaint } from '@/types'
import { formatDate, getComplaintStatusText } from '@/utils'
import { usePagination } from '@/hooks'
import { useAppSelector } from '@/store/hooks'

const ComplaintList: React.FC = () => {
  const { userInfo } = useAppSelector((state) => state.auth)
  const { page, pageSize, total, setPage, setPageSize, setTotal } = usePagination()
  const [complaints, setComplaints] = useState<Complaint[]>([])
  const [loading, setLoading] = useState(false)
  const [createModalVisible, setCreateModalVisible] = useState(false)
  const [handleModalVisible, setHandleModalVisible] = useState(false)
  const [handleComplaint, setHandleComplaint] = useState<Complaint | null>(null)
  const [handleResult, setHandleResult] = useState('')
  const [handleStatus, setHandleStatus] = useState('')
  const [form] = Form.useForm()

  useEffect(() => {
    loadComplaints()
  }, [page, pageSize])

  const loadComplaints = async () => {
    setLoading(true)
    try {
      const res = await complaintApi.getList({
        page,
        page_size: pageSize,
      })
      setComplaints(res.list)
      setTotal(res.total)
    } catch (error) {
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  const handleCreate = async () => {
    try {
      const values = await form.validateFields()
      await complaintApi.create(values)
      message.success('投诉已提交')
      setCreateModalVisible(false)
      form.resetFields()
      loadComplaints()
    } catch (error) {
      console.error(error)
    }
  }

  const handleProcess = async () => {
    if (!handleComplaint || !handleResult || !handleStatus) {
      message.warning('请填写处理结果')
      return
    }
    try {
      await complaintApi.handle(handleComplaint.id, {
        status: handleStatus,
        result: handleResult,
      })
      message.success('处理成功')
      setHandleModalVisible(false)
      setHandleComplaint(null)
      setHandleResult('')
      setHandleStatus('')
      loadComplaints()
    } catch (error) {
      console.error(error)
    }
  }

  const columns = [
    {
      title: '投诉ID',
      dataIndex: 'id',
      key: 'id',
    },
    {
      title: '标题',
      dataIndex: 'title',
      key: 'title',
      ellipsis: true,
    },
    {
      title: '订单号',
      dataIndex: ['order', 'order_no'],
      key: 'order_no',
      render: (_: any, record: Complaint) => record.order?.order_no || '-',
    },
    {
      title: '投诉人',
      key: 'complainant',
      render: (_: any, record: Complaint) => record.complainant?.nickname || '-',
    },
    {
      title: '被投诉人',
      key: 'respondent',
      render: (_: any, record: Complaint) => record.respondent?.nickname || '-',
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (text: string) => <Tag>{getComplaintStatusText(text)}</Tag>,
    },
    {
      title: '投诉时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (text: string) => formatDate(text),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Complaint) => (
        <div>
          {userInfo?.role === 'admin' && record.status === 'pending' && (
            <Button
              type="link"
              onClick={() => {
                setHandleComplaint(record)
                setHandleModalVisible(true)
              }}
            >
              处理
            </Button>
          )}
        </div>
      ),
    },
  ]

  return (
    <div>
      <div className="page-header">
        <h1 className="page-title">投诉管理</h1>
        <Button type="primary" onClick={() => setCreateModalVisible(true)}>
          发起投诉
        </Button>
      </div>

      <Card>
        <Table
          rowKey="id"
          loading={loading}
          dataSource={complaints}
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
        title="发起投诉"
        open={createModalVisible}
        onOk={handleCreate}
        onCancel={() => {
          setCreateModalVisible(false)
          form.resetFields()
        }}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="order_id"
            label="订单ID"
            rules={[{ required: true, message: '请输入订单ID' }]}
          >
            <Input placeholder="请输入订单ID" />
          </Form.Item>
          <Form.Item
            name="title"
            label="投诉标题"
            rules={[{ required: true, message: '请输入投诉标题' }]}
          >
            <Input placeholder="请输入投诉标题" />
          </Form.Item>
          <Form.Item
            name="content"
            label="投诉内容"
            rules={[{ required: true, message: '请输入投诉内容' }]}
          >
            <Input.TextArea rows={4} placeholder="请详细描述投诉内容" />
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title="处理投诉"
        open={handleModalVisible}
        onOk={handleProcess}
        onCancel={() => {
          setHandleModalVisible(false)
          setHandleComplaint(null)
          setHandleResult('')
          setHandleStatus('')
        }}
      >
        {handleComplaint && (
          <div style={{ marginBottom: 16 }}>
            <div>投诉标题：{handleComplaint.title}</div>
            <div>投诉内容：{handleComplaint.content}</div>
          </div>
        )}
        <Select
          placeholder="选择处理状态"
          style={{ width: '100%', marginBottom: 16 }}
          value={handleStatus}
          onChange={setHandleStatus}
          options={[
            { label: '已解决', value: 'resolved' },
            { label: '已驳回', value: 'rejected' },
          ]}
        />
        <Input.TextArea
          rows={3}
          placeholder="请输入处理结果"
          value={handleResult}
          onChange={(e) => setHandleResult(e.target.value)}
        />
      </Modal>
    </div>
  )
}

export default ComplaintList
