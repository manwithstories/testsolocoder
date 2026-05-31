import { useState } from 'react'
import { Card, Form, Input, Button, Avatar, Upload, message, Tabs, Divider } from 'antd'
import { UserOutlined, UploadOutlined } from '@ant-design/icons'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { authApi, uploadApi } from '@/services/api'
import { useAuthStore } from '@/context/AuthContext'
import type { User } from '@/types'

export default function Profile() {
  const queryClient = useQueryClient()
  const { user, setUser } = useAuthStore()
  const [profileForm] = Form.useForm()
  const [passwordForm] = Form.useForm()
  const [storeForm] = Form.useForm()
  const [keeperForm] = Form.useForm()

  const { data: profileData } = useQuery({
    queryKey: ['profile'],
    queryFn: () => authApi.getProfile(),
  })

  const profile: User | undefined = profileData?.data

  const updateProfileMutation = useMutation({
    mutationFn: (values: any) => authApi.updateProfile(values),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['profile'] })
      message.success('更新成功')
    },
    onError: (err: any) => message.error(err.message || '更新失败'),
  })

  const changePasswordMutation = useMutation({
    mutationFn: (values: any) => authApi.changePassword(values),
    onSuccess: () => {
      message.success('密码修改成功')
      passwordForm.resetFields()
    },
    onError: (err: any) => message.error(err.message || '修改失败'),
  })

  const updateStoreMutation = useMutation({
    mutationFn: (values: any) => authApi.updateStoreInfo(values),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['profile'] })
      message.success('门店信息更新成功')
    },
    onError: (err: any) => message.error(err.message || '更新失败'),
  })

  const updateKeeperMutation = useMutation({
    mutationFn: (values: any) => authApi.updateKeeperInfo(values),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['profile'] })
      message.success('管家信息更新成功')
    },
    onError: (err: any) => message.error(err.message || '更新失败'),
  })

  const handleAvatarUpload = async (file: File) => {
    const formData = new FormData()
    formData.append('file', file)
    try {
      const res = await uploadApi.upload(formData)
      profileForm.setFieldsValue({ avatar_url: res.data.url })
      message.success('头像上传成功')
    } catch (err: any) {
      message.error(err.message || '上传失败')
    }
    return false
  }

  if (!profile) {
    return <div className="text-center py-10">加载中...</div>
  }

  return (
    <div className="space-y-4">
      <Card>
        <div className="flex items-center gap-6">
          <Avatar size={80} icon={<UserOutlined />} src={profile.avatar_url}>
            {profile.username.charAt(0).toUpperCase()}
          </Avatar>
          <div>
            <h2 className="text-xl font-semibold m-0">{profile.real_name || profile.username}</h2>
            <p className="text-gray-500 m-0">
              {profile.role === 'owner' && '宠物主人'}
              {profile.role === 'store' && '寄养门店'}
              {profile.role === 'keeper' && '宠物管家'}
            </p>
            <p className="text-gray-400 m-0 text-sm">{profile.email}</p>
          </div>
        </div>
      </Card>

      <Card title="个人信息">
        <Form
          form={profileForm}
          layout="vertical"
          initialValues={{
            phone: profile.phone,
            avatar_url: profile.avatar_url,
            real_name: profile.real_name,
          }}
          onFinish={(values) => updateProfileMutation.mutate(values)}
        >
          <div className="grid grid-cols-2 gap-4">
            <Form.Item name="phone" label="手机号">
              <Input placeholder="手机号" />
            </Form.Item>
            <Form.Item name="real_name" label="真实姓名">
              <Input placeholder="真实姓名" />
            </Form.Item>
          </div>
          <Form.Item name="avatar_url" label="头像URL">
            <div className="flex gap-3 items-center">
              <Input placeholder="头像URL" />
              <Upload beforeUpload={handleAvatarUpload} showUploadList={false}>
                <Button icon={<UploadOutlined />}>上传</Button>
              </Upload>
            </div>
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={updateProfileMutation.isPending}>
              保存修改
            </Button>
          </Form.Item>
        </Form>
      </Card>

      <Card title="修改密码">
        <Form
          form={passwordForm}
          layout="vertical"
          onFinish={(values) => changePasswordMutation.mutate(values)}
        >
          <div className="grid grid-cols-2 gap-4">
            <Form.Item name="old_password" label="当前密码" rules={[{ required: true }]}>
              <Input.Password placeholder="当前密码" />
            </Form.Item>
            <Form.Item name="new_password" label="新密码" rules={[{ required: true, min: 6 }]}>
              <Input.Password placeholder="新密码（至少6位）" />
            </Form.Item>
          </div>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={changePasswordMutation.isPending}>
              修改密码
            </Button>
          </Form.Item>
        </Form>
      </Card>

      {profile.role === 'store' && (
        <Card title="门店信息">
          <Form
            form={storeForm}
            layout="vertical"
            initialValues={profile.store_info}
            onFinish={(values) => updateStoreMutation.mutate(values)}
          >
            <div className="grid grid-cols-2 gap-4">
              <Form.Item name="store_name" label="门店名称" rules={[{ required: true }]}>
                <Input placeholder="门店名称" />
              </Form.Item>
              <Form.Item name="license_no" label="营业执照号">
                <Input placeholder="营业执照号" />
              </Form.Item>
            </div>
            <Form.Item name="address" label="地址">
              <Input placeholder="地址" />
            </Form.Item>
            <Form.Item name="business_hours" label="营业时间">
              <Input placeholder="如：09:00-21:00" />
            </Form.Item>
            <Form.Item name="description" label="门店描述">
              <Input.TextArea rows={3} placeholder="门店简介" />
            </Form.Item>
            <Form.Item>
              <Button type="primary" htmlType="submit" loading={updateStoreMutation.isPending}>
                保存修改
              </Button>
            </Form.Item>
          </Form>
        </Card>
      )}

      {profile.role === 'keeper' && (
        <Card title="管家信息">
          <Form
            form={keeperForm}
            layout="vertical"
            initialValues={profile.keeper_info}
            onFinish={(values) => updateKeeperMutation.mutate(values)}
          >
            <div className="grid grid-cols-2 gap-4">
              <Form.Item name="nick_name" label="昵称">
                <Input placeholder="昵称" />
              </Form.Item>
              <Form.Item name="experience" label="从业年限">
                <Input type="number" placeholder="从业年限" />
              </Form.Item>
            </div>
            <Form.Item name="specialty" label="专长">
              <Input placeholder="如：训犬、猫护理等" />
            </Form.Item>
            <Form.Item name="certification" label="证书">
              <Input placeholder="相关证书" />
            </Form.Item>
            <Form.Item>
              <Button type="primary" htmlType="submit" loading={updateKeeperMutation.isPending}>
                保存修改
              </Button>
            </Form.Item>
          </Form>
        </Card>
      )}
    </div>
  )
}
