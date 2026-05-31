import React, { useEffect, useState } from 'react';
import {
  Card,
  Row,
  Col,
  Statistic,
  Typography,
  DatePicker,
  Select,
  Button,
  Space,
  Tabs,
  Table,
  Tag,
  Empty,
  List,
  Avatar,
  Progress,
} from 'antd';
import {
  BarChartOutlined,
  LineChartOutlined,
  PieChartOutlined,
  RiseOutlined,
  ShoppingCartOutlined,
  AppstoreOutlined,
  PrinterOutlined,
  DownloadOutlined,
  ExportOutlined,
  FileExcelOutlined,
} from '@ant-design/icons';
import ReactECharts from 'echarts-for-react';
import { statsApi, modelApi } from '@/services/api';
import { formatPrice, formatFileSize } from '@/utils/format';
import dayjs from 'dayjs';

const { Title, Text } = Typography;
const { RangePicker } = DatePicker;
const { Option } = Select;
const { TabPane } = Tabs;

interface StatsData {
  period: {
    start_date: string;
    end_date: string;
  };
  order_stats: {
    total_orders: number;
    total_revenue: number;
    status_stats: Record<string, number>;
  };
  model_stats: {
    total_models: number;
    total_downloads: number;
    total_purchases: number;
    total_revenue: number;
  };
  material_stats: {
    material_usage: Array<{
      material_id: string;
      material_name: string;
      total_weight: number;
      order_count: number;
    }>;
  };
  hot_models: any[];
}

