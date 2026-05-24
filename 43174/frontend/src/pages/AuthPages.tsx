import React, { useState } from 'react';
import { Form, Input, Button, Card, Typography, message, Tabs } from 'antd';
import { UserOutlined, LockOutlined, MailOutlined, PhoneOutlined } from '@ant-design/icons';
import { useNavigate, Link } from 'react-router-dom';
import { authApi } from '../services/api';
import { useAuthStore } from '../context/authStore';

const { Title, Text } = Typography;

export const LoginPage: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
  const { login } = useAuthStore();

  const onFinish = async (values: { username: string; password: string }) => {
    setLoading(true);
    try {
      const response: any = await authApi.login(values);
      const { token, user } = response.data;
      localStorage.setItem('token', token);
      login(token, user);
      message.success('登录成功');
      navigate('/');
    } catch (error: any) {
      message.error(error.message || '登录失败');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center gradient-bg p-4">
      <Card className="w-full max-w-md shadow-xl">
        <div className="text-center mb-8">
          <Title level={2} className="!mb-2">
            欢迎回来
          </Title>
          <Text type="secondary">登录您的账户继续使用</Text>
        </div>
        <Form
          name="login"
          onFinish={onFinish}
          autoComplete="off"
          layout="vertical"
        >
          <Form.Item
            name="username"
            rules={[{ required: true, message: '请输入用户名' }]}
          >
            <Input prefix={<UserOutlined />} placeholder="用户名" size="large" />
          </Form.Item>

          <Form.Item
            name="password"
            rules={[{ required: true, message: '请输入密码' }]}
          >
            <Input.Password prefix={<LockOutlined />} placeholder="密码" size="large" />
          </Form.Item>

          <Form.Item>
            <Button
              type="primary"
              htmlType="submit"
              loading={loading}
              block
              size="large"
            >
              登录
            </Button>
          </Form.Item>
        </Form>
        <div className="text-center">
          <Text>还没有账户？ </Text>
          <Link to="/register" className="text-blue-500 hover:text-blue-600">
            立即注册
          </Link>
        </div>
      </Card>
    </div>
  );
};

export const RegisterPage: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [activeTab, setActiveTab] = useState('student');
  const navigate = useNavigate();

  const onFinish = async (values: any) => {
    setLoading(true);
    try {
      await authApi.register({
        ...values,
        role: activeTab,
      });
      message.success('注册成功，请等待审核');
      navigate('/login');
    } catch (error: any) {
      message.error(error.message || '注册失败');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center gradient-bg p-4 py-8">
      <Card className="w-full max-w-md shadow-xl">
        <div className="text-center mb-6">
          <Title level={2} className="!mb-2">
            创建账户
          </Title>
          <Text type="secondary">加入我们的校园教材交易平台</Text>
        </div>

        <Tabs
          activeKey={activeTab}
          onChange={setActiveTab}
          centered
          className="mb-4"
          items={[
            { key: 'student', label: '学生' },
            { key: 'merchant', label: '书商' },
          ]}
        />

        <Form
          name="register"
          onFinish={onFinish}
          autoComplete="off"
          layout="vertical"
        >
          <Form.Item
            name="username"
            label="用户名"
            rules={[
              { required: true, message: '请输入用户名' },
              { min: 3, message: '用户名至少3个字符' },
            ]}
          >
            <Input prefix={<UserOutlined />} placeholder="请输入用户名" />
          </Form.Item>

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
              { min: 6, message: '密码至少6个字符' },
            ]}
          >
            <Input.Password prefix={<LockOutlined />} placeholder="请输入密码" />
          </Form.Item>

          <Form.Item name="phone" label="手机号">
            <Input prefix={<PhoneOutlined />} placeholder="请输入手机号（选填）" />
          </Form.Item>

          {activeTab === 'student' && (
            <>
              <Form.Item name="real_name" label="真实姓名">
                <Input placeholder="请输入真实姓名" />
              </Form.Item>
              <Form.Item name="school_name" label="学校名称">
                <Input placeholder="请输入学校名称" />
              </Form.Item>
              <Form.Item name="student_id" label="学号">
                <Input placeholder="请输入学号" />
              </Form.Item>
            </>
          )}

          {activeTab === 'merchant' && (
            <Form.Item
              name="business_license"
              label="营业执照"
              rules={[{ required: true, message: '请上传营业执照' }]}
            >
              <Input placeholder="请上传营业执照图片URL" />
            </Form.Item>
          )}

          <Form.Item>
            <Button
              type="primary"
              htmlType="submit"
              loading={loading}
              block
              size="large"
            >
              注册
            </Button>
          </Form.Item>
        </Form>
        <div className="text-center">
          <Text>已有账户？ </Text>
          <Link to="/login" className="text-blue-500 hover:text-blue-600">
            立即登录
          </Link>
        </div>
      </Card>
    </div>
  );
};
