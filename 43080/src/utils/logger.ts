import { MAX_LOG_ENTRIES } from '../constants/config';

export type LogLevel = 'debug' | 'info' | 'warn' | 'error';

interface LogEntry {
  timestamp: string;
  level: LogLevel;
  message: string;
  data?: unknown;
}

const logs: LogEntry[] = [];

const formatTimestamp = (date: Date): string => {
  return date.toISOString();
};

const addLog = (level: LogLevel, message: string, data?: unknown): void => {
  const entry: LogEntry = {
    timestamp: formatTimestamp(new Date()),
    level,
    message,
    data,
  };

  logs.push(entry);

  if (logs.length > MAX_LOG_ENTRIES) {
    logs.shift();
  }

  if (import.meta.env.DEV) {
    const consoleMethod = level === 'debug' ? 'log' : level;
    const logFn = console[consoleMethod as 'log' | 'info' | 'warn' | 'error'];
    if (typeof logFn === 'function') {
      logFn(`[${level.toUpperCase()}] ${message}`, data ?? '');
    }
  }
};

export const logger = {
  debug: (message: string, data?: unknown) => addLog('debug', message, data),
  info: (message: string, data?: unknown) => addLog('info', message, data),
  warn: (message: string, data?: unknown) => addLog('warn', message, data),
  error: (message: string, data?: unknown) => addLog('error', message, data),
  getLogs: () => [...logs],
  clearLogs: () => {
    logs.length = 0;
  },
  exportLogs: () => JSON.stringify(logs, null, 2),
};
