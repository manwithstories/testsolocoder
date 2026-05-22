import React, { useState, useEffect } from 'react'
import {
  Card,
  Table,
  Button,
  Space,
  Typography,
  Tag,
  Modal,
  Descriptions,
  message,
  Spin,
  Empty,
  Pagination,
  Select,
  DatePicker,
  Form,
  Row,
  Col,
  Statistic
} from 'antd'
import {
  CreditCardOutlined,
  EyeOutlined,
  FileExcelOutlined,
  CloseOutlined,
  SearchOutlined,
  ReloadOutlined
} from '@ant-design/icons'
import type { ColumnsType } from 'antd/es/table'
import type { Payment } from '@/types'
import { paymentAPI } from '@/services/api'
import { saveAs } from 'file-saver'
import dayjs, { Dayjs } from 'dayjs'

const { Title, Text } = Typography
const { RangePicker } = DatePicker

const statusMap: Record<string, { text: string; color: string }> = {
  pending: { text: '待支付', color: 'orange' },
  paid: { text: '已支付', color: 'green' },
  failed: { text: '支付失败', color: 'red' },
  refunded: { text: '已退款', color: 'default' }
}

const methodMap: Record<string, string> = {
  wechat: '微信支付',
  alipay: '支付宝',
  card: '银行卡',
  cash: '现金'
}

interface FilterFormData {
  status?: string
  dateRange?: [Dayjs, Dayjs]
}

