import {
  format,
  startOfMonth,
  endOfMonth,
  startOfWeek,
  endOfWeek,
  eachDayOfInterval,
  isSameMonth,
  isSameDay,
  addMonths,
  subMonths,
  addWeeks,
  subWeeks,
  isToday,
  isWithinInterval,
  parseISO,
} from 'date-fns'
import { zhCN } from 'date-fns/locale'

export function formatDate(date: Date | string, pattern: string = 'yyyy-MM-dd'): string {
  const d = typeof date === 'string' ? parseISO(date) : date
  return format(d, pattern, { locale: zhCN })
}

export function formatDateReadable(date: Date | string): string {
  const d = typeof date === 'string' ? parseISO(date) : date
  return format(d, 'yyyy年MM月dd日', { locale: zhCN })
}

export function getMonthDays(date: Date): Date[] {
  const start = startOfWeek(startOfMonth(date), { weekStartsOn: 1 })
  const end = endOfWeek(endOfMonth(date), { weekStartsOn: 1 })
  return eachDayOfInterval({ start, end })
}

export function getWeekDays(date: Date): Date[] {
  const start = startOfWeek(date, { weekStartsOn: 1 })
  const end = endOfWeek(date, { weekStartsOn: 1 })
  return eachDayOfInterval({ start, end })
}

export function isCurrentMonth(day: Date, currentDate: Date): boolean {
  return isSameMonth(day, currentDate)
}

export function isTodayDate(day: Date): boolean {
  return isToday(day)
}

export function isSameDayDate(day1: Date, day2: Date): boolean {
  return isSameDay(day1, day2)
}

export function nextMonth(date: Date): Date {
  return addMonths(date, 1)
}

export function prevMonth(date: Date): Date {
  return subMonths(date, 1)
}

export function nextWeek(date: Date): Date {
  return addWeeks(date, 1)
}

export function prevWeek(date: Date): Date {
  return subWeeks(date, 1)
}

export function getMonthLabel(date: Date): string {
  return format(date, 'yyyy年MM月', { locale: zhCN })
}

export function getWeekLabel(date: Date): string {
  const start = startOfWeek(date, { weekStartsOn: 1 })
  const end = endOfWeek(date, { weekStartsOn: 1 })
  if (start.getMonth() === end.getMonth()) {
    return format(start, 'yyyy年MM月第W周', { locale: zhCN })
  }
  return `${format(start, 'MM月dd日', { locale: zhCN })} - ${format(end, 'MM月dd日', { locale: zhCN })}`
}

export function getDayOfWeekLabel(date: Date): string {
  return format(date, 'EEE', { locale: zhCN })
}

export function isDateInRange(date: Date, start: Date, end: Date): boolean {
  return isWithinInterval(date, { start, end })
}

export function getStartOfDay(date: Date): Date {
  const d = new Date(date)
  d.setHours(0, 0, 0, 0)
  return d
}

export function getEndOfDay(date: Date): Date {
  const d = new Date(date)
  d.setHours(23, 59, 59, 999)
  return d
}
