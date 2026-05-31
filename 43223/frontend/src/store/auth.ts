import { create } from 'zustand'
import { User, UserRole } from '@/types'
import { authApi } from '@/api/auth'

interface AuthState {
  user: User | null
  token: string | null
  isAuthenticated: boolean
  isLoading: boolean
  login: (account: string, password: string) => Promise<void>
  register: (data: { username: string; email: string; password: string; phone?: string; nickname?: string }) => Promise<void>
  logout: () => void
  setUser: (user: User) => void
  loadUser: () => Promise<void>
  hasRole: (role: UserRole | UserRole[]) => boolean
}

export const useAuthStore = create<AuthState>((set, get) => ({
  user: null,
  token: localStorage.getItem('token'),
  isAuthenticated: !!localStorage.getItem('token'),
  isLoading: false,

  login: async (account, password) => {
    set({ isLoading: true })
    try {
      const res = await authApi.login({ account, password })
      if (res.data) {
        localStorage.setItem('token', res.data.token)
        localStorage.setItem('user', JSON.stringify(res.data.user))
        set({
          user: res.data.user,
          token: res.data.token,
          isAuthenticated: true,
        })
      }
    } finally {
      set({ isLoading: false })
    }
  },

  register: async (data) => {
    set({ isLoading: true })
    try {
      const res = await authApi.register(data)
      if (res.data) {
        localStorage.setItem('token', res.data.token)
        localStorage.setItem('user', JSON.stringify(res.data.user))
        set({
          user: res.data.user,
          token: res.data.token,
          isAuthenticated: true,
        })
      }
    } finally {
      set({ isLoading: false })
    }
  },

  logout: () => {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    set({ user: null, token: null, isAuthenticated: false })
  },

  setUser: (user) => {
    set({ user })
    localStorage.setItem('user', JSON.stringify(user))
  },

  loadUser: async () => {
    try {
      const res = await authApi.getProfile()
      if (res.data) {
        set({ user: res.data })
        localStorage.setItem('user', JSON.stringify(res.data))
      }
    } catch {
      get().logout()
    }
  },

  hasRole: (role) => {
    const user = get().user
    if (!user) return false
    if (Array.isArray(role)) {
      return role.includes(user.role)
    }
    return user.role === role
  },
}))
