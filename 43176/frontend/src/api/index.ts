import request from '@/utils/request'
import type { ApiResponse, LoginResponse, User, PaginatedData, PaginationParams } from '@/types'

export const authApi = {
  register: (data: { phone: string; password: string; nickname: string; role: string }) => {
    return request.post<any, ApiResponse>('/auth/register', data)
  },

  login: (data: { phone: string; password: string }) => {
    return request.post<any, ApiResponse<LoginResponse>>('/auth/login', data)
  },

  refreshToken: (refreshToken: string) => {
    return request.post<any, ApiResponse<LoginResponse>>('/auth/refresh', { refresh_token: refreshToken })
  }
}

export const userApi = {
  getProfile: () => {
    return request.get<any, ApiResponse<User>>('/user/profile')
  },

  updateProfile: (data: { nickname?: string; avatar?: string }) => {
    return request.put<any, ApiResponse>('/user/profile', data)
  },

  submitVerification: (data: { id_card_no: string; id_card_name: string; id_card_front: string; id_card_back: string }) => {
    return request.post<any, ApiResponse>('/user/verify', data)
  }
}

export const taskApi = {
  create: (data: any) => {
    return request.post<any, ApiResponse>('/tasks', data)
  },

  list: (params?: { status?: string; type?: string } & PaginationParams) => {
    return request.get<any, ApiResponse<PaginatedData<any>>>('/tasks', { params })
  },

  get: (id: number) => {
    return request.get<any, ApiResponse>(`/tasks/${id}`)
  },

  update: (id: number, data: any) => {
    return request.put<any, ApiResponse>(`/tasks/${id}`, data)
  },

  delete: (id: number) => {
    return request.delete<any, ApiResponse>(`/tasks/${id}`)
  },

  accept: (id: number) => {
    return request.post<any, ApiResponse>(`/tasks/${id}/accept`)
  },

  complete: (id: number, data?: { proof_images?: string[] }) => {
    return request.post<any, ApiResponse>(`/tasks/${id}/complete`, data || {})
  },

  cancel: (id: number, data?: { reason?: string }) => {
    return request.post<any, ApiResponse>(`/tasks/${id}/cancel`, data || {})
  },

  getNearby: (params: { lat: number; lng: number; radius?: number }) => {
    return request.get<any, ApiResponse>('/tasks/nearby', { params })
  }
}

export const orderApi = {
  list: (params?: { status?: string } & PaginationParams) => {
    return request.get<any, ApiResponse<PaginatedData<any>>>('/orders', { params })
  },

  get: (id: number) => {
    return request.get<any, ApiResponse>(`/orders/${id}`)
  },

  track: (id: number, data: { latitude: number; longitude: number; address?: string; message?: string; event_type?: string }) => {
    return request.post<any, ApiResponse>(`/orders/${id}/track`, data)
  }
}

export const paymentApi = {
  deposit: (data: { amount: number; payment_method: string }) => {
    return request.post<any, ApiResponse>('/payments/deposit', data)
  },

  withdraw: (data: { amount: number; account_type: string; account_no: string; account_name: string }) => {
    return request.post<any, ApiResponse>('/payments/withdraw', data)
  },

  refund: (data: { order_id: number; reason: string; amount?: number }) => {
    return request.post<any, ApiResponse>('/payments/refund', data)
  },

  history: (params?: { type?: string } & PaginationParams) => {
    return request.get<any, ApiResponse<PaginatedData<any>>>('/payments/history', { params })
  }
}

export const reviewApi = {
  create: (data: { order_id: number; review_type: string; rating: number; content?: string; tags?: string }) => {
    return request.post<any, ApiResponse>('/reviews', data)
  },

  list: (params?: { user_id?: number; type?: string } & PaginationParams) => {
    return request.get<any, ApiResponse>('/reviews', { params })
  }
}

export const courierApi = {
  apply: (data: { id_card_no: string; id_card_name: string; id_card_front: string; id_card_back: string; experience?: string; vehicle?: string }) => {
    return request.post<any, ApiResponse>('/courier/apply', data)
  },

  getMyTasks: (params?: { status?: string } & PaginationParams) => {
    return request.get<any, ApiResponse<PaginatedData<any>>>('/courier/tasks', { params })
  }
}

export const adminApi = {
  listUsers: (params?: { role?: string; status?: string } & PaginationParams) => {
    return request.get<any, ApiResponse<PaginatedData<User>>>('/admin/users', { params })
  },

  freezeUser: (id: number, data: { reason: string }) => {
    return request.put<any, ApiResponse>(`/admin/users/${id}/freeze`, data)
  },

  unfreezeUser: (id: number) => {
    return request.put<any, ApiResponse>(`/admin/users/${id}/unfreeze`)
  },

  approveCourier: (id: number) => {
    return request.put<any, ApiResponse>(`/admin/couriers/${id}/approve`)
  },

  rejectCourier: (id: number, data: { reason: string }) => {
    return request.put<any, ApiResponse>(`/admin/couriers/${id}/reject`, data)
  }
}
