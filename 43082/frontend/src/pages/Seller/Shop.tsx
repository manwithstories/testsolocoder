import { useState, useEffect } from 'react'
import { Form, Input, Button, Card, Typography, message, Upload, Descriptions, Tag } from 'antd'
import { UploadOutlined, EditOutlined, SaveOutlined } from '@ant-design/icons'
import { shopAPI } from '@/api'
import { Shop } from '@/types'
import type { UploadProps } from 'antd'

const { Title, Text } = Typography
const { TextArea } = Input

const SellerShop = () => {
  const [shop, setShop] = useState<Shop | null>(null)
  const [form] = Form.useForm()
  const [editing, setEditing] = useState(false)
  const [loading, setLoading] = useState(false)
  const [logoUrl, setLogoUrl] = useState('')

  useEffect(() => {
    loadMyShop()
  }, [])

  const loadMyShop = async () => {
    try {
      const res = await shopAPI.getMyShop() as any
      setShop(res.data)
      form.setFieldsValue(res.data)
      setLogoUrl(res.data.logo || '')
    } catch (err) {
      console.error('加载店铺信息失败', err)
    }
  }

  const handleSubmit = async (values: any) => {
    setLoading(true)
    try {
      const data = {
        ...values,
        logo: logoUrl,
      }
      if (shop) {
        await shopAPI.update(data)
      } else {
        await shopAPI.apply(data)
      }
      message.success(shop ? '更新成功' : '申请已提交，等待审核')
      setEditing(false)
      loadMyShop()
    } catch (err: any) {
      message.error(err.message || '操作失败')
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
        setLogoUrl(info.file.response.data.url)
        message.success('上传成功')
      } else if (info.file.status === 'error') {
        message.error('上传失败')
      }
    },
  }

  const getStatusTag = (status: string) => {
    const statusMap: Record<string, { color: string; text: string }> = {
      pending: { color: 'orange', text: '审核中' },
      approved: { color: 'green', text: '已通过' },
      rejected: { color: 'red', text: '已拒绝' },
    }
    const info = statusMap[status] || { color: 'default', text: status }
    return <Tag color={info.color}>{info.text}</Tag>
  }

  if (!shop) {
    return (
      <Card title="申请入驻">
        <Form form={form} layout="vertical" onFinish={handleSubmit}>
          <Form.Item label="店铺Logo">
            <Upload {...uploadProps}>
              {logoUrl ? (
                <img src={logoUrl} alt="logo" style={{ width: 100, height: 100, objectFit: 'cover', borderRadius: '50%' }} />
              ) : (
                <Button icon={<UploadOutlined />}>上传Logo</Button>
              )}
            </Upload>
          </Form.Item>
          <Form.Item name="name" label="店铺名称" rules={[{ required: true, message: '请输入店铺名称' }]}>
            <Input placeholder="请输入店铺名称" />
          </Form.Item>
          <Form.Item name="description" label="店铺描述" rules={[{ required: true, message: '请输入店铺描述' }]}>
            <TextArea rows={4} placeholder="请输入店铺描述" />
          </Form.Item>
          <Form.Item name="businessLicense" label="营业执照号" rules={[{ required: true, message: '请输入营业执照号' }]}>
            <Input placeholder="请输入营业执照号" />
          </Form.Item>
          <Form.Item name="contactName" label="联系人姓名" rules={[{ required: true, message: '请输入联系人姓名' }]}>
            <Input placeholder="请输入联系人姓名" />
          </Form.Item>
          <Form.Item name="contactPhone" label="联系电话" rules={[{ required: true, message: '请输入联系电话' }]}>
            <Input placeholder="请输入联系电话" />
          </Form.Item>
          <Form.Item name="idCardFront" label="身份证正面">
            <Input placeholder="请上传身份证正面照片URL" />
          </Form.Item>
          <Form.Item name="idCardBack" label="身份证背面">
            <Input placeholder="请上传身份证背面照片URL" />
          </Form.Item>
          <Form.Item name="address" label="店铺地址">
            <Input placeholder="请输入店铺地址" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading}>
              提交申请
            </Button>
          </Form.Item>
        </Form>
      </Card>
    )
  }

  return (
    <div>
      <div className="page-header">
        <Title level={3} style={{ margin: 0 }}>
          店铺管理
          {getStatusTag(shop.status)}
        </Title>
        {shop.status === 'approved' && (
          <Button
            icon={editing ? <SaveOutlined /> : <EditOutlined />}
            onClick={() => editing ? form.submit() : setEditing(true)}
            type={editing ? 'primary' : 'default'}
          >
            {editing ? '保存' : '编辑'}
          </Button>
        )}
      </div>

      <Card>
        {editing ? (
          <Form form={form} layout="vertical" onFinish={handleSubmit}>
            <Form.Item label="店铺Logo">
              <Upload {...uploadProps}>
                {logoUrl ? (
                  <img src={logoUrl} alt="logo" style={{ width: 100, height: 100, objectFit: 'cover', borderRadius: '50%' }} />
                ) : (
                  <Button icon={<UploadOutlined />}>上传Logo</Button>
                )}
              </Upload>
            </Form.Item>
            <Form.Item name="name" label="店铺名称" rules={[{ required: true, message: '请输入店铺名称' }]}>
              <Input placeholder="请输入店铺名称" />
            </Form.Item>
            <Form.Item name="description" label="店铺描述">
              <TextArea rows={4} placeholder="请输入店铺描述" />
            </Form.Item>
            <Form.Item name="contactPhone" label="联系电话">
              <Input placeholder="请输入联系电话" />
            </Form.Item>
            <Form.Item>
              <Button type="primary" htmlType="submit" loading={loading}>
                保存修改
              </Button>
              <Button onClick={() => setEditing(false)} style={{ marginLeft: 12 }}>
                取消
              </Button>
            </Form.Item>
          </Form>
        ) : (
          <Descriptions column={1}>
            <Descriptions.Item label="店铺Logo">
              {shop.logo ? (
                <img src={shop.logo} alt="logo" style={{ width: 100, height: 100, objectFit: 'cover', borderRadius: '50%' }} />
              ) : (
                <Text type="secondary">暂无</Text>
              )}
            </Descriptions.Item>
            <Descriptions.Item label="店铺名称">{shop.name}</Descriptions.Item>
            <Descriptions.Item label="店铺描述">{shop.description || '-'}</Descriptions.Item>
            <Descriptions.Item label="联系电话">{shop.contactPhone || '-'}</Descriptions.Item>
            <Descriptions.Item label="好评率">{shop.rating || 5}%</Descriptions.Item>
            <Descriptions.Item label="商品数量">{shop.productCount || 0}</Descriptions.Item>
            <Descriptions.Item label="创建时间">{shop.createdAt}</Descriptions.Item>
          </Descriptions>
        )}
      </Card>
    </div>
  )
}

export default SellerShop
