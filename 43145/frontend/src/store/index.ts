import { defineStore } from 'pinia'
import type { User } from '@/types'

interface UserState {
  user: User | null
  token: string
}

export const useUserStore = defineStore('user', {
  state: (): UserState => ({
    user: null,
    token: localStorage.getItem('token') || ''
  }),

  getters: {
    isLoggedIn: (state) => !!state.token,
    userRole: (state) => state.user?.role || '',
    isAdmin: (state) => state.user?.role === 'admin',
    isEditor: (state) => state.user?.role === 'editor',
    isViewer: (state) => state.user?.role === 'viewer'
  },

  actions: {
    setUser(user: User, token: string) {
      this.user = user
      this.token = token
      localStorage.setItem('token', token)
    },

    clearUser() {
      this.user = null
      this.token = ''
      localStorage.removeItem('token')
    },

    setUserInfo(user: User) {
      this.user = user
    }
  }
})
