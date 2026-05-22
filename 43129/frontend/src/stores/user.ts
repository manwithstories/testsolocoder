import { defineStore } from 'pinia'
import { login, register, getProfile } from '@/api/auth'
import type { User } from '@/types'

export const useUserStore = defineStore('user', {
  state: () => ({
    token: localStorage.getItem('token') || '',
    user: JSON.parse(localStorage.getItem('user') || 'null') as User | null
  }),

  getters: {
    isLoggedIn: (state) => !!state.token,
    isAdmin: (state) => state.user?.role === 'admin',
    isTechnician: (state) => state.user?.role === 'technician',
    isCustomer: (state) => state.user?.role === 'customer'
  },

  actions: {
    async login(phone: string, password: string) {
      const res = await login({ phone, password })
      this.token = res.data.token
      this.user = res.data.user
      localStorage.setItem('token', this.token)
      localStorage.setItem('user', JSON.stringify(this.user))
      return res
    },

    async register(data: { phone: string; password: string; nickname?: string; role?: string }) {
      const res = await register(data)
      this.token = res.data.token
      this.user = res.data.user
      localStorage.setItem('token', this.token)
      localStorage.setItem('user', JSON.stringify(this.user))
      return res
    },

    async fetchProfile() {
      const res = await getProfile()
      this.user = res.data
      localStorage.setItem('user', JSON.stringify(this.user))
      return res
    },

    logout() {
      this.token = ''
      this.user = null
      localStorage.removeItem('token')
      localStorage.removeItem('user')
    }
  }
})
