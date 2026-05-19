import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { useAuthStore } from '@/store/authStore'
import MainLayout from '@/components/MainLayout'
import Login from '@/pages/Login'
import ForgotPassword from '@/pages/ForgotPassword'
import Home from '@/pages/Home'
import VenueList from '@/pages/venues/VenueList'
import VenueDetail from '@/pages/venues/VenueDetail'
import DeviceList from '@/pages/devices/DeviceList'
import CalendarPage from '@/pages/booking/CalendarPage'
import OrderList from '@/pages/orders/OrderList'
import PaymentList from '@/pages/payments/PaymentList'
import ReviewList from '@/pages/reviews/ReviewList'
import StatsPage from '@/pages/stats/StatsPage'
import UserList from '@/pages/users/UserList'
import Profile from '@/pages/Profile'

const PrivateRoute = ({ children }: { children: React.ReactNode }) => {
  const { isAuthenticated } = useAuthStore()
  return isAuthenticated ? <>{children}</> : <Navigate to="/login" />
}

const AdminRoute = ({ children }: { children: React.ReactNode }) => {
  const { user } = useAuthStore()
  const isAdmin = user?.role === 'admin' || user?.role === 'super_admin'
  return isAdmin ? <>{children}</> : <Navigate to="/" />
}

const SuperAdminRoute = ({ children }: { children: React.ReactNode }) => {
  const { user } = useAuthStore()
  return user?.role === 'super_admin' ? <>{children}</> : <Navigate to="/" />
}

const AppRoutes = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/forgot-password" element={<ForgotPassword />} />
        <Route
          path="/"
          element={
            <PrivateRoute>
              <MainLayout />
            </PrivateRoute>
          }
        >
          <Route index element={<Home />} />
          <Route path="venues" element={<VenueList />} />
          <Route path="venues/:id" element={<VenueDetail />} />
          <Route path="devices" element={<DeviceList />} />
          <Route path="calendar" element={<CalendarPage />} />
          <Route path="orders" element={<OrderList />} />
          <Route
            path="payments"
            element={
              <AdminRoute>
                <PaymentList />
              </AdminRoute>
            }
          />
          <Route path="reviews" element={<ReviewList />} />
          <Route
            path="stats"
            element={
              <AdminRoute>
                <StatsPage />
              </AdminRoute>
            }
          />
          <Route
            path="users"
            element={
              <SuperAdminRoute>
                <UserList />
              </SuperAdminRoute>
            }
          />
          <Route path="profile" element={<Profile />} />
        </Route>
        <Route path="*" element={<Navigate to="/" />} />
      </Routes>
    </BrowserRouter>
  )
}

export default AppRoutes
