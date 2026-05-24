// Token 存储键
const TOKEN_KEY = 'admin_token'
// 用户信息存储键
const USER_INFO_KEY = 'admin_user_info'

// 获取 Token
export const getToken = (): string => {
  return localStorage.getItem(TOKEN_KEY) || ''
}

// 设置 Token
export const setToken = (token: string): void => {
  localStorage.setItem(TOKEN_KEY, token)
}

// 移除 Token
export const removeToken = (): void => {
  localStorage.removeItem(TOKEN_KEY)
}

// 获取用户信息
export const getUserInfo = <T = unknown>(): T | null => {
  const data = localStorage.getItem(USER_INFO_KEY)
  if (data) {
    try {
      return JSON.parse(data) as T
    } catch {
      return null
    }
  }
  return null
}

// 设置用户信息
export const setUserInfo = <T = unknown>(userInfo: T): void => {
  localStorage.setItem(USER_INFO_KEY, JSON.stringify(userInfo))
}

// 移除用户信息
export const removeUserInfo = (): void => {
  localStorage.removeItem(USER_INFO_KEY)
}

// 通用 localStorage 封装
export const storage = {
  get<T = unknown>(key: string): T | null {
    const data = localStorage.getItem(key)
    if (data) {
      try {
        return JSON.parse(data) as T
      } catch {
        return data as unknown as T
      }
    }
    return null
  },
  set(key: string, value: unknown): void {
    localStorage.setItem(key, typeof value === 'string' ? value : JSON.stringify(value))
  },
  remove(key: string): void {
    localStorage.removeItem(key)
  },
  clear(): void {
    localStorage.clear()
  }
}

export default storage