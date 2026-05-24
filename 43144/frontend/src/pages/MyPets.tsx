import React, { useEffect, useState } from 'react'
import { Row, Col, Button, Modal, Form, Input, Select, message, Table, Space, Tag } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { getMyPets, createPet, deletePet, updatePetStatus } from '../api/pet'
import { Pet, CreatePetRequest } from '../types'

const { Option } = Select
const { TextArea } = Input

const MyPets: React.FC = () => {
  const navigate = useNavigate()
  const [pets, setPets] = useState<Pet[]>([])
  const [loading, setLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [editingPet, setEditingPet] = useState<Pet | null>(null)
  const [form] = Form.useForm()
  const [submitting, setSubmitting] = useState(false)

  useEffect(() => {
    loadPets()
  }, [])

  const loadPets = async () => {
    setLoading(true)
    try {
      const response = await getMyPets({ page_size: 100 })
      if (response.code === 0 && response.data) {
        setPets((response.data as any).items || [])
      }
    } catch (error) {
      console.error('Failed to load pets:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      setSubmitting(true)

      if (editingPet) {
        message.success('请前往详情页编辑')
      } else {
        await createPet(values as CreatePetRequest)
        message.success('宠物添加成功')
      }

      setModalVisible(false)
      form.resetFields()
      setEditingPet(null)
      loadPets()
    } catch (error: any) {
      message.error(error.message || '操作失败')
    } finally {
      setSubmitting(false)
    }
  }

  const handleDelete = async (id: number) => {
    try {
      await deletePet(id)
      message.success('删除成功')
      loadPets()
    } catch (error: any) {
      message.error('删除失败')
    }
  }

  const handleStatusChange = async (id: number, status: string) => {
    try {
      await updatePetStatus(id, status)
      message.success('状态更新成功')
      loadPets()
    } catch (error: any) {
      message.error('更新失败')
    }
  }

  const columns = [
    { title: '编号', dataIndex: 'archive_number', key: 'archive_number' },
    { title: '名称', dataIndex: 'name', key: 'name' },
    { title: '物种', dataIndex: 'species', key: 'species' },
    { title: '品种', dataIndex: 'breed', key: 'breed' },
    { title: '性别', dataIndex: 'gender', key: 'gender', render: (g: string) => g === 'male' ? '公' : g === 'female' ? '母' : '未知' },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string, record: Pet) => (
        <Select
          value={status}
          size="small"
          style={{ width: 100 }}
          onChange={(val) => handleStatusChange(record.id, val)}
        >
          <Option value="adoptable">待领养</Option>
          <Option value="adopted">已领养</Option>
          <Option value="treatment">治疗中</Option>
          <Option value="deceased">已去世</Option>
        </Select>
      ),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Pet) => (
        <Space>
          <Button type="link" size="small" onClick={() => navigate(`/pets/${record.id}`)}>
            查看
          </Button>
          <Button type="link" size="small" danger onClick={() => handleDelete(record.id)}>
            删除
          </Button>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <h2 style={{ margin: 0 }}>宠物管理</h2>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => setModalVisible(true)}>
          添加宠物
        </Button>
      </div>

      <Table
        dataSource={pets}
        columns={columns}
        rowKey="id"
        loading={loading}
        pagination={{ pageSize: 10 }}
      />

      <Modal
        title={editingPet ? '编辑宠物' : '添加宠物'}
        open={modalVisible}
        onCancel={() => {
          setModalVisible(false)
          setEditingPet(null)
          form.resetFields()
        }}
        footer={null}
        width={600}
      >
        <Form form={form} onFinish={handleSubmit} layout="vertical">
          <Form.Item name="name" label="名称" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="species" label="物种" rules={[{ required: true }]}>
                <Select>
                  <Option value="dog">狗</Option>
                  <Option value="cat">猫</Option>
                  <Option value="rabbit">兔子</Option>
                  <Option value="bird">鸟</Option>
                  <Option value="other">其他</Option>
                </Select>
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="breed" label="品种">
                <Input />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={16}>
            <Col span={8}>
              <Form.Item name="age" label="年龄">
                <Input placeholder="如: 2岁" />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item name="gender" label="性别" rules={[{ required: true }]}>
                <Select>
                  <Option value="male">公</Option>
                  <Option value="female">母</Option>
                  <Option value="unknown">未知</Option>
                </Select>
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item name="weight" label="体重(kg)">
                <Input type="number" />
              </Form.Item>
            </Col>
          </Row>
          <Form.Item name="description" label="描述">
            <TextArea rows={3} />
          </Form.Item>
          <Form.Item name="personality" label="性格特点">
            <TextArea rows={2} />
          </Form.Item>
          <Form.Item name="special_needs" label="特殊需求">
            <TextArea rows={2} />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={submitting} block>
              {editingPet ? '保存' : '添加'}
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default MyPets
