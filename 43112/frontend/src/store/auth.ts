import { create } from 'zustand'
import { User } from '@/types'
import { authApi } from '@/services'

interface AuthState {
  token: string | null
  user: User | null
  isAuthenticated: boolean
  setAuth: (token: string, user: User) => void
  clearAuth: () => void
  login: (account: string, password: string) => Promise<void>
  register: (data: any) => Promise<void>
  fetchProfile: () => Promise<void>
  logout: () => void
  checkAuth: () => boolean
}

export const useAuthStore = create<AuthState>((set, get) => ({
  token: localStorage.getItem('token'),
  user: JSON.parse(localStorage.getItem('user') || 'null'),
  isAuthenticated: !!localStorage.getItem('token'),

  setAuth: (token: string, user: User) => {
    localStorage.setItem('token', token)
    localStorage.setItem('user', JSON.stringify(user))
    set({ token, user, isAuthenticated: true })
  },

  clearAuth: () => {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    set({ token: null, user: null, isAuthenticated: false })
  },

  login: async (account: string, password: string) => {
    const res = await authApi.login({ account, password })
    if (res.code === 0 && res.data) {
      get().setAuth(res.data.token, res.data.user)
    } else {
      throw new Error(res.message || 'Login failed')
    }
  },

  register: async (data: any) => {
    const res = await authApi.register(data)
    if (res.code !== 0) {
      throw new Error(res.message || 'Registration failed')
    }
  },

  fetchProfile: async () => {
    try {
      const res = await authApi.getProfile()
      if (res.code === 0 && res.data) {
        localStorage.setItem('user', JSON.stringify(res.data))
        set({ user: res.data })
      }
    } catch (error) {
      console.error('Failed to fetch profile:', error)
    }
  },

  logout: () => {
    get().clearAuth()
  },

  checkAuth: () => {
    return get().isAuthenticated
  },
}))
