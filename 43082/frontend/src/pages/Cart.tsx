import { useState, useEffect } from 'react'
import {
  Table, Button, Card, Typography, Checkbox, InputNumber,
  Empty, message, Tag, Space, Popconfirm, Form, Input, Modal
} from 'antd'
import { DeleteOutlined, ShoppingOutlined } from '@ant-design/icons'
import { Link, useNavigate } from 'react-router-dom'
import { cartAPI, orderAPI } from '@/api'
import { useAppStore } from '@/store'
import { CartItem } from '@/types'

const { Title, Text } = Typography

const Cart = () => {
  const navigate = useNavigate()
  const { cart, updateCart, removeFromCart, clearCart } = useAppStore()
  const [selectedRowKeys, setSelectedRowKeys] = useState<number[]>([])
  const [loading, setLoading] = useState(false)
  const [items, setItems] = useState<CartItem[]>([])
  const [checkoutModalVisible, setCheckoutModalVisible] = useState(false)
  const [form] = Form.useForm()

  useEffect(() => {
    loadCart()
  }, [])

  const loadCart = async () => {
    try {
      const res = await cartAPI.list() as any
      setItems(res.data)
    } catch (err) {
      console.error('加载购物车失败', err)
    }
  }

  const handleQuantityChange = async (id: number, quantity: number) => {
    try {
      await cartAPI.update(id, { quantity })
      updateCart(id, quantity)
      setItems(items.map(item => item.id === id ? { ...item, quantity } : item))
    } catch (err: any) {
      message.error(err.message || '更新失败')
    }
  }

  const handleRemove = async (id: number) => {
    try {
      await cartAPI.delete(id)
      removeFromCart(id)
      setItems(items.filter(item => item.id !== id))
      message.success('已移除')
    } catch (err: any) {
      message.error(err.message || '删除失败')
    }
  }

  const handleCheckout = () => {
    if (selectedRowKeys.length === 0) {
      message.warning('请选择要结算的商品')
      return
    }
    setCheckoutModalVisible(true)
  }

  const handleSubmitOrder = async (values: any) => {
    setLoading(true)
    try {
      const selectedItems = items.filter(item => selectedRowKeys.includes(item.id))
      const res = await orderAPI.create({
        cartIds: selectedItems.map(item => item.id),
        receiverName: values.receiverName,
        receiverPhone: values.receiverPhone,
        receiverAddress: values.receiverAddress,
        remark: values.remark,
      })
      message.success('订单创建成功')
      clearCart()
      setCheckoutModalVisible(false)
      navigate(`/orders`)
    } catch (err: any) {
      message.error(err.message || '下单失败')
    } finally {
      setLoading(false)
    }
  }

  const totalPrice = items
    .filter(item => selectedRowKeys.includes(item.id))
    .reduce((sum, item) => sum + item.price * item.quantity, 0)

  const columns = [
    {
      title: '商品',
      dataIndex: 'productName',
      key: 'productName',
      render: (text: string, record: CartItem) => (
        <Link to={`/products/${record.productId}`} style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
          <img
            src={record.productImage || 'https://via.placeholder.com/80'}
            alt={text}
            style={{ width: 60, height: 60, objectFit: 'cover', borderRadius: 4 }}
          />
          <div>
            <div>{text}</div>
            {record.specValues && (
              <Tag color="blue" style={{ marginTop: 4 }}>
                {record.specValues.join(' / ')}
              </Tag>
            )}
          </div>
        </Link>
      ),
    },
    {
      title: '单价',
      dataIndex: 'price',
      key: 'price',
      width: 120,
      render: (price: number) => <span className="price">¥{price}</span>,
    },
    {
      title: '数量',
      dataIndex: 'quantity',
      key: 'quantity',
      width: 160,
      render: (quantity: number, record: CartItem) => (
        <InputNumber
          min={1}
          max={record.stock || 99}
          value={quantity}
          onChange={(val) => handleQuantityChange(record.id, val || 1)}
        />
      ),
    },
    {
      title: '小计',
      key: 'subtotal',
      width: 120,
      render: (_: any, record: CartItem) => (
        <span className="price">¥{(record.price * record.quantity).toFixed(2)}</span>
      ),
    },
    {
      title: '操作',
      key: 'action',
      width: 80,
      render: (_: any, record: CartItem) => (
        <Popconfirm title="确定移除该商品？" onConfirm={() => handleRemove(record.id)}>
          <Button type="text" danger icon={<DeleteOutlined />} />
        </Popconfirm>
      ),
    },
  ]

  const rowSelection = {
    selectedRowKeys,
    onChange: (keys: React.Key[]) => setSelectedRowKeys(keys as number[]),
  }

  return (
    <div>
      <Title level={3}>购物车</Title>

      <Card>
        {items.length > 0 ? (
          <>
            <Table
              rowSelection={rowSelection}
              columns={columns}
              dataSource={items}
              rowKey="id"
              pagination={false}
            />

            <div style={{
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'center',
              padding: '16px 24px',
              background: '#fafafa',
              borderRadius: 4,
              marginTop: 16,
            }}>
              <div>
                <Checkbox
                  checked={selectedRowKeys.length === items.length && items.length > 0}
                  onChange={(e) => {
                    if (e.target.checked) {
                      setSelectedRowKeys(items.map(item => item.id))
                    } else {
                      setSelectedRowKeys([])
                    }
                  }}
                >
                  全选
                </Checkbox>
                <Button type="text" danger style={{ marginLeft: 16 }}>
                  删除选中
                </Button>
              </div>
              <div style={{ display: 'flex', alignItems: 'center', gap: 24 }}>
                <div>
                  <Text>已选 {selectedRowKeys.length} 件商品</Text>
                </div>
                <div>
                  <Text>合计: </Text>
                  <span className="price" style={{ fontSize: 24 }}>¥{totalPrice.toFixed(2)}</span>
                </div>
                <Button
                  type="primary"
                  size="large"
                  icon={<ShoppingOutlined />}
                  onClick={handleCheckout}
                  disabled={selectedRowKeys.length === 0}
                >
                  结算
                </Button>
              </div>
            </div>
          </>
        ) : (
          <Empty description="购物车是空的">
            <Link to="/products">
              <Button type="primary">去逛逛</Button>
            </Link>
          </Empty>
        )}
      </Card>

      <Modal
        title="填写收货地址"
        open={checkoutModalVisible}
        onCancel={() => setCheckoutModalVisible(false)}
        footer={null}
        width={500}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSubmitOrder}
        >
          <Form.Item
            name="receiverName"
            label="收货人姓名"
            rules={[{ required: true, message: '请输入收货人姓名' }]}
          >
            <Input placeholder="请输入收货人姓名" />
          </Form.Item>
          <Form.Item
            name="receiverPhone"
            label="收货人电话"
            rules={[
              { required: true, message: '请输入收货人电话' },
              { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号' }
            ]}
          >
            <Input placeholder="请输入收货人电话" />
          </Form.Item>
          <Form.Item
            name="receiverAddress"
            label="收货地址"
            rules={[{ required: true, message: '请输入收货地址' }]}
          >
            <Input.TextArea rows={3} placeholder="请输入详细收货地址" />
          </Form.Item>
          <Form.Item
            name="remark"
            label="订单备注"
          >
            <Input.TextArea rows={2} placeholder="选填，如有特殊要求请备注" />
          </Form.Item>
          <Form.Item style={{ marginBottom: 0, textAlign: 'right' }}>
            <Space>
              <Button onClick={() => setCheckoutModalVisible(false)}>取消</Button>
              <Button type="primary" htmlType="submit" loading={loading}>
                提交订单
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default Cart
