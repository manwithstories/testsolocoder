import { http } from './request'
import {
  ProfessionalStats,
  AdminStats,
} from '@/types'

export const statisticsApi = {
  getProfessionalStats: (params?: { start_date?: string; end_date?: string }) =>
    http.get<ProfessionalStats>('/statistics/professional', { params }),
  getAdminStats: (params?: { start_date?: string; end_date?: string }) =>
    http.get<AdminStats>('/statistics/admin', { params }),
  exportAppointments: (params?: { start_date?: string; end_date?: string; status?: string }) =>
    http.get<Blob>('/statistics/professional/export/appointments', {
      params,
      responseType: 'blob',
    }),
  exportRevenue: (params?: { start_date?: string; end_date?: string }) =>
    http.get<Blob>('/statistics/professional/export/revenue', {
      params,
      responseType: 'blob',
    }),
}

export { ProfessionalStats, AdminStats }
