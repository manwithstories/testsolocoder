import { request } from '@/utils/request'
import { ApiResponse, DashboardStats } from '@/types'

export const dashboardApi = {
  getStats: (params?: { wedding_id?: number }) => 
    request.get<ApiResponse<DashboardStats>>('/dashboard/stats', { params }),
  
  getBudgetChart: (params?: { wedding_id?: number }) => 
    request.get<ApiResponse<any[]>>('/dashboard/budget-chart', { params }),
  
  getTaskProgress: (params?: { wedding_id?: number }) => 
    request.get<ApiResponse<any[]>>('/dashboard/task-progress', { params }),
  
  getUpcomingTasks: (params?: { wedding_id?: number }) => 
    request.get<ApiResponse<any[]>>('/dashboard/upcoming-tasks', { params }),
  
  getVendorStats: () => 
    request.get<ApiResponse<any[]>>('/dashboard/vendor-stats'),
  
  exportReport: (params?: { wedding_id?: number }) => 
    request.get('/dashboard/export', { params, responseType: 'blob' })
}

export const notificationApi = {
  getList: () => 
    request.get<ApiResponse<{ notifications: any[]; unread_count: number }>>('/notifications'),
  
  markAsRead: (id: number) => 
    request.put<ApiResponse>(`/notifications/${id}/read`),
  
  markAllAsRead: () => 
    request.put<ApiResponse>('/notifications/read-all')
}

export const logApi = {
  getLogs: (params?: { module?: string; action?: string; page?: number; page_size?: number }) => 
    request.get<ApiResponse<any>>('/logs', { params })
}
