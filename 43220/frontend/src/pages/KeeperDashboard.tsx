import { useState } from 'react'
import { Card, Row, Col, Statistic, Table, Button, Tag, Modal, Form, Input, DatePicker, message } from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { dailyRecordApi, reservationApi } from '@/services/api'
import { useAuthStore } from '@/context/AuthContext'
import type { DailyRecord, Reservation } from '@/types'
import dayjs from 'dayjs'

export default function KeeperDashboard() {
  const queryClient = useQueryClient()
  const { user } = useAuthStore()
  const [recordModal, setRecordModal] = useState(false)
  const [recordForm] = Form.useForm()

  const { data: recordsData } = useQuery({
    queryKey: ['daily-records', 'today'],
    queryFn: () => dailyRecordApi.listByReservation({ page_size: 10 }),
  })

  const { data: reservationsData } = useQuery({
    queryKey: ['reservations', 'active'],
    queryFn: () => reservationApi.list({ status: 'checked_in', page_size: 10 }),
  })

  const records: DailyRecord[] = recordsData?.data?.items || []
  const reservations: Reservation[] = reservationsData?.data?.items || []

  const createRecordMutation = useMutation({
    mutationFn: (values: any) =>
      dailyRecordApi.create({
        ...values,
        record_date: values.record_date.toISOString(),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['daily-records'] })
      message.success('记录成功')
      setRecordModal(false)
      recordForm.resetFields()
    },
    onError: (err: any) => message.error(err.message || '记录失败'),
  })

  const recordColumns = [
    {
      title: '日期',
      dataIndex: 'record_date',
      key: 'record_date',
      render: (d: string) => dayjs(d).format('YYYY-MM-DD'),
    },
    { title: '饮食', dataIndex: 'feed_status', key: 'feed_status' },
    { title: '活动', dataIndex: 'activity', key: 'activity' },
    { title: '健康', dataIndex: 'health_status', key: 'health_status' },
    { title: '心情', dataIndex: 'mood', key: 'mood' },
  ]

  const resColumns = [
    { title: '订单号', dataIndex: 'order_no', key: 'order_no' },
    { title: '宠物ID', dataIndex: 'pet_id', key: 'pet_id' },
    {
      title: '入住',
      dataIndex: 'check_in_date',
      key: 'check_in_date',
      render: (d: string) => dayjs(d).format('YYYY-MM-DD'),
    },
    {
      title: '退房',
      dataIndex: 'check_out_date',
      key: 'check_out_date',
      render: (d: string) => dayjs(d).format('YYYY-MM-DD'),
    },
  ]

  return (
    <div className="space-y-4">
      <Card title="管家工作台">
        <Row gutter={16} className="mb-4">
          <Col xs={12} md={8}>
            <Card size="small">
              <Statistic title="今日任务" value={records.length} />
            </Card>
          </Col>
          <Col xs={12} md={8}>
            <Card size="small">
              <Statistic title="照顾中宠物" value={reservations.length} />
            </Card>
          </Col>
          <Col xs={12} md={8}>
            <Card size="small">
              <Statistic title="待记录" value={Math.max(0, reservations.length - records.length)} />
            </Card>
          </Col>
        </Row>

        <Card
          title="今日记录"
          size="small"
          extra={
            <Button type="primary" icon={<PlusOutlined />} onClick={() => setRecordModal(true)}>
              快速记录
            </Button>
          }
        >
          <Table columns={recordColumns} dataSource={records} rowKey="id" size="small" pagination={false} />
        </Card>

        <Card title="我的任务" size="small" className="mt-4">
          <Table columns={resColumns} dataSource={reservations} rowKey="id" size="small" pagination={false} />
        </Card>
      </Card>

      <Modal
        title="快速添加记录"
        open={recordModal}
        onCancel={() => setRecordModal(false)}
        onOk={() => recordForm.submit()}
        confirmLoading={createRecordMutation.isPending}
      >
        <Form form={recordForm} layout="vertical" onFinish={(v) => createRecordMutation.mutate(v)} className="mt-4">
          <Form.Item name="reservation_id" label="选择预约" rules={[{ required: true }]}>
            <Input placeholder="预约ID" />
          </Form.Item>
          <Form.Item name="pet_id" label="宠物ID" rules={[{ required: true }]}>
            <Input placeholder="宠物ID" />
          </Form.Item>
          <Form.Item name="record_date" label="记录日期" rules={[{ required: true }]}>
            <DatePicker className="w-full" />
          </Form.Item>
          <Form.Item name="feed_status" label="饮食情况">
            <Input placeholder="饮食情况" />
          </Form.Item>
          <Form.Item name="activity" label="活动情况">
            <Input.TextArea rows={2} placeholder="活动情况" />
          </Form.Item>
          <Form.Item name="health_status" label="健康状况">
            <Input.TextArea rows={2} placeholder="健康状况" />
          </Form.Item>
          <Form.Item name="mood" label="心情">
            <Input placeholder="心情描述" />
          </Form.Item>
          <Form.Item name="remark" label="备注">
            <Input.TextArea rows={2} placeholder="备注" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
