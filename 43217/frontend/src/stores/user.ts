import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User, LoginResponse } from '@/types'
import { login as apiLogin, logout as apiLogout, getUserInfo } from '@/api/auth'

export const useUserStore = defineStore('user', () => {
  const token = ref<string>('')
  const userInfo = ref<User | null>(null)

  const isLoggedIn = computed(() => !!token.value)
  const userRole = computed(() => userInfo.value?.role || '')
  const companyId = computed(() => userInfo.value?.company_id)
  const agencyId = computed(() => userInfo.value?.agency_id)

  const setToken = (newToken: string) => {
    token.value = newToken
  }

  const setUserInfo = (info: User) => {
    userInfo.value = info
  }

  const login = async (username: string, password: string): Promise<LoginResponse> => {
    const response = await apiLogin(username, password)
    setToken(response.token)
    setUserInfo(response.user)
    return response
  }

  const logout = () => {
    token.value = ''
    userInfo.value = null
  }

  const fetchUserInfo = async () => {
    try {
      const info = await getUserInfo()
      setUserInfo(info)
      return info
    } catch (error) {
      console.error('Failed to fetch user info:', error)
      return null
    }
  }

  return {
    token,
    userInfo,
    isLoggedIn,
    userRole,
    companyId,
    agencyId,
    setToken,
    setUserInfo,
    login,
    logout,
    fetchUserInfo
  }
}, {
  persist: {
    paths: ['token', 'userInfo']
  }
})
