import { useState, useEffect } from 'react';
import { Table, Tag, message } from 'antd';
import { useNavigate } from 'react-router-dom';
import { jobApi } from '../../services/api';
import { JobApplication } from '../../types';
import dayjs from 'dayjs';

export const MyApplicationsPage = () => {
  const navigate = useNavigate();
  const [applications, setApplications] = useState<JobApplication[]>([]);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);

  const fetchApplications = async () => {
    setLoading(true);
    try {
      const res = await jobApi.getMyApplications({ page, page_size: pageSize });
      if (res.data.code === 200) {
        setApplications(res.data.data.data);
        setTotal(res.data.data.total);
      }
    } catch (error) {
      message.error('获取申请列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchApplications();
  }, [page, pageSize]);

  const columns = [
    {
      title: '岗位',
      key: 'job',
      render: (_: any, record: JobApplication) => (
        <a onClick={() => navigate(`/jobs/${record.job_id}`)}>
          {record.job_posting?.title}
        </a>
      ),
    },
    {
      title: '雇主',
      dataIndex: ['job_posting', 'employer', 'real_name'],
      key: 'employer',
    },
    {
      title: '申请留言',
      dataIndex: 'message',
      key: 'message',
      ellipsis: true,
    },
    {
      title: '申请时间',
      dataIndex: 'applied_at',
      key: 'applied_at',
      render: (text: string) => dayjs(text).format('YYYY-MM-DD HH:mm'),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => {
        const colorMap: Record<string, string> = {
          pending: 'orange',
          approved: 'green',
          rejected: 'red',
        };
        const labelMap: Record<string, string> = {
          pending: '待审核',
          approved: '已通过',
          rejected: '已拒绝',
        };
        return <Tag color={colorMap[status]}>{labelMap[status] || status}</Tag>;
      },
    },
    {
      title: '审核时间',
      dataIndex: 'reviewed_at',
      key: 'reviewed_at',
      render: (text: string | null) => text ? dayjs(text).format('YYYY-MM-DD HH:mm') : '-',
    },
  ];

  return (
    <div>
      <h2 className="text-2xl font-bold mb-6">我的申请</h2>
      <Table
        columns={columns}
        dataSource={applications}
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
    </div>
  );
};