const StatsPage: React.FC = () => {
  const [stats, setStats] = useState<StatsData | null>(null);
  const [loading, setLoading] = useState(false);
  const [dateRange, setDateRange] = useState<[dayjs.Dayjs, dayjs.Dayjs]>([
    dayjs().subtract(30, 'day'),
    dayjs(),
  ]);
  const [exportType, setExportType] = useState('csv');

  const fetchStats = async () => {
    setLoading(true);
    try {
      const params = {
        start_date: dateRange[0].format('YYYY-MM-DD'),
        end_date: dateRange[1].format('YYYY-MM-DD'),
      };
      const response = await statsApi.getPlatformStats(params);
      setStats(response.data);
    } catch (error) {
      console.error('Failed to fetch stats:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchStats();
  }, [dateRange]);

  const handleExport = async () => {
    try {
      const params = {
        start_date: dateRange[0].format('YYYY-MM-DD'),
        end_date: dateRange[1].format('YYYY-MM-DD'),
        type: exportType,
      };
      const response = await statsApi.exportStats(params);
      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', `stats_report.${exportType}`);
      document.body.appendChild(link);
      link.click();
    } catch (error: any) {
      console.error('Export failed:', error);
    }
  };

  const getRevenueChartOption = () => {
    if (!stats) return {};

    return {
      tooltip: {
        trigger: 'axis',
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true,
      },
      xAxis: {
        type: 'category',
        boundaryGap: false,
        data: ['1月', '2月', '3月', '4月', '5月', '6月'],
      },
      yAxis: {
        type: 'value',
        axisLabel: {
          formatter: '¥{value}',
        },
      },
      series: [
        {
          name: '订单收入',
          type: 'line',
          smooth: true,
          data: [12000, 19000, 15000, 22000, 28000, stats?.order_stats.total_revenue || 0],
          areaStyle: {
            color: {
              type: 'linear',
              x: 0,
              y: 0,
              x2: 0,
              y2: 1,
              colorStops: [
                { offset: 0, color: 'rgba(59, 130, 246, 0.3)' },
                { offset: 1, color: 'rgba(59, 130, 246, 0.05)' },
              ],
            },
          },
          lineStyle: {
            color: '#3b82f6',
            width: 3,
          },
          itemStyle: {
            color: '#3b82f6',
          },
        },
      ],
    };
  };

  const getOrderStatusChartOption = () => {
    if (!stats) return {};

    const statusStats = stats.order_stats.status_stats || {};
    const statusMap: Record<string, string> = {
      pending: '待付款',
      paid: '已付款',
      printing: '打印中',
      quality_check: '质检中',
      shipped: '已发货',
      delivered: '已送达',
      completed: '已完成',
      cancelled: '已取消',
      refunded: '已退款',
    };

    const data = Object.entries(statusStats).map(([key, value]) => ({
      name: statusMap[key] || key,
      value,
    }));

    return {
      tooltip: {
        trigger: 'item',
        formatter: '{a} <br/>{b}: {c} ({d}%)',
      },
      legend: {
        orient: 'vertical',
        left: 'left',
      },
      series: [
        {
          name: '订单状态',
          type: 'pie',
          radius: ['40%', '70%'],
          avoidLabelOverlap: false,
          itemStyle: {
            borderRadius: 10,
            borderColor: '#fff',
            borderWidth: 2,
          },
          label: {
            show: false,
            position: 'center',
          },
          emphasis: {
            label: {
              show: true,
              fontSize: 20,
              fontWeight: 'bold',
            },
          },
          labelLine: {
            show: false,
          },
          data,
        },
      ],
    };
  };

  const getMaterialUsageChartOption = () => {
    if (!stats) return {};

    const materials = stats.material_stats.material_usage || [];
    return {
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'shadow',
        },
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true,
      },
      xAxis: {
        type: 'category',
        data: materials.map((m) => m.material_name),
        axisLabel: {
          rotate: 45,
        },
      },
      yAxis: {
        type: 'value',
        axisLabel: {
          formatter: '{value}g',
        },
      },
      series: [
        {
          name: '使用量',
          type: 'bar',
          data: materials.map((m) => m.total_weight),
          itemStyle: {
            color: {
              type: 'linear',
              x: 0,
              y: 0,
              x2: 0,
              y2: 1,
              colorStops: [
                { offset: 0, color: '#10b981' },
                { offset: 1, color: '#059669' },
              ],
            },
            borderRadius: [8, 8, 0, 0],
          },
        },
      ],
    };
  };

  if (!stats) {
    return (
      <Card>
        <div className="text-center py-12">
          <div className="animate-spin w-8 h-8 border-4 border-blue-500 border-t-transparent rounded-full mx-auto"></div>
        </div>
      </Card>
    );
  }

  return (
    <div className="space-y-6">
      <Card>
        <div className="flex justify-between items-center flex-wrap gap-4">
          <Title level={3} className="!mb-0">
            <BarChartOutlined className="mr-2" />
            数据统计
          </Title>
          <Space wrap>
            <RangePicker
              value={dateRange}
              onChange={(dates) => dates && setDateRange(dates as any)}
            />
            <Select
              value={exportType}
              onChange={setExportType}
              style={{ width: 120 }}
            >
              <Option value="csv">CSV 格式</Option>
              <Option value="json">JSON 格式</Option>
            </Select>
            <Button type="primary" icon={<ExportOutlined />} onClick={handleExport}>
              导出报表
            </Button>
          </Space>
        </div>

        <div className="mt-2 text-sm text-gray-500">
          统计周期: {stats.period.start_date} 至 {stats.period.end_date}
        </div>
      </Card>

      {/* 核心指标 */}
      <Row gutter={16}>
        <Col xs={12} md={6}>
          <Card>
            <Statistic
              title="订单总数"
              value={stats.order_stats.total_orders}
              prefix={<ShoppingCartOutlined />}
              valueStyle={{ color: '#3b82f6' }}
            />
          </Card>
        </Col>
        <Col xs={12} md={6}>
          <Card>
            <Statistic
              title="订单总营收"
              value={stats.order_stats.total_revenue}
              precision={2}
              prefix="¥"
              valueStyle={{ color: '#10b981' }}
            />
          </Card>
        </Col>
        <Col xs={12} md={6}>
          <Card>
            <Statistic
              title="新增模型"
              value={stats.model_stats.total_models}
              prefix={<AppstoreOutlined />}
              valueStyle={{ color: '#8b5cf6' }}
            />
          </Card>
        </Col>
        <Col xs={12} md={6}>
          <Card>
            <Statistic
              title="模型总下载"
              value={stats.model_stats.total_downloads}
              prefix={<DownloadOutlined />}
              valueStyle={{ color: '#f59e0b' }}
            />
          </Card>
        </Col>
      </Row>

      {/* 图表区域 */}
      <Tabs defaultActiveKey="revenue">
        <TabPane
          tab={
            <span>
              <LineChartOutlined /> 营收趋势
            </span>
          }
          key="revenue"
        >
          <Card>
            <ReactECharts option={getRevenueChartOption()} style={{ height: 400 }} />
          </Card>
        </TabPane>

        <TabPane
          tab={
            <span>
              <PieChartOutlined /> 订单状态分布
            </span>
          }
          key="orders"
        >
          <Card>
            <ReactECharts option={getOrderStatusChartOption()} style={{ height: 400 }} />
          </Card>
        </TabPane>

        <TabPane
          tab={
            <span>
              <BarChartOutlined /> 材料消耗
            </span>
          }
          key="materials"
        >
          <Card>
            <ReactECharts option={getMaterialUsageChartOption()} style={{ height: 400 }} />
          </Card>
        </TabPane>
      </Tabs>

      {/* 热门模型 */}
      <Card
        title={
          <Space>
            <RiseOutlined className="text-red-500" />
            热门模型排行
          </Space>
        }
      >
        {stats.hot_models.length === 0 ? (
          <Empty description="暂无数据" />
        ) : (
          <List
            dataSource={stats.hot_models}
            renderItem={(model, index) => (
              <List.Item key={model.id}>
                <List.Item.Meta
                  avatar={
                    <div
                      className={`w-8 h-8 rounded-full flex items-center justify-center text-white font-bold ${
                        index === 0
                          ? 'bg-yellow-500'
                          : index === 1
                          ? 'bg-gray-400'
                          : index === 2
                          ? 'bg-orange-600'
                          : 'bg-gray-200 text-gray-600'
                      }`}
                    >
                      {index + 1}
                    </div>
                  }
                  title={
                    <Space>
                      <span className="font-medium">{model.title}</span>
                      <Tag color="blue">下载 {model.download_count}</Tag>
                      <Tag color="green">浏览 {model.view_count}</Tag>
                    </Space>
                  }
                  description={
                    <Space size="large">
                      <Text type="secondary">
                        设计: {model.designer?.designer_profile?.nickname || model.designer?.username}
                      </Text>
                      <Text type="secondary">评分: {model.rating?.toFixed(1)}</Text>
                      <Text strong className="text-red-500">
                        {formatPrice(model.price)}
                      </Text>
                    </Space>
                  }
                />
                <Progress
                  percent={Math.min((model.download_count / stats.hot_models[0].download_count) * 100, 100)}
                  showInfo={false}
                  style={{ width: 200 }}
                  strokeColor="#3b82f6"
                />
              </List.Item>
            )}
          />
        )}
      </Card>

      {/* 材料消耗明细 */}
      <Card
        title={
          <Space>
            <PrinterOutlined />
            材料消耗明细
          </Space>
        }
      >
        <Table
          rowKey="material_id"
          dataSource={stats.material_stats.material_usage || []}
          pagination={false}
          columns={[
            {
              title: '材料名称',
              dataIndex: 'material_name',
              key: 'material_name',
            },
            {
              title: '使用重量',
              dataIndex: 'total_weight',
              key: 'total_weight',
              render: (val: number) => `${val.toFixed(2)}g`,
              sorter: (a: any, b: any) => a.total_weight - b.total_weight,
            },
            {
              title: '订单数量',
              dataIndex: 'order_count',
              key: 'order_count',
              sorter: (a: any, b: any) => a.order_count - b.order_count,
            },
          ]}
        />
      </Card>
    </div>
  );
};

export default StatsPage;
