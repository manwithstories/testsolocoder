import React, { createContext, useContext, useReducer, useEffect, ReactNode } from 'react';
import type { AppState, AppAction } from '../types';
import { initialState, appReducer } from './reducer';
import { saveState, loadState } from '../utils/storage';
import { logAction } from '../utils/logger';

interface AppContextType {
  state: AppState;
  dispatch: React.Dispatch<AppAction>;
}

const AppContext = createContext<AppContextType | undefined>(undefined);

interface AppProviderProps {
  children: ReactNode;
}

export const AppProvider: React.FC<AppProviderProps> = ({ children }) => {
  const [state, dispatch] = useReducer(appReducer, initialState);

  useEffect(() => {
    const savedState = loadState();
    if (savedState) {
      dispatch({ type: 'LOAD_STATE', payload: savedState });
      logAction('STATE_RESTORED_FROM_STORAGE');
    }
  }, []);

  useEffect(() => {
    if (state.cards.length > 0 || state.tags.length > 0 || state !== initialState) {
      saveState(state);
    }
  }, [state]);

  return (
    <AppContext.Provider value={{ state, dispatch }}>
      {children}
    </AppContext.Provider>
  );
};

export const useApp = (): AppContextType => {
  const context = useContext(AppContext);
  if (context === undefined) {
    throw new Error('useApp must be used within an AppProvider');
  }
  return context;
};
