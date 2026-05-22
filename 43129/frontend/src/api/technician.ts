import request from '@/utils/request'
import type { ApiResponse, Technician, PageResult, PageParams } from '@/types'

export const getTechnicians = (params: PageParams & { status?: number }) => {
  return request.get<ApiResponse<PageResult<Technician>>>('/technicians', { params })
}

export const getAllTechnicians = () => {
  return request.get<ApiResponse<Technician[]>>('/technicians/all')
}

export const getTechnician = (id: number) => {
  return request.get<ApiResponse<Technician>>(`/technicians/${id}`)
}

export const getMyTechnician = () => {
  return request.get<ApiResponse<Technician>>('/technicians/my')
}

export const createTechnician = (data: any) => {
  return request.post<ApiResponse<Technician>>('/technicians', data)
}

export const updateTechnician = (id: number, data: any) => {
  return request.put<ApiResponse<Technician>>(`/technicians/${id}`, data)
}

export const addTechnicianLeave = (id: number, data: { date: string; reason?: string }) => {
  return request.post<ApiResponse<null>>(`/technicians/${id}/leave`, data)
}

export const getTechnicianLeaves = (id: number, params?: { month?: number }) => {
  return request.get<ApiResponse<any[]>>(`/technicians/${id}/leaves`, { params })
}

export const getTechnicianSchedule = (id: number, params: { date: string }) => {
  return request.get<ApiResponse<any[]>>(`/technicians/${id}/schedule`, { params })
}
