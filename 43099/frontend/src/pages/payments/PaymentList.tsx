import { useState, useEffect } from 'react'
import { Table, Button, Tag, Space, Modal, Form, Input, Select, DatePicker, message, Card } from 'antd'
import { CheckCircleOutlined, DownloadOutlined } from '@ant-design/icons'
import { paymentApi } from '@/api'
import type { Payment } from '@/types'
import dayjs from 'dayjs'

const PaymentList = () => {
  const [payments, setPayments] = useState<Payment[]>([])
  const [loading, setLoading] = useState(false)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [confirmModalVisible, setConfirmModalVisible] = useState(false)
  const [currentPayment, setCurrentPayment] = useState<Payment | null>(null)
  const [form] = Form.useForm()
  const [exportStartDate, setExportStartDate] = useState<string>()
  const [exportEndDate, setExportEndDate] = useState<string>()

  useEffect(() => {
    loadPayments()
  }, [page, pageSize])

  const loadPayments = async () => {
    setLoading(true)
    try {
      const data: any = await paymentApi.list({ page, page_size: pageSize })
      setPayments(data.list)
      setTotal(data.total)
    } catch (error: any) {
      message.error(error.message || '加载失败')
    } finally {
      setLoading(false)
    }
  }

  const handleConfirm = (payment: Payment) => {
    setCurrentPayment(payment)
    form.resetFields()
    setConfirmModalVisible(true)
  }

  const handleSubmitConfirm = async (values: any) => {
    if (!currentPayment) return
    try {
      await paymentApi.confirm(currentPayment.order_id, values)
      message.success('支付确认成功')
      setConfirmModalVisible(false)
      loadPayments()
    } catch (error: any) {
      message.error(error.message || '确认失败')
    }
  }

  const handleExport = async () => {
    if (!exportStartDate || !exportEndDate) {
      message.error('请选择导出日期范围')
      return
    }
    try {
      const blob: any = await paymentApi.export({
        start_date: exportStartDate,
        end_date: exportEndDate,
      })
      const url = window.URL.createObjectURL(new Blob([blob]))
      const link = document.createElement('a')
      link.href = url
      link.download = `payment_records_${exportStartDate}_${exportEndDate}.xlsx`
      link.click()
      window.URL.revokeObjectURL(url)
      message.success('导出成功')
    } catch (error: any) {
      message.error(error.message || '导出失败')
    }
  }

  const getStatusTag = (status: string) => {
    const map: Record<string, { text: string; color: string }> = {
      pending: { text: '待支付', color: 'orange' },
      success: { text: '支付成功', color: 'green' },
      failed: { text: '支付失败', color: 'red' },
      refunded: { text: '已退款', color: 'default' },
    }
    const info = map[status] || { text: status, color: 'default' }
    return <Tag color={info.color}>{info.text}</Tag>
  }

  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 60,
    },
    {
      title: '订单号',
      dataIndex: ['order', 'order_no'],
      key: 'order_no',
    },
    {
      title: '交易号',
      dataIndex: 'transaction_no',
      key: 'transaction_no',
    },
    {
      title: '金额',
      dataIndex: 'amount',
      key: 'amount',
      render: (val: number) => `¥${val.toFixed(2)}`,
    },
    {
      title: '支付方式',
      dataIndex: 'payment_method',
      key: 'payment_method',
      render: (method: string) => {
        const map: Record<string, string> = {
          wechat: '微信支付',
          alipay: '支付宝',
          cash: '现金',
        }
        return map[method] || method
      },
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: getStatusTag,
    },
    {
      title: '支付时间',
      dataIndex: 'paid_at',
      key: 'paid_at',
      render: (val: string) => (val ? dayjs(val).format('YYYY-MM-DD HH:mm') : '-'),
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (val: string) => dayjs(val).format('YYYY-MM-DD HH:mm'),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Payment) => (
        <Space>
          {record.status === 'pending' && (
            <Button type="link" icon={<CheckCircleOutlined />} onClick={() => handleConfirm(record)}>
              确认支付
            </Button>
          )}
        </Space>
      ),
    },
  ]

  return (
    <Card
      title="支付管理"
      extra={
        <Space>
          <DatePicker.RangePicker
            onChange={(dates) => {
              if (dates && dates[0] && dates[1]) {
                setExportStartDate(dates[0].format('YYYY-MM-DD'))
                setExportEndDate(dates[1].format('YYYY-MM-DD'))
              }
            }}
          />
          <Button icon={<DownloadOutlined />} onClick={handleExport}>
            导出对账单
          </Button>
        </Space>
      }
    >
      <Table
        columns={columns}
        dataSource={payments}
        rowKey="id"
        loading={loading}
        pagination={{
          current: page,
          pageSize,
          total,
          onChange: (p, ps) => {
            setPage(p)
            setPageSize(ps)
          },
        }}
      />

      <Modal
        title="确认支付"
        open={confirmModalVisible}
        onCancel={() => setConfirmModalVisible(false)}
        footer={null}
      >
        <Form form={form} layout="vertical" onFinish={handleSubmitConfirm}>
          <Form.Item
            name="transaction_no"
            label="交易流水号"
            rules={[{ required: true, message: '请输入交易流水号' }]}
          >
            <Input placeholder="请输入交易流水号" />
          </Form.Item>
          <Form.Item
            name="payment_method"
            label="支付方式"
            rules={[{ required: true, message: '请选择支付方式' }]}
          >
            <Select placeholder="请选择支付方式">
              <Select.Option value="wechat">微信支付</Select.Option>
              <Select.Option value="alipay">支付宝</Select.Option>
              <Select.Option value="cash">现金</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item
            name="amount"
            label="支付金额"
            rules={[{ required: true, message: '请输入支付金额' }]}
          >
            <Input type="number" step="0.01" placeholder="请输入支付金额" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block>
              确认支付
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </Card>
  )
}

export default PaymentList
