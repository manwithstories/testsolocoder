export const DEFAULT_CATEGORIES = [
  '文学',
  '科幻',
  '历史',
  '哲学',
  '心理学',
  '商业',
  '技术',
  '艺术',
  '传记',
  '教育',
  '健康',
  '旅行',
  '悬疑',
  '浪漫',
  '科普',
  '经济',
  '政治',
  '宗教',
  '社会科学',
  '自然科学',
] as const;

export const BOOK_STATUS = [
  { value: 'unread', label: '未读' },
  { value: 'reading', label: '正在阅读' },
  { value: 'completed', label: '已完成' },
  { value: 'paused', label: '已暂停' },
] as const;

export const STATUS_LABELS: Record<string, string> = {
  unread: '未读',
  reading: '正在阅读',
  completed: '已完成',
  paused: '已暂停',
};

export const STATUS_COLORS: Record<string, string> = {
  unread: 'gray',
  reading: 'blue',
  completed: 'green',
  paused: 'warning',
};
