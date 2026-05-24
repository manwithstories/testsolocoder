import React, { useEffect, useState } from 'react'
import {
  Row,
  Col,
  Card,
  Tag,
  Descriptions,
  Button,
  Modal,
  Form,
  Input,
  Select,
  Upload,
  message,
  Tabs,
  Timeline,
  Space,
  Spin,
  Checkbox,
  DatePicker,
  TimePicker,
  Table,
  Empty,
  Statistic,
} from 'antd'
import {
  EditOutlined,
  UploadOutlined,
  DeleteOutlined,
  HeartOutlined,
  CalendarOutlined,
  FileTextOutlined,
} from '@ant-design/icons'
import { useNavigate, useParams } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'
import { getPet, updatePet, uploadPetPhotos } from '../api/pet'
import { createAdoptionApplication } from '../api/adoption'
import { listHealthRecords, getPetHealthSummary, getHealthReminders } from '../api/health'
import { getPetAdoptionHistory } from '../api/pet'
import { createAppointment } from '../api/appointment'
import { Pet, HealthRecord, AdoptionApplication } from '../types'

const { Option } = Select
const { TextArea } = Input

const statusColorMap: Record<string, string> = {
  adoptable: 'green',
  adopted: 'blue',
  treatment: 'orange',
  deceased: 'default',
}

const statusTextMap: Record<string, string> = {
  adoptable: '待领养',
  adopted: '已领养',
  treatment: '治疗中',
  deceased: '已去世',
}

const genderMap: Record<string, string> = {
  male: '公',
  female: '母',
  unknown: '未知',
}

const PetDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>()
  const { user } = useAuth()
  const navigate = useNavigate()
  const [pet, setPet] = useState<Pet | null>(null)
  const [healthRecords, setHealthRecords] = useState<HealthRecord[]>([])
  const [adoptionHistory, setAdoptionHistory] = useState<AdoptionApplication[]>([])
  const [healthSummary, setHealthSummary] = useState<any>(null)
  const [reminders, setReminders] = useState<any[]>([])
  const [loading, setLoading] = useState(false)
  const [adoptionModal, setAdoptionModal] = useState(false)
  const [appointmentModal, setAppointmentModal] = useState(false)
  const [editModal, setEditModal] = useState(false)
  const [adoptionForm] = Form.useForm()
  const [appointmentForm] = Form.useForm()
  const [editForm] = Form.useForm()
  const [submitting, setSubmitting] = useState(false)

  const isRescue = user?.role === 'rescue'
  const isOwner = isRescue && pet?.rescue_id === user?.rescue_id
  const canAdopt = user?.role === 'adopter' && pet?.status === 'adoptable'

  useEffect(() => {
    loadData()
  }, [id])

  const loadData = async () => {
    if (!id) return
    setLoading(true)
    try {
      const [petRes, healthRes, historyRes, summaryRes, remindersRes] = await Promise.all([
        getPet(Number(id)),
        listHealthRecords({ pet_id: Number(id), page_size: 100 }),
        getPetAdoptionHistory(Number(id)),
        getPetHealthSummary(Number(id)),
        getHealthReminders(Number(id)),
      ])

      if (petRes.code === 0) setPet(petRes.data as Pet)
      if (healthRes.code === 0) setHealthRecords((healthRes.data as any)?.items || [])
      if (historyRes.code === 0) setAdoptionHistory(historyRes.data as AdoptionApplication[])
      if (summaryRes.code === 0) setHealthSummary(summaryRes.data)
      if (remindersRes.code === 0) setReminders(remindersRes.data as any[] || [])
    } catch (error) {
      console.error('Failed to load pet detail:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleAdoptionSubmit = async () => {
    try {
      const values = await adoptionForm.validateFields()
      setSubmitting(true)
      await createAdoptionApplication({
        pet_id: Number(id),
        ...values,
      })
      message.success('领养申请已提交')
      setAdoptionModal(false)
      adoptionForm.resetFields()
    } catch (error: any) {
      message.error(error.message || '提交失败')
    } finally {
      setSubmitting(false)
    }
  }

  const handleAppointmentSubmit = async () => {
    try {
      const values = await appointmentForm.validateFields()
      setSubmitting(true)
      await createAppointment({
        pet_id: Number(id),
        appointment_type: 'visit',
        ...values,
      })
      message.success('预约已提交')
      setAppointmentModal(false)
      appointmentForm.resetFields()
    } catch (error: any) {
      message.error(error.message || '提交失败')
    } finally {
      setSubmitting(false)
    }
  }

  const handleEditSubmit = async () => {
    try {
      const values = await editForm.validateFields()
      setSubmitting(true)
      await updatePet(Number(id), values)
      message.success('信息已更新')
      setEditModal(false)
      loadData()
    } catch (error: any) {
      message.error(error.message || '更新失败')
    } finally {
      setSubmitting(false)
    }
  }

  const handlePhotoUpload = async (file: File) => {
    const formData = new FormData()
    formData.append('photos', file)
    try {
      await uploadPetPhotos(Number(id), formData)
      message.success('照片上传成功')
      loadData()
    } catch (error: any) {
      message.error('上传失败')
    }
  }

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: 60 }}>
        <Spin size="large" />
      </div>
    )
  }

  if (!pet) {
    return <div>宠物不存在</div>
  }

  const photos = pet.photos ? pet.photos.split(',').filter(Boolean) : []

  return (
    <div>
      <Card style={{ marginBottom: 16 }}>
        <div className="pet-detail-header">
          <div>
            {photos.length > 0 ? (
              <img
                src={photos[0]}
                alt={pet.name}
                className="pet-detail-photo"
              />
            ) : (
              <div
                style={{
                  width: 300,
                  height: 300,
                  background: '#f0f0f0',
                  borderRadius: 12,
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  fontSize: 80,
                }}
              >
                🐾
              </div>
            )}
            {photos.length > 1 && (
              <div className="pet-gallery" style={{ marginTop: 12 }}>
                {photos.slice(1, 5).map((photo, index) => (
                  <img key={index} src={photo} alt={`${pet.name} ${index + 2}`} />
                ))}
              </div>
            )}
            {isOwner && (
              <Upload
                showUploadList={false}
                beforeUpload={(file) => {
                  handlePhotoUpload(file)
                  return false
                }}
                accept="image/*"
                multiple
              >
                <Button icon={<UploadOutlined />} style={{ marginTop: 12 }}>
                  上传照片
                </Button>
              </Upload>
            )}
          </div>
          <div className="pet-detail-info">
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <h1 style={{ margin: 0 }}>{pet.name}</h1>
              <Tag color={statusColorMap[pet.status]} style={{ fontSize: 14, padding: '4px 12px' }}>
                {statusTextMap[pet.status]}
              </Tag>
            </div>
            <p style={{ color: '#999', marginTop: 8 }}>档案编号: {pet.archive_number}</p>

            <Descriptions column={2} size="small" style={{ marginTop: 16 }}>
              <Descriptions.Item label="物种">{pet.species}</Descriptions.Item>
              <Descriptions.Item label="品种">{pet.breed || '-'}</Descriptions.Item>
              <Descriptions.Item label="年龄">{pet.age || '-'}</Descriptions.Item>
              <Descriptions.Item label="性别">{genderMap[pet.gender] || pet.gender}</Descriptions.Item>
              <Descriptions.Item label="体重">{pet.weight ? `${pet.weight}kg` : '-'}</Descriptions.Item>
              <Descriptions.Item label="颜色">{pet.color || '-'}</Descriptions.Item>
              <Descriptions.Item label="已接种疫苗">{pet.vaccinated ? '是' : '否'}</Descriptions.Item>
              <Descriptions.Item label="已绝育">{pet.neutered ? '是' : '否'}</Descriptions.Item>
            </Descriptions>

            {pet.description && (
              <div style={{ marginTop: 16 }}>
                <h4>简介</h4>
                <p>{pet.description}</p>
              </div>
            )}

            {pet.personality && (
              <div style={{ marginTop: 16 }}>
                <h4>性格特点</h4>
                <p>{pet.personality}</p>
              </div>
            )}

            {pet.special_needs && (
              <div style={{ marginTop: 16 }}>
                <h4>特殊需求</h4>
                <p>{pet.special_needs}</p>
              </div>
            )}

            <Space style={{ marginTop: 24 }}>
              {canAdopt && (
                <Button
                  type="primary"
                  icon={<HeartOutlined />}
                  onClick={() => setAdoptionModal(true)}
                >
                  申请领养
                </Button>
              )}
              {user?.role === 'adopter' && (
                <Button
                  icon={<CalendarOutlined />}
                  onClick={() => setAppointmentModal(true)}
                >
                  预约看望
                </Button>
              )}
              {isOwner && (
                <Button
                  icon={<EditOutlined />}
                  onClick={() => {
                    editForm.setFieldsValue(pet)
                    setEditModal(true)
                  }}
                >
                  编辑信息
                </Button>
              )}
            </Space>
          </div>
        </div>
      </Card>

      <Tabs
        items={[
          {
            key: 'health',
            label: '健康档案',
            children: (
              <div>
                {healthSummary && (
                  <Row gutter={16} style={{ marginBottom: 16 }}>
                    <Col span={6}>
                      <Card size="small">
                        <Statistic title="健康记录" value={healthSummary.total_records} />
                      </Card>
                    </Col>
                    <Col span={6}>
                      <Card size="small">
                        <Statistic title="疫苗接种" value={healthSummary.vaccines} />
                      </Card>
                    </Col>
                    <Col span={6}>
                      <Card size="small">
                        <Statistic title="体检次数" value={healthSummary.checkups} />
                      </Card>
                    </Col>
                    <Col span={6}>
                      <Card size="small">
                        <Statistic title="驱虫次数" value={healthSummary.dewormings} />
                      </Card>
                    </Col>
                  </Row>
                )}

                {reminders.length > 0 && (
                  <Card title="待办提醒" size="small" style={{ marginBottom: 16 }}>
                    <Timeline
                      items={reminders.map((r) => ({
                        color: 'blue',
                        children: (
                          <div>
                            <strong>{r.title}</strong> - {r.reminder_date}
                          </div>
                        ),
                      }))}
                    />
                  </Card>
                )}

                <Card title="健康记录时间线">
                  {healthRecords.length > 0 ? (
                    <Timeline
                      className="health-timeline"
                      items={healthRecords.map((record) => ({
                        color: record.record_type === 'vaccine' ? 'green' :
                               record.record_type === 'disease' ? 'red' : 'blue',
                        children: (
                          <div>
                            <div>
                              <strong>{record.title}</strong>
                              <Tag style={{ marginLeft: 8 }}>{record.record_type}</Tag>
                            </div>
                            <div style={{ color: '#999', fontSize: 12 }}>
                              {record.record_date} | {record.vet_name || '-'} | {record.hospital || '-'}
                            </div>
                            {record.description && <p style={{ marginTop: 8 }}>{record.description}</p>}
                          </div>
                        ),
                      }))}
                    />
                  ) : (
                    <Empty description="暂无健康记录" />
                  )}
                </Card>
              </div>
            ),
          },
          {
            key: 'history',
            label: '领养历史',
            children: (
              <Card>
                {adoptionHistory.length > 0 ? (
                  <Table
                    dataSource={adoptionHistory}
                    rowKey="id"
                    columns={[
                      { title: '申请时间', dataIndex: 'created_at', key: 'created_at' },
                      { title: '领养人', key: 'adopter', render: (_, r) => r.adopter?.name || '-' },
                      {
                        title: '状态',
                        dataIndex: 'status',
                        key: 'status',
                        render: (status) => <Tag>{status}</Tag>,
                      },
                    ]}
                  />
                ) : (
                  <Empty description="暂无领养记录" />
                )}
              </Card>
            ),
          },
        ]}
      />

      <Modal
        title="申请领养"
        open={adoptionModal}
        onCancel={() => setAdoptionModal(false)}
        footer={null}
        width={600}
      >
        <Form form={adoptionForm} onFinish={handleAdoptionSubmit} layout="vertical">
          <Form.Item name="reason" label="领养原因" rules={[{ required: true }]}>
            <TextArea rows={3} placeholder="请说明您想领养这只宠物的原因" />
          </Form.Item>
          <Form.Item name="living_situation" label="居住情况" rules={[{ required: true }]}>
            <TextArea rows={2} placeholder="请描述您的居住环境" />
          </Form.Item>
          <Form.Item name="housing_type" label="住房类型" rules={[{ required: true }]}>
            <Select>
              <Option value="apartment">公寓</Option>
              <Option value="house">独栋</Option>
              <Option value="villa">别墅</Option>
              <Option value="other">其他</Option>
            </Select>
          </Form.Item>
          <Form.Item name="pet_experience" label="养宠经验">
            <TextArea rows={2} placeholder="请描述您的养宠经验" />
          </Form.Item>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="family_members" label="家庭成员数量">
                <Input type="number" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="income_level" label="收入水平">
                <Select>
                  <Option value="low">一般</Option>
                  <Option value="medium">中等</Option>
                  <Option value="high">较高</Option>
                </Select>
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={16}>
            <Col span={8}>
              <Form.Item name="has_children" label="有小孩" valuePropName="checked">
                <Checkbox />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item name="has_other_pets" label="有其他宠物" valuePropName="checked">
                <Checkbox />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item name="can_afford_vet" label="可承担医疗费用" valuePropName="checked" initialValue={true}>
                <Checkbox />
              </Form.Item>
            </Col>
          </Row>
          <Form.Item name="other_pets_desc" label="其他宠物描述">
            <TextArea rows={2} placeholder="如有其他宠物，请描述种类和数量" />
          </Form.Item>
          <Form.Item name="agree_to_visit" label="同意定期回访" rules={[{ required: true }]} valuePropName="checked">
            <Checkbox />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={submitting} block>
              提交申请
            </Button>
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title="预约看望"
        open={appointmentModal}
        onCancel={() => setAppointmentModal(false)}
        footer={null}
      >
        <Form form={appointmentForm} onFinish={handleAppointmentSubmit} layout="vertical">
          <Form.Item name="appointment_date" label="预约日期" rules={[{ required: true }]}>
            <DatePicker style={{ width: '100%' }} />
          </Form.Item>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="start_time" label="开始时间" rules={[{ required: true }]}>
                <TimePicker format="HH:mm" style={{ width: '100%' }} />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="end_time" label="结束时间" rules={[{ required: true }]}>
                <TimePicker format="HH:mm" style={{ width: '100%' }} />
              </Form.Item>
            </Col>
          </Row>
          <Form.Item name="notes" label="备注">
            <TextArea rows={3} />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={submitting} block>
              提交预约
            </Button>
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title="编辑宠物信息"
        open={editModal}
        onCancel={() => setEditModal(false)}
        footer={null}
        width={600}
      >
        <Form form={editForm} onFinish={handleEditSubmit} layout="vertical">
          <Form.Item name="name" label="名称" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="species" label="物种" rules={[{ required: true }]}>
                <Select>
                  <Option value="dog">狗</Option>
                  <Option value="cat">猫</Option>
                  <Option value="rabbit">兔子</Option>
                  <Option value="bird">鸟</Option>
                  <Option value="other">其他</Option>
                </Select>
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="breed" label="品种">
                <Input />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={16}>
            <Col span={8}>
              <Form.Item name="age" label="年龄">
                <Input />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item name="gender" label="性别" rules={[{ required: true }]}>
                <Select>
                  <Option value="male">公</Option>
                  <Option value="female">母</Option>
                  <Option value="unknown">未知</Option>
                </Select>
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item name="weight" label="体重(kg)">
                <Input type="number" />
              </Form.Item>
            </Col>
          </Row>
          <Form.Item name="description" label="描述">
            <TextArea rows={3} />
          </Form.Item>
          <Form.Item name="personality" label="性格特点">
            <TextArea rows={2} />
          </Form.Item>
          <Form.Item name="special_needs" label="特殊需求">
            <TextArea rows={2} />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={submitting} block>
              保存
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default PetDetail
