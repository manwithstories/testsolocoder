import React, { useState, useEffect } from 'react';
import {
  Card,
  Row,
  Col,
  Statistic,
  Table,
  Tag,
  Button,
  Modal,
  Select,
  message,
  Space,
  DatePicker,
} from 'antd';
import {
  BookOutlined,
  UserOutlined,
  ShoppingCartOutlined,
  CheckOutlined,
  CloseOutlined,
  ExportOutlined,
  DownloadOutlined,
} from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import {
  statisticsApi,
  userApi,
  textbookApi,
  noteApi,
  orderApi,
  reviewApi,
} from '../services/api';
import { User, Textbook, Note, Order, Review } from '../types';
import { Loading } from '../components/Loading';
import dayjs from 'dayjs';
import { saveAs } from 'file-saver';

const { Option } = Select;
const { RangePicker } = DatePicker;

export const AdminDashboard: React.FC = () => {
  const [loading, setLoading] = useState(true);
  const [stats, setStats] = useState<any>(null);
  const [monthlyStats, setMonthlyStats] = useState<any[]>([]);
  const navigate = useNavigate();

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    setLoading(true);
    try {
      const [textbookRes, userRes, orderRes, monthlyRes]: any = await Promise.all([
        statisticsApi.getTextbookStats(),
        statisticsApi.getUserStats(),
        statisticsApi.getOrderStats(),
        statisticsApi.getMonthlyStats(6),
      ]);
      setStats({
        textbooks: textbookRes.data,
        users: userRes.data,
        orders: orderRes.data,
      });
      setMonthlyStats(monthlyRes.data || []);
    } catch (error) {
      console.error('Failed to load statistics:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) return <Loading />;

  return (
    <div>
      <h1 className="text-2xl font-bold mb-6">数据概览</h1>

      <Row gutter={[16, 16]} className="mb-6">
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="教材总数"
              value={stats?.textbooks?.total_count || 0}
              prefix={<BookOutlined className="text-blue-500" />}
            />
            <div className="mt-2 text-sm text-gray-500">
              在售: {stats?.textbooks?.available_count || 0} |
              已售: {stats?.textbooks?.sold_count || 0}
            </div>
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="用户总数"
              value={stats?.users?.total_count || 0}
              prefix={<UserOutlined className="text-green-500" />}
            />
            <div className="mt-2 text-sm text-gray-500">
              学生: {stats?.users?.student_count || 0} |
              书商: {stats?.users?.merchant_count || 0}
            </div>
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="订单总数"
              value={stats?.orders?.total_count || 0}
              prefix={<ShoppingCartOutlined className="text-orange-500" />}
            />
            <div className="mt-2 text-sm text-gray-500">
              待处理: {stats?.orders?.pending_count || 0}
            </div>
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="总交易额"
              value={stats?.orders?.total_revenue || 0}
              precision={2}
              prefix="¥"
              valueStyle={{ color: '#cf1322' }}
            />
            <div className="mt-2 text-sm text-gray-500">
              已完成: {stats?.orders?.completed_count || 0} 单
            </div>
          </Card>
        </Col>
      </Row>

      <Card title="月度统计" className="mb-6">
        <Table
          dataSource={monthlyStats}
          rowKey="month"
          pagination={false}
          columns={[
            { title: '月份', dataIndex: 'month', key: 'month' },
            { title: '订单数', dataIndex: 'order_count', key: 'order_count' },
            {
              title: '交易额',
              dataIndex: 'revenue',
              key: 'revenue',
              render: (value: number) => `¥${value?.toFixed(2) || '0.00'}`,
            },
            { title: '新增用户', dataIndex: 'new_users', key: 'new_users' },
            { title: '新增教材', dataIndex: 'new_textbooks', key: 'new_textbooks' },
          ]}
        />
        <div className="mt-4 text-right">
          <Button
            icon={<ExportOutlined />}
            onClick={async () => {
              const month = dayjs().format('YYYY-MM');
              try {
                const response: any = await statisticsApi.exportReport(month);
                const blob = new Blob([response], {
                  type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
                });
                saveAs(blob, `monthly_report_${month}.xlsx`);
                message.success('报表导出成功');
              } catch (error: any) {
                message.error(error.message || '导出失败');
              }
            }}
          >
            导出本月报表
          </Button>
        </div>
      </Card>

      <Row gutter={[16, 16]}>
        <Col span={12}>
          <Card title="快速操作">
            <Space direction="vertical" className="w-full">
              <Button block onClick={() => navigate('/admin/users')}>
                用户管理
              </Button>
              <Button block onClick={() => navigate('/admin/textbooks')}>
                教材管理
              </Button>
              <Button block onClick={() => navigate('/admin/notes')}>
                笔记管理
              </Button>
              <Button block onClick={() => navigate('/admin/orders')}>
                订单管理
              </Button>
              <Button block onClick={() => navigate('/admin/reviews')}>
                评价管理
              </Button>
            </Space>
          </Card>
        </Col>
        <Col span={12}>
          <Card title="快捷导出">
            <Space direction="vertical" className="w-full">
              <div className="flex items-center gap-2">
                <RangePicker
                  className="flex-1"
                  picker="month"
                  format="YYYY-MM"
                />
                <Button
                  icon={<DownloadOutlined />}
                  onClick={async () => {
                    message.info('请选择月份后导出报表');
                  }}
                >
                  导出
                </Button>
              </div>
              <p className="text-gray-500 text-sm">选择月份后点击导出按钮导出该月数据报表</p>
            </Space>
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export const AdminUsersPage: React.FC = () => {
  const [loading, setLoading] = useState(true);
  const [users, setUsers] = useState<User[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize] = useState(10);
  const [roleFilter, setRoleFilter] = useState<string | undefined>();
  const [statusFilter, setStatusFilter] = useState<string | undefined>();

  useEffect(() => {
    loadUsers();
  }, [page, roleFilter, statusFilter]);

  const loadUsers = async () => {
    setLoading(true);
    try {
      const params: any = { page, page_size: pageSize };
      if (roleFilter) params.role = roleFilter;
      if (statusFilter) params.status = statusFilter;
      const response: any = await userApi.getUsers(params);
      setUsers(response.data || []);
      setTotal(response.total || 0);
    } catch (error) {
      console.error('Failed to load users:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleApprove = async (userId: string) => {
    try {
      await userApi.updateUserStatus(userId, 'active');
      message.success('审核通过');
      loadUsers();
    } catch (error: any) {
      message.error(error.message || '操作失败');
    }
  };

  const handleReject = async (userId: string) => {
    Modal.confirm({
      title: '拒绝用户',
      content: '确定要拒绝该用户的注册申请吗？',
      onOk: async () => {
        try {
          await userApi.updateUserStatus(userId, 'rejected');
          message.success('已拒绝');
          loadUsers();
        } catch (error: any) {
          message.error(error.message || '操作失败');
        }
      },
    });
  };

  const handleDelete = async (userId: string) => {
    Modal.confirm({
      title: '删除用户',
      content: '确定要删除该用户吗？此操作不可恢复。',
      okText: '确定删除',
      okButtonProps: { danger: true },
      onOk: async () => {
        try {
          await userApi.deleteUser(userId);
          message.success('删除成功');
          loadUsers();
        } catch (error: any) {
          message.error(error.message || '删除失败');
        }
      },
    });
  };

  const columns = [
    {
      title: '用户名',
      dataIndex: 'username',
      key: 'username',
    },
    {
      title: '邮箱',
      dataIndex: 'email',
      key: 'email',
    },
    {
      title: '角色',
      dataIndex: 'role',
      key: 'role',
      render: (role: string) => (
        <Tag color={role === 'admin' ? 'purple' : role === 'merchant' ? 'blue' : 'green'}>
          {role === 'admin' ? '管理员' : role === 'merchant' ? '书商' : '学生'}
        </Tag>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag
          color={
            status === 'active'
              ? 'success'
              : status === 'pending'
              ? 'orange'
              : status === 'rejected'
              ? 'red'
              : 'default'
          }
        >
          {status === 'active'
            ? '已认证'
            : status === 'pending'
            ? '待审核'
            : status === 'rejected'
            ? '已拒绝'
            : '已禁用'}
        </Tag>
      ),
    },
    {
      title: '评分',
      dataIndex: 'rating',
      key: 'rating',
      render: (rating: number) => rating?.toFixed(1) || '5.0',
    },
    {
      title: '注册时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date: string) => new Date(date).toLocaleDateString(),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: User) => (
        <Space>
          {record.status === 'pending' && (
            <>
              <Button
                size="small"
                type="primary"
                icon={<CheckOutlined />}
                onClick={() => handleApprove(record.id)}
              >
                通过
              </Button>
              <Button
                size="small"
                danger
                icon={<CloseOutlined />}
                onClick={() => handleReject(record.id)}
              >
                拒绝
              </Button>
            </>
          )}
          <Button size="small" danger onClick={() => handleDelete(record.id)}>
            删除
          </Button>
        </Space>
      ),
    },
  ];

  return (
    <div>
      <h1 className="text-2xl font-bold mb-6">用户管理</h1>
      <Card>
        <Space className="mb-4">
          <span>角色:</span>
          <Select
            placeholder="全部角色"
            allowClear
            style={{ width: 120 }}
            value={roleFilter}
            onChange={(value) => {
              setRoleFilter(value);
              setPage(1);
            }}
          >
            <Option value="student">学生</Option>
            <Option value="merchant">书商</Option>
          </Select>
          <span>状态:</span>
          <Select
            placeholder="全部状态"
            allowClear
            style={{ width: 120 }}
            value={statusFilter}
            onChange={(value) => {
              setStatusFilter(value);
              setPage(1);
            }}
          >
            <Option value="pending">待审核</Option>
            <Option value="active">已认证</Option>
            <Option value="rejected">已拒绝</Option>
          </Select>
        </Space>
        <Table
          columns={columns}
          dataSource={users}
          rowKey="id"
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

export const AdminTextbooksPage: React.FC = () => {
  const [loading, setLoading] = useState(true);
  const [textbooks, setTextbooks] = useState<Textbook[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize] = useState(10);
  const [statusFilter, setStatusFilter] = useState<string | undefined>();

  useEffect(() => {
    loadTextbooks();
  }, [page, statusFilter]);

  const loadTextbooks = async () => {
    setLoading(true);
    try {
      const params: any = { page, page_size: pageSize };
      if (statusFilter) params.status = statusFilter;
      const response: any = await textbookApi.getAll(params);
      setTextbooks(response.data || []);
      setTotal(response.total || 0);
    } catch (error) {
      console.error('Failed to load textbooks:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (textbookId: string) => {
    Modal.confirm({
      title: '删除教材',
      content: '确定要删除该教材吗？',
      onOk: async () => {
        try {
          await textbookApi.delete(textbookId);
          message.success('删除成功');
          loadTextbooks();
        } catch (error: any) {
          message.error(error.message || '删除失败');
        }
      },
    });
  };

  const columns = [
    {
      title: '教材名称',
      dataIndex: 'title',
      key: 'title',
    },
    {
      title: 'ISBN',
      dataIndex: 'isbn',
      key: 'isbn',
    },
    {
      title: '价格',
      dataIndex: 'price',
      key: 'price',
      render: (price: number) => <span className="text-red-500">¥{price}</span>,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={status === 'available' ? 'green' : status === 'reserved' ? 'orange' : 'red'}>
          {status === 'available' ? '在售' : status === 'reserved' ? '已预定' : '已售出'}
        </Tag>
      ),
    },
    {
      title: '发布者',
      dataIndex: ['seller', 'username'],
      key: 'seller',
    },
    {
      title: '浏览量',
      dataIndex: 'view_count',
      key: 'view_count',
    },
    {
      title: '发布时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date: string) => new Date(date).toLocaleDateString(),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Textbook) => (
        <Button size="small" danger onClick={() => handleDelete(record.id)}>
          删除
        </Button>
      ),
    },
  ];

  return (
    <div>
      <h1 className="text-2xl font-bold mb-6">教材管理</h1>
      <Card>
        <Space className="mb-4">
          <span>状态:</span>
          <Select
            placeholder="全部状态"
            allowClear
            style={{ width: 120 }}
            value={statusFilter}
            onChange={(value) => {
              setStatusFilter(value);
              setPage(1);
            }}
          >
            <Option value="available">在售</Option>
            <Option value="reserved">已预定</Option>
            <Option value="sold">已售出</Option>
          </Select>
        </Space>
        <Table
          columns={columns}
          dataSource={textbooks}
          rowKey="id"
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

export const AdminNotesPage: React.FC = () => {
  const [loading, setLoading] = useState(true);
  const [notes, setNotes] = useState<Note[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize] = useState(10);

  useEffect(() => {
    loadNotes();
  }, [page]);

  const loadNotes = async () => {
    setLoading(true);
    try {
      const response: any = await noteApi.getAll({ page, page_size: pageSize });
      setNotes(response.data || []);
      setTotal(response.total || 0);
    } catch (error) {
      console.error('Failed to load notes:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleToggleFeatured = async (noteId: string, currentStatus: boolean) => {
    Modal.confirm({
      title: currentStatus ? '取消精选' : '设为精选',
      content: `确定要${currentStatus ? '取消' : '设为'}精选笔记吗？`,
      onOk: async () => {
        try {
          await noteApi.setFeatured(noteId, !currentStatus);
          message.success('操作成功');
          loadNotes();
        } catch (error: any) {
          message.error(error.message || '操作失败');
        }
      },
    });
  };

  const handleDelete = async (noteId: string) => {
    Modal.confirm({
      title: '删除笔记',
      content: '确定要删除该笔记吗？',
      onOk: async () => {
        try {
          await noteApi.delete(noteId);
          message.success('删除成功');
          loadNotes();
        } catch (error: any) {
          message.error(error.message || '删除失败');
        }
      },
    });
  };

  const columns = [
    {
      title: '笔记标题',
      dataIndex: 'title',
      key: 'title',
    },
    {
      title: '科目',
      dataIndex: 'subject',
      key: 'subject',
    },
    {
      title: '精选',
      dataIndex: 'is_featured',
      key: 'is_featured',
      render: (isFeatured: boolean) =>
        isFeatured ? <Tag color="gold">精选</Tag> : null,
    },
    {
      title: '评分',
      dataIndex: 'rating',
      key: 'rating',
    },
    {
      title: '下载量',
      dataIndex: 'download_count',
      key: 'download_count',
    },
    {
      title: '上传者',
      dataIndex: ['uploader', 'username'],
      key: 'uploader',
    },
    {
      title: '上传时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date: string) => new Date(date).toLocaleDateString(),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Note) => (
        <Space>
          <Button
            size="small"
            type={record.is_featured ? 'default' : 'primary'}
            onClick={() => handleToggleFeatured(record.id, record.is_featured)}
          >
            {record.is_featured ? '取消精选' : '设为精选'}
          </Button>
          <Button size="small" danger onClick={() => handleDelete(record.id)}>
            删除
          </Button>
        </Space>
      ),
    },
  ];

  return (
    <div>
      <h1 className="text-2xl font-bold mb-6">笔记管理</h1>
      <Card>
        <Table
          columns={columns}
          dataSource={notes}
          rowKey="id"
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

export const AdminOrdersPage: React.FC = () => {
  const [loading, setLoading] = useState(true);
  const [orders, setOrders] = useState<Order[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize] = useState(10);
  const [statusFilter, setStatusFilter] = useState<string | undefined>();

  useEffect(() => {
    loadOrders();
  }, [page, statusFilter]);

  const loadOrders = async () => {
    setLoading(true);
    try {
      const params: any = { page, page_size: pageSize };
      if (statusFilter) params.status = statusFilter;
      const response: any = await orderApi.getAll(params);
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
      render: (status: string) => <Tag color={getStatusColor(status)}>{getStatusText(status)}</Tag>,
    },
    {
      title: '买家',
      dataIndex: ['buyer', 'username'],
      key: 'buyer',
    },
    {
      title: '卖家',
      dataIndex: ['seller', 'username'],
      key: 'seller',
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date: string) => new Date(date).toLocaleDateString(),
    },
  ];

  return (
    <div>
      <h1 className="text-2xl font-bold mb-6">订单管理</h1>
      <Card>
        <Space className="mb-4">
          <span>状态:</span>
          <Select
            placeholder="全部状态"
            allowClear
            style={{ width: 120 }}
            value={statusFilter}
            onChange={(value) => {
              setStatusFilter(value);
              setPage(1);
            }}
          >
            <Option value="pending">待支付</Option>
            <Option value="paid">已支付</Option>
            <Option value="shipped">已发货</Option>
            <Option value="delivered">已送达</Option>
            <Option value="completed">已完成</Option>
            <Option value="cancelled">已取消</Option>
          </Select>
        </Space>
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
          }}
        />
      </Card>
    </div>
  );
};

export const AdminReviewsPage: React.FC = () => {
  const [loading, setLoading] = useState(true);
  const [reviews, setReviews] = useState<Review[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize] = useState(10);
  const [isMaliciousFilter, setIsMaliciousFilter] = useState(false);

  useEffect(() => {
    loadReviews();
  }, [page, isMaliciousFilter]);

  const loadReviews = async () => {
    setLoading(true);
    try {
      const params: any = { page, page_size: pageSize };
      if (isMaliciousFilter) params.is_malicious = true;
      const response: any = await reviewApi.getAll(params);
      setReviews(response.data || []);
      setTotal(response.total || 0);
    } catch (error) {
      console.error('Failed to load reviews:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleHide = async (reviewId: string) => {
    try {
      await reviewApi.hide(reviewId);
      message.success('已隐藏');
      loadReviews();
    } catch (error: any) {
      message.error(error.message || '操作失败');
    }
  };

  const handleMarkMalicious = async (reviewId: string, isMalicious: boolean) => {
    Modal.confirm({
      title: isMalicious ? '取消恶意标记' : '标记为恶意评价',
      content: `确定要${isMalicious ? '取消' : '标记为'}恶意评价吗？`,
      onOk: async () => {
        try {
          await reviewApi.markMalicious(reviewId, !isMalicious);
          message.success('操作成功');
          loadReviews();
        } catch (error: any) {
          message.error(error.message || '操作失败');
        }
      },
    });
  };

  const columns = [
    {
      title: '评价内容',
      dataIndex: 'content',
      key: 'content',
      render: (content: string) => content || '（无内容）',
    },
    {
      title: '评分',
      dataIndex: 'rating',
      key: 'rating',
      render: (rating: number) => `${rating} / 5`,
    },
    {
      title: '用户',
      dataIndex: ['user', 'username'],
      key: 'user',
    },
    {
      title: '恶意',
      dataIndex: 'is_malicious',
      key: 'is_malicious',
      render: (isMalicious: boolean) =>
        isMalicious ? <Tag color="red">恶意</Tag> : null,
    },
    {
      title: '评价时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date: string) => new Date(date).toLocaleDateString(),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Review) => (
        <Space>
          <Button size="small" onClick={() => handleHide(record.id)}>
            隐藏
          </Button>
          <Button
            size="small"
            danger={!record.is_malicious}
            onClick={() => handleMarkMalicious(record.id, record.is_malicious)}
          >
            {record.is_malicious ? '取消标记' : '标记恶意'}
          </Button>
        </Space>
      ),
    },
  ];

  return (
    <div>
      <h1 className="text-2xl font-bold mb-6">评价管理</h1>
      <Card>
        <Space className="mb-4">
          <Button
            type={isMaliciousFilter ? 'primary' : 'default'}
            onClick={() => {
              setIsMaliciousFilter(!isMaliciousFilter);
              setPage(1);
            }}
          >
            只看恶意评价
          </Button>
        </Space>
        <Table
          columns={columns}
          dataSource={reviews}
          rowKey="id"
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
