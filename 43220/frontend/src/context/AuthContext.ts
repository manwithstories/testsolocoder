import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { User, LoginData } from '@/types'

interface AuthState {
  token: string | null
  user: User | null
  setAuth: (data: LoginData) => void
  setUser: (user: User) => void
  logout: () => void
  isAuthenticated: () => boolean
  hasRole: (role: string) => boolean
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      token: null,
      user: null,
      setAuth: (data: LoginData) => {
        set({
          token: data.token,
          user: {
            id: data.user_id,
            username: data.username,
            email: '',
            role: data.role as User['role'],
            status: 'active',
            created_at: '',
          },
        })
      },
      setUser: (user: User) => set({ user }),
      logout: () => set({ token: null, user: null }),
      isAuthenticated: () => !!get().token,
      hasRole: (role: string) => get().user?.role === role,
    }),
    {
      name: 'pet-board-auth',
    }
  )
)
