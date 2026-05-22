import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { User, UserRole } from '@/types'

interface AuthState {
  user: User | null
  token: string | null
  isAuthenticated: boolean
  login: (token: string, user: User) => void
  logout: () => void
  updateUser: (user: Partial<User>) => void
}

interface AppState {
  theme: 'light' | 'dark'
  collapsed: boolean
  setTheme: (theme: 'light' | 'dark') => void
  toggleCollapsed: () => void
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      token: null,
      isAuthenticated: false,
      login: (token: string, user: User) => {
        set({ token, user, isAuthenticated: true })
        localStorage.setItem('token', token)
        localStorage.setItem('user', JSON.stringify(user))
      },
      logout: () => {
        set({ token: null, user: null, isAuthenticated: false })
        localStorage.removeItem('token')
        localStorage.removeItem('user')
      },
      updateUser: (userData: Partial<User>) =>
        set((state) => ({
          user: state.user ? { ...state.user, ...userData } : null
        }))
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({
        user: state.user,
        token: state.token,
        isAuthenticated: state.isAuthenticated
      })
    }
  )
)

export const useAppStore = create<AppState>((set) => ({
  theme: 'light',
  collapsed: false,
  setTheme: (theme: 'light' | 'dark') => set({ theme }),
  toggleCollapsed: () => set((state) => ({ collapsed: !state.collapsed }))
}))

export const hasRole = (requiredRoles: UserRole | UserRole[]): boolean => {
  const user = useAuthStore.getState().user
  if (!user) return false
  const roles = Array.isArray(requiredRoles) ? requiredRoles : [requiredRoles]
  return roles.includes(user.role)
}
