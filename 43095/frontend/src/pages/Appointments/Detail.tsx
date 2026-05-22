import React, { useState, useEffect } from 'react'
import {
  Card,
  Descriptions,
  Button,
  Spin,
  Alert,
  Tag,
  Avatar,
  Space,
  Row,
  Col,
  Statistic,
  Divider,
  List,
  Empty,
  Table,
  Rate,
  Typography,
  Modal,
  Form,
  Input,
  message
} from 'antd'
import {
  ArrowLeftOutlined,
  UserOutlined,
  FileTextOutlined,
  MedicineBoxOutlined,
  FileSearchOutlined,
  CreditCardOutlined,
  StarOutlined,
  EditOutlined
} from '@ant-design/icons'
import { useParams, useNavigate } from 'react-router-dom'
import dayjs from 'dayjs'
import { appointmentAPI, reviewAPI } from '@/services/api'
import type { Appointment, AppointmentStatus, PaymentStatus } from '@/types'

const { Title, Text, Paragraph } = Typography
const { TextArea } = Input

const statusMap: Record<AppointmentStatus, { text: string; color: string }> = {
  pending: { text: '待确认', color: 'orange' },
  confirmed: { text: '已确认', color: 'blue' },
  completed: { text: '已完成', color: 'green' },
  cancelled: { text: '已取消', color: 'red' },
  no_show: { text: '未就诊', color: 'default' }
}

const paymentStatusMap: Record<PaymentStatus, { text: string; color: string }> = {
  pending: { text: '待支付', color: 'orange' },
  paid: { text: '已支付', color: 'green' },
  failed: { text: '支付失败', color: 'red' },
  refunded: { text: '已退款', color: 'default' }
}

const AppointmentDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [appointment, setAppointment] = useState<Appointment | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const [reviewModalVisible, setReviewModalVisible] = useState(false)
  const [reviewForm] = Form.useForm()
  const [reviewLoading, setReviewLoading] = useState(false)

  useEffect(() => {
    if (id) {
      fetchAppointmentDetail(parseInt(id))
    }
  }, [id])

  const fetchAppointmentDetail = async (appointmentId: number) => {
    setLoading(true)
    setError(null)
    try {
      const response = await appointmentAPI.getDetail(appointmentId)
      setAppointment(response)
    } catch (err: any) {
      setError(err.message || '获取预约详情失败')
    } finally {
      setLoading(false)
    }
  }

  const handleBack = () => {
    navigate('/appointments')
  }

  const handleEdit = () => {
    navigate(`/appointments/create?appointment_id=${id}`)
  }

  const handleReview = () => {
    if (appointment?.review) {
      message.info('您已经评价过此预约')
      return
    }
    setReviewModalVisible(true)
    reviewForm.resetFields()
  }

  const handleReviewSubmit = async () => {
    if (!id) return
    setReviewLoading(true)
    try {
      const values = await reviewForm.validateFields()
      await reviewAPI.create({
        appointment_id: parseInt(id),
        rating: values.rating,
        content: values.content,
        is_anonymous: values.is_anonymous
      })
      message.success('评价提交成功')
      setReviewModalVisible(false)
      fetchAppointmentDetail(parseInt(id))
    } catch (err: any) {
      message.error(err.message || '评价提交失败')
    } finally {
      setReviewLoading(false)
    }
  }

  const prescriptionColumns = [
    {
      title: '药品名称',
      dataIndex: 'drug_name',
      key: 'drug_name'
    },
    {
      title: '规格',
      dataIndex: 'specification',
      key: 'specification'
    },
    {
      title: '用法用量',
      key: 'dosage',
      render: (_: unknown, record: any) => `${record.dosage} ${record.frequency}`
    },
    {
      title: '数量',
      dataIndex: 'quantity',
      key: 'quantity'
    },
    {
      title: '单价',
      dataIndex: 'unit_price',
      key: 'unit_price',
      render: (price: number) => `¥${price.toFixed(2)}`
    },
    {
      title: '小计',
      dataIndex: 'subtotal',
      key: 'subtotal',
      render: (price: number) => `¥${price.toFixed(2)}`
    }
  ]

  return (
    <div className="p-6 max-w-6xl mx-auto">
      <div className="flex justify-between items-center mb-4">
        <Button
          icon={<ArrowLeftOutlined />}
          onClick={handleBack}
          type="text"
        >
          返回预约列表
        </Button>
        {appointment && (appointment.status === 'pending' || appointment.status === 'confirmed') && (
          <Button
            type="primary"
            icon={<EditOutlined />}
            onClick={handleEdit}
          >
            改签预约
          </Button>
        )}
      </div>

      {error && (
        <Alert
          message="错误"
          description={error}
          type="error"
          showIcon
          className="mb-4"
          closable
          onClose={() => setError(null)}
        />
      )}

      <Spin spinning={loading} tip="加载中...">
        {appointment && (
          <>
            <Card className="mb-6">
              <Row gutter={24} align="middle">
                <Col flex="auto">
                  <Space size="large" align="center" className="mb-4">
                    <Title level={3} className="mb-0">
                      预约详情
                    </Title>
                    <Tag color={statusMap[appointment.status].color} className="text-base py-1 px-3">
                      {statusMap[appointment.status].text}
                    </Tag>
                    <Text type="secondary">
                      预约编号：{appointment.id}
                    </Text>
                  </Space>
                  <Row gutter={32}>
                    <Col>
                      <Statistic
                        title="预约日期"
                        value={dayjs(appointment.appointment_date).format('YYYY-MM-DD')}
                      />
                    </Col>
                    <Col>
                      <Statistic
                        title="就诊时间"
                        value={`${appointment.start_time} - ${appointment.end_time}`}
                      />
                    </Col>
                    <Col>
                      <Statistic
                        title="创建时间"
                        value={dayjs(appointment.created_at).format('YYYY-MM-DD HH:mm')}
                      />
                    </Col>
                  </Row>
                </Col>
              </Row>
            </Card>

            <Card
              title={
                <Space>
                  <UserOutlined />
                  患者信息
                </Space>
              }
              className="mb-6"
            >
              <Descriptions column={2} bordered>
                <Descriptions.Item label="姓名">
                  {appointment.patient?.user?.full_name || '-'}
                </Descriptions.Item>
                <Descriptions.Item label="性别">
                  {appointment.patient?.user?.gender || '-'}
                </Descriptions.Item>
                <Descriptions.Item label="手机号">
                  {appointment.patient?.user?.phone || '-'}
                </Descriptions.Item>
                <Descriptions.Item label="邮箱">
                  {appointment.patient?.user?.email || '-'}
                </Descriptions.Item>
                <Descriptions.Item label="身份证号">
                  {appointment.patient?.id_card_no || '-'}
                </Descriptions.Item>
                <Descriptions.Item label="紧急联系人">
                  {appointment.patient?.emergency_contact_name || '-'}
                  {appointment.patient?.emergency_contact_phone &&
                    ` (${appointment.patient.emergency_contact_phone})`}
                </Descriptions.Item>
              </Descriptions>
            </Card>

            <Card
              title={
                <Space>
                  <UserOutlined />
                  医生信息
                </Space>
              }
              className="mb-6"
            >
              <Row gutter={24}>
                <Col flex="100px">
                  <Avatar
                    size={80}
                    src={appointment.doctor?.user?.avatar_url}
                    icon={!appointment.doctor?.user?.avatar_url && <UserOutlined />}
                  />
                </Col>
                <Col flex="auto">
                  <Space size="large" align="center" className="mb-2">
                    <span className="text-xl font-semibold">
                      {appointment.doctor?.user?.full_name}
                    </span>
                    <Tag color="blue">{appointment.doctor?.title}</Tag>
                    <Tag color="green">{appointment.doctor?.department?.name}</Tag>
                  </Space>
                  <p className="text-gray-500 mb-2">
                    擅长：{appointment.doctor?.specialty || '暂无信息'}
                  </p>
                  <p className="text-gray-600">
                    {appointment.doctor?.introduction || '暂无介绍'}
                  </p>
                </Col>
              </Row>
            </Card>

            {appointment.symptoms && (
              <Card
                title={
                  <Space>
                    <FileTextOutlined />
                    症状描述
                  </Space>
                }
                className="mb-6"
              >
                <Paragraph className="mb-0">
                  {appointment.symptoms}
                </Paragraph>
              </Card>
            )}

            {appointment.consultation && (
              <Card
                title={
                  <Space>
                    <FileSearchOutlined />
                    问诊记录
                  </Space>
                }
                className="mb-6"
              >
                <Descriptions column={1} bordered>
                  <Descriptions.Item label="诊断结果">
                    {appointment.consultation.diagnosis || '-'}
                  </Descriptions.Item>
                  <Descriptions.Item label="治疗方案">
                    {appointment.consultation.treatment_plan || '-'}
                  </Descriptions.Item>
                  <Descriptions.Item label="医生建议">
                    {appointment.consultation.doctor_notes || '-'}
                  </Descriptions.Item>
                  <Descriptions.Item label="复诊日期">
                    {appointment.consultation.follow_up_date
                      ? dayjs(appointment.consultation.follow_up_date).format('YYYY-MM-DD')
                      : '-'}
                  </Descriptions.Item>
                </Descriptions>
              </Card>
            )}

            {appointment.consultation?.prescription && (
              <Card
                title={
                  <Space>
                    <MedicineBoxOutlined />
                    处方信息
                  </Space>
                }
                className="mb-6"
                extra={
                  <Tag color={appointment.consultation.prescription.is_fulfilled ? 'green' : 'orange'}>
                    {appointment.consultation.prescription.is_fulfilled ? '已取药' : '待取药'}
                  </Tag>
                }
              >
                <div className="mb-4">
                  <Text type="secondary">处方编号：</Text>
                  <Text strong>{appointment.consultation.prescription.prescription_no}</Text>
                </div>
                {appointment.consultation.prescription.items &&
                appointment.consultation.prescription.items.length > 0 ? (
                  <>
                    <Table
                      columns={prescriptionColumns}
                      dataSource={appointment.consultation.prescription.items}
                      rowKey="id"
                      pagination={false}
                      size="small"
                    />
                    <div className="mt-4 text-right">
                      <Text strong className="text-lg">
                        总计：¥
                        {appointment.consultation.prescription.items
                          .reduce((sum, item) => sum + item.subtotal, 0)
                          .toFixed(2)}
                      </Text>
                    </div>
                  </>
                ) : (
                  <Empty description="暂无处方药品" image={Empty.PRESENTED_IMAGE_SIMPLE} />
                )}
                {appointment.consultation.prescription.notes && (
                  <div className="mt-4">
                    <Text type="secondary">医嘱：</Text>
                    <Text>{appointment.consultation.prescription.notes}</Text>
                  </div>
                )}
              </Card>
            )}

            {appointment.consultation?.reports && appointment.consultation.reports.length > 0 && (
              <Card
                title={
                  <Space>
                    <FileSearchOutlined />
                    检查报告
                  </Space>
                }
                className="mb-6"
              >
                <List
                  dataSource={appointment.consultation.reports}
                  renderItem={(report) => (
                    <List.Item key={report.id}>
                      <List.Item.Meta
                        avatar={<FileSearchOutlined className="text-2xl text-blue-500" />}
                        title={
                          <Space>
                            <span className="font-medium">{report.report_name}</span>
                            <Tag>{report.report_type}</Tag>
                          </Space>
                        }
                        description={
                          <div>
                            <p className="text-gray-600 mb-1">
                              检查结论：{report.conclusion || '暂无'}
                            </p>
                            <Text type="secondary" className="text-xs">
                              上传时间：{dayjs(report.created_at).format('YYYY-MM-DD HH:mm')}
                            </Text>
                          </div>
                        }
                      />
                      <Button
                        type="link"
                        onClick={() => window.open(report.file_url, '_blank')}
                      >
                        查看报告
                      </Button>
                    </List.Item>
                  )}
                />
              </Card>
            )}

            {appointment.payment && (
              <Card
                title={
                  <Space>
                    <CreditCardOutlined />
                    支付信息
                  </Space>
                }
                className="mb-6"
                extra={
                  <Tag color={paymentStatusMap[appointment.payment.status].color}>
                    {paymentStatusMap[appointment.payment.status].text}
                  </Tag>
                }
              >
                <Row gutter={16}>
                  <Col span={8}>
                    <Statistic title="挂号费" value={appointment.payment.registration_fee} prefix="¥" />
                  </Col>
                  <Col span={8}>
                    <Statistic title="问诊费" value={appointment.payment.consultation_fee} prefix="¥" />
                  </Col>
                  <Col span={8}>
                    <Statistic title="药品费" value={appointment.payment.drug_fee} prefix="¥" />
                  </Col>
                  <Col span={8}>
                    <Statistic title="检查费" value={appointment.payment.examination_fee} prefix="¥" />
                  </Col>
                  <Col span={8}>
                    <Statistic title="其他费用" value={appointment.payment.other_fee} prefix="¥" />
                  </Col>
                  <Col span={8}>
                    <Statistic
                      title="总计金额"
                      value={appointment.payment.total_amount}
                      prefix="¥"
                      valueStyle={{ color: '#cf1322' }}
                    />
                  </Col>
                </Row>
                <Divider />
                <Descriptions column={2} size="small">
                  <Descriptions.Item label="交易单号">
                    {appointment.payment.transaction_no || '-'}
                  </Descriptions.Item>
                  <Descriptions.Item label="支付方式">
                    {appointment.payment.method || '-'}
                  </Descriptions.Item>
                  <Descriptions.Item label="支付时间">
                    {appointment.payment.paid_at
                      ? dayjs(appointment.payment.paid_at).format('YYYY-MM-DD HH:mm:ss')
                      : '-'}
                  </Descriptions.Item>
                  <Descriptions.Item label="退款时间">
                    {appointment.payment.refunded_at
                      ? dayjs(appointment.payment.refunded_at).format('YYYY-MM-DD HH:mm:ss')
                      : '-'}
                  </Descriptions.Item>
                </Descriptions>
              </Card>
            )}

            <Card
              title={
                <Space>
                  <StarOutlined />
                  患者评价
                </Space>
              }
              className="mb-6"
              extra={
                appointment.status === 'completed' && !appointment.review && (
                  <Button type="primary" size="small" onClick={handleReview}>
                    去评价
                  </Button>
                )
              }
            >
              {appointment.review ? (
                <div>
                  <Space className="mb-2">
                    <Rate disabled value={appointment.review.rating} allowHalf />
                    <Text type="secondary">
                      {dayjs(appointment.review.created_at).format('YYYY-MM-DD HH:mm')}
                    </Text>
                    {appointment.review.is_anonymous && <Tag>匿名</Tag>}
                  </Space>
                  <Paragraph className="mb-0">
                    {appointment.review.content || '暂无评价内容'}
                  </Paragraph>
                </div>
              ) : (
                <Empty
                  description={
                    appointment.status === 'completed'
                      ? '您还没有评价此次就诊'
                      : '就诊完成后可进行评价'
                  }
                  image={Empty.PRESENTED_IMAGE_SIMPLE}
                />
              )}
            </Card>

            {appointment.cancel_reason && (
              <Card
                title="取消原因"
                className="mb-6"
                size="small"
                style={{ borderColor: '#ffccc7', background: '#fff2f0' }}
              >
                <Paragraph className="mb-0">{appointment.cancel_reason}</Paragraph>
              </Card>
            )}
          </>
        )}
      </Spin>

      <Modal
        title="评价医生"
        open={reviewModalVisible}
        onCancel={() => setReviewModalVisible(false)}
        footer={[
          <Button key="back" onClick={() => setReviewModalVisible(false)}>
            取消
          </Button>,
          <Button
            key="submit"
            type="primary"
            loading={reviewLoading}
            onClick={handleReviewSubmit}
          >
            提交评价
          </Button>
        ]}
      >
        <Form form={reviewForm} layout="vertical">
          <Form.Item
            name="rating"
            label="评分"
            rules={[{ required: true, message: '请选择评分' }]}
          >
            <Rate />
          </Form.Item>
          <Form.Item name="content" label="评价内容">
            <TextArea rows={4} placeholder="请分享您的就诊体验（选填）" maxLength={500} showCount />
          </Form.Item>
          <Form.Item name="is_anonymous" label="匿名评价" valuePropName="checked">
            <Input type="checkbox" />
            <span className="ml-2 text-gray-500">匿名评价将不显示您的个人信息</span>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default AppointmentDetail
