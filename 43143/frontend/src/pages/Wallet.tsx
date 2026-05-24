import React, { useEffect, useState } from 'react';
import { Card, Row, Col, Statistic, List, Button, Modal, Form, Input, message, Spin, Tag } from 'antd';
import { WalletOutlined, PlusOutlined, MinusOutlined, ExportOutlined } from '@ant-design/icons';
import { paymentApi } from '@/api/message';
import { Wallet as WalletType, Payment } from '@/types';

const Wallet: React.FC = () => {
  const [loading, setLoading] = useState(true);
  const [wallet, setWallet] = useState<WalletType | null>(null);
  const [payments, setPayments] = useState<Payment[]>([]);
  const [withdrawModalVisible, setWithdrawModalVisible] = useState(false);
  const [withdrawLoading, setWithdrawLoading] = useState(false);
  const [withdrawForm] = Form.useForm();

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    setLoading(true);
    try {
      const [walletData, paymentsData] = await Promise.all([
        paymentApi.getWallet(),
        paymentApi.list({ page_size: 20 }),
      ]);

      setWallet(walletData);
      setPayments(paymentsData.items || []);
    } catch (error: any) {
      message.error(error.message || '获取钱包信息失败');
    } finally {
      setLoading(false);
    }
  };

  const handleWithdraw = async () => {
    try {
      const values = await withdrawForm.validateFields();
      setWithdrawLoading(true);

      await paymentApi.withdraw(values.amount);

      message.success('提现申请已提交');
      setWithdrawModalVisible(false);
      withdrawForm.resetFields();
      fetchData();
    } catch (error: any) {
      if (!error.errorFields) {
        message.error(error.message || '提现失败');
      }
    } finally {
      setWithdrawLoading(false);
    }
  };

  const getPaymentTypeText = (type: string) => {
    const types: Record<string, string> = {
      payment: '支付',
      refund: '退款',
      withdraw: '提现',
      deposit: '充值',
      fee: '手续费',
    };
    return types[type] || type;
  };

  const getPaymentStatusColor = (status: string) => {
    const colors: Record<string, string> = {
      pending: 'orange',
      processing: 'blue',
      completed: 'green',
      failed: 'red',
      cancelled: 'default',
    };
    return colors[status] || 'default';
  };

  const getPaymentStatusText = (status: string) => {
    const texts: Record<string, string> = {
      pending: '处理中',
      processing: '处理中',
      completed: '已完成',
      failed: '失败',
      cancelled: '已取消',
    };
    return texts[status] || status;
  };

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: 100 }}>
        <Spin size="large" />
      </div>
    );
  }

  return (
    <div>
      <Row gutter={16} style={{ marginBottom: 16 }}>
        <Col span={8}>
          <Card>
            <Statistic
              title="账户余额"
              value={wallet?.balance || 0}
              precision={2}
              prefix={<WalletOutlined />}
              suffix="元"
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col span={8}>
          <Card>
            <Statistic
              title="冻结金额"
              value={wallet?.frozen || 0}
              precision={2}
              suffix="元"
              valueStyle={{ color: '#faad14' }}
            />
          </Card>
        </Col>
        <Col span={8}>
          <Card>
            <Button
              type="primary"
              icon={<PlusOutlined />}
              block
              size="large"
              style={{ marginBottom: 8 }}
              onClick={() => message.info('充值功能开发中')}
            >
              充值
            </Button>
            <Button
              icon={<MinusOutlined />}
              block
              size="large"
              onClick={() => setWithdrawModalVisible(true)}
            >
              提现
            </Button>
          </Card>
        </Col>
      </Row>

      <Card
        title="交易记录"
        extra={
          <Button icon={<ExportOutlined />} size="small">
            导出
          </Button>
        }
      >
        {payments.length === 0 ? (
          <div style={{ textAlign: 'center', padding: 40, color: '#999' }}>
            暂无交易记录
          </div>
        ) : (
          <List
            itemLayout="horizontal"
            dataSource={payments}
            renderItem={(item) => (
              <List.Item>
                <List.Item.Meta
                  title={
                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                      <span>{getPaymentTypeText(item.type)}</span>
                      <Tag color={getPaymentStatusColor(item.status)}>
                        {getPaymentStatusText(item.status)}
                      </Tag>
                    </div>
                  }
                  description={
                    <div>
                      <div>订单号: {item.transaction_id || item.id}</div>
                      <div style={{ color: '#999', fontSize: 12 }}>
                        {new Date(item.created_at).toLocaleString()}
                      </div>
                    </div>
                  }
                />
                <div
                  style={{
                    fontSize: 18,
                    fontWeight: 'bold',
                    color: item.type === 'payment' || item.type === 'deposit' ? '#52c41a' : '#ff4d4f',
                  }}
                >
                  {item.type === 'payment' || item.type === 'deposit' ? '+' : '-'}
                  ¥{item.amount.toFixed(2)}
                </div>
              </List.Item>
            )}
          />
        )}
      </Card>

      <Modal
        title="申请提现"
        open={withdrawModalVisible}
        onCancel={() => {
          setWithdrawModalVisible(false);
          withdrawForm.resetFields();
        }}
        onOk={handleWithdraw}
        confirmLoading={withdrawLoading}
        okText="确认提现"
        cancelText="取消"
      >
        <Form form={withdrawForm} layout="vertical">
          <Form.Item
            label="提现金额"
            name="amount"
            rules={[
              { required: true, message: '请输入提现金额' },
              {
                validator: (_, value) => {
                  if (value && value > (wallet?.balance || 0)) {
                    return Promise.reject(new Error('提现金额不能超过余额'));
                  }
                  if (value && value <= 0) {
                    return Promise.reject(new Error('提现金额必须大于0'));
                  }
                  return Promise.resolve();
                },
              },
            ]}
          >
            <Input
              type="number"
              prefix="¥"
              placeholder="请输入提现金额"
              addonAfter={`余额: ¥${wallet?.balance?.toFixed(2) || '0.00'}`}
            />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default Wallet;
