import React, { useState, useEffect } from 'react'
import { Card, Form, Input, Button, Avatar, Upload, message, Row, Col } from 'antd'
import { UploadOutlined } from '@ant-design/icons'
import { authApi } from '@/services/auth'
import { useAppSelector, useAppDispatch } from '@/store/hooks'
import { setUserInfo } from '@/store'

const Profile: React.FC = () => {
  const { userInfo } = useAppSelector((state) => state.auth)
  const dispatch = useAppDispatch()
  const [form] = Form.useForm()
  const [loading, setLoading] = useState(false)
  const [passwordForm] = Form.useForm()
  const [changingPassword, setChangingPassword] = useState(false)

  useEffect(() => {
    if (userInfo) {
      form.setFieldsValue({
        nickname: userInfo.nickname,
        phone: userInfo.phone,
      })
    }
  }, [userInfo])

  const handleUpdateProfile = async () => {
    try {
      const values = await form.validateFields()
      setLoading(true)
      const res = await authApi.updateProfile(values)
      dispatch(setUserInfo(res))
      message.success('更新成功')
    } catch (error) {
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  const handleChangePassword = async () => {
    try {
      const values = await passwordForm.validateFields()
      if (values.new_password !== values.confirm_password) {
        message.error('两次输入的密码不一致')
        return
      }
      setChangingPassword(true)
      await authApi.changePassword(values)
      message.success('密码修改成功')
      passwordForm.resetFields()
    } catch (error) {
      console.error(error)
    } finally {
      setChangingPassword(false)
    }
  }

  return (
    <div>
      <div className="page-header">
        <h1 className="page-title">个人中心</h1>
      </div>

      <Row gutter={24}>
        <Col span={8}>
          <Card style={{ textAlign: 'center' }}>
            <Avatar size={100} src={userInfo?.avatar}>
              {userInfo?.nickname?.charAt(0)}
            </Avatar>
            <h2 style={{ marginTop: 16 }}>{userInfo?.nickname}</h2>
            <p style={{ color: '#999' }}>
              {userInfo?.role === 'customer' && '客户'}
              {userInfo?.role === 'service_provider' && '服务人员'}
              {userInfo?.role === 'admin' && '管理员'}
            </p>
            <Upload
              name="avatar"
              showUploadList={false}
              action="/api/upload/avatar"
              onChange={async (info) => {
                if (info.file.status === 'done') {
                  message.success('头像上传成功')
                  try {
                    const res = await authApi.getCurrentUser()
                    dispatch(setUserInfo(res))
                    localStorage.setItem('userInfo', JSON.stringify(res))
                  } catch (error) {
                    console.error(error)
                  }
                }
              }}
            >
              <Button icon={<UploadOutlined />}>更换头像</Button>
            </Upload>
          </Card>
        </Col>

        <Col span={16}>
          <Card title="基本信息" style={{ marginBottom: 24 }}>
            <Form form={form} layout="vertical">
              <Form.Item
                name="nickname"
                label="昵称"
                rules={[{ required: true, message: '请输入昵称' }]}
              >
                <Input placeholder="请输入昵称" />
              </Form.Item>
              <Form.Item name="phone" label="手机号">
                <Input disabled />
              </Form.Item>
              <Form.Item>
                <Button
                  type="primary"
                  loading={loading}
                  onClick={handleUpdateProfile}
                >
                  保存修改
                </Button>
              </Form.Item>
            </Form>
          </Card>

          <Card title="修改密码">
            <Form form={passwordForm} layout="vertical">
              <Form.Item
                name="old_password"
                label="当前密码"
                rules={[{ required: true, message: '请输入当前密码' }]}
              >
                <Input.Password placeholder="请输入当前密码" />
              </Form.Item>
              <Form.Item
                name="new_password"
                label="新密码"
                rules={[{ required: true, message: '请输入新密码' }]}
              >
                <Input.Password placeholder="请输入新密码" />
              </Form.Item>
              <Form.Item
                name="confirm_password"
                label="确认新密码"
                rules={[{ required: true, message: '请再次输入新密码' }]}
              >
                <Input.Password placeholder="请再次输入新密码" />
              </Form.Item>
              <Form.Item>
                <Button
                  type="primary"
                  loading={changingPassword}
                  onClick={handleChangePassword}
                >
                  修改密码
                </Button>
              </Form.Item>
            </Form>
          </Card>
        </Col>
      </Row>
    </div>
  )
}

export default Profile
