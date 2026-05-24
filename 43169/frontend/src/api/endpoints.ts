import { apiGet, apiPost, apiPut, apiDelete, apiUpload, PageData, ApiResponse } from './index'

export interface UserInfo {
  id: number
  username: string
  phone: string
  email: string
  role: string
  status: string
  verify_status: string
  real_name: string
  avatar: string
  member_level: string
  member_expire: string | null
  profile?: ProfileInfo | null
}

export interface ProfileInfo {
  id: number
  user_id: number
  nickname: string
  gender: string
  birthday: string | null
  age: number
  height: number
  weight: number
  education: string
  occupation: string
  income: string
  city: string
  district: string
  address: string
  latitude: number
  longitude: number
  intro: string
  hobbies: string
  tags: string
  photos: string[]
  min_age: number
  max_age: number
  min_height: number
  max_height: number
  prefer_education: string
  prefer_income: string
  prefer_city: string
}

export interface LoginResponse {
  token: string
  user: UserInfo
}

export interface MatchResultItem {
  user_id: number
  profile: ProfileInfo
  match_score: number
  match_reason: string
  is_favorited: boolean
  is_blocked: boolean
}

export const userApi = {
  register: (data: { username: string; password: string; phone: string; email?: string; code: string }) =>
    apiPost<LoginResponse>('/auth/register', data),

  login: (data: { account: string; password: string }) =>
    apiPost<LoginResponse>('/auth/login', data),

  sendSmsCode: (phone: string) =>
    apiGet<{ phone: string; code: string; expires_in: number }>(`/auth/sms-code?phone=${phone}`),

  getUserInfo: () =>
    apiGet<UserInfo>('/user/info'),

  getUserProfile: (id: number) =>
    apiGet<UserInfo>(`/user/${id}`),

  updateProfile: (data: any) =>
    apiPut('/user/profile', data),

  verify: (data: { real_name: string; id_card: string; id_card_front: string; id_card_back: string; phone: string; sms_code: string }) =>
    apiPost('/user/verify', data),

  uploadAvatar: (file: File) =>
    apiUpload('/user/avatar', file),

  uploadPhotos: (files: File[]) => {
    const formData = new FormData()
    files.forEach((f) => formData.append('files', f))
    return apiPost<{ photos: string[] }>('/user/photos', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
  },

  listUsers: (params: { page?: number; page_size?: number; keyword?: string }) =>
    apiGet<PageData<UserInfo>>('/admin/users', { params }),

  disableUser: (id: number) =>
    apiPut(`/admin/users/${id}/disable`),

  enableUser: (id: number) =>
    apiPut(`/admin/users/${id}/enable`),

  approveVerify: (id: number) =>
    apiPut(`/admin/users/${id}/verify/approve`),

  rejectVerify: (id: number) =>
    apiPut(`/admin/users/${id}/verify/reject`),
}

export const matchApi = {
  smartMatch: (params: { page?: number; page_size?: number }) =>
    apiGet<PageData<MatchResultItem>>('/match/smart', { params }),

  filterMatch: (params: any) =>
    apiGet<PageData<MatchResultItem>>('/match/filter', { params }),

  favorite: (id: number) =>
    apiPost(`/match/${id}/favorite`),

  block: (id: number) =>
    apiPost(`/match/${id}/block`),

  getFavorites: (params: { page?: number; page_size?: number }) =>
    apiGet<PageData<MatchResultItem>>('/match/favorites', { params }),

  getBlocked: (params: { page?: number; page_size?: number }) =>
    apiGet<PageData<MatchResultItem>>('/match/blocked', { params }),
}

export interface DateRecord {
  id: number
  initiator_id: number
  receiver_id: number
  title: string
  location: string
  date_at: string
  duration: number
  status: string
  note: string
  created_at: string
}

export interface DateReview {
  id: number
  date_id: number
  reviewer_id: number
  target_id: number
  rating: number
  content: string
  created_at: string
}

export const dateApi = {
  createInvite: (data: { receiver_id: number; title: string; location?: string; date_at: string; duration?: number; note?: string }) =>
    apiPost<DateRecord>('/dates', data),

  accept: (id: number) =>
    apiPost(`/dates/${id}/accept`),

  reject: (id: number) =>
    apiPost(`/dates/${id}/reject`),

  cancel: (id: number) =>
    apiPost(`/dates/${id}/cancel`),

  complete: (id: number) =>
    apiPost(`/dates/${id}/complete`),

  list: (params: { page?: number; page_size?: number }) =>
    apiGet<PageData<DateRecord>>('/dates', { params }),

  createReview: (data: { date_id: number; target_id: number; rating: number; content?: string }) =>
    apiPost('/dates/reviews', data),

  getReviews: (params: { page?: number; page_size?: number }) =>
    apiGet<PageData<DateReview>>('/dates/reviews', { params }),
}

export interface ChatMessage {
  id: number
  sender_id: number
  receiver_id: number
  type: string
  content: string
  is_read: boolean
  created_at: string
}

export interface ChatSession {
  id: number
  user_a_id: number
  user_b_id: number
  last_message: string
  last_time: string
  unread_a: number
  unread_b: number
}

export const chatApi = {
  sendMessage: (data: { receiver_id: number; type?: string; content: string }) =>
    apiPost<ChatMessage>('/chat/send', data),

  getHistory: (userId: number, params: { page?: number; page_size?: number }) =>
    apiGet<PageData<ChatMessage>>(`/chat/history/${userId}`, { params }),

  getUnreadCount: () =>
    apiGet<{ unread_count: number }>('/chat/unread'),

  getSessions: (params: { page?: number; page_size?: number }) =>
    apiGet<PageData<ChatSession>>('/chat/sessions', { params }),

  markAsRead: (userId: number) =>
    apiPost(`/chat/${userId}/read`),

  uploadFile: (file: File, type: string) => {
    const formData = new FormData()
    formData.append('file', file)
    return apiPost<{ url: string; type: string }>(`/chat/upload?type=${type}`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
  },
}

export interface MemberBenefit {
  id: number
  level: string
  daily_interact: number
  unlimited_chat: boolean
  view_who_liked: boolean
  priority_match: boolean
  advanced_filter: boolean
  video_chat: boolean
  hide_online: boolean
  no_ads: boolean
  matchmaker_assist: boolean
  price_per_month: number
  description: string
}

export interface MemberOrder {
  id: number
  user_id: number
  level: string
  months: number
  amount: number
  status: string
  paid_at: string | null
  expire_at: string | null
  created_at: string
}

export const memberApi = {
  getBenefits: () =>
    apiGet<MemberBenefit[]>('/member/benefits'),

  createOrder: (data: { level: string; months: number }) =>
    apiPost<MemberOrder>('/member/orders', data),

  payOrder: (id: number) =>
    apiPost(`/member/orders/${id}/pay`),

  getOrders: (params: { page?: number; page_size?: number }) =>
    apiGet<PageData<MemberOrder>>('/member/orders', { params }),

  checkInteractLimit: () =>
    apiGet<{ today_count: number }>('/member/interact-limit'),
}

export interface MatchmakerMember {
  id: number
  matchmaker_id: number
  member_id: number
  status: string
  joined_at: string
}

export interface MatchmakerService {
  id: number
  matchmaker_id: number
  member_a_id: number
  member_b_id: number
  date_id: number | null
  service_type: string
  note: string
  status: string
  progress: number
  created_at: string
}

export interface MatchmakerStats {
  id: number
  matchmaker_id: number
  total_members: number
  total_services: number
  total_dates: number
  success_dates: number
  avg_rating: number
  total_rating: number
}

export const matchmakerApi = {
  addMember: (data: { member_id: number }) =>
    apiPost('/matchmaker/members', data),

  removeMember: (id: number) =>
    apiDelete(`/matchmaker/members/${id}`),

  listMembers: (params: { page?: number; page_size?: number }) =>
    apiGet<PageData<MatchmakerMember>>('/matchmaker/members', { params }),

  createService: (data: { member_a_id: number; member_b_id: number; service_type: string; note?: string; date_id?: number | null }) =>
    apiPost('/matchmaker/services', data),

  updateProgress: (id: number, progress: number) =>
    apiPut(`/matchmaker/services/${id}/progress`, { progress }),

  listServices: (params: { page?: number; page_size?: number }) =>
    apiGet<PageData<MatchmakerService>>('/matchmaker/services', { params }),

  getStats: () =>
    apiGet<MatchmakerStats>('/matchmaker/stats'),

  listAll: (params: { page?: number; page_size?: number }) =>
    apiGet<PageData<UserInfo>>('/matchmakers', { params }),
}

export interface PlatformStats {
  total_users: number
  active_today: number
  verified_users: number
  total_matchmakers: number
  total_dates: number
  completed_dates: number
  match_success_rate: number
  total_messages: number
  new_users_today: number
}

export interface DailyStats {
  date: string
  new_users: number
  active_users: number
  dates_created: number
  dates_completed: number
  messages_count: number
}

export interface SystemLog {
  id: number
  user_id: number | null
  module: string
  action: string
  ip: string
  detail: string
  created_at: string
}

export const statsApi = {
  getPlatformStats: () =>
    apiGet<PlatformStats>('/admin/stats/platform'),

  getDailyStats: (params: { start_date?: string; end_date?: string }) =>
    apiGet<DailyStats[]>('/admin/stats/daily', { params }),

  getMatchmakerStats: (params: { start_date?: string; end_date?: string; matchmaker_id?: number }) =>
    apiGet<any[]>('/admin/stats/matchmaker', { params }),

  exportExcel: (params: { start_date?: string; end_date?: string }) =>
    apiGet<Blob>('/admin/stats/export/excel', { params, responseType: 'blob' }),

  exportPDF: (params: { start_date?: string; end_date?: string }) =>
    apiGet<Blob>('/admin/stats/export/pdf', { params, responseType: 'blob' }),

  getSystemLogs: (params: { page?: number; page_size?: number; module?: string }) =>
    apiGet<PageData<SystemLog>>('/admin/logs', { params }),
}

export type { PageData } from './index'
