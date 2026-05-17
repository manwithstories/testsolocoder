import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { statsAPI, groupAPI } from '../api/services';
import { StatisticsSummary, Group } from '../types';
import dayjs from 'dayjs';

const Dashboard: React.FC = () => {
  const [summary, setSummary] = useState<StatisticsSummary | null>(null);
  const [groups, setGroups] = useState<Group[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [summaryRes, groupsRes] = await Promise.all([
          statsAPI.getSummary(),
          groupAPI.getUserGroups(),
        ]);
        setSummary(summaryRes.data);
        setGroups(groupsRes.data);
      } catch (err) {
        console.error('Failed to fetch dashboard data:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  if (loading) {
    return <div className="text-center py-12">加载中...</div>;
  }

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-2xl font-bold text-gray-800">首页概览</h1>
        <p className="text-gray-500 mt-1">查看您的账单统计概览</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="bg-white rounded-xl shadow-sm p-6 border border-gray-100">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-gray-500 text-sm">总支出</p>
              <p className="text-3xl font-bold text-gray-800 mt-2">
                ¥{summary?.totalPaid.toFixed(2) || '0.00'}
              </p>
            </div>
            <div className="w-14 h-14 bg-blue-100 rounded-full flex items-center justify-center text-2xl">
              💰
            </div>
          </div>
        </div>

        <div className="bg-white rounded-xl shadow-sm p-6 border border-gray-100">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-gray-500 text-sm">待还金额</p>
              <p className="text-3xl font-bold text-red-500 mt-2">
                ¥{summary?.totalOwed.toFixed(2) || '0.00'}
              </p>
            </div>
            <div className="w-14 h-14 bg-red-100 rounded-full flex items-center justify-center text-2xl">
              💸
            </div>
          </div>
        </div>

        <div className="bg-white rounded-xl shadow-sm p-6 border border-gray-100">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-gray-500 text-sm">净余额</p>
              <p className={`text-3xl font-bold mt-2 ${
                (summary?.netBalance || 0) >= 0 ? 'text-green-600' : 'text-red-500'
              }`}>
                {(summary?.netBalance || 0) >= 0 ? '+' : ''}¥{summary?.netBalance.toFixed(2) || '0.00'}
              </p>
            </div>
            <div className={`w-14 h-14 rounded-full flex items-center justify-center text-2xl ${
              (summary?.netBalance || 0) >= 0 ? 'bg-green-100' : 'bg-red-100'
            }`}>
              {(summary?.netBalance || 0) >= 0 ? '📈' : '📉'}
            </div>
          </div>
        </div>
      </div>

      <div className="bg-white rounded-xl shadow-sm border border-gray-100">
        <div className="p-6 border-b border-gray-100 flex items-center justify-between">
          <h2 className="text-lg font-semibold text-gray-800">我的群组</h2>
          <Link
            to="/groups"
            className="text-primary-600 hover:text-primary-700 text-sm font-medium"
          >
            查看全部 →
          </Link>
        </div>
        <div className="p-6">
          {groups.length === 0 ? (
            <div className="text-center py-8">
              <p className="text-gray-500">暂无群组</p>
              <Link
                to="/groups"
                className="inline-block mt-4 px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition-colors"
              >
                创建群组
              </Link>
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {groups.slice(0, 6).map((group) => (
                <Link
                  key={group.id}
                  to={`/groups/${group.id}`}
                  className="block p-4 border border-gray-200 rounded-lg hover:border-primary-300 hover:bg-primary-50 transition-colors"
                >
                  <div className="flex items-center gap-3">
                    <div className="w-10 h-10 bg-primary-100 rounded-full flex items-center justify-center text-primary-600 font-bold">
                      {group.name.charAt(0).toUpperCase()}
                    </div>
                    <div className="flex-1 min-w-0">
                      <h3 className="font-medium text-gray-800 truncate">{group.name}</h3>
                      <p className="text-sm text-gray-500 truncate">{group.description || '暂无描述'}</p>
                    </div>
                  </div>
                  <div className="mt-3 flex items-center justify-between text-sm text-gray-500">
                    <span>{group.members?.length || 0} 位成员</span>
                    <span>创建于 {dayjs(group.createdAt).format('YYYY-MM-DD')}</span>
                  </div>
                </Link>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
