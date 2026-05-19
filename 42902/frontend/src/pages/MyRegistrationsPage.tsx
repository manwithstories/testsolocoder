import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { authAPI } from '../services/api';
import { Registration } from '../types';
import { useAuth } from '../context/AuthContext';

const MyRegistrationsPage: React.FC = () => {
  const { user } = useAuth();
  const navigate = useNavigate();
  const [registrations, setRegistrations] = useState<Registration[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    if (!user) {
      navigate('/login');
      return;
    }
    fetchRegistrations();
  }, [user, navigate]);

  const fetchRegistrations = async () => {
    try {
      const res = await authAPI.getMyRegistrations();
      setRegistrations(res.data);
    } catch (err: any) {
      setError(err.response?.data?.error || '获取报名记录失败');
    } finally {
      setLoading(false);
    }
  };

  const formatDate = (dateStr: string) => {
    return new Date(dateStr).toLocaleString('zh-CN');
  };

  const getStatusBadge = (status: string) => {
    switch (status) {
      case 'registered':
        return <span className="px-2 py-1 text-xs rounded bg-green-200 text-green-700">已报名</span>;
      case 'cancelled':
        return <span className="px-2 py-1 text-xs rounded bg-gray-200 text-gray-600">已取消</span>;
      default:
        return <span className="px-2 py-1 text-xs rounded bg-gray-200 text-gray-600">{status}</span>;
    }
  };

  if (loading) return <div className="text-center py-8">加载中...</div>;
  if (error) return <div className="text-center py-8 text-red-500">{error}</div>;

  return (
    <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="mb-6">
        <Link to="/" className="text-blue-600 hover:text-blue-800">← 返回活动列表</Link>
      </div>

      <h1 className="text-2xl font-bold text-gray-900 mb-6">我的报名记录</h1>

      {registrations.length === 0 ? (
        <div className="text-center py-16 text-gray-500">
          <p className="text-xl">暂无报名记录</p>
          <Link to="/" className="mt-4 inline-block text-blue-600 hover:text-blue-800">
            去浏览活动 →
          </Link>
        </div>
      ) : (
        <div className="bg-white rounded-lg shadow-md overflow-hidden">
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">活动名称</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">地点</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">开始时间</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">状态</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">取消次数</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">报名时间</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">操作</th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {registrations.map((reg) => (
                  <tr key={reg.id}>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <Link to={`/events/${reg.event_id}`} className="text-blue-600 hover:text-blue-800">
                        {reg.event.title}
                      </Link>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{reg.event.location}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{formatDate(reg.event.start_time)}</td>
                    <td className="px-6 py-4 whitespace-nowrap">{getStatusBadge(reg.status)}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{reg.cancel_count}次</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{formatDate(reg.created_at)}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm">
                      <Link to={`/events/${reg.event_id}`} className="text-blue-600 hover:text-blue-800">
                        查看详情
                      </Link>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}
    </div>
  );
};

export default MyRegistrationsPage;
