import React, { useState } from 'react'
import { Form, Input, Button, Card, message, Select, Radio } from 'antd'
import { UserOutlined, LockOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'

const { Option } = Select

const Register: React.FC = () => {
  const [form] = Form.useForm()
  const { register } = useAuth()
  const navigate = useNavigate()
  const [loading, setLoading] = useState(false)
  const [role, setRole] = useState<string>('adopter')

  const onFinish = async (values: any) => {
    setLoading(true)
    try {
      await register({
        email: values.email,
        password: values.password,
        name: values.name,
        phone: values.phone,
        role: values.role,
        rescue_name: values.rescue_name,
      })
      message.success('注册成功，请登录')
      navigate('/login')
    } catch (error: any) {
      message.error(error.message || '注册失败')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="register-container">
      <Card className="register-card">
        <div style={{ textAlign: 'center', marginBottom: 24 }}>
          <h1 style={{ margin: 0 }}>🐾 注册账号</h1>
          <p style={{ color: '#999', marginTop: 8 }}>创建您的宠物领养平台账号</p>
        </div>
        <Form
          form={form}
          name="register"
          onFinish={onFinish}
          autoComplete="off"
          size="large"
        >
          <Form.Item
            name="role"
            rules={[{ required: true, message: '请选择角色' }]}
            initialValue="adopter"
          >
            <Radio.Group onChange={(e) => setRole(e.target.value)}>
              <Radio value="adopter">领养人</Radio>
              <Radio value="rescue">救助站</Radio>
            </Radio.Group>
          </Form.Item>

          {role === 'rescue' && (
            <Form.Item
              name="rescue_name"
              rules={[{ required: true, message: '请输入救助站名称' }]}
            >
              <Input placeholder="救助站名称" />
            </Form.Item>
          )}

          <Form.Item
            name="name"
            rules={[{ required: true, message: '请输入姓名' }]}
          >
            <Input prefix={<UserOutlined />} placeholder={role === 'rescue' ? '联系人姓名' : '姓名'} />
          </Form.Item>

          <Form.Item
            name="email"
            rules={[{ required: true, message: '请输入邮箱' }, { type: 'email', message: '请输入有效邮箱' }]}
          >
            <Input prefix={<UserOutlined />} placeholder="邮箱" />
          </Form.Item>

          <Form.Item
            name="password"
            rules={[{ required: true, message: '请输入密码' }, { min: 6, message: '密码至少6位' }]}
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
                  return Promise.reject(new Error('两次密码输入不一致'))
                },
              }),
            ]}
          >
            <Input.Password prefix={<LockOutlined />} placeholder="确认密码" />
          </Form.Item>

          <Form.Item name="phone">
            <Input placeholder="手机号（选填）" />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading} block>
              注册
            </Button>
          </Form.Item>
          <div style={{ textAlign: 'center' }}>
            已有账号？ <a onClick={() => navigate('/login')}>立即登录</a>
          </div>
        </Form>
      </Card>
    </div>
  )
}

export default Register
