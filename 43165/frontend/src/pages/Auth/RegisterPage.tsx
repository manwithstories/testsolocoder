import { useState } from 'react';
import { Form, Input, Button, Card, Typography, Select, message } from 'antd';
import { UserOutlined, LockOutlined, MailOutlined, PhoneOutlined, IdcardOutlined, BankOutlined, EnvironmentOutlined } from '@ant-design/icons';
import { useNavigate, Link } from 'react-router-dom';
import { authApi } from '../../services/api';
import { RegisterRequest } from '../../types';

const { Title, Text } = Typography;
const { Option } = Select;

export const RegisterPage = () => {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);

  const handleRegister = async (values: RegisterRequest) => {
    setLoading(true);
    try {
      const response = await authApi.register(values);
      if (response.data.code === 201) {
        message.success('注册成功');
        navigate('/login');
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
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100 py-8">
      <Card className="w-full max-w-md shadow-lg" style={{ borderRadius: '12px' }}>
        <div className="text-center mb-8">
          <Title level={2} style={{ marginBottom: '8px' }}>用户注册</Title>
          <Text type="secondary">创建您的账号</Text>
        </div>

        <Form
          name="register"
          onFinish={handleRegister}
          layout="vertical"
          size="large"
        >
          <Form.Item
            name="username"
            label="用户名"
            rules={[{ required: true, message: '请输入用户名' }, { min: 3, message: '用户名至少3个字符' }]}
          >
            <Input prefix={<UserOutlined />} placeholder="请输入用户名" />
          </Form.Item>
          <Form.Item
            name="email"
            label="邮箱"
            rules={[{ required: true, message: '请输入邮箱' }, { type: 'email', message: '请输入有效的邮箱' }]}
          >
            <Input prefix={<MailOutlined />} placeholder="请输入邮箱" />
          </Form.Item>
          <Form.Item
            name="phone"
            label="手机号"
            rules={[{ pattern: /^1[3-9]\d{9}$/, message: '请输入有效的手机号' }]}
          >
            <Input prefix={<PhoneOutlined />} placeholder="请输入手机号" />
          </Form.Item>
          <Form.Item
            name="password"
            label="密码"
            rules={[
              { required: true, message: '请输入密码' },
              { min: 8, message: '密码至少8个字符' },
            ]}
          >
            <Input.Password prefix={<LockOutlined />} placeholder="请输入密码" />
          </Form.Item>
          <Form.Item
            name="real_name"
            label="真实姓名"
            rules={[{ required: true, message: '请输入真实姓名' }]}
          >
            <Input prefix={<IdcardOutlined />} placeholder="请输入真实姓名" />
          </Form.Item>
          <Form.Item
            name="role"
            label="用户角色"
            rules={[{ required: true, message: '请选择用户角色' }]}
          >
            <Select placeholder="请选择角色">
              <Option value="employer">雇主</Option>
              <Option value="agent">中介</Option>
              <Option value="temporary">临时工</Option>
            </Select>
          </Form.Item>
          <Form.Item
            name="company"
            label="公司名称"
          >
            <Input prefix={<BankOutlined />} placeholder="请输入公司名称（选填）" />
          </Form.Item>
          <Form.Item
            name="address"
            label="地址"
          >
            <Input prefix={<EnvironmentOutlined />} placeholder="请输入地址（选填）" />
          </Form.Item>
          <Form.Item>
            <Button
              type="primary"
              htmlType="submit"
              loading={loading}
              block
            >
              注册
            </Button>
          </Form.Item>
          <div className="text-center">
            <Text type="secondary">已有账号？</Text>
            <Link to="/login" className="ml-2">立即登录</Link>
          </div>
        </Form>
      </Card>
    </div>
  );
};
