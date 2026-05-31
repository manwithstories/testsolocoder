import { request } from '@/utils/request'

export const getCompanyStatistics = (params: { year?: number; start_date?: string; end_date?: string }): Promise<any> => {
  return request.get('/statistics/company', { params })
}

export const getDepartmentStatistics = (departmentId: number, year?: number): Promise<any> => {
  return request.get(`/statistics/department/${departmentId}`, { params: { year } })
}

export const getAgeDistribution = (): Promise<any> => {
  return request.get('/statistics/age-distribution')
}

export const getGenderDistribution = (): Promise<any> => {
  return request.get('/statistics/gender-distribution')
}

export const getAbnormalDistribution = (year?: number): Promise<any> => {
  return request.get('/statistics/abnormal-distribution', { params: { year } })
}

export const getAgencyRating = (): Promise<any> => {
  return request.get('/statistics/agency-rating')
}

export const getPackageRanking = (limit?: number): Promise<any> => {
  return request.get('/statistics/package-ranking', { params: { limit } })
}

export const exportStatistics = (params: { year?: number }): Promise<any> => {
  return request.get('/statistics/export', { 
    params,
    responseType: 'blob'
  })
}
