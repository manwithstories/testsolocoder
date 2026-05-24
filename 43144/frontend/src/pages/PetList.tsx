import React, { useEffect, useState } from 'react'
import { Row, Col, Input, Select, Pagination, Card, Empty, Spin } from 'antd'
import PetCard from '../components/PetCard'
import { listPets } from '../api/pet'
import { Pet, PetListQuery } from '../types'

const { Option } = Select

const PetList: React.FC = () => {
  const [pets, setPets] = useState<Pet[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize] = useState(12)
  const [loading, setLoading] = useState(false)
  const [filters, setFilters] = useState<PetListQuery>({})

  useEffect(() => {
    loadPets()
  }, [page, filters])

  const loadPets = async () => {
    setLoading(true)
    try {
      const response = await listPets({
        page,
        page_size: pageSize,
        ...filters,
      })
      if (response.code === 0 && response.data) {
        const data = response.data as any
        setPets(data.items || [])
        setTotal(data.total || 0)
      }
    } catch (error) {
      console.error('Failed to load pets:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleSearch = (value: string) => {
    setPage(1)
    setFilters({ ...filters, search: value || undefined })
  }

  const handleStatusChange = (value: string) => {
    setPage(1)
    setFilters({ ...filters, status: value || undefined })
  }

  const handleSpeciesChange = (value: string) => {
    setPage(1)
    setFilters({ ...filters, species: value || undefined })
  }

  const handleGenderChange = (value: string) => {
    setPage(1)
    setFilters({ ...filters, gender: value || undefined })
  }

  return (
    <div>
      <Card style={{ marginBottom: 16 }}>
        <Row gutter={16}>
          <Col xs={24} sm={12} md={6}>
            <Input.Search
              placeholder="搜索宠物名称、品种..."
              allowClear
              onSearch={handleSearch}
            />
          </Col>
          <Col xs={24} sm={12} md={6}>
            <Select
              placeholder="选择状态"
              allowClear
              style={{ width: '100%' }}
              onChange={handleStatusChange}
            >
              <Option value="adoptable">待领养</Option>
              <Option value="adopted">已领养</Option>
              <Option value="treatment">治疗中</Option>
            </Select>
          </Col>
          <Col xs={24} sm={12} md={6}>
            <Select
              placeholder="选择物种"
              allowClear
              style={{ width: '100%' }}
              onChange={handleSpeciesChange}
            >
              <Option value="dog">狗</Option>
              <Option value="cat">猫</Option>
              <Option value="rabbit">兔子</Option>
              <Option value="bird">鸟</Option>
              <Option value="other">其他</Option>
            </Select>
          </Col>
          <Col xs={24} sm={12} md={6}>
            <Select
              placeholder="选择性别"
              allowClear
              style={{ width: '100%' }}
              onChange={handleGenderChange}
            >
              <Option value="male">公</Option>
              <Option value="female">母</Option>
            </Select>
          </Col>
        </Row>
      </Card>

      <Spin spinning={loading}>
        {pets.length > 0 ? (
          <>
            <Row gutter={[16, 16]}>
              {pets.map((pet) => (
                <Col xs={24} sm={12} md={8} lg={6} key={pet.id}>
                  <PetCard pet={pet} />
                </Col>
              ))}
            </Row>
            <div style={{ textAlign: 'center', marginTop: 24 }}>
              <Pagination
                current={page}
                pageSize={pageSize}
                total={total}
                onChange={setPage}
                showSizeChanger={false}
              />
            </div>
          </>
        ) : (
          !loading && (
            <Empty description="暂无宠物数据" style={{ marginTop: 60 }} />
          )
        )}
      </Spin>
    </div>
  )
}

export default PetList
