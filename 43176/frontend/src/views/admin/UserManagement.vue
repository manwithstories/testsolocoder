<template>
  <div class="user-management-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <el-icon><User /></el-icon>
          <span>用户管理</span>
        </div>
      </template>

      <div class="filter-bar">
        <el-select v-model="filterRole" placeholder="选择角色" style="width: 140px" @change="fetchUsers">
          <el-option label="全部" value="" />
          <el-option label="发布者" value="publisher" />
          <el-option label="跑腿员" value="courier" />
          <el-option label="管理员" value="admin" />
        </el-select>
        <el-select v-model="filterStatus" placeholder="选择状态" style="width: 140px" @change="fetchUsers">
          <el-option label="全部" value="" />
          <el-option label="正常" value="active" />
          <el-option label="已冻结" value="frozen" />
          <el-option label="已认证" value="verified" />
        </el-select>
      </div>

      <el-table :data="users" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="头像" width="80">
          <template #default="{ row }">
            <el-avatar :size="40" :src="row.avatar">
              {{ row.nickname?.charAt(0) }}
            </el-avatar>
          </template>
        </el-table-column>
        <el-table-column prop="nickname" label="昵称" />
        <el-table-column prop="phone" label="手机号" />
        <el-table-column label="角色" width="100">
          <template #default="{ row }">
            <el-tag :type="getRoleTagType(row.role)" size="small">
              {{ getRoleLabel(row.role) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusTagType(row.status)" size="small">
              {{ getStatusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="balance" label="余额" width="100">
          <template #default="{ row }">
            ¥{{ row.balance?.toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column label="评分" width="120">
          <template #default="{ row }">
            <el-rate :model-value="row.rating || 5" disabled size="small" />
          </template>
        </el-table-column>
        <el-table-column prop="order_count" label="订单数" width="80" />
        <el-table-column label="操作" width="180">
          <template #default="{ row }">
            <el-button
              v-if="row.status !== 'frozen'"
              type="danger"
              size="small"
              @click="handleFreeze(row)"
            >
              冻结
            </el-button>
            <el-button
              v-else
              type="success"
              size="small"
              @click="handleUnfreeze(row)"
            >
              解冻
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          :page-size="20"
          :total="total"
          layout="total, prev, pager, next"
          @current-change="fetchUsers"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { User } from '@element-plus/icons-vue'
import { adminApi } from '@/api'
import type { UserRole, UserStatus } from '@/types'

const loading = ref(false)
const users = ref<any[]>([])
const total = ref(0)
const currentPage = ref(1)
const filterRole = ref('')
const filterStatus = ref('')

const getRoleLabel = (role: UserRole) => {
  const labels: Record<UserRole, string> = {
    publisher: '发布者',
    courier: '跑腿员',
    admin: '管理员'
  }
  return labels[role] || role
}

const getRoleTagType = (role: UserRole): 'primary' | 'success' | 'warning' | 'danger' | 'info' => {
  const types: Record<UserRole, 'primary' | 'success' | 'warning' | 'danger' | 'info'> = {
    publisher: 'primary',
    courier: 'success',
    admin: 'danger'
  }
  return types[role] || 'primary'
}

const getStatusLabel = (status: UserStatus) => {
  const labels: Record<UserStatus, string> = {
    active: '正常',
    frozen: '已冻结',
    verified: '已认证'
  }
  return labels[status] || status
}

const getStatusTagType = (status: UserStatus): 'primary' | 'success' | 'warning' | 'danger' | 'info' => {
  const types: Record<UserStatus, 'primary' | 'success' | 'warning' | 'danger' | 'info'> = {
    active: 'success',
    frozen: 'danger',
    verified: 'warning'
  }
  return types[status] || 'primary'
}

const fetchUsers = async () => {
  loading.value = true
  try {
    const res = await adminApi.listUsers({
      role: filterRole.value || undefined,
      status: filterStatus.value || undefined,
      page: currentPage.value,
      page_size: 20
    })
    if (res.code === 200) {
      users.value = res.data.items
      total.value = res.data.total
    }
  } catch (error) {
    console.error('Failed to fetch users:', error)
  } finally {
    loading.value = false
  }
}

const handleFreeze = async (user: any) => {
  try {
    await ElMessageBox.prompt(`确定要冻结用户 ${user.nickname} 吗？请输入冻结原因`, '冻结确认', {
      confirmButtonText: '确定冻结',
      cancelButtonText: '取消',
      inputPattern: /.+/,
      inputErrorMessage: '请输入冻结原因'
    }).then(async ({ value }) => {
      const res = await adminApi.freezeUser(user.id, { reason: value })
      if (res.code === 200) {
        ElMessage.success('已冻结用户')
        fetchUsers()
      }
    })
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to freeze user:', error)
    }
  }
}

const handleUnfreeze = async (user: any) => {
  try {
    await ElMessageBox.confirm(`确定要解冻用户 ${user.nickname} 吗？`, '解冻确认', {
      type: 'warning'
    })
    const res = await adminApi.unfreezeUser(user.id)
    if (res.code === 200) {
      ElMessage.success('已解冻用户')
      fetchUsers()
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to unfreeze user:', error)
    }
  }
}

onMounted(() => {
  fetchUsers()
})
</script>

<style lang="scss" scoped>
.user-management-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;

  .card-header {
    display: flex;
    align-items: center;
    gap: 8px;
    font-weight: 600;
  }

  .filter-bar {
    display: flex;
    gap: 12px;
    margin-bottom: 20px;
  }

  .pagination {
    display: flex;
    justify-content: center;
    margin-top: 20px;
  }
}
</style>
