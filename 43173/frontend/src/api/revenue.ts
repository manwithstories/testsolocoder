import request, { PageData, PageParams } from '@/utils/request'

export interface RevenueRecord {
  id: number
  user_id: number
  artist_id: number
  work_id?: number
  order_id?: number
  type: string
  amount: number
  play_count: number
  rate: number
  status: number
  period?: string
  settled_at?: string
  created_at: string
  updated_at: string
}

export interface WithdrawRequest {
  id: number
  user_id: number
  artist_id: number
  amount: number
  fee: number
  actual_amount: number
  method: string
  account: string
  account_name?: string
  bank_name?: string
  status: number
  remark?: string
  approved_at?: string
  approved_by?: number
  reject_reason?: string
  paid_at?: string
  transaction_no?: string
  created_at: string
  updated_at: string
}

export interface DailyStats {
  id: number
  date: string
  user_id: number
  artist_id: number
  new_followers: number
  new_plays: number
  new_likes: number
  new_shares: number
  new_comments: number
  revenue: number
  created_at: string
  updated_at: string
}

export interface OperationLog {
  id: number
  user_id?: number
  username?: string
  module: string
  operation: string
  method: string
  path: string
  ip: string
  user_agent: string
  params?: string
  result?: string
  status: number
  duration: number
  error_msg?: string
  created_at: string
}

export interface WithdrawParams {
  amount: number
  method: string
  account: string
  account_name?: string
  bank_name?: string
}

export const revenueApi = {
  getRecords: (params: PageParams & { start_date?: string; end_date?: string }) => 
    request.get<PageData<RevenueRecord>>('/revenue/records', params),
  getArtistRecords: (artistId: number, params: PageParams & { start_date?: string; end_date?: string }) => 
    request.get<PageData<RevenueRecord>>(`/revenue/records/artist/${artistId}`, params),
  getTotal: (params?: { start_date?: string; end_date?: string }) => 
    request.get<{ total_revenue: number }>('/revenue/total', params),
  getSummary: () => request.get('/revenue/summary'),
  
  requestWithdraw: (data: WithdrawParams) => request.post<WithdrawRequest>('/withdraw', data),
  getWithdrawList: (params: PageParams & { status?: number }) => 
    request.get<PageData<WithdrawRequest>>('/withdraw', params),
  getWithdrawStatusList: () => request.get('/withdraw/status-list'),
  
  getSubscriptions: (params: PageParams) => request.get('/subscriptions', params),
  getArtistSubscribers: (artistId: number, params: PageParams) => 
    request.get(`/subscriptions/artist/${artistId}`, params),
  
  getDailyStats: (params: { start_date: string; end_date: string }) => 
    request.get<DailyStats[]>('/stats/daily', params),
  getArtistDailyStats: (artistId: number, params: { start_date: string; end_date: string }) => 
    request.get<DailyStats[]>(`/stats/daily/artist/${artistId}`, params),
  getArtistStats: (artistId: number) => request.get(`/stats/artist/${artistId}`),
  
  getAdminWithdrawList: (params: PageParams & { status?: number }) => 
    request.get<PageData<WithdrawRequest>>('/admin/withdraw', params),
  approveWithdraw: (id: number) => request.put(`/admin/withdraw/${id}/approve`),
  rejectWithdraw: (id: number, data: { reason: string }) => 
    request.put(`/admin/withdraw/${id}/reject`, data),
  markWithdrawPaid: (id: number, data?: { transaction_no?: string }) => 
    request.put(`/admin/withdraw/${id}/paid`, data || {}),
  settleRevenue: (period: string) => request.post('/admin/revenue/settle', { period }),
  
  getOperationLogs: (params: PageParams & { user_id?: number; module?: string; keyword?: string }) => 
    request.get<PageData<OperationLog>>('/admin/operation-logs', params),
  
  exportRevenue: (params?: { start_date?: string; end_date?: string }) => 
    request.get('/export/revenue', params, { responseType: 'blob' }),
  exportWithdraw: (params?: { status?: number }) => 
    request.get('/export/withdraw', params, { responseType: 'blob' })
}
