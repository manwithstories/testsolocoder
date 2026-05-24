<template>
  <div class="admin-categories-page">
    <div class="page-header">
      <h2 class="page-title">分类管理</h2>
      <el-button type="primary" @click="showCreateDialog = true">
        新增分类
      </el-button>
    </div>

    <el-card>
      <el-table :data="categories" style="width: 100%">
        <el-table-column prop="icon" label="图标" width="100">
          <template #default="{ row }">
            <span style="font-size: 24px;">{{ row.icon }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="name" label="名称" width="150" />
        <el-table-column prop="code" label="编码" width="150" />
        <el-table-column prop="sort" label="排序" width="100" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status ? 'success' : 'danger'">
              {{ row.status ? '正常' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150">
          <template #default="{ row }">
            <el-button size="small" @click="editCategory(row)">编辑</el-button>
            <el-button size="small" type="danger" @click="deleteCategory(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="showCreateDialog" :title="editingCategory ? '编辑分类' : '新增分类'" width="400px">
      <el-form :model="categoryForm" label-width="80px">
        <el-form-item label="名称">
          <el-input v-model="categoryForm.name" />
        </el-form-item>
        <el-form-item label="编码">
          <el-input v-model="categoryForm.code" />
        </el-form-item>
        <el-form-item label="图标">
          <el-input v-model="categoryForm.icon" placeholder="输入emoji或图标URL" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="categoryForm.sort" :min="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" @click="saveCategory" :loading="saving">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { categoryApi } from '@/api/category'
import type { Category } from '@/types'

const categories = ref<Category[]>([])
const showCreateDialog = ref(false)
const editingCategory = ref<Category | null>(null)
const saving = ref(false)

const categoryForm = reactive({
  name: '',
  code: '',
  icon: '',
  sort: 0
})

onMounted(() => {
  loadCategories()
})

async function loadCategories() {
  try {
    const res = await categoryApi.getCategories()
    categories.value = res.data || []
  } catch (error) {
    console.error('Failed to load categories:', error)
  }
}

function editCategory(category: Category) {
  editingCategory.value = category
  categoryForm.name = category.name
  categoryForm.code = category.code
  categoryForm.icon = category.icon
  categoryForm.sort = category.sort
  showCreateDialog.value = true
}

async function saveCategory() {
  if (!categoryForm.name || !categoryForm.code) {
    ElMessage.error('请填写名称和编码')
    return
  }

  saving.value = true
  try {
    if (editingCategory.value) {
      await categoryApi.updateCategory(editingCategory.value.id, categoryForm)
    } else {
      await categoryApi.createCategory(categoryForm)
    }
    ElMessage.success('保存成功')
    showCreateDialog.value = false
    resetForm()
    loadCategories()
  } catch (error) {
    console.error('Failed to save category:', error)
  } finally {
    saving.value = false
  }
}

function resetForm() {
  editingCategory.value = null
  categoryForm.name = ''
  categoryForm.code = ''
  categoryForm.icon = ''
  categoryForm.sort = 0
}

async function deleteCategory(category: Category) {
  try {
    await ElMessageBox.confirm(`确定要删除分类"${category.name}"吗？`, '确认', { type: 'warning' })
    await categoryApi.deleteCategory(category.id)
    ElMessage.success('删除成功')
    loadCategories()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('Failed to delete category:', error)
    }
  }
}
</script>

<style scoped>
.admin-categories-page {
  padding: 0;
}
</style>
