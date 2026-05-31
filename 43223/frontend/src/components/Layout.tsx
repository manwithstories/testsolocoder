import { Link, Outlet, useNavigate } from 'react-router-dom'
import { useAuthStore } from '@/store/auth'
import { useCartStore } from '@/store/cart'
import { useEffect } from 'react'

export default function Layout() {
  const { isAuthenticated, user, logout, hasRole } = useAuthStore()
  const { totalCount, loadCart } = useCartStore()
  const navigate = useNavigate()

  useEffect(() => {
    if (isAuthenticated) {
      loadCart()
    }
  }, [isAuthenticated])

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  return (
    <div className="min-h-screen flex flex-col">
      <header className="bg-white shadow-sm sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-4 py-4 flex items-center justify-between">
          <Link to="/" className="text-2xl font-bold text-coffee-700 flex items-center gap-2">
            <span>☕</span>
            <span>Coffee Hub</span>
          </Link>

          <nav className="hidden md:flex items-center gap-6">
            <Link to="/" className="text-gray-700 hover:text-coffee-600">首页</Link>
            <Link to="/products" className="text-gray-700 hover:text-coffee-600">咖啡豆</Link>
            <Link to="/roasters" className="text-gray-700 hover:text-coffee-600">烘焙师</Link>
            <Link to="/search" className="text-gray-700 hover:text-coffee-600">搜索</Link>
          </nav>

          <div className="flex items-center gap-4">
            {isAuthenticated ? (
              <>
                <Link to="/cart" className="relative text-gray-700 hover:text-coffee-600">
                  <span className="text-xl">🛒</span>
                  {totalCount > 0 && (
                    <span className="absolute -top-2 -right-2 bg-coffee-600 text-white text-xs rounded-full w-5 h-5 flex items-center justify-center">
                      {totalCount}
                    </span>
                  )}
                </Link>
                <div className="relative group">
                  <button className="flex items-center gap-2 text-gray-700 hover:text-coffee-600">
                    <div className="w-8 h-8 bg-coffee-200 rounded-full flex items-center justify-center">
                      {user?.nickname?.[0] || user?.username[0] || 'U'}
                    </div>
                    <span className="hidden md:inline">{user?.nickname || user?.username}</span>
                  </button>
                  <div className="absolute right-0 mt-2 w-48 bg-white rounded-lg shadow-lg py-2 hidden group-hover:block">
                    <Link to="/profile" className="block px-4 py-2 text-gray-700 hover:bg-gray-100">个人中心</Link>
                    <Link to="/orders" className="block px-4 py-2 text-gray-700 hover:bg-gray-100">我的订单</Link>
                    <Link to="/cupping" className="block px-4 py-2 text-gray-700 hover:bg-gray-100">我的杯测</Link>
                    {hasRole(['admin', 'roaster']) && (
                      <>
                        <hr className="my-1" />
                        {hasRole(['roaster']) && (
                          <>
                            <Link to="/roaster/dashboard" className="block px-4 py-2 text-gray-700 hover:bg-gray-100">烘焙师工作台</Link>
                            <Link to="/certification" className="block px-4 py-2 text-gray-700 hover:bg-gray-100">认证申请</Link>
                          </>
                        )}
                        {hasRole(['admin']) && (
                          <Link to="/admin" className="block px-4 py-2 text-gray-700 hover:bg-gray-100">管理后台</Link>
                        )}
                      </>
                    )}
                    <hr className="my-1" />
                    <button onClick={handleLogout} className="block w-full text-left px-4 py-2 text-red-600 hover:bg-gray-100">
                      退出登录
                    </button>
                  </div>
                </div>
              </>
            ) : (
              <>
                <Link to="/login" className="text-gray-700 hover:text-coffee-600">登录</Link>
                <Link to="/register" className="btn btn-primary">注册</Link>
              </>
            )}
          </div>
        </div>
      </header>

      <main className="flex-1 max-w-7xl w-full mx-auto px-4 py-6">
        <Outlet />
      </main>

      <footer className="bg-gray-800 text-white py-8">
        <div className="max-w-7xl mx-auto px-4 text-center">
          <p className="text-lg font-bold mb-2">☕ Coffee Hub</p>
          <p className="text-gray-400 text-sm">精品咖啡豆交易与烘焙管理平台</p>
          <p className="text-gray-500 text-xs mt-4">© 2024 Coffee Hub. All rights reserved.</p>
        </div>
      </footer>
    </div>
  )
}
