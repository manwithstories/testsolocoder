import { Card, List, Avatar, Button, Tag, Modal, Input, message, Space } from 'antd'
import { CheckOutlined, CloseOutlined, UserOutlined } from '@ant-design/icons'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { userApi } from '@/services/auth'
import { User } from '@/types'

export function VerificationManagement() {
  const queryClient = useQueryClient()

  const { data, isLoading } = useQuery({
    queryKey: ['pending-verifications'],
    queryFn: () => userApi.getPendingVerifications({ page: 1, page_size: 50 }),
  })

  const verifyMutation = useMutation({
    mutationFn: ({ id, status, note }: { id: string; status: string; note?: string }) =>
      userApi.verifyProfessional(id, { status, note }),
    onSuccess: () => {
      message.success('操作成功')
      queryClient.invalidateQueries({ queryKey: ['pending-verifications'] })
    },
    onError: (error: any) => {
      message.error(error.message || '操作失败')
    },
  })

  const handleReject = (user: User) => {
    Modal.confirm({
      title: '拒绝认证',
      content: (
        <Input.TextArea
          rows={4}
          placeholder="请输入拒绝原因"
          onChange={(e) => {
            Modal.confirm({
              title: '确认拒绝',
              content: `拒绝原因：${e.target.value}`,
              onOk: () => verifyMutation.mutate({ id: user.id, status: 'rejected', note: e.target.value }),
            })
          }}
        />
      ),
    })
  }

  return (
    <div className="page-container">
      <h2 style={{ marginBottom: 24 }}>资质审核</h2>

      <Card>
        {data?.items && data.items.length > 0 ? (
          <List
            dataSource={data.items}
            loading={isLoading}
            renderItem={(user: User) => (
              <List.Item
                actions={[
                  <Space key="actions">
                    <Button
                      type="primary"
                      icon={<CheckOutlined />}
                      onClick={() => verifyMutation.mutate({ id: user.id, status: 'approved' })}
                      loading={verifyMutation.isPending}
                    >
                      通过
                    </Button>
                    <Button
                      danger
                      icon={<CloseOutlined />}
                      onClick={() => handleReject(user)}
                    >
                      拒绝
                    </Button>
                  </Space>,
                ]}
              >
                <List.Item.Meta
                  avatar={<Avatar icon={<UserOutlined />} src={user.avatar} />}
                  title={
                    <Space>
                      <span style={{ fontSize: 16, fontWeight: 500 }}>{user.full_name}</span>
                      <Tag color="orange">待审核</Tag>
                    </Space>
                  }
                  description={
                    <div>
                      <div>用户名：@{user.username}</div>
                      <div>邮箱：{user.email}</div>
                      <div>手机号：{user.phone || '未填写'}</div>
                      {user.verification_docs && (
                        <div style={{ marginTop: 8, padding: 12, background: '#f5f5f5', borderRadius: 4 }}>
                          <strong>资质证明：</strong>
                          <div style={{ marginTop: 4, whiteSpace: 'pre-wrap' }}>
                            {user.verification_docs}
                          </div>
                        </div>
                      )}
                    </div>
                  }
                />
              </List.Item>
            )}
          />
        ) : (
          <div className="empty-state">暂无待审核的专业人士</div>
        )}
      </Card>
    </div>
  )
}
