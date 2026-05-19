import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Member } from '@/types'

export const useUserStore = defineStore('user', () => {
  const token = ref<string>(localStorage.getItem('token') || '')
  const userInfo = ref<Member | null>(null)

  const setToken = (t: string) => {
    token.value = t
    localStorage.setItem('token', t)
  }

  const setUserInfo = (user: Member | null) => {
    userInfo.value = user
    if (user) {
      localStorage.setItem('userInfo', JSON.stringify(user))
    } else {
      localStorage.removeItem('userInfo')
    }
  }

  const logout = () => {
    token.value = ''
    userInfo.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('userInfo')
  }

  const isLoggedIn = () => {
    return !!token.value
  }

  return {
    token,
    userInfo,
    setToken,
    setUserInfo,
    logout,
    isLoggedIn
  }
})
