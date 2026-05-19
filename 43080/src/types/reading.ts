export interface ReadingPlan {
  id: string;
  bookId: string;
  dailyPages: number;
  startDate: string;
  targetDate: string;
}

export interface ReadingSession {
  id: string;
  bookId: string;
  startPage: number;
  endPage: number;
  duration: number;
  date: string;
  notes?: string;
}

export interface ReadingSessionFormData {
  startPage: number;
  endPage: number;
  duration: number;
  date: string;
  notes?: string;
}

export interface ReadingPlanFormData {
  dailyPages: number;
  startDate: string;
  targetDate: string;
}

export interface DailyReadingStats {
  date: string;
  duration: number;
  pagesRead: number;
}

export interface WeeklyStats {
  weekStart: string;
  totalDuration: number;
  totalPages: number;
  booksCompleted: number;
}

export interface MonthlyStats {
  month: string;
  totalDuration: number;
  totalPages: number;
  booksCompleted: number;
}

export interface ReadingStats {
  totalBooksRead: number;
  totalPagesRead: number;
  totalReadingTime: number;
  averageReadingSpeed: number;
  currentStreak: number;
  longestStreak: number;
  weeklyStats: WeeklyStats[];
  monthlyStats: MonthlyStats[];
  dailyStats: DailyReadingStats[];
}
