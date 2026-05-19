import { useState, useEffect } from 'react'
import { Card, Form, Input, Button, message, Avatar, Descriptions } from 'antd'
import { UserOutlined } from '@ant-design/icons'
import { useAuthStore } from '@/store/authStore'
import { userApi } from '@/api'
import type { User } from '@/types'
import dayjs from 'dayjs'

const Profile = () => {
  const { user, fetchUser } = useAuthStore()
  const [form] = Form.useForm()
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (user) {
      form.setFieldsValue({
        real_name: user.real_name,
        phone: user.phone,
        avatar: user.avatar,
      })
    }
  }, [user])

  const handleSubmit = async (values: any) => {
    setLoading(true)
    try {
      await userApi.updateMe(values)
      await fetchUser()
      message.success('更新成功')
    } catch (error: any) {
      message.error(error.message || '更新失败')
    } finally {
      setLoading(false)
    }
  }

  const getRoleText = (role: string) => {
    const map: Record<string, string> = {
      user: '普通用户',
      admin: '管理员',
      super_admin: '超级管理员',
    }
    return map[role] || role
  }

  if (!user) return null

  return (
    <div>
      <Card title="个人信息" style={{ marginBottom: 16 }}>
        <div style={{ display: 'flex', alignItems: 'center', marginBottom: 24 }}>
          <Avatar size={80} src={user.avatar} icon={<UserOutlined />} />
          <div style={{ marginLeft: 24 }}>
            <h2 style={{ margin: 0 }}>{user.real_name || user.username}</h2>
            <p style={{ color: '#999', margin: '8px 0' }}>@{user.username}</p>
            <p style={{ color: '#999', margin: 0 }}>角色: {getRoleText(user.role)}</p>
          </div>
        </div>

        <Descriptions bordered column={2}>
          <Descriptions.Item label="邮箱">{user.email}</Descriptions.Item>
          <Descriptions.Item label="电话">{user.phone || '-'}</Descriptions.Item>
          <Descriptions.Item label="邮箱验证">
            {user.email_verified ? '已验证' : '未验证'}
          </Descriptions.Item>
          <Descriptions.Item label="注册时间">
            {dayjs(user.created_at).format('YYYY-MM-DD HH:mm')}
          </Descriptions.Item>
        </Descriptions>
      </Card>

      <Card title="编辑信息">
        <Form form={form} layout="vertical" onFinish={handleSubmit} style={{ maxWidth: 500 }}>
          <Form.Item
            name="real_name"
            label="真实姓名"
            rules={[{ max: 50, message: '最多50个字符' }]}
          >
            <Input placeholder="请输入真实姓名" />
          </Form.Item>
          <Form.Item
            name="phone"
            label="联系电话"
            rules={[{ max: 20, message: '最多20个字符' }]}
          >
            <Input placeholder="请输入联系电话" />
          </Form.Item>
          <Form.Item
            name="avatar"
            label="头像URL"
          >
            <Input placeholder="请输入头像图片URL" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading}>
              保存修改
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  )
}

export default Profile
