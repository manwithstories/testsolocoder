import { useState, useEffect } from 'react'
import {
  Table, Button, Card, Typography, Tag, Space, Input, Select,
  Pagination, message, Popconfirm, Switch
} from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined, SearchOutlined } from '@ant-design/icons'
import { Link, useNavigate } from 'react-router-dom'
import { productAPI } from '@/api'
import { Product } from '@/types'

const { Title } = Typography
const { Option } = Select

const SellerProducts = () => {
  const navigate = useNavigate()
  const [products, setProducts] = useState<Product[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize] = useState(10)
  const [keyword, setKeyword] = useState('')
  const [status, setStatus] = useState<string>('')

  useEffect(() => {
    loadProducts()
  }, [page, keyword, status])

  const loadProducts = async () => {
    try {
      const params: any = { page, pageSize }
      if (keyword) params.keyword = keyword
      if (status) params.status = status
      const res = await productAPI.myProducts(params)
      setProducts(res.data.data)
      setTotal(res.data.pagination.total)
    } catch (err) {
      console.error('加载商品失败', err)
    }
  }

  const handleStatusChange = async (id: number, checked: boolean) => {
    try {
      await productAPI.update(id, { status: checked ? 'on_sale' : 'off_shelf' })
      message.success('状态更新成功')
      loadProducts()
    } catch (err: any) {
      message.error(err.message || '更新失败')
    }
  }

  const handleDelete = async (id: number) => {
    try {
      await productAPI.delete(id)
      message.success('删除成功')
      loadProducts()
    } catch (err: any) {
      message.error(err.message || '删除失败')
    }
  }

  const columns = [
    {
      title: '商品',
      dataIndex: 'name',
      key: 'name',
      render: (text: string, record: Product) => (
        <div style={{ display: 'flex', gap: 12, alignItems: 'center' }}>
          <img
            src={record.images?.[0] || 'https://via.placeholder.com/60'}
            style={{ width: 60, height: 60, objectFit: 'cover', borderRadius: 4 }}
          />
          <div>
            <div className="truncate" style={{ maxWidth: 200 }}>{text}</div>
            <div style={{ color: '#999', fontSize: 12 }}>库存: {record.stock || 0}</div>
          </div>
        </div>
      ),
    },
    {
      title: '分类',
      dataIndex: ['category', 'name'],
      key: 'category',
      width: 120,
    },
    {
      title: '价格',
      dataIndex: 'price',
      key: 'price',
      width: 120,
      render: (price: number) => <span className="price">¥{price}</span>,
    },
    {
      title: '销量',
      dataIndex: 'sales',
      key: 'sales',
      width: 100,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 120,
      render: (status: string, record: Product) => (
        <Switch
          checked={status === 'on_sale'}
          onChange={(checked) => handleStatusChange(record.id, checked)}
          checkedChildren="上架"
          unCheckedChildren="下架"
        />
      ),
    },
    {
      title: '操作',
      key: 'action',
      width: 160,
      render: (_: any, record: Product) => (
        <Space>
          <Button
            size="small"
            icon={<EditOutlined />}
            onClick={() => navigate(`/seller/products/${record.id}/edit`)}
          >
            编辑
          </Button>
          <Popconfirm title="确定删除该商品？" onConfirm={() => handleDelete(record.id)}>
            <Button size="small" danger icon={<DeleteOutlined />}>删除</Button>
          </Popconfirm>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <div className="page-header">
        <Title level={3} style={{ margin: 0 }}>商品管理</Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => navigate('/seller/products/create')}>
          发布商品
        </Button>
      </div>

      <Card>
        <div style={{ display: 'flex', gap: 16, marginBottom: 16 }}>
          <Input.Search
            placeholder="搜索商品名称..."
            value={keyword}
            onChange={(e) => setKeyword(e.target.value)}
            onSearch={() => setPage(1)}
            style={{ width: 300 }}
            allowClear
          />
          <Select
            placeholder="选择状态"
            allowClear
            style={{ width: 150 }}
            value={status}
            onChange={(val) => { setStatus(val); setPage(1) }}
          >
            <Option value="on_sale">上架中</Option>
            <Option value="off_shelf">已下架</Option>
          </Select>
        </div>

        <Table
          columns={columns}
          dataSource={products}
          rowKey="id"
          pagination={false}
        />

        <div style={{ textAlign: 'center', marginTop: 24 }}>
          <Pagination
            current={page}
            pageSize={pageSize}
            total={total}
            onChange={setPage}
          />
        </div>
      </Card>
    </div>
  )
}

export default SellerProducts
