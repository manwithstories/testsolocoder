<template>
  <div class="user-list">
    <div class="page-header">
      <el-input
        v-model="searchQuery"
        placeholder="搜索用户..."
        :prefix-icon="Search"
        clearable
        style="width: 200px"
        @clear="fetchUsers"
        @keyup.enter="fetchUsers"
      />
      <el-select v-model="roleFilter" placeholder="角色" clearable style="width: 140px" @change="fetchUsers">
        <el-option label="准新人" value="couple" />
        <el-option label="婚礼策划师" value="planner" />
        <el-option label="婚庆服务商" value="vendor" />
        <el-option label="管理员" value="admin" />
      </el-select>
    </div>

    <el-table :data="users" v-loading="loading" stripe>
      <el-table-column label="用户信息" min-width="200">
        <template #default="{ row }">
          <div class="user-cell">
            <el-avatar :size="36" :src="row.avatar">
              {{ row.full_name?.charAt(0) || row.username.charAt(0) }}
            </el-avatar>
            <div class="user-info">
              <div class="username">{{ row.username }}</div>
              <div class="email">{{ row.email }}</div>
            </div>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="full_name" label="姓名" width="120" />
      <el-table-column label="角色" width="120">
        <template #default="{ row }">
          <el-tag :type="roleType(row.role)" size="small">{{ roleText(row.role) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 'active' ? 'success' : 'danger'" size="small">
            {{ row.status === 'active' ? '正常' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="注册时间" width="180">
        <template #default="{ row }">
          {{ formatDate(row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="150" fixed="right">
        <template #default="{ row }">
          <el-button
            v-if="row.status === 'active'"
            type="warning"
            link
            @click="toggleStatus(row)"
          >
            禁用
          </el-button>
          <el-button
            v-else
            type="success"
            link
            @click="toggleStatus(row)"
          >
            启用
          </el-button>
          <el-button type="danger" link @click="deleteUser(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      v-model:current-page="page"
      v-model:page-size="pageSize"
      :page-sizes="[10, 20, 50]"
      :total="total"
      layout="total, sizes, prev, pager, next, jumper"
      @size-change="fetchUsers"
      @current-change="fetchUsers"
      class="pagination"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { userApi } from '@/api/auth'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import type { User } from '@/types'

const loading = ref(false)
const users = ref<User[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const searchQuery = ref('')
const roleFilter = ref('')

function roleText(role: string) {
  const texts: Record<string, string> = {
    couple: '准新人',
    planner: '婚礼策划师',
    vendor: '婚庆服务商',
    admin: '管理员'
  }
  return texts[role] || role
}

function roleType(role: string) {
  const types: Record<string, string> = {
    couple: 'primary',
    planner: 'success',
    vendor: 'warning',
    admin: 'danger'
  }
  return types[role] || 'info'
}

function formatDate(date: string) {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

async function fetchUsers() {
  loading.value = true
  try {
    const res = await userApi.getUsers({
      search: searchQuery.value,
      role: roleFilter.value,
      page: page.value,
      page_size: pageSize.value
    })
    users.value = res.data.list
    total.value = res.data.total
  } catch (error) {
    console.error('Failed to fetch users:', error)
  } finally {
    loading.value = false
  }
}

async function toggleStatus(user: User) {
  const newStatus = user.status === 'active' ? 'disabled' : 'active'
  try {
    await ElMessageBox.confirm(`确定要${newStatus === 'active' ? '启用' : '禁用'}用户"${user.username}"吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await userApi.updateStatus(user.id, newStatus)
    ElMessage.success('状态更新成功')
    fetchUsers()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to update status:', error)
    }
  }
}

async function deleteUser(user: User) {
  try {
    await ElMessageBox.confirm(`确定要删除用户"${user.username}"吗？此操作不可恢复！`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await userApi.deleteUser(user.id)
    ElMessage.success('删除成功')
    fetchUsers()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete user:', error)
    }
  }
}

onMounted(fetchUsers)
</script>

<style scoped>
.user-list {
  padding: 0;
}

.page-header {
  display: flex;
  gap: 16px;
  margin-bottom: 20px;
}

.user-cell {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-info {
  display: flex;
  flex-direction: column;
}

.username {
  font-weight: 600;
}

.email {
  color: #909399;
  font-size: 12px;
}

.pagination {
  margin-top: 20px;
  justify-content: flex-end;
}
</style>
