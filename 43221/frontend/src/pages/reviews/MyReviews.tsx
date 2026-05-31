import { useState } from 'react'
import { Card, List, Avatar, Rate, Tag, Pagination, Empty } from 'antd'
import { UserOutlined } from '@ant-design/icons'
import { useQuery } from '@tanstack/react-query'
import { useAuthContext } from '@/contexts/AuthContext'
import { recordApi } from '@/services/record'
import { Review } from '@/types'

export function MyReviews() {
  const { user } = useAuthContext()
  const [page, setPage] = useState(1)
  const [pageSize] = useState(10)

  const { data, isLoading } = useQuery({
    queryKey: ['my-reviews', page, pageSize, user?.role],
    queryFn: () => {
      if (user?.role === 'professional') {
        return recordApi.getProfessionalReviews({ page, page_size: pageSize, status: 'approved' })
      }
      return recordApi.getProfessionalReviews({ page, page_size: pageSize })
    },
    enabled: !!user,
  })

  return (
    <div className="page-container">
      <h2 style={{ marginBottom: 24 }}>
        {user?.role === 'professional' ? '我的评价' : '我的评价'}
      </h2>

      <Card>
        {data?.items && data.items.length > 0 ? (
          <List
            dataSource={data.items}
            loading={isLoading}
            renderItem={(review: Review) => (
              <List.Item key={review.id}>
                <List.Item.Meta
                  avatar={<Avatar icon={<UserOutlined />} src={review.client?.avatar} />}
                  title={
                    <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                      <span>{review.client?.full_name || '匿名用户'}</span>
                      <Rate disabled value={review.rating} style={{ fontSize: 14 }} />
                      {user?.role === 'admin' && (
                        <Tag color={
                          review.status === 'approved' ? 'green' :
                          review.status === 'pending' ? 'orange' : 'red'
                        }>
                          {review.status === 'approved' ? '已通过' :
                           review.status === 'pending' ? '待审核' : '已拒绝'}
                        </Tag>
                      )}
                    </div>
                  }
                  description={
                    <div>
                      <div style={{ color: '#666', marginBottom: 8 }}>
                        服务：{review.service?.title}
                      </div>
                      {review.content && <div>{review.content}</div>}
                      <div style={{ color: '#999', fontSize: 12, marginTop: 8 }}>
                        {new Date(review.created_at).toLocaleString()}
                      </div>
                    </div>
                  }
                />
              </List.Item>
            )}
          />
        ) : (
          <Empty description={isLoading ? '加载中...' : '暂无评价'} />
        )}

        {data && data.total > pageSize && (
          <div style={{ marginTop: 16, textAlign: 'center' }}>
            <Pagination
              current={page}
              pageSize={pageSize}
              total={data.total}
              onChange={setPage}
            />
          </div>
        )}
      </Card>
    </div>
  )
}
