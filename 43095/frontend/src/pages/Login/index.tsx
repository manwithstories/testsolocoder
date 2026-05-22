import { useState } from 'react'
import { Form, Input, Button, Card, Checkbox, message } from 'antd'
import { UserOutlined, LockOutlined } from '@ant-design/icons'
import { Link, useNavigate } from 'react-router-dom'
import { useAuthStore } from '@/store'
import { authAPI } from '@/services/api'
import type { LoginRequest } from '@/types'

const Login = () => {
  const [loading, setLoading] = useState(false)
  const [form] = Form.useForm()
  const navigate = useNavigate()
  const { login } = useAuthStore()

  const onFinish = async (values: LoginRequest) => {
    setLoading(true)
    try {
      const data = await authAPI.login(values)
      login(data.token, data.user)
      message.success('登录成功')
      navigate('/')
    } catch (error) {
      message.error('登录失败，请稍后重试')
    } finally {
      setLoading(false)
    }
  }

  return (
    <Card className="shadow-xl">
      <h2 className="text-2xl font-bold text-center mb-6">用户登录</h2>
      <Form
        form={form}
        name="login"
        onFinish={onFinish}
        autoComplete="off"
        size="large"
      >
        <Form.Item
          name="username"
          rules={[{ required: true, message: '请输入用户名' }]}
        >
          <Input prefix={<UserOutlined />} placeholder="用户名" />
        </Form.Item>

        <Form.Item
          name="password"
          rules={[{ required: true, message: '请输入密码' }]}
        >
          <Input.Password prefix={<LockOutlined />} placeholder="密码" />
        </Form.Item>

        <Form.Item>
          <div className="flex justify-between items-center">
            <Form.Item name="remember" valuePropName="checked" noStyle>
              <Checkbox>记住我</Checkbox>
            </Form.Item>
            <a href="#" className="text-blue-500">忘记密码？</a>
          </div>
        </Form.Item>

        <Form.Item>
          <Button type="primary" htmlType="submit" block loading={loading}>
            登录
          </Button>
        </Form.Item>

        <div className="text-center">
          还没有账号？ <Link to="/register" className="text-blue-500">立即注册</Link>
        </div>
      </Form>
    </Card>
  )
}

export default Login
