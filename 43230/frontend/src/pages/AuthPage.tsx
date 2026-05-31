import React, { useEffect } from 'react';
import { Form, Input, Button, Card, Radio, Typography, Divider } from 'antd';
import { UserOutlined, LockOutlined, MailOutlined, PhoneOutlined } from '@ant-design/icons';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import { login, register } from '@/store/authSlice';
import { UserRole } from '@/types';

const { Title, Text } = Typography;

interface AuthPageProps {
  mode: 'login' | 'register';
}

const AuthPage: React.FC<AuthPageProps> = ({ mode }) => {
  const [form] = Form.useForm();
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const location = useLocation();
  const { loading, isAuthenticated } = useAppSelector((state) => state.auth);

  const from = (location.state as any)?.from?.pathname || '/';

  useEffect(() => {
    if (isAuthenticated) {
      navigate(from, { replace: true });
    }
  }, [isAuthenticated, navigate, from]);

  const onFinish = async (values: any) => {
    if (mode === 'login') {
      await dispatch(login(values));
    } else {
      await dispatch(register(values));
      navigate('/login');
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center p-4">
      <Card className="w-full max-w-md shadow-xl">
        <div className="text-center mb-8">
          <Title level={2} className="!mb-2">
            {mode === 'login' ? '欢迎回来' : '创建账号'}
          </Title>
          <Text type="secondary">
            {mode === 'login'
              ? '登录您的3D打印平台账号'
              : '加入我们，开启3D打印之旅'}
          </Text>
        </div>

        <Form form={form} layout="vertical" onFinish={onFinish} size="large">
          {mode === 'register' && (
            <>
              <Form.Item
                name="username"
                label="用户名"
                rules={[
                  { required: true, message: '请输入用户名' },
                  { min: 3, max: 50, message: '用户名长度为3-50个字符' },
                ]}
              >
                <Input prefix={<UserOutlined />} placeholder="请输入用户名" />
              </Form.Item>

              <Form.Item
                name="role"
                label="用户角色"
                rules={[{ required: true, message: '请选择用户角色' }]}
                initialValue="customer"
              >
                <Radio.Group className="w-full">
                  <Radio.Button value="customer" className="flex-1 text-center">
                    客户
                  </Radio.Button>
                  <Radio.Button value="designer" className="flex-1 text-center">
                    建模师
                  </Radio.Button>
                  <Radio.Button value="printer" className="flex-1 text-center">
                    打印商
                  </Radio.Button>
                </Radio.Group>
              </Form.Item>
            </>
          )}

          <Form.Item
            name="email"
            label="邮箱"
            rules={[
              { required: true, message: '请输入邮箱' },
              { type: 'email', message: '请输入有效的邮箱地址' },
            ]}
          >
            <Input prefix={<MailOutlined />} placeholder="请输入邮箱" />
          </Form.Item>

          <Form.Item
            name="password"
            label="密码"
            rules={[
              { required: true, message: '请输入密码' },
              { min: 6, max: 128, message: '密码长度为6-128个字符' },
            ]}
          >
            <Input.Password prefix={<LockOutlined />} placeholder="请输入密码" />
          </Form.Item>

          {mode === 'register' && (
            <>
              <Form.Item
                name="phone"
                label="手机号"
                rules={[{ pattern: /^1[3-9]\d{9}$/, message: '请输入有效的手机号' }]}
              >
                <Input prefix={<PhoneOutlined />} placeholder="请输入手机号" />
              </Form.Item>

              <Form.Item
                noStyle
                shouldUpdate={(prevValues, curValues) => prevValues.role !== curValues.role}
              >
                {({ getFieldValue }) => {
                  const role = getFieldValue('role') as UserRole;
                  if (role === 'printer') {
                    return (
                      <Form.Item
                        name="company_name"
                        label="公司名称"
                        rules={[{ required: true, message: '请输入公司名称' }]}
                      >
                        <Input placeholder="请输入公司名称" />
                      </Form.Item>
                    );
                  }
                  if (role === 'designer' || role === 'customer') {
                    return (
                      <Form.Item name="real_name" label="真实姓名">
                        <Input placeholder="请输入真实姓名" />
                      </Form.Item>
                    );
                  }
                  return null;
                }}
              </Form.Item>
            </>
          )}

          <Form.Item>
            <Button
              type="primary"
              htmlType="submit"
              block
              size="large"
              loading={loading}
            >
              {mode === 'login' ? '登录' : '注册'}
            </Button>
          </Form.Item>
        </Form>

        <Divider plain className="my-4">
          <Text type="secondary">
            {mode === 'login' ? '还没有账号？' : '已有账号？'}
          </Text>
        </Divider>

        <div className="text-center">
          <Link to={mode === 'login' ? '/register' : '/login'}>
            <Button type="link">
              {mode === 'login' ? '立即注册' : '立即登录'}
            </Button>
          </Link>
        </div>
      </Card>
    </div>
  );
};

export default AuthPage;
