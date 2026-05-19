import { create } from 'zustand';
import type { Note, NoteFormData, Bookmark, BookmarkFormData, NoteType } from '../types/note';
import { STORAGE_KEYS } from '../constants/config';
import { safeStorage } from '../utils/storage';
import { logger } from '../utils/logger';
import { validateData, noteSchema } from '../utils/validation';

interface NoteState {
  notes: Note[];
  bookmarks: Bookmark[];
  filterBookId: string | null;
  filterType: NoteType | 'all';
  searchQuery: string;
}

interface NoteActions {
  addNote: (bookId: string, noteData: NoteFormData) => Note;
  updateNote: (noteId: string, noteData: Partial<NoteFormData>) => void;
  deleteNote: (noteId: string) => void;
  addBookmark: (bookId: string, bookmarkData: BookmarkFormData) => Bookmark;
  updateBookmark: (bookmarkId: string, bookmarkData: Partial<BookmarkFormData>) => void;
  deleteBookmark: (bookmarkId: string) => void;
  getNotesByBookId: (bookId: string) => Note[];
  getBookmarksByBookId: (bookId: string) => Bookmark[];
  getHighlightsByBookId: (bookId: string) => Note[];
  setFilterBookId: (bookId: string | null) => void;
  setFilterType: (type: NoteType | 'all') => void;
  setSearchQuery: (query: string) => void;
  getFilteredNotes: () => Note[];
  loadNotes: () => void;
  saveNotes: () => void;
}

const generateId = (): string => {
  return `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
};

const initialState: NoteState = {
  notes: [],
  bookmarks: [],
  filterBookId: null,
  filterType: 'all',
  searchQuery: '',
};

export const useNoteStore = create<NoteState & NoteActions>((set, get) => ({
  ...initialState,

  addNote: (bookId, noteData) => {
    const newNote: Note = {
      id: generateId(),
      bookId,
      ...noteData,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    };

    const validated = validateData(noteSchema, newNote);
    if (!validated) {
      throw new Error('笔记数据验证失败');
    }

    set((state) => ({
      notes: [...state.notes, newNote],
    }));

    logger.info('Note added', newNote);
    get().saveNotes();
    return newNote;
  },

  updateNote: (noteId, noteData) => {
    set((state) => ({
      notes: state.notes.map((note) =>
        note.id === noteId
          ? { ...note, ...noteData, updatedAt: new Date().toISOString() }
          : note
      ),
    }));
    logger.info(`Note updated: ${noteId}`, noteData);
    get().saveNotes();
  },

  deleteNote: (noteId) => {
    set((state) => ({
      notes: state.notes.filter((note) => note.id !== noteId),
    }));
    logger.info(`Note deleted: ${noteId}`);
    get().saveNotes();
  },

  addBookmark: (bookId, bookmarkData) => {
    const newBookmark: Bookmark = {
      id: generateId(),
      bookId,
      ...bookmarkData,
      createdAt: new Date().toISOString(),
    };

    set((state) => ({
      bookmarks: [...state.bookmarks, newBookmark],
    }));

    logger.info('Bookmark added', newBookmark);
    get().saveNotes();
    return newBookmark;
  },

  updateBookmark: (bookmarkId, bookmarkData) => {
    set((state) => ({
      bookmarks: state.bookmarks.map((bookmark) =>
        bookmark.id === bookmarkId ? { ...bookmark, ...bookmarkData } : bookmark
      ),
    }));
    logger.info(`Bookmark updated: ${bookmarkId}`, bookmarkData);
    get().saveNotes();
  },

  deleteBookmark: (bookmarkId) => {
    set((state) => ({
      bookmarks: state.bookmarks.filter((bookmark) => bookmark.id !== bookmarkId),
    }));
    logger.info(`Bookmark deleted: ${bookmarkId}`);
    get().saveNotes();
  },

  getNotesByBookId: (bookId) => {
    return get().notes
      .filter((n) => n.bookId === bookId && n.type === 'note')
      .sort((a, b) => b.page - a.page);
  },

  getHighlightsByBookId: (bookId) => {
    return get().notes
      .filter((n) => n.bookId === bookId && n.type === 'highlight')
      .sort((a, b) => b.page - a.page);
  },

  getBookmarksByBookId: (bookId) => {
    return get().bookmarks
      .filter((b) => b.bookId === bookId)
      .sort((a, b) => b.page - a.page);
  },

  setFilterBookId: (bookId) => set({ filterBookId: bookId }),
  setFilterType: (type) => set({ filterType: type }),
  setSearchQuery: (query) => set({ searchQuery: query }),

  getFilteredNotes: () => {
    const { notes, filterBookId, filterType, searchQuery } = get();
    let filtered = [...notes];

    if (filterBookId) {
      filtered = filtered.filter((n) => n.bookId === filterBookId);
    }

    if (filterType !== 'all') {
      filtered = filtered.filter((n) => n.type === filterType);
    }

    if (searchQuery) {
      const query = searchQuery.toLowerCase();
      filtered = filtered.filter(
        (n) =>
          n.content.toLowerCase().includes(query) ||
          (n.chapter && n.chapter.toLowerCase().includes(query))
      );
    }

    return filtered.sort(
      (a, b) => new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime()
    );
  },

  loadNotes: () => {
    try {
      const storedNotes = safeStorage.get<Note[]>(STORAGE_KEYS.NOTES, []);
      const storedBookmarks = safeStorage.get<Bookmark[]>(STORAGE_KEYS.BOOKMARKS, []);

      const validatedNotes = storedNotes.filter((n) => validateData(noteSchema, n) !== null);

      set({ notes: validatedNotes, bookmarks: storedBookmarks });
      logger.info(`Loaded ${validatedNotes.length} notes and ${storedBookmarks.length} bookmarks`);
    } catch (error) {
      logger.error('Failed to load notes', error);
    }
  },

  saveNotes: () => {
    const { notes, bookmarks } = get();
    const notesSuccess = safeStorage.set(STORAGE_KEYS.NOTES, notes);
    const bookmarksSuccess = safeStorage.set(STORAGE_KEYS.BOOKMARKS, bookmarks);
    if (!notesSuccess || !bookmarksSuccess) {
      logger.warn('Failed to save notes to storage');
    }
  },
}));
