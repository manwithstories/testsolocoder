import { request } from '@/utils/request'
import { ApiResponse, PaginatedData, User, LoginRequest, RegisterRequest } from '@/types'

export const authApi = {
  login: (data: LoginRequest) => 
    request.post<ApiResponse>('/auth/login', data),
  
  register: (data: RegisterRequest) => 
    request.post<ApiResponse>('/auth/register', data)
}

export const userApi = {
  getProfile: () => 
    request.get<ApiResponse<User>>('/users/profile'),
  
  updateProfile: (data: Partial<User>) => 
    request.put<ApiResponse>('/users/profile', data),
  
  changePassword: (data: { old_password: string; new_password: string }) => 
    request.put<ApiResponse>('/users/password', data),
  
  getUsers: (params?: { search?: string; role?: string; page?: number; page_size?: number }) => 
    request.get<ApiResponse<PaginatedData<User>>>('/users', { params }),
  
  updateStatus: (id: number, status: string) => 
    request.put<ApiResponse>(`/users/${id}/status`, { status }),
  
  deleteUser: (id: number) => 
    request.delete<ApiResponse>(`/users/${id}`)
}
