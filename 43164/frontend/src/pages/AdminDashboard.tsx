import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { toast } from 'sonner'
import { adminApi } from '@/services/api'
import { Users, DollarSign, Calendar, CheckCircle, TrendingUp } from 'lucide-react'

export default function AdminDashboard() {
  const navigate = useNavigate()
  const [stats, setStats] = useState<any>(null)
  const [pendingApprovals, setPendingApprovals] = useState<any[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    try {
      setLoading(true)
      const [statsRes, approvalsRes] = await Promise.all([
        adminApi.getStats(),
        adminApi.getPendingApprovals(),
      ])
      setStats(statsRes.data)
      setPendingApprovals(approvalsRes.data)
    } catch (error) {
      console.error('Failed to load data:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleApprove = async (id: string) => {
    try {
      await adminApi.approveTeacher(id)
      toast.success('老师审核通过')
      loadData()
    } catch (error: any) {
      toast.error(error.response?.data?.error || '审核失败')
    }
  }

  const handleReject = async (id: string) => {
    const reason = prompt('请输入拒绝原因:')
    if (reason) {
      try {
        await adminApi.rejectTeacher(id, { notes: reason })
        toast.success('已拒绝该老师申请')
        loadData()
      } catch (error: any) {
        toast.error(error.response?.data?.error || '操作失败')
      }
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">管理后台</h1>
        <p className="text-gray-500">管理平台运营数据和用户</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-5 gap-6">
        <div className="card">
          <div className="flex items-center gap-3">
            <div className="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center">
              <Users className="h-6 w-6 text-blue-600" />
            </div>
            <div>
              <div className="text-sm text-gray-500">总用户数</div>
              <div className="text-2xl font-bold text-gray-900">{stats?.totalUsers || 0}</div>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center gap-3">
            <div className="w-12 h-12 bg-green-100 rounded-lg flex items-center justify-center">
              <CheckCircle className="h-6 w-6 text-green-600" />
            </div>
            <div>
              <div className="text-sm text-gray-500">老师总数</div>
              <div className="text-2xl font-bold text-gray-900">{stats?.totalTeachers || 0}</div>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center gap-3">
            <div className="w-12 h-12 bg-purple-100 rounded-lg flex items-center justify-center">
              <Users className="h-6 w-6 text-purple-600" />
            </div>
            <div>
              <div className="text-sm text-gray-500">学生总数</div>
              <div className="text-2xl font-bold text-gray-900">{stats?.totalStudents || 0}</div>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center gap-3">
            <div className="w-12 h-12 bg-yellow-100 rounded-lg flex items-center justify-center">
              <Calendar className="h-6 w-6 text-yellow-600" />
            </div>
            <div>
              <div className="text-sm text-gray-500">课程总数</div>
              <div className="text-2xl font-bold text-gray-900">{stats?.totalBookings || 0}</div>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center gap-3">
            <div className="w-12 h-12 bg-green-100 rounded-lg flex items-center justify-center">
              <DollarSign className="h-6 w-6 text-green-600" />
            </div>
            <div>
              <div className="text-sm text-gray-500">平台收入</div>
              <div className="text-2xl font-bold text-gray-900">${stats?.totalRevenue?.toFixed(2) || '0.00'}</div>
            </div>
          </div>
        </div>
      </div>

      <div className="card">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-semibold text-gray-900">待审核老师 ({stats?.pendingApprovals || 0})</h2>
          <button
            onClick={() => loadData()}
            className="text-primary-600 hover:text-primary-700 text-sm"
          >
            刷新
          </button>
        </div>
        {pendingApprovals.length > 0 ? (
          <div className="space-y-4">
            {pendingApprovals.map((approval) => (
              <div key={approval.id} className="p-4 border border-gray-200 rounded-lg">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <div className="w-12 h-12 rounded-full bg-primary-600 flex items-center justify-center text-white text-lg font-medium">
                      {approval.user?.firstName?.[0]}{approval.user?.lastName?.[0]}
                    </div>
                    <div>
                      <div className="font-medium text-gray-900">
                        {approval.user?.firstName} {approval.user?.lastName}
                      </div>
                      <div className="text-sm text-gray-500">{approval.user?.email}</div>
                    </div>
                  </div>
                  <div className="flex gap-2">
                    <button
                      onClick={() => navigate(`/teachers/${approval.userId}`)}
                      className="btn-secondary"
                    >
                      查看详情
                    </button>
                    <button
                      onClick={() => handleApprove(approval.userId)}
                      className="btn-primary"
                    >
                      通过
                    </button>
                    <button
                      onClick={() => handleReject(approval.userId)}
                      className="btn-danger"
                    >
                      拒绝
                    </button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        ) : (
          <p className="text-gray-500 text-center py-8">暂无待审核的老师申请</p>
        )}
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div className="card">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">快捷操作</h2>
          <div className="space-y-3">
            <button
              onClick={() => navigate('/admin/dashboard')}
              className="w-full p-4 bg-blue-50 rounded-xl hover:bg-blue-100 transition-colors text-left flex items-center gap-3"
            >
              <Users className="h-6 w-6 text-blue-600" />
              <div>
                <div className="font-medium text-gray-900">用户管理</div>
                <div className="text-sm text-gray-500">查看和管理用户</div>
              </div>
            </button>
            <button
              onClick={() => navigate('/wallet')}
              className="w-full p-4 bg-green-50 rounded-xl hover:bg-green-100 transition-colors text-left flex items-center gap-3"
            >
              <DollarSign className="h-6 w-6 text-green-600" />
              <div>
                <div className="font-medium text-gray-900">财务管理</div>
                <div className="text-sm text-gray-500">处理提现申请</div>
              </div>
            </button>
            <button
              onClick={() => navigate('/profile')}
              className="w-full p-4 bg-purple-50 rounded-xl hover:bg-purple-100 transition-colors text-left flex items-center gap-3"
            >
              <TrendingUp className="h-6 w-6 text-purple-600" />
              <div>
                <div className="font-medium text-gray-900">系统设置</div>
                <div className="text-sm text-gray-500">配置系统参数</div>
              </div>
            </button>
          </div>
        </div>

        <div className="card">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">系统信息</h2>
          <div className="space-y-3">
            <div className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
              <span className="text-gray-600">平台佣金率</span>
              <span className="font-medium text-gray-900">10%</span>
            </div>
            <div className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
              <span className="text-gray-600">最低提现金额</span>
              <span className="font-medium text-gray-900">$100.00</span>
            </div>
            <div className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
              <span className="text-gray-600">支持的支付方式</span>
              <span className="font-medium text-gray-900">信用卡, 支付宝, 微信</span>
            </div>
            <div className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
              <span className="text-gray-600">默认时区</span>
              <span className="font-medium text-gray-900">UTC</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
