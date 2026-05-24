import React, { useEffect, useState } from 'react';
import { Card, Row, Col, Tag, Descriptions, Button, Rate, Avatar, Spin, message, Modal, Form, DatePicker, Input } from 'antd';
import { useParams, useNavigate } from 'react-router-dom';
import { ArrowLeftOutlined, StarOutlined, BookOutlined } from '@ant-design/icons';
import dayjs, { Dayjs } from 'dayjs';
import { postingApi } from '@/api/skill';
import { bookingApi } from '@/api/booking';
import { SkillPosting } from '@/types';

const { RangePicker } = DatePicker;
const { TextArea } = Input;

const PostingDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(true);
  const [posting, setPosting] = useState<SkillPosting | null>(null);
  const [bookingModalVisible, setBookingModalVisible] = useState(false);
  const [bookingLoading, setBookingLoading] = useState(false);
  const [bookingForm] = Form.useForm();

  useEffect(() => {
    const fetchData = async () => {
      if (!id) return;

      setLoading(true);
      try {
        const data = await postingApi.get(id);
        setPosting(data);
      } catch (error: any) {
        message.error(error.message || '获取课程详情失败');
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [id]);

  const handleBooking = async () => {
    try {
      const values = await bookingForm.validateFields();
      setBookingLoading(true);

      const [start, end] = values.timeRange as [Dayjs, Dayjs];

      await bookingApi.create({
        posting_id: id!,
        scheduled_start: start.toISOString(),
        scheduled_end: end.toISOString(),
        note: values.note,
      });

      message.success('预约请求已发送');
      setBookingModalVisible(false);
      bookingForm.resetFields();
    } catch (error: any) {
      if (error.errorFields) {
        return;
      }
      message.error(error.message || '预约失败');
    } finally {
      setBookingLoading(false);
    }
  };

  const teachingMethodLabels: Record<string, string> = {
    online: '线上',
    offline: '线下',
    both: '线上/线下',
  };

  const teachingModeLabels: Record<string, string> = {
    one_to_one: '一对一',
    small_class: '小班教学',
    group: '大班教学',
  };

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: 100 }}>
        <Spin size="large" />
      </div>
    );
  }

  if (!posting) {
    return <div>课程不存在</div>;
  }

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
          <Card title={posting.title} extra={<Tag color="blue">{posting.skill?.title}</Tag>}>
            <Descriptions column={2} bordered size="small">
              <Descriptions.Item label="授课老师">
                <Avatar src={posting.teacher?.avatar} style={{ marginRight: 8 }} />
                {posting.teacher?.nickname}
                <Rate
                  disabled
                  defaultValue={posting.teacher?.rating || 0}
                  style={{ fontSize: 12, marginLeft: 8 }}
                />
              </Descriptions.Item>
              <Descriptions.Item label="授课方式">
                {teachingMethodLabels[posting.teaching_method]}
              </Descriptions.Item>
              <Descriptions.Item label="教学模式">
                {teachingModeLabels[posting.teaching_mode]}
              </Descriptions.Item>
              <Descriptions.Item label="最多学员">
                {posting.max_students}人
              </Descriptions.Item>
              <Descriptions.Item label="课程时长">
                {posting.session_duration}分钟
              </Descriptions.Item>
              <Descriptions.Item label="课程费用">
                <span style={{ color: '#1890ff', fontWeight: 'bold', fontSize: 18 }}>
                  ¥{posting.price_per_hour}
                </span>
                /小时
              </Descriptions.Item>
              {posting.location && (
                <Descriptions.Item label="地点" span={2}>
                  {posting.location}
                </Descriptions.Item>
              )}
            </Descriptions>

            {posting.description && (
              <div style={{ marginTop: 24 }}>
                <h3>课程介绍</h3>
                <p style={{ whiteSpace: 'pre-wrap' }}>{posting.description}</p>
              </div>
            )}

            <div style={{ marginTop: 24, textAlign: 'center' }}>
              <Button
                type="primary"
                size="large"
                icon={<BookOutlined />}
                onClick={() => setBookingModalVisible(true)}
              >
                立即预约
              </Button>
            </div>
          </Card>
        </Col>

        <Col span={8}>
          <Card title="老师信息">
            <div style={{ textAlign: 'center', marginBottom: 16 }}>
              <Avatar
                size={80}
                src={posting.teacher?.avatar}
                icon={<StarOutlined />}
              />
              <h3 style={{ marginTop: 12 }}>{posting.teacher?.nickname}</h3>
              <div>
                <Rate disabled defaultValue={posting.teacher?.rating || 0} />
                <span style={{ marginLeft: 8 }}>
                  ({posting.teacher?.review_count}条评价)
                </span>
              </div>
            </div>
            <Descriptions column={1} size="small">
              <Descriptions.Item label="授课时长">
                {posting.teacher?.teaching_hours?.toFixed(1) || 0}小时
              </Descriptions.Item>
              <Descriptions.Item label="学员数量">
                {posting.teacher?.student_count || 0}人
              </Descriptions.Item>
              {posting.teacher?.location && (
                <Descriptions.Item label="地点">
                  {posting.teacher.location}
                </Descriptions.Item>
              )}
            </Descriptions>
          </Card>

          <Card title="课程数据" style={{ marginTop: 16 }}>
            <Descriptions column={1} size="small">
              <Descriptions.Item label="预约次数">
                {posting.booking_count}次
              </Descriptions.Item>
              <Descriptions.Item label="总授课时长">
                {posting.total_hours?.toFixed(1) || 0}小时
              </Descriptions.Item>
              <Descriptions.Item label="课程评分">
                <span style={{ color: '#faad14' }}>★</span> {posting.rating.toFixed(1)}
              </Descriptions.Item>
              <Descriptions.Item label="评价数量">
                {posting.review_count}条
              </Descriptions.Item>
            </Descriptions>
          </Card>
        </Col>
      </Row>

      <Modal
        title="预约课程"
        open={bookingModalVisible}
        onCancel={() => setBookingModalVisible(false)}
        onOk={handleBooking}
        confirmLoading={bookingLoading}
        okText="确认预约"
        cancelText="取消"
      >
        <Form form={bookingForm} layout="vertical">
          <Form.Item
            label="预约时间"
            name="timeRange"
            rules={[{ required: true, message: '请选择预约时间' }]}
          >
            <RangePicker
              showTime={{ format: 'HH:mm' }}
              format="YYYY-MM-DD HH:mm"
              style={{ width: '100%' }}
              disabledDate={(current) =>
                current && current < dayjs().startOf('day')
              }
            />
          </Form.Item>
          <Form.Item label="备注" name="note">
            <TextArea rows={4} placeholder="请输入备注信息（可选）" maxLength={500} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default PostingDetail;
