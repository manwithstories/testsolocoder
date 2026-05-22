import { createBrowserRouter, Navigate } from 'react-router-dom'
import App from '@/App'
import Login from '@/pages/Login'
import Register from '@/pages/Register'
import Home from '@/pages/Home'
import ServiceList from '@/pages/ServiceList'
import ServiceDetail from '@/pages/ServiceDetail'
import OrderList from '@/pages/OrderList'
import OrderDetail from '@/pages/OrderDetail'
import MyInvitations from '@/pages/MyInvitations'
import ReviewList from '@/pages/ReviewList'
import ComplaintList from '@/pages/ComplaintList'
import BillList from '@/pages/BillList'
import WithdrawList from '@/pages/WithdrawList'
import MessageCenter from '@/pages/MessageCenter'
import Profile from '@/pages/Profile'
import AddressManage from '@/pages/AddressManage'
import Certification from '@/pages/Certification'
import MyServices from '@/pages/MyServices'
import Dashboard from '@/pages/Dashboard'
import UserManage from '@/pages/UserManage'
import ServiceCategoryManage from '@/pages/ServiceCategoryManage'
import ComplaintHandle from '@/pages/ComplaintHandle'

const ProtectedRoute = ({ children }: { children: React.ReactNode }) => {
  const token = localStorage.getItem('token')
  if (!token) {
    return <Navigate to="/login" replace />
  }
  return <>{children}</>
}

const RoleRoute = ({ children, roles }: { children: React.ReactNode; roles: string[] }) => {
  const userInfo = JSON.parse(localStorage.getItem('userInfo') || 'null')
  if (!userInfo || !roles.includes(userInfo.role)) {
    return <Navigate to="/" replace />
  }
  return <>{children}</>
}

export const router = createBrowserRouter([
  {
    path: '/',
    element: <App />,
    children: [
      {
        index: true,
        element: <Home />,
      },
      {
        path: 'services',
        element: <ServiceList />,
      },
      {
        path: 'services/:id',
        element: <ServiceDetail />,
      },
      {
        path: 'orders',
        element: (
          <ProtectedRoute>
            <OrderList />
          </ProtectedRoute>
        ),
      },
      {
        path: 'orders/:id',
        element: (
          <ProtectedRoute>
            <OrderDetail />
          </ProtectedRoute>
        ),
      },
      {
        path: 'invitations',
        element: (
          <ProtectedRoute>
            <RoleRoute roles={['service_provider']}>
              <MyInvitations />
            </RoleRoute>
          </ProtectedRoute>
        ),
      },
      {
        path: 'reviews',
        element: (
          <ProtectedRoute>
            <ReviewList />
          </ProtectedRoute>
        ),
      },
      {
        path: 'complaints',
        element: (
          <ProtectedRoute>
            <ComplaintList />
          </ProtectedRoute>
        ),
      },
      {
        path: 'bills',
        element: (
          <ProtectedRoute>
            <RoleRoute roles={['service_provider']}>
              <BillList />
            </RoleRoute>
          </ProtectedRoute>
        ),
      },
      {
        path: 'withdraws',
        element: (
          <ProtectedRoute>
            <RoleRoute roles={['service_provider']}>
              <WithdrawList />
            </RoleRoute>
          </ProtectedRoute>
        ),
      },
      {
        path: 'messages',
        element: (
          <ProtectedRoute>
            <MessageCenter />
          </ProtectedRoute>
        ),
      },
      {
        path: 'profile',
        element: (
          <ProtectedRoute>
            <Profile />
          </ProtectedRoute>
        ),
      },
      {
        path: 'addresses',
        element: (
          <ProtectedRoute>
            <AddressManage />
          </ProtectedRoute>
        ),
      },
      {
        path: 'certification',
        element: (
          <ProtectedRoute>
            <RoleRoute roles={['service_provider']}>
              <Certification />
            </RoleRoute>
          </ProtectedRoute>
        ),
      },
      {
        path: 'my-services',
        element: (
          <ProtectedRoute>
            <RoleRoute roles={['service_provider']}>
              <MyServices />
            </RoleRoute>
          </ProtectedRoute>
        ),
      },
      {
        path: 'admin/dashboard',
        element: (
          <ProtectedRoute>
            <RoleRoute roles={['admin']}>
              <Dashboard />
            </RoleRoute>
          </ProtectedRoute>
        ),
      },
      {
        path: 'admin/users',
        element: (
          <ProtectedRoute>
            <RoleRoute roles={['admin']}>
              <UserManage />
            </RoleRoute>
          </ProtectedRoute>
        ),
      },
      {
        path: 'admin/categories',
        element: (
          <ProtectedRoute>
            <RoleRoute roles={['admin']}>
              <ServiceCategoryManage />
            </RoleRoute>
          </ProtectedRoute>
        ),
      },
      {
        path: 'admin/withdraws',
        element: (
          <ProtectedRoute>
            <RoleRoute roles={['admin']}>
              <WithdrawList />
            </RoleRoute>
          </ProtectedRoute>
        ),
      },
      {
        path: 'admin/complaints',
        element: (
          <ProtectedRoute>
            <RoleRoute roles={['admin']}>
              <ComplaintHandle />
            </RoleRoute>
          </ProtectedRoute>
        ),
      },
    ],
  },
  {
    path: '/login',
    element: <Login />,
  },
  {
    path: '/register',
    element: <Register />,
  },
  {
    path: '*',
    element: <Navigate to="/" replace />,
  },
])
