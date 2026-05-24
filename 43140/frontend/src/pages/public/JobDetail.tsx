import { useState, useEffect } from 'react'
import { Link, useParams, useNavigate } from 'react-router-dom'
import { MapPin, Briefcase, DollarSign, Building2, ArrowLeft, FileText, Clock } from 'lucide-react'
import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { getJobById, selectJobs } from '@/store/slices/jobSlice'
import { applicationApi, resumeApi, Resume } from '@/api/jobs'
import { selectAuth } from '@/store/slices/authSlice'
import dayjs from 'dayjs'

export default function JobDetail() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const dispatch = useAppDispatch()
  const { currentJob, loading } = useAppSelector(selectJobs)
  const { user } = useAppSelector(selectAuth)
  const [resumes, setResumes] = useState<Resume[]>([])
  const [selectedResume, setSelectedResume] = useState<number | null>(null)
  const [showApplyModal, setShowApplyModal] = useState(false)
  const [applying, setApplying] = useState(false)
  const [coverLetter, setCoverLetter] = useState('')
  const [error, setError] = useState('')

  useEffect(() => {
    if (id) {
      dispatch(getJobById(parseInt(id)))
    }
  }, [id, dispatch])

  const handleApply = () => {
    if (!user) {
      navigate('/login')
      return
    }
    if (user.role !== 'jobseeker') {
      setError('Only job seekers can apply for jobs')
      return
    }
    setShowApplyModal(true)
    loadResumes()
  }

  const loadResumes = async () => {
    try {
      const response = await resumeApi.list()
      setResumes(response.data)
      const defaultResume = response.data.find((r) => r.is_default)
      if (defaultResume) setSelectedResume(defaultResume.id)
    } catch {
      setError('Failed to load resumes')
    }
  }

  const submitApplication = async () => {
    if (!selectedResume) {
      setError('Please select a resume')
      return
    }
    setApplying(true)
    setError('')
    try {
      await applicationApi.create({
        job_id: parseInt(id!),
        resume_id: selectedResume,
        cover_letter: coverLetter,
      })
      setShowApplyModal(false)
      alert('Application submitted successfully!')
    } catch (err: any) {
      setError(err.message || 'Failed to submit application')
    } finally {
      setApplying(false)
    }
  }

  if (loading || !currentJob) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin w-8 h-8 border-4 border-primary-600 border-t-transparent rounded-full"></div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white shadow-sm">
        <div className="max-w-4xl mx-auto px-4 py-4">
          <Link
            to="/jobs"
            className="flex items-center gap-2 text-gray-600 hover:text-primary-600"
          >
            <ArrowLeft className="w-5 h-5" />
            Back to Jobs
          </Link>
        </div>
      </header>

      <div className="max-w-4xl mx-auto px-4 py-8">
        <div className="bg-white rounded-lg shadow-md p-8 mb-6">
          <div className="flex items-start justify-between">
            <div>
              <h1 className="text-3xl font-bold text-gray-900">{currentJob.title}</h1>
              <p className="text-xl text-gray-600 mt-1">
                {currentJob.company?.company_name || 'Unknown Company'}
              </p>
            </div>
            <button onClick={handleApply} className="btn-primary text-lg px-8">
              Apply Now
            </button>
          </div>

          <div className="flex flex-wrap items-center gap-6 mt-6 text-gray-600">
            <span className="flex items-center gap-2">
              <MapPin className="w-5 h-5" />
              {currentJob.location}
            </span>
            <span className="flex items-center gap-2">
              <Briefcase className="w-5 h-5" />
              {currentJob.job_type}
            </span>
            {currentJob.salary_min > 0 && (
              <span className="flex items-center gap-2 text-green-600 font-medium">
                <DollarSign className="w-5 h-5" />
                ${currentJob.salary_min.toLocaleString()} - ${currentJob.salary_max.toLocaleString()}
              </span>
            )}
            <span className="flex items-center gap-2">
              <Clock className="w-5 h-5" />
              Posted {dayjs(currentJob.created_at).format('MMM D, YYYY')}
            </span>
          </div>
        </div>

        <div className="bg-white rounded-lg shadow-md p-8 mb-6">
          <h2 className="text-xl font-bold mb-4">Job Description</h2>
          <div className="prose max-w-none text-gray-700 whitespace-pre-wrap">
            {currentJob.description}
          </div>
        </div>

        {currentJob.requirements && (
          <div className="bg-white rounded-lg shadow-md p-8 mb-6">
            <h2 className="text-xl font-bold mb-4">Requirements</h2>
            <div className="prose max-w-none text-gray-700 whitespace-pre-wrap">
              {currentJob.requirements}
            </div>
          </div>
        )}

        {currentJob.skills && (
          <div className="bg-white rounded-lg shadow-md p-8">
            <h2 className="text-xl font-bold mb-4">Skills</h2>
            <div className="flex flex-wrap gap-2">
              {currentJob.skills.split(',').map((skill, index) => (
                <span
                  key={index}
                  className="px-3 py-1 bg-primary-100 text-primary-700 rounded-full text-sm"
                >
                  {skill.trim()}
                </span>
              ))}
            </div>
          </div>
        )}

        {error && (
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mt-4">
            {error}
          </div>
        )}
      </div>

      {showApplyModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-lg shadow-xl max-w-lg w-full max-h-[90vh] overflow-y-auto">
            <div className="p-6">
              <h2 className="text-2xl font-bold mb-4">Apply for {currentJob.title}</h2>

              {resumes.length === 0 ? (
                <div className="text-center py-8">
                  <FileText className="w-16 h-16 text-gray-300 mx-auto mb-4" />
                  <p className="text-gray-500 mb-4">You don't have any resumes yet.</p>
                  <Link to="/jobseeker/resumes" className="btn-primary">
                    Create a Resume
                  </Link>
                </div>
              ) : (
                <>
                  <div className="mb-4">
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      Select Resume
                    </label>
                    <select
                      value={selectedResume || ''}
                      onChange={(e) => setSelectedResume(parseInt(e.target.value))}
                      className="input-field"
                    >
                      <option value="">Select a resume...</option>
                      {resumes.map((resume) => (
                        <option key={resume.id} value={resume.id}>
                          {resume.title} ({resume.full_name})
                          {resume.is_default && ' - Default'}
                        </option>
                      ))}
                    </select>
                  </div>

                  <div className="mb-4">
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      Cover Letter (Optional)
                    </label>
                    <textarea
                      value={coverLetter}
                      onChange={(e) => setCoverLetter(e.target.value)}
                      rows={5}
                      className="input-field"
                      placeholder="Write a brief introduction..."
                    />
                  </div>

                  <div className="flex gap-3">
                    <button
                      onClick={() => setShowApplyModal(false)}
                      className="btn-secondary flex-1"
                    >
                      Cancel
                    </button>
                    <button
                      onClick={submitApplication}
                      disabled={applying}
                      className="btn-primary flex-1 disabled:opacity-50"
                    >
                      {applying ? 'Submitting...' : 'Submit Application'}
                    </button>
                  </div>
                </>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
