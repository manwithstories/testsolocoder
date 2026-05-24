import React, { useEffect, useState } from 'react';
import { Card, Row, Col, Button, List, Tag, Modal, Form, Input, Select, TimePicker, message, Spin, Empty } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { scheduleApi } from '@/api/message';
import { Schedule as ScheduleType, ScheduleType as ScheduleTypeEnum } from '@/types';
import dayjs from 'dayjs';

const { Option } = Select;

const SchedulePage: React.FC = () => {
  const [loading, setLoading] = useState(true);
  const [schedules, setSchedules] = useState<ScheduleType[]>([]);
  const [modalVisible, setModalVisible] = useState(false);
  const [editingSchedule, setEditingSchedule] = useState<ScheduleType | null>(null);
  const [modalLoading, setModalLoading] = useState(false);
  const [form] = Form.useForm();

  const dayOfWeekOptions = [
    { value: 'monday', label: '周一' },
    { value: 'tuesday', label: '周二' },
    { value: 'wednesday', label: '周三' },
    { value: 'thursday', label: '周四' },
    { value: 'friday', label: '周五' },
    { value: 'saturday', label: '周六' },
    { value: 'sunday', label: '周日' },
  ];

  useEffect(() => {
    fetchSchedules();
  }, []);

  const fetchSchedules = async () => {
    setLoading(true);
    try {
      const data = await scheduleApi.list();
      setSchedules(data);
    } catch (error: any) {
      message.error(error.message || '获取日程失败');
    } finally {
      setLoading(false);
    }
  };

  const handleAdd = () => {
    setEditingSchedule(null);
    form.resetFields();
    setModalVisible(true);
  };

  const handleEdit = (schedule: ScheduleType) => {
    setEditingSchedule(schedule);
    form.setFieldsValue({
      type: schedule.type,
      day_of_week: schedule.day_of_week,
      start_time: schedule.start_time ? dayjs(schedule.start_time, 'HH:mm') : null,
      end_time: schedule.end_time ? dayjs(schedule.end_time, 'HH:mm') : null,
      is_recurring: schedule.is_recurring,
      title: schedule.title,
      description: schedule.description,
    });
    setModalVisible(true);
  };

  const handleDelete = async (id: string) => {
    try {
      await scheduleApi.delete(id);
      message.success('删除成功');
      fetchSchedules();
    } catch (error: any) {
      message.error(error.message || '删除失败');
    }
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      setModalLoading(true);

      const data = {
        type: values.type,
        day_of_week: values.day_of_week,
        start_time: values.start_time.format('HH:mm'),
        end_time: values.end_time.format('HH:mm'),
        is_recurring: values.is_recurring,
        title: values.title,
        description: values.description,
      };

      if (editingSchedule) {
        await scheduleApi.update(editingSchedule.id, data);
        message.success('更新成功');
      } else {
        await scheduleApi.create(data);
        message.success('创建成功');
      }

      setModalVisible(false);
      form.resetFields();
      fetchSchedules();
    } catch (error: any) {
      if (!error.errorFields) {
        message.error(error.message || '保存失败');
      }
    } finally {
      setModalLoading(false);
    }
  };

  const getTypeColor = (type: ScheduleTypeEnum) => {
    const colors: Record<ScheduleTypeEnum, string> = {
      availability: 'green',
      busy: 'red',
      booking: 'blue',
    };
    return colors[type];
  };

  const getTypeText = (type: ScheduleTypeEnum) => {
    const texts: Record<ScheduleTypeEnum, string> = {
      availability: '可用',
      busy: '忙碌',
      booking: '预约',
    };
    return texts[type];
  };

  const groupedSchedules = dayOfWeekOptions.reduce((acc, day) => {
    acc[day.value] = schedules.filter((s) => s.day_of_week === day.value);
    return acc;
  }, {} as Record<string, ScheduleType[]>);

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: 100 }}>
        <Spin size="large" />
      </div>
    );
  }

  return (
    <div>
      <Card
        title="日程管理"
        extra={
          <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
            添加日程
          </Button>
        }
      >
        {schedules.length === 0 ? (
          <Empty description="暂无日程安排" style={{ padding: 40 }} />
        ) : (
          <Row gutter={[16, 16]}>
            {dayOfWeekOptions.map((day) => (
              <Col span={8} key={day.value}>
                <Card
                  size="small"
                  title={day.label}
                  style={{ marginBottom: 16 }}
                >
                  {groupedSchedules[day.value]?.length === 0 ? (
                    <div style={{ color: '#999', textAlign: 'center', padding: 8 }}>
                      无安排
                    </div>
                  ) : (
                    <List
                      size="small"
                      dataSource={groupedSchedules[day.value]}
                      renderItem={(item) => (
                        <List.Item
                          actions={[
                            <Button
                              key="edit"
                              type="link"
                              size="small"
                              icon={<EditOutlined />}
                              onClick={() => handleEdit(item)}
                            />,
                            item.type !== 'booking' && (
                              <Button
                                key="delete"
                                type="link"
                                size="small"
                                danger
                                icon={<DeleteOutlined />}
                                onClick={() => handleDelete(item.id)}
                              />
                            ),
                          ]}
                        >
                          <List.Item.Meta
                            title={
                              <div>
                                <Tag color={getTypeColor(item.type)}>
                                  {getTypeText(item.type)}
                                </Tag>
                                {item.title}
                              </div>
                            }
                            description={`${item.start_time} - ${item.end_time}`}
                          />
                        </List.Item>
                      )}
                    />
                  )}
                </Card>
              </Col>
            ))}
          </Row>
        )}
      </Card>

      <Modal
        title={editingSchedule ? '编辑日程' : '添加日程'}
        open={modalVisible}
        onCancel={() => {
          setModalVisible(false);
          form.resetFields();
        }}
        onOk={handleSubmit}
        confirmLoading={modalLoading}
        okText="保存"
        cancelText="取消"
      >
        <Form form={form} layout="vertical">
          <Form.Item name="type" label="类型" rules={[{ required: true, message: '请选择类型' }]}>
            <Select>
              <Option value="availability">可用时间</Option>
              <Option value="busy">忙碌时间</Option>
            </Select>
          </Form.Item>
          <Form.Item name="day_of_week" label="星期" rules={[{ required: true, message: '请选择星期' }]}>
            <Select>
              {dayOfWeekOptions.map((day) => (
                <Option key={day.value} value={day.value}>
                  {day.label}
                </Option>
              ))}
            </Select>
          </Form.Item>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="start_time" label="开始时间" rules={[{ required: true, message: '请选择开始时间' }]}>
                <TimePicker format="HH:mm" style={{ width: '100%' }} />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="end_time" label="结束时间" rules={[{ required: true, message: '请选择结束时间' }]}>
                <TimePicker format="HH:mm" style={{ width: '100%' }} />
              </Form.Item>
            </Col>
          </Row>
          <Form.Item name="is_recurring" label="是否重复" valuePropName="checked">
            <Select>
              <Option value={true}>每周重复</Option>
              <Option value={false}>仅一次</Option>
            </Select>
          </Form.Item>
          <Form.Item name="title" label="标题">
            <Input placeholder="请输入标题" />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input.TextArea rows={3} placeholder="请输入描述" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default SchedulePage;
