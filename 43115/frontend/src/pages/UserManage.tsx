import React, { useState, useEffect } from 'react'
import { Table, Card, Button, Tag, Modal, Input, Select, message, Avatar } from 'antd'
import { userApi } from '@/services/auth'
import { User } from '@/types'
import { formatDate, getRoleText, getCertificationStatusText } from '@/utils'
import { usePagination } from '@/hooks'

const UserManage: React.FC = () => {
  const { page, pageSize, total, setPage, setPageSize, setTotal } = usePagination()
  const [users, setUsers] = useState<User[]>([])
  const [loading, setLoading] = useState(false)
  const [role, setRole] = useState<string | undefined>()
  const [certModalVisible, setCertModalVisible] = useState(false)
  const [selectedCert, setSelectedCert] = useState<any>(null)

  useEffect(() => {
    loadUsers()
  }, [page, pageSize, role])

  const loadUsers = async () => {
    setLoading(true)
    try {
      const res = await userApi.getList({
        page,
        page_size: pageSize,
        role,
      })
      setUsers(res.list)
      setTotal(res.total)
    } catch (error) {
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  const handleApproveCert = async (userId: number) => {
    try {
      await userApi.reviewCertification(userId, { approved: true })
      message.success('审核通过')
      loadUsers()
    } catch (error) {
      console.error(error)
    }
  }

  const handleRejectCert = async (userId: number, reason: string) => {
    try {
      await userApi.reviewCertification(userId, { approved: false, reject_reason: reason })
      message.success('已拒绝')
      setCertModalVisible(false)
      setSelectedCert(null)
      loadUsers()
    } catch (error) {
      console.error(error)
    }
  }

  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
    },
    {
      title: '头像',
      dataIndex: 'avatar',
      key: 'avatar',
      render: (text: string, record: User) => (
        <Avatar src={text}>{record.nickname?.charAt(0)}</Avatar>
      ),
    },
    {
      title: '昵称',
      dataIndex: 'nickname',
      key: 'nickname',
    },
    {
      title: '手机号',
      dataIndex: 'phone',
      key: 'phone',
    },
    {
      title: '角色',
      dataIndex: 'role',
      key: 'role',
      render: (text: string) => <Tag color="blue">{getRoleText(text)}</Tag>,
    },
    {
      title: '认证状态',
      dataIndex: 'provider_status',
      key: 'cert_status',
      render: (_: any, record: User) =>
        record.role === 'service_provider' ? (
          <Tag color={record.provider_status === 'approved' ? 'green' : record.provider_status === 'pending' ? 'orange' : 'red'}>
            {getCertificationStatusText(record.provider_status)}
          </Tag>
        ) : (
          '-'
        ),
    },
    {
      title: '注册时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (text: string) => formatDate(text),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: User) => (
        <div style={{ display: 'flex', gap: 8 }}>
          {record.role === 'service_provider' && record.provider_status === 'pending' && (
            <>
              <Button type="link" size="small" onClick={() => handleApproveCert(record.id)}>
                通过
              </Button>
              <Button
                type="link"
                size="small"
                danger
                onClick={() => {
                  setSelectedCert(record)
                  setCertModalVisible(true)
                }}
              >
                拒绝
              </Button>
            </>
          )}
        </div>
      ),
    },
  ]

  return (
    <div>
      <div className="page-header">
        <h1 className="page-title">用户管理</h1>
      </div>

      <div className="filter-bar">
        <Select
          placeholder="用户角色"
          style={{ width: 150 }}
          allowClear
          value={role}
          onChange={(value) => {
            setRole(value)
            setPage(1)
          }}
          options={[
            { label: '客户', value: 'customer' },
            { label: '服务人员', value: 'service_provider' },
            { label: '管理员', value: 'admin' },
          ]}
        />
      </div>

      <Card>
        <Table
          rowKey="id"
          loading={loading}
          dataSource={users}
          columns={columns}
          pagination={{
            current: page,
            pageSize,
            total,
            showSizeChanger: true,
            onChange: (p, ps) => {
              setPage(p)
              setPageSize(ps)
            },
          }}
        />
      </Card>

      <Modal
        title="拒绝认证"
        open={certModalVisible}
        onOk={() => handleRejectCert(selectedCert?.id, '不符合要求')}
        onCancel={() => setCertModalVisible(false)}
      >
        <p>确定要拒绝该用户的认证申请吗？</p>
        <Input.TextArea rows={3} placeholder="请输入拒绝原因" />
      </Modal>
    </div>
  )
}

export default UserManage
