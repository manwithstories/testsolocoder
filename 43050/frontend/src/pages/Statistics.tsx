import React, { useState, useEffect } from 'react';
import { statsAPI, groupAPI } from '../api/services';
import { StatisticsSummary, MonthlyStats, Expense, Group, GroupMember } from '../types';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer, PieChart, Pie, Cell } from 'recharts';
import dayjs from 'dayjs';

const Statistics: React.FC = () => {
  const [summary, setSummary] = useState<StatisticsSummary | null>(null);
  const [monthlyStats, setMonthlyStats] = useState<MonthlyStats[]>([]);
  const [history, setHistory] = useState<Expense[]>([]);
  const [groups, setGroups] = useState<Group[]>([]);
  const [members, setMembers] = useState<GroupMember[]>([]);
  const [loading, setLoading] = useState(true);
  const [filters, setFilters] = useState({
    startDate: '',
    endDate: '',
    groupId: '',
    memberId: '',
  });

  const fetchData = async () => {
    setLoading(true);
    try {
      const params: any = {};
      if (filters.startDate) params.startDate = filters.startDate;
      if (filters.endDate) params.endDate = filters.endDate;
      if (filters.groupId) params.groupId = filters.groupId;
      if (filters.memberId) params.memberId = filters.memberId;

      const [summaryRes, monthlyRes, historyRes, groupsRes] = await Promise.all([
        statsAPI.getSummary(params),
        statsAPI.getMonthlyStats(),
        statsAPI.getHistory(params),
        groupAPI.getUserGroups(),
      ]);
      setSummary(summaryRes.data);
      setMonthlyStats(monthlyRes.data);
      setHistory(historyRes.data);
      setGroups(groupsRes.data);
    } catch (err) {
      console.error('Failed to fetch statistics:', err);
    } finally {
      setLoading(false);
    }
  };

  const fetchGroupMembers = async (groupId: string) => {
    if (!groupId) {
      setMembers([]);
      setFilters((prev) => ({ ...prev, memberId: '' }));
      return;
    }
    try {
      const res = await groupAPI.getGroupMembers(parseInt(groupId, 10));
      setMembers(res.data.filter((m) => m.isActive));
    } catch (err) {
      console.error('Failed to fetch group members:', err);
    }
  };

  useEffect(() => {
    fetchData();
  }, [filters]);

  useEffect(() => {
    if (filters.groupId) {
      fetchGroupMembers(filters.groupId);
    } else {
      setMembers([]);
    }
  }, [filters.groupId]);

  const chartData = monthlyStats.map((m) => ({
    month: m.month,
    支出: m.totalPaid,
    待还: m.totalOwed,
  })).reverse();

  const categoryData = React.useMemo(() => {
    const map = new Map<string, number>();
    history.forEach((expense) => {
      const key = expense.title;
      map.set(key, (map.get(key) || 0) + expense.amount);
    });
    return Array.from(map.entries())
      .map(([name, value]) => ({ name, value }))
      .sort((a, b) => b.value - a.value)
      .slice(0, 8);
  }, [history]);

  const COLORS = ['#3b82f6', '#10b981', '#f59e0b', '#ef4444', '#8b5cf6', '#ec4899', '#06b6d4', '#84cc16'];

  if (loading) {
    return <div className="text-center py-12">加载中...</div>;
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-800">统计分析</h1>
        <p className="text-gray-500 mt-1">查看您的收支明细和统计数据</p>
      </div>

      <div className="bg-white rounded-xl shadow-sm border border-gray-100 p-6">
        <h3 className="text-sm font-medium text-gray-700 mb-4">筛选条件</h3>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          <div>
            <label className="block text-sm text-gray-600 mb-1">开始日期</label>
            <input
              type="date"
              value={filters.startDate}
              onChange={(e) => setFilters({ ...filters, startDate: e.target.value })}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
            />
          </div>
          <div>
            <label className="block text-sm text-gray-600 mb-1">结束日期</label>
            <input
              type="date"
              value={filters.endDate}
              onChange={(e) => setFilters({ ...filters, endDate: e.target.value })}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
            />
          </div>
          <div>
            <label className="block text-sm text-gray-600 mb-1">群组</label>
            <select
              value={filters.groupId}
              onChange={(e) => setFilters({ ...filters, groupId: e.target.value })}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
            >
              <option value="">全部群组</option>
              {groups.map((g) => (
                <option key={g.id} value={g.id.toString()}>
                  {g.name}
                </option>
              ))}
            </select>
          </div>
          <div>
            <label className="block text-sm text-gray-600 mb-1">成员</label>
            <select
              value={filters.memberId}
              onChange={(e) => setFilters({ ...filters, memberId: e.target.value })}
              disabled={!filters.groupId}
              className={`w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none ${
                !filters.groupId ? 'bg-gray-100 cursor-not-allowed' : ''
              }`}
            >
              <option value="">全部成员</option>
              {members.map((m) => (
                <option key={m.user.id} value={m.user.id.toString()}>
                  {m.user.username}
                </option>
              ))}
            </select>
            {!filters.groupId && (
              <p className="text-xs text-gray-400 mt-1">请先选择群组</p>
            )}
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div className="bg-white rounded-xl shadow-sm border border-gray-100 p-6">
          <p className="text-sm text-gray-500">总支出</p>
          <p className="text-2xl font-bold text-gray-800 mt-2">
            ¥{summary?.totalPaid.toFixed(2) || '0.00'}
          </p>
        </div>
        <div className="bg-white rounded-xl shadow-sm border border-gray-100 p-6">
          <p className="text-sm text-gray-500">待还金额</p>
          <p className="text-2xl font-bold text-red-500 mt-2">
            ¥{summary?.totalOwed.toFixed(2) || '0.00'}
          </p>
        </div>
        <div className="bg-white rounded-xl shadow-sm border border-gray-100 p-6">
          <p className="text-sm text-gray-500">净余额</p>
          <p className={`text-2xl font-bold mt-2 ${
            (summary?.netBalance || 0) >= 0 ? 'text-green-600' : 'text-red-500'
          }`}>
            {(summary?.netBalance || 0) >= 0 ? '+' : ''}¥{summary?.netBalance.toFixed(2) || '0.00'}
          </p>
        </div>
        <div className="bg-white rounded-xl shadow-sm border border-gray-100 p-6">
          <p className="text-sm text-gray-500">账单笔数</p>
          <p className="text-2xl font-bold text-gray-800 mt-2">
            {summary?.expenseCount || 0}
          </p>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-white rounded-xl shadow-sm border border-gray-100 p-6">
          <h3 className="text-lg font-semibold text-gray-800 mb-4">月度收支趋势</h3>
          {chartData.length > 0 ? (
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={chartData}>
                <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" />
                <XAxis dataKey="month" tick={{ fontSize: 12 }} />
                <YAxis tick={{ fontSize: 12 }} />
                <Tooltip />
                <Legend />
                <Bar dataKey="支出" fill="#3b82f6" radius={[4, 4, 0, 0]} />
                <Bar dataKey="待还" fill="#ef4444" radius={[4, 4, 0, 0]} />
              </BarChart>
            </ResponsiveContainer>
          ) : (
            <div className="h-[300px] flex items-center justify-center text-gray-500">
              暂无数据
            </div>
          )}
        </div>

        <div className="bg-white rounded-xl shadow-sm border border-gray-100 p-6">
          <h3 className="text-lg font-semibold text-gray-800 mb-4">消费分类占比</h3>
          {categoryData.length > 0 ? (
            <ResponsiveContainer width="100%" height={300}>
              <PieChart>
                <Pie
                  data={categoryData}
                  cx="50%"
                  cy="50%"
                  labelLine={false}
                  label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
                  outerRadius={100}
                  fill="#8884d8"
                  dataKey="value"
                >
                  {categoryData.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                  ))}
                </Pie>
                <Tooltip />
              </PieChart>
            </ResponsiveContainer>
          ) : (
            <div className="h-[300px] flex items-center justify-center text-gray-500">
              暂无数据
            </div>
          )}
        </div>
      </div>

      <div className="bg-white rounded-xl shadow-sm border border-gray-100">
        <div className="p-6 border-b border-gray-100">
          <h3 className="text-lg font-semibold text-gray-800">收支明细</h3>
        </div>
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="bg-gray-50">
                <th className="text-left py-3 px-6 text-sm font-medium text-gray-600">日期</th>
                <th className="text-left py-3 px-6 text-sm font-medium text-gray-600">账单</th>
                <th className="text-left py-3 px-6 text-sm font-medium text-gray-600">金额</th>
                <th className="text-left py-3 px-6 text-sm font-medium text-gray-600">付款人</th>
                <th className="text-left py-3 px-6 text-sm font-medium text-gray-600">分摊方式</th>
              </tr>
            </thead>
            <tbody>
              {history.length === 0 ? (
                <tr>
                  <td colSpan={5} className="py-12 text-center text-gray-500">
                    暂无账单记录
                  </td>
                </tr>
              ) : (
                history.map((expense) => (
                  <tr key={expense.id} className="border-b border-gray-100 hover:bg-gray-50">
                    <td className="py-4 px-6 text-sm text-gray-600">
                      {dayjs(expense.expenseDate).format('YYYY-MM-DD')}
                    </td>
                    <td className="py-4 px-6 text-sm font-medium text-gray-800">
                      {expense.title}
                    </td>
                    <td className="py-4 px-6 text-sm text-gray-800">
                      ¥{expense.amount.toFixed(2)}
                    </td>
                    <td className="py-4 px-6 text-sm text-gray-600">
                      {expense.payer?.username}
                    </td>
                    <td className="py-4 px-6">
                      <span className={`px-2 py-1 text-xs rounded-full ${
                        expense.splitType === 'equal' ? 'bg-blue-100 text-blue-700' :
                        expense.splitType === 'ratio' ? 'bg-purple-100 text-purple-700' :
                        'bg-orange-100 text-orange-700'
                      }`}>
                        {expense.splitType === 'equal' ? '均摊' :
                         expense.splitType === 'ratio' ? '按比例' : '自定义'}
                      </span>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};

export default Statistics;
