import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Category } from '@/types'
import { storage } from '@/utils/storage'
import { logger } from '@/utils/logger'
import { validateCategoryName, validateColor } from '@/utils/validator'
import { useToastStore } from './toast'

export const useCategoriesStore = defineStore('categories', () => {
  const categories = ref<Category[]>([])
  const selectedCategoryId = ref<string | null>(null)
  const toast = useToastStore()

  const loadCategories = () => {
    try {
      categories.value = storage.getCategories()
      if (categories.value.length > 0 && !selectedCategoryId.value) {
        selectedCategoryId.value = categories.value[0].id
      }
      logger.info('加载分类成功', { count: categories.value.length })
    } catch (error) {
      logger.error('加载分类失败', { error })
      toast.error('加载分类失败')
    }
  }

  const categoryMap = computed(() => {
    return new Map(categories.value.map(c => [c.id, c]))
  })

  const selectedCategory = computed(() => {
    if (!selectedCategoryId.value) return null
    return categoryMap.value.get(selectedCategoryId.value) || null
  })

  const getCategoryById = (id: string): Category | undefined => {
    return categoryMap.value.get(id)
  }

  const addCategory = (data: Omit<Category, 'id' | 'createdAt'>): Category | null => {
    const existingNames = categories.value.map(c => c.name)
    const nameValidation = validateCategoryName(data.name, existingNames)
    if (!nameValidation.valid) {
      toast.error(nameValidation.error || '分类名称无效')
      return null
    }

    const colorValidation = validateColor(data.color)
    if (!colorValidation.valid) {
      toast.error(colorValidation.error || '颜色格式无效')
      return null
    }

    const newCategory: Category = {
      ...data,
      id: Date.now().toString(36) + Math.random().toString(36).substr(2, 9),
      createdAt: Date.now()
    }

    const updated = [...categories.value, newCategory]
    if (storage.saveCategories(updated)) {
      categories.value = updated
      if (!selectedCategoryId.value) {
        selectedCategoryId.value = newCategory.id
      }
      logger.info('添加分类', newCategory)
      toast.success(`已添加分类「${data.name}」`)
      return newCategory
    }
    toast.error('添加分类失败')
    return null
  }

  const updateCategory = (id: string, data: Partial<Omit<Category, 'id' | 'createdAt'>>): boolean => {
    const index = categories.value.findIndex(c => c.id === id)
    if (index === -1) {
      toast.error('分类不存在')
      return false
    }

    if (data.name !== undefined) {
      const existingNames = categories.value.filter(c => c.id !== id).map(c => c.name)
      const nameValidation = validateCategoryName(data.name, existingNames, categories.value[index].name)
      if (!nameValidation.valid) {
        toast.error(nameValidation.error || '分类名称无效')
        return false
      }
    }

    if (data.color !== undefined) {
      const colorValidation = validateColor(data.color)
      if (!colorValidation.valid) {
        toast.error(colorValidation.error || '颜色格式无效')
        return false
      }
    }

    const updated = [...categories.value]
    updated[index] = { ...updated[index], ...data }

    if (storage.saveCategories(updated)) {
      categories.value = updated
      logger.info('更新分类', { id, data })
      toast.success('分类已更新')
      return true
    }
    toast.error('更新分类失败')
    return false
  }

  const deleteCategory = (id: string): boolean => {
    if (categories.value.length <= 1) {
      toast.error('至少保留一个分类')
      return false
    }

    const category = categories.value.find(c => c.id === id)
    const updated = categories.value.filter(c => c.id !== id)

    if (storage.saveCategories(updated)) {
      categories.value = updated
      if (selectedCategoryId.value === id) {
        selectedCategoryId.value = updated[0]?.id || null
      }
      logger.info('删除分类', { id, name: category?.name })
      toast.success(`已删除分类「${category?.name}」`)
      return true
    }
    toast.error('删除分类失败')
    return false
  }

  const selectCategory = (id: string) => {
    if (categoryMap.value.has(id)) {
      selectedCategoryId.value = id
      logger.info('选择分类', { id })
    }
  }

  return {
    categories,
    selectedCategoryId,
    selectedCategory,
    categoryMap,
    loadCategories,
    getCategoryById,
    addCategory,
    updateCategory,
    deleteCategory,
    selectCategory
  }
})
