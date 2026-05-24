import { get } from '@/utils/request'
import type { ApiResponse, DailyStatistics } from '@/types'
import axios from 'axios'
import { useUserStore } from '@/stores/user'

export function getDailyStats(date?: string) {
  return get<ApiResponse<DailyStatistics>>('/stats/daily', { params: { date } })
}

export function getDateRangeStats(startDate: string, endDate: string) {
  return get<ApiResponse<DailyStatistics[]>>('/stats/range', { params: { start_date: startDate, end_date: endDate } })
}

export function getApplicationStats(startDate?: string, endDate?: string) {
  return get<ApiResponse<any>>('/stats/applications', { params: { start_date: startDate, end_date: endDate } })
}

export function getJobStats(startDate?: string, endDate?: string) {
  return get<ApiResponse<any[]>>('/company/stats/jobs', { params: { start_date: startDate, end_date: endDate } })
}

export function getRecruitmentCycleStats() {
  return get<ApiResponse<any[]>>('/company/stats/recruitment-cycle')
}

function getAuthHeaders() {
  const userStore = useUserStore()
  return userStore.token ? { Authorization: `Bearer ${userStore.token}` } : {}
}

export async function exportJobStats(startDate?: string, endDate?: string) {
  const params: Record<string, string> = {}
  if (startDate) params.start_date = startDate
  if (endDate) params.end_date = endDate
  const response = await axios.get('/api/v1/company/stats/export/jobs', {
    params,
    responseType: 'blob',
    headers: getAuthHeaders()
  })
  downloadBlob(response.data, '职位投递统计.xlsx')
}

export async function exportApplicationStats(startDate?: string, endDate?: string) {
  const params: Record<string, string> = {}
  if (startDate) params.start_date = startDate
  if (endDate) params.end_date = endDate
  const response = await axios.get('/api/v1/company/stats/export/applications', {
    params,
    responseType: 'blob',
    headers: getAuthHeaders()
  })
  downloadBlob(response.data, '投递数据统计.xlsx')
}

function downloadBlob(blob: Blob, filename: string) {
  const url = window.URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = filename
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  window.URL.revokeObjectURL(url)
}
