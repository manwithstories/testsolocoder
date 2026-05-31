import { Routes, Route, Navigate } from 'react-router-dom'
import { useAuthStore } from '@/context/AuthContext'
import MainLayout from '@/components/Layout/MainLayout'
import Login from '@/pages/Login'
import Register from '@/pages/Register'
import Dashboard from '@/pages/Dashboard'
import Pets from '@/pages/Pets'
import PetDetail from '@/pages/PetDetail'
import Reservations from '@/pages/Reservations'
import ReservationDetail from '@/pages/ReservationDetail'
import NewReservation from '@/pages/NewReservation'
import DailyRecords from '@/pages/DailyRecords'
import Reviews from '@/pages/Reviews'
import Orders from '@/pages/Orders'
import Alerts from '@/pages/Alerts'
import Statistics from '@/pages/Statistics'
import Profile from '@/pages/Profile'
import StoreDashboard from '@/pages/StoreDashboard'
import KeeperDashboard from '@/pages/KeeperDashboard'

function PrivateRoute({ children }: { children: React.ReactNode }) {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated())
  return isAuthenticated ? <>{children}</> : <Navigate to="/login" replace />
}

function App() {
  return (
    <Routes>
      <Route path="/login" element={<Login />} />
      <Route path="/register" element={<Register />} />
      <Route
        path="/"
        element={
          <PrivateRoute>
            <MainLayout />
          </PrivateRoute>
        }
      >
        <Route index element={<Dashboard />} />
        <Route path="pets" element={<Pets />} />
        <Route path="pets/:id" element={<PetDetail />} />
        <Route path="reservations" element={<Reservations />} />
        <Route path="reservations/new" element={<NewReservation />} />
        <Route path="reservations/:id" element={<ReservationDetail />} />
        <Route path="daily-records" element={<DailyRecords />} />
        <Route path="reviews" element={<Reviews />} />
        <Route path="orders" element={<Orders />} />
        <Route path="alerts" element={<Alerts />} />
        <Route path="statistics" element={<Statistics />} />
        <Route path="profile" element={<Profile />} />
        <Route path="store-dashboard" element={<StoreDashboard />} />
        <Route path="keeper-dashboard" element={<KeeperDashboard />} />
      </Route>
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  )
}

export default App
