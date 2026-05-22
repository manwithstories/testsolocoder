import request from '@/utils/request'
import { CheckIn, PageParams, PageResult, CheckInStatus } from '@/types'

export const getCheckInList = (params: PageParams & { status?: CheckInStatus; startDate?: string; endDate?: string }) => {
  return request.get<PageResult<CheckIn>>('/checkins', { params })
}

export const getCheckInById = (id: number) => {
  return request.get<CheckIn>(`/checkins/${id}`)
}

export const createCheckIn = (data: Omit<CheckIn, 'id' | 'checkInNo' | 'createdAt' | 'updatedAt'>) => {
  return request.post<CheckIn>('/checkins', data)
}

export const updateCheckIn = (id: number, data: Partial<CheckIn>) => {
  return request.put<CheckIn>(`/checkins/${id}`, data)
}

export const checkOut = (id: number, data?: { actualCheckOutTime?: string; extraCharges?: number }) => {
  return request.post(`/checkins/${id}/checkout`, data)
}

export const extendStay = (id: number, data: { newCheckOutTime: string; additionalAmount: number }) => {
  return request.post(`/checkins/${id}/extend`, data)
}

export const getTodayCheckIns = () => {
  return request.get<CheckIn[]>('/checkins/today')
}

export const getTodayCheckOuts = () => {
  return request.get<CheckIn[]>('/checkins/today-checkouts')
}

export const getCheckedInRooms = () => {
  return request.get<CheckIn[]>('/checkins/checked-in')
}
