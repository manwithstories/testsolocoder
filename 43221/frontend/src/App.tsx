import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { ConfigProvider } from 'antd'
import zhCN from 'antd/locale/zh_CN'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { useAuth } from '@/hooks/useAuth'
import { AuthProvider } from '@/contexts/AuthContext'
import { MainLayout } from '@/components/layout/MainLayout'
import { Login } from '@/pages/auth/Login'
import { Register } from '@/pages/auth/Register'
import { Dashboard } from '@/pages/dashboard/Dashboard'
import { ServiceList } from '@/pages/services/ServiceList'
import { ServiceDetail } from '@/pages/services/ServiceDetail'
import { MyServices } from '@/pages/services/MyServices'
import { CreateService } from '@/pages/services/CreateService'
import { MyAppointments } from '@/pages/appointments/MyAppointments'
import { AppointmentDetail } from '@/pages/appointments/AppointmentDetail'
import { MyReviews } from '@/pages/reviews/MyReviews'
import { ReviewManagement } from '@/pages/reviews/ReviewManagement'
import { Statistics } from '@/pages/statistics/Statistics'
import { Notifications } from '@/pages/notifications/Notifications'
import { UserManagement } from '@/pages/admin/UserManagement'
import { VerificationManagement } from '@/pages/admin/VerificationManagement'
import { TemplateManagement } from '@/pages/admin/TemplateManagement'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: 1,
    },
  },
})

function PrivateRoute({ children, requiredRole }: { children: React.ReactNode; requiredRole?: string[] }) {
  const { user, loading, isAuthenticated } = useAuth()

  if (loading) {
    return <div>加载中...</div>
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />
  }

  if (requiredRole && user && !requiredRole.includes(user.role)) {
    return <Navigate to="/" replace />
  }

  return <>{children}</>
}

function AppRoutes() {
  const auth = useAuth()

  return (
    <AuthProvider value={auth}>
      <ConfigProvider locale={zhCN}>
        <QueryClientProvider client={queryClient}>
          <BrowserRouter>
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

                <Route path="services">
                  <Route index element={<ServiceList />} />
                  <Route path=":id" element={<ServiceDetail />} />
                  <Route
                    path="mine"
                    element={
                      <PrivateRoute requiredRole={['professional']}>
                        <MyServices />
                      </PrivateRoute>
                    }
                  />
                  <Route
                    path="create"
                    element={
                      <PrivateRoute requiredRole={['professional']}>
                        <CreateService />
                      </PrivateRoute>
                    }
                  />
                </Route>

                <Route path="appointments">
                  <Route index element={<MyAppointments />} />
                  <Route path=":id" element={<AppointmentDetail />} />
                </Route>

                <Route path="reviews">
                  <Route index element={<MyReviews />} />
                  <Route
                    path="management"
                    element={
                      <PrivateRoute requiredRole={['admin']}>
                        <ReviewManagement />
                      </PrivateRoute>
                    }
                  />
                </Route>

                <Route
                  path="statistics"
                  element={
                    <PrivateRoute requiredRole={['professional', 'admin']}>
                      <Statistics />
                    </PrivateRoute>
                  }
                />

                <Route path="notifications" element={<Notifications />} />

                <Route
                  path="admin"
                  element={
                    <PrivateRoute requiredRole={['admin']}>
                      <UserManagement />
                    </PrivateRoute>
                  }
                />
                <Route
                  path="admin/verifications"
                  element={
                    <PrivateRoute requiredRole={['admin']}>
                      <VerificationManagement />
                    </PrivateRoute>
                  }
                />
                <Route
                  path="admin/templates"
                  element={
                    <PrivateRoute requiredRole={['admin']}>
                      <TemplateManagement />
                    </PrivateRoute>
                  }
                />
              </Route>

              <Route path="*" element={<Navigate to="/" replace />} />
            </Routes>
          </BrowserRouter>
        </QueryClientProvider>
      </ConfigProvider>
    </AuthProvider>
  )
}

export default AppRoutes
