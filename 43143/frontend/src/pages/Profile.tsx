import React, { useEffect, useState } from 'react';
import { Card, Row, Col, Avatar, Descriptions, Tag, Button, List, Rate, Spin, message, Tabs } from 'antd';
import { UserOutlined, EditOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { useSelector } from 'react-redux';
import { userApi } from '@/api/user';
import { reviewApi } from '@/api/booking';
import { User, Review } from '@/types';
import { RootState } from '@/store';

const Profile: React.FC = () => {
  const navigate = useNavigate();
  const currentUser = useSelector((state: RootState) => state.auth.user);
  const [loading, setLoading] = useState(true);
  const [profile, setProfile] = useState<User | null>(null);
  const [reviews, setReviews] = useState<Review[]>([]);

  useEffect(() => {
    const fetchData = async () => {
      if (!currentUser?.id) return;

      setLoading(true);
      try {
        const [userData, reviewsData] = await Promise.all([
          userApi.getProfile(),
          reviewApi.getByUser(currentUser.id, { page_size: 10 }),
        ]);

        setProfile(userData);
        setReviews(reviewsData.items || []);
      } catch (error: any) {
        message.error(error.message || '获取用户信息失败');
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [currentUser?.id]);

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: 100 }}>
        <Spin size="large" />
      </div>
    );
  }

  if (!profile) {
    return <div>用户信息加载失败</div>;
  }

  const tabItems = [
    {
      key: 'info',
      label: '基本信息',
      children: (
        <Card>
          <Descriptions column={2} bordered>
            <Descriptions.Item label="昵称">{profile.nickname}</Descriptions.Item>
            <Descriptions.Item label="邮箱">
              {profile.email || '未设置'}
              {profile.email_verified && <Tag color="green" style={{ marginLeft: 8 }}>已验证</Tag>}
            </Descriptions.Item>
            <Descriptions.Item label="手机号">
              {profile.phone || '未设置'}
              {profile.phone_verified && <Tag color="green" style={{ marginLeft: 8 }}>已验证</Tag>}
            </Descriptions.Item>
            <Descriptions.Item label="性别">{profile.gender || '未设置'}</Descriptions.Item>
            <Descriptions.Item label="生日">{profile.birthday ? new Date(profile.birthday).toLocaleDateString() : '未设置'}</Descriptions.Item>
            <Descriptions.Item label="地区">{profile.location || '未设置'}</Descriptions.Item>
            <Descriptions.Item label="角色" span={2}>
              <Tag color="blue">{profile.role === 'learner' ? '学员' : profile.role === 'teacher' ? '老师' : '学员/老师'}</Tag>
            </Descriptions.Item>
            <Descriptions.Item label="个人简介" span={2}>{profile.bio || '暂无简介'}</Descriptions.Item>
          </Descriptions>
        </Card>
      ),
    },
    {
      key: 'skills',
      label: '技能标签',
      children: (
        <Card title="我的技能">
          <div style={{ marginBottom: 16 }}>
            <h4>我擅长的技能</h4>
            <div>
              {profile.teach_tags?.length ? profile.teach_tags.map((tag) => (
                <Tag key={tag.id} color="blue" style={{ marginBottom: 8 }}>
                  {tag.name}
                </Tag>
              )) : <span style={{ color: '#999' }}>暂无</span>}
            </div>
          </div>
          <div>
            <h4>我想学的技能</h4>
            <div>
              {profile.learn_tags?.length ? profile.learn_tags.map((tag) => (
                <Tag key={tag.id} color="green" style={{ marginBottom: 8 }}>
                  {tag.name}
                </Tag>
              )) : <span style={{ color: '#999' }}>暂无</span>}
            </div>
          </div>
        </Card>
      ),
    },
    {
      key: 'reviews',
      label: '收到的评价',
      children: (
        <Card title="评价列表">
          {reviews.length === 0 ? (
            <div style={{ textAlign: 'center', padding: 40, color: '#999' }}>
              暂无评价
            </div>
          ) : (
            <List
              itemLayout="horizontal"
              dataSource={reviews}
              renderItem={(item) => (
                <List.Item>
                  <List.Item.Meta
                    avatar={<Avatar src={item.reviewer?.avatar} icon={<UserOutlined />} />}
                    title={
                      <div>
                        <span>{item.reviewer?.nickname}</span>
                        <Rate disabled defaultValue={item.rating} style={{ marginLeft: 16, fontSize: 14 }} />
                      </div>
                    }
                    description={
                      <div>
                        <div style={{ color: '#999', fontSize: 12, marginBottom: 4 }}>
                          {new Date(item.created_at).toLocaleString()}
                        </div>
                        <div>{item.content}</div>
                      </div>
                    }
                  />
                </List.Item>
              )}
            />
          )}
        </Card>
      ),
    },
  ];

  return (
    <div>
      <Row gutter={24}>
        <Col span={8}>
          <Card style={{ textAlign: 'center' }}>
            <Avatar
              size={100}
              src={profile.avatar}
              icon={<UserOutlined />}
              style={{ marginBottom: 16 }}
            />
            <h2 style={{ marginBottom: 8 }}>{profile.nickname}</h2>
            <div style={{ marginBottom: 16 }}>
              <Rate disabled defaultValue={profile.rating} />
              <span style={{ marginLeft: 8 }}>({profile.review_count}条评价)</span>
            </div>
            <div style={{ marginBottom: 16 }}>
              <Tag color="blue">{profile.role === 'learner' ? '学员' : profile.role === 'teacher' ? '老师' : '学员/老师'}</Tag>
            </div>
            <Row gutter={16} style={{ marginBottom: 16 }}>
              <Col span={8}>
                <div style={{ fontSize: 24, fontWeight: 'bold' }}>{profile.teaching_hours?.toFixed(1) || 0}</div>
                <div style={{ color: '#999', fontSize: 12 }}>授课时长</div>
              </Col>
              <Col span={8}>
                <div style={{ fontSize: 24, fontWeight: 'bold' }}>{profile.student_count || 0}</div>
                <div style={{ color: '#999', fontSize: 12 }}>学员数量</div>
              </Col>
              <Col span={8}>
                <div style={{ fontSize: 24, fontWeight: 'bold' }}>¥{profile.balance?.toFixed(2) || '0.00'}</div>
                <div style={{ color: '#999', fontSize: 12 }}>账户余额</div>
              </Col>
            </Row>
            <Button
              type="primary"
              icon={<EditOutlined />}
              block
              onClick={() => navigate('/profile/edit')}
            >
              编辑资料
            </Button>
          </Card>
        </Col>
        <Col span={16}>
          <Tabs defaultActiveKey="info" items={tabItems} />
        </Col>
      </Row>
    </div>
  );
};

export default Profile;
