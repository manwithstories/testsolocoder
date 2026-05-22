import React, { useState, useEffect } from 'react'
import { Table, Card, Statistic, Row, Col, Button, Modal, Input, Select, message, Tag, Form } from 'antd'
import { billApi, withdrawApi } from '@/services/bill'
import { Bill } from '@/types'
import { formatPrice, formatDate, getBillTypeText } from '@/utils'
import { usePagination } from '@/hooks'
import { useAppSelector } from '@/store/hooks'

const BillList: React.FC = () => {
  const { userInfo } = useAppSelector((state) => state.auth)
  const { page, pageSize, total, setPage, setPageSize, setTotal } = usePagination()
  const [bills, setBills] = useState<Bill[]>([])
  const [loading, setLoading] = useState(false)
  const [summary, setSummary] = useState({ total_income: 0, total_commission: 0, net_income: 0 })
  const [withdrawModalVisible, setWithdrawModalVisible] = useState(false)
  const [withdrawForm] = Form.useForm()

  useEffect(() => {
    loadBills()
    loadSummary()
  }, [page, pageSize])

  const loadBills = async () => {
    setLoading(true)
    try {
      const res = await billApi.getList({
        page,
        page_size: pageSize,
      })
      setBills(res.list)
      setTotal(res.total)
      setSummary({
        total_income: res.income_total || 0,
        total_commission: 0,
        net_income: res.income_total || 0,
      })
    } catch (error) {
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  const loadSummary = async () => {
    try {
      const res = await billApi.getBalance()
      setSummary({
        total_income: res.total_income || 0,
        total_commission: 0,
        net_income: res.balance || 0,
      })
    } catch (error) {
      console.error(error)
    }
  }

  const handleWithdraw = async () => {
    try {
      const values = await withdrawForm.validateFields()
      await withdrawApi.create(values)
      message.success('提现申请已提交')
      setWithdrawModalVisible(false)
      withdrawForm.resetFields()
      loadBills()
    } catch (error) {
      console.error(error)
    }
  }

  const columns = [
    {
      title: '账单ID',
      dataIndex: 'id',
      key: 'id',
    },
    {
      title: '交易类型',
      dataIndex: 'bill_type',
      key: 'bill_type',
      render: (text: string) => {
        return <Tag>{getBillTypeText(text)}</Tag>
      },
    },
    {
      title: '金额',
      dataIndex: 'amount',
      key: 'amount',
      render: (text: number, record: Bill) => (
        <span style={{ color: record.bill_type === 'income' ? 'green' : record.bill_type === 'withdraw' ? 'red' : '#ff4d4f' }}>
          {record.bill_type === 'income' ? '+' : '-'}{formatPrice(text)}
        </span>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (text: string) => {
        const statusMap: Record<string, string> = {
          pending: '待结算',
          completed: '已结算',
          failed: '已冻结',
        }
        return <Tag color={text === 'completed' ? 'green' : 'orange'}>{statusMap[text] || text}</Tag>
      },
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (text: string) => formatDate(text),
    },
  ]

  return (
    <div>
      <div className="page-header">
        <h1 className="page-title">账单中心</h1>
        {userInfo?.role === 'service_provider' && (
          <Button type="primary" onClick={() => setWithdrawModalVisible(true)}>
            申请提现
          </Button>
        )}
      </div>

      <Row gutter={16} style={{ marginBottom: 24 }}>
        <Col span={8}>
          <Card>
            <Statistic
              title="总收入"
              value={summary.total_income}
              precision={2}
              valueStyle={{ color: '#3f8600' }}
              prefix="¥"
            />
          </Card>
        </Col>
        <Col span={8}>
          <Card>
            <Statistic
              title="平台佣金"
              value={summary.total_commission}
              precision={2}
              valueStyle={{ color: '#ff4d4f' }}
              prefix="¥"
            />
          </Card>
        </Col>
        <Col span={8}>
          <Card>
            <Statistic
              title="实际收入"
              value={summary.net_income}
              precision={2}
              valueStyle={{ color: '#1890ff' }}
              prefix="¥"
            />
          </Card>
        </Col>
      </Row>

      <Card>
        <Table
          rowKey="id"
          loading={loading}
          dataSource={bills}
          columns={columns}
          pagination={{
            current: page,
            pageSize,
            total,
            showSizeChanger: true,
            onChange: (p, ps) => {
              setPage(p)
              setPageSize(ps)
            },
          }}
        />
      </Card>

      <Modal
        title="申请提现"
        open={withdrawModalVisible}
        onOk={handleWithdraw}
        onCancel={() => {
          setWithdrawModalVisible(false)
          withdrawForm.resetFields()
        }}
      >
        <Form form={withdrawForm} layout="vertical">
          <Form.Item
            name="amount"
            label="提现金额"
            rules={[{ required: true, message: '请输入提现金额' }]}
          >
            <Input prefix="¥" placeholder="请输入提现金额" />
          </Form.Item>
          <Form.Item
            name="bank_name"
            label="银行名称"
            rules={[{ required: true, message: '请选择银行' }]}
          >
            <Select
              placeholder="请选择银行"
              options={[
                { label: '中国工商银行', value: '工商银行' },
                { label: '中国建设银行', value: '建设银行' },
                { label: '中国农业银行', value: '农业银行' },
                { label: '中国银行', value: '中国银行' },
                { label: '交通银行', value: '交通银行' },
                { label: '招商银行', value: '招商银行' },
                { label: '支付宝', value: '支付宝' },
                { label: '微信', value: '微信' },
              ]}
            />
          </Form.Item>
          <Form.Item
            name="bank_account"
            label="银行账号"
            rules={[{ required: true, message: '请输入银行账号' }]}
          >
            <Input placeholder="请输入银行账号" />
          </Form.Item>
          <Form.Item
            name="account_name"
            label="账户名称"
            rules={[{ required: true, message: '请输入账户名称' }]}
          >
            <Input placeholder="请输入账户名称" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default BillList
