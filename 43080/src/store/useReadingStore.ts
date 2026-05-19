import { create } from 'zustand';
import type {
  ReadingPlan,
  ReadingSession,
  ReadingPlanFormData,
  ReadingSessionFormData,
  ReadingStats,
  DailyReadingStats,
  WeeklyStats,
  MonthlyStats,
} from '../types/reading';
import { STORAGE_KEYS } from '../constants/config';
import { safeStorage } from '../utils/storage';
import { logger } from '../utils/logger';
import { validateData, readingPlanSchema, readingSessionSchema, validateReadingSession } from '../utils/validation';
import { calculateStreak, getWeekRange, getMonthRange, isDateInRange, generateDateRange, formatDate } from '../utils/date';
import { useBookStore } from './useBookStore';

interface ReadingState {
  plans: ReadingPlan[];
  sessions: ReadingSession[];
}

interface ReadingActions {
  addPlan: (bookId: string, planData: ReadingPlanFormData) => ReadingPlan;
  updatePlan: (planId: string, planData: Partial<ReadingPlanFormData>) => void;
  deletePlan: (planId: string) => void;
  getPlanByBookId: (bookId: string) => ReadingPlan | undefined;
  addSession: (bookId: string, sessionData: ReadingSessionFormData) => ReadingSession;
  updateSession: (sessionId: string, sessionData: Partial<ReadingSessionFormData>) => void;
  deleteSession: (sessionId: string) => void;
  getSessionsByBookId: (bookId: string) => ReadingSession[];
  getReadingStats: () => ReadingStats;
  calculateProgress: (bookId: string) => number;
  calculateEstimatedDays: (bookId: string) => number | null;
  loadReadingData: () => void;
  saveReadingData: () => void;
}

