import React, { useEffect, useState } from 'react';
import {
  Row,
  Col,
  Card,
  Descriptions,
  Tag,
  Rate,
  Button,
  Space,
  Typography,
  Modal,
  Form,
  Select,
  InputNumber,
  Radio,
  Divider,
  Statistic,
  Empty,
  List,
  Avatar,
  message,
  Tabs,
  Image,
} from 'antd';
import {
  ShoppingCartOutlined,
  HeartOutlined,
  DownloadOutlined,
  PrinterOutlined,
  StarOutlined,
  UserOutlined,
  FileTextOutlined,
  ClockCircleOutlined,
  BoxPlotOutlined,
  EditOutlined,
} from '@ant-design/icons';
import { useNavigate, useParams } from 'react-router-dom';
import { modelApi, orderApi, printerApi } from '@/services/api';
import {
  Model3D,
  Material,
  PrintQuality,
  PriceEstimateRequest,
  PriceEstimateResponse,
  Review,
} from '@/types';
import {
  formatPrice,
  formatDuration,
  getMaterialTypeText,
  getQualityText,
  formatFileSize,
  generateModelThumbnail,
} from '@/utils/format';
import { formatDate, getRelativeTime } from '@/utils/date';
import { useAuth } from '@/hooks/useAuth';
import { useChunkedUpload } from '@/hooks/useChunkedUpload';

const { Title, Text, Paragraph } = Typography;
const { Option } = Select;
const { TabPane } = Tabs;

const ModelDetailPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [model, setModel] = useState<Model3D | null>(null);
  const [materials, setMaterials] = useState<Material[]>([]);
  const [reviews, setReviews] = useState<Review[]>([]);
  const [loading, setLoading] = useState(true);
  const [orderModalVisible, setOrderModalVisible] = useState(false);
  const [purchaseModalVisible, setPurchaseModalVisible] = useState(false);
  const [estimate, setEstimate] = useState<PriceEstimateResponse | null>(null);
  const [estimating, setEstimating] = useState(false);
  const [isPurchased, setIsPurchased] = useState(false);
  const [isFavorited, setIsFavorited] = useState(false);
  const { userRole, isAuthenticated } = useAuth();
  const { uploadFile } = useChunkedUpload();

  const [orderForm] = Form.useForm();
  const [purchaseForm] = Form.useForm();

  useEffect(() => {
    const fetchData = async () => {
      if (!id) return;

      try {
        const [modelRes, materialsRes, reviewsRes] = await Promise.all([
          modelApi.get(id),
          printerApi.getMaterials(),
          printerApi.getModelReviews(id, { page_size: 10 }),
        ]);

        setModel(modelRes.data);
        setMaterials(materialsRes.data);
        setReviews(reviewsRes.data?.data || []);
      } catch (error) {
        console.error('Failed to fetch model details:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [id]);

  const handleEstimate = async (values: any) => {
    if (!id) return;

    setEstimating(true);
    try {
      const req: PriceEstimateRequest = {
        model_id: id,
        quantity: values.quantity,
        material_id: values.material_id,
        quality: values.quality,
        layer_height: values.layer_height || 0.2,
        infill_percent: values.infill_percent || 20,
        supports: values.supports || false,
      };

      const response = await orderApi.estimate(req);
      setEstimate(response.data);
    } catch (error: any) {
      message.error(error.response?.data?.error || '估价失败');
    } finally {
      setEstimating(false);
    }
  };

  const handleCreateOrder = async () => {
    try {
      const values = await orderForm.validateFields();
      const orderData = {
        ...values,
        model_id: id,
        shipping_address: values.shipping_address || '请填写收货地址',
      };

      const response = await orderApi.create(orderData);
      message.success('订单创建成功');
      setOrderModalVisible(false);
      navigate(`/my-orders/${response.data.id}`);
    } catch (error: any) {
      message.error(error.response?.data?.error || '创建订单失败');
    }
  };

  const handlePurchaseModel = async () => {
    if (!id) return;

    try {
      const values = await purchaseForm.validateFields();
      await modelApi.purchase(id, { purchase_type: values.purchase_type });
      message.success('购买成功');
      setPurchaseModalVisible(false);
      setIsPurchased(true);
    } catch (error: any) {
      message.error(error.response?.data?.error || '购买失败');
    }
  };

  const handleDownload = async () => {
    if (!id) return;

    try {
      const response = await modelApi.download(id);
      const fileURL = response.data.file_url;
      window.open(fileURL, '_blank');
      message.success('开始下载');
    } catch (error: any) {
      message.error(error.response?.data?.error || '下载失败');
    }
  };

  const handleFavorite = async () => {
    if (!id) return;

    try {
      if (isFavorited) {
        await modelApi.removeFavorite(id);
        setIsFavorited(false);
        message.success('已取消收藏');
      } else {
        await modelApi.addFavorite(id);
        setIsFavorited(true);
        message.success('收藏成功');
      }
    } catch (error: any) {
      message.error(error.response?.data?.error || '操作失败');
    }
  };

  const handleOrderValuesChange = (_: any, allValues: any) => {
    if (allValues.quantity && allValues.material_id && allValues.quality) {
      handleEstimate(allValues);
    }
  };

  if (loading) {
    return (
      <Card>
        <div className="text-center py-12">
          <div className="animate-spin w-8 h-8 border-4 border-blue-500 border-t-transparent rounded-full mx-auto"></div>
        </div>
      </Card>
    );
  }

  if (!model) {
    return (
      <Card>
        <Empty description="模型不存在或已被删除" />
      </Card>
    );
  }

  return (
    <div className="space-y-6">
      <Card>
        <Row gutter={32}>
          <Col xs={24} md={10}>
            <div className="sticky top-6">
              <div className="bg-gray-100 rounded-lg overflow-hidden aspect-square flex items-center justify-center">
                <Image
                  src={model.thumbnail_url || generateModelThumbnail(model.id)}
                  alt={model.title}
                  className="w-full h-full object-contain"
                  preview
                />
              </div>

              <Space className="w-full mt-4">
                {(isPurchased || model.designer_id || userRole === 'designer') && (
                  <Button
                    type="primary"
                    icon={<DownloadOutlined />}
                    onClick={handleDownload}
                    className="flex-1"
                    size="large"
                  >
                    下载模型
                  </Button>
                )}
                {!isPurchased && model.price > 0 && (
                  <Button
                    type="primary"
                    icon={<ShoppingCartOutlined />}
                    onClick={() => setPurchaseModalVisible(true)}
                    className="flex-1"
                    size="large"
                  >
                    购买模型 {formatPrice(model.price)}
                  </Button>
                )}
                <Button
                  icon={<PrinterOutlined />}
                  onClick={() => setOrderModalVisible(true)}
                  className="flex-1"
                  size="large"
                >
                  立即打印
                </Button>
                <Button
                  icon={<HeartOutlined />}
                  onClick={handleFavorite}
                  danger={isFavorited}
                  size="large"
                >
                  {isFavorited ? '已收藏' : '收藏'}
                </Button>
              </Space>
            </div>
          </Col>

          <Col xs={24} md={14}>
            <Space direction="vertical" size="large" className="w-full">
              <div>
                <Space className="mb-2">
                  {model.license_type === 'subscription' ? (
                    <Tag color="purple">订阅下载</Tag>
                  ) : (
                    <Tag color="blue">按件购买</Tag>
                  )}
                  {model.category && <Tag color="geekblue">{model.category}</Tag>}
                  {model.is_featured && <Tag color="gold">精选</Tag>}
                </Space>
                <Title level={2} className="!mb-2">
                  {model.title}
                </Title>
                <div className="flex items-center gap-4 mb-4">
                  <Space>
                    <Rate disabled value={model.rating} />
                    <Text type="secondary">
                      {model.rating.toFixed(1)} ({model.rating_count}条评价)
                    </Text>
                  </Space>
                  <Text type="secondary">
                    <StarOutlined className="mr-1" />
                    {model.view_count} 浏览
                  </Text>
                  <Text type="secondary">
                    <DownloadOutlined className="mr-1" />
                    {model.download_count} 下载
                  </Text>
                  <Text type="secondary">
                    <ShoppingCartOutlined className="mr-1" />
                    {model.purchase_count} 购买
                  </Text>
                </div>
                <Title level={4} type="success">
                  {formatPrice(model.price)}
                  {model.license_type === 'subscription' && (
                    <Text type="secondary" className="ml-2 text-sm font-normal">
                      / 月
                    </Text>
                  )}
                </Title>
              </div>

              <Divider />

              <Tabs defaultActiveKey="description">
                <TabPane tab="商品详情" key="description">
                  <Paragraph className="whitespace-pre-wrap">{model.description}</Paragraph>
                </TabPane>

                <TabPane tab="技术参数" key="specs">
                  <Row gutter={[16, 16]}>
                    <Col xs={12} md={8}>
                      <Card size="small">
                        <Statistic
                          title="文件格式"
                          value={model.file_type?.toUpperCase() || 'STL'}
                          prefix={<FileTextOutlined />}
                        />
                      </Card>
                    </Col>
                    <Col xs={12} md={8}>
                      <Card size="small">
                        <Statistic
                          title="文件大小"
                          value={formatFileSize(model.file_size || 0)}
                          prefix={<BoxPlotOutlined />}
                        />
                      </Card>
                    </Col>
                    <Col xs={12} md={8}>
                      <Card size="small">
                        <Statistic
                          title="模型体积"
                          value={model.volume?.toFixed(2) || 0}
                          suffix="cm³"
                          prefix={<BoxPlotOutlined />}
                        />
                      </Card>
                    </Col>
                    <Col xs={12} md={8}>
                      <Card size="small">
                        <Statistic
                          title="预估打印时间"
                          value={formatDuration(model.print_time_hours || 0)}
                          prefix={<ClockCircleOutlined />}
                        />
                      </Card>
                    </Col>
                    <Col xs={12} md={8}>
                      <Card size="small">
                        <Statistic
                          title="多边形数"
                          value={model.polygon_count || 0}
                          prefix={<EditOutlined />}
                        />
                      </Card>
                    </Col>
                    <Col xs={12} md={8}>
                      <Card size="small">
                        <Statistic
                          title="版本"
                          value={model.version || '1.0.0'}
                          prefix={<FileTextOutlined />}
                        />
                      </Card>
                    </Col>
                  </Row>

                  <Divider orientation="left">推荐打印材料</Divider>
                  <Space wrap>
                    {model.recommended_materials?.map((mat, idx) => (
                      <Tag key={idx} color="green" size="large">
                        {getMaterialTypeText(mat)}
                      </Tag>
                    ))}
                  </Space>

                  <Divider orientation="left">打印参数建议</Divider>
                  <Descriptions column={1} size="small">
                    <Descriptions.Item label="外边框尺寸">
                      {model.bounding_box || '未提供'}
                    </Descriptions.Item>
                    <Descriptions.Item label="层高建议">
                      0.1mm - 0.3mm
                    </Descriptions.Item>
                    <Descriptions.Item label="填充密度建议">
                      15% - 30%
                    </Descriptions.Item>
                    <Descriptions.Item label="是否需要支撑">
                      根据摆放角度确定
                    </Descriptions.Item>
                  </Descriptions>
                </TabPane>

                <TabPane tab={`用户评价 (${reviews.length})`} key="reviews">
                  {reviews.length === 0 ? (
                    <Empty description="暂无评价" />
                  ) : (
                    <List
                      dataSource={reviews}
                      renderItem={(review) => (
                        <List.Item key={review.id}>
                          <List.Item.Meta
                            avatar={
                              <Avatar src={review.customer?.avatar}>
                                {review.customer?.username?.charAt(0).toUpperCase()}
                              </Avatar>
                            }
                            title={
                              <Space>
                                <span className="font-medium">
                                  {review.is_anonymous ? '匿名用户' : review.customer?.username}
                                </span>
                                <Text type="secondary" className="text-sm">
                                  {getRelativeTime(review.created_at)}
                                </Text>
                              </Space>
                            }
                            description={
                              <div className="space-y-2">
                                <div className="flex gap-4">
                                  <Space>
                                    <Text type="secondary">模型质量:</Text>
                                    <Rate disabled value={review.model_rating} size="small" />
                                  </Space>
                                  <Space>
                                    <Text type="secondary">打印效果:</Text>
                                    <Rate disabled value={review.print_rating} size="small" />
                                  </Space>
                                </div>
                                {review.model_comment && (
                                  <Paragraph className="!mb-0">
                                    <Text type="secondary">模型评价: </Text>
                                    {review.model_comment}
                                  </Paragraph>
                                )}
                                {review.print_comment && (
                                  <Paragraph className="!mb-0">
                                    <Text type="secondary">打印评价: </Text>
                                    {review.print_comment}
                                  </Paragraph>
                                )}
                              </div>
                            }
                          />
                        </List.Item>
                      )}
                    />
                  )}
                </TabPane>
              </Tabs>
            </Space>
          </Col>
        </Row>
      </Card>

      {/* 设计师信息 */}
      <Card title="设计师信息">
        <Space align="start" size="large">
          <Avatar size={64} src={model.designer?.avatar}>
            {model.designer?.username?.charAt(0).toUpperCase()}
          </Avatar>
          <div className="flex-1">
            <div className="flex items-center gap-3 mb-2">
              <Title level={4} className="!mb-0">
                {model.designer?.designer_profile?.nickname || model.designer?.username}
              </Title>
              <Tag color="purple">建模师</Tag>
            </div>
            <Paragraph className="mb-4">
              {model.designer?.designer_profile?.bio || '该设计师暂无简介'}
            </Paragraph>
            <Space size="large">
              <Space>
                <StarOutlined className="text-yellow-500" />
                <Text>
                  {model.designer?.designer_profile?.rating?.toFixed(1) || '5.0'}
                  <Text type="secondary" className="ml-1">
                    ({model.designer?.designer_profile?.rating_count || 0}条评价)
                  </Text>
                </Text>
              </Space>
              <Text>
                模型总数: {model.designer?.designer_profile?.total_models || 0}
              </Text>
              <Text>
                总销售额: {formatPrice(model.designer?.designer_profile?.total_sales || 0)}
              </Text>
            </Space>
          </div>
          <Button type="primary">进入设计师主页</Button>
        </Space>
      </Card>

      {/* 打印下单弹窗 */}
      <Modal
        title="3D打印下单"
        open={orderModalVisible}
        width={700}
        onCancel={() => setOrderModalVisible(false)}
        footer={[
          <Button key="cancel" onClick={() => setOrderModalVisible(false)}>
            取消
          </Button>,
          <Button
            key="submit"
            type="primary"
            onClick={handleCreateOrder}
            disabled={!estimate}
          >
            确认下单 {estimate && `(${formatPrice(estimate.total_amount)})`}
          </Button>,
        ]}
      >
        <Form
          form={orderForm}
          layout="vertical"
          onValuesChange={handleOrderValuesChange}
          initialValues={{
            quantity: 1,
            quality: 'standard',
            layer_height: 0.2,
            infill_percent: 20,
            supports: false,
          }}
        >
          <Row gutter={16}>
            <Col xs={24} md={12}>
              <Form.Item
                name="quantity"
                label="打印数量"
                rules={[{ required: true, message: '请输入打印数量' }]}
              >
                <InputNumber min={1} max={100} className="w-full" />
              </Form.Item>
            </Col>
            <Col xs={24} md={12}>
              <Form.Item
                name="material_id"
                label="打印材料"
                rules={[{ required: true, message: '请选择打印材料' }]}
              >
                <Select placeholder="请选择材料">
                  {materials.map((mat) => (
                    <Option key={mat.id} value={mat.id}>
                      {mat.name} - {formatPrice(mat.price_per_gram)}/g
                    </Option>
                  ))}
                </Select>
              </Form.Item>
            </Col>
            <Col xs={24} md={12}>
              <Form.Item
                name="quality"
                label="打印精度"
                rules={[{ required: true, message: '请选择打印精度' }]}
              >
                <Select>
                  <Option value="draft">{getQualityText('draft')}</Option>
                  <Option value="standard">{getQualityText('standard')}</Option>
                  <Option value="high">{getQualityText('high')}</Option>
                  <Option value="ultra">{getQualityText('ultra')}</Option>
                </Select>
              </Form.Item>
            </Col>
            <Col xs={24} md={12}>
              <Form.Item
                name="color"
                label="颜色"
                rules={[{ required: true, message: '请选择颜色' }]}
              >
                <Select placeholder="请选择颜色">
                  <Option value="white">白色</Option>
                  <Option value="black">黑色</Option>
                  <Option value="gray">灰色</Option>
                  <Option value="red">红色</Option>
                  <Option value="blue">蓝色</Option>
                  <Option value="green">绿色</Option>
                  <Option value="yellow">黄色</Option>
                  <Option value="transparent">透明</Option>
                </Select>
              </Form.Item>
            </Col>
            <Col xs={24} md={12}>
              <Form.Item name="layer_height" label="层高 (mm)">
                <Select>
                  <Option value={0.3}>0.3mm (粗糙)</Option>
                  <Option value={0.2}>0.2mm (标准)</Option>
                  <Option value={0.15}>0.15mm (精细)</Option>
                  <Option value={0.1}>0.1mm (超高精度)</Option>
                </Select>
              </Form.Item>
            </Col>
            <Col xs={24} md={12}>
              <Form.Item name="infill_percent" label="填充密度 (%)">
                <Select>
                  <Option value={10}>10% (轻质)</Option>
                  <Option value={20}>20% (标准)</Option>
                  <Option value={40}>40% (坚固)</Option>
                  <Option value={60}>60% (高强度)</Option>
                  <Option value={100}>100% (实心)</Option>
                </Select>
              </Form.Item>
            </Col>
            <Col xs={24}>
              <Form.Item name="supports" label="是否需要添加支撑">
                <Radio.Group>
                  <Radio value={false}>不需要</Radio>
                  <Radio value={true}>需要（增加材料费用）</Radio>
                </Radio.Group>
              </Form.Item>
            </Col>
            <Col xs={24}>
              <Form.Item
                name="shipping_address"
                label="收货地址"
                rules={[{ required: true, message: '请输入收货地址' }]}
              >
                <Input.TextArea rows={2} placeholder="请输入详细的收货地址" />
              </Form.Item>
            </Col>
            <Col xs={24}>
              <Form.Item name="notes" label="备注说明">
                <Input.TextArea rows={2} placeholder="如有特殊要求请备注" />
              </Form.Item>
            </Col>
          </Row>

          {estimate && (
            <Card size="small" className="bg-blue-50">
              <Title level={5} className="!mb-3">
                价格预估
              </Title>
              <Row gutter={16}>
                <Col xs={12} md={6}>
                  <Statistic
                    title="模型费用"
                    value={estimate.base_price}
                    prefix="¥"
                    precision={2}
                    size="small"
                  />
                </Col>
                <Col xs={12} md={6}>
                  <Statistic
                    title="材料费用"
                    value={estimate.material_cost}
                    prefix="¥"
                    precision={2}
                    size="small"
                  />
                </Col>
                <Col xs={12} md={6}>
                  <Statistic
                    title="服务费用"
                    value={estimate.service_fee}
                    prefix="¥"
                    precision={2}
                    size="small"
                  />
                </Col>
                <Col xs={12} md={6}>
                  <Statistic
                    title="总费用"
                    value={estimate.total_amount}
                    prefix="¥"
                    precision={2}
                    valueStyle={{ color: '#cf1322' }}
                    size="small"
                  />
                </Col>
              </Row>
              <Divider className="my-3" />
              <div className="flex justify-between text-sm">
                <Text type="secondary">
                  总重量: {estimate.estimated_weight.toFixed(2)}g
                </Text>
                <Text type="secondary">
                  预估时间: {formatDuration(estimate.estimated_print_time)}
                </Text>
              </div>
            </Card>
          )}

          {estimating && (
            <div className="text-center py-4">
              <div className="animate-spin w-6 h-6 border-4 border-blue-500 border-t-transparent rounded-full mx-auto"></div>
              <Text type="secondary" className="mt-2">计算中...</Text>
            </div>
          )}
        </Form>
      </Modal>

      {/* 购买模型弹窗 */}
      <Modal
        title="购买模型"
        open={purchaseModalVisible}
        onCancel={() => setPurchaseModalVisible(false)}
        onOk={handlePurchaseModel}
        okText="确认购买"
      >
        <Form form={purchaseForm} layout="vertical">
          <div className="mb-4 p-4 bg-gray-50 rounded-lg">
            <div className="flex justify-between items-center mb-2">
              <Text strong>{model.title}</Text>
              <Tag color="blue">{model.license_type === 'subscription' ? '订阅' : '按件'}</Tag>
            </div>
            <div className="flex justify-between">
              <Text type="secondary">价格</Text>
              <Text strong className="text-red-500 text-lg">
                {formatPrice(model.price)}
                {model.license_type === 'subscription' && '/月'}
              </Text>
            </div>
          </div>

          <Form.Item
            name="purchase_type"
            label="授权方式"
            rules={[{ required: true, message: '请选择授权方式' }]}
            initialValue={model.license_type}
          >
            <Radio.Group>
              <Radio value="per_purchase">
                按件购买 - 单次授权，永久使用
              </Radio>
              {model.subscription_price > 0 && (
                <Radio value="subscription">
                  订阅下载 - {formatPrice(model.subscription_price)}/月，全库畅下
                </Radio>
              )}
            </Radio.Group>
          </Form.Item>

          <div className="text-sm text-gray-500">
            <p>• 购买后可下载STL/OBJ源文件</p>
            <p>• 支持个人和商业用途（根据授权类型）</p>
            <p>• 购买后请在30天内下载</p>
          </div>
        </Form>
      </Modal>
    </div>
  );
};

export default ModelDetailPage;
