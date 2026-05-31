import { useState } from 'react'
import { Card, Table, Tag, Button, Select, Input, Space, Modal, message } from 'antd'
import { SearchOutlined, UserOutlined } from '@ant-design/icons'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { userApi } from '@/services/auth'
import { User, UserRole, VerificationStatus } from '@/types'

export function UserManagement() {
  const queryClient = useQueryClient()
  const [page, setPage] = useState(1)
  const [pageSize] = useState(10)
  const [role, setRole] = useState<string>('')
  const [searchText, setSearchText] = useState('')

  const { data, isLoading } = useQuery({
    queryKey: ['users', page, pageSize, role],
    queryFn: () => userApi.getUsers({ page, page_size: pageSize, role: role || undefined }),
  })

  const filteredUsers = (data?.items || []).filter(
    (user: User) =>
      !searchText ||
      user.username.toLowerCase().includes(searchText.toLowerCase()) ||
      user.email.toLowerCase().includes(searchText.toLowerCase()) ||
      user.full_name.toLowerCase().includes(searchText.toLowerCase())
  )

  const columns = [
    {
      title: '用户信息',
      key: 'user',
      render: (_: any, record: User) => (
        <Space>
          <UserOutlined />
          <div>
            <div style={{ fontWeight: 500 }}>{record.full_name}</div>
            <div style={{ color: '#999', fontSize: 12 }}>@{record.username}</div>
          </div>
        </Space>
      ),
    },
    {
      title: '邮箱',
      dataIndex: 'email',
      key: 'email',
    },
    {
      title: '角色',
      dataIndex: 'role',
      key: 'role',
      render: (role: UserRole) => {
        const colorMap: Record<string, string> = {
          client: 'blue',
          professional: 'green',
          admin: 'purple',
        }
        const textMap: Record<string, string> = {
          client: '客户',
          professional: '专业人士',
          admin: '管理员',
        }
        return <Tag color={colorMap[role]}>{textMap[role]}</Tag>
      },
    },
    {
      title: '认证状态',
      dataIndex: 'verification_status',
      key: 'verification_status',
      render: (status: VerificationStatus) => {
        const colorMap: Record<string, string> = {
          pending: 'orange',
          approved: 'green',
          rejected: 'red',
        }
        const textMap: Record<string, string> = {
          pending: '待审核',
          approved: '已通过',
          rejected: '已拒绝',
        }
        return <Tag color={colorMap[status]}>{textMap[status]}</Tag>
      },
    },
    {
      title: '状态',
      dataIndex: 'is_active',
      key: 'is_active',
      render: (isActive: boolean) => (
        <Tag color={isActive ? 'green' : 'red'}>
          {isActive ? '活跃' : '已禁用'}
        </Tag>
      ),
    },
    {
      title: '注册时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date: string) => new Date(date).toLocaleDateString(),
    },
  ]

  return (
    <div className="page-container">
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 24 }}>
        <h2>用户管理</h2>
        <Space>
          <Input
            placeholder="搜索用户"
            prefix={<SearchOutlined />}
            value={searchText}
            onChange={(e) => setSearchText(e.target.value)}
            style={{ width: 200 }}
            allowClear
          />
          <Select
            style={{ width: 140 }}
            value={role}
            onChange={setRole}
            allowClear
            placeholder="筛选角色"
            options={[
              { value: '', label: '全部' },
              { value: 'client', label: '客户' },
              { value: 'professional', label: '专业人士' },
              { value: 'admin', label: '管理员' },
            ]}
          />
        </Space>
      </div>

      <Card>
        <Table
          columns={columns}
          dataSource={filteredUsers}
          loading={isLoading}
          rowKey="id"
          pagination={{
            current: page,
            pageSize,
            total: data?.total || 0,
            onChange: setPage,
          }}
        />
      </Card>
    </div>
  )
}
