import { useState } from 'react';
import { Form, Input, Button, Card, Typography, message, Tabs } from 'antd';
import { UserOutlined, LockOutlined, MailOutlined, PhoneOutlined } from '@ant-design/icons';
import { useNavigate, Link } from 'react-router-dom';
import { authApi } from '../../services/api';
import { useAuthStore } from '../../context/AuthContext';
import { LoginRequest, RegisterRequest } from '../../types';

const { Title, Text } = Typography;

export const LoginPage = () => {
  const navigate = useNavigate();
  const { login } = useAuthStore();
  const [loading, setLoading] = useState(false);

  const handleLogin = async (values: LoginRequest) => {
    setLoading(true);
    try {
      const response = await authApi.login(values);
      if (response.data.code === 200) {
        login(response.data.data.token, response.data.data.user);
        message.success('登录成功');
        navigate('/dashboard');
      } else {
        message.error(response.data.message || '登录失败');
      }
    } catch (error: any) {
      message.error(error.response?.data?.message || '登录失败');
    } finally {
      setLoading(false);
    }
  };

  const handleRegister = async (values: RegisterRequest) => {
    setLoading(true);
    try {
      const response = await authApi.register(values);
      if (response.data.code === 201) {
        message.success('注册成功，请登录');
      } else {
        message.error(response.data.message || '注册失败');
      }
    } catch (error: any) {
      message.error(error.response?.data?.message || '注册失败');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100">
      <Card className="w-full max-w-md shadow-lg" style={{ borderRadius: '12px' }}>
        <div className="text-center mb-8">
          <Title level={2} style={{ marginBottom: '8px' }}>临时工招聘平台</Title>
          <Text type="secondary">活动临时工招聘与排班管理系统</Text>
        </div>

        <Tabs
          defaultActiveKey="login"
          centered
          items={[
            {
              key: 'login',
              label: '登录',
              children: (
                <Form
                  name="login"
                  onFinish={handleLogin}
                  layout="vertical"
                  size="large"
                >
                  <Form.Item
                    name="login"
                    rules={[{ required: true, message: '请输入用户名/邮箱/手机号' }]}
                  >
                    <Input prefix={<UserOutlined />} placeholder="用户名/邮箱/手机号" />
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
                      loading={loading}
                      block
                    >
                      登录
                    </Button>
                  </Form.Item>
                  <div className="text-center">
                    <Text type="secondary">还没有账号？</Text>
                    <Link to="/register" className="ml-2">立即注册</Link>
                  </div>
                </Form>
              ),
            },
          ]}
        />
      </Card>
    </div>
  );
};
