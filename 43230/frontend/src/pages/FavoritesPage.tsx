import React, { useEffect, useState } from 'react';
import {
  Card,
  Table,
  Button,
  Space,
  Typography,
  Empty,
  Image,
  Rate,
  Tag,
  message,
} from 'antd';
import {
  EyeOutlined,
  HeartOutlined,
  ShoppingCartOutlined,
  DownloadOutlined,
  StarOutlined,
} from '@ant-design/icons';
import { Link, useNavigate } from 'react-router-dom';
import { modelApi } from '@/services/api';
import { Model3D } from '@/types';
import { formatPrice, generateModelThumbnail } from '@/utils/format';

const { Title, Text } = Typography;

const FavoritesPage: React.FC = () => {
  const [models, setModels] = useState<Model3D[]>([]);
  const [loading, setLoading] = useState(false);
  const [page, setPage] = useState(1);
  const [pageSize] = useState(20);
  const [total, setTotal] = useState(0);
  const navigate = useNavigate();

  const fetchData = async () => {
    setLoading(true);
    try {
      const response = await modelApi.getFavorites({ page, page_size: pageSize });
      setModels(response.data.data || []);
      setTotal(response.data.total || 0);
    } catch (error) {
      console.error('Failed to fetch favorites:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, [page]);

  const handleRemoveFavorite = async (modelId: string) => {
    try {
      await modelApi.removeFavorite(modelId);
      message.success('已取消收藏');
      fetchData();
    } catch (error: any) {
      message.error(error.response?.data?.error || '操作失败');
    }
  };

  const columns = [
    {
      title: '模型',
      key: 'model',
      render: (_: any, record: Model3D) => (
        <div className="flex items-center gap-3">
          <div className="w-16 h-16 bg-gray-100 rounded overflow-hidden">
            <Image
              src={record.thumbnail_url || generateModelThumbnail(record.id)}
              alt={record.title}
              className="w-full h-full object-cover"
              preview={false}
            />
          </div>
          <div>
            <Link to={`/models/${record.id}`} className="font-medium hover:text-blue-600">
              {record.title}
            </Link>
            <div className="text-sm text-gray-500">
              设计师: {record.designer?.designer_profile?.nickname || record.designer?.username}
            </div>
          </div>
        </div>
      ),
    },
    {
      title: '分类',
      dataIndex: 'category',
      key: 'category',
      render: (cat: string) => <Tag color="blue">{cat}</Tag>,
    },
    {
      title: '价格',
      dataIndex: 'price',
      key: 'price',
      render: (val: number) => (
        <div className="text-red-500 font-medium">{formatPrice(val)}</div>
      ),
    },
    {
      title: '评分',
      dataIndex: 'rating',
      key: 'rating',
      render: (val: number) => (
        <Space>
          <Rate disabled value={val} size="small" />
          <span>{val.toFixed(1)}</span>
        </Space>
      ),
    },
    {
      title: '数据',
      key: 'stats',
      render: (_: any, record: Model3D) => (
        <div className="text-sm space-y-1">
          <div>
            <DownloadOutlined className="mr-1 text-gray-400" />
            {record.download_count} 下载
          </div>
          <div>
            <ShoppingCartOutlined className="mr-1 text-gray-400" />
            {record.purchase_count} 购买
          </div>
        </div>
      ),
    },
    {
      title: '操作',
      key: 'actions',
      render: (_: any, record: Model3D) => (
        <Space size="small">
          <Link to={`/models/${record.id}`}>
            <Button type="primary" size="small">
              <EyeOutlined /> 查看
            </Button>
          </Link>
          <Button
            size="small"
            danger
            onClick={() => handleRemoveFavorite(record.id)}
          >
            <HeartOutlined /> 取消收藏
          </Button>
        </Space>
      ),
    },
  ];

  return (
    <div className="space-y-6">
      <Card
        title={
          <Space>
            <HeartOutlined className="text-red-500" />
            我的收藏
          </Space>
        }
      >
        <Table
          rowKey="id"
          columns={columns}
          dataSource={models}
          loading={loading}
          pagination={{
            current: page,
            pageSize,
            total,
            onChange: setPage,
          }}
          locale={{ emptyText: <Empty description="暂无收藏" /> }}
        />
      </Card>
    </div>
  );
};

export default FavoritesPage;
