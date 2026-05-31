<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElTable, ElTableColumn, ElButton, ElDialog, ElForm, ElFormItem, ElInput, ElMessage, ElPopconfirm, ElCard, ElIcon } from 'element-plus'
import { Plus, Edit, Delete } from '@element-plus/icons-vue'
import { getDepartments, addDepartment, updateDepartment } from '@/api/company'
import type { Department } from '@/types'

const loading = ref(false)
const departments = ref<Department[]>([])
const dialogVisible = ref(false)
const dialogType = ref<'add' | 'edit'>('add')
const currentDepartment = ref<Department | null>(null)

const departmentForm = ref({
  name: '',
  manager_name: '',
  manager_phone: ''
})

const formRules = {
  name: [{ required: true, message: '请输入部门名称', trigger: 'blur' }]
}

const fetchDepartments = async () => {
  loading.value = true
  try {
    departments.value = await getDepartments()
  } catch (error) {
    console.error('Failed to fetch departments:', error)
  } finally {
    loading.value = false
  }
}

const handleAdd = () => {
  dialogType.value = 'add'
  currentDepartment.value = null
  departmentForm.value = {
    name: '',
    manager_name: '',
    manager_phone: ''
  }
  dialogVisible.value = true
}

const handleEdit = (row: Department) => {
  dialogType.value = 'edit'
  currentDepartment.value = row
  departmentForm.value = {
    name: row.name,
    manager_name: row.manager_name || '',
    manager_phone: row.manager_phone || ''
  }
  dialogVisible.value = true
}

const handleSubmit = async () => {
  try {
    if (dialogType.value === 'add') {
      await addDepartment(departmentForm.value)
      ElMessage.success('添加成功')
    } else if (currentDepartment.value) {
      await updateDepartment(currentDepartment.value.id, departmentForm.value)
      ElMessage.success('更新成功')
    }
    dialogVisible.value = false
    fetchDepartments()
  } catch (error) {
    console.error('Failed to submit department:', error)
  }
}

onMounted(() => {
  fetchDepartments()
})
</script>

<template>
  <div class="department-list">
    <ElCard>
      <template #header>
        <div class="card-header">
          <span>部门管理</span>
          <ElButton type="primary" :icon="Plus" @click="handleAdd">
            添加部门
          </ElButton>
        </div>
      </template>

      <ElTable :data="departments" v-loading="loading" border stripe>
        <ElTableColumn prop="name" label="部门名称" />
        <ElTableColumn prop="manager_name" label="负责人" width="120" />
        <ElTableColumn prop="manager_phone" label="联系电话" width="150" />
        <ElTableColumn label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </ElTableColumn>
        <ElTableColumn label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <ElButton type="primary" link :icon="Edit" @click="handleEdit(row)">
              编辑
            </ElButton>
            <ElPopconfirm title="确定删除该部门？">
              <ElButton type="danger" link :icon="Delete">
                删除
              </ElButton>
            </ElPopconfirm>
          </template>
        </ElTableColumn>
      </ElTable>
    </ElCard>

    <ElDialog
      v-model="dialogVisible"
      :title="dialogType === 'add' ? '添加部门' : '编辑部门'"
      width="500px"
    >
      <ElForm :model="departmentForm" :rules="formRules" label-width="100px">
        <ElFormItem label="部门名称" prop="name">
          <ElInput v-model="departmentForm.name" placeholder="请输入部门名称" />
        </ElFormItem>
        <ElFormItem label="负责人">
          <ElInput v-model="departmentForm.manager_name" placeholder="请输入负责人姓名" />
        </ElFormItem>
        <ElFormItem label="联系电话">
          <ElInput v-model="departmentForm.manager_phone" placeholder="请输入联系电话" />
        </ElFormItem>
      </ElForm>
      <template #footer>
        <ElButton @click="dialogVisible = false">取消</ElButton>
        <ElButton type="primary" @click="handleSubmit">确定</ElButton>
      </template>
    </ElDialog>
  </div>
</template>

<style scoped lang="scss">
.department-list {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
}
</style>
