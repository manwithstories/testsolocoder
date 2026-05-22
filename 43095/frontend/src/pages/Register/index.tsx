import { useState } from 'react'
import { Form, Input, Button, Card, Select, message } from 'antd'
import { UserOutlined, LockOutlined, MailOutlined, PhoneOutlined } from '@ant-design/icons'
import { Link, useNavigate } from 'react-router-dom'
import { authAPI } from '@/services/api'
import type { RegisterRequest } from '@/types'

const Register = () => {
  const [loading, setLoading] = useState(false)
  const [form] = Form.useForm()
  const navigate = useNavigate()

  const onFinish = async (values: RegisterRequest & { confirmPassword: string }) => {
    setLoading(true)
    try {
      const { confirmPassword, ...data } = values
      await authAPI.register(data)
      message.success('注册成功，请登录')
      navigate('/login')
    } catch (error) {
      message.error('注册失败，请稍后重试')
    } finally {
      setLoading(false)
    }
  }

  return (
    <Card className="shadow-xl">
      <h2 className="text-2xl font-bold text-center mb-6">用户注册</h2>
      <Form
        form={form}
        name="register"
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
          name="full_name"
          rules={[{ required: true, message: '请输入姓名' }]}
        >
          <Input placeholder="真实姓名" />
        </Form.Item>

        <Form.Item
          name="email"
          rules={[
            { required: true, message: '请输入邮箱' },
            { type: 'email', message: '请输入有效的邮箱地址' }
          ]}
        >
          <Input prefix={<MailOutlined />} placeholder="邮箱" />
        </Form.Item>

        <Form.Item
          name="phone"
          rules={[{ required: true, message: '请输入手机号' }]}
        >
          <Input prefix={<PhoneOutlined />} placeholder="手机号" />
        </Form.Item>

        <Form.Item
          name="role"
          rules={[{ required: true, message: '请选择角色' }]}
          initialValue="patient"
        >
          <Select placeholder="选择角色">
            <Select.Option value="patient">患者</Select.Option>
            <Select.Option value="doctor">医生</Select.Option>
          </Select>
        </Form.Item>

        <Form.Item
          name="password"
          rules={[
            { required: true, message: '请输入密码' },
            { min: 6, message: '密码至少6位' }
          ]}
        >
          <Input.Password prefix={<LockOutlined />} placeholder="密码" />
        </Form.Item>

        <Form.Item
          name="confirmPassword"
          dependencies={['password']}
          rules={[
            { required: true, message: '请确认密码' },
            ({ getFieldValue }) => ({
              validator(_, value) {
                if (!value || getFieldValue('password') === value) {
                  return Promise.resolve()
                }
                return Promise.reject(new Error('两次输入的密码不一致'))
              }
            })
          ]}
        >
          <Input.Password prefix={<LockOutlined />} placeholder="确认密码" />
        </Form.Item>

        <Form.Item>
          <Button type="primary" htmlType="submit" block loading={loading}>
            注册
          </Button>
        </Form.Item>

        <div className="text-center">
          已有账号？ <Link to="/login" className="text-blue-500">立即登录</Link>
        </div>
      </Form>
    </Card>
  )
}

export default Register
