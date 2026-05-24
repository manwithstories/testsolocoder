import { useState } from 'react'
import { Form, Input, Button, Card, Tabs, message } from 'antd'
import { UserOutlined, LockOutlined } from '@ant-design/icons'
import api from '../api'
import type { User } from '../types'

export default function Login({ onLogin }: { onLogin: (user: User) => void }) {
  const [loading, setLoading] = useState(false)
  const [tab, setTab] = useState('login')

  const handleLogin = async (values: any) => {
    setLoading(true)
    try {
      const res = await api.post('/auth/login', values)
      localStorage.setItem('token', res.data.token)
      localStorage.setItem('user', JSON.stringify(res.data.user))
      onLogin(res.data.user)
      message.success('登录成功')
    } catch (e: any) {
      message.error(e.response?.data?.error || '登录失败')
    } finally {
      setLoading(false)
    }
  }

  const handleRegister = async (values: any) => {
    setLoading(true)
    try {
      const res = await api.post('/auth/register', values)
      localStorage.setItem('token', res.data.token)
      localStorage.setItem('user', JSON.stringify(res.data.user))
      onLogin(res.data.user)
      message.success('注册成功')
    } catch (e: any) {
      message.error(e.response?.data?.error || '注册失败')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '100vh', background: '#f0f2f5' }}>
      <Card style={{ width: 420 }} title="体育联赛管理平台">
        <Tabs
          activeKey={tab}
          onChange={setTab}
          items={[
            {
              key: 'login',
              label: '登录',
              children: (
                <Form onFinish={handleLogin} layout="vertical">
                  <Form.Item name="email" rules={[{ required: true, type: 'email', message: '请输入邮箱' }]}>
                    <Input prefix={<UserOutlined />} placeholder="邮箱" />
                  </Form.Item>
                  <Form.Item name="password" rules={[{ required: true, message: '请输入密码' }]}>
                    <Input.Password prefix={<LockOutlined />} placeholder="密码" />
                  </Form.Item>
                  <Button type="primary" htmlType="submit" loading={loading} block>
                    登录
                  </Button>
                </Form>
              )
            },
            {
              key: 'register',
              label: '注册',
              children: (
                <Form onFinish={handleRegister} layout="vertical">
                  <Form.Item name="full_name" rules={[{ required: true, message: '请输入姓名' }]}>
                    <Input placeholder="姓名" />
                  </Form.Item>
                  <Form.Item name="email" rules={[{ required: true, type: 'email', message: '请输入邮箱' }]}>
                    <Input prefix={<UserOutlined />} placeholder="邮箱" />
                  </Form.Item>
                  <Form.Item name="password" rules={[{ required: true, min: 6, message: '密码至少6位' }]}>
                    <Input.Password prefix={<LockOutlined />} placeholder="密码" />
                  </Form.Item>
                  <Form.Item name="phone">
                    <Input placeholder="手机号" />
                  </Form.Item>
                  <Form.Item name="role" initialValue="player">
                    <Input placeholder="角色 (player/captain/referee)" />
                  </Form.Item>
                  <Button type="primary" htmlType="submit" loading={loading} block>
                    注册
                  </Button>
                </Form>
              )
            }
          ]}
        />
      </Card>
    </div>
  )
}
