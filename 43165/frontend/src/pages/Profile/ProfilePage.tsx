import { useState, useEffect } from 'react';
import { Card, Form, Input, Button, Avatar, message, Row, Col, Descriptions, Tag, Progress, Modal, Divider } from 'antd';
import { UserOutlined, MailOutlined, PhoneOutlined, BankOutlined, EnvironmentOutlined } from '@ant-design/icons';
import { authApi, evaluationApi } from '../../services/api';
import { useAuthStore } from '../../context/AuthContext';
import { User } from '../../types';
import dayjs from 'dayjs';

const { TextArea } = Input;

export const ProfilePage = () => {
  const { user, setUser, logout } = useAuthStore();
  const [loading, setLoading] = useState(false);
  const [stats, setStats] = useState<any>(null);
  const [passwordModalVisible, setPasswordModalVisible] = useState(false);
  const [profileForm] = Form.useForm();
  const [passwordForm] = Form.useForm();

  useEffect(() => {
    if (user) {
      profileForm.setFieldsValue({
        real_name: user.real_name,
        phone: user.phone,
        company: user.company,
        address: user.address,
      });
      fetchStats();
    }
  }, [user]);

  const fetchStats = async () => {
    if (user?.role === 'temporary') {
      try {
        const res = await evaluationApi.getEvaluationStats(user.id);
        if (res.data.code === 200) {
          setStats(res.data.data);
        }
      } catch (error) {
        console.error('获取统计数据失败');
      }
    }
  };

  const handleUpdateProfile = async (values: any) => {
    setLoading(true);
    try {
      const res = await authApi.updateProfile(values);
      if (res.data.code === 200) {
        message.success('更新成功');
        setUser(res.data.data as User);
      }
    } catch (error: any) {
      message.error(error.response?.data?.message || '更新失败');
    } finally {
      setLoading(false);
    }
  };

  const handleChangePassword = async (values: any) => {
    setLoading(true);
    try {
      const res = await authApi.changePassword(values);
      if (res.data.code === 200) {
        message.success('密码修改成功');
        setPasswordModalVisible(false);
        passwordForm.resetFields();
      }
    } catch (error: any) {
      message.error(error.response?.data?.message || '修改失败');
    } finally {
      setLoading(false);
    }
  };

  if (!user) {
    return <div>用户信息未加载</div>;
  }

  const roleMap: Record<string, { label: string; color: string }> = {
    employer: { label: '雇主', color: 'blue' },
    agent: { label: '中介', color: 'green' },
    temporary: { label: '临时工', color: 'orange' },
  };

  return (
    <div>
      <h2 className="text-2xl font-bold mb-6">个人中心</h2>

      <Row gutter={[16, 16]}>
        <Col xs={24} lg={8}>
          <Card className="text-center">
            <Avatar
              size={100}
              src={user.avatar}
              icon={!user.avatar && <UserOutlined />}
              className="mb-4"
            />
            <h3 className="text-xl font-bold">{user.real_name}</h3>
            <p className="text-gray-500">@{user.username}</p>
            <Tag color={roleMap[user.role]?.color} className="mt-2">
              {roleMap[user.role]?.label}
            </Tag>

            {user.role === 'temporary' && stats && (
              <div className="mt-6">
                <Divider />
                <div className="mb-4">
                  <div className="flex justify-between mb-2">
                    <span>信用分</span>
                    <span className="font-semibold text-green-600">{stats.credit_score}/100</span>
                  </div>
                  <Progress percent={stats.credit_score} strokeColor="#52c41a" />
                </div>
                <Descriptions column={1} size="small">
                  <Descriptions.Item label="收到评价">{stats.total_count} 次</Descriptions.Item>
                  <Descriptions.Item label="平均评分">{stats.average_rating?.toFixed(1)} / 5</Descriptions.Item>
                </Descriptions>
              </div>
            )}

            <Divider />
            <Button danger onClick={() => logout()}>退出登录</Button>
          </Card>
        </Col>

        <Col xs={24} lg={16}>
          <Card title="基本信息" className="mb-4">
            <Form
              form={profileForm}
              onFinish={handleUpdateProfile}
              layout="vertical"
            >
              <Descriptions column={2} bordered className="mb-4">
                <Descriptions.Item label="邮箱">
                  <MailOutlined className="mr-1" />{user.email}
                </Descriptions.Item>
                <Descriptions.Item label="账号状态">
                  <Tag color={user.status === 'active' ? 'green' : 'red'}>
                    {user.status === 'active' ? '正常' : '禁用'}
                  </Tag>
                </Descriptions.Item>
                <Descriptions.Item label="注册时间">
                  {dayjs(user.created_at).format('YYYY-MM-DD HH:mm')}
                </Descriptions.Item>
                <Descriptions.Item label="最后登录">
                  {user.last_login_at ? dayjs(user.last_login_at).format('YYYY-MM-DD HH:mm') : '-'}
                </Descriptions.Item>
              </Descriptions>

              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item name="real_name" label="真实姓名">
                    <Input prefix={<UserOutlined />} />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="phone" label="手机号">
                    <Input prefix={<PhoneOutlined />} />
                  </Form.Item>
                </Col>
              </Row>
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item name="company" label="公司名称">
                    <Input prefix={<BankOutlined />} />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="address" label="地址">
                    <Input prefix={<EnvironmentOutlined />} />
                  </Form.Item>
                </Col>
              </Row>

              <Form.Item>
                <Button type="primary" htmlType="submit" loading={loading}>
                  保存修改
                </Button>
                <Button className="ml-2" onClick={() => setPasswordModalVisible(true)}>
                  修改密码
                </Button>
              </Form.Item>
            </Form>
          </Card>
        </Col>
      </Row>

      <Modal
        title="修改密码"
        open={passwordModalVisible}
        onCancel={() => setPasswordModalVisible(false)}
        footer={null}
      >
        <Form form={passwordForm} onFinish={handleChangePassword} layout="vertical">
          <Form.Item
            name="old_password"
            label="当前密码"
            rules={[{ required: true, message: '请输入当前密码' }]}
          >
            <Input.Password />
          </Form.Item>
          <Form.Item
            name="new_password"
            label="新密码"
            rules={[
              { required: true, message: '请输入新密码' },
              { min: 8, message: '密码至少8个字符' },
            ]}
          >
            <Input.Password />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading} block>
              确认修改
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};
