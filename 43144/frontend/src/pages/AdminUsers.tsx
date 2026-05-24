import React, { useEffect, useState } from 'react'
import { Table, Tag, Card, Button, Modal, message, Space, Select, Empty, Spin } from 'antd'
import { CheckOutlined } from '@ant-design/icons'
import { getUsers, verifyUser } from '../api/auth'
import { User } from '../types'

const { Option } = Select

const AdminUsers: React.FC = () => {
  const [users, setUsers] = useState<User[]>([])
  const [loading, setLoading] = useState(false)
  const [roleFilter, setRoleFilter] = useState<string>('')

  useEffect(() => {
    loadUsers()
  }, [roleFilter])

  const loadUsers = async () => {
    setLoading(true)
    try {
      const response = await getUsers({
        role: roleFilter || undefined,
        page_size: 100,
      })
      if (response.code === 0 && response.data) {
        setUsers((response.data as any).items || [])
      }
    } catch (error) {
      console.error('Failed to load users:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleVerify = async (id: number) => {
    Modal.confirm({
      title: '确认认证',
      content: '确定要认证该用户吗？',
      onOk: async () => {
        try {
          await verifyUser(id)
          message.success('已认证')
          loadUsers()
        } catch (error: any) {
          message.error('认证失败')
        }
      },
    })
  }

  const roleText: Record<string, string> = {
    admin: '管理员',
    rescue: '救助站',
    adopter: '领养人',
  }

  const columns = [
    { title: '编号', dataIndex: 'id', key: 'id' },
    { title: '姓名', dataIndex: 'name', key: 'name' },
    { title: '邮箱', dataIndex: 'email', key: 'email' },
    { title: '手机号', dataIndex: 'phone', key: 'phone', render: (p: string) => p || '-' },
    {
      title: '角色',
      dataIndex: 'role',
      key: 'role',
      render: (role: string) => <Tag>{roleText[role] || role}</Tag>,
    },
    {
      title: '认证状态',
      dataIndex: 'is_verified',
      key: 'verified',
      render: (v: boolean) => (
        v ? <Tag color="green">已认证</Tag> : <Tag color="orange">未认证</Tag>
      ),
    },
    {
      title: '注册时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date: string) => new Date(date).toLocaleDateString(),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: User) => (
        !record.is_verified && record.role !== 'admin' && (
          <Button
            type="link"
            size="small"
            icon={<CheckOutlined />}
            onClick={() => handleVerify(record.id)}
          >
            认证
          </Button>
        )
      ),
    },
  ]

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h2 style={{ margin: 0 }}>用户管理</h2>
        <Select
          placeholder="筛选角色"
          allowClear
          style={{ width: 150 }}
          value={roleFilter || undefined}
          onChange={setRoleFilter}
        >
          <Option value="adopter">领养人</Option>
          <Option value="rescue">救助站</Option>
          <Option value="admin">管理员</Option>
        </Select>
      </div>

      <Spin spinning={loading}>
        {users.length > 0 ? (
          <Card>
            <Table
              dataSource={users}
              columns={columns}
              rowKey="id"
              pagination={{ pageSize: 10 }}
            />
          </Card>
        ) : (
          !loading && (
            <Card>
              <Empty description="暂无用户" />
            </Card>
          )
        )}
      </Spin>
    </div>
  )
}

export default AdminUsers
