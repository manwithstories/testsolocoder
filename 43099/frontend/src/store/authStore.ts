import { create } from 'zustand'
import type { User } from '@/types'
import { userApi, authApi } from '@/api'

interface AuthState {
  token: string | null
  user: User | null
  isAuthenticated: boolean
  login: (username: string, password: string) => Promise<void>
  register: (data: {
    username: string
    email: string
    password: string
    real_name?: string
    phone?: string
  }) => Promise<void>
  logout: () => void
  fetchUser: () => Promise<void>
}

export const useAuthStore = create<AuthState>((set) => ({
  token: localStorage.getItem('token'),
  user: localStorage.getItem('user') ? JSON.parse(localStorage.getItem('user')!) : null,
  isAuthenticated: !!localStorage.getItem('token'),

  login: async (username, password) => {
    const data: any = await authApi.login({ username, password })
    const { token, user } = data
    localStorage.setItem('token', token)
    localStorage.setItem('user', JSON.stringify(user))
    set({ token, user, isAuthenticated: true })
  },

  register: async (data) => {
    const result: any = await authApi.register(data)
    const { token, user } = result
    localStorage.setItem('token', token)
    localStorage.setItem('user', JSON.stringify(user))
    set({ token, user, isAuthenticated: true })
  },

  logout: async () => {
    try {
      await authApi.logout()
    } catch (e) {
      console.error('Logout error:', e)
    }
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    set({ token: null, user: null, isAuthenticated: false })
  },

  fetchUser: async () => {
    try {
      const user = await userApi.getMe()
      localStorage.setItem('user', JSON.stringify(user))
      set({ user })
    } catch (e) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      set({ token: null, user: null, isAuthenticated: false })
    }
  },
}))
