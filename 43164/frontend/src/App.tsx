import { Routes, Route, Navigate } from 'react-router-dom'
import { useAuthStore } from '@/store/auth'
import MainLayout from '@/layouts/MainLayout'
import Login from '@/pages/Login'
import Register from '@/pages/Register'
import Home from '@/pages/Home'
import TeacherList from '@/pages/TeacherList'
import TeacherDetail from '@/pages/TeacherDetail'
import BookingPage from '@/pages/Booking'
import MyBookings from '@/pages/MyBookings'
import VideoSession from '@/pages/VideoSession'
import Wallet from '@/pages/Wallet'
import Messages from '@/pages/Messages'
import Profile from '@/pages/Profile'
import Reviews from '@/pages/Reviews'
import LearningProgress from '@/pages/LearningProgress'
import TeacherDashboard from '@/pages/TeacherDashboard'
import StudentDashboard from '@/pages/StudentDashboard'
import AdminDashboard from '@/pages/AdminDashboard'

function PrivateRoute({ children, role }: { children: React.ReactNode; role?: string }) {
  const { isAuthenticated, user } = useAuthStore()

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />
  }

  if (role && user?.role !== role && user?.role !== 'admin') {
    return <Navigate to="/" replace />
  }

  return <>{children}</>
}

function App() {
  return (
    <Routes>
      <Route path="/login" element={<Login />} />
      <Route path="/register" element={<Register />} />

      <Route path="/" element={<MainLayout />}>
        <Route index element={<Home />} />
        <Route path="teachers" element={<TeacherList />} />
        <Route path="teachers/:id" element={<TeacherDetail />} />

        <Route path="booking" element={
          <PrivateRoute role="student">
            <BookingPage />
          </PrivateRoute>
        } />

        <Route path="my-bookings" element={
          <PrivateRoute>
            <MyBookings />
          </PrivateRoute>
        } />

        <Route path="video/:sessionId" element={
          <PrivateRoute>
            <VideoSession />
          </PrivateRoute>
        } />

        <Route path="wallet" element={
          <PrivateRoute>
            <Wallet />
          </PrivateRoute>
        } />

        <Route path="messages" element={
          <PrivateRoute>
            <Messages />
          </PrivateRoute>
        } />

        <Route path="profile" element={
          <PrivateRoute>
            <Profile />
          </PrivateRoute>
        } />

        <Route path="reviews" element={
          <PrivateRoute>
            <Reviews />
          </PrivateRoute>
        } />

        <Route path="learning" element={
          <PrivateRoute role="student">
            <LearningProgress />
          </PrivateRoute>
        } />

        <Route path="teacher/dashboard" element={
          <PrivateRoute role="teacher">
            <TeacherDashboard />
          </PrivateRoute>
        } />

        <Route path="student/dashboard" element={
          <PrivateRoute role="student">
            <StudentDashboard />
          </PrivateRoute>
        } />

        <Route path="admin/dashboard" element={
          <PrivateRoute role="admin">
            <AdminDashboard />
          </PrivateRoute>
        } />
      </Route>
    </Routes>
  )
}

export default App
