import React, { useEffect } from 'react'
import { Routes, Route, Navigate } from 'react-router-dom'
import { Layout } from 'antd'
import MainLayout from '@/components/Layout/MainLayout'
import ProtectedRoute from '@/components/ProtectedRoute'
import LoginPage from '@/pages/Login'
import RegisterPage from '@/pages/Register'
import HomePage from '@/pages/Home'
import CourseDetailPage from '@/pages/CourseDetail'
import CourseListPage from '@/pages/CourseList'
import MyCoursesPage from '@/pages/MyCourses'
import OrdersPage from '@/pages/Orders'
import ProfilePage from '@/pages/Profile'
import InstructorDashboard from '@/pages/instructor/Dashboard'
import InstructorCourses from '@/pages/instructor/Courses'
import InstructorCourseEdit from '@/pages/instructor/CourseEdit'
import InstructorAnalytics from '@/pages/instructor/Analytics'
import AdminDashboard from '@/pages/admin/Dashboard'
import AdminUsers from '@/pages/admin/Users'
import AdminOrders from '@/pages/admin/Orders'
import AdminCoupons from '@/pages/admin/Coupons'
import AdminApplications from '@/pages/admin/Applications'
import AdminAnalytics from '@/pages/admin/Analytics'
import { useAuthStore } from '@/store/auth'

const { Content } = Layout

const App: React.FC = () => {
  const fetchProfile = useAuthStore((state) => state.fetchProfile)
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated)

  useEffect(() => {
    if (isAuthenticated) {
      fetchProfile()
    }
  }, [])

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Content>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />

          <Route element={<MainLayout />}>
            <Route path="/" element={<HomePage />} />
            <Route path="/courses" element={<CourseListPage />} />
            <Route path="/courses/:id" element={<CourseDetailPage />} />

            <Route element={<ProtectedRoute roles={['student', 'instructor', 'admin']} />}>
              <Route path="/my-courses" element={<MyCoursesPage />} />
              <Route path="/orders" element={<OrdersPage />} />
              <Route path="/profile" element={<ProfilePage />} />
            </Route>

            <Route element={<ProtectedRoute roles={['instructor', 'admin']} />}>
              <Route path="/instructor/dashboard" element={<InstructorDashboard />} />
              <Route path="/instructor/courses" element={<InstructorCourses />} />
              <Route path="/instructor/courses/new" element={<InstructorCourseEdit />} />
              <Route path="/instructor/courses/:id/edit" element={<InstructorCourseEdit />} />
              <Route path="/instructor/analytics" element={<InstructorAnalytics />} />
            </Route>

            <Route element={<ProtectedRoute roles={['admin']} />}>
              <Route path="/admin/dashboard" element={<AdminDashboard />} />
              <Route path="/admin/users" element={<AdminUsers />} />
              <Route path="/admin/orders" element={<AdminOrders />} />
              <Route path="/admin/coupons" element={<AdminCoupons />} />
              <Route path="/admin/applications" element={<AdminApplications />} />
              <Route path="/admin/analytics" element={<AdminAnalytics />} />
            </Route>

            <Route path="*" element={<Navigate to="/" replace />} />
          </Route>
        </Routes>
      </Content>
    </Layout>
  )
}

export default App
