import React, { useState, useEffect } from 'react'
import {
  Card,
  Form,
  Input,
  Button,
  Space,
  Typography,
  Descriptions,
  message,
  Spin,
  Tag,
  Divider,
  Row,
  Col
} from 'antd'
import { EditOutlined, SaveOutlined, CloseOutlined, UserOutlined } from '@ant-design/icons'
import type { HealthRecord, User } from '@/types'
import { healthRecordAPI, authAPI } from '@/services/api'

const { Title, Text } = Typography
const { TextArea } = Input

interface HealthRecordFormData {
  medical_history: string
  allergies: string
  medications: string
  vaccinations: string
  family_history: string
  life_habits: string
  remarks: string
}

const HealthRecordsPage: React.FC = () => {
  const [form] = Form.useForm<HealthRecordFormData>()
  const [loading, setLoading] = useState(false)
  const [saving, setSaving] = useState(false)
  const [editing, setEditing] = useState(false)
  const [healthRecord, setHealthRecord] = useState<HealthRecord | null>(null)
  const [currentUser, setCurrentUser] = useState<User | null>(null)

  useEffect(() => {
    fetchCurrentUser()
  }, [])

  const fetchCurrentUser = async () => {
    try {
      const user = await authAPI.getCurrentUser()
      setCurrentUser(user)
      if (user.id) {
        fetchHealthRecord(user.id)
      }
    } catch (error) {
      console.error('获取用户信息失败:', error)
    }
  }

  const fetchHealthRecord = async (userId: number) => {
    setLoading(true)
    try {
      const record = await healthRecordAPI.get(userId)
      setHealthRecord(record)
      form.setFieldsValue({
        medical_history: record.medical_history as unknown as string || '',
        allergies: record.allergies as unknown as string || '',
        medications: record.medications as unknown as string || '',
        vaccinations: record.vaccinations as unknown as string || '',
        family_history: record.family_history || '',
        life_habits: record.life_habits || '',
        remarks: record.remarks || ''
      })
    } catch (error) {
      console.error('获取健康档案失败:', error)
      message.error('获取健康档案失败')
    } finally {
      setLoading(false)
    }
  }

  const handleEdit = () => {
    setEditing(true)
  }

  const handleCancel = () => {
    setEditing(false)
    if (healthRecord) {
      form.setFieldsValue({
        medical_history: healthRecord.medical_history as unknown as string || '',
        allergies: healthRecord.allergies as unknown as string || '',
        medications: healthRecord.medications as unknown as string || '',
        vaccinations: healthRecord.vaccinations as unknown as string || '',
        family_history: healthRecord.family_history || '',
        life_habits: healthRecord.life_habits || '',
        remarks: healthRecord.remarks || ''
      })
    }
  }

  const handleSave = async (values: HealthRecordFormData) => {
    if (!currentUser) return

    setSaving(true)
    try {
      const data: Partial<HealthRecord> = {
        medical_history: values.medical_history as unknown as Record<string, unknown>,
        allergies: values.allergies as unknown as Record<string, unknown>,
        medications: values.medications as unknown as Record<string, unknown>,
        vaccinations: values.vaccinations as unknown as Record<string, unknown>,
        family_history: values.family_history,
        life_habits: values.life_habits,
        remarks: values.remarks
      }
      const updatedRecord = await healthRecordAPI.update(currentUser.id, data)
      setHealthRecord(updatedRecord)
      setEditing(false)
      message.success('健康档案保存成功')
    } catch (error) {
      console.error('保存健康档案失败:', error)
      message.error('保存健康档案失败')
    } finally {
      setSaving(false)
    }
  }

  const getAge = (birthDate: string | null) => {
    if (!birthDate) return '-'
    const today = new Date()
    const birth = new Date(birthDate)
    let age = today.getFullYear() - birth.getFullYear()
    const monthDiff = today.getMonth() - birth.getMonth()
    if (monthDiff < 0 || (monthDiff === 0 && today.getDate() < birth.getDate())) {
      age--
    }
    return age
  }

  if (loading) {
    return (
      <div className="flex justify-center items-center min-h-[400px]">
        <Spin size="large" />
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <Card>
        <div className="flex items-center justify-between mb-4">
          <Title level={3} style={{ margin: 0 }}>
            <UserOutlined className="mr-2" />
            健康档案
          </Title>
          {!editing && (
            <Button type="primary" icon={<EditOutlined />} onClick={handleEdit}>
              编辑
            </Button>
          )}
        </div>

        <Divider />

        <Title level={5} style={{ marginBottom: 16 }}>患者基本信息</Title>
        <Descriptions bordered column={{ xs: 1, sm: 2, md: 3 }} size="middle">
          <Descriptions.Item label="姓名">{currentUser?.full_name || '-'}</Descriptions.Item>
          <Descriptions.Item label="性别">{currentUser?.gender || '-'}</Descriptions.Item>
          <Descriptions.Item label="年龄">{getAge(currentUser?.birth_date || null)}</Descriptions.Item>
          <Descriptions.Item label="手机号">{currentUser?.phone || '-'}</Descriptions.Item>
          <Descriptions.Item label="邮箱">{currentUser?.email || '-'}</Descriptions.Item>
          <Descriptions.Item label="档案更新时间">
            {healthRecord?.updated_at ? new Date(healthRecord.updated_at).toLocaleString() : '-'}
          </Descriptions.Item>
        </Descriptions>
      </Card>

      <Card>
        <div className="flex items-center justify-between mb-4">
          <Title level={5} style={{ margin: 0 }}>健康信息</Title>
          {editing && (
            <Space>
              <Button icon={<CloseOutlined />} onClick={handleCancel}>
                取消
              </Button>
              <Button type="primary" icon={<SaveOutlined />} loading={saving} onClick={() => form.submit()}>
                保存
              </Button>
            </Space>
          )}
        </div>

        <Form
          form={form}
          layout="vertical"
          onFinish={handleSave}
          disabled={!editing}
        >
          <Row gutter={24}>
            <Col xs={24} md={12}>
              <Form.Item
                name="medical_history"
                label="既往病史"
                rules={[{ required: true, message: '请输入既往病史' }]}
              >
                <TextArea rows={4} placeholder="请输入既往病史，如：高血压、糖尿病等" />
              </Form.Item>
            </Col>
            <Col xs={24} md={12}>
              <Form.Item
                name="allergies"
                label="过敏信息"
                rules={[{ required: true, message: '请输入过敏信息' }]}
              >
                <TextArea rows={4} placeholder="请输入过敏信息，如：青霉素、花粉等" />
              </Form.Item>
            </Col>
          </Row>

          <Row gutter={24}>
            <Col xs={24} md={12}>
              <Form.Item
                name="medications"
                label="常用药物"
                rules={[{ required: true, message: '请输入常用药物' }]}
              >
                <TextArea rows={4} placeholder="请输入常用药物，如：阿司匹林、二甲双胍等" />
              </Form.Item>
            </Col>
            <Col xs={24} md={12}>
              <Form.Item
                name="vaccinations"
                label="疫苗接种"
                rules={[{ required: true, message: '请输入疫苗接种情况' }]}
              >
                <TextArea rows={4} placeholder="请输入疫苗接种情况，如：新冠疫苗、乙肝疫苗等" />
              </Form.Item>
            </Col>
          </Row>

          <Row gutter={24}>
            <Col xs={24} md={12}>
              <Form.Item
                name="family_history"
                label="家族病史"
              >
                <TextArea rows={3} placeholder="请输入家族病史" />
              </Form.Item>
            </Col>
            <Col xs={24} md={12}>
              <Form.Item
                name="life_habits"
                label="生活习惯"
              >
                <TextArea rows={3} placeholder="请输入生活习惯，如：吸烟、饮酒、运动等" />
              </Form.Item>
            </Col>
          </Row>

          <Form.Item
            name="remarks"
            label="备注"
          >
            <TextArea rows={3} placeholder="其他需要说明的信息" />
          </Form.Item>
        </Form>
      </Card>

      {!editing && healthRecord && (
        <Card>
          <Title level={5} style={{ marginBottom: 16 }}>健康概览</Title>
          <Row gutter={[16, 16]}>
            <Col xs={12} sm={6}>
              <Card size="small" className="text-center">
                <Text type="secondary">既往病史</Text>
                <div className="mt-2">
                  {healthRecord.medical_history ? (
                    <Tag color="orange">已填写</Tag>
                  ) : (
                    <Tag color="default">未填写</Tag>
                  )}
                </div>
              </Card>
            </Col>
            <Col xs={12} sm={6}>
              <Card size="small" className="text-center">
                <Text type="secondary">过敏信息</Text>
                <div className="mt-2">
                  {healthRecord.allergies ? (
                    <Tag color="red">已填写</Tag>
                  ) : (
                    <Tag color="default">未填写</Tag>
                  )}
                </div>
              </Card>
            </Col>
            <Col xs={12} sm={6}>
              <Card size="small" className="text-center">
                <Text type="secondary">常用药物</Text>
                <div className="mt-2">
                  {healthRecord.medications ? (
                    <Tag color="blue">已填写</Tag>
                  ) : (
                    <Tag color="default">未填写</Tag>
                  )}
                </div>
              </Card>
            </Col>
            <Col xs={12} sm={6}>
              <Card size="small" className="text-center">
                <Text type="secondary">疫苗接种</Text>
                <div className="mt-2">
                  {healthRecord.vaccinations ? (
                    <Tag color="green">已填写</Tag>
                  ) : (
                    <Tag color="default">未填写</Tag>
                  )}
                </div>
              </Card>
            </Col>
          </Row>
        </Card>
      )}
    </div>
  )
}

export default HealthRecordsPage
