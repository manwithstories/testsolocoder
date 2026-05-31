import { useEffect, useState } from 'react'
import { userApi } from '@/api/auth'
import { User, UserRole, UserStatus } from '@/types'

export default function AdminUsers() {
  const [users, setUsers] = useState<User[]>([])
  const [loading, setLoading] = useState(true)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [filters, setFilters] = useState({
    role: '',
    status: '',
    keyword: '',
  })

  useEffect(() => {
    loadUsers()
  }, [page, filters])

  const loadUsers = async () => {
    setLoading(true)
    try {
      const res = await userApi.list({ page, page_size: 10, ...filters })
      setUsers(res.data?.items || [])
      setTotal(res.data?.total || 0)
    } catch {
      setUsers([])
    } finally {
      setLoading(false)
    }
  }

  const handleStatusChange = async (id: number, status: UserStatus) => {
    try {
      await userApi.updateStatus(id, status)
      alert('状态更新成功')
      loadUsers()
    } catch (err: any) {
      alert(err.message || '更新失败')
    }
  }

  const handleRoleChange = async (id: number, role: UserRole) => {
    try {
      await userApi.updateRole(id, role)
      alert('角色更新成功')
      loadUsers()
    } catch (err: any) {
      alert(err.message || '更新失败')
    }
  }

  const handleDelete = async (id: number) => {
    if (!confirm('确定要删除该用户吗？')) return
    try {
      await userApi.delete(id)
      alert('删除成功')
      loadUsers()
    } catch (err: any) {
      alert(err.message || '删除失败')
    }
  }

  if (loading) {
    return <div className="text-center py-12">加载中...</div>
  }

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-gray-800">用户管理</h1>

      <div className="card p-4">
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <input
            type="text"
            placeholder="搜索用户名/邮箱/昵称"
            value={filters.keyword}
            onChange={(e) => setFilters({ ...filters, keyword: e.target.value })}
            className="input"
          />
          <select
            value={filters.role}
            onChange={(e) => setFilters({ ...filters, role: e.target.value })}
            className="input"
          >
            <option value="">全部角色</option>
            <option value="admin">管理员</option>
            <option value="roaster">烘焙师</option>
            <option value="user">普通用户</option>
          </select>
          <select
            value={filters.status}
            onChange={(e) => setFilters({ ...filters, status: e.target.value })}
            className="input"
          >
            <option value="">全部状态</option>
            <option value="active">活跃</option>
            <option value="disabled">禁用</option>
            <option value="pending">待审核</option>
          </select>
          <button onClick={() => setPage(1)} className="btn btn-primary">搜索</button>
        </div>
      </div>

      <div className="card overflow-hidden">
        <table className="w-full">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-4 py-3 text-left text-sm font-medium text-gray-600">用户</th>
              <th className="px-4 py-3 text-left text-sm font-medium text-gray-600">邮箱</th>
              <th className="px-4 py-3 text-left text-sm font-medium text-gray-600">角色</th>
              <th className="px-4 py-3 text-left text-sm font-medium text-gray-600">状态</th>
              <th className="px-4 py-3 text-left text-sm font-medium text-gray-600">认证</th>
              <th className="px-4 py-3 text-left text-sm font-medium text-gray-600">操作</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-100">
            {users.map((user) => (
              <tr key={user.id}>
                <td className="px-4 py-3">
                  <div className="flex items-center gap-3">
                    <div className="w-8 h-8 bg-coffee-200 rounded-full flex items-center justify-center text-sm">
                      {user.nickname?.[0] || user.username[0]}
                    </div>
                    <div>
                      <p className="font-medium">{user.nickname || user.username}</p>
                      <p className="text-xs text-gray-500">@{user.username}</p>
                    </div>
                  </div>
                </td>
                <td className="px-4 py-3 text-sm">{user.email}</td>
                <td className="px-4 py-3">
                  <select
                    value={user.role}
                    onChange={(e) => handleRoleChange(user.id, e.target.value as UserRole)}
                    className="input text-sm py-1"
                  >
                    <option value="admin">管理员</option>
                    <option value="roaster">烘焙师</option>
                    <option value="user">普通用户</option>
                  </select>
                </td>
                <td className="px-4 py-3">
                  <select
                    value={user.status}
                    onChange={(e) => handleStatusChange(user.id, e.target.value as UserStatus)}
                    className="input text-sm py-1"
                  >
                    <option value="active">活跃</option>
                    <option value="disabled">禁用</option>
                    <option value="pending">待审核</option>
                  </select>
                </td>
                <td className="px-4 py-3">
                  {user.is_certified ? (
                    <span className="badge bg-green-100 text-green-700">已认证</span>
                  ) : (
                    <span className="badge bg-gray-100 text-gray-700">未认证</span>
                  )}
                </td>
                <td className="px-4 py-3">
                  <button
                    onClick={() => handleDelete(user.id)}
                    className="text-red-500 hover:underline text-sm"
                  >
                    删除
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {total > 10 && (
        <div className="flex justify-center items-center gap-2">
          <button
            onClick={() => setPage(Math.max(1, page - 1))}
            disabled={page === 1}
            className="btn btn-secondary disabled:opacity-50"
          >
            上一页
          </button>
          <span className="text-gray-600">{page} / {Math.ceil(total / 10)}</span>
          <button
            onClick={() => setPage(page + 1)}
            disabled={page >= Math.ceil(total / 10)}
            className="btn btn-secondary disabled:opacity-50"
          >
            下一页
          </button>
        </div>
      )}
    </div>
  )
}
