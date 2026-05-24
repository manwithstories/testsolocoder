import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User, UserRole } from '@/types'
import { login as apiLogin, register as apiRegister, getProfile } from '@/api/user'

export const useUserStore = defineStore('user', () => {
  const token = ref<string>(localStorage.getItem('token') || '')
  const user = ref<User | null>(null)
  const loading = ref(false)

  const isLoggedIn = computed(() => !!token.value)
  const userRole = computed(() => user.value?.role || '')
  const userId = computed(() => user.value?.id || 0)
  const userName = computed(() => {
    if (user.value?.profile?.full_name) return user.value.profile.full_name
    if (user.value?.company?.name) return user.value.company.name
    return user.value?.email || ''
  })

  async function login(email: string, password: string) {
    loading.value = true
    try {
      const res = await apiLogin({ email, password })
      if (res.data) {
        token.value = res.data.token
        user.value = res.data.user
        localStorage.setItem('token', res.data.token)
        localStorage.setItem('user', JSON.stringify(res.data.user))
      }
      return res
    } finally {
      loading.value = false
    }
  }

  async function register(data: any) {
    loading.value = true
    try {
      const res = await apiRegister(data)
      if (res.data) {
        token.value = res.data.token
        user.value = res.data.user
        localStorage.setItem('token', res.data.token)
        localStorage.setItem('user', JSON.stringify(res.data.user))
      }
      return res
    } finally {
      loading.value = false
    }
  }

  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  async function fetchProfile() {
    try {
      const res = await getProfile()
      if (res.data) {
        user.value = res.data
        localStorage.setItem('user', JSON.stringify(res.data))
      }
    } catch (e) {
      console.error('获取用户信息失败', e)
    }
  }

  function initFromStorage() {
    const storedUser = localStorage.getItem('user')
    if (storedUser) {
      try {
        user.value = JSON.parse(storedUser)
      } catch (e) {
        localStorage.removeItem('user')
      }
    }
  }

  function hasRole(role: UserRole | UserRole[]) {
    if (!user.value) return false
    if (Array.isArray(role)) {
      return role.includes(user.value.role)
    }
    return user.value.role === role
  }

  return {
    token,
    user,
    loading,
    isLoggedIn,
    userRole,
    userId,
    userName,
    login,
    register,
    logout,
    fetchProfile,
    initFromStorage,
    hasRole
  }
})
