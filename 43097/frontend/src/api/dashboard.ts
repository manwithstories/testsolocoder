import request from '@/utils/request'
import { DashboardStats } from '@/types'

export const getDashboardStats = () => {
  return request.get<DashboardStats>('/dashboard/stats')
}

export const getRecentBookings = (limit: number = 10) => {
  return request.get('/dashboard/recent-bookings', { params: { limit } })
}

export const getRecentCheckIns = (limit: number = 10) => {
  return request.get('/dashboard/recent-checkins', { params: { limit } })
}

export const getRoomStatusOverview = () => {
  return request.get('/dashboard/room-status')
}

export const getRevenueTrend = (days: number = 7) => {
  return request.get('/dashboard/revenue-trend', { params: { days } })
}

export const getBookingTrend = (days: number = 7) => {
  return request.get('/dashboard/booking-trend', { params: { days } })
}

export const getTopRoomTypes = (limit: number = 5) => {
  return request.get('/dashboard/top-room-types', { params: { limit } })
}

export const getMemberGrowth = (months: number = 6) => {
  return request.get('/dashboard/member-growth', { params: { months } })
}
