import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login as apiLogin, register as apiRegister, getProfile } from '@/api/auth'

export interface User {
  id: number
  email: string
  username: string
  avatar?: string
  phone?: string
}

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const user = ref<User | null>(null)

  const isLoggedIn = computed(() => !!token.value)

  async function login(email: string, password: string) {
    const res = await apiLogin({ email, password })
    token.value = res.token
    localStorage.setItem('token', res.token)
    user.value = {
      id: res.userId,
      email: res.email,
      username: res.username,
      avatar: res.avatar
    }
    return res
  }

  async function register(email: string, username: string, password: string, phone?: string) {
    return await apiRegister({ email, username, password, phone })
  }

  async function fetchProfile() {
    try {
      const res = await getProfile()
      user.value = res
    } catch (e) {
      console.error('Failed to fetch profile', e)
    }
  }

  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
  }

  return { token, user, isLoggedIn, login, register, fetchProfile, logout }
})
