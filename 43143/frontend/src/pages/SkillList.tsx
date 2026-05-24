import React, { useEffect, useState } from 'react';
import { Row, Col, Card, Tag, Input, Select, Button, Spin, message, Pagination } from 'antd';
import { SearchOutlined, AppstoreOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { skillApi, categoryApi } from '@/api/skill';
import { Skill, SkillCategory, DifficultyLevel } from '@/types';

const { Option } = Select;

const SkillList: React.FC = () => {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(true);
  const [skills, setSkills] = useState<Skill[]>([]);
  const [categories, setCategories] = useState<SkillCategory[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize] = useState(12);
  const [keyword, setKeyword] = useState('');
  const [selectedCategory, setSelectedCategory] = useState<string | undefined>();
  const [selectedDifficulty, setSelectedDifficulty] = useState<DifficultyLevel | undefined>();

  useEffect(() => {
    const fetchCategories = async () => {
      try {
        const data = await categoryApi.list();
        setCategories(data);
      } catch (error) {
        console.error('获取分类失败:', error);
      }
    };

    fetchCategories();
  }, []);

  useEffect(() => {
    const fetchSkills = async () => {
      setLoading(true);
      try {
        const data = await skillApi.list({
          page,
          page_size: pageSize,
          category_id: selectedCategory,
          keyword: keyword || undefined,
        });
        setSkills(data.items || []);
        setTotal(data.total || 0);
      } catch (error: any) {
        message.error(error.message || '获取技能列表失败');
      } finally {
        setLoading(false);
      }
    };

    fetchSkills();
  }, [page, pageSize, selectedCategory, keyword]);

  const difficultyColors: Record<DifficultyLevel, string> = {
    beginner: 'green',
    intermediate: 'blue',
    advanced: 'orange',
    expert: 'red',
  };

  const difficultyLabels: Record<DifficultyLevel, string> = {
    beginner: '入门',
    intermediate: '进阶',
    advanced: '高级',
    expert: '专家',
  };

  if (loading && skills.length === 0) {
    return (
      <div style={{ textAlign: 'center', padding: 100 }}>
        <Spin size="large" />
      </div>
    );
  }

  return (
    <div>
      <Card style={{ marginBottom: 16 }}>
        <Row gutter={16} align="middle">
          <Col flex="auto">
            <Input
              placeholder="搜索技能..."
              prefix={<SearchOutlined />}
              value={keyword}
              onChange={(e) => setKeyword(e.target.value)}
              onPressEnter={() => setPage(1)}
              style={{ maxWidth: 400 }}
            />
          </Col>
          <Col>
            <Select
              placeholder="选择分类"
              allowClear
              style={{ width: 150 }}
              value={selectedCategory}
              onChange={(value) => {
                setSelectedCategory(value);
                setPage(1);
              }}
            >
              {categories.map((cat) => (
                <Option key={cat.id} value={cat.id}>
                  {cat.name}
                </Option>
              ))}
            </Select>
          </Col>
          <Col>
            <Select
              placeholder="难度等级"
              allowClear
              style={{ width: 120 }}
              value={selectedDifficulty}
              onChange={(value) => {
                setSelectedDifficulty(value);
                setPage(1);
              }}
            >
              <Option value="beginner">入门</Option>
              <Option value="intermediate">进阶</Option>
              <Option value="advanced">高级</Option>
              <Option value="expert">专家</Option>
            </Select>
          </Col>
          <Col>
            <Button type="primary" onClick={() => setPage(1)}>
              搜索
            </Button>
          </Col>
        </Row>
      </Card>

      <Row gutter={[16, 16]}>
        {skills.map((skill) => (
          <Col span={6} key={skill.id}>
            <Card
              hoverable
              cover={
                skill.cover_image ? (
                  <img
                    alt={skill.title}
                    src={skill.cover_image}
                    style={{ height: 140, objectFit: 'cover' }}
                  />
                ) : (
                  <div
                    style={{
                      height: 140,
                      background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                      color: '#fff',
                      fontSize: 40,
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
                    <Tag color={difficultyColors[skill.difficulty]}>
                      {difficultyLabels[skill.difficulty]}
                    </Tag>
                    {skill.category && <Tag color="geekblue">{skill.category.name}</Tag>}
                    <div style={{ marginTop: 8, color: '#666' }}>
                      {skill.description?.substring(0, 50)}...
                    </div>
                    <div style={{ marginTop: 8 }}>
                      <span style={{ color: '#faad14' }}>★</span> {skill.rating.toFixed(1)}
                      <span style={{ color: '#999', marginLeft: 8 }}>
                        {skill.review_count}条评价
                      </span>
                    </div>
                  </div>
                }
              />
            </Card>
          </Col>
        ))}
      </Row>

      {total > 0 && (
        <div style={{ textAlign: 'center', marginTop: 24 }}>
          <Pagination
            current={page}
            pageSize={pageSize}
            total={total}
            onChange={(p) => setPage(p)}
            showSizeChanger={false}
          />
        </div>
      )}
    </div>
  );
};

export default SkillList;
