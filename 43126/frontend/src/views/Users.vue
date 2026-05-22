<template>
  <div class="users">
    <el-card class="filter-card">
      <el-form :inline="true" :model="filterForm">
        <el-form-item label="角色">
          <el-select v-model="filterForm.role" placeholder="全部角色" clearable @change="loadUsers">
            <el-option label="管理员" value="admin" />
            <el-option label="空间管理员" value="space_admin" />
            <el-option label="普通用户" value="user" />
          </el-select>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card style="margin-top: 20px; border-radius: 8px">
      <el-table :data="users" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column prop="email" label="邮箱" width="200" />
        <el-table-column prop="real_name" label="真实姓名" width="120" />
        <el-table-column prop="department" label="部门" width="120" />
        <el-table-column prop="phone" label="电话" width="140" />
        <el-table-column label="角色" width="120">
          <template #default="{ row }">
            <el-tag :type="getRoleType(row.role)">{{ getRoleText(row.role) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">
              {{ row.status === 1 ? '正常' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-dropdown @command="(cmd: string) => handleRoleChange(row, cmd)">
              <el-button type="primary" size="small">
                修改角色<el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="admin">管理员</el-dropdown-item>
                  <el-dropdown-item command="space_admin">空间管理员</el-dropdown-item>
                  <el-dropdown-item command="user">普通用户</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
            <el-button type="danger" size="small" @click="deleteUser(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        v-if="total > 0"
        style="margin-top: 16px; justify-content: flex-end"
        layout="total, prev, pager, next"
        :total="total"
        :current-page="filterForm.page"
        :page-size="filterForm.page_size"
        @current-change="handlePageChange"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { api } from '@/api'

const loading = ref(false)
const users = ref<any[]>([])
const total = ref(0)

const filterForm = reactive({
  page: 1,
  page_size: 10,
  role: ''
})

onMounted(() => {
  loadUsers()
})

async function loadUsers() {
  loading.value = true
  try {
    const res: any = await api.listUsers(filterForm)
    users.value = res.data?.users || []
    total.value = res.data?.total || 0
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function handlePageChange(page: number) {
  filterForm.page = page
  loadUsers()
}

async function handleRoleChange(row: any, role: string) {
  try {
    await ElMessageBox.confirm(`确定将用户 "${row.username}" 角色修改为 "${getRoleText(role)}" 吗？`, '提示', {
      type: 'warning'
    })
    await api.updateUserRole(row.id, role)
    ElMessage.success('角色修改成功')
    loadUsers()
  } catch (e: any) {
    if (e !== 'cancel') console.error(e)
  }
}

async function deleteUser(row: any) {
  try {
    await ElMessageBox.confirm(`确定删除用户 "${row.username}" 吗？此操作不可恢复！`, '警告', {
      type: 'error'
    })
    await api.deleteUser(row.id)
    ElMessage.success('删除成功')
    loadUsers()
  } catch (e: any) {
    if (e !== 'cancel') console.error(e)
  }
}

function getRoleType(role: string) {
  const map: Record<string, string> = { admin: 'danger', space_admin: 'warning', user: '' }
  return map[role] || ''
}

function getRoleText(role: string) {
  const map: Record<string, string> = { admin: '管理员', space_admin: '空间管理员', user: '普通用户' }
  return map[role] || role
}
</script>

<style scoped>
.filter-card {
  border-radius: 8px;
}
</style>
