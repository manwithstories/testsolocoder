import { createBrowserRouter, Navigate } from 'react-router-dom'
import { useSelector } from 'react-redux'
import type { RootState } from '@/store'
import { Suspense, lazy, ReactNode } from 'react'

const Login = lazy(() => import('@/pages/Login'))
const Register = lazy(() => import('@/pages/Register'))
const Dashboard = lazy(() => import('@/pages/Dashboard'))
const PlanList = lazy(() => import('@/pages/PlanList'))
const PlanDetail = lazy(() => import('@/pages/PlanDetail'))
const Activities = lazy(() => import('@/pages/Activities'))
const Budget = lazy(() => import('@/pages/Budget'))
const Files = lazy(() => import('@/pages/Files'))
const Checklist = lazy(() => import('@/pages/Checklist'))
const MapView = lazy(() => import('@/pages/MapView'))
const Reminders = lazy(() => import('@/pages/Reminders'))
const Layout = lazy(() => import('@/components/Layout'))

function PrivateRoute({ children }: { children: ReactNode }) {
  const isAuthenticated = useSelector((state: RootState) => state.auth.isAuthenticated)
  return isAuthenticated ? children : <Navigate to="/login" replace />
}

export const router = createBrowserRouter([
  {
    path: '/login',
    element: <Login />,
  },
  {
    path: '/register',
    element: <Register />,
  },
  {
    path: '/',
    element: (
      <PrivateRoute>
        <Layout />
      </PrivateRoute>
    ),
    children: [
      {
        index: true,
        element: <Navigate to="/dashboard" replace />,
      },
      {
        path: 'dashboard',
        element: <Dashboard />,
      },
      {
        path: 'plans',
        element: <PlanList />,
      },
      {
        path: 'plans/:id',
        element: <PlanDetail />,
      },
      {
        path: 'plans/:id/activities',
        element: <Activities />,
      },
      {
        path: 'plans/:id/budget',
        element: <Budget />,
      },
      {
        path: 'plans/:id/files',
        element: <Files />,
      },
      {
        path: 'plans/:id/checklist',
        element: <Checklist />,
      },
      {
        path: 'plans/:id/map',
        element: <MapView />,
      },
      {
        path: 'reminders',
        element: <Reminders />,
      },
    ],
  },
])
