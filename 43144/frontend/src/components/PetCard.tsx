import React from 'react'
import { Card, Tag, Button } from 'antd'
import { EyeOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { Pet } from '../types'

interface PetCardProps {
  pet: Pet
  showActions?: boolean
}

const statusColorMap: Record<string, string> = {
  adoptable: 'green',
  adopted: 'blue',
  treatment: 'orange',
  deceased: 'default',
}

const statusTextMap: Record<string, string> = {
  adoptable: '待领养',
  adopted: '已领养',
  treatment: '治疗中',
  deceased: '已去世',
}

const genderMap: Record<string, string> = {
  male: '公',
  female: '母',
  unknown: '未知',
}

const PetCard: React.FC<PetCardProps> = ({ pet, showActions = true }) => {
  const navigate = useNavigate()

  const photos = pet.photos ? pet.photos.split(',').filter(Boolean) : []
  const firstPhoto = photos.length > 0 ? photos[0] : ''

  return (
    <Card
      hoverable
      cover={
        firstPhoto ? (
          <img
            alt={pet.name}
            src={firstPhoto}
            className="pet-photo"
            onError={(e) => {
              (e.target as HTMLImageElement).style.display = 'none'
            }}
          />
        ) : (
          <div
            style={{
              height: 200,
              background: '#f0f0f0',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              fontSize: 48,
            }}
          >
            🐾
          </div>
        )
      }
      actions={
        showActions
          ? [
              <Button
                type="link"
                key="view"
                icon={<EyeOutlined />}
                onClick={() => navigate(`/pets/${pet.id}`)}
              >
                查看详情
              </Button>,
            ]
          : undefined
      }
    >
      <Card.Meta
        title={
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <span>{pet.name}</span>
            <Tag color={statusColorMap[pet.status]}>{statusTextMap[pet.status]}</Tag>
          </div>
        }
        description={
          <div>
            <div style={{ marginBottom: 8 }}>
              <span>编号: {pet.archive_number}</span>
            </div>
            <div style={{ marginBottom: 8 }}>
              {pet.breed && <span>品种: {pet.breed} </span>}
              {pet.age && <span>| 年龄: {pet.age} </span>}
              <span>| 性别: {genderMap[pet.gender] || pet.gender}</span>
            </div>
            {pet.description && (
              <div style={{ color: '#666', fontSize: 12 }}>
                {pet.description.length > 50 ? pet.description.slice(0, 50) + '...' : pet.description}
              </div>
            )}
          </div>
        }
      />
    </Card>
  )
}

export default PetCard
