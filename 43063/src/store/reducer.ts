import type { AppState, AppAction, Card } from '../types';
import { calculateNextReview } from '../utils/spacedRepetition';
import { logAction } from '../utils/logger';

export const initialState: AppState = {
  cards: [],
  tags: [],
  selectedCardId: null,
  selectedTagIds: [],
  searchQuery: '',
  theme: { mode: 'light' },
  view: 'cards',
};

export const appReducer = (state: AppState, action: AppAction): AppState => {
  logAction(`ACTION: ${action.type}`, action.payload);

  switch (action.type) {
    case 'LOAD_STATE':
      return action.payload;

    case 'ADD_CARD': {
      const newCard = {
        ...action.payload,
        order: state.cards.length,
      };
      return {
        ...state,
        cards: [...state.cards, newCard],
      };
    }

    case 'UPDATE_CARD': {
      return {
        ...state,
        cards: state.cards.map((card) =>
          card.id === action.payload.id ? action.payload : card
        ),
      };
    }

    case 'DELETE_CARD': {
      const cardId = action.payload;
      return {
        ...state,
        cards: state.cards
          .filter((card) => card.id !== cardId)
          .map((card) => ({
            ...card,
            linkedCardIds: card.linkedCardIds.filter((id) => id !== cardId),
          })),
        selectedCardId: state.selectedCardId === cardId ? null : state.selectedCardId,
      };
    }

    case 'SELECT_CARD':
      return {
        ...state,
        selectedCardId: action.payload,
      };

    case 'ADD_TAG':
      return {
        ...state,
        tags: [...state.tags, action.payload],
      };

    case 'UPDATE_TAG':
      return {
        ...state,
        tags: state.tags.map((tag) =>
          tag.id === action.payload.id ? action.payload : tag
        ),
      };

    case 'DELETE_TAG': {
      const tagId = action.payload;
      return {
        ...state,
        tags: state.tags.filter((tag) => tag.id !== tagId),
        cards: state.cards.map((card) => ({
          ...card,
          tags: card.tags.filter((t) => t !== tagId),
        })),
        selectedTagIds: state.selectedTagIds.filter((id) => id !== tagId),
      };
    }

    case 'MERGE_TAGS': {
      const { sourceId, targetId } = action.payload;
      return {
        ...state,
        tags: state.tags.filter((tag) => tag.id !== sourceId),
        cards: state.cards.map((card) => {
          const hasSource = card.tags.includes(sourceId);
          const hasTarget = card.tags.includes(targetId);
          if (!hasSource) return card;
          return {
            ...card,
            tags: hasTarget
              ? card.tags.filter((t) => t !== sourceId)
              : card.tags.map((t) => (t === sourceId ? targetId : t)),
          };
        }),
        selectedTagIds: state.selectedTagIds.map((id) =>
          id === sourceId ? targetId : id
        ),
      };
    }

    case 'TOGGLE_TAG_SELECTION': {
      const tagId = action.payload;
      return {
        ...state,
        selectedTagIds: state.selectedTagIds.includes(tagId)
          ? state.selectedTagIds.filter((id) => id !== tagId)
          : [...state.selectedTagIds, tagId],
      };
    }

    case 'SET_SEARCH_QUERY':
      return {
        ...state,
        searchQuery: action.payload,
      };

    case 'SET_THEME':
      return {
        ...state,
        theme: action.payload,
      };

    case 'SET_VIEW':
      return {
        ...state,
        view: action.payload,
      };

    case 'LINK_CARDS': {
      const { cardId1, cardId2 } = action.payload;
      return {
        ...state,
        cards: state.cards.map((card) => {
          if (card.id === cardId1 && !card.linkedCardIds.includes(cardId2)) {
            return {
              ...card,
              linkedCardIds: [...card.linkedCardIds, cardId2],
            };
          }
          if (card.id === cardId2 && !card.linkedCardIds.includes(cardId1)) {
            return {
              ...card,
              linkedCardIds: [...card.linkedCardIds, cardId1],
            };
          }
          return card;
        }),
      };
    }

    case 'UNLINK_CARDS': {
      const { cardId1, cardId2 } = action.payload;
      return {
        ...state,
        cards: state.cards.map((card) => {
          if (card.id === cardId1) {
            return {
              ...card,
              linkedCardIds: card.linkedCardIds.filter((id) => id !== cardId2),
            };
          }
          if (card.id === cardId2) {
            return {
              ...card,
              linkedCardIds: card.linkedCardIds.filter((id) => id !== cardId1),
            };
          }
          return card;
        }),
      };
    }

    case 'REVIEW_CARD': {
      const { cardId, quality } = action.payload;
      return {
        ...state,
        cards: state.cards.map((card) => {
          if (card.id !== cardId) return card;
          const newReview = calculateNextReview(card.review, quality);
          return {
            ...card,
            review: newReview,
            updatedAt: Date.now(),
          };
        }),
      };
    }

    case 'REORDER_CARDS': {
      const { sourceId, targetId } = action.payload;
      const cards = [...state.cards];
      const sourceIndex = cards.findIndex((c) => c.id === sourceId);
      const targetIndex = cards.findIndex((c) => c.id === targetId);
      if (sourceIndex === -1 || targetIndex === -1) return state;

      const [removed] = cards.splice(sourceIndex, 1);
      cards.splice(targetIndex, 0, removed);

      const reorderedCards = cards.map((card, index) => ({
        ...card,
        order: index,
      })) as Card[];

      return {
        ...state,
        cards: reorderedCards,
      };
    }

    default:
      return state;
  }
};
