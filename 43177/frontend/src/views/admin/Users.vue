<template>
  <div class="admin-users-page">
    <div class="page-header">
      <h2 class="page-title">用户管理</h2>
    </div>

    <el-card>
      <div class="filter-bar">
        <el-select v-model="roleFilter" placeholder="选择角色" clearable style="width: 150px;" @change="loadUsers">
          <el-option label="客户" value="customer" />
          <el-option label="技师" value="technician" />
          <el-option label="管理员" value="admin" />
        </el-select>
        <el-select v-model="statusFilter" placeholder="选择状态" clearable style="width: 150px;" @change="loadUsers">
          <el-option label="正常" value="active" />
          <el-option label="禁用" value="disabled" />
          <el-option label="待审核" value="pending" />
        </el-select>
      </div>

      <el-table :data="users" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户名" width="150" />
        <el-table-column prop="real_name" label="真实姓名" width="150" />
        <el-table-column prop="phone" label="手机号" width="150" />
        <el-table-column prop="email" label="邮箱" width="200" />
        <el-table-column prop="role" label="角色" width="100">
          <template #default="{ row }">
            <el-tag :type="getRoleType(row.role)">{{ getRoleText(row.role) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="balance" label="余额" width="120">
          <template #default="{ row }">¥{{ row.balance }}</template>
        </el-table-column>
        <el-table-column prop="created_at" label="注册时间" width="180">
          <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="150">
          <template #default="{ row }">
            <el-button size="small" @click="viewUser(row)">查看</el-button>
            <el-button
              size="small"
              :type="row.status === 'active' ? 'danger' : 'success'"
              @click="toggleStatus(row)"
            >
              {{ row.status === 'active' ? '禁用' : '启用' }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div v-if="total > pageSize" class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          :page-size="pageSize"
          :total="total"
          layout="prev, pager, next"
          @current-change="loadUsers"
        />
      </div>
    </el-card>

    <el-dialog v-model="showDetailDialog" title="用户详情" width="600px">
      <el-descriptions :column="2" border v-if="selectedUser">
        <el-descriptions-item label="ID">{{ selectedUser.id }}</el-descriptions-item>
        <el-descriptions-item label="用户名">{{ selectedUser.username }}</el-descriptions-item>
        <el-descriptions-item label="真实姓名">{{ selectedUser.real_name || '-' }}</el-descriptions-item>
        <el-descriptions-item label="手机号">{{ selectedUser.phone }}</el-descriptions-item>
        <el-descriptions-item label="邮箱">{{ selectedUser.email || '-' }}</el-descriptions-item>
        <el-descriptions-item label="角色">
          <el-tag :type="getRoleType(selectedUser.role)">{{ getRoleText(selectedUser.role) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(selectedUser.status)">{{ getStatusText(selectedUser.status) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="余额">¥{{ selectedUser.balance }}</el-descriptions-item>
        <el-descriptions-item label="地址">{{ selectedUser.address || '-' }}</el-descriptions-item>
        <el-descriptions-item label="注册时间">{{ formatTime(selectedUser.created_at) }}</el-descriptions-item>
      </el-descriptions>

      <el-card v-if="selectedUser?.technician_profile" class="mt-20">
        <template #header>技师信息</template>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="专业">
            {{ selectedUser.technician_profile.specialty || '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="评分">
            {{ selectedUser.technician_profile.rating }}
          </el-descriptions-item>
          <el-descriptions-item label="完成订单">
            {{ selectedUser.technician_profile.completed_orders }}
          </el-descriptions-item>
          <el-descriptions-item label="认证状态">
            <el-tag :type="selectedUser.technician_profile.is_verified ? 'success' : 'warning'">
              {{ selectedUser.technician_profile.is_verified ? '已认证' : '未认证' }}
            </el-tag>
          </el-descriptions-item>
        </el-descriptions>
      </el-card>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { adminApi } from '@/api/admin'
import type { User } from '@/types'
import dayjs from 'dayjs'

const users = ref<User[]>([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const roleFilter = ref('')
const statusFilter = ref('')
const showDetailDialog = ref(false)
const selectedUser = ref<any>(null)

onMounted(() => {
  loadUsers()
})

async function loadUsers() {
  try {
    const res = await adminApi.getUsers({
      page: currentPage.value,
      page_size: pageSize.value,
      role: roleFilter.value || undefined,
      status: statusFilter.value || undefined
    })
    users.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch (error) {
    console.error('Failed to load users:', error)
  }
}

function getRoleType(role: string): string {
  const typeMap: Record<string, string> = {
    customer: 'info',
    technician: 'success',
    admin: 'warning'
  }
  return typeMap[role] || 'info'
}

function getRoleText(role: string): string {
  const textMap: Record<string, string> = {
    customer: '客户',
    technician: '技师',
    admin: '管理员'
  }
  return textMap[role] || role
}

function getStatusType(status: string): string {
  const typeMap: Record<string, string> = {
    active: 'success',
    disabled: 'danger',
    pending: 'warning'
  }
  return typeMap[status] || 'info'
}

function getStatusText(status: string): string {
  const textMap: Record<string, string> = {
    active: '正常',
    disabled: '禁用',
    pending: '待审核'
  }
  return textMap[status] || status
}

function formatTime(time: string): string {
  return dayjs(time).format('YYYY-MM-DD HH:mm')
}

async function viewUser(user: User) {
  try {
    const res = await adminApi.getUserDetail(user.id)
    selectedUser.value = res.data
    showDetailDialog.value = true
  } catch (error) {
    console.error('Failed to load user detail:', error)
  }
}

async function toggleStatus(user: User) {
  const newStatus = user.status === 'active' ? 'disabled' : 'active'
  try {
    await ElMessageBox.confirm(
      `确定要${newStatus === 'active' ? '启用' : '禁用'}该用户吗？`,
      '确认',
      { type: 'warning' }
    )
    await adminApi.updateUserStatus(user.id, { status: newStatus })
    ElMessage.success('操作成功')
    loadUsers()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('Failed to update status:', error)
    }
  }
}
</script>

<style scoped>
.admin-users-page {
  padding: 0;
}

.filter-bar {
  display: flex;
  gap: 15px;
  margin-bottom: 20px;
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

.mt-20 {
  margin-top: 20px;
}
</style>
