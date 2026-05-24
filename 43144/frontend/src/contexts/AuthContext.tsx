import React, { createContext, useContext, useState, useEffect, ReactNode, useCallback } from 'react'
import { User, LoginRequest, RegisterRequest, ApiResponse, LoginResponse } from '../types'
import { login as apiLogin, register as apiRegister, getProfile } from '../api/auth'

interface AuthContextType {
  user: User | null
  token: string | null
  isAuthenticated: boolean
  login: (email: string, password: string) => Promise<void>
  register: (data: RegisterRequest) => Promise<void>
  logout: () => void
  loadUser: () => Promise<void>
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export const AuthProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null)
  const [token, setToken] = useState<string | null>(localStorage.getItem('token'))

  const isAuthenticated = !!token && !!user

  const loadUser = useCallback(async () => {
    try {
      const response = await getProfile()
      if (response.code === 0 && response.data) {
        setUser(response.data)
        localStorage.setItem('user', JSON.stringify(response.data))
      }
    } catch (error) {
      console.error('Failed to load user:', error)
      logout()
    }
  }, [])

  const login = async (email: string, password: string) => {
    const response = await apiLogin({ email, password })
    if (response.code === 0 && response.data) {
      const data = response.data as LoginResponse
      setToken(data.token)
      localStorage.setItem('token', data.token)
      await loadUser()
    } else {
      throw new Error(response.message || '登录失败')
    }
  }

  const register = async (data: RegisterRequest) => {
    const response = await apiRegister(data)
    if (response.code === 0) {
      return
    }
    throw new Error(response.message || '注册失败')
  }

  const logout = () => {
    setUser(null)
    setToken(null)
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  useEffect(() => {
    if (token && !user) {
      loadUser()
    }
  }, [token])

  return (
    <AuthContext.Provider
      value={{
        user,
        token,
        isAuthenticated,
        login,
        register,
        logout,
        loadUser,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export const useAuth = () => {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
