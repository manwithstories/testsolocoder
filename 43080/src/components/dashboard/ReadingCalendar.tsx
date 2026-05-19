import React, { useMemo } from 'react';
import { Card, CardHeader, CardTitle, CardContent } from '../ui/Card';
import type { DailyReadingStats } from '../../types/reading';
import { generateDateRange, formatDate } from '../../utils/date';
import { subDays } from 'date-fns';

interface ReadingCalendarProps {
  dailyStats: DailyReadingStats[];
}

export const ReadingCalendar: React.FC<ReadingCalendarProps> = ({ dailyStats }) => {
  const days = useMemo(() => {
    const today = new Date();
    const startDate = subDays(today, 83);
    return generateDateRange(startDate, today);
  }, []);

  const statsMap = useMemo(() => {
    const map = new Map<string, number>();
    dailyStats.forEach((stat) => {
      map.set(stat.date, stat.duration);
    });
    return map;
  }, [dailyStats]);

  const getIntensity = (duration: number): string => {
    if (duration === 0) return 'bg-gray-100';
    if (duration < 30) return 'bg-blue-200';
    if (duration < 60) return 'bg-blue-300';
    if (duration < 120) return 'bg-blue-400';
    return 'bg-blue-600';
  };

  const weekDays = ['一', '二', '三', '四', '五', '六', '日'];

  return (
    <Card>
      <CardHeader>
        <CardTitle>阅读日历</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="flex gap-1 mb-2">
          {weekDays.map((day, index) => (
            <div
              key={day}
              className="flex-1 text-center text-xs text-gray-400"
              style={{ visibility: index < 7 ? 'visible' : 'hidden' }}
            >
              {day}
            </div>
          ))}
        </div>
        <div className="grid grid-cols-7 gap-1">
          {days.map((day) => {
            const dateStr = formatDate(day);
            const duration = statsMap.get(dateStr) || 0;
            return (
              <div
                key={dateStr}
                className={`aspect-square rounded-sm ${getIntensity(duration)} transition-colors hover:ring-2 hover:ring-blue-500`}
                title={`${dateStr}: ${duration > 0 ? `${duration} 分钟` : '未阅读'}`}
              />
            );
          })}
        </div>
        <div className="flex items-center justify-end gap-2 mt-4 text-xs text-gray-500">
          <span>少</span>
          <div className="w-3 h-3 rounded-sm bg-gray-100" />
          <div className="w-3 h-3 rounded-sm bg-blue-200" />
          <div className="w-3 h-3 rounded-sm bg-blue-300" />
          <div className="w-3 h-3 rounded-sm bg-blue-400" />
          <div className="w-3 h-3 rounded-sm bg-blue-600" />
          <span>多</span>
        </div>
      </CardContent>
    </Card>
  );
};
