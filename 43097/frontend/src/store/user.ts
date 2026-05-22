import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login as loginApi, getProfile as getProfileApi, logout as logoutApi } from '@/api/auth'
import type { User, LoginRequest } from '@/types'

const TOKEN_KEY = 'hotel_token'
const USER_KEY = 'hotel_user'

export const useUserStore = defineStore('user', () => {
  const token = ref<string>(localStorage.getItem(TOKEN_KEY) || '')
  const userInfo = ref<User | null>(null)

  const isLoggedIn = computed(() => !!token.value)
  const userRole = computed(() => userInfo.value?.role)
  const userName = computed(() => userInfo.value?.name || '')

  const initFromStorage = () => {
    const storedUser = localStorage.getItem(USER_KEY)
    if (storedUser) {
      try {
        userInfo.value = JSON.parse(storedUser)
      } catch (e) {
        console.error('Failed to parse user info from storage:', e)
      }
    }
  }

  const login = async (credentials: LoginRequest) => {
    const res = await loginApi(credentials)
    token.value = res.token
    userInfo.value = res.user
    localStorage.setItem(TOKEN_KEY, res.token)
    localStorage.setItem(USER_KEY, JSON.stringify(res.user))
    return res
  }

  const getProfile = async () => {
    const user = await getProfileApi()
    userInfo.value = user
    localStorage.setItem(USER_KEY, JSON.stringify(user))
    return user
  }

  const updateUserInfo = (user: Partial<User>) => {
    if (userInfo.value) {
      userInfo.value = { ...userInfo.value, ...user }
      localStorage.setItem(USER_KEY, JSON.stringify(userInfo.value))
    }
  }

  const logout = async () => {
    try {
      await logoutApi()
    } catch (e) {
      console.error('Logout API error:', e)
    } finally {
      clearUser()
    }
  }

  const clearUser = () => {
    token.value = ''
    userInfo.value = null
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(USER_KEY)
  }

  return {
    token,
    userInfo,
    isLoggedIn,
    userRole,
    userName,
    initFromStorage,
    login,
    getProfile,
    updateUserInfo,
    logout,
    clearUser
  }
})
