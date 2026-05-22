<template>
  <div class="user-list-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>用户管理</span>
          <el-input
            v-model="keyword"
            placeholder="搜索用户"
            style="width: 200px"
            clearable
            @keyup.enter="fetchUsers"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </div>
      </template>

      <el-table :data="users" style="width: 100%">
        <el-table-column label="用户" width="200">
          <template #default="{ row }">
            <div class="user-cell">
              <el-avatar :size="32" :src="row.avatar">
                {{ row.nickname?.charAt(0) || 'U' }}
              </el-avatar>
              <div class="user-info">
                <span class="username">{{ row.nickname || row.username }}</span>
                <el-tag v-if="row.isExpert" type="primary" size="small">专家</el-tag>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="email" label="邮箱" />
        <el-table-column label="等级" width="100">
          <template #default="{ row }">
            Lv.{{ row.level }}
          </template>
        </el-table-column>
        <el-table-column prop="points" label="积分" width="100" />
        <el-table-column label="角色" width="100">
          <template #default="{ row }">
            <el-tag :type="row.role === 'admin' ? 'danger' : 'info'" size="small">
              {{ row.role === 'admin' ? '管理员' : '普通用户' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusName(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="注册时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button
              v-if="row.status === 'active'"
              type="warning"
              size="small"
              @click="updateStatus(row.id, 'disabled')"
            >
              禁用
            </el-button>
            <el-button
              v-else
              type="success"
              size="small"
              @click="updateStatus(row.id, 'active')"
            >
              启用
            </el-button>
            <el-button
              v-if="row.status !== 'banned'"
              type="danger"
              size="small"
              @click="updateStatus(row.id, 'banned')"
            >
              封禁
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        class="pagination"
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        layout="prev, pager, next, total"
        @current-change="fetchUsers"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { userApi } from '@/api/user'
import type { User } from '@/types'
import dayjs from 'dayjs'

const users = ref<User[]>([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const keyword = ref('')

const fetchUsers = async () => {
  try {
    const res = await userApi.getUserList({
      page: page.value,
      pageSize: pageSize.value,
      keyword: keyword.value
    })
    users.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch (e) {
    console.error(e)
  }
}

const updateStatus = async (id: number, status: string) => {
  try {
    await userApi.updateUserStatus(id, status)
    fetchUsers()
  } catch (e) {
    console.error(e)
  }
}

const getStatusType = (status: string) => {
  const map: Record<string, string> = {
    active: 'success',
    disabled: 'warning',
    banned: 'danger'
  }
  return map[status] || 'info'
}

const getStatusName = (status: string) => {
  const map: Record<string, string> = {
    active: '正常',
    disabled: '禁用',
    banned: '封禁'
  }
  return map[status] || status
}

const formatTime = (time: string) => {
  return dayjs(time).format('YYYY-MM-DD HH:mm:ss')
}

onMounted(() => {
  fetchUsers()
})
</script>

<style scoped lang="scss">
.user-list-page {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .user-cell {
    display: flex;
    align-items: center;
    gap: 12px;

    .user-info {
      display: flex;
      flex-direction: column;
      gap: 4px;

      .username {
        font-weight: 500;
      }
    }
  }

  .pagination {
    justify-content: center;
    margin-top: 20px;
  }
}
</style>