const PaymentsPage: React.FC = () => {
  const [form] = Form.useForm<FilterFormData>()
  const [loading, setLoading] = useState(false)
  const [exporting, setExporting] = useState(false)
  const [payments, setPayments] = useState<Payment[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize] = useState(10)
  const [detailModalVisible, setDetailModalVisible] = useState(false)
  const [selectedPayment, setSelectedPayment] = useState<Payment | null>(null)
  const [stats, setStats] = useState({
    totalAmount: 0,
    paidAmount: 0,
    pendingAmount: 0,
    count: 0
  })

  useEffect(() => {
    fetchPayments()
  }, [page, pageSize])

  const fetchPayments = async (filters?: FilterFormData) => {
    setLoading(true)
    try {
      const params: any = {
        page,
        pageSize
      }
      if (filters?.status) {
        params.status = filters.status
      }
      if (filters?.dateRange && filters.dateRange.length === 2) {
        params.start_date = filters.dateRange[0].format('YYYY-MM-DD')
        params.end_date = filters.dateRange[1].format('YYYY-MM-DD')
      }
      const result = await paymentAPI.getList(params)
      setPayments(result.list)
      setTotal(result.total)
      calculateStats(result.list)
    } catch (error) {
      console.error('获取费用列表失败:', error)
      message.error('获取费用列表失败')
    } finally {
      setLoading(false)
    }
  }

  const calculateStats = (paymentList: Payment[]) => {
    const totalAmount = paymentList.reduce((sum, p) => sum + p.total_amount, 0)
    const paidAmount = paymentList
      .filter((p) => p.status === 'paid')
      .reduce((sum, p) => sum + p.total_amount, 0)
    const pendingAmount = paymentList
      .filter((p) => p.status === 'pending')
      .reduce((sum, p) => sum + p.total_amount, 0)
    setStats({
      totalAmount,
      paidAmount,
      pendingAmount,
      count: paymentList.length
    })
  }

  const handleSearch = (values: FilterFormData) => {
    setPage(1)
    fetchPayments(values)
  }

  const handleReset = () => {
    form.resetFields()
    setPage(1)
    fetchPayments()
  }

  const handleViewDetail = (payment: Payment) => {
    setSelectedPayment(payment)
    setDetailModalVisible(true)
  }

  const handleExport = async () => {
    setExporting(true)
    try {
      const values = form.getFieldsValue()
      const params: any = {
        start_date: values.dateRange?.[0]?.format('YYYY-MM-DD') || dayjs().subtract(30, 'day').format('YYYY-MM-DD'),
        end_date: values.dateRange?.[1]?.format('YYYY-MM-DD') || dayjs().format('YYYY-MM-DD')
      }
      if (values.status) {
        params.status = values.status
      }
      const blob = await paymentAPI.exportReport(params)
      saveAs(blob, `费用报表_${dayjs().format('YYYYMMDDHHmmss')}.xlsx`)
      message.success('导出成功')
    } catch (error) {
      console.error('导出失败:', error)
      message.error('导出失败')
    } finally {
      setExporting(false)
    }
  }

  const handlePageChange = (newPage: number) => {
    setPage(newPage)
    fetchPayments(form.getFieldsValue())
  }

  const columns: ColumnsType<Payment> = [
    {
      title: '订单号',
      dataIndex: 'transaction_no',
      key: 'transaction_no',
      width: 180,
      render: (text: string) => text || '-'
    },
    {
      title: '就诊日期',
      key: 'appointment_date',
      width: 120,
      render: (_, record) =>
        record.appointment?.appointment_date
          ? new Date(record.appointment.appointment_date).toLocaleDateString()
          : '-'
    },
    {
      title: '患者',
      key: 'patient',
      width: 100,
      render: (_, record) => record.appointment?.patient?.user?.full_name || '-'
    },
    {
      title: '医生',
      key: 'doctor',
      width: 100,
      render: (_, record) => record.appointment?.doctor?.user?.full_name || '-'
    },
    {
      title: '科室',
      key: 'department',
      width: 100,
      render: (_, record) => record.appointment?.doctor?.department?.name || '-'
    },
    {
      title: '总金额',
      dataIndex: 'total_amount',
      key: 'total_amount',
      width: 100,
      render: (amount: number) => `¥${amount.toFixed(2)}`
    },
    {
      title: '支付方式',
      dataIndex: 'method',
      key: 'method',
      width: 100,
      render: (method: string) => methodMap[method] || method
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (status: string) => {
        const statusInfo = statusMap[status] || { text: status, color: 'default' }
        return <Tag color={statusInfo.color}>{statusInfo.text}</Tag>
      }
    },
    {
      title: '支付时间',
      dataIndex: 'paid_at',
      key: 'paid_at',
      width: 160,
      render: (date: string | null) => (date ? new Date(date).toLocaleString() : '-')
    },
    {
      title: '操作',
      key: 'action',
      width: 100,
      fixed: 'right',
      render: (_, record) => (
        <Button
          type="link"
          icon={<EyeOutlined />}
          onClick={() => handleViewDetail(record)}
        >
          详情
        </Button>
      )
    }
  ]

  return (
    <div className="space-y-6">
      <Card>
        <div className="flex items-center justify-between mb-4">
          <Title level={3} style={{ margin: 0 }}>
            <CreditCardOutlined className="mr-2" />
            费用管理
          </Title>
          <Button
            type="primary"
            icon={<FileExcelOutlined />}
            loading={exporting}
            onClick={handleExport}
          >
            导出报表
          </Button>
        </div>

        <Row gutter={[16, 16]} className="mb-4">
          <Col xs={24} sm={12} md={6}>
            <Card size="small">
              <Statistic
                title="总金额"
                value={stats.totalAmount}
                prefix="¥"
                precision={2}
                valueStyle={{ color: '#1890ff' }}
              />
            </Card>
          </Col>
          <Col xs={24} sm={12} md={6}>
            <Card size="small">
              <Statistic
                title="已支付"
                value={stats.paidAmount}
                prefix="¥"
                precision={2}
                valueStyle={{ color: '#52c41a' }}
              />
            </Card>
          </Col>
          <Col xs={24} sm={12} md={6}>
            <Card size="small">
              <Statistic
                title="待支付"
                value={stats.pendingAmount}
                prefix="¥"
                precision={2}
                valueStyle={{ color: '#fa8c16' }}
              />
            </Card>
          </Col>
          <Col xs={24} sm={12} md={6}>
            <Card size="small">
              <Statistic
                title="订单数"
                value={stats.count}
                valueStyle={{ color: '#722ed1' }}
              />
            </Card>
          </Col>
        </Row>

        <Form
          form={form}
          layout="inline"
          onFinish={handleSearch}
          className="mb-4"
        >
          <Form.Item name="status" label="状态">
            <Select
              placeholder="全部状态"
              allowClear
              style={{ width: 150 }}
              options={[
                { value: 'pending', label: '待支付' },
                { value: 'paid', label: '已支付' },
                { value: 'failed', label: '支付失败' },
                { value: 'refunded', label: '已退款' }
              ]}
            />
          </Form.Item>
          <Form.Item name="dateRange" label="日期范围">
            <RangePicker />
          </Form.Item>
          <Form.Item>
            <Space>
              <Button type="primary" htmlType="submit" icon={<SearchOutlined />}>
                查询
              </Button>
              <Button icon={<ReloadOutlined />} onClick={handleReset}>
                重置
              </Button>
            </Space>
          </Form.Item>
        </Form>

        {loading && payments.length === 0 ? (
          <div className="flex justify-center items-center min-h-[400px]">
            <Spin size="large" />
          </div>
        ) : payments.length === 0 ? (
          <Empty description="暂无费用记录" />
        ) : (
          <>
            <Table
              columns={columns}
              dataSource={payments}
              rowKey="id"
              pagination={false}
              loading={loading}
              scroll={{ x: 1300 }}
            />
            <div className="flex justify-end mt-4">
              <Pagination
                current={page}
                pageSize={pageSize}
                total={total}
                onChange={handlePageChange}
                showSizeChanger={false}
                showQuickJumper
                showTotal={(total) => `共 ${total} 条记录`}
              />
            </div>
          </>
        )}
      </Card>

      <Modal
        title="费用详情"
        open={detailModalVisible}
        onCancel={() => setDetailModalVisible(false)}
        footer={[
          <Button
            key="close"
            icon={<CloseOutlined />}
            onClick={() => setDetailModalVisible(false)}
          >
            关闭
          </Button>
        ]}
        width={700}
      >
        {selectedPayment && (
          <div className="space-y-4">
            <Descriptions bordered column={2} size="small">
              <Descriptions.Item label="订单号">
                {selectedPayment.transaction_no || '-'}
              </Descriptions.Item>
              <Descriptions.Item label="状态">
                {statusMap[selectedPayment.status] && (
                  <Tag color={statusMap[selectedPayment.status].color}>
                    {statusMap[selectedPayment.status].text}
                  </Tag>
                )}
              </Descriptions.Item>
              <Descriptions.Item label="支付方式">
                {methodMap[selectedPayment.method] || selectedPayment.method}
              </Descriptions.Item>
              <Descriptions.Item label="支付时间">
                {selectedPayment.paid_at ? new Date(selectedPayment.paid_at).toLocaleString() : '-'}
              </Descriptions.Item>
            </Descriptions>

            <Card size="small" title="费用明细">
              <Descriptions bordered column={2} size="small">
                <Descriptions.Item label="挂号费">
                  ¥{selectedPayment.registration_fee.toFixed(2)}
                </Descriptions.Item>
                <Descriptions.Item label="诊疗费">
                  ¥{selectedPayment.consultation_fee.toFixed(2)}
                </Descriptions.Item>
                <Descriptions.Item label="药品费">
                  ¥{selectedPayment.drug_fee.toFixed(2)}
                </Descriptions.Item>
                <Descriptions.Item label="检查费">
                  ¥{selectedPayment.examination_fee.toFixed(2)}
                </Descriptions.Item>
                <Descriptions.Item label="其他费用">
                  ¥{selectedPayment.other_fee.toFixed(2)}
                </Descriptions.Item>
                <Descriptions.Item label="总金额" labelStyle={{ fontWeight: 'bold' }}>
                  <Text strong type="success" className="text-lg">
                    ¥{selectedPayment.total_amount.toFixed(2)}
                  </Text>
                </Descriptions.Item>
              </Descriptions>
            </Card>

            {selectedPayment.appointment && (
              <Card size="small" title="就诊信息">
                <Descriptions bordered column={2} size="small">
                  <Descriptions.Item label="就诊日期">
                    {new Date(selectedPayment.appointment.appointment_date).toLocaleDateString()}
                  </Descriptions.Item>
                  <Descriptions.Item label="时间段">
                    {selectedPayment.appointment.start_time} - {selectedPayment.appointment.end_time}
                  </Descriptions.Item>
                  <Descriptions.Item label="患者">
                    {selectedPayment.appointment.patient?.user?.full_name || '-'}
                  </Descriptions.Item>
                  <Descriptions.Item label="医生">
                    {selectedPayment.appointment.doctor?.user?.full_name || '-'}
                  </Descriptions.Item>
                  <Descriptions.Item label="科室" span={2}>
                    {selectedPayment.appointment.doctor?.department?.name || '-'}
                  </Descriptions.Item>
                </Descriptions>
              </Card>
            )}

            {selectedPayment.notes && (
              <Card size="small" title="备注">
                <Text>{selectedPayment.notes}</Text>
              </Card>
            )}
          </div>
        )}
      </Modal>
    </div>
  )
}

export default PaymentsPage
