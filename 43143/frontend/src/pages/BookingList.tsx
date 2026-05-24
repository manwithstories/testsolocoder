import React, { useEffect, useState } from 'react';
import { Card, Tabs, List, Tag, Button, Avatar, Spin, message, Modal, Form, Input, Select, Statistic, Row, Col } from 'antd';
import { CalendarOutlined, UserOutlined, CheckCircleOutlined, CloseCircleOutlined, ClockCircleOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { bookingApi } from '@/api/booking';
import { Booking, BookingStatus } from '@/types';

const { Option } = Select;

const BookingList: React.FC = () => {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState('student');
  const [bookings, setBookings] = useState<Booking[]>([]);
  const [page, setPage] = useState(1);
  const [pageSize] = useState(10);
  const [statusFilter, setStatusFilter] = useState<string | undefined>();
  const [actionModal, setActionModal] = useState<{ visible: boolean; booking: Booking | null; action: string }>({
    visible: false,
    booking: null,
    action: '',
  });
  const [actionForm] = Form.useForm();

  const fetchBookings = async () => {
    setLoading(true);
    try {
      const data = await bookingApi.list({
        role: activeTab,
        page,
        page_size: pageSize,
        status: statusFilter,
      });
      setBookings(data.items || []);
    } catch (error: any) {
      message.error(error.message || '获取预约列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchBookings();
  }, [activeTab, page, pageSize, statusFilter]);

  const handleAction = async () => {
    try {
      const values = await actionForm.validateFields();

      switch (actionModal.action) {
        case 'confirm':
          await bookingApi.confirm(actionModal.booking!.id);
          message.success('预约已确认');
          break;
        case 'reject':
          await bookingApi.reject(actionModal.booking!.id, values.reason);
          message.success('预约已拒绝');
          break;
        case 'cancel':
          await bookingApi.cancel(actionModal.booking!.id, values.reason);
          message.success('预约已取消');
          break;
        case 'complete':
          await bookingApi.complete(actionModal.booking!.id);
          message.success('课程已完成');
          break;
      }

      setActionModal({ visible: false, booking: null, action: '' });
      actionForm.resetFields();
      fetchBookings();
    } catch (error: any) {
      if (!error.errorFields) {
        message.error(error.message || '操作失败');
      }
    }
  };

  const getStatusColor = (status: BookingStatus) => {
    const colors: Record<BookingStatus, string> = {
      pending: 'orange',
      confirmed: 'blue',
      rejected: 'red',
      cancelled: 'default',
      completed: 'green',
      no_show: 'red',
    };
    return colors[status];
  };

  const getStatusText = (status: BookingStatus) => {
    const texts: Record<BookingStatus, string> = {
      pending: '待确认',
      confirmed: '已确认',
      rejected: '已拒绝',
      cancelled: '已取消',
      completed: '已完成',
      no_show: '未出席',
    };
    return texts[status];
  };

  const tabItems = [
    {
      key: 'student',
      label: '我预约的',
    },
    {
      key: 'teacher',
      label: '我的授课',
    },
  ];

  return (
    <div>
      <Row gutter={16} style={{ marginBottom: 16 }}>
        <Col span={6}>
          <Card>
            <Statistic title="待确认" value={bookings.filter((b) => b.status === 'pending').length} prefix={<ClockCircleOutlined />} />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic title="已确认" value={bookings.filter((b) => b.status === 'confirmed').length} prefix={<CheckCircleOutlined />} />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic title="已完成" value={bookings.filter((b) => b.status === 'completed').length} prefix={<CalendarOutlined />} />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic title="已取消" value={bookings.filter((b) => b.status === 'cancelled').length} prefix={<CloseCircleOutlined />} />
          </Card>
        </Col>
      </Row>

      <Card
        title={
          <div style={{ display: 'flex', alignItems: 'center', gap: 16 }}>
            <Tabs
              activeKey={activeTab}
              onChange={setActiveTab}
              items={tabItems}
              style={{ marginBottom: 0 }}
            />
            <Select
              placeholder="筛选状态"
              allowClear
              style={{ width: 120 }}
              value={statusFilter}
              onChange={(value) => {
                setStatusFilter(value);
                setPage(1);
              }}
            >
              <Option value="pending">待确认</Option>
              <Option value="confirmed">已确认</Option>
              <Option value="completed">已完成</Option>
              <Option value="cancelled">已取消</Option>
            </Select>
          </div>
        }
      >
        {loading ? (
          <div style={{ textAlign: 'center', padding: 40 }}>
            <Spin />
          </div>
        ) : bookings.length === 0 ? (
          <div style={{ textAlign: 'center', padding: 40, color: '#999' }}>
            暂无预约记录
          </div>
        ) : (
          <List
            itemLayout="horizontal"
            dataSource={bookings}
            renderItem={(item) => (
              <List.Item
                actions={[
                  <Tag key="status" color={getStatusColor(item.status)}>
                    {getStatusText(item.status)}
                  </Tag>,
                  item.status === 'pending' && activeTab === 'teacher' && (
                    <Button
                      key="confirm"
                      type="link"
                      onClick={() =>
                        setActionModal({ visible: true, booking: item, action: 'confirm' })
                      }
                    >
                      确认
                    </Button>
                  ),
                  item.status === 'pending' && activeTab === 'teacher' && (
                    <Button
                      key="reject"
                      type="link"
                      danger
                      onClick={() =>
                        setActionModal({ visible: true, booking: item, action: 'reject' })
                      }
                    >
                      拒绝
                    </Button>
                  ),
                  (item.status === 'pending' || item.status === 'confirmed') && (
                    <Button
                      key="cancel"
                      type="link"
                      danger
                      onClick={() =>
                        setActionModal({ visible: true, booking: item, action: 'cancel' })
                      }
                    >
                      取消
                    </Button>
                  ),
                  item.status === 'confirmed' && activeTab === 'teacher' && (
                    <Button
                      key="complete"
                      type="link"
                      onClick={() =>
                        setActionModal({ visible: true, booking: item, action: 'complete' })
                      }
                    >
                      完成
                    </Button>
                  ),
                  <Button
                    key="detail"
                    type="link"
                    onClick={() => navigate(`/bookings/${item.id}`)}
                  >
                    详情
                  </Button>,
                ]}
              >
                <List.Item.Meta
                  avatar={<Avatar icon={<UserOutlined />} src={activeTab === 'student' ? item.teacher?.avatar : item.student?.avatar} />}
                  title={
                    <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                      <span>{item.posting?.title}</span>
                    </div>
                  }
                  description={
                    <div>
                      <div>
                        {activeTab === 'student' ? item.teacher?.nickname : item.student?.nickname}
                      </div>
                      <div style={{ color: '#666', marginTop: 4 }}>
                        <CalendarOutlined /> {new Date(item.scheduled_start).toLocaleString()} -{' '}
                        {new Date(item.scheduled_end).toLocaleTimeString()}
                      </div>
                      <div style={{ color: '#1890ff', fontWeight: 'bold', marginTop: 4 }}>
                        ¥{item.price.toFixed(2)}
                      </div>
                    </div>
                  }
                />
              </List.Item>
            )}
          />
        )}
      </Card>

      <Modal
        title={
          actionModal.action === 'confirm'
            ? '确认预约'
            : actionModal.action === 'reject'
            ? '拒绝预约'
            : actionModal.action === 'cancel'
            ? '取消预约'
            : '完成课程'
        }
        open={actionModal.visible}
        onCancel={() => {
          setActionModal({ visible: false, booking: null, action: '' });
          actionForm.resetFields();
        }}
        onOk={handleAction}
        okText="确认"
        cancelText="取消"
      >
        {actionModal.action === 'confirm' || actionModal.action === 'complete' ? (
          <p>确认要{actionModal.action === 'confirm' ? '确认' : '完成'}此预约吗？</p>
        ) : (
          <Form form={actionForm} layout="vertical">
            <Form.Item
              label="原因"
              name="reason"
              rules={[{ required: true, message: '请输入原因' }]}
            >
              <Input.TextArea rows={4} placeholder="请输入原因" />
            </Form.Item>
          </Form>
        )}
      </Modal>
    </div>
  );
};

export default BookingList;
