import { useState } from 'react';
import { Form, Input, Select, DatePicker, InputNumber, Button, Card, message, Row, Col } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { jobApi } from '../../services/api';
import { CreateJobRequest } from '../../types';

const { Option } = Select;
const { TextArea } = Input;
const { RangePicker } = DatePicker;

export const CreateJobPage = () => {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [form] = Form.useForm();

  const handleSubmit = async (values: any) => {
    setLoading(true);
    try {
      const data: CreateJobRequest = {
        ...values,
        start_date: values.dateRange[0].format('YYYY-MM-DD'),
        end_date: values.dateRange[1].format('YYYY-MM-DD'),
      };
      delete (data as any).dateRange;

      const res = await jobApi.createJob(data);
      if (res.data.code === 201) {
        message.success('发布成功');
        navigate(`/jobs/${res.data.data.id}`);
      }
    } catch (error: any) {
      message.error(error.response?.data?.message || '发布失败');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <h2 className="text-2xl font-bold mb-6">发布岗位</h2>
      <Card>
        <Form
          form={form}
          onFinish={handleSubmit}
          layout="vertical"
          initialValues={{
            salary_type: 'hourly',
            is_urgent: false,
          }}
        >
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="title"
                label="岗位标题"
                rules={[{ required: true, message: '请输入岗位标题' }]}
              >
                <Input placeholder="例如：展会临时工作人员" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                name="position"
                label="具体职位"
                rules={[{ required: true, message: '请输入具体职位' }]}
              >
                <Input placeholder="例如：安保、礼仪、促销" />
              </Form.Item>
            </Col>
          </Row>

          <Form.Item
            name="description"
            label="工作描述"
            rules={[{ required: true, message: '请输入工作描述' }]}
          >
            <TextArea rows={4} placeholder="请详细描述工作内容..." />
          </Form.Item>

          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="activity_type" label="活动类型">
                <Select placeholder="请选择活动类型">
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
              <Form.Item
                name="location"
                label="工作地点"
                rules={[{ required: true, message: '请输入工作地点' }]}
              >
                <Input placeholder="请输入工作地点" />
              </Form.Item>
            </Col>
          </Row>

          <Form.Item
            name="dateRange"
            label="活动时间"
            rules={[{ required: true, message: '请选择活动时间' }]}
          >
            <RangePicker style={{ width: '100%' }} />
          </Form.Item>

          <Row gutter={16}>
            <Col span={8}>
              <Form.Item
                name="salary_per_hour"
                label="时薪 (元)"
                rules={[{ required: true, message: '请输入时薪' }]}
              >
                <InputNumber min={0} style={{ width: '100%' }} />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item
                name="headcount"
                label="招聘人数"
                rules={[{ required: true, message: '请输入招聘人数' }]}
              >
                <InputNumber min={1} style={{ width: '100%' }} />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item name="work_hours" label="每日工作时间">
                <Input placeholder="例如：09:00-18:00" />
              </Form.Item>
            </Col>
          </Row>

          <Form.Item name="requirements" label="任职要求">
            <TextArea rows={3} placeholder="请输入任职要求..." />
          </Form.Item>

          <Form.Item name="benefits" label="福利待遇">
            <TextArea rows={3} placeholder="请输入福利待遇..." />
          </Form.Item>

          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="contact_person" label="联系人">
                <Input />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="contact_phone" label="联系电话">
                <Input />
              </Form.Item>
            </Col>
          </Row>

          <Form.Item name="is_urgent" label="紧急招聘" valuePropName="checked">
            <Select>
              <Option value={false}>否</Option>
              <Option value={true}>是</Option>
            </Select>
          </Form.Item>

          <Form.Item name="tags" label="标签 (用逗号分隔)">
            <Input placeholder="例如：日结,包吃住,交通补贴" />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading} icon={<PlusOutlined />}>
              发布岗位
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  );
};
