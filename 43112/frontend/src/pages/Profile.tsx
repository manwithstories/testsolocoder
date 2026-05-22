import React, { useEffect, useState } from 'react'
import { Card, Form, Input, Button, Avatar, message, Upload, Descriptions, Tabs } from 'antd'
import { UserOutlined, UploadOutlined } from '@ant-design/icons'
import { authApi, uploadApi } from '@/services'
import { useAuthStore } from '@/store/auth'
import { User } from '@/types'

const ProfilePage: React.FC = () => {
  const { user, fetchProfile, setAuth } = useAuthStore()
  const [form] = Form.useForm()
  const [passwordForm] = Form.useForm()
  const [loading, setLoading] = useState(false)
  const [avatarUrl, setAvatarUrl] = useState<string>('')

  useEffect(() => {
    if (user) {
      form.setFieldsValue({
        nickname: user.nickname,
        phone: user.phone,
        bio: user.bio,
      })
      setAvatarUrl(user.avatar || '')
    }
  }, [user])

  const handleUpdateProfile = async (values: any) => {
    setLoading(true)
    try {
      const data = { ...values, avatar: avatarUrl }
      const res = await authApi.updateProfile(data)
      if (res.code === 0) {
        message.success('个人信息更新成功')
        await fetchProfile()
      }
    } catch (error: any) {
      message.error(error.message || '更新失败')
    } finally {
      setLoading(false)
    }
  }

  const handleChangePassword = async (values: any) => {
    setLoading(true)
    try {
      const res = await authApi.changePassword(values)
      if (res.code === 0) {
        message.success('密码修改成功')
        passwordForm.resetFields()
      }
    } catch (error: any) {
      message.error(error.message || '密码修改失败')
    } finally {
      setLoading(false)
    }
  }

  const handleAvatarUpload = async (file: File) => {
    try {
      const res = await uploadApi.upload(file, 'image')
      if (res.code === 0 && res.data) {
        setAvatarUrl(res.data.url)
        message.success('头像上传成功')
      }
    } catch (error: any) {
      message.error(error.message || '上传失败')
    }
    return false
  }

  const tabs = [
    {
      key: 'info',
      label: '基本信息',
      children: (
        <Card>
          <div style={{ textAlign: 'center', marginBottom: 24 }}>
            <Avatar size={96} icon={<UserOutlined />} src={avatarUrl} />
            <div style={{ marginTop: 8 }}>
              <Upload beforeUpload={handleAvatarUpload} showUploadList={false}>
                <Button icon={<UploadOutlined />}>更换头像</Button>
              </Upload>
            </div>
          </div>
          <Form form={form} onFinish={handleUpdateProfile} layout="vertical">
            <Form.Item name="nickname" label="昵称">
              <Input />
            </Form.Item>
            <Form.Item name="phone" label="手机号">
              <Input />
            </Form.Item>
            <Form.Item name="bio" label="个人简介">
              <Input.TextArea rows={3} />
            </Form.Item>
            <Form.Item>
              <Button type="primary" htmlType="submit" loading={loading}>
                保存修改
              </Button>
            </Form.Item>
          </Form>
        </Card>
      ),
    },
    {
      key: 'password',
      label: '修改密码',
      children: (
        <Card>
          <Form form={passwordForm} onFinish={handleChangePassword} layout="vertical">
            <Form.Item name="old_password" label="当前密码" rules={[{ required: true }]}>
              <Input.Password />
            </Form.Item>
            <Form.Item name="new_password" label="新密码" rules={[{ required: true, min: 6, message: '密码至少6位' }]}>
              <Input.Password />
            </Form.Item>
            <Form.Item>
              <Button type="primary" htmlType="submit" loading={loading}>
                修改密码
              </Button>
            </Form.Item>
          </Form>
        </Card>
      ),
    },
  ]

  if (!user) {
    return <div>加载中...</div>
  }

  return (
    <div>
      <h2>个人中心</h2>
      <Card style={{ marginBottom: 16 }}>
        <Descriptions column={2} bordered size="small">
          <Descriptions.Item label="用户名">{user.username}</Descriptions.Item>
          <Descriptions.Item label="邮箱">{user.email}</Descriptions.Item>
          <Descriptions.Item label="角色">
            {user.role === 'student' ? '学员' : user.role === 'instructor' ? '讲师' : '管理员'}
          </Descriptions.Item>
          <Descriptions.Item label="注册时间">{new Date(user.created_at).toLocaleDateString()}</Descriptions.Item>
        </Descriptions>
      </Card>
      <Tabs defaultActiveKey="info" items={tabs} />
    </div>
  )
}

export default ProfilePage
