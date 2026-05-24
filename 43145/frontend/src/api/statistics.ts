import { get, put, del } from '@/utils/request'
import type { Statistics } from '@/types'

export const statisticsApi = {
  getStatistics(params: { survey_id: number; start_date?: string; end_date?: string; channel?: string }) {
    return get<Statistics>('/statistics', { params })
  },

  crossAnalysis(params: { survey_id: number; row_question_id: number; col_question_id: number }) {
    return get('/statistics/cross-analysis', { params })
  }
}

export const exportApi = {
  exportExcel(surveyId: number) {
    return get(`/export/${surveyId}/excel`, { responseType: 'blob' } as any)
  },

  exportPDF(surveyId: number) {
    return get(`/export/${surveyId}/pdf`)
  },

  exportCharts(surveyId: number) {
    return get(`/export/${surveyId}/charts`)
  }
}

export const userApi = {
  list(params: { page?: number; page_size?: number; keyword?: string; status?: number }) {
    return get('/users', { params })
  },

  getById(id: number) {
    return get(`/users/${id}`)
  },

  updateRole(id: number, roleId: number) {
    return put(`/users/${id}/role`, { role_id: roleId })
  },

  updateStatus(id: number, status: number) {
    return put(`/users/${id}/status`, { status })
  },

  remove(id: number) {
    return del(`/users/${id}`)
  }
}
