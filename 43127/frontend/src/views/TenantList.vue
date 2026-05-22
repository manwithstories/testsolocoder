<template>
  <div class="tenant-list">
    <div class="page-header">
      <h2 class="page-title">租户管理</h2>
      <el-button type="primary" @click="showCreateDialog">
        <el-icon><Plus /></el-icon>添加租户
      </el-button>
    </div>

    <div class="search-bar card">
      <el-input
        v-model="keyword"
        placeholder="搜索租户姓名/手机号"
        clearable
        style="width: 200px"
        @keyup.enter="loadData"
      />
      <el-button type="primary" @click="loadData">搜索</el-button>
    </div>

    <div class="card">
      <el-table :data="tenants" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="姓名" width="120" />
        <el-table-column prop="phone" label="手机号" width="130" />
        <el-table-column prop="idCard" label="身份证号" width="200" />
        <el-table-column prop="email" label="邮箱" width="180" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">
              {{ row.status === 1 ? '正常' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="showEditDialog(row)">编辑</el-button>
            <el-button link type="danger" @click="deleteTenant(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        class="pagination"
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next, jumper"
        @current-change="loadData"
        @size-change="loadData"
      />
    </div>

    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑租户' : '添加租户'"
      width="500px"
    >
      <el-form :model="tenantForm" :rules="tenantRules" ref="tenantFormRef" label-width="80px">
        <el-form-item label="姓名" prop="name">
          <el-input v-model="tenantForm.name" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="tenantForm.phone" />
        </el-form-item>
        <el-form-item label="身份证号" prop="idCard">
          <el-input v-model="tenantForm.idCard" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="tenantForm.email" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitTenant">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import type { Tenant } from '@/types'
import { getTenants, createTenant, updateTenant, deleteTenant as deleteTenantApi } from '@/api/tenant'

const loading = ref(false)
const tenants = ref<Tenant[]>([])
const keyword = ref('')
const dialogVisible = ref(false)
const isEdit = ref(false)
const editId = ref<number | null>(null)
const tenantFormRef = ref<FormInstance>()

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const tenantForm = reactive({
  name: '',
  phone: '',
  idCard: '',
  email: ''
})

const tenantRules: FormRules = {
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  phone: [{ required: true, message: '请输入手机号', trigger: 'blur' }]
}

onMounted(() => {
  loadData()
})

async function loadData() {
  loading.value = true
  try {
    const res = await getTenants({
      page: pagination.page,
      pageSize: pagination.pageSize,
      keyword: keyword.value
    })
    tenants.value = res.data.list
    pagination.total = res.data.total
  } catch (error) {
    console.error('Failed to load tenants:', error)
  } finally {
    loading.value = false
  }
}

function showCreateDialog() {
  isEdit.value = false
  editId.value = null
  tenantForm.name = ''
  tenantForm.phone = ''
  tenantForm.idCard = ''
  tenantForm.email = ''
  dialogVisible.value = true
}

function showEditDialog(row: Tenant) {
  isEdit.value = true
  editId.value = row.id
  tenantForm.name = row.name
  tenantForm.phone = row.phone
  tenantForm.idCard = row.idCard
  tenantForm.email = row.email
  dialogVisible.value = true
}

async function submitTenant() {
  if (!tenantFormRef.value) return
  
  await tenantFormRef.value.validate(async (valid) => {
    if (valid) {
      try {
        if (isEdit.value && editId.value) {
          await updateTenant(editId.value, tenantForm)
          ElMessage.success('更新成功')
        } else {
          await createTenant(tenantForm)
          ElMessage.success('创建成功')
        }
        dialogVisible.value = false
        loadData()
      } catch (error) {
        console.error('Submit failed:', error)
      }
    }
  })
}

async function deleteTenant(row: Tenant) {
  try {
    await ElMessageBox.confirm('确定要删除该租户吗？', '提示', {
      type: 'warning'
    })
    await deleteTenantApi(row.id)
    ElMessage.success('删除成功')
    loadData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error(error)
    }
  }
}
</script>

<style scoped>
.tenant-list {
  padding: 0;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
