import React, { useState, useEffect } from 'react'
import { Card, Input, Select, Pagination, Spin, Alert, Avatar, Rate, Button, Empty, Tag } from 'antd'
import { SearchOutlined, UserOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { doctorAPI, departmentAPI } from '@/services/api'
import type { Doctor, Department } from '@/types'

const { Search } = Input
const { Option } = Select

const DoctorList: React.FC = () => {
  const navigate = useNavigate()
  const [doctors, setDoctors] = useState<Doctor[]>([])
  const [departments, setDepartments] = useState<Department[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [keyword, setKeyword] = useState('')
  const [departmentId, setDepartmentId] = useState<number | undefined>()
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(8)
  const [total, setTotal] = useState(0)

  useEffect(() => {
    fetchDepartments()
  }, [])

  useEffect(() => {
    fetchDoctors()
  }, [keyword, departmentId, page, pageSize])

  const fetchDepartments = async () => {
    try {
      const response = await departmentAPI.getList({ page: 1, pageSize: 100 })
      setDepartments(response.list)
    } catch (err) {
      console.error('获取科室列表失败:', err)
    }
  }

  const fetchDoctors = async () => {
    setLoading(true)
    setError(null)
    try {
      const params: any = { page, pageSize }
      if (keyword) params.keyword = keyword
      if (departmentId) params.department_id = departmentId
      const response = await doctorAPI.getList(params)
      setDoctors(response.list)
      setTotal(response.total)
    } catch (err: any) {
      setError(err.message || '获取医生列表失败')
    } finally {
      setLoading(false)
    }
  }

  const handleSearch = (value: string) => {
    setKeyword(value)
    setPage(1)
  }

  const handleDepartmentChange = (value: number | undefined) => {
    setDepartmentId(value)
    setPage(1)
  }

  const handleCardClick = (doctorId: number) => {
    navigate(`/doctors/${doctorId}`)
  }

  return (
    <div className="p-6 max-w-7xl mx-auto">
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-gray-800 mb-2">医生列表</h1>
        <p className="text-gray-500">选择您需要的医生进行预约挂号</p>
      </div>

      <div className="mb-6 flex flex-wrap gap-4 items-center">
        <Search
          placeholder="搜索医生姓名..."
          allowClear
          enterButton={<SearchOutlined />}
          size="large"
          onSearch={handleSearch}
          style={{ width: 300 }}
        />
        <Select
          placeholder="选择科室"
          allowClear
          size="large"
          style={{ width: 200 }}
          onChange={handleDepartmentChange}
          value={departmentId}
        >
          {departments.map((dept) => (
            <Option key={dept.id} value={dept.id}>
              {dept.name}
            </Option>
          ))}
        </Select>
      </div>

      {error && (
        <Alert
          message="错误"
          description={error}
          type="error"
          showIcon
          className="mb-4"
          closable
          onClose={() => setError(null)}
        />
      )}

      <Spin spinning={loading} tip="加载中...">
        {doctors.length > 0 ? (
          <>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
              {doctors.map((doctor) => (
                <Card
                  key={doctor.id}
                  hoverable
                  className="cursor-pointer transition-all duration-300 hover:shadow-lg"
                  onClick={() => handleCardClick(doctor.id)}
                >
                  <div className="flex items-start gap-4">
                    <Avatar
                      size={64}
                      src={doctor.user?.avatar_url}
                      icon={!doctor.user?.avatar_url && <UserOutlined />}
                      className="flex-shrink-0"
                    />
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2 mb-1">
                        <span className="font-semibold text-lg truncate">
                          {doctor.user?.full_name}
                        </span>
                        <Tag color="blue" className="flex-shrink-0">
                          {doctor.title}
                        </Tag>
                      </div>
                      <div className="text-sm text-gray-500 mb-2">
                        {doctor.department?.name}
                      </div>
                      <div className="flex items-center gap-1 mb-2">
                        <Rate
                          disabled
                          value={doctor.average_rating}
                          allowHalf
                          className="text-xs"
                        />
                        <span className="text-sm text-gray-500">
                          {doctor.average_rating.toFixed(1)}
                        </span>
                        <span className="text-xs text-gray-400 ml-1">
                          ({doctor.review_count}条评价)
                        </span>
                      </div>
                      <p className="text-sm text-gray-600 line-clamp-2 mb-3">
                        {doctor.introduction || '暂无简介'}
                      </p>
                      <div className="flex items-center justify-between">
                        <span className="text-orange-500 font-medium">
                          ¥{doctor.registration_fee}
                        </span>
                        <Button type="primary" size="small">
                          查看详情
                        </Button>
                      </div>
                    </div>
                  </div>
                </Card>
              ))}
            </div>

            <div className="mt-8 flex justify-center">
              <Pagination
                current={page}
                pageSize={pageSize}
                total={total}
                showSizeChanger
                showQuickJumper
                showTotal={(total) => `共 ${total} 位医生`}
                onChange={(page, pageSize) => {
                  setPage(page)
                  setPageSize(pageSize)
                }}
              />
            </div>
          </>
        ) : (
          !loading && (
            <Empty
              description="暂无符合条件的医生"
              image={Empty.PRESENTED_IMAGE_SIMPLE}
              className="py-16"
            />
          )
        )}
      </Spin>
    </div>
  )
}

export default DoctorList
