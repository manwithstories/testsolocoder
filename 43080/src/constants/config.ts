export const STORAGE_KEYS = {
  BOOKS: 'reading_tracker_books',
  READING_PLANS: 'reading_tracker_plans',
  READING_SESSIONS: 'reading_tracker_sessions',
  NOTES: 'reading_tracker_notes',
  BOOKMARKS: 'reading_tracker_bookmarks',
  SETTINGS: 'reading_tracker_settings',
} as const;

export const READING_STATUS = {
  UNREAD: 'unread' as const,
  READING: 'reading' as const,
  COMPLETED: 'completed' as const,
  PAUSED: 'paused' as const,
};

export const NOTE_TYPES = {
  NOTE: 'note' as const,
  HIGHLIGHT: 'highlight' as const,
  BOOKMARK: 'bookmark' as const,
};

export const HIGHLIGHT_COLORS = [
  '#fef08a',
  '#bbf7d0',
  '#bfdbfe',
  '#fecaca',
  '#e9d5ff',
  '#fed7aa',
] as const;

export const DEFAULT_SETTINGS = {
  notificationEnabled: false,
  notificationTime: '20:00',
  theme: 'light' as const,
  defaultView: 'grid' as const,
};

export const ISBN_API_ENDPOINTS = [
  'https://openlibrary.org/api/books?bibkeys=ISBN:{isbn}&format=json&jscmd=data',
  'https://www.googleapis.com/books/v1/volumes?q=isbn:{isbn}',
] as const;

export const MAX_LOG_ENTRIES = 1000;

export const DEBOUNCE_DELAY = 300;
