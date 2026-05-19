import { useState, useEffect, useRef } from 'react'
import { Table, Button, Tag, Space, Modal, Form, Input, Select, Switch, message, Card, Upload } from 'antd'
import { PlusOutlined, EditOutlined, ImportOutlined } from '@ant-design/icons'
import type { UploadProps } from 'antd'
import { deviceApi } from '@/api'
import { useAuthStore } from '@/store/authStore'
import type { Device, DeviceCategory } from '@/types'
import dayjs from 'dayjs'

const DeviceList = () => {
  const { user } = useAuthStore()
  const isAdmin = user?.role === 'admin' || user?.role === 'super_admin'
  const [devices, setDevices] = useState<Device[]>([])
  const [categories, setCategories] = useState<DeviceCategory[]>([])
  const [loading, setLoading] = useState(false)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [categoryFilter, setCategoryFilter] = useState<number | undefined>()
  const [isModalVisible, setIsModalVisible] = useState(false)
  const [editingDevice, setEditingDevice] = useState<Device | null>(null)
  const [importModalVisible, setImportModalVisible] = useState(false)
  const [importText, setImportText] = useState('')
  const [form] = Form.useForm()

  useEffect(() => {
    loadCategories()
  }, [])

  useEffect(() => {
    loadDevices()
  }, [page, pageSize, categoryFilter])

  const loadCategories = async () => {
    try {
      const data = await deviceApi.listCategories()
      setCategories(data)
    } catch (error: any) {
      message.error(error.message || '加载分类失败')
    }
  }

  const loadDevices = async () => {
    setLoading(true)
    try {
      const params: any = { page, page_size: pageSize }
      if (categoryFilter) params.category_id = categoryFilter
      const data: any = await deviceApi.list(params)
      setDevices(data.list)
      setTotal(data.total)
    } catch (error: any) {
      message.error(error.message || '加载失败')
    } finally {
      setLoading(false)
    }
  }

  const handleAdd = () => {
    setEditingDevice(null)
    form.resetFields()
    setIsModalVisible(true)
  }

  const handleEdit = (device: Device) => {
    setEditingDevice(device)
    form.setFieldsValue(device)
    setIsModalVisible(true)
  }

  const handleSubmit = async (values: any) => {
    try {
      if (editingDevice) {
        await deviceApi.update(editingDevice.id, values)
        message.success('更新成功')
      } else {
        await deviceApi.create(values)
        message.success('创建成功')
      }
      setIsModalVisible(false)
      loadDevices()
    } catch (error: any) {
      message.error(error.message || '操作失败')
    }
  }

  const handleStatusChange = async (id: number, checked: boolean) => {
    try {
      await deviceApi.updateStatus(id, checked ? 'online' : 'offline')
      message.success('状态更新成功')
      loadDevices()
    } catch (error: any) {
      message.error(error.message || '更新失败')
    }
  }

  const handleBatchImport = () => {
    setImportText('')
    setImportModalVisible(true)
  }

  const handleImportSubmit = async () => {
    try {
      const items = JSON.parse(importText)
      if (!Array.isArray(items)) {
        throw new Error('请输入有效的JSON数组')
      }
      const result = await deviceApi.batchImport(items)
      message.success(`导入完成：成功 ${result.success_count} 条，失败 ${result.fail_count} 条`)
      if (result.errors && result.errors.length > 0) {
        message.warning(`失败详情：${result.errors.join(', ')}`)
      }
      setImportModalVisible(false)
      loadDevices()
    } catch (error: any) {
      message.error(error.message || '导入失败')
    }
  }

  const handleFileImport: UploadProps['onChange'] = (info) => {
    if (info.file.status === 'done') {
      const reader = new FileReader()
      reader.onload = (e) => {
        try {
          const text = e.target?.result as string
          setImportText(text)
          message.success('文件读取成功')
        } catch (err) {
          message.error('文件读取失败')
        }
      }
      reader.readAsText(info.file.originFileObj as File)
    }
  }

  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 60,
    },
    {
      title: '名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '分类',
      dataIndex: ['category', 'name'],
      key: 'category',
    },
    {
      title: '库存',
      dataIndex: 'stock_quantity',
      key: 'stock_quantity',
    },
    {
      title: '可用',
      dataIndex: 'available_quantity',
      key: 'available_quantity',
    },
    {
      title: '租赁价格',
      dataIndex: 'rental_price',
      key: 'rental_price',
      render: (val: number) => `¥${val.toFixed(2)}/小时`,
    },
    {
      title: '押金',
      dataIndex: 'deposit_amount',
      key: 'deposit_amount',
      render: (val: number) => `¥${val.toFixed(2)}`,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => {
        const color = status === 'online' ? 'green' : 'red'
        const text = status === 'online' ? '上架' : '下架'
        return <Tag color={color}>{text}</Tag>
      },
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (val: string) => dayjs(val).format('YYYY-MM-DD HH:mm'),
    },
    isAdmin && {
      title: '操作',
      key: 'action',
      render: (_: any, record: Device) => (
        <Space>
          <Button type="link" icon={<EditOutlined />} onClick={() => handleEdit(record)}>
            编辑
          </Button>
          <Switch
            checked={record.status === 'online'}
            onChange={(checked) => handleStatusChange(record.id, checked)}
          />
        </Space>
      ),
    },
  ].filter(Boolean)

  return (
    <Card
      title="设备管理"
      extra={
        <Space>
          <Select
            placeholder="选择分类"
            style={{ width: 150 }}
            allowClear
            value={categoryFilter}
            onChange={setCategoryFilter}
          >
            {categories.map((cat) => (
              <Select.Option key={cat.id} value={cat.id}>
                {cat.name}
              </Select.Option>
            ))}
          </Select>
          {isAdmin && (
            <>
              <Button icon={<ImportOutlined />} onClick={handleBatchImport}>
                批量导入
              </Button>
              <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
              添设备
            </Button>
            </>
          )}
        </Space>
      }
    >
      <Table
        columns={columns}
        dataSource={devices}
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
        title={editingDevice ? '编辑设备' : '添设备'}
        open={isModalVisible}
        onCancel={() => setIsModalVisible(false)}
        footer={null}
      >
        <Form form={form} layout="vertical" onFinish={handleSubmit}>
          <Form.Item
            name="category_id"
            label="分类"
            rules={[{ required: true, message: '请选择分类' }]}
          >
            <Select placeholder="请选择分类">
              {categories.map((cat) => (
                <Select.Option key={cat.id} value={cat.id}>
                  {cat.name}
                </Select.Option>
              ))}
            </Select>
          </Form.Item>
          <Form.Item
            name="name"
            label="设备名称"
            rules={[{ required: true, message: '请输入设备名称' }]}
          >
            <Input placeholder="请输入设备名称" />
          </Form.Item>
          <Form.Item
            name="stock_quantity"
            label="库存数量"
            rules={[{ required: true, message: '请输入库存数量' }]}
          >
            <Input type="number" placeholder="请输入库存数量" />
          </Form.Item>
          <Form.Item
            name="rental_price"
            label="租赁单价(元/小时)"
            rules={[{ required: true, message: '请输入租赁单价' }]}
          >
            <Input type="number" step="0.01" placeholder="请输入租赁单价" />
          </Form.Item>
          <Form.Item
            name="deposit_amount"
            label="押金金额"
            rules={[{ required: true, message: '请输入押金金额' }]}
          >
            <Input type="number" step="0.01" placeholder="请输入押金金额" />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input.TextArea rows={3} placeholder="请输入描述" />
          </Form.Item>
          <Form.Item name="specification" label="规格参数">
            <Input.TextArea rows={3} placeholder="请输入规格参数" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block>
              {editingDevice ? '更新' : '创建'}
            </Button>
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title="批量导入设备"
        open={importModalVisible}
        onCancel={() => setImportModalVisible(false)}
        footer={[
          <Button key="cancel" onClick={() => setImportModalVisible(false)}>
            取消
          </Button>,
          <Button key="submit" type="primary" onClick={handleImportSubmit}>
            导入
          </Button>,
        ]}
        width={700}
      >
        <p style={{ marginBottom: 12, color: '#666' }}>
          请输入JSON格式的设备数据数组，格式示例：
          <pre style={{ background: '#f5f5f5', padding: 8, borderRadius: 4, marginTop: 8 }}>
{`[
  {
    "category_name": "投影仪",
    "name": "EPSON投影仪",
    "description": "高清投影仪",
    "specification": "1080P",
    "stock_quantity": 10,
    "rental_price": 50,
    "deposit_amount": 200
  }
]`}
          </pre>
        </p>
        <div style={{ marginBottom: 12 }}>
          <Upload
            accept=".json,.txt"
            showUploadList={false}
            customRequest={({ file, onSuccess }) => {
              const reader = new FileReader()
              reader.onload = (e) => {
                setImportText(e.target?.result as string)
                message.success('文件读取成功')
                onSuccess?.(null)
              }
              reader.readAsText(file as File)
            }}
          >
            <Button icon={<ImportOutlined />}>上传JSON文件</Button>
          </Upload>
        </div>
        <Input.TextArea
          value={importText}
          onChange={(e) => setImportText(e.target.value)}
          rows={12}
          placeholder="请粘贴JSON格式的设备数据..."
        />
      </Modal>
    </Card>
  )
}

export default DeviceList
