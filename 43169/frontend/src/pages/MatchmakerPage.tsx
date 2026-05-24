import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { Card, Table, Button, Modal, Form, Input, InputNumber, Tabs, Tag, Statistic, Row, Col, message } from 'antd'
import { PlusOutlined, UserAddOutlined, TeamOutlined, StarOutlined } from '@ant-design/icons'
import { matchmakerApi, MatchmakerMember, MatchmakerService, MatchmakerStats } from '@/api/endpoints'
import { useAuthStore } from '@/store/authStore'

export default function MatchmakerPage() {
  const { user } = useAuthStore()
  const isMatchmaker = user?.role === 'matchmaker' || user?.role === 'admin'
  const [memberModal, setMemberModal] = useState(false)
  const [serviceModal, setServiceModal] = useState(false)
  const [memberForm] = Form.useForm()
  const [serviceForm] = Form.useForm()
  const queryClient = useQueryClient()

  const { data: statsData } = useQuery({
    queryKey: ['matchmakerStats'],
    queryFn: matchmakerApi.getStats,
    enabled: isMatchmaker,
  })

  const { data: membersData } = useQuery({
    queryKey: ['matchmakerMembers'],
    queryFn: () => matchmakerApi.listMembers({ page: 1, page_size: 100 }),
    enabled: isMatchmaker,
  })

  const { data: servicesData } = useQuery({
    queryKey: ['matchmakerServices'],
    queryFn: () => matchmakerApi.listServices({ page: 1, page_size: 100 }),
    enabled: isMatchmaker,
  })

  const { data: allMatchmakers } = useQuery({
    queryKey: ['allMatchmakers'],
    queryFn: () => matchmakerApi.listAll({ page: 1, page_size: 50 }),
  })

  const addMemberMutation = useMutation({
    mutationFn: matchmakerApi.addMember,
    onSuccess: () => {
      message.success('添加成功')
      setMemberModal(false)
      memberForm.resetFields()
      queryClient.invalidateQueries({ queryKey: ['matchmakerMembers'] })
      queryClient.invalidateQueries({ queryKey: ['matchmakerStats'] })
    },
  })

  const removeMemberMutation = useMutation({
    mutationFn: matchmakerApi.removeMember,
    onSuccess: () => {
      message.success('已移除')
      queryClient.invalidateQueries({ queryKey: ['matchmakerMembers'] })
      queryClient.invalidateQueries({ queryKey: ['matchmakerStats'] })
    },
  })

  const createServiceMutation = useMutation({
    mutationFn: matchmakerApi.createService,
    onSuccess: () => {
      message.success('服务创建成功')
      setServiceModal(false)
      serviceForm.resetFields()
      queryClient.invalidateQueries({ queryKey: ['matchmakerServices'] })
    },
  })

  const updateProgressMutation = useMutation({
    mutationFn: ({ id, progress }: { id: number; progress: number }) =>
      matchmakerApi.updateProgress(id, progress),
    onSuccess: () => {
      message.success('进度已更新')
      queryClient.invalidateQueries({ queryKey: ['matchmakerServices'] })
    },
  })

  const stats = statsData as MatchmakerStats | undefined

  if (!isMatchmaker) {
    return (
      <Card title="红娘列表">
        <Table
          dataSource={allMatchmakers?.list || []}
          rowKey="id"
          columns={[
            { title: '用户名', dataIndex: 'username', key: 'username' },
            { title: '头像', dataIndex: 'avatar', key: 'avatar', render: (v: string) => v ? <img src={v} alt="" style={{ width: 40, height: 40, borderRadius: '50%' }} /> : '-' },
            { title: '角色', dataIndex: 'role', key: 'role' },
            { title: '认证状态', dataIndex: 'verify_status', key: 'verify_status', render: (v: string) => <Tag color={v === 'verified' ? 'green' : 'orange'}>{v}</Tag> },
          ]}
          pagination={false}
        />
      </Card>
    )
  }

  const memberColumns = [
    { title: '会员ID', dataIndex: 'member_id', key: 'member_id' },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => <Tag color={status === 'active' ? 'green' : 'default'}>{status}</Tag>,
    },
    { title: '加入时间', dataIndex: 'joined_at', key: 'joined_at' },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: MatchmakerMember) => (
        <Button danger size="small" onClick={() => removeMemberMutation.mutate(record.member_id)}>
          移除
        </Button>
      ),
    },
  ]

  const serviceColumns = [
    { title: '会员A', dataIndex: 'member_a_id', key: 'member_a_id' },
    { title: '会员B', dataIndex: 'member_b_id', key: 'member_b_id' },
    { title: '服务类型', dataIndex: 'service_type', key: 'service_type' },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => <Tag color={status === 'completed' ? 'green' : 'blue'}>{status}</Tag>,
    },
    {
      title: '进度',
      dataIndex: 'progress',
      key: 'progress',
      render: (progress: number, record: MatchmakerService) => (
        <InputNumber
          size="small"
          min={0}
          max={100}
          value={progress}
          onChange={(v) => v && updateProgressMutation.mutate({ id: record.id, progress: v })}
          suffix="%"
        />
      ),
    },
    { title: '备注', dataIndex: 'note', key: 'note' },
  ]

  return (
    <div>
      <Card style={{ marginBottom: 24 }}>
        <h2 style={{ marginBottom: 16 }}>业绩统计</h2>
        <Row gutter={16}>
          <Col span={6}>
            <Statistic title="会员总数" value={stats?.total_members || 0} prefix={<TeamOutlined />} />
          </Col>
          <Col span={6}>
            <Statistic title="服务总数" value={stats?.total_services || 0} />
          </Col>
          <Col span={6}>
            <Statistic title="约会总数" value={stats?.total_dates || 0} />
          </Col>
          <Col span={6}>
            <Statistic title="成功约会" value={stats?.success_dates || 0} />
          </Col>
          <Col span={6}>
            <Statistic title="平均评分" value={stats?.avg_rating || 0} prefix={<StarOutlined />} precision={1} />
          </Col>
        </Row>
      </Card>

      <Card>
        <Tabs
          items={[
            {
              key: 'members',
              label: '我的会员',
              children: (
                <>
                  <Button
                    type="primary"
                    icon={<UserAddOutlined />}
                    style={{ marginBottom: 16 }}
                    onClick={() => setMemberModal(true)}
                  >
                    添加会员
                  </Button>
                  <Table
                    columns={memberColumns}
                    dataSource={membersData?.list || []}
                    rowKey="id"
                    pagination={false}
                  />
                </>
              ),
            },
            {
              key: 'services',
              label: '服务记录',
              children: (
                <>
                  <Button
                    type="primary"
                    icon={<PlusOutlined />}
                    style={{ marginBottom: 16 }}
                    onClick={() => setServiceModal(true)}
                  >
                    创建服务
                  </Button>
                  <Table
                    columns={serviceColumns}
                    dataSource={servicesData?.list || []}
                    rowKey="id"
                    pagination={false}
                  />
                </>
              ),
            },
          ]}
        />
      </Card>

      <Modal
        title="添加会员"
        open={memberModal}
        onCancel={() => setMemberModal(false)}
        footer={null}
      >
        <Form form={memberForm} onFinish={(v) => addMemberMutation.mutate(v)} layout="vertical">
          <Form.Item name="member_id" label="会员ID" rules={[{ required: true }]}>
            <InputNumber style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={addMemberMutation.isPending} block>
              添加
            </Button>
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title="创建服务"
        open={serviceModal}
        onCancel={() => setServiceModal(false)}
        footer={null}
      >
        <Form form={serviceForm} onFinish={(v) => createServiceMutation.mutate(v)} layout="vertical">
          <Form.Item name="member_a_id" label="会员A ID" rules={[{ required: true }]}>
            <InputNumber style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="member_b_id" label="会员B ID" rules={[{ required: true }]}>
            <InputNumber style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="service_type" label="服务类型" rules={[{ required: true }]}>
            <Input placeholder="如：初次牵线、约会安排等" />
          </Form.Item>
          <Form.Item name="note" label="备注">
            <Input.TextArea rows={3} />
          </Form.Item>
          <Form.Item name="date_id" label="关联约会ID（选填）">
            <InputNumber style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={createServiceMutation.isPending} block>
              创建
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
