import { Routes, Route, Navigate } from 'react-router-dom'
import { useEffect } from 'react'
import { useAuthStore } from '@/store/auth'
import MainLayout from '@/layouts/MainLayout'
import LoginPage from '@/pages/LoginPage'
import RegisterPage from '@/pages/RegisterPage'
import HomePage from '@/pages/HomePage'
import PlotsPage from '@/pages/PlotsPage'
import PlantsPage from '@/pages/PlantsPage'
import CalendarPage from '@/pages/CalendarPage'
import GrowthPage from '@/pages/GrowthPage'
import DiseasePage from '@/pages/DiseasePage'
import CommunityPage from '@/pages/CommunityPage'
import PostDetailPage from '@/pages/PostDetailPage'
import ExchangePage from '@/pages/ExchangePage'
import ShopPage from '@/pages/ShopPage'
import CartPage from '@/pages/CartPage'
import OrdersPage from '@/pages/OrdersPage'
import ProfilePage from '@/pages/ProfilePage'

function App() {
  const initialize = useAuthStore((state) => state.initialize)

  useEffect(() => {
    initialize()
  }, [initialize])

  const isAuthenticated = useAuthStore((state) => state.isAuthenticated)

  const ProtectedRoute = ({ children }: { children: React.ReactNode }) => {
    if (!isAuthenticated) {
      return <Navigate to="/login" replace />
    }
    return <>{children}</>
  }

  return (
    <Routes>
      <Route path="/login" element={!isAuthenticated ? <LoginPage /> : <Navigate to="/" />} />
      <Route path="/register" element={!isAuthenticated ? <RegisterPage /> : <Navigate to="/" />} />

      <Route
        path="/"
        element={
          <ProtectedRoute>
            <MainLayout />
          </ProtectedRoute>
        }
      >
        <Route index element={<HomePage />} />
        <Route path="plots" element={<PlotsPage />} />
        <Route path="plants" element={<PlantsPage />} />
        <Route path="calendar" element={<CalendarPage />} />
        <Route path="growth" element={<GrowthPage />} />
        <Route path="disease" element={<DiseasePage />} />
        <Route path="community" element={<CommunityPage />} />
        <Route path="community/:id" element={<PostDetailPage />} />
        <Route path="exchange" element={<ExchangePage />} />
        <Route path="shop" element={<ShopPage />} />
        <Route path="cart" element={<CartPage />} />
        <Route path="orders" element={<OrdersPage />} />
        <Route path="profile" element={<ProfilePage />} />
      </Route>

      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  )
}

export default App
