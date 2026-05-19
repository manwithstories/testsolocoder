import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { eventAPI } from '../services/api';
import { EventWithStatus } from '../types';
import { useAuth } from '../context/AuthContext';

const EventsPage: React.FC = () => {
  const [events, setEvents] = useState<EventWithStatus[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const { user } = useAuth();

  useEffect(() => {
    fetchEvents();
  }, []);

  const fetchEvents = async () => {
    try {
      const res = await eventAPI.getEvents();
      setEvents(res.data);
    } catch (err: any) {
      setError(err.response?.data?.error || '获取活动列表失败');
    } finally {
      setLoading(false);
    }
  };

  const handleRegister = async (eventId: number) => {
    try {
      await eventAPI.registerEvent(eventId);
      alert('报名成功！');
      fetchEvents();
    } catch (err: any) {
      alert(err.response?.data?.error || '报名失败');
    }
  };

  const handleCancel = async (eventId: number) => {
    if (!window.confirm('确定要取消报名吗？')) return;
    try {
      await eventAPI.cancelRegistration(eventId);
      alert('取消报名成功');
      fetchEvents();
    } catch (err: any) {
      alert(err.response?.data?.error || '取消失败');
    }
  };

  const formatDate = (dateStr: string) => {
    return new Date(dateStr).toLocaleString('zh-CN');
  };

  const getStatusBadge = (event: EventWithStatus) => {
    if (event.is_deadline_passed) {
      return <span className="px-2 py-1 text-xs rounded bg-gray-200 text-gray-600">已截止</span>;
    }
    if (event.is_full) {
      return <span className="px-2 py-1 text-xs rounded bg-orange-200 text-orange-700">已满员</span>;
    }
    if (event.is_registered) {
      return <span className="px-2 py-1 text-xs rounded bg-green-200 text-green-700">已报名</span>;
    }
    return <span className="px-2 py-1 text-xs rounded bg-blue-200 text-blue-700">报名中</span>;
  };

  if (loading) return <div className="text-center py-8">加载中...</div>;
  if (error) return <div className="text-center py-8 text-red-500">{error}</div>;

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-3xl font-bold text-gray-900">活动列表</h1>
        {user && (
          <Link
            to="/create"
            className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md"
          >
            + 发布活动
          </Link>
        )}
      </div>

      {events.length === 0 ? (
        <div className="text-center py-16 text-gray-500">
          <p className="text-xl">暂无活动</p>
          {user && <p className="mt-2">点击上方按钮发布第一个活动吧！</p>}
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {events.map((event) => (
            <div key={event.id} className="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-lg transition-shadow">
              <div className="p-6">
                <div className="flex justify-between items-start mb-3">
                  <h3 className="text-xl font-semibold text-gray-900 flex-1">{event.title}</h3>
                  {getStatusBadge(event)}
                </div>
                <p className="text-gray-600 mb-4 line-clamp-2">{event.description || '暂无描述'}</p>
                <div className="space-y-2 text-sm text-gray-500">
                  <p>📍 {event.location}</p>
                  <p>🕐 {formatDate(event.start_time)}</p>
                  <p>⏰ 报名截止: {formatDate(event.deadline)}</p>
                  <p>👥 {event.registered_count}/{event.capacity} 人</p>
                  <p>👤 组织者: {event.organizer.username}</p>
                </div>
                {event.cancel_count > 0 && (
                  <p className="text-xs text-orange-600 mt-2">已取消报名 {event.cancel_count} 次（最多3次）</p>
                )}
              </div>
              <div className="px-6 py-4 bg-gray-50 flex gap-2">
                <Link
                  to={`/events/${event.id}`}
                  className="flex-1 text-center px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-100"
                >
                  查看详情
                </Link>
                {user && !event.is_organizer && (
                  event.is_registered ? (
                    <button
                      onClick={() => handleCancel(event.id)}
                      className="flex-1 px-4 py-2 bg-red-500 hover:bg-red-600 text-white rounded-md"
                    >
                      取消报名
                    </button>
                  ) : (
                    <button
                      onClick={() => handleRegister(event.id)}
                      disabled={!event.can_register}
                      className="flex-1 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-md disabled:bg-gray-400 disabled:cursor-not-allowed"
                    >
                      {event.cancel_count >= 3 ? '取消次数超限' : '立即报名'}
                    </button>
                  )
                )}
                {user && event.is_organizer && (
                  <Link
                    to={`/events/${event.id}`}
                    className="flex-1 text-center px-4 py-2 bg-green-600 hover:bg-green-700 text-white rounded-md"
                  >
                    管理活动
                  </Link>
                )}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default EventsPage;
