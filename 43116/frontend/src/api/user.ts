import request from './request'
import type {
  User, LoginRequest, RegisterRequest, LoginResponse, UserInfo,
  PageResult, ApiResponse
} from '@/types'

export const userApi = {
  register: (data: RegisterRequest) =>
    request.post<any, ApiResponse<User>>('/register', data),

  login: (data: LoginRequest) =>
    request.post<any, ApiResponse<LoginResponse>>('/login', data),

  refreshToken: (refreshToken: string) =>
    request.post<any, ApiResponse<LoginResponse>>('/refresh-token', { refresh_token: refreshToken }),

  getProfile: () =>
    request.get<any, ApiResponse<User>>('/profile'),

  updateProfile: (data: Record<string, any>) =>
    request.put<any, ApiResponse<null>>('/profile', data),

  changePassword: (oldPassword: string, newPassword: string) =>
    request.put<any, ApiResponse<null>>('/password', { old_password: oldPassword, new_password: newPassword }),

  uploadLicense: (file: File, field: string) => {
    const formData = new FormData()
    formData.append('file', file)
    return request.post<any, ApiResponse<{ url: string }>>('/upload-license', formData, {
      params: { field }
    })
  },

  uploadAvatar: (file: File) => {
    const formData = new FormData()
    formData.append('file', file)
    return request.post<any, ApiResponse<{ url: string }>>('/upload-avatar', formData)
  },

  getUsers: (params: { page?: number; page_size?: number; keyword?: string }) =>
    request.get<any, ApiResponse<PageResult<User>>>('/users', { params }),

  getUserById: (id: number) =>
    request.get<any, ApiResponse<User>>(`/users/${id}`),

  updateAuthStatus: (id: number, status: string) =>
    request.put<any, ApiResponse<null>>(`/users/${id}/auth-status`, { status }),

  updateUserStatus: (id: number, status: string) =>
    request.put<any, ApiResponse<null>>(`/users/${id}/status`, { status }),

  deleteUser: (id: number) =>
    request.delete<any, ApiResponse<null>>(`/users/${id}`)
}
