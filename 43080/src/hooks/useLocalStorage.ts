import { useState, useEffect, useCallback } from 'react';
import { safeStorage } from '../utils/storage';
import { logger } from '../utils/logger';

export function useLocalStorage<T>(key: string, initialValue: T): [T, (value: T | ((prev: T) => T)) => void] {
  const [storedValue, setStoredValue] = useState<T>(initialValue);

  useEffect(() => {
    try {
      const item = safeStorage.get<T | null>(key, null);
      if (item !== null) {
        setStoredValue(item);
      }
    } catch (error) {
      logger.error(`Error reading localStorage key "${key}"`, error);
    }
  }, [key]);

  const setValue = useCallback(
    (value: T | ((prev: T) => T)) => {
      try {
        setStoredValue((prev) => {
          const valueToStore = value instanceof Function ? value(prev) : value;
          safeStorage.set(key, valueToStore);
          return valueToStore;
        });
      } catch (error) {
        logger.error(`Error setting localStorage key "${key}"`, error);
      }
    },
    [key]
  );

  return [storedValue, setValue];
}
