export interface Card {
  id: string;
  title: string;
  content: string;
  tags: string[];
  createdAt: number;
  updatedAt: number;
  linkedCardIds: string[];
  review: ReviewInfo;
  order: number;
}

export interface ReviewInfo {
  nextReviewAt: number;
  interval: number;
  easeFactor: number;
  repetitions: number;
  lastReviewedAt: number | null;
}

export interface Tag {
  id: string;
  name: string;
  parentId: string | null;
  color: string;
}

export interface Theme {
  mode: 'light' | 'dark';
}

export interface AppState {
  cards: Card[];
  tags: Tag[];
  selectedCardId: string | null;
  selectedTagIds: string[];
  searchQuery: string;
  theme: Theme;
  view: 'cards' | 'graph' | 'review' | 'stats' | 'tags';
}

export type AppAction =
  | { type: 'ADD_CARD'; payload: Card }
  | { type: 'UPDATE_CARD'; payload: Card }
  | { type: 'DELETE_CARD'; payload: string }
  | { type: 'SELECT_CARD'; payload: string | null }
  | { type: 'ADD_TAG'; payload: Tag }
  | { type: 'UPDATE_TAG'; payload: Tag }
  | { type: 'DELETE_TAG'; payload: string }
  | { type: 'MERGE_TAGS'; payload: { sourceId: string; targetId: string } }
  | { type: 'TOGGLE_TAG_SELECTION'; payload: string }
  | { type: 'SET_SEARCH_QUERY'; payload: string }
  | { type: 'SET_THEME'; payload: Theme }
  | { type: 'SET_VIEW'; payload: AppState['view'] }
  | { type: 'LINK_CARDS'; payload: { cardId1: string; cardId2: string } }
  | { type: 'UNLINK_CARDS'; payload: { cardId1: string; cardId2: string } }
  | { type: 'REVIEW_CARD'; payload: { cardId: string; quality: number } }
  | { type: 'REORDER_CARDS'; payload: { sourceId: string; targetId: string } }
  | { type: 'LOAD_STATE'; payload: AppState };

export interface ValidationResult {
  isValid: boolean;
  errors: string[];
}

export interface LogEntry {
  timestamp: number;
  action: string;
  details: unknown;
}
