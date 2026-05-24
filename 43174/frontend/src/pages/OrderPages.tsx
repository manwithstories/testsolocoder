import React, { useState, useEffect } from 'react';
import { Table, Tag, Button, Space, Modal, message, Card, Descriptions, Timeline } from 'antd';
import { EyeOutlined, CheckOutlined } from '@ant-design/icons';
import { orderApi } from '../services/api';
import { Order } from '../types';
import { Loading } from '../components/Loading';
import { useAuthStore } from '../context/authStore';

export const OrderListPage: React.FC = () => {
  const [loading, setLoading] = useState(true);
  const [orders, setOrders] = useState<Order[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize] = useState(10);
  const [statusFilter, setStatusFilter] = useState<string | undefined>();
  const [detailModalVisible, setDetailModalVisible] = useState(false);
  const [selectedOrder, setSelectedOrder] = useState<Order | null>(null);
  const { user: _user } = useAuthStore();

  useEffect(() => {
    loadOrders();
  }, [page, statusFilter]);

  const loadOrders = async () => {
    setLoading(true);
    try {
      const params: any = { page, page_size: pageSize };
      if (statusFilter) params.status = statusFilter;
      const response: any = await orderApi.getMyOrders(params);
      setOrders(response.data || []);
      setTotal(response.total || 0);
    } catch (error) {
      console.error('Failed to load orders:', error);
    } finally {
      setLoading(false);
    }
  };

  const getStatusColor = (status: string) => {
    const colors: Record<string, string> = {
      pending: 'orange',
      paid: 'blue',
      shipped: 'cyan',
      delivered: 'green',
      completed: 'success',
      cancelled: 'default',
      refunded: 'red',
    };
    return colors[status] || 'default';
  };

  const getStatusText = (status: string) => {
    const texts: Record<string, string> = {
      pending: '待支付',
      paid: '已支付',
      shipped: '已发货',
      delivered: '已送达',
      completed: '已完成',
      cancelled: '已取消',
      refunded: '已退款',
    };
    return texts[status] || status;
  };

  const handleViewDetail = (order: Order) => {
    setSelectedOrder(order);
    setDetailModalVisible(true);
  };

  const handlePay = async (order: Order) => {
    Modal.confirm({
      title: '确认支付',
      content: `您将要支付 ¥${order.total_amount}`,
      onOk: async () => {
        try {
          await orderApi.pay(order.id);
          message.success('支付成功');
          loadOrders();
        } catch (error: any) {
          message.error(error.message || '支付失败');
        }
      },
    });
  };

  const handleDeliver = async (order: Order) => {
    Modal.confirm({
      title: '确认收货',
      content: '确认您已收到商品？',
      onOk: async () => {
        try {
          await orderApi.deliver(order.id);
          message.success('收货成功');
          loadOrders();
        } catch (error: any) {
          message.error(error.message || '操作失败');
        }
      },
    });
  };

  const handleComplete = async (order: Order) => {
    Modal.confirm({
      title: '确认完成',
      content: '确认完成此订单？',
      onOk: async () => {
        try {
          await orderApi.complete(order.id);
          message.success('订单已完成');
          loadOrders();
        } catch (error: any) {
          message.error(error.message || '操作失败');
        }
      },
    });
  };

  const handleCancel = async (order: Order) => {
    Modal.confirm({
      title: '取消订单',
      content: '确定要取消此订单吗？',
      onOk: async () => {
        try {
          await orderApi.cancel(order.id);
          message.success('订单已取消');
          loadOrders();
        } catch (error: any) {
          message.error(error.message || '操作失败');
        }
      },
    });
  };

  const columns = [
    {
      title: '订单号',
      dataIndex: 'order_no',
      key: 'order_no',
    },
    {
      title: '金额',
      dataIndex: 'total_amount',
      key: 'total_amount',
      render: (amount: number) => <span className="text-red-500 font-bold">¥{amount}</span>,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={getStatusColor(status)}>{getStatusText(status)}</Tag>
      ),
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date: string) => new Date(date).toLocaleString(),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Order) => (
        <Space>
          <Button
            size="small"
            icon={<EyeOutlined />}
            onClick={() => handleViewDetail(record)}
          >
            详情
          </Button>
          {record.status === 'pending' && (
            <Button
              size="small"
              type="primary"
              onClick={() => handlePay(record)}
            >
              支付
            </Button>
          )}
          {record.status === 'pending' && (
            <Button
              size="small"
              danger
              onClick={() => handleCancel(record)}
            >
              取消
            </Button>
          )}
          {record.status === 'shipped' && (
            <Button
              size="small"
              type="primary"
              onClick={() => handleDeliver(record)}
            >
              确认收货
            </Button>
          )}
          {record.status === 'delivered' && (
            <Button
              size="small"
              type="primary"
              icon={<CheckOutlined />}
              onClick={() => handleComplete(record)}
            >
              完成
            </Button>
          )}
        </Space>
      ),
    },
  ];

  if (loading) return <Loading />;

  return (
    <div className="max-w-6xl mx-auto px-4 py-8">
      <h1 className="text-2xl font-bold mb-6">我的订单</h1>

      <Card className="mb-6">
        <Space>
          <span>状态筛选:</span>
          <Button
            type={!statusFilter ? 'primary' : 'default'}
            onClick={() => { setStatusFilter(undefined); setPage(1); }}
          >
            全部
          </Button>
          <Button
            type={statusFilter === 'pending' ? 'primary' : 'default'}
            onClick={() => { setStatusFilter('pending'); setPage(1); }}
          >
            待支付
          </Button>
          <Button
            type={statusFilter === 'paid' ? 'primary' : 'default'}
            onClick={() => { setStatusFilter('paid'); setPage(1); }}
          >
            已支付
          </Button>
          <Button
            type={statusFilter === 'shipped' ? 'primary' : 'default'}
            onClick={() => { setStatusFilter('shipped'); setPage(1); }}
          >
            已发货
          </Button>
          <Button
            type={statusFilter === 'completed' ? 'primary' : 'default'}
            onClick={() => { setStatusFilter('completed'); setPage(1); }}
          >
            已完成
          </Button>
        </Space>
      </Card>

      <Table
        columns={columns}
        dataSource={orders}
        rowKey="id"
        loading={loading}
        pagination={{
          current: page,
          pageSize,
          total,
          onChange: setPage,
          showSizeChanger: false,
        }}
      />

      <Modal
        title="订单详情"
        open={detailModalVisible}
        onCancel={() => setDetailModalVisible(false)}
        footer={null}
        width={600}
      >
        {selectedOrder && (
          <div>
            <Descriptions column={2} bordered size="small" className="mb-4">
              <Descriptions.Item label="订单号">{selectedOrder.order_no}</Descriptions.Item>
              <Descriptions.Item label="状态">
                <Tag color={getStatusColor(selectedOrder.status)}>
                  {getStatusText(selectedOrder.status)}
                </Tag>
              </Descriptions.Item>
              <Descriptions.Item label="金额">¥{selectedOrder.total_amount}</Descriptions.Item>
              <Descriptions.Item label="支付方式">{selectedOrder.payment_method || '-'}</Descriptions.Item>
              <Descriptions.Item label="收货地址" span={2}>
                {selectedOrder.shipping_address || '-'}
              </Descriptions.Item>
              <Descriptions.Item label="物流单号">
                {selectedOrder.tracking_number || '-'}
              </Descriptions.Item>
              <Descriptions.Item label="创建时间">
                {new Date(selectedOrder.created_at).toLocaleString()}
              </Descriptions.Item>
            </Descriptions>

            {selectedOrder.items && selectedOrder.items.length > 0 && (
              <div className="mb-4">
                <h4 className="font-semibold mb-2">商品列表</h4>
                {selectedOrder.items.map((item) => (
                  <div key={item.id} className="flex justify-between py-2 border-b">
                    <span>{item.textbook?.title || '教材'}</span>
                    <span>
                      ¥{item.price} x {item.quantity} = ¥{item.subtotal}
                    </span>
                  </div>
                ))}
              </div>
            )}

            {selectedOrder.status_history && selectedOrder.status_history.length > 0 && (
              <div>
                <h4 className="font-semibold mb-2">订单状态</h4>
                <Timeline
                  items={selectedOrder.status_history.map((h) => ({
                    color: h.status === 'completed' ? 'green' : 'blue',
                    children: (
                      <div>
                        <p className="font-medium">{getStatusText(h.status)}</p>
                        <p className="text-sm text-gray-500">{h.remark}</p>
                        <p className="text-xs text-gray-400">
                          {new Date(h.created_at).toLocaleString()}
                        </p>
                      </div>
                    ),
                  }))}
                />
              </div>
            )}
          </div>
        )}
      </Modal>
    </div>
  );
};
