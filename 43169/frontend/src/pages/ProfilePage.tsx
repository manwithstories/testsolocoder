import { useState, useRef } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { Form, Input, Select, DatePicker, InputNumber, Button, Card, Row, Col, Avatar, Upload, message, Tag, Divider } from 'antd'
import { PlusOutlined, UserOutlined } from '@ant-design/icons'
import dayjs from 'dayjs'
import { userApi } from '@/api/endpoints'
import { useAuthStore } from '@/store/authStore'

const { TextArea } = Input
const { Option } = Select

export default function ProfilePage() {
  const { user, setUser } = useAuthStore()
  const queryClient = useQueryClient()
  const [editing, setEditing] = useState(false)
  const [form] = Form.useForm()
  const avatarRef = useRef<HTMLInputElement>(null)

  const { data: userInfo } = useQuery({
    queryKey: ['userInfo'],
    queryFn: userApi.getUserInfo,
  })

  const profile = userInfo?.profile

  const updateMutation = useMutation({
    mutationFn: (values: any) => {
      const data = { ...values }
      if (data.birthday) {
        data.birthday = dayjs(data.birthday).format('YYYY-MM-DD')
      }
      return userApi.updateProfile(data)
    },
    onSuccess: () => {
      message.success('更新成功')
      queryClient.invalidateQueries({ queryKey: ['userInfo'] })
      setEditing(false)
    },
  })

  const handleAvatarClick = () => {
    avatarRef.current?.click()
  }

  const handleAvatarChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (!file) return

    try {
      const res = await userApi.uploadAvatar(file)
      const data = res as any
      if (data?.avatar) {
        setUser({ ...user!, avatar: data.avatar })
        message.success('头像更新成功')
      }
    } catch {
      // handled
    }
  }

  const handlePhotosUpload = async (info: any) => {
    if (info.fileList.length > 0) {
      try {
        const files = info.fileList.map((f: any) => f.originFileObj || f)
        await userApi.uploadPhotos(files)
        message.success('照片上传成功')
        queryClient.invalidateQueries({ queryKey: ['userInfo'] })
      } catch {
        // handled
      }
    }
  }

  const handleSubmit = (values: any) => {
    updateMutation.mutate(values)
  }

  return (
    <Row gutter={24}>
      <Col xs={24} md={8}>
        <Card>
          <div style={{ textAlign: 'center', marginBottom: 24 }}>
            <Avatar
              size={120}
              src={user?.avatar}
              icon={<UserOutlined />}
              style={{ cursor: 'pointer', backgroundColor: '#ff6b81' }}
              onClick={handleAvatarClick}
            />
            <input
              ref={avatarRef}
              type="file"
              accept="image/*"
              style={{ display: 'none' }}
              onChange={handleAvatarChange}
            />
            <h3 style={{ marginTop: 16 }}>{user?.username}</h3>
            <Tag color={user?.verify_status === 'verified' ? 'green' : 'orange'}>
              {user?.verify_status === 'verified' ? '已认证' : '未认证'}
            </Tag>
            <Tag color="gold">{user?.member_level}</Tag>
          </div>
        </Card>

        <Card title="我的照片" style={{ marginTop: 16 }}>
          <Upload
            listType="picture-card"
            multiple
            customRequest={handlePhotosUpload}
            maxCount={9}
          >
            <div>
              <PlusOutlined />
              <div style={{ marginTop: 8 }}>上传照片</div>
            </div>
          </Upload>
        </Card>
      </Col>

      <Col xs={24} md={16}>
        <Card
          title="个人资料"
          extra={
            <Button type="link" onClick={() => setEditing(!editing)}>
              {editing ? '取消' : '编辑'}
            </Button>
          }
        >
          {!editing ? (
            <div>
              <Row gutter={16}>
                <Col span={12}><p><strong>昵称：</strong>{profile?.nickname || '-'}</p></Col>
                <Col span={12}><p><strong>性别：</strong>{profile?.gender === 'male' ? '男' : profile?.gender === 'female' ? '女' : '-'}</p></Col>
                <Col span={12}><p><strong>年龄：</strong>{profile?.age || '-'}</p></Col>
                <Col span={12}><p><strong>身高：</strong>{profile?.height || '-'} cm</p></Col>
                <Col span={12}><p><strong>体重：</strong>{profile?.weight || '-'} kg</p></Col>
                <Col span={12}><p><strong>学历：</strong>{profile?.education || '-'}</p></Col>
                <Col span={12}><p><strong>职业：</strong>{profile?.occupation || '-'}</p></Col>
                <Col span={12}><p><strong>收入：</strong>{profile?.income || '-'}</p></Col>
                <Col span={12}><p><strong>城市：</strong>{profile?.city || '-'}</p></Col>
                <Col span={12}><p><strong>区域：</strong>{profile?.district || '-'}</p></Col>
              </Row>
              <Divider />
              <p><strong>个人介绍：</strong></p>
              <p>{profile?.intro || '暂无介绍'}</p>
              <p><strong>兴趣爱好：</strong></p>
              <p>{profile?.hobbies || '暂无'}</p>
              <p><strong>个人标签：</strong></p>
              <p>{profile?.tags || '暂无'}</p>
              <Divider />
              <h4>择偶标准</h4>
              <Row gutter={16}>
                <Col span={12}><p><strong>年龄范围：</strong>{profile?.min_age || 0} - {profile?.max_age || 0}</p></Col>
                <Col span={12}><p><strong>身高范围：</strong>{profile?.min_height || 0} - {profile?.max_height || 0} cm</p></Col>
                <Col span={12}><p><strong>学历要求：</strong>{profile?.prefer_education || '-'}</p></Col>
                <Col span={12}><p><strong>收入要求：</strong>{profile?.prefer_income || '-'}</p></Col>
                <Col span={12}><p><strong>城市偏好：</strong>{profile?.prefer_city || '-'}</p></Col>
              </Row>
            </div>
          ) : (
            <Form
              form={form}
              layout="vertical"
              initialValues={{
                nickname: profile?.nickname,
                gender: profile?.gender,
                birthday: profile?.birthday ? dayjs(profile.birthday) : null,
                height: profile?.height,
                weight: profile?.weight,
                education: profile?.education,
                occupation: profile?.occupation,
                income: profile?.income,
                city: profile?.city,
                district: profile?.district,
                address: profile?.address,
                intro: profile?.intro,
                hobbies: profile?.hobbies,
                tags: profile?.tags,
                min_age: profile?.min_age,
                max_age: profile?.max_age,
                min_height: profile?.min_height,
                max_height: profile?.max_height,
                prefer_education: profile?.prefer_education,
                prefer_income: profile?.prefer_income,
                prefer_city: profile?.prefer_city,
              }}
              onFinish={handleSubmit}
            >
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item name="nickname" label="昵称">
                    <Input />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="gender" label="性别">
                    <Select>
                      <Option value="male">男</Option>
                      <Option value="female">女</Option>
                    </Select>
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="birthday" label="生日">
                    <DatePicker style={{ width: '100%' }} />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="height" label="身高(cm)">
                    <InputNumber min={100} max={250} style={{ width: '100%' }} />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="weight" label="体重(kg)">
                    <InputNumber min={30} max={200} style={{ width: '100%' }} />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="education" label="学历">
                    <Select>
                      <Option value="high_school">高中</Option>
                      <Option value="college">大专</Option>
                      <Option value="bachelor">本科</Option>
                      <Option value="master">硕士</Option>
                      <Option value="phd">博士</Option>
                    </Select>
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="occupation" label="职业">
                    <Input />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="income" label="收入">
                    <Select>
                      <Option value="low">3万以下</Option>
                      <Option value="mid">3-10万</Option>
                      <Option value="high">10-30万</Option>
                      <Option value="luxury">30万以上</Option>
                    </Select>
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="city" label="城市">
                    <Input />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="district" label="区域">
                    <Input />
                  </Form.Item>
                </Col>
              </Row>
              <Form.Item name="intro" label="个人介绍">
                <TextArea rows={3} />
              </Form.Item>
              <Form.Item name="hobbies" label="兴趣爱好（逗号分隔）">
                <TextArea rows={2} />
              </Form.Item>
              <Form.Item name="tags" label="个人标签（逗号分隔）">
                <TextArea rows={2} />
              </Form.Item>
              <Divider>择偶标准</Divider>
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item name="min_age" label="最小年龄">
                    <InputNumber min={18} max={60} style={{ width: '100%' }} />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="max_age" label="最大年龄">
                    <InputNumber min={18} max={60} style={{ width: '100%' }} />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="min_height" label="最小身高">
                    <InputNumber min={100} max={250} style={{ width: '100%' }} />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="max_height" label="最大身高">
                    <InputNumber min={100} max={250} style={{ width: '100%' }} />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="prefer_education" label="学历要求">
                    <Select>
                      <Option value="">不限</Option>
                      <Option value="high_school">高中</Option>
                      <Option value="college">大专</Option>
                      <Option value="bachelor">本科</Option>
                      <Option value="master">硕士</Option>
                      <Option value="phd">博士</Option>
                    </Select>
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="prefer_income" label="收入要求">
                    <Select>
                      <Option value="">不限</Option>
                      <Option value="low">3万以下</Option>
                      <Option value="mid">3-10万</Option>
                      <Option value="high">10-30万</Option>
                      <Option value="luxury">30万以上</Option>
                    </Select>
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="prefer_city" label="城市偏好">
                    <Input />
                  </Form.Item>
                </Col>
              </Row>
              <Form.Item>
                <Button type="primary" htmlType="submit" loading={updateMutation.isPending}>
                  保存
                </Button>
              </Form.Item>
            </Form>
          )}
        </Card>
      </Col>
    </Row>
  )
}
