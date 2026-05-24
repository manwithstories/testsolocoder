import { useState, useEffect } from 'react';
import { Table, Tag, Card, Button, Space, Modal, Form, Select, DatePicker, InputNumber, message, Descriptions, List, Divider, Input } from 'antd';
import { DollarOutlined, ExportOutlined, PlusOutlined, SendOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { salaryApi } from '../../services/api';
import { useAuthStore } from '../../context/AuthContext';
import { SalaryRecord, SalaryDetail } from '../../types';
import dayjs from 'dayjs';

const { RangePicker } = DatePicker;

export const SalaryListPage = () => {
  const navigate = useNavigate();
  const { user } = useAuthStore();
  const [salaries, setSalaries] = useState<SalaryRecord[]>([]);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [dateRange, setDateRange] = useState<any>(null);
  const [status, setStatus] = useState<string | undefined>();
  const [createModalVisible, setCreateModalVisible] = useState(false);
  const [form] = Form.useForm();

  const fetchSalaries = async () => {
    setLoading(true);
    try {
      const params: any = { page, page_size: pageSize };
      if (dateRange && dateRange.length === 2) {
        params.start_date = dateRange[0].format('YYYY-MM-DD');
        params.end_date = dateRange[1].format('YYYY-MM-DD');
      }
      if (status) params.status = status;

      const res = await salaryApi.getSalaries(params);
      if (res.data.code === 200) {
        setSalaries(res.data.data.data || []);
        setTotal(res.data.data.total || 0);
      }
    } catch (error) {
      message.error('获取薪资列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchSalaries();
  }, [page, pageSize]);

  const handleCalculate = async (values: any) => {
    setLoading(true);
    try {
      const data = {
        ...values,
        period_start: values.period[0].format('YYYY-MM-DD'),
        period_end: values.period[1].format('YYYY-MM-DD'),
      };
      delete data.period;

      const res = await salaryApi.calculateSalary(data);
      if (res.data.code === 201) {
        message.success('薪资计算成功');
        setCreateModalVisible(false);
        form.resetFields();
        fetchSalaries();
      }
    } catch (error: any) {
      message.error(error.response?.data?.message || '计算失败');
    } finally {
      setLoading(false);
    }
  };

  const handlePay = async (id: string) => {
    try {
      const res = await salaryApi.paySalary(id);
      if (res.data.code === 200) {
        message.success('支付成功');
        fetchSalaries();
      }
    } catch (error: any) {
      message.error(error.response?.data?.message || '支付失败');
    }
  };

  const handleExport = async (id: string) => {
    try {
      const res = await salaryApi.exportSalary(id);
      const blob = new Blob([res.data]);
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `salary_${id}.pdf`;
      a.click();
      window.URL.revokeObjectURL(url);
      message.success('导出成功');
    } catch (error) {
      message.error('导出失败');
    }
  };

  const columns = [
    {
      title: '临时工',
      dataIndex: ['temporary', 'real_name'],
      key: 'temporary',
    },
    {
      title: '岗位',
      dataIndex: ['job_posting', 'position'],
      key: 'position',
    },
    {
      title: '周期',
      key: 'period',
      render: (_: any, record: SalaryRecord) => (
        <span>{dayjs(record.period_start).format('MM-DD')} ~ {dayjs(record.period_end).format('MM-DD')}</span>
      ),
    },
    {
      title: '总工时',
      dataIndex: 'total_hours',
      key: 'total_hours',
      render: (hours: number) => `${hours.toFixed(1)} 小时`,
    },
    {
      title: '基本工资',
      dataIndex: 'base_salary',
      key: 'base_salary',
      render: (amount: number) => `¥${amount.toFixed(2)}`,
    },
    {
      title: '加班费',
      dataIndex: 'overtime_pay',
      key: 'overtime_pay',
      render: (amount: number) => amount > 0 ? `¥${amount.toFixed(2)}` : '-',
    },
    {
      title: '扣款',
      dataIndex: 'deductions',
      key: 'deductions',
      render: (amount: number) => amount > 0 ? `¥${amount.toFixed(2)}` : '-',
    },
    {
      title: '实发工资',
      dataIndex: 'total_salary',
      key: 'total_salary',
      render: (amount: number) => (
        <span className="font-semibold text-green-600">¥{amount.toFixed(2)}</span>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => {
        const colorMap: Record<string, string> = {
          pending: 'orange',
          paid: 'green',
        };
        const labelMap: Record<string, string> = {
          pending: '待支付',
          paid: '已支付',
        };
        return <Tag color={colorMap[status]}>{labelMap[status] || status}</Tag>;
      },
    },
    {
      title: '操作',
      key: 'actions',
      render: (_: any, record: SalaryRecord) => (
        <Space>
          <Button size="small" onClick={() => navigate(`/salaries/${record.id}`)}>
            详情
          </Button>
          {record.status === 'pending' && (user?.role === 'employer' || user?.role === 'agent') && (
            <Button size="small" type="primary" icon={<SendOutlined />} onClick={() => handlePay(record.id)}>
              支付
            </Button>
          )}
          <Button size="small" icon={<ExportOutlined />} onClick={() => handleExport(record.id)}>
            导出
          </Button>
        </Space>
      ),
    },
  ];

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-2xl font-bold">薪资管理</h2>
        {(user?.role === 'employer' || user?.role === 'agent') && (
          <Button type="primary" icon={<PlusOutlined />} onClick={() => setCreateModalVisible(true)}>
            计算薪资
          </Button>
        )}
      </div>

      <Card>
        <div className="flex gap-3 mb-4 flex-wrap">
          <RangePicker
            value={dateRange}
            onChange={(dates) => setDateRange(dates)}
          />
          <Select
            placeholder="状态"
            value={status}
            onChange={(val) => setStatus(val)}
            style={{ width: 140 }}
            allowClear
          >
            <Select.Option value="pending">待支付</Select.Option>
            <Select.Option value="paid">已支付</Select.Option>
          </Select>
          <Button type="primary" onClick={fetchSalaries}>搜索</Button>
        </div>

        <Table
          columns={columns}
          dataSource={salaries}
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
      </Card>

      <Modal
        title="计算薪资"
        open={createModalVisible}
        onCancel={() => setCreateModalVisible(false)}
        footer={null}
        width={600}
      >
        <Form form={form} onFinish={handleCalculate} layout="vertical">
          <Form.Item name="temporary_id" label="临时工ID" rules={[{ required: true, message: '请输入临时工ID' }]}>
            <Input placeholder="请输入临时工ID" />
          </Form.Item>
          <Form.Item name="employer_id" label="雇主ID" rules={[{ required: true, message: '请输入雇主ID' }]}>
            <Input placeholder="请输入雇主ID" />
          </Form.Item>
          <Form.Item name="job_id" label="岗位ID" rules={[{ required: true, message: '请输入岗位ID' }]}>
            <Input placeholder="请输入岗位ID" />
          </Form.Item>
          <Form.Item name="period" label="薪资周期" rules={[{ required: true, message: '请选择薪资周期' }]}>
            <RangePicker style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="overtime_rate" label="加班倍率" initialValue={1.5}>
            <InputNumber min={1} step={0.5} style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="deductions" label="扣款金额">
            <InputNumber min={0} style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="deduction_note" label="扣款说明">
            <Input.TextArea rows={2} />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block loading={loading}>
              计算薪资
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};
