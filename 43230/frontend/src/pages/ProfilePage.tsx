import React, { useEffect, useState } from 'react';
import {
  Card,
  Form,
  Input,
  Button,
  Space,
  Typography,
  Avatar,
  Row,
  Col,
  Statistic,
  Descriptions,
  Tag,
  Tabs,
  Upload,
  message,
  Divider,
  Select,
} from 'antd';
import {
  UserOutlined,
  WalletOutlined,
  ShoppingCartOutlined,
  StarOutlined,
  SaveOutlined,
  CameraOutlined,
  EditOutlined,
} from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import { getProfile, updateProfile, updateDesignerProfile, updatePrinterProfile } from '@/store/authSlice';
import { authApi } from '@/services/api';
import { User } from '@/types';
import { formatPrice, getRoleText } from '@/utils/format';
import { formatDate } from '@/utils/date';

const { Title, Text, Paragraph } = Typography;
const { TextArea } = Input;
const { Option } = Select;
const { TabPane } = Tabs;

const ProfilePage: React.FC = () => {
  const dispatch = useAppDispatch();
  const { user, loading } = useAppSelector((state) => state.auth);
  const [form] = Form.useForm();
  const [designerForm] = Form.useForm();
  const [printerForm] = Form.useForm();
  const [stats, setStats] = useState<any>(null);
  const [activeTab, setActiveTab] = useState('profile');
  const navigate = useNavigate();

  useEffect(() => {
    if (!user) {
      dispatch(getProfile() as any);
    }
    fetchStats();
  }, [dispatch, user]);

  const fetchStats = async () => {
    try {
      const response = await authApi.getUserStats();
      setStats(response.data);
    } catch (error) {
      console.error('Failed to fetch stats:', error);
    }
  };

  useEffect(() => {
    if (user) {
      form.setFieldsValue(user);
      if (user.designer_profile) {
        designerForm.setFieldsValue(user.designer_profile);
      }
      if (user.printer_profile) {
        printerForm.setFieldsValue(user.printer_profile);
      }
    }
  }, [user, form, designerForm, printerForm]);

  const handleUpdateProfile = async () => {
    try {
      const values = await form.validateFields();
      await dispatch(updateProfile(values) as any);
      message.success('个人信息更新成功');
    } catch (error: any) {
      message.error(error.message || '更新失败');
    }
  };

  const handleUpdateDesignerProfile = async () => {
    try {
      const values = await designerForm.validateFields();
      await dispatch(updateDesignerProfile(values) as any);
      message.success('设计师信息更新成功');
    } catch (error: any) {
      message.error(error.message || '更新失败');
    }
  };

  const handleUpdatePrinterProfile = async () => {
    try {
      const values = await printerForm.validateFields();
      await dispatch(updatePrinterProfile(values) as any);
      message.success('打印商信息更新成功');
    } catch (error: any) {
      message.error(error.message || '更新失败');
    }
  };

  if (!user && loading) {
    return (
      <Card>
        <div className="text-center py-12">
          <div className="animate-spin w-8 h-8 border-4 border-blue-500 border-t-transparent rounded-full mx-auto"></div>
        </div>
      </Card>
    );
  }

  return (
    <div className="space-y-6">
      {/* 用户信息卡片 */}
      <Card>
        <Row gutter={24} align="middle">
          <Col xs={24} md={6} className="text-center">
            <div className="mb-4">
              <Avatar size={120} src={user?.avatar} icon={<UserOutlined />} />
              <div className="mt-3">
                <Upload
                  showUploadList={false}
                  beforeUpload={() => {
                    message.info('头像上传功能待实现');
                    return false;
                  }}
                >
                  <Button type="link" icon={<CameraOutlined />}>
                    更换头像
                  </Button>
                </Upload>
              </div>
            </div>
          </Col>
          <Col xs={24} md={18}>
            <Space direction="vertical" size="small" className="w-full">
              <Space size="middle">
                <Title level={3} className="!mb-0">
                  {user?.username}
                </Title>
                <Tag color="blue">{getRoleText(user?.role || '')}</Tag>
                {user?.email_verified && <Tag color="green">已认证</Tag>}
              </Space>
              <Text type="secondary">{user?.email}</Text>
              <Space size="large" wrap>
                <Text>
                  <UserOutlined className="mr-1" />
                  {user?.real_name || '未设置'}
                </Text>
                <Text>
                  <WalletOutlined className="mr-1" />
                  余额: <span className="text-red-500 font-medium">{formatPrice(user?.balance || 0)}</span>
                </Text>
                <Text>
                  <StarOutlined className="mr-1 text-yellow-500" />
                  信用分: {user?.credit_score?.toFixed(1)}
                </Text>
                <Text type="secondary">
                  注册于 {formatDate(user?.created_at || '')}
                </Text>
              </Space>
            </Space>
          </Col>
        </Row>

        <Divider />

        <Row gutter={16}>
          <Col xs={12} md={6}>
            <Card size="small">
              <Statistic
                title="订单总数"
                value={stats?.order_count || 0}
                prefix={<ShoppingCartOutlined />}
              />
            </Card>
          </Col>
          <Col xs={12} md={6}>
            <Card size="small">
              <Statistic
                title="模型总数"
                value={stats?.model_count || 0}
                prefix={<EditOutlined />}
              />
            </Card>
          </Col>
          <Col xs={12} md={6}>
            <Card size="small">
              <Statistic
                title="总支出"
                value={stats?.total_spent || 0}
                precision={2}
                prefix="¥"
                valueStyle={{ color: '#cf1322' }}
              />
            </Card>
          </Col>
          <Col xs={12} md={6}>
            <Card size="small">
              <Statistic
                title="总收入"
                value={stats?.total_earned || 0}
                precision={2}
                prefix="¥"
                valueStyle={{ color: '#3f8600' }}
              />
            </Card>
          </Col>
        </Row>
      </Card>

      <Card>
        <Tabs defaultActiveKey="profile" activeKey={activeTab} onChange={setActiveTab}>
          <TabPane tab="基本信息" key="profile">
            <Form form={form} layout="vertical" className="max-w-2xl">
              <Row gutter={16}>
                <Col xs={24} md={12}>
                  <Form.Item name="username" label="用户名">
                    <Input />
                  </Form.Item>
                </Col>
                <Col xs={24} md={12}>
                  <Form.Item name="email" label="邮箱">
                    <Input disabled />
                  </Form.Item>
                </Col>
                <Col xs={24} md={12}>
                  <Form.Item name="phone" label="手机号">
                    <Input placeholder="请输入手机号" />
                  </Form.Item>
                </Col>
                <Col xs={24} md={12}>
                  <Form.Item name="real_name" label="真实姓名">
                    <Input placeholder="请输入真实姓名" />
                  </Form.Item>
                </Col>
              </Row>
              <Button type="primary" onClick={handleUpdateProfile} icon={<SaveOutlined />}>
                保存修改
              </Button>
            </Form>
          </TabPane>

          {user?.role === 'designer' && (
            <TabPane tab="设计师信息" key="designer">
              <Form form={designerForm} layout="vertical" className="max-w-2xl">
                <Row gutter={16}>
                  <Col xs={24} md={12}>
                    <Form.Item name="nickname" label="昵称">
                      <Input placeholder="请输入展示昵称" />
                    </Form.Item>
                  </Col>
                  <Col xs={24} md={12}>
                    <Form.Item name="experience_years" label="从业年限">
                      <InputNumber className="w-full" min={0} placeholder="请输入从业年限" />
                    </Form.Item>
                  </Col>
                  <Col xs={24}>
                    <Form.Item name="bio" label="个人简介">
                      <TextArea rows={3} placeholder="介绍一下你自己吧" />
                    </Form.Item>
                  </Col>
                  <Col xs={24}>
                    <Form.Item name="portfolio_url" label="作品集链接">
                      <Input placeholder="https://..." />
                    </Form.Item>
                  </Col>
                  <Col xs={24}>
                    <Form.Item name="specialties" label="擅长领域">
                      <Select mode="tags" placeholder="输入擅长领域后按回车添加">
                        {['建筑', '人物', '机械', '艺术', '玩具'].map((s) => (
                          <Option key={s} value={s}>
                            {s}
                          </Option>
                        ))}
                      </Select>
                    </Form.Item>
                  </Col>
                </Row>
                <Button type="primary" onClick={handleUpdateDesignerProfile} icon={<SaveOutlined />}>
                  保存修改
                </Button>
              </Form>
            </TabPane>
          )}

          {user?.role === 'printer' && (
            <TabPane tab="打印商信息" key="printer">
              <Form form={printerForm} layout="vertical" className="max-w-2xl">
                <Row gutter={16}>
                  <Col xs={24}>
                    <Form.Item name="company_name" label="公司名称">
                      <Input placeholder="请输入公司名称" />
                    </Form.Item>
                  </Col>
                  <Col xs={24}>
                    <Form.Item name="address" label="公司地址">
                      <TextArea rows={2} placeholder="请输入详细地址" />
                    </Form.Item>
                  </Col>
                  <Col xs={24}>
                    <Form.Item name="max_print_size" label="最大打印尺寸">
                      <Input placeholder="如：300x300x400mm" />
                    </Form.Item>
                  </Col>
                  <Col xs={24}>
                    <Form.Item name="supported_materials" label="支持材料">
                      <Select mode="tags" placeholder="输入支持的材料后按回车添加">
                        {['PLA', 'ABS', 'PETG', 'TPU', '树脂'].map((m) => (
                          <Option key={m} value={m.toLowerCase()}>
                            {m}
                          </Option>
                        ))}
                      </Select>
                    </Form.Item>
                  </Col>
                </Row>
                <Button type="primary" onClick={handleUpdatePrinterProfile} icon={<SaveOutlined />}>
                  保存修改
                </Button>
              </Form>
            </TabPane>
          )}

          <TabPane tab="账号设置" key="settings">
            <Descriptions column={1} size="small" bordered>
              <Descriptions.Item label="账号状态">
                <Tag color="green">正常</Tag>
              </Descriptions.Item>
              <Descriptions.Item label="邮箱验证">
                {user?.email_verified ? <Tag color="green">已验证</Tag> : <Tag color="orange">未验证</Tag>}
              </Descriptions.Item>
              <Descriptions.Item label="实名认证">
                {user?.id_card_verified ? <Tag color="green">已认证</Tag> : <Tag color="orange">未认证</Tag>}
              </Descriptions.Item>
              <Descriptions.Item label="上次登录">
                {user?.last_login_at ? formatDate(user.last_login_at) : '首次登录'}
                {user?.last_login_ip && ` (${user.last_login_ip})`}
              </Descriptions.Item>
            </Descriptions>
            <Divider />
            <Space direction="vertical" className="w-full">
              <Button type="primary" onClick={() => navigate('/wallet')}>
                <WalletOutlined /> 充值/提现
              </Button>
              <Button danger>
                修改密码
              </Button>
            </Space>
          </TabPane>
        </Tabs>
      </Card>
    </div>
  );
};

export default ProfilePage;
