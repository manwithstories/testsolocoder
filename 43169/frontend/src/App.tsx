import { Routes, Route, Navigate } from 'react-router-dom'
import { useAuthStore } from '@/store/authStore'
import MainLayout from '@/layouts/MainLayout'
import LoginPage from '@/pages/LoginPage'
import RegisterPage from '@/pages/RegisterPage'
import HomePage from '@/pages/HomePage'
import ProfilePage from '@/pages/ProfilePage'
import MatchPage from '@/pages/MatchPage'
import DatePage from '@/pages/DatePage'
import ChatPage from '@/pages/ChatPage'
import MemberPage from '@/pages/MemberPage'
import MatchmakerPage from '@/pages/MatchmakerPage'
import AdminDashboard from '@/pages/AdminDashboard'
import VerifyPage from '@/pages/VerifyPage'

function App() {
  const { user } = useAuthStore()

  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route path="/register" element={<RegisterPage />} />

      <Route path="/" element={user ? <MainLayout /> : <Navigate to="/login" />}>
        <Route index element={<HomePage />} />
        <Route path="profile" element={<ProfilePage />} />
        <Route path="verify" element={<VerifyPage />} />
        <Route path="match" element={<MatchPage />} />
        <Route path="dates" element={<DatePage />} />
        <Route path="chat" element={<ChatPage />} />
        <Route path="chat/:userId" element={<ChatPage />} />
        <Route path="member" element={<MemberPage />} />
        <Route path="matchmaker" element={<MatchmakerPage />} />
        <Route path="admin" element={<AdminDashboard />} />
      </Route>

      <Route path="*" element={<Navigate to="/" />} />
    </Routes>
  )
}

export default App
