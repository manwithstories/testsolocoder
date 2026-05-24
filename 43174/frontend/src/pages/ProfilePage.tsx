import React, { useState, useEffect } from 'react';
import { Card, Avatar, Button, Form, Input, Upload, message, Tabs, Tag, Rate, List, Empty, Modal, Statistic, Row, Col } from 'antd';
import { UserOutlined, EditOutlined, CameraOutlined } from '@ant-design/icons';
import { userApi, textbookApi, orderApi, noteApi } from '../services/api';
import { User, Textbook, Order, Note } from '../types';
import { Loading } from '../components/Loading';
import { useAuthStore } from '../context/authStore';

export const ProfilePage: React.FC = () => {
  const [loading, setLoading] = useState(true);
  const [user, setUser] = useState<User | null>(null);
  const [textbooks, setTextbooks] = useState<Textbook[]>([]);
  const [orders, setOrders] = useState<Order[]>([]);
  const [notes, setNotes] = useState<Note[]>([]);
  const [editModalVisible, setEditModalVisible] = useState(false);
  const [passwordModalVisible, setPasswordModalVisible] = useState(false);
  const [editForm] = Form.useForm();
  const [passwordForm] = Form.useForm();
  const { user: authUser, updateUser } = useAuthStore();

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    setLoading(true);
    try {
      const [userRes, textbooksRes, ordersRes, notesRes]: any = await Promise.all([
        userApi.getProfile(),
        textbookApi.getBySeller(authUser?.id || '', { page: 1, page_size: 5 }),
        orderApi.getMyOrders({ page: 1, page_size: 5 }),
        noteApi.getByUploader(authUser?.id || '', { page: 1, page_size: 5 }),
      ]);
      setUser(userRes.data);
      setTextbooks(textbooksRes.data || []);
      setOrders(ordersRes.data || []);
      setNotes(notesRes.data || []);
    } catch (error) {
      console.error('Failed to load profile data:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleAvatarUpload = async (file: File) => {
    const formData = new FormData();
    formData.append('file', file);
    try {
      const response: any = await userApi.uploadAvatar(formData);
      const newAvatar = response.data?.avatar_url;
      if (newAvatar && user) {
        setUser({ ...user, avatar: newAvatar });
        updateUser({ avatar: newAvatar });
        message.success('头像更新成功');
      }
    } catch (error: any) {
      message.error(error.message || '上传失败');
    }
  };

  const handleEditProfile = async (values: any) => {
    try {
      await userApi.updateProfile(values);
      message.success('个人信息更新成功');
      setEditModalVisible(false);
      loadData();
    } catch (error: any) {
      message.error(error.message || '更新失败');
    }
  };

  const handleChangePassword = async (values: any) => {
    try {
      await userApi.changePassword(values);
      message.success('密码修改成功');
      setPasswordModalVisible(false);
      passwordForm.resetFields();
    } catch (error: any) {
      message.error(error.message || '修改失败');
    }
  };

  if (loading) return <Loading />;
  if (!user) return <div className="text-center py-16">请先登录</div>;

  return (
    <div className="max-w-6xl mx-auto px-4 py-8">
      <Card className="mb-6">
        <div className="flex items-start gap-6">
          <div className="relative">
            <Avatar
              size={120}
              src={user.avatar}
              icon={<UserOutlined />}
              className="border-4 border-white shadow-lg"
            />
            <Upload
              showUploadList={false}
              accept="image/*"
              customRequest={async ({ file }: any) => {
                await handleAvatarUpload(file);
              }}
              className="absolute bottom-0 right-0"
            >
              <Button
                size="small"
                icon={<CameraOutlined />}
                className="rounded-full"
              />
            </Upload>
          </div>
          <div className="flex-1">
            <div className="flex items-center gap-3 mb-2">
              <h1 className="text-2xl font-bold">{user.username}</h1>
              <Tag color={user.role === 'admin' ? 'purple' : user.role === 'merchant' ? 'blue' : 'green'}>
                {user.role === 'admin' ? '管理员' : user.role === 'merchant' ? '书商' : '学生'}
              </Tag>
              <Tag color={user.status === 'active' ? 'success' : user.status === 'pending' ? 'orange' : 'default'}>
                {user.status === 'active' ? '已认证' : user.status === 'pending' ? '待审核' : '已禁用'}
              </Tag>
            </div>
            <p className="text-gray-500 mb-2">{user.email}</p>
            {user.school_name && <p className="text-gray-500 mb-2">{user.school_name}</p>}
            <div className="flex items-center gap-2 mb-4">
              <Rate disabled allowHalf value={user.rating} />
              <span className="text-gray-500">({user.rating_count} 评价)</span>
            </div>
            <div className="flex gap-2">
              <Button icon={<EditOutlined />} onClick={() => {
                editForm.setFieldsValue({
                  phone: user.phone,
                  real_name: user.real_name,
                  school_name: user.school_name,
                });
                setEditModalVisible(true);
              }}>
                编辑资料
              </Button>
              <Button onClick={() => setPasswordModalVisible(true)}>
                修改密码
              </Button>
            </div>
          </div>
        </div>

        <Row gutter={16} className="mt-6">
          <Col span={6}>
            <Card>
              <Statistic title="发布教材" value={textbooks.length} />
            </Card>
          </Col>
          <Col span={6}>
            <Card>
              <Statistic title="订单数" value={orders.length} />
            </Card>
          </Col>
          <Col span={6}>
            <Card>
              <Statistic title="上传笔记" value={notes.length} />
            </Card>
          </Col>
          <Col span={6}>
            <Card>
              <Statistic title="评分" value={user.rating} precision={1} suffix="/ 5" />
            </Card>
          </Col>
        </Row>
      </Card>

      <Card>
        <Tabs
          items={[
            {
              key: 'textbooks',
              label: '我的教材',
              children: textbooks.length > 0 ? (
                <List
                  grid={{ gutter: 16, column: 4 }}
                  dataSource={textbooks}
                  renderItem={(item) => (
                    <List.Item>
                      <Card
                        hoverable
                        cover={
                          <div className="h-32 bg-gray-100 flex items-center justify-center">
                            {item.cover_image ? (
                              <img src={item.cover_image} alt={item.title} className="h-full object-contain" />
                            ) : (
                              <span className="text-4xl">📚</span>
                            )}
                          </div>
                        }
                      >
                        <Card.Meta
                          title={<span className="text-sm truncate block">{item.title}</span>}
                          description={<span className="text-red-500">¥{item.price}</span>}
                        />
                      </Card>
                    </List.Item>
                  )}
                />
              ) : (
                <Empty description="暂无发布的教材" />
              ),
            },
            {
              key: 'notes',
              label: '我的笔记',
              children: notes.length > 0 ? (
                <List
                  grid={{ gutter: 16, column: 4 }}
                  dataSource={notes}
                  renderItem={(item) => (
                    <List.Item>
                      <Card hoverable>
                        <Card.Meta
                          title={<span className="text-sm truncate block">{item.title}</span>}
                          description={<span className="text-gray-500">{item.subject}</span>}
                        />
                      </Card>
                    </List.Item>
                  )}
                />
              ) : (
                <Empty description="暂无上传的笔记" />
              ),
            },
            {
              key: 'orders',
              label: '我的订单',
              children: orders.length > 0 ? (
                <List
                  dataSource={orders}
                  renderItem={(item) => (
                    <List.Item>
                      <div className="flex justify-between items-center w-full">
                        <div>
                          <p className="font-medium">{item.order_no}</p>
                          <p className="text-gray-500 text-sm">
                            {new Date(item.created_at).toLocaleDateString()}
                          </p>
                        </div>
                        <div className="text-right">
                          <p className="text-red-500 font-bold">¥{item.total_amount}</p>
                          <Tag>{item.status}</Tag>
                        </div>
                      </div>
                    </List.Item>
                  )}
                />
              ) : (
                <Empty description="暂无订单" />
              ),
            },
          ]}
        />
      </Card>

      <Modal
        title="编辑个人资料"
        open={editModalVisible}
        onCancel={() => setEditModalVisible(false)}
        footer={null}
      >
        <Form form={editForm} layout="vertical" onFinish={handleEditProfile}>
          <Form.Item name="phone" label="手机号">
            <Input placeholder="请输入手机号" />
          </Form.Item>
          <Form.Item name="real_name" label="真实姓名">
            <Input placeholder="请输入真实姓名" />
          </Form.Item>
          {user.role === 'student' && (
            <Form.Item name="school_name" label="学校名称">
              <Input placeholder="请输入学校名称" />
            </Form.Item>
          )}
          <Form.Item>
            <Button type="primary" htmlType="submit" block>
              保存
            </Button>
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title="修改密码"
        open={passwordModalVisible}
        onCancel={() => {
          setPasswordModalVisible(false);
          passwordForm.resetFields();
        }}
        footer={null}
      >
        <Form form={passwordForm} layout="vertical" onFinish={handleChangePassword}>
          <Form.Item
            name="old_password"
            label="当前密码"
            rules={[{ required: true, message: '请输入当前密码' }]}
          >
            <Input.Password placeholder="请输入当前密码" />
          </Form.Item>
          <Form.Item
            name="new_password"
            label="新密码"
            rules={[
              { required: true, message: '请输入新密码' },
              { min: 6, message: '密码至少6个字符' },
            ]}
          >
            <Input.Password placeholder="请输入新密码" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block>
              确认修改
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};
