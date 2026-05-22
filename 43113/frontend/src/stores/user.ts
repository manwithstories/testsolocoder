import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User } from '@/types'
import { authApi } from '@/api/user'

export const useUserStore = defineStore('user', () => {
  const token = ref<string>(localStorage.getItem('token') || '')
  const userInfo = ref<User | null>(null)

  const isLoggedIn = computed(() => !!token.value)
  const isAdmin = computed(() => userInfo.value?.role === 'admin')
  const isExpert = computed(() => userInfo.value?.isExpert === true)

  async function login(username: string, password: string) {
    const res = await authApi.login({ username, password })
    if (res.data) {
      token.value = res.data.token
      userInfo.value = res.data.user
      localStorage.setItem('token', res.data.token)
      localStorage.setItem('userInfo', JSON.stringify(res.data.user))
    }
  }

  async function register(username: string, email: string, password: string, nickname?: string) {
    const res = await authApi.register({ username, email, password, nickname })
    if (res.data) {
      token.value = res.data.token
      userInfo.value = res.data.user
      localStorage.setItem('token', res.data.token)
      localStorage.setItem('userInfo', JSON.stringify(res.data.user))
    }
  }

  function logout() {
    token.value = ''
    userInfo.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('userInfo')
  }

  function loadUserInfo() {
    const stored = localStorage.getItem('userInfo')
    if (stored) {
      try {
        userInfo.value = JSON.parse(stored)
      } catch (e) {
        console.error('Failed to parse user info:', e)
      }
    }
  }

  function setUserInfo(user: User) {
    userInfo.value = user
    localStorage.setItem('userInfo', JSON.stringify(user))
  }

  return {
    token,
    userInfo,
    isLoggedIn,
    isAdmin,
    isExpert,
    login,
    register,
    logout,
    loadUserInfo,
    setUserInfo
  }
})
