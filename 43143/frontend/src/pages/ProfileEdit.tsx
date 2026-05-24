import React, { useEffect, useState } from 'react';
import { Card, Form, Input, Select, Button, DatePicker, message, Spin, Row, Col } from 'antd';
import { UserOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { useSelector, useDispatch } from 'react-redux';
import dayjs from 'dayjs';
import { userApi } from '@/api/user';
import { tagApi } from '@/api/skill';
import { setUser } from '@/store';
import { SkillTag } from '@/types';
import { RootState } from '@/store';

const { Option } = Select;
const { TextArea } = Input;

const ProfileEdit: React.FC = () => {
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const currentUser = useSelector((state: RootState) => state.auth.user);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [allTags, setAllTags] = useState<SkillTag[]>([]);
  const [form] = Form.useForm();

  useEffect(() => {
    const fetchData = async () => {
      if (!currentUser?.id) return;

      setLoading(true);
      try {
        const [userData, tagsData] = await Promise.all([
          userApi.getProfile(),
          tagApi.list(),
        ]);

        setAllTags(tagsData);

        form.setFieldsValue({
          nickname: userData.nickname,
          email: userData.email,
          phone: userData.phone,
          gender: userData.gender,
          birthday: userData.birthday ? dayjs(userData.birthday) : null,
          location: userData.location,
          bio: userData.bio,
          role: userData.role,
          teach_tags: userData.teach_tags?.map((t) => t.id) || [],
          learn_tags: userData.learn_tags?.map((t) => t.id) || [],
        });
      } catch (error: any) {
        message.error(error.message || '获取用户信息失败');
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [currentUser?.id]);

  const handleSubmit = async (values: any) => {
    setSaving(true);
    try {
      const updates: any = {
        nickname: values.nickname,
        gender: values.gender,
        location: values.location,
        bio: values.bio,
        role: values.role,
      };

      if (values.birthday) {
        updates.birthday = values.birthday.format('YYYY-MM-DD');
      }

      const user = await userApi.updateProfile(updates);

      if (values.teach_tags?.length) {
        await userApi.addSkillTags(values.teach_tags);
      }

      dispatch(setUser(user));
      message.success('保存成功');
      navigate('/profile');
    } catch (error: any) {
      message.error(error.message || '保存失败');
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: 100 }}>
        <Spin size="large" />
      </div>
    );
  }

  return (
    <div>
      <Card title="编辑个人资料">
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
          style={{ maxWidth: 600 }}
        >
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="nickname"
                label="昵称"
                rules={[{ required: true, message: '请输入昵称' }]}
              >
                <Input prefix={<UserOutlined />} placeholder="昵称" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="email" label="邮箱">
                <Input placeholder="邮箱" disabled />
              </Form.Item>
            </Col>
          </Row>

          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="phone" label="手机号">
                <Input placeholder="手机号" disabled />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="gender" label="性别">
                <Select placeholder="选择性别">
                  <Option value="male">男</Option>
                  <Option value="female">女</Option>
                  <Option value="other">其他</Option>
                </Select>
              </Form.Item>
            </Col>
          </Row>

          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="birthday" label="生日">
                <DatePicker style={{ width: '100%' }} />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="location" label="地区">
                <Input placeholder="请输入地区" />
              </Form.Item>
            </Col>
          </Row>

          <Form.Item name="role" label="身份角色">
            <Select>
              <Option value="learner">学员</Option>
              <Option value="teacher">老师</Option>
              <Option value="both">学员/老师</Option>
            </Select>
          </Form.Item>

          <Form.Item name="teach_tags" label="我擅长的技能">
            <Select mode="multiple" placeholder="选择擅长的技能">
              {allTags.map((tag) => (
                <Option key={tag.id} value={tag.id}>
                  {tag.name}
                </Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item name="learn_tags" label="我想学的技能">
            <Select mode="multiple" placeholder="选择想学的技能">
              {allTags.map((tag) => (
                <Option key={tag.id} value={tag.id}>
                  {tag.name}
                </Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item name="bio" label="个人简介">
            <TextArea rows={4} placeholder="介绍一下自己" maxLength={500} />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" loading={saving}>
              保存
            </Button>
            <Button style={{ marginLeft: 8 }} onClick={() => navigate('/profile')}>
              取消
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  );
};

export default ProfileEdit;
