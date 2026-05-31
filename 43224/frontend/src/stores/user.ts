import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login as apiLogin, register as apiRegister, getCurrentUser } from '@/api/auth'

export interface UserInfo {
  id: number
  username: string
  email: string
  role: string
  real_name?: string
  avatar?: string
  language_pairs?: any[]
  expertise_tags?: any[]
  rating?: number
  completed_count?: number
  current_workload?: number
  daily_capacity?: number
}

export const useUserStore = defineStore('user', () => {
  const token = ref<string>(localStorage.getItem('token') || '')
  const userInfo = ref<UserInfo | null>(null)

  const role = computed(() => userInfo.value?.role || '')
  const isLoggedIn = computed(() => !!token.value)

  async function login(username: string, password: string) {
    const res = await apiLogin({ username, password })
    token.value = res.token
    userInfo.value = res.user
    localStorage.setItem('token', res.token)
    return res
  }

  async function register(data: any) {
    return await apiRegister(data)
  }

  async function fetchUserInfo() {
    if (!token.value) return
    try {
      const res = await getCurrentUser()
      userInfo.value = res
    } catch (e) {
      console.error('获取用户信息失败', e)
    }
  }

  function logout() {
    token.value = ''
    userInfo.value = null
    localStorage.removeItem('token')
  }

  function hasRole(roles: string[]) {
    return roles.includes(role.value)
  }

  return {
    token,
    userInfo,
    role,
    isLoggedIn,
    login,
    register,
    fetchUserInfo,
    logout,
    hasRole
  }
})
