import request from './request'
import type { Course, CourseSchedule, Coach, ApiResponse, PaginatedResponse } from '@/types'

export const courseApi = {
  create: (data: {
    name: string
    description?: string
    coach_id: number
    capacity: number
    duration: number
    type: string
    weekdays?: string
    start_date: string
    end_date?: string
    start_time: string
    location?: string
  }) => {
    return request.post<any, ApiResponse<Course>>('/courses', data)
  },

  getList: (params: { page?: number; page_size?: number; keyword?: string; coach_id?: number; type?: string }) => {
    return request.get<any, PaginatedResponse<Course[]>>('/courses', { params })
  },

  getById: (id: number) => {
    return request.get<any, ApiResponse<Course>>(`/courses/${id}`)
  },

  update: (id: number, data: Partial<Course>) => {
    return request.put<any, ApiResponse>(`/courses/${id}`, data)
  },

  delete: (id: number) => {
    return request.delete<any, ApiResponse>(`/courses/${id}`)
  },

  updateStatus: (id: number, status: number) => {
    return request.patch<any, ApiResponse>(`/courses/${id}/status`, { status })
  },

  generateSchedules: (id: number) => {
    return request.post<any, ApiResponse>(`/courses/${id}/generate-schedules`)
  }
}

export const scheduleApi = {
  getList: (params: { page?: number; page_size?: number; course_id?: number; start_date?: string; end_date?: string }) => {
    return request.get<any, PaginatedResponse<CourseSchedule[]>>('/courses/schedules', { params })
  },

  getAvailable: () => {
    return request.get<any, ApiResponse<CourseSchedule[]>>('/courses/schedules/available')
  },

  getById: (id: number) => {
    return request.get<any, ApiResponse<CourseSchedule>>(`/courses/schedules/${id}`)
  },

  updateStatus: (id: number, status: number) => {
    return request.patch<any, ApiResponse>(`/courses/schedules/${id}/status`, { status })
  }
}

export const coachApi = {
  create: (data: { name: string; phone: string; specialty?: string; description?: string; photo?: string }) => {
    return request.post<any, ApiResponse<Coach>>('/coaches', data)
  },

  getList: (params: { page?: number; page_size?: number; keyword?: string }) => {
    return request.get<any, PaginatedResponse<Coach[]>>('/coaches', { params })
  },

  getById: (id: number) => {
    return request.get<any, ApiResponse<Coach>>(`/coaches/${id}`)
  },

  update: (id: number, data: Partial<Coach>) => {
    return request.put<any, ApiResponse>(`/coaches/${id}`, data)
  },

  delete: (id: number) => {
    return request.delete<any, ApiResponse>(`/coaches/${id}`)
  }
}
