import { useEffect, useState } from 'react'
import { Table, Button, Modal, Form, Input, Select, Space, message, Card, Tag, InputNumber } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import api from '../../api'
import type { Product, Inventory } from '../../types'

function MyProductsPage() {
  const [data, setData] = useState<Product[]>([])
  const [loading, setLoading] = useState(false)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [isModalVisible, setIsModalVisible] = useState(false)
  const [editingRecord, setEditingRecord] = useState<Product | null>(null)
  const [inventoryList, setInventoryList] = useState<Inventory[]>([])
  const [form] = Form.useForm()

  const fetchData = async () => {
    setLoading(true)
    try {
      const response = await api.get('/my-products', {
        params: { page, page_size: pageSize },
      })
      setData(response.data as Product[])
      setTotal(response.total || 0)
    } catch (error: any) {
      message.error(error.message || '获取数据失败')
    } finally {
      setLoading(false)
    }
  }

  const fetchInventory = async () => {
    try {
      const response = await api.get('/inventory', { params: { page: 1, page_size: 100 } })
      setInventoryList(response.data || [])
    } catch (error) {
      console.error(error)
    }
  }

  useEffect(() => {
    fetchData()
    fetchInventory()
  }, [page, pageSize])

  const handleAdd = () => {
    setEditingRecord(null)
    form.resetFields()
    setIsModalVisible(true)
  }

  const handleEdit = (record: Product) => {
    setEditingRecord(record)
    form.setFieldsValue({
      title: record.title,
      description: record.description,
      price: record.price,
      stock: record.stock,
      status: record.status,
    })
    setIsModalVisible(true)
  }

  const handleDelete = async (id: number) => {
    try {
      await api.delete(`/products/${id}`)
      message.success('删除成功')
      fetchData()
    } catch (error: any) {
      message.error(error.message || '删除失败')
    }
  }

  const handleSubmit = async (values: any) => {
    try {
      if (editingRecord) {
        await api.put(`/products/${editingRecord.id}`, values)
        message.success('更新成功')
      } else {
        const selectedInventory = inventoryList.find((inv) => inv.id === values.inventory_id)
        await api.post('/products', {
          ...values,
          honey_type: selectedInventory?.honey_type || '',
          batch_code: selectedInventory?.batch_code || '',
          grade: selectedInventory?.grade || 'ungraded',
        })
        message.success('创建成功')
      }
      setIsModalVisible(false)
      fetchData()
    } catch (error: any) {
      message.error(error.message || '操作失败')
    }
  }

  const statusColorMap: Record<string, string> = {
    on_sale: 'green',
    off_shelf: 'default',
    sold_out: 'red',
  }

  const columns = [
    { title: '商品名称', dataIndex: 'title', key: 'title' },
    { title: '蜂蜜类型', dataIndex: 'honey_type', key: 'honey_type' },
    { title: '批次号', dataIndex: 'batch_code', key: 'batch_code' },
    {
      title: '价格',
      dataIndex: 'price',
      key: 'price',
      render: (price: number) => `¥${price}`,
    },
    {
      title: '库存',
      dataIndex: 'stock',
      key: 'stock',
      render: (stock: number, record: Product) => `${stock} ${record.unit}`,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => <Tag color={statusColorMap[status]}>{status}</Tag>,
    },
    { title: '浏览量', dataIndex: 'view_count', key: 'view_count' },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Product) => (
        <Space>
          <Button type="link" icon={<EditOutlined />} onClick={() => handleEdit(record)}>编辑</Button>
          <Button type="link" danger icon={<DeleteOutlined />} onClick={() => handleDelete(record.id)}>删除</Button>
        </Space>
      ),
    },
  ]

  return (
    <Card title="我的商品" extra={<Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>发布商品</Button>}>
      <Table
        columns={columns}
        dataSource={data}
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
        title={editingRecord ? '编辑商品' : '发布商品'}
        open={isModalVisible}
        onCancel={() => setIsModalVisible(false)}
        footer={null}
      >
        <Form form={form} onFinish={handleSubmit} layout="vertical">
          {!editingRecord && (
            <Form.Item name="inventory_id" label="选择库存" rules={[{ required: true, message: '请选择库存' }]}>
              <Select placeholder="请选择要发布的库存">
                {inventoryList.map((inv) => (
                  <Select.Option key={inv.id} value={inv.id}>
                    {inv.honey_type} - 批次{inv.batch_code} - 库存{inv.quantity}kg
                  </Select.Option>
                ))}
              </Select>
            </Form.Item>
          )}
          <Form.Item name="title" label="商品名称" rules={[{ required: true, message: '请输入商品名称' }]}>
            <Input placeholder="商品名称" />
          </Form.Item>
          <Form.Item name="description" label="商品描述">
            <Input.TextArea rows={4} placeholder="商品描述" />
          </Form.Item>
          <Form.Item name="price" label="价格(¥/kg)" rules={[{ required: true, message: '请输入价格' }]}>
            <InputNumber min={0} step={0.01} style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="stock" label="上架数量(kg)" rules={[{ required: true, message: '请输入上架数量' }]}>
            <InputNumber min={0} step={0.1} style={{ width: '100%' }} />
          </Form.Item>
          {editingRecord && (
            <Form.Item name="status" label="状态">
              <Select>
                <Select.Option value="on_sale">上架</Select.Option>
                <Select.Option value="off_shelf">下架</Select.Option>
              </Select>
            </Form.Item>
          )}
          <Form.Item>
            <Button type="primary" htmlType="submit" block>提交</Button>
          </Form.Item>
        </Form>
      </Modal>
    </Card>
  )
}

export default MyProductsPage
