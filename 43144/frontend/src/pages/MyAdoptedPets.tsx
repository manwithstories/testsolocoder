import React, { useEffect, useState } from 'react'
import { Card, Row, Col, Table, Tag, Empty, Spin, Button } from 'antd'
import { EyeOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { getMyAdoptedPets } from '../api/pet'
import { Pet } from '../types'

const MyAdoptedPets: React.FC = () => {
  const navigate = useNavigate()
  const [pets, setPets] = useState<Pet[]>([])
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    loadPets()
  }, [])

  const loadPets = async () => {
    setLoading(true)
    try {
      const response = await getMyAdoptedPets()
      if (response.code === 0) {
        setPets(response.data || [])
      }
    } catch (error) {
      console.error('Failed to load adopted pets:', error)
    } finally {
      setLoading(false)
    }
  }

  const columns = [
    { title: '编号', dataIndex: 'archive_number', key: 'archive_number' },
    { title: '名称', dataIndex: 'name', key: 'name' },
    { title: '物种', dataIndex: 'species', key: 'species' },
    { title: '品种', dataIndex: 'breed', key: 'breed' },
    {
      title: '领养日期',
      dataIndex: 'adopted_date',
      key: 'adopted_date',
      render: (date: string) => date || '-',
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => <Tag color="blue">{status}</Tag>,
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Pet) => (
        <Button
          type="link"
          icon={<EyeOutlined />}
          onClick={() => navigate(`/pets/${record.id}`)}
        >
          查看详情
        </Button>
      ),
    },
  ]

  return (
    <div>
      <h2 style={{ marginBottom: 16 }}>我的领养</h2>
      <Spin spinning={loading}>
        {pets.length > 0 ? (
          <Card>
            <Table
              dataSource={pets}
              columns={columns}
              rowKey="id"
              pagination={{ pageSize: 10 }}
            />
          </Card>
        ) : (
          !loading && (
            <Card>
              <Empty description="您还没有领养任何宠物" style={{ marginTop: 40 }} />
            </Card>
          )
        )}
      </Spin>
    </div>
  )
}

export default MyAdoptedPets
