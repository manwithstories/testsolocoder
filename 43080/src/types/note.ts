export type NoteType = 'note' | 'highlight' | 'bookmark';

export interface Note {
  id: string;
  bookId: string;
  page: number;
  chapter?: string;
  content: string;
  type: NoteType;
  highlightColor?: string;
  createdAt: string;
  updatedAt: string;
}

export interface Bookmark {
  id: string;
  bookId: string;
  page: number;
  label?: string;
  description?: string;
  createdAt: string;
}

export interface NoteFormData {
  page: number;
  chapter?: string;
  content: string;
  type: NoteType;
  highlightColor?: string;
}

export interface BookmarkFormData {
  page: number;
  label?: string;
  description?: string;
}
