import { create } from 'zustand'
import type { Cart } from '@/types'

interface CartState {
  items: Cart[]
  totalCount: number
  totalAmount: number
  setItems: (items: Cart[]) => void
  addItem: (item: Cart) => void
  updateItem: (id: string, quantity: number) => void
  removeItem: (id: string) => void
  clearCart: () => void
  calculateTotals: () => void
}

export const useCartStore = create<CartState>((set, get) => ({
  items: [],
  totalCount: 0,
  totalAmount: 0,

  setItems: (items) => {
    set({ items })
    get().calculateTotals()
  },

  addItem: (item) => {
    const { items } = get()
    const existing = items.find((i) => i.product_id === item.product_id)
    if (existing) {
      set({
        items: items.map((i) =>
          i.id === existing.id ? { ...i, quantity: i.quantity + item.quantity } : i
        ),
      })
    } else {
      set({ items: [...items, item] })
    }
    get().calculateTotals()
  },

  updateItem: (id, quantity) => {
    if (quantity <= 0) {
      get().removeItem(id)
    } else {
      set({
        items: get().items.map((i) => (i.id === id ? { ...i, quantity } : i)),
      })
      get().calculateTotals()
    }
  },

  removeItem: (id) => {
    set({ items: get().items.filter((i) => i.id !== id) })
    get().calculateTotals()
  },

  clearCart: () => {
    set({ items: [], totalCount: 0, totalAmount: 0 })
  },

  calculateTotals: () => {
    const { items } = get()
    let totalCount = 0
    let totalAmount = 0
    items.forEach((item) => {
      totalCount += item.quantity
      if (item.product) {
        totalAmount += item.product.price * item.quantity
      }
    })
    set({ totalCount, totalAmount })
  },
}))
