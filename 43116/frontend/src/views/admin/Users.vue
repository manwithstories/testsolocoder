<template>
  <div>
    <div class="page-header">
      <h1 class="page-title">用户管理</h1>
    </div>

    <div class="search-bar">
      <el-input
        v-model="filters.keyword"
        placeholder="搜索用户名/邮箱/姓名"
        clearable
        style="width: 250px"
        @keyup.enter="loadUsers"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
      <el-button type="primary" @click="loadUsers">搜索</el-button>
    </div>

    <el-table :data="users" v-loading="loading" style="width: 100%">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="username" label="用户名" width="150" />
      <el-table-column prop="email" label="邮箱" width="200" />
      <el-table-column prop="phone" label="手机号" width="130" />
      <el-table-column prop="real_name" label="真实姓名" width="120" />
      <el-table-column label="认证状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getAuthStatusType(row.auth_status)">
            {{ getAuthStatusText(row.auth_status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="账户状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 'active' ? 'success' : 'danger'">
            {{ row.status === 'active' ? '正常' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="注册时间" width="170">
        <template #default="{ row }">
          {{ formatDate(row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button
            v-if="row.auth_status === 'pending'"
            type="success"
            size="small"
            @click="handleAuth(row.id, 'active')"
          >
            通过认证
          </el-button>
          <el-button
            v-if="row.auth_status === 'pending'"
            type="danger"
            size="small"
            @click="handleAuth(row.id, 'rejected')"
          >
            拒绝
          </el-button>
          <el-button
            v-if="row.status === 'active'"
            type="warning"
            size="small"
            @click="handleStatus(row.id, 'disabled')"
          >
            禁用
          </el-button>
          <el-button
            v-else
            type="success"
            size="small"
            @click="handleStatus(row.id, 'active')"
          >
            启用
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <div style="text-align: center; margin-top: 20px;">
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next, jumper"
        :page-sizes="[10, 20, 50]"
        @size-change="loadUsers"
        @current-change="loadUsers"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import { userApi } from '@/api'
import type { User } from '@/types'

const users = ref<User[]>([])
const loading = ref(false)

const filters = reactive({
  keyword: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

onMounted(() => {
  loadUsers()
})

const loadUsers = async () => {
  loading.value = true
  try {
    const res = await userApi.getUsers({
      page: pagination.page,
      page_size: pagination.pageSize,
      keyword: filters.keyword
    })
    users.value = res.data.items
    pagination.total = res.data.total
  } catch {
    // ignore
  } finally {
    loading.value = false
  }
}

const handleAuth = async (id: number, status: string) => {
  try {
    const action = status === 'active' ? '通过' : '拒绝'
    await ElMessageBox.confirm(`确定要${action}该用户的认证吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await userApi.updateAuthStatus(id, status)
    ElMessage.success('操作成功')
    loadUsers()
  } catch (err: any) {
    if (err !== 'cancel') {
      ElMessage.error(err.message || '操作失败')
    }
  }
}

const handleStatus = async (id: number, status: string) => {
  try {
    await userApi.updateUserStatus(id, status)
    ElMessage.success('操作成功')
    loadUsers()
  } catch (err: any) {
    ElMessage.error(err.message || '操作失败')
  }
}

const formatDate = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

const getAuthStatusType = (status: string) => {
  const map: Record<string, string> = {
    pending: 'warning',
    active: 'success',
    disabled: 'info',
    rejected: 'danger'
  }
  return map[status] || 'info'
}

const getAuthStatusText = (status: string) => {
  const map: Record<string, string> = {
    pending: '待审核',
    active: '已认证',
    disabled: '已禁用',
    rejected: '已拒绝'
  }
  return map[status] || status
}
</script>
