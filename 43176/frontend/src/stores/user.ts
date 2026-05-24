import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User, UserRole } from '@/types'
import { authApi, userApi } from '@/api'

export const useUserStore = defineStore('user', () => {
  const token = ref<string>(localStorage.getItem('token') || '')
  const refreshToken = ref<string>(localStorage.getItem('refreshToken') || '')
  const userInfo = ref<User | null>(null)

  const isLoggedIn = computed(() => !!token.value)
  const userRole = computed<UserRole | ''>(() => userInfo.value?.role || '')
  const isPublisher = computed(() => userRole.value === 'publisher')
  const isCourier = computed(() => userRole.value === 'courier')
  const isAdmin = computed(() => userRole.value === 'admin')

  const setToken = (accessToken: string, refToken: string) => {
    token.value = accessToken
    refreshToken.value = refToken
    localStorage.setItem('token', accessToken)
    localStorage.setItem('refreshToken', refToken)
  }

  const login = async (phone: string, password: string) => {
    const res = await authApi.login({ phone, password })
    if (res.code === 200 && res.data) {
      setToken(res.data.access_token, res.data.refresh_token)
      userInfo.value = {
        id: res.data.user.id,
        phone: res.data.user.phone,
        nickname: res.data.user.nickname,
        avatar: res.data.user.avatar,
        role: res.data.user.role,
        status: res.data.user.status,
        balance: 0,
        rating: 0,
        order_count: 0,
        created_at: ''
      }
      return true
    }
    return false
  }

  const register = async (phone: string, password: string, nickname: string, role: string) => {
    const res = await authApi.register({ phone, password, nickname, role })
    return res.code === 200
  }

  const fetchUserInfo = async () => {
    try {
      const res = await userApi.getProfile()
      if (res.code === 200 && res.data) {
        userInfo.value = res.data as unknown as User
      }
    } catch (error) {
      console.error('Failed to fetch user info:', error)
    }
  }

  const updateUserInfo = (data: Partial<User>) => {
    if (userInfo.value) {
      userInfo.value = { ...userInfo.value, ...data }
    }
  }

  const logout = () => {
    token.value = ''
    refreshToken.value = ''
    userInfo.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('refreshToken')
  }

  return {
    token,
    refreshToken,
    userInfo,
    isLoggedIn,
    userRole,
    isPublisher,
    isCourier,
    isAdmin,
    login,
    register,
    fetchUserInfo,
    updateUserInfo,
    logout,
    setToken
  }
})
