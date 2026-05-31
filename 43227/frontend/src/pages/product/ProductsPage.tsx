import { useEffect, useState } from 'react'
import { Table, Card, Tag, Input, Select, Button, Space, Image, Rate, message } from 'antd'
import { ShoppingCartOutlined, SearchOutlined, EyeOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import api from '../../api'
import type { Product } from '../../types'

function ProductsPage() {
  const navigate = useNavigate()
  const [data, setData] = useState<Product[]>([])
  const [loading, setLoading] = useState(false)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [keyword, setKeyword] = useState('')
  const [honeyType, setHoneyType] = useState<string | undefined>()
  const [grade, setGrade] = useState<string | undefined>()

  const fetchData = async () => {
    setLoading(true)
    try {
      const params: any = { page, page_size: pageSize }
      if (keyword) params.keyword = keyword
      if (honeyType) params.honey_type = honeyType
      if (grade) params.grade = grade
      const response = await api.get('/products', { params })
      setData(response.data as Product[])
      setTotal(response.total || 0)
    } catch (error: any) {
      message.error(error.message || '获取数据失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchData()
  }, [page, pageSize])

  const handleSearch = () => {
    setPage(1)
    fetchData()
  }

  const columns = [
    {
      title: '图片',
      dataIndex: 'images',
      key: 'images',
      width: 100,
      render: (images: string[]) => (
        images && images.length > 0 ? (
          <Image src={images[0]} width={60} height={60} style={{ objectFit: 'cover' }} />
        ) : <div style={{ width: 60, height: 60, background: '#f5f5f5', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>无图</div>
      ),
    },
    { title: '商品名称', dataIndex: 'title', key: 'title' },
    { title: '蜂蜜类型', dataIndex: 'honey_type', key: 'honey_type' },
    {
      title: '等级',
      dataIndex: 'grade',
      key: 'grade',
      render: (g: string) => {
        if (!g) return '-'
        const colorMap: Record<string, string> = { grade_a: 'gold', grade_b: 'blue', grade_c: 'green' }
        return <Tag color={colorMap[g]}>{g}</Tag>
      },
    },
    {
      title: '价格',
      dataIndex: 'price',
      key: 'price',
      render: (price: number) => <span style={{ color: '#f5222d', fontWeight: 'bold' }}>¥{price}/kg</span>,
    },
    {
      title: '库存',
      dataIndex: 'stock',
      key: 'stock',
      render: (stock: number, record: Product) => `${stock} ${record.unit}`,
    },
    {
      title: '卖家',
      key: 'seller',
      render: (_: any, record: Product) => (
        <Space>
          <span>{record.user?.username}</span>
          {record.user?.reputation != null && <Rate disabled allowHalf value={record.user.reputation} style={{ fontSize: 12 }} />}
        </Space>
      ),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Product) => (
        <Space>
          <Button type="link" icon={<EyeOutlined />} onClick={() => navigate(`/products/${record.id}`)}>
            查看
          </Button>
          <Button type="primary" icon={<ShoppingCartOutlined />} onClick={() => navigate(`/products/${record.id}`)}>
            购买
          </Button>
        </Space>
      ),
    },
  ]

  return (
    <Card title="蜂蜜交易市场">
      <Space style={{ marginBottom: 16 }}>
        <Input
          placeholder="搜索商品名称"
          prefix={<SearchOutlined />}
          value={keyword}
          onChange={(e) => setKeyword(e.target.value)}
          style={{ width: 200 }}
          onPressEnter={handleSearch}
        />
        <Select
          placeholder="蜂蜜类型"
          value={honeyType}
          onChange={setHoneyType}
          style={{ width: 150 }}
          allowClear
        >
          <Select.Option value="洋槐蜜">洋槐蜜</Select.Option>
          <Select.Option value="枣花蜜">枣花蜜</Select.Option>
          <Select.Option value="百花蜜">百花蜜</Select.Option>
          <Select.Option value="荆条蜜">荆条蜜</Select.Option>
        </Select>
        <Select
          placeholder="等级"
          value={grade}
          onChange={setGrade}
          style={{ width: 120 }}
          allowClear
        >
          <Select.Option value="grade_a">A级</Select.Option>
          <Select.Option value="grade_b">B级</Select.Option>
          <Select.Option value="grade_c">C级</Select.Option>
        </Select>
        <Button type="primary" onClick={handleSearch}>搜索</Button>
      </Space>
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
    </Card>
  )
}

export default ProductsPage
