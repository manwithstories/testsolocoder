import React, { useEffect, useState } from 'react';
import { Card, Form, Input, Select, InputNumber, Button, message, Spin, Row, Col, TimePicker } from 'antd';
import { useNavigate } from 'react-router-dom';
import { skillApi, postingApi } from '@/api/skill';
import { Skill } from '@/types';

const { Option } = Select;
const { TextArea } = Input;

const PostingCreate: React.FC = () => {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [skills, setSkills] = useState<Skill[]>([]);
  const [form] = Form.useForm();

  useEffect(() => {
    const fetchData = async () => {
      try {
        const skillsData = await skillApi.list({ page_size: 100 });
        setSkills(skillsData.items || []);
      } catch (error: any) {
        message.error(error.message || '获取数据失败');
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  const handleSubmit = async (values: any) => {
    setSubmitting(true);
    try {
      const data = {
        skill_id: values.skill_id,
        title: values.title,
        description: values.description,
        teaching_method: values.teaching_method,
        teaching_mode: values.teaching_mode,
        max_students: values.max_students,
        price_per_hour: values.price_per_hour,
        session_duration: values.session_duration,
        location: values.location,
        availability: JSON.stringify({
          days: values.available_days || [],
          start_time: values.start_time?.format('HH:mm'),
          end_time: values.end_time?.format('HH:mm'),
        }),
      };

      await postingApi.create(data);
      message.success('课程发布成功');
      navigate('/skills');
    } catch (error: any) {
      message.error(error.message || '发布失败');
    } finally {
      setSubmitting(false);
    }
  };

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: 100 }}>
        <Spin size="large" />
      </div>
    );
  }

  return (
    <div>
      <Card title="发布技能课程">
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
          style={{ maxWidth: 800 }}
        >
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="skill_id"
                label="选择技能"
                rules={[{ required: true, message: '请选择技能' }]}
              >
                <Select placeholder="请选择要教授的技能" showSearch optionFilterProp="children">
                  {skills.map((skill) => (
                    <Option key={skill.id} value={skill.id}>
                      {skill.title}
                    </Option>
                  ))}
                </Select>
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                name="title"
                label="课程标题"
                rules={[{ required: true, message: '请输入课程标题' }]}
              >
                <Input placeholder="请输入课程标题" maxLength={200} />
              </Form.Item>
            </Col>
          </Row>

          <Form.Item
            name="description"
            label="课程介绍"
            rules={[{ required: true, message: '请输入课程介绍' }]}
          >
            <TextArea
              rows={6}
              placeholder="详细介绍你的课程内容、教学方法等"
              maxLength={2000}
            />
          </Form.Item>

          <Row gutter={16}>
            <Col span={8}>
              <Form.Item
                name="teaching_method"
                label="授课方式"
                rules={[{ required: true, message: '请选择授课方式' }]}
              >
                <Select placeholder="请选择授课方式">
                  <Option value="online">线上</Option>
                  <Option value="offline">线下</Option>
                  <Option value="both">线上/线下</Option>
                </Select>
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item
                name="teaching_mode"
                label="教学模式"
                rules={[{ required: true, message: '请选择教学模式' }]}
              >
                <Select placeholder="请选择教学模式">
                  <Option value="one_to_one">一对一</Option>
                  <Option value="small_class">小班教学</Option>
                  <Option value="group">大班教学</Option>
                </Select>
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item
                name="max_students"
                label="最多学员数"
                rules={[{ required: true, message: '请输入最多学员数' }]}
              >
                <InputNumber min={1} max={50} style={{ width: '100%' }} placeholder="请输入" />
              </Form.Item>
            </Col>
          </Row>

          <Row gutter={16}>
            <Col span={8}>
              <Form.Item
                name="price_per_hour"
                label="每小时价格"
                rules={[{ required: true, message: '请输入价格' }]}
              >
                <InputNumber
                  min={0}
                  precision={2}
                  style={{ width: '100%' }}
                  placeholder="请输入价格"
                  addonBefore="¥"
                  addonAfter="/小时"
                />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item
                name="session_duration"
                label="单次时长(分钟)"
                rules={[{ required: true, message: '请输入单次时长' }]}
              >
                <InputNumber min={30} max={480} step={30} style={{ width: '100%' }} placeholder="请输入" />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item name="location" label="地点(可选)">
                <Input placeholder="线下授课地点" />
              </Form.Item>
            </Col>
          </Row>

          <Row gutter={16}>
            <Col span={8}>
              <Form.Item
                name="available_days"
                label="可用日期"
                rules={[{ required: true, message: '请选择可用日期' }]}
              >
                <Select mode="multiple" placeholder="选择可用日期">
                  <Option value="monday">周一</Option>
                  <Option value="tuesday">周二</Option>
                  <Option value="wednesday">周三</Option>
                  <Option value="thursday">周四</Option>
                  <Option value="friday">周五</Option>
                  <Option value="saturday">周六</Option>
                  <Option value="sunday">周日</Option>
                </Select>
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item
                name="start_time"
                label="可用开始时间"
                rules={[{ required: true, message: '请选择开始时间' }]}
              >
                <TimePicker format="HH:mm" style={{ width: '100%' }} />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item
                name="end_time"
                label="可用结束时间"
                rules={[{ required: true, message: '请选择结束时间' }]}
              >
                <TimePicker format="HH:mm" style={{ width: '100%' }} />
              </Form.Item>
            </Col>
          </Row>

          <Form.Item>
            <Button type="primary" htmlType="submit" loading={submitting} size="large">
              发布课程
            </Button>
            <Button style={{ marginLeft: 8 }} size="large" onClick={() => navigate(-1)}>
              取消
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  );
};

export default PostingCreate;
