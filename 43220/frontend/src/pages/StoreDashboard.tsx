import { Card, Row, Col, Statistic, Table, Button, Tag, Space, Modal, Form, Input, InputNumber, Select, message } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useNavigate } from 'react-router-dom'
import { packageApi, reservationApi } from '@/services/api'
import { useAuthStore } from '@/context/AuthContext'
import type { BoardingPackage, Reservation } from '@/types'
import { useState } from 'react'
import dayjs from 'dayjs'

export default function StoreDashboard() {
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const { user } = useAuthStore()
  const [packageModal, setPackageModal] = useState(false)
  const [editingPkg, setEditingPkg] = useState<BoardingPackage | null>(null)
  const [pkgForm] = Form.useForm()

  const { data: packagesData } = useQuery({
    queryKey: ['packages', 'store'],
    queryFn: () => packageApi.list({ store_id: user?.id, page_size: 100 }),
  })

  const { data: reservationsData } = useQuery({
    queryKey: ['reservations', 'today'],
    queryFn: () =>
      reservationApi.list({
        status: 'checked_in',
        page_size: 10,
        start_date: new Date().toISOString(),
      }),
  })

  const packages: BoardingPackage[] = packagesData?.data?.items || []
  const reservations: Reservation[] = reservationsData?.data?.items || []

  const createPkgMutation = useMutation({
    mutationFn: (values: any) =>
      editingPkg ? packageApi.update(editingPkg.id, values) : packageApi.create(values),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['packages'] })
      message.success(editingPkg ? '更新成功' : '创建成功')
      setPackageModal(false)
      pkgForm.resetFields()
      setEditingPkg(null)
    },
    onError: (err: any) => message.error(err.message || '操作失败'),
  })

  const deletePkgMutation = useMutation({
    mutationFn: (id: string) => packageApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['packages'] })
      message.success('删除成功')
    },
    onError: (err: any) => message.error(err.message || '删除失败'),
  })

  const handleAddPkg = () => {
    setEditingPkg(null)
    pkgForm.resetFields()
    setPackageModal(true)
  }

  const handleEditPkg = (pkg: BoardingPackage) => {
    setEditingPkg(pkg)
    pkgForm.setFieldsValue(pkg)
    setPackageModal(true)
  }

  const pkgColumns = [
    { title: '名称', dataIndex: 'name', key: 'name' },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
      render: (t: string) => <Tag color={t === 'daycare' ? 'blue' : 'green'}>{t === 'daycare' ? '日托' : '寄养'}</Tag>,
    },
    { title: '价格/天', dataIndex: 'price_per_day', key: 'price_per_day', render: (p: number) => `¥${p}` },
    { title: '容量', dataIndex: 'capacity', key: 'capacity' },
    {
      title: '操作',
      key: 'actions',
      render: (_: any, pkg: BoardingPackage) => (
        <Space>
          <Button icon={<EditOutlined />} size="small" onClick={() => handleEditPkg(pkg)}>
            编辑
          </Button>
          <Button icon={<DeleteOutlined />} size="small" danger onClick={() => deletePkgMutation.mutate(pkg.id)}>
            删除
          </Button>
        </Space>
      ),
    },
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
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (s: string) => <Tag color="green">已入住</Tag>,
    },
  ]

  return (
    <div className="space-y-4">
      <Card title="门店管理">
        <Row gutter={16} className="mb-4">
          <Col xs={12} md={6}>
            <Card size="small">
              <Statistic title="套餐数量" value={packages.length} />
            </Card>
          </Col>
          <Col xs={12} md={6}>
            <Card size="small">
              <Statistic title="今日入住" value={reservations.length} />
            </Card>
          </Col>
          <Col xs={12} md={6}>
            <Card size="small">
              <Statistic title="今日退房" value={0} />
            </Card>
          </Col>
          <Col xs={12} md={6}>
            <Card size="small">
              <Statistic title="待确认" value={0} />
            </Card>
          </Col>
        </Row>

        <Card
          title="套餐管理"
          size="small"
          extra={
            <Button type="primary" icon={<PlusOutlined />} onClick={handleAddPkg}>
              添加套餐
            </Button>
          }
        >
          <Table columns={pkgColumns} dataSource={packages} rowKey="id" size="small" pagination={false} />
        </Card>

        <Card title="今日入住" size="small" className="mt-4">
          <Table columns={resColumns} dataSource={reservations} rowKey="id" size="small" pagination={false} />
        </Card>
      </Card>

      <Modal
        title={editingPkg ? '编辑套餐' : '添加套餐'}
        open={packageModal}
        onCancel={() => {
          setPackageModal(false)
          setEditingPkg(null)
          pkgForm.resetFields()
        }}
        onOk={() => pkgForm.submit()}
        confirmLoading={createPkgMutation.isPending}
      >
        <Form form={pkgForm} layout="vertical" onFinish={(v) => createPkgMutation.mutate(v)} className="mt-4">
          <Form.Item name="name" label="套餐名称" rules={[{ required: true }]}>
            <Input placeholder="套餐名称" />
          </Form.Item>
          <Form.Item name="type" label="类型" rules={[{ required: true }]}>
            <Select
              options={[
                { value: 'daycare', label: '日托' },
                { value: 'boarding', label: '寄养' },
              ]}
            />
          </Form.Item>
          <Form.Item name="price_per_day" label="每日价格" rules={[{ required: true }]}>
            <InputNumber className="w-full" min={0} step={1} />
          </Form.Item>
          <Form.Item name="capacity" label="容量" rules={[{ required: true }]}>
            <InputNumber className="w-full" min={1} step={1} />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input.TextArea rows={3} placeholder="套餐描述" />
          </Form.Item>
          <Form.Item name="features" label="特色">
            <Input.TextArea rows={2} placeholder="套餐特色" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
