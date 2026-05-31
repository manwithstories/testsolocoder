import request from './index'
import type { ApiResponse } from './index'

export function statsOverview() {
  return request.get<ApiResponse<any>>('/admin/stats/overview')
}

export function statsRevenue() {
  return request.get<ApiResponse<any[]>>('/admin/stats/revenue')
}

export function statsCategory() {
  return request.get<ApiResponse<any[]>>('/admin/stats/category')
}

export function statsStaff() {
  return request.get<ApiResponse<any[]>>('/admin/stats/staff')
}

export function companyDashboard() {
  return request.get<ApiResponse<any>>('/company/dashboard')
}
