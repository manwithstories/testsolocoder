import { create } from 'zustand'

interface User {
  id: number
  username: string
  phone: string
  email: string
  role: string
  status: string
  verify_status: string
  real_name: string
  avatar: string
  member_level: string
  member_expire: string | null
}

interface AuthState {
  token: string | null
  user: User | null
  setToken: (token: string) => void
  setUser: (user: User) => void
  logout: () => void
  isVerified: () => boolean
  isMatchmaker: () => boolean
  isAdmin: () => boolean
}

export const useAuthStore = create<AuthState>((set, get) => ({
  token: localStorage.getItem('token'),
  user: JSON.parse(localStorage.getItem('user') || 'null'),
  setToken: (token) => {
    localStorage.setItem('token', token)
    set({ token })
  },
  setUser: (user) => {
    localStorage.setItem('user', JSON.stringify(user))
    set({ user })
  },
  logout: () => {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    set({ token: null, user: null })
  },
  isVerified: () => {
    const user = get().user
    return user?.verify_status === 'verified'
  },
  isMatchmaker: () => {
    const user = get().user
    return user?.role === 'matchmaker' || user?.role === 'admin'
  },
  isAdmin: () => {
    const user = get().user
    return user?.role === 'admin'
  },
}))
