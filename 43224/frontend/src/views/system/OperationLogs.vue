<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>操作日志</span>
        </div>
      </template>

      <div class="filter-bar">
        <el-input v-model="module" placeholder="模块" clearable style="width: 140px" @keyup.enter="loadData">
          <template #prefix><el-icon><Menu /></el-icon></template>
        </el-input>
        <el-input v-model="action" placeholder="操作" clearable style="width: 200px" @keyup.enter="loadData">
          <template #prefix><el-icon><Search /></el-icon></template>
        </el-input>
        <el-button type="primary" @click="loadData">搜索</el-button>
      </div>

      <el-table :data="logs" stripe v-loading="loading">
        <el-table-column prop="action" label="操作" min-width="200" />
        <el-table-column prop="module" label="模块" width="120" />
        <el-table-column prop="user.username" label="操作用户" width="120" />
        <el-table-column prop="ip_address" label="IP地址" width="140" />
        <el-table-column prop="user_agent" label="浏览器" min-width="200" show-overflow-tooltip />
        <el-table-column prop="created_at" label="时间" width="160">
          <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          layout="total, prev, pager, next"
          @change="loadData"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { listOperationLogs } from '@/api/statistics'
import dayjs from 'dayjs'

const logs = ref<any[]>([])
const loading = ref(false)
const module = ref('')
const action = ref('')
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

async function loadData() {
  loading.value = true
  try {
    const res = await listOperationLogs({
      module: module.value,
      action: action.value,
      page: page.value,
      page_size: pageSize.value
    }) as any
    if (Array.isArray(res)) {
      logs.value = res
    } else {
      logs.value = res?.list || []
      total.value = res?.total || 0
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function formatDate(date: string) {
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}

onMounted(loadData)
</script>

<style lang="scss" scoped>
.page-container {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .filter-bar {
    display: flex;
    gap: 12px;
    margin-bottom: 16px;
  }

  .pagination {
    margin-top: 16px;
    display: flex;
    justify-content: flex-end;
  }
}
</style>
