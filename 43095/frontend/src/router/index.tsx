import { createBrowserRouter, Navigate, type RouteObject } from 'react-router-dom'
import { useAuthStore } from '@/store'
import MainLayout from '@/layouts/MainLayout'
import AuthLayout from '@/layouts/AuthLayout'
import Login from '@/pages/Login'
import Register from '@/pages/Register'
import Dashboard from '@/pages/Dashboard'
import DoctorList from '@/pages/Doctors'
import DoctorDetail from '@/pages/Doctors/Detail'
import AppointmentList from '@/pages/Appointments'
import CreateAppointment from '@/pages/Appointments/Create'
import AppointmentDetail from '@/pages/Appointments/Detail'
import HealthRecords from '@/pages/HealthRecords'
import HealthRecordsHistory from '@/pages/HealthRecords/History'
import Notifications from '@/pages/Notifications'
import Payments from '@/pages/Payments'
import AdminDepartments from '@/pages/Admin/Departments'
import AdminDoctors from '@/pages/Admin/Doctors'
import type { UserRole } from '@/types'

const ProtectedRoute = ({ children, roles }: { children: React.ReactNode; roles?: UserRole[] }) => {
  const { isAuthenticated, user } = useAuthStore()

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />
  }

  if (roles && user && !roles.includes(user.role)) {
    return <Navigate to="/403" replace />
  }

  return <>{children}</>
}

const RoleBasedRedirect = () => {
  const { user } = useAuthStore()
  if (!user) return <Navigate to="/login" replace />

  switch (user.role) {
    case 'admin':
      return <Navigate to="/admin" replace />
    case 'doctor':
      return <Navigate to="/doctor" replace />
    case 'patient':
    default:
      return <Navigate to="/dashboard" replace />
  }
}

const adminRoutes: RouteObject[] = [
  {
    path: 'departments',
    element: <AdminDepartments />
  },
  {
    path: 'doctors',
    element: <AdminDoctors />
  },
  {
    path: 'patients',
    element: <div>患者管理</div>
  }
]

const doctorRoutes: RouteObject[] = [
  {
    path: 'schedule',
    element: <div>我的排班</div>
  },
  {
    path: 'appointments',
    element: <div>我的预约</div>
  },
  {
    path: 'patients',
    element: <div>我的患者</div>
  }
]

const routes: RouteObject[] = [
  {
    path: '/',
    element: <ProtectedRoute><MainLayout /></ProtectedRoute>,
    children: [
      {
        index: true,
        element: <RoleBasedRedirect />
      },
      {
        path: 'dashboard',
        element: <Dashboard />
      },
      {
        path: 'doctors',
        element: <DoctorList />
      },
      {
        path: 'doctors/:id',
        element: <DoctorDetail />
      },
      {
        path: 'appointments',
        element: <AppointmentList />
      },
      {
        path: 'appointments/create',
        element: <CreateAppointment />
      },
      {
        path: 'appointments/:id',
        element: <AppointmentDetail />
      },
      {
        path: 'health-records',
        element: <HealthRecords />
      },
      {
        path: 'health-records/history',
        element: <HealthRecordsHistory />
      },
      {
        path: 'notifications',
        element: <Notifications />
      },
      {
        path: 'payments',
        element: <Payments />
      },
      {
        path: 'admin',
        element: <ProtectedRoute roles={['admin']}><Dashboard /></ProtectedRoute>,
        children: adminRoutes
      },
      {
        path: 'doctor',
        element: <ProtectedRoute roles={['doctor']}><Dashboard /></ProtectedRoute>,
        children: doctorRoutes
      }
    ]
  },
  {
    path: '/login',
    element: <AuthLayout><Login /></AuthLayout>
  },
  {
    path: '/register',
    element: <AuthLayout><Register /></AuthLayout>
  },
  {
    path: '/403',
    element: <div>403 无权限访问</div>
  },
  {
    path: '*',
    element: <div>404 页面不存在</div>
  }
]

const router = createBrowserRouter(routes)

export default router
