import React, { useEffect, useState } from 'react';
import { Row, Col, Card, Statistic, Tag, Button, Spin, message } from 'antd';
import {
  AppstoreOutlined,
  CalendarOutlined,
  StarOutlined,
  UserOutlined,
  ArrowRightOutlined,
} from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { skillApi, postingApi } from '@/api/skill';
import { Skill, SkillPosting } from '@/types';

const Home: React.FC = () => {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(true);
  const [popularSkills, setPopularSkills] = useState<Skill[]>([]);
  const [recentPostings, setRecentPostings] = useState<SkillPosting[]>([]);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [skillsData, postingsData] = await Promise.all([
          skillApi.getPopular({ limit: 6 }),
          postingApi.list({ page_size: 8 }),
        ]);

        setPopularSkills(skillsData);
        setRecentPostings(postingsData.items || []);
      } catch (error: any) {
        message.error(error.message || '加载数据失败');
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: 100 }}>
        <Spin size="large" />
      </div>
    );
  }

  return (
    <div>
      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        <Col span={6}>
          <Card>
            <Statistic
              title="热门技能"
              value={popularSkills.length}
              prefix={<AppstoreOutlined />}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="今日课程"
              value={recentPostings.length}
              prefix={<CalendarOutlined />}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="平均评分"
              value={4.8}
              precision={1}
              prefix={<StarOutlined />}
              suffix="/ 5"
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="活跃用户"
              value={1258}
              prefix={<UserOutlined />}
            />
          </Card>
        </Col>
      </Row>

      <Card
        title="热门技能"
        extra={
          <Button type="link" onClick={() => navigate('/skills')}>
            查看更多 <ArrowRightOutlined />
          </Button>
        }
        style={{ marginBottom: 24 }}
      >
        <Row gutter={[16, 16]}>
          {popularSkills.map((skill) => (
            <Col span={8} key={skill.id}>
              <Card
                hoverable
                cover={
                  skill.cover_image ? (
                    <img
                      alt={skill.title}
                      src={skill.cover_image}
                      style={{ height: 160, objectFit: 'cover' }}
                    />
                  ) : (
                    <div
                      style={{
                        height: 160,
                        background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        color: '#fff',
                        fontSize: 48,
                      }}
                    >
                      <AppstoreOutlined />
                    </div>
                  )
                }
                onClick={() => navigate(`/skills/${skill.id}`)}
              >
                <Card.Meta
                  title={skill.title}
                  description={
                    <div>
                      <Tag color="blue">{skill.difficulty}</Tag>
                      <Tag color="green">{skill.category?.name}</Tag>
                      <div style={{ marginTop: 8 }}>
                        <StarOutlined style={{ color: '#faad14' }} /> {skill.rating.toFixed(1)}
                        <span style={{ marginLeft: 8, color: '#999' }}>
                          ({skill.review_count}条评价)
                        </span>
                      </div>
                    </div>
                  }
                />
              </Card>
            </Col>
          ))}
        </Row>
      </Card>

      <Card
        title="最新课程"
        extra={
          <Button type="link" onClick={() => navigate('/skills')}>
            查看更多 <ArrowRightOutlined />
          </Button>
        }
      >
        <Row gutter={[16, 16]}>
          {recentPostings.map((posting) => (
            <Col span={12} key={posting.id}>
              <Card
                hoverable
                onClick={() => navigate(`/postings/${posting.id}`)}
              >
                <Row align="middle" justify="space-between">
                  <Col>
                    <div style={{ fontWeight: 500, fontSize: 16 }}>{posting.title}</div>
                    <div style={{ color: '#666', marginTop: 4 }}>
                      {posting.teacher?.nickname}
                    </div>
                    <div style={{ marginTop: 8 }}>
                      <Tag color="blue">{posting.teaching_method}</Tag>
                      <Tag color="purple">{posting.teaching_mode}</Tag>
                    </div>
                  </Col>
                  <Col>
                    <div style={{ fontSize: 24, fontWeight: 'bold', color: '#1890ff' }}>
                      ¥{posting.price_per_hour}
                    </div>
                    <div style={{ color: '#999', textAlign: 'right' }}>/小时</div>
                  </Col>
                </Row>
              </Card>
            </Col>
          ))}
        </Row>
      </Card>
    </div>
  );
};

export default Home;
