import request from './request'
import type { Book, PaginatedResponse, BookInfo } from '@/types'

export interface BookQueryParams {
  page?: number
  page_size?: number
  search?: string
  status?: string
  tag_ids?: number[]
  category_ids?: number[]
  sort_by?: string
  sort_order?: string
}

export interface CreateBookData {
  title: string
  author?: string
  publisher?: string
  isbn?: string
  summary?: string
  total_pages?: number
  tag_ids?: number[]
  category_ids?: number[]
}

export interface UpdateBookData {
  title?: string
  author?: string
  publisher?: string
  isbn?: string
  summary?: string
  total_pages?: number
  tag_ids?: number[]
  category_ids?: number[]
}

export const getBooks = (params: BookQueryParams = {}) => {
  return request.get<any, PaginatedResponse<Book>>('/books', { params })
}

export const getBook = (id: number) => {
  return request.get<any, Book>(`/books/${id}`)
}

export const createBook = (data: CreateBookData) => {
  return request.post<any, Book>('/books', data)
}

export const updateBook = (id: number, data: UpdateBookData) => {
  return request.put<any, Book>(`/books/${id}`, data)
}

export const deleteBook = (id: number) => {
  return request.delete(`/books/${id}`)
}

export const uploadCover = (id: number, file: File) => {
  const formData = new FormData()
  formData.append('cover', file)
  return request.post<any, { cover_image: string }>(`/books/${id}/cover`, formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

export const fetchBookByISBN = (isbn: string) => {
  return request.get<any, BookInfo>(`/books/isbn/${isbn}`)
}

export const updateReadingProgress = (id: number, currentPage: number) => {
  return request.patch<any, {
    current_page: number
    reading_progress: number
    reading_status: string
  }>(`/books/${id}/progress`, { current_page: currentPage })
}

export const updateBookStatus = (id: number, status: string) => {
  return request.patch<any, { reading_status: string }>(`/books/${id}/status?status=${status}`)
}

export const getCurrentlyReading = () => {
  return request.get<any, Book[]>('/books/reading')
}
