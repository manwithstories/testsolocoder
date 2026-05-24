import { useState, useEffect } from 'react';
import { Card, List, Avatar, Rate, Tag, Statistic, Row, Col, Button, Modal, Form, Input, Select, message } from 'antd';
import { StarOutlined } from '@ant-design/icons';
import { evaluationApi } from '../../services/api';
import { useAuthStore } from '../../context/AuthContext';
import { Evaluation } from '../../types';
import dayjs from 'dayjs';

const { TextArea } = Input;

export const MyEvaluationsPage = () => {
  const { user } = useAuthStore();
  const [evaluations, setEvaluations] = useState<Evaluation[]>([]);
  const [stats, setStats] = useState<any>(null);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [form] = Form.useForm();

  const fetchData = async () => {
    setLoading(true);
    try {
      const [evalRes, statsRes] = await Promise.all([
        evaluationApi.getMyEvaluations(),
        user ? evaluationApi.getEvaluationStats(user.id) : Promise.resolve({ data: { data: null } }),
      ]);

      if (evalRes.data.code === 200) {
        setEvaluations(evalRes.data.data.data || []);
      }
      if (statsRes.data.code === 200) {
        setStats(statsRes.data.data);
      }
    } catch (error) {
      message.error('获取数据失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  const handleSubmit = async (values: any) => {
    try {
      const res = await evaluationApi.createEvaluation(values);
      if (res.data.code === 201) {
        message.success('评价提交成功');
        setModalVisible(false);
        form.resetFields();
        fetchData();
      }
    } catch (error: any) {
      message.error(error.response?.data?.message || '提交失败');
    }
  };

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-2xl font-bold">我的评价</h2>
        <Button type="primary" onClick={() => setModalVisible(true)}>
          提交评价
        </Button>
      </div>

      {stats && (
        <Row gutter={[16, 16]} className="mb-6">
          <Col xs={12} sm={8}>
            <Card>
              <Statistic
                title="收到评价数"
                value={stats.total_count}
                prefix={<StarOutlined />}
              />
            </Card>
          </Col>
          <Col xs={12} sm={8}>
            <Card>
              <Statistic
                title="平均评分"
                value={stats.average_rating}
                precision={1}
                suffix="/ 5"
                valueStyle={{ color: '#faad14' }}
              />
            </Card>
          </Col>
          <Col xs={12} sm={8}>
            <Card>
              <Statistic
                title="信用分"
                value={stats.credit_score}
                suffix="/ 100"
                valueStyle={{ color: '#52c41a' }}
              />
            </Card>
          </Col>
        </Row>
      )}

      <Card>
        {evaluations.length === 0 ? (
          <div className="text-center py-12 text-gray-500">暂无评价</div>
        ) : (
          <List
            itemLayout="horizontal"
            dataSource={evaluations}
            renderItem={(item) => (
              <List.Item>
                <List.Item.Meta
                avatar={<Avatar src={item.from_user?.avatar}>{item.is_anonymous ? '?' : item.from_user?.real_name?.[0]}</Avatar>}
                title={
                  <div className="flex items-center gap-2">
                  <span>{item.is_anonymous ? '匿名用户' : item.from_user?.real_name}</span>
                  <Rate disabled value={item.rating} />
                  <Tag>{item.type === 'employer_to_temp' ? '雇主评价' : '临时工评价'}</Tag>
                </div>
              }
                description={
                  <div>
                    <p className="mb-1">{item.content || '暂无评价内容'}</p>
                    <p className="text-xs text-gray-400">
                      {item.job_posting?.position} | {dayjs(item.created_at).format('YYYY-MM-DD HH:mm')}
                    </p>
                  </div>
                }
              />
            </List.Item>
          )}
        />
        )}
      </Card>

      <Modal
        title="提交评价"
        open={modalVisible}
        onCancel={() => setModalVisible(false)}
        footer={null}
        width={500}
      >
        <Form form={form} onFinish={handleSubmit} layout="vertical">
          <Form.Item name="job_id" label="岗位ID" rules={[{ required: true, message: '请输入岗位ID' }]}>
            <Input placeholder="请输入岗位ID" />
          </Form.Item>
          <Form.Item name="to_user_id" label="被评价人ID" rules={[{ required: true, message: '请输入被评价人ID' }]}>
            <Input placeholder="请输入被评价人ID" />
          </Form.Item>
          <Form.Item name="type" label="评价类型" rules={[{ required: true, message: '请选择评价类型' }]}>
            <Select>
              <Select.Option value="employer_to_temp">雇主评临时工</Select.Option>
              <Select.Option value="temp_to_employer">临时工评雇主</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item name="rating" label="评分" rules={[{ required: true, message: '请选择评分' }]}>
            <Rate />
          </Form.Item>
          <Form.Item name="content" label="评价内容">
            <TextArea rows={4} placeholder="请输入评价内容..." />
          </Form.Item>
          <Form.Item name="is_anonymous" label="匿名评价" valuePropName="checked">
            <Select>
              <Select.Option value={false}>否</Select.Option>
              <Select.Option value={true}>是</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block>提交评价</Button>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};
