import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { getToken, setToken, removeToken, getUserInfo, setUserInfo, removeUserInfo } from '@/utils/storage'
import type { User, LoginRequest, Role } from '@/types'

export const useUserStore = defineStore('user', () => {
  const token = ref<string>(getToken())
  const userInfo = ref<User | null>(getUserInfo<User>())

  const isLoggedIn = computed(() => !!token.value)
  const userRole = computed<Role | ''>(() => userInfo.value?.role || '')

  // 登录
  async function login(params: LoginRequest): Promise<void> {
    // TODO: 调用真实登录接口
    // const data = await request.post<LoginResponse>('/auth/login', params)
    // token.value = data.token
    // userInfo.value = data.userInfo
    // setToken(data.token)
    // setUserInfo(data.userInfo)
    console.log('login params:', params)
  }

  // 退出登录
  function logout(): void {
    token.value = ''
    userInfo.value = null
    removeToken()
    removeUserInfo()
  }

  // 恢复用户信息（从 localStorage）
  function restoreUserInfo(): void {
    const stored = getUserInfo<User>()
    if (stored) {
      userInfo.value = stored
    }
  }

  // 检查是否有指定角色
  function hasRole(role: Role): boolean {
    return userInfo.value?.role === role
  }

  // 检查是否有任一指定角色
  function hasAnyRole(roles: Role[]): boolean {
    if (!userInfo.value) return false
    return roles.includes(userInfo.value.role)
  }

  return {
    token,
    userInfo,
    isLoggedIn,
    userRole,
    login,
    logout,
    restoreUserInfo,
    hasRole,
    hasAnyRole
  }
})