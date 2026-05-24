import { useState, useEffect } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { Briefcase, FileText, Users, Eye, TrendingUp, BarChart3 } from 'lucide-react'
import { statisticsApi } from '@/api/jobs'

interface Stats {
  total_jobs: number
  open_jobs: number
  total_applications: number
  total_interviews: number
  total_offers: number
  interview_rate: number
  offer_rate: number
  total_views: number
  conversion_rate: number
}

export default function CompanyDashboard() {
  const navigate = useNavigate()
  const [stats, setStats] = useState<Stats | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadStats()
  }, [])

  const loadStats = async () => {
    try {
      const response = await statisticsApi.getCompanyStats()
      setStats(response.data)
    } catch (err) {
      console.error('Failed to load stats:', err)
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin w-8 h-8 border-4 border-primary-600 border-t-transparent rounded-full"></div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-gray-900">Dashboard</h1>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <div className="card">
          <div className="flex items-center gap-3">
            <div className="p-3 bg-blue-100 rounded-lg">
              <Briefcase className="w-6 h-6 text-blue-600" />
            </div>
            <div>
              <p className="text-sm text-gray-500">Total Jobs</p>
              <p className="text-2xl font-bold">{stats?.total_jobs || 0}</p>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center gap-3">
            <div className="p-3 bg-green-100 rounded-lg">
              <Eye className="w-6 h-6 text-green-600" />
            </div>
            <div>
              <p className="text-sm text-gray-500">Total Views</p>
              <p className="text-2xl font-bold">{stats?.total_views || 0}</p>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center gap-3">
            <div className="p-3 bg-purple-100 rounded-lg">
              <FileText className="w-6 h-6 text-purple-600" />
            </div>
            <div>
              <p className="text-sm text-gray-500">Applications</p>
              <p className="text-2xl font-bold">{stats?.total_applications || 0}</p>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center gap-3">
            <div className="p-3 bg-orange-100 rounded-lg">
              <Users className="w-6 h-6 text-orange-600" />
            </div>
            <div>
              <p className="text-sm text-gray-500">Interviews</p>
              <p className="text-2xl font-bold">{stats?.total_interviews || 0}</p>
            </div>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
        <div className="card">
          <div className="flex items-center gap-2 mb-4">
            <TrendingUp className="w-5 h-5 text-primary-600" />
            <h3 className="font-semibold">Conversion Rate</h3>
          </div>
          <p className="text-3xl font-bold text-primary-600">
            {stats?.conversion_rate?.toFixed(1) || 0}%
          </p>
          <p className="text-sm text-gray-500 mt-1">Views to Applications</p>
        </div>

        <div className="card">
          <div className="flex items-center gap-2 mb-4">
            <BarChart3 className="w-5 h-5 text-green-600" />
            <h3 className="font-semibold">Interview Rate</h3>
          </div>
          <p className="text-3xl font-bold text-green-600">
            {stats?.interview_rate?.toFixed(1) || 0}%
          </p>
          <p className="text-sm text-gray-500 mt-1">Applications to Interviews</p>
        </div>

        <div className="card">
          <div className="flex items-center gap-2 mb-4">
            <TrendingUp className="w-5 h-5 text-purple-600" />
            <h3 className="font-semibold">Offer Rate</h3>
          </div>
          <p className="text-3xl font-bold text-purple-600">
            {stats?.offer_rate?.toFixed(1) || 0}%
          </p>
          <p className="text-sm text-gray-500 mt-1">Applications to Offers</p>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <Link to="/company/jobs" className="card hover:shadow-lg transition-shadow">
          <h3 className="font-semibold text-lg mb-2">Manage Jobs</h3>
          <p className="text-gray-500">Create, edit, and manage your job postings</p>
        </Link>
        <Link to="/company/applications" className="card hover:shadow-lg transition-shadow">
          <h3 className="font-semibold text-lg mb-2">Review Applications</h3>
          <p className="text-gray-500">View and manage incoming applications</p>
        </Link>
      </div>
    </div>
  )
}
