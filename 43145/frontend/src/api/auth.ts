import { post, get, put } from '@/utils/request'
import type { User } from '@/types'

export const authApi = {
  login(data: { email: string; password: string }) {
    return post<{ token: string; user_id: number; email: string; nickname: string; role: string }>('/auth/login', data)
  },

  register(data: { email: string; password: string; nickname: string }) {
    return post<{ token: string; user_id: number; email: string; nickname: string; role: string }>('/auth/register', data)
  },

  getProfile() {
    return get<User>('/profile')
  },

  updateProfile(data: { nickname?: string; avatar?: string }) {
    return put('/profile', data)
  },

  changePassword(data: { old_password: string; new_password: string }) {
    return put('/password', data)
  }
}
