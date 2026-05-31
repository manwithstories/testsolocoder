import { useEffect, useState } from 'react'
import { Card, Descriptions, Tag, Spin, message, Button, Space } from 'antd'
import { useParams, useNavigate } from 'react-router-dom'
import { ArrowLeftOutlined } from '@ant-design/icons'
import api from '../../api'
import type { Beehive } from '../../types'

function BeehiveDetailPage() {
  const { id } = useParams()
  const navigate = useNavigate()
  const [loading, setLoading] = useState(false)
  const [data, setData] = useState<Beehive | null>(null)

  useEffect(() => {
    if (id) {
      fetchData()
    }
  }, [id])

  const fetchData = async () => {
    setLoading(true)
    try {
      const response = await api.get(`/beehives/${id}`)
      setData(response.data as Beehive)
    } catch (error: any) {
      message.error(error.message || '获取数据失败')
    } finally {
      setLoading(false)
    }
  }

  const statusColorMap: Record<string, string> = {
    active: 'green',
    inactive: 'default',
    harvesting: 'orange',
  }

  const healthColorMap: Record<string, string> = {
    healthy: 'green',
    warning: 'orange',
    critical: 'red',
  }

  const queenColorMap: Record<string, string> = {
    normal: 'green',
    old: 'orange',
    missing: 'red',
    new: 'blue',
  }

  if (loading) {
    return <Spin />
  }

  if (!data) {
    return <Card>数据不存在</Card>
  }

  return (
    <div>
      <Space style={{ marginBottom: 16 }}>
        <Button icon={<ArrowLeftOutlined />} onClick={() => navigate('/beehives')}>
          返回
        </Button>
      </Space>
      <Card title={`蜂箱详情 - ${data.name}`}>
        <Descriptions bordered column={2}>
          <Descriptions.Item label="编号">{data.code}</Descriptions.Item>
          <Descriptions.Item label="名称">{data.name}</Descriptions.Item>
          <Descriptions.Item label="纬度">{data.latitude}</Descriptions.Item>
          <Descriptions.Item label="经度">{data.longitude}</Descriptions.Item>
          <Descriptions.Item label="区域">{data.region || '-'}</Descriptions.Item>
          <Descriptions.Item label="蜜蜂品种">{data.bee_species || '-'}</Descriptions.Item>
          <Descriptions.Item label="分组">{data.group_name || '-'}</Descriptions.Item>
          <Descriptions.Item label="工蜂数量">{data.worker_count}</Descriptions.Item>
          <Descriptions.Item label="状态">
            <Tag color={statusColorMap[data.status]}>{data.status}</Tag>
          </Descriptions.Item>
          <Descriptions.Item label="健康状态">
            <Tag color={healthColorMap[data.health_status]}>{data.health_status}</Tag>
          </Descriptions.Item>
          <Descriptions.Item label="蜂王状态">
            <Tag color={queenColorMap[data.queen_status]}>{data.queen_status}</Tag>
          </Descriptions.Item>
          <Descriptions.Item label="最近检查">
            {data.last_inspection || '-'}
          </Descriptions.Item>
          <Descriptions.Item label="创建时间">{data.created_at}</Descriptions.Item>
          <Descriptions.Item label="更新时间">{data.updated_at}</Descriptions.Item>
          <Descriptions.Item label="备注" span={2}>{data.notes || '-'}</Descriptions.Item>
        </Descriptions>
      </Card>
    </div>
  )
}

export default BeehiveDetailPage
