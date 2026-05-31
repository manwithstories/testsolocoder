import { useState } from 'react'
import { Card, Table, Button, Modal, Form, Input, Select, DatePicker, Upload, message, Space } from 'antd'
import { PlusOutlined, UploadOutlined } from '@ant-design/icons'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useNavigate } from 'react-router-dom'
import { dailyRecordApi, reservationApi, petApi } from '@/services/api'
import { useAuthStore } from '@/context/AuthContext'
import type { DailyRecord } from '@/types'
import dayjs from 'dayjs'

export default function DailyRecords() {
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const { user } = useAuthStore()
  const [modalVisible, setModalVisible] = useState(false)
  const [form] = Form.useForm()

  const { data: recordsData, isLoading } = useQuery({
    queryKey: ['daily-records'],
    queryFn: () => dailyRecordApi.listByReservation({ page_size: 100 }),
  })

  const { data: reservationsData } = useQuery({
    queryKey: ['reservations', 'all'],
    queryFn: () => reservationApi.list({ page_size: 100, status: 'checked_in' }),
  })

  const records: DailyRecord[] = recordsData?.data?.items || []
  const reservations = reservationsData?.data?.items || []

  const createMutation = useMutation({
    mutationFn: (values: any) =>
      dailyRecordApi.create({
        ...values,
        record_date: values.record_date.toISOString(),
        photos: values.photos?.join(',') || '',
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['daily-records'] })
      message.success('创建成功')
      setModalVisible(false)
      form.resetFields()
    },
    onError: (err: any) => message.error(err.message || '创建失败'),
  })

  const columns = [
    {
      title: '记录日期',
      dataIndex: 'record_date',
      key: 'record_date',
      render: (d: string) => dayjs(d).format('YYYY-MM-DD'),
    },
    { title: '饮食', dataIndex: 'feed_status', key: 'feed_status' },
    { title: '活动', dataIndex: 'activity', key: 'activity' },
    { title: '健康', dataIndex: 'health_status', key: 'health_status' },
    { title: '心情', dataIndex: 'mood', key: 'mood' },
    { title: '备注', dataIndex: 'remark', key: 'remark' },
  ]

  return (
    <div className="space-y-4">
      <Card
        title="实时动态"
        extra={
          user?.role === 'keeper' && (
            <Button type="primary" icon={<PlusOutlined />} onClick={() => setModalVisible(true)}>
              添加记录
            </Button>
          )
        }
      >
        <Table columns={columns} dataSource={records} rowKey="id" loading={isLoading} pagination={{ pageSize: 10 }} />
      </Card>

      <Modal
        title="添加动态记录"
        open={modalVisible}
        onCancel={() => setModalVisible(false)}
        onOk={() => form.submit()}
        confirmLoading={createMutation.isPending}
        width={600}
      >
        <Form form={form} layout="vertical" onFinish={(v) => createMutation.mutate(v)} className="mt-4">
          <Form.Item name="reservation_id" label="选择预约" rules={[{ required: true }]}>
            <Select
              placeholder="选择进行中的预约"
              options={reservations.map((r: any) => ({
                value: r.id,
                label: `${r.order_no} - ${dayjs(r.check_in_date).format('MM-DD')} ~ ${dayjs(r.check_out_date).format('MM-DD')}`,
              }))}
            />
          </Form.Item>
          <Form.Item name="pet_id" label="选择宠物" rules={[{ required: true }]}>
            <Input placeholder="宠物ID（可从预约中获取）" />
          </Form.Item>
          <Form.Item name="record_date" label="记录日期" rules={[{ required: true }]}>
            <DatePicker className="w-full" />
          </Form.Item>
          <div className="grid grid-cols-2 gap-4">
            <Form.Item name="feed_status" label="饮食情况">
              <Input placeholder="饮食情况描述" />
            </Form.Item>
            <Form.Item name="mood" label="心情">
              <Input placeholder="心情描述" />
            </Form.Item>
          </div>
          <Form.Item name="activity" label="活动情况">
            <Input.TextArea rows={2} placeholder="活动情况描述" />
          </Form.Item>
          <Form.Item name="health_status" label="健康状况">
            <Input.TextArea rows={2} placeholder="健康状况描述" />
          </Form.Item>
          <Form.Item name="remark" label="备注">
            <Input.TextArea rows={2} placeholder="其他备注" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
