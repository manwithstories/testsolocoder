import { create } from 'zustand'
import { CartItem, Product } from '@/types'
import { cartApi } from '@/api/order'

interface CartState {
  items: CartItem[]
  totalAmount: number
  totalCount: number
  isLoading: boolean
  loadCart: () => Promise<void>
  addToCart: (productId: number, quantity: number) => Promise<void>
  updateQuantity: (id: number, quantity: number) => Promise<void>
  removeFromCart: (id: number) => Promise<void>
  clearCart: () => Promise<void>
}

export const useCartStore = create<CartState>((set) => ({
  items: [],
  totalAmount: 0,
  totalCount: 0,
  isLoading: false,

  loadCart: async () => {
    set({ isLoading: true })
    try {
      const res = await cartApi.get()
      if (res.data) {
        set({
          items: res.data.items,
          totalAmount: res.data.total_amount,
          totalCount: res.data.total_count,
        })
      }
    } catch {
      set({ items: [], totalAmount: 0, totalCount: 0 })
    } finally {
      set({ isLoading: false })
    }
  },

  addToCart: async (productId, quantity) => {
    set({ isLoading: true })
    try {
      await cartApi.add(productId, quantity)
      await useCartStore.getState().loadCart()
    } finally {
      set({ isLoading: false })
    }
  },

  updateQuantity: async (id, quantity) => {
    set({ isLoading: true })
    try {
      await cartApi.update(id, quantity)
      await useCartStore.getState().loadCart()
    } finally {
      set({ isLoading: false })
    }
  },

  removeFromCart: async (id) => {
    set({ isLoading: true })
    try {
      await cartApi.remove(id)
      await useCartStore.getState().loadCart()
    } finally {
      set({ isLoading: false })
    }
  },

  clearCart: async () => {
    set({ isLoading: true })
    try {
      await cartApi.clear()
      set({ items: [], totalAmount: 0, totalCount: 0 })
    } finally {
      set({ isLoading: false })
    }
  },
}))
