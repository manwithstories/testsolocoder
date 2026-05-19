import { logger } from './logger';

export const safeStorage = {
  get: <T>(key: string, defaultValue: T): T => {
    try {
      const item = localStorage.getItem(key);
      if (!item) return defaultValue;
      return JSON.parse(item) as T;
    } catch (error) {
      logger.error(`Failed to read from localStorage: ${key}`, error);
      return defaultValue;
    }
  },

  set: <T>(key: string, value: T): boolean => {
    try {
      localStorage.setItem(key, JSON.stringify(value));
      return true;
    } catch (error) {
      logger.error(`Failed to write to localStorage: ${key}`, error);
      if (error instanceof Error && error.name === 'QuotaExceededError') {
        logger.warn('LocalStorage quota exceeded');
      }
      return false;
    }
  },

  remove: (key: string): boolean => {
    try {
      localStorage.removeItem(key);
      return true;
    } catch (error) {
      logger.error(`Failed to remove from localStorage: ${key}`, error);
      return false;
    }
  },

  clear: (): boolean => {
    try {
      localStorage.clear();
      return true;
    } catch (error) {
      logger.error('Failed to clear localStorage', error);
      return false;
    }
  },

  getSize: (): number => {
    try {
      let total = 0;
      for (let i = 0; i < localStorage.length; i++) {
        const key = localStorage.key(i);
        if (key) {
          total += localStorage.getItem(key)?.length ?? 0;
        }
      }
      return total;
    } catch {
      return 0;
    }
  },
};

export const createInMemoryStorage = <T>(initialValue: T) => {
  let data: T = initialValue;

  return {
    get: (): T => data,
    set: (value: T): void => {
      data = value;
    },
  };
};
