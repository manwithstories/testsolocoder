import { useState, useEffect } from 'react';
import { Card, Row, Col, Spin, Tabs, DatePicker, Select, message } from 'antd';
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
  LineChart,
  Line,
} from 'recharts';
import { statsApi } from '../../services/api';

const COLORS = ['#3b82f6', '#10b981', '#f59e0b', '#ef4444', '#8b5cf6', '#06b6d4'];

export const StatsPage = () => {
  const [loading, setLoading] = useState(true);
  const [activityData, setActivityData] = useState<any>(null);
  const [personnelData, setPersonnelData] = useState<any>(null);
  const [salaryData, setSalaryData] = useState<any>(null);

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    setLoading(true);
    try {
      const [activityRes, personnelRes, salaryRes] = await Promise.all([
        statsApi.getActivityStats(),
        statsApi.getPersonnelStats(),
        statsApi.getSalaryStats(),
      ]);

      if (activityRes.data.code === 200) setActivityData(activityRes.data.data);
      if (personnelRes.data.code === 200) setPersonnelData(personnelRes.data.data);
      if (salaryRes.data.code === 200) setSalaryData(salaryRes.data.data);
    } catch (error) {
      message.error('获取统计数据失败');
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return <div className="flex justify-center items-center h-96"><Spin size="large" /></div>;
  }

  const activityTypeData = activityData?.by_activity_type?.map((item: any) => ({
    name: item.activity_type || '未分类',
    count: item.count,
  })) || [];

  const creditDistData = personnelData?.credit_distribution?.map((item: any) => ({
    name: item.range,
    value: item.count,
  })) || [];

  const salaryByPositionData = salaryData?.salary_by_position?.map((item: any) => ({
    name: item.position,
    total: item.total_salary,
    avg: item.avg_salary,
  })) || [];

  const monthlySalaryData = salaryData?.monthly_salary || [];

  return (
    <div>
      <h2 className="text-2xl font-bold mb-6">数据统计</h2>

      <Tabs
        defaultActiveKey="activity"
        items={[
          {
            key: 'activity',
            label: '活动统计',
            children: (
              <Row gutter={[16, 16]}>
                <Col xs={24} lg={12}>
                  <Card title="活动类型分布">
                    <ResponsiveContainer width="100%" height={300}>
                      <PieChart>
                        <Pie
                          data={activityTypeData}
                          cx="50%"
                          cy="50%"
                          labelLine={false}
                          label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
                          outerRadius={80}
                          fill="#8884d8"
                          dataKey="count"
                        >
                          {activityTypeData.map((_: any, index: number) => (
                            <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                          ))}
                        </Pie>
                        <Tooltip />
                      </PieChart>
                    </ResponsiveContainer>
                  </Card>
                </Col>
                <Col xs={24} lg={12}>
                  <Card title="月度岗位统计">
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
                <Col xs={24}>
                  <Card title="岗位统计概览">
                    <Row gutter={16}>
                      <Col span={6}>
                        <Card size="small">
                          <div className="text-center">
                            <div className="text-3xl font-bold text-blue-500">{activityData?.total_jobs || 0}</div>
                            <div className="text-gray-500">总岗位数</div>
                          </div>
                        </Card>
                      </Col>
                      <Col span={6}>
                        <Card size="small">
                          <div className="text-center">
                            <div className="text-3xl font-bold text-green-500">{activityData?.recruiting_jobs || 0}</div>
                            <div className="text-gray-500">招聘中</div>
                          </div>
                        </Card>
                      </Col>
                      <Col span={6}>
                        <Card size="small">
                          <div className="text-center">
                            <div className="text-3xl font-bold text-gray-500">{activityData?.completed_jobs || 0}</div>
                            <div className="text-gray-500">已完成</div>
                          </div>
                        </Card>
                      </Col>
                      <Col span={6}>
                        <Card size="small">
                          <div className="text-center">
                            <div className="text-3xl font-bold text-orange-500">{activityData?.total_hired || 0}</div>
                            <div className="text-gray-500">已招聘人数</div>
                          </div>
                        </Card>
                      </Col>
                    </Row>
                  </Card>
                </Col>
              </Row>
            ),
          },
          {
            key: 'personnel',
            label: '人员统计',
            children: (
              <Row gutter={[16, 16]}>
                <Col xs={24} lg={12}>
                  <Card title="用户角色分布">
                    <ResponsiveContainer width="100%" height={300}>
                      <PieChart>
                        <Pie
                          data={[
                            { name: '雇主', value: personnelData?.employer_count || 0 },
                            { name: '中介', value: personnelData?.agent_count || 0 },
                            { name: '临时工', value: personnelData?.temporary_count || 0 },
                          ]}
                          cx="50%"
                          cy="50%"
                          labelLine={false}
                          label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
                          outerRadius={80}
                          fill="#8884d8"
                          dataKey="value"
                        >
                          <Cell fill="#3b82f6" />
                          <Cell fill="#10b981" />
                          <Cell fill="#f59e0b" />
                        </Pie>
                        <Tooltip />
                      </PieChart>
                    </ResponsiveContainer>
                  </Card>
                </Col>
                <Col xs={24} lg={12}>
                  <Card title="临时工信用分分布">
                    <ResponsiveContainer width="100%" height={300}>
                      <BarChart data={creditDistData}>
                        <CartesianGrid strokeDasharray="3 3" />
                        <XAxis dataKey="name" />
                        <YAxis />
                        <Tooltip />
                        <Legend />
                        <Bar dataKey="value" fill="#52c41a" name="人数" radius={[4, 4, 0, 0]} />
                      </BarChart>
                    </ResponsiveContainer>
                  </Card>
                </Col>
                <Col xs={24}>
                  <Card title="高信用分临时工TOP10">
                    <div className="space-y-2">
                      {(personnelData?.top_rated_temps || []).slice(0, 10).map((temp: any, index: number) => (
                        <div key={temp.id} className="flex items-center justify-between py-2 px-4 bg-gray-50 rounded">
                          <div className="flex items-center gap-3">
                            <span className={`font-bold ${index < 3 ? 'text-yellow-500' : 'text-gray-400'}`}>
                              #{index + 1}
                            </span>
                            <span>{temp.real_name}</span>
                          </div>
                          <div className="flex items-center gap-4">
                            <span className="text-gray-500">评价: {temp.rating_count}次</span>
                            <span className="text-green-600 font-semibold">信用分: {temp.credit_score}</span>
                          </div>
                        </div>
                      ))}
                    </div>
                  </Card>
                </Col>
              </Row>
            ),
          },
          {
            key: 'salary',
            label: '薪资统计',
            children: (
              <Row gutter={[16, 16]}>
                <Col xs={24} lg={12}>
                  <Card title="月度薪资趋势">
                    <ResponsiveContainer width="100%" height={300}>
                      <LineChart data={monthlySalaryData}>
                        <CartesianGrid strokeDasharray="3 3" />
                        <XAxis dataKey="month" />
                        <YAxis />
                        <Tooltip />
                        <Legend />
                        <Line type="monotone" dataKey="total" stroke="#3b82f6" name="薪资总额" strokeWidth={2} />
                        <Line type="monotone" dataKey="count" stroke="#10b981" name="薪资笔数" strokeWidth={2} />
                      </LineChart>
                    </ResponsiveContainer>
                  </Card>
                </Col>
                <Col xs={24} lg={12}>
                  <Card title="按职位薪资统计">
                    <ResponsiveContainer width="100%" height={300}>
                      <BarChart data={salaryByPositionData}>
                        <CartesianGrid strokeDasharray="3 3" />
                        <XAxis dataKey="name" />
                        <YAxis />
                        <Tooltip />
                        <Legend />
                        <Bar dataKey="total" fill="#3b82f6" name="总薪资" radius={[4, 4, 0, 0]} />
                        <Bar dataKey="avg" fill="#10b981" name="平均薪资" radius={[4, 4, 0, 0]} />
                      </BarChart>
                    </ResponsiveContainer>
                  </Card>
                </Col>
                <Col xs={24}>
                  <Card title="薪资统计概览">
                    <Row gutter={16}>
                      <Col span={8}>
                        <Card size="small">
                          <div className="text-center">
                            <div className="text-3xl font-bold text-green-500">¥{salaryData?.total_paid_salary?.toFixed(2) || '0.00'}</div>
                            <div className="text-gray-500">已支付薪资</div>
                          </div>
                        </Card>
                      </Col>
                      <Col span={8}>
                        <Card size="small">
                          <div className="text-center">
                            <div className="text-3xl font-bold text-orange-500">¥{salaryData?.pending_salary?.toFixed(2) || '0.00'}</div>
                            <div className="text-gray-500">待支付薪资</div>
                          </div>
                        </Card>
                      </Col>
                      <Col span={8}>
                        <Card size="small">
                          <div className="text-center">
                            <div className="text-3xl font-bold text-blue-500">{salaryData?.total_hours?.toFixed(1) || '0'}</div>
                            <div className="text-gray-500">总工时</div>
                          </div>
                        </Card>
                      </Col>
                    </Row>
                  </Card>
                </Col>
              </Row>
            ),
          },
        ]}
      />
    </div>
  );
};
