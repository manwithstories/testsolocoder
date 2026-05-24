import { useState, useEffect } from 'react';
import { Table, Tag, Card, DatePicker, Select, Button, Space, message, Statistic, Row, Col } from 'antd';
import { ExportOutlined, SearchOutlined } from '@ant-design/icons';
import { checkInApi } from '../../services/api';
import { useAuthStore } from '../../context/AuthContext';
import { CheckIn } from '../../types';
import dayjs from 'dayjs';

const { RangePicker } = DatePicker;

export const CheckInRecordsPage = () => {
  const { user } = useAuthStore();
  const [records, setRecords] = useState<CheckIn[]>([]);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [dateRange, setDateRange] = useState<any>(null);
  const [status, setStatus] = useState<string | undefined>();
  const [stats, setStats] = useState<any>(null);

  const fetchRecords = async () => {
    setLoading(true);
    try {
      const params: any = { page, page_size: pageSize };
      if (dateRange && dateRange.length === 2) {
        params.start_date = dateRange[0].format('YYYY-MM-DD');
        params.end_date = dateRange[1].format('YYYY-MM-DD');
      }
      if (status) params.status = status;

      const res = await checkInApi.getCheckIns(params);
      if (res.data.code === 200) {
        setRecords(res.data.data.data || []);
        setTotal(res.data.data.total || 0);
      }
    } catch (error) {
      message.error('获取签到记录失败');
    } finally {
      setLoading(false);
    }
  };

  const fetchStats = async () => {
    try {
      const res = await checkInApi.getStats();
      if (res.data.code === 200) {
        setStats(res.data.data);
      }
    } catch (error) {
      console.error('获取统计数据失败');
    }
  };

  useEffect(() => {
    fetchRecords();
    fetchStats();
  }, [page, pageSize]);

  const handleSearch = () => {
    setPage(1);
    fetchRecords();
  };

  const handleExport = () => {
    message.info('导出功能开发中');
  };

  const columns = [
    {
      title: '临时工',
      dataIndex: ['temporary', 'real_name'],
      key: 'temporary',
    },
    {
      title: '岗位',
      dataIndex: ['schedule', 'job_posting', 'position'],
      key: 'position',
    },
    {
      title: '签到方式',
      dataIndex: 'check_in_type',
      key: 'check_in_type',
      render: (type: string) => {
        const typeMap: Record<string, string> = {
          qr: '扫码',
          location: '定位',
          face: '人脸识别',
        };
        return typeMap[type] || type;
      },
    },
    {
      title: '签到时间',
      dataIndex: 'check_in_time',
      key: 'check_in_time',
      render: (text: string) => dayjs(text).format('YYYY-MM-DD HH:mm'),
    },
    {
      title: '签退时间',
      dataIndex: 'check_out_time',
      key: 'check_out_time',
      render: (text: string | null) => text ? dayjs(text).format('YYYY-MM-DD HH:mm') : '-',
    },
    {
      title: '工时',
      dataIndex: 'work_hours',
      key: 'work_hours',
      render: (hours: number) => hours ? `${hours.toFixed(1)} 小时` : '-',
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => {
        const colorMap: Record<string, string> = {
          checked_in: 'orange',
          checked_out: 'green',
        };
        const labelMap: Record<string, string> = {
          checked_in: '签到中',
          checked_out: '已签退',
        };
        return <Tag color={colorMap[status]}>{labelMap[status] || status}</Tag>;
      },
    },
    {
      title: '人脸识别',
      dataIndex: 'face_verified',
      key: 'face_verified',
      render: (verified: boolean) => (
        <Tag color={verified ? 'green' : 'default'}>{verified ? '已验证' : '未验证'}</Tag>
      ),
    },
  ];

  return (
    <div>
      <h2 className="text-2xl font-bold mb-6">签到记录</h2>

      {stats && (
        <Row gutter={[16, 16]} className="mb-6">
          <Col xs={12} sm={6}>
            <Card>
              <Statistic title="总签到次数" value={stats.total_check_ins} />
            </Card>
          </Col>
          <Col xs={12} sm={6}>
            <Card>
              <Statistic title="当前签到中" value={stats.currently_checked_in} valueStyle={{ color: '#faad14' }} />
            </Card>
          </Col>
          <Col xs={12} sm={6}>
            <Card>
              <Statistic title="已完成签退" value={stats.completed_check_outs} valueStyle={{ color: '#52c41a' }} />
            </Card>
          </Col>
          <Col xs={12} sm={6}>
            <Card>
              <Statistic title="总工时" value={stats.total_work_hours} suffix="小时" precision={2} />
            </Card>
          </Col>
        </Row>
      )}

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
            <Select.Option value="checked_in">签到中</Select.Option>
            <Select.Option value="checked_out">已签退</Select.Option>
          </Select>
          <Button type="primary" icon={<SearchOutlined />} onClick={handleSearch}>
            搜索
          </Button>
          <Button icon={<ExportOutlined />} onClick={handleExport}>
            导出
          </Button>
        </div>

        <Table
          columns={columns}
          dataSource={records}
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
    </div>
  );
};
