import React, { useState, useEffect } from 'react'
import { Card, Button, List, Modal, Form, Input, message, Empty } from 'antd'
import { addressApi } from '@/services/auth'
import { Address } from '@/types'

const AddressManage: React.FC = () => {
  const [addresses, setAddresses] = useState<Address[]>([])
  const [loading, setLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [editingAddress, setEditingAddress] = useState<Address | null>(null)
  const [form] = Form.useForm()

  useEffect(() => {
    loadAddresses()
  }, [])

  const loadAddresses = async () => {
    setLoading(true)
    try {
      const res = await addressApi.getList()
      setAddresses(res)
    } catch (error) {
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  const handleAdd = () => {
    setEditingAddress(null)
    form.resetFields()
    form.setFieldsValue({ is_default: false })
    setModalVisible(true)
  }

  const handleEdit = (address: Address) => {
    setEditingAddress(address)
    form.setFieldsValue(address)
    setModalVisible(true)
  }

  const handleDelete = async (id: number) => {
    try {
      await addressApi.delete(id)
      message.success('删除成功')
      loadAddresses()
    } catch (error) {
      console.error(error)
    }
  }

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      if (editingAddress) {
        await addressApi.update(editingAddress.id, values)
        message.success('更新成功')
      } else {
        await addressApi.create(values)
        message.success('添加成功')
      }
      setModalVisible(false)
      loadAddresses()
    } catch (error) {
      console.error(error)
    }
  }

  const handleSetDefault = async (id: number) => {
    try {
      await addressApi.setDefault(id)
      message.success('设置成功')
      loadAddresses()
    } catch (error) {
      console.error(error)
    }
  }

  return (
    <div>
      <div className="page-header">
        <h1 className="page-title">地址管理</h1>
        <Button type="primary" onClick={handleAdd}>
          添加地址
        </Button>
      </div>

      <Card loading={loading}>
        {addresses.length === 0 ? (
          <Empty description="暂无地址，请添加" />
        ) : (
          <List
            dataSource={addresses}
            renderItem={(item) => (
              <List.Item
                key={item.id}
                actions={[
                  <Button type="link" onClick={() => handleEdit(item)}>编辑</Button>,
                  <Button type="link" danger onClick={() => handleDelete(item.id)}>删除</Button>,
                  !item.is_default && (
                    <Button type="link" onClick={() => handleSetDefault(item.id)}>
                      设为默认
                    </Button>
                  ),
                ]}
              >
                <List.Item.Meta
                  title={
                    <div>
                      {item.contact_name}
                      <span style={{ marginLeft: 16 }}>{item.contact_phone}</span>
                      {item.is_default && (
                        <span style={{ marginLeft: 8, color: 'blue' }}>[默认]</span>
                      )}
                    </div>
                  }
                  description={`${item.province}${item.city}${item.district}${item.address}`}
                />
              </List.Item>
            )}
          />
        )}
      </Card>

      <Modal
        title={editingAddress ? '编辑地址' : '添加地址'}
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => setModalVisible(false)}
        width={500}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="contact_name"
            label="联系人"
            rules={[{ required: true, message: '请输入联系人' }]}
          >
            <Input placeholder="请输入联系人" />
          </Form.Item>
          <Form.Item
            name="contact_phone"
            label="联系电话"
            rules={[{ required: true, message: '请输入联系电话' }]}
          >
            <Input placeholder="请输入联系电话" />
          </Form.Item>
          <Form.Item
            name="province"
            label="省份"
            rules={[{ required: true, message: '请输入省份' }]}
          >
            <Input placeholder="请输入省份" />
          </Form.Item>
          <Form.Item
            name="city"
            label="城市"
            rules={[{ required: true, message: '请输入城市' }]}
          >
            <Input placeholder="请输入城市" />
          </Form.Item>
          <Form.Item
            name="district"
            label="区县"
            rules={[{ required: true, message: '请输入区县' }]}
          >
            <Input placeholder="请输入区县" />
          </Form.Item>
          <Form.Item
            name="address"
            label="详细地址"
            rules={[{ required: true, message: '请输入详细地址' }]}
          >
            <Input.TextArea rows={2} placeholder="请输入详细地址" />
          </Form.Item>
          <Form.Item name="is_default" label="设为默认地址" valuePropName="checked">
            <Input type="checkbox" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default AddressManage
