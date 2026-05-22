<template>
  <div class="user-management">
    <h3>用户管理</h3>
    <el-card class="filter-card">
      <el-input
        v-model="keyword"
        placeholder="搜索用户名/邮箱"
        style="width: 300px"
        clearable
        @keyup.enter="fetchUsers"
      />
    </el-card>
    <el-table :data="users" v-loading="loading">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column label="头像" width="80">
        <template #default="{ row }">
          <el-avatar :size="40" :src="row.avatar">{{ row.username?.charAt(0)?.toUpperCase() }}</el-avatar>
        </template>
      </el-table-column>
      <el-table-column prop="username" label="用户名" width="120" />
      <el-table-column prop="email" label="邮箱" width="200" />
      <el-table-column prop="phone" label="手机号" width="130" />
      <el-table-column prop="balance" label="余额" width="120">
        <template #default="{ row }">¥{{ row.balance.toFixed(2) }}</template>
      </el-table-column>
      <el-table-column prop="role" label="角色" width="100">
        <template #default="{ row }">
          <el-tag :type="row.role === 'admin' ? 'danger' : 'primary'" size="small">{{ row.role }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
            {{ row.status === 1 ? '正常' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="注册时间" width="160">
        <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="120">
        <template #default="{ row }">
          <el-switch
            v-model="row.status"
            :active-value="1"
            :inactive-value="0"
            @change="updateStatus(row)"
          />
        </template>
      </el-table-column>
    </el-table>
    <el-pagination
      class="pagination"
      v-model:current-page="page"
      v-model:page-size="pageSize"
      :total="total"
      layout="total, prev, pager, next"
      @current-change="fetchUsers"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import dayjs from 'dayjs'
import type { User } from '@/types'
import { adminApi } from '@/api'

const users = ref<User[]>([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const keyword = ref('')

const fetchUsers = async () => {
  loading.value = true
  try {
    const res = await adminApi.getUsers({ page: page.value, page_size: pageSize.value, keyword: keyword.value })
    users.value = res.list
    total.value = res.total
  } finally {
    loading.value = false
  }
}

const updateStatus = async (user: User) => {
  try {
    await adminApi.updateUserStatus(user.id, user.status)
  } catch (e) {
    user.status = user.status === 1 ? 0 : 1
  }
}

const formatTime = (time: string) => dayjs(time).format('YYYY-MM-DD HH:mm')

onMounted(() => {
  fetchUsers()
})
</script>

<style scoped>
h3 {
  margin: 0 0 20px;
}

.filter-card {
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}
</style>
