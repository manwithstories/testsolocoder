import React, { useState, useEffect } from 'react';
import { Row, Col, Card, Typography, Button, Tag, Input, Rate } from 'antd';
import { SearchOutlined, BookOutlined, FileTextOutlined, UserOutlined, ArrowRightOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { textbookApi, noteApi, userApi } from '../services/api';
import { Textbook } from '../types';
import { TextbookCard } from '../components/TextbookCard';
import { NoteCard } from '../components/NoteCard';
import { Loading } from '../components/Loading';

const { Title, Text } = Typography;

export const HomePage: React.FC = () => {
  const [loading, setLoading] = useState(true);
  const [popularTextbooks, setPopularTextbooks] = useState<Textbook[]>([]);
  const [featuredNotes, setFeaturedNotes] = useState<any[]>([]);
  const [topUsers, setTopUsers] = useState<any[]>([]);
  const navigate = useNavigate();

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      const [textbooksRes, notesRes, usersRes]: any = await Promise.all([
        textbookApi.getPopular(8),
        noteApi.getFeatured(4),
        userApi.getTopRated(5),
      ]);
      setPopularTextbooks(textbooksRes.data || []);
      setFeaturedNotes(notesRes.data || []);
      setTopUsers(usersRes.data || []);
    } catch (error) {
      console.error('Failed to load data:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = (value: string) => {
    if (value) {
      navigate(`/textbooks?keyword=${encodeURIComponent(value)}`);
    }
  };

  if (loading) return <Loading />;

  return (
    <div className="bg-gray-50">
      <div className="gradient-bg py-16 px-4">
        <div className="max-w-4xl mx-auto text-center text-white">
          <Title level={1} className="!text-white !mb-4">
            校园二手教材交易平台
          </Title>
          <Text className="text-white text-lg block mb-8">
            循环利用教材，共享学习笔记，共建绿色校园
          </Text>
          <Input.Search
            size="large"
            placeholder="搜索教材、课程、笔记..."
            enterButton={<SearchOutlined />}
            onSearch={handleSearch}
            className="max-w-xl"
          />
          <div className="mt-6 flex justify-center gap-4">
            <Button size="large" onClick={() => navigate('/textbooks')}>
              浏览教材
            </Button>
            <Button size="large" type="primary" onClick={() => navigate('/notes')}>
              浏览笔记
            </Button>
          </div>
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-4 py-12">
        <Row gutter={[24, 24]} className="mb-12">
          <Col xs={24} md={8}>
            <Card className="text-center h-full">
              <BookOutlined className="text-5xl text-blue-500 mb-4" />
              <Title level={4}>教材交易</Title>
              <Text type="secondary">买卖二手教材，循环利用资源</Text>
            </Card>
          </Col>
          <Col xs={24} md={8}>
            <Card className="text-center h-full">
              <FileTextOutlined className="text-5xl text-green-500 mb-4" />
              <Title level={4}>笔记共享</Title>
              <Text type="secondary">分享学习笔记，共同进步</Text>
            </Card>
          </Col>
          <Col xs={24} md={8}>
            <Card className="text-center h-full">
              <UserOutlined className="text-5xl text-purple-500 mb-4" />
              <Title level={4}>安全交易</Title>
              <Text type="secondary">身份认证，交易有保障</Text>
            </Card>
          </Col>
        </Row>

        <div className="mb-12">
          <div className="flex justify-between items-center mb-6">
            <Title level={3} className="!mb-0">热门教材</Title>
            <Button type="link" onClick={() => navigate('/textbooks')}>
              查看更多 <ArrowRightOutlined />
            </Button>
          </div>
          <Row gutter={[16, 16]}>
            {popularTextbooks.slice(0, 8).map((textbook) => (
              <Col xs={12} sm={12} md={8} lg={6} key={textbook.id}>
                <TextbookCard textbook={textbook} />
              </Col>
            ))}
          </Row>
        </div>

        <div className="mb-12">
          <div className="flex justify-between items-center mb-6">
            <Title level={3} className="!mb-0">精选笔记</Title>
            <Button type="link" onClick={() => navigate('/notes')}>
              查看更多 <ArrowRightOutlined />
            </Button>
          </div>
          <Row gutter={[16, 16]}>
            {featuredNotes.slice(0, 4).map((note) => (
              <Col xs={24} sm={12} md={12} lg={6} key={note.id}>
                <NoteCard note={note} />
              </Col>
            ))}
          </Row>
        </div>

        <div>
          <Title level={3} className="!mb-6">优质用户</Title>
          <Row gutter={[16, 16]}>
            {topUsers.map((user) => (
              <Col xs={12} sm={8} md={6} lg={4} key={user.id}>
                <Card className="text-center">
                  <img
                    src={user.avatar || 'https://api.dicebear.com/7.x/identicon/svg?seed=' + user.username}
                    alt={user.username}
                    className="w-16 h-16 rounded-full mx-auto mb-3"
                  />
                  <Text strong className="block">{user.username}</Text>
                  <Rate disabled allowHalf value={user.rating} className="text-xs" />
                  <div className="mt-2">
                    <Tag color={user.role === 'merchant' ? 'purple' : 'blue'}>
                      {user.role === 'merchant' ? '书商' : '学生'}
                    </Tag>
                  </div>
                </Card>
              </Col>
            ))}
          </Row>
        </div>
      </div>
    </div>
  );
};
