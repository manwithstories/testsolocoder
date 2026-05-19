import type { LogEntry } from '@/types'

class Logger {
  private logs: LogEntry[] = []
  private maxLogs = 1000
  private enabled = true

  setEnabled(enabled: boolean) {
    this.enabled = enabled
  }

  private addLog(level: LogEntry['level'], message: string, data?: any) {
    if (!this.enabled) return

    const entry: LogEntry = {
      id: Date.now().toString(36) + Math.random().toString(36).substr(2, 9),
      timestamp: Date.now(),
      level,
      message,
      data
    }

    this.logs.push(entry)
    if (this.logs.length > this.maxLogs) {
      this.logs = this.logs.slice(-this.maxLogs)
    }

    if (level === 'error') {
      console.error(`[${level.toUpperCase()}] ${message}`, data || '')
    } else if (level === 'warn') {
      console.warn(`[${level.toUpperCase()}] ${message}`, data || '')
    } else {
      console.log(`[${level.toUpperCase()}] ${message}`, data || '')
    }
  }

  info(message: string, data?: any) {
    this.addLog('info', message, data)
  }

  warn(message: string, data?: any) {
    this.addLog('warn', message, data)
  }

  error(message: string, data?: any) {
    this.addLog('error', message, data)
  }

  getLogs(): LogEntry[] {
    return [...this.logs]
  }

  getLogsByLevel(level: LogEntry['level']): LogEntry[] {
    return this.logs.filter(log => log.level === level)
  }

  clear() {
    this.logs = []
  }

  export(): string {
    return JSON.stringify(this.logs, null, 2)
  }
}

export const logger = new Logger()
