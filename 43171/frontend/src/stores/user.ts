import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import request from '@/utils/request'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const userInfo = ref<UserInfo | null>(null)

  const isLoggedIn = computed(() => !!token.value)
  const role = computed(() => userInfo.value?.role || '')

  async function login(username: string, password: string) {
    const res: ApiResponse<LoginResp> = await request.post('/auth/login', { username, password })
    token.value = res.data.token
    localStorage.setItem('token', res.data.token)
    userInfo.value = {
      user_id: res.data.user_id,
      username: res.data.username,
      nickname: res.data.nickname,
      avatar: res.data.avatar,
      role: res.data.role as 'client' | 'pilot' | 'owner',
      verify_status: res.data.verify_status as 'pending' | 'approved' | 'rejected'
    }
    localStorage.setItem('userInfo', JSON.stringify(userInfo.value))
    return res.data
  }

  async function register(data: { username: string; password: string; role: string; nickname?: string }) {
    return request.post('/auth/register', data)
  }

  async function fetchUserInfo() {
    const res: ApiResponse<UserInfo> = await request.get('/user/profile')
    userInfo.value = res.data
    localStorage.setItem('userInfo', JSON.stringify(res.data))
    return res.data
  }

  function logout() {
    token.value = ''
    userInfo.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('userInfo')
  }

  function restoreSession() {
    const stored = localStorage.getItem('userInfo')
    if (stored) {
      userInfo.value = JSON.parse(stored)
    }
  }

  return {
    token,
    userInfo,
    isLoggedIn,
    role,
    login,
    register,
    fetchUserInfo,
    logout,
    restoreSession
  }
})

interface LoginResp {
  token: string
  user_id: number
  username: string
  role: string
  nickname: string
  avatar: string
  verify_status: string
}
