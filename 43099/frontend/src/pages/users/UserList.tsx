import { useState, useEffect } from 'react'
import { Table, Tag, Select, message, Card, Modal, Form, Input, Button } from 'antd'
import { EditOutlined } from '@ant-design/icons'
import { userApi } from '@/api'
import type { User } from '@/types'
import dayjs from 'dayjs'

const UserList = () => {
  const [users, setUsers] = useState<User[]>([])
  const [loading, setLoading] = useState(false)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [editModalVisible, setEditModalVisible] = useState(false)
  const [editingUser, setEditingUser] = useState<User | null>(null)
  const [form] = Form.useForm()

  useEffect(() => {
    loadUsers()
  }, [page, pageSize])

  const loadUsers = async () => {
    setLoading(true)
    try {
      const data: any = await userApi.list({ page, page_size: pageSize })
      setUsers(data.list)
      setTotal(data.total)
    } catch (error: any) {
      message.error(error.message || '加载失败')
    } finally {
      setLoading(false)
    }
  }

  const handleEditRole = (user: User) => {
    setEditingUser(user)
    form.setFieldsValue({ role: user.role })
    setEditModalVisible(true)
  }

  const handleSubmitRole = async (values: any) => {
    if (!editingUser) return
    try {
      await userApi.updateRole(editingUser.id, values.role)
      message.success('角色更新成功')
      setEditModalVisible(false)
      loadUsers()
    } catch (error: any) {
      message.error(error.message || '更新失败')
    }
  }

  const getRoleTag = (role: string) => {
    const map: Record<string, { text: string; color: string }> = {
      user: { text: '普通用户', color: 'blue' },
      admin: { text: '管理员', color: 'orange' },
      super_admin: { text: '超级管理员', color: 'red' },
    }
    const info = map[role] || { text: role, color: 'default' }
    return <Tag color={info.color}>{info.text}</Tag>
  }

  const getStatusTag = (status: string) => {
    const color = status === 'active' ? 'green' : 'red'
    const text = status === 'active' ? '正常' : '禁用'
    return <Tag color={color}>{text}</Tag>
  }

  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 60,
    },
    {
      title: '用户名',
      dataIndex: 'username',
      key: 'username',
    },
    {
      title: '邮箱',
      dataIndex: 'email',
      key: 'email',
    },
    {
      title: '真实姓名',
      dataIndex: 'real_name',
      key: 'real_name',
    },
    {
      title: '电话',
      dataIndex: 'phone',
      key: 'phone',
    },
    {
      title: '角色',
      dataIndex: 'role',
      key: 'role',
      render: getRoleTag,
    },
    {
      title: '邮箱验证',
      dataIndex: 'email_verified',
      key: 'email_verified',
      render: (verified: boolean) => (
        <Tag color={verified ? 'green' : 'orange'}>
          {verified ? '已验证' : '未验证'}
        </Tag>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: getStatusTag,
    },
    {
      title: '注册时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (val: string) => dayjs(val).format('YYYY-MM-DD HH:mm'),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: User) => (
        <Button type="link" icon={<EditOutlined />} onClick={() => handleEditRole(record)}>
          修改角色
        </Button>
      ),
    },
  ]

  return (
    <Card title="用户管理">
      <Table
        columns={columns}
        dataSource={users}
        rowKey="id"
        loading={loading}
        pagination={{
          current: page,
          pageSize,
          total,
          onChange: (p, ps) => {
            setPage(p)
            setPageSize(ps)
          },
        }}
      />

      <Modal
        title="修改用户角色"
        open={editModalVisible}
        onCancel={() => setEditModalVisible(false)}
        footer={null}
      >
        <Form form={form} layout="vertical" onFinish={handleSubmitRole}>
          <Form.Item label="用户">
            <Input value={editingUser?.username} disabled />
          </Form.Item>
          <Form.Item
            name="role"
            label="角色"
            rules={[{ required: true, message: '请选择角色' }]}
          >
            <Select placeholder="请选择角色">
              <Select.Option value="user">普通用户</Select.Option>
              <Select.Option value="admin">管理员</Select.Option>
              <Select.Option value="super_admin">超级管理员</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block>
              保存
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </Card>
  )
}

export default UserList
