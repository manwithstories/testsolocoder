import { apiGet, apiPut } from './client'
import { RescueStation, RescueStats, RescueListQuery, ReviewRescueRequest, ApiResponse } from '../types'

export const listRescueStations = (params?: RescueListQuery): Promise<ApiResponse<any>> => {
  return apiGet('/rescue', params)
}

export const getRescueStation = (id: number): Promise<ApiResponse<RescueStation>> => {
  return apiGet<RescueStation>(`/rescue/${id}`)
}

export const reviewRescueStation = (id: number, data: ReviewRescueRequest): Promise<ApiResponse<RescueStation>> => {
  return apiPut<RescueStation>(`/rescue/${id}/review`, data)
}

export const getRescueStats = (): Promise<ApiResponse<RescueStats>> => {
  return apiGet<RescueStats>('/rescue/me/stats')
}

export const getRescueStatsById = (id: number): Promise<ApiResponse<RescueStats>> => {
  return apiGet<RescueStats>(`/rescue/${id}/stats`)
}

export const getAllRescuesStats = (): Promise<ApiResponse<any>> => {
  return apiGet('/rescue/stats/all')
}
