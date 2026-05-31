import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi, userApi } from '@/api'
import type { User } from '@/types'

export const useUserStore = defineStore('user', () => {
  const token = ref<string>('')
  const refreshToken = ref<string>('')
  const userInfo = ref<User | null>(null)

  const isLoggedIn = computed(() => !!token.value)
  const userRole = computed(() => userInfo.value?.role || '')
  const isAdmin = computed(() => userInfo.value?.role === 'admin')
  const isSeller = computed(() => userInfo.value?.role === 'seller')
  const isBuyer = computed(() => userInfo.value?.role === 'buyer')
  const isTechnician = computed(() => userInfo.value?.role === 'technician')

  async function login(username: string, password: string) {
    try {
      const res = await authApi.login({ username, password })
      token.value = res.data.accessToken
      refreshToken.value = res.data.refreshToken
      userInfo.value = res.data.user
      localStorage.setItem('token', res.data.accessToken)
      localStorage.setItem('refreshToken', res.data.refreshToken)
      localStorage.setItem('userInfo', JSON.stringify(res.data.user))
      return res.data
    } catch (error) {
      throw error
    }
  }

  async function register(data: {
    username: string
    password: string
    email?: string
    phone?: string
    role: 'seller' | 'buyer' | 'technician'
    nickname?: string
  }) {
    try {
      const res = await authApi.register(data)
      return res.data
    } catch (error) {
      throw error
    }
  }

  async function logout() {
    try {
      await authApi.logout()
    } catch (error) {
      console.error('Logout error:', error)
    } finally {
      clearUser()
    }
  }

  function clearUser() {
    token.value = ''
    refreshToken.value = ''
    userInfo.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('refreshToken')
    localStorage.removeItem('userInfo')
  }

  async function fetchUserInfo() {
    try {
      const res = await userApi.getProfile()
      userInfo.value = res.data
      localStorage.setItem('userInfo', JSON.stringify(res.data))
      return res.data
    } catch (error) {
      throw error
    }
  }

  async function updateProfile(data: {
    nickname?: string
    email?: string
    phone?: string
    avatar?: string
  }) {
    try {
      await userApi.updateProfile(data)
      if (userInfo.value) {
        Object.assign(userInfo.value, data)
        localStorage.setItem('userInfo', JSON.stringify(userInfo.value))
      }
    } catch (error) {
      throw error
    }
  }

  function initUser() {
    const savedToken = localStorage.getItem('token')
    const savedRefreshToken = localStorage.getItem('refreshToken')
    const savedUserInfo = localStorage.getItem('userInfo')

    if (savedToken) {
      token.value = savedToken
    }
    if (savedRefreshToken) {
      refreshToken.value = savedRefreshToken
    }
    if (savedUserInfo) {
      try {
        userInfo.value = JSON.parse(savedUserInfo)
      } catch (error) {
        console.error('Parse user info error:', error)
      }
    }
  }

  async function refreshTokenFn() {
    if (!refreshToken.value) return
    try {
      const res = await authApi.refreshToken({ refreshToken: refreshToken.value })
      token.value = res.data.accessToken
      localStorage.setItem('token', res.data.accessToken)
    } catch (error) {
      clearUser()
      throw error
    }
  }

  return {
    token,
    refreshToken,
    userInfo,
    isLoggedIn,
    userRole,
    isAdmin,
    isSeller,
    isBuyer,
    isTechnician,
    login,
    register,
    logout,
    clearUser,
    fetchUserInfo,
    updateProfile,
    initUser,
    refreshTokenFn
  }
}, {
  persist: false
})
