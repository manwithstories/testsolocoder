import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Book, Tag, Category } from '@/types'
import { getBooks, getCurrentlyReading } from '@/api/book'
import { getTags } from '@/api/common'

export const useBookStore = defineStore('book', () => {
  const books = ref<Book[]>([])
  const currentlyReading = ref<Book[]>([])
  const tags = ref<Tag[]>([])
  const categories = ref<Category[]>([])
  const loading = ref(false)
  const total = ref(0)

  const tagMap = computed(() => {
    const map: Record<number, Tag> = {}
    tags.value.forEach(t => map[t.id] = t)
    return map
  })

  const categoryMap = computed(() => {
    const map: Record<number, Category> = {}
    categories.value.forEach(c => map[c.id] = c)
    return map
  })

  const fetchBooks = async (params: any = {}) => {
    loading.value = true
    try {
      const res = await getBooks(params)
      books.value = res.data
      total.value = res.total
    } finally {
      loading.value = false
    }
  }

  const fetchCurrentlyReading = async () => {
    try {
      const res = await getCurrentlyReading()
      currentlyReading.value = res
    } catch (e) {
      console.error(e)
    }
  }

  const fetchTags = async () => {
    try {
      const res = await getTags()
      tags.value = res
    } catch (e) {
      console.error(e)
    }
  }

  return {
    books,
    currentlyReading,
    tags,
    categories,
    loading,
    total,
    tagMap,
    categoryMap,
    fetchBooks,
    fetchCurrentlyReading,
    fetchTags
  }
})
