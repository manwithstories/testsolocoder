import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User } from '@/types'

export const useUserStore = defineStore('user', () => {
  const token = ref<string>(localStorage.getItem('token') || '')
  const user = ref<User | null>(null)

  const isLoggedIn = computed(() => !!token.value)
  const userRole = computed(() => user.value?.role || '')
  const isVerified = computed(() => user.value?.verified || false)

  function setToken(newToken: string) {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  function setUser(newUser: User) {
    user.value = newUser
    localStorage.setItem('user', JSON.stringify(newUser))
  }

  function loadUser() {
    const savedUser = localStorage.getItem('user')
    if (savedUser) {
      user.value = JSON.parse(savedUser)
    }
  }

  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  function isOwner() {
    return userRole.value === 'owner' || userRole.value === 'admin'
  }

  function isAdmin() {
    return userRole.value === 'admin'
  }

  return {
    token,
    user,
    isLoggedIn,
    userRole,
    isVerified,
    setToken,
    setUser,
    loadUser,
    logout,
    isOwner,
    isAdmin
  }
})
