import React, { useEffect, useState } from 'react';
import {
  Card,
  Table,
  Tag,
  Button,
  Space,
  Typography,
  Select,
  Empty,
  Image,
  Rate,
  Modal,
  Form,
  Input,
  message,
  Descriptions,
} from 'antd';
import {
  EyeOutlined,
  EditOutlined,
  DeleteOutlined,
  PlusOutlined,
  UploadOutlined,
  AppstoreOutlined,
  StarOutlined,
  DownloadOutlined,
  ShoppingCartOutlined,
} from '@ant-design/icons';
import { Link, useNavigate } from 'react-router-dom';
import { modelApi } from '@/services/api';
import { Model3D, ModelPurchase } from '@/types';
import {
  formatPrice,
  formatFileSize,
  generateModelThumbnail,
} from '@/utils/format';
import { formatDate } from '@/utils/date';

const { Title, Text } = Typography;
const { Option } = Select;

const MyModelsPage: React.FC = () => {
  const [models, setModels] = useState<Model3D[]>([]);
  const [purchases, setPurchases] = useState<ModelPurchase[]>([]);
  const [favorites, setFavorites] = useState<Model3D[]>([]);
  const [loading, setLoading] = useState(false);
  const [activeTab, setActiveTab] = useState<'my' | 'purchased' | 'favorite'>('my');
  const [page, setPage] = useState(1);
  const [pageSize] = useState(10);
  const [total, setTotal] = useState(0);
  const navigate = useNavigate();

  const fetchData = async () => {
    setLoading(true);
    try {
      if (activeTab === 'my') {
        const response = await modelApi.getMyModels({ page, page_size: pageSize });
        setModels(response.data.data || []);
        setTotal(response.data.total || 0);
      } else if (activeTab === 'purchased') {
        const response = await modelApi.getPurchases({ page, page_size: pageSize });
        setPurchases(response.data.data || []);
        setTotal(response.data.total || 0);
      } else if (activeTab === 'favorite') {
        const response = await modelApi.getFavorites({ page, page_size: pageSize });
        setFavorites(response.data.data || []);
        setTotal(response.data.total || 0);
      }
    } catch (error) {
      console.error('Failed to fetch data:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, [activeTab, page]);

  const handleDelete = (id: string) => {
    Modal.confirm({
      title: '确认删除',
      content: '删除模型后将无法恢复，确定要删除吗？',
      okType: 'danger',
      onOk: async () => {
        try {
          await modelApi.delete(id);
          message.success('删除成功');
          fetchData();
        } catch (error: any) {
          message.error(error.response?.data?.error || '删除失败');
        }
      },
    });
  };

  const handleDownload = async (modelId: string) => {
    try {
      const response = await modelApi.download(modelId);
      window.open(response.data.file_url, '_blank');
      message.success('开始下载');
    } catch (error: any) {
      message.error(error.response?.data?.error || '下载失败');
    }
  };

  const myModelColumns = [
    {
      title: '缩略图',
      dataIndex: 'thumbnail_url',
      key: 'thumbnail',
      width: 100,
      render: (url: string, record: Model3D) => (
        <div className="w-16 h-16 bg-gray-100 rounded overflow-hidden">
          <Image
            src={url || generateModelThumbnail(record.id)}
            alt={record.title}
            className="w-full h-full object-cover"
            preview={false}
          />
        </div>
      ),
    },
    {
      title: '模型信息',
      key: 'info',
      render: (_: any, record: Model3D) => (
        <div>
          <Link to={`/models/${record.id}`} className="font-medium hover:text-blue-600">
            {record.title}
          </Link>
          <div className="text-sm text-gray-500">
            {record.category} · 版本 {record.version}
          </div>
          <div className="text-xs text-gray-400">
            创建于 {formatDate(record.created_at)}
          </div>
        </div>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => {
        const statusMap: Record<string, { text: string; color: string }> = {
          draft: { text: '草稿', color: 'default' },
          published: { text: '已发布', color: 'green' },
          rejected: { text: '已拒绝', color: 'red' },
          banned: { text: '已封禁', color: 'red' },
        };
        const info = statusMap[status] || statusMap.draft;
        return <Tag color={info.color}>{info.text}</Tag>;
      },
    },
    {
      title: '价格',
      dataIndex: 'price',
      key: 'price',
      render: (price: number) => (
        <div className="text-red-500 font-bold">{formatPrice(price)}</div>
      ),
    },
    {
      title: '数据',
      key: 'stats',
      render: (_: any, record: Model3D) => (
        <div className="text-sm space-y-1">
          <div>
            <EyeOutlined className="mr-1 text-gray-400" />
            {record.view_count}
            <DownloadOutlined className="ml-3 mr-1 text-gray-400" />
            {record.download_count}
          </div>
          <div>
            <ShoppingCartOutlined className="mr-1 text-gray-400" />
            {record.purchase_count}
            <StarOutlined className="ml-3 mr-1 text-yellow-500" />
            {record.rating.toFixed(1)}
          </div>
        </div>
      ),
    },
    {
      title: '操作',
      key: 'actions',
      render: (_: any, record: Model3D) => (
        <Space size="small" direction="vertical">
          <Link to={`/models/${record.id}`}>
            <Button type="link" size="small">
              <EyeOutlined /> 预览
            </Button>
          </Link>
          <Button type="link" size="small">
            <EditOutlined /> 编辑
          </Button>
          <Button
            type="link"
            size="small"
            danger
            onClick={() => handleDelete(record.id)}
          >
            <DeleteOutlined /> 删除
          </Button>
        </Space>
      ),
    },
  ];

  const purchasedColumns = [
    {
      title: '模型',
      key: 'model',
      render: (_: any, record: ModelPurchase) => (
        <div className="flex items-center gap-3">
          <div className="w-16 h-16 bg-gray-100 rounded overflow-hidden">
            <Image
              src={record.model?.thumbnail_url || generateModelThumbnail(record.model_id)}
              alt={record.model?.title}
              className="w-full h-full object-cover"
              preview={false}
            />
          </div>
          <div>
            <Link to={`/models/${record.model_id}`} className="font-medium hover:text-blue-600">
              {record.model?.title}
            </Link>
            <div className="text-sm text-gray-500">
              设计师: {record.model?.designer?.username}
            </div>
          </div>
        </div>
      ),
    },
    {
      title: '购买类型',
      dataIndex: 'purchase_type',
      key: 'purchase_type',
      render: (type: string) => (
        <Tag color={type === 'subscription' ? 'purple' : 'blue'}>
          {type === 'subscription' ? '订阅' : '按件'}
        </Tag>
      ),
    },
    {
      title: '金额',
      dataIndex: 'amount',
      key: 'amount',
      render: (val: number) => (
        <div className="text-red-500 font-medium">{formatPrice(val)}</div>
      ),
    },
    {
      title: '购买时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date: string) => formatDate(date),
    },
    {
      title: '操作',
      key: 'actions',
      render: (_: any, record: ModelPurchase) => (
        <Space>
          <Button
            type="primary"
            size="small"
            onClick={() => handleDownload(record.model_id)}
          >
            <DownloadOutlined /> 下载
          </Button>
          <Link to={`/models/${record.model_id}/print`}>
            <Button size="small">
              <AppstoreOutlined /> 打印
            </Button>
          </Link>
        </Space>
      ),
    },
  ];

  const favoriteColumns = [
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
              {record.designer?.designer_profile?.nickname || record.designer?.username}
            </div>
          </div>
        </div>
      ),
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
      title: '操作',
      key: 'actions',
      render: (_: any, record: Model3D) => (
        <Space>
          <Link to={`/models/${record.id}`}>
            <Button type="primary" size="small">
              <EyeOutlined /> 查看
            </Button>
          </Link>
          <Button
            size="small"
            danger
            onClick={async () => {
              try {
                await modelApi.removeFavorite(record.id);
                message.success('已取消收藏');
                fetchData();
              } catch (error: any) {
                message.error(error.response?.data?.error || '操作失败');
              }
            }}
          >
            取消收藏
          </Button>
        </Space>
      ),
    },
  ];

  return (
    <div className="space-y-6">
      <Card
        title="模型管理"
        tabList={[
          { key: 'my', tab: '我的模型' },
          { key: 'purchased', tab: '已购模型' },
          { key: 'favorite', tab: '我的收藏' },
        ]}
        activeTabKey={activeTab}
        onTabChange={(key) => {
          setActiveTab(key as any);
          setPage(1);
        }}
        extra={
          activeTab === 'my' && (
            <Button type="primary" icon={<PlusOutlined />} onClick={() => navigate('/model-upload')}>
              上传模型
            </Button>
          )
        }
      >
        <Table
          rowKey="id"
          columns={
            activeTab === 'my'
              ? myModelColumns
              : activeTab === 'purchased'
              ? purchasedColumns
              : favoriteColumns
          }
          dataSource={
            activeTab === 'my'
              ? models
              : activeTab === 'purchased'
              ? purchases
              : favorites
          }
          loading={loading}
          pagination={{
            current: page,
            pageSize,
            total,
            onChange: setPage,
          }}
          locale={{ emptyText: <Empty description="暂无数据" /> }}
        />
      </Card>
    </div>
  );
};

export default MyModelsPage;
