import React, { useState } from 'react'
import { Form, Input, Button, Card, Typography, message, Tabs, Radio } from 'antd'
import { UserOutlined, LockOutlined, MailOutlined } from '@ant-design/icons'
import { useNavigate, Link } from 'react-router-dom'
import { useAuthStore } from '@/store/auth'

const { Title, Text } = Typography

const LoginPage: React.FC = () => {
  const navigate = useNavigate()
  const { login } = useAuthStore()
  const [loading, setLoading] = useState(false)

  const handleLogin = async (values: { account: string; password: string }) => {
    setLoading(true)
    try {
      await login(values.account, values.password)
      message.success('登录成功')
      navigate('/')
    } catch (error: any) {
      message.error(error.message || '登录失败')
    } finally {
      setLoading(false)
    }
  }

  const handleRegister = async (values: any) => {
    setLoading(true)
    try {
      const { register } = useAuthStore.getState()
      await register(values)
      message.success('注册成功，请登录')
    } catch (error: any) {
      message.error(error.message || '注册失败')
    } finally {
      setLoading(false)
    }
  }

  const loginItems = [
    {
      key: 'login',
      label: '登录',
      children: (
        <Form name="login" onFinish={handleLogin} autoComplete="off" size="large">
          <Form.Item name="account" rules={[{ required: true, message: '请输入用户名或邮箱' }]}>
            <Input prefix={<UserOutlined />} placeholder="用户名或邮箱" />
          </Form.Item>
          <Form.Item name="password" rules={[{ required: true, message: '请输入密码' }]}>
            <Input.Password prefix={<LockOutlined />} placeholder="密码" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading} block>
              登录
            </Button>
          </Form.Item>
          <div style={{ textAlign: 'center' }}>
            <Text>还没有账号？ </Text>
            <Link to="/register">立即注册</Link>
          </div>
        </Form>
      ),
    },
  ]

  const registerItems = [
    {
      key: 'register',
      label: '注册',
      children: (
        <Form name="register" onFinish={handleRegister} autoComplete="off" size="large">
          <Form.Item name="username" rules={[{ required: true, message: '请输入用户名' }, { min: 3, max: 50 }]}>
            <Input prefix={<UserOutlined />} placeholder="用户名" />
          </Form.Item>
          <Form.Item name="email" rules={[{ required: true, message: '请输入邮箱' }, { type: 'email' }]}>
            <Input prefix={<MailOutlined />} placeholder="邮箱" />
          </Form.Item>
          <Form.Item name="password" rules={[{ required: true, message: '请输入密码' }, { min: 6, max: 50 }]}>
            <Input.Password prefix={<LockOutlined />} placeholder="密码" />
          </Form.Item>
          <Form.Item name="nickname">
            <Input placeholder="昵称（可选）" />
          </Form.Item>
          <Form.Item name="role" rules={[{ required: true }]} initialValue="student">
            <Radio.Group>
              <Radio value="student">学员</Radio>
              <Radio value="instructor">讲师</Radio>
            </Radio.Group>
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading} block>
              注册
            </Button>
          </Form.Item>
          <div style={{ textAlign: 'center' }}>
            <Text>已有账号？ </Text>
            <Link to="/login">立即登录</Link>
          </div>
        </Form>
      ),
    },
  ]

  return (
    <div
      style={{
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
      }}
    >
      <Card style={{ width: 420, boxShadow: '0 4px 20px rgba(0,0,0,0.15)' }}>
        <div style={{ textAlign: 'center', marginBottom: 24 }}>
          <Title level={3} style={{ marginBottom: 8 }}>在线学习平台</Title>
          <Text type="secondary">开启你的学习之旅</Text>
        </div>
        <Tabs
          defaultActiveKey="login"
          centered
          items={[...loginItems, ...registerItems]}
        />
      </Card>
    </div>
  )
}

export default LoginPage
