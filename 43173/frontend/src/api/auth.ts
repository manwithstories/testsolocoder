import request, { PageData, PageParams } from '@/utils/request'

export interface LoginParams {
  account: string
  password: string
}

export interface RegisterParams {
  username: string
  email: string
  password: string
  nickname?: string
  role?: string
  phone?: string
}

export interface UserInfo {
  id: number
  username: string
  nickname: string
  avatar: string
  role: string
  email?: string
  phone?: string
  bio?: string
  artist_info?: ArtistInfo
}

export interface ArtistInfo {
  id: number
  user_id: number
  artist_name: string
  genre: string
  label: string
  website?: string
  is_verified: boolean
  total_plays: number
  total_followers: number
  total_works: number
  balance: number
  frozen_balance: number
}

export const authApi = {
  login: (data: LoginParams) => request.post('/auth/login', data),
  register: (data: RegisterParams) => request.post('/auth/register', data),
  getProfile: () => request.get<UserInfo>('/auth/profile'),
  updateProfile: (data: Partial<UserInfo>) => request.put('/auth/profile', data),
  getArtistInfo: () => request.get<ArtistInfo>('/auth/artist-info'),
  updateArtistInfo: (data: Partial<ArtistInfo>) => request.put('/auth/artist-info', data),
  getBalance: () => request.get<{ balance: number; frozen_balance: number }>('/auth/balance')
}

export const userApi = {
  getById: (id: number) => request.get<UserInfo>(`/users/${id}`),
  list: (params: PageParams & { keyword?: string; role?: string }) => request.get<PageData<UserInfo>>('/users', params),
  updateStatus: (id: number, data: { status: number }) => request.put(`/users/${id}/status`, data),
  updateRole: (id: number, data: { role: string }) => request.put(`/users/${id}/role`, data),
  verifyArtist: (id: number) => request.post(`/users/${id}/verify`),
  getOperationLogs: (params: PageParams & { keyword?: string; start_date?: string; end_date?: string }) => 
    request.get<PageData<any>>('/admin/operation-logs', params)
}
