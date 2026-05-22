import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User } from '@/types'
import { authApi } from '@/api'
import { ElMessage } from 'element-plus'

export const useUserStore = defineStore('user', () => {
  const token = ref<string>(localStorage.getItem('token') || '')
  const userInfo = ref<User | null>(null)
  const unreadCount = ref(0)

  const isLoggedIn = computed(() => !!token.value)
  const isAdmin = computed(() => userInfo.value?.role === 'admin')

  const setToken = (t: string) => {
    token.value = t
    localStorage.setItem('token', t)
  }

  const setUserInfo = (user: User) => {
    userInfo.value = user
    localStorage.setItem('userInfo', JSON.stringify(user))
  }

  const clearUser = () => {
    token.value = ''
    userInfo.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('userInfo')
  }

  const login = async (data: { username: string; password: string }) => {
    const res = await authApi.login(data)
    setToken(res.token)
    setUserInfo(res.user_info)
    ElMessage.success('登录成功')
    return res
  }

  const register = async (data: any) => {
    const res = await authApi.register(data)
    ElMessage.success('注册成功')
    return res
  }

  const logout = () => {
    clearUser()
    ElMessage.success('已退出登录')
  }

  const fetchUserInfo = async () => {
    if (!token.value) return
    try {
      const res = await authApi.getUserInfo()
      setUserInfo(res)
    } catch (e) {
      clearUser()
    }
  }

  const initFromStorage = () => {
    const storedUser = localStorage.getItem('userInfo')
    if (storedUser) {
      try {
        userInfo.value = JSON.parse(storedUser)
      } catch (e) {}
    }
  }

  return {
    token,
    userInfo,
    unreadCount,
    isLoggedIn,
    isAdmin,
    setToken,
    setUserInfo,
    clearUser,
    login,
    register,
    logout,
    fetchUserInfo,
    initFromStorage,
  }
})
