import { useState } from 'react'
import { Form, Input, Button, Card, Typography, message, Tabs } from 'antd'
import { UserOutlined, LockOutlined } from '@ant-design/icons'
import { useNavigate, Link } from 'react-router-dom'
import { useMutation } from '@tanstack/react-query'
import { authApi } from '@/services/api'
import { useAuthStore } from '@/context/AuthContext'
import type { LoginData, RegisterData } from '@/types'

const { Title, Text } = Typography

export default function Login() {
  const navigate = useNavigate()
  const setAuth = useAuthStore((state) => state.setAuth)
  const [activeTab, setActiveTab] = useState('login')

  const loginMutation = useMutation({
    mutationFn: (data: { account: string; password: string }) => authApi.login(data),
    onSuccess: (res: any) => {
      const data: LoginData = res.data
      setAuth(data)
      message.success('登录成功')
      navigate('/', { replace: true })
    },
    onError: (err: any) => {
      message.error(err.message || '登录失败')
    },
  })

  const registerMutation = useMutation({
    mutationFn: (data: RegisterData) => authApi.register(data),
    onSuccess: (res: any) => {
      const data: LoginData = res.data
      setAuth(data)
      message.success('注册成功')
      navigate('/', { replace: true })
    },
    onError: (err: any) => {
      message.error(err.message || '注册失败')
    },
  })

  const handleLogin = (values: any) => {
    loginMutation.mutate(values)
  }

  const handleRegister = (values: any) => {
    registerMutation.mutate(values)
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-sky-100 to-blue-50 p-4">
      <Card className="w-full max-w-md shadow-lg" bordered={false}>
        <div className="text-center mb-6">
          <div className="text-4xl mb-2">🐾</div>
          <Title level={3} style={{ marginBottom: 8 }}>
            宠物寄养与日托平台
          </Title>
          <Text type="secondary">为您的宠物提供专业的寄养和日托服务</Text>
        </div>

        <Tabs
          activeKey={activeTab}
          onChange={setActiveTab}
          centered
          items={[
            {
              key: 'login',
              label: '登录',
              children: (
                <Form onFinish={handleLogin} layout="vertical" size="large">
                  <Form.Item
                    name="account"
                    rules={[{ required: true, message: '请输入用户名或邮箱' }]}
                  >
                    <Input prefix={<UserOutlined />} placeholder="用户名或邮箱" />
                  </Form.Item>
                  <Form.Item
                    name="password"
                    rules={[{ required: true, message: '请输入密码' }]}
                  >
                    <Input.Password prefix={<LockOutlined />} placeholder="密码" />
                  </Form.Item>
                  <Form.Item>
                    <Button
                      type="primary"
                      htmlType="submit"
                      block
                      loading={loginMutation.isPending}
                    >
                      登录
                    </Button>
                  </Form.Item>
                </Form>
              ),
            },
            {
              key: 'register',
              label: '注册',
              children: (
                <Form onFinish={handleRegister} layout="vertical" size="large">
                  <Form.Item
                    name="username"
                    rules={[
                      { required: true, message: '请输入用户名' },
                      { min: 3, message: '用户名至少3个字符' },
                    ]}
                  >
                    <Input prefix={<UserOutlined />} placeholder="用户名" />
                  </Form.Item>
                  <Form.Item
                    name="email"
                    rules={[
                      { required: true, message: '请输入邮箱' },
                      { type: 'email', message: '请输入有效的邮箱地址' },
                    ]}
                  >
                    <Input placeholder="邮箱" />
                  </Form.Item>
                  <Form.Item name="phone">
                    <Input placeholder="手机号（选填）" />
                  </Form.Item>
                  <Form.Item
                    name="password"
                    rules={[
                      { required: true, message: '请输入密码' },
                      { min: 6, message: '密码至少6个字符' },
                    ]}
                  >
                    <Input.Password prefix={<LockOutlined />} placeholder="密码" />
                  </Form.Item>
                  <Form.Item
                    name="role"
                    rules={[{ required: true, message: '请选择角色' }]}
                  >
                    <select
                      className="w-full h-10 px-3 border rounded-md"
                      defaultValue=""
                    >
                      <option value="" disabled>
                        选择角色
                      </option>
                      <option value="owner">宠物主人</option>
                      <option value="store">寄养门店</option>
                      <option value="keeper">宠物管家</option>
                    </select>
                  </Form.Item>
                  <Form.Item name="real_name">
                    <Input placeholder="真实姓名（选填）" />
                  </Form.Item>
                  <Form.Item>
                    <Button
                      type="primary"
                      htmlType="submit"
                      block
                      loading={registerMutation.isPending}
                    >
                      注册
                    </Button>
                  </Form.Item>
                </Form>
              ),
            },
          ]}
        />

        <div className="text-center mt-4">
          <Text type="secondary">
            {activeTab === 'login' ? '还没有账号？' : '已有账号？'}
            <Link
              to="#"
              onClick={(e) => {
                e.preventDefault()
                setActiveTab(activeTab === 'login' ? 'register' : 'login')
              }}
              className="text-sky-500 ml-1"
            >
              {activeTab === 'login' ? '立即注册' : '立即登录'}
            </Link>
          </Text>
        </div>
      </Card>
    </div>
  )
}
