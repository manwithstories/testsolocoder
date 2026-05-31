<template>
  <div class="dashboard-categories">
    <div class="card-shadow p-20">
      <div class="flex-between mb-20">
        <h2 class="page-title">分类管理</h2>
        <el-button type="primary" @click="handleAdd">
          <el-icon><Plus /></el-icon> 新增分类
        </el-button>
      </div>

      <el-table :data="categories" v-loading="loading" border row-key="id">
        <el-table-column prop="name" label="分类名称" />
        <el-table-column prop="parent.name" label="父级分类" width="150">
          <template #default="{ row }">
            {{ row.parent?.name || '顶级分类' }}
          </template>
        </el-table-column>
        <el-table-column prop="sort_order" label="排序" width="100" />
        <el-table-column label="操作" width="180">
          <template #default="{ row }">
            <el-button type="primary" size="small" link @click="handleAddChild(row)">添加子分类</el-button>
            <el-button type="primary" size="small" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" size="small" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-dialog v-model="showDialog" :title="isEdit ? '编辑分类' : '新增分类'" width="500px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="分类名称" prop="name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="父级分类">
          <el-select v-model="form.parent_id" clearable placeholder="顶级分类" style="width: 100%">
            <el-option
              v-for="cat in availableParents"
              :key="cat.id"
              :label="cat.name"
              :value="cat.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sort_order" :min="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import * as collectionApi from '@/api/collection'
import type { CollectionCategory } from '@/types'

const loading = ref(false)
const submitting = ref(false)
const categories = ref<CollectionCategory[]>([])
const showDialog = ref(false)
const isEdit = ref(false)
const editId = ref(0)
const formRef = ref<FormInstance>()

const form = reactive({
  name: '',
  parent_id: undefined as number | undefined,
  sort_order: 0
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入分类名称', trigger: 'blur' }]
}

const availableParents = computed(() => {
  if (isEdit.value) {
    return categories.value.filter(c => c.id !== editId.value)
  }
  return categories.value
})

const fetchList = async () => {
  try {
    loading.value = true
    const res = await collectionApi.listCategories()
    categories.value = res.data
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const handleAdd = () => {
  isEdit.value = false
  editId.value = 0
  Object.assign(form, { name: '', parent_id: undefined, sort_order: 0 })
  showDialog.value = true
}

const handleAddChild = (row: CollectionCategory) => {
  isEdit.value = false
  editId.value = 0
  Object.assign(form, { name: '', parent_id: row.id, sort_order: 0 })
  showDialog.value = true
}

const handleEdit = (row: CollectionCategory) => {
  isEdit.value = true
  editId.value = row.id
  Object.assign(form, {
    name: row.name,
    parent_id: row.parent_id,
    sort_order: row.sort_order
  })
  showDialog.value = true
}

const handleDelete = (row: CollectionCategory) => {
  ElMessageBox.confirm('确定要删除该分类吗？子分类也将受到影响', '提示', {
    type: 'warning'
  }).then(async () => {
    await collectionApi.deleteCategory(row.id)
    ElMessage.success('删除成功')
    fetchList()
  }).catch(() => {})
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate()
  try {
    submitting.value = true
    if (isEdit.value) {
      await collectionApi.updateCategory(editId.value, form)
      ElMessage.success('更新成功')
    } else {
      await collectionApi.createCategory(form)
      ElMessage.success('创建成功')
    }
    showDialog.value = false
    fetchList()
  } catch (e) {
    console.error(e)
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  fetchList()
})
</script>
