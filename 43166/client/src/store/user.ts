import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { User, type UserRole, LoginRequest, RegisterRequest } from '@/types'
import { authApi } from '@/api/auth'

export const useUserStore = defineStore('user', () => {
  const token = ref<string>(localStorage.getItem('token') || '')
  const userInfo = ref<User | null>(null)

  const isLoggedIn = computed(() => !!token.value)
  const userRole = computed<UserRole | null>(() => userInfo.value?.role || null)

  const setToken = (newToken: string) => {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  const setUserInfo = (user: User | null) => {
    userInfo.value = user
  }

  const login = async (data: LoginRequest) => {
    const res = await authApi.login(data)
    if (res) {
      setToken(res.token)
      setUserInfo(res.userInfo)
    }
    return res
  }

  const register = async (data: RegisterRequest) => {
    return await authApi.register(data)
  }

  const logout = () => {
    token.value = ''
    userInfo.value = null
    localStorage.removeItem('token')
  }

  const fetchUserInfo = async () => {
    try {
      const res = await authApi.getProfile()
      setUserInfo(res || null)
      return res
    } catch (error) {
      logout()
      throw error
    }
  }

  const updateProfile = async (data: Partial<User>) => {
    const res = await authApi.updateProfile(data)
    if (userInfo.value) {
      Object.assign(userInfo.value, data)
    }
    return res
  }

  const changePassword = async (oldPassword: string, newPassword: string) => {
    return await authApi.changePassword(oldPassword, newPassword)
  }

  return {
    token,
    userInfo,
    isLoggedIn,
    userRole,
    setToken,
    setUserInfo,
    login,
    register,
    logout,
    fetchUserInfo,
    updateProfile,
    changePassword
  }
})
