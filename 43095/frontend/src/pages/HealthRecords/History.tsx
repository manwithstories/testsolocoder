import React, { useState, useEffect } from 'react'
import {
  Card,
  Table,
  Button,
  Typography,
  Tag,
  Modal,
  Descriptions,
  message,
  Spin,
  Empty,
  Pagination
} from 'antd'
import {
  HistoryOutlined,
  EyeOutlined,
  FilePdfOutlined,
  CloseOutlined
} from '@ant-design/icons'
import type { ColumnsType } from 'antd/es/table'
import type { Appointment, User } from '@/types'
import { healthRecordAPI, authAPI } from '@/services/api'
import { saveAs } from 'file-saver'

const { Title, Text } = Typography

const statusMap: Record<string, { text: string; color: string }> = {
  pending: { text: '待确认', color: 'orange' },
  confirmed: { text: '已确认', color: 'blue' },
  completed: { text: '已完成', color: 'green' },
  cancelled: { text: '已取消', color: 'red' },
  no_show: { text: '未就诊', color: 'default' }
}

const VisitHistoryPage: React.FC = () => {
  const [loading, setLoading] = useState(false)
  const [exporting, setExporting] = useState(false)
  const [appointments, setAppointments] = useState<Appointment[]>([])
  const [currentUser, setCurrentUser] = useState<User | null>(null)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize] = useState(10)
  const [detailModalVisible, setDetailModalVisible] = useState(false)
  const [selectedAppointment, setSelectedAppointment] = useState<Appointment | null>(null)

  useEffect(() => {
    fetchCurrentUser()
  }, [])

  useEffect(() => {
    if (currentUser) {
      fetchVisitHistory()
    }
  }, [currentUser, page, pageSize])

  const fetchCurrentUser = async () => {
    try {
      const user = await authAPI.getCurrentUser()
      setCurrentUser(user)
    } catch (error) {
      console.error('获取用户信息失败:', error)
    }
  }

  const fetchVisitHistory = async () => {
    if (!currentUser) return

    setLoading(true)
    try {
      const result = await healthRecordAPI.getVisitHistory(currentUser.id, {
        page,
        pageSize
      })
      setAppointments(result.list)
      setTotal(result.total)
    } catch (error) {
      console.error('获取就诊历史失败:', error)
      message.error('获取就诊历史失败')
    } finally {
      setLoading(false)
    }
  }

  const handleViewDetail = (appointment: Appointment) => {
    setSelectedAppointment(appointment)
    setDetailModalVisible(true)
  }

  const handleExportPDF = async () => {
    if (!currentUser) return

    setExporting(true)
    try {
      const blob = await healthRecordAPI.export(currentUser.id)
      saveAs(blob, `就诊记录_${currentUser.full_name}_${new Date().toLocaleDateString()}.pdf`)
      message.success('导出成功')
    } catch (error) {
      console.error('导出PDF失败:', error)
      message.error('导出PDF失败')
    } finally {
      setExporting(false)
    }
  }

  const handlePageChange = (newPage: number) => {
    setPage(newPage)
  }

  const columns: ColumnsType<Appointment> = [
    {
      title: '就诊日期',
      dataIndex: 'appointment_date',
      key: 'appointment_date',
      width: 120,
      render: (date: string) => new Date(date).toLocaleDateString()
    },
    {
      title: '时间段',
      key: 'time',
      width: 120,
      render: (_, record) => `${record.start_time} - ${record.end_time}`
    },
    {
      title: '医生',
      key: 'doctor',
      width: 150,
      render: (_, record) => record.doctor?.user?.full_name || '-'
    },
    {
      title: '科室',
      key: 'department',
      width: 120,
      render: (_, record) => record.doctor?.department?.name || '-'
    },
    {
      title: '职称',
      key: 'title',
      width: 120,
      render: (_, record) => record.doctor?.title || '-'
    },
    {
      title: '症状描述',
      dataIndex: 'symptoms',
      key: 'symptoms',
      ellipsis: true,
      render: (text: string) => text || '-'
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
            <HistoryOutlined className="mr-2" />
            就诊历史
          </Title>
          <Button
            type="primary"
            icon={<FilePdfOutlined />}
            loading={exporting}
            onClick={handleExportPDF}
            disabled={appointments.length === 0}
          >
            导出PDF
          </Button>
        </div>

        {loading && appointments.length === 0 ? (
          <div className="flex justify-center items-center min-h-[400px]">
            <Spin size="large" />
          </div>
        ) : appointments.length === 0 ? (
          <Empty description="暂无就诊记录" />
        ) : (
          <>
            <Table
              columns={columns}
              dataSource={appointments}
              rowKey="id"
              pagination={false}
              loading={loading}
              scroll={{ x: 1000 }}
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
        title="就诊详情"
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
        {selectedAppointment && (
          <div className="space-y-4">
            <Descriptions bordered column={2} size="small">
              <Descriptions.Item label="就诊日期">
                {new Date(selectedAppointment.appointment_date).toLocaleDateString()}
              </Descriptions.Item>
              <Descriptions.Item label="时间段">
                {selectedAppointment.start_time} - {selectedAppointment.end_time}
              </Descriptions.Item>
              <Descriptions.Item label="医生">
                {selectedAppointment.doctor?.user?.full_name || '-'}
              </Descriptions.Item>
              <Descriptions.Item label="科室">
                {selectedAppointment.doctor?.department?.name || '-'}
              </Descriptions.Item>
              <Descriptions.Item label="职称">
                {selectedAppointment.doctor?.title || '-'}
              </Descriptions.Item>
              <Descriptions.Item label="状态">
                {statusMap[selectedAppointment.status] && (
                  <Tag color={statusMap[selectedAppointment.status].color}>
                    {statusMap[selectedAppointment.status].text}
                  </Tag>
                )}
              </Descriptions.Item>
              <Descriptions.Item label="挂号费" span={2}>
                ¥{selectedAppointment.doctor?.registration_fee || 0}
              </Descriptions.Item>
            </Descriptions>

            <Card size="small" title="症状描述">
              <Text>{selectedAppointment.symptoms || '无'}</Text>
            </Card>

            {selectedAppointment.notes && (
              <Card size="small" title="医生备注">
                <Text>{selectedAppointment.notes}</Text>
              </Card>
            )}

            {selectedAppointment.cancel_reason && (
              <Card size="small" title="取消原因">
                <Text type="danger">{selectedAppointment.cancel_reason}</Text>
              </Card>
            )}

            {selectedAppointment.consultation && (
              <Card size="small" title="诊断信息">
                <Descriptions bordered column={1} size="small">
                  <Descriptions.Item label="诊断结果">
                    {selectedAppointment.consultation.diagnosis || '-'}
                  </Descriptions.Item>
                  <Descriptions.Item label="治疗方案">
                    {selectedAppointment.consultation.treatment_plan || '-'}
                  </Descriptions.Item>
                  <Descriptions.Item label="医生建议">
                    {selectedAppointment.consultation.doctor_notes || '-'}
                  </Descriptions.Item>
                </Descriptions>
              </Card>
            )}

            {selectedAppointment.payment && (
              <Card size="small" title="费用信息">
                <Descriptions bordered column={2} size="small">
                  <Descriptions.Item label="总金额">
                    ¥{selectedAppointment.payment.total_amount.toFixed(2)}
                  </Descriptions.Item>
                  <Descriptions.Item label="支付状态">
                    <Tag color={selectedAppointment.payment.status === 'paid' ? 'green' : 'orange'}>
                      {selectedAppointment.payment.status === 'paid' ? '已支付' : '待支付'}
                    </Tag>
                  </Descriptions.Item>
                </Descriptions>
              </Card>
            )}
          </div>
        )}
      </Modal>
    </div>
  )
}

export default VisitHistoryPage
