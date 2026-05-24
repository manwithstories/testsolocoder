import React, { useEffect, useState } from 'react';
import { Card, Row, Col, Descriptions, Tag, Button, Rate, Avatar, Spin, message, Modal, Form, Input, List } from 'antd';
import { useParams, useNavigate } from 'react-router-dom';
import { ArrowLeftOutlined, CalendarOutlined, StarOutlined } from '@ant-design/icons';
import { bookingApi, reviewApi } from '@/api/booking';
import { Booking } from '@/types';
import { useSelector } from 'react-redux';
import { RootState } from '@/store';

const { TextArea } = Input;

const BookingDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const currentUser = useSelector((state: RootState) => state.auth.user);
  const [loading, setLoading] = useState(true);
  const [booking, setBooking] = useState<Booking | null>(null);
  const [reviewModalVisible, setReviewModalVisible] = useState(false);
  const [reviewLoading, setReviewLoading] = useState(false);
  const [reviewForm] = Form.useForm();

  useEffect(() => {
    const fetchData = async () => {
      if (!id) return;

      setLoading(true);
      try {
        const data = await bookingApi.get(id);
        setBooking(data);
      } catch (error: any) {
        message.error(error.message || '获取预约详情失败');
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [id]);

  const handleSubmitReview = async () => {
    try {
      const values = await reviewForm.validateFields();
      setReviewLoading(true);

      await reviewApi.create({
        booking_id: id!,
        rating: values.rating,
        content: values.content,
        is_public: true,
      });

      message.success('评价提交成功');
      setReviewModalVisible(false);
      reviewForm.resetFields();
    } catch (error: any) {
      if (!error.errorFields) {
        message.error(error.message || '评价失败');
      }
    } finally {
      setReviewLoading(false);
    }
  };

  const getStatusColor = (status: string) => {
    const colors: Record<string, string> = {
      pending: 'orange',
      confirmed: 'blue',
      rejected: 'red',
      cancelled: 'default',
      completed: 'green',
      no_show: 'red',
    };
    return colors[status] || 'default';
  };

  const getStatusText = (status: string) => {
    const texts: Record<string, string> = {
      pending: '待确认',
      confirmed: '已确认',
      rejected: '已拒绝',
      cancelled: '已取消',
      completed: '已完成',
      no_show: '未出席',
    };
    return texts[status] || status;
  };

  const getPaymentStatusText = (status: string) => {
    const texts: Record<string, string> = {
      pending: '待支付',
      paid: '已支付',
      held: '已托管',
      released: '已结算',
      refunded: '已退款',
      failed: '支付失败',
    };
    return texts[status] || status;
  };

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: 100 }}>
        <Spin size="large" />
      </div>
    );
  }

  if (!booking) {
    return <div>预约不存在</div>;
  }

  const isStudent = currentUser?.id === booking.student_id;
  const isTeacher = currentUser?.id === booking.teacher_id;
  const canReview = booking.status === 'completed' && 
    ((isStudent && !booking.reviewed_by_student) || (isTeacher && !booking.reviewed_by_teacher));

  return (
    <div>
      <Button
        icon={<ArrowLeftOutlined />}
        onClick={() => navigate(-1)}
        style={{ marginBottom: 16 }}
      >
        返回
      </Button>

      <Row gutter={[24, 24]}>
        <Col span={16}>
          <Card
            title="预约详情"
            extra={<Tag color={getStatusColor(booking.status)}>{getStatusText(booking.status)}</Tag>}
          >
            <Descriptions column={2} bordered size="small">
              <Descriptions.Item label="课程名称" span={2}>
                {booking.posting?.title}
              </Descriptions.Item>
              <Descriptions.Item label="学员">
                <Avatar src={booking.student?.avatar} style={{ marginRight: 8 }} />
                {booking.student?.nickname}
              </Descriptions.Item>
              <Descriptions.Item label="老师">
                <Avatar src={booking.teacher?.avatar} style={{ marginRight: 8 }} />
                {booking.teacher?.nickname}
              </Descriptions.Item>
              <Descriptions.Item label="预约时间" span={2}>
                <CalendarOutlined /> {new Date(booking.scheduled_start).toLocaleString()} -{' '}
                {new Date(booking.scheduled_end).toLocaleTimeString()}
              </Descriptions.Item>
              {booking.actual_start && (
                <Descriptions.Item label="实际时间" span={2}>
                  {new Date(booking.actual_start).toLocaleString()} -{' '}
                  {new Date(booking.actual_end).toLocaleTimeString()}
                </Descriptions.Item>
              )}
              <Descriptions.Item label="课程费用">
                <span style={{ color: '#1890ff', fontWeight: 'bold', fontSize: 18 }}>
                  ¥{booking.price.toFixed(2)}
                </span>
              </Descriptions.Item>
              <Descriptions.Item label="平台服务费">
                ¥{booking.platform_fee?.toFixed(2)}
              </Descriptions.Item>
              <Descriptions.Item label="老师收入">
                <span style={{ color: '#52c41a' }}>¥{booking.teacher_earnings?.toFixed(2)}</span>
              </Descriptions.Item>
              <Descriptions.Item label="支付状态">
                <Tag color="blue">{getPaymentStatusText(booking.payment_status)}</Tag>
              </Descriptions.Item>
              {booking.note && (
                <Descriptions.Item label="备注" span={2}>
                  {booking.note}
                </Descriptions.Item>
              )}
              {booking.cancel_reason && (
                <Descriptions.Item label="取消原因" span={2}>
                  {booking.cancel_reason}
                </Descriptions.Item>
              )}
              {booking.reject_reason && (
                <Descriptions.Item label="拒绝原因" span={2}>
                  {booking.reject_reason}
                </Descriptions.Item>
              )}
            </Descriptions>

            {canReview && (
              <div style={{ marginTop: 24, textAlign: 'center' }}>
                <Button
                  type="primary"
                  icon={<StarOutlined />}
                  onClick={() => setReviewModalVisible(true)}
                >
                  提交评价
                </Button>
              </div>
            )}
          </Card>
        </Col>

        <Col span={8}>
          <Card title="操作记录">
            <List
              size="small"
              dataSource={[
                { time: booking.created_at, action: '创建预约' },
                ...(booking.status === 'confirmed' ? [{ time: booking.updated_at, action: '确认预约' }] : []),
                ...(booking.status === 'completed' ? [{ time: booking.actual_end || booking.updated_at, action: '完成课程' }] : []),
                ...(booking.status === 'cancelled' ? [{ time: booking.updated_at, action: '取消预约' }] : []),
              ]}
              renderItem={(item) => (
                <List.Item>
                  <List.Item.Meta
                    title={item.action}
                    description={new Date(item.time).toLocaleString()}
                  />
                </List.Item>
              )}
            />
          </Card>
        </Col>
      </Row>

      <Modal
        title="提交评价"
        open={reviewModalVisible}
        onCancel={() => {
          setReviewModalVisible(false);
          reviewForm.resetFields();
        }}
        onOk={handleSubmitReview}
        confirmLoading={reviewLoading}
        okText="提交评价"
        cancelText="取消"
      >
        <Form form={reviewForm} layout="vertical">
          <Form.Item
            label="评分"
            name="rating"
            rules={[{ required: true, message: '请选择评分' }]}
          >
            <Rate />
          </Form.Item>
          <Form.Item
            label="评价内容"
            name="content"
            rules={[{ required: true, message: '请输入评价内容' }]}
          >
            <TextArea rows={4} placeholder="请输入您的评价" maxLength={1000} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default BookingDetail;
