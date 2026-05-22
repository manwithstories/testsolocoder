import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { api } from '@/api'

export interface UserInfo {
  id: number
  username: string
  email: string
  real_name: string
  phone: string
  department: string
  role: string
  status: number
}

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const userInfo = ref<UserInfo | null>(null)

  const role = computed(() => userInfo.value?.role || '')
  const isAdmin = computed(() => userInfo.value?.role === 'admin')
  const isSpaceAdmin = computed(() => userInfo.value?.role === 'space_admin')
  const isLoggedIn = computed(() => !!token.value)

  async function login(email: string, password: string) {
    const res: any = await api.login({ email, password })
    token.value = res.data.token
    userInfo.value = res.data.user
    localStorage.setItem('token', res.data.token)
    localStorage.setItem('userInfo', JSON.stringify(res.data.user))
    return res.data
  }

  async function register(data: any) {
    const res: any = await api.register(data)
    return res.data
  }

  async function fetchProfile() {
    try {
      const res: any = await api.getProfile()
      userInfo.value = res.data
      localStorage.setItem('userInfo', JSON.stringify(res.data))
    } catch (e) {
      console.error(e)
    }
  }

  function logout() {
    token.value = ''
    userInfo.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('userInfo')
  }

  function initFromStorage() {
    const saved = localStorage.getItem('userInfo')
    if (saved) {
      try {
        userInfo.value = JSON.parse(saved)
      } catch (e) {
        console.error(e)
      }
    }
  }

  return {
    token,
    userInfo,
    role,
    isAdmin,
    isSpaceAdmin,
    isLoggedIn,
    login,
    register,
    fetchProfile,
    logout,
    initFromStorage
  }
})
