<template>
  <div class="admin-service-items-page">
    <div class="page-header">
      <h2 class="page-title">服务项目管理</h2>
      <el-button type="primary" @click="showCreateDialog = true">
        新增服务项目
      </el-button>
    </div>

    <el-card>
      <div class="filter-bar">
        <el-select v-model="categoryFilter" placeholder="选择分类" clearable style="width: 150px;" @change="loadServiceItems">
          <el-option
            v-for="category in categories"
            :key="category.id"
            :label="category.name"
            :value="category.id"
          />
        </el-select>
      </div>

      <el-table :data="serviceItems" style="width: 100%">
        <el-table-column prop="name" label="名称" width="150" />
        <el-table-column label="分类" width="120">
          <template #default="{ row }">{{ row.category?.name }}</template>
        </el-table-column>
        <el-table-column prop="description" label="描述" min-width="200" />
        <el-table-column label="价格区间" width="150">
          <template #default="{ row }">¥{{ row.min_price }} - ¥{{ row.max_price }}</template>
        </el-table-column>
        <el-table-column prop="estimated_time" label="预计时长" width="120">
          <template #default="{ row }">{{ row.estimated_time }}分钟</template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status ? 'success' : 'danger'">
              {{ row.status ? '正常' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150">
          <template #default="{ row }">
            <el-button size="small" @click="editItem(row)">编辑</el-button>
            <el-button size="small" type="danger" @click="deleteItem(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="showCreateDialog" :title="editingItem ? '编辑服务项目' : '新增服务项目'" width="500px">
      <el-form :model="itemForm" label-width="100px">
        <el-form-item label="所属分类">
          <el-select v-model="itemForm.category_id" placeholder="选择分类">
            <el-option
              v-for="category in categories"
              :key="category.id"
              :label="category.name"
              :value="category.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="名称">
          <el-input v-model="itemForm.name" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="itemForm.description" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item label="最低价格">
          <el-input-number v-model="itemForm.min_price" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="最高价格">
          <el-input-number v-model="itemForm.max_price" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="预计时长(分钟)">
          <el-input-number v-model="itemForm.estimated_time" :min="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" @click="saveItem" :loading="saving">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { categoryApi, serviceItemApi } from '@/api/category'
import type { Category, ServiceItem } from '@/types'

const categories = ref<Category[]>([])
const serviceItems = ref<ServiceItem[]>([])
const categoryFilter = ref<number | ''>('')
const showCreateDialog = ref(false)
const editingItem = ref<ServiceItem | null>(null)
const saving = ref(false)

const itemForm = reactive({
  category_id: null as number | null,
  name: '',
  description: '',
  min_price: 0,
  max_price: 0,
  estimated_time: 0
})

onMounted(() => {
  loadCategories()
  loadServiceItems()
})

async function loadCategories() {
  try {
    const res = await categoryApi.getCategories()
    categories.value = res.data || []
  } catch (error) {
    console.error('Failed to load categories:', error)
  }
}

async function loadServiceItems() {
  try {
    const res = await serviceItemApi.getServiceItems(
      categoryFilter.value ? { category_id: Number(categoryFilter.value) } : undefined
    )
    serviceItems.value = res.data || []
  } catch (error) {
    console.error('Failed to load service items:', error)
  }
}

function editItem(item: ServiceItem) {
  editingItem.value = item
  itemForm.category_id = item.category_id
  itemForm.name = item.name
  itemForm.description = item.description
  itemForm.min_price = item.min_price
  itemForm.max_price = item.max_price
  itemForm.estimated_time = item.estimated_time
  showCreateDialog.value = true
}

async function saveItem() {
  if (!itemForm.category_id || !itemForm.name) {
    ElMessage.error('请填写分类和名称')
    return
  }

  saving.value = true
  try {
    if (editingItem.value) {
      await serviceItemApi.updateServiceItem(editingItem.value.id, itemForm as any)
    } else {
      await serviceItemApi.createServiceItem(itemForm as any)
    }
    ElMessage.success('保存成功')
    showCreateDialog.value = false
    resetForm()
    loadServiceItems()
  } catch (error) {
    console.error('Failed to save item:', error)
  } finally {
    saving.value = false
  }
}

function resetForm() {
  editingItem.value = null
  itemForm.category_id = null
  itemForm.name = ''
  itemForm.description = ''
  itemForm.min_price = 0
  itemForm.max_price = 0
  itemForm.estimated_time = 0
}

async function deleteItem(item: ServiceItem) {
  try {
    await ElMessageBox.confirm(`确定要删除服务项目"${item.name}"吗？`, '确认', { type: 'warning' })
    await serviceItemApi.deleteServiceItem(item.id)
    ElMessage.success('删除成功')
    loadServiceItems()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('Failed to delete item:', error)
    }
  }
}
</script>

<style scoped>
.admin-service-items-page {
  padding: 0;
}

.filter-bar {
  display: flex;
  gap: 15px;
  margin-bottom: 20px;
}
</style>
