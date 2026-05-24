import React, { useEffect, useState } from 'react'
import { Row, Col, Card, Statistic, Button } from 'antd'
import {
  SmileOutlined,
  HomeOutlined,
  CalendarOutlined,
  FileTextOutlined,
} from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'
import { listPets } from '../api/pet'

const Home: React.FC = () => {
  const { user } = useAuth()
  const navigate = useNavigate()
  const [stats, setStats] = useState({ totalPets: 0, adoptablePets: 0 })

  useEffect(() => {
    loadStats()
  }, [])

  const loadStats = async () => {
    try {
      const [allRes, adoptableRes] = await Promise.all([
        listPets({ page: 1, page_size: 1 }),
        listPets({ page: 1, page_size: 1, status: 'adoptable' }),
      ])
      setStats({
        totalPets: (allRes.data as any)?.total || 0,
        adoptablePets: (adoptableRes.data as any)?.total || 0,
      })
    } catch (error) {
      console.error('Failed to load stats:', error)
    }
  }

  const actionCards = [
    {
      title: '浏览宠物',
      description: '查看所有待领养的宠物',
      icon: <SmileOutlined style={{ fontSize: 48, color: '#1890ff' }} />,
      action: () => navigate('/pets'),
    },
    {
      title: '我的预约',
      description: '管理您的预约记录',
      icon: <CalendarOutlined style={{ fontSize: 48, color: '#52c41a' }} />,
      action: () => navigate('/appointments'),
    },
    {
      title: '健康档案',
      description: '查看宠物健康记录',
      icon: <FileTextOutlined style={{ fontSize: 48, color: '#faad14' }} />,
      action: () => navigate('/health-records'),
    },
    {
      title: '个人中心',
      description: '管理您的账号信息',
      icon: <HomeOutlined style={{ fontSize: 48, color: '#722ed1' }} />,
      action: () => navigate('/profile'),
    },
  ]

  return (
    <div>
      <Card style={{ marginBottom: 24 }}>
        <h2 style={{ margin: 0 }}>欢迎回来，{user?.name}！</h2>
        <p style={{ color: '#666', marginTop: 8 }}>
          {user?.role === 'adopter' && '在这里您可以浏览待领养宠物、预约看望、管理领养记录'}
          {user?.role === 'rescue' && '在这里您可以管理宠物信息、审核领养申请、查看统计数据'}
          {user?.role === 'admin' && '在这里您可以管理用户、审核救助站、查看平台数据'}
        </p>
      </Card>

      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        <Col xs={24} sm={12} md={6}>
          <Card className="stat-card">
            <Statistic title="平台宠物总数" value={stats.totalPets} prefix={<SmileOutlined />} />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <Card className="stat-card">
            <Statistic title="待领养宠物" value={stats.adoptablePets} prefix={<SmileOutlined />} valueStyle={{ color: '#52c41a' }} />
          </Card>
        </Col>
      </Row>

      <Row gutter={[16, 16]}>
        {actionCards.map((card, index) => (
          <Col xs={24} sm={12} md={6} key={index}>
            <Card
              hoverable
              onClick={card.action}
              style={{ textAlign: 'center', cursor: 'pointer' }}
            >
              {card.icon}
              <h3 style={{ marginTop: 16 }}>{card.title}</h3>
              <p style={{ color: '#666' }}>{card.description}</p>
            </Card>
          </Col>
        ))}
      </Row>
    </div>
  )
}

export default Home
