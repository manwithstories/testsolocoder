import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import api from '../api'

interface User {
  id: number
  username: string
  email: string
  role: string
  avatar: string
  reputation: number
}

interface AuthState {
  token: string | null
  user: User | null
  isAuthenticated: boolean
  login: (username: string, password: string) => Promise<void>
  register: (data: RegisterData) => Promise<void>
  logout: () => void
  fetchProfile: () => Promise<void>
}

interface RegisterData {
  username: string
  password: string
  email: string
  phone?: string
  role: string
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      token: null,
      user: null,
      isAuthenticated: false,
      
      login: async (username: string, password: string) => {
        const response = await api.post('/auth/login', { username, password })
        const { token, user_info } = response.data
        set({ token, user: user_info, isAuthenticated: true })
        localStorage.setItem('token', token)
      },
      
      register: async (data: RegisterData) => {
        await api.post('/auth/register', data)
      },
      
      logout: () => {
        set({ token: null, user: null, isAuthenticated: false })
        localStorage.removeItem('token')
      },
      
      fetchProfile: async () => {
        const response = await api.get('/auth/profile')
        set({ user: response.data })
      },
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({ token: state.token, user: state.user, isAuthenticated: state.isAuthenticated }),
    }
  )
)
