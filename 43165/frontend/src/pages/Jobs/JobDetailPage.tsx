import { useState, useEffect } from 'react';
import { Card, Descriptions, Tag, Button, Space, Modal, message, List, Avatar, Form, Input, Select, DatePicker, InputNumber, Divider, Row, Col } from 'antd';
import { useParams, useNavigate } from 'react-router-dom';
import {
  EnvironmentOutlined,
  DollarOutlined,
  CalendarOutlined,
  TeamOutlined,
  ClockCircleOutlined,
  ArrowLeftOutlined,
  EditOutlined,
} from '@ant-design/icons';
import { jobApi } from '../../services/api';
import { useAuthStore } from '../../context/AuthContext';
import { JobPosting, JobApplication } from '../../types';
import dayjs from 'dayjs';

const { Option } = Select;
const { TextArea } = Input;
const { RangePicker } = DatePicker;

export const JobDetailPage = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { user } = useAuthStore();
  const [job, setJob] = useState<JobPosting | null>(null);
  const [loading, setLoading] = useState(true);
  const [applyModalVisible, setApplyModalVisible] = useState(false);
  const [editModalVisible, setEditModalVisible] = useState(false);
  const [applications, setApplications] = useState<JobApplication[]>([]);
  const [form] = Form.useForm();
  const [editForm] = Form.useForm();

  const fetchJob = async () => {
    try {
      const res = await jobApi.getJob(id!);
      if (res.data.code === 200) {
        setJob(res.data.data);
        editForm.setFieldsValue({
          ...res.data.data,
          dateRange: [dayjs(res.data.data.start_date), dayjs(res.data.data.end_date)],
        });
      }
    } catch (error) {
      message.error('获取岗位详情失败');
    } finally {
      setLoading(false);
    }
  };

  const fetchApplications = async () => {
    if (user?.role === 'employer' || user?.role === 'agent') {
      try {
        const res = await jobApi.getApplications(id!);
        if (res.data.code === 200) {
          setApplications(res.data.data.data);
        }
      } catch (error) {
        console.error('获取申请列表失败');
      }
    }
  };

  useEffect(() => {
    fetchJob();
    fetchApplications();
  }, [id]);

  const handleApply = async (values: any) => {
    try {
      const res = await jobApi.applyJob(id!, values);
      if (res.data.code === 201) {
        message.success('申请成功');
        setApplyModalVisible(false);
        form.resetFields();
        fetchJob();
      }
    } catch (error: any) {
      message.error(error.response?.data?.message || '申请失败');
    }
  };

  const handleEdit = async (values: any) => {
    try {
      const data = {
        ...values,
        start_date: values.dateRange[0].format('YYYY-MM-DD'),
        end_date: values.dateRange[1].format('YYYY-MM-DD'),
      };
      delete data.dateRange;

      const res = await jobApi.updateJob(id!, data);
      if (res.data.code === 200) {
        message.success('更新成功');
        setEditModalVisible(false);
        fetchJob();
      }
    } catch (error: any) {
      message.error(error.response?.data?.message || '更新失败');
    }
  };

  const handleReviewApplication = async (appId: string, approved: boolean) => {
    try {
      const res = await jobApi.reviewApplication(appId, {
        status: approved ? 'approved' : 'rejected',
        review_note: approved ? '已通过审核' : '未通过审核',
      });
      if (res.data.code === 200) {
        message.success('审核成功');
        fetchApplications();
      }
    } catch (error: any) {
      message.error(error.response?.data?.message || '审核失败');
    }
  };

  if (loading) {
    return <div className="flex justify-center items-center h-96"><div className="animate-spin text-4xl text-blue-500">⚙️</div></div>;
  }

  if (!job) {
    return <div>岗位不存在</div>;
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <Button icon={<ArrowLeftOutlined />} onClick={() => navigate(-1)}>
          返回
        </Button>
        <Space>
          {job.employer_id === user?.id && (
            <Button icon={<EditOutlined />} onClick={() => setEditModalVisible(true)}>
              编辑
            </Button>
          )}
          {user?.role === 'temporary' && job.status === 'recruiting' && (
            <Button type="primary" onClick={() => setApplyModalVisible(true)}>
              申请岗位
            </Button>
          )}
        </Space>
      </div>

      <Card title={job.title} extra={job.is_urgent && <Tag color="red">急招</Tag>}>
        <Descriptions column={2} bordered>
          <Descriptions.Item label="职位">{job.position}</Descriptions.Item>
          <Descriptions.Item label="活动类型">
            <Tag>{job.activity_type || '未分类'}</Tag>
          </Descriptions.Item>
          <Descriptions.Item label="地点">
            <EnvironmentOutlined className="mr-1" />{job.location}
          </Descriptions.Item>
          <Descriptions.Item label="薪资">
            <span className="text-green-600 font-semibold">
              <DollarOutlined className="mr-1" />¥{job.salary_per_hour}/小时
            </span>
          </Descriptions.Item>
          <Descriptions.Item label="时间">
            <CalendarOutlined className="mr-1" />
            {dayjs(job.start_date).format('YYYY-MM-DD')} ~ {dayjs(job.end_date).format('YYYY-MM-DD')}
          </Descriptions.Item>
          <Descriptions.Item label="招聘人数">
            <TeamOutlined className="mr-1" />
            {job.hired_count}/{job.headcount} 人
          </Descriptions.Item>
          <Descriptions.Item label="工作时间">{job.work_hours}</Descriptions.Item>
          <Descriptions.Item label="薪资类型">{job.salary_type}</Descriptions.Item>
          <Descriptions.Item label="状态" span={2}>
            <Tag color={job.status === 'recruiting' ? 'green' : 'default'}>{job.status}</Tag>
          </Descriptions.Item>
        </Descriptions>

        <Divider />

        <div className="mb-4">
          <h4>工作描述</h4>
          <p className="text-gray-600 whitespace-pre-wrap">{job.description}</p>
        </div>

        {job.requirements && (
          <div className="mb-4">
            <h4>任职要求</h4>
            <p className="text-gray-600 whitespace-pre-wrap">{job.requirements}</p>
          </div>
        )}

        {job.benefits && (
          <div className="mb-4">
            <h4>福利待遇</h4>
            <p className="text-gray-600 whitespace-pre-wrap">{job.benefits}</p>
          </div>
        )}

        {job.contact_person && (
          <div>
            <h4>联系方式</h4>
            <p>联系人：{job.contact_person}</p>
            <p>电话：{job.contact_phone}</p>
          </div>
        )}
      </Card>

      {(user?.role === 'employer' || user?.role === 'agent') && (
        <Card title="申请列表" className="mt-4">
          {applications.length === 0 ? (
            <div className="text-center py-8 text-gray-500">暂无申请</div>
          ) : (
            <List
              itemLayout="horizontal"
              dataSource={applications}
              renderItem={(item) => (
                <List.Item
                  actions={
                item.status === 'pending' ? [
                  <Button key="approve" type="primary" size="small" onClick={() => handleReviewApplication(item.id, true)}>
                    通过
                  </Button>,
                  <Button key="reject" size="small" danger onClick={() => handleReviewApplication(item.id, false)}>
                    拒绝
                  </Button>,
                ] : undefined
              }
                >
                  <List.Item.Meta
                    avatar={<Avatar src={item.temporary?.avatar}>{item.temporary?.real_name?.[0]}</Avatar>}
                    title={item.temporary?.real_name}
                    description={
                      <Space>
                        <Tag color={item.status === 'approved' ? 'green' : item.status === 'rejected' ? 'red' : 'orange'}>
                          {item.status}
                        </Tag>
                        <span>{item.message}</span>
                      </Space>
                    }
                  />
                </List.Item>
              )}
            />
          )}
        </Card>
      )}

      <Modal
        title="申请岗位"
        open={applyModalVisible}
        onCancel={() => setApplyModalVisible(false)}
        footer={null}
      >
        <Form form={form} onFinish={handleApply} layout="vertical">
          <Form.Item name="message" label="申请留言">
            <TextArea rows={4} placeholder="请输入您的申请留言..." />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block>
              提交申请
            </Button>
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title="编辑岗位"
        open={editModalVisible}
        onCancel={() => setEditModalVisible(false)}
        footer={null}
        width={700}
      >
        <Form form={editForm} onFinish={handleEdit} layout="vertical">
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="title" label="岗位标题" rules={[{ required: true }]}>
                <Input />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="position" label="职位" rules={[{ required: true }]}>
                <Input />
              </Form.Item>
            </Col>
          </Row>
          <Form.Item name="description" label="工作描述" rules={[{ required: true }]}>
            <TextArea rows={4} />
          </Form.Item>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="activity_type" label="活动类型">
                <Select>
                  <Option value="exhibition">展会</Option>
                  <Option value="conference">会议</Option>
                  <Option value="performance">演出</Option>
                  <Option value="promotion">促销</Option>
                  <Option value="wedding">婚礼</Option>
                  <Option value="other">其他</Option>
                </Select>
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="location" label="地点" rules={[{ required: true }]}>
                <Input />
              </Form.Item>
            </Col>
          </Row>
          <Form.Item name="dateRange" label="活动时间" rules={[{ required: true }]}>
            <RangePicker style={{ width: '100%' }} />
          </Form.Item>
          <Row gutter={16}>
            <Col span={8}>
              <Form.Item name="salary_per_hour" label="时薪(元)" rules={[{ required: true }]}>
                <InputNumber min={0} style={{ width: '100%' }} />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item name="headcount" label="招聘人数" rules={[{ required: true }]}>
                <InputNumber min={1} style={{ width: '100%' }} />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item name="is_urgent" label="紧急招聘" valuePropName="checked">
                <Select>
                  <Option value={true}>是</Option>
                  <Option value={false}>否</Option>
                </Select>
              </Form.Item>
            </Col>
          </Row>
          <Form.Item>
            <Button type="primary" htmlType="submit" block>
              保存修改
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};
