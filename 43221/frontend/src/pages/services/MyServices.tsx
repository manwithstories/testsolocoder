import { useState } from 'react'
import { Card, List, Button, Tag, Space, Modal, message, Popconfirm } from 'antd'
import { EditOutlined, DeleteOutlined, PlusOutlined, CalendarOutlined } from '@ant-design/icons'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useNavigate } from 'react-router-dom'
import { useAuthContext } from '@/contexts/AuthContext'
import { serviceApi } from '@/services/service'
import { Service } from '@/types'

export function MyServices() {
  const navigate = useNavigate()
  const { user } = useAuthContext()
  const queryClient = useQueryClient()
  const [page, setPage] = useState(1)
  const [pageSize] = useState(10)

  const { data, isLoading } = useQuery({
    queryKey: ['my-services', page, pageSize],
    queryFn: () => serviceApi.getProfessionalServices({ page, page_size: pageSize }),
    enabled: !!user,
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => serviceApi.delete(id),
    onSuccess: () => {
      message.success('删除成功')
      queryClient.invalidateQueries({ queryKey: ['my-services'] })
    },
    onError: (error: any) => {
      message.error(error.message || '删除失败')
    },
  })

  const updateStatusMutation = useMutation({
    mutationFn: ({ id, status }: { id: string; status: string }) =>
      serviceApi.update(id, { status }),
    onSuccess: () => {
      message.success('状态更新成功')
      queryClient.invalidateQueries({ queryKey: ['my-services'] })
    },
    onError: (error: any) => {
      message.error(error.message || '状态更新失败')
    },
  })

  return (
    <div className="page-container">
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 24 }}>
        <h2>我的服务</h2>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => navigate('/services/create')}>
          发布服务
        </Button>
      </div>

      <Card>
        <List
          dataSource={data?.items || []}
          loading={isLoading}
          renderItem={(service: Service) => (
            <List.Item
              actions={[
                <Button
                  key="edit"
                  type="link"
                  icon={<CalendarOutlined />}
                  onClick={() => navigate(`/services/${service.id}/schedules`)}
                >
                  时段管理
                </Button>,
                <Button
                  key="toggle"
                  type="link"
                  onClick={() =>
                    updateStatusMutation.mutate({
                      id: service.id,
                      status: service.status === 'active' ? 'inactive' : 'active',
                    })
                  }
                >
                  {service.status === 'active' ? '下架' : '上架'}
                </Button>,
                <Popconfirm
                  key="delete"
                  title="确定要删除这个服务吗？"
                  onConfirm={() => deleteMutation.mutate(service.id)}
                >
                  <Button type="link" danger icon={<DeleteOutlined />}>
                    删除
                  </Button>
                </Popconfirm>,
              ]}
            >
              <List.Item.Meta
                title={
                  <Space>
                    <span style={{ fontSize: 16, fontWeight: 500 }}>{service.title}</span>
                    <Tag color={service.status === 'active' ? 'green' : 'default'}>
                      {service.status === 'active' ? '已上架' : '已下架'}
                    </Tag>
                    <Tag color="blue">{service.service_type}</Tag>
                  </Space>
                }
                description={
                  <div>
                    <div style={{ color: '#666', marginBottom: 4 }}>{service.description}</div>
                    <div style={{ display: 'flex', gap: 16 }}>
                      <span>价格：¥{service.price}</span>
                      <span>时长：{service.duration_minutes}分钟</span>
                      <span>评分：{service.average_rating.toFixed(1)}</span>
                      <span>评价数：{service.review_count}</span>
                    </div>
                  </div>
                }
              />
            </List.Item>
          )}
        />
      </Card>
    </div>
  )
}
