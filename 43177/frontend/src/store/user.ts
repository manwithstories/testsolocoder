import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User } from '@/types'
import { authApi, userApi } from '@/api/user'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const userInfo = ref<User | null>(null)

  const isLoggedIn = computed(() => !!token.value)
  const userRole = computed(() => userInfo.value?.role || '')
  const isAdmin = computed(() => userRole.value === 'admin')
  const isCustomer = computed(() => userRole.value === 'customer')
  const isTechnician = computed(() => userRole.value === 'technician')

  async function login(username: string, password: string) {
    const res = await authApi.login({ username, password })
    if (res.data) {
      token.value = res.data.token
      localStorage.setItem('token', res.data.token)
      localStorage.setItem('userInfo', JSON.stringify({
        id: res.data.user_id,
        username: res.data.username,
        role: res.data.role,
        status: res.data.status
      }))
      await fetchUserInfo()
    }
    return res
  }

  async function register(data: {
    username: string
    password: string
    phone: string
    email?: string
    real_name?: string
    role: string
  }) {
    return await authApi.register(data)
  }

  async function fetchUserInfo() {
    try {
      const res = await userApi.getProfile()
      userInfo.value = res.data || null
    } catch (error) {
      console.error('Failed to fetch user info:', error)
    }
  }

  function logout() {
    token.value = ''
    userInfo.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('userInfo')
  }

  return {
    token,
    userInfo,
    isLoggedIn,
    userRole,
    isAdmin,
    isCustomer,
    isTechnician,
    login,
    register,
    fetchUserInfo,
    logout
  }
})
