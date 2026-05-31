import request from './index'
import type { ApiResponse } from './index'

export interface StaffProfile {
  id: number
  user_id: number
  id_card?: string
  cert_files?: string
  health_files?: string
  cert_verified: boolean
  intro?: string
}

export interface StaffDetail {
  id: number
  username: string
  real_name: string
  role: string
  rating: number
  level: number
  suspended: boolean
  skills?: string
  company_id?: number
  profile: StaffProfile
  reviews: any[]
}

export function listStaff(params?: Record<string, string | number>) {
  return request.get<ApiResponse<StaffDetail[]>>('/staff', { params })
}

export function getStaffDetail(id: number) {
  return request.get<ApiResponse<StaffDetail>>(`/staff/${id}`)
}

export function listStaffReviews(id: number) {
  return request.get<ApiResponse<any[]>>(`/staff/${id}/reviews`)
}

export function updateMe(payload: Record<string, string>) {
  return request.put<ApiResponse<string>>('/me', payload)
}

export function getMe() {
  return request.get<ApiResponse<any>>('/me')
}

export function updateStaffCert(payload: {
  cert_files?: string
  health_files?: string
  id_card?: string
  intro?: string
}) {
  return request.put<ApiResponse<string>>('/staff-area/certs', payload)
}
