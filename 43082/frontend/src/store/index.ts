import { create } from 'zustand'
import { User, CartItem } from '@/types'
import { authAPI, cartAPI } from '@/api'
import { message } from 'antd'

interface AppState {
  user: User | null
  token: string | null
  cart: CartItem[]
  cartCount: number
  setUser: (user: User | null) => void
  setToken: (token: string | null) => void
  login: (data: { token: string; user: User }) => void
  logout: () => void
  loadUser: () => Promise<void>
  loadCart: () => Promise<void>
  addToCart: (data: { productId: number; skuId?: number; quantity: number }) => Promise<void>
  updateCart: (id: number, quantity: number) => Promise<void>
  removeFromCart: (id: number) => Promise<void>
  clearCart: () => void
}

export const useAppStore = create<AppState>((set, get) => ({
  user: null,
  token: localStorage.getItem('token'),
  cart: [],
  cartCount: 0,

  setUser: (user) => set({ user }),
  setToken: (token) => set({ token }),

  login: (data) => {
    localStorage.setItem('token', data.token)
    localStorage.setItem('user', JSON.stringify(data.user))
    set({ token: data.token, user: data.user })
  },

  logout: () => {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    set({ token: null, user: null, cart: [], cartCount: 0 })
  },

  loadUser: async () => {
    try {
      const res = await authAPI.getProfile() as any
      set({ user: res.data })
    } catch (error) {
      get().logout()
    }
  },

  loadCart: async () => {
    try {
      const res = await cartAPI.list() as any
      const cart = res.data
      const count = cart.reduce((sum: number, item: CartItem) => sum + item.quantity, 0)
      set({ cart, cartCount: count })
    } catch (error) {}
  },

  addToCart: async (data) => {
    await cartAPI.add(data)
    await get().loadCart()
    message.success('已加入购物车')
  },

  updateCart: async (id, quantity) => {
    await cartAPI.update(id, { quantity })
    await get().loadCart()
  },

  removeFromCart: async (id) => {
    await cartAPI.delete(id)
    await get().loadCart()
    message.success('已移除')
  },

  clearCart: () => {
    set({ cart: [], cartCount: 0 })
  },
}))
