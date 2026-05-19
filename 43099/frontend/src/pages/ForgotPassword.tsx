import { useState } from 'react'
import { Form, Input, Button, Card, message, Tabs } from 'antd'
import { MailOutlined, LockOutlined } from '@ant-design/icons'
import { Link, useNavigate } from 'react-router-dom'
import { authApi } from '@/api'

const ForgotPassword = () => {
  const navigate = useNavigate()
  const [form] = Form.useForm()
  const [resetForm] = Form.useForm()
  const [loading, setLoading] = useState(false)
  const [resetLoading, setResetLoading] = useState(false)

  const handleSendCode = async (values: any) => {
    setLoading(true)
    try {
      await authApi.forgotPassword(values.email)
      message.success('验证码已发送到您的邮箱')
    } catch (error: any) {
      message.error(error.message || '发送失败')
    } finally {
      setLoading(false)
    }
  }

  const handleResetPassword = async (values: any) => {
    setResetLoading(true)
    try {
      await authApi.resetPassword(values)
      message.success('密码重置成功，请重新登录')
      navigate('/login')
    } catch (error: any) {
      message.error(error.message || '重置失败')
    } finally {
      setResetLoading(false)
    }
  }

  const items = [
    {
      key: 'send',
      label: '发送验证码',
      children: (
        <Form form={form} onFinish={handleSendCode} layout="vertical">
          <Form.Item
            name="email"
            rules={[
              { required: true, message: '请输入邮箱' },
              { type: 'email', message: '请输入有效的邮箱地址' },
            ]}
          >
            <Input prefix={<MailOutlined />} placeholder="请输入注册邮箱" size="large" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block size="large" loading={loading}>
              发送验证码
            </Button>
          </Form.Item>
        </Form>
      ),
    },
    {
      key: 'reset',
      label: '重置密码',
      children: (
        <Form form={resetForm} onFinish={handleResetPassword} layout="vertical">
          <Form.Item
            name="email"
            rules={[
              { required: true, message: '请输入邮箱' },
              { type: 'email', message: '请输入有效的邮箱地址' },
            ]}
          >
            <Input prefix={<MailOutlined />} placeholder="请输入邮箱" size="large" />
          </Form.Item>
          <Form.Item
            name="code"
            rules={[{ required: true, message: '请输入验证码' }]}
          >
            <Input placeholder="请输入验证码" size="large" />
          </Form.Item>
          <Form.Item
            name="new_password"
            rules={[
              { required: true, message: '请输入新密码' },
              { min: 6, message: '密码至少6个字符' },
            ]}
          >
            <Input.Password prefix={<LockOutlined />} placeholder="请输入新密码" size="large" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block size="large" loading={resetLoading}>
              重置密码
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
          <h2 style={{ marginBottom: 8 }}>忘记密码</h2>
          <p style={{ color: '#999' }}>请输入您的注册邮箱</p>
        </div>
        <Tabs items={items} centered defaultActiveKey="send" />
        <div style={{ textAlign: 'center', marginTop: 16 }}>
          <Link to="/login">返回登录</Link>
        </div>
      </Card>
    </div>
  )
}

export default ForgotPassword
