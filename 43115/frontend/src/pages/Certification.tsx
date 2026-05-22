import React, { useState, useEffect } from 'react'
import { Card, Form, Input, Button, Upload, message, Select, Row, Col, Alert } from 'antd'
import { UploadOutlined } from '@ant-design/icons'
import { certificationApi } from '@/services/auth'
import { useAppSelector } from '@/store/hooks'

const Certification: React.FC = () => {
  const { userInfo } = useAppSelector((state) => state.auth)
  const [form] = Form.useForm()
  const [loading, setLoading] = useState(false)
  const [status, setStatus] = useState<'none' | 'pending' | 'approved' | 'rejected'>('none')

  useEffect(() => {
    if (userInfo?.provider_status) {
      setStatus(userInfo.provider_status as any)
    }
  }, [userInfo])

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      setLoading(true)
      await certificationApi.submit(values)
      message.success('认证提交成功，请等待审核')
      setStatus('pending')
    } catch (error) {
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  const renderStatus = () => {
    switch (status) {
      case 'pending':
        return <Alert message="认证审核中，请耐心等待" type="warning" showIcon />
      case 'approved':
        return <Alert message="认证已通过" type="success" showIcon />
      case 'rejected':
        return (
          <Alert
            message={`认证被拒绝：${userInfo?.reject_reason || '请重新提交'}`}
            type="error"
            showIcon
          />
        )
      default:
        return null
    }
  }

  return (
    <div>
      <div className="page-header">
        <h1 className="page-title">服务人员认证</h1>
      </div>

      {renderStatus()}

      <Card title="身份认证" style={{ marginBottom: 24 }}>
        <Form form={form} layout="vertical">
          <Row gutter={24}>
            <Col span={12}>
              <Form.Item
                name="real_name"
                label="真实姓名"
                rules={[{ required: true, message: '请输入真实姓名' }]}
              >
                <Input placeholder="请输入真实姓名" disabled={status === 'pending' || status === 'approved'} />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                name="id_card"
                label="身份证号"
                rules={[{ required: true, message: '请输入身份证号' }]}
              >
                <Input placeholder="请输入身份证号" disabled={status === 'pending' || status === 'approved'} />
              </Form.Item>
            </Col>
          </Row>
          <Form.Item
            name="id_card_front"
            label="身份证正面"
            rules={[{ required: true, message: '请上传身份证正面' }]}
          >
            <Upload
              name="file"
              action="/api/upload"
              listType="picture"
              maxCount={1}
              disabled={status === 'pending' || status === 'approved'}
            >
              <Button icon={<UploadOutlined />}>上传身份证正面</Button>
            </Upload>
          </Form.Item>
          <Form.Item
            name="id_card_back"
            label="身份证背面"
            rules={[{ required: true, message: '请上传身份证背面' }]}
          >
            <Upload
              name="file"
              action="/api/upload"
              listType="picture"
              maxCount={1}
              disabled={status === 'pending' || status === 'approved'}
            >
              <Button icon={<UploadOutlined />}>上传身份证背面</Button>
            </Upload>
          </Form.Item>
        </Form>
      </Card>

      <Card title="技能证书">
        <Form form={form} layout="vertical">
          <Form.Item name="certifications" label="技能证书列表">
            <Upload
              name="file"
              action="/api/upload"
              listType="picture-card"
              multiple
              disabled={status === 'pending' || status === 'approved'}
            >
              <div>
                <UploadOutlined />
                <div style={{ marginTop: 8 }}>上传证书</div>
              </div>
            </Upload>
          </Form.Item>
          <Form.Item
            name="skills"
            label="服务技能"
            rules={[{ required: true, message: '请选择服务技能' }]}
          >
            <Select
              mode="multiple"
              placeholder="请选择服务技能"
              disabled={status === 'pending' || status === 'approved'}
              options={[
                { label: '保洁', value: 'cleaning' },
                { label: '月嫂', value: 'nanny' },
                { label: '护工', value: 'caregiver' },
                { label: '育婴师', value: 'babysitter' },
                { label: '厨师', value: 'cook' },
                { label: '维修', value: 'repair' },
              ]}
            />
          </Form.Item>
          <Form.Item
            name="experience"
            label="工作经验"
            rules={[{ required: true, message: '请输入工作经验' }]}
          >
            <Input.TextArea
              rows={4}
              placeholder="请描述您的工作经验"
              disabled={status === 'pending' || status === 'approved'}
            />
          </Form.Item>
          {status !== 'pending' && status !== 'approved' && (
            <Form.Item>
              <Button
                type="primary"
                size="large"
                loading={loading}
                onClick={handleSubmit}
              >
                提交认证
              </Button>
            </Form.Item>
          )}
        </Form>
      </Card>
    </div>
  )
}

export default Certification
