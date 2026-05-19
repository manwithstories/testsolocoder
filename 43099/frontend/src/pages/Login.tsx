import { useState } from 'react'
import { Form, Input, Button, Card, message, Tabs } from 'antd'
import { UserOutlined, LockOutlined, MailOutlined, PhoneOutlined } from '@ant-design/icons'
import { Link, useNavigate } from 'react-router-dom'
import { useAuthStore } from '@/store/authStore'

const Login = () => {
  const navigate = useNavigate()
  const { login, register } = useAuthStore()
  const [form] = Form.useForm()
  const [registerForm] = Form.useForm()
  const [loading, setLoading] = useState(false)
  const [registerLoading, setRegisterLoading] = useState(false)

  const handleLogin = async (values: any) => {
    setLoading(true)
    try {
      await login(values.username, values.password)
      message.success('登录成功')
      navigate('/')
    } catch (error: any) {
      message.error(error.message || '登录失败')
    } finally {
      setLoading(false)
    }
  }

  const handleRegister = async (values: any) => {
    setRegisterLoading(true)
    try {
      await register(values)
      message.success('注册成功')
      navigate('/')
    } catch (error: any) {
      message.error(error.message || '注册失败')
    } finally {
      setRegisterLoading(false)
    }
  }

  const items = [
    {
      key: 'login',
      label: '登录',
      children: (
        <Form form={form} onFinish={handleLogin} layout="vertical">
          <Form.Item
            name="username"
            rules={[{ required: true, message: '请输入用户名或邮箱' }]}
          >
            <Input prefix={<UserOutlined />} placeholder="用户名或邮箱" size="large" />
          </Form.Item>
          <Form.Item
            name="password"
            rules={[{ required: true, message: '请输入密码' }]}
          >
            <Input.Password prefix={<LockOutlined />} placeholder="密码" size="large" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block size="large" loading={loading}>
              登录
            </Button>
          </Form.Item>
          <div style={{ textAlign: 'center' }}>
            <Link to="/forgot-password">忘记密码？</Link>
          </div>
        </Form>
      ),
    },
    {
      key: 'register',
      label: '注册',
      children: (
        <Form form={registerForm} onFinish={handleRegister} layout="vertical">
          <Form.Item
            name="username"
            rules={[
              { required: true, message: '请输入用户名' },
              { min: 3, max: 50, message: '用户名长度3-50个字符' },
            ]}
          >
            <Input prefix={<UserOutlined />} placeholder="用户名" size="large" />
          </Form.Item>
          <Form.Item
            name="email"
            rules={[
              { required: true, message: '请输入邮箱' },
              { type: 'email', message: '请输入有效的邮箱地址' },
            ]}
          >
            <Input prefix={<MailOutlined />} placeholder="邮箱" size="large" />
          </Form.Item>
          <Form.Item
            name="password"
            rules={[
              { required: true, message: '请输入密码' },
              { min: 6, message: '密码至少6个字符' },
            ]}
          >
            <Input.Password prefix={<LockOutlined />} placeholder="密码" size="large" />
          </Form.Item>
          <Form.Item
            name="real_name"
            rules={[{ max: 50, message: '姓名最多50个字符' }]}
          >
            <Input placeholder="真实姓名（选填）" size="large" />
          </Form.Item>
          <Form.Item
            name="phone"
            rules={[{ max: 20, message: '电话最多20个字符' }]}
          >
            <Input prefix={<PhoneOutlined />} placeholder="联系电话（选填）" size="large" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block size="large" loading={registerLoading}>
              注册
            </Button>
          </Form.Item>
        </Form>
      ),
    },
  ]

  return (
    <div style={{
      minHeight: '100vh',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
    }}>
      <Card style={{ width: 450, boxShadow: '0 10px 40px rgba(0,0,0,0.1)' }}>
        <div style={{ textAlign: 'center', marginBottom: 24 }}>
          <h1 style={{ marginBottom: 8 }}>场地预约系统</h1>
          <p style={{ color: '#999' }}>Venue Booking Management System</p>
        </div>
        <Tabs items={items} centered defaultActiveKey="login" />
      </Card>
    </div>
  )
}

export default Login
