import React, { useEffect, useState } from 'react';
import { Card, Row, Col, Tag, Descriptions, Button, List, Avatar, Spin, message } from 'antd';
import { useParams, useNavigate } from 'react-router-dom';
import { ArrowLeftOutlined, CalendarOutlined } from '@ant-design/icons';
import { skillApi, postingApi } from '@/api/skill';
import { Skill, SkillPosting } from '@/types';

const SkillDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(true);
  const [skill, setSkill] = useState<Skill | null>(null);
  const [postings, setPostings] = useState<SkillPosting[]>([]);

  useEffect(() => {
    const fetchData = async () => {
      if (!id) return;

      setLoading(true);
      try {
        const [skillData, postingsData] = await Promise.all([
          skillApi.get(id),
          postingApi.list({ skill_id: id, page_size: 10 }),
        ]);

        setSkill(skillData);
        setPostings(postingsData.items || []);
      } catch (error: any) {
        message.error(error.message || '获取技能详情失败');
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [id]);

  const difficultyLabels: Record<string, string> = {
    beginner: '入门',
    intermediate: '进阶',
    advanced: '高级',
    expert: '专家',
  };

  const difficultyColors: Record<string, string> = {
    beginner: 'green',
    intermediate: 'blue',
    advanced: 'orange',
    expert: 'red',
  };

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: 100 }}>
        <Spin size="large" />
      </div>
    );
  }

  if (!skill) {
    return <div>技能不存在</div>;
  }

  return (
    <div>
      <Button
        icon={<ArrowLeftOutlined />}
        onClick={() => navigate(-1)}
        style={{ marginBottom: 16 }}
      >
        返回
      </Button>

      <Row gutter={[24, 24]}>
        <Col span={16}>
          <Card
            cover={
              skill.cover_image ? (
                <img
                  alt={skill.title}
                  src={skill.cover_image}
                  style={{ maxHeight: 300, objectFit: 'cover' }}
                />
              ) : (
                <div
                  style={{
                    height: 200,
                    background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
                  }}
                />
              )
            }
          >
            <h1 style={{ marginBottom: 16 }}>{skill.title}</h1>
            <div style={{ marginBottom: 16 }}>
              <Tag color={difficultyColors[skill.difficulty]}>
                {difficultyLabels[skill.difficulty]}
              </Tag>
              {skill.category && <Tag color="geekblue">{skill.category.name}</Tag>}
              {skill.tags?.map((tag) => (
                <Tag key={tag.id} color="purple">
                  {tag.name}
                </Tag>
              ))}
            </div>
            <Descriptions column={2} size="small">
              <Descriptions.Item label="评分">
                <span style={{ color: '#faad14' }}>★</span> {skill.rating.toFixed(1)}
              </Descriptions.Item>
              <Descriptions.Item label="评价数">{skill.review_count}</Descriptions.Item>
              <Descriptions.Item label="课程数">{skill.posting_count}</Descriptions.Item>
            </Descriptions>

            {skill.description && (
              <div style={{ marginTop: 16 }}>
                <h3>技能描述</h3>
                <p style={{ whiteSpace: 'pre-wrap' }}>{skill.description}</p>
              </div>
            )}

            {skill.prerequisites && (
              <div style={{ marginTop: 16 }}>
                <h3>前置条件</h3>
                <p style={{ whiteSpace: 'pre-wrap' }}>{skill.prerequisites}</p>
              </div>
            )}

            {skill.outcomes && (
              <div style={{ marginTop: 16 }}>
                <h3>学习成果</h3>
                <p style={{ whiteSpace: 'pre-wrap' }}>{skill.outcomes}</p>
              </div>
            )}
          </Card>
        </Col>

        <Col span={8}>
          <Card title="相关课程" style={{ marginBottom: 16 }}>
            {postings.length === 0 ? (
              <div style={{ textAlign: 'center', color: '#999', padding: 20 }}>
                暂无相关课程
              </div>
            ) : (
              <List
                itemLayout="horizontal"
                dataSource={postings}
                renderItem={(item) => (
                  <List.Item
                    onClick={() => navigate(`/postings/${item.id}`)}
                    style={{ cursor: 'pointer' }}
                  >
                    <List.Item.Meta
                      avatar={<Avatar icon={<CalendarOutlined />} />}
                      title={item.title}
                      description={
                        <div>
                          <div>{item.teacher?.nickname}</div>
                          <div style={{ color: '#1890ff', fontWeight: 'bold' }}>
                            ¥{item.price_per_hour}/小时
                          </div>
                        </div>
                      }
                    />
                  </List.Item>
                )}
              />
            )}
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default SkillDetail;
