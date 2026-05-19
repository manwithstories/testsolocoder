<script setup lang="ts">
import { ref } from 'vue'
import { Plus } from 'lucide-vue-next'
import { useCategoriesStore } from '@/stores/categories'
import type { Category } from '@/types'
import CategoryCard from '@/components/CategoryCard.vue'
import CategoryForm from '@/components/CategoryForm.vue'

const categoriesStore = useCategoriesStore()

const showForm = ref(false)
const editingCategory = ref<Category | null>(null)

const openCreateForm = () => {
  editingCategory.value = null
  showForm.value = true
}

const openEditForm = (category: Category) => {
  editingCategory.value = category
  showForm.value = true
}

const handleSubmit = (data: { name: string; color: string; icon: string }) => {
  if (editingCategory.value) {
    categoriesStore.updateCategory(editingCategory.value.id, data)
  } else {
    categoriesStore.addCategory(data)
  }
  showForm.value = false
}

const handleDelete = (id: string) => {
  if (confirm('确定要删除这个分类吗？相关的记录不会被删除。')) {
    categoriesStore.deleteCategory(id)
  }
}
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl lg:text-3xl font-bold text-gray-900 dark:text-white mb-2">
          分类管理
        </h1>
        <p class="text-gray-500 dark:text-gray-400">
          管理你的专注分类
        </p>
      </div>
      <button
        @click="openCreateForm"
        class="flex items-center gap-2 px-4 py-2.5 bg-indigo-500 hover:bg-indigo-600 text-white rounded-xl font-medium transition-colors shadow-lg shadow-indigo-500/25"
      >
        <Plus class="w-5 h-5" />
        <span class="hidden sm:inline">新建分类</span>
      </button>
    </div>

    <div v-if="categoriesStore.categories.length === 0" class="text-center py-16 bg-white dark:bg-gray-800 rounded-2xl shadow-lg">
      <p class="text-gray-500 dark:text-gray-400">暂无分类</p>
      <button
        @click="openCreateForm"
        class="mt-4 text-indigo-500 hover:text-indigo-600 font-medium"
      >
        创建第一个分类
      </button>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <CategoryCard
        v-for="category in categoriesStore.categories"
        :key="category.id"
        :category="category"
        @edit="openEditForm"
        @delete="handleDelete"
      />
    </div>

    <CategoryForm
      :visible="showForm"
      :category="editingCategory"
      @submit="handleSubmit"
      @cancel="showForm = false"
    />
  </div>
</template>
