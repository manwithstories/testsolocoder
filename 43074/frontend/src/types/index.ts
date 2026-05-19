export type ReadingStatus = 'to_read' | 'reading' | 'completed' | 'abandoned'

export interface Book {
  id: number
  title: string
  author: string
  publisher: string
  isbn: string
  cover_image: string
  summary: string
  total_pages: number
  reading_status: ReadingStatus
  current_page: number
  reading_progress: number
  start_date?: string
  end_date?: string
  total_read_time: number
  tags?: Tag[]
  categories?: Category[]
  reading_notes?: ReadingNote[]
  borrow_record?: BorrowRecord
  created_at: string
  updated_at: string
}

export interface Tag {
  id: number
  name: string
  color: string
  created_at: string
}

export interface Category {
  id: number
  name: string
  parent_id?: number
  parent?: Category
  children?: Category[]
  created_at: string
}

export interface ReadingNote {
  id: number
  book_id: number
  page: number
  content: string
  created_at: string
  updated_at: string
}

export interface BorrowRecord {
  id: number
  book_id: number
  borrower_name: string
  borrower_phone: string
  borrower_email: string
  borrow_date: string
  expected_return_date?: string
  return_date?: string
  status: string
  notes: string
  reminder_sent: boolean
  created_at: string
  updated_at: string
}

export interface ReadingGoal {
  id: number
  year: number
  month?: number
  target_books: number
  target_pages: number
  created_at: string
  updated_at: string
}

export interface GoalProgress extends ReadingGoal {
  completed_books: number
  completed_pages: number
  book_progress: number
  page_progress: number
}

export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  size: number
}

export interface ApiError {
  code: number
  message: string
  details?: any
}

export interface BookInfo {
  title: string
  author: string
  publisher: string
  isbn: string
  cover_url: string
  summary: string
  pages: number
}

export interface OverviewStats {
  total_books: number
  reading_books: number
  completed_books: number
  total_pages_read: number
  total_read_time: number
  currently_borrowed: number
}

export interface MonthlyStats {
  month: string
  completed: number
  pages_read: number
}

export interface HeatmapData {
  date: string
  count: number
}

export interface DurationStats {
  range: string
  count: number
  category: string
}

export interface CategoryStats {
  name: string
  count: number
}

export interface TagStats {
  name: string
  color: string
  count: number
}
