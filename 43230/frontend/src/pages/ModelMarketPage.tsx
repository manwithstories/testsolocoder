import React, { useEffect, useState } from 'react';
import {
  Row,
  Col,
  Card,
  Input,
  Select,
  Slider,
  Button,
  Tag,
  Rate,
  Pagination,
  Space,
  Typography,
  Empty,
  Image,
} from 'antd';
import {
  SearchOutlined,
  FilterOutlined,
  ShoppingCartOutlined,
  HeartOutlined,
  StarOutlined,
} from '@ant-design/icons';
import { Link } from 'react-router-dom';
import { modelApi } from '@/services/api';
import { Model3D, LicenseType } from '@/types';
import { formatPrice, formatFileSize, generateModelThumbnail } from '@/utils/format';

const { Title, Text } = Typography;
const { Search } = Input;
const { Option } = Select;

const categories = [
  '全部',
  '建筑',
  '人物',
  '交通工具',
  '家居',
  '电子产品',
  '玩具',
  '艺术',
  '其他',
];

const ModelMarketPage: React.FC = () => {
  const [models, setModels] = useState<Model3D[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [page, setPage] = useState(1);
  const [pageSize] = useState(20);
  const [filters, setFilters] = useState({
    keyword: '',
    category: '',
    minPrice: 0,
    maxPrice: 1000,
    sortBy: 'created_at_desc',
    licenseType: '',
  });

  const fetchModels = async () => {
    setLoading(true);
    try {
      const params: any = {
        page,
        page_size: pageSize,
        sort_by: filters.sortBy,
      };

      if (filters.keyword) params.keyword = filters.keyword;
      if (filters.category && filters.category !== '全部') params.category = filters.category;
      if (filters.minPrice > 0) params.min_price = filters.minPrice;
      if (filters.maxPrice < 1000) params.max_price = filters.maxPrice;
      if (filters.licenseType) params.license_type = filters.licenseType;

      const response = await modelApi.list(params);
      setModels(response.data.data || []);
      setTotal(response.data.total || 0);
    } catch (error) {
      console.error('Failed to fetch models:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchModels();
  }, [page, filters]);

  const handleSearch = (value: string) => {
    setFilters((prev) => ({ ...prev, keyword: value }));
    setPage(1);
  };

  const handleCategoryChange = (value: string) => {
    setFilters((prev) => ({ ...prev, category: value }));
    setPage(1);
  };

  const handlePriceChange = (value: [number, number]) => {
    setFilters((prev) => ({ ...prev, minPrice: value[0], maxPrice: value[1] }));
    setPage(1);
  };

  const handleSortChange = (value: string) => {
    setFilters((prev) => ({ ...prev, sortBy: value }));
    setPage(1);
  };

  const handleLicenseChange = (value: string) => {
    setFilters((prev) => ({ ...prev, licenseType: value }));
    setPage(1);
  };

  return (
    <div className="space-y-6">
      <Card>
        <Space direction="vertical" size="large" className="w-full">
          <div className="flex justify-between items-center">
            <Title level={3} className="!mb-0">
              模型市场
            </Title>
            <Search
              placeholder="搜索模型名称或描述"
              allowClear
              enterButton={<SearchOutlined />}
              size="large"
              onSearch={handleSearch}
              style={{ width: 400 }}
            />
          </div>

          <div className="flex flex-wrap gap-4 items-center">
            <Space>
              <Text type="secondary">分类:</Text>
              <Select
                defaultValue="全部"
                style={{ width: 150 }}
                onChange={handleCategoryChange}
              >
                {categories.map((cat) => (
                  <Option key={cat} value={cat}>
                    {cat}
                  </Option>
                ))}
              </Select>
            </Space>

            <Space>
              <Text type="secondary">授权方式:</Text>
              <Select
                placeholder="全部"
                style={{ width: 150 }}
                onChange={handleLicenseChange}
                allowClear
              >
                <Option value="per_purchase">按件购买</Option>
                <Option value="subscription">订阅下载</Option>
              </Select>
            </Space>

            <Space>
              <Text type="secondary">排序:</Text>
              <Select defaultValue="created_at_desc" style={{ width: 150 }} onChange={handleSortChange}>
                <Option value="created_at_desc">最新上传</Option>
                <Option value="downloads_desc">下载最多</Option>
                <Option value="rating_desc">评分最高</Option>
                <Option value="price_asc">价格从低到高</Option>
                <Option value="price_desc">价格从高到低</Option>
              </Select>
            </Space>
          </div>

          <div className="flex items-center gap-4">
            <FilterOutlined className="text-gray-400" />
            <Text type="secondary">价格区间:</Text>
            <Slider
              range
              min={0}
              max={1000}
              step={10}
              value={[filters.minPrice, filters.maxPrice]}
              onChange={handlePriceChange}
              style={{ width: 300 }}
            />
            <Text>
              ¥{filters.minPrice} - ¥{filters.maxPrice}
            </Text>
          </div>
        </Space>
      </Card>

      {loading ? (
        <Card>
          <div className="text-center py-12">
            <div className="animate-spin w-8 h-8 border-4 border-blue-500 border-t-transparent rounded-full mx-auto mb-4"></div>
            <Text type="secondary">加载中...</Text>
          </div>
        </Card>
      ) : models.length === 0 ? (
        <Card>
          <Empty description="暂无模型" />
        </Card>
      ) : (
        <>
          <Row gutter={[16, 16]}>
            {models.map((model) => (
              <Col xs={24} sm={12} md={8} lg={6} key={model.id}>
                <Card
                  hoverable
                  className="h-full transition-transform hover:-translate-y-1"
                  cover={
                    <div className="h-48 bg-gray-100 relative overflow-hidden group">
                      <Image
                        src={model.thumbnail_url || generateModelThumbnail(model.id)}
                        alt={model.title}
                        className="w-full h-full object-cover transition-transform group-hover:scale-105"
                        preview={false}
                      />
                      <div className="absolute top-2 right-2">
                        {model.license_type === 'subscription' ? (
                          <Tag color="purple">订阅</Tag>
                        ) : (
                          <Tag color="blue">按件</Tag>
                        )}
                      </div>
                    </div>
                  }
                  actions={[
                    <Link to={`/models/${model.id}`} key="view">
                      <Button type="text" icon={<ShoppingCartOutlined />}>
                        查看详情
                      </Button>
                    </Link>,
                    <Button type="text" icon={<HeartOutlined />} key="favorite">
                      收藏
                    </Button>,
                  ]}
                >
                  <Card.Meta
                    title={
                      <Link to={`/models/${model.id}`} className="hover:text-blue-600">
                        <div className="font-medium truncate">{model.title}</div>
                      </Link>
                    }
                    description={
                      <div className="space-y-2 mt-2">
                        <div className="flex justify-between items-center">
                          <Text strong className="text-red-500 text-lg">
                            {formatPrice(model.price)}
                          </Text>
                          <Rate disabled value={model.rating} size="small" />
                        </div>
                        <div className="flex gap-1 flex-wrap">
                          {model.tags?.slice(0, 3).map((tag, idx) => (
                            <Tag key={idx} color="geekblue" size="small">
                              {tag}
                            </Tag>
                          ))}
                        </div>
                        <div className="flex justify-between text-xs text-gray-400">
                          <span>
                            <StarOutlined className="mr-1" />
                            {model.rating.toFixed(1)} ({model.rating_count})
                          </span>
                          <span>下载 {model.download_count}</span>
                        </div>
                        {model.file_size > 0 && (
                          <div className="text-xs text-gray-400">
                            文件大小: {formatFileSize(model.file_size)}
                          </div>
                        )}
                      </div>
                    }
                  />
                </Card>
              </Col>
            ))}
          </Row>

          <div className="flex justify-center mt-8">
            <Pagination
              current={page}
              pageSize={pageSize}
              total={total}
              showSizeChanger={false}
              onChange={setPage}
            />
          </div>
        </>
      )}
    </div>
  );
};

export default ModelMarketPage;
