import React, { useEffect, useState } from 'react';
import { Row, Col, Card, Statistic, Typography, List, Tag, Button, Rate, Avatar, Space } from 'antd';
import {
  AppstoreOutlined,
  ShoppingCartOutlined,
  PrinterOutlined,
  RiseOutlined,
  StarOutlined,
  ArrowRightOutlined,
} from '@ant-design/icons';
import { Link } from 'react-router-dom';
import { modelApi, printerApi, statsApi } from '@/services/api';
import { Model3D, User } from '@/types';
import { formatPrice, getRoleText } from '@/utils/format';
import { generateModelThumbnail } from '@/utils/format';
import { useAuth } from '@/hooks/useAuth';

const { Title, Text, Paragraph } = Typography;

const HomePage: React.FC = () => {
  const [hotModels, setHotModels] = useState<Model3D[]>([]);
  const [topDesigners, setTopDesigners] = useState<User[]>([]);
  const [topPrinters, setTopPrinters] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const { userRole } = useAuth();

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [modelsRes, designersRes, printersRes] = await Promise.all([
          modelApi.getHot(8),
          printerApi.listDesigners({ page_size: 6 }),
          printerApi.listPrinters({ page_size: 6 }),
        ]);

        setHotModels(modelsRes.data || []);
        setTopDesigners(designersRes.data?.data || []);
        setTopPrinters(printersRes.data?.data || []);
      } catch (error) {
        console.error('Failed to fetch home data:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  return (
    <div className="space-y-8">
      {/* Hero Section */}
      <Card
        className="bg-gradient-to-r from-blue-600 to-indigo-700 text-white border-0"
        bodyStyle={{ padding: '60px 40px' }}
      >
        <Row align="middle" gutter={32}>
          <Col xs={24} md={16}>
            <Title level={1} className="!text-white !mb-4">
              一站式3D打印服务与模型交易平台
            </Title>
            <Paragraph className="text-blue-100 text-lg mb-8">
              连接3D建模师、打印服务商和需要定制打印的客户，
              为您提供高质量的3D模型下载和打印服务
            </Paragraph>
            <Space size="large">
              <Link to="/models">
                <Button size="large" type="primary">
                  浏览模型市场
                </Button>
              </Link>
              {!userRole && (
                <Link to="/register">
                  <Button size="large" ghost className="text-white border-white hover:!text-white hover:!bg-white/10">
                    免费注册
                  </Button>
                </Link>
              )}
            </Space>
          </Col>
          <Col xs={0} md={8}>
            <div className="text-center">
              <div className="w-48 h-48 mx-auto bg-white/10 rounded-3xl flex items-center justify-center">
                <AppstoreOutlined className="text-8xl text-white/80" />
              </div>
            </div>
          </Col>
        </Row>
      </Card>

      {/* Stats */}
      <Row gutter={16}>
        <Col xs={12} md={6}>
          <Card>
            <Statistic
              title="3D模型总数"
              value={12580}
              prefix={<AppstoreOutlined />}
              valueStyle={{ color: '#3b82f6' }}
            />
          </Card>
        </Col>
        <Col xs={12} md={6}>
          <Card>
            <Statistic
              title="注册建模师"
              value={856}
              prefix={<StarOutlined />}
              valueStyle={{ color: '#8b5cf6' }}
            />
          </Card>
        </Col>
        <Col xs={12} md={6}>
          <Card>
            <Statistic
              title="合作打印商"
              value={234}
              prefix={<PrinterOutlined />}
              valueStyle={{ color: '#10b981' }}
            />
          </Card>
        </Col>
        <Col xs={12} md={6}>
          <Card>
            <Statistic
              title="完成订单"
              value={45680}
              prefix={<ShoppingCartOutlined />}
              valueStyle={{ color: '#f59e0b' }}
            />
          </Card>
        </Col>
      </Row>

      {/* Hot Models */}
      <Card
        title={
          <Space>
            <RiseOutlined className="text-red-500" />
            <span>热门模型</span>
          </Space>
        }
        extra={
          <Link to="/models">
            <Button type="link">
              查看全部 <ArrowRightOutlined />
            </Button>
          </Link>
        }
        loading={loading}
      >
        <Row gutter={[16, 16]}>
          {hotModels.map((model) => (
            <Col xs={24} sm={12} md={6} key={model.id}>
              <Card
                hoverable
                className="h-full"
                cover={
                  <div className="h-48 bg-gray-100 flex items-center justify-center overflow-hidden">
                    {model.thumbnail_url ? (
                      <img
                        src={model.thumbnail_url}
                        alt={model.title}
                        className="w-full h-full object-cover"
                      />
                    ) : (
                      <img
                        src={generateModelThumbnail(model.id)}
                        alt={model.title}
                        className="w-full h-full object-cover"
                      />
                    )}
                  </div>
                }
              >
                <Card.Meta
                  title={
                    <Link to={`/models/${model.id}`} className="hover:text-blue-600">
                      <div className="truncate">{model.title}</div>
                    </Link>
                  }
                  description={
                    <div className="space-y-2">
                      <div className="flex justify-between items-center">
                        <Text strong className="text-red-500">
                          {formatPrice(model.price)}
                        </Text>
                        <Rate disabled value={model.rating} size="small" />
                      </div>
                      <div className="flex gap-1 flex-wrap">
                        {model.tags?.slice(0, 2).map((tag, idx) => (
                          <Tag key={idx} color="blue" size="small">
                            {tag}
                          </Tag>
                        ))}
                      </div>
                      <div className="flex justify-between text-xs text-gray-400">
                        <span>下载 {model.download_count}</span>
                        <span>浏览 {model.view_count}</span>
                      </div>
                    </div>
                  }
                />
              </Card>
            </Col>
          ))}
        </Row>
      </Card>

      <Row gutter={16}>
        {/* Top Designers */}
        <Col xs={24} lg={12}>
          <Card
            title="优秀建模师"
            extra={
              <Link to="/designers">
                <Button type="link">
                  更多 <ArrowRightOutlined />
                </Button>
              </Link>
            }
            loading={loading}
          >
            <List
              dataSource={topDesigners}
              renderItem={(designer) => (
                <List.Item key={designer.id} className="border-b-0">
                  <List.Item.Meta
                    avatar={
                      <Avatar size="large" src={designer.avatar}>
                        {designer.username?.charAt(0).toUpperCase()}
                      </Avatar>
                    }
                    title={
                      <Space>
                        <span className="font-medium">{designer.designer_profile?.nickname || designer.username}</span>
                        <Tag color="purple">{getRoleText(designer.role)}</Tag>
                      </Space>
                    }
                    description={
                      <Space size="large">
                        <Text type="secondary">
                          <StarOutlined className="text-yellow-500 mr-1" />
                          {designer.designer_profile?.rating?.toFixed(1) || '5.0'}
                        </Text>
                        <Text type="secondary">
                          模型 {designer.designer_profile?.total_models || 0}
                        </Text>
                        <Text type="secondary">
                          销售额 {formatPrice(designer.designer_profile?.total_sales || 0)}
                        </Text>
                      </Space>
                    }
                  />
                </List.Item>
              )}
            />
          </Card>
        </Col>

        {/* Top Printers */}
        <Col xs={24} lg={12}>
          <Card
            title="优质打印商"
            extra={
              <Link to="/printers">
                <Button type="link">
                  更多 <ArrowRightOutlined />
                </Button>
              </Link>
            }
            loading={loading}
          >
            <List
              dataSource={topPrinters}
              renderItem={(printer) => (
                <List.Item key={printer.id} className="border-b-0">
                  <List.Item.Meta
                    avatar={
                      <Avatar size="large" src={printer.avatar}>
                        {printer.username?.charAt(0).toUpperCase()}
                      </Avatar>
                    }
                    title={
                      <Space>
                        <span className="font-medium">{printer.printer_profile?.company_name || printer.username}</span>
                        <Tag color="green">{getRoleText(printer.role)}</Tag>
                      </Space>
                    }
                    description={
                      <Space size="large">
                        <Text type="secondary">
                          <StarOutlined className="text-yellow-500 mr-1" />
                          {printer.printer_profile?.rating?.toFixed(1) || '5.0'}
                        </Text>
                        <Text type="secondary">
                          订单 {printer.printer_profile?.total_orders || 0}
                        </Text>
                        <Text type="secondary">
                          营收 {formatPrice(printer.printer_profile?.total_revenue || 0)}
                        </Text>
                      </Space>
                    }
                  />
                </List.Item>
              )}
            />
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default HomePage;
