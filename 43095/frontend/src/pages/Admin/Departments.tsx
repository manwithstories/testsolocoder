import React, { useState, useEffect } from 'react'
import {
  Card,
  Table,
  Button,
  Space,
  Typography,
  Modal,
  Form,
  Input,
  message,
  Spin,
  Empty,
  Pagination,
  Popconfirm,
  Row,
  Col,
  Descriptions,
  Tag
} from 'antd'
import {
  ApartmentOutlined,
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  CloseOutlined,
  EyeOutlined
} from '@ant-design/icons'
import type { ColumnsType } from 'antd/es/table'
import type { Department } from '@/types'
import { departmentAPI, doctorAPI } from '@/services/api'

const { Title, Text, Paragraph } = Typography
const { TextArea } = Input

interface DepartmentFormData {
  name: string
  description: string
  location: string
}

const AdminDepartmentsPage: React.FC = () => {
  const [form] = Form.useForm<DepartmentFormData>()
  const [loading, setLoading] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [departments, setDepartments] = useState<Department[]>([])
  const [doctorCounts, setDoctorCounts] = useState<Record<number, number>>({})
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize] = useState(10)
  const [modalVisible, setModalVisible] = useState(false)
  const [detailModalVisible, setDetailModalVisible] = useState(false)
  const [editingDepartment, setEditingDepartment] = useState<Department | null>(null)
  const [selectedDepartment, setSelectedDepartment] = useState<Department | null>(null)

  useEffect(() => {
    fetchDepartments()
  }, [page, pageSize])

  const fetchDepartments = async () => {
    setLoading(true)
    try {
      const result = await departmentAPI.getList({ page, pageSize })
      setDepartments(result.list)
      setTotal(result.total)
      result.list.forEach((dept) => {
        fetchDoctorCount(dept.id)
      })
    } catch (error) {
      console.error('获取科室列表失败:', error)
      message.error('获取科室列表失败')
    } finally {
      setLoading(false)
    }
  }

  const fetchDoctorCount = async (departmentId: number) => {
    try {
      const result = await doctorAPI.getList({ department_id: departmentId, page: 1, pageSize: 1 })
      setDoctorCounts((prev) => ({
        ...prev,
        [departmentId]: result.total
      }))
    } catch (error) {
      console.error('获取医生数量失败:', error)
    }
  }

  const handleAdd = () => {
    setEditingDepartment(null)
    form.resetFields()
    setModalVisible(true)
  }

  const handleEdit = (department: Department) => {
    setEditingDepartment(department)
    form.setFieldsValue({
      name: department.name,
      description: department.description,
      location: department.location
    })
    setModalVisible(true)
  }

  const handleViewDetail = (department: Department) => {
    setSelectedDepartment(department)
    setDetailModalVisible(true)
  }

  const handleDelete = async (id: number) => {
    try {
      await departmentAPI.delete(id)
      message.success('删除成功')
      fetchDepartments()
    } catch (error) {
      console.error('删除失败:', error)
      message.error('删除失败')
    }
  }

  const handleSubmit = async (values: DepartmentFormData) => {
    setSubmitting(true)
    try {
      if (editingDepartment) {
        await departmentAPI.update(editingDepartment.id, values)
        message.success('更新成功')
      } else {
        await departmentAPI.create(values)
        message.success('创建成功')
      }
      setModalVisible(false)
      fetchDepartments()
    } catch (error) {
      console.error('提交失败:', error)
      message.error('提交失败')
    } finally {
      setSubmitting(false)
    }
  }

  const handlePageChange = (newPage: number) => {
    setPage(newPage)
  }

  const columns: ColumnsType<Department> = [
    {
      title: '科室名称',
      dataIndex: 'name',
      key: 'name',
      width: 150,
      render: (text: string) => (
        <Space>
          <ApartmentOutlined className="text-blue-500" />
          <Text strong>{text}</Text>
        </Space>
      )
    },
    {
      title: '位置',
      dataIndex: 'location',
      key: 'location',
      width: 150
    },
    {
      title: '医生数量',
      key: 'doctor_count',
      width: 100,
      render: (_, record) => (
        <Tag color="blue">
          {doctorCounts[record.id] !== undefined ? doctorCounts[record.id] : '-'}
        </Tag>
      )
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
      ellipsis: true
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 160,
      render: (date: string) => new Date(date).toLocaleString()
    },
    {
      title: '操作',
      key: 'action',
      width: 200,
      fixed: 'right',
      render: (_, record) => (
        <Space size="small">
          <Button
            type="link"
            size="small"
            icon={<EyeOutlined />}
            onClick={() => handleViewDetail(record)}
          >
            详情
          </Button>
          <Button
            type="link"
            size="small"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
          >
            编辑
          </Button>
          <Popconfirm
            title="确定要删除这个科室吗？"
            description="删除科室将同时移除该科室下所有医生的关联关系"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button type="link" size="small" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      )
    }
  ]

  return (
    <div className="space-y-6">
      <Card>
        <div className="flex items-center justify-between mb-4">
          <Title level={3} style={{ margin: 0 }}>
            <ApartmentOutlined className="mr-2" />
            科室管理
          </Title>
          <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
            添加科室
          </Button>
        </div>

        <Row gutter={[16, 16]} className="mb-4">
          <Col xs={24} sm={8}>
            <Card size="small" className="text-center">
              <Text type="secondary">科室总数</Text>
              <div className="text-2xl font-bold text-blue-500 mt-1">{total}</div>
            </Card>
          </Col>
          <Col xs={24} sm={8}>
            <Card size="small" className="text-center">
              <Text type="secondary">医生总数</Text>
              <div className="text-2xl font-bold text-green-500 mt-1">
                {Object.values(doctorCounts).reduce((sum, count) => sum + count, 0)}
              </div>
            </Card>
          </Col>
          <Col xs={24} sm={8}>
            <Card size="small" className="text-center">
              <Text type="secondary">平均医生数</Text>
              <div className="text-2xl font-bold text-orange-500 mt-1">
                {total > 0
                  ? (
                      Object.values(doctorCounts).reduce((sum, count) => sum + count, 0) / total
                    ).toFixed(1)
                  : '0'}
              </div>
            </Card>
          </Col>
        </Row>

        {loading && departments.length === 0 ? (
          <div className="flex justify-center items-center min-h-[400px]">
            <Spin size="large" />
          </div>
        ) : departments.length === 0 ? (
          <Empty description="暂无科室数据" />
        ) : (
          <>
            <Table
              columns={columns}
              dataSource={departments}
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
        title={editingDepartment ? '编辑科室' : '添加科室'}
        open={modalVisible}
        onCancel={() => setModalVisible(false)}
        footer={[
          <Button key="cancel" onClick={() => setModalVisible(false)}>
            取消
          </Button>,
          <Button
            key="submit"
            type="primary"
            loading={submitting}
            onClick={() => form.submit()}
          >
            确定
          </Button>
        ]}
        width={600}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
        >
          <Form.Item
            name="name"
            label="科室名称"
            rules={[
              { required: true, message: '请输入科室名称' },
              { max: 50, message: '科室名称不能超过50个字符' }
            ]}
          >
            <Input placeholder="请输入科室名称" />
          </Form.Item>
          <Form.Item
            name="location"
            label="科室位置"
            rules={[
              { required: true, message: '请输入科室位置' },
              { max: 100, message: '科室位置不能超过100个字符' }
            ]}
          >
            <Input placeholder="请输入科室位置，如：门诊楼3层" />
          </Form.Item>
          <Form.Item
            name="description"
            label="科室描述"
            rules={[
              { required: true, message: '请输入科室描述' },
              { max: 500, message: '科室描述不能超过500个字符' }
            ]}
          >
            <TextArea rows={4} placeholder="请输入科室描述，包括诊疗范围、特色专科等" />
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title="科室详情"
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
        {selectedDepartment && (
          <div className="space-y-4">
            <Descriptions bordered column={2} size="small">
              <Descriptions.Item label="科室名称">
                {selectedDepartment.name}
              </Descriptions.Item>
              <Descriptions.Item label="科室位置">
                {selectedDepartment.location}
              </Descriptions.Item>
              <Descriptions.Item label="医生数量">
                <Tag color="blue">
                  {doctorCounts[selectedDepartment.id] !== undefined
                    ? doctorCounts[selectedDepartment.id]
                    : 0}
                </Tag>
              </Descriptions.Item>
              <Descriptions.Item label="创建时间">
                {new Date(selectedDepartment.created_at).toLocaleString()}
              </Descriptions.Item>
              <Descriptions.Item label="更新时间" span={2}>
                {new Date(selectedDepartment.updated_at).toLocaleString()}
              </Descriptions.Item>
            </Descriptions>

            <Card size="small" title="科室描述">
              <Paragraph className="mb-0">
                {selectedDepartment.description}
              </Paragraph>
            </Card>
          </div>
        )}
      </Modal>
    </div>
  )
}

export default AdminDepartmentsPage
