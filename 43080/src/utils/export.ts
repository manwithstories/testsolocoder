import type { Book } from '../types/book';
import type { ReadingSession, ReadingPlan } from '../types/reading';
import type { Note, Bookmark } from '../types/note';
import type { ExportOptions } from '../types/settings';
import { logger } from './logger';
import { formatDate } from './date';

interface ExportData {
  books?: Book[];
  readingPlans?: ReadingPlan[];
  readingSessions?: ReadingSession[];
  notes?: Note[];
  bookmarks?: Bookmark[];
  exportedAt: string;
  version: string;
}

export const exportToJSON = (data: ExportData, fileName?: string): void => {
  try {
    const content = JSON.stringify(data, null, 2);
    if (fileName) {
      downloadFile(content, `${fileName}.json`, 'application/json');
    }
  } catch (error) {
    logger.error('Failed to export to JSON', error);
    throw new Error('导出JSON失败');
  }
};

export const exportToCSV = (data: ExportData, fileName?: string): void => {
  try {
    const parts: string[] = [];
    if (data.books && data.books.length > 0) {
      parts.push('=== 书籍 ===\n' + booksToCSV(data.books));
    }
    if (data.readingSessions && data.readingSessions.length > 0) {
      parts.push('\n=== 阅读记录 ===\n' + readingSessionsToCSV(data.readingSessions, data.books || []));
    }
    if (data.notes && data.notes.length > 0) {
      parts.push('\n=== 笔记 ===\n' + notesToCSV(data.notes, data.books || []));
    }
    if (fileName) {
      downloadFile(parts.join('\n\n'), `${fileName}.csv`, 'text/csv');
    }
  } catch (error) {
    logger.error('Failed to export to CSV', error);
    throw new Error('导出CSV失败');
  }
};

const escapeCSV = (value: unknown): string => {
  if (value === null || value === undefined) return '';
  const str = String(value);
  if (str.includes(',') || str.includes('"') || str.includes('\n')) {
    return `"${str.replace(/"/g, '""')}"`;
  }
  return str;
};

export const booksToCSV = (books: Book[]): string => {
  const headers = ['ID', '书名', '作者', 'ISBN', '总页数', '分类', '状态', '当前页', '收藏', '阅读队列', '创建时间', '更新时间'];
  const rows = books.map(book => [
    book.id,
    book.title,
    book.author,
    book.isbn ?? '',
    book.totalPages,
    book.categories.join('; '),
    book.status,
    book.currentPage,
    book.isFavorite ? '是' : '否',
    book.inReadingQueue ? '是' : '否',
    formatDate(book.createdAt, 'yyyy-MM-dd HH:mm:ss'),
    formatDate(book.updatedAt, 'yyyy-MM-dd HH:mm:ss'),
  ]);
  return [headers, ...rows].map(row => row.map(escapeCSV).join(',')).join('\n');
};

export const readingSessionsToCSV = (sessions: ReadingSession[], books: Book[]): string => {
  const bookMap = new Map(books.map(b => [b.id, b]));
  const headers = ['ID', '书名', '起始页', '结束页', '阅读页数', '阅读时长(分钟)', '日期', '笔记'];
  const rows = sessions.map(session => {
    const book = bookMap.get(session.bookId);
    return [
      session.id,
      book?.title ?? session.bookId,
      session.startPage,
      session.endPage,
      session.endPage - session.startPage,
      session.duration,
      formatDate(session.date, 'yyyy-MM-dd HH:mm:ss'),
      session.notes ?? '',
    ];
  });
  return [headers, ...rows].map(row => row.map(escapeCSV).join(',')).join('\n');
};

export const notesToCSV = (notes: Note[], books: Book[]): string => {
  const bookMap = new Map(books.map(b => [b.id, b]));
  const headers = ['ID', '书名', '页码', '章节', '类型', '内容', '高亮颜色', '创建时间', '更新时间'];
  const rows = notes.map(note => {
    const book = bookMap.get(note.bookId);
    return [
      note.id,
      book?.title ?? note.bookId,
      note.page,
      note.chapter ?? '',
      note.type,
      note.content.replace(/\n/g, ' '),
      note.highlightColor ?? '',
      formatDate(note.createdAt, 'yyyy-MM-dd HH:mm:ss'),
      formatDate(note.updatedAt, 'yyyy-MM-dd HH:mm:ss'),
    ];
  });
  return [headers, ...rows].map(row => row.map(escapeCSV).join(',')).join('\n');
};

export const downloadFile = (content: string, filename: string, mimeType: string): void => {
  try {
    const blob = new Blob([content], { type: mimeType });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = filename;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    URL.revokeObjectURL(url);
    logger.info(`File downloaded: ${filename}`);
  } catch (error) {
    logger.error('Failed to download file', error);
    throw new Error('文件下载失败');
  }
};

export const generateExportFilename = (type: string, format: 'json' | 'csv'): string => {
  const date = formatDate(new Date(), 'yyyyMMdd_HHmmss');
  return `reading_tracker_${type}_${date}.${format}`;
};

export const exportData = (
  options: ExportOptions,
  data: {
    books: Book[];
    readingPlans: ReadingPlan[];
    readingSessions: ReadingSession[];
    notes: Note[];
    bookmarks: Bookmark[];
  }
): void => {
  const exportData: ExportData = {
    exportedAt: new Date().toISOString(),
    version: '1.0.0',
  };

  if (options.includeBooks) exportData.books = data.books;
  if (options.includeReadingSessions) {
    exportData.readingPlans = data.readingPlans;
    exportData.readingSessions = data.readingSessions;
  }
  if (options.includeNotes) {
    exportData.notes = data.notes;
    exportData.bookmarks = data.bookmarks;
  }

  if (options.format === 'json') {
    exportToJSON(exportData, generateExportFilename('full', 'json'));
  } else {
    const parts: string[] = [];
    if (options.includeBooks) {
      parts.push('=== 书籍 ===\n' + booksToCSV(data.books));
    }
    if (options.includeReadingSessions) {
      parts.push('\n=== 阅读记录 ===\n' + readingSessionsToCSV(data.readingSessions, data.books));
    }
    if (options.includeNotes) {
      parts.push('\n=== 笔记 ===\n' + notesToCSV(data.notes, data.books));
    }
    downloadFile(parts.join('\n\n'), generateExportFilename('full', 'csv'), 'text/csv');
  }
};
