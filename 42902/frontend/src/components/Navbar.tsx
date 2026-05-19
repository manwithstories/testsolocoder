import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

const Navbar: React.FC = () => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <nav className="bg-blue-600 text-white shadow-lg">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between h-16 items-center">
          <Link to="/" className="text-xl font-bold">活动管理系统</Link>
          <div className="flex items-center space-x-4">
            {user ? (
              <>
                <Link to="/" className="hover:text-blue-200">活动列表</Link>
                <Link to="/create" className="hover:text-blue-200">发布活动</Link>
                <Link to="/my-registrations" className="hover:text-blue-200">我的报名</Link>
                <span className="text-blue-200">欢迎, {user.username}</span>
                <button
                  onClick={handleLogout}
                  className="bg-red-500 hover:bg-red-600 px-4 py-2 rounded"
                >
                  退出
                </button>
              </>
            ) : (
              <>
                <Link to="/login" className="hover:text-blue-200">登录</Link>
                <Link to="/register" className="bg-white text-blue-600 px-4 py-2 rounded hover:bg-blue-50">
                  注册
                </Link>
              </>
            )}
          </div>
        </div>
      </div>
    </nav>
  );
};

export default Navbar;
