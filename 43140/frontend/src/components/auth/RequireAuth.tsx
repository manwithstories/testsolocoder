import { ReactNode } from 'react'
import { Navigate, useLocation } from 'react-router-dom'
import { useAppSelector } from '@/hooks/redux'
import { selectAuth } from '@/store/slices/authSlice'

interface RequireAuthProps {
  children: ReactNode
  roles?: string[]
}

export default function RequireAuth({ children, roles }: RequireAuthProps) {
  const { user, token } = useAppSelector(selectAuth)
  const location = useLocation()

  if (!token) {
    return <Navigate to="/login" state={{ from: location }} replace />
  }

  if (roles && user && !roles.includes(user.role)) {
    return <Navigate to="/" replace />
  }

  return <>{children}</>
}
