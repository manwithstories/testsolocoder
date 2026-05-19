import request from './request'
import type { Tag, Category, ReadingNote, BorrowRecord, ReadingGoal, GoalProgress } from '@/types'

export const getTags = () => request.get<any, Tag[]>('/tags')
export const getTag = (id: number) => request.get<any, Tag>(`/tags/${id}`)
export const createTag = (data: { name: string; color?: string }) => request.post<any, Tag>('/tags', data)
export const updateTag = (id: number, data: { name: string; color?: string }) => request.put<any, Tag>(`/tags/${id}`, data)
export const deleteTag = (id: number) => request.delete(`/tags/${id}`)

export const getCategories = (flat = false) => request.get<any, Category[]>(`/categories?flat=${flat}`)
export const getCategory = (id: number) => request.get<any, Category>(`/categories/${id}`)
export const createCategory = (data: { name: string; parent_id?: number }) => request.post<any, Category>('/categories', data)
export const updateCategory = (id: number, data: { name: string; parent_id?: number }) => request.put<any, Category>(`/categories/${id}`, data)
export const deleteCategory = (id: number) => request.delete(`/categories/${id}`)

export const getNotesByBook = (bookId: number) => request.get<any, ReadingNote[]>(`/notes/book/${bookId}`)
export const getNote = (id: number) => request.get<any, ReadingNote>(`/notes/${id}`)
export const createNote = (data: { book_id: number; page?: number; content: string }) => request.post<any, ReadingNote>('/notes', data)
export const updateNote = (id: number, data: { page?: number; content?: string }) => request.put<any, ReadingNote>(`/notes/${id}`, data)
export const deleteNote = (id: number) => request.delete(`/notes/${id}`)

export const getBorrows = (status?: string) => request.get<any, BorrowRecord[]>(`/borrows${status ? `?status=${status}` : ''}`)
export const getOverdueBorrows = () => request.get<any, BorrowRecord[]>('/borrows/overdue')
export const getBorrowByBook = (bookId: number) => request.get<any, BorrowRecord>(`/borrows/book/${bookId}`)
export const getBorrow = (id: number) => request.get<any, BorrowRecord>(`/borrows/${id}`)
export const createBorrow = (data: {
  book_id: number
  borrower_name: string
  borrower_phone?: string
  borrower_email?: string
  borrow_date?: string
  expected_return_date?: string
  notes?: string
}) => request.post<any, BorrowRecord>('/borrows', data)
export const returnBook = (id: number, returnDate?: string) => request.post<any, BorrowRecord>(`/borrows/${id}/return`, { return_date: returnDate })
export const deleteBorrow = (id: number) => request.delete(`/borrows/${id}`)

export const getGoals = (year?: number) => request.get<any, ReadingGoal[]>(`/goals${year ? `?year=${year}` : ''}`)
export const getGoal = (id: number) => request.get<any, GoalProgress>(`/goals/${id}`)
export const createGoal = (data: { year: number; month?: number; target_books?: number; target_pages?: number }) => request.post<any, ReadingGoal>('/goals', data)
export const updateGoal = (id: number, data: { target_books?: number; target_pages?: number }) => request.put<any, ReadingGoal>(`/goals/${id}`, data)
export const deleteGoal = (id: number) => request.delete(`/goals/${id}`)
export const getYearlyGoalProgress = (year: number) => request.get<any, {
  yearly: GoalProgress
  monthly: GoalProgress[]
}>(`/goals/yearly-progress?year=${year}`)
