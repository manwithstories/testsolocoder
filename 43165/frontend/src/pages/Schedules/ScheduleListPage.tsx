import { useState, useEffect } from 'react';
import { Table, Tag, Button, Modal, Form, Input, Select, DatePicker, TimePicker, Space, message, Card, Popconfirm } from 'antd';
import { PlusOutlined, DeleteOutlined, EditOutlined, ExportOutlined } from '@ant-design/icons';
import { scheduleApi } from '../../services/api';
import { useAuthStore } from '../../context/AuthContext';
import { Schedule } from '../../types';
import dayjs from 'dayjs';

export const ScheduleListPage = () => {
  const { user } = useAuthStore();
  const [schedules, setSchedules] = useState<Schedule[]>([]);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [createModalVisible, setCreateModalVisible] = useState(false);
  const [editModalVisible, setEditModalVisible] = useState(false);
  const [currentSchedule, setCurrentSchedule] = useState<Schedule | null>(null);
  const [form] = Form.useForm();
  const [editForm] = Form.useForm();

  const fetchSchedules = async () => {
    setLoading(true);
    try {
      const params: any = { page, page_size: pageSize };
      const res = user?.role === 'temporary'
        ? await scheduleApi.getMySchedules(params)
        : await scheduleApi.getSchedules(params);
      if (res.data.code === 200) {
        setSchedules(res.data.data.data || []);
        setTotal(res.data.data.total || 0);
      }
    } catch (error) {
      message.error('获取排班列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchSchedules();
  }, [page, pageSize]);

  const handleCreate = async (values: any) => {
    try {
      const data = {
        ...values,
        shift_date: values.shift_date.format('YYYY-MM-DD'),
        start_time: values.start_time.format('HH:mm'),
        end_time: values.end_time.format('HH:mm'),
      };
      const res = await scheduleApi.createSchedule(data);
      if (res.data.code === 201) {
        message.success('创建成功');
        setCreateModalVisible(false);
        form.resetFields();
        fetchSchedules();
      }
    } catch (error: any) {
      message.error(error.response?.data?.message || '创建失败');
    }
  };

  const handleEdit = async (values: any) => {
    try {
      const data = { ...values };
      if (data.start_time) data.start_time = data.start_time.format('HH:mm');
      if (data.end_time) data.end_time = data.end_time.format('HH:mm');
      
      const res = await scheduleApi.updateSchedule(currentSchedule!.id, data);
      if (res.data.code === 200) {
        message.success('更新成功');
        setEditModalVisible(false);
        editForm.resetFields();
        fetchSchedules();
      }
    } catch (error: any) {
      message.error(error.response?.data?.message || '更新失败');
    }
  };

  const handleDelete = async (id: string) => {
    try {
      const res = await scheduleApi.deleteSchedule(id);
      if (res.data.code === 200) {
        message.success('删除成功');
        fetchSchedules();
      }
    } catch (error: any) {
      message.error(error.response?.data?.message || '删除失败');
    }
  };

  const handleExport = async () => {
    try {
      const res = await scheduleApi.exportSchedules();
      const blob = new Blob([res.data]);
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `schedules_${dayjs().format('YYYYMMDD')}.xlsx`;
      a.click();
      window.URL.revokeObjectURL(url);
      message.success('导出成功');
    } catch (error) {
      message.error('导出失败');
    }
  };

  const columns = [
    {
      title: '日期',
      dataIndex: 'shift_date',
      key: 'shift_date',
      render: (text: string) => dayjs(text).format('YYYY-MM-DD'),
    },
    {
      title: '时间',
      key: 'time',
      render: (_: any, record: Schedule) => `${record.start_time} - ${record.end_time}`,
    },
    {
      title: '岗位',
      dataIndex: ['job_posting', 'position'],
      key: 'position',
    },
    {
      title: '临时工',
      dataIndex: ['temporary', 'real_name'],
      key: 'temporary',
    },
    {
      title: '地点',
      dataIndex: 'location',
      key: 'location',
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => {
        const colorMap: Record<string, string> = {
          scheduled: 'blue',
          in_progress: 'orange',
          completed: 'green',
        };
        const labelMap: Record<string, string> = {
          scheduled: '已排班',
          in_progress: '进行中',
          completed: '已完成',
        };
        return <Tag color={colorMap[status]}>{labelMap[status] || status}</Tag>;
      },
    },
    {
      title: '操作',
      key: 'actions',
      render: (_: any, record: Schedule) => (
        <Space>
          {user?.role !== 'temporary' && (
            <>
              <Button
                size="small"
                icon={<EditOutlined />}
                onClick={() => {
                  setCurrentSchedule(record);
                  editForm.setFieldsValue({
                    start_time: dayjs(record.start_time, 'HH:mm'),
                    end_time: dayjs(record.end_time, 'HH:mm'),
                    location: record.location,
                    notes: record.notes,
                    status: record.status,
                  });
                  setEditModalVisible(true);
                }}
              >
                编辑
              </Button>
              <Popconfirm title="确定删除?" onConfirm={() => handleDelete(record.id)}>
                <Button size="small" danger icon={<DeleteOutlined />}>删除</Button>
              </Popconfirm>
            </>
          )}
        </Space>
      ),
    },
  ];

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-2xl font-bold">排班管理</h2>
        <Space>
          {user?.role !== 'temporary' && (
            <Button icon={<ExportOutlined />} onClick={handleExport}>
              导出
            </Button>
          )}
          {user?.role !== 'temporary' && (
            <Button type="primary" icon={<PlusOutlined />} onClick={() => setCreateModalVisible(true)}>
              新建排班
            </Button>
          )}
        </Space>
      </div>

      <Table
        columns={columns}
        dataSource={schedules}
        rowKey="id"
        loading={loading}
        pagination={{
          current: page,
          pageSize,
          total,
          showSizeChanger: true,
          onChange: (p, ps) => {
            setPage(p);
            setPageSize(ps);
          },
        }}
      />

      <Modal
        title="新建排班"
        open={createModalVisible}
        onCancel={() => setCreateModalVisible(false)}
        footer={null}
        width={600}
      >
        <Form form={form} onFinish={handleCreate} layout="vertical">
          <Form.Item name="job_id" label="岗位" rules={[{ required: true, message: '请选择岗位' }]}>
            <Input placeholder="请输入岗位ID" />
          </Form.Item>
          <Form.Item name="temporary_id" label="临时工" rules={[{ required: true, message: '请选择临时工' }]}>
            <Input placeholder="请输入临时工ID" />
          </Form.Item>
          <Form.Item name="shift_date" label="日期" rules={[{ required: true, message: '请选择日期' }]}>
            <DatePicker style={{ width: '100%' }} />
          </Form.Item>
          <div className="flex gap-4">
            <Form.Item name="start_time" label="开始时间" rules={[{ required: true, message: '请选择开始时间' }]} style={{ flex: 1 }}>
              <TimePicker format="HH:mm" style={{ width: '100%' }} />
            </Form.Item>
            <Form.Item name="end_time" label="结束时间" rules={[{ required: true, message: '请选择结束时间' }]} style={{ flex: 1 }}>
              <TimePicker format="HH:mm" style={{ width: '100%' }} />
            </Form.Item>
          </div>
          <Form.Item name="location" label="地点">
            <Input />
          </Form.Item>
          <Form.Item name="notes" label="备注">
            <Input.TextArea rows={3} />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block>创建</Button>
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title="编辑排班"
        open={editModalVisible}
        onCancel={() => setEditModalVisible(false)}
        footer={null}
        width={600}
      >
        <Form form={editForm} onFinish={handleEdit} layout="vertical">
          <div className="flex gap-4">
            <Form.Item name="start_time" label="开始时间" style={{ flex: 1 }}>
              <TimePicker format="HH:mm" style={{ width: '100%' }} />
            </Form.Item>
            <Form.Item name="end_time" label="结束时间" style={{ flex: 1 }}>
              <TimePicker format="HH:mm" style={{ width: '100%' }} />
            </Form.Item>
          </div>
          <Form.Item name="location" label="地点">
            <Input />
          </Form.Item>
          <Form.Item name="notes" label="备注">
            <Input.TextArea rows={3} />
          </Form.Item>
          <Form.Item name="status" label="状态">
            <Select>
              <Select.Option value="scheduled">已排班</Select.Option>
              <Select.Option value="in_progress">进行中</Select.Option>
              <Select.Option value="completed">已完成</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block>保存</Button>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};
