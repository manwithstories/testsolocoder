import React, { useCallback } from 'react';
import { AppProvider, useApp } from './store/context';
import { ToastProvider } from './components/common/Toast';
import ErrorBoundary from './components/common/ErrorBoundary';
import Sidebar from './components/layout/Sidebar';
import Header from './components/layout/Header';
import CardList from './components/cards/CardList';
import TagManager from './components/tags/TagManager';
import KnowledgeGraph from './components/graph/KnowledgeGraph';
import ReviewSystem from './components/review/ReviewSystem';
import StatsPanel from './components/stats/StatsPanel';
import { initialState } from './store/reducer';
import { logAction } from './utils/logger';

const AppContent: React.FC = () => {
  const { state, dispatch } = useApp();

  const handleSetView = useCallback((view: string) => {
    dispatch({ type: 'SET_VIEW', payload: view as typeof state.view });
    dispatch({ type: 'SELECT_CARD', payload: null });
  }, [dispatch]);

  const handleToggleTheme = useCallback(() => {
    const newMode = state.theme.mode === 'light' ? 'dark' : 'light';
    dispatch({ type: 'SET_THEME', payload: { mode: newMode } });
    logAction('THEME_TOGGLED', { mode: newMode });
  }, [dispatch, state.theme.mode]);

  const handleClearData = useCallback(() => {
    dispatch({ type: 'LOAD_STATE', payload: initialState });
  }, [dispatch]);

  const renderMainContent = () => {
    switch (state.view) {
      case 'cards':
        return (
          <CardList
            cards={state.cards}
            tags={state.tags}
            selectedTagIds={state.selectedTagIds}
            searchQuery={state.searchQuery}
            selectedCardId={state.selectedCardId}
            onSelectCard={(id) => dispatch({ type: 'SELECT_CARD', payload: id })}
            onAddCard={(card) => dispatch({ type: 'ADD_CARD', payload: card })}
            onUpdateCard={(card) => dispatch({ type: 'UPDATE_CARD', payload: card })}
            onDeleteCard={(id) => dispatch({ type: 'DELETE_CARD', payload: id })}
            onLinkCards={(cardId1, cardId2) =>
              dispatch({ type: 'LINK_CARDS', payload: { cardId1, cardId2 } })
            }
            onUnlinkCards={(cardId1, cardId2) =>
              dispatch({ type: 'UNLINK_CARDS', payload: { cardId1, cardId2 } })
            }
            onReviewCard={(cardId, quality) =>
              dispatch({ type: 'REVIEW_CARD', payload: { cardId, quality } })
            }
            onReorderCards={(sourceId, targetId) =>
              dispatch({ type: 'REORDER_CARDS', payload: { sourceId, targetId } })
            }
          />
        );
      case 'graph':
        return (
          <KnowledgeGraph
            cards={state.cards}
            tags={state.tags}
            selectedCardId={state.selectedCardId}
            onSelectCard={(id) => {
              dispatch({ type: 'SET_VIEW', payload: 'cards' });
              dispatch({ type: 'SELECT_CARD', payload: id });
            }}
          />
        );
      case 'review':
        return (
          <ReviewSystem
            cards={state.cards}
            tags={state.tags}
            onReviewCard={(cardId, quality) =>
              dispatch({ type: 'REVIEW_CARD', payload: { cardId, quality } })
            }
            onSelectCard={(id) => {
              dispatch({ type: 'SET_VIEW', payload: 'cards' });
              dispatch({ type: 'SELECT_CARD', payload: id });
            }}
          />
        );
      case 'tags':
        return (
          <TagManager
            tags={state.tags}
            cards={state.cards}
            selectedTagIds={state.selectedTagIds}
            onAddTag={(tag) => dispatch({ type: 'ADD_TAG', payload: tag })}
            onUpdateTag={(tag) => dispatch({ type: 'UPDATE_TAG', payload: tag })}
            onDeleteTag={(id) => dispatch({ type: 'DELETE_TAG', payload: id })}
            onMergeTags={(sourceId, targetId) =>
              dispatch({ type: 'MERGE_TAGS', payload: { sourceId, targetId } })
            }
            onToggleTagSelection={(id) =>
              dispatch({ type: 'TOGGLE_TAG_SELECTION', payload: id })
            }
          />
        );
      case 'stats':
        return <StatsPanel cards={state.cards} tags={state.tags} />;
      default:
        return null;
    }
  };

  return (
    <div className={`app theme-${state.theme.mode}`}>
      <Sidebar
        tags={state.tags}
        selectedTagIds={state.selectedTagIds}
        currentView={state.view}
        onToggleTag={(id) => dispatch({ type: 'TOGGLE_TAG_SELECTION', payload: id })}
        onSetView={handleSetView}
      />
      <div className="main-container">
        <Header
          searchQuery={state.searchQuery}
          theme={state.theme.mode}
          cards={state.cards}
          tags={state.tags}
          onSearchChange={(query) => dispatch({ type: 'SET_SEARCH_QUERY', payload: query })}
          onToggleTheme={handleToggleTheme}
          onClearData={handleClearData}
        />
        <main className="main-content">
          <ErrorBoundary>{renderMainContent()}</ErrorBoundary>
        </main>
      </div>
    </div>
  );
};

const App: React.FC = () => {
  return (
    <ToastProvider>
      <AppProvider>
        <AppContent />
      </AppProvider>
    </ToastProvider>
  );
};

export default App;
