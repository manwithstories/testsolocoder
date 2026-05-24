import { useState } from 'react'
import { Form, Input, Button, Card, message } from 'antd'
import { UserOutlined, LockOutlined, MobileOutlined, MailOutlined, SafetyOutlined } from '@ant-design/icons'
import { useNavigate, Link } from 'react-router-dom'
import { userApi } from '@/api/endpoints'
import { useAuthStore } from '@/store/authStore'

export default function RegisterPage() {
  const [loading, setLoading] = useState(false)
  const [countdown, setCountdown] = useState(0)
  const navigate = useNavigate()
  const { setToken, setUser } = useAuthStore()
  const [form] = Form.useForm()

  const sendCode = async () => {
    const phone = form.getFieldValue('phone')
    if (!phone) {
      message.warning('请先输入手机号')
      return
    }
    try {
      await userApi.sendSmsCode(phone)
      message.success('验证码已发送')
      setCountdown(60)
      const timer = setInterval(() => {
        setCountdown((c) => {
          if (c <= 1) {
            clearInterval(timer)
            return 0
          }
          return c - 1
        })
      }, 1000)
    } catch {
      // handled
    }
  }

  const onFinish = async (values: any) => {
    if (values.password !== values.confirmPassword) {
      message.error('两次输入的密码不一致')
      return
    }
    setLoading(true)
    try {
      const res = await userApi.register({
        username: values.username,
        password: values.password,
        phone: values.phone,
        email: values.email,
        code: values.code,
      })
      setToken(res.token)
      setUser(res.user)
      message.success('注册成功')
      navigate('/')
    } catch {
      // handled
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="auth-container">
      <Card className="auth-card" style={{ width: 440 }}>
        <div className="auth-title">
          <h1>💕 创建账号</h1>
          <p>加入我们，开启你的缘分之旅</p>
        </div>
        <Form form={form} onFinish={onFinish} layout="vertical">
          <Form.Item name="username" rules={[{ required: true, message: '请输入用户名' }, { min: 3, message: '至少3个字符' }]}>
            <Input prefix={<UserOutlined />} placeholder="用户名" size="large" />
          </Form.Item>
          <Form.Item name="phone" rules={[{ required: true, message: '请输入手机号' }, { pattern: /^1[3-9]\d{9}$/, message: '手机号格式不正确' }]}>
            <Input prefix={<MobileOutlined />} placeholder="手机号" size="large" />
          </Form.Item>
          <Form.Item name="email" rules={[{ type: 'email', message: '邮箱格式不正确' }]}>
            <Input prefix={<MailOutlined />} placeholder="邮箱（选填）" size="large" />
          </Form.Item>
          <Form.Item name="password" rules={[{ required: true, message: '请输入密码' }, { min: 6, message: '至少6个字符' }]}>
            <Input.Password prefix={<LockOutlined />} placeholder="密码" size="large" />
          </Form.Item>
          <Form.Item name="confirmPassword" rules={[{ required: true, message: '请确认密码' }]}>
            <Input.Password prefix={<LockOutlined />} placeholder="确认密码" size="large" />
          </Form.Item>
          <Form.Item name="code" rules={[{ required: true, message: '请输入验证码' }]}>
            <Input
              prefix={<SafetyOutlined />}
              placeholder="短信验证码"
              size="large"
              suffix={
                <Button size="small" onClick={sendCode} disabled={countdown > 0}>
                  {countdown > 0 ? `${countdown}s` : '获取验证码'}
                </Button>
              }
            />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading} block size="large">
              注册
            </Button>
          </Form.Item>
          <div style={{ textAlign: 'center' }}>
            <span>已有账号？</span>
            <Link to="/login">返回登录</Link>
          </div>
        </Form>
      </Card>
    </div>
  )
}