const generateId = (): string => {
  return `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
};

const initialState: ReadingState = {
  plans: [],
  sessions: [],
};

export const useReadingStore = create<ReadingState & ReadingActions>((set, get) => ({
  ...initialState,

  addPlan: (bookId, planData) => {
    const existingPlan = get().getPlanByBookId(bookId);
    if (existingPlan) {
      throw new Error('该书已有阅读计划');
    }

    const newPlan: ReadingPlan = {
      id: generateId(),
      bookId,
      ...planData,
    };

    const validated = validateData(readingPlanSchema, newPlan);
    if (!validated) {
      throw new Error('阅读计划数据验证失败');
    }

    set((state) => ({
      plans: [...state.plans, newPlan],
    }));

    logger.info('Reading plan added', newPlan);
    get().saveReadingData();
    return newPlan;
  },

  updatePlan: (planId, planData) => {
    set((state) => ({
      plans: state.plans.map((plan) =>
        plan.id === planId ? { ...plan, ...planData } : plan
      ),
    }));
    logger.info(`Reading plan updated: ${planId}`, planData);
    get().saveReadingData();
  },

  deletePlan: (planId) => {
    set((state) => ({
      plans: state.plans.filter((plan) => plan.id !== planId),
    }));
    logger.info(`Reading plan deleted: ${planId}`);
    get().saveReadingData();
  },

  getPlanByBookId: (bookId) => {
    return get().plans.find((plan) => plan.bookId === bookId);
  },

  addSession: (bookId, sessionData) => {
    const book = useBookStore.getState().getBookById(bookId);
    if (!book) {
      throw new Error('书籍不存在');
    }

    const validation = validateReadingSession(
      sessionData.startPage,
      sessionData.endPage,
      book.totalPages
    );
    if (!validation.valid) {
      throw new Error(validation.error);
    }

    const newSession: ReadingSession = {
      id: generateId(),
      bookId,
      ...sessionData,
    };

    const validated = validateData(readingSessionSchema, newSession);
    if (!validated) {
      throw new Error('阅读记录数据验证失败');
    }

    set((state) => ({
      sessions: [...state.sessions, newSession],
    }));

    useBookStore.getState().updateCurrentPage(bookId, sessionData.endPage);

    logger.info('Reading session added', newSession);
    get().saveReadingData();
    return newSession;
  },

  updateSession: (sessionId, sessionData) => {
    const session = get().sessions.find((s) => s.id === sessionId);
    if (!session) return;

    const book = useBookStore.getState().getBookById(session.bookId);
    if (book && sessionData.endPage !== undefined) {
      const validation = validateReadingSession(
        sessionData.startPage ?? session.startPage,
        sessionData.endPage,
        book.totalPages
      );
      if (!validation.valid) {
        throw new Error(validation.error);
      }
    }

    set((state) => ({
      sessions: state.sessions.map((s) =>
        s.id === sessionId ? { ...s, ...sessionData } : s
      ),
    }));
    logger.info(`Reading session updated: ${sessionId}`, sessionData);
    get().saveReadingData();
  },

  deleteSession: (sessionId) => {
    set((state) => ({
      sessions: state.sessions.filter((s) => s.id !== sessionId),
    }));
    logger.info(`Reading session deleted: ${sessionId}`);
    get().saveReadingData();
  },

  getSessionsByBookId: (bookId) => {
    return get().sessions
      .filter((s) => s.bookId === bookId)
      .sort((a, b) => new Date(b.date).getTime() - new Date(a.date).getTime());
  },

  calculateProgress: (bookId) => {
    const book = useBookStore.getState().getBookById(bookId);
    if (!book || book.totalPages <= 0) return 0;
    return Math.min(Math.round((book.currentPage / book.totalPages) * 100), 100);
  },

  calculateEstimatedDays: (bookId) => {
    const book = useBookStore.getState().getBookById(bookId);
    const plan = get().getPlanByBookId(bookId);
    if (!book || !plan) return null;

    const remainingPages = book.totalPages - book.currentPage;
    if (remainingPages <= 0) return 0;

    return Math.ceil(remainingPages / plan.dailyPages);
  },

  getReadingStats: () => {
    const { sessions } = get();
    const { books } = useBookStore.getState();

    const totalPagesRead = sessions.reduce(
      (sum, s) => sum + (s.endPage - s.startPage),
      0
    );
    const totalReadingTime = sessions.reduce((sum, s) => sum + s.duration, 0);
    const booksCompleted = books.filter((b) => b.status === 'completed').length;
    const averageReadingSpeed = totalReadingTime > 0 ? Math.round(totalPagesRead / (totalReadingTime / 60)) : 0;

    const readingDates = sessions.map((s) => s.date);
    const { current: currentStreak, longest: longestStreak } = calculateStreak(readingDates);

    const dailyStatsMap = new Map<string, DailyReadingStats>();
    sessions.forEach((session) => {
      const date = formatDate(session.date);
      const existing = dailyStatsMap.get(date) || {
        date,
        duration: 0,
        pagesRead: 0,
      };
      dailyStatsMap.set(date, {
        date,
        duration: existing.duration + session.duration,
        pagesRead: existing.pagesRead + (session.endPage - session.startPage),
      });
    });
    const dailyStats = Array.from(dailyStatsMap.values()).sort(
      (a, b) => new Date(a.date).getTime() - new Date(b.date).getTime()
    );

    const now = new Date();
    const weeklyStats: WeeklyStats[] = [];
    for (let i = 11; i >= 0; i--) {
      const weekDate = new Date(now);
      weekDate.setDate(weekDate.getDate() - i * 7);
      const { start, end } = getWeekRange(weekDate);
      const weekSessions = sessions.filter((s) =>
        isDateInRange(new Date(s.date), start, end)
      );
      const completedInWeek = books.filter((b) =>
        b.status === 'completed' && isDateInRange(new Date(b.updatedAt), start, end)
      ).length;

      weeklyStats.push({
        weekStart: formatDate(start),
        totalDuration: weekSessions.reduce((sum, s) => sum + s.duration, 0),
        totalPages: weekSessions.reduce(
          (sum, s) => sum + (s.endPage - s.startPage),
          0
        ),
        booksCompleted: completedInWeek,
      });
    }

    const monthlyStats: MonthlyStats[] = [];
    for (let i = 5; i >= 0; i--) {
      const monthDate = new Date(now);
      monthDate.setMonth(monthDate.getMonth() - i);
      const { start, end } = getMonthRange(monthDate);
      const monthSessions = sessions.filter((s) =>
        isDateInRange(new Date(s.date), start, end)
      );
      const completedInMonth = books.filter((b) =>
        b.status === 'completed' && isDateInRange(new Date(b.updatedAt), start, end)
      ).length;

      monthlyStats.push({
        month: formatDate(start, 'yyyy-MM'),
        totalDuration: monthSessions.reduce((sum, s) => sum + s.duration, 0),
        totalPages: monthSessions.reduce(
          (sum, s) => sum + (s.endPage - s.startPage),
          0
        ),
        booksCompleted: completedInMonth,
      });
    }

    return {
      totalBooksRead: booksCompleted,
      totalPagesRead,
      totalReadingTime,
      averageReadingSpeed,
      currentStreak,
      longestStreak,
      weeklyStats,
      monthlyStats,
      dailyStats,
    };
  },

  loadReadingData: () => {
    try {
      const storedPlans = safeStorage.get<ReadingPlan[]>(STORAGE_KEYS.READING_PLANS, []);
      const storedSessions = safeStorage.get<ReadingSession[]>(STORAGE_KEYS.READING_SESSIONS, []);

      const validatedPlans = storedPlans.filter((p) => validateData(readingPlanSchema, p) !== null);
      const validatedSessions = storedSessions.filter((s) => validateData(readingSessionSchema, s) !== null);

      set({ plans: validatedPlans, sessions: validatedSessions });
      logger.info(`Loaded ${validatedPlans.length} plans and ${validatedSessions.length} sessions`);
    } catch (error) {
      logger.error('Failed to load reading data', error);
    }
  },

  saveReadingData: () => {
    const { plans, sessions } = get();
    const plansSuccess = safeStorage.set(STORAGE_KEYS.READING_PLANS, plans);
    const sessionsSuccess = safeStorage.set(STORAGE_KEYS.READING_SESSIONS, sessions);
    if (!plansSuccess || !sessionsSuccess) {
      logger.warn('Failed to save reading data to storage');
    }
  },
}));
