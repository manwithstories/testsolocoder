import { useState, useEffect } from 'react'
import {
  Table, Button, Card, Typography, Tag, Space, Pagination,
  message, Input, Select
} from 'antd'
import { SearchOutlined } from '@ant-design/icons'
import { adminAPI } from '@/api'
import { User } from '@/types'

const { Title } = Typography
const { Option } = Select

const AdminUsers = () => {
  const [users, setUsers] = useState<User[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize] = useState(10)
  const [keyword, setKeyword] = useState('')
  const [role, setRole] = useState<string>('')

  useEffect(() => {
    loadUsers()
  }, [page, keyword, role])

  const loadUsers = async () => {
    try {
      const params: any = { page, pageSize }
      if (keyword) params.keyword = keyword
      if (role) params.role = role
      const res = await adminAPI.getUsers(params)
      setUsers(res.data.data)
      setTotal(res.data.pagination.total)
    } catch (err) {
      console.error('加载用户失败', err)
    }
  }

  const getRoleTag = (role: string) => {
    const roleMap: Record<string, { color: string; text: string }> = {
      buyer: { color: 'blue', text: '买家' },
      seller: { color: 'green', text: '商家' },
      admin: { color: 'gold', text: '管理员' },
    }
    const info = roleMap[role] || { color: 'default', text: role }
    return <Tag color={info.color}>{info.text}</Tag>
  }

  const getStatusTag = (status: string) => {
    const statusMap: Record<string, { color: string; text: string }> = {
      active: { color: 'green', text: '正常' },
      banned: { color: 'red', text: '已封禁' },
    }
    const info = statusMap[status] || { color: 'default', text: status }
    return <Tag color={info.color}>{info.text}</Tag>
  }

  // 后端暂未提供更新用户状态接口，功能已禁用
  // const handleBanUser = async (id: number, status: string) => {
  //   try {
  //     await adminAPI.updateUserStatus(id, { status })
  //     message.success(status === 'banned' ? '已封禁' : '已解封')
  //     loadUsers()
  //   } catch (err: any) {
  //     message.error(err.message || '操作失败')
  //   }
  // }

  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '用户名',
      dataIndex: 'username',
      key: 'username',
    },
    {
      title: '昵称',
      dataIndex: 'nickname',
      key: 'nickname',
    },
    {
      title: '邮箱',
      dataIndex: 'email',
      key: 'email',
    },
    {
      title: '手机号',
      dataIndex: 'phone',
      key: 'phone',
    },
    {
      title: '角色',
      dataIndex: 'role',
      key: 'role',
      render: (role: string) => getRoleTag(role),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => getStatusTag(status),
    },
    {
      title: '注册时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
    },
    // 后端暂未提供更新用户状态接口，操作列已禁用
    // {
    //   title: '操作',
    //   key: 'action',
    //   width: 150,
    //   render: (_: any, record: User) => (
    //     <Space>
    //       {record.status === 'active' ? (
    //         <Button
    //           danger
    //           size="small"
    //           onClick={() => handleBanUser(record.id, 'banned')}
    //         >
    //           封禁
    //         </Button>
    //       ) : (
    //         <Button
    //           type="primary"
    //           size="small"
    //           onClick={() => handleBanUser(record.id, 'active')}
    //         >
    //           解封
    //         </Button>
    //       )}
    //     </Space>
    //   ),
    // },
  ]

  return (
    <div>
      <Title level={3}>用户管理</Title>

      <Card>
        <div style={{ display: 'flex', gap: 16, marginBottom: 16 }}>
          <Input.Search
            placeholder="搜索用户名/昵称..."
            value={keyword}
            onChange={(e) => setKeyword(e.target.value)}
            onSearch={() => setPage(1)}
            style={{ width: 300 }}
            allowClear
          />
          <Select
            placeholder="选择角色"
            allowClear
            style={{ width: 150 }}
            value={role}
            onChange={(val) => { setRole(val); setPage(1) }}
          >
            <Option value="buyer">买家</Option>
            <Option value="seller">商家</Option>
            <Option value="admin">管理员</Option>
          </Select>
        </div>

        <Table
          columns={columns}
          dataSource={users}
          rowKey="id"
          pagination={false}
        />

        <div style={{ textAlign: 'center', marginTop: 24 }}>
          <Pagination
            current={page}
            pageSize={pageSize}
            total={total}
            onChange={setPage}
          />
        </div>
      </Card>
    </div>
  )
}

export default AdminUsers
