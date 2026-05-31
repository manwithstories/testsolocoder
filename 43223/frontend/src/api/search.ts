import { request } from './client'
import { Product, PaginatedData } from '@/types'

export const searchApi = {
  searchProducts: (params?: {
    q?: string
    origin?: string
    roast_level?: string
    process_method?: string
    min_price?: string
    max_price?: string
    min_score?: string
    max_score?: string
    roaster_id?: string
    sort_by?: string
    sort_order?: string
    page?: number
    page_size?: number
  }) =>
    request.get<{
      items: Product[]
      total: number
      page: number
      pageSize: number
      filters: {
        origins: string[]
        roast_levels: string[]
        process_methods: string[]
        price_range: { min: number; max: number }
        score_range: { min: number; max: number }
      }
    }>('/search', { params }),

  suggest: (q: string) =>
    request.get<{ id: number; name: string; type: string; image?: string }[]>('/search/suggest', {
      params: { q },
    }),

  advancedSearch: (data: {
    keywords?: string[]
    filters?: Record<string, string>
    price_range?: [number, number]
    score_range?: [number, number]
    sort_by?: string
    sort_order?: string
    page?: number
    page_size?: number
  }) =>
    request.post<PaginatedData<Product>>('/search/advanced', data),

  getHistory: () =>
    request.get('/search/history'),
}

export const statsApi = {
  getSalesStats: () =>
    request.get('/stats/sales'),

  getSalesTrend: (params?: { days?: number }) =>
    request.get('/stats/sales-trend', { params }),

  getOriginDistribution: () =>
    request.get('/stats/origins'),

  getUserActivity: (params?: { days?: number }) =>
    request.get('/stats/user-activity', { params }),

  getTopProducts: (params?: { limit?: number }) =>
    request.get('/stats/top-products', { params }),

  exportExcel: (type: string, params?: Record<string, string>) =>
    request.get('/stats/export/excel', { params: { type, ...params }, responseType: 'blob' }),

  exportPDF: (type: string) =>
    request.get('/stats/export/pdf', { params: { type }, responseType: 'blob' }),
}
