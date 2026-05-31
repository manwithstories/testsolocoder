import { useState } from 'react'
import { Form, Input, Button, message, Tabs } from 'antd'
import { UserOutlined, LockOutlined, MailOutlined, PhoneOutlined } from '@ant-design/icons'
import { useNavigate, Link } from 'react-router-dom'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useAuthContext } from '@/contexts/AuthContext'
import { authApi } from '@/services/auth'

const loginSchema = z.object({
  username: z.string().min(1, '请输入用户名'),
  password: z.string().min(1, '请输入密码'),
})

const registerSchema = z.object({
  username: z.string().min(3, '用户名至少3个字符').max(50, '用户名最多50个字符'),
  email: z.string().email('请输入有效的邮箱地址'),
  password: z.string().min(6, '密码至少6个字符').max(100, '密码最多100个字符'),
  role: z.enum(['client', 'professional'], { required_error: '请选择用户类型' }),
  full_name: z.string().min(1, '请输入姓名'),
  phone: z.string().optional(),
})

type LoginFormData = z.infer<typeof loginSchema>
type RegisterFormData = z.infer<typeof registerSchema>

export function Login() {
  const navigate = useNavigate()
  const { login } = useAuthContext()
  const [loading, setLoading] = useState(false)

  const {
    register: registerLogin,
    handleSubmit: handleLoginSubmit,
    formState: { errors: loginErrors },
  } = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
  })

  const {
    register: registerUser,
    handleSubmit: handleRegisterSubmit,
    formState: { errors: registerErrors },
    watch,
  } = useForm<RegisterFormData>({
    resolver: zodResolver(registerSchema),
    defaultValues: {
      role: 'client',
    },
  })

  const watchRole = watch('role')

  const onLogin = async (data: LoginFormData) => {
    try {
      setLoading(true)
      const result = await authApi.login(data)
      login(result.access_token, result.refresh_token, result.user)
      message.success('登录成功')
      navigate('/')
    } catch (error: any) {
      message.error(error.message || '登录失败')
    } finally {
      setLoading(false)
    }
  }

  const onRegister = async (data: RegisterFormData) => {
    try {
      setLoading(true)
      await authApi.register(data)
      message.success('注册成功，请登录')
    } catch (error: any) {
      message.error(error.message || '注册失败')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="auth-container">
      <div className="auth-card">
        <h2>专业咨询服务平台</h2>
        <Tabs
          defaultActiveKey="login"
          centered
          items={[
            {
              key: 'login',
              label: '登录',
              children: (
                <Form layout="vertical" onFinish={handleLoginSubmit(onLogin)}>
                  <Form.Item
                    label="用户名"
                    validateStatus={loginErrors.username ? 'error' : ''}
                    help={loginErrors.username?.message}
                  >
                    <Input
                      prefix={<UserOutlined />}
                      placeholder="请输入用户名"
                      {...registerLogin('username')}
                    />
                  </Form.Item>
                  <Form.Item
                    label="密码"
                    validateStatus={loginErrors.password ? 'error' : ''}
                    help={loginErrors.password?.message}
                  >
                    <Input.Password
                      prefix={<LockOutlined />}
                      placeholder="请输入密码"
                      {...registerLogin('password')}
                    />
                  </Form.Item>
                  <Form.Item>
                    <Button type="primary" htmlType="submit" loading={loading} block>
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
                <Form layout="vertical" onFinish={handleRegisterSubmit(onRegister)}>
                  <Form.Item
                    label="用户名"
                    validateStatus={registerErrors.username ? 'error' : ''}
                    help={registerErrors.username?.message}
                  >
                    <Input
                      prefix={<UserOutlined />}
                      placeholder="请输入用户名"
                      {...registerUser('username')}
                    />
                  </Form.Item>
                  <Form.Item
                    label="邮箱"
                    validateStatus={registerErrors.email ? 'error' : ''}
                    help={registerErrors.email?.message}
                  >
                    <Input
                      prefix={<MailOutlined />}
                      placeholder="请输入邮箱"
                      {...registerUser('email')}
                    />
                  </Form.Item>
                  <Form.Item
                    label="密码"
                    validateStatus={registerErrors.password ? 'error' : ''}
                    help={registerErrors.password?.message}
                  >
                    <Input.Password
                      prefix={<LockOutlined />}
                      placeholder="请输入密码"
                      {...registerUser('password')}
                    />
                  </Form.Item>
                  <Form.Item
                    label="姓名"
                    validateStatus={registerErrors.full_name ? 'error' : ''}
                    help={registerErrors.full_name?.message}
                  >
                    <Input placeholder="请输入姓名" {...registerUser('full_name')} />
                  </Form.Item>
                  <Form.Item
                    label="手机号"
                    validateStatus={registerErrors.phone ? 'error' : ''}
                    help={registerErrors.phone?.message}
                  >
                    <Input
                      prefix={<PhoneOutlined />}
                      placeholder="请输入手机号（选填）"
                      {...registerUser('phone')}
                    />
                  </Form.Item>
                  <Form.Item
                    label="用户类型"
                    validateStatus={registerErrors.role ? 'error' : ''}
                    help={registerErrors.role?.message}
                  >
                    <Input.Group compact>
                      <Form.Item noStyle>
                        <Input.Radio.Group
                          value={watchRole || 'client'}
                          onChange={(e) => registerUser('role').onChange(e)}
                        >
                          <Input.Radio value="client">客户</Input.Radio>
                          <Input.Radio value="professional">专业人士</Input.Radio>
                        </Input.Radio.Group>
                      </Form.Item>
                    </Input.Group>
                  </Form.Item>
                  {watchRole === 'professional' && (
                    <Form.Item label="资质证明">
                      <Input.TextArea
                        placeholder="请上传或填写您的资质证明信息"
                        rows={3}
                        {...registerUser('verification_docs' as any)}
                      />
                    </Form.Item>
                  )}
                  <Form.Item>
                    <Button type="primary" htmlType="submit" loading={loading} block>
                      注册
                    </Button>
                  </Form.Item>
                  <div style={{ textAlign: 'center' }}>
                    <span>已有账号？</span>
                    <Link to="/login">立即登录</Link>
                  </div>
                </Form>
              ),
            },
          ]}
        />
      </div>
    </div>
  )
}

export function Register() {
  return <Login />
}
