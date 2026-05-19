import React, { useState } from 'react';
import {
  AreaChart,
  Area,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  BarChart,
  Bar,
} from 'recharts';
import type { WeeklyStats, MonthlyStats } from '../../types/reading';
import { Card, CardHeader, CardTitle, CardContent } from '../ui/Card';
import { Tabs } from '../ui/Tabs';
import { formatDuration } from '../../utils/date';

interface ReadingChartProps {
  weeklyStats: WeeklyStats[];
  monthlyStats: MonthlyStats[];
}

export const ReadingChart: React.FC<ReadingChartProps> = ({ weeklyStats, monthlyStats }) => {
  const [chartType, setChartType] = useState<'weekly' | 'monthly'>('weekly');

  const weeklyData = weeklyStats.map((stat) => ({
    name: stat.weekStart.slice(5),
    时长: Math.round(stat.totalDuration / 60),
    页数: stat.totalPages,
  }));

  const monthlyData = monthlyStats.map((stat) => ({
    name: stat.month.slice(5),
    时长: Math.round(stat.totalDuration / 60),
    页数: stat.totalPages,
  }));

  const data = chartType === 'weekly' ? weeklyData : monthlyData;

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle>阅读趋势</CardTitle>
          <div className="flex gap-2">
            <button
              onClick={() => setChartType('weekly')}
              className={`px-3 py-1 text-sm rounded-lg transition-colors ${
                chartType === 'weekly'
                  ? 'bg-blue-100 text-blue-700'
                  : 'text-gray-500 hover:bg-gray-100'
              }`}
            >
              周
            </button>
            <button
              onClick={() => setChartType('monthly')}
              className={`px-3 py-1 text-sm rounded-lg transition-colors ${
                chartType === 'monthly'
                  ? 'bg-blue-100 text-blue-700'
                  : 'text-gray-500 hover:bg-gray-100'
              }`}
            >
              月
            </button>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <div className="h-80">
          <ResponsiveContainer width="100%" height="100%">
            <AreaChart data={data}>
              <defs>
                <linearGradient id="colorDuration" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="#3b82f6" stopOpacity={0.3} />
                  <stop offset="95%" stopColor="#3b82f6" stopOpacity={0} />
                </linearGradient>
              </defs>
              <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" />
              <XAxis dataKey="name" stroke="#9ca3af" fontSize={12} />
              <YAxis stroke="#9ca3af" fontSize={12} />
              <Tooltip
                contentStyle={{
                  backgroundColor: 'white',
                  border: '1px solid #e5e7eb',
                  borderRadius: '8px',
                  boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)',
                }}
                formatter={(value: number, name: string) => [
                  name === '时长' ? `${value} 小时` : `${value} 页`,
                  name,
                ]}
              />
              <Area
                type="monotone"
                dataKey="时长"
                stroke="#3b82f6"
                strokeWidth={2}
                fill="url(#colorDuration)"
              />
            </AreaChart>
          </ResponsiveContainer>
        </div>
      </CardContent>
    </Card>
  );
};
