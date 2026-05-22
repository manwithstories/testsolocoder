import React, { useEffect, useState } from 'react'
import { Table, Tag, Button, Space, Select, Input, Modal, Form, message, Popconfirm } from 'antd'
import { userApi } from '@/services'
import { User, PaginatedData } from '@/types'

const roleMap: Record<string, string> = {
  student: '学员',
  instructor: '讲师',
  admin: '管理员',
}

const AdminUsers: React.FC = () => {
  const [data, setData] = useState<PaginatedData<User>>({
    items: [], total: 0, page: 1, page_size: 10, total_pages: 0,
  })
  const [loading, setLoading] = useState(false)
  const [roleFilter, setRoleFilter] = useState('')
  const [statusFilter, setStatusFilter] = useState('')
  const [search, setSearch] = useState('')

  const loadUsers = async (page = 1) => {
    setLoading(true)
    try {
      const params: any = { page, page_size: 10 }
      if (roleFilter) params.role = roleFilter
      if (statusFilter) params.status = statusFilter
      if (search) params.search = search
      const res = await userApi.list(params)
      if (res.code === 0 && res.data) {
        setData(res.data)
      }
    } catch (error) {
      console.error('Failed to load users:', error)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadUsers()
  }, [])

  const handleUpdateStatus = async (id: string, status: string) => {
    try {
      const res = await userApi.updateStatus(id, status)
      if (res.code === 0) {
        message.success('状态更新成功')
        loadUsers()
      }
    } catch (error: any) {
      message.error(error.message || '操作失败')
    }
  }

  const columns = [
    {
      title: '头像',
      dataIndex: 'avatar',
      key: 'avatar',
      render: (val: string) => val ? <img src={val} alt="" style={{ width: 32, height: 32, borderRadius: '50%' }} /> : <div style={{ width: 32, height: 32, borderRadius: '50%', background: '#e8e8e8' }} />,
    },
    { title: '用户名', dataIndex: 'username', key: 'username' },
    { title: '邮箱', dataIndex: 'email', key: 'email' },
    { title: '昵称', dataIndex: 'nickname', key: 'nickname' },
    {
      title: '角色',
      dataIndex: 'role',
      key: 'role',
      render: (role: string) => <Tag color={role === 'admin' ? 'red' : role === 'instructor' ? 'blue' : 'green'}>{roleMap[role]}</Tag>,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={status === 'active' ? 'green' : 'red'}>
          {status === 'active' ? '正常' : '禁用'}
        </Tag>
      ),
    },
    {
      title: '注册时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (val: string) => new Date(val).toLocaleDateString(),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: User) => (
        <Space>
          {record.status === 'active' ? (
            <Popconfirm title="确定禁用该用户？" onConfirm={() => handleUpdateStatus(record.id, 'disabled')}>
              <Button type="link" danger>禁用</Button>
            </Popconfirm>
          ) : (
            <Popconfirm title="确定启用该用户？" onConfirm={() => handleUpdateStatus(record.id, 'active')}>
              <Button type="link">启用</Button>
            </Popconfirm>
          )}
        </Space>
      ),
    },
  ]

  return (
    <div>
      <h2>用户管理</h2>
      <Space style={{ marginBottom: 16 }}>
        <Select
          placeholder="全部角色"
          style={{ width: 120 }}
          allowClear
          value={roleFilter || undefined}
          onChange={(val) => { setRoleFilter(val || ''); loadUsers(1) }}
        >
          <Select.Option value="student">学员</Select.Option>
          <Select.Option value="instructor">讲师</Select.Option>
          <Select.Option value="admin">管理员</Select.Option>
        </Select>
        <Select
          placeholder="全部状态"
          style={{ width: 120 }}
          allowClear
          value={statusFilter || undefined}
          onChange={(val) => { setStatusFilter(val || ''); loadUsers(1) }}
        >
          <Select.Option value="active">正常</Select.Option>
          <Select.Option value="disabled">禁用</Select.Option>
        </Select>
        <Input.Search
          placeholder="搜索用户"
          style={{ width: 200 }}
          onSearch={(val) => { setSearch(val); loadUsers(1) }}
          allowClear
        />
      </Space>
      <Table
        rowKey="id"
        columns={columns}
        dataSource={data.items}
        loading={loading}
        pagination={{
          current: data.page,
          total: data.total,
          pageSize: data.page_size,
          onChange: (page) => loadUsers(page),
        }}
      />
    </div>
  )
}

export default AdminUsers
