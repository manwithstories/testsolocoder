import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { FileText, Briefcase, Calendar, Star, TrendingUp, BarChart3 } from 'lucide-react'
import { statisticsApi } from '@/api/jobs'

interface Stats {
  total_applications: number
  interview_count: number
  offer_count: number
  rejection_count: number
  pending_count: number
  success_rate: number
  resume_count: number
}

export default function JobSeekerDashboard() {
  const [stats, setStats] = useState<Stats | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadStats()
  }, [])

  const loadStats = async () => {
    try {
      const response = await statisticsApi.getJobSeekerStats()
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
        <div className="animate-spin w-8 h-8 border-4 border-primary-600 border-t-transparent rounded-full mx-auto"></div>
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
              <p className="text-sm text-gray-500">Applications</p>
              <p className="text-2xl font-bold">{stats?.total_applications || 0}</p>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center gap-3">
            <div className="p-3 bg-purple-100 rounded-lg">
              <Calendar className="w-6 h-6 text-purple-600" />
            </div>
            <div>
              <p className="text-sm text-gray-500">Interviews</p>
              <p className="text-2xl font-bold">{stats?.interview_count || 0}</p>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center gap-3">
            <div className="p-3 bg-green-100 rounded-lg">
              <Star className="w-6 h-6 text-green-600" />
            </div>
            <div>
              <p className="text-sm text-gray-500">Offers</p>
              <p className="text-2xl font-bold">{stats?.offer_count || 0}</p>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center gap-3">
            <div className="p-3 bg-orange-100 rounded-lg">
              <FileText className="w-6 h-6 text-orange-600" />
            </div>
            <div>
              <p className="text-sm text-gray-500">Resumes</p>
              <p className="text-2xl font-bold">{stats?.resume_count || 0}</p>
            </div>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <div className="card">
          <div className="flex items-center gap-2 mb-4">
            <TrendingUp className="w-5 h-5 text-primary-600" />
            <h3 className="font-semibold">Success Rate</h3>
          </div>
          <p className="text-3xl font-bold text-primary-600">
            {stats?.success_rate?.toFixed(1) || 0}%
          </p>
          <p className="text-sm text-gray-500 mt-1">Applications to Offers</p>
        </div>

        <div className="card">
          <div className="flex items-center gap-2 mb-4">
            <BarChart3 className="w-5 h-5 text-yellow-600" />
            <h3 className="font-semibold">Pending Applications</h3>
          </div>
          <p className="text-3xl font-bold text-yellow-600">{stats?.pending_count || 0}</p>
          <p className="text-sm text-gray-500 mt-1">Awaiting response</p>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <Link to="/jobseeker/jobs" className="card hover:shadow-lg transition-shadow">
          <h3 className="font-semibold text-lg mb-2">Browse Jobs</h3>
          <p className="text-gray-500">Search and apply for new opportunities</p>
        </Link>
        <Link to="/jobseeker/resumes" className="card hover:shadow-lg transition-shadow">
          <h3 className="font-semibold text-lg mb-2">Manage Resumes</h3>
          <p className="text-gray-500">Create and edit your professional resumes</p>
        </Link>
      </div>
    </div>
  )
}
