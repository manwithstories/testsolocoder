import { useState, useEffect } from 'react'
import { Link, useSearchParams } from 'react-router-dom'
import { Search, MapPin, Briefcase, Filter, ChevronLeft, ChevronRight } from 'lucide-react'
import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { searchJobs, selectJobs } from '@/store/slices/jobSlice'
import { Job } from '@/api/jobs'
import dayjs from 'dayjs'

export default function PublicJobSearch() {
  const dispatch = useAppDispatch()
  const { jobs, loading, pagination } = useAppSelector(selectJobs)
  const [searchParams, setSearchParams] = useSearchParams()

  const [keyword, setKeyword] = useState(searchParams.get('keyword') || '')
  const [location, setLocation] = useState(searchParams.get('location') || '')
  const [jobType, setJobType] = useState(searchParams.get('job_type') || '')
  const [salaryMin, setSalaryMin] = useState(searchParams.get('salary_min') || '')
  const [showFilters, setShowFilters] = useState(false)

  const fetchJobs = (page = 1) => {
    const params: Record<string, string> = {
      page: page.toString(),
    }
    if (keyword) params.keyword = keyword
    if (location) params.location = location
    if (jobType) params.job_type = jobType
    if (salaryMin) params.salary_min = salaryMin

    dispatch(searchJobs(params))
  }

  useEffect(() => {
    fetchJobs(1)
  }, [])

  const handleSearch = () => {
    fetchJobs(1)
  }

  const handlePageChange = (page: number) => {
    fetchJobs(page)
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4 py-4">
          <div className="flex items-center justify-between">
            <Link to="/" className="flex items-center gap-2">
              <Briefcase className="w-8 h-8 text-primary-600" />
              <span className="text-xl font-bold">RecruitPro</span>
            </Link>
            <div className="flex items-center gap-4">
              <Link to="/login" className="text-gray-600 hover:text-primary-600">
                Sign In
              </Link>
              <Link to="/register" className="btn-primary">
                Sign Up
              </Link>
            </div>
          </div>
        </div>
      </header>

      <div className="max-w-7xl mx-auto px-4 py-8">
        <div className="bg-white rounded-lg shadow-md p-6 mb-6">
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
            <div className="mt-4 pt-4 border-t grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
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
              <div className="md:col-span-2 flex items-end gap-2">
                <button onClick={handleSearch} className="btn-primary flex-1">
                  Apply Filters
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

        {loading ? (
          <div className="text-center py-12">
            <div className="animate-spin w-8 h-8 border-4 border-primary-600 border-t-transparent rounded-full mx-auto"></div>
            <p className="mt-4 text-gray-500">Loading jobs...</p>
          </div>
        ) : jobs.length === 0 ? (
          <div className="text-center py-12">
            <Briefcase className="w-16 h-16 text-gray-300 mx-auto" />
            <p className="mt-4 text-gray-500">No jobs found matching your criteria</p>
          </div>
        ) : (
          <>
            <div className="grid gap-4">
              {jobs.map((job: Job) => (
                <Link
                  key={job.id}
                  to={`/jobs/${job.id}`}
                  className="bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow"
                >
                  <div className="flex items-start justify-between">
                    <div className="flex-1">
                      <h3 className="text-lg font-semibold text-gray-900">{job.title}</h3>
                      <p className="text-gray-600">{job.company?.company_name || 'Unknown Company'}</p>
                      <div className="flex items-center gap-4 mt-2 text-sm text-gray-500">
                        <span className="flex items-center gap-1">
                          <MapPin className="w-4 h-4" />
                          {job.location}
                        </span>
                        <span className="flex items-center gap-1">
                          <Briefcase className="w-4 h-4" />
                          {job.job_type}
                        </span>
                        {job.salary_min > 0 && (
                          <span className="text-green-600 font-medium">
                            ${job.salary_min.toLocaleString()} - ${job.salary_max.toLocaleString()}
                          </span>
                        )}
                      </div>
                    </div>
                    <span className="text-sm text-gray-400">
                      {dayjs(job.created_at).fromNow()}
                    </span>
                  </div>
                </Link>
              ))}
            </div>

            {pagination.total_pages > 1 && (
              <div className="flex items-center justify-center gap-2 mt-8">
                <button
                  onClick={() => handlePageChange(pagination.page - 1)}
                  disabled={pagination.page === 1}
                  className="p-2 rounded-lg hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  <ChevronLeft className="w-5 h-5" />
                </button>
                {Array.from({ length: pagination.total_pages }, (_, i) => i + 1).map((page) => (
                  <button
                    key={page}
                    onClick={() => handlePageChange(page)}
                    className={`w-10 h-10 rounded-lg font-medium ${
                      page === pagination.page
                        ? 'bg-primary-600 text-white'
                        : 'hover:bg-gray-100'
                    }`}
                  >
                    {page}
                  </button>
                ))}
                <button
                  onClick={() => handlePageChange(pagination.page + 1)}
                  disabled={pagination.page === pagination.total_pages}
                  className="p-2 rounded-lg hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  <ChevronRight className="w-5 h-5" />
                </button>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  )
}
