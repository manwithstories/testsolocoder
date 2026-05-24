import { useState, useEffect } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { Search, MapPin, Filter, Star, Check } from 'lucide-react'
import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { searchJobs, selectJobs } from '@/store/slices/jobSlice'
import { recommendationApi, Job } from '@/api/jobs'
import dayjs from 'dayjs'

export default function JobSeekerJobs() {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()
  const { jobs, loading, pagination } = useAppSelector(selectJobs)
  const [recommendations, setRecommendations] = useState<(Job & { match_score: number; matched_skills: string[] })[]>([])
  const [showRecommendations, setShowRecommendations] = useState(true)
  const [keyword, setKeyword] = useState('')
  const [location, setLocation] = useState('')
  const [jobType, setJobType] = useState('')
  const [salaryMin, setSalaryMin] = useState('')
  const [showFilters, setShowFilters] = useState(false)

  useEffect(() => {
    fetchJobs()
    loadRecommendations()
  }, [])

  const fetchJobs = (page = 1) => {
    const params: Record<string, string> = { page: page.toString() }
    if (keyword) params.keyword = keyword
    if (location) params.location = location
    if (jobType) params.job_type = jobType
    if (salaryMin) params.salary_min = salaryMin
    dispatch(searchJobs(params))
  }

  const loadRecommendations = async () => {
    try {
      const response = await recommendationApi.getRecommended()
      setRecommendations(response.data)
    } catch {
      console.log('No recommendations available')
    }
  }

  const handleSearch = () => {
    fetchJobs(1)
  }

  const handleApply = (jobId: number) => {
    navigate(`/jobs/${jobId}`)
  }

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-gray-900">Job Search</h1>

      <div className="card">
        <div className="flex flex-col md:flex-row gap-4">
          <div className="flex-1">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
              <input
                type="text"
                value={keyword}
                onChange={(e) => setKeyword(e.target.value)}
                placeholder="Search jobs, keywords, or companies"
                className="input-field pl-10"
                onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
              />
            </div>
          </div>
          <div className="flex-1">
            <div className="relative">
              <MapPin className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
              <input
                type="text"
                value={location}
                onChange={(e) => setLocation(e.target.value)}
                placeholder="Location"
                className="input-field pl-10"
                onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
              />
            </div>
          </div>
          <button onClick={handleSearch} className="btn-primary">
            Search
          </button>
          <button
            onClick={() => setShowFilters(!showFilters)}
            className="btn-secondary flex items-center gap-2"
          >
            <Filter className="w-4 h-4" />
            Filters
          </button>
        </div>

        {showFilters && (
          <div className="mt-4 pt-4 border-t grid grid-cols-1 md:grid-cols-3 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Job Type</label>
              <select
                value={jobType}
                onChange={(e) => setJobType(e.target.value)}
                className="input-field"
              >
                <option value="">All Types</option>
                <option value="full-time">Full Time</option>
                <option value="part-time">Part Time</option>
                <option value="contract">Contract</option>
                <option value="internship">Internship</option>
                <option value="remote">Remote</option>
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Min Salary</label>
              <input
                type="number"
                value={salaryMin}
                onChange={(e) => setSalaryMin(e.target.value)}
                placeholder="Min salary"
                className="input-field"
              />
            </div>
            <div className="flex items-end gap-2">
              <button onClick={handleSearch} className="btn-primary flex-1">
                Apply
              </button>
              <button
                onClick={() => {
                  setJobType('')
                  setSalaryMin('')
                }}
                className="btn-secondary"
              >
                Clear
              </button>
            </div>
          </div>
        )}
      </div>

      {showRecommendations && recommendations.length > 0 && (
        <div className="card bg-gradient-to-r from-primary-50 to-primary-100">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-lg font-semibold flex items-center gap-2">
              <Star className="w-5 h-5 text-yellow-500" />
              Recommended for You
            </h2>
            <button
              onClick={() => setShowRecommendations(false)}
              className="text-sm text-gray-500 hover:text-gray-700"
            >
              Hide
            </button>
          </div>
          <div className="space-y-2">
            {recommendations.slice(0, 5).map((job) => (
              <Link
                key={job.id}
                to={`/jobs/${job.id}`}
                className="block bg-white rounded-lg p-4 hover:shadow-md transition-shadow"
              >
                <div className="flex items-center justify-between">
                  <div>
                    <h3 className="font-medium">{job.title}</h3>
                    <p className="text-sm text-gray-500">{job.company?.company_name}</p>
                  </div>
                  {job.match_score > 0 && (
                    <div className="text-right">
                      <span className="text-xs text-gray-500">Match</span>
                      <p className="font-semibold text-green-600">{job.match_score * 20}%</p>
                    </div>
                  )}
                </div>
              </Link>
            ))}
          </div>
        </div>
      )}

      {loading ? (
        <div className="text-center py-12">
          <div className="animate-spin w-8 h-8 border-4 border-primary-600 border-t-transparent rounded-full mx-auto"></div>
        </div>
      ) : jobs.length === 0 ? (
        <div className="text-center py-12">
          <Search className="w-16 h-16 text-gray-300 mx-auto" />
          <p className="mt-4 text-gray-500">No jobs found</p>
        </div>
      ) : (
        <div className="space-y-4">
          {jobs.map((job: Job) => (
            <div key={job.id} className="card">
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <h3 className="text-lg font-semibold">{job.title}</h3>
                  <p className="text-gray-600">{job.company?.company_name}</p>
                  <div className="flex items-center gap-4 mt-2 text-sm text-gray-500">
                    <span className="flex items-center gap-1">
                      <MapPin className="w-4 h-4" />
                      {job.location}
                    </span>
                    <span className="flex items-center gap-1">
                      <Filter className="w-4 h-4" />
                      {job.job_type}
                    </span>
                    {job.salary_min > 0 && (
                      <span className="text-green-600 font-medium">
                        ${job.salary_min.toLocaleString()} - ${job.salary_max.toLocaleString()}
                      </span>
                    )}
                    <span className="text-xs text-gray-400">
                      {dayjs(job.created_at).fromNow()}
                    </span>
                  </div>
                </div>
                <div className="flex items-center gap-2">
                  <Link to={`/jobs/${job.id}`} className="btn-secondary">
                    View Details
                  </Link>
                  <button onClick={() => handleApply(job.id)} className="btn-primary">
                    Apply
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      {pagination.total_pages > 1 && (
        <div className="flex items-center justify-center gap-2">
          {Array.from({ length: Math.min(pagination.total_pages, 5) }, (_, i) => i + 1).map(
            (page) => (
              <button
                key={page}
                onClick={() => fetchJobs(page)}
                className={`w-10 h-10 rounded-lg font-medium ${
                  page === pagination.page
                    ? 'bg-primary-600 text-white'
                    : 'hover:bg-gray-100'
                }`}
              >
                {page}
              </button>
            )
          )}
        </div>
      )}
    </div>
  )
}
