<template>
  <div class="page-container">
    <div class="card">
      <div class="page-header">
        <h2>用户管理</h2>
        <div class="filter-bar">
          <el-select v-model="role" placeholder="全部角色" clearable style="width: 150px">
            <el-option label="买家" value="buyer" />
            <el-option label="卖家" value="seller" />
            <el-option label="鉴定师" value="authenticator" />
          </el-select>
          <el-select v-model="status" placeholder="全部状态" clearable style="width: 150px">
            <el-option label="活跃" value="active" />
            <el-option label="待审核" value="pending" />
            <el-option label="禁用" value="banned" />
          </el-select>
        </div>
      </div>
    </div>

    <div v-if="loading" class="loading-wrapper">
      <el-icon class="is-loading" :size="32"><Loading /></el-icon>
    </div>

    <el-table
      v-else-if="users.length > 0"
      :data="users"
      style="width: 100%"
      stripe
    >
      <el-table-column label="ID" prop="id" width="80" />
      <el-table-column label="用户名" prop="username" width="150" />
      <el-table-column label="邮箱" prop="email" width="200" />
      <el-table-column label="角色" width="100">
        <template #default="{ row }">
          <el-tag>{{ getRoleLabel(row.role) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)">{{ getStatusLabel(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="信用评分" prop="credit_score" width="100" />
      <el-table-column label="注册时间" width="180">
        <template #default="{ row }">
          {{ formatDate(row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="150" fixed="right">
        <template #default="{ row }">
          <el-button
            v-if="row.status === 'active'"
            type="danger"
            size="small"
            @click="handleBan(row)"
          >
            禁用
          </el-button>
          <el-button
            v-if="row.status === 'banned'"
            type="success"
            size="small"
            @click="handleUnban(row)"
          >
            启用
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <div v-else class="empty-state">
      <el-icon :size="64"><Box /></el-icon>
      <p>暂无用户</p>
    </div>

    <div v-if="!loading && total > 0" class="pagination-wrapper">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="loadUsers"
        @size-change="loadUsers"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { userApi } from '@/api/user'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { User } from '@/types'
import dayjs from 'dayjs'
import { Loading, Box } from '@element-plus/icons-vue'

const users = ref<User[]>([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const role = ref('')
const status = ref('')

const loadUsers = async () => {
  loading.value = true
  try {
    const res = await userApi.listUsers({
      page: page.value,
      page_size: pageSize.value,
      role: role.value || undefined,
      status: status.value || undefined
    })
    if (res.code === 200) {
      const data = res.data as any
      users.value = data?.list || []
      total.value = data?.total || 0
    }
  } catch (error) {
    console.error('Load users error:', error)
  } finally {
    loading.value = false
  }
}

const getRoleLabel = (role: string) => {
  const labels: Record<string, string> = {
    buyer: '买家',
    seller: '卖家',
    authenticator: '鉴定师',
    admin: '管理员'
  }
  return labels[role] || role
}

const getStatusType = (status: string) => {
  const types: Record<string, string> = {
    active: 'success',
    pending: 'warning',
    inactive: 'info',
    banned: 'danger'
  }
  return types[status] || 'info'
}

const getStatusLabel = (status: string) => {
  const labels: Record<string, string> = {
    active: '活跃',
    pending: '待审核',
    inactive: '未激活',
    banned: '已禁用'
  }
  return labels[status] || status
}

const formatDate = (date: string) => dayjs(date).format('YYYY-MM-DD HH:mm:ss')

const handleBan = async (user: User) => {
  try {
    await ElMessageBox.confirm(`确认禁用用户 "${user.username}" 吗？`, '禁用确认', {
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      type: 'warning'
    })
    ElMessage.success('已禁用')
    loadUsers()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Ban error:', error)
    }
  }
}

const handleUnban = async (user: User) => {
  try {
    await ElMessageBox.confirm(`确认启用用户 "${user.username}" 吗？`, '启用确认', {
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      type: 'warning'
    })
    ElMessage.success('已启用')
    loadUsers()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Unban error:', error)
    }
  }
}

onMounted(() => {
  loadUsers()
})
</script>

<style lang="scss" scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  
  h2 {
    font-size: 20px;
    font-weight: 600;
  }
}

.filter-bar {
  display: flex;
  gap: 12px;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 24px;
}
</style>
