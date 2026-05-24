import React, { useEffect, useState } from 'react'
import { Form, Input, Button, Card, message, Tabs } from 'antd'
import { UserOutlined, LockOutlined } from '@ant-design/icons'
import { useNavigate, useLocation } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'

const Login: React.FC = () => {
  const [form] = Form.useForm()
  const { login, isAuthenticated } = useAuth()
  const navigate = useNavigate()
  const location = useLocation()
  const [loading, setLoading] = useState(false)

  const from = (location.state as any)?.from?.pathname || '/'

  useEffect(() => {
    if (isAuthenticated) {
      navigate(from, { replace: true })
    }
  }, [isAuthenticated, navigate, from])

  const onFinish = async (values: { email: string; password: string }) => {
    setLoading(true)
    try {
      await login(values.email, values.password)
      message.success('登录成功')
      navigate(from, { replace: true })
    } catch (error: any) {
      message.error(error.message || '登录失败')
    } finally {
      setLoading(false)
    }
  }

  const loginTab = (
    <Form
      form={form}
      name="login"
      onFinish={onFinish}
      autoComplete="off"
      size="large"
    >
      <Form.Item
        name="email"
        rules={[{ required: true, message: '请输入邮箱' }, { type: 'email', message: '请输入有效邮箱' }]}
      >
        <Input prefix={<UserOutlined />} placeholder="邮箱" />
      </Form.Item>
      <Form.Item
        name="password"
        rules={[{ required: true, message: '请输入密码' }]}
      >
        <Input.Password prefix={<LockOutlined />} placeholder="密码" />
      </Form.Item>
      <Form.Item>
        <Button type="primary" htmlType="submit" loading={loading} block>
          登录
        </Button>
      </Form.Item>
      <div style={{ textAlign: 'center' }}>
        还没有账号？ <a onClick={() => navigate('/register')}>立即注册</a>
      </div>
    </Form>
  )

  return (
    <div className="login-container">
      <Card className="login-card">
        <div style={{ textAlign: 'center', marginBottom: 24 }}>
          <h1 style={{ margin: 0 }}>🐾 宠物领养平台</h1>
          <p style={{ color: '#999', marginTop: 8 }}>欢迎回来，请登录您的账号</p>
        </div>
        {loginTab}
      </Card>
    </div>
  )
}

export default Login
