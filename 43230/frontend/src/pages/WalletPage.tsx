import React, { useEffect, useState } from 'react';
import {
  Card,
  Table,
  Tag,
  Button,
  Space,
  Typography,
  Statistic,
  Row,
  Col,
  Modal,
  Form,
  InputNumber,
  message,
  Empty,
  Tabs,
} from 'antd';
import {
  WalletOutlined,
  PlusOutlined,
  ArrowUpOutlined,
  ArrowDownOutlined,
  ShoppingCartOutlined,
  GiftOutlined,
  CheckCircleOutlined,
} from '@ant-design/icons';
import { authApi } from '@/services/api';
import { Transaction } from '@/types';
import { formatPrice } from '@/utils/format';
import { formatDate } from '@/utils/date';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import { getProfile } from '@/store/authSlice';

const { Title, Text } = Typography;
const { TabPane } = Tabs;

const WalletPage: React.FC = () => {
  const dispatch = useAppDispatch();
  const { user } = useAppSelector((state) => state.auth);
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [loading, setLoading] = useState(false);
  const [page, setPage] = useState(1);
  const [pageSize] = useState(20);
  const [total, setTotal] = useState(0);
  const [rechargeVisible, setRechargeVisible] = useState(false);
  const [withdrawVisible, setWithdrawVisible] = useState(false);
  const [form] = Form.useForm();

  const fetchData = async () => {
    setLoading(true);
    try {
      const response = await authApi.getTransactions({ page, page_size: pageSize });
      setTransactions(response.data.data || []);
      setTotal(response.data.total || 0);
    } catch (error) {
      console.error('Failed to fetch transactions:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    dispatch(getProfile() as any);
    fetchData();
  }, [page]);

  const handleRecharge = async () => {
    try {
      const values = await form.validateFields();
      // TODO: 调用充值接口
      message.success(`充值 ${formatPrice(values.amount)} 成功！`);
      setRechargeVisible(false);
      form.resetFields();
      dispatch(getProfile() as any);
      fetchData();
    } catch (error) {
      console.error('Recharge failed:', error);
    }
  };

  const handleWithdraw = async () => {
    try {
      const values = await form.validateFields();
      if (values.amount > (user?.balance || 0)) {
        message.error('余额不足');
        return;
      }
      // TODO: 调用提现接口
      message.success(`提现 ${formatPrice(values.amount)} 申请已提交！`);
      setWithdrawVisible(false);
      form.resetFields();
      dispatch(getProfile() as any);
      fetchData();
    } catch (error) {
      console.error('Withdraw failed:', error);
    }
  };

  const typeMap: Record<string, { text: string; color: string; icon: React.ReactNode }> = {
    income: { text: '收入', color: 'success', icon: <ArrowUpOutlined /> },
    expense: { text: '支出', color: 'error', icon: <ArrowDownOutlined /> },
    refund: { text: '退款', color: 'warning', icon: <GiftOutlined /> },
  };

  const columns = [
    {
      title: '时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date: string) => formatDate(date),
    },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
      render: (type: string) => {
        const info = typeMap[type] || typeMap.expense;
        return (
          <Tag color={info.color} icon={info.icon}>
            {info.text}
          </Tag>
        );
      },
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
      render: (text: string, record: Transaction) => (
        <div>
          <div>{text}</div>
          <div className="text-xs text-gray-400">{record.transaction_no}</div>
        </div>
      ),
    },
    {
      title: '金额',
      dataIndex: 'amount',
      key: 'amount',
      render: (amount: number, record: Transaction) => (
        <div className={record.type === 'income' ? 'text-green-600 font-medium' : 'text-red-500 font-medium'}>
          {record.type === 'income' ? '+' : '-'}
          {formatPrice(amount)}
        </div>
      ),
    },
    {
      title: '余额',
      dataIndex: 'balance_after',
      key: 'balance_after',
      render: (val: number) => formatPrice(val),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => {
        const statusMap: Record<string, { text: string; color: string }> = {
          completed: { text: '已完成', color: 'success' },
          pending: { text: '处理中', color: 'warning' },
          failed: { text: '失败', color: 'error' },
        };
        const info = statusMap[status] || statusMap.completed;
        return <Tag color={info.color}>{info.text}</Tag>;
      },
    },
  ];

  const totalIncome = transactions
    .filter((t) => t.type === 'income')
    .reduce((sum, t) => sum + t.amount, 0);
  const totalExpense = transactions
    .filter((t) => t.type === 'expense')
    .reduce((sum, t) => sum + t.amount, 0);

  return (
    <div className="space-y-6">
      <Card>
        <Row gutter={24} align="middle">
          <Col xs={24} md={12}>
            <div className="flex items-center gap-4">
              <div className="w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center">
                <WalletOutlined className="text-3xl text-blue-600" />
              </div>
              <div>
                <Text type="secondary">账户余额</Text>
                <Title level={2} className="!mb-0 !text-blue-600">
                  {formatPrice(user?.balance || 0)}
                </Title>
              </div>
            </div>
          </Col>
          <Col xs={24} md={12} className="text-right">
            <Space size="large">
              <Button type="primary" size="large" icon={<PlusOutlined />} onClick={() => setRechargeVisible(true)}>
                充值
              </Button>
              <Button size="large" onClick={() => setWithdrawVisible(true)}>
                提现
              </Button>
            </Space>
          </Col>
        </Row>
      </Card>

      <Row gutter={16}>
        <Col xs={12} md={8}>
          <Card>
            <Statistic
              title="总收入"
              value={totalIncome}
              precision={2}
              prefix="¥"
              valueStyle={{ color: '#3f8600' }}
              prefixCls="ant-statistic"
            />
          </Card>
        </Col>
        <Col xs={12} md={8}>
          <Card>
            <Statistic
              title="总支出"
              value={totalExpense}
              precision={2}
              prefix="¥"
              valueStyle={{ color: '#cf1322' }}
            />
          </Card>
        </Col>
        <Col xs={12} md={8}>
          <Card>
            <Statistic
              title="交易笔数"
              value={total}
              suffix="笔"
            />
          </Card>
        </Col>
      </Row>

      <Card title="交易记录">
        <Table
          rowKey="id"
          columns={columns}
          dataSource={transactions}
          loading={loading}
          pagination={{
            current: page,
            pageSize,
            total,
            onChange: setPage,
          }}
          locale={{ emptyText: <Empty description="暂无交易记录" /> }}
        />
      </Card>

      {/* 充值弹窗 */}
      <Modal
        title="账户充值"
        open={rechargeVisible}
        onCancel={() => setRechargeVisible(false)}
        onOk={handleRecharge}
        okText="确认充值"
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="amount"
            label="充值金额"
            rules={[
              { required: true, message: '请输入充值金额' },
              { type: 'number', min: 1, message: '最小充值金额为1元' },
            ]}
          >
            <InputNumber
              className="w-full"
              min={1}
              step={10}
              prefix="¥"
              placeholder="请输入充值金额"
            />
          </Form.Item>
          <div className="text-sm text-gray-500">
            <p>• 支持支付宝、微信支付</p>
            <p>• 充值金额即时到账</p>
            <p>• 余额可用于购买模型和支付打印订单</p>
          </div>
        </Form>
      </Modal>

      {/* 提现弹窗 */}
      <Modal
        title="账户提现"
        open={withdrawVisible}
        onCancel={() => setWithdrawVisible(false)}
        onOk={handleWithdraw}
        okText="确认提现"
      >
        <Form form={form} layout="vertical">
          <div className="mb-4 p-4 bg-gray-50 rounded-lg">
            <Text type="secondary">可提现金额</Text>
            <div className="text-2xl font-bold text-green-600">
              {formatPrice(user?.balance || 0)}
            </div>
          </div>
          <Form.Item
            name="amount"
            label="提现金额"
            rules={[
              { required: true, message: '请输入提现金额' },
              { type: 'number', min: 10, message: '最小提现金额为10元' },
            ]}
          >
            <InputNumber
              className="w-full"
              min={10}
              step={10}
              max={user?.balance || 0}
              prefix="¥"
              placeholder="请输入提现金额"
            />
          </Form.Item>
          <div className="text-sm text-gray-500">
            <p>• 提现将在1-3个工作日内到账</p>
            <p>• 单笔提现最低10元，最高10000元</p>
            <p>• 提现手续费2元/笔</p>
          </div>
        </Form>
      </Modal>
    </div>
  );
};

export default WalletPage;
