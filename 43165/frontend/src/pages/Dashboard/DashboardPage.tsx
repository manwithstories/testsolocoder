import { useEffect, useState } from 'react';
import { Row, Col, Card, Statistic, Spin, message } from 'antd';
import {
  UserOutlined,
  UnorderedListOutlined,
  CalendarOutlined,
  ScanOutlined,
  DollarOutlined,
  RiseOutlined,
} from '@ant-design/icons';
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
  PieChart,
  Pie,
  Cell,
} from 'recharts';
import { statsApi } from '../../services/api';

export const DashboardPage = () => {
  const [loading, setLoading] = useState(true);
  const [overview, setOverview] = useState<any>(null);
  const [activityData, setActivityData] = useState<any>(null);
  const [personnelData, setPersonnelData] = useState<any>(null);

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    setLoading(true);
    try {
      const [overviewRes, activityRes, personnelRes] = await Promise.all([
        statsApi.getOverview(),
        statsApi.getActivityStats(),
        statsApi.getPersonnelStats(),
      ]);

      if (overviewRes.data.code === 200) {
        setOverview(overviewRes.data.data);
      }
      if (activityRes.data.code === 200) {
        setActivityData(activityRes.data.data);
      }
      if (personnelRes.data.code === 200) {
        setPersonnelData(personnelRes.data.data);
      }
    } catch (error) {
      message.error('获取数据失败');
    } finally {
      setLoading(false);
    }
  };

  const COLORS = ['#3b82f6', '#10b981', '#f59e0b', '#ef4444', '#8b5cf6'];

  if (loading) {
    return (
      <div className="flex justify-center items-center h-96">
        <Spin size="large" />
      </div>
    );
  }

  const pieData = personnelData?.credit_distribution?.map((item: any) => ({
    name: item.range,
    value: item.count,
  })) || [];

  return (
    <div>
      <h2 className="text-2xl font-bold mb-6">仪表盘</h2>

      <Row gutter={[16, 16]} className="mb-6">
        <Col xs={24} sm={12} md={8} lg={6}>
          <Card>
            <Statistic
              title="总用户数"
              value={overview?.total_users || 0}
              prefix={<UserOutlined />}
              valueStyle={{ color: '#3b82f6' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={8} lg={6}>
          <Card>
            <Statistic
              title="岗位总数"
              value={overview?.total_jobs || 0}
              prefix={<UnorderedListOutlined />}
              valueStyle={{ color: '#10b981' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={8} lg={6}>
          <Card>
            <Statistic
              title="排班总数"
              value={overview?.total_schedules || 0}
              prefix={<CalendarOutlined />}
              valueStyle={{ color: '#f59e0b' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={8} lg={6}>
          <Card>
            <Statistic
              title="今日签到"
              value={overview?.today_check_ins || 0}
              prefix={<ScanOutlined />}
              valueStyle={{ color: '#8b5cf6' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={8} lg={6}>
          <Card>
            <Statistic
              title="本月薪资"
              value={overview?.monthly_salary || 0}
              prefix={<DollarOutlined />}
              precision={2}
              valueStyle={{ color: '#ef4444' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={8} lg={6}>
          <Card>
            <Statistic
              title="活跃岗位"
              value={overview?.active_jobs || 0}
              prefix={<RiseOutlined />}
              valueStyle={{ color: '#06b6d4' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={8} lg={6}>
          <Card>
            <Statistic
              title="本月新增用户"
              value={overview?.new_users_month || 0}
              prefix={<UserOutlined />}
              valueStyle={{ color: '#ec4899' }}
            />
          </Card>
        </Col>
      </Row>

      <Row gutter={[16, 16]}>
        <Col xs={24} lg={16}>
          <Card title="岗位月度统计">
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={activityData?.by_month || []}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="month" />
                <YAxis />
                <Tooltip />
                <Legend />
                <Bar dataKey="count" fill="#3b82f6" name="岗位数" radius={[4, 4, 0, 0]} />
              </BarChart>
            </ResponsiveContainer>
          </Card>
        </Col>
        <Col xs={24} lg={8}>
          <Card title="信用分分布">
            <ResponsiveContainer width="100%" height={300}>
              <PieChart>
                <Pie
                  data={pieData}
                  cx="50%"
                  cy="50%"
                  labelLine={false}
                  label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
                  outerRadius={80}
                  fill="#8884d8"
                  dataKey="value"
                >
                  {pieData.map((_: any, index: number) => (
                    <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                  ))}
                </Pie>
                <Tooltip />
              </PieChart>
            </ResponsiveContainer>
          </Card>
        </Col>
      </Row>
    </div>
  );
};
