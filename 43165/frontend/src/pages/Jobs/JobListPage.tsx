import { useState, useEffect } from 'react';
import { Table, Tag, Input, Select, Button, Space, Modal, message, Card, Pagination } from 'antd';
import {
  PlusOutlined,
  SearchOutlined,
  FilterOutlined,
  EnvironmentOutlined,
  DollarOutlined,
  CalendarOutlined,
} from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { jobApi } from '../../services/api';
import { useAuthStore } from '../../context/AuthContext';
import { JobPosting } from '../../types';
import dayjs from 'dayjs';

const { Option } = Select;

export const JobListPage = () => {
  const navigate = useNavigate();
  const { user } = useAuthStore();
  const [jobs, setJobs] = useState<JobPosting[]>([]);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [keyword, setKeyword] = useState('');
  const [activityType, setActivityType] = useState<string | undefined>();
  const [status, setStatus] = useState<string | undefined>();

  const fetchJobs = async () => {
    setLoading(true);
    try {
      const params: any = {
        page,
        page_size: pageSize,
      };
      if (keyword) params.keyword = keyword;
      if (activityType) params.activity_type = activityType;
      if (status) params.status = status;

      const res = await jobApi.getJobs(params);
      if (res.data.code === 200) {
        setJobs(res.data.data.data);
        setTotal(res.data.data.total);
      }
    } catch (error) {
      message.error('获取岗位列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchJobs();
  }, [page, pageSize]);

  const handleSearch = () => {
    setPage(1);
    fetchJobs();
  };

  const columns = [
    {
      title: '岗位名称',
      dataIndex: 'title',
      key: 'title',
      render: (text: string, record: JobPosting) => (
        <a onClick={() => navigate(`/jobs/${record.id}`)}>{text}</a>
      ),
    },
    {
      title: '职位',
      dataIndex: 'position',
      key: 'position',
    },
    {
      title: '地点',
      dataIndex: 'location',
      key: 'location',
      render: (text: string) => (
        <span><EnvironmentOutlined className="mr-1" />{text}</span>
      ),
    },
    {
      title: '薪资',
      dataIndex: 'salary_per_hour',
      key: 'salary_per_hour',
      render: (value: number, record: JobPosting) => (
        <span className="font-semibold text-green-600">
          <DollarOutlined className="mr-1" />¥{value}/小时
        </span>
      ),
    },
    {
      title: '招聘人数',
      key: 'headcount',
      render: (_: any, record: JobPosting) => (
        <span>{record.hired_count}/{record.headcount}</span>
      ),
    },
    {
      title: '时间',
      key: 'date',
      render: (_: any, record: JobPosting) => (
        <span>
          <CalendarOutlined className="mr-1" />
          {dayjs(record.start_date).format('MM-DD')} ~ {dayjs(record.end_date).format('MM-DD')}
        </span>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string, record: JobPosting) => {
        const colorMap: Record<string, string> = {
          recruiting: 'green',
          paused: 'orange',
          completed: 'default',
        };
        return (
          <Space>
            <Tag color={colorMap[status] || 'blue'}>{status}</Tag>
            {record.is_urgent && <Tag color="red">急招</Tag>}
          </Space>
        );
      },
    },
  ];

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-2xl font-bold">岗位列表</h2>
        {(user?.role === 'employer' || user?.role === 'agent') && (
          <Button type="primary" icon={<PlusOutlined />} onClick={() => navigate('/jobs/create')}>
            发布岗位
          </Button>
        )}
      </div>

      <Card>
        <div className="flex gap-3 mb-4 flex-wrap">
          <Input
            placeholder="搜索岗位"
            value={keyword}
            onChange={(e) => setKeyword(e.target.value)}
            onPressEnter={handleSearch}
            prefix={<SearchOutlined />}
            style={{ width: 240 }}
          />
          <Select
            placeholder="活动类型"
            value={activityType}
            onChange={(val) => setActivityType(val)}
            style={{ width: 160 }}
            allowClear
          >
            <Option value="exhibition">展会</Option>
            <Option value="conference">会议</Option>
            <Option value="performance">演出</Option>
            <Option value="promotion">促销</Option>
            <Option value="wedding">婚礼</Option>
            <Option value="other">其他</Option>
          </Select>
          <Select
            placeholder="状态"
            value={status}
            onChange={(val) => setStatus(val)}
            style={{ width: 140 }}
            allowClear
          >
            <Option value="recruiting">招聘中</Option>
            <Option value="paused">已暂停</Option>
            <Option value="completed">已完成</Option>
          </Select>
          <Button type="primary" icon={<SearchOutlined />} onClick={handleSearch}>
            搜索
          </Button>
        </div>

        <Table
          columns={columns}
          dataSource={jobs}
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
