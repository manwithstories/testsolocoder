import React, { useState, useEffect } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { eventAPI } from '../services/api';
import { Registration } from '../types';
import { useAuth } from '../context/AuthContext';

const EventDetailPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { user } = useAuth();
  const [eventData, setEventData] = useState<any>(null);
  const [registrations, setRegistrations] = useState<Registration[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    if (id) {
      fetchEvent();
    }
  }, [id]);

  const fetchEvent = async () => {
    try {
      const res = await eventAPI.getEvent(Number(id));
      setEventData(res.data);
      if (res.data.is_organizer) {
        const regRes = await eventAPI.getRegistrations(Number(id));
        setRegistrations(regRes.data);
      }
    } catch (err: any) {
      setError(err.response?.data?.error || '获取活动详情失败');
    } finally {
      setLoading(false);
    }
  };

  const handleRegister = async () => {
    try {
      await eventAPI.registerEvent(Number(id));
      alert('报名成功！');
      fetchEvent();
    } catch (err: any) {
      alert(err.response?.data?.error || '报名失败');
    }
  };

  const handleCancel = async () => {
    if (!window.confirm('确定要取消报名吗？')) return;
    try {
      await eventAPI.cancelRegistration(Number(id));
      alert('取消报名成功');
      fetchEvent();
    } catch (err: any) {
      alert(err.response?.data?.error || '取消失败');
    }
  };

  const handleDelete = async () => {
    if (!window.confirm('确定要删除这个活动吗？此操作不可撤销。')) return;
    try {
      await eventAPI.deleteEvent(Number(id));
      alert('活动已删除');
      navigate('/');
    } catch (err: any) {
      alert(err.response?.data?.error || '删除失败');
    }
  };

  const handleExport = () => {
    eventAPI.exportRegistrations(Number(id));
  };

  const formatDate = (dateStr: string) => {
    return new Date(dateStr).toLocaleString('zh-CN');
  };

  if (loading) return <div className="text-center py-8">加载中...</div>;
  if (error) return <div className="text-center py-8 text-red-500">{error}</div>;
  if (!eventData) return <div className="text-center py-8">活动不存在</div>;

  const { event, registered_count, is_registered, is_full, is_deadline_passed, can_register, is_organizer, cancel_count } = eventData;

  return (
    <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="mb-6">
        <Link to="/" className="text-blue-600 hover:text-blue-800">← 返回活动列表</Link>
      </div>

      <div className="bg-white rounded-lg shadow-md overflow-hidden">
        <div className="p-8">
          <div className="flex justify-between items-start mb-6">
            <div>
              <h1 className="text-3xl font-bold text-gray-900 mb-2">{event.title}</h1>
              <p className="text-gray-500">组织者: {event.organizer.username}</p>
            </div>
            <div className="flex gap-2">
              {is_deadline_passed ? (
                <span className="px-3 py-1 text-sm rounded bg-gray-200 text-gray-600">已截止</span>
              ) : is_full ? (
                <span className="px-3 py-1 text-sm rounded bg-orange-200 text-orange-700">已满员</span>
              ) : is_registered ? (
                <span className="px-3 py-1 text-sm rounded bg-green-200 text-green-700">已报名</span>
              ) : (
                <span className="px-3 py-1 text-sm rounded bg-blue-200 text-blue-700">报名中</span>
              )}
            </div>
          </div>

          <div className="bg-gray-50 rounded-lg p-6 mb-6">
            <h3 className="font-semibold text-gray-900 mb-4">活动信息</h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
              <div>
                <span className="text-gray-500">📍 地点:</span>
                <span className="ml-2 text-gray-900">{event.location}</span>
              </div>
              <div>
                <span className="text-gray-500">👥 人数:</span>
                <span className="ml-2 text-gray-900">{registered_count}/{event.capacity} 人</span>
              </div>
              <div>
                <span className="text-gray-500">🕐 开始时间:</span>
                <span className="ml-2 text-gray-900">{formatDate(event.start_time)}</span>
              </div>
              <div>
                <span className="text-gray-500">⏰ 报名截止:</span>
                <span className="ml-2 text-gray-900">{formatDate(event.deadline)}</span>
              </div>
              {event.end_time && (
                <div>
                  <span className="text-gray-500">📅 结束时间:</span>
                  <span className="ml-2 text-gray-900">{formatDate(event.end_time)}</span>
                </div>
              )}
            </div>
          </div>

          {event.description && (
            <div className="mb-6">
              <h3 className="font-semibold text-gray-900 mb-2">活动描述</h3>
              <p className="text-gray-700 whitespace-pre-wrap">{event.description}</p>
            </div>
          )}

          {cancel_count > 0 && (
            <p className="text-orange-600 text-sm mb-4">您已取消报名 {cancel_count} 次（最多可取消3次）</p>
          )}

          <div className="flex gap-4">
            {user && !is_organizer && (
              is_registered ? (
                <button
                  onClick={handleCancel}
                  className="px-6 py-2 bg-red-500 hover:bg-red-600 text-white rounded-md"
                >
                  取消报名
                </button>
              ) : (
                <button
                  onClick={handleRegister}
                  disabled={!can_register}
                  className="px-6 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-md disabled:bg-gray-400 disabled:cursor-not-allowed"
                >
                  {cancel_count >= 3 ? '取消次数超限，无法报名' : is_deadline_passed ? '报名已截止' : is_full ? '活动已满' : '立即报名'}
                </button>
              )
            )}
            {is_organizer && (
              <>
                <button
                  onClick={handleExport}
                  className="px-6 py-2 bg-green-600 hover:bg-green-700 text-white rounded-md"
                >
                  导出报名名单
                </button>
                <button
                  onClick={handleDelete}
                  className="px-6 py-2 bg-red-500 hover:bg-red-600 text-white rounded-md"
                >
                  删除活动
                </button>
              </>
            )}
          </div>
        </div>
      </div>

      {is_organizer && registrations.length > 0 && (
        <div className="mt-8 bg-white rounded-lg shadow-md overflow-hidden">
          <div className="p-6">
            <h2 className="text-xl font-bold text-gray-900 mb-4">报名名单 ({registrations.length}人)</h2>
            <div className="overflow-x-auto">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">ID</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">用户名</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">邮箱</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text text-gray-500 uppercase">报名时间</th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {registrations.map((reg) => (
                    <tr key={reg.id}>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{reg.id}</td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{reg.user.username}</td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{reg.user.email}</td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{formatDate(reg.created_at)}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default EventDetailPage;
