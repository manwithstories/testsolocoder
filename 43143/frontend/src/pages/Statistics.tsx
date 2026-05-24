import React, { useEffect, useState } from 'react';
import { Card, Row, Col, Statistic, DatePicker, Button, Table, message, Spin } from 'antd';
import { CalendarOutlined, UserOutlined, StarOutlined, ExportOutlined } from '@ant-design/icons';
import { statsApi } from '@/api/message';
import { TeacherStats, MonthlyReport } from '@/types';
import dayjs, { Dayjs } from 'dayjs';

const { RangePicker } = DatePicker;

const Statistics: React.FC = () => {
  const [loading, setLoading] = useState(true);
  const [stats, setStats] = useState<TeacherStats | null>(null);
  const [monthlyReports, setMonthlyReports] = useState<MonthlyReport[]>([]);
  const [dateRange, setDateRange] = useState<[Dayjs, Dayjs]>([
    dayjs().subtract(1, 'month'),
    dayjs(),
  ]);

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    setLoading(true);
    try {
      const [statsData, monthlyData] = await Promise.all([
        statsApi.getTeacherStats({
          start_date: dateRange[0]?.format('YYYY-MM-DD'),
          end_date: dateRange[1]?.format('YYYY-MM-DD'),
        }),
        Promise.all([
          statsApi.getMonthlyReport({ year: dayjs().year(), month: dayjs().month() + 1 }),
          statsApi.getMonthlyReport({ year: dayjs().subtract(1, 'month').year(), month: dayjs().subtract(1, 'month').month() + 1 }),
          statsApi.getMonthlyReport({ year: dayjs().subtract(2, 'month').year(), month: dayjs().subtract(2, 'month').month() + 1 }),
        ]),
      ]);

      setStats(statsData);
      setMonthlyReports(monthlyData);
    } catch (error: any) {
      message.error(error.message || '获取统计数据失败');
    } finally {
      setLoading(false);
    }
  };

  const handleDateChange = (dates: any) => {
    if (dates && dates.length === 2) {
      setDateRange(dates);
    }
  };

  const handleExport = async () => {
    try {
      await statsApi.exportReport();
      message.success('报表导出成功');
    } catch (error: any) {
      message.error(error.message || '导出失败');
    }
  };

  const columns = [
    {
      title: '月份',
      dataIndex: 'month',
      key: 'month',
    },
    {
      title: '授课时长',
      dataIndex: 'teaching_hours',
      key: 'teaching_hours',
      render: (value: number) => `${value.toFixed(1)} 小时`,
    },
    {
      title: '学员数量',
      dataIndex: 'student_count',
      key: 'student_count',
    },
    {
      title: '预约次数',
      dataIndex: 'bookings',
      key: 'bookings',
    },
    {
      title: '收入',
      dataIndex: 'income',
      key: 'income',
      render: (value: number) => `¥${value.toFixed(2)}`,
    },
  ];

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
        title="数据统计"
        extra={
          <div style={{ display: 'flex', gap: 8 }}>
            <RangePicker
              value={dateRange}
              onChange={handleDateChange}
              style={{ marginRight: 8 }}
            />
            <Button icon={<ExportOutlined />} onClick={handleExport}>
              导出报表
            </Button>
          </div>
        }
      >
        <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
          <Col span={6}>
            <Card>
              <Statistic
                title="授课时长"
                value={stats?.teaching_hours || 0}
                precision={1}
                suffix="小时"
                prefix={<CalendarOutlined />}
                valueStyle={{ color: '#1890ff' }}
              />
            </Card>
          </Col>
          <Col span={6}>
            <Card>
              <Statistic
                title="学员数量"
                value={stats?.student_count || 0}
                prefix={<UserOutlined />}
                valueStyle={{ color: '#52c41a' }}
              />
            </Card>
          </Col>
          <Col span={6}>
            <Card>
              <Statistic
                title="总收入"
                value={stats?.total_income || 0}
                precision={2}
                prefix="¥"
                valueStyle={{ color: '#faad14' }}
              />
            </Card>
          </Col>
          <Col span={6}>
            <Card>
              <Statistic
                title="平均评分"
                value={stats?.avg_rating || 0}
                precision={1}
                suffix="/ 5"
                prefix={<StarOutlined />}
                valueStyle={{ color: '#722ed1' }}
              />
            </Card>
          </Col>
        </Row>

        <Card title="月度报告" style={{ marginBottom: 16 }}>
          <Table
            dataSource={monthlyReports}
            columns={columns}
            rowKey="month"
            pagination={false}
            size="small"
          />
        </Card>
      </Card>
    </div>
  );
};

export default Statistics;
