import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { Row, Col, Card, Avatar, Button, Tag, Pagination, Form, Select, InputNumber, Space, message, Tabs } from 'antd'
import { HeartOutlined, HeartFilled, StopOutlined, UserOutlined, FilterOutlined } from '@ant-design/icons'
import { matchApi } from '@/api/endpoints'
import { MatchResultItem } from '@/api/endpoints'

export default function MatchPage() {
  const [page, setPage] = useState(1)
  const [pageSize] = useState(12)
  const [filterParams, setFilterParams] = useState<any>({})
  const queryClient = useQueryClient()

  const { data: smartData, isLoading: smartLoading } = useQuery({
    queryKey: ['smartMatch', page],
    queryFn: () => matchApi.smartMatch({ page, page_size: pageSize }),
  })

  const { data: filterData, isLoading: filterLoading } = useQuery({
    queryKey: ['filterMatch', page, filterParams],
    queryFn: () => matchApi.filterMatch({ ...filterParams, page, page_size: pageSize }),
    enabled: Object.keys(filterParams).length > 0,
  })

  const { data: favoritesData } = useQuery({
    queryKey: ['favorites'],
    queryFn: () => matchApi.getFavorites({ page: 1, page_size: 100 }),
  })

  const { data: blockedData } = useQuery({
    queryKey: ['blocked'],
    queryFn: () => matchApi.getBlocked({ page: 1, page_size: 100 }),
  })

  const favoriteMutation = useMutation({
    mutationFn: (id: number) => matchApi.favorite(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['smartMatch'] })
      queryClient.invalidateQueries({ queryKey: ['favorites'] })
      message.success('操作成功')
    },
  })

  const blockMutation = useMutation({
    mutationFn: (id: number) => matchApi.block(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['smartMatch'] })
      queryClient.invalidateQueries({ queryKey: ['blocked'] })
      message.success('操作成功')
    },
  })

  const handleFilter = (values: any) => {
    setFilterParams(values)
    setPage(1)
  }

  const renderMatchCard = (item: MatchResultItem) => (
    <Col xs={12} sm={8} md={6} lg={4} key={item.user_id}>
      <Card
        hoverable
        className="match-card"
        cover={
          <div style={{ height: 180, background: '#f5f5f5', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
            <Avatar size={80} src={item.profile?.photos?.[0]} icon={<UserOutlined />} />
          </div>
        }
        actions={[
          <Button
            type={item.is_favorited ? 'primary' : 'text'}
            icon={item.is_favorited ? <HeartFilled /> : <HeartOutlined />}
            onClick={() => favoriteMutation.mutate(item.user_id)}
          />,
          <Button
            type={item.is_blocked ? 'primary' : 'text'}
            danger
            icon={<StopOutlined />}
            onClick={() => blockMutation.mutate(item.user_id)}
          />,
        ]}
      >
        <Card.Meta
          title={item.profile?.nickname || `用户${item.user_id}`}
          description={
            <div>
              <Tag color="magenta">匹配度 {item.match_score}%</Tag>
              {item.match_reason && <p style={{ fontSize: 12, color: '#888', marginTop: 8 }}>{item.match_reason}</p>}
            </div>
          }
        />
      </Card>
    </Col>
  )

  return (
    <div>
      <Card style={{ marginBottom: 16 }}>
        <Form layout="inline" onFinish={handleFilter}>
          <Space wrap>
            <Form.Item name="gender">
              <Select placeholder="性别" allowClear style={{ width: 100 }}>
                <Select.Option value="male">男</Select.Option>
                <Select.Option value="female">女</Select.Option>
              </Select>
            </Form.Item>
            <Form.Item name="min_age">
              <InputNumber placeholder="最小年龄" min={18} max={60} style={{ width: 100 }} />
            </Form.Item>
            <Form.Item name="max_age">
              <InputNumber placeholder="最大年龄" min={18} max={60} style={{ width: 100 }} />
            </Form.Item>
            <Form.Item name="city">
              <Select placeholder="城市" allowClear showSearch style={{ width: 120 }}>
                <Select.Option value="北京">北京</Select.Option>
                <Select.Option value="上海">上海</Select.Option>
                <Select.Option value="广州">广州</Select.Option>
                <Select.Option value="深圳">深圳</Select.Option>
                <Select.Option value="杭州">杭州</Select.Option>
              </Select>
            </Form.Item>
            <Form.Item name="education">
              <Select placeholder="学历" allowClear style={{ width: 100 }}>
                <Select.Option value="high_school">高中</Select.Option>
                <Select.Option value="college">大专</Select.Option>
                <Select.Option value="bachelor">本科</Select.Option>
                <Select.Option value="master">硕士</Select.Option>
                <Select.Option value="phd">博士</Select.Option>
              </Select>
            </Form.Item>
            <Form.Item>
              <Button type="primary" htmlType="submit" icon={<FilterOutlined />}>
                筛选
              </Button>
            </Form.Item>
          </Space>
        </Form>
      </Card>

      <Tabs
        items={[
          {
            key: 'smart',
            label: '智能匹配',
            children: (
              <>
                <Row gutter={[16, 16]}>
                  {smartData?.list.map(renderMatchCard)}
                </Row>
                <div style={{ textAlign: 'center', marginTop: 24 }}>
                  <Pagination
                    current={page}
                    pageSize={pageSize}
                    total={smartData?.total || 0}
                    onChange={setPage}
                  />
                </div>
              </>
            ),
          },
          {
            key: 'favorites',
            label: `我的收藏 (${favoritesData?.total || 0})`,
            children: (
              <Row gutter={[16, 16]}>
                {favoritesData?.list.map(renderMatchCard)}
              </Row>
            ),
          },
          {
            key: 'blocked',
            label: `屏蔽列表 (${blockedData?.total || 0})`,
            children: (
              <Row gutter={[16, 16]}>
                {blockedData?.list.map(renderMatchCard)}
              </Row>
            ),
          },
        ]}
      />
    </div>
  )
}
