import { useState, useEffect } from 'react'
import { Table, Button, Tag, Space, Modal, Form, Input, Select, message, Switch, Card } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined, SettingOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { venueApi } from '@/api'
import { useAuthStore } from '@/store/authStore'
import type { Venue } from '@/types'
import dayjs from 'dayjs'

const VenueList = () => {
  const navigate = useNavigate()
  const { user } = useAuthStore()
  const isAdmin = user?.role === 'admin' || user?.role === 'super_admin'
  const [venues, setVenues] = useState<Venue[]>([])
  const [loading, setLoading] = useState(false)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [isModalVisible, setIsModalVisible] = useState(false)
  const [editingVenue, setEditingVenue] = useState<Venue | null>(null)
  const [priceModalVisible, setPriceModalVisible] = useState(false)
  const [currentVenue, setCurrentVenue] = useState<Venue | null>(null)
  const [form] = Form.useForm()
  const [priceForm] = Form.useForm()

  useEffect(() => {
    loadVenues()
  }, [page, pageSize])

  const loadVenues = async () => {
    setLoading(true)
    try {
      const data: any = await venueApi.list({ page, page_size: pageSize })
      setVenues(data.list)
      setTotal(data.total)
    } catch (error: any) {
      message.error(error.message || '加载失败')
    } finally {
      setLoading(false)
    }
  }

  const handleAdd = () => {
    setEditingVenue(null)
    form.resetFields()
    setIsModalVisible(true)
  }

  const handleEdit = (venue: Venue) => {
    setEditingVenue(venue)
    form.setFieldsValue({
      name: venue.name,
      location: venue.location,
      capacity: venue.capacity,
      facilities: venue.facilities,
      description: venue.description,
      cover_image: venue.cover_image,
    })
    setIsModalVisible(true)
  }

  const handleSubmit = async (values: any) => {
    try {
      if (editingVenue) {
        await venueApi.update(editingVenue.id, values)
        message.success('更新成功')
      } else {
        await venueApi.create(values)
        message.success('创建成功')
      }
      setIsModalVisible(false)
      loadVenues()
    } catch (error: any) {
      message.error(error.message || '操作失败')
    }
  }

  const handleDelete = (id: number) => {
    Modal.confirm({
      title: '确认删除',
      content: '删除后无法恢复，确认要删除这个场地吗？',
      onOk: async () => {
        try {
          await venueApi.delete(id)
          message.success('删除成功')
          loadVenues()
        } catch (error: any) {
          message.error(error.message || '删除失败')
        }
      },
    })
  }

  const handleStatusChange = async (id: number, checked: boolean) => {
    try {
      await venueApi.updateStatus(id, checked ? 'online' : 'offline')
      message.success('状态更新成功')
      loadVenues()
    } catch (error: any) {
      message.error(error.message || '更新失败')
    }
  }

  const handleSetPrice = (venue: Venue) => {
    setCurrentVenue(venue)
    priceForm.resetFields()
    setPriceModalVisible(true)
  }

  const handlePriceSubmit = async (values: any) => {
    if (!currentVenue) return
    try {
      await venueApi.setPrice(currentVenue.id, {
        day_of_week: values.day_of_week,
        time_slots: values.time_slots || [],
      })
      message.success('价格设置成功')
      setPriceModalVisible(false)
    } catch (error: any) {
      message.error(error.message || '设置失败')
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
      render: (text: string, record: Venue) => (
        <a onClick={() => navigate(`/venues/${record.id}`)}>{text}</a>
      ),
    },
    {
      title: '位置',
      dataIndex: 'location',
      key: 'location',
    },
    {
      title: '容纳人数',
      dataIndex: 'capacity',
      key: 'capacity',
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
      render: (_: any, record: Venue) => (
        <Space>
          <Button type="link" icon={<EditOutlined />} onClick={() => handleEdit(record)}>
            编辑
          </Button>
          <Button type="link" icon={<SettingOutlined />} onClick={() => handleSetPrice(record)}>
            价格
          </Button>
          <Switch
            checked={record.status === 'online'}
            onChange={(checked) => handleStatusChange(record.id, checked)}
          />
          <Button type="link" danger icon={<DeleteOutlined />} onClick={() => handleDelete(record.id)}>
            删除
          </Button>
        </Space>
      ),
    },
  ].filter(Boolean)

  return (
    <Card
      title="场地管理"
      extra={isAdmin ? (
        <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
          添加上地
        </Button>
      ) : null}
    >
      <Table
        columns={columns}
        dataSource={venues}
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
        title={editingVenue ? '编辑场地' : '添加上地'}
        open={isModalVisible}
        onCancel={() => setIsModalVisible(false)}
        footer={null}
      >
        <Form form={form} layout="vertical" onFinish={handleSubmit}>
          <Form.Item
            name="name"
            label="场地名称"
            rules={[{ required: true, message: '请输入场地名称' }]}
          >
            <Input placeholder="请输入场地名称" />
          </Form.Item>
          <Form.Item
            name="location"
            label="位置"
            rules={[{ required: true, message: '请输入位置' }]}
          >
            <Input placeholder="请输入位置" />
          </Form.Item>
          <Form.Item
            name="capacity"
            label="容纳人数"
            rules={[{ required: true, message: '请输入容纳人数' }]}
          >
            <Input type="number" placeholder="请输入容纳人数" />
          </Form.Item>
          <Form.Item name="facilities" label="设施设备">
            <Input.TextArea rows={3} placeholder="请输入设施设备" />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input.TextArea rows={3} placeholder="请输入描述" />
          </Form.Item>
          <Form.Item name="cover_image" label="封面图片URL">
            <Input placeholder="请输入封面图片URL" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block>
              {editingVenue ? '更新' : '创建'}
            </Button>
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title="设置价格"
        open={priceModalVisible}
        onCancel={() => setPriceModalVisible(false)}
        footer={null}
        width={600}
      >
        <Form form={priceForm} layout="vertical" onFinish={handlePriceSubmit}>
          <Form.Item
            name="day_of_week"
            label="星期"
            rules={[{ required: true, message: '请选择星期' }]}
          >
            <Select placeholder="请选择星期">
              {['周日', '周一', '周二', '周三', '周四', '周五', '周六'].map((day, index) => (
                <Select.Option key={index} value={index}>
                  {day}
                </Select.Option>
              ))}
            </Select>
          </Form.Item>
          <Form.Item
            name="time_slots"
            label="时间段价格"
            tooltip="可添加多个时间段，格式：开始时间-结束时间-价格"
          >
            <Form.List name="time_slots">
              {(fields, { add, remove }) => (
                <>
                  {fields.map(({ key, name, ...restField }) => (
                    <Space key={key} style={{ display: 'flex', marginBottom: 8 }} align="baseline">
                      <Form.Item
                        {...restField}
                        name={[name, 'start']}
                        rules={[{ required: true, message: '开始时间' }]}
                        style={{ marginBottom: 0 }}
                      >
                        <Input placeholder="开始时间 如: 09:00" style={{ width: 120 }} />
                      </Form.Item>
                      <Form.Item
                        {...restField}
                        name={[name, 'end']}
                        rules={[{ required: true, message: '结束时间' }]}
                        style={{ marginBottom: 0 }}
                      >
                        <Input placeholder="结束时间 如: 12:00" style={{ width: 120 }} />
                      </Form.Item>
                      <Form.Item
                        {...restField}
                        name={[name, 'price']}
                        rules={[{ required: true, message: '价格' }]}
                        style={{ marginBottom: 0 }}
                      >
                        <Input type="number" placeholder="价格" style={{ width: 100 }} />
                      </Form.Item>
                      <Button type="text" danger onClick={() => remove(name)}>
                        删除
                      </Button>
                    </Space>
                  ))}
                  <Button type="dashed" onClick={() => add()} block icon={<PlusOutlined />}>
                    添加时间段
                  </Button>
                </>
              )}
            </Form.List>
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block>
              保存
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </Card>
  )
}

export default VenueList
