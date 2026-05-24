<template>
  <div class="page-container">
    <h1 class="page-title">用户管理</h1>

    <div class="card">
      <div v-loading="loading" class="user-management">
        <el-table :data="users" stripe style="width: 100%">
          <el-table-column label="用户ID" prop="id" width="100" />
          <el-table-column label="用户名" min-width="120">
            <template #default="{ row }">
              <div class="user-cell">
                <el-avatar :size="36" :src="row.avatar">
                  {{ row.username?.charAt(0)?.toUpperCase() }}
                </el-avatar>
                <span class="username">{{ row.username }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="邮箱" prop="email" min-width="180" />
          <el-table-column label="真实姓名" prop="realName" min-width="100">
            <template #default="{ row }">
              {{ row.realName || '-' }}
            </template>
          </el-table-column>
          <el-table-column label="角色" width="100">
            <template #default="{ row }">
              <el-tag :type="getRoleTagType(row.role)">
                {{ getRoleText(row.role) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="认证状态" width="100">
            <template #default="{ row }">
              <el-switch
                :model-value="row.verified"
                @change="(val: boolean) => handleVerifyChange(row, val)"
              />
            </template>
          </el-table-column>
          <el-table-column label="注册时间" width="180">
            <template #default="{ row }">
              {{ formatDate(row.createdAt) }}
            </template>
          </el-table-column>
        </el-table>

        <div v-if="!loading && users.length === 0" class="empty-state">
          <el-icon><User /></el-icon>
          <p>暂无用户</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import dayjs from 'dayjs'
import { userApi } from '@/api/equipment'
import type { User } from '@/types'

const loading = ref(false)
const users = ref<User[]>([])

onMounted(() => {
  loadUsers()
})

async function loadUsers() {
  loading.value = true
  try {
    const response = await userApi.getAllUsers()
    users.value = response.data
  } catch (error) {
    console.error('Failed to load users:', error)
    ElMessage.error('加载用户列表失败')
  } finally {
    loading.value = false
  }
}

async function handleVerifyChange(user: User, verified: boolean) {
  try {
    await ElMessageBox.confirm(
      `确定要${verified ? '通过' : '取消'}用户 "${user.username}" 的认证吗？`,
      '认证确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    await userApi.verifyUser(user.id, verified)
    ElMessage.success('操作成功')
    user.verified = verified
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to update verification:', error)
      loadUsers()
    }
  }
}

function getRoleText(role: string) {
  const textMap: Record<string, string> = {
    renter: '租借方',
    owner: '出租方',
    admin: '管理员'
  }
  return textMap[role] || role
}

function getRoleTagType(role: string) {
  const typeMap: Record<string, string> = {
    renter: '',
    owner: 'success',
    admin: 'danger'
  }
  return typeMap[role] || ''
}

function formatDate(date: string) {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}
</script>

<style scoped>
.user-cell {
  display: flex;
  align-items: center;
  gap: 10px;
}

.username {
  font-weight: 500;
}
</style>
