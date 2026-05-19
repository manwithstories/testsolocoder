import {
  format,
  parseISO,
  differenceInDays,
  startOfWeek,
  endOfWeek,
  startOfMonth,
  endOfMonth,
  isSameDay,
  isToday,
  addDays,
} from 'date-fns';
import { zhCN } from 'date-fns/locale';

export const formatDate = (date: string | Date, pattern: string = 'yyyy-MM-dd'): string => {
  const d = typeof date === 'string' ? parseISO(date) : date;
  return format(d, pattern, { locale: zhCN });
};

export const formatDateTime = (date: string | Date): string => {
  return formatDate(date, 'yyyy-MM-dd HH:mm');
};

export const formatRelativeTime = (date: string | Date): string => {
  const d = typeof date === 'string' ? parseISO(date) : date;
  const now = new Date();
  const diffDays = differenceInDays(now, d);

  if (diffDays === 0) return '今天';
  if (diffDays === 1) return '昨天';
  if (diffDays < 7) return `${diffDays}天前`;
  if (diffDays < 30) return `${Math.floor(diffDays / 7)}周前`;
  return formatDate(d);
};

export const formatDuration = (minutes: number): string => {
  if (minutes < 60) return `${minutes}分钟`;
  const hours = Math.floor(minutes / 60);
  const mins = minutes % 60;
  return mins > 0 ? `${hours}小时${mins}分钟` : `${hours}小时`;
};

export const getWeekRange = (date: Date = new Date()): { start: Date; end: Date } => {
  return {
    start: startOfWeek(date, { weekStartsOn: 1 }),
    end: endOfWeek(date, { weekStartsOn: 1 }),
  };
};

export const getMonthRange = (date: Date = new Date()): { start: Date; end: Date } => {
  return {
    start: startOfMonth(date),
    end: endOfMonth(date),
  };
};

export const isDateInRange = (date: Date, start: Date, end: Date): boolean => {
  return date >= start && date <= end;
};

export const generateDateRange = (start: Date, end: Date): Date[] => {
  const dates: Date[] = [];
  let current = new Date(start);
  while (current <= end) {
    dates.push(new Date(current));
    current = addDays(current, 1);
  }
  return dates;
};

export const calculateStreak = (dates: string[]): { current: number; longest: number } => {
  if (dates.length === 0) return { current: 0, longest: 0 };

  const uniqueDates = [...new Set(dates.map(d => formatDate(d, 'yyyy-MM-dd')))]
    .map(d => parseISO(d))
    .sort((a, b) => b.getTime() - a.getTime());

  let currentStreak = 0;
  let longestStreak = 0;
  let tempStreak = 0;

  for (let i = 0; i < uniqueDates.length; i++) {
    const date = uniqueDates[i];
    const expectedDate = addDays(new Date(), -i);

    if (isSameDay(date, expectedDate)) {
      currentStreak++;
    } else if (i === 0 && !isToday(date)) {
      break;
    } else {
      break;
    }
  }

  for (let i = 0; i < uniqueDates.length; i++) {
    if (i === 0) {
      tempStreak = 1;
    } else {
      const diff = differenceInDays(uniqueDates[i - 1], uniqueDates[i]);
      if (diff === 1) {
        tempStreak++;
      } else {
        longestStreak = Math.max(longestStreak, tempStreak);
        tempStreak = 1;
      }
    }
  }
  longestStreak = Math.max(longestStreak, tempStreak);

  return { current: currentStreak, longest: longestStreak };
};

export const getDateInputValue = (date: Date = new Date()): string => {
  return format(date, 'yyyy-MM-dd');
};
