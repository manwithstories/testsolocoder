export type BookStatus = 'unread' | 'reading' | 'completed' | 'paused';

export interface Book {
  id: string;
  title: string;
  author: string;
  isbn?: string;
  coverUrl?: string;
  totalPages: number;
  categories: string[];
  description?: string;
  publishDate?: string;
  publisher?: string;
  status: BookStatus;
  isFavorite: boolean;
  inReadingQueue: boolean;
  currentPage: number;
  createdAt: string;
  updatedAt: string;
}

export interface BookFormData {
  title: string;
  author: string;
  isbn?: string;
  coverUrl?: string;
  totalPages: number;
  categories: string[];
  description?: string;
  publishDate?: string;
  publisher?: string;
}

export interface ISBNBookResult {
  title: string;
  author: string;
  isbn: string;
  coverUrl?: string;
  totalPages?: number;
  categories?: string[];
}
