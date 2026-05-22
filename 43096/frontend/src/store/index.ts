import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User } from '@/types'
import { userApi } from '@/api'

export const useUserStore = defineStore('user', () => {
  const token = ref<string>(localStorage.getItem('token') || '')
  const user = ref<User | null>(localStorage.getItem('user') ? JSON.parse(localStorage.getItem('user')!) : null)

  const isLoggedIn = computed(() => !!token.value && !!user.value)
  const isAdmin = computed(() => user.value?.role === 'admin')

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

  async function fetchUserInfo() {
    if (!token.value) return
    try {
      const res = await userApi.getInfo()
      user.value = res as User
      localStorage.setItem('user', JSON.stringify(res))
    } catch (err) {
      console.error(err)
    }
  }

  return {
    token,
    user,
    isLoggedIn,
    isAdmin,
    setToken,
    setUser,
    logout,
    fetchUserInfo
  }
})

export const useSeatStore = defineStore('seat', () => {
  const lockedSeats = ref<Set<number>>(new Set())
  const selectedSeats = ref<Set<number>>(new Set())

  function toggleSeat(seatId: number) {
    if (selectedSeats.value.has(seatId)) {
      selectedSeats.value.delete(seatId)
    } else {
      selectedSeats.value.add(seatId)
    }
  }

  function clearSelected() {
    selectedSeats.value.clear()
  }

  function lockSeats(seatIds: number[]) {
    seatIds.forEach(id => lockedSeats.value.add(id))
  }

  function unlockSeats(seatIds: number[]) {
    seatIds.forEach(id => lockedSeats.value.delete(id))
  }

  return {
    lockedSeats,
    selectedSeats,
    toggleSeat,
    clearSelected,
    lockSeats,
    unlockSeats
  }
})
