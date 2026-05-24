import { useState, useEffect } from 'react';
import { Table, Tag, Card, Rate, Avatar, List, Button, Space, Modal, Form, Select, message } from 'antd';
import { StarOutlined } from '@ant-design/icons';
import { evaluationApi } from '../../services/api';
import { Evaluation } from '../../types';
import dayjs from 'dayjs';

export const EvaluationListPage = () => {
  const [evaluations, setEvaluations] = useState<Evaluation[]>([]);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [evalType, setEvalType] = useState<string | undefined>();
  const [minRating, setMinRating] = useState<number | undefined>();

  const fetchEvaluations = async () => {
    setLoading(true);
    try {
      const params: any = { page, page_size: pageSize };
      if (evalType) params.type = evalType;
      if (minRating) params.min_rating = minRating;

      const res = await evaluationApi.getEvaluations(params);
      if (res.data.code === 200) {
        setEvaluations(res.data.data.data || []);
        setTotal(res.data.data.total || 0);
      }
    } catch (error) {
      message.error('获取评价列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchEvaluations();
  }, [page, pageSize, evalType, minRating]);

  const columns = [
    {
      title: '评价人',
      dataIndex: ['from_user', 'real_name'],
      key: 'from_user',
      render: (name: string, record: Evaluation) => (
        record.is_anonymous ? '匿名用户' : name
      ),
    },
    {
      title: '被评价人',
      dataIndex: ['to_user', 'real_name'],
      key: 'to_user',
    },
    {
      title: '岗位',
      dataIndex: ['job_posting', 'position'],
      key: 'position',
    },
    {
      title: '评分',
      dataIndex: 'rating',
      key: 'rating',
      render: (rating: number) => <Rate disabled defaultValue={rating} />,
    },
    {
      title: '评价内容',
      dataIndex: 'content',
      key: 'content',
      ellipsis: true,
    },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
      render: (type: string) => {
        const labelMap: Record<string, string> = {
          employer_to_temp: '雇主评临时工',
          temp_to_employer: '临时工评雇主',
        };
        return <Tag>{labelMap[type] || type}</Tag>;
      },
    },
    {
      title: '评价时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (text: string) => dayjs(text).format('YYYY-MM-DD HH:mm'),
    },
  ];

  return (
    <div>
      <h2 className="text-2xl font-bold mb-6">评价管理</h2>

      <Card>
        <div className="flex gap-3 mb-4">
          <Select
            placeholder="评价类型"
            value={evalType}
            onChange={(val) => setEvalType(val)}
            style={{ width: 180 }}
            allowClear
          >
            <Select.Option value="employer_to_temp">雇主评临时工</Select.Option>
            <Select.Option value="temp_to_employer">临时工评雇主</Select.Option>
          </Select>
          <Select
            placeholder="最低评分"
            value={minRating}
            onChange={(val) => setMinRating(val)}
            style={{ width: 140 }}
            allowClear
          >
            <Select.Option value={5}>5星</Select.Option>
            <Select.Option value={4}>4星以上</Select.Option>
            <Select.Option value={3}>3星以上</Select.Option>
            <Select.Option value={2}>2星以上</Select.Option>
            <Select.Option value={1}>1星以上</Select.Option>
          </Select>
        </div>

        <Table
          columns={columns}
          dataSource={evaluations}
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
