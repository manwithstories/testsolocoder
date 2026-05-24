import { defineStore } from 'pinia'
import type { User } from '@/types'
import { authApi } from '@/api/auth'
import { ElMessage } from 'element-plus'

interface UserState {
  token: string
  user: User | null
}

export const useUserStore = defineStore('user', {
  state: (): UserState => ({
    token: localStorage.getItem('token') || '',
    user: JSON.parse(localStorage.getItem('user') || 'null')
  }),

  getters: {
    isLoggedIn: (state) => !!state.token && !!state.user,
    userRole: (state) => state.user?.role || '',
    userID: (state) => state.user?.id || 0,
    username: (state) => state.user?.username || ''
  },

  actions: {
    async login(username: string, password: string) {
      const res = await authApi.login({ username, password })
      if (res.code === 200 && res.data) {
        this.token = res.data.token
        this.user = res.data.user
        localStorage.setItem('token', res.data.token)
        localStorage.setItem('user', JSON.stringify(res.data.user))
        ElMessage.success('登录成功')
        return true
      }
      return false
    },

    async register(data: {
      username: string
      email: string
      password: string
      phone?: string
      real_name?: string
      role: string
    }) {
      const res = await authApi.register(data)
      if (res.code === 201 || res.code === 200) {
        ElMessage.success('注册成功，请登录')
        return true
      }
      return false
    },

    async registerAuthenticator(data: {
      username: string
      email: string
      password: string
      phone?: string
      real_name?: string
      role: string
      license_number: string
      license_file: string
      certifications?: string
      specialties?: string
    }) {
      const res = await authApi.registerAuthenticator(data)
      if (res.code === 201 || res.code === 200) {
        ElMessage.success('注册成功，请等待审核')
        return true
      }
      return false
    },

    logout() {
      this.token = ''
      this.user = null
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      ElMessage.success('已退出登录')
    },

    async fetchUserProfile() {
      const res = await authApi.getProfile()
      if (res.code === 200 && res.data) {
        this.user = res.data
        localStorage.setItem('user', JSON.stringify(res.data))
      }
    },

    async updateProfile(data: Record<string, any>) {
      const res = await authApi.updateProfile(data)
      if (res.code === 200 && res.data) {
        this.user = res.data
        localStorage.setItem('user', JSON.stringify(res.data))
        ElMessage.success('更新成功')
        return true
      }
      return false
    }
  }
})
