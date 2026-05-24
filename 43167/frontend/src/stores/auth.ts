import { defineStore } from 'pinia'
import request from '@/utils/request'
import type { User } from '@/types'

interface AuthState {
  token: string
  user: User | null
}

export const useAuthStore = defineStore('auth', {
  state: (): AuthState => ({
    token: localStorage.getItem('token') || '',
    user: JSON.parse(localStorage.getItem('user') || 'null')
  }),
  actions: {
    async register(payload: { username: string; password: string; role: string; email?: string }) {
      return request.post('/auth/register', payload)
    },
    async login(payload: { username: string; password: string }) {
      const res: any = await request.post('/auth/login', payload)
      this.token = res.token
      this.user = res.user
      localStorage.setItem('token', res.token)
      localStorage.setItem('user', JSON.stringify(res.user))
      return res
    },
    async me() {
      const u: any = await request.get('/me')
      this.user = u
      localStorage.setItem('user', JSON.stringify(u))
      return u
    },
    logout() {
      this.token = ''
      this.user = null
      localStorage.removeItem('token')
      localStorage.removeItem('user')
    }
  }
})
