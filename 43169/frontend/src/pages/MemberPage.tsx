import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { Row, Col, Card, Button, Table, Tag, Modal, Select, message, Checkbox } from 'antd'
import { CrownOutlined, CheckOutlined } from '@ant-design/icons'
import { memberApi, MemberBenefit, MemberOrder } from '@/api/endpoints'
import { useAuthStore } from '@/store/authStore'

const memberColors: Record<string, string> = {
  free: '#bfbfbf',
  silver: '#c0c0c0',
  gold: '#ffd700',
  diamond: '#b9f2ff',
}

const memberNames: Record<string, string> = {
  free: '免费会员',
  silver: '白银会员',
  gold: '黄金会员',
  diamond: '钻石会员',
}

const benefitNames: Record<string, string> = {
  unlimited_chat: '无限聊天',
  view_who_liked: '查看谁喜欢我',
  priority_match: '优先匹配',
  advanced_filter: '高级筛选',
  video_chat: '视频聊天',
  hide_online: '隐藏在线状态',
  no_ads: '无广告',
  matchmaker_assist: '红娘服务',
}

export default function MemberPage() {
  const { user } = useAuthStore()
  const [orderModal, setOrderModal] = useState(false)
  const [selectedBenefit, setSelectedBenefit] = useState<MemberBenefit | null>(null)
  const [months, setMonths] = useState(1)
  const queryClient = useQueryClient()

  const { data: benefitsData } = useQuery({
    queryKey: ['benefits'],
    queryFn: memberApi.getBenefits,
  })

  const { data: ordersData } = useQuery({
    queryKey: ['orders'],
    queryFn: () => memberApi.getOrders({ page: 1, page_size: 20 }),
  })

  const orderMutation = useMutation({
    mutationFn: (params: { level: string; months: number }) =>
      memberApi.createOrder(params),
    onSuccess: (data: any) => {
      const order = data?.data || data
      if (order?.id) {
        memberApi.payOrder(order.id).then(() => {
          message.success('支付成功！会员已升级')
          queryClient.invalidateQueries({ queryKey: ['orders'] })
          queryClient.invalidateQueries({ queryKey: ['userInfo'] })
        })
      }
      setOrderModal(false)
    },
  })

  const benefits = benefitsData || []

  const columns = [
    { title: '会员等级', dataIndex: 'level', key: 'level', render: (level: string) => memberNames[level] },
    { title: '时长', dataIndex: 'months', key: 'months', render: (m: number) => `${m}个月` },
    { title: '金额', dataIndex: 'amount', key: 'amount', render: (a: number) => `¥${a}` },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={status === 'paid' ? 'green' : status === 'pending' ? 'orange' : 'red'}>
          {status === 'paid' ? '已支付' : status === 'pending' ? '待支付' : '已取消'}
        </Tag>
      ),
    },
    {
      title: '有效期至',
      dataIndex: 'expire_at',
      key: 'expire_at',
      render: (v: string) => v || '-',
    },
    { title: '创建时间', dataIndex: 'created_at', key: 'created_at' },
  ]

  return (
    <div>
      <Card style={{ marginBottom: 24 }}>
        <h2 style={{ marginBottom: 16 }}>当前会员状态</h2>
        <Row gutter={16}>
          <Col span={8}>
            <Tag color={memberColors[user?.member_level || 'free']} style={{ fontSize: 16, padding: '8px 16px' }}>
              <CrownOutlined /> {memberNames[user?.member_level || 'free']}
            </Tag>
          </Col>
          <Col span={8}>
            <p style={{ color: '#888' }}>有效期至：{user?.member_expire || '永久（免费会员）'}</p>
          </Col>
        </Row>
      </Card>

      <Card title="会员权益" style={{ marginBottom: 24 }}>
        <Row gutter={[16, 16]}>
          {benefits.map((benefit) => (
            <Col xs={24} sm={12} lg={6} key={benefit.id}>
              <Card
                className={`member-benefit-card ${benefit.level === 'gold' ? 'featured' : ''}`}
                style={{
                  borderColor: memberColors[benefit.level],
                  background: `linear-gradient(135deg, ${memberColors[benefit.level]}20 0%, #ffffff 100%)`,
                }}
              >
                <CrownOutlined style={{ fontSize: 48, color: memberColors[benefit.level] }} />
                <h3 style={{ margin: '16px 0' }}>{memberNames[benefit.level]}</h3>
                <p style={{ fontSize: 24, fontWeight: 'bold', color: memberColors[benefit.level] }}>
                  ¥{benefit.price_per_month}/月
                </p>
                <p style={{ color: '#888', marginBottom: 16 }}>{benefit.description}</p>
                <div style={{ textAlign: 'left', marginBottom: 16 }}>
                  <p><Checkbox checked disabled /> 每天{benefit.daily_interact}次互动</p>
                  {Object.entries(benefitNames).map(([key, name]) => (
                    (benefit as any)[key] && (
                      <p key={key}><Checkbox checked disabled /> {name}</p>
                    )
                  ))}
                </div>
                {benefit.level !== 'free' && (
                  <Button
                    type="primary"
                    block
                    onClick={() => {
                      setSelectedBenefit(benefit)
                      setMonths(1)
                      setOrderModal(true)
                    }}
                  >
                    立即开通
                  </Button>
                )}
              </Card>
            </Col>
          ))}
        </Row>
      </Card>

      <Card title="我的订单">
        <Table
          columns={columns}
          dataSource={ordersData?.list || []}
          rowKey="id"
          pagination={false}
        />
      </Card>

      <Modal
        title="开通会员"
        open={orderModal}
        onCancel={() => setOrderModal(false)}
        footer={null}
      >
        {selectedBenefit && (
          <div>
            <p><strong>会员等级：</strong>{memberNames[selectedBenefit.level]}</p>
            <p><strong>单价：</strong>¥{selectedBenefit.price_per_month}/月</p>
            <p><strong>选择时长：</strong></p>
            <Select value={months} onChange={setMonths} style={{ width: '100%', marginBottom: 16 }}>
              {[1, 3, 6, 12].map((m) => (
                <Select.Option key={m} value={m}>
                  {m}个月 - ¥{(selectedBenefit.price_per_month * m).toFixed(2)}
                </Select.Option>
              ))}
            </Select>
            <p><strong>总计：</strong>¥{(selectedBenefit.price_per_month * months).toFixed(2)}</p>
            <Button
              type="primary"
              block
              icon={<CheckOutlined />}
              loading={orderMutation.isPending}
              onClick={() => orderMutation.mutate({ level: selectedBenefit.level, months })}
            >
              确认支付
            </Button>
          </div>
        )}
      </Modal>
    </div>
  )
}
