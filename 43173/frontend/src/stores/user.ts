import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { ArtistInfo } from '@/api/auth'

export interface User {
  id: number
  username: string
  nickname: string
  avatar: string
  role: string
  email?: string
  phone?: string
  bio?: string
  artist_info?: ArtistInfo
}

export const useUserStore = defineStore('user', () => {
  const token = ref<string>(localStorage.getItem('token') || '')
  const user = ref<User | null>(null)
  
  const userStr = localStorage.getItem('user')
  if (userStr) {
    try {
      user.value = JSON.parse(userStr)
    } catch (e) {
      console.error('Failed to parse user data')
    }
  }

  const isLoggedIn = computed(() => !!token.value)
  const isAdmin = computed(() => user.value?.role === 'admin')
  const isArtist = computed(() => user.value?.role === 'artist' || user.value?.role === 'label')

  function setToken(newToken: string) {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  function setUser(newUser: User) {
    user.value = newUser
    localStorage.setItem('user', JSON.stringify(newUser))
  }

  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  function updateUser(data: Partial<User>) {
    if (user.value) {
      user.value = { ...user.value, ...data }
      localStorage.setItem('user', JSON.stringify(user.value))
    }
  }

  return {
    token,
    user,
    isLoggedIn,
    isAdmin,
    isArtist,
    setToken,
    setUser,
    logout,
    updateUser
  }
}, {
  persist: {
    key: 'user-store',
    storage: localStorage
  }
})
