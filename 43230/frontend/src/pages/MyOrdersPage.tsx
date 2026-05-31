import React, { useEffect, useState } from 'react';
import {
  Table,
  Tag,
  Button,
  Space,
  Card,
  Typography,
  Select,
  Modal,
  Descriptions,
  Timeline,
  Empty,
  Input,
  Form,
  Rate,
  message,
  Steps,
  Row,
  Col,
  Checkbox,
} from 'antd';
import {
  EyeOutlined,
  PrinterOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  StarOutlined,
  PackageOutlined,
  ClockCircleOutlined,
  SendOutlined,
} from '@ant-design/icons';
import { Link } from 'react-router-dom';
import { orderApi, printerApi } from '@/services/api';
import { PrintOrder, OrderStatus } from '@/types';
import {
  formatPrice,
  getOrderStatusText,
  getOrderStatusColor,
  formatFileSize,
  formatDuration,
  getQualityText,
  getMaterialTypeText,
} from '@/utils/format';
import { formatDate, getRelativeTime } from '@/utils/date';

const { Title, Text, Paragraph } = Typography;
const { Option } = Select;
const { Step } = Steps;

const MyOrdersPage: React.FC = () => {
  const [orders, setOrders] = useState<PrintOrder[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [page, setPage] = useState(1);
  const [pageSize] = useState(10);
  const [status, setStatus] = useState<string>('');
  const [detailVisible, setDetailVisible] = useState(false);
  const [selectedOrder, setSelectedOrder] = useState<PrintOrder | null>(null);
  const [reviewVisible, setReviewVisible] = useState(false);
  const [history, setHistory] = useState<any[]>([]);
  const [reviewForm] = Form.useForm();

  const fetchOrders = async () => {
    setLoading(true);
    try {
      const params: any = {
        page,
        page_size: pageSize,
      };
      if (status) {
        params.status = status;
      }

      const response = await orderApi.listCustomerOrders(params);
      setOrders(response.data.data || []);
      setTotal(response.data.total || 0);
    } catch (error) {
      console.error('Failed to fetch orders:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchOrders();
  }, [page, status]);

  const handleViewDetail = async (order: PrintOrder) => {
    setSelectedOrder(order);
    try {
      const historyRes = await orderApi.getHistory(order.id);
      setHistory(historyRes.data || []);
    } catch (error) {
      console.error('Failed to fetch order history:', error);
    }
    setDetailVisible(true);
  };

  const handleConfirmDelivery = async (orderId: string) => {
    try {
      await orderApi.deliverOrder(orderId);
      message.success('确认收货成功');
      fetchOrders();
    } catch (error: any) {
      message.error(error.response?.data?.error || '操作失败');
    }
  };

  const handleCompleteOrder = async (orderId: string) => {
    try {
      await orderApi.completeOrder(orderId);
      message.success('订单已完成，系统已自动结算');
      fetchOrders();
    } catch (error: any) {
      message.error(error.response?.data?.error || '操作失败');
    }
  };

  const handleCancelOrder = async (orderId: string) => {
    Modal.confirm({
      title: '确认取消订单',
      content: '取消订单后款项将原路退回，确定要取消吗？',
      okText: '确定取消',
      cancelText: '再想想',
      okType: 'danger',
      onOk: async () => {
        try {
          await orderApi.cancelOrder(orderId, { reason: '客户取消' });
          message.success('订单已取消');
          fetchOrders();
        } catch (error: any) {
          message.error(error.response?.data?.error || '取消失败');
        }
      },
    });
  };

  const handleSubmitReview = async (values: any) => {
    try {
      await printerApi.createReview({
        ...values,
        order_id: selectedOrder?.id,
      });
      message.success('评价提交成功');
      setReviewVisible(false);
      reviewForm.resetFields();
      fetchOrders();
    } catch (error: any) {
      message.error(error.response?.data?.error || '评价失败');
    }
  };

  const getOrderStep = (status: OrderStatus): number => {
    const steps: OrderStatus[] = [
      'pending',
      'paid',
      'printing',
      'quality_check',
      'shipped',
      'delivered',
      'completed',
    ];
    const idx = steps.indexOf(status);
    return idx >= 0 ? idx : 0;
  };

  const columns = [
    {
      title: '订单信息',
      dataIndex: 'order_no',
      key: 'order_no',
      render: (_: string, record: PrintOrder) => (
        <div>
          <div className="font-medium">{record.order_no}</div>
          <div className="text-sm text-gray-500">
            下单时间: {formatDate(record.created_at)}
          </div>
        </div>
      ),
    },
    {
      title: '模型信息',
      key: 'model',
      render: (_: any, record: PrintOrder) => (
        <div className="flex items-center gap-3">
          <div className="w-16 h-16 bg-gray-100 rounded flex items-center justify-center">
            {record.model?.thumbnail_url ? (
              <img
                src={record.model.thumbnail_url}
                alt={record.model?.title}
                className="w-full h-full object-cover rounded"
              />
            ) : (
              <PackageOutlined className="text-2xl text-gray-400" />
            )}
          </div>
          <div>
            <Link to={`/models/${record.model_id}`} className="hover:text-blue-600">
              {record.model?.title}
            </Link>
            <div className="text-sm text-gray-500">
              {record.quantity}件 × {formatPrice(record.model?.price || 0)}
            </div>
          </div>
        </div>
      ),
    },
    {
      title: '打印参数',
      key: 'params',
      render: (_: any, record: PrintOrder) => (
        <div className="text-sm space-y-1">
          <div>
            <Text type="secondary">材料: </Text>
            {getMaterialTypeText(record.material?.type || '')} / {record.color}
          </div>
          <div>
            <Text type="secondary">精度: </Text>
            {getQualityText(record.quality)}
          </div>
          <div>
            <Text type="secondary">填充: </Text>
            {record.infill_percent}%
          </div>
        </div>
      ),
    },
    {
      title: '金额',
      dataIndex: 'total_amount',
      key: 'total_amount',
      render: (value: number) => (
        <div className="text-red-500 font-bold">{formatPrice(value)}</div>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: OrderStatus) => (
        <Tag color={getOrderStatusColor(status)}>{getOrderStatusText(status)}</Tag>
      ),
    },
    {
      title: '操作',
      key: 'actions',
      render: (_: any, record: PrintOrder) => (
        <Space direction="vertical" size="small">
          <Button type="link" size="small" onClick={() => handleViewDetail(record)}>
            <EyeOutlined /> 详情
          </Button>
          {record.status === 'pending' && (
            <Button
              type="link"
              size="small"
              danger
              onClick={() => handleCancelOrder(record.id)}
            >
              <CloseCircleOutlined /> 取消
            </Button>
          )}
          {record.status === 'shipped' && (
            <Button
              type="link"
              size="small"
              type="primary"
              onClick={() => handleConfirmDelivery(record.id)}
            >
              <CheckCircleOutlined /> 确认收货
            </Button>
          )}
          {record.status === 'delivered' && (
            <Button
              type="link"
              size="small"
              type="primary"
              onClick={() => handleCompleteOrder(record.id)}
            >
              <CheckCircleOutlined /> 确认完成
            </Button>
          )}
          {record.status === 'completed' && (
            <Button
              type="link"
              size="small"
              onClick={() => {
                setSelectedOrder(record);
                setReviewVisible(true);
              }}
            >
              <StarOutlined /> 评价
            </Button>
          )}
        </Space>
      ),
    },
  ];

  return (
    <div className="space-y-6">
      <Card
        title="我的订单"
        extra={
          <Select
            placeholder="全部状态"
            style={{ width: 150 }}
            allowClear
            onChange={setStatus}
            value={status || undefined}
          >
            <Option value="pending">待付款</Option>
            <Option value="paid">已付款</Option>
            <Option value="printing">打印中</Option>
            <Option value="quality_check">质检中</Option>
            <Option value="shipped">已发货</Option>
            <Option value="delivered">已送达</Option>
            <Option value="completed">已完成</Option>
            <Option value="cancelled">已取消</Option>
          </Select>
        }
      >
        <Table
          rowKey="id"
          columns={columns}
          dataSource={orders}
          loading={loading}
          pagination={{
            current: page,
            pageSize,
            total,
            onChange: setPage,
          }}
          locale={{ emptyText: <Empty description="暂无订单" /> }}
        />
      </Card>

      {/* 订单详情弹窗 */}
      <Modal
        title={`订单详情 - ${selectedOrder?.order_no}`}
        open={detailVisible}
        width={800}
        onCancel={() => setDetailVisible(false)}
        footer={null}
      >
        {selectedOrder && (
          <div className="space-y-6">
            <Card size="small">
              <Steps
                current={getOrderStep(selectedOrder.status)}
                labelPlacement="vertical"
                size="small"
              >
                <Step title="创建订单" icon={<ClockCircleOutlined />} />
                <Step title="已付款" icon={<CheckCircleOutlined />} />
                <Step title="打印中" icon={<PrinterOutlined />} />
                <Step title="质检中" icon={<CheckCircleOutlined />} />
                <Step title="已发货" icon={<SendOutlined />} />
                <Step title="已送达" icon={<PackageOutlined />} />
                <Step title="已完成" icon={<CheckCircleOutlined />} />
              </Steps>
            </Card>

            <Descriptions column={2} size="small" bordered>
              <Descriptions.Item label="订单编号">
                {selectedOrder.order_no}
              </Descriptions.Item>
              <Descriptions.Item label="订单状态">
                <Tag color={getOrderStatusColor(selectedOrder.status)}>
                  {getOrderStatusText(selectedOrder.status)}
                </Tag>
              </Descriptions.Item>
              <Descriptions.Item label="下单时间">
                {formatDate(selectedOrder.created_at)}
              </Descriptions.Item>
              <Descriptions.Item label="打印商">
                {selectedOrder.printer?.printer_profile?.company_name || '待分配'}
              </Descriptions.Item>
              <Descriptions.Item label="收货地址" span={2}>
                {selectedOrder.shipping_address}
              </Descriptions.Item>
              {selectedOrder.tracking_number && (
                <Descriptions.Item label="物流单号" span={2}>
                  {selectedOrder.tracking_number}
                </Descriptions.Item>
              )}
              <Descriptions.Item label="备注" span={2}>
                {selectedOrder.notes || '无'}
              </Descriptions.Item>
            </Descriptions>

            <Card title="打印信息" size="small">
              <Row gutter={16}>
                <Col span={12}>
                  <Descriptions column={1} size="small">
                    <Descriptions.Item label="模型">
                      {selectedOrder.model?.title}
                    </Descriptions.Item>
                    <Descriptions.Item label="数量">
                      {selectedOrder.quantity}件
                    </Descriptions.Item>
                    <Descriptions.Item label="材料">
                      {getMaterialTypeText(selectedOrder.material?.type || '')}
                    </Descriptions.Item>
                    <Descriptions.Item label="颜色">
                      {selectedOrder.color}
                    </Descriptions.Item>
                    <Descriptions.Item label="打印精度">
                      {getQualityText(selectedOrder.quality)}
                    </Descriptions.Item>
                  </Descriptions>
                </Col>
                <Col span={12}>
                  <Descriptions column={1} size="small">
                    <Descriptions.Item label="层高">
                      {selectedOrder.layer_height}mm
                    </Descriptions.Item>
                    <Descriptions.Item label="填充密度">
                      {selectedOrder.infill_percent}%
                    </Descriptions.Item>
                    <Descriptions.Item label="添加支撑">
                      {selectedOrder.supports ? '是' : '否'}
                    </Descriptions.Item>
                    <Descriptions.Item label="预估重量">
                      {selectedOrder.estimated_weight?.toFixed(2)}g
                    </Descriptions.Item>
                    <Descriptions.Item label="预估时间">
                      {formatDuration(selectedOrder.estimated_print_time || 0)}
                    </Descriptions.Item>
                  </Descriptions>
                </Col>
              </Row>
            </Card>

            <Card title="费用明细" size="small">
              <Row gutter={16}>
                <Col span={6}>
                  <Text type="secondary">模型费用</Text>
                  <div className="text-lg font-medium">
                    {formatPrice(selectedOrder.base_price)}
                  </div>
                </Col>
                <Col span={6}>
                  <Text type="secondary">材料费用</Text>
                  <div className="text-lg font-medium">
                    {formatPrice(selectedOrder.material_cost)}
                  </div>
                </Col>
                <Col span={6}>
                  <Text type="secondary">服务费用</Text>
                  <div className="text-lg font-medium">
                    {formatPrice(selectedOrder.service_fee)}
                  </div>
                </Col>
                <Col span={6}>
                  <Text type="secondary">总金额</Text>
                  <div className="text-xl font-bold text-red-500">
                    {formatPrice(selectedOrder.total_amount)}
                  </div>
                </Col>
              </Row>
            </Card>

            <Card title="订单进度" size="small">
              <Timeline
                items={history.map((h) => ({
                  color:
                    h.status === 'cancelled'
                      ? 'red'
                      : h.status === 'completed'
                      ? 'green'
                      : 'blue',
                  children: (
                    <div>
                      <div className="font-medium">
                        {getOrderStatusText(h.status)}
                      </div>
                      <div className="text-sm text-gray-500">{h.description}</div>
                      <div className="text-xs text-gray-400">
                        {formatDate(h.created_at)}
                      </div>
                    </div>
                  ),
                }))}
              />
            </Card>
          </div>
        )}
      </Modal>

      {/* 评价弹窗 */}
      <Modal
        title="订单评价"
        open={reviewVisible}
        onCancel={() => setReviewVisible(false)}
        onOk={reviewForm.submit}
        okText="提交评价"
      >
        <Form form={reviewForm} layout="vertical" onFinish={handleSubmitReview}>
          <div className="mb-6 p-4 bg-gray-50 rounded-lg">
            <div className="font-medium mb-2">{selectedOrder?.model?.title}</div>
            <div className="text-sm text-gray-500">
              订单号: {selectedOrder?.order_no}
            </div>
          </div>

          <Form.Item
            name="model_rating"
            label="模型质量评分"
            rules={[{ required: true, message: '请给模型质量评分' }]}
          >
            <Rate count={5} />
          </Form.Item>

          <Form.Item
            name="model_comment"
            label="模型评价"
          >
            <Input.TextArea rows={3} placeholder="请输入对模型的评价（可选）" />
          </Form.Item>

          <Form.Item
            name="print_rating"
            label="打印效果评分"
            rules={[{ required: true, message: '请给打印效果评分' }]}
          >
            <Rate count={5} />
          </Form.Item>

          <Form.Item
            name="print_comment"
            label="打印评价"
          >
            <Input.TextArea rows={3} placeholder="请输入对打印效果的评价（可选）" />
          </Form.Item>

          <Form.Item name="is_anonymous" valuePropName="checked">
            <Checkbox>匿名评价</Checkbox>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default MyOrdersPage;
