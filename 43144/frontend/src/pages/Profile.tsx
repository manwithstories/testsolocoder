import React, { useEffect } from 'react'
import { Card, Form, Input, Button, Avatar, message, Descriptions, Row, Col } from 'antd'
import { UserOutlined } from '@ant-design/icons'
import { useAuth } from '../contexts/AuthContext'
import { updateProfile } from '../api/auth'

const Profile: React.FC = () => {
  const { user, loadUser } = useAuth()
  const [form] = Form.useForm()
  const [loading, setLoading] = React.useState(false)

  useEffect(() => {
    if (user) {
      form.setFieldsValue({
        name: user.name,
        phone: user.phone,
        address: user.address,
      })
    }
  }, [user])

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      setLoading(true)
      await updateProfile(values)
      message.success('更新成功')
      loadUser()
    } catch (error: any) {
      message.error(error.message || '更新失败')
    } finally {
      setLoading(false)
    }
  }

  if (!user) {
    return <div>加载中...</div>
  }

  const roleText: Record<string, string> = {
    admin: '管理员',
    rescue: '救助站',
    adopter: '领养人',
  }

  return (
    <div>
      <h2 style={{ marginBottom: 16 }}>个人设置</h2>
      <Row gutter={16}>
        <Col xs={24} md={8}>
          <Card style={{ textAlign: 'center' }}>
            <Avatar size={96} icon={<UserOutlined />} src={user.avatar} />
            <h3 style={{ marginTop: 16 }}>{user.name}</h3>
            <p style={{ color: '#999' }}>{roleText[user.role] || user.role}</p>
            <p>{user.email}</p>
          </Card>
        </Col>
        <Col xs={24} md={16}>
          <Card title="基本信息">
            <Form form={form} onFinish={handleSubmit} layout="vertical">
              <Form.Item name="name" label="姓名" rules={[{ required: true }]}>
                <Input />
              </Form.Item>
              <Form.Item name="phone" label="手机号">
                <Input />
              </Form.Item>
              <Form.Item name="address" label="地址">
                <Input.TextArea rows={2} />
              </Form.Item>
              <Form.Item>
                <Button type="primary" htmlType="submit" loading={loading}>
                  保存
                </Button>
              </Form.Item>
            </Form>
          </Card>
        </Col>
      </Row>
    </div>
  )
}

export default Profile
