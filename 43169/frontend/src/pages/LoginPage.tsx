import { useState } from 'react'
import { Form, Input, Button, Card, message, Tabs } from 'antd'
import { UserOutlined, LockOutlined, MobileOutlined } from '@ant-design/icons'
import { useNavigate, Link } from 'react-router-dom'
import { userApi } from '@/api/endpoints'
import { useAuthStore } from '@/store/authStore'

export default function LoginPage() {
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()
  const { setToken, setUser } = useAuthStore()

  const onFinish = async (values: any) => {
    setLoading(true)
    try {
      const res = await userApi.login(values)
      setToken(res.token)
      setUser(res.user)
      message.success('登录成功')
      navigate('/')
    } catch {
      // error handled by interceptor
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="auth-container">
      <Card className="auth-card">
        <div className="auth-title">
          <h1>💕 相亲交友</h1>
          <p>遇见你的命中注定</p>
        </div>
        <Tabs
          defaultActiveKey="login"
          centered
          items={[
            {
              key: 'login',
              label: '登录',
              children: (
                <Form onFinish={onFinish} layout="vertical">
                  <Form.Item name="account" rules={[{ required: true, message: '请输入账号' }]}>
                    <Input prefix={<UserOutlined />} placeholder="用户名/手机号/邮箱" size="large" />
                  </Form.Item>
                  <Form.Item name="password" rules={[{ required: true, message: '请输入密码' }]}>
                    <Input.Password prefix={<LockOutlined />} placeholder="密码" size="large" />
                  </Form.Item>
                  <Form.Item>
                    <Button type="primary" htmlType="submit" loading={loading} block size="large">
                      登录
                    </Button>
                  </Form.Item>
                  <div style={{ textAlign: 'center' }}>
                    <span>还没有账号？</span>
                    <Link to="/register">立即注册</Link>
                  </div>
                </Form>
              ),
            },
            {
              key: 'phone',
              label: '手机登录',
              children: (
                <Form onFinish={onFinish} layout="vertical">
                  <Form.Item name="account" rules={[{ required: true, message: '请输入手机号' }]}>
                    <Input prefix={<MobileOutlined />} placeholder="手机号" size="large" />
                  </Form.Item>
                  <Form.Item name="password" rules={[{ required: true, message: '请输入密码' }]}>
                    <Input.Password prefix={<LockOutlined />} placeholder="密码" size="large" />
                  </Form.Item>
                  <Form.Item>
                    <Button type="primary" htmlType="submit" loading={loading} block size="large">
                      登录
                    </Button>
                  </Form.Item>
                </Form>
              ),
            },
          ]}
        />
      </Card>
    </div>
  )
}
