import { useState, useEffect } from 'react'
import {
  Form, Input, Button, Card, Typography, InputNumber, Select,
  Upload, message, List, Space
} from 'antd'
import { ArrowLeftOutlined, PlusOutlined, UploadOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { productAPI, categoryAPI } from '@/api'
import { Category } from '@/types'
import type { UploadProps } from 'antd'

const { Title } = Typography
const { TextArea } = Input
const { Option } = Select

const ProductCreate = () => {
  const navigate = useNavigate()
  const [form] = Form.useForm()
  const [categories, setCategories] = useState<Category[]>([])
  const [images, setImages] = useState<string[]>([])
  const [loading, setLoading] = useState(false)
  const [specs, setSpecs] = useState<{ name: string; values: string[] }[]>([])
  const [newSpecName, setNewSpecName] = useState('')
  const [newSpecValues, setNewSpecValues] = useState('')

  useEffect(() => {
    loadCategories()
  }, [])

  const loadCategories = async () => {
    try {
      const res = await categoryAPI.list() as any
      setCategories(res.data)
    } catch (err) {
      console.error('加载分类失败', err)
    }
  }

  const handleSubmit = async (values: any) => {
    if (images.length === 0) {
      message.warning('请至少上传一张商品图片')
      return
    }
    setLoading(true)
    try {
      const { originalPrice, ...restValues } = values
      let skus: any[] | undefined
      if (specs.length > 0) {
        skus = specs.flatMap(spec => 
          spec.values.map(value => ({
            specValues: [value],
            price: restValues.price,
            stock: restValues.stock || 0,
          }))
        )
      }
      const data = {
        ...restValues,
        images,
        mainImage: images[0],
        specs: specs.length > 0 ? specs : undefined,
        skus,
      }
      await productAPI.create(data)
      message.success('商品创建成功')
      navigate('/seller/products')
    } catch (err: any) {
      message.error(err.message || '创建失败')
    } finally {
      setLoading(false)
    }
  }

  const uploadProps: UploadProps = {
    name: 'file',
    action: '/api/upload',
    headers: {
      Authorization: `Bearer ${localStorage.getItem('token')}`,
    },
    showUploadList: false,
    onChange(info) {
      if (info.file.status === 'done') {
        setImages([...images, info.file.response.data.url])
        message.success('上传成功')
      } else if (info.file.status === 'error') {
        message.error('上传失败')
      }
    },
  }

  const removeImage = (index: number) => {
    setImages(images.filter((_, i) => i !== index))
  }

  const addSpec = () => {
    if (!newSpecName.trim() || !newSpecValues.trim()) {
      message.warning('请输入规格名称和规格值')
      return
    }
    setSpecs([...specs, {
      name: newSpecName.trim(),
      values: newSpecValues.split(',').map(v => v.trim()).filter(v => v),
    }])
    setNewSpecName('')
    setNewSpecValues('')
  }

  const removeSpec = (index: number) => {
    setSpecs(specs.filter((_, i) => i !== index))
  }

  return (
    <div>
      <div className="page-header">
        <Title level={3} style={{ margin: 0 }}>
          <Button
            type="text"
            icon={<ArrowLeftOutlined />}
            onClick={() => navigate('/seller/products')}
            style={{ marginRight: 16 }}
          />
          发布商品
        </Title>
      </div>

      <Card>
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
          initialValues={{ status: 'on_sale' }}
        >
          <Form.Item
            name="name"
            label="商品名称"
            rules={[{ required: true, message: '请输入商品名称' }]}
          >
            <Input placeholder="请输入商品名称" />
          </Form.Item>

          <Form.Item
            name="categoryId"
            label="商品分类"
            rules={[{ required: true, message: '请选择商品分类' }]}
          >
            <Select placeholder="请选择分类">
              {categories.map((cat) => (
                <Option key={cat.id} value={cat.id}>{cat.name}</Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item
            name="description"
            label="商品描述"
            rules={[{ required: true, message: '请输入商品描述' }]}
          >
            <TextArea rows={4} placeholder="请输入商品描述" />
          </Form.Item>

          <Form.Item label="商品图片">
            <div style={{ display: 'flex', gap: 12, flexWrap: 'wrap', marginBottom: 12 }}>
              {images.map((img, index) => (
                <div key={index} style={{ position: 'relative' }}>
                  <img
                    src={img}
                    alt=""
                    style={{ width: 120, height: 120, objectFit: 'cover', borderRadius: 4 }}
                  />
                  <Button
                    type="text"
                    danger
                    size="small"
                    style={{ position: 'absolute', top: 0, right: 0 }}
                    onClick={() => removeImage(index)}
                  >
                    ×
                  </Button>
                </div>
              ))}
              <Upload {...uploadProps}>
                <Button icon={<UploadOutlined />} style={{ width: 120, height: 120 }}>
                  上传图片
                </Button>
              </Upload>
            </div>
          </Form.Item>

          <Form.Item
            name="price"
            label="商品价格"
            rules={[{ required: true, message: '请输入商品价格' }]}
          >
            <InputNumber
              min={0.01}
              step={0.01}
              style={{ width: '100%' }}
              placeholder="请输入商品价格"
              prefix="¥"
            />
          </Form.Item>

          <Form.Item
            name="stock"
            label="库存数量"
            rules={[{ required: true, message: '请输入库存数量' }]}
          >
            <InputNumber
              min={0}
              style={{ width: '100%' }}
              placeholder="请输入库存数量"
            />
          </Form.Item>

          <Form.Item label="商品规格（可选）">
            <div style={{ marginBottom: 12 }}>
              <Space>
                <Input
                  placeholder="规格名称（如：颜色）"
                  value={newSpecName}
                  onChange={(e) => setNewSpecName(e.target.value)}
                  style={{ width: 150 }}
                />
                <Input
                  placeholder="规格值，逗号分隔（如：红色,蓝色,绿色）"
                  value={newSpecValues}
                  onChange={(e) => setNewSpecValues(e.target.value)}
                  style={{ width: 300 }}
                />
                <Button icon={<PlusOutlined />} onClick={addSpec}>添加规格</Button>
              </Space>
            </div>
            <List
              dataSource={specs}
              locale={{ emptyText: '暂无规格' }}
              renderItem={(spec, index) => (
                <List.Item
                  key={index}
                  actions={[
                    <Button type="text" danger onClick={() => removeSpec(index)}>删除</Button>
                  ]}
                >
                  <List.Item.Meta
                    title={spec.name}
                    description={spec.values.join(' / ')}
                  />
                </List.Item>
              )}
            />
          </Form.Item>

          <Form.Item
            name="status"
            label="商品状态"
            rules={[{ required: true, message: '请选择商品状态' }]}
          >
            <Select>
              <Option value="on_sale">立即上架</Option>
              <Option value="off_sale">暂不上架</Option>
            </Select>
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading} size="large">
              发布商品
            </Button>
            <Button onClick={() => navigate('/seller/products')} style={{ marginLeft: 12 }} size="large">
              取消
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  )
}

export default ProductCreate
