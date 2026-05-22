<template>
  <div class="log-page">
    <div class="page-header">
      <el-input
        v-model="searchQuery"
        placeholder="搜索模块..."
        :prefix-icon="Search"
        clearable
        style="width: 200px"
        @clear="fetchLogs"
        @keyup.enter="fetchLogs"
      />
      <el-select v-model="actionFilter" placeholder="操作类型" clearable style="width: 140px" @change="fetchLogs">
        <el-option label="创建" value="create" />
        <el-option label="更新" value="update" />
        <el-option label="删除" value="delete" />
      </el-select>
    </div>

    <el-table :data="logs" v-loading="loading" stripe>
      <el-table-column prop="module" label="模块" width="120">
        <template #default="{ row }">
          <el-tag size="small">{{ moduleText(row.module) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="action" label="操作" width="100">
        <template #default="{ row }">
          <el-tag :type="actionType(row.action)" size="small">{{ actionText(row.action) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="detail" label="详情" min-width="200" show-overflow-tooltip />
      <el-table-column prop="ip_address" label="IP地址" width="140" />
      <el-table-column label="操作时间" width="180">
        <template #default="{ row }">
          {{ formatDate(row.created_at) }}
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      v-model:current-page="page"
      v-model:page-size="pageSize"
      :page-sizes="[10, 20, 50, 100]"
      :total="total"
      layout="total, sizes, prev, pager, next, jumper"
      @size-change="fetchLogs"
      @current-change="fetchLogs"
      class="pagination"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { logApi } from '@/api/dashboard'
import { Search } from '@element-plus/icons-vue'
import dayjs from 'dayjs'

const loading = ref(false)
const logs = ref<any[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const searchQuery = ref('')
const actionFilter = ref('')

function moduleText(module: string) {
  const texts: Record<string, string> = {
    auth: '认证',
    user: '用户',
    wedding: '婚礼',
    vendor: '供应商',
    guest: '嘉宾',
    budget: '预算',
    task: '任务',
    document: '文档',
    notification: '通知'
  }
  return texts[module] || module
}

function actionText(action: string) {
  const texts: Record<string, string> = {
    create: '创建',
    update: '更新',
    delete: '删除'
  }
  return texts[action] || action
}

function actionType(action: string) {
  const types: Record<string, string> = {
    create: 'success',
    update: 'warning',
    delete: 'danger'
  }
  return types[action] || 'info'
}

function formatDate(date: string) {
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}

async function fetchLogs() {
  loading.value = true
  try {
    const res = await logApi.getLogs({
      module: searchQuery.value,
      action: actionFilter.value,
      page: page.value,
      page_size: pageSize.value
    })
    logs.value = res.data.list
    total.value = res.data.total
  } catch (error) {
    console.error('Failed to fetch logs:', error)
  } finally {
    loading.value = false
  }
}

onMounted(fetchLogs)
</script>

<style scoped>
.log-page {
  padding: 0;
}

.page-header {
  display: flex;
  gap: 16px;
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  justify-content: flex-end;
}
</style>
