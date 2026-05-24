import React, { useEffect, useState } from 'react'
import { Table, Tag, Card, Button, Modal, Form, Input, Select, DatePicker, Upload, message, Space, Empty, Spin, Tabs, Row, Col } from 'antd'
import { PlusOutlined, UploadOutlined, DeleteOutlined, FileTextOutlined } from '@ant-design/icons'
import { useAuth } from '../contexts/AuthContext'
import { listHealthRecords, createHealthRecord, deleteHealthRecord, uploadHealthReport } from '../api/health'
import { HealthRecord, CreateHealthRecordRequest } from '../types'
import dayjs from 'dayjs'

const { Option } = Select
const { TextArea } = Input

const HealthRecords: React.FC = () => {
  const { user } = useAuth()
  const [records, setRecords] = useState<HealthRecord[]>([])
  const [loading, setLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [form] = Form.useForm()
  const [submitting, setSubmitting] = useState(false)
  const [typeFilter, setTypeFilter] = useState<string>('')

  const isRescue = user?.role === 'rescue' || user?.role === 'admin'

  useEffect(() => {
    loadRecords()
  }, [typeFilter])

  const loadRecords = async () => {
    setLoading(true)
    try {
      const response = await listHealthRecords({
        record_type: typeFilter || undefined,
        page_size: 100,
      })
      if (response.code === 0 && response.data) {
        setRecords((response.data as any).items || [])
      }
    } catch (error) {
      console.error('Failed to load health records:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      setSubmitting(true)

      const data: CreateHealthRecordRequest = {
        pet_id: values.pet_id,
        record_type: values.record_type,
        title: values.title,
        description: values.description,
        vaccine_name: values.vaccine_name,
        record_date: values.record_date.format('YYYY-MM-DD'),
        next_date: values.next_date?.format('YYYY-MM-DD'),
        weight: values.weight,
        temperature: values.temperature,
        vet_name: values.vet_name,
        hospital: values.hospital,
        notes: values.notes,
      }

      await createHealthRecord(data)
      message.success('记录添加成功')
      setModalVisible(false)
      form.resetFields()
      loadRecords()
    } catch (error: any) {
      message.error(error.message || '添加失败')
    } finally {
      setSubmitting(false)
    }
  }

  const handleDelete = async (id: number) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这条记录吗？',
      onOk: async () => {
        try {
          await deleteHealthRecord(id)
          message.success('已删除')
          loadRecords()
        } catch (error: any) {
          message.error('删除失败')
        }
      },
    })
  }

  const handleUploadReport = async (recordId: number, file: File) => {
    const formData = new FormData()
    formData.append('file', file)
    try {
      await uploadHealthReport(recordId, formData)
      message.success('报告上传成功')
      loadRecords()
    } catch (error: any) {
      message.error('上传失败')
    }
  }

  const columns = [
    { title: '编号', dataIndex: 'id', key: 'id' },
    {
      title: '宠物',
      key: 'pet',
      render: (_: any, r: HealthRecord) => r.pet?.name || '-',
    },
    { title: '类型', dataIndex: 'record_type', key: 'type', render: (t: string) => <Tag>{t}</Tag> },
    { title: '标题', dataIndex: 'title', key: 'title' },
    { title: '疫苗', dataIndex: 'vaccine_name', key: 'vaccine_name', render: (v: string) => v || '-' },
    {
      title: '记录日期',
      dataIndex: 'record_date',
      key: 'record_date',
      render: (date: string) => date ? new Date(date).toLocaleDateString() : '-',
    },
    {
      title: '下次日期',
      dataIndex: 'next_date',
      key: 'next_date',
      render: (date: string) => date ? new Date(date).toLocaleDateString() : '-',
    },
    { title: '兽医', dataIndex: 'vet_name', key: 'vet_name', render: (v: string) => v || '-' },
    {
      title: '报告',
      key: 'report',
      render: (_: any, r: HealthRecord) => (
        r.report_file ? (
          <a href={r.report_file} target="_blank" rel="noreferrer">查看</a>
        ) : isRescue ? (
          <Upload
            showUploadList={false}
            beforeUpload={(file) => {
              handleUploadReport(r.id, file)
              return false
            }}
            accept=".pdf"
          >
            <Button type="link" size="small" icon={<UploadOutlined />}>上传</Button>
          </Upload>
        ) : (
          '-'
        )
      ),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: HealthRecord) => (
        isRescue && (
          <Button type="link" size="small" danger icon={<DeleteOutlined />} onClick={() => handleDelete(record.id)}>
            删除
          </Button>
        )
      ),
    },
  ]

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h2 style={{ margin: 0 }}>健康档案</h2>
        <Space>
          <Select
            placeholder="筛选类型"
            allowClear
            style={{ width: 150 }}
            value={typeFilter || undefined}
            onChange={setTypeFilter}
          >
            <Option value="vaccine">疫苗</Option>
            <Option value="deworm">驱虫</Option>
            <Option value="checkup">体检</Option>
            <Option value="disease">病史</Option>
            <Option value="surgery">手术</Option>
            <Option value="other">其他</Option>
          </Select>
          {isRescue && (
            <Button type="primary" icon={<PlusOutlined />} onClick={() => setModalVisible(true)}>
              添加记录
            </Button>
          )}
        </Space>
      </div>

      <Spin spinning={loading}>
        {records.length > 0 ? (
          <Card>
            <Table
              dataSource={records}
              columns={columns}
              rowKey="id"
              pagination={{ pageSize: 10 }}
            />
          </Card>
        ) : (
          !loading && (
            <Card>
              <Empty description="暂无健康记录" />
            </Card>
          )
        )}
      </Spin>

      <Modal
        title="添加健康记录"
        open={modalVisible}
        onCancel={() => {
          setModalVisible(false)
          form.resetFields()
        }}
        footer={null}
        width={600}
      >
        <Form form={form} onFinish={handleSubmit} layout="vertical">
          <Form.Item name="pet_id" label="宠物ID" rules={[{ required: true }]}>
            <Input type="number" placeholder="请输入宠物ID" />
          </Form.Item>
          <Form.Item name="record_type" label="记录类型" rules={[{ required: true }]}>
            <Select>
              <Option value="vaccine">疫苗接种</Option>
              <Option value="deworm">驱虫</Option>
              <Option value="checkup">体检</Option>
              <Option value="disease">病史</Option>
              <Option value="surgery">手术</Option>
              <Option value="other">其他</Option>
            </Select>
          </Form.Item>
          <Form.Item name="title" label="标题" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <TextArea rows={3} />
          </Form.Item>
          <Form.Item name="vaccine_name" label="疫苗名称">
            <Input />
          </Form.Item>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="record_date" label="记录日期" rules={[{ required: true }]}>
                <DatePicker style={{ width: '100%' }} />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="next_date" label="下次日期">
                <DatePicker style={{ width: '100%' }} />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="weight" label="体重(kg)">
                <Input type="number" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="temperature" label="体温">
                <Input type="number" />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="vet_name" label="兽医">
                <Input />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="hospital" label="医院">
                <Input />
              </Form.Item>
            </Col>
          </Row>
          <Form.Item name="notes" label="备注">
            <TextArea rows={2} />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={submitting} block>
              保存
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default HealthRecords
