import React, { useMemo } from 'react';
import {
  PieChart,
  Pie,
  Cell,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  LineChart,
  Line,
} from 'recharts';
import type { Card, Tag } from '../../types';
import { format, eachDayOfInterval, startOfWeek, endOfWeek, isSameDay } from 'date-fns';

interface StatsPanelProps {
  cards: Card[];
  tags: Tag[];
}

const StatsPanel: React.FC<StatsPanelProps> = ({ cards, tags }) => {
  const tagDistribution = useMemo(() => {
    return tags.map((tag) => ({
      name: tag.name,
      value: cards.filter((c) => c.tags.includes(tag.id)).length,
      color: tag.color,
    })).filter((t) => t.value > 0);
  }, [cards, tags]);

  const weeklyActivity = useMemo(() => {
    const now = new Date();
    const weekStart = startOfWeek(now);
    const weekEnd = endOfWeek(now);
    const days = eachDayOfInterval({ start: weekStart, end: weekEnd });

    return days.map((day) => {
      const createdCount = cards.filter((c) =>
        isSameDay(new Date(c.createdAt), day)
      ).length;
      const updatedCount = cards.filter((c) =>
        isSameDay(new Date(c.updatedAt), day) && !isSameDay(new Date(c.createdAt), day)
      ).length;

      return {
        name: format(day, 'EEE'),
        创建: createdCount,
        更新: updatedCount,
      };
    });
  }, [cards]);

  const reviewStatus = useMemo(() => {
    const overdue = cards.filter((c) => Date.now() > c.review.nextReviewAt).length;
    const dueSoon = cards.filter((c) => {
      const nextWeek = Date.now() + 7 * 24 * 60 * 60 * 1000;
      return Date.now() <= c.review.nextReviewAt && c.review.nextReviewAt <= nextWeek;
    }).length;
    const upcoming = cards.length - overdue - dueSoon;

    return [
      { name: '已逾期', value: overdue, color: '#FF6B6B' },
      { name: '一周内', value: dueSoon, color: '#FFD93D' },
      { name: '已掌握', value: upcoming, color: '#6BCB77' },
    ];
  }, [cards]);

  const totalCards = cards.length;
  const totalTags = tags.length;
  const totalLinks = cards.reduce((sum, c) => sum + c.linkedCardIds.length, 0) / 2;
  const avgLinks = totalCards > 0 ? (totalLinks / totalCards).toFixed(1) : '0';

  const last7DaysCards = cards.filter(
    (c) => Date.now() - c.createdAt < 7 * 24 * 60 * 60 * 1000
  ).length;

  const avgContentLength = totalCards > 0
    ? Math.round(cards.reduce((sum, c) => sum + c.content.length, 0) / totalCards)
    : 0;

  return (
    <div className="stats-panel">
      <div className="list-header">
        <h3>统计分析</h3>
      </div>

      <div className="stats-overview">
        <div className="stat-card">
          <span className="stat-number">{totalCards}</span>
          <span className="stat-label">知识卡片</span>
        </div>
        <div className="stat-card">
          <span className="stat-number">{totalTags}</span>
          <span className="stat-label">标签数量</span>
        </div>
        <div className="stat-card">
          <span className="stat-number">{totalLinks}</span>
          <span className="stat-label">关联数量</span>
        </div>
        <div className="stat-card">
          <span className="stat-number">{avgLinks}</span>
          <span className="stat-label">平均关联</span>
        </div>
        <div className="stat-card">
          <span className="stat-number">{last7DaysCards}</span>
          <span className="stat-label">近7天新增</span>
        </div>
        <div className="stat-card">
          <span className="stat-number">{avgContentLength}</span>
          <span className="stat-label">平均字数</span>
        </div>
      </div>

      <div className="charts-grid">
        <div className="chart-card">
          <h4>标签分布</h4>
          {tagDistribution.length === 0 ? (
            <p className="text-muted center">暂无数据</p>
          ) : (
            <ResponsiveContainer width="100%" height={300}>
              <PieChart>
                <Pie
                  data={tagDistribution}
                  cx="50%"
                  cy="50%"
                  labelLine={false}
                  label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
                  outerRadius={100}
                  dataKey="value"
                >
                  {tagDistribution.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={entry.color} />
                  ))}
                </Pie>
                <Tooltip />
              </PieChart>
            </ResponsiveContainer>
          )}
        </div>

        <div className="chart-card">
          <h4>复习状态</h4>
          <ResponsiveContainer width="100%" height={300}>
            <PieChart>
              <Pie
                data={reviewStatus}
                cx="50%"
                cy="50%"
                labelLine={false}
                label={({ name, value }) => `${name}: ${value}`}
                outerRadius={100}
                dataKey="value"
              >
                {reviewStatus.map((entry, index) => (
                  <Cell key={`cell-${index}`} fill={entry.color} />
                ))}
              </Pie>
              <Tooltip />
            </PieChart>
          </ResponsiveContainer>
        </div>

        <div className="chart-card full-width">
          <h4>本周活跃度</h4>
          <ResponsiveContainer width="100%" height={300}>
            <BarChart data={weeklyActivity}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="name" />
              <YAxis allowDecimals={false} />
              <Tooltip />
              <Bar dataKey="创建" fill="#4A90D9" />
              <Bar dataKey="更新" fill="#96CEB4" />
            </BarChart>
          </ResponsiveContainer>
        </div>

        <div className="chart-card full-width">
          <h4>卡片创建趋势</h4>
          <ResponsiveContainer width="100%" height={300}>
            <LineChart
              data={(() => {
                if (cards.length === 0) return [];
                const sorted = [...cards].sort((a, b) => a.createdAt - b.createdAt);
                const data: { date: string; 累计: number }[] = [];
                let count = 0;
                sorted.forEach((card) => {
                  count++;
                  data.push({
                    date: format(new Date(card.createdAt), 'MM-dd'),
                    累计: count,
                  });
                });
                return data;
              })()}
            >
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="date" />
              <YAxis allowDecimals={false} />
              <Tooltip />
              <Line type="monotone" dataKey="累计" stroke="#4A90D9" strokeWidth={2} dot={false} />
            </LineChart>
          </ResponsiveContainer>
        </div>
      </div>
    </div>
  );
};

export default StatsPanel;
