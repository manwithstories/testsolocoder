<template>
  <div class="page-container">
    <div class="page-header">
      <h1 class="page-title">用户管理</h1>
    </div>

    <div class="card">
      <div class="card-body">
        <el-form :inline="true" :model="filterForm">
          <el-form-item label="关键词">
            <el-input v-model="filterForm.keyword" placeholder="搜索邮箱/昵称" clearable @input="loadUsers" />
          </el-form-item>
          <el-form-item label="状态">
            <el-select v-model="filterForm.status" placeholder="全部" clearable @change="loadUsers">
              <el-option label="正常" :value="1" />
              <el-option label="禁用" :value="2" />
            </el-select>
          </el-form-item>
        </el-form>
      </div>
    </div>

    <div class="card">
      <el-table :data="users" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="头像" width="80">
          <template #default="{ row }">
            <el-avatar :size="32" :src="row.avatar">
              {{ row.nickname?.charAt(0) || 'U' }}
            </el-avatar>
          </template>
        </el-table-column>
        <el-table-column prop="email" label="邮箱" />
        <el-table-column prop="nickname" label="昵称" />
        <el-table-column label="角色" width="120">
          <template #default="{ row }">
            <el-tag :type="getRoleTagType(row.role)">{{ getRoleText(row.role) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">
              {{ row.status === 1 ? '正常' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="注册时间" width="180" />
        <el-table-column label="操作" width="240">
          <template #default="{ row }">
            <el-button size="small" @click="handleUpdateRole(row)">
              分配角色
            </el-button>
            <el-button
              size="small"
              :type="row.status === 1 ? 'warning' : 'success'"
              @click="handleToggleStatus(row)"
            >
              {{ row.status === 1 ? '禁用' : '启用' }}
            </el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="card-body">
        <el-pagination
          v-model:current-page="filterForm.page"
          v-model:page-size="filterForm.page_size"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="loadUsers"
          @current-change="loadUsers"
        />
      </div>
    </div>

    <el-dialog v-model="showRoleDialog" title="分配角色" width="400px">
      <el-form label-width="80px">
        <el-form-item label="用户">
          <span>{{ selectedUser?.email }}</span>
        </el-form-item>
        <el-form-item label="角色">
          <el-select v-model="selectedRole" style="width: 100%;">
            <el-option label="管理员" :value="1" />
            <el-option label="编辑员" :value="2" />
            <el-option label="查看员" :value="3" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showRoleDialog = false">取消</el-button>
        <el-button type="primary" @click="confirmUpdateRole">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { userApi } from '@/api/statistics'
import type { User } from '@/types'

const loading = ref(false)
const users = ref<User[]>([])
const total = ref(0)
const showRoleDialog = ref(false)
const selectedUser = ref<User>()
const selectedRole = ref(1)

const filterForm = reactive({
  page: 1,
  page_size: 10,
  keyword: '',
  status: undefined as number | undefined
})

const getRoleText = (role: string) => {
  const map: Record<string, string> = {
    admin: '管理员',
    editor: '编辑员',
    viewer: '查看员'
  }
  return map[role] || role
}

const getRoleTagType = (role: string) => {
  const map: Record<string, string> = {
    admin: 'danger',
    editor: 'warning',
    viewer: 'info'
  }
  return map[role] || ''
}

const loadUsers = async () => {
  loading.value = true
  try {
    const res = await userApi.list(filterForm) as any
    users.value = res.items
    total.value = res.total
  } catch (error) {
    console.error('Failed to load users')
  } finally {
    loading.value = false
  }
}

const handleUpdateRole = (user: User) => {
  selectedUser.value = user
  const roleMap: Record<string, number> = { admin: 1, editor: 2, viewer: 3 }
  selectedRole.value = roleMap[user.role] || 3
  showRoleDialog.value = true
}

const confirmUpdateRole = async () => {
  if (!selectedUser.value) return

  try {
    await userApi.updateRole(selectedUser.value.id, selectedRole.value)
    ElMessage.success('角色已更新')
    showRoleDialog.value = false
    loadUsers()
  } catch (e: any) {
    ElMessage.error(e.message || '更新失败')
  }
}

const handleToggleStatus = async (user: User) => {
  const newStatus = user.status === 1 ? 2 : 1
  try {
    await userApi.updateStatus(user.id, newStatus)
    ElMessage.success(`已${newStatus === 1 ? '启用' : '禁用'}用户`)
    loadUsers()
  } catch (e: any) {
    ElMessage.error(e.message || '操作失败')
  }
}

const handleDelete = async (user: User) => {
  try {
    await ElMessageBox.confirm(`确定要删除用户 ${user.email} 吗？`, '确认删除', {
      type: 'warning'
    })
    await userApi.remove(user.id)
    ElMessage.success('用户已删除')
    loadUsers()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '删除失败')
    }
  }
}

onMounted(loadUsers)
</script>
