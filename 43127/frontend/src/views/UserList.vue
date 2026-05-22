<template>
  <div class="user-list">
    <div class="page-header">
      <h2 class="page-title">用户管理</h2>
    </div>

    <div class="search-bar card">
      <el-select v-model="roleFilter" placeholder="角色" clearable style="width: 150px">
        <el-option label="管理员" value="admin" />
        <el-option label="房东" value="landlord" />
        <el-option label="维修人员" value="maintainer" />
        <el-option label="租户" value="tenant" />
      </el-select>
      <el-button type="primary" @click="loadData">搜索</el-button>
    </div>

    <div class="card">
      <el-table :data="users" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户名" width="150" />
        <el-table-column prop="realName" label="姓名" width="120" />
        <el-table-column prop="phone" label="手机号" width="130" />
        <el-table-column prop="email" label="邮箱" width="180" />
        <el-table-column prop="role" label="角色" width="100">
          <template #default="{ row }">
            <el-tag :type="getRoleType(row.role)">{{ getRoleText(row.role) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button
              link
              :type="row.status === 1 ? 'warning' : 'success'"
              @click="toggleStatus(row)"
            >{{ row.status === 1 ? '禁用' : '启用' }}</el-button>
            <el-button
              link
              type="danger"
              @click="deleteUser(row)"
              v-if="row.username !== 'admin'"
            >删除</el-button>
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
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { User } from '@/types'
import { getUsers, updateUserStatus, deleteUser as deleteUserApi } from '@/api/auth'

const loading = ref(false)
const users = ref<User[]>([])
const roleFilter = ref('')

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

onMounted(() => {
  loadData()
})

async function loadData() {
  loading.value = true
  try {
    const res = await getUsers({
      page: pagination.page,
      pageSize: pagination.pageSize,
      role: roleFilter.value || undefined
    })
    users.value = res.data.list
    pagination.total = res.data.total
  } catch (error) {
    console.error('Failed to load users:', error)
  } finally {
    loading.value = false
  }
}

function getRoleType(role: string) {
  switch (role) {
    case 'admin': return 'danger'
    case 'landlord': return 'primary'
    case 'maintainer': return 'warning'
    default: return 'info'
  }
}

function getRoleText(role: string) {
  switch (role) {
    case 'admin': return '管理员'
    case 'landlord': return '房东'
    case 'maintainer': return '维修人员'
    default: return '租户'
  }
}

async function toggleStatus(row: User) {
  try {
    await updateUserStatus(row.id, row.status === 1 ? 0 : 1)
    ElMessage.success('操作成功')
    loadData()
  } catch (error) {
    console.error(error)
  }
}

async function deleteUser(row: User) {
  try {
    await ElMessageBox.confirm(`确定要删除用户 ${row.username} 吗？`, '提示', { type: 'warning' })
    await deleteUserApi(row.id)
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
.user-list {
  padding: 0;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
