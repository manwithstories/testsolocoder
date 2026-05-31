import dayjs from 'dayjs';

export const formatDate = (date: string | Date, format: string = 'YYYY-MM-DD HH:mm:ss'): string => {
  return dayjs(date).format(format);
};

export const formatDateShort = (date: string | Date): string => {
  return dayjs(date).format('YYYY-MM-DD');
};

export const formatDateTimeRange = (start: string | Date, end: string | Date): string => {
  return `${formatDateShort(start)} 至 ${formatDateShort(end)}`;
};

export const getRelativeTime = (date: string | Date): string => {
  const diff = dayjs().diff(dayjs(date), 'minute');
  if (diff < 1) return '刚刚';
  if (diff < 60) return `${diff} 分钟前`;
  if (diff < 24 * 60) return `${Math.floor(diff / 60)} 小时前`;
  if (diff < 30 * 24 * 60) return `${Math.floor(diff / (24 * 60))} 天前`;
  return formatDateShort(date);
};

export const getTimeRange = (days: number) => {
  const end = dayjs();
  const start = end.subtract(days, 'day');
  return {
    startDate: start.format('YYYY-MM-DD'),
    endDate: end.format('YYYY-MM-DD'),
  };
};

export const isExpired = (date: string | Date): boolean => {
  return dayjs(date).isBefore(dayjs());
};

export const getDaysRemaining = (date: string | Date): number => {
  return dayjs(date).diff(dayjs(), 'day');
};
