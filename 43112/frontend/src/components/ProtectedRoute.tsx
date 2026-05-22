import React from 'react'
import { Navigate, Outlet } from 'react-router-dom'
import { Result, Button } from 'antd'
import { useAuthStore } from '@/store/auth'

interface Props {
  roles?: string[]
}

const ProtectedRoute: React.FC<Props> = ({ roles }) => {
  const { isAuthenticated, user } = useAuthStore()

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />
  }

  if (roles && user && !roles.includes(user.role)) {
    return (
      <Result
        status="403"
        title="403"
        subTitle="抱歉，您没有访问此页面的权限"
        extra={<Button type="primary" href="/">返回首页</Button>}
      />
    )
  }

  return <Outlet />
}

export default ProtectedRoute
