import { useState } from 'react'
import { Table, Button, Modal, Form, Input, Select, DatePicker, InputNumber, Card, message, Popconfirm, Tag, Space } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined, EyeOutlined } from '@ant-design/icons'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useNavigate } from 'react-router-dom'
import { petApi } from '@/services/api'
import type { Pet } from '@/types'
import dayjs from 'dayjs'

export default function Pets() {
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const [modalVisible, setModalVisible] = useState(false)
  const [editingPet, setEditingPet] = useState<Pet | null>(null)
  const [form] = Form.useForm()

  const { data, isLoading } = useQuery({
    queryKey: ['pets', 'list'],
    queryFn: () => petApi.list({ page_size: 100 }),
  })

  const pets: Pet[] = data?.data?.items || []
  const total = data?.data?.total || 0

  const createMutation = useMutation({
    mutationFn: (values: any) => {
      const data = {
        ...values,
        birth_date: values.birth_date ? values.birth_date.toISOString() : undefined,
      }
      if (editingPet) {
        return petApi.update(editingPet.id, data)
      }
      return petApi.create(data)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['pets', 'list'] })
      message.success(editingPet ? '更新成功' : '创建成功')
      setModalVisible(false)
      form.resetFields()
      setEditingPet(null)
    },
    onError: (err: any) => {
      message.error(err.message || '操作失败')
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => petApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['pets', 'list'] })
      message.success('删除成功')
    },
    onError: (err: any) => {
      message.error(err.message || '删除失败')
    },
  })

  const handleAdd = () => {
    setEditingPet(null)
    form.resetFields()
    setModalVisible(true)
  }

  const handleEdit = (pet: Pet) => {
    setEditingPet(pet)
    form.setFieldsValue({
      ...pet,
      birth_date: pet.birth_date ? dayjs(pet.birth_date) : undefined,
    })
    setModalVisible(true)
  }

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      createMutation.mutate(values)
    } catch {}
  }

  const columns = [
    {
      title: '照片',
      dataIndex: 'avatar_url',
      key: 'avatar_url',
      render: (url: string) =>
        url ? (
          <img src={url} alt="" className="w-10 h-10 rounded-full object-cover" />
        ) : (
          <div className="w-10 h-10 rounded-full bg-gray-200 flex items-center justify-center text-xl">
            🐾
          </div>
        ),
    },
    { title: '名字', dataIndex: 'name', key: 'name' },
    { title: '物种', dataIndex: 'species', key: 'species' },
    { title: '品种', dataIndex: 'breed', key: 'breed' },
    { title: '性别', dataIndex: 'gender', key: 'gender', render: (g: string) => g === 'male' ? '公' : '母' },
    {
      title: '体重',
      dataIndex: 'weight',
      key: 'weight',
      render: (w: number) => `${w} kg`,
    },
    {
      title: '操作',
      key: 'actions',
      render: (_: any, pet: Pet) => (
        <Space>
          <Button icon={<EyeOutlined />} size="small" onClick={() => navigate(`/pets/${pet.id}`)}>
            详情
          </Button>
          <Button icon={<EditOutlined />} size="small" onClick={() => handleEdit(pet)}>
            编辑
          </Button>
          <Popconfirm title="确定删除？" onConfirm={() => deleteMutation.mutate(pet.id)}>
            <Button icon={<DeleteOutlined />} size="small" danger>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ]

  return (
    <div className="space-y-4">
      <Card title="宠物档案" extra={<Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>添加宠物</Button>}>
        <Table
          columns={columns}
          dataSource={pets}
          rowKey="id"
          loading={isLoading}
          pagination={{ pageSize: 10, total }}
        />
      </Card>

      <Modal
        title={editingPet ? '编辑宠物' : '添加宠物'}
        open={modalVisible}
        onCancel={() => {
          setModalVisible(false)
          setEditingPet(null)
          form.resetFields()
        }}
        onOk={handleSubmit}
        confirmLoading={createMutation.isPending}
        width={600}
      >
        <Form form={form} layout="vertical" className="mt-4">
          <div className="grid grid-cols-2 gap-4">
            <Form.Item name="name" label="名字" rules={[{ required: true }]}>
              <Input placeholder="宠物名字" />
            </Form.Item>
            <Form.Item name="species" label="物种" rules={[{ required: true }]}>
              <Select
                placeholder="选择物种"
                options={[
                  { value: 'dog', label: '狗' },
                  { value: 'cat', label: '猫' },
                  { value: 'rabbit', label: '兔子' },
                  { value: 'bird', label: '鸟' },
                  { value: 'other', label: '其他' },
                ]}
              />
            </Form.Item>
            <Form.Item name="breed" label="品种">
              <Input placeholder="品种" />
            </Form.Item>
            <Form.Item name="gender" label="性别" rules={[{ required: true }]}>
              <Select
                placeholder="选择性别"
                options={[
                  { value: 'male', label: '公' },
                  { value: 'female', label: '母' },
                ]}
              />
            </Form.Item>
            <Form.Item name="birth_date" label="出生日期">
              <DatePicker className="w-full" />
            </Form.Item>
            <Form.Item name="weight" label="体重(kg)">
              <InputNumber className="w-full" min={0} step={0.1} />
            </Form.Item>
            <Form.Item name="color" label="毛色">
              <Input placeholder="毛色" />
            </Form.Item>
            <Form.Item name="avatar_url" label="照片URL">
              <Input placeholder="宠物照片链接" />
            </Form.Item>
          </div>
          <Form.Item name="allergies" label="过敏史">
            <Input.TextArea rows={2} placeholder="过敏史" />
          </Form.Item>
          <Form.Item name="diet_habit" label="饮食习惯">
            <Input.TextArea rows={2} placeholder="饮食习惯" />
          </Form.Item>
          <Form.Item name="temperament" label="性格特点">
            <Input.TextArea rows={2} placeholder="性格特点" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
