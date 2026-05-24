import React from 'react';
import { Outlet, NavLink, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

const navItems = [
  { to: '/', label: '数据看板', icon: '📊' },
  { to: '/templates', label: '商品模板', icon: '📐' },
  { to: '/orders', label: '订单管理', icon: '📋' },
  { to: '/orders/new', label: '在线定制', icon: '✏️' },
  { to: '/pricing', label: '智能报价', icon: '💰' },
  { to: '/production', label: '生产排程', icon: '🏭' },
  { to: '/customers', label: '客户管理', icon: '👥' },
  { to: '/invoices', label: '财务对账', icon: '📄' },
];

export default function Layout() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <div style={{ display: 'flex', minHeight: '100vh' }}>
      <aside style={{
        width: 220,
        background: '#1f2937',
        color: '#fff',
        padding: '16px 0',
        flexShrink: 0,
      }}>
        <div style={{ padding: '0 20px 20px', fontSize: 18, fontWeight: 600, borderBottom: '1px solid #374151' }}>
          印刷管理平台
        </div>
        <nav style={{ marginTop: 16 }}>
          {navItems.map((item) => (
            <NavLink
              key={item.to}
              to={item.to}
              end={item.to === '/'}
              style={({ isActive }) => ({
                display: 'flex',
                alignItems: 'center',
                padding: '12px 20px',
                color: isActive ? '#fff' : '#9ca3af',
                background: isActive ? '#3b82f6' : 'transparent',
                borderLeft: isActive ? '3px solid #60a5fa' : '3px solid transparent',
                textDecoration: 'none',
                fontSize: 14,
              })}
            >
              <span style={{ marginRight: 10 }}>{item.icon}</span>
              {item.label}
            </NavLink>
          ))}
        </nav>
      </aside>

      <div style={{ flex: 1, display: 'flex', flexDirection: 'column' }}>
        <header style={{
          background: '#fff',
          padding: '14px 24px',
          borderBottom: '1px solid #e5e7eb',
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
        }}>
          <div style={{ fontSize: 16, fontWeight: 500 }}>
            {navItems.find(n => window.location.pathname.endsWith(n.to.slice(1)) || (n.to === '/' && window.location.pathname === '/'))?.label || ''}
          </div>
          <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
            <span style={{ fontSize: 14, color: '#6b7280' }}>{user?.real_name || user?.username}</span>
            <button onClick={handleLogout} style={btnStyle}>退出</button>
          </div>
        </header>
        <main style={{ flex: 1, padding: 24, overflow: 'auto' }}>
          <Outlet />
        </main>
      </div>
    </div>
  );
}

const btnStyle: React.CSSProperties = {
  padding: '6px 14px',
  border: '1px solid #d1d5db',
  borderRadius: 6,
  background: '#fff',
  cursor: 'pointer',
  fontSize: 13,
};
