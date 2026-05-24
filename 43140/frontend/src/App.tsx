import { useEffect } from 'react'
import { Routes, Route, Navigate } from 'react-router-dom'
import { useAppSelector, useAppDispatch } from '@/hooks/redux'
import { selectAuth, setUser } from '@/store/slices/authSlice'
import Login from '@/pages/auth/Login'
import Register from '@/pages/auth/Register'
import CompanyLayout from '@/layouts/CompanyLayout'
import JobSeekerLayout from '@/layouts/JobSeekerLayout'
import AdminLayout from '@/layouts/AdminLayout'
import PublicJobSearch from '@/pages/public/JobSearch'
import JobDetail from '@/pages/public/JobDetail'
import RequireAuth from '@/components/auth/RequireAuth'
import CompanyDashboard from '@/pages/company/Dashboard'
import CompanyJobs from '@/pages/company/Jobs'
import CompanyApplications from '@/pages/company/Applications'
import CompanyInterviews from '@/pages/company/Interviews'
import CompanyReviews from '@/pages/company/Reviews'
import CompanyExport from '@/pages/company/Export'
import CompanyProfile from '@/pages/company/Profile'
import JobSeekerDashboard from '@/pages/jobseeker/Dashboard'
import JobSeekerJobs from '@/pages/jobseeker/Jobs'
import JobSeekerResumes from '@/pages/jobseeker/Resumes'
import JobSeekerApplications from '@/pages/jobseeker/Applications'
import JobSeekerInterviews from '@/pages/jobseeker/Interviews'
import JobSeekerReviews from '@/pages/jobseeker/Reviews'
import JobSeekerProfile from '@/pages/jobseeker/Profile'
import AdminDashboard from '@/pages/admin/Dashboard'
import AdminUsers from '@/pages/admin/Users'

function App() {
  const { user, token } = useAppSelector(selectAuth)
  const dispatch = useAppDispatch()

  useEffect(() => {
    const savedUser = localStorage.getItem('user')
    const savedToken = localStorage.getItem('token')
    if (savedToken && !user && savedUser) {
      try {
        dispatch(setUser(JSON.parse(savedUser)))
      } catch {
        localStorage.removeItem('user')
        localStorage.removeItem('token')
      }
    }
  }, [user, dispatch])

  const getHomeRoute = () => {
    if (!user || !token) return '/login'
    if (user.role === 'company') return '/company/dashboard'
    if (user.role === 'jobseeker') return '/jobseeker/dashboard'
    if (user.role === 'admin') return '/admin/dashboard'
    return '/login'
  }

  return (
    <Routes>
      <Route path="/" element={<Navigate to={getHomeRoute()} replace />} />
      <Route path="/login" element={<Login />} />
      <Route path="/register" element={<Register />} />
      <Route path="/jobs" element={<PublicJobSearch />} />
      <Route path="/jobs/:id" element={<JobDetail />} />

      <Route
        path="/company"
        element={
          <RequireAuth roles={['company', 'admin']}>
            <CompanyLayout />
          </RequireAuth>
        }
      >
        <Route path="" element={<Navigate to="dashboard" replace />} />
        <Route path="dashboard" element={<CompanyDashboard />} />
        <Route path="jobs" element={<CompanyJobs />} />
        <Route path="applications" element={<CompanyApplications />} />
        <Route path="interviews" element={<CompanyInterviews />} />
        <Route path="reviews" element={<CompanyReviews />} />
        <Route path="export" element={<CompanyExport />} />
        <Route path="profile" element={<CompanyProfile />} />
      </Route>

      <Route
        path="/jobseeker"
        element={
          <RequireAuth roles={['jobseeker', 'admin']}>
            <JobSeekerLayout />
          </RequireAuth>
        }
      >
        <Route path="" element={<Navigate to="dashboard" replace />} />
        <Route path="dashboard" element={<JobSeekerDashboard />} />
        <Route path="jobs" element={<JobSeekerJobs />} />
        <Route path="resumes" element={<JobSeekerResumes />} />
        <Route path="applications" element={<JobSeekerApplications />} />
        <Route path="interviews" element={<JobSeekerInterviews />} />
        <Route path="reviews" element={<JobSeekerReviews />} />
        <Route path="profile" element={<JobSeekerProfile />} />
      </Route>

      <Route
        path="/admin"
        element={
          <RequireAuth roles={['admin']}>
            <AdminLayout />
          </RequireAuth>
        }
      >
        <Route path="" element={<Navigate to="dashboard" replace />} />
        <Route path="dashboard" element={<AdminDashboard />} />
        <Route path="users" element={<AdminUsers />} />
      </Route>

      <Route path="*" element={<div className="p-8 text-center">404 - Page Not Found</div>} />
    </Routes>
  )
}

export default App
