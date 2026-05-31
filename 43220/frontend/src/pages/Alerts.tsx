import { Card, Table, Tag, Button, message, Space } from 'antd'
import { CheckOutlined } from '@ant-design/icons'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { alertApi } from '@/services/api'
import { useAuthStore } from '@/context/AuthContext'
import type { HealthAlert } from '@/types'
import dayjs from 'dayjs'

export default function Alerts() {
  const queryClient = useQueryClient()
  const { user } = useAuthStore()

  const { data, isLoading } = useQuery({
    queryKey: ['alerts', 'list'],
    queryFn: () => alertApi.list({ page_size: 100 }),
  })

  const alerts: HealthAlert[] = data?.data?.items || []
  const total = data?.data?.total || 0

  const markReadMutation = useMutation({
    mutationFn: (id: string) => alertApi.markAsRead(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['alerts'] })
      message.success('已标记为已读')
    },
    onError: (err: any) => message.error(err.message || '操作失败'),
  })

  const markAllReadMutation = useMutation({
    mutationFn: () => alertApi.markAllAsRead(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['alerts'] })
      message.success('全部标记为已读')
    },
    onError: (err: any) => message.error(err.message || '操作失败'),
  })

  const columns = [
    {
      title: '标题',
      dataIndex: 'title',
      key: 'title',
      render: (title: string, record: HealthAlert) => (
        <span className={record.is_read ? 'text-gray-400' : 'font-medium'}>{title}</span>
      ),
    },
    {
      title: '类型',
      dataIndex: 'alert_type',
      key: 'alert_type',
      render: (type: string) => {
        const map: Record<string, { color: string; text: string }> = {
          vaccine_expire: { color: 'orange', text: '疫苗到期' },
          deworm_expire: { color: 'cyan', text: '驱虫到期' },
        }
        const t = map[type] || { color: 'default', text: type }
        return <Tag color={t.color}>{t.text}</Tag>
      },
    },
    { title: '内容', dataIndex: 'content', key: 'content', ellipsis: true },
    {
      title: '到期时间',
      dataIndex: 'expire_at',
      key: 'expire_at',
      render: (d: string) => dayjs(d).format('YYYY-MM-DD'),
    },
    {
      title: '状态',
      dataIndex: 'is_read',
      key: 'is_read',
      render: (read: boolean) => (
        <Tag color={read ? 'default' : 'red'}>{read ? '已读' : '未读'}</Tag>
      ),
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (d: string) => dayjs(d).format('YYYY-MM-DD HH:mm'),
    },
    {
      title: '操作',
      key: 'actions',
      render: (_: any, alert: HealthAlert) =>
      !alert.is_read ? (
        <Button size="small" icon={<CheckOutlined />} onClick={() => markReadMutation.mutate(alert.id)}>
          标记已读
        </Button>
      ) : null,
    },
  ]

  return (
    <div className="space-y-4">
      <Card
        title="健康提醒"
        extra={
          <Button type="primary" onClick={() => markAllReadMutation.mutate()}>
            全部已读
          </Button>
        }
      >
        <Table columns={columns} dataSource={alerts} rowKey="id" loading={isLoading} pagination={{ pageSize: 10, total }} />
      </Card>
    </div>
  )
}
