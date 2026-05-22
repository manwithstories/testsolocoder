import request from './request'

export const statisticsApi = {
  getDashboard: (params?: { start_date?: string; end_date?: string }) => {
    return request.get<any, {
      total_orders: number
      completed_orders: number
      total_revenue: number
      total_platform_fee: number
      new_customers: number
      active_providers: number
      order_trend: Array<{ date: string; order_count: number; revenue: number }>
      service_type_distribution: Array<{ category_name: string; order_count: number; percentage: number }>
      top_providers: Array<{
        id: number
        nickname: string
        rating: number
        order_count: number
        total_income: number
      }>
    }>('/admin/dashboard', { params })
  },

  getOrderStatistics: (params?: {
    start_date?: string
    end_date?: string
    status?: string
    service_type?: string
  }) => {
    return request.get<any, {
      total_orders: number
      total_amount: number
      status_distribution: Array<{ status: string; order_count: number }>
      hourly_distribution: Array<{ hour: number; order_count: number }>
    }>('/admin/statistics/orders', { params })
  },

  getUserStatistics: (params?: { start_date?: string; end_date?: string }) => {
    return request.get<any, {
      total_customers: number
      total_providers: number
      new_customers: number
      new_providers: number
      active_customers: number
      active_providers: number
      customer_growth: Array<{ date: string; new_users: number }>
      rating_distribution: Array<{ rating: number; provider_count: number }>
    }>('/admin/statistics/users', { params })
  },

  export: (params?: { type?: string; start_date?: string; end_date?: string }) => {
    return request.get<any, any>('/admin/statistics/export', {
      params,
      responseType: 'blob',
    } as any)
  },
}
