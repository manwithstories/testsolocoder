import React, { useState, useEffect } from 'react'
import { Table, Card, Tag, Button, Modal, Input, message, Select } from 'antd'
import { complaintApi } from '@/services/review'
import { Complaint } from '@/types'
import { formatDate, getComplaintStatusText } from '@/utils'
import { usePagination } from '@/hooks'

const ComplaintHandle: React.FC = () => {
  const { page, pageSize, total, setPage, setPageSize, setTotal } = usePagination()
  const [complaints, setComplaints] = useState<Complaint[]>([])
  const [loading, setLoading] = useState(false)
  const [handleModalVisible, setHandleModalVisible] = useState(false)
  const [selectedComplaint, setSelectedComplaint] = useState<Complaint | null>(null)
  const [handleResult, setHandleResult] = useState('')
  const [handleStatus, setHandleStatus] = useState('')
  const [status, setStatus] = useState<string | undefined>()

  useEffect(() => {
    loadComplaints()
  }, [page, pageSize, status])

  const loadComplaints = async () => {
    setLoading(true)
    try {
      const res = await complaintApi.getList({
        page,
        page_size: pageSize,
        status,
      })
      setComplaints(res.list)
      setTotal(res.total)
    } catch (error) {
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  const handleProcess = async () => {
    if (!selectedComplaint || !handleResult || !handleStatus) {
      message.warning('请填写处理结果')
      return
    }
    try {
      await complaintApi.handle(selectedComplaint.id, {
        status: handleStatus,
        result: handleResult,
      })
      message.success('处理成功')
      setHandleModalVisible(false)
      setSelectedComplaint(null)
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
      title: '订单号',
      dataIndex: 'order_no',
      key: 'order_no',
      render: (_: any, record: Complaint) => record.order?.order_no || '-',
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (text: string) => <Tag color={text === 'pending' ? 'orange' : text === 'resolved' ? 'green' : 'red'}>{getComplaintStatusText(text)}</Tag>,
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
        record.status === 'pending' && (
          <Button
            type="link"
            onClick={() => {
              setSelectedComplaint(record)
              setHandleModalVisible(true)
            }}
          >
            处理
          </Button>
        )
      ),
    },
  ]

  return (
    <div>
      <div className="page-header">
        <h1 className="page-title">投诉处理</h1>
      </div>

      <div className="filter-bar">
        <Select
          placeholder="投诉状态"
          style={{ width: 150 }}
          allowClear
          value={status}
          onChange={(value) => {
            setStatus(value)
            setPage(1)
          }}
          options={[
            { label: '待处理', value: 'pending' },
            { label: '已解决', value: 'resolved' },
            { label: '已驳回', value: 'rejected' },
          ]}
        />
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
        title="处理投诉"
        open={handleModalVisible}
        onOk={handleProcess}
        onCancel={() => {
          setHandleModalVisible(false)
          setSelectedComplaint(null)
          setHandleResult('')
          setHandleStatus('')
        }}
      >
        {selectedComplaint && (
          <div style={{ marginBottom: 16 }}>
            <div style={{ marginBottom: 8 }}>
              <strong>投诉标题：</strong>{selectedComplaint.title}
            </div>
            <div style={{ marginBottom: 8 }}>
              <strong>投诉内容：</strong>{selectedComplaint.content}
            </div>
            {selectedComplaint.images && selectedComplaint.images.length > 0 && (
              <div style={{ marginBottom: 8 }}>
                <strong>证据图片：</strong>
                <div style={{ display: 'flex', gap: 8, marginTop: 8 }}>
                  {selectedComplaint.images.split(',').map((img, index) => (
                    <img
                      key={index}
                      src={img}
                      alt="证据"
                      style={{ width: 80, height: 80, objectFit: 'cover' }}
                    />
                  ))}
                </div>
              </div>
            )}
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

export default ComplaintHandle
