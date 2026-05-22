import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { userApi } from '@/api'
import type { LoginResp, User } from '@/types'

export const useUserStore = defineStore('user', () => {
  const token = ref<string>(localStorage.getItem('token') || '')
  const userInfo = ref<User | null>(null)

  const isLogin = computed(() => !!token.value)
  const isAdmin = computed(() => userInfo.value?.role === 'admin')

  function setToken(t: string) {
    token.value = t
    localStorage.setItem('token', t)
  }
  function logout() {
    token.value = ''
    userInfo.value = null
    localStorage.removeItem('token')
  }
  async function fetchProfile() {
    try {
      const res = await userApi.profile()
      userInfo.value = res.data as User
    } catch (_) {
      // ignore
    }
  }
  async function login(username: string, password: string): Promise<LoginResp> {
    const res = await userApi.login({ username, password })
    const data = res.data as LoginResp
    setToken(data.token)
    userInfo.value = {
      id: data.user_id,
      username: data.username,
      role: data.role,
      real_name: data.real_name,
      verified: data.verified,
    } as User
    return data
  }
  return { token, userInfo, isLogin, isAdmin, setToken, logout, fetchProfile, login }
})
