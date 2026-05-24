import { useState } from 'react'
import { Form, Input, Button, Card, message, Upload } from 'antd'
import { SafetyOutlined, UserOutlined, IdcardOutlined } from '@ant-design/icons'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { userApi } from '@/api/endpoints'
import { useAuthStore } from '@/store/authStore'

export default function VerifyPage() {
  const { user } = useAuthStore()
  const queryClient = useQueryClient()
  const [countdown, setCountdown] = useState(0)
  const [form] = Form.useForm()

  const sendCode = async () => {
    const phone = form.getFieldValue('phone')
    if (!phone) {
      message.warning('请先输入手机号')
      return
    }
    try {
      await userApi.sendSmsCode(phone)
      message.success('验证码已发送')
      setCountdown(60)
      const timer = setInterval(() => {
        setCountdown((c) => {
          if (c <= 1) {
            clearInterval(timer)
            return 0
          }
          return c - 1
        })
      }, 1000)
    } catch {
      // handled
    }
  }

  const verifyMutation = useMutation({
    mutationFn: userApi.verify,
    onSuccess: () => {
      message.success('认证信息已提交，等待审核')
      queryClient.invalidateQueries({ queryKey: ['userInfo'] })
    },
  })

  const onFinish = (values: any) => {
    const files = values.id_card_files || []
    verifyMutation.mutate({
      real_name: values.real_name,
      id_card: values.id_card,
      id_card_front: files[0]?.url || '',
      id_card_back: files[1]?.url || '',
      phone: values.phone,
      sms_code: values.sms_code,
    })
  }

  if (user?.verify_status === 'verified') {
    return (
      <Card>
        <div style={{ textAlign: 'center', padding: 40 }}>
          <SafetyOutlined style={{ fontSize: 64, color: '#52c41a' }} />
          <h2 style={{ marginTop: 16 }}>您已完成实名认证</h2>
          <p style={{ color: '#888' }}>感谢您的配合，现在可以自由使用平台功能了</p>
        </div>
      </Card>
    )
  }

  return (
    <Card title="实名认证">
      <p style={{ color: '#888', marginBottom: 24 }}>
        完成实名认证后才能与其他用户互动，请如实填写以下信息
      </p>
      <Form form={form} onFinish={onFinish} layout="vertical" style={{ maxWidth: 500 }}>
        <Form.Item name="real_name" label="真实姓名" rules={[{ required: true, message: '请输入真实姓名' }]}>
          <Input prefix={<UserOutlined />} placeholder="请输入真实姓名" />
        </Form.Item>
        <Form.Item name="id_card" label="身份证号" rules={[{ required: true, message: '请输入身份证号' }]}>
          <Input prefix={<IdcardOutlined />} placeholder="请输入身份证号" />
        </Form.Item>
        <Form.Item name="id_card_files" label="身份证照片" rules={[{ required: true, message: '请上传身份证照片' }]}>
          <Upload
            listType="picture"
            maxCount={2}
            beforeUpload={() => false}
          >
            <Button>上传身份证（正反面）</Button>
          </Upload>
        </Form.Item>
        <Form.Item name="phone" label="手机号" rules={[{ required: true, message: '请输入手机号' }]}>
          <Input placeholder="请输入手机号" />
        </Form.Item>
        <Form.Item name="sms_code" label="短信验证码" rules={[{ required: true, message: '请输入验证码' }]}>
          <Input
            placeholder="短信验证码"
            suffix={
              <Button size="small" onClick={sendCode} disabled={countdown > 0}>
                {countdown > 0 ? `${countdown}s` : '获取验证码'}
              </Button>
            }
          />
        </Form.Item>
        <Form.Item>
          <Button type="primary" htmlType="submit" loading={verifyMutation.isPending}>
            提交认证
          </Button>
        </Form.Item>
      </Form>
    </Card>
  )
}
