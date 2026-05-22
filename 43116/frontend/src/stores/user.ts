import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { userApi } from '@/api/user'
import type { UserInfo, LoginRequest, RegisterRequest } from '@/types'

export const useUserStore = defineStore('user', () => {
  const accessToken = ref<string>('')
  const refreshToken = ref<string>('')
  const userInfo = ref<UserInfo | null>(null)

  const isLoggedIn = computed(() => !!accessToken.value)
  const isAdmin = computed(() => userInfo.value?.role_id === 1)

  const login = async (data: LoginRequest) => {
    const res = await userApi.login(data)
    accessToken.value = res.data.access_token
    refreshToken.value = res.data.refresh_token
    userInfo.value = res.data.user
  }

  const register = async (data: RegisterRequest) => {
    await userApi.register(data)
  }

  const logout = () => {
    accessToken.value = ''
    refreshToken.value = ''
    userInfo.value = null
  }

  const getProfile = async () => {
    const res = await userApi.getProfile()
    if (res.data) {
      userInfo.value = {
        id: res.data.id,
        username: res.data.username,
        email: res.data.email,
        phone: res.data.phone,
        real_name: res.data.real_name,
        auth_status: res.data.auth_status,
        status: res.data.status,
        role_id: res.data.role_id,
        role_name: res.data.role?.name || '用户',
        avatar: res.data.avatar
      }
    }
  }

  const refreshTokenAction = async () => {
    if (!refreshToken.value) return
    try {
      const res = await userApi.refreshToken(refreshToken.value)
      accessToken.value = res.data.access_token
      refreshToken.value = res.data.refresh_token
    } catch {
      logout()
    }
  }

  return {
    accessToken,
    refreshToken,
    userInfo,
    isLoggedIn,
    isAdmin,
    login,
    register,
    logout,
    getProfile,
    refreshTokenAction
  }
}, {
  persist: {
    key: 'car_rental_user',
    paths: ['accessToken', 'refreshToken', 'userInfo']
  }
})
