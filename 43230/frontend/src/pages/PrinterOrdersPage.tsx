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
  Form,
  Input,
  InputNumber,
  message,
  Descriptions,
  Empty,
  Statistic,
  Row,
  Col,
  List,
  Progress,
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  PrinterOutlined,
  SettingOutlined,
  CheckCircleOutlined,
  PlayCircleOutlined,
  ClockCircleOutlined,
  CloseCircleOutlined,
  DatabaseOutlined,
} from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { orderApi, printerApi } from '@/services/api';
import { PrintOrder, PrinterDevice, MaterialInventory, PrintSchedule } from '@/types';
import {
  formatPrice,
  getOrderStatusText,
  getOrderStatusColor,
  getPrinterStatusText,
  getPrinterStatusColor,
  getMaterialTypeText,
} from '@/utils/format';
import { formatDate } from '@/utils/date';
import { useAuth } from '@/hooks/useAuth';

const { Title, Text } = Typography;
const { Option } = Select;

// 需要添加一个工具函数
const getPrinterStatusText = (status: string): string => {
  const map: Record<string, string> = {
    idle: '空闲',
    printing: '打印中',
    maintenance: '维护中',
    offline: '离线',
  };
  return map[status] || status;
};

const getPrinterStatusColor = (status: string): string => {
  const map: Record<string, string> = {
    idle: 'green',
    printing: 'blue',
    maintenance: 'orange',
    offline: 'default',
  };
  return map[status] || 'default';
};

const PrinterOrdersPage: React.FC = () => {
  const [orders, setOrders] = useState<PrintOrder[]>([]);
  const [pendingOrders, setPendingOrders] = useState<PrintOrder[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [page, setPage] = useState(1);
  const [pageSize] = useState(10);
  const [status, setStatus] = useState<string>('');
  const { userRole } = useAuth();

  const fetchData = async () => {
    setLoading(true);
    try {
      const params: any = { page, page_size: pageSize };
      if (status) params.status = status;

      const [ordersRes, pendingRes] = await Promise.all([
        orderApi.listPrinterOrders(params),
        orderApi.getPending(),
      ]);

      setOrders(ordersRes.data.data || []);
      setTotal(ordersRes.data.total || 0);
      setPendingOrders(pendingRes.data || []);
    } catch (error) {
      console.error('Failed to fetch printer orders:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, [page, status]);

  const handleAcceptOrder = async (orderId: string) => {
    try {
      await orderApi.assignPrinter(orderId);
      message.success('接单成功');
      fetchData();
    } catch (error: any) {
      message.error(error.response?.data?.error || '接单失败');
    }
  };

  const handleStartPrinting = async (orderId: string) => {
    try {
      await orderApi.startPrinting(orderId);
      message.success('开始打印');
      fetchData();
    } catch (error: any) {
      message.error(error.response?.data?.error || '操作失败');
    }
  };

  const handleCompletePrinting = async (orderId: string) => {
    try {
      await orderApi.completePrinting(orderId);
      message.success('打印完成，进入质检');
      fetchData();
    } catch (error: any) {
      message.error(error.response?.data?.error || '操作失败');
    }
  };

  const handleApproveQuality = async (orderId: string) => {
    try {
      await orderApi.approveQuality(orderId);
      message.success('质检通过');
      fetchData();
    } catch (error: any) {
      message.error(error.response?.data?.error || '操作失败');
    }
  };

  const handleShipOrder = async (orderId: string) => {
    Modal.confirm({
      title: '确认发货',
      content: (
        <Form layout="vertical">
          <Form.Item
            name="tracking_number"
            label="物流单号"
            rules={[{ required: true, message: '请输入物流单号' }]}
          >
            <Input placeholder="请输入物流单号" />
          </Form.Item>
        </Form>
      ),
      onOk: async () => {
        try {
          await orderApi.shipOrder(orderId, { tracking_number: 'SF1234567890' });
          message.success('已发货');
          fetchData();
        } catch (error: any) {
          message.error(error.response?.data?.error || '操作失败');
        }
      },
    });
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
            {formatDate(record.created_at)}
          </div>
          <div className="text-sm text-gray-500">
            客户: {record.customer?.username}
          </div>
        </div>
      ),
    },
    {
      title: '模型',
      key: 'model',
      render: (_: any, record: PrintOrder) => (
        <div>
          <div className="font-medium">{record.model?.title}</div>
          <div className="text-sm text-gray-500">
            数量: {record.quantity}件
          </div>
          <div className="text-sm text-gray-500">
            材料: {getMaterialTypeText(record.material?.type || '')}
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
      render: (status: string) => (
        <Tag color={getOrderStatusColor(status as any)}>
          {getOrderStatusText(status as any)}
        </Tag>
      ),
    },
    {
      title: '操作',
      key: 'actions',
      render: (_: any, record: PrintOrder) => (
        <Space direction="vertical" size="small">
          {record.status === 'paid' && (
            <Button
              type="primary"
              size="small"
              onClick={() => handleStartPrinting(record.id)}
            >
              <PlayCircleOutlined /> 开始打印
            </Button>
          )}
          {record.status === 'printing' && (
            <Button
              type="primary"
              size="small"
              onClick={() => handleCompletePrinting(record.id)}
            >
              <CheckCircleOutlined /> 打印完成
            </Button>
          )}
          {record.status === 'quality_check' && (
            <Button
              type="primary"
              size="small"
              onClick={() => handleApproveQuality(record.id)}
            >
              <CheckCircleOutlined /> 质检通过
            </Button>
          )}
          {record.status === 'shipped' && (
            <Button type="link" size="small">
              已发货，等待签收
            </Button>
          )}
        </Space>
      ),
    },
  ];

  return (
    <div className="space-y-6">
      {/* 待接单列表 */}
      {pendingOrders.length > 0 && (
        <Card title="待接单订单" extra={<Tag color="red">{pendingOrders.length}个待处理</Tag>}>
          <List
            dataSource={pendingOrders}
            renderItem={(order) => (
              <List.Item
                key={order.id}
                actions={[
                  <Button type="primary" onClick={() => handleAcceptOrder(order.id)}>
                    接单
                  </Button>,
                ]}
              >
                <List.Item.Meta
                  title={order.order_no}
                  description={
                    <Space size="large">
                      <span>模型: {order.model?.title}</span>
                      <span>数量: {order.quantity}件</span>
                      <span className="text-red-500 font-bold">
                        {formatPrice(order.total_amount)}
                      </span>
                    </Space>
                  }
                />
              </List.Item>
            )}
          />
        </Card>
      )}

      <Card
        title="我的打印订单"
        extra={
          <Select
            placeholder="全部状态"
            style={{ width: 150 }}
            allowClear
            onChange={setStatus}
            value={status || undefined}
          >
            <Option value="paid">已付款待打印</Option>
            <Option value="printing">打印中</Option>
            <Option value="quality_check">质检中</Option>
            <Option value="shipped">已发货</Option>
            <Option value="delivered">已送达</Option>
            <Option value="completed">已完成</Option>
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
        />
      </Card>
    </div>
  );
};

export default PrinterOrdersPage;
