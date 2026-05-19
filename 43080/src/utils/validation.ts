import { z } from 'zod';
import { logger } from './logger';

export const bookSchema = z.object({
  id: z.string(),
  title: z.string().min(1, '书名不能为空'),
  author: z.string().min(1, '作者不能为空'),
  isbn: z.string().optional(),
  coverUrl: z.string().url().optional().or(z.literal('')),
  totalPages: z.number().int().positive('总页数必须大于0'),
  categories: z.array(z.string()),
  status: z.enum(['unread', 'reading', 'completed', 'paused']),
  isFavorite: z.boolean(),
  inReadingQueue: z.boolean(),
  currentPage: z.number().int().nonnegative(),
  createdAt: z.string(),
  updatedAt: z.string(),
});

export const bookFormSchema = z.object({
  title: z.string().min(1, '书名不能为空'),
  author: z.string().min(1, '作者不能为空'),
  isbn: z.string().optional(),
  coverUrl: z.string().url('请输入有效的URL').optional().or(z.literal('')),
  totalPages: z.number().int().positive('总页数必须大于0'),
  categories: z.array(z.string()),
});

export const readingSessionSchema = z.object({
  id: z.string(),
  bookId: z.string(),
  startPage: z.number().int().nonnegative('起始页不能为负数'),
  endPage: z.number().int().positive('结束页必须大于0'),
  duration: z.number().int().positive('阅读时长必须大于0'),
  date: z.string(),
  notes: z.string().optional(),
});

export const readingPlanSchema = z.object({
  id: z.string(),
  bookId: z.string(),
  dailyPages: z.number().int().positive('每日页数必须大于0'),
  startDate: z.string(),
  targetDate: z.string(),
});

export const noteSchema = z.object({
  id: z.string(),
  bookId: z.string(),
  page: z.number().int().nonnegative('页码不能为负数'),
  chapter: z.string().optional(),
  content: z.string().min(1, '内容不能为空'),
  type: z.enum(['note', 'highlight', 'bookmark']),
  highlightColor: z.string().optional(),
  createdAt: z.string(),
  updatedAt: z.string(),
});

export const settingsSchema = z.object({
  notificationEnabled: z.boolean(),
  notificationTime: z.string().regex(/^([01]?[0-9]|2[0-3]):[0-5][0-9]$/, '时间格式必须为 HH:mm'),
  theme: z.enum(['light', 'dark', 'system']),
  defaultView: z.enum(['grid', 'list']),
});

export const validateData = <T>(schema: z.ZodSchema<T>, data: unknown): T | null => {
  try {
    return schema.parse(data);
  } catch (error) {
    if (error instanceof z.ZodError) {
      logger.error('Validation failed', error.errors);
    }
    return null;
  }
};

export const validateReadingSession = (
  startPage: number,
  endPage: number,
  totalPages: number
): { valid: boolean; error?: string } => {
  if (startPage < 0) {
    return { valid: false, error: '起始页不能为负数' };
  }
  if (endPage <= startPage) {
    return { valid: false, error: '结束页必须大于起始页' };
  }
  if (endPage > totalPages) {
    return { valid: false, error: `结束页不能超过总页数 ${totalPages}` };
  }
  return { valid: true };
};

export const validateReadingPlan = (
  dailyPages: number,
  startDate: string,
  targetDate: string
): { valid: boolean; error?: string } => {
  if (dailyPages <= 0) {
    return { valid: false, error: '每日阅读页数必须大于0' };
  }
  const start = new Date(startDate);
  const target = new Date(targetDate);
  const today = new Date();
  today.setHours(0, 0, 0, 0);

  if (start < today) {
    return { valid: false, error: '开始日期不能早于今天' };
  }
  if (target <= start) {
    return { valid: false, error: '目标日期必须晚于开始日期' };
  }
  return { valid: true };
};

export const validateISBN = (isbn: string): boolean => {
  const cleaned = isbn.replace(/[-\s]/g, '');
  if (cleaned.length === 10) {
    let sum = 0;
    for (let i = 0; i < 9; i++) {
      sum += parseInt(cleaned[i]) * (10 - i);
    }
    const check = cleaned[9] === 'X' ? 10 : parseInt(cleaned[9]);
    return (sum + check) % 11 === 0;
  }
  if (cleaned.length === 13) {
    let sum = 0;
    for (let i = 0; i < 12; i++) {
      sum += parseInt(cleaned[i]) * (i % 2 === 0 ? 1 : 3);
    }
    const check = parseInt(cleaned[12]);
    return (10 - (sum % 10)) % 10 === check;
  }
  return false;
};
