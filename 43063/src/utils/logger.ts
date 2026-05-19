import type { LogEntry } from '../types';

export const logAction = (action: string, details: unknown = {}): void => {
  const entry: LogEntry = {
    timestamp: Date.now(),
    action,
    details,
  };
  console.log('[KnowledgeNetwork]', entry);
};

export const logError = (action: string, error: unknown): void => {
  console.error('[KnowledgeNetwork][Error]', {
    timestamp: Date.now(),
    action,
    error,
  });
};

export const logWarning = (action: string, message: string): void => {
  console.warn('[KnowledgeNetwork][Warning]', {
    timestamp: Date.now(),
    action,
    message,
  });
};
