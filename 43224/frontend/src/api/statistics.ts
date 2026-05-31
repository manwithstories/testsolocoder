import request from '@/utils/request'

export function listReviewTasks(params?: any) {
  return request.get('/review-tasks', { params })
}

export function createReviewTask(data: any) {
  return request.post('/review-tasks', data)
}

export function processReview(id: number, data: any) {
  return request.put(`/review-tasks/${id}`, data)
}

export function batchReview(data: any) {
  return request.post('/review-tasks/batch', data)
}

export function getProjectReviewSummary(projectId: number) {
  return request.get(`/projects/${projectId}/review-summary`)
}

export function listPayments(params?: any) {
  return request.get('/payments', { params })
}

export function getPayment(id: number) {
  return request.get(`/payments/${id}`)
}

export function confirmPayment(id: number) {
  return request.put(`/payments/${id}/confirm`)
}

export function calculateFee(data: any) {
  return request.post('/fee/calculate', data)
}

export function getPaymentStatistics() {
  return request.get('/payments/statistics')
}

export function getProjectStatistics() {
  return request.get('/statistics/projects')
}

export function getTranslatorStatistics() {
  return request.get('/statistics/translators')
}

export function getRevenueTrend(params?: any) {
  return request.get('/statistics/revenue-trend', { params })
}

export function getLanguagePairStatistics() {
  return request.get('/statistics/language-pairs')
}

export function exportExcel() {
  return request.get('/statistics/export/excel', { responseType: 'blob' })
}

export function exportPDF() {
  return request.get('/statistics/export/pdf', { responseType: 'blob' })
}

export function listOperationLogs(params?: any) {
  return request.get('/operation-logs', { params })
}
