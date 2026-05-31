import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User, UserRole } from '@/types/user'
import { loginApi, getProfileApi, logoutApi } from '@/api/user'

export const useUserStore = defineStore('user', () => {
  const token = ref<string>('')
  const user = ref<User | null>(null)

  const isLoggedIn = computed(() => !!token.value)
  const userRole = computed(() => user.value?.role || 'tenant')
  const userId = computed(() => user.value?.id || '')

  const setToken = (newToken: string) => {
    token.value = newToken
  }

  const setUser = (newUser: User | null) => {
    user.value = newUser
  }

  const login = async (email: string, password: string) => {
    const res = await loginApi({ email, password })
    if (res.data) {
      token.value = res.data.token
      user.value = res.data.user
    }
    return res
  }

  const fetchProfile = async () => {
    try {
      const res = await getProfileApi()
      if (res.data) {
        user.value = res.data
      }
      return res
    } catch (error) {
      logout()
      throw error
    }
  }

  const logout = () => {
    token.value = ''
    user.value = null
  }

  const hasRole = (roles: UserRole[]) => {
    if (!user.value) return false
    return roles.includes(user.value.role)
  }

  return {
    token,
    user,
    isLoggedIn,
    userRole,
    userId,
    setToken,
    setUser,
    login,
    fetchProfile,
    logout,
    hasRole
  }
}, {
  persist: {
    key: 'ship-rental-user',
    storage: localStorage,
    paths: ['token', 'user']
  }
})
