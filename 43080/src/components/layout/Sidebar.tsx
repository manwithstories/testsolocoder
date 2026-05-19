import React from 'react';
import { NavLink, useLocation } from 'react-router-dom';
import {
  LayoutDashboard,
  BookOpen,
  PlusCircle,
  StickyNote,
  Settings,
  BookMarked,
  ListOrdered,
} from 'lucide-react';
import { cn } from '../../lib/utils';

const menuItems = [
  { path: '/', icon: LayoutDashboard, label: '仪表盘' },
  { path: '/shelf', icon: BookOpen, label: '我的书架' },
  { path: '/shelf?queue=true', icon: ListOrdered, label: '阅读队列' },
  { path: '/shelf?favorite=true', icon: BookMarked, label: '收藏夹' },
  { path: '/add-book', icon: PlusCircle, label: '添加书籍' },
  { path: '/notes', icon: StickyNote, label: '笔记管理' },
  { path: '/settings', icon: Settings, label: '设置' },
];

export const Sidebar: React.FC = () => {
  const location = useLocation();

  const isActive = (path: string) => {
    if (path === '/') return location.pathname === '/';
    return location.pathname === path.split('?')[0];
  };

  return (
    <aside className="w-64 bg-white border-r border-gray-200 h-screen fixed left-0 top-0 flex flex-col">
      <div className="p-6 border-b border-gray-100">
        <h1 className="text-2xl font-bold text-gray-900 flex items-center gap-2">
          <BookOpen className="w-8 h-8 text-blue-600" />
          阅读助手
        </h1>
        <p className="text-sm text-gray-500 mt-1">记录你的阅读之旅</p>
      </div>

      <nav className="flex-1 p-4 space-y-1 overflow-y-auto">
        {menuItems.map((item) => (
          <NavLink
            key={item.path}
            to={item.path}
            className={cn(
              'flex items-center gap-3 px-4 py-3 rounded-lg text-sm font-medium transition-colors',
              isActive(item.path)
                ? 'bg-blue-50 text-blue-700'
                : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'
            )}
          >
            <item.icon className="w-5 h-5" />
            {item.label}
          </NavLink>
        ))}
      </nav>

      <div className="p-4 border-t border-gray-100">
        <div className="bg-gradient-to-br from-blue-500 to-blue-600 rounded-xl p-4 text-white">
          <h3 className="font-semibold text-sm">今日阅读</h3>
          <p className="text-blue-100 text-xs mt-1">保持阅读习惯，每天进步一点</p>
        </div>
      </div>
    </aside>
  );
};
