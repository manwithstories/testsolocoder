import { request } from './index'
import type { DashboardStats, ApiResponse } from '@/types'

export const statisticApi = {
  getDashboardStats(days: number = 30) {
    return request.get<DashboardStats>('/statistics/dashboard', { params: { days } })
  },

  invalidateCache() {
    return request.post('/statistics/invalidate-cache')
  }
}

export default statisticApi
