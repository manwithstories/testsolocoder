import request from './request'
import type { Booking, Waitlist, CheckIn, ApiResponse, PaginatedResponse } from '@/types'

export const bookingApi = {
  book: (data: { member_id: number; schedule_id: number }) => {
    return request.post<any, ApiResponse<Booking>>('/bookings', data)
  },

  getList: (params: { page?: number; page_size?: number; member_id?: number; schedule_id?: number }) => {
    return request.get<any, PaginatedResponse<Booking[]>>('/bookings', { params })
  },

  getById: (id: number) => {
    return request.get<any, ApiResponse<Booking>>(`/bookings/${id}`)
  },

  cancel: (id: number) => {
    return request.delete<any, ApiResponse>(`/bookings/${id}`)
  },

  getByMember: (memberId: number, params?: { page?: number; page_size?: number }) => {
    return request.get<any, PaginatedResponse<Booking[]>>(`/bookings/member/${memberId}`, { params })
  },

  getBySchedule: (scheduleId: number, params?: { page?: number; page_size?: number }) => {
    return request.get<any, PaginatedResponse<Booking[]>>(`/bookings/schedule/${scheduleId}`, { params })
  },

  addToWaitlist: (data: { member_id: number; schedule_id: number }) => {
    return request.post<any, ApiResponse<Waitlist>>('/bookings/waitlist', data)
  },

  removeFromWaitlist: (id: number) => {
    return request.delete<any, ApiResponse>(`/bookings/waitlist/${id}`)
  }
}

export const checkInApi = {
  checkIn: (data: { member_id: number; schedule_id?: number }) => {
    return request.post<any, ApiResponse<CheckIn>>('/check-ins', data)
  },

  getList: (params: { page?: number; page_size?: number; member_id?: number; date?: string }) => {
    return request.get<any, PaginatedResponse<CheckIn[]>>('/check-ins', { params })
  },

  getById: (id: number) => {
    return request.get<any, ApiResponse<CheckIn>>(`/check-ins/${id}`)
  },

  getByMember: (memberId: number, params?: { page?: number; page_size?: number }) => {
    return request.get<any, PaginatedResponse<CheckIn[]>>(`/check-ins/member/${memberId}`, { params })
  },

  getByDate: (date: string) => {
    return request.get<any, ApiResponse<CheckIn[]>>(`/check-ins/date/${date}`)
  }
}

export const statsApi = {
  getDashboard: () => {
    return request.get<any, ApiResponse<any>>('/stats/dashboard')
  },

  getMemberStats: (params?: { start_date?: string; end_date?: string }) => {
    return request.get<any, ApiResponse<any>>('/stats/members', { params })
  },

  getCourseStats: (params?: { start_date?: string; end_date?: string }) => {
    return request.get<any, ApiResponse<any>>('/stats/courses', { params })
  },

  getCoachStats: (params?: { start_date?: string; end_date?: string }) => {
    return request.get<any, ApiResponse<any>>('/stats/coaches', { params })
  }
}
