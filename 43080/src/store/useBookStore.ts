import { create } from 'zustand';
import type { Book, BookFormData, BookStatus } from '../types/book';
import { STORAGE_KEYS } from '../constants/config';
import { safeStorage } from '../utils/storage';
import { logger } from '../utils/logger';
import { validateData, bookSchema } from '../utils/validation';

interface BookState {
  books: Book[];
  searchQuery: string;
  filterStatus: BookStatus | 'all';
  filterCategory: string | null;
  sortBy: 'title' | 'author' | 'createdAt' | 'currentPage';
  sortOrder: 'asc' | 'desc';
  showFavoritesOnly: boolean;
  showReadingQueueOnly: boolean;
}

interface BookActions {
  addBook: (bookData: BookFormData) => Book;
  updateBook: (id: string, bookData: Partial<BookFormData>) => void;
  deleteBook: (id: string) => void;
  toggleFavorite: (id: string) => void;
  toggleReadingQueue: (id: string) => void;
  updateBookStatus: (id: string, status: BookStatus) => void;
  updateCurrentPage: (id: string, page: number) => void;
  setSearchQuery: (query: string) => void;
  setFilterStatus: (status: BookStatus | 'all') => void;
  setFilterCategory: (category: string | null) => void;
  setSortBy: (sortBy: BookState['sortBy']) => void;
  setSortOrder: (sortOrder: 'asc' | 'desc') => void;
  toggleFavoritesOnly: () => void;
  toggleReadingQueueOnly: () => void;
  getBookById: (id: string) => Book | undefined;
  getFilteredBooks: () => Book[];
  getCategories: () => string[];
  loadBooks: () => void;
  saveBooks: () => void;
}

const generateId = (): string => {
  return `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
};

const initialState: BookState = {
  books: [],
  searchQuery: '',
  filterStatus: 'all',
  filterCategory: null,
  sortBy: 'createdAt',
  sortOrder: 'desc',
  showFavoritesOnly: false,
  showReadingQueueOnly: false,
};

export const useBookStore = create<BookState & BookActions>((set, get) => ({
  ...initialState,

  addBook: (bookData) => {
    const newBook: Book = {
      id: generateId(),
      ...bookData,
      status: 'unread',
      isFavorite: false,
      inReadingQueue: false,
      currentPage: 0,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    };

    const validated = validateData(bookSchema, newBook);
    if (!validated) {
      throw new Error('书籍数据验证失败');
    }

    set((state) => ({
      books: [...state.books, newBook],
    }));

    logger.info('Book added', newBook);
    get().saveBooks();
    return newBook;
  },

  updateBook: (id, bookData) => {
    set((state) => ({
      books: state.books.map((book) =>
        book.id === id
          ? { ...book, ...bookData, updatedAt: new Date().toISOString() }
          : book
      ),
    }));
    logger.info(`Book updated: ${id}`, bookData);
    get().saveBooks();
  },

  deleteBook: (id) => {
    set((state) => ({
      books: state.books.filter((book) => book.id !== id),
    }));
    logger.info(`Book deleted: ${id}`);
    get().saveBooks();
  },

  toggleFavorite: (id) => {
    set((state) => ({
      books: state.books.map((book) =>
        book.id === id
          ? { ...book, isFavorite: !book.isFavorite, updatedAt: new Date().toISOString() }
          : book
      ),
    }));
    get().saveBooks();
  },

  toggleReadingQueue: (id) => {
    set((state) => ({
      books: state.books.map((book) =>
        book.id === id
          ? { ...book, inReadingQueue: !book.inReadingQueue, updatedAt: new Date().toISOString() }
          : book
      ),
    }));
    get().saveBooks();
  },

  updateBookStatus: (id, status) => {
    set((state) => ({
      books: state.books.map((book) =>
        book.id === id
          ? { ...book, status, updatedAt: new Date().toISOString() }
          : book
      ),
    }));
    get().saveBooks();
  },

  updateCurrentPage: (id, page) => {
    set((state) => ({
      books: state.books.map((book) => {
        if (book.id !== id) return book;
        const newStatus = page >= book.totalPages ? 'completed' : book.status === 'unread' ? 'reading' : book.status;
        return {
          ...book,
          currentPage: Math.min(page, book.totalPages),
          status: newStatus,
          updatedAt: new Date().toISOString(),
        };
      }),
    }));
    get().saveBooks();
  },

  setSearchQuery: (query) => set({ searchQuery: query }),
  setFilterStatus: (status) => set({ filterStatus: status }),
  setFilterCategory: (category) => set({ filterCategory: category }),
  setSortBy: (sortBy) => set({ sortBy }),
  setSortOrder: (sortOrder) => set({ sortOrder }),
  toggleFavoritesOnly: () => set((state) => ({ showFavoritesOnly: !state.showFavoritesOnly })),
  toggleReadingQueueOnly: () => set((state) => ({ showReadingQueueOnly: !state.showReadingQueueOnly })),

  getBookById: (id) => {
    return get().books.find((book) => book.id === id);
  },

  getFilteredBooks: () => {
    const { books, searchQuery, filterStatus, filterCategory, sortBy, sortOrder, showFavoritesOnly, showReadingQueueOnly } = get();

    let filtered = [...books];

    if (showFavoritesOnly) {
      filtered = filtered.filter((book) => book.isFavorite);
    }

    if (showReadingQueueOnly) {
      filtered = filtered.filter((book) => book.inReadingQueue);
    }

    if (filterStatus !== 'all') {
      filtered = filtered.filter((book) => book.status === filterStatus);
    }

    if (filterCategory) {
      filtered = filtered.filter((book) => book.categories.includes(filterCategory));
    }

    if (searchQuery) {
      const query = searchQuery.toLowerCase();
      filtered = filtered.filter(
        (book) =>
          book.title.toLowerCase().includes(query) ||
          book.author.toLowerCase().includes(query)
      );
    }

    filtered.sort((a, b) => {
      let comparison = 0;
      switch (sortBy) {
        case 'title':
          comparison = a.title.localeCompare(b.title);
          break;
        case 'author':
          comparison = a.author.localeCompare(b.author);
          break;
        case 'createdAt':
          comparison = new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime();
          break;
        case 'currentPage':
          comparison = a.currentPage - b.currentPage;
          break;
      }
      return sortOrder === 'asc' ? comparison : -comparison;
    });

    return filtered;
  },

  getCategories: () => {
    const categories = new Set<string>();
    get().books.forEach((book) => {
      book.categories.forEach((cat) => categories.add(cat));
    });
    return Array.from(categories).sort();
  },

  loadBooks: () => {
    try {
      const storedBooks = safeStorage.get<Book[]>(STORAGE_KEYS.BOOKS, []);
      const validatedBooks = storedBooks.filter((book) => validateData(bookSchema, book) !== null);
      set({ books: validatedBooks });
      logger.info(`Loaded ${validatedBooks.length} books from storage`);
    } catch (error) {
      logger.error('Failed to load books', error);
    }
  },

  saveBooks: () => {
    const { books } = get();
    const success = safeStorage.set(STORAGE_KEYS.BOOKS, books);
    if (!success) {
      logger.warn('Failed to save books to storage');
    }
  },
}));
