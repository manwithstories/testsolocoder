import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface UserInfo {
  id: number
  username: string
  role: 'company' | 'staff' | 'customer' | 'admin'
  real_name?: string
  phone?: string
  expires_at?: number
}

export const useUserStore = defineStore('user', () => {
  const token = ref<string>(localStorage.getItem('token') || '')
  const user = ref<UserInfo | null>(null)
  try {
    const raw = localStorage.getItem('user')
    if (raw) user.value = JSON.parse(raw)
  } catch {}

  const isLoggedIn = computed(() => !!token.value)
  const role = computed(() => user.value?.role || '')

  function setAuth(t: string, u: UserInfo) {
    token.value = t
    user.value = u
    localStorage.setItem('token', t)
    localStorage.setItem('user', JSON.stringify(u))
  }

  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  function updateUser(patch: Partial<UserInfo>) {
    if (user.value) {
      user.value = { ...user.value, ...patch }
      localStorage.setItem('user', JSON.stringify(user.value))
    }
  }

  return { token, user, isLoggedIn, role, setAuth, logout, updateUser }
})
